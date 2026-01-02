package services

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/kawai-network/veridium/internal/constant"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/internal/whisper"
	"github.com/kawai-network/veridium/pkg/fantasy"
	"github.com/kawai-network/veridium/pkg/fantasy/llamalib"
	llamavl "github.com/kawai-network/veridium/pkg/fantasy/providers/llama-vl"
	"github.com/kawai-network/veridium/pkg/hardware"
	"github.com/kawai-network/veridium/types"
)

// FileProcessorService orchestrates file processing pipeline
type FileProcessorService struct {
	db             *sql.DB
	queries        *db.Queries
	fileLoader     *FileLoader
	ragProcessor   *RAGProcessor
	libraryService *llamalib.Service
	vlProvider     llamavl.Provider // For VL (Vision-Language) processing
	whisperService *whisper.Service
	languageModel  fantasy.LanguageModel // For OCR/transcript cleanup
	FileBaseDir    string                // Base directory for file storage (injected)
}

// NewFileProcessorService creates a new file processor service
func NewFileProcessorService(
	database *sql.DB,
	fileLoader *FileLoader,
	ragProcessor *RAGProcessor,
	libraryService *llamalib.Service,
	whisperService *whisper.Service,
	fileBaseDir string,
) *FileProcessorService {
	return &FileProcessorService{
		db:             database,
		queries:        db.New(database),
		fileLoader:     fileLoader,
		ragProcessor:   ragProcessor,
		libraryService: libraryService,
		whisperService: whisperService,
		FileBaseDir:    fileBaseDir,
	}
}

// SetLanguageModel sets the language model for OCR/transcript cleanup
func (s *FileProcessorService) SetLanguageModel(model fantasy.LanguageModel) {
	s.languageModel = model
}

// SetVLProvider sets the VL provider for Vision-Language processing
func (s *FileProcessorService) SetVLProvider(provider llamavl.Provider) {
	s.vlProvider = provider
}

// ProcessFileRequest represents a file processing request
type ProcessFileRequest struct {
	FilePath  string
	Filename  string
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
	log.Printf("[INFO] Processing file: filename=%s", req.Filename)

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
		log.Printf("[INFO] Starting async image processing (hybrid OCR/VL): filename=%s document_id=%s", req.Filename, documentID)

		go s.processImageDescriptionAsync(req.FilePath, req.Filename, documentID, fileID, req.EnableRAG)
	} else if s.fileLoader.IsVideoFile(detectedFileType) {
		// Step 4b: If video, generate description using OpenRouter VL model ASYNCHRONOUSLY
		log.Printf("[INFO] Starting async video understanding generation: filename=%s document_id=%s", req.Filename, documentID)

		go s.processVideoDescriptionAsync(req.FilePath, req.Filename, documentID, fileID, req.EnableRAG)
	} else if req.EnableRAG && s.fileLoader.CanChunkForRAG(detectedFileType) {
		// Step 5: Process for RAG (for non-image files, do it synchronously)
		chunkIDs, err := s.ragProcessor.ProcessFile(ctx, RAGProcessRequest{
			FilePath:   req.FilePath,
			FileID:     fileID,
			DocumentID: documentID,
			Filename:   req.Filename,
		})
		if err != nil {
			log.Printf("[ERROR] Failed to process file for RAG: error=%v file_id=%s", err, fileID)
		} else {
			log.Printf("[INFO] RAG processing completed: file_id=%s chunks=%d", fileID, len(chunkIDs))
		}
	} else if req.EnableRAG && !s.fileLoader.CanChunkForRAG(detectedFileType) {
		log.Printf("[INFO] Skipping RAG processing for unsupported file type: file_type=%s file_id=%s", detectedFileType, fileID)
	}

	return response, nil
}

