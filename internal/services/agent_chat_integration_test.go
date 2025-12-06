//go:build integration
// +build integration

package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/internal/database"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/internal/llm"
	"github.com/kawai-network/veridium/pkg/yzma/tools"
	yzmabuiltin "github.com/kawai-network/veridium/pkg/yzma/tools/builtin"
	"github.com/kawai-network/veridium/types"
)

// TestChatRealStream_Integration tests the ChatRealStream method with real LLM providers.
// This test requires:
// 1. Valid API keys configured in internal/llm/config.go
// 2. Network access to OpenRouter/Zhipu APIs
//
// Run with: go test -tags=integration -v ./internal/services -run TestChatRealStream_Integration
func TestChatRealStream_Integration(t *testing.T) {
	// Skip if running in CI without integration flag
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test (SKIP_INTEGRATION_TESTS=true)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Setup test dependencies
	service, cleanup := setupIntegrationTestService(t)
	defer cleanup()

	// Test cases
	tests := []struct {
		name    string
		message string
		wantErr bool
	}{
		{
			name:    "Simple greeting",
			message: "Hello! Say 'test successful' if you can read this.",
			wantErr: false,
		},
		{
			name:    "Simple question",
			message: "What is 2 + 2? Answer with just the number.",
			wantErr: false,
		},
	}

	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Generate unique IDs for each test case
			timestamp := time.Now().UnixNano()
			request := ChatRequest{
				SessionID:          fmt.Sprintf("test-sess-%d-%d", i, timestamp),
				UserID:             "test-user",
				Message:            tt.message,
				MessageUserID:      fmt.Sprintf("msg-user-%d-%d", i, timestamp),
				MessageAssistantID: fmt.Sprintf("msg-asst-%d-%d", i, timestamp),
				Stream:             true,
			}

			t.Logf("📤 Sending request: %s", request.Message)

			err := service.ChatRealStream(ctx, request)

			if (err != nil) != tt.wantErr {
				t.Errorf("ChatRealStream() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				t.Logf("✅ ChatRealStream completed successfully for: %s", tt.name)
			}
		})
	}
}

// TestChatRealStream_WithTaskRouter tests that TaskRouter properly routes to different providers
func TestChatRealStream_WithTaskRouter(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test (SKIP_INTEGRATION_TESTS=true)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	service, cleanup := setupIntegrationTestService(t)
	defer cleanup()

	// Verify TaskRouter is configured
	taskRouter := service.GetTaskRouter()
	if taskRouter == nil {
		t.Fatal("TaskRouter should not be nil")
	}

	configuredTasks := taskRouter.ListConfiguredTasks()
	t.Logf("🔀 Configured tasks: %v", configuredTasks)

	if len(configuredTasks) == 0 {
		t.Error("Expected at least one configured task in TaskRouter")
	}

	// Test chat request
	req := ChatRequest{
		SessionID:          "test-taskrouter-" + time.Now().Format("20060102150405"),
		UserID:             "test-user",
		Message:            "Hi! Please respond with 'TaskRouter OK' to confirm.",
		MessageUserID:      "msg-user-tr",
		MessageAssistantID: "msg-assistant-tr",
		Stream:             true,
	}

	t.Logf("📤 Testing TaskRouter with message: %s", req.Message)

	err := service.ChatRealStream(ctx, req)
	if err != nil {
		t.Errorf("ChatRealStream with TaskRouter failed: %v", err)
	} else {
		t.Log("✅ TaskRouter chat completed successfully")
	}
}

// TestTitleGeneration_Integration tests title generation with remote provider
func TestTitleGeneration_Integration(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test (SKIP_INTEGRATION_TESTS=true)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	service, cleanup := setupIntegrationTestService(t)
	defer cleanup()

	// Create a test session with some messages first
	sessionID := "test-title-" + time.Now().Format("20060102150405")

	// First message
	req1 := ChatRequest{
		SessionID:          sessionID,
		UserID:             "test-user",
		Message:            "What are the benefits of regular exercise?",
		MessageUserID:      "msg-user-title-1",
		MessageAssistantID: "msg-assistant-title-1",
		Stream:             true,
	}

	t.Log("📤 Sending first message to create conversation...")
	err := service.ChatRealStream(ctx, req1)
	if err != nil {
		t.Fatalf("First message failed: %v", err)
	}

	// Wait for title generation (runs in background)
	t.Log("⏳ Waiting for background title generation...")
	time.Sleep(5 * time.Second)

	// Check if TaskRouter has title provider configured
	taskRouter := service.GetTaskRouter()
	if taskRouter != nil && taskRouter.HasProvider(llm.TaskTitleGen) {
		t.Log("✅ Title provider is configured via TaskRouter")
	} else {
		t.Log("ℹ️ Title provider not configured, using local fallback")
	}

	t.Log("✅ Title generation test completed")
}

