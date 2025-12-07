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

	"github.com/kawai-network/veridium/fantasy"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/internal/whisper"
	"github.com/kawai-network/veridium/pkg/xlog"
	"github.com/kawai-network/veridium/types"
)

// FileProcessorService orchestrates file processing pipeline
type FileProcessorService struct {
	db             *sql.DB
	queries        *db.Queries
	fileLoader     *FileLoader
	ragProcessor   *RAGProcessor
	libraryService *llama.LibraryService
	whisperService *whisper.Service
	languageModel  fantasy.LanguageModel // For OCR/transcript cleanup
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
		db:             database,
		queries:        db.New(database),
		fileLoader:     fileLoader,
		ragProcessor:   ragProcessor,
		libraryService: libraryService,
		whisperService: whisperService,
	}
}

// SetLanguageModel sets the language model for OCR/transcript cleanup
func (s *FileProcessorService) SetLanguageModel(model fantasy.LanguageModel) {
	s.languageModel = model
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
	// Always try image processing - OCR can work without VL model
	if s.fileLoader.IsImageFile(detectedFileType) {
		xlog.Info("Starting async image processing (hybrid OCR/VL)", "filename", req.Filename, "document_id", documentID)

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

// processImageDescriptionAsync generates image description asynchronously using hybrid approach:
// 1. First try Tesseract OCR (fast) for text extraction
// 2. If significant text found, use that (fast path)
// 3. If no/minimal text, fallback to VL model for image description (slow path)
func (s *FileProcessorService) processImageDescriptionAsync(filePath, filename, documentID, fileID, userID string, enableRAG bool) {
	ctx := context.Background()

	xlog.Info("Async: Starting hybrid image processing", "filename", filename)

	var finalContent string
	var contentType string

	// Step 1: Try Tesseract OCR first (fast path)
	ocrText, err := s.extractTextWithTesseract(filePath)
	if err != nil {
		xlog.Warn("Async: Tesseract OCR failed, will try VL model", "error", err, "filename", filename)
	}

	// Check if we got meaningful text (more than 20 chars, excluding whitespace)
	// Lower threshold to catch short but valid text like logos, labels, etc.
	cleanedText := strings.TrimSpace(ocrText)
	if len(cleanedText) > 20 {
		// Fast path: sufficient text extracted via OCR
		xlog.Info("Async: OCR extracted sufficient text", "length", len(cleanedText), "filename", filename)

		// Step 2: Clean up OCR text using LLM (if available)
		if s.languageModel != nil {
			xlog.Info("Async: Cleaning up OCR text with LLM", "filename", filename)
			cleanedContent, err := s.cleanupOCRText(ctx, cleanedText, filename)
			if err != nil {
				xlog.Warn("Async: LLM cleanup failed, using raw OCR", "error", err, "filename", filename)
				finalContent = cleanedText
				contentType = "OCR Text (Tesseract - raw)"
			} else {
				finalContent = cleanedContent
				contentType = "OCR Text (Tesseract + LLM cleanup)"
			}
		} else {
			finalContent = cleanedText
			contentType = "OCR Text (Tesseract)"
		}
	} else {
		// Slow path: use VL model for image description
		xlog.Info("Async: Minimal text from OCR, using VL model", "ocr_length", len(cleanedText), "filename", filename)

		// Ensure VL model is loaded
		if s.libraryService != nil {
			if !s.libraryService.IsVLModelLoaded() {
				if err := s.libraryService.LoadVLModel(""); err != nil {
					xlog.Error("Async: Failed to load VL model", "error", err)
					// If VL fails but we have some OCR text, use that
					if len(cleanedText) > 0 {
						finalContent = cleanedText
						contentType = "OCR Text (Tesseract - partial)"
					} else {
						return
					}
				}
			}

			if s.libraryService.IsVLModelLoaded() && finalContent == "" {
				prompt := "Describe this image in detail. Include all visible text, objects, people, colors, and layout."
				description, err := s.libraryService.ProcessImageWithText(filePath, prompt, 2048)
				if err != nil {
					xlog.Error("Async: VL model processing failed", "error", err, "filename", filename)
					// Fallback to OCR text if available
					if len(cleanedText) > 0 {
						finalContent = cleanedText
						contentType = "OCR Text (Tesseract - VL fallback failed)"
					} else {
						return
					}
				} else {
					finalContent = description
					contentType = "Image Description (VL Model)"

					// If we also have OCR text, append it
					if len(cleanedText) > 0 {
						finalContent = fmt.Sprintf("%s\n\n**Extracted Text (OCR):**\n%s", description, cleanedText)
					}
				}
			}
		} else if len(cleanedText) > 0 {
			// No VL service but have some OCR text
			finalContent = cleanedText
			contentType = "OCR Text (Tesseract - no VL available)"
		} else {
			xlog.Error("Async: No VL model available and no OCR text", "filename", filename)
			return
		}
	}

	if finalContent == "" {
		xlog.Warn("Async: No content extracted from image", "filename", filename)
		return
	}

	xlog.Info("Async: Image processing completed", "type", contentType, "length", len(finalContent), "filename", filename)

	// Format content as markdown
	contentMarkdown := fmt.Sprintf("\n\n### %s\n\n%s", contentType, finalContent)

	// Update document with content
	err = s.appendContentToDocument(ctx, documentID, userID, contentMarkdown)
	if err != nil {
		xlog.Error("Async: Failed to update document", "error", err, "document_id", documentID)
		return
	}

	xlog.Info("Async: Document updated", "document_id", documentID, "type", contentType)

	// Process for RAG if enabled
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

// extractTextWithTesseract extracts text from image using Tesseract OCR
func (s *FileProcessorService) extractTextWithTesseract(imagePath string) (string, error) {
	tesseractPath, err := exec.LookPath("tesseract")
	if err != nil {
		return "", fmt.Errorf("tesseract not found: %w", err)
	}

	// Run tesseract: tesseract <image> stdout
	cmd := exec.Command(tesseractPath, imagePath, "stdout", "-l", "eng+ind+jpn+chi_sim")
	output, err := cmd.Output()
	if err != nil {
		// Try without language packs (just default)
		cmd = exec.Command(tesseractPath, imagePath, "stdout")
		output, err = cmd.Output()
		if err != nil {
			return "", fmt.Errorf("tesseract failed: %w", err)
		}
	}

	return string(output), nil
}

// cleanupOCRText uses LLM to clean up raw OCR text and format it as markdown
func (s *FileProcessorService) cleanupOCRText(ctx context.Context, rawText, filename string) (string, error) {
	if s.languageModel == nil {
		return rawText, nil
	}

	// Determine likely document type from filename
	ext := strings.ToLower(filepath.Ext(filename))
	docHint := ""
	switch ext {
	case ".pdf":
		docHint = "This appears to be from a PDF document."
	case ".png", ".jpg", ".jpeg":
		docHint = "This is from a screenshot or image."
	default:
		docHint = "This is from an image file."
	}

	systemPrompt := `You are an OCR text editor. Clean up raw OCR output:
- Fix obvious typos and character recognition errors
- Preserve original structure and formatting
- Format as proper markdown when appropriate
- Output ONLY the cleaned text without explanations.`

	userPrompt := fmt.Sprintf(`Clean up the following OCR-extracted text and format it as clean markdown.

%s

Instructions:
1. Fix obvious OCR errors and typos
2. Remove artifacts like random characters, broken words
3. Preserve the original meaning and structure
4. Format as proper markdown:
   - Use headers (##, ###) for section titles
   - Use bullet points or numbered lists where appropriate
   - Use **bold** for emphasis or important terms
   - Use code blocks for any code or technical content
5. If it's a table, format it as a markdown table
6. Keep the content concise but complete
7. Output ONLY the cleaned markdown, no explanations

Raw OCR text:
---
%s
---

Cleaned markdown:`, docHint, rawText)

	xlog.Info("Async: Calling LLM for OCR cleanup", "prompt_len", len(userPrompt))

	// Use timeout context - 60s should be enough for OCR cleanup
	timeoutCtx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	resp, err := s.languageModel.Generate(timeoutCtx, fantasy.Call{
		Prompt: []fantasy.Message{
			fantasy.NewSystemMessage(systemPrompt),
			fantasy.NewUserMessage(userPrompt),
		},
	})
	if err != nil {
		return "", fmt.Errorf("LLM cleanup failed: %w", err)
	}

	xlog.Info("Async: LLM OCR response received", "response_len", len(resp.Content.Text()))

	result := resp.Content.Text()

	// Trim any leading/trailing whitespace
	result = strings.TrimSpace(result)

	// If result is empty or too short, return original
	if len(result) < 10 {
		return rawText, nil
	}

	xlog.Info("Async: OCR text cleaned up", "original_len", len(rawText), "cleaned_len", len(result))
	return result, nil
}

// cleanupTranscription uses LLM to clean up and correct Whisper transcription
// Fixes common Whisper errors like misheard words, typos, and formatting issues
func (s *FileProcessorService) cleanupTranscription(ctx context.Context, rawTranscript string) (string, error) {
	if s.languageModel == nil {
		return rawTranscript, nil // Return original if no LLM available
	}

	// Skip if transcript is too short
	if len(rawTranscript) < 50 {
		return rawTranscript, nil
	}

	systemPrompt := `You are a transcription editor for Indonesian language content.
Your task is to fix common speech-to-text errors. Apply these corrections:
- terlion/terliunan → triliun
- stengah → setengah
- merogikan → merugikan
- ngontongnya → menguntungkan
- meleksaham → melek saham
- persubstritp → oversubscribe
- SOJK/SCUJK → SEOJK
- IHSK → IHSG
- rebu → ribu
- lokasinya → alokasinya
- timis → tipis

Keep segment markers (**[Segment X]**). Keep stock codes (RLCO, GOTO, PGHB).
Output ONLY the corrected text without explanations.`

	userPrompt := fmt.Sprintf(`Fix this transcription from Whisper speech-to-text:

%s`, rawTranscript)

	xlog.Info("Async: Calling LLM for transcript cleanup", "prompt_len", len(userPrompt))

	// Use timeout context - LLM calls for long transcripts can take a while
	timeoutCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
	defer cancel()

	resp, err := s.languageModel.Generate(timeoutCtx, fantasy.Call{
		Prompt: []fantasy.Message{
			fantasy.NewSystemMessage(systemPrompt),
			fantasy.NewUserMessage(userPrompt),
		},
	})
	if err != nil {
		xlog.Warn("Async: Transcript cleanup failed, using original", "error", err)
		return rawTranscript, nil // Return original on error, don't fail
	}

	xlog.Info("Async: LLM response received", "response_len", len(resp.Content.Text()))

	result := strings.TrimSpace(resp.Content.Text())

	// If result is empty or significantly shorter, return original
	if len(result) < len(rawTranscript)/2 {
		xlog.Warn("Async: Transcript cleanup result too short, using original",
			"original_len", len(rawTranscript), "result_len", len(result))
		return rawTranscript, nil
	}

	xlog.Info("Async: Transcript cleaned up", "original_len", len(rawTranscript), "cleaned_len", len(result))
	return result, nil
}

// processVideoDescriptionAsync extracts audio from video using ffmpeg and transcribes using whisper
// Uses parallel chunked transcription for faster processing with progressive DB updates
func (s *FileProcessorService) processVideoDescriptionAsync(filePath, filename, documentID, fileID, userID string, enableRAG bool) {
	ctx := context.Background()

	xlog.Info("Async: Starting video transcription (parallel)", "filename", filename)

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

	// Check if we have a model - prefer tiny for speed in parallel mode
	models, err := s.whisperService.ListModels()
	if err != nil || len(models) == 0 {
		xlog.Error("Async: No whisper models available", "error", err, "filename", filename)
		return
	}

	// Select Whisper model based on available RAM
	// Model sizes: tiny (~75MB), base (~150MB), small (~500MB), medium (~1.5GB), large (~3GB)
	// RAM requirements: tiny=1GB, base=2GB, small=4GB, medium=8GB, large=16GB
	modelName := selectWhisperModelByHardware(models)
	xlog.Info("Async: Selected Whisper model based on hardware", "model", modelName)

	// Get video duration
	duration, err := s.getVideoDuration(filePath)
	if err != nil {
		xlog.Warn("Async: Could not get video duration, using sequential mode", "error", err)
		duration = 0
	}

	xlog.Info("Async: Video info", "duration_sec", duration, "model", modelName, "filename", filename)

	// Initialize transcription header in document
	err = s.appendContentToDocument(ctx, documentID, userID, "\n\n### Video Transcription (AI Generated via Whisper)\n\n*Transcribing...*\n")
	if err != nil {
		xlog.Error("Async: Failed to initialize transcription header", "error", err)
	}

	var fullTranscription string

	// Use parallel chunked transcription for videos > 2 minutes
	if duration > 120 {
		// Calculate number of chunks (each chunk ~5 minutes = 300 seconds)
		chunkDuration := 300.0
		numChunks := int(duration/chunkDuration) + 1
		if numChunks > 8 {
			numChunks = 8 // Cap at 8 parallel workers
		}
		if numChunks < 2 {
			numChunks = 2
		}

		xlog.Info("Async: Using parallel transcription", "chunks", numChunks, "chunk_duration", chunkDuration)

		fullTranscription, err = s.transcribeVideoParallel(ctx, filePath, modelName, numChunks, chunkDuration, documentID, userID)
		if err != nil {
			xlog.Error("Async: Parallel transcription failed", "error", err, "filename", filename)
			return
		}
	} else {
		// Short video - use sequential transcription
		xlog.Info("Async: Using sequential transcription (short video)")

		audioPath, err := s.extractAudioFromVideo(filePath, 0, 0) // Full video
		if err != nil {
			xlog.Error("Async: Failed to extract audio", "error", err)
			return
		}
		defer os.Remove(audioPath)

		fullTranscription, err = s.whisperService.Transcribe(ctx, modelName, audioPath)
		if err != nil {
			xlog.Error("Async: Transcription failed", "error", err)
			return
		}
	}

	if fullTranscription == "" {
		fullTranscription = "(No speech detected in video)"
	}

	xlog.Info("Async: Video transcription completed", "length", len(fullTranscription), "filename", filename)

	// Post-process transcription with LLM to fix common Whisper errors
	if s.languageModel != nil && len(fullTranscription) > 50 {
		xlog.Info("Async: Cleaning up transcription with LLM", "filename", filename)
		cleanedTranscription, err := s.cleanupTranscription(ctx, fullTranscription)
		if err == nil && len(cleanedTranscription) > 0 {
			fullTranscription = cleanedTranscription
			xlog.Info("Async: Transcription cleanup completed", "length", len(fullTranscription))
		}
	}

	// Update document with final complete transcription
	finalMarkdown := fmt.Sprintf("\n\n### Video Transcription (AI Generated via Whisper + LLM Cleanup)\n\n%s", fullTranscription)
	err = s.replaceTranscriptionInDocument(ctx, documentID, userID, finalMarkdown)
	if err != nil {
		xlog.Error("Async: Failed to update document with final transcription", "error", err)
	}

	xlog.Info("Async: Document updated with complete transcription", "document_id", documentID)

	// Process for RAG if enabled
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

// transcribeVideoParallel transcribes video in parallel chunks with progressive updates
func (s *FileProcessorService) transcribeVideoParallel(ctx context.Context, videoPath, modelName string, numChunks int, chunkDuration float64, documentID, userID string) (string, error) {
	// Create temp directory for chunks
	tempDir, err := os.MkdirTemp("", "whisper_chunks_*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Extract audio chunks in parallel
	type chunkResult struct {
		index int
		path  string
		err   error
	}

	extractChan := make(chan chunkResult, numChunks)

	xlog.Info("Async: Extracting audio chunks", "count", numChunks)

	for i := 0; i < numChunks; i++ {
		go func(idx int) {
			startTime := float64(idx) * chunkDuration
			chunkPath := filepath.Join(tempDir, fmt.Sprintf("chunk_%d.wav", idx))

			err := s.extractAudioChunk(videoPath, chunkPath, startTime, chunkDuration)
			extractChan <- chunkResult{index: idx, path: chunkPath, err: err}
		}(i)
	}

	// Collect extracted chunks
	chunkPaths := make([]string, numChunks)
	for i := 0; i < numChunks; i++ {
		result := <-extractChan
		if result.err != nil {
			xlog.Warn("Async: Failed to extract chunk", "index", result.index, "error", result.err)
			continue
		}
		chunkPaths[result.index] = result.path
	}

	// Transcribe chunks in parallel with progressive updates
	type transcribeResult struct {
		index         int
		transcription string
		err           error
	}

	transcribeChan := make(chan transcribeResult, numChunks)
	transcriptions := make([]string, numChunks)
	completedChunks := 0

	xlog.Info("Async: Starting parallel transcription", "chunks", numChunks)

	for i := 0; i < numChunks; i++ {
		if chunkPaths[i] == "" {
			transcribeChan <- transcribeResult{index: i, transcription: "", err: nil}
			continue
		}

		go func(idx int, audioPath string) {
			text, err := s.whisperService.Transcribe(ctx, modelName, audioPath)
			transcribeChan <- transcribeResult{index: idx, transcription: text, err: err}
		}(i, chunkPaths[i])
	}

	// Collect transcriptions with progressive updates
	for i := 0; i < numChunks; i++ {
		result := <-transcribeChan
		if result.err != nil {
			xlog.Warn("Async: Chunk transcription failed", "index", result.index, "error", result.err)
			continue
		}

		transcriptions[result.index] = result.transcription
		completedChunks++

		// Progressive update: update DB with current progress
		progress := fmt.Sprintf("*Transcribing... (%d/%d segments completed)*\n\n", completedChunks, numChunks)
		partialTranscription := s.combineTranscriptions(transcriptions)
		progressMarkdown := fmt.Sprintf("\n\n### Video Transcription (AI Generated via Whisper)\n\n%s%s", progress, partialTranscription)

		if err := s.replaceTranscriptionInDocument(ctx, documentID, userID, progressMarkdown); err != nil {
			xlog.Warn("Async: Failed to update progress", "error", err)
		} else {
			xlog.Info("Async: Progress updated", "completed", completedChunks, "total", numChunks)
		}
	}

	return s.combineTranscriptions(transcriptions), nil
}

// combineTranscriptions combines transcription segments in order
// Also deduplicates repetitive sentences common in Whisper hallucinations
func (s *FileProcessorService) combineTranscriptions(transcriptions []string) string {
	var result strings.Builder
	for i, t := range transcriptions {
		if t == "" {
			continue
		}
		if result.Len() > 0 {
			result.WriteString("\n\n")
		}
		// Deduplicate repetitive lines before adding
		// TODO: Uncomment after testing with different Whisper model
		// cleaned := deduplicateRepetitiveText(strings.TrimSpace(t))
		// result.WriteString(fmt.Sprintf("**[Segment %d]**\n%s", i+1, cleaned))
		result.WriteString(fmt.Sprintf("**[Segment %d]**\n%s", i+1, strings.TrimSpace(t)))
	}
	return result.String()
}

// deduplicateRepetitiveText removes consecutive duplicate sentences
// This helps clean up Whisper hallucinations where it repeats the same sentence
func deduplicateRepetitiveText(text string) string {
	// Split by sentence boundaries
	sentences := splitIntoSentences(text)
	if len(sentences) == 0 {
		return text
	}

	var result []string
	seen := make(map[string]int)
	maxOccurrences := 2 // Allow max 2 occurrences of same sentence

	for _, sentence := range sentences {
		trimmed := strings.TrimSpace(sentence)
		if trimmed == "" {
			continue
		}
		// Normalize for comparison
		normalized := strings.ToLower(trimmed)
		if seen[normalized] < maxOccurrences {
			result = append(result, trimmed)
			seen[normalized]++
		}
	}

	if len(result) == 0 {
		return text
	}
	return strings.Join(result, " ")
}

// splitIntoSentences splits text into sentences
func splitIntoSentences(text string) []string {
	// Replace newlines with spaces first
	text = strings.ReplaceAll(text, "\n", " ")

	// Split by sentence-ending punctuation
	var sentences []string
	var current strings.Builder

	for i, r := range text {
		current.WriteRune(r)
		// Check for sentence boundary
		if r == '.' || r == '!' || r == '?' {
			// Check if followed by space or end of text
			if i+1 >= len(text) || text[i+1] == ' ' {
				sentence := strings.TrimSpace(current.String())
				if len(sentence) > 0 {
					sentences = append(sentences, sentence)
				}
				current.Reset()
			}
		}
	}

	// Add any remaining text
	remaining := strings.TrimSpace(current.String())
	if len(remaining) > 0 {
		sentences = append(sentences, remaining)
	}

	return sentences
}

// selectWhisperModelByHardware selects the best Whisper model based on available RAM
// Model accuracy: large > medium > small > base > tiny
// IMPORTANT: Minimum recommended model is "small" for acceptable accuracy
// Model RAM requirements (approximate during inference):
//   - tiny:   ~1GB RAM   (not recommended - poor accuracy)
//   - base:   ~2GB RAM   (not recommended - poor accuracy)
//   - small:  ~4GB RAM   (minimum recommended)
//   - medium: ~8GB RAM   (good accuracy)
//   - large:  ~16GB RAM  (best accuracy)
func selectWhisperModelByHardware(availableModels []string) string {
	if len(availableModels) == 0 {
		return "small" // fallback to minimum recommended
	}

	// Detect hardware
	specs := llama.DetectHardwareSpecs()
	availableRAM := specs.AvailableRAM

	// Determine best model based on RAM
	// Minimum model is "small" for acceptable transcription quality
	var preferredModels []string
	switch {
	case availableRAM >= 16:
		// 16GB+ RAM: can use large or medium
		preferredModels = []string{"large", "large-v3", "large-v2", "medium", "medium.en", "small", "small.en"}
	case availableRAM >= 8:
		// 8-16GB RAM: use medium or small
		preferredModels = []string{"medium", "medium.en", "small", "small.en"}
	default:
		// <8GB RAM: use small (minimum recommended)
		preferredModels = []string{"small", "small.en"}
	}

	// Find first available model from preference list
	for _, preferred := range preferredModels {
		for _, available := range availableModels {
			if available == preferred {
				return available
			}
		}
	}

	// Fallback to first available
	return availableModels[0]
}

// getVideoDuration gets video duration in seconds using ffprobe
func (s *FileProcessorService) getVideoDuration(videoPath string) (float64, error) {
	ffprobePath, err := exec.LookPath("ffprobe")
	if err != nil {
		// Try ffmpeg location
		ffmpegPath, _ := exec.LookPath("ffmpeg")
		if ffmpegPath != "" {
			ffprobePath = filepath.Join(filepath.Dir(ffmpegPath), "ffprobe")
		}
		if _, statErr := os.Stat(ffprobePath); statErr != nil {
			return 0, fmt.Errorf("ffprobe not found")
		}
	}

	cmd := exec.Command(ffprobePath,
		"-v", "error",
		"-show_entries", "format=duration",
		"-of", "default=noprint_wrappers=1:nokey=1",
		videoPath,
	)

	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	var duration float64
	_, err = fmt.Sscanf(strings.TrimSpace(string(output)), "%f", &duration)
	return duration, err
}

// extractAudioChunk extracts a specific chunk of audio from video
func (s *FileProcessorService) extractAudioChunk(videoPath, outputPath string, startTime, duration float64) error {
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		homeDir, _ := os.UserHomeDir()
		localFFmpeg := filepath.Join(homeDir, ".local", "bin", "ffmpeg")
		if _, statErr := os.Stat(localFFmpeg); statErr == nil {
			ffmpegPath = localFFmpeg
		} else {
			return fmt.Errorf("ffmpeg not found: %w", err)
		}
	}

	args := []string{
		"-i", videoPath,
		"-vn",
		"-acodec", "pcm_s16le",
		"-ar", "16000",
		"-ac", "1",
		"-ss", fmt.Sprintf("%.2f", startTime),
		"-t", fmt.Sprintf("%.2f", duration),
		"-y",
		outputPath,
	}

	cmd := exec.Command(ffmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ffmpeg chunk extraction failed: %w, output: %s", err, string(output))
	}

	return nil
}

// replaceTranscriptionInDocument replaces the transcription section in document
func (s *FileProcessorService) replaceTranscriptionInDocument(ctx context.Context, documentID, userID, newContent string) error {
	// Get current document by document ID (not file ID)
	var docID string
	var currentContent sql.NullString

	err := s.db.QueryRowContext(ctx,
		"SELECT id, content FROM documents WHERE id = ?",
		documentID,
	).Scan(&docID, &currentContent)
	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	content := currentContent.String

	// Find and replace transcription section
	marker := "### Video Transcription (AI Generated via Whisper)"
	idx := strings.Index(content, marker)
	if idx >= 0 {
		// Replace from marker to end
		content = content[:idx] + strings.TrimPrefix(newContent, "\n\n")
	} else {
		// Append if not found
		content += newContent
	}

	// Update document using raw SQL
	_, err = s.db.ExecContext(ctx,
		"UPDATE documents SET content = ?, updated_at = ? WHERE id = ?",
		content, time.Now().Unix(), docID,
	)
	return err
}

// extractAudioFromVideo extracts audio from video file using ffmpeg
// If startTime and duration are both 0, extracts full audio
// Returns path to temporary WAV file (caller must clean up)
func (s *FileProcessorService) extractAudioFromVideo(videoPath string, startTime, duration float64) (string, error) {
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

	// Build ffmpeg arguments
	args := []string{
		"-i", videoPath,
		"-vn",                  // No video
		"-acodec", "pcm_s16le", // 16-bit PCM
		"-ar", "16000", // 16kHz sample rate
		"-ac", "1", // Mono
	}

	// Add time range if specified
	if startTime > 0 || duration > 0 {
		if startTime > 0 {
			args = append(args, "-ss", fmt.Sprintf("%.2f", startTime))
		}
		if duration > 0 {
			args = append(args, "-t", fmt.Sprintf("%.2f", duration))
		}
	}

	args = append(args, "-y", audioPath) // Overwrite

	cmd := exec.Command(ffmpegPath, args...)
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
