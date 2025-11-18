package contextengine

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewContextEngineService tests service creation
func TestNewContextEngineService(t *testing.T) {
	service := NewContextEngineService()
	assert.NotNil(t, service, "Service should not be nil")
}

// TestProcessMessages_BasicScenario tests basic message processing
func TestProcessMessages_BasicScenario(t *testing.T) {
	service := NewContextEngineService()

	request := ContextEngineeringRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: "Hello, world!",
			},
		},
		Model:    "gpt-4",
		Provider: "openai",
	}

	response := service.ProcessMessages(request)

	assert.Empty(t, response.Error, "Should not have error")
	assert.NotEmpty(t, response.Messages, "Should have messages")
	assert.GreaterOrEqual(t, len(response.Messages), 1, "Should have at least one message")
}

// TestProcessMessages_WithSystemRole tests processing with system role
func TestProcessMessages_WithSystemRole(t *testing.T) {
	service := NewContextEngineService()

	request := ContextEngineeringRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
		Model:      "gpt-4",
		Provider:   "openai",
		SystemRole: "You are a helpful assistant",
	}

	response := service.ProcessMessages(request)

	assert.Empty(t, response.Error, "Should not have error")
	assert.NotEmpty(t, response.Messages, "Should have messages")

	// Check if system role was injected
	hasSystemRole := false
	for _, msg := range response.Messages {
		if msg.Role == "system" {
			hasSystemRole = true
			break
		}
	}
	assert.True(t, hasSystemRole, "Should have system role message")
}

// TestProcessMessages_WithHistoryTruncation tests history truncation
func TestProcessMessages_WithHistoryTruncation(t *testing.T) {
	service := NewContextEngineService()

	// Create 10 messages
	messages := make([]Message, 10)
	for i := 0; i < 10; i++ {
		messages[i] = Message{
			Role:    "user",
			Content: "Message " + string(rune(i)),
		}
	}

	request := ContextEngineeringRequest{
		Messages:           messages,
		Model:              "gpt-4",
		Provider:           "openai",
		EnableHistoryCount: true,
		HistoryCount:       5,
	}

	response := service.ProcessMessages(request)

	assert.Empty(t, response.Error, "Should not have error")
	assert.NotEmpty(t, response.Messages, "Should have messages")
	// History truncation should keep last N messages
	assert.LessOrEqual(t, len(response.Messages), 5+2, "Should truncate to history count + buffer")
}

// TestProcessMessages_WithTools tests processing with tools
func TestProcessMessages_WithTools(t *testing.T) {
	service := NewContextEngineService()

	request := ContextEngineeringRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: "Search for information",
			},
		},
		Model:    "gpt-4",
		Provider: "openai",
		Tools:    []string{"web-search", "calculator"},
	}

	response := service.ProcessMessages(request)

	assert.Empty(t, response.Error, "Should not have error")
	assert.NotEmpty(t, response.Messages, "Should have messages")
}

// TestProcessMessages_WithInputTemplate tests input template processing
func TestProcessMessages_WithInputTemplate(t *testing.T) {
	service := NewContextEngineService()

	request := ContextEngineeringRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: "What is TypeScript?",
			},
		},
		Model:         "gpt-4",
		Provider:      "openai",
		InputTemplate: "User query: {{input}}\n\nPlease provide a detailed answer.",
	}

	response := service.ProcessMessages(request)

	assert.Empty(t, response.Error, "Should not have error")
	assert.NotEmpty(t, response.Messages, "Should have messages")
}

// TestProcessMessages_WithVariables tests placeholder variables
func TestProcessMessages_WithVariables(t *testing.T) {
	service := NewContextEngineService()

	request := ContextEngineeringRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: "Hello {{name}}",
			},
		},
		Model:    "gpt-4",
		Provider: "openai",
		Variables: map[string]interface{}{
			"name": "John",
		},
	}

	response := service.ProcessMessages(request)

	assert.Empty(t, response.Error, "Should not have error")
	assert.NotEmpty(t, response.Messages, "Should have messages")
}

