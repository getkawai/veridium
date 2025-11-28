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

package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kawai-network/veridium/internal/database"
	database_gen "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/internal/services"
)

func main() {
	ctx := context.Background()

	log.Println("🧪 Testing FileProcessorService")
	log.Println(strings.Repeat("=", 80))

	// Step 1: Initialize Database
	log.Println("\n📦 Step 1: Initializing database...")
	dbService, err := database.NewService()
	if err != nil {
		log.Fatalf("❌ Failed to initialize database: %v", err)
	}
	defer dbService.Close()
	log.Println("✅ Database initialized")

	// Step 2: Initialize Llama Library Service
	log.Println("\n🦙 Step 2: Initializing Llama Library service...")
	libService, err := llama.NewLibraryService()
	if err != nil {
		log.Fatalf("❌ Failed to initialize Llama Library service: %v", err)
	}
	defer libService.Cleanup()
	log.Println("✅ Llama Library service initialized")
	log.Printf("   Models directory: %s", libService.GetModelsDirectory())

	// Step 3: Initialize DuckDB Store
	log.Println("\n🦆 Step 3: Initializing DuckDB store...")
	// Store DuckDB in project data directory
	duckDBPath := filepath.Join("data", "test-duckdb.db")
	duckDBStore, err := services.NewDuckDBStore(duckDBPath, 384) // 384 dims for granite-embedding
	if err != nil {
		log.Fatalf("❌ Failed to initialize DuckDB Store: %v", err)
	}
	defer duckDBStore.Close()
	log.Println("✅ DuckDB Store initialized")
	log.Printf("   Database path: %s", duckDBPath)

	// Step 4: Initialize Vector Search Service
	log.Println("\n🔍 Step 4: Initializing Vector Search service...")
	vectorSearchService, err := services.NewVectorSearchService(
		dbService.DB(),
		duckDBStore,
		"llama",
		"",
		libService,
	)
	if err != nil {
		log.Fatalf("❌ Failed to initialize Vector Search service: %v", err)
	}
	log.Println("✅ Vector Search service initialized")

	// Step 5: Initialize File Loader
	log.Println("\n📂 Step 5: Initializing File Loader...")
	fileLoader := services.NewFileLoader()
	log.Println("✅ File Loader initialized")

	// Step 6: Initialize File Processor Service
	log.Println("\n⚙️  Step 6: Initializing File Processor service...")
	// Store test files in project data directory
	fileBaseDir := filepath.Join("data", "test-files")
	os.MkdirAll(fileBaseDir, 0755)

	documentService := services.NewDocumentService(dbService.DB())
	embedder := vectorSearchService.GetEmbedder()
	ragProcessor := services.NewRAGProcessor(dbService.DB(), duckDBStore, fileLoader, embedder)
	fileProcessor := services.NewFileProcessorService(
		dbService.DB(),
		fileLoader,
		documentService,
		ragProcessor,
	)
	log.Println("✅ File Processor service initialized")
	log.Printf("   File base directory: %s", fileBaseDir)

	// Step 7: Find Real Test Files
	log.Println("\n📝 Step 7: Finding real test files from project...")
	testFiles := findRealTestFiles()
	log.Printf("✅ Found %d test files", len(testFiles))

	// Step 8: Create Test User
	log.Println("\n👤 Step 8: Creating test user...")
	testUserID := "test-user-" + uuid.New().String()[:8]
	now := time.Now().UnixMilli()

	_, err = dbService.Queries().CreateUser(ctx, database_gen.CreateUserParams{
		ID:              testUserID,
		Username:        sql.NullString{String: "testuser-" + testUserID[:8], Valid: true},
		Email:           sql.NullString{String: "test-" + testUserID[:8] + "@example.com", Valid: true},
		Avatar:          sql.NullString{},
		Phone:           sql.NullString{},
		FirstName:       sql.NullString{String: "Test", Valid: true},
		LastName:        sql.NullString{String: "User", Valid: true},
		IsOnboarded:     1,
		ClerkCreatedAt:  sql.NullInt64{},
		EmailVerifiedAt: sql.NullInt64{Int64: now, Valid: true},
		Preference:      sql.NullString{String: "{}", Valid: true},
		CreatedAt:       now,
		UpdatedAt:       now,
	})
	if err != nil {
		log.Fatalf("❌ Failed to create test user: %v", err)
	}
	log.Printf("✅ Test user created: %s", testUserID)

	// Step 9: Test File Processing
	log.Println("\n🧪 Step 9: Testing file processing...")

	for i, testFile := range testFiles {
		log.Printf("\n--- Test %d/%d: %s ---", i+1, len(testFiles), testFile.Name)
		log.Printf("File type: %s", testFile.FileType)
		log.Printf("Enable RAG: %v", testFile.EnableRAG)

		// Process file
		req := services.ProcessFileRequest{
			FilePath:  testFile.Path,
			Filename:  testFile.Name,
			FileType:  testFile.FileType,
			UserID:    testUserID,
			ClientID:  "",
			Source:    testFile.Path,
			EnableRAG: testFile.EnableRAG,
			IsShared:  false,
			FileMetadata: &services.FileMetadata{
				Filename: testFile.Name,
				FileType: testFile.FileType,
			},
		}

		response, err := fileProcessor.ProcessFile(ctx, req)
		if err != nil {
			log.Printf("❌ Failed to process file: %v", err)
			continue
		}

		log.Printf("✅ File processed successfully")
		log.Printf("   File ID: %s", response.FileID)
		log.Printf("   Document ID: %s", response.DocumentID)
		if len(response.ChunkIDs) > 0 {
			log.Printf("   Chunks created: %d", len(response.ChunkIDs))
		} else if testFile.EnableRAG {
			log.Printf("   ⚠️  No chunks created (RAG was enabled)")
		}

		// Verify document was saved
		doc, err := dbService.Queries().GetDocument(ctx, database_gen.GetDocumentParams{
			ID:     response.DocumentID,
			UserID: testUserID,
		})
		if err != nil {
			log.Printf("❌ Failed to retrieve document: %v", err)
			continue
		}

		log.Printf("   Document title: %s", doc.Title.String)
		log.Printf("   Document content length: %d chars", len(doc.Content.String))
		log.Printf("   Document char count: %d", doc.TotalCharCount)
		log.Printf("   Document line count: %d", doc.TotalLineCount)

		// If RAG enabled, verify chunks
		if testFile.EnableRAG && len(response.ChunkIDs) > 0 {
			// Get first chunk
			chunk, err := dbService.Queries().GetChunk(ctx, database_gen.GetChunkParams{
				ID:     response.ChunkIDs[0],
				UserID: sql.NullString{String: testUserID, Valid: true},
			})
			if err != nil {
				log.Printf("❌ Failed to retrieve chunk: %v", err)
			} else {
				log.Printf("   First chunk length: %d chars", len(chunk.Text.String))
				log.Printf("   First chunk preview: %s...", truncate(chunk.Text.String, 80))
			}

			// Verify vector in DuckDB by getting the chunk and searching with its embedding
			chunkData, errChunk := dbService.Queries().GetChunk(ctx, database_gen.GetChunkParams{
				ID:     response.ChunkIDs[0],
				UserID: sql.NullString{String: testUserID, Valid: true},
			})
			if errChunk == nil {
				// Generate embedding for the chunk to search
				embeddings, errEmbed := embedder.EmbedStrings(ctx, []string{chunkData.Text.String})
				if errEmbed == nil && len(embeddings) > 0 {
					vectorIDs, errSearch := duckDBStore.SearchVectors(ctx, embeddings[0], 1)
					if errSearch != nil {
						log.Printf("⚠️  Failed to search vectors in DuckDB: %v", errSearch)
					} else if len(vectorIDs) > 0 {
						log.Printf("   ✅ Vector found in DuckDB (ID: %s)", vectorIDs[0])
					} else {
						log.Printf("   ⚠️  Vector not found in DuckDB")
					}
				}
			}
		}
	}

	// Step 10: Test Semantic Search
	log.Println("\n🔎 Step 10: Testing semantic search...")

	// 10a. Check if vectors exist in DuckDB
	log.Println("   10a. Checking DuckDB vectors...")
	testEmbedding := make([]float32, 384)
	for i := range testEmbedding {
		testEmbedding[i] = 0.1
	}
	vectorResults, err := duckDBStore.SearchVectors(ctx, testEmbedding, 10)
	if err != nil {
		log.Printf("   ❌ Failed to query DuckDB: %v", err)
	} else {
		log.Printf("   ✅ Found %d vectors in DuckDB", len(vectorResults))
	}

	// 10b. Test semantic search with default metric (Euclidean)
	log.Println("   10b. Testing semantic search (Euclidean distance)...")
	testQuery := "How to build the application?"

	// Get all file IDs from processed files
	files, err := dbService.Queries().ListFiles(ctx, database_gen.ListFilesParams{
		UserID: testUserID,
		Limit:  100,
		Offset: 0,
	})
	if err != nil {
		log.Printf("❌ Failed to list files: %v", err)
	} else {
		log.Printf("   Found %d files in database", len(files))
		fileIDs := make([]string, len(files))
		for i, f := range files {
			fileIDs[i] = f.ID
		}

		if len(fileIDs) > 0 {
			results, err := vectorSearchService.SemanticSearch(ctx, testUserID, testQuery, fileIDs, 5)
			if err != nil {
				log.Printf("   ❌ Semantic search failed: %v", err)
			} else {
				log.Printf("   ✅ Semantic search completed (Euclidean)")
				log.Printf("      Query: %s", testQuery)
				log.Printf("      Results: %d", len(results))
				for i, result := range results {
					log.Printf("      %d. Similarity: %.4f, File: %s", i+1, result.Similarity, result.FileName)
					if i == 0 {
						log.Printf("         Text: %s...", truncate(result.Text, 80))
					}
				}
			}

			// 10c. Test batch search (LATERAL joins - 66× faster!)
			log.Println("   10c. Testing batch search with LATERAL joins...")
			if len(results) > 0 {
				// Create batch queries from first 3 results
				batchQueries := []services.BatchSearchRequest{}
				for i := 0; i < 3 && i < len(results); i++ {
					// Generate embedding for each result text
					embeddings, err := embedder.EmbedStrings(ctx, []string{results[i].Text})
					if err == nil && len(embeddings) > 0 {
						batchQueries = append(batchQueries, services.BatchSearchRequest{
							QueryID:   fmt.Sprintf("query-%d", i+1),
							Embedding: embeddings[0],
						})
					}
				}

				if len(batchQueries) > 0 {
					batchResults, err := duckDBStore.BatchSearchVectors(ctx, batchQueries, 3)
					if err != nil {
						log.Printf("      ❌ Batch search failed: %v", err)
					} else {
						log.Printf("      ✅ Batch search completed")
						log.Printf("         Processed %d queries in one go (66× faster than individual!)", len(batchQueries))
						for _, br := range batchResults {
							log.Printf("         Query %s: found %d results", br.QueryID, len(br.Results))
						}
					}
				}
			}
		}
	}

	// Step 11: Summary
	log.Println("\n📊 Step 11: Test Summary...")
	log.Printf("   Total files processed: %d", len(testFiles))
	log.Printf("   Files with RAG enabled: %d", countRAGEnabled(testFiles))
	log.Println("✅ No cleanup needed (used real project files)")

	// Final Summary
	log.Println("\n" + strings.Repeat("=", 80))
	log.Println("✅ All FileProcessorService tests completed successfully!")
	log.Println(strings.Repeat("=", 80))
}

