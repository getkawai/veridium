package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kawai-network/veridium/langchaingo/embeddings"
	"github.com/kawai-network/veridium/langchaingo/llms/ollama"
	"github.com/kawai-network/veridium/langchaingo/schema"
	"github.com/kawai-network/veridium/langchaingo/vectorstores"
	chromemgo "github.com/kawai-network/veridium/pkg/chromem"
)

// chromemStore implements a simple vector store using chromem-go with Ollama embedder
type chromemStore struct {
	collection *chromemgo.Collection
	embedder   *embeddings.EmbedderImpl
	docIndex   int
}

func newChromemStore(ctx context.Context, dbPath, collectionName string, embedder *embeddings.EmbedderImpl) (*chromemStore, error) {
	// Create embedding function using the Ollama embedder
	embeddingFunc := func(ctx context.Context, text string) ([]float32, error) {
		embedding, err := embedder.EmbedQuery(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("failed to create embedding: %w", err)
		}
		return embedding, nil
	}

	// Create persistent chromem database
	db, err := chromemgo.NewPersistentDB(dbPath, false)
	if err != nil {
		return nil, fmt.Errorf("failed to create chromem DB: %w", err)
	}

	// Get or create collection with custom embedding function
	collection, err := db.GetOrCreateCollection(collectionName, nil, embeddingFunc)
	if err != nil {
		return nil, fmt.Errorf("failed to create collection: %w", err)
	}

	return &chromemStore{
		collection: collection,
		embedder:   embedder,
		docIndex:   collection.Count() + 1,
	}, nil
}

func (s *chromemStore) AddDocuments(ctx context.Context, docs []schema.Document) ([]string, error) {
	chromemDocs := make([]chromemgo.Document, len(docs))
	ids := make([]string, len(docs))

	for i, doc := range docs {
		id := fmt.Sprintf("%d", s.docIndex+i)
		chromemDocs[i] = chromemgo.Document{
			ID:      id,
			Content: doc.PageContent,
		}
		ids[i] = id
	}

	err := s.collection.AddDocuments(ctx, chromemDocs, runtime.NumCPU())
	if err != nil {
		return nil, fmt.Errorf("failed to add documents: %w", err)
	}

	s.docIndex += len(docs)
	return ids, nil
}

func (s *chromemStore) SimilaritySearch(ctx context.Context, query string, numDocuments int, options ...vectorstores.Option) ([]schema.Document, error) {
	// Parse options
	opts := &vectorstores.Options{}
	for _, opt := range options {
		opt(opts)
	}

	// Perform search
	results, err := s.collection.Query(ctx, query, numDocuments, nil, nil)
	if err != nil {
		return nil, fmt.Errorf("similarity search failed: %w", err)
	}

	// Convert results to schema.Document
	docs := make([]schema.Document, 0, len(results))
	for _, result := range results {
		// Apply score threshold if set
		if opts.ScoreThreshold > 0 && result.Similarity < opts.ScoreThreshold {
			continue
		}

		docs = append(docs, schema.Document{
			PageContent: result.Content,
			Score:       result.Similarity,
		})
	}

	return docs, nil
}

func main() {
	ollamaLLM, err := ollama.New(ollama.WithModel("llama2"))
	if err != nil {
		log.Fatal(err)
	}
	ollamaEmbedder, err := embeddings.NewEmbedder(ollamaLLM)
	if err != nil {
		log.Fatal(err)
	}

	// Create a temporary directory for the chromem database
	tempDir := filepath.Join(os.TempDir(), "chromem-ollama-example")
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir) // Clean up after example

	// Create a new Chromem vector store with Ollama embeddings
	ctx := context.Background()
	store, err := newChromemStore(ctx, tempDir, "ollama-example", ollamaEmbedder)
	if err != nil {
		log.Fatalf("new: %v\n", err)
	}

	type meta = map[string]any

	// Add documents to the vector store.
	_, errAd := store.AddDocuments(ctx, []schema.Document{
		{PageContent: "Tokyo", Metadata: meta{"population": 9.7, "area": 622}},
		{PageContent: "Kyoto", Metadata: meta{"population": 1.46, "area": 828}},
		{PageContent: "Hiroshima", Metadata: meta{"population": 1.2, "area": 905}},
		{PageContent: "Kazuno", Metadata: meta{"population": 0.04, "area": 707}},
		{PageContent: "Nagoya", Metadata: meta{"population": 2.3, "area": 326}},
		{PageContent: "Toyota", Metadata: meta{"population": 0.42, "area": 918}},
		{PageContent: "Fukuoka", Metadata: meta{"population": 1.59, "area": 341}},
		{PageContent: "Paris", Metadata: meta{"population": 11, "area": 105}},
		{PageContent: "London", Metadata: meta{"population": 9.5, "area": 1572}},
		{PageContent: "Santiago", Metadata: meta{"population": 6.9, "area": 641}},
		{PageContent: "Buenos Aires", Metadata: meta{"population": 15.5, "area": 203}},
		{PageContent: "Rio de Janeiro", Metadata: meta{"population": 13.7, "area": 1200}},
		{PageContent: "Sao Paulo", Metadata: meta{"population": 22.6, "area": 1523}},
	})
	if errAd != nil {
		log.Fatalf("AddDocument: %v\n", errAd)
	}

	type exampleCase struct {
		name         string
		query        string
		numDocuments int
		options      []vectorstores.Option
	}

	type filter = map[string]any

	exampleCases := []exampleCase{
		{
			name:         "Up to 5 Cities in Japan",
			query:        "Which of these are cities are located in Japan?",
			numDocuments: 5,
			options: []vectorstores.Option{
				vectorstores.WithScoreThreshold(0.8),
			},
		},
		{
			name:         "A City in South America",
			query:        "Which of these are cities are located in South America?",
			numDocuments: 1,
			options: []vectorstores.Option{
				vectorstores.WithScoreThreshold(0.8),
			},
		},
		{
			name:         "Large Cities in South America",
			query:        "Which of these are cities are located in South America?",
			numDocuments: 100,
			options: []vectorstores.Option{
				// Note: Chromem currently doesn't support filters in the vectorstore interface
				vectorstores.WithScoreThreshold(0.7),
			},
		},
	}

	// run the example cases
	results := make([][]schema.Document, len(exampleCases))
	for ecI, ec := range exampleCases {
		docs, errSs := store.SimilaritySearch(ctx, ec.query, ec.numDocuments, ec.options...)
		if errSs != nil {
			log.Fatalf("query1: %v\n", errSs)
		}
		results[ecI] = docs
	}

	// print out the results of the run
	fmt.Printf("Results:\n")
	for ecI, ec := range exampleCases {
		texts := make([]string, len(results[ecI]))
		for docI, doc := range results[ecI] {
			texts[docI] = doc.PageContent
		}
		fmt.Printf("%d. case: %s\n", ecI+1, ec.name)
		fmt.Printf("    result: %s\n", strings.Join(texts, ", "))
	}
}
