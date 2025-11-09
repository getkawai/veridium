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
	// Create an OpenAI client. Requires environment variable OPENAI_API_KEY to be set.
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	llmClient := openai.NewClient(apiKey)

	// Create a temporary directory for the chromem database
	tempDir := filepath.Join(os.TempDir(), "chromem-alloydb-example")
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

	vs := chromem.New(persistentKB)

	_, err := vs.AddDocuments(ctx, []schema.Document{
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
			PageContent: "Sao Paulo",
			Metadata: map[string]any{
				"population": 22.6,
				"area":       1523,
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	docs, err := vs.SimilaritySearch(ctx, "Japan", 5)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Docs:", docs)

	// Note: Chromem currently doesn't support metadata filters in the vectorstore interface
	// Searching without filters
	filteredDocs, err := vs.SimilaritySearch(ctx, "Japan", 5, vectorstores.WithScoreThreshold(0.7))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("FilteredDocs:", filteredDocs)
}