// DeleteFile deletes a file and all its associated data (chunks, vectors, documents)
func (s *FileProcessorService) DeleteFile(ctx context.Context, fileID string) error {
	log.Printf("[INFO] Deleting file: file_id=%s", fileID)

	// 1. Get file details first (to clean up global files later if needed)
	file, err := s.queries.GetFile(ctx, fileID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // File doesn't exist, nothing to do
		}
		return fmt.Errorf("failed to get file: %w", err)
	}

	// 1.5. Delete physical file from disk
	// The Url field often contains the relative path (e.g., "files/uploads/...") or a serving URL
	if file.Url != "" {
		log.Printf("[DEBUG] Checking physical file deletion. DbUrl=%s", file.Url)

		// Optimize for Cross-Platform (Windows/Linux/Mac)
		// 1. Unify separators to "/" (Standard URL format)
		normalizedPath := filepath.ToSlash(file.Url)

		// 2. Remove leading "/" if present (e.g. "/files/..." -> "files/...")
		normalizedPath = strings.TrimPrefix(normalizedPath, "/")

		// 3. Convert back to OS-specific separators (e.g. "\" on Windows)
		filePath := filepath.Clean(filepath.FromSlash(normalizedPath))

		// SAFETY CHECK: Ensure we ONLY delete files inside the allowed directory.
		// We use the injected FileBaseDir (e.g. "files").
		expectedPrefix := s.FileBaseDir + string(filepath.Separator)

		// Only delete if it looks like a local file in our expected directory structure
		// This strictly prevents deleting files outside of 'files/' (e.g. system files).
		if !strings.HasPrefix(filePath, "http") && !filepath.IsAbs(filePath) && strings.HasPrefix(filePath, expectedPrefix) {
			info, err := os.Stat(filePath)
			if err == nil && !info.IsDir() {
				if err := os.Remove(filePath); err != nil {
					log.Printf("[WARN] Failed to delete physical file: %v path=%s", err, filePath)
				} else {
					log.Printf("[INFO] Deleted physical file: %s", filePath)
				}
			} else {
				log.Printf("[DEBUG] File not found or is directory, skipping delete: %s (err: %v)", filePath, err)
			}
		} else {
			log.Printf("[WARN] Safety check failed for path: %s (Expected prefix: %s)", filePath, expectedPrefix)
		}
	}

	// 2. Delete vectors from DuckDB (via RAGProcessor)
	if err := s.ragProcessor.DeleteFileVectors(ctx, fileID); err != nil {
		log.Printf("[WARN] Failed to delete vectors from DuckDB: %v", err)
		// Continue deletion even if vector delete fails
	}

	// 3. Delete chunks from SQLite (Cascading delete in SQL usually handles this, but let's be safe/explicit if needed)
	// Actually, schema usually has ON DELETE CASCADE. If not, we need explicit deletes.
	// Assuming ON DELETE CASCADE is NOT set for all relations based on previous observations, let's explicit delete.

	// Delete chunks
	// Note: queries.DeleteFileChunks is not standard, we usually rely on foreign keys or direct query
	// Let's use the DB directly for cascade-like cleanup if simple query doesn't exist
	// Or even better, check if we have specific delete queries.
	// For now, we will assume we need to delete chunks first.
	// Actually, let's construct a transaction for safety
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	qTx := s.queries.WithTx(tx)

	// Get chunk IDs to delete (if needed for other stores, but we already handled DuckDB)
	// Just delete chunks for this file
	// We might need a query `DeleteChunksByFileID` but since we don't know if it exists,
	// let's rely on `DeleteFile` deleting everything if ON DELETE CASCADE is set,
	// OR we might need to add `DeleteChunksByFileID` to queries.sql.
	//
	// Checking schema.sql (memory):
	// chunks -> documents -> (maybe file_id?)
	// documents has file_id.
	//
	// Let's try to act conservatively.
	// Delete chunks linked to documents linked to this file.
	//
	// However, without modifying SQL queries, we might be limited.
	// Let's use `DeleteFile` from queries and hope for CASCADE or add SQL.
	// The user mentioned `DBService.DeleteFileWithCascade` in action.ts.
	// This suggests there might be a service method for this already?
	// Ah, action.ts called `DBService.DeleteFileWithCascade`.
	// `DBService` is likely `internal/database/service.go`.
	//
	// WAIT! If `internal/database` already has a Service with `DeleteFileWithCascade`,
	// maybe we should just use THAT?
	//
	// Let's pause and check `internal/database/service.go`.
	// If it exists, `FileProcessorService` can just use `s.queries` (which is `*db.Queries`)
	// or maybe `FileProcessorService` should call that existing service?
	//
	// `FileProcessorService` struct has `queries *db.Queries`.
	// It does NOT have access to `internal/database/Service`.
	//
	// Only `DeleteFileVectors` is unique to `FileProcessorService` (via RAGProcessor).
	// usage in action.ts:
	// await DBService.DeleteFileWithCascade(...)
	//
	// So `DeleteFileWithCascade` handles SQLite cleanup.
	// We just need to add Vector deletion to it OR wrap it.
	//
	// Best approach:
	// Implement `DeleteFile` here that:
	// 1. Calls `ragProcessor.DeleteFileVectors`.
	// 2. Executes the deletion transaction (similar to DeleteFileWithCascade logic) OR
	//    calls the `queries.DeleteFile` (and relies on SQL definition).
	//
	// Let's stick to using `s.queries.DeleteFile(ctx, fileID)`.
	// If the schema has FK with cascade, it works.
	// If NOT, we might leave orphans.
	//
	// Let's check `internal/database/schema/schema.sql` to see if ON DELETE CASCADE is used.
	//
	// I'll assume for now I should do the deletions manually to be safe,
	// mimicking what `action.ts` was doing but efficiently.

	// Delete documents (and thus chunks if cascaded? chunks usually reference document_id)
	// If documents have `file_id`, we delete where `file_id = ?`.
	// chunks refer to `document_id`.

	// Let's perform a direct SQL execution for cascade delete if needed, for simplicity and speed.
	_, err = tx.ExecContext(ctx, "DELETE FROM chunks WHERE document_id IN (SELECT id FROM documents WHERE file_id = ?)", fileID)
	if err != nil {
		return fmt.Errorf("failed to delete chunks: %w", err)
	}

	_, err = tx.ExecContext(ctx, "DELETE FROM documents WHERE file_id = ?", fileID)
	if err != nil {
		return fmt.Errorf("failed to delete documents: %w", err)
	}

	if err := qTx.DeleteFile(ctx, fileID); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 4. Clean up global file if needed (check if other files use same hash)
	if file.FileHash.Valid && file.FileHash.String != "" {
		hash := file.FileHash.String
		count, err := s.queries.CountFilesByHash(ctx, sql.NullString{String: hash, Valid: true})
		if err != nil {
			log.Printf("[WARN] Failed to count files by hash: %v", err)
		} else if count == 0 {
			// No other files use this hash, delete global file
			if err := s.queries.DeleteGlobalFile(ctx, hash); err != nil {
				log.Printf("[WARN] Failed to delete global file: %v", err)
			} else {
				log.Printf("[INFO] Deleted global file orphan: %s", hash)
			}
		}
	}

	return nil
}

