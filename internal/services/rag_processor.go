package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/pkg/xlog"
	"github.com/kawai-network/veridium/pkg/yzma/embedding"
)

// RAGProcessor handles RAG processing (chunking + embedding)
type RAGProcessor struct {
	queries     *db.Queries
	duckDB      *DuckDBStore
	fileLoader  *FileLoader
	embedder    embedding.Embedder
	chunkSize   int
	overlapSize int
}

// NewRAGProcessor creates a new RAG processor
func NewRAGProcessor(database *sql.DB, duckDB *DuckDBStore, fileLoader *FileLoader, embedder embedding.Embedder) *RAGProcessor {
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

	fileDoc := &FileDocument{
		Content:  doc.Content.String,
		FileType: doc.FileType, // FileType is string, not sql.NullString
		Filename: req.Filename,
		// Pages: pages, // Populate if we implement JSON unmarshal
	}

	// Configure chunking
	chunkConfig := ChunkingConfig{
		Enabled:     true,
		ChunkSize:   r.chunkSize,
		OverlapSize: r.overlapSize,
	}

	// Use FileLoader's superior chunking logic (Eino, etc.)
	chunks := r.fileLoader.ChunkDocument(fileDoc, chunkConfig)
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
			ClientID:   sql.NullString{String: "", Valid: false}, // Optional
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

// chunkText splits text into overlapping chunks using recursive splitting strategy
// Inspired by CloudWeGo Eino's recursive splitter for production-grade text chunking
func (r *RAGProcessor) chunkText(text string) []string {
	if text == "" {
		return []string{}
	}

	// Define separators in order of preference (coarse to fine)
	// Try to keep semantic units together as much as possible
	separators := []string{
		"\n\n", // Paragraph breaks (highest priority)
		"\n",   // Line breaks
		". ",   // Sentences
		"? ",   // Questions
		"! ",   // Exclamations
		"; ",   // Semicolons
		", ",   // Commas
		" ",    // Words (last resort)
	}

	return r.recursiveSplit(text, separators, 0)
}

// recursiveSplit implements recursive text splitting with multiple separators
func (r *RAGProcessor) recursiveSplit(text string, separators []string, depth int) []string {
	// Base case: text fits within chunk size
	if len(text) <= r.chunkSize {
		return []string{text}
	}

	// If we've exhausted all separators, force split by character count
	if depth >= len(separators) {
		return r.forceSplitBySize(text)
	}

	separator := separators[depth]

	// Check if separator exists in text
	if !strings.Contains(text, separator) {
		// Try next separator
		return r.recursiveSplit(text, separators, depth+1)
	}

	// Split by current separator
	parts := strings.Split(text, separator)

	var finalChunks []string
	var goodParts []string

	// Process each part
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}

		if len(part) > r.chunkSize {
			// Part is too large, need to process accumulated good parts first
			if len(goodParts) > 0 {
				merged := r.mergeParts(goodParts, separator)
				finalChunks = append(finalChunks, merged...)
				goodParts = nil
			}

			// Recursively split the large part with next separator
			subChunks := r.recursiveSplit(part, separators, depth+1)
			finalChunks = append(finalChunks, subChunks...)
		} else {
			// Part is small enough, accumulate it
			goodParts = append(goodParts, part)
		}
	}

	// Process remaining good parts
	if len(goodParts) > 0 {
		merged := r.mergeParts(goodParts, separator)
		finalChunks = append(finalChunks, merged...)
	}

	return finalChunks
}

// mergeParts merges small parts into chunks with overlap support
func (r *RAGProcessor) mergeParts(parts []string, separator string) []string {
	var chunks []string
	var currentChunk strings.Builder

	for _, part := range parts {
		partLen := len(part)
		sepLen := 0
		if currentChunk.Len() > 0 {
			sepLen = len(separator)
		}

		// Check if adding this part would exceed chunk size
		if currentChunk.Len() > 0 && currentChunk.Len()+sepLen+partLen > r.chunkSize {
			// Save current chunk
			chunks = append(chunks, currentChunk.String())

			// Start new chunk with overlap
			currentChunk.Reset()
			if r.overlapSize > 0 && len(chunks) > 0 {
				prevChunk := chunks[len(chunks)-1]
				overlapStart := len(prevChunk) - r.overlapSize
				if overlapStart < 0 {
					overlapStart = 0
				}
				currentChunk.WriteString(prevChunk[overlapStart:])
				currentChunk.WriteString(separator)
			}
		}

		// Add part to current chunk
		if currentChunk.Len() > 0 {
			currentChunk.WriteString(separator)
		}
		currentChunk.WriteString(part)
	}

	// Add final chunk
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks
}

// forceSplitBySize splits text by character count when no separator works
// This is a last resort to ensure we never exceed chunk size
func (r *RAGProcessor) forceSplitBySize(text string) []string {
	var chunks []string

	for len(text) > 0 {
		if len(text) <= r.chunkSize {
			chunks = append(chunks, text)
			break
		}

		// Take chunk size worth of text
		chunk := text[:r.chunkSize]
		chunks = append(chunks, chunk)

		// Move forward with overlap
		step := r.chunkSize - r.overlapSize
		if step <= 0 {
			step = r.chunkSize
		}
		text = text[step:]
	}

	return chunks
}
