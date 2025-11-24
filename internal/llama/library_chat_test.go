package llama

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/wailsapp/wails/v3/pkg/application"
)

// TestChatCompletionRequest_Validation tests request validation
func TestChatCompletionRequest_Validation(t *testing.T) {
	tests := []struct {
		name        string
		req         ChatCompletionRequest
		shouldError bool
		errorMsg    string
	}{
		{
			name: "Valid request with all fields",
			req: ChatCompletionRequest{
				Model: "qwen2.5-0.5b-instruct-q8_0.gguf",
				Messages: []ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				MaxTokens:   100,
				Temperature: 0.7,
				TopP:        0.9,
				TopK:        40,
			},
			shouldError: false,
		},
		{
			name: "Valid request with minimal fields",
			req: ChatCompletionRequest{
				Messages: []ChatMessage{
					{Role: "user", Content: "Hello"},
				},
			},
			shouldError: false,
		},
		{
			name: "Empty messages",
			req: ChatCompletionRequest{
				Model:    "test-model",
				Messages: []ChatMessage{},
			},
			shouldError: true,
			errorMsg:    "messages cannot be empty",
		},
		{
			name: "Invalid role",
			req: ChatCompletionRequest{
				Messages: []ChatMessage{
					{Role: "invalid", Content: "Hello"},
				},
			},
			shouldError: true,
			errorMsg:    "invalid role",
		},
		{
			name: "Empty content",
			req: ChatCompletionRequest{
				Messages: []ChatMessage{
					{Role: "user", Content: ""},
				},
			},
			shouldError: true,
			errorMsg:    "message content cannot be empty",
		},
		{
			name: "Negative max tokens",
			req: ChatCompletionRequest{
				Messages: []ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				MaxTokens: -1,
			},
			shouldError: true,
			errorMsg:    "max_tokens must be positive",
		},
		{
			name: "Invalid temperature",
			req: ChatCompletionRequest{
				Messages: []ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				Temperature: -0.5,
			},
			shouldError: true,
			errorMsg:    "temperature must be between 0 and 2",
		},
		{
			name: "Invalid top_p",
			req: ChatCompletionRequest{
				Messages: []ChatMessage{
					{Role: "user", Content: "Hello"},
				},
				TopP: 1.5,
			},
			shouldError: true,
			errorMsg:    "top_p must be between 0 and 1",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateChatRequest(tt.req)
			if tt.shouldError {
				if err == nil {
					t.Errorf("Expected error containing '%s', got nil", tt.errorMsg)
				} else if !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("Expected error containing '%s', got: %v", tt.errorMsg, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, got: %v", err)
				}
			}
		})
	}
}

// TestBuildPrompt tests prompt building from messages
func TestBuildPrompt(t *testing.T) {
	// Create a mock chat service (no real model needed for prompt building)
	service := &LibraryChatService{
		libService: &LibraryService{},
	}

	tests := []struct {
		name     string
		messages []ChatMessage
		contains []string // Strings that should be in the prompt
	}{
		{
			name: "Single user message",
			messages: []ChatMessage{
				{Role: "user", Content: "Hello, how are you?"},
			},
			contains: []string{"User:", "Hello, how are you?", "Assistant:"},
		},
		{
			name: "System and user messages",
			messages: []ChatMessage{
				{Role: "system", Content: "You are a helpful assistant."},
				{Role: "user", Content: "What is 2+2?"},
			},
			contains: []string{"System:", "You are a helpful assistant", "User:", "What is 2+2?", "Assistant:"},
		},
		{
			name: "Multi-turn conversation",
			messages: []ChatMessage{
				{Role: "user", Content: "Hi"},
				{Role: "assistant", Content: "Hello!"},
				{Role: "user", Content: "How are you?"},
			},
			contains: []string{"User:", "Hi", "Assistant:", "Hello!", "How are you?"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := service.buildPrompt(tt.messages)

			if prompt == "" {
				t.Error("Prompt should not be empty")
			}

			for _, expected := range tt.contains {
				if !strings.Contains(prompt, expected) {
					t.Errorf("Prompt should contain '%s', got: %s", expected, prompt)
				}
			}

			t.Logf("Generated prompt:\n%s", prompt)
		})
	}
}

