package chromem

import (
	"context"
	"fmt"
	"runtime"

	"github.com/cloudwego/eino/components/indexer"
	"github.com/cloudwego/eino/schema"
	"github.com/kawai-network/veridium/pkg/chromem"
)

// Indexer is an Eino-compatible indexer that wraps chromem.Collection.
// It implements the indexer.Indexer interface from Eino.
type Indexer struct {
	collection *chromem.Collection
}

// NewIndexer creates a new Eino-compatible indexer wrapping a chromem collection.
func NewIndexer(collection *chromem.Collection) *Indexer {
	return &Indexer{
		collection: collection,
	}
}

// Store implements the indexer.Indexer interface.
// It stores documents in the chromem collection.
func (i *Indexer) Store(ctx context.Context, docs []*schema.Document, opts ...indexer.Option) ([]string, error) {
	if len(docs) == 0 {
		return []string{}, nil
	}

	// Convert Eino documents to chromem documents
	chromemDocs := make([]chromem.Document, len(docs))
	for idx, doc := range docs {
		// Generate ID if not present
		id := doc.ID
		if id == "" {
			id = fmt.Sprintf("doc-%d", idx)
		}

		// Convert metadata from map[string]any to map[string]string
		metadata := make(map[string]string)
		if doc.MetaData != nil {
			for k, v := range doc.MetaData {
				if str, ok := v.(string); ok {
					metadata[k] = str
				} else {
					metadata[k] = fmt.Sprintf("%v", v)
				}
			}
		}

		chromemDocs[idx] = chromem.Document{
			ID:       id,
			Content:  doc.Content,
			Metadata: metadata,
			// Embedding will be created by chromem if not provided
		}
	}

	// Add documents to chromem collection
	err := i.collection.AddDocuments(ctx, chromemDocs, runtime.NumCPU())
	if err != nil {
		return nil, fmt.Errorf("failed to add documents to chromem: %w", err)
	}

	// Return document IDs
	ids := make([]string, len(chromemDocs))
	for idx, doc := range chromemDocs {
		ids[idx] = doc.ID
	}

	return ids, nil
}

// GetCollection returns the underlying chromem collection.
// This is useful for accessing chromem-specific functionality.
func (i *Indexer) GetCollection() *chromem.Collection {
	return i.collection
}

// Compile-time check that Indexer implements indexer.Indexer
var _ indexer.Indexer = (*Indexer)(nil)
