package services

import (
	"context"
	"database/sql"
	"fmt"

	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/pkg/chromem"
	einochromem "github.com/kawai-network/veridium/pkg/eino-adapters/chromem"
	"github.com/kawai-network/veridium/pkg/xlog"
)

// RAGProcessor handles RAG processing (chunking + embedding)
type RAGProcessor struct {
	queries     *db.Queries
	chromemDB   *chromem.DB
	assetDir    string
	chunkSize   int
	overlapSize int
}

// NewRAGProcessor creates a new RAG processor
func NewRAGProcessor(database *sql.DB, chromemDB *chromem.DB, assetDir string) *RAGProcessor {
	return &RAGProcessor{
		queries:     db.New(database),
		chromemDB:   chromemDB,
		assetDir:    assetDir,
		chunkSize:   1000,
		overlapSize: 200,
	}
}

// RAGProcessRequest represents a RAG processing request
type RAGProcessRequest struct {
	FilePath   string
	FileID     string
	DocumentID string
	UserID     string
	Filename   string
}

// ProcessFile processes a file for RAG (chunking + embedding)
func (r *RAGProcessor) ProcessFile(ctx context.Context, req RAGProcessRequest) ([]string, error) {
	xlog.Info("Starting RAG processing", "file_id", req.FileID, "document_id", req.DocumentID)

	// 1. Get or create user-specific chromem collection
	collectionName := fmt.Sprintf("user-%s-kb", req.UserID)
	collection, _ := r.chromemDB.GetOrCreateCollection(collectionName, nil, nil)

	// 2. Initialize eino-adapters FileManager
	indexer := einochromem.NewIndexer(collection)
	fileManager, err := einochromem.NewFileManager(ctx, &einochromem.FileManagerConfig{
		Indexer:     indexer,
		AssetDir:    fmt.Sprintf("%s/%s", r.assetDir, req.UserID),
		ChunkSize:   r.chunkSize,
		OverlapSize: r.overlapSize,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create file manager: %w", err)
	}

	// 3. Store file (parse + chunk + embed)
	err = fileManager.StoreFile(ctx, req.FilePath, map[string]any{
		"file_id":     req.FileID,
		"document_id": req.DocumentID,
		"user_id":     req.UserID,
		"filename":    req.Filename,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store file in chromem: %w", err)
	}

	xlog.Info("RAG processing completed successfully", "file_id", req.FileID, "document_id", req.DocumentID)

	// Note: Chunk metadata is stored in chromem with file_id, document_id, and user_id
	// SQLite chunks table is not used for eino-adapters approach
	// Chunks can be queried directly from chromem using metadata filters

	return []string{}, nil
}
