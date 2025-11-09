package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kawai-network/veridium/langchaingo/schema"
	"github.com/kawai-network/veridium/langchaingo/vectorstores"
	"github.com/kawai-network/veridium/langchaingo/vectorstores/chromem"
	"github.com/sashabaranov/go-openai"
)

func main() {
	// Get OpenAI API key from environment
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENAI_API_KEY environment variable is required")
	}

	// Create OpenAI client
	llmClient := openai.NewClient(apiKey)

	// Create a temporary directory for the chromem database
	tempDir := filepath.Join(os.TempDir(), "chromem-example")
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir) // Clean up after example

	// Create a new Chromem vector store
	collectionName := "cities-example"
	embeddingModel := "text-embedding-ada-002"

	persistentKB := chromem.NewPersistentChromeCollection(
		llmClient,
		collectionName,
		tempDir,
		filepath.Join(tempDir, "assets"),
		embeddingModel,
		1000, // max chunk size
	)

	store := chromem.New(persistentKB)

	type meta = map[string]any

	// Add documents to the vector store.
	_, errAd := store.AddDocuments(context.Background(), []schema.Document{
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

	ctx := context.TODO()

	type exampleCase struct {
		name         string
		query        string
		numDocuments int
		options      []vectorstores.Option
	}

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
				// The underlying chromem-go library supports filters, but they're not yet exposed
				// This query will return all matching cities without filtering by area/population
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
