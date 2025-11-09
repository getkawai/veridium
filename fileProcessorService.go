package main

import (
	"context"
	"database/sql"

	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/pkg/chromem"
)

// loadFileServiceAdapter adapts *LoadFileService to services.LoadFileService interface
type loadFileServiceAdapter struct {
	service *LoadFileService
}

func (a *loadFileServiceAdapter) LoadFile(filePath string, metadata *services.FileMetadata) (*services.FileDocument, error) {
	// Convert services.FileMetadata to main.FileMetadata
	var mainMetadata *FileMetadata
	if metadata != nil {
		mainMetadata = &FileMetadata{
			Source:       metadata.Source,
			Filename:     metadata.Filename,
			FileType:     metadata.FileType,
			CreatedTime:  metadata.CreatedTime,
			ModifiedTime: metadata.ModifiedTime,
			Error:        metadata.Error,
		}
	}

	// Call the actual service
	result, err := a.service.LoadFile(filePath, mainMetadata)
	if err != nil {
		return nil, err
	}

	// Convert main.FileDocument to services.FileDocument
	servicesDoc := &services.FileDocument{
		Content:        result.Content,
		CreatedTime:    result.CreatedTime,
		FileType:       result.FileType,
		Filename:       result.Filename,
		Metadata: services.FileMetadata{
			Source:       result.Metadata.Source,
			Filename:     result.Metadata.Filename,
			FileType:     result.Metadata.FileType,
			CreatedTime:  result.Metadata.CreatedTime,
			ModifiedTime: result.Metadata.ModifiedTime,
			Error:        result.Metadata.Error,
		},
		ModifiedTime:   result.ModifiedTime,
		Pages:          convertPages(result.Pages),
		Source:         result.Source,
		TotalCharCount: result.TotalCharCount,
		TotalLineCount: result.TotalLineCount,
	}

	return servicesDoc, nil
}

func convertPages(pages []DocumentPage) []services.DocumentPage {
	result := make([]services.DocumentPage, len(pages))
	for i, page := range pages {
		result[i] = services.DocumentPage{
			CharCount:   page.CharCount,
			LineCount:   page.LineCount,
			Metadata:    page.Metadata,
			PageContent: page.PageContent,
		}
	}
	return result
}

// FileProcessorService is the Wails-exposed service
type FileProcessorService struct {
	processor *services.FileProcessorService
}

// NewFileProcessorService creates a new Wails file processor service
func NewFileProcessorService(
	db *sql.DB,
	loadFileService *LoadFileService,
	chromemDB *chromem.DB,
) *FileProcessorService {
	// Initialize sub-services
	documentService := services.NewDocumentService(db)
	ragProcessor := services.NewRAGProcessor(db, chromemDB, "./assets")

	// Create adapter
	adapter := &loadFileServiceAdapter{service: loadFileService}

	// Create file processor
	processor := services.NewFileProcessorService(
		db,
		adapter,
		documentService,
		ragProcessor,
	)

	return &FileProcessorService{
		processor: processor,
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

	req := services.ProcessFileRequest{
		FilePath:  filePath,
		Filename:  filename,
		FileType:  fileType,
		UserID:    userID,
		ClientID:  "", // Optional
		Source:    filePath,
		EnableRAG: enableRAG,
		IsShared:  false,
		FileMetadata: &services.FileMetadata{
			Filename: filename,
			FileType: fileType,
		},
	}

	return f.processor.ProcessFile(ctx, req)
}

