package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"log/slog"

	unillm "github.com/getkawai/unillm"
	llamaembed "github.com/getkawai/unillm/providers/llama-embed"
	"github.com/kawai-network/veridium/internal/services"
)

// FileProcessorService is the Wails-exposed service
type FileProcessorService struct {
	processor   *services.FileProcessorService
	fileBaseDir string // Base directory for file storage
}

// NewFileProcessorService creates a new Wails file processor service
func NewFileProcessorService(
	db *sql.DB,
	fileLoader *services.FileLoader,
	vectorSearchService *services.VectorSearchService,
	duckDB *services.DuckDBStore,
	fileBaseDir string,
) *FileProcessorService {
	// Get embedding function from vector search service (may be nil if embedder failed)
	var embedder llamaembed.Embedder
	if vectorSearchService != nil {
		embedder = vectorSearchService.GetEmbedder()
	}

	// Create RAG processor with Eino embedder and file loader
	ragProcessor := services.NewRAGProcessor(db, duckDB, fileLoader, embedder)

	// Create file processor service
	processor := services.NewFileProcessorService(
		db,
		fileLoader,
		ragProcessor,
		fileBaseDir,
	)

	return &FileProcessorService{
		processor:   processor,
		fileBaseDir: fileBaseDir,
	}
}

// SetLanguageModel sets the language model for OCR/transcript cleanup
//
//wails:ignore
func (f *FileProcessorService) SetLanguageModel(model unillm.LanguageModel) {
	f.processor.SetLanguageModel(model)
}

// ProcessFileFromPath processes a file from absolute path (e.g., from file dialog)
// It copies the file to local storage and processes it for RAG
func (f *FileProcessorService) ProcessFileFromPath(
	absolutePath string,
) (*ProcessFileFromPathResponse, error) {
	// TODO: Replace with valid userid
	ctx := context.Background()

	// Security check: ensure path is absolute and exists
	if !filepath.IsAbs(absolutePath) {
		return nil, fmt.Errorf("path must be absolute: %s", absolutePath)
	}

	// Check if file exists
	if _, err := os.Stat(absolutePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file does not exist: %s", absolutePath)
	}

	// Extract filename
	filename := filepath.Base(absolutePath)

	// Generate unique filename with timestamp
	timestamp := time.Now().UnixMilli()
	uniqueFileName := fmt.Sprintf("%d-%s", timestamp, filename)
	relativeKey := filepath.Join("uploads", uniqueFileName)

	// Construct destination path
	destPath := filepath.Join(f.fileBaseDir, relativeKey)

	// Ensure directory exists
	destDir := filepath.Dir(destPath)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create directory %s: %w", destDir, err)
	}

	// Copy file
	sourceFile, err := os.Open(absolutePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open source file: %w", err)
	}
	defer sourceFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create destination file: %w", err)
	}
	defer destFile.Close()

	// Copy content
	written, err := io.Copy(destFile, sourceFile)
	if err != nil {
		return nil, fmt.Errorf("failed to copy file: %w", err)
	}

	slog.InfoContext(ctx, "File copied successfully", "filename", filename, "destPath", relativeKey, "bytes", written)

	// Process file for storage and RAG
	req := services.ProcessFileRequest{
		FilePath:  destPath,
		Filename:  filename,
		Source:    destPath,
		EnableRAG: true,
		IsShared:  false,
	}

	result, err := f.processor.ProcessFile(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to process file: %w", err)
	}

	return &ProcessFileFromPathResponse{
		FileID:       result.FileID,
		DocumentID:   result.DocumentID,
		ChunkIDs:     result.ChunkIDs,
		GlobalFileID: result.GlobalFileID,
		Processing:   result.Processing,
		RelativeURL:  "/files/" + relativeKey,
		Filename:     filename,
	}, nil
}

// ProcessFileFromPathResponse represents the result of processing a file from absolute path
type ProcessFileFromPathResponse struct {
	FileID       string   `json:"fileId"`
	DocumentID   string   `json:"documentId"`
	ChunkIDs     []string `json:"chunkIds,omitempty"`
	GlobalFileID string   `json:"globalFileId,omitempty"`
	Processing   bool     `json:"processing,omitempty"`
	RelativeURL  string   `json:"relativeUrl"` // URL for frontend to display (e.g., /files/uploads/123-file.pdf)
	Filename     string   `json:"filename"`    // Original filename
}

// RemoveFiles deletes multiple files and their associated data
// This replaces the complex frontend logic with a robust backend implementation
func (f *FileProcessorService) RemoveFiles(ids []string) error {
	ctx := context.Background()
	slog.InfoContext(ctx, "RemoveFiles called", "count", len(ids))

	for _, id := range ids {
		if err := f.processor.DeleteFile(ctx, id); err != nil {
			slog.ErrorContext(ctx, "Failed to delete file", "id", id, "error", err)
			// We continue deleting other files instead of returning immediately
			// but we could collect errors if strictly needed.
			// For UI, best effort is usually preferred.
		}
	}

	return nil
}
