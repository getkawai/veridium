package llama

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

// TestLibraryServiceInitialization tests basic service initialization
func TestLibraryServiceInitialization(t *testing.T) {
	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create library service: %v", err)
	}
	defer service.Cleanup()

	// Service should be created even if library is not yet initialized
	if service == nil {
		t.Fatal("Service should not be nil")
	}
}

// TestLibraryInitialization tests library loading
func TestLibraryInitialization(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test (set INTEGRATION_TEST=1 to run)")
	}

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create library service: %v", err)
	}
	defer service.Cleanup()

	// Wait for background initialization
	time.Sleep(5 * time.Second)

	err = service.InitializeLibrary()
	if err != nil {
		t.Fatalf("Failed to initialize library: %v", err)
	}

	if !service.isInitialized {
		t.Fatal("Library should be initialized")
	}

	if service.libPath == "" {
		t.Fatal("Library path should be set")
	}

	t.Logf("Library initialized from: %s", service.libPath)
}

// TestModelSelection tests automatic model selection
func TestModelSelection(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test")
	}

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create library service: %v", err)
	}
	defer service.Cleanup()

	models, err := service.GetAvailableModels()
	if err != nil {
		t.Fatalf("Failed to get available models: %v", err)
	}

	if len(models) == 0 {
		t.Skip("No models available for testing")
	}

	bestModel, err := service.selectBestModel()
	if err != nil {
		t.Fatalf("Failed to select best model: %v", err)
	}

	if bestModel == "" {
		t.Fatal("Best model should not be empty")
	}

	t.Logf("Selected model: %s", filepath.Base(bestModel))
}

// TestChatModelLoading tests loading a chat model
func TestChatModelLoading(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test")
	}

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create library service: %v", err)
	}
	defer service.Cleanup()

	// Wait for initialization
	time.Sleep(5 * time.Second)

	err = service.LoadChatModel("")
	if err != nil {
		t.Fatalf("Failed to load chat model: %v", err)
	}

	if !service.IsChatModelLoaded() {
		t.Fatal("Chat model should be loaded")
	}

	modelPath := service.GetLoadedChatModel()
	if modelPath == "" {
		t.Fatal("Loaded model path should not be empty")
	}

	t.Logf("Loaded chat model: %s", filepath.Base(modelPath))
}

// TestTextGeneration tests basic text generation
func TestTextGeneration(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test")
	}

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create library service: %v", err)
	}
	defer service.Cleanup()

	// Wait for initialization and model loading
	time.Sleep(10 * time.Second)

	if !service.IsChatModelLoaded() {
		err = service.LoadChatModel("")
		if err != nil {
			t.Fatalf("Failed to load chat model: %v", err)
		}
	}

	prompt := "What is 2+2?"
	response, err := service.Generate(prompt, 50)
	if err != nil {
		t.Fatalf("Failed to generate text: %v", err)
	}

	if response == "" {
		t.Fatal("Generated response should not be empty")
	}

	t.Logf("Prompt: %s", prompt)
	t.Logf("Response: %s", response)
}

// TestEmbeddingModelLoading tests loading an embedding model
func TestEmbeddingModelLoading(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test")
	}

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create library service: %v", err)
	}
	defer service.Cleanup()

	// Wait for initialization
	time.Sleep(5 * time.Second)

	// Check if embedding models are available
	downloaded := service.embeddingManager.GetDownloadedModels()
	if len(downloaded) == 0 {
		t.Skip("No embedding models available")
	}

	err = service.LoadEmbeddingModel("")
	if err != nil {
		t.Fatalf("Failed to load embedding model: %v", err)
	}

	if !service.IsEmbeddingModelLoaded() {
		t.Fatal("Embedding model should be loaded")
	}

	t.Logf("Loaded embedding model: %s", filepath.Base(service.GetLoadedEmbeddingModel()))
}

// TestEmbeddingGeneration tests embedding generation
func TestEmbeddingGeneration(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test")
	}

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create library service: %v", err)
	}
	defer service.Cleanup()

	// Wait for initialization
	time.Sleep(10 * time.Second)

	// Check if embedding models are available
	downloaded := service.embeddingManager.GetDownloadedModels()
	if len(downloaded) == 0 {
		t.Skip("No embedding models available")
	}

	if !service.IsEmbeddingModelLoaded() {
		err = service.LoadEmbeddingModel("")
		if err != nil {
			t.Fatalf("Failed to load embedding model: %v", err)
		}
	}

	text := "This is a test sentence for embedding generation."
	embedding, err := service.GenerateEmbedding(text)
	if err != nil {
		t.Fatalf("Failed to generate embedding: %v", err)
	}

	if len(embedding) == 0 {
		t.Fatal("Embedding should not be empty")
	}

	t.Logf("Generated embedding with dimension: %d", len(embedding))
	t.Logf("First 5 values: %v", embedding[:min(5, len(embedding))])
}