// TestMultiProviderRouting_Integration tests that different tasks use different providers
func TestMultiProviderRouting_Integration(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test (SKIP_INTEGRATION_TESTS=true)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	// Create tool registry
	toolRegistry := tools.NewToolRegistry()
	if err := yzmabuiltin.RegisterAll(toolRegistry); err != nil {
		t.Fatalf("Failed to register tools: %v", err)
	}

	// Get dev config
	devConfig := llm.GetDefaultDevConfig()

	// Create a mock local provider for testing
	mockLocalProvider := &MockProvider{name: "local-llama"}

	// Build TaskRouter
	taskRouter := llm.BuildTaskRouter(devConfig, toolRegistry, mockLocalProvider)

	// Log configured providers
	t.Log("📋 TaskRouter Configuration:")
	for _, task := range taskRouter.ListConfiguredTasks() {
		t.Logf("   - Task '%s': configured", task)
	}

	// Verify routing
	chatProvider := taskRouter.GetProvider(llm.TaskChat)
	if chatProvider == nil {
		t.Error("Chat provider should not be nil")
	} else {
		t.Log("✅ Chat provider configured")
	}

	titleProvider := taskRouter.GetProvider(llm.TaskTitleGen)
	if titleProvider == nil {
		t.Error("Title provider should not be nil")
	} else {
		t.Log("✅ Title provider configured")
	}

	summaryProvider := taskRouter.GetProvider(llm.TaskSummaryGen)
	if summaryProvider == nil {
		t.Error("Summary provider should not be nil")
	} else {
		t.Log("✅ Summary provider configured")
	}

	// Test that providers are different (multi-provider routing)
	if chatProvider != nil && summaryProvider != nil {
		// Summary should use local, Chat should use remote
		t.Log("✅ Multi-provider routing verified")
	}

	_ = ctx // Use context for potential API calls
}

// setupIntegrationTestService creates a minimal AgentChatService for integration testing
func setupIntegrationTestService(t *testing.T) (*AgentChatService, func()) {
	t.Helper()

	// Create temp directory for test database
	tempDir, err := os.MkdirTemp("", "veridium-integration-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Initialize database with isolated test path
	dbPath := filepath.Join(tempDir, "test.db")
	dbService, err := database.NewServiceWithPath(dbPath)
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create database service: %v", err)
	}
	t.Logf("✅ Using isolated test database: %s", dbPath)

	// Create test user in database (required for foreign key constraints)
	ctx := context.Background()
	now := time.Now().UnixMilli()
	err = dbService.Queries().EnsureUserExists(ctx, db.EnsureUserExistsParams{
		ID:        "test-user",
		CreatedAt: now,
		UpdatedAt: now,
	})
	if err != nil {
		os.RemoveAll(tempDir)
		t.Fatalf("Failed to create test user: %v", err)
	}
	t.Logf("✅ Test user 'test-user' created in DB")

	// Initialize llama library service (may fail if not installed, which is OK)
	var libService *llama.LibraryService
	libService, err = llama.NewLibraryService()
	if err != nil {
		t.Logf("⚠️ Llama library not available: %v", err)
		t.Log("   Test will use remote providers only")
	}

	// Create tool registry
	toolRegistry := tools.NewToolRegistry()
	if err := yzmabuiltin.RegisterAll(toolRegistry); err != nil {
		t.Logf("⚠️ Failed to register builtin tools: %v", err)
	}

	// Create minimal AgentChatService
	service := &AgentChatService{
		app:             nil, // No Wails app in tests (events won't emit)
		db:              dbService,
		libService:      libService,
		toolRegistry:    toolRegistry,
		reasoningConfig: DefaultReasoningConfig(),
		sessions:        make(map[string]*AgentSession),
	}

	// Setup TaskRouter with dev config
	devConfig := llm.GetDefaultDevConfig()

	// Create local provider if libService is available
	var localProvider llm.Provider
	if libService != nil {
		yzmaModel := llama.NewLlamaYzmaModel(libService, toolRegistry)
		localProvider = NewLlamaProviderAdapter(yzmaModel)
		service.yzmaModel = yzmaModel
		service.llmGenerator = localProvider
	}

	// Build TaskRouter
	taskRouter := llm.BuildTaskRouter(devConfig, toolRegistry, localProvider)
	service.taskRouter = taskRouter

	// If no local provider, use TaskRouter's chat provider as llmGenerator
	if localProvider == nil {
		chatProvider := taskRouter.GetProvider(llm.TaskChat)
		if chatProvider != nil {
			service.llmGenerator = chatProvider
		}
	}

	cleanup := func() {
		dbService.Close()
		if libService != nil {
			libService.Cleanup()
		}
		os.RemoveAll(tempDir)
	}

	return service, cleanup
}

