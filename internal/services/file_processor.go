package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/internal/whisper"
	"github.com/kawai-network/veridium/pkg/xlog"
	"github.com/kawai-network/veridium/types"
)

// FileProcessorService orchestrates file processing pipeline
type FileProcessorService struct {
	queries        *db.Queries
	fileLoader     *FileLoader
	ragProcessor   *RAGProcessor
	libraryService *llama.LibraryService
	whisperService *whisper.Service
}

// NewFileProcessorService creates a new file processor service
func NewFileProcessorService(
	database *sql.DB,
	fileLoader *FileLoader,
	ragProcessor *RAGProcessor,
	libraryService *llama.LibraryService,
	whisperService *whisper.Service,
) *FileProcessorService {
	return &FileProcessorService{
		queries:        db.New(database),
		fileLoader:     fileLoader,
		ragProcessor:   ragProcessor,
		libraryService: libraryService,
		whisperService: whisperService,
	}
}



// ProcessFileRequest represents a file processing request
type ProcessFileRequest struct {
	FilePath  string
	Filename  string
	UserID    string
	Source    string
	EnableRAG bool // Whether to process for RAG
	IsShared  bool // Whether to store in global_files
}

// ProcessFileResponse represents the result of file processing
type ProcessFileResponse struct {
	FileID       string   `json:"fileId"`
	DocumentID   string   `json:"documentId"`
	ChunkIDs     []string `json:"chunkIds,omitempty"`
	GlobalFileID string   `json:"globalFileId,omitempty"`
	Processing   bool     `json:"processing,omitempty"` // True if async processing is in progress
}

// ProcessFile is the main entry point for file processing
// It handles: files → global_files (optional) → documents → chunks (optional)
// For images, VL model processing is done asynchronously to avoid blocking the UI
func (s *FileProcessorService) ProcessFile(ctx context.Context, req ProcessFileRequest) (*ProcessFileResponse, error) {
	xlog.Info("Processing file", "filename", req.Filename, "user_id", req.UserID)

	response := &ProcessFileResponse{}

	// Step 1: Parse file using FileLoader (detects file type from extension)
	fileDoc, err := s.fileLoader.LoadFile(req.FilePath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to load file: %w", err)
	}

	// Get detected file type from FileLoader
	detectedFileType := fileDoc.FileType

	// Step 2: Save to files table (always)
	fileID, globalFileID, err := s.saveFileMetadata(ctx, req, detectedFileType)
	if err != nil {
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}
	response.FileID = fileID
	response.GlobalFileID = globalFileID

	// Step 3: Save to documents table (immediately, without waiting for VL processing)
	documentID, err := s.saveDocument(ctx, fileDoc, fileID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to save document: %w", err)
	}
	response.DocumentID = documentID

	// Step 4: If image, generate description using VL model ASYNCHRONOUSLY
	// This allows the UI to show the image immediately while description is being generated
	// Note: detectedFileType is the file extension (e.g., "jpg"), not the SupportedFileType ("image")
	// We don't set Processing=true because the file IS ready for use - only the AI description is pending
	if s.fileLoader.IsImageFile(detectedFileType) && s.libraryService != nil {
		xlog.Info("Starting async image description generation", "filename", req.Filename, "document_id", documentID)

		go s.processImageDescriptionAsync(req.FilePath, req.Filename, documentID, fileID, req.UserID, req.EnableRAG)
	} else if s.fileLoader.IsVideoFile(detectedFileType) {
		// Step 4b: If video, generate description using OpenRouter VL model ASYNCHRONOUSLY
		xlog.Info("Starting async video understanding generation", "filename", req.Filename, "document_id", documentID)

		go s.processVideoDescriptionAsync(req.FilePath, req.Filename, documentID, fileID, req.UserID, req.EnableRAG)
	} else if req.EnableRAG && s.fileLoader.CanChunkForRAG(detectedFileType) {
		// Step 5: Process for RAG (for non-image files, do it synchronously)
		chunkIDs, err := s.ragProcessor.ProcessFile(ctx, RAGProcessRequest{
			FilePath:   req.FilePath,
			FileID:     fileID,
			DocumentID: documentID,
			UserID:     req.UserID,
			Filename:   req.Filename,
		})
		if err != nil {
			xlog.Error("Failed to process file for RAG", "error", err, "file_id", fileID)
		} else {
			xlog.Info("RAG processing completed", "file_id", fileID, "chunks", len(chunkIDs))
		}
	} else if req.EnableRAG && !s.fileLoader.CanChunkForRAG(detectedFileType) {
		xlog.Info("Skipping RAG processing for unsupported file type", "file_type", detectedFileType, "file_id", fileID)
	}

	return response, nil
}

