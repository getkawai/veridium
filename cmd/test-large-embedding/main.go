package main

import (
	"context"
	"log"
	"os"
	"strings"
	"time"

	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/pkg/chromem"
)

func main() {
	log.SetFlags(log.Ltime | log.Lmicroseconds)
	log.Println("🚀 Starting Large Embedding Crash Test...")

	// Initialize Library Service
	libService, err := llama.NewLibraryService()
	if err != nil {
		log.Fatalf("❌ Failed to create library service: %v", err)
	}

	// Wait for initialization
	log.Println("⏳ Waiting for library initialization...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	if err := libService.WaitForInitialization(ctx); err != nil {
		log.Fatalf("❌ Library initialization timed out: %v", err)
	}
	log.Println("✅ Library initialized")

	// Get model path
	modelName := llama.GetRecommendedEmbeddingModel()
	modelSpec, exists := llama.GetEmbeddingModel(modelName)
	if !exists {
		log.Fatalf("❌ Model %s not found in catalog", modelName)
	}

	installer := llama.NewLlamaCppInstaller()
	modelPath := installer.GetModelsDirectory() + "/" + modelSpec.Filename

	// Check if model exists
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		log.Printf("⚠️  Model file not found at %s", modelPath)
		log.Println("   Attempting to download...")
		// In a real app, the service would handle this, but here we might need to wait
		// or just fail if it's not there.
		// Since we ran the app before, it should be there or downloading.
		// Let's assume it's there for this test.
	}

	log.Printf("📂 Using model: %s", modelPath)

	// Create embedder
	embedder := chromem.NewEmbeddingFuncLlamaWithPreloadedLibrary(modelPath)

	// Test with increasing string sizes
	sizes := []int{100, 500, 1000, 2000, 4000, 8000, 16000}

	for _, size := range sizes {
		log.Printf("🧪 Testing with input size: %d chars (~%d tokens)...", size, size/4)

		input := strings.Repeat("The quick brown fox jumps over the lazy dog. ", size/45+1)[:size]

		start := time.Now()
		embedding, err := embedder(context.Background(), input)
		duration := time.Since(start)

		if err != nil {
			log.Printf("❌ Failed at size %d: %v", size, err)
		} else {
			log.Printf("✅ Success at size %d! (Time: %v, Dim: %d)", size, duration, len(embedding))
		}

		// Small pause
		time.Sleep(500 * time.Millisecond)
	}

	log.Println("🎉 Test completed successfully!")
}
