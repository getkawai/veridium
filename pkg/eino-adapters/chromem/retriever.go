package chromem

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	"github.com/kawai-network/veridium/pkg/chromem"
)

// Retriever is an Eino-compatible retriever that wraps chromem.Collection.
// It implements the retriever.Retriever interface from Eino.
type Retriever struct {
	collection *chromem.Collection
	topK       int
}

// RetrieverConfig holds configuration for creating a Retriever.
type RetrieverConfig struct {
	// Collection is the chromem collection to query.
	// Required.
	Collection *chromem.Collection

	// TopK is the default number of results to return.
	// Optional, defaults to 5.
	TopK int
}

// NewRetriever creates a new Eino-compatible retriever wrapping a chromem collection.
func NewRetriever(config *RetrieverConfig) (*Retriever, error) {
	if config.Collection == nil {
		return nil, fmt.Errorf("collection is required")
	}

	topK := config.TopK
	if topK <= 0 {
		topK = 5 // default
	}

	return &Retriever{
		collection: config.Collection,
		topK:       topK,
	}, nil
}

// Retrieve implements the retriever.Retriever interface.
// It performs similarity search on the chromem collection.
func (r *Retriever) Retrieve(ctx context.Context, query string, opts ...retriever.Option) ([]*schema.Document, error) {
	// Parse options using Eino's helper
	config := retriever.GetCommonOptions(nil, opts...)

	// Determine number of results
	nResults := r.topK
	if config.TopK != nil && *config.TopK > 0 {
		nResults = *config.TopK
	}

	// Note: Chromem doesn't support Index/SubIndex, so we ignore those options
	// For metadata filtering, we would need to extend chromem's Query API
	// For now, we use chromem's native Query method

	// Perform query using chromem's native API
	results, err := r.collection.Query(ctx, query, nResults, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("chromem query failed: %w", err)
	}

	// Convert chromem results to Eino documents
	docs := make([]*schema.Document, 0, len(results))
	for _, result := range results {
		// Apply score threshold if set
		if config.ScoreThreshold != nil && result.Similarity < float32(*config.ScoreThreshold) {
			continue
		}

		// Convert metadata from map[string]string to map[string]any
		metadata := make(map[string]any)
		for k, v := range result.Metadata {
			metadata[k] = v
		}

		// Add similarity score to metadata
		metadata["similarity"] = result.Similarity

		doc := &schema.Document{
			ID:       result.ID,
			Content:  result.Content,
			MetaData: metadata,
		}

		docs = append(docs, doc)
	}

	return docs, nil
}

// GetCollection returns the underlying chromem collection.
// This is useful for accessing chromem-specific functionality.
func (r *Retriever) GetCollection() *chromem.Collection {
	return r.collection
}

// Compile-time check that Retriever implements retriever.Retriever
var _ retriever.Retriever = (*Retriever)(nil)
