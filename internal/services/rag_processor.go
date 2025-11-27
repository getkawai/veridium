package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/google/uuid"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/pkg/chromem"
	"github.com/kawai-network/veridium/pkg/xlog"
)

// RAGProcessor handles RAG processing (chunking + embedding)
type RAGProcessor struct {
	queries     *db.Queries
	chromemDB   *chromem.DB
	embedFunc   chromem.EmbeddingFunc
	chunkSize   int
	overlapSize int
}

// NewRAGProcessor creates a new RAG processor
func NewRAGProcessor(database *sql.DB, chromemDB *chromem.DB, embedFunc chromem.EmbeddingFunc) *RAGProcessor {
	return &RAGProcessor{
		queries:     db.New(database),
		chromemDB:   chromemDB,
		embedFunc:   embedFunc,
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

	// 2. Get or create user-specific chromem collection
	collectionName := fmt.Sprintf("user-%s-kb", req.UserID)
	collection, err := r.chromemDB.GetOrCreateCollection(collectionName, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create collection: %w", err)
	}

	// 3. Chunk the already-parsed content
	chunks := r.chunkText(doc.Content.String)
	if len(chunks) == 0 {
		xlog.Warn("No chunks created from document", "document_id", req.DocumentID)
		return []string{}, nil
	}

	xlog.Info("Created chunks", "count", len(chunks), "document_id", req.DocumentID)

	// 4. Generate embeddings and store in ChromemDB
	chunkIDs := make([]string, 0, len(chunks))
	for i, chunk := range chunks {
		chunkID := uuid.New().String()

		// Generate embedding using chromem embedding function
		embedding, err := r.embedFunc(ctx, chunk)
		if err != nil {
			xlog.Error("Failed to generate embedding", "error", err, "chunk_index", i)
			continue
		}

		if len(embedding) == 0 {
			xlog.Error("Empty embedding response", "chunk_index", i)
			continue
		}

		// Store in ChromemDB
		err = collection.AddDocument(ctx, chromem.Document{
			ID:        chunkID,
			Content:   chunk,
			Embedding: embedding,
			Metadata: map[string]string{
				"file_id":     req.FileID,
				"document_id": req.DocumentID,
				"user_id":     req.UserID,
				"filename":    req.Filename,
				"chunk_index": fmt.Sprintf("%d", i),
			},
		})
		if err != nil {
			xlog.Error("Failed to store chunk in ChromemDB", "error", err, "chunk_id", chunkID)
			continue
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
