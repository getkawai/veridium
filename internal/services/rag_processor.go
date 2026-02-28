package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"log/slog"

	llamaembed "github.com/getkawai/unillm/providers/llama-embed"
	"github.com/google/uuid"
	db "github.com/getkawai/database/db"
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
	Filename   string
}

// ProcessFile processes a file for RAG (chunking + embedding)
// Uses already-parsed content from database to avoid double parsing
func (r *RAGProcessor) ProcessFile(ctx context.Context, req RAGProcessRequest) ([]string, error) {
	slog.InfoContext(ctx, "Starting RAG processing", "file_id", req.FileID, "document_id", req.DocumentID)

	// 1. Get already-parsed document content from database
	doc, err := r.queries.GetDocument(ctx, req.DocumentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get document: %w", err)
	}

	// Check if document has content
	if !doc.Content.Valid || doc.Content.String == "" {
		slog.WarnContext(ctx, "Document has no content", "document_id", req.DocumentID)
		return []string{}, nil
	}

	// 2. Reconstruct FileDocument for chunking
	// Unmarshal pages from JSON if available (needed for PDF chunking)
	var pages []types.DocumentPage
	if doc.Pages.Valid && doc.Pages.String != "" {
		if err := json.Unmarshal([]byte(doc.Pages.String), &pages); err != nil {
			slog.WarnContext(ctx, "Failed to unmarshal pages from database", "error", err, "document_id", req.DocumentID)
			// Continue without pages - will fallback to content-based chunking
			pages = nil
		} else {
			slog.InfoContext(ctx, "Unmarshaled pages from database", "count", len(pages), "document_id", req.DocumentID)
		}
	}

	fileDoc := &types.FileDocument{
		Content:  doc.Content.String,
		FileType: doc.FileType,
		Filename: req.Filename,
		Pages:    pages, // Now properly populated for PDF chunking
	}

	// Use FileLoader's superior chunking logic (Eino, etc.)
	chunks := r.fileLoader.ChunkDocument(fileDoc, types.ChunkingConfig{
		Enabled:     true,
		ChunkSize:   r.chunkSize,
		OverlapSize: r.overlapSize,
	})
	if len(chunks) == 0 {
		slog.WarnContext(ctx, "No chunks created from document", "document_id", req.DocumentID)
		return []string{}, nil
	}

	slog.InfoContext(ctx, "Created chunks", "count", len(chunks), "document_id", req.DocumentID)

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
		})
		if err != nil {
			slog.ErrorContext(ctx, "Failed to save chunk to SQLite", "error", err, "chunk_index", i)
			continue
		}

		// Generate embedding
		embeddings, err := r.embedder.Embed(ctx, []string{chunkContent})
		if err != nil {
			slog.ErrorContext(ctx, "Failed to generate embedding", "error", err, "chunk_index", i)
			continue
		}

		if len(embeddings) == 0 || len(embeddings[0]) == 0 {
			slog.ErrorContext(ctx, "Empty embedding response", "chunk_index", i)
			continue
		}

		embedding := embeddings[0]

		// B. Save embedding to DuckDB (Vector Engine)
		if r.duckDB != nil {
			if err := r.duckDB.UpsertVector(ctx, chunkID, req.FileID, embedding); err != nil {
				slog.ErrorContext(ctx, "Failed to save vector to DuckDB", "error", err, "chunk_id", chunkID)
				// We continue even if DuckDB fails, as SQLite is the source of truth
			}
		}

		chunkIDs = append(chunkIDs, chunkID)
	}

	slog.InfoContext(ctx, "RAG processing completed successfully",
		"file_id", req.FileID,
		"document_id", req.DocumentID,
		"chunks_stored", len(chunkIDs))

	// Update file stats with chunk count and status
	chunkCount := int64(len(chunkIDs))
	chunkingStatus := "success"
	embeddingStatus := "success"
	if chunkCount == 0 {
		chunkingStatus = "empty"
		embeddingStatus = "empty"
	}

	err = r.queries.UpdateFileChunkStats(ctx, db.UpdateFileChunkStatsParams{
		ChunkCount:      sql.NullInt64{Int64: chunkCount, Valid: true},
		ChunkingStatus:  sql.NullString{String: chunkingStatus, Valid: true},
		EmbeddingStatus: sql.NullString{String: embeddingStatus, Valid: true},
		UpdatedAt:       time.Now().UnixMilli(),
		ID:              req.FileID,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Failed to update file chunk stats", "error", err, "file_id", req.FileID)
		// Don't fail the whole process if stats update fails
	} else {
		slog.InfoContext(ctx, "Updated file chunk stats", "file_id", req.FileID, "chunk_count", chunkCount)
	}

	return chunkIDs, nil
}

// DeleteFileVectors deletes all vectors associated with a file from DuckDB
func (r *RAGProcessor) DeleteFileVectors(ctx context.Context, fileID string) error {
	if r.duckDB == nil {
		return nil
	}
	return r.duckDB.DeleteVectorsByFileID(ctx, fileID)
}