// TestChatCompletionResponse_JSON tests response serialization
func TestChatCompletionResponse_JSON(t *testing.T) {
	response := &ChatCompletionResponse{
		ID:      "chatcmpl-123",
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   "test-model",
		Choices: []ChatCompletionChoice{
			{
				Index: 0,
				Message: ChatMessage{
					Role:    "assistant",
					Content: "Hello! How can I help you?",
				},
				FinishReason: "stop",
			},
		},
		Usage: &ChatUsage{
			PromptTokens:     10,
			CompletionTokens: 8,
			TotalTokens:      18,
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// Test JSON unmarshaling
	var decoded ChatCompletionResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// Verify fields
	if decoded.ID != response.ID {
		t.Errorf("ID mismatch: got %s, want %s", decoded.ID, response.ID)
	}
	if decoded.Object != response.Object {
		t.Errorf("Object mismatch: got %s, want %s", decoded.Object, response.Object)
	}
	if len(decoded.Choices) != 1 {
		t.Errorf("Choices count mismatch: got %d, want 1", len(decoded.Choices))
	}
	if decoded.Choices[0].Message.Content != response.Choices[0].Message.Content {
		t.Errorf("Content mismatch: got %s, want %s",
			decoded.Choices[0].Message.Content,
			response.Choices[0].Message.Content)
	}

	t.Logf("✅ JSON serialization works correctly")
}

// TestChatCompletionChunk_JSON tests streaming chunk serialization
func TestChatCompletionChunk_JSON(t *testing.T) {
	chunk := ChatCompletionChunk{
		ID:      "chatcmpl-123",
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   "test-model",
		Choices: []ChatCompletionChunkChoice{
			{
				Index: 0,
				Delta: ChatMessageDelta{
					Role:    "assistant",
					Content: "Hello",
				},
			},
		},
	}

	// Test JSON marshaling
	data, err := json.Marshal(chunk)
	if err != nil {
		t.Fatalf("Failed to marshal chunk: %v", err)
	}

	// Verify it's valid JSON
	var decoded ChatCompletionChunk
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal chunk: %v", err)
	}

	if decoded.Choices[0].Delta.Content != "Hello" {
		t.Errorf("Content mismatch: got %s, want Hello", decoded.Choices[0].Delta.Content)
	}

	t.Logf("✅ Streaming chunk JSON works correctly")
}

// TestNewLibraryChatService tests service creation
func TestNewLibraryChatService(t *testing.T) {
	libService := &LibraryService{}
	app := &application.App{}

	chatService := NewLibraryChatService(libService, app)

	if chatService == nil {
		t.Fatal("Expected non-nil chat service")
	}

	if chatService.libService != libService {
		t.Error("LibraryService not set correctly")
	}

	if chatService.app != app {
		t.Error("App not set correctly")
	}

	t.Logf("✅ Service creation works correctly")
}

// Integration test - requires real model and library
func TestChatCompletion_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create library service
	libService, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create library service: %v", err)
	}
	defer libService.Cleanup()

	// Wait for initialization and model download (up to 10 minutes for large models)
	maxWait := 10 * time.Minute
	checkInterval := 5 * time.Second
	waited := 0 * time.Second

	t.Log("⏳ Waiting for library initialization and model download...")
	for waited < maxWait {
		// Check if model is loaded
		if libService.IsChatModelLoaded() {
			t.Logf("✅ Model loaded after %v", waited)
			break
		}

		// Check if any models are available
		installer := libService.installer
		models, _ := installer.GetAvailableChatModels()
		if len(models) > 0 {
			// Model exists but not loaded, try to load it
			t.Logf("📦 Found %d model(s), attempting to load...", len(models))
			if err := libService.LoadChatModel(""); err == nil {
				t.Logf("✅ Model loaded successfully after %v", waited)
				break
			} else {
				t.Logf("⚠️  Failed to load model: %v (will retry)", err)
			}
		} else {
			// No models found, trigger download if not already downloading
			modelsDir := installer.GetModelsDirectory()
			files, _ := os.ReadDir(modelsDir)
			downloading := false
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".tmp") {
					downloading = true
					break
				}
			}
			if downloading {
				t.Logf("📥 Model download in progress... (%v/%v)", waited, maxWait)
			} else if waited == 0 {
				// Trigger auto-download on first iteration
				t.Log("🚀 No chat models found, triggering auto-download...")
				go func() {
					if err := installer.AutoDownloadRecommendedChatModel(); err != nil {
						t.Logf("⚠️  Auto-download failed: %v", err)
					}
				}()
			}
		}

		time.Sleep(checkInterval)
		waited += checkInterval

		if waited%(1*time.Minute) == 0 {
			t.Logf("⏳ Still waiting... (%v/%v)", waited, maxWait)
		}
	}

	// Verify model is ready - FAIL if not ready after timeout
	if !libService.IsChatModelLoaded() {
		installer := libService.installer
		models, _ := installer.GetAvailableChatModels()
		if len(models) == 0 {
			t.Fatalf("❌ Test FAILED: No chat models available after %v wait time. Model download may have failed or timed out.", maxWait)
		}
		// Try one more time to load model
		if err := libService.LoadChatModel(""); err != nil {
			t.Fatalf("❌ Test FAILED: Failed to load chat model after %v wait time: %v", maxWait, err)
		}
	}

	// Final verification
	if !libService.IsChatModelLoaded() {
		t.Fatal("❌ Test FAILED: Chat model is not loaded after all attempts")
	}

	// Create chat service
	app := &application.App{}
	chatService := NewLibraryChatService(libService, app)

	// Test request
	req := ChatCompletionRequest{
		Messages: []ChatMessage{
			{Role: "user", Content: "Say 'test' and nothing else."},
		},
		MaxTokens:   10,
		Temperature: 0.1,
	}

	// Validate request
	if err := validateChatRequest(req); err != nil {
		t.Fatalf("Request validation failed: %v", err)
	}

	// Make request
	ctx := context.Background()
	resp, err := chatService.ChatCompletion(ctx, req)
	if err != nil {
		t.Fatalf("Chat completion failed: %v", err)
	}

	// Verify response
	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	if resp.Object != "chat.completion" {
		t.Errorf("Expected object='chat.completion', got: %s", resp.Object)
	}

	if len(resp.Choices) == 0 {
		t.Fatal("Expected at least one choice")
	}

	if resp.Choices[0].Message.Role != "assistant" {
		t.Errorf("Expected role='assistant', got: %s", resp.Choices[0].Message.Role)
	}

	if resp.Choices[0].Message.Content == "" {
		t.Error("Expected non-empty content")
	}

	if resp.Usage == nil {
		t.Error("Expected usage information")
	}

	t.Logf("✅ Integration test passed")
	t.Logf("Response: %s", resp.Choices[0].Message.Content)
	t.Logf("Usage: %+v", resp.Usage)
}

