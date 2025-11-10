package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kawai-network/veridium/internal/llama"
)

func main() {
	fmt.Println("🚀 Testing Library-based llama.cpp Service")
	fmt.Println("==========================================")
	
	// Create library service
	fmt.Println("\n📦 Creating library service...")
	libService, err := llama.NewLibraryService()
	if err != nil {
		log.Fatalf("❌ Failed to create service: %v", err)
	}
	defer libService.Cleanup()
	
	fmt.Println("✅ Service created!")
	
	// Wait for background initialization
	fmt.Println("\n⏳ Waiting for initialization (10 seconds)...")
	time.Sleep(10 * time.Second)
	
	// Check available models
	fmt.Println("\n📋 Checking available models...")
	models, err := libService.GetAvailableModels()
	if err != nil {
		log.Printf("⚠️  Failed to get models: %v", err)
	} else {
		fmt.Printf("✅ Found %d model(s):\n", len(models))
		for i, model := range models {
			fmt.Printf("   %d. %s\n", i+1, model)
		}
	}
	
	// Try to load a model
	if len(models) > 0 {
		fmt.Println("\n🔄 Loading chat model...")
		err := libService.LoadChatModel("")
		if err != nil {
			log.Printf("⚠️  Failed to load model: %v", err)
			fmt.Println("\n💡 To test with a model:")
			fmt.Println("   1. Download a GGUF model from Hugging Face")
			fmt.Println("   2. Place it in ~/.llama-cpp/models/")
			fmt.Println("   3. Run this program again")
		} else {
			fmt.Println("✅ Model loaded successfully!")
			
			// Test generation
			fmt.Println("\n💬 Testing text generation...")
			prompt := "What is 2+2?"
			fmt.Printf("   Prompt: %s\n", prompt)
			
			response, err := libService.Generate(prompt, 50)
			if err != nil {
				log.Printf("⚠️  Generation failed: %v", err)
			} else {
				fmt.Printf("   Response: %s\n", response)
				fmt.Println("\n🎉 SUCCESS! Library-based service is working!")
			}
		}
	} else {
		fmt.Println("\n⚠️  No models available for testing")
		fmt.Println("\n💡 To download a model automatically:")
		fmt.Println("   The service will auto-download on first use")
		fmt.Println("   Or manually download from:")
		fmt.Println("   https://huggingface.co/models?library=gguf")
	}
	
	// Check embedding models
	fmt.Println("\n📋 Checking embedding models...")
	embModels := libService.GetEmbeddingManager().GetDownloadedModels()
	if len(embModels) > 0 {
		fmt.Printf("✅ Found %d embedding model(s):\n", len(embModels))
		for i, model := range embModels {
			fmt.Printf("   %d. %s (%s)\n", i+1, model.Name, model.Description)
		}
		
		// Try embedding generation
		fmt.Println("\n🔄 Loading embedding model...")
		err := libService.LoadEmbeddingModel("")
		if err != nil {
			log.Printf("⚠️  Failed to load embedding model: %v", err)
		} else {
			fmt.Println("✅ Embedding model loaded!")
			
			fmt.Println("\n🔢 Testing embedding generation...")
			text := "This is a test sentence."
			fmt.Printf("   Text: %s\n", text)
			
			embedding, err := libService.GenerateEmbedding(text)
			if err != nil {
				log.Printf("⚠️  Embedding generation failed: %v", err)
			} else {
				fmt.Printf("   Generated embedding: %d dimensions\n", len(embedding))
				fmt.Printf("   First 5 values: %v\n", embedding[:min(5, len(embedding))])
				fmt.Println("\n🎉 SUCCESS! Embedding generation is working!")
			}
		}
	} else {
		fmt.Println("⚠️  No embedding models available")
		fmt.Println("   Embedding models will be auto-downloaded on first use")
	}
	
	fmt.Println("\n==========================================")
	fmt.Println("✅ Library-based service test completed!")
	fmt.Println("\n📚 For more info, see:")
	fmt.Println("   - internal/llama/QUICK_START.md")
	fmt.Println("   - internal/llama/LIBRARY_USAGE.md")
	fmt.Println("   - internal/llama/README_MIGRASI.md")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