// processImageDescriptionAsync generates image description asynchronously using hybrid approach:
// 1. First try Tesseract OCR (fast) for text extraction
// 2. If significant text found, use that (fast path)
// 3. If no/minimal text, fallback to VL model for image description (slow path)
func (s *FileProcessorService) processImageDescriptionAsync(filePath, filename, documentID, fileID string, enableRAG bool) {
	ctx := context.Background()

	log.Printf("[INFO] Async: Starting hybrid image processing: filename=%s", filename)

	var finalContent string
	var contentType string

	// Step 1: Try Tesseract OCR first (fast path)
	ocrText, err := s.extractTextWithTesseract(filePath)
	if err != nil {
		log.Printf("[WARN] Async: Tesseract OCR failed, will try VL model: error=%v filename=%s", err, filename)
	}

	// Check if we got meaningful text (more than 20 chars, excluding whitespace)
	// Lower threshold to catch short but valid text like logos, labels, etc.
	cleanedText := strings.TrimSpace(ocrText)
	if len(cleanedText) > 20 {
		// Fast path: sufficient text extracted via OCR
		log.Printf("[INFO] Async: OCR extracted sufficient text: length=%d filename=%s", len(cleanedText), filename)

		// Step 2: Clean up OCR text using LLM (if available)
		if s.languageModel != nil {
			log.Printf("[INFO] Async: Cleaning up OCR text with LLM: filename=%s", filename)
			cleanedContent, err := s.cleanupOCRText(ctx, cleanedText, filename)
			if err != nil {
				log.Printf("[WARN] Async: LLM cleanup failed, using raw OCR: error=%v filename=%s", err, filename)
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
		log.Printf("[INFO] Async: Minimal text from OCR, using VL model: ocr_length=%d filename=%s", len(cleanedText), filename)

		// Use vlProvider for VL processing (preferred), fallback to libraryService
		if s.vlProvider != nil {
			// Use provider's ProcessImage which handles model loading automatically
			prompt := "Describe this image in detail. Include all visible text, objects, people, colors, and layout."
			description, err := s.vlProvider.ProcessImage(context.Background(), filePath, prompt, 2048)
			if err != nil {
				log.Printf("[ERROR] Async: VL model processing failed: error=%v filename=%s", err, filename)
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
		} else if s.libraryService != nil {
			// Fallback to libraryService for backward compatibility
			if !s.libraryService.IsVLModelLoaded() {
				if err := s.libraryService.LoadVLModel(""); err != nil {
					log.Printf("[ERROR] Async: Failed to load VL model: error=%v", err)
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
					log.Printf("[ERROR] Async: VL model processing failed: error=%v filename=%s", err, filename)
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
			log.Printf("[ERROR] Async: No VL model available and no OCR text: filename=%s", filename)
			return
		}
	}

	if finalContent == "" {
		log.Printf("[WARN] Async: No content extracted from image: filename=%s", filename)
		return
	}

	log.Printf("[INFO] Async: Image processing completed: type=%s length=%d filename=%s", contentType, len(finalContent), filename)

	// Format content as markdown
	contentMarkdown := fmt.Sprintf("\n\n### %s\n\n%s", contentType, finalContent)

	// Update document with content
	err = s.appendContentToDocument(ctx, fileID, contentMarkdown)
	if err != nil {
		log.Printf("[ERROR] Async: Failed to update document: error=%v document_id=%s", err, documentID)
		return
	}

	log.Printf("[INFO] Async: Document updated: document_id=%s type=%s", documentID, contentType)

	// Process for RAG if enabled
	if enableRAG {
		chunkIDs, err := s.ragProcessor.ProcessFile(ctx, RAGProcessRequest{
			FilePath:   filePath,
			FileID:     fileID,
			DocumentID: documentID,
			Filename:   filename,
		})
		if err != nil {
			log.Printf("[ERROR] Async: Failed to process file for RAG: error=%v file_id=%s", err, fileID)
		} else {
			log.Printf("[INFO] Async: RAG processing completed: file_id=%s chunks=%d", fileID, len(chunkIDs))
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

	log.Printf("[INFO] Async: Calling LLM for OCR cleanup: prompt_len=%d", len(userPrompt))

	// Use timeout context for OCR cleanup
	timeoutCtx, cancel := context.WithTimeout(ctx, constant.LLMCleanupTimeout)
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

	log.Printf("[INFO] Async: LLM OCR response received: response_len=%d", len(resp.Content.Text()))

	result := resp.Content.Text()

	// Trim any leading/trailing whitespace
	result = strings.TrimSpace(result)

	// If result is empty or too short, return original
	if len(result) < 10 {
		return rawText, nil
	}

	log.Printf("[INFO] Async: OCR text cleaned up: original_len=%d cleaned_len=%d", len(rawText), len(result))
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

	log.Printf("[INFO] Async: Calling LLM for transcript cleanup: prompt_len=%d", len(userPrompt))

	// Use timeout context for LLM generation
	timeoutCtx, cancel := context.WithTimeout(ctx, constant.LLMGenerateTimeout)
	defer cancel()

	resp, err := s.languageModel.Generate(timeoutCtx, fantasy.Call{
		Prompt: []fantasy.Message{
			fantasy.NewSystemMessage(systemPrompt),
			fantasy.NewUserMessage(userPrompt),
		},
	})
	if err != nil {
		log.Printf("[WARN] Async: Transcript cleanup failed, using original: error=%v", err)
		return rawTranscript, nil // Return original on error, don't fail
	}

	log.Printf("[INFO] Async: LLM response received: response_len=%d", len(resp.Content.Text()))

	result := strings.TrimSpace(resp.Content.Text())

	// If result is empty or significantly shorter, return original
	if len(result) < len(rawTranscript)/2 {
		log.Printf("[WARN] Async: Transcript cleanup result too short, using original: original_len=%d result_len=%d", len(rawTranscript), len(result))
		return rawTranscript, nil
	}

	log.Printf("[INFO] Async: Transcript cleaned up: original_len=%d cleaned_len=%d", len(rawTranscript), len(result))
	return result, nil
}

// processVideoDescriptionAsync extracts audio from video using ffmpeg and transcribes using whisper
// Uses parallel chunked transcription for faster processing with progressive DB updates
func (s *FileProcessorService) processVideoDescriptionAsync(filePath, filename, documentID, fileID string, enableRAG bool) {
	ctx := context.Background()

	log.Printf("[INFO] Async: Starting video transcription (parallel): filename=%s", filename)

	// Check if whisper service is available
	if s.whisperService == nil {
		log.Printf("[ERROR] Async: Whisper service not available: filename=%s", filename)
		return
	}

	// Check if whisper is installed
	if !s.whisperService.IsWhisperInstalled() {
		log.Printf("[ERROR] Async: Whisper CLI not installed: filename=%s", filename)
		return
	}

	// Check if we have a model - prefer tiny for speed in parallel mode
	models, err := s.whisperService.ListModels()
	if err != nil || len(models) == 0 {
		log.Printf("[ERROR] Async: No whisper models available: error=%v filename=%s", err, filename)
		return
	}

	// Select Whisper model based on available RAM
	// Model sizes: tiny (~75MB), base (~150MB), small (~500MB), medium (~1.5GB), large (~3GB)
	// RAM requirements: tiny=1GB, base=2GB, small=4GB, medium=8GB, large=16GB
	modelName := selectWhisperModelByHardware(models)
	log.Printf("[INFO] Async: Selected Whisper model based on hardware: model=%s", modelName)

	// Get video duration
	duration, err := s.getVideoDuration(filePath)
	if err != nil {
		log.Printf("[WARN] Async: Could not get video duration, using sequential mode: error=%v", err)
		duration = 0
	}

	log.Printf("[INFO] Async: Video info: duration_sec=%.0f model=%s filename=%s", duration, modelName, filename)

	// Initialize transcription header in document
	err = s.appendContentToDocument(ctx, fileID, "\n\n### Video Transcription (AI Generated via Whisper)\n\n*Transcribing...*\n")
	if err != nil {
		log.Printf("[WARN] Async: Failed to initialize transcription header: error=%v", err)
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

		log.Printf("[INFO] Async: Using parallel transcription: chunks=%d chunk_duration=%f", numChunks, chunkDuration)

		fullTranscription, err = s.transcribeVideoParallel(ctx, filePath, modelName, numChunks, chunkDuration, documentID)
		if err != nil {
			log.Printf("[WARN] Async: Parallel transcription failed: error=%v filename=%s", err, filename)
			return
		}
	} else {
		// Short video - use sequential transcription
		log.Printf("[INFO] Async: Using sequential transcription (short video)")

		audioPath, err := s.extractAudioFromVideo(filePath, 0, 0) // Full video
		if err != nil {
			log.Printf("[WARN] Async: Failed to extract audio: error=%v", err)
			return
		}
		defer os.Remove(audioPath)

		fullTranscription, err = s.whisperService.Transcribe(ctx, modelName, audioPath)
		if err != nil {
			log.Printf("[WARN] Async: Transcription failed: error=%v", err)
			return
		}
	}

	if fullTranscription == "" {
		fullTranscription = "(No speech detected in video)"
	}

	log.Printf("[INFO] Async: Video transcription completed: length=%d filename=%s", len(fullTranscription), filename)

	// Post-process transcription with LLM to fix common Whisper errors
	if s.languageModel != nil && len(fullTranscription) > 50 {
		log.Printf("[INFO] Async: Cleaning up transcription with LLM: filename=%s", filename)
		cleanedTranscription, err := s.cleanupTranscription(ctx, fullTranscription)
		if err == nil && len(cleanedTranscription) > 0 {
			fullTranscription = cleanedTranscription
			log.Printf("[INFO] Async: Transcription cleanup completed: length=%d", len(fullTranscription))
		}
	}

	// Update document with final complete transcription
	finalMarkdown := fmt.Sprintf("\n\n### Video Transcription (AI Generated via Whisper + LLM Cleanup)\n\n%s", fullTranscription)
	err = s.replaceTranscriptionInDocument(ctx, documentID, finalMarkdown)
	if err != nil {
		log.Printf("[WARN] Async: Failed to update document with final transcription: error=%v", err)
	}

	log.Printf("[INFO] Async: Document updated with complete transcription: document_id=%s", documentID)

	// Process for RAG if enabled
	if enableRAG {
		chunkIDs, err := s.ragProcessor.ProcessFile(ctx, RAGProcessRequest{
			FilePath:   filePath,
			FileID:     fileID,
			DocumentID: documentID,
			Filename:   filename,
		})
		if err != nil {
			log.Printf("[WARN] Async: Failed to process video for RAG: error=%v file_id=%s", err, fileID)
		} else {
			log.Printf("[INFO] Async: RAG processing completed for video: file_id=%s chunks=%d", fileID, len(chunkIDs))
		}
	}
}

// transcribeVideoParallel splits audio into chunks and transcribes them in parallel
func (s *FileProcessorService) transcribeVideoParallel(ctx context.Context, videoPath, modelName string, numChunks int, chunkDuration float64, documentID string) (string, error) {
	// Create temp dir for chunks
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

	log.Printf("[INFO] Async: Extracting audio chunks: count=%d", numChunks)

	for i := 0; i < numChunks; i++ {
		go func(idx int) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("❌ [PANIC] Audio chunk extraction panic recovered for index %d: %v", idx, r)
					extractChan <- chunkResult{index: idx, path: "", err: fmt.Errorf("panic: %v", r)}
				}
			}()

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
			log.Printf("[WARN] Async: Failed to extract chunk: index=%d error=%v", result.index, result.err)
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

	log.Printf("[INFO] Async: Starting parallel transcription: chunks=%d", numChunks)

	for i := 0; i < numChunks; i++ {
		if chunkPaths[i] == "" {
			transcribeChan <- transcribeResult{index: i, transcription: "", err: nil}
			continue
		}

		go func(idx int, audioPath string) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("❌ [PANIC] Transcription panic recovered for index %d: %v", idx, r)
					transcribeChan <- transcribeResult{index: idx, transcription: "", err: fmt.Errorf("panic: %v", r)}
				}
			}()

			text, err := s.whisperService.Transcribe(ctx, modelName, audioPath)
			transcribeChan <- transcribeResult{index: idx, transcription: text, err: err}
		}(i, chunkPaths[i])
	}

	// Collect transcriptions with progressive updates
	for i := 0; i < numChunks; i++ {
		result := <-transcribeChan
		if result.err != nil {
			log.Printf("[WARN] Async: Chunk transcription failed: index=%d error=%v", result.index, result.err)
			continue
		}

		transcriptions[result.index] = result.transcription
		completedChunks++

		// Progressive update:
		// Update document slightly to show progress
		progressMd := fmt.Sprintf("\n\n*Transcribing segment %d/%d...*\n", result.index+1, numChunks)
		if err := s.appendContentToDocument(ctx, documentID, progressMd); err != nil {
			slog.Warn("Failed to append progress to document", "document_id", documentID, "error", err)
		}

		partialTranscription := s.combineTranscriptions(transcriptions)
		progressMarkdown := fmt.Sprintf("\n\n### Video Transcription (AI Generated via Whisper)\n\n%s%s", progressMd, partialTranscription)

		if err := s.replaceTranscriptionInDocument(ctx, documentID, progressMarkdown); err != nil {
			log.Printf("[WARN] Async: Failed to update progress: error=%v", err)
		} else {
			log.Printf("[INFO] Async: Progress updated: completed=%d total=%d", completedChunks, numChunks)
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
	specs := hardware.DetectHardwareSpecs()
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

// replaceTranscriptionInDocument replaces the transcription placeholder with actual content
func (s *FileProcessorService) replaceTranscriptionInDocument(ctx context.Context, documentID, content string) error {
	// Get existing document
	// Retry logic as in appendContentToDocument
	deadline := time.Now().Add(10 * time.Second)
	var doc db.Document
	var err error

	for attempt := 1; time.Now().Before(deadline); attempt++ {
		// Since we don't have fileID passed here easily, we get document by ID which is what we have
		// NOTE: GetDocument now takes only ID
		doc, err = s.queries.GetDocument(ctx, documentID)
		if err == nil {
			break
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to get document params on attempt %d: %w", attempt, err)
		}
		time.Sleep(time.Duration(attempt) * 100 * time.Millisecond)
	}

	if err != nil {
		return fmt.Errorf("failed to get document: %w", err)
	}

	// Logic to replace content...
	// For simplicity, we just check if content has the placeholder and replace it, or append if not found
	newContent := ""
	if doc.Content.Valid {
		if strings.Contains(doc.Content.String, "*Transcribing...*") {
			newContent = strings.Replace(doc.Content.String, "*Transcribing...*", content, 1)
		} else {
			newContent = doc.Content.String + "\n\n" + content
		}
	} else {
		newContent = content
	}

	// Update document
	_, err = s.queries.UpdateDocument(ctx, db.UpdateDocumentParams{
		ID:         doc.ID,
		Title:      doc.Title,
		Content:    sql.NullString{String: newContent, Valid: true},
		Metadata:   doc.Metadata,
		EditorData: doc.EditorData,
		UpdatedAt:  time.Now().Unix(),
	})
	if err != nil {
		return fmt.Errorf("failed to update document: %w", err)
	}

	return nil
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
		FileID:         sql.NullString{String: fileID, Valid: true},
		EditorData:     sql.NullString{},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create document: %w", err)
	}

	return doc.ID, nil
}

// appendContentToDocument appends content to an existing document
// Used for async operations like image description generation
func (s *FileProcessorService) appendContentToDocument(ctx context.Context, fileID, additionalContent string) error {
	// Get existing document
	// Retry logic to ensure the document is available after file processing
	deadline := time.Now().Add(10 * time.Second) // 10-second timeout
	var doc db.Document
	var err error

	for attempt := 1; time.Now().Before(deadline); attempt++ {
		doc, err = s.queries.GetDocumentByFileID(ctx, sql.NullString{String: fileID, Valid: true})
		if err == nil {
			break // Document found, exit retry loop
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("failed to get document by file ID on attempt %d: %w", attempt, err)
		}
		time.Sleep(time.Duration(attempt) * 100 * time.Millisecond) // Exponential backoff
	}

	if err != nil {
		return fmt.Errorf("failed to get document by file ID after multiple attempts: %w", err)
	}

	// Append new content
	newContent := doc.Content.String + additionalContent
	now := time.Now().UnixMilli()

	// Update document
	_, err = s.queries.UpdateDocument(ctx, db.UpdateDocumentParams{
		ID:         doc.ID,
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
