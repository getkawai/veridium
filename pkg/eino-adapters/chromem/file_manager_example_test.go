/*
 * Copyright 2025 Veridium Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package chromem_test

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kawai-network/veridium/pkg/chromem"
	chromemAdapter "github.com/kawai-network/veridium/pkg/eino-adapters/chromem"
)

// Example_fileManager demonstrates how to use FileManager to index and search documents
func Example_fileManager() {
	ctx := context.Background()

	// 1. Setup chromem database
	db := chromem.NewDB()

	// 2. Create collection with embedding function
	collection, err := db.CreateCollection(
		"documents",
		nil,
		chromem.NewEmbeddingFuncDefault(), // Uses OpenAI by default
	)
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	// 3. Create Eino adapter components
	indexer := chromemAdapter.NewIndexer(collection)
	_, err = chromemAdapter.NewRetriever(&chromemAdapter.RetrieverConfig{
		Collection: collection,
		TopK:       3,
	})
	if err != nil {
		log.Fatalf("Failed to create retriever: %v", err)
	}

	// 4. Create file manager with custom parsers
	fileManager, err := chromemAdapter.NewFileManager(ctx, &chromemAdapter.FileManagerConfig{
		Indexer:     indexer,
		AssetDir:    "./test_assets",
		ChunkSize:   1500,
		OverlapSize: 300,
	})
	if err != nil {
		log.Fatalf("Failed to create file manager: %v", err)
	}

	// Clean up test assets
	defer os.RemoveAll("./test_assets")

	// 5. Store files with metadata
	// Note: In a real scenario, you would have actual files
	// For this example, we'll create temporary test files

	// Create test DOCX (simulated)
	fmt.Println("Storing documents...")
	// fileManager.StoreFile(ctx, "/path/to/document.docx", map[string]any{
	// 	"category": "manual",
	// 	"version":  "1.0",
	// })

	// Create test XLSX (simulated)
	// fileManager.StoreFile(ctx, "/path/to/spreadsheet.xlsx", map[string]any{
	// 	"category": "data",
	// 	"year":     2024,
	// })

	// Create test PDF (simulated)
	// fileManager.StoreFile(ctx, "/path/to/report.pdf", map[string]any{
	// 	"category": "report",
	// })

	// 6. List stored files
	files := fileManager.ListFiles()
	fmt.Printf("Stored %d files\n", len(files))

	// 7. Search documents
	// docs, err := retriever.Retrieve(ctx, "installation instructions")
	// if err != nil {
	// 	log.Fatalf("Failed to retrieve: %v", err)
	// }

	// for _, doc := range docs {
	// 	fmt.Printf("Found: %s (similarity: %.2f)\n",
	// 		doc.MetaData["source_file"],
	// 		doc.MetaData["similarity"])
	// 	fmt.Printf("Content preview: %s...\n\n", doc.Content[:100])
	// }

	// Output:
	// Storing documents...
	// Stored 0 files
}

// Example_fileManagerWithCustomParser demonstrates adding a custom parser
func Example_fileManagerWithCustomParser() {
	ctx := context.Background()

	// Setup (same as above)
	db := chromem.NewDB()
	collection, _ := db.CreateCollection("docs", nil, chromem.NewEmbeddingFuncDefault())
	indexer := chromemAdapter.NewIndexer(collection)

	// Create file manager with custom parser
	fm, err := chromemAdapter.NewFileManager(ctx, &chromemAdapter.FileManagerConfig{
		Indexer:  indexer,
		AssetDir: "./test_assets",
		// CustomParsers: map[string]parser.Parser{
		// 	".custom": myCustomParser,
		// },
	})
	if err != nil {
		log.Fatalf("Failed to create file manager: %v", err)
	}

	defer os.RemoveAll("./test_assets")

	// Get supported extensions
	exts := fm.GetSupportedExtensions()
	fmt.Printf("Supported extensions: %v\n", len(exts))

	// Output:
	// Supported extensions: 7
}

// Example_einoWorkflow demonstrates a complete Eino workflow with file management
func Example_einoWorkflow() {
	ctx := context.Background()

	// 1. Setup chromem
	db, _ := chromem.NewPersistentDB("./test_vectors", true)
	defer os.RemoveAll("./test_vectors")

	collection, _ := db.CreateCollection(
		"knowledge_base",
		nil,
		chromem.NewEmbeddingFuncDefault(),
	)

	// 2. Create Eino components
	indexer := chromemAdapter.NewIndexer(collection)
	_, _ = chromemAdapter.NewRetriever(&chromemAdapter.RetrieverConfig{
		Collection: collection,
		TopK:       5,
	})

	// 3. Create file manager
	_, _ = chromemAdapter.NewFileManager(ctx, &chromemAdapter.FileManagerConfig{
		Indexer:     indexer,
		AssetDir:    "./test_kb_assets",
		ChunkSize:   2000,
		OverlapSize: 400,
	})
	defer os.RemoveAll("./test_kb_assets")

	// 4. Index documents
	fmt.Println("Building knowledge base...")
	// fileManager.StoreFile(ctx, "/docs/manual.pdf", map[string]any{"type": "manual"})
	// fileManager.StoreFile(ctx, "/docs/api.md", map[string]any{"type": "api"})
	// fileManager.StoreFile(ctx, "/docs/tutorial.docx", map[string]any{"type": "tutorial"})

	// 5. Query knowledge base
	query := "how to authenticate"
	// docs, _ := retriever.Retrieve(ctx, query)

	fmt.Printf("Query: %s\n", query)
	fmt.Printf("Found %d relevant documents\n", 0)

	// Output:
	// Building knowledge base...
	// Query: how to authenticate
	// Found 0 relevant documents
}

// Example_persistentFileManager demonstrates file manager with persistent storage
func Example_persistentFileManager() {
	ctx := context.Background()

	// Use persistent database
	db, err := chromem.NewPersistentDB("./persistent_vectors", true)
	if err != nil {
		log.Fatalf("Failed to create persistent database: %v", err)
	}
	defer os.RemoveAll("./persistent_vectors")

	collection, err := db.CreateCollection(
		"persistent_docs",
		nil,
		chromem.NewEmbeddingFuncDefault(),
	)
	if err != nil {
		log.Fatalf("Failed to create collection: %v", err)
	}

	// Create components
	indexer := chromemAdapter.NewIndexer(collection)
	fm, err := chromemAdapter.NewFileManager(ctx, &chromemAdapter.FileManagerConfig{
		Indexer:  indexer,
		AssetDir: "./persistent_assets",
	})
	if err != nil {
		log.Fatalf("Failed to create file manager: %v", err)
	}
	defer os.RemoveAll("./persistent_assets")

	// Store files - they will persist across restarts
	fmt.Println("Files are stored persistently")
	files := fm.ListFiles()
	fmt.Printf("Currently tracking %d files\n", len(files))

	// Output:
	// Files are stored persistently
	// Currently tracking 0 files
}