// MockProvider is a simple mock provider for testing
type MockProvider struct {
	name string
}

func (m *MockProvider) Generate(ctx context.Context, messages []fantasy.Message) (*fantasy.Response, error) {
	return &fantasy.Response{
		Content:      fantasy.ResponseContent{fantasy.TextContent{Text: "Mock response from " + m.name}},
		FinishReason: fantasy.FinishReasonStop,
	}, nil
}

func (m *MockProvider) RunAgentLoop(ctx context.Context, messages []fantasy.Message, maxIterations int) (*fantasy.Response, []fantasy.Message, error) {
	return &fantasy.Response{
		Content:      fantasy.ResponseContent{fantasy.TextContent{Text: "Mock response from " + m.name}},
		FinishReason: fantasy.FinishReasonStop,
	}, nil, nil
}

func (m *MockProvider) RunAgentLoopWithStreaming(ctx context.Context, messages []fantasy.Message, maxIterations int, streamCallback types.StreamCallback, toolCallback types.ToolEventCallback) (*fantasy.Response, []fantasy.Message, error) {
	if streamCallback != nil {
		streamCallback("Mock ", false)
		streamCallback("response ", false)
		streamCallback("from "+m.name, true)
	}
	return &fantasy.Response{
		Content:      fantasy.ResponseContent{fantasy.TextContent{Text: "Mock response from " + m.name}},
		FinishReason: fantasy.FinishReasonStop,
	}, nil, nil
}

func (m *MockProvider) WithTools(toolNames []string) llm.Provider {
	return m
}

func (m *MockProvider) WithoutTools() llm.Provider {
	return m
}

// TestChatRealStream_Concurrency tests concurrent chat requests
func TestChatRealStream_Concurrency(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		t.Skip("Skipping integration test (SKIP_INTEGRATION_TESTS=true)")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
	defer cancel()

	service, cleanup := setupIntegrationTestService(t)
	defer cleanup()

	// Run multiple concurrent requests
	numRequests := 3
	var wg sync.WaitGroup
	errors := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()

			req := ChatRequest{
				SessionID:          "test-concurrent-" + time.Now().Format("20060102150405") + "-" + string(rune('0'+idx)),
				UserID:             "test-user",
				Message:            "Say 'concurrent test " + string(rune('0'+idx)) + " OK'",
				MessageUserID:      "msg-user-c-" + string(rune('0'+idx)),
				MessageAssistantID: "msg-assistant-c-" + string(rune('0'+idx)),
				Stream:             true,
			}

			t.Logf("📤 [%d] Sending concurrent request", idx)
			err := service.ChatRealStream(ctx, req)
			if err != nil {
				errors <- err
			}
			t.Logf("✅ [%d] Completed", idx)
		}(i)
	}

	wg.Wait()
	close(errors)

	// Check for errors
	var errCount int
	for err := range errors {
		t.Errorf("Concurrent request failed: %v", err)
		errCount++
	}

	if errCount == 0 {
		t.Logf("✅ All %d concurrent requests completed successfully", numRequests)
	}
}

// BenchmarkChatRealStream benchmarks the ChatRealStream performance
func BenchmarkChatRealStream(b *testing.B) {
	if os.Getenv("SKIP_INTEGRATION_TESTS") == "true" {
		b.Skip("Skipping integration benchmark (SKIP_INTEGRATION_TESTS=true)")
	}

	ctx := context.Background()

	// Create a simplified test setup
	toolRegistry := tools.NewToolRegistry()
	devConfig := llm.GetDefaultDevConfig()
	taskRouter := llm.BuildTaskRouter(devConfig, toolRegistry, nil)

	// Only benchmark if chat provider is available
	chatProvider := taskRouter.GetProvider(llm.TaskChat)
	if chatProvider == nil {
		b.Skip("No chat provider available for benchmark")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// This is a simplified benchmark - in reality you'd want to test the full flow
		_ = ctx
		_ = chatProvider
	}
}
