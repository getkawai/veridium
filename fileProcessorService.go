package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/internal/whisper"
	"github.com/kawai-network/veridium/pkg/xlog"
)

// FileProcessorService is the Wails-exposed service
type FileProcessorService struct {
	processor      *services.FileProcessorService
	libraryService *llama.LibraryService
	fileBaseDir    string // Base directory for file storage
}

// NewFileProcessorService creates a new Wails file processor service
func NewFileProcessorService(
	db *sql.DB,
	fileLoader *services.FileLoader,
	vectorSearchService *services.VectorSearchService,
	duckDB *services.DuckDBStore,
	libraryService *llama.LibraryService,
	whisperService *whisper.Service,
	fileBaseDir string,
) *FileProcessorService {
	// Get embedding function from vector search service
	embedder := vectorSearchService.GetEmbedder()

	// Create RAG processor with Eino embedder and file loader
	ragProcessor := services.NewRAGProcessor(db, duckDB, fileLoader, embedder)

	// Create file processor with whisper service for video transcription
	processor := services.NewFileProcessorService(
		db,
		fileLoader,
		ragProcessor,
		libraryService,
		whisperService,
	)

	return &FileProcessorService{
		processor:      processor,
		libraryService: libraryService,
		fileBaseDir:    fileBaseDir,
	}
}

// ProcessFileFromPath processes a file from absolute path (e.g., from file dialog)
// It copies the file to local storage and processes it for RAG
// Returns the processed file response with the relative URL for frontend display
func (f *FileProcessorService) ProcessFileFromPath(
	absolutePath string,
	userID string,
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

	xlog.Info("File copied successfully", "filename", filename, "destPath", relativeKey, "bytes", written)

	// Process file for storage and RAG
	req := services.ProcessFileRequest{
		FilePath:  destPath,
		Filename:  filename,
		UserID:    userID,
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
