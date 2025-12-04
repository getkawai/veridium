package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/pkg/xlog"
	"github.com/kawai-network/veridium/types"
)

// FileProcessorService orchestrates file processing pipeline
type FileProcessorService struct {
	queries         *db.Queries
	fileLoader      *FileLoader
	documentService *DocumentService
	ragProcessor    *RAGProcessor
	libraryService  *llama.LibraryService
}

// NewFileProcessorService creates a new file processor service
func NewFileProcessorService(
	database *sql.DB,
	fileLoader *FileLoader,
	documentService *DocumentService,
	ragProcessor *RAGProcessor,
	libraryService *llama.LibraryService,
) *FileProcessorService {
	return &FileProcessorService{
		queries:         db.New(database),
		fileLoader:      fileLoader,
		documentService: documentService,
		ragProcessor:    ragProcessor,
		libraryService:  libraryService,
	}
}

// ProcessFileRequest represents a file processing request
type ProcessFileRequest struct {
	FilePath     string
	Filename     string
	FileType     string
	UserID       string
	ClientID     string
	Source       string
	EnableRAG    bool // Whether to process for RAG
	IsShared     bool // Whether to store in global_files
	FileMetadata *types.FileMetadata
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

	// Step 1: Save to files table (always)
	fileID, globalFileID, err := s.saveFileMetadata(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to save file metadata: %w", err)
	}
	response.FileID = fileID
	response.GlobalFileID = globalFileID

	// Step 2: Parse file using FileLoader
	fileDoc, err := s.fileLoader.LoadFile(req.FilePath, req.FileMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to load file: %w", err)
	}

	// Step 3: Save to documents table (immediately, without waiting for VL processing)
	documentID, err := s.saveDocument(ctx, fileDoc, fileID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to save document: %w", err)
	}
	response.DocumentID = documentID

	// Step 4: If image, generate description using VL model ASYNCHRONOUSLY
	// This allows the UI to show the image immediately while description is being generated
	if req.FileType == string(types.FileTypeImage) && s.libraryService != nil {
		response.Processing = true
		xlog.Info("Starting async image description generation", "filename", req.Filename, "document_id", documentID)

		go s.processImageDescriptionAsync(req.FilePath, req.Filename, documentID, fileID, req.UserID, req.EnableRAG)
	} else if req.EnableRAG && s.fileLoader.CanChunkForRAG(req.FileType) {
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
	} else if req.EnableRAG && !s.fileLoader.CanChunkForRAG(req.FileType) {
		xlog.Info("Skipping RAG processing for unsupported file type", "file_type", req.FileType, "file_id", fileID)
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

	// Generate description
	prompt := "Describe this image in detail. Include all visible text, objects, and layout."
	description, err := s.libraryService.ProcessImageWithText(filePath, prompt, 512)
	if err != nil {
		xlog.Error("Async: Failed to process image with VL model", "error", err, "filename", filename)
		return
	}

	xlog.Info("Async: Image description generated", "length", len(description), "filename", filename)

	// Format description as markdown
	descriptionMarkdown := fmt.Sprintf("\n\n### Image Description (AI Generated)\n\n%s", description)

	// Update document with description
	err = s.documentService.AppendContentToDocument(ctx, documentID, userID, descriptionMarkdown)
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

// saveFileMetadata saves file metadata to files and optionally global_files
func (s *FileProcessorService) saveFileMetadata(ctx context.Context, req ProcessFileRequest) (string, string, error) {
	fileID := uuid.New().String()
	now := time.Now().UnixMilli()

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
				HashID:    fileHash,
				FileType:  req.FileType,
				Size:      fileInfo.Size,
				Url:       req.FilePath,
				Metadata:  sql.NullString{Valid: false},
				Creator:   req.UserID,
				CreatedAt: now,
			})
			if err != nil {
				return "", "", fmt.Errorf("failed to create global file: %w", err)
			}
		}
	}

	// Save to files table
	fileInfo, err := getFileInfo(req.FilePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to get file info: %w", err)
	}

	_, err = s.queries.CreateFile(ctx, db.CreateFileParams{
		ID:       fileID,
		UserID:   req.UserID,
		FileType: req.FileType,
		FileHash: sql.NullString{String: globalFileID, Valid: globalFileID != ""},
		Name:     req.Filename,
		Size:     fileInfo.Size,
		Url:      req.FilePath,
		Source:   sql.NullString{String: req.Source, Valid: req.Source != ""},
		ClientID: sql.NullString{String: req.ClientID, Valid: req.ClientID != ""},
		Metadata: sql.NullString{Valid: false},
	})

	if err != nil {
		return "", "", fmt.Errorf("failed to create file: %w", err)
	}

	return fileID, globalFileID, nil
}

// saveDocument saves parsed file content to documents table
func (s *FileProcessorService) saveDocument(ctx context.Context, fileDoc *types.FileDocument, fileID string, req ProcessFileRequest) (string, error) {
	// Convert FileDocument.Metadata to map[string]interface{}
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

	// Convert FileDocument.Pages to []DocumentPage
	pages := make([]types.DocumentPage, len(fileDoc.Pages))
	for i, page := range fileDoc.Pages {
		pages[i] = types.DocumentPage{
			CharCount:   page.CharCount,
			LineCount:   page.LineCount,
			Metadata:    page.Metadata,
			PageContent: page.PageContent,
		}
	}

	// Create document
	documentID, err := s.documentService.CreateDocument(ctx, CreateDocumentParams{
		Title:          fileDoc.Filename,
		Content:        fileDoc.Content,
		FileType:       fileDoc.FileType,
		Filename:       fileDoc.Filename,
		TotalCharCount: fileDoc.TotalCharCount,
		TotalLineCount: fileDoc.TotalLineCount,
		Metadata:       metadata,
		Pages:          pages,
		SourceType:     "file",
		Source:         fileDoc.Source,
		FileID:         fileID,
		UserID:         req.UserID,
		ClientID:       req.ClientID,
	})

	return documentID, err
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
