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
	doc, err := r.queries.GetDocument(ctx, db.GetDocumentParams{ID: req.DocumentID})
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

// chunkText splits text into overlapping chunks
func (r *RAGProcessor) chunkText(text string) []string {
	if text == "" {
		return []string{}
	}

	// Split by sentences (simple approach using periods)
	sentences := strings.Split(text, ".")

	var chunks []string
	var currentChunk strings.Builder
	currentSize := 0

	for _, sentence := range sentences {
		sentence = strings.TrimSpace(sentence)
		if sentence == "" {
			continue
		}

		sentenceLen := len(sentence)

		// If adding this sentence exceeds chunk size, save current chunk
		if currentSize > 0 && currentSize+sentenceLen > r.chunkSize {
			chunks = append(chunks, currentChunk.String())

			// Start new chunk with overlap
			currentChunk.Reset()
			// Add last part of previous chunk for overlap
			if r.overlapSize > 0 && len(chunks) > 0 {
				prevChunk := chunks[len(chunks)-1]
				if len(prevChunk) > r.overlapSize {
					currentChunk.WriteString(prevChunk[len(prevChunk)-r.overlapSize:])
					currentChunk.WriteString(" ")
				}
			}
		}

		// Add sentence to current chunk
		if currentChunk.Len() > 0 {
			currentChunk.WriteString(". ")
		}
		currentChunk.WriteString(sentence)
		currentSize = currentChunk.Len()
	}

	// Add final chunk
	if currentChunk.Len() > 0 {
		chunks = append(chunks, currentChunk.String())
	}

	return chunks
}