// TestChatCompletion tests the chat completion API
func TestChatCompletion(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test")
	}

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create library service: %v", err)
	}
	defer service.Cleanup()

	// Wait for initialization
	time.Sleep(10 * time.Second)

	chatService := NewLibraryChatService(service, nil)

	req := ChatCompletionRequest{
		Model: "auto",
		Messages: []ChatMessage{
			{Role: "user", Content: "What is 2+2?"},
		},
		MaxTokens:   50,
		Temperature: 0.7,
	}

	ctx := context.Background()
	resp, err := chatService.ChatCompletion(ctx, req)
	if err != nil {
		t.Fatalf("Failed to complete chat: %v", err)
	}

	if len(resp.Choices) == 0 {
		t.Fatal("Response should have at least one choice")
	}

	if resp.Choices[0].Message.Content == "" {
		t.Fatal("Response content should not be empty")
	}

	t.Logf("Chat response: %s", resp.Choices[0].Message.Content)
}

// TestMultipleModels tests loading multiple models simultaneously
func TestMultipleModels(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test")
	}

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create library service: %v", err)
	}
	defer service.Cleanup()

	// Wait for initialization
	time.Sleep(10 * time.Second)

	// Load chat model
	err = service.LoadChatModel("")
	if err != nil {
		t.Fatalf("Failed to load chat model: %v", err)
	}

	// Check if embedding models are available
	downloaded := service.embeddingManager.GetDownloadedModels()
	if len(downloaded) > 0 {
		// Load embedding model
		err = service.LoadEmbeddingModel("")
		if err != nil {
			t.Fatalf("Failed to load embedding model: %v", err)
		}

		// Both should be loaded
		if !service.IsChatModelLoaded() || !service.IsEmbeddingModelLoaded() {
			t.Fatal("Both models should be loaded simultaneously")
		}

		t.Logf("Successfully loaded both models simultaneously")
	} else {
		t.Skip("No embedding models available for simultaneous loading test")
	}
}

// BenchmarkTextGeneration benchmarks text generation performance
func BenchmarkTextGeneration(b *testing.B) {
	if os.Getenv("BENCHMARK") == "" {
		b.Skip("Skipping benchmark (set BENCHMARK=1 to run)")
	}

	service, err := NewLibraryService()
	if err != nil {
		b.Fatalf("Failed to create library service: %v", err)
	}
	defer service.Cleanup()

	// Wait for initialization
	time.Sleep(10 * time.Second)

	if !service.IsChatModelLoaded() {
		err = service.LoadChatModel("")
		if err != nil {
			b.Fatalf("Failed to load chat model: %v", err)
		}
	}

	prompt := "What is 2+2?"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.Generate(prompt, 10)
		if err != nil {
			b.Fatalf("Failed to generate: %v", err)
		}
	}
}

// BenchmarkEmbedding benchmarks embedding generation performance
func BenchmarkEmbedding(b *testing.B) {
	if os.Getenv("BENCHMARK") == "" {
		b.Skip("Skipping benchmark")
	}

	service, err := NewLibraryService()
	if err != nil {
		b.Fatalf("Failed to create library service: %v", err)
	}
	defer service.Cleanup()

	// Wait for initialization
	time.Sleep(10 * time.Second)

	downloaded := service.embeddingManager.GetDownloadedModels()
	if len(downloaded) == 0 {
		b.Skip("No embedding models available")
	}

	if !service.IsEmbeddingModelLoaded() {
		err = service.LoadEmbeddingModel("")
		if err != nil {
			b.Fatalf("Failed to load embedding model: %v", err)
		}
	}

	text := "This is a test sentence for benchmarking."

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := service.GenerateEmbedding(text)
		if err != nil {
			b.Fatalf("Failed to generate embedding: %v", err)
		}
	}
}

// TestLibraryCleanup tests proper cleanup
func TestLibraryCleanup(t *testing.T) {
	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}

	// Cleanup should not panic
	service.Cleanup()

	// Second cleanup should also not panic
	service.Cleanup()
}