// TestProcessMessages_WithHistorySummary tests history summary injection
func TestProcessMessages_WithHistorySummary(t *testing.T) {
	service := NewContextEngineService()

	request := ContextEngineeringRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: "Continue our discussion",
			},
		},
		Model:          "gpt-4",
		Provider:       "openai",
		HistorySummary: "Previous conversation was about React hooks and state management.",
	}

	response := service.ProcessMessages(request)

	assert.Empty(t, response.Error, "Should not have error")
	assert.NotEmpty(t, response.Messages, "Should have messages")
}

// TestProcessMessages_WithSessionID tests session-specific processing
func TestProcessMessages_WithSessionID(t *testing.T) {
	service := NewContextEngineService()

	request := ContextEngineeringRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
		Model:     "gpt-4",
		Provider:  "openai",
		SessionID: "session-123",
	}

	response := service.ProcessMessages(request)

	assert.Empty(t, response.Error, "Should not have error")
	assert.NotEmpty(t, response.Messages, "Should have messages")
}

// TestProcessMessages_WelcomeQuestion tests welcome question scenario
func TestProcessMessages_WelcomeQuestion(t *testing.T) {
	service := NewContextEngineService()

	request := ContextEngineeringRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: "Hello, I'm new here",
			},
		},
		Model:             "gpt-4",
		Provider:          "openai",
		SessionID:         "inbox",
		IsWelcomeQuestion: true,
	}

	response := service.ProcessMessages(request)

	assert.Empty(t, response.Error, "Should not have error")
	assert.NotEmpty(t, response.Messages, "Should have messages")
}

// TestProcessMessages_EmptyMessages tests empty message list
func TestProcessMessages_EmptyMessages(t *testing.T) {
	service := NewContextEngineService()

	request := ContextEngineeringRequest{
		Messages: []Message{},
		Model:    "gpt-4",
		Provider: "openai",
	}

	response := service.ProcessMessages(request)

	// Should handle empty messages gracefully
	assert.NotNil(t, response, "Response should not be nil")
}

// TestProcessMessages_ComplexContent tests complex message content
func TestProcessMessages_ComplexContent(t *testing.T) {
	service := NewContextEngineService()

	// Complex content with multiple parts
	complexContent := []map[string]interface{}{
		{
			"type": "text",
			"text": "What's in this image?",
		},
		{
			"type": "image_url",
			"image_url": map[string]string{
				"url": "https://example.com/image.jpg",
			},
		},
	}

	request := ContextEngineeringRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: complexContent,
			},
		},
		Model:    "gpt-4-vision",
		Provider: "openai",
	}

	response := service.ProcessMessages(request)

	assert.Empty(t, response.Error, "Should not have error")
	assert.NotEmpty(t, response.Messages, "Should have messages")
}

// TestProcessMessages_WithMetadata tests message metadata handling
func TestProcessMessages_WithMetadata(t *testing.T) {
	service := NewContextEngineService()

	request := ContextEngineeringRequest{
		Messages: []Message{
			{
				ID:        "msg-123",
				Role:      "user",
				Content:   "Hello",
				CreatedAt: 1234567890,
				UpdatedAt: 1234567890,
				Meta: map[string]interface{}{
					"source": "web",
					"lang":   "en",
				},
			},
		},
		Model:    "gpt-4",
		Provider: "openai",
	}

	response := service.ProcessMessages(request)

	assert.Empty(t, response.Error, "Should not have error")
	assert.NotEmpty(t, response.Messages, "Should have messages")
}

