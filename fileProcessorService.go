package main

import (
	"context"
	"database/sql"
	"path/filepath"

	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/pkg/chromem"
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
	chromemDB *chromem.DB,
	fileBaseDir string,
) *FileProcessorService {
	// Initialize sub-services
	documentService := services.NewDocumentService(db)
	ragProcessor := services.NewRAGProcessor(db, chromemDB, "./assets")

	// Create file processor
	processor := services.NewFileProcessorService(
		db,
		fileLoader,
		documentService,
		ragProcessor,
	)

	return &FileProcessorService{
		processor:   processor,
		fileBaseDir: fileBaseDir,
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

	req := services.ProcessFileRequest{
		FilePath:  absolutePath,
		Filename:  filename,
		FileType:  fileType,
		UserID:    userID,
		ClientID:  "", // Optional
		Source:    absolutePath,
		EnableRAG: enableRAG,
		IsShared:  false,
		FileMetadata: &services.FileMetadata{
			Filename: filename,
			FileType: fileType,
		},
	}

	return f.processor.ProcessFile(ctx, req)
}