// TestCustomSamplerCreation tests that custom sampler is created with request parameters
func TestCustomSamplerCreation(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping sampler creation test in short mode")
	}

	// Create library service
	libService, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create library service: %v", err)
	}
	defer libService.Cleanup()

	// Wait for initialization and model download (up to 10 minutes)
	maxWait := 10 * time.Minute
	checkInterval := 5 * time.Second
	waited := 0 * time.Second

	for waited < maxWait {
		if libService.IsChatModelLoaded() {
			break
		}
		installer := libService.installer
		models, _ := installer.GetAvailableChatModels()
		if len(models) > 0 {
			if err := libService.LoadChatModel(""); err == nil {
				break
			}
		}
		time.Sleep(checkInterval)
		waited += checkInterval
	}

	// Verify model is ready - FAIL if not ready after timeout
	if !libService.IsChatModelLoaded() {
		installer := libService.installer
		models, _ := installer.GetAvailableChatModels()
		if len(models) == 0 {
			t.Fatalf("❌ Test FAILED: No chat models available after %v wait time", maxWait)
		}
		if err := libService.LoadChatModel(""); err != nil {
			t.Fatalf("❌ Test FAILED: Failed to load chat model after %v wait time: %v", maxWait, err)
		}
	}

	if !libService.IsChatModelLoaded() {
		t.Fatal("❌ Test FAILED: Chat model is not loaded after all attempts")
	}

	// Create chat service
	app := &application.App{}
	chatService := NewLibraryChatService(libService, app)

	// Test with custom parameters
	req := ChatCompletionRequest{
		Messages: []ChatMessage{
			{Role: "user", Content: "Test"},
		},
		MaxTokens:   10,
		Temperature: 0.5,
		TopP:        0.9,
		TopK:        20,
	}

	// This should trigger custom sampler creation
	ctx := context.Background()
	resp, err := chatService.ChatCompletion(ctx, req)
	if err != nil {
		t.Fatalf("Chat completion with custom sampler failed: %v", err)
	}

	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	t.Logf("✅ Custom sampler creation works")
	t.Logf("Response with temp=0.5, top_p=0.9, top_k=20: %s", resp.Choices[0].Message.Content)
}