// TestProcessMessages_MultipleMessages tests processing multiple messages
func TestProcessMessages_MultipleMessages(t *testing.T) {
	service := NewContextEngineService()

	request := ContextEngineeringRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a helpful assistant",
			},
			{
				Role:    "user",
				Content: "What is Go?",
			},
			{
				Role:    "assistant",
				Content: "Go is a programming language.",
			},
			{
				Role:    "user",
				Content: "Tell me more",
			},
		},
		Model:    "gpt-4",
		Provider: "openai",
	}

	response := service.ProcessMessages(request)

	assert.Empty(t, response.Error, "Should not have error")
	assert.NotEmpty(t, response.Messages, "Should have messages")
}

// TestProcessMessages_AllOptions tests all options combined
func TestProcessMessages_AllOptions(t *testing.T) {
	service := NewContextEngineService()

	request := ContextEngineeringRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: "Hello {{name}}, search for {{topic}}",
			},
		},
		Model:              "gpt-4",
		Provider:           "openai",
		SystemRole:         "You are a helpful assistant",
		InputTemplate:      "Query: {{input}}",
		EnableHistoryCount: true,
		HistoryCount:       10,
		HistorySummary:     "Previous discussion about programming",
		SessionID:          "session-456",
		IsWelcomeQuestion:  false,
		Tools:              []string{"web-search", "calculator"},
		Variables: map[string]interface{}{
			"name":  "Alice",
			"topic": "TypeScript",
		},
	}

	response := service.ProcessMessages(request)

	assert.Empty(t, response.Error, "Should not have error")
	assert.NotEmpty(t, response.Messages, "Should have messages")
}

// TestProcessMessages_DifferentProviders tests different providers
func TestProcessMessages_DifferentProviders(t *testing.T) {
	service := NewContextEngineService()

	providers := []string{"openai", "anthropic", "azure", "google", "ollama"}

	for _, provider := range providers {
		t.Run(provider, func(t *testing.T) {
			request := ContextEngineeringRequest{
				Messages: []Message{
					{
						Role:    "user",
						Content: "Hello",
					},
				},
				Model:    "test-model",
				Provider: provider,
			}

			response := service.ProcessMessages(request)

			assert.NotNil(t, response, "Response should not be nil for provider: "+provider)
		})
	}
}

// TestGetEngineStats tests engine statistics retrieval
func TestGetEngineStats(t *testing.T) {
	service := NewContextEngineService()

	stats := service.GetEngineStats()

	assert.NotNil(t, stats, "Stats should not be nil")
	assert.Contains(t, stats, "processorCount", "Should contain processorCount")
	assert.Contains(t, stats, "customProcessorCount", "Should contain customProcessorCount")
	assert.Contains(t, stats, "disabledProcessorCount", "Should contain disabledProcessorCount")
	assert.Contains(t, stats, "processorNames", "Should contain processorNames")
	assert.Contains(t, stats, "customProcessorNames", "Should contain customProcessorNames")
	assert.Contains(t, stats, "disabledProcessorNames", "Should contain disabledProcessorNames")

	// Verify types
	assert.IsType(t, 0, stats["processorCount"], "processorCount should be int")
	assert.IsType(t, []string{}, stats["processorNames"], "processorNames should be []string")
}

// TestGetEngineStats_ProcessorCount tests processor count
func TestGetEngineStats_ProcessorCount(t *testing.T) {
	service := NewContextEngineService()

	stats := service.GetEngineStats()

	processorCount := stats["processorCount"].(int)
	assert.Greater(t, processorCount, 0, "Should have at least one processor")

	processorNames := stats["processorNames"].([]string)
	assert.Equal(t, processorCount, len(processorNames), "Processor count should match names length")
}

