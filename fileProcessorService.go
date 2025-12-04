package main

import (
	"context"
	"database/sql"
	"path/filepath"
	"strings"

	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/types"
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
	fileBaseDir string,
) *FileProcessorService {
	// Initialize sub-services
	documentService := services.NewDocumentService(db)

	// Get embedding function from vector search service
	embedder := vectorSearchService.GetEmbedder()

	// Create RAG processor with Eino embedder and file loader
	ragProcessor := services.NewRAGProcessor(db, duckDB, fileLoader, embedder)

	// Create file processor
	processor := services.NewFileProcessorService(
		db,
		fileLoader,
		documentService,
		ragProcessor,
		libraryService,
	)

	return &FileProcessorService{
		processor:      processor,
		libraryService: libraryService,
		fileBaseDir:    fileBaseDir,
	}
}

// ProcessFileForStorage processes a file and saves to database
// This is called from frontend after file upload
func (f *FileProcessorService) ProcessFileForStorage(
	filePath string,
	filename string,
	fileType string,
	userID string,
	enableRAG bool,
) (*services.ProcessFileResponse, error) {
	ctx := context.Background()

	// Convert relative path to absolute path if needed
	absolutePath := filePath
	if !filepath.IsAbs(filePath) {
		absolutePath = filepath.Join(f.fileBaseDir, filePath)
	}

	// Normalize file type for images to ensure Qwen-VL processing is triggered
	// The frontend sends "image/jpeg" etc, but the internal service expects "image"
	normalizedFileType := fileType
	if strings.HasPrefix(fileType, "image/") {
		normalizedFileType = string(types.FileTypeImage)
	}

	req := services.ProcessFileRequest{
		FilePath:  absolutePath,
		Filename:  filename,
		FileType:  normalizedFileType,
		UserID:    userID,
		ClientID:  "", // Optional
		Source:    absolutePath,
		EnableRAG: enableRAG,
		IsShared:  false,
		FileMetadata: &types.FileMetadata{
			Filename: filename,
			FileType: normalizedFileType,
		},
	}

	return f.processor.ProcessFile(ctx, req)
}