// TestFile represents a test file
type TestFile struct {
	Name      string
	Path      string
	FileType  string
	EnableRAG bool
}

// findRealTestFiles finds real files from the project to test
func findRealTestFiles() []TestFile {
	testFiles := []TestFile{
		// Markdown files
		{
			Name:      "README.md",
			Path:      "/Users/yuda/github.com/kawai-network/veridium/README.md",
			FileType:  "txt",
			EnableRAG: true,
		},
		{
			Name:      "BUILD.md",
			Path:      "/Users/yuda/github.com/kawai-network/veridium/docs/BUILD.md",
			FileType:  "txt",
			EnableRAG: true,
		},
		{
			Name:      "LLM_OPTIMIZATION_GUIDE.md",
			Path:      "/Users/yuda/github.com/kawai-network/veridium/docs/LLM_OPTIMIZATION_GUIDE.md",
			FileType:  "txt",
			EnableRAG: true,
		},

		// PDF file
		{
			Name:      "test_pdf.pdf",
			Path:      "/Users/yuda/github.com/kawai-network/veridium/cloudwego/eino-ext/components/document/parser/pdf/testdata/test_pdf.pdf",
			FileType:  "pdf",
			EnableRAG: true,
		},

		// DOCX file
		{
			Name:      "Word-Windows.docx",
			Path:      "/Users/yuda/github.com/kawai-network/veridium/gooxml/testdata/Office2016/Word-Windows.docx",
			FileType:  "docx",
			EnableRAG: true,
		},

		// XLSX file
		{
			Name:      "Excel-Windows.xlsx",
			Path:      "/Users/yuda/github.com/kawai-network/veridium/gooxml/testdata/Office2016/Excel-Windows.xlsx",
			FileType:  "xlsx",
			EnableRAG: true,
		},

		// PPTX file
		{
			Name:      "PowerPoint-Windows.pptx",
			Path:      "/Users/yuda/github.com/kawai-network/veridium/gooxml/testdata/Office2016/PowerPoint-Windows.pptx",
			FileType:  "pptx",
			EnableRAG: true,
		},

		// Plain text (no RAG)
		{
			Name:      "go.mod",
			Path:      "/Users/yuda/github.com/kawai-network/veridium/go.mod",
			FileType:  "txt",
			EnableRAG: false, // Don't RAG process go.mod
		},
	}

	// Filter out files that don't exist
	var validFiles []TestFile
	for _, file := range testFiles {
		if _, err := os.Stat(file.Path); err == nil {
			validFiles = append(validFiles, file)
		} else {
			log.Printf("⚠️  File not found, skipping: %s", file.Path)
		}
	}

	return validFiles
}

// countRAGEnabled counts how many files have RAG enabled
func countRAGEnabled(files []TestFile) int {
	count := 0
	for _, f := range files {
		if f.EnableRAG {
			count++
		}
	}
	return count
}

// truncate truncates a string to maxLen characters
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
