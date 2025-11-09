package llamacpp

import (
	"context"
	"fmt"

	"github.com/kawai-network/veridium/langchaingo/embeddings"
)

// LlamaCPP is an embedder using llama.cpp with GGUF models.
// It implements the embeddings.Embedder interface.
type LlamaCPP struct {
	vectorizer *Vectorizer
}

// New creates a new llama.cpp embedder using a GGUF model.
func New(modelPath string, gpuLayers int) (*LlamaCPP, error) {
	vectorizer, err := NewVectorizer(modelPath, gpuLayers)
	if err != nil {
		return nil, fmt.Errorf("failed to create vectorizer: %w", err)
	}

	return &LlamaCPP{vectorizer: vectorizer}, nil
}

// NewFromVectorizer creates a LlamaCPP embedder from an existing Vectorizer.
func NewFromVectorizer(vectorizer *Vectorizer) *LlamaCPP {
	return &LlamaCPP{vectorizer: vectorizer}
}

// EmbedDocuments creates embeddings for multiple documents.
// This implements the embeddings.Embedder interface.
func (l *LlamaCPP) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	results := make([][]float32, len(texts))

	for i, text := range texts {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		embedding, err := l.vectorizer.EmbedText(text)
		if err != nil {
			return nil, fmt.Errorf("failed to embed document %d: %w", i, err)
		}
		results[i] = embedding
	}

	return results, nil
}

// EmbedQuery creates an embedding for a single query text.
// This implements the embeddings.Embedder interface.
func (l *LlamaCPP) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	// Check context cancellation
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	embedding, err := l.vectorizer.EmbedText(text)
	if err != nil {
		return nil, fmt.Errorf("failed to embed query: %w", err)
	}

	return embedding, nil
}

// Close releases resources held by the vectorizer.
func (l *LlamaCPP) Close() error {
	return l.vectorizer.Close()
}

// GetVectorizer returns the underlying Vectorizer instance.
// Useful for accessing custom functionality.
func (l *LlamaCPP) GetVectorizer() *Vectorizer {
	return l.vectorizer
}

// Compile-time check that LlamaCPP implements embeddings.Embedder
var _ embeddings.Embedder = (*LlamaCPP)(nil)
