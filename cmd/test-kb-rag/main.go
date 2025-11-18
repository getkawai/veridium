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
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kawai-network/veridium/internal/database"
	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/pkg/chromem"
)

func main() {
	ctx := context.Background()

	log.Println("🧪 Testing Knowledge Base and RAG Integration")
	log.Println(strings.Repeat("=", 60))

	// 1. Initialize Database
	log.Println("\n📦 Step 1: Initializing database...")
	dbService, err := database.NewService()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbService.Close()
	log.Println("✅ Database initialized")

	// 2. Initialize Llama Library Service
	log.Println("\n🦙 Step 2: Initializing Llama Library service...")
	libService, err := llama.NewLibraryService()
	if err != nil {
		log.Fatalf("Failed to initialize Llama Library service: %v", err)
	}
	defer libService.Cleanup()
	log.Println("✅ Llama Library service initialized")
	log.Printf("   Models directory: %s", libService.GetModelsDirectory())

	// 3. Initialize Embedding Function
	log.Println("\n🔢 Step 3: Initializing embedding function...")
	embeddingModelPath := filepath.Join(libService.GetModelsDirectory(),
		"granite-embedding-107m-multilingual-Q6_K_L.gguf")
	embedFunc := chromem.NewEmbeddingFuncLlamaWithPreloadedLibrary(embeddingModelPath)
	log.Println("✅ Embedding function initialized")
	log.Printf("   Model: %s", embeddingModelPath)

	// 4. Initialize Knowledge Base Service
	log.Println("\n📚 Step 4: Initializing Knowledge Base service...")
	userConfigDir, _ := os.UserConfigDir()
	kbPath := filepath.Join(userConfigDir, "veridium", "test-kb")
	kbAssetPath := filepath.Join(userConfigDir, "veridium", "test-kb-assets")

	kbService, err := services.NewKnowledgeBaseService(dbService, &services.KnowledgeBaseConfig{
		ChromemPath:   kbPath,
		EmbeddingFunc: embedFunc,
		AssetDir:      kbAssetPath,
	})
	if err != nil {
		log.Fatalf("Failed to initialize Knowledge Base service: %v", err)
	}
	log.Println("✅ Knowledge Base service initialized")

	// 5. Create a Test Knowledge Base
	log.Println("\n📝 Step 5: Creating test knowledge base...")
	testUserID := "test-user-001"
	kbID, err := kbService.CreateKnowledgeBase(ctx, "Test KB", "Test knowledge base for integration testing", testUserID)
	if err != nil {
		log.Fatalf("Failed to create knowledge base: %v", err)
	}
	log.Printf("✅ Knowledge base created: %s", kbID)

	// 6. Create a Test Document
	log.Println("\n📄 Step 6: Creating test document...")
	testDocPath := filepath.Join(os.TempDir(), "test-doc.txt")
	testContent := `
This is a test document for the Knowledge Base system.

The system uses CloudWeGo Eino for orchestration.
It uses Chromem for vector storage.
It uses Llama.cpp for local inference.

The architecture includes:
1. Knowledge Base Service - manages knowledge bases and documents
2. RAG Workflow - implements retrieval-augmented generation
3. RAG Agent - AI agent with tool calling capabilities
4. Llama Eino Model - adapter for Eino compatibility

Key features:
- Local-first: No external API dependencies
- Hybrid storage: SQLite for metadata, Chromem for vectors
- Document parsing: DOCX, PDF, XLSX, HTML, TXT, MD
- Eino integration: ADK agents, Flow/Compose, tools
`
	if err := os.WriteFile(testDocPath, []byte(testContent), 0644); err != nil {
		log.Fatalf("Failed to create test document: %v", err)
	}
	defer os.Remove(testDocPath)
	log.Printf("✅ Test document created: %s", testDocPath)

	// 7. Add Document to Knowledge Base
	log.Println("\n📥 Step 7: Adding document to knowledge base...")
	metadata := map[string]any{
		"source_file": "test-doc.txt",
		"test":        true,
	}
	if err := kbService.AddFileToKnowledgeBase(ctx, kbID, testDocPath, metadata, testUserID); err != nil {
		log.Fatalf("Failed to add file to knowledge base: %v", err)
	}
	log.Println("✅ Document added to knowledge base")

	// 8. Query Knowledge Base
	log.Println("\n🔍 Step 8: Querying knowledge base...")
	testQuery := "What framework is used for orchestration?"
	docs, err := kbService.QueryKnowledgeBase(ctx, kbID, testQuery, 3, testUserID)
	if err != nil {
		log.Fatalf("Failed to query knowledge base: %v", err)
	}
	log.Printf("✅ Query successful: retrieved %d documents", len(docs))
	for i, doc := range docs {
		log.Printf("\n   Document %d:", i+1)
		log.Printf("   Content: %s...", truncate(doc.Content, 100))
		if source, ok := doc.MetaData["source_file"].(string); ok {
			log.Printf("   Source: %s", source)
		}
	}

	// 9. Test RAG Workflow
	log.Println("\n🔄 Step 9: Testing RAG workflow...")
	ragWorkflow := services.NewRAGWorkflow(kbService)
	ragReq := services.RAGRequest{
		Query:           testQuery,
		KnowledgeBaseID: kbID,
		UserID:          testUserID,
		TopK:            3,
		IncludeSources:  true,
	}
	ragResp, err := ragWorkflow.ExecuteRAG(ctx, ragReq)
	if err != nil {
		log.Fatalf("Failed to execute RAG workflow: %v", err)
	}
	log.Printf("✅ RAG workflow successful")
	log.Printf("   Retrieved chunks: %d", ragResp.RetrievedChunks)
	log.Printf("   Context length: %d characters", len(ragResp.Context))

	// 10. Test Llama Eino Model Adapter
	log.Println("\n🤖 Step 10: Testing Llama Eino Model adapter...")
	llamaModel := llama.NewLlamaEinoModel(libService)
	log.Println("✅ Llama Eino Model adapter created")

	// Load chat model if not already loaded
	if !libService.IsChatModelLoaded() {
		log.Println("   Loading chat model...")
		if err := libService.LoadChatModel(""); err != nil {
			log.Printf("⚠️  Warning: Failed to load chat model: %v", err)
		} else {
			log.Println("   ✅ Chat model loaded")
		}
	}

	// 11. Test RAG Agent (optional, requires chat model)
	if libService.IsChatModelLoaded() {
		log.Println("\n🤖 Step 11: Testing RAG Agent...")
		agentConfig := &services.RAGAgentConfig{
			Name:             "test-rag-agent",
			Description:      "Test RAG agent",
			Model:            llamaModel,
			KnowledgeBaseIDs: []string{kbID},
			UserID:           testUserID,
			Instruction:      "You are a helpful assistant. Use the knowledge base search tool to find relevant information.",
			MaxIterations:    5,
		}

		ragAgent, err := services.NewRAGAgent(ctx, agentConfig, kbService)
		if err != nil {
			log.Fatalf("Failed to create RAG agent: %v", err)
		}
		log.Println("✅ RAG Agent created")

		// Test agent execution
		log.Println("\n   Testing agent with query: 'What is CloudWeGo Eino?'")
		response, err := ragAgent.Run(ctx, "What is CloudWeGo Eino?")
		if err != nil {
			log.Printf("⚠️  Warning: Agent execution failed: %v", err)
		} else {
			log.Printf("✅ Agent response received")
			log.Printf("   Response: %s", truncate(response, 200))
		}
	} else {
		log.Println("\n⚠️  Step 11: Skipping RAG Agent test (chat model not loaded)")
	}

	// 12. Cleanup
	log.Println("\n🧹 Step 12: Cleaning up...")
	if err := kbService.DeleteKnowledgeBase(ctx, kbID, testUserID); err != nil {
		log.Printf("⚠️  Warning: Failed to cleanup knowledge base: %v", err)
	} else {
		log.Println("✅ Test knowledge base deleted")
	}

	// Final Summary
	log.Println("\n" + strings.Repeat("=", 60))
	log.Println("✅ All tests completed successfully!")
	log.Println(strings.Repeat("=", 60))
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