// TestValidateConfig_ValidConfig tests valid configuration
func TestValidateConfig_ValidConfig(t *testing.T) {
	service := NewContextEngineService()

	// Create config without function fields (they can't be JSON marshaled)
	configMap := map[string]interface{}{
		"systemRole":         "You are helpful",
		"enableHistoryCount": true,
		"historyCount":       10,
		"model":              "gpt-4",
		"provider":           "openai",
	}

	configJSON, err := json.Marshal(configMap)
	require.NoError(t, err, "Should marshal config")

	result := service.ValidateConfig(string(configJSON))

	assert.NotNil(t, result, "Result should not be nil")
	assert.Contains(t, result, "valid", "Should contain valid field")
	assert.Contains(t, result, "errors", "Should contain errors field")
	assert.True(t, result["valid"].(bool), "Config should be valid")
	assert.Empty(t, result["errors"].([]string), "Should have no errors")
}

// TestValidateConfig_InvalidJSON tests invalid JSON
func TestValidateConfig_InvalidJSON(t *testing.T) {
	service := NewContextEngineService()

	result := service.ValidateConfig("invalid json {{{")

	assert.NotNil(t, result, "Result should not be nil")
	assert.False(t, result["valid"].(bool), "Should be invalid")
	assert.NotEmpty(t, result["errors"], "Should have errors")
}

// TestValidateConfig_NegativeHistoryCount tests negative history count
func TestValidateConfig_NegativeHistoryCount(t *testing.T) {
	service := NewContextEngineService()

	configMap := map[string]interface{}{
		"enableHistoryCount": true,
		"historyCount":       -5, // Invalid
	}

	configJSON, err := json.Marshal(configMap)
	require.NoError(t, err, "Should marshal config")

	result := service.ValidateConfig(string(configJSON))

	assert.NotNil(t, result, "Result should not be nil")
	assert.False(t, result["valid"].(bool), "Should be invalid")
	assert.NotEmpty(t, result["errors"], "Should have errors")
}

// TestValidateConfig_EmptyConfig tests empty configuration
func TestValidateConfig_EmptyConfig(t *testing.T) {
	service := NewContextEngineService()

	configMap := map[string]interface{}{}
	configJSON, err := json.Marshal(configMap)
	require.NoError(t, err, "Should marshal config")

	result := service.ValidateConfig(string(configJSON))

	assert.NotNil(t, result, "Result should not be nil")
	assert.Contains(t, result, "valid", "Should contain valid field")
}

// TestValidateConfig_WithTools tests configuration with tools
func TestValidateConfig_WithTools(t *testing.T) {
	service := NewContextEngineService()

	configMap := map[string]interface{}{
		"tools": []map[string]interface{}{
			{
				"id":          "tool-1",
				"name":        "Web Search",
				"description": "Search the web",
			},
			{
				"id":          "tool-2",
				"name":        "Calculator",
				"description": "Perform calculations",
			},
		},
	}

	configJSON, err := json.Marshal(configMap)
	require.NoError(t, err, "Should marshal config")

	result := service.ValidateConfig(string(configJSON))

	assert.NotNil(t, result, "Result should not be nil")
	assert.True(t, result["valid"].(bool), "Config should be valid")
}

// TestProcessMessages_Idempotency tests idempotency
func TestProcessMessages_Idempotency(t *testing.T) {
	service := NewContextEngineService()

	request := ContextEngineeringRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
		Model:    "gpt-4",
		Provider: "openai",
	}

	// Process same request twice
	response1 := service.ProcessMessages(request)
	response2 := service.ProcessMessages(request)

	// Should produce consistent results
	assert.Empty(t, response1.Error, "First call should not have error")
	assert.Empty(t, response2.Error, "Second call should not have error")
	assert.Equal(t, len(response1.Messages), len(response2.Messages), "Should have same message count")
}

// TestProcessMessages_ConcurrentCalls tests concurrent processing
func TestProcessMessages_ConcurrentCalls(t *testing.T) {
	service := NewContextEngineService()

	request := ContextEngineeringRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: "Hello",
			},
		},
		Model:    "gpt-4",
		Provider: "openai",
	}

	// Process concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			response := service.ProcessMessages(request)
			assert.Empty(t, response.Error, "Should not have error")
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 10; i++ {
		<-done
	}
}

