package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kawai-network/veridium/langchaingo/embeddings"
	"github.com/kawai-network/veridium/langchaingo/embeddings/cybertron"
	"github.com/kawai-network/veridium/langchaingo/schema"
	"github.com/kawai-network/veridium/langchaingo/vectorstores"
	"github.com/kawai-network/veridium/pkg/chromem"
)

func cosineSimilarity(x, y []float32) float32 {
	if len(x) != len(y) {
		log.Fatal("x and y have different lengths")
	}

	var dot, nx, ny float32

	for i := range x {
		nx += x[i] * x[i]
		ny += y[i] * y[i]
		dot += x[i] * y[i]
	}

	return dot / (float32(math.Sqrt(float64(nx))) * float32(math.Sqrt(float64(ny))))
}

// chromemStore implements a simple vector store using chromem-go with custom embedder
type chromemStore struct {
	collection *chromem.Collection
	embedder   embeddings.Embedder
	docIndex   int
}

func newChromemStore(ctx context.Context, dbPath, collectionName string, embedder embeddings.Embedder) (*chromemStore, error) {
	// Create embedding function using the Cybertron embedder
	embeddingFunc := func(ctx context.Context, text string) ([]float32, error) {
		embedding, err := embedder.EmbedQuery(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("failed to create embedding: %w", err)
		}
		return embedding, nil
	}

	// Create persistent chromem database
	db, err := chromem.NewPersistentDB(dbPath, false)
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
	chromemDocs := make([]chromem.Document, len(docs))
	ids := make([]string, len(docs))

	for i, doc := range docs {
		id := fmt.Sprintf("%d", s.docIndex+i)
		chromemDocs[i] = chromem.Document{
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

func exampleInMemory(ctx context.Context, emb embeddings.Embedder) {
	// We're going to create embeddings for the following strings, then calculate the similarity
	// between them using cosine-simularity.
	docs := []string{
		"tokyo",
		"japan",
		"potato",
	}

	vecs, err := emb.EmbedDocuments(ctx, docs)
	if err != nil {
		log.Fatal("embed query", err)
	}

	fmt.Println("Similarities:")

	for i := range docs {
		for j := range docs {
			fmt.Printf("%6s ~ %6s = %0.2f\n", docs[i], docs[j], cosineSimilarity(vecs[i], vecs[j]))
		}
	}
}

func exampleChromem(ctx context.Context, emb embeddings.Embedder) {
	// Create a temporary directory for the chromem database
	tempDir := filepath.Join(os.TempDir(), "chromem-cybertron-example")
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir) // Clean up after example

	// Create a new Chromem vector store with the Cybertron Embedder to generate embeddings.
	store, err := newChromemStore(ctx, tempDir, "cybertron-example", emb)
	if err != nil {
		log.Fatal("create chromem store", err)
	}

	// Add some documents to the vector store. This will use the Cybertron Embedder to create
	// embeddings for the documents.
	_, err = store.AddDocuments(ctx, []schema.Document{
		{PageContent: "tokyo"},
		{PageContent: "japan"},
		{PageContent: "potato"},
	})
	if err != nil {
		log.Fatal("add documents", err)
	}

	// Perform a similarity search, returning at most three results with similarity scores of
	// at least 0.8. This again uses the Cybertron Embedder to create an embedding for the
	// search query.
	matches, err := store.SimilaritySearch(ctx, "japan", 3,
		vectorstores.WithScoreThreshold(0.8),
	)
	if err != nil {
		log.Fatal("similarity search", err)
	}

	fmt.Println("Matches:")
	for _, match := range matches {
		fmt.Printf(" japan ~ %6s = %0.2f\n", match.PageContent, match.Score)
	}
}

func main() {
	ctx := context.Background()

	// Create an embedder client that uses the "BAAI/bge-small-en-v1.5" model and caches it in
	// the "models" directory. Cybertron will automatically download the model from HuggingFace
	// and convert it when needed.
	//
	// Note that not all models are supported and that Cybertron executes the model locally on
	// the CPU, so larger models will be quite slow!
	emc, err := cybertron.NewCybertron(
		cybertron.WithModelsDir("models"),
		cybertron.WithModel("BAAI/bge-small-en-v1.5"),
	)
	if err != nil {
		log.Fatal("create embedder client", err)
	}

	// Create an embedder from the previously created client.
	emb, err := embeddings.NewEmbedder(emc,
		embeddings.WithStripNewLines(false),
	)
	if err != nil {
		log.Fatal("create embedder", err)
	}

	// Example: use the Embedder to do an in-memory comparison between some documents.
	exampleInMemory(ctx, emb)

	// Example: use the Embedder together with a Chromem Vector Store.
	exampleChromem(ctx, emb)
}
