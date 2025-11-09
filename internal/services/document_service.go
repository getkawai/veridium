package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	db "github.com/kawai-network/veridium/internal/database/generated"
)

// DocumentService handles document CRUD operations
type DocumentService struct {
	queries *db.Queries
}

// NewDocumentService creates a new document service
func NewDocumentService(database *sql.DB) *DocumentService {
	return &DocumentService{
		queries: db.New(database),
	}
}

// CreateDocumentParams represents parameters for creating a document
type CreateDocumentParams struct {
	Title          string
	Content        string
	FileType       string
	Filename       string
	TotalCharCount int
	TotalLineCount int
	Metadata       map[string]interface{}
	Pages          []DocumentPage // from LoadFileService
	SourceType     string         // "file", "web", "api"
	Source         string
	FileID         string
	UserID         string
	ClientID       string
}

// DocumentPage matches LoadFileService.DocumentPage
type DocumentPage struct {
	CharCount   int                    `json:"charCount"`
	LineCount   int                    `json:"lineCount"`
	Metadata    map[string]interface{} `json:"metadata"`
	PageContent string                 `json:"pageContent"`
}

// CreateDocument creates a new document from LoadFileService output
func (s *DocumentService) CreateDocument(ctx context.Context, params CreateDocumentParams) (string, error) {
	documentID := uuid.New().String()
	now := time.Now().UnixMilli()

	// Marshal metadata to JSON
	metadataJSON, err := json.Marshal(params.Metadata)
	if err != nil {
		return "", fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Marshal pages to JSON
	pagesJSON, err := json.Marshal(params.Pages)
	if err != nil {
		return "", fmt.Errorf("failed to marshal pages: %w", err)
	}

	// Create document record
	_, err = s.queries.CreateDocument(ctx, db.CreateDocumentParams{
		ID:             documentID,
		Title:          sql.NullString{String: params.Title, Valid: params.Title != ""},
		Content:        sql.NullString{String: params.Content, Valid: params.Content != ""},
		FileType:       params.FileType,
		Filename:       sql.NullString{String: params.Filename, Valid: params.Filename != ""},
		TotalCharCount: int64(params.TotalCharCount),
		TotalLineCount: int64(params.TotalLineCount),
		Metadata:       sql.NullString{String: string(metadataJSON), Valid: true},
		Pages:          sql.NullString{String: string(pagesJSON), Valid: true},
		SourceType:     params.SourceType,
		Source:         params.Source,
		FileID:         sql.NullString{String: params.FileID, Valid: params.FileID != ""},
		UserID:         params.UserID,
		ClientID:       sql.NullString{String: params.ClientID, Valid: params.ClientID != ""},
		CreatedAt:      now,
		UpdatedAt:      now,
	})

	if err != nil {
		return "", fmt.Errorf("failed to create document: %w", err)
	}

	return documentID, nil
}

// GetDocument retrieves a document by ID
func (s *DocumentService) GetDocument(ctx context.Context, documentID, userID string) (db.Document, error) {
	return s.queries.GetDocument(ctx, db.GetDocumentParams{
		ID:     documentID,
		UserID: userID,
	})
}

// DeleteDocument deletes a document
func (s *DocumentService) DeleteDocument(ctx context.Context, documentID, userID string) error {
	return s.queries.DeleteDocument(ctx, db.DeleteDocumentParams{
		ID:     documentID,
		UserID: userID,
	})
}