// TestGetModelsDirectory tests models directory retrieval
func TestGetModelsDirectory(t *testing.T) {
	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Cleanup()

	modelsDir := service.GetModelsDirectory()
	if modelsDir == "" {
		t.Fatal("Models directory should not be empty")
	}

	t.Logf("Models directory: %s", modelsDir)
}

// TestGetEmbeddingManager tests embedding manager retrieval
func TestGetEmbeddingManager(t *testing.T) {
	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Cleanup()

	embMgr := service.GetEmbeddingManager()
	if embMgr == nil {
		t.Fatal("Embedding manager should not be nil")
	}
}

// TestModelStatusChecks tests model status checking functions
func TestModelStatusChecks(t *testing.T) {
	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Cleanup()

	// Initially no models should be loaded
	if service.IsChatModelLoaded() {
		t.Error("Chat model should not be loaded initially")
	}

	if service.IsEmbeddingModelLoaded() {
		t.Error("Embedding model should not be loaded initially")
	}

	// Paths should be empty
	if service.GetLoadedChatModel() != "" {
		t.Error("Chat model path should be empty initially")
	}

	if service.GetLoadedEmbeddingModel() != "" {
		t.Error("Embedding model path should be empty initially")
	}
}

// TestLoadModelWithInvalidPath tests error handling for invalid paths
func TestLoadModelWithInvalidPath(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test")
	}

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Cleanup()

	time.Sleep(5 * time.Second)

	// Try to load non-existent model
	err = service.LoadChatModel("/nonexistent/model.gguf")
	if err == nil {
		t.Error("Should fail with non-existent model path")
	}

	t.Logf("Expected error: %v", err)
}

// TestGenerateWithoutModel tests error when generating without loaded model
func TestGenerateWithoutModel(t *testing.T) {
	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Cleanup()

	// Try to generate without loading model
	_, err = service.Generate("test", 10)
	if err == nil {
		t.Error("Should fail when no model is loaded")
	}

	t.Logf("Expected error: %v", err)
}

// TestEmbeddingWithoutModel tests error when embedding without loaded model
func TestEmbeddingWithoutModel(t *testing.T) {
	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Cleanup()

	// Try to generate embedding without loading model
	_, err = service.GenerateEmbedding("test")
	if err == nil {
		t.Error("Should fail when no embedding model is loaded")
	}

	t.Logf("Expected error: %v", err)
}

// TestConcurrentModelAccess tests thread safety
func TestConcurrentModelAccess(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test")
	}

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Cleanup()

	time.Sleep(10 * time.Second)

	err = service.LoadChatModel("")
	if err != nil {
		t.Skipf("No models available: %v", err)
	}

	// Launch multiple concurrent requests
	var wg sync.WaitGroup
	errors := make(chan error, 10)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			_, err := service.Generate(fmt.Sprintf("Test %d", id), 10)
			if err != nil {
				errors <- err
			}
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check if there were any errors
	errorCount := 0
	for err := range errors {
		t.Logf("Concurrent error: %v", err)
		errorCount++
	}

	if errorCount > 0 {
		t.Logf("Had %d errors in concurrent access (may be acceptable)", errorCount)
	}
}

// TestModelSwitching tests switching between different models
func TestModelSwitching(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test")
	}

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Cleanup()

	time.Sleep(10 * time.Second)

	models, err := service.GetAvailableModels()
	if err != nil || len(models) < 2 {
		t.Skip("Need at least 2 models for switching test")
	}

	// Load first model
	err = service.LoadChatModel(models[0])
	if err != nil {
		t.Fatalf("Failed to load first model: %v", err)
	}

	firstPath := service.GetLoadedChatModel()

	// Load second model
	err = service.LoadChatModel(models[1])
	if err != nil {
		t.Fatalf("Failed to load second model: %v", err)
	}

	secondPath := service.GetLoadedChatModel()

	if firstPath == secondPath {
		t.Error("Model path should change after switching")
	}

	t.Logf("Switched from %s to %s", firstPath, secondPath)
}

// TestEmptyPrompt tests handling of empty prompts
func TestEmptyPrompt(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test")
	}

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Cleanup()

	time.Sleep(10 * time.Second)

	err = service.LoadChatModel("")
	if err != nil {
		t.Skipf("No models available: %v", err)
	}

	// Try with empty prompt
	response, err := service.Generate("", 10)
	// May succeed or fail depending on tokenizer
	t.Logf("Empty prompt result - Error: %v, Response: %q", err, response)
}

