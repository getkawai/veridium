package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kawai-network/veridium/langchaingo/schema"
	"github.com/kawai-network/veridium/langchaingo/vectorstores"
	"github.com/kawai-network/veridium/langchaingo/vectorstores/chromem"
	"github.com/sashabaranov/go-openai"
)

func main() {
	// Note: Mistral embeddings require Mistral API, but chromem needs OpenAI-compatible client
	// We'll use OpenAI embeddings for this example
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	llmClient := openai.NewClient(apiKey)

	// Create a temporary directory for the chromem database
	tempDir := filepath.Join(os.TempDir(), "chromem-mistral-example")
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir) // Clean up after example

	// Create a new Chromem vector store
	ctx := context.Background()
	persistentKB := chromem.NewPersistentChromeCollection(
		llmClient,
		"cities-collection",
		tempDir,
		filepath.Join(tempDir, "assets"),
		"text-embedding-ada-002",
		1000, // max chunk size
	)

	store := chromem.New(persistentKB)

	// Add documents to the chromem store.
	_, err := store.AddDocuments(context.Background(), []schema.Document{
		{
			PageContent: "Tokyo",
			Metadata: map[string]any{
				"population": 38,
				"area":       2190,
			},
		},
		{
			PageContent: "Paris",
			Metadata: map[string]any{
				"population": 11,
				"area":       105,
			},
		},
		{
			PageContent: "London",
			Metadata: map[string]any{
				"population": 9.5,
				"area":       1572,
			},
		},
		{
			PageContent: "Santiago",
			Metadata: map[string]any{
				"population": 6.9,
				"area":       641,
			},
		},
		{
			PageContent: "Buenos Aires",
			Metadata: map[string]any{
				"population": 15.5,
				"area":       203,
			},
		},
		{
			PageContent: "Rio de Janeiro",
			Metadata: map[string]any{
				"population": 13.7,
				"area":       1200,
			},
		},
		{
			PageContent: "Sao Paulo",
			Metadata: map[string]any{
				"population": 22.6,
				"area":       1523,
			},
		},
	})
	if err != nil {
		log.Fatal("store.AddDocuments:\n", err)
	}

	// Search for similar documents.
	docs, err := store.SimilaritySearch(ctx, "japan", 1)
	if err != nil {
		log.Fatal("store.SimilaritySearch1:\n", err)
	}
	fmt.Println("store.SimilaritySearch1:\n", docs)

	// Search for similar documents using score threshold.
	docs, err = store.SimilaritySearch(ctx, "only cities in south america", 3, vectorstores.WithScoreThreshold(0.50))
	if err != nil {
		log.Fatal("store.SimilaritySearch2:\n", err)
	}
	fmt.Println("store.SimilaritySearch2:\n", docs)

	// Note: Chromem currently doesn't support metadata filters in the vectorstore interface
	// Searching without filters
	docs, err = store.SimilaritySearch(ctx, "only cities in south america",
		3,
		vectorstores.WithScoreThreshold(0.50),
	)
	if err != nil {
		log.Fatal("store.SimilaritySearch3:\n", err)
	}
	fmt.Println("store.SimilaritySearch3:\n", docs)
}
