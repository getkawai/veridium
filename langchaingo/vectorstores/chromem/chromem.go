package chromem

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync/atomic"

	"github.com/kawai-network/veridium/langchaingo/schema"
	"github.com/kawai-network/veridium/langchaingo/vectorstores"
	"github.com/kawai-network/veridium/pkg/xlog"
)

// Store is a vector store using ChromemDB with persistence.
// It implements the vectorstores.VectorStore interface.
type Store struct {
	kb       *PersistentKB
	tempDir  string
	docIndex atomic.Uint64
}

// New creates a new ChromemDB vector store with persistence.
func New(kb *PersistentKB) *Store {
	tempDir := filepath.Join(os.TempDir(), "chromem-vectorstore")
	os.MkdirAll(tempDir, 0755)

	return &Store{
		kb:      kb,
		tempDir: tempDir,
	}
}

// AddDocuments adds documents to the vector store.
// This implements the vectorstores.VectorStore interface.
func (s *Store) AddDocuments(ctx context.Context, docs []schema.Document, options ...vectorstores.Option) ([]string, error) {
	ids := make([]string, len(docs))

	for i, doc := range docs {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// Convert metadata to string map
		metadata := make(map[string]string)
		for k, v := range doc.Metadata {
			if str, ok := v.(string); ok {
				metadata[k] = str
			} else {
				metadata[k] = fmt.Sprintf("%v", v)
			}
		}

		// Create temporary file for content
		docID := s.docIndex.Add(1)
		tmpFile := filepath.Join(s.tempDir, fmt.Sprintf("doc-%d.txt", docID))

		if err := os.WriteFile(tmpFile, []byte(doc.PageContent), 0644); err != nil {
			return nil, fmt.Errorf("failed to write temp file: %w", err)
		}

		// Store using existing RAG system
		if err := s.kb.Store(tmpFile, metadata); err != nil {
			os.Remove(tmpFile)
			return nil, fmt.Errorf("failed to store document: %w", err)
		}

		// Clean up temp file after storing
		os.Remove(tmpFile)

		ids[i] = fmt.Sprintf("%d", docID)
		xlog.Debug("Added document to vector store", "id", ids[i])
	}

	return ids, nil
}

// SimilaritySearch performs a similarity search in the vector store.
// This implements the vectorstores.VectorStore interface.
func (s *Store) SimilaritySearch(ctx context.Context, query string, numDocuments int, options ...vectorstores.Option) ([]schema.Document, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// Parse options
	opts := &vectorstores.Options{}
	for _, opt := range options {
		opt(opts)
	}

	// Perform search using existing RAG system
	results, err := s.kb.Search(query, numDocuments)
	if err != nil {
		return nil, fmt.Errorf("similarity search failed: %w", err)
	}

	// Convert to LangChain documents
	docs := make([]schema.Document, 0, len(results))
	for _, result := range results {
		// Apply score threshold if set
		if opts.ScoreThreshold > 0 && result.Similarity < opts.ScoreThreshold {
			continue
		}

		// Convert metadata
		metadata := make(map[string]any)
		for k, v := range result.Metadata {
			metadata[k] = v
		}

		docs = append(docs, schema.Document{
			PageContent: result.Content,
			Metadata:    metadata,
			Score:       result.Similarity,
		})
	}

	xlog.Debug("Similarity search completed", "query", query, "results", len(docs))
	return docs, nil
}

// GetPersistentKB returns the underlying PersistentKB instance.
// Useful for accessing custom functionality.
func (s *Store) GetPersistentKB() *PersistentKB {
	return s.kb
}

// Reset clears all documents from the vector store.
func (s *Store) Reset() error {
	return s.kb.Reset()
}

// Count returns the number of documents in the vector store.
func (s *Store) Count() int {
	return s.kb.Count()
}

// Cleanup removes temporary files.
func (s *Store) Cleanup() error {
	return os.RemoveAll(s.tempDir)
}

// Compile-time check that Store implements vectorstores.VectorStore
var _ vectorstores.VectorStore = (*Store)(nil)