// processImageDescriptionAsync generates image description asynchronously
// and updates the document when complete
func (s *FileProcessorService) processImageDescriptionAsync(filePath, filename, documentID, fileID, userID string, enableRAG bool) {
	ctx := context.Background()

	xlog.Info("Async: Starting VL model image processing", "filename", filename)

	// Ensure VL model is loaded
	if !s.libraryService.IsVLModelLoaded() {
		if err := s.libraryService.LoadVLModel(""); err != nil {
			xlog.Error("Async: Failed to load VL model for image processing", "error", err)
			return
		}
	}

	if !s.libraryService.IsVLModelLoaded() {
		xlog.Error("Async: VL model not available", "filename", filename)
		return
	}

	// Generate description with sufficient tokens for document OCR
	prompt := "Describe this image in detail. Include all visible text, objects, and layout."
	description, err := s.libraryService.ProcessImageWithText(filePath, prompt, 2048)
	if err != nil {
		xlog.Error("Async: Failed to process image with VL model", "error", err, "filename", filename)
		return
	}

	xlog.Info("Async: Image description generated", "length", len(description), "filename", filename)

	// Format description as markdown
	descriptionMarkdown := fmt.Sprintf("\n\n### Image Description (AI Generated)\n\n%s", description)

	// Update document with description
	err = s.appendContentToDocument(ctx, documentID, userID, descriptionMarkdown)
	if err != nil {
		xlog.Error("Async: Failed to update document with description", "error", err, "document_id", documentID)
		return
	}

	xlog.Info("Async: Document updated with image description", "document_id", documentID)

	// Process for RAG if enabled (now that we have the description)
	if enableRAG {
		chunkIDs, err := s.ragProcessor.ProcessFile(ctx, RAGProcessRequest{
			FilePath:   filePath,
			FileID:     fileID,
			DocumentID: documentID,
			UserID:     userID,
			Filename:   filename,
		})
		if err != nil {
			xlog.Error("Async: Failed to process file for RAG", "error", err, "file_id", fileID)
		} else {
			xlog.Info("Async: RAG processing completed", "file_id", fileID, "chunks", len(chunkIDs))
		}
	}
}

// processVideoDescriptionAsync extracts audio from video using ffmpeg and transcribes using whisper
func (s *FileProcessorService) processVideoDescriptionAsync(filePath, filename, documentID, fileID, userID string, enableRAG bool) {
	ctx := context.Background()

	xlog.Info("Async: Starting video transcription", "filename", filename)

	// Check if whisper service is available
	if s.whisperService == nil {
		xlog.Error("Async: Whisper service not available", "filename", filename)
		return
	}

	// Check if whisper is installed
	if !s.whisperService.IsWhisperInstalled() {
		xlog.Error("Async: Whisper CLI not installed", "filename", filename)
		return
	}

	// Check if we have a model
	models, err := s.whisperService.ListModels()
	if err != nil || len(models) == 0 {
		xlog.Error("Async: No whisper models available", "error", err, "filename", filename)
		return
	}
	modelName := models[0] // Use first available model

	// Step 1: Extract audio from video using ffmpeg
	audioPath, err := s.extractAudioFromVideo(filePath)
	if err != nil {
		xlog.Error("Async: Failed to extract audio from video", "error", err, "filename", filename)
		return
	}
	defer os.Remove(audioPath) // Clean up temp audio file

	xlog.Info("Async: Audio extracted from video", "audio_path", audioPath, "filename", filename)

	// Step 2: Transcribe audio using whisper
	xlog.Info("Async: Starting whisper transcription", "model", modelName, "filename", filename)
	transcription, err := s.whisperService.Transcribe(ctx, modelName, audioPath)
	if err != nil {
		xlog.Error("Async: Failed to transcribe audio", "error", err, "filename", filename)
		return
	}

	if transcription == "" {
		xlog.Warn("Async: Empty transcription result", "filename", filename)
		transcription = "(No speech detected in video)"
	}

	xlog.Info("Async: Video transcription completed", "length", len(transcription), "filename", filename)

	// Format transcription as markdown
	transcriptionMarkdown := fmt.Sprintf("\n\n### Video Transcription (AI Generated via Whisper)\n\n%s", transcription)

	// Update document with transcription
	err = s.appendContentToDocument(ctx, documentID, userID, transcriptionMarkdown)
	if err != nil {
		xlog.Error("Async: Failed to update document with transcription", "error", err, "document_id", documentID)
		return
	}

	xlog.Info("Async: Document updated with video transcription", "document_id", documentID)

	// Process for RAG if enabled (now that we have the transcription)
	if enableRAG {
		chunkIDs, err := s.ragProcessor.ProcessFile(ctx, RAGProcessRequest{
			FilePath:   filePath,
			FileID:     fileID,
			DocumentID: documentID,
			UserID:     userID,
			Filename:   filename,
		})
		if err != nil {
			xlog.Error("Async: Failed to process video for RAG", "error", err, "file_id", fileID)
		} else {
			xlog.Info("Async: RAG processing completed for video", "file_id", fileID, "chunks", len(chunkIDs))
		}
	}
}