// TestProcessMessages_LargeMessageCount tests processing many messages
func TestProcessMessages_LargeMessageCount(t *testing.T) {
	service := NewContextEngineService()

	// Create 100 messages
	messages := make([]Message, 100)
	for i := 0; i < 100; i++ {
		messages[i] = Message{
			Role:    "user",
			Content: "Message number " + string(rune(i)),
		}
	}

	request := ContextEngineeringRequest{
		Messages: messages,
		Model:    "gpt-4",
		Provider: "openai",
	}

	response := service.ProcessMessages(request)

	assert.Empty(t, response.Error, "Should handle large message count")
	assert.NotEmpty(t, response.Messages, "Should have messages")
}

// TestProcessMessages_SpecialCharacters tests special characters handling
func TestProcessMessages_SpecialCharacters(t *testing.T) {
	service := NewContextEngineService()

	specialChars := []string{
		"Hello 世界",                  // Unicode
		"Test\nNew\nLine",           // Newlines
		"Tab\tSeparated",            // Tabs
		"Quote\"Test\"",             // Quotes
		"Emoji 🚀 🎉 ✨",               // Emojis
		"<html>Tags</html>",         // HTML
		"JSON {\"key\": \"value\"}", // JSON
		"Math: ∑∫∂∇",                // Math symbols
	}

	for _, content := range specialChars {
		t.Run(content, func(t *testing.T) {
			request := ContextEngineeringRequest{
				Messages: []Message{
					{
						Role:    "user",
						Content: content,
					},
				},
				Model:    "gpt-4",
				Provider: "openai",
			}

			response := service.ProcessMessages(request)
			assert.Empty(t, response.Error, "Should handle special characters: "+content)
		})
	}
}

// TestProcessMessages_DifferentRoles tests different message roles
func TestProcessMessages_DifferentRoles(t *testing.T) {
	service := NewContextEngineService()

	roles := []string{"user", "assistant", "system", "tool"}

	for _, role := range roles {
		t.Run(role, func(t *testing.T) {
			request := ContextEngineeringRequest{
				Messages: []Message{
					{
						Role:    role,
						Content: "Test message",
					},
				},
				Model:    "gpt-4",
				Provider: "openai",
			}

			response := service.ProcessMessages(request)
			assert.NotNil(t, response, "Should handle role: "+role)
		})
	}
}

// BenchmarkProcessMessages benchmarks message processing
func BenchmarkProcessMessages(b *testing.B) {
	service := NewContextEngineService()

	request := ContextEngineeringRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: "Hello, world!",
			},
		},
		Model:    "gpt-4",
		Provider: "openai",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ProcessMessages(request)
	}
}

// BenchmarkProcessMessages_WithOptions benchmarks with all options
func BenchmarkProcessMessages_WithOptions(b *testing.B) {
	service := NewContextEngineService()

	request := ContextEngineeringRequest{
		Messages: []Message{
			{
				Role:    "user",
				Content: "Hello {{name}}",
			},
		},
		Model:              "gpt-4",
		Provider:           "openai",
		SystemRole:         "You are helpful",
		InputTemplate:      "Query: {{input}}",
		EnableHistoryCount: true,
		HistoryCount:       10,
		Tools:              []string{"web-search"},
		Variables: map[string]interface{}{
			"name": "User",
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ProcessMessages(request)
	}
}

// BenchmarkGetEngineStats benchmarks stats retrieval
func BenchmarkGetEngineStats(b *testing.B) {
	service := NewContextEngineService()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.GetEngineStats()
	}
}

// BenchmarkValidateConfig benchmarks config validation
func BenchmarkValidateConfig(b *testing.B) {
	service := NewContextEngineService()

	config := Config{
		SystemRole:   "You are helpful",
		HistoryCount: 10,
	}
	configJSON, _ := json.Marshal(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		service.ValidateConfig(string(configJSON))
	}
}