// TestSamplerParameterVariations tests different sampler parameter combinations
func TestSamplerParameterVariations(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping parameter variations test in short mode")
	}

	// Create library service
	libService, err := NewLibraryService()
	if err != nil {
		t.Fatalf("Failed to create library service: %v", err)
	}
	defer libService.Cleanup()

	// Wait for initialization and model download (up to 10 minutes)
	maxWait := 10 * time.Minute
	checkInterval := 5 * time.Second
	waited := 0 * time.Second

	for waited < maxWait {
		if libService.IsChatModelLoaded() {
			break
		}
		installer := libService.installer
		models, _ := installer.GetAvailableChatModels()
		if len(models) > 0 {
			if err := libService.LoadChatModel(""); err == nil {
				break
			}
		}
		time.Sleep(checkInterval)
		waited += checkInterval
	}

	// Verify model is ready - FAIL if not ready after timeout
	if !libService.IsChatModelLoaded() {
		installer := libService.installer
		models, _ := installer.GetAvailableChatModels()
		if len(models) == 0 {
			t.Fatalf("❌ Test FAILED: No chat models available after %v wait time", maxWait)
		}
		if err := libService.LoadChatModel(""); err != nil {
			t.Fatalf("❌ Test FAILED: Failed to load chat model after %v wait time: %v", maxWait, err)
		}
	}

	if !libService.IsChatModelLoaded() {
		t.Fatal("❌ Test FAILED: Chat model is not loaded after all attempts")
	}

	// Create chat service
	app := &application.App{}
	chatService := NewLibraryChatService(libService, app)

	tests := []struct {
		name        string
		temperature float32
		topP        float32
		topK        int32
	}{
		{"Low temperature", 0.1, 0.9, 40},
		{"High temperature", 1.5, 0.9, 40},
		{"Low top_p", 0.7, 0.5, 40},
		{"High top_p", 0.7, 0.99, 40},
		{"Low top_k", 0.7, 0.9, 10},
		{"High top_k", 0.7, 0.9, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := ChatCompletionRequest{
				Messages: []ChatMessage{
					{Role: "user", Content: "Say 'test'"},
				},
				MaxTokens:   5,
				Temperature: tt.temperature,
				TopP:        tt.topP,
				TopK:        tt.topK,
			}

			ctx := context.Background()
			resp, err := chatService.ChatCompletion(ctx, req)
			if err != nil {
				t.Errorf("Failed with temp=%.1f, top_p=%.2f, top_k=%d: %v",
					tt.temperature, tt.topP, tt.topK, err)
				return
			}

			if resp == nil || len(resp.Choices) == 0 {
				t.Error("Expected valid response")
				return
			}

			t.Logf("✅ temp=%.1f, top_p=%.2f, top_k=%d → %s",
				tt.temperature, tt.topP, tt.topK, resp.Choices[0].Message.Content)
		})
	}
}

// Benchmark prompt building
func BenchmarkBuildPrompt(b *testing.B) {
	service := &LibraryChatService{
		libService: &LibraryService{},
	}

	messages := []ChatMessage{
		{Role: "system", Content: "You are a helpful assistant."},
		{Role: "user", Content: "Hello"},
		{Role: "assistant", Content: "Hi there!"},
		{Role: "user", Content: "How are you?"},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = service.buildPrompt(messages)
	}
}