// extractAudioFromVideo extracts audio from video file using ffmpeg
// Returns path to temporary WAV file (caller must clean up)
func (s *FileProcessorService) extractAudioFromVideo(videoPath string) (string, error) {
	// Create temp file for audio
	tempDir := os.TempDir()
	audioPath := filepath.Join(tempDir, fmt.Sprintf("audio_%d.wav", time.Now().UnixNano()))

	// Find ffmpeg binary
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		// Try common locations
		homeDir, _ := os.UserHomeDir()
		localFFmpeg := filepath.Join(homeDir, ".local", "bin", "ffmpeg")
		if _, statErr := os.Stat(localFFmpeg); statErr == nil {
			ffmpegPath = localFFmpeg
		} else {
			return "", fmt.Errorf("ffmpeg not found: %w", err)
		}
	}

	// Extract audio: convert to 16kHz mono WAV (whisper's preferred format)
	cmd := exec.Command(ffmpegPath,
		"-i", videoPath,
		"-vn",           // No video
		"-acodec", "pcm_s16le", // 16-bit PCM
		"-ar", "16000",  // 16kHz sample rate
		"-ac", "1",      // Mono
		"-y",            // Overwrite
		audioPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("ffmpeg failed: %w, output: %s", err, string(output))
	}

	// Verify audio file was created
	if _, err := os.Stat(audioPath); err != nil {
		return "", fmt.Errorf("audio file not created: %w", err)
	}

	return audioPath, nil
}

// getVideoMimeType returns the MIME type for a video file based on extension
func getVideoMimeType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".mp4":
		return "video/mp4"
	case ".mkv":
		return "video/x-matroska"
	case ".avi":
		return "video/x-msvideo"
	case ".mov":
		return "video/quicktime"
	case ".wmv":
		return "video/x-ms-wmv"
	case ".flv":
		return "video/x-flv"
	case ".webm":
		return "video/webm"
	case ".m4v":
		return "video/x-m4v"
	case ".mpeg", ".mpg":
		return "video/mpeg"
	case ".3gp":
		return "video/3gpp"
	default:
		return "video/mp4"
	}
}

