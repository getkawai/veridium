package services

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/pkg/xlog"
	llamaembed "github.com/kawai-network/veridium/fantasy/providers/llama-embed"
	"github.com/kawai-network/veridium/types"
)

// RAGProcessor handles RAG processing (chunking + embedding)
type RAGProcessor struct {
	queries     *db.Queries
	duckDB      *DuckDBStore
	fileLoader  *FileLoader
	embedder    llamaembed.Embedder
	chunkSize   int
	overlapSize int
}

// NewRAGProcessor creates a new RAG processor
func NewRAGProcessor(database *sql.DB, duckDB *DuckDBStore, fileLoader *FileLoader, embedder llamaembed.Embedder) *RAGProcessor {
	return &RAGProcessor{
		queries:     db.New(database),
		duckDB:      duckDB,
		fileLoader:  fileLoader,
		embedder:    embedder,
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
// Uses already-parsed content from database to avoid double parsing
func (r *RAGProcessor) ProcessFile(ctx context.Context, req RAGProcessRequest) ([]string, error) {
	xlog.Info("Starting RAG processing", "file_id", req.FileID, "document_id", req.DocumentID)

	// 1. Get already-parsed document content from database
	doc, err := r.queries.GetDocument(ctx, db.GetDocumentParams{
		ID:     req.DocumentID,
		UserID: req.UserID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	// Check if document has content
	if !doc.Content.Valid || doc.Content.String == "" {
		xlog.Warn("Document has no content", "document_id", req.DocumentID)
		return []string{}, nil
	}

	// 2. Reconstruct FileDocument for chunking
	// We need to unmarshal pages if available (for PDF chunking)
	// TODO: Unmarshal pages from doc.Pages (JSON) if needed for PDF
	// For now, we'll rely on content-based chunking for non-PDFs or if pages missing

	fileDoc := &types.FileDocument{
		Content:  doc.Content.String,
		FileType: doc.FileType, // FileType is string, not sql.NullString
		Filename: req.Filename,
		// Pages: pages, // Populate if we implement JSON unmarshal
	}

	// Use FileLoader's superior chunking logic (Eino, etc.)
	chunks := r.fileLoader.ChunkDocument(fileDoc, types.ChunkingConfig{
		Enabled:     true,
		ChunkSize:   r.chunkSize,
		OverlapSize: r.overlapSize,
	})
	if len(chunks) == 0 {
		xlog.Warn("No chunks created from document", "document_id", req.DocumentID)
		return []string{}, nil
	}

	xlog.Info("Created chunks", "count", len(chunks), "document_id", req.DocumentID)

	// 4. Generate embeddings and store in both SQLite and DuckDB (and Chromem for backward compat/fallback)
	chunkIDs := make([]string, 0, len(chunks))
	for i, chunk := range chunks {
		chunkID := uuid.New().String()
		chunkContent := chunk.Content

		// A. Save chunk to SQLite (Source of Truth)
		_, err := r.queries.CreateChunk(ctx, db.CreateChunkParams{
			ID:         chunkID,
			DocumentID: sql.NullString{String: req.DocumentID, Valid: true},
			Text:       sql.NullString{String: chunkContent, Valid: true},
			Metadata:   sql.NullString{String: fmt.Sprintf(`{"filename": "%s", "chunk_index": %d, "type": "%s"}`, req.Filename, i, chunk.Metadata["type"]), Valid: true},
			ChunkIndex: sql.NullInt64{Int64: int64(i), Valid: true},
			Type:       sql.NullString{String: "text", Valid: true},
			UserID:     sql.NullString{String: req.UserID, Valid: true},
		})
		if err != nil {
			xlog.Error("Failed to save chunk to SQLite", "error", err, "chunk_index", i)
			continue
		}

		// Generate embedding
		embeddings, err := r.embedder.Embed(ctx, []string{chunkContent})
		if err != nil {
			xlog.Error("Failed to generate embedding", "error", err, "chunk_index", i)
			continue
		}

		if len(embeddings) == 0 || len(embeddings[0]) == 0 {
			xlog.Error("Empty embedding response", "chunk_index", i)
			continue
		}

		embedding := embeddings[0]

		// B. Save embedding to DuckDB (Vector Engine)
		if r.duckDB != nil {
			if err := r.duckDB.UpsertVector(ctx, chunkID, req.FileID, embedding); err != nil {
				xlog.Error("Failed to save vector to DuckDB", "error", err, "chunk_id", chunkID)
				// We continue even if DuckDB fails, as SQLite is the source of truth
			}
		}

		chunkIDs = append(chunkIDs, chunkID)
	}

	xlog.Info("RAG processing completed successfully",
		"file_id", req.FileID,
		"document_id", req.DocumentID,
		"chunks_stored", len(chunkIDs))

	return chunkIDs, nil
}