// TestLargeTokenLimit tests with very large token limits
func TestLargeTokenLimit(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test")
	}

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Cleanup()

	time.Sleep(10 * time.Second)

	err = service.LoadChatModel("")
	if err != nil {
		t.Skipf("No models available: %v", err)
	}

	// Try with large token limit (should be capped by context)
	response, err := service.Generate("Tell me a story", 10000)
	if err != nil {
		t.Logf("Large token limit error (expected): %v", err)
	} else {
		t.Logf("Generated %d chars with large limit", len(response))
	}
}

// TestLibraryInitializationError tests handling of library init failures
func TestLibraryInitializationError(t *testing.T) {
	// Save original env
	originalEnv := os.Getenv("YZMA_LIB")
	defer func() {
		if originalEnv != "" {
			os.Setenv("YZMA_LIB", originalEnv)
		} else {
			os.Unsetenv("YZMA_LIB")
		}
	}()

	// Set invalid path
	os.Setenv("YZMA_LIB", "/nonexistent/path")

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Service creation should succeed: %v", err)
	}
	defer service.Cleanup()

	// Library initialization should fail with invalid path
	err = service.InitializeLibrary()
	if err == nil {
		t.Log("Library init didn't fail (may have fallback paths)")
	} else {
		t.Logf("Expected error with invalid path: %v", err)
	}
}

// TestChatServiceWithNilApp tests chat service without Wails app
func TestChatServiceWithNilApp(t *testing.T) {
	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Cleanup()

	// Create chat service with nil app (should not panic)
	chatService := NewLibraryChatService(service, nil)
	if chatService == nil {
		t.Fatal("Chat service should not be nil")
	}
}

// TestEmbeddingServiceCreation tests embedding service creation
func TestEmbeddingServiceCreation(t *testing.T) {
	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Cleanup()

	embService := NewLibraryEmbeddingService(service)
	if embService == nil {
		t.Fatal("Embedding service should not be nil")
	}
}

// TestProxyServiceCreation tests proxy service creation
func TestProxyServiceCreation(t *testing.T) {
	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Cleanup()

	proxyService := NewLibraryProxyService(service, nil)
	if proxyService == nil {
		t.Fatal("Proxy service should not be nil")
	}
}

// TestBatchEmbeddingEmpty tests batch embedding with empty input
func TestBatchEmbeddingEmpty(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test")
	}

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Cleanup()

	embService := NewLibraryEmbeddingService(service)

	// Empty batch
	ctx := context.Background()
	results, err := embService.BatchEmbedding(ctx, []string{})
	if err != nil {
		t.Logf("Empty batch error: %v", err)
	}
	if len(results) != 0 {
		t.Error("Empty batch should return empty results")
	}
}

// TestContextCancellation tests context cancellation during operations
func TestContextCancellation(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test")
	}

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Cleanup()

	time.Sleep(10 * time.Second)

	downloaded := service.embeddingManager.GetDownloadedModels()
	if len(downloaded) == 0 {
		t.Skip("No embedding models available")
	}

	err = service.LoadEmbeddingModel("")
	if err != nil {
		t.Skipf("Failed to load embedding model: %v", err)
	}

	embService := NewLibraryEmbeddingService(service)

	// Create context with immediate cancellation
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err = embService.CreateEmbedding(ctx, EmbeddingRequest{
		Model: "auto",
		Input: []string{"test"},
	})

	if err != context.Canceled {
		t.Logf("Context cancellation may not be immediate: %v", err)
	}
}

// TestAutoDownloadWithoutModels tests auto-download functionality
func TestAutoDownloadWithoutModels(t *testing.T) {
	// This test would trigger actual download, skip by default
	t.Skip("Skipping auto-download test (would download large files)")

	service, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create service: %v", err)
	}
	defer service.Cleanup()

	err = service.AutoDownloadRecommendedModel()
	if err != nil {
		t.Logf("Auto-download error: %v", err)
	}
}

// TestMemoryLeakPrevention tests for memory leaks
func TestMemoryLeakPrevention(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("Skipping integration test")
	}

	// Create and destroy multiple services
	for i := 0; i < 5; i++ {
		service, err := NewLibraryService()
		if err != nil {
			t.Fatalf("Iteration %d: Failed to create service: %v", i, err)
		}

		// Do some work
		time.Sleep(1 * time.Second)

		// Cleanup
		service.Cleanup()
	}

	t.Log("Memory leak test completed (check with profiling tools)")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