// saveFileMetadata saves file metadata to files and optionally global_files
func (s *FileProcessorService) saveFileMetadata(ctx context.Context, req ProcessFileRequest, fileType string) (string, string, error) {
	var globalFileID string

	// If shared, save to global_files first
	if req.IsShared {
		// Calculate file hash
		fileHash, err := calculateFileHash(req.FilePath)
		if err != nil {
			return "", "", fmt.Errorf("failed to calculate file hash: %w", err)
		}

		// Check if global file already exists
		existing, err := s.queries.GetGlobalFileByHash(ctx, fileHash)
		if err == nil && existing.HashID != "" {
			globalFileID = existing.HashID
		} else {
			// Create new global file
			fileInfo, err := getFileInfo(req.FilePath)
			if err != nil {
				return "", "", fmt.Errorf("failed to get file info: %w", err)
			}
			globalFileID = fileHash

			_, err = s.queries.CreateGlobalFile(ctx, db.CreateGlobalFileParams{
				HashID:   fileHash,
				FileType: fileType,
				Size:     fileInfo.Size,
				Url:      req.FilePath,
				Metadata: sql.NullString{Valid: false},
				Creator:  req.UserID,
			})
			if err != nil {
				return "", "", fmt.Errorf("failed to create global file: %w", err)
			}
		}
	}

	// Save to files table (ID and timestamps generated by SQLite)
	fileInfo, err := getFileInfo(req.FilePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to get file info: %w", err)
	}

	file, err := s.queries.CreateFile(ctx, db.CreateFileParams{
		UserID:   req.UserID,
		FileType: fileType,
		FileHash: sql.NullString{String: globalFileID, Valid: globalFileID != ""},
		Name:     req.Filename,
		Size:     fileInfo.Size,
		Url:      req.FilePath,
		Source:   sql.NullString{String: req.Source, Valid: req.Source != ""},
		Metadata: sql.NullString{Valid: false},
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to create file: %w", err)
	}

	return file.ID, globalFileID, nil
}

// saveDocument saves parsed file content to documents table
func (s *FileProcessorService) saveDocument(ctx context.Context, fileDoc *types.FileDocument, fileID string, req ProcessFileRequest) (string, error) {
	// Convert FileDocument.Metadata to JSON
	metadata := map[string]interface{}{
		"source":       fileDoc.Metadata.Source,
		"filename":     fileDoc.Metadata.Filename,
		"fileType":     fileDoc.Metadata.FileType,
		"createdTime":  fileDoc.Metadata.CreatedTime,
		"modifiedTime": fileDoc.Metadata.ModifiedTime,
	}
	if fileDoc.Metadata.Error != "" {
		metadata["error"] = fileDoc.Metadata.Error
	}
	metadataJSON, _ := json.Marshal(metadata)

	// Convert FileDocument.Pages to JSON
	pagesJSON, _ := json.Marshal(fileDoc.Pages)

	// Create document (ID and timestamps generated by SQLite)
	doc, err := s.queries.CreateDocument(ctx, db.CreateDocumentParams{
		Title:          sql.NullString{String: fileDoc.Filename, Valid: fileDoc.Filename != ""},
		Content:        sql.NullString{String: fileDoc.Content, Valid: fileDoc.Content != ""},
		FileType:       fileDoc.FileType,
		Filename:       sql.NullString{String: fileDoc.Filename, Valid: fileDoc.Filename != ""},
		TotalCharCount: int64(fileDoc.TotalCharCount),
		TotalLineCount: int64(fileDoc.TotalLineCount),
		Metadata:       sql.NullString{String: string(metadataJSON), Valid: true},
		Pages:          sql.NullString{String: string(pagesJSON), Valid: true},
		SourceType:     "file",
		Source:         fileDoc.Source,
		FileID:         sql.NullString{String: fileID, Valid: fileID != ""},
		UserID:         req.UserID,
		EditorData:     sql.NullString{},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create document: %w", err)
	}

	return doc.ID, nil
}

// appendContentToDocument appends content to an existing document
// Used for async operations like image description generation
func (s *FileProcessorService) appendContentToDocument(ctx context.Context, documentID, userID, additionalContent string) error {
	// Get existing document
	doc, err := s.queries.GetDocument(ctx, db.GetDocumentParams{
		ID:     documentID,
		UserID: userID,
	})
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	// Append new content
	newContent := doc.Content.String + additionalContent
	now := time.Now().UnixMilli()

	// Update document
	_, err = s.queries.UpdateDocument(ctx, db.UpdateDocumentParams{
		ID:         documentID,
		UserID:     userID,
		Title:      doc.Title,
		Content:    sql.NullString{String: newContent, Valid: true},
		Metadata:   doc.Metadata,
		EditorData: doc.EditorData,
		UpdatedAt:  now,
	})
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	return nil
}

// Helper functions

// calculateFileHash calculates SHA256 hash of a file
func calculateFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// FileInfo represents file information
type FileInfo struct {
	Size int64
}

// getFileInfo gets file information
func getFileInfo(filePath string) (FileInfo, error) {
	stat, err := os.Stat(filePath)
	if err != nil {
		return FileInfo{}, err
	}
	return FileInfo{Size: stat.Size()}, nil
}
