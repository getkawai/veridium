package chromem_test

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/components/retriever"
	"github.com/cloudwego/eino/schema"
	"github.com/kawai-network/veridium/pkg/chromem"
	chromemAdapter "github.com/kawai-network/veridium/pkg/eino-adapters/chromem"
)

// Example demonstrates basic usage of chromem with Eino adapters
func Example_basic() {
	ctx := context.Background()

	// 1. Create chromem database and collection
	db := chromem.NewDB()
	collection, err := db.CreateCollection(
		"my_documents",
		nil,                               // metadata
		chromem.NewEmbeddingFuncDefault(), // uses OpenAI by default
	)
	if err != nil {
		log.Fatal(err)
	}

	// 2. Create Eino indexer
	indexer := chromemAdapter.NewIndexer(collection)

	// 3. Add documents using Eino interface
	docs := []*schema.Document{
		{
			ID:      "doc1",
			Content: "Chromem is a local vector database for Go",
			MetaData: map[string]any{
				"source": "documentation",
				"type":   "intro",
			},
		},
		{
			ID:      "doc2",
			Content: "Eino is a graph orchestration framework",
			MetaData: map[string]any{
				"source": "documentation",
				"type":   "intro",
			},
		},
	}

	ids, err := indexer.Store(ctx, docs)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Stored %d documents\n", len(ids))

	// 4. Create Eino retriever
	ret, err := chromemAdapter.NewRetriever(&chromemAdapter.RetrieverConfig{
		Collection: collection,
		TopK:       5,
	})
	if err != nil {
		log.Fatal(err)
	}

	// 5. Query using Eino interface
	results, err := ret.Retrieve(ctx, "vector database")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d results\n", len(results))
	for i, doc := range results {
		fmt.Printf("Result %d: %s (similarity: %.2f)\n",
			i+1, doc.Content, doc.MetaData["similarity"])
	}

	// Output:
	// Stored 2 documents
	// Found 2 results
}

// Example demonstrates using chromem adapter in an Eino graph
func Example_einoGraph() {
	ctx := context.Background()

	// Setup chromem
	db := chromem.NewDB()
	collection, _ := db.CreateCollection("docs", nil, chromem.NewEmbeddingFuncDefault())

	// Add some documents
	indexer := chromemAdapter.NewIndexer(collection)
	docs := []*schema.Document{
		{Content: "Go is a programming language"},
		{Content: "Python is also a programming language"},
		{Content: "Rust is a systems programming language"},
	}
	indexer.Store(ctx, docs)

	// Create retriever
	ret, _ := chromemAdapter.NewRetriever(&chromemAdapter.RetrieverConfig{
		Collection: collection,
		TopK:       2,
	})

	// Use retriever directly (graph integration example simplified)
	results, err := ret.Retrieve(ctx, "programming language")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Retrieved %d documents\n", len(results))

	// Output:
	// Retrieved 2 documents
}

// Example demonstrates using retriever options
func Example_retrieverOptions() {
	ctx := context.Background()

	// Setup
	db := chromem.NewDB()
	collection, _ := db.CreateCollection("docs", nil, chromem.NewEmbeddingFuncDefault())

	indexer := chromemAdapter.NewIndexer(collection)
	docs := []*schema.Document{
		{
			Content: "Machine learning basics",
			MetaData: map[string]any{
				"category": "ml",
				"level":    "beginner",
			},
		},
		{
			Content: "Advanced deep learning",
			MetaData: map[string]any{
				"category": "ml",
				"level":    "advanced",
			},
		},
		{
			Content: "Web development guide",
			MetaData: map[string]any{
				"category": "web",
				"level":    "beginner",
			},
		},
	}
	indexer.Store(ctx, docs)

	// Create retriever
	ret, _ := chromemAdapter.NewRetriever(&chromemAdapter.RetrieverConfig{
		Collection: collection,
		TopK:       10,
	})

	// Query with options
	results, err := ret.Retrieve(ctx, "learning",
		retriever.WithTopK(2),
		retriever.WithScoreThreshold(0.5),
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Found %d documents about learning\n", len(results))

	// Output:
	// Found 2 documents about learning
}

// Example demonstrates persistent chromem with Eino
func Example_persistent() {
	ctx := context.Background()

	// Create persistent chromem database
	db, err := chromem.NewPersistentDB("./data/chromem-eino", true)
	if err != nil {
		log.Fatal(err)
	}

	// Get or create collection
	collection, err := db.GetOrCreateCollection(
		"persistent_docs",
		nil,
		chromem.NewEmbeddingFuncDefault(),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Use with Eino adapters
	indexer := chromemAdapter.NewIndexer(collection)
	retriever, _ := chromemAdapter.NewRetriever(&chromemAdapter.RetrieverConfig{
		Collection: collection,
		TopK:       5,
	})

	// Add documents (will be persisted)
	docs := []*schema.Document{
		{Content: "This will be persisted to disk"},
	}
	indexer.Store(ctx, docs)

	// Query (works even after restart)
	results, _ := retriever.Retrieve(ctx, "persisted")
	fmt.Printf("Found %d persisted documents\n", len(results))

	// Output:
	// Found 1 persisted documents
}
