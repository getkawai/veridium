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
	"github.com/kawai-network/veridium/pkg/xlog"
)

// Type aliases for external types (from root package)
type LoadFileService interface {
	LoadFile(filePath string, metadata *FileMetadata) (*FileDocument, error)
}

type FileMetadata struct {
	Source       string
	Filename     string
	FileType     string
	CreatedTime  time.Time
	ModifiedTime time.Time
	Error        string
}

type FileDocument struct {
	Content        string
	CreatedTime    time.Time
	FileType       string
	Filename       string
	Metadata       FileMetadata
	ModifiedTime   time.Time
	Pages          []DocumentPage
	Source         string
	TotalCharCount int
	TotalLineCount int
}

// FileProcessorService orchestrates file processing pipeline
type FileProcessorService struct {
	queries         *db.Queries
	loadFileService LoadFileService
	documentService *DocumentService
	ragProcessor    *RAGProcessor
}

// NewFileProcessorService creates a new file processor service
func NewFileProcessorService(
	database *sql.DB,
	loadFileService LoadFileService,
	documentService *DocumentService,
	ragProcessor *RAGProcessor,
) *FileProcessorService {
	return &FileProcessorService{
		queries:         db.New(database),
		loadFileService: loadFileService,
		documentService: documentService,
		ragProcessor:    ragProcessor,
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
	FileMetadata *FileMetadata
}

// ProcessFileResponse represents the result of file processing
type ProcessFileResponse struct {
	FileID       string   `json:"fileId"`
	DocumentID   string   `json:"documentId"`
	ChunkIDs     []string `json:"chunkIds,omitempty"`
	GlobalFileID string   `json:"globalFileId,omitempty"`
}

// ProcessFile is the main entry point for file processing
// It handles: files → global_files (optional) → documents → chunks (optional)
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

	// Step 2: Parse file using LoadFileService
	fileDoc, err := s.loadFileService.LoadFile(req.FilePath, req.FileMetadata)
	if err != nil {
		return nil, fmt.Errorf("failed to load file: %w", err)
	}

	// Step 3: Save to documents table
	documentID, err := s.saveDocument(ctx, fileDoc, fileID, req)
	if err != nil {
		return nil, fmt.Errorf("failed to save document: %w", err)
	}
	response.DocumentID = documentID

	// Step 4: Process for RAG (optional, background)
	if req.EnableRAG {
		go func() {
			chunkIDs, err := s.ragProcessor.ProcessFile(context.Background(), RAGProcessRequest{
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
		}()
	}

	return response, nil
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
		ID:              fileID,
		UserID:          req.UserID,
		FileType:        req.FileType,
		FileHash:        sql.NullString{String: globalFileID, Valid: globalFileID != ""},
		Name:            req.Filename,
		Size:            fileInfo.Size,
		Url:             req.FilePath,
		Source:          sql.NullString{String: req.Source, Valid: req.Source != ""},
		ClientID:        sql.NullString{String: req.ClientID, Valid: req.ClientID != ""},
		Metadata:        sql.NullString{Valid: false},
		ChunkTaskID:     sql.NullString{Valid: false},
		EmbeddingTaskID: sql.NullString{Valid: false},
		CreatedAt:       now,
		UpdatedAt:       now,
	})

	if err != nil {
		return "", "", fmt.Errorf("failed to create file: %w", err)
	}

	return fileID, globalFileID, nil
}

// saveDocument saves parsed file content to documents table
func (s *FileProcessorService) saveDocument(ctx context.Context, fileDoc *FileDocument, fileID string, req ProcessFileRequest) (string, error) {
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
	pages := make([]DocumentPage, len(fileDoc.Pages))
	for i, page := range fileDoc.Pages {
		pages[i] = DocumentPage{
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
