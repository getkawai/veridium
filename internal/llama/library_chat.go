package llama

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/kawai-network/veridium/pkg/yzma/llama"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// LibraryChatService provides chat functionality using llama library
type LibraryChatService struct {
	libService *LibraryService
	app        *application.App
}

// NewLibraryChatService creates a new library-based chat service
func NewLibraryChatService(libService *LibraryService, app *application.App) *LibraryChatService {
	return &LibraryChatService{
		libService: libService,
		app:        app,
	}
}

// ChatCompletionRequest represents an OpenAI-compatible chat completion request
type ChatCompletionRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	MaxTokens   int32         `json:"max_tokens,omitempty"`
	Temperature float32       `json:"temperature,omitempty"`
	TopP        float32       `json:"top_p,omitempty"`
	TopK        int32         `json:"top_k,omitempty"`
	Stream      bool          `json:"stream,omitempty"`
}

// ChatMessage represents a chat message
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionResponse represents the response
type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []ChatCompletionChoice `json:"choices"`
	Usage   *ChatUsage             `json:"usage,omitempty"`
}

// ChatCompletionChoice represents a single choice in the response
type ChatCompletionChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finish_reason"`
}

// ChatCompletionChunk represents a streaming chunk
type ChatCompletionChunk struct {
	ID      string                      `json:"id"`
	Object  string                      `json:"object"`
	Created int64                       `json:"created"`
	Model   string                      `json:"model"`
	Choices []ChatCompletionChunkChoice `json:"choices"`
}

// ChatCompletionChunkChoice represents a streaming choice
type ChatCompletionChunkChoice struct {
	Index        int              `json:"index"`
	Delta        ChatMessageDelta `json:"delta"`
	FinishReason string           `json:"finish_reason,omitempty"`
}

// ChatMessageDelta represents delta content in streaming
type ChatMessageDelta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

// ChatUsage represents token usage
type ChatUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletion handles a chat completion request (non-streaming)
func (c *LibraryChatService) ChatCompletion(ctx context.Context, req ChatCompletionRequest) (*ChatCompletionResponse, error) {
	// Validate request
	if err := validateChatRequest(req); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	// Ensure chat model is loaded
	if !c.libService.IsChatModelLoaded() {
		if err := c.libService.LoadChatModel(""); err != nil {
			return nil, fmt.Errorf("failed to load chat model: %w", err)
		}
	}

	// Build prompt from messages
	prompt := c.buildPrompt(req.Messages)

	// Set default parameters
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 512
	}

	// Update sampler if needed
	c.updateSampler(req)

	// Generate response
	response, err := c.libService.Generate(prompt, maxTokens)
	if err != nil {
		return nil, fmt.Errorf("generation failed: %w", err)
	}

	// Build response
	return &ChatCompletionResponse{
		ID:      fmt.Sprintf("chatcmpl-%d", time.Now().Unix()),
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   c.libService.GetLoadedChatModel(),
		Choices: []ChatCompletionChoice{
			{
				Index: 0,
				Message: ChatMessage{
					Role:    "assistant",
					Content: response,
				},
				FinishReason: "stop",
			},
		},
		Usage: &ChatUsage{
			PromptTokens:     len(prompt) / 4, // Rough estimate
			CompletionTokens: len(response) / 4,
			TotalTokens:      (len(prompt) + len(response)) / 4,
		},
	}, nil
}

// ChatCompletionStream handles a streaming chat completion request
func (c *LibraryChatService) ChatCompletionStream(ctx context.Context, requestID string, req ChatCompletionRequest) error {
	// Validate request
	if err := validateChatRequest(req); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	// Ensure chat model is loaded
	if !c.libService.IsChatModelLoaded() {
		if err := c.libService.LoadChatModel(""); err != nil {
			return fmt.Errorf("failed to load chat model: %w", err)
		}
	}

	// Build prompt from messages
	prompt := c.buildPrompt(req.Messages)

	// Set default parameters
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 512
	}

	// Update sampler
	c.updateSampler(req)

	// Generate response with streaming
	err := c.generateStreaming(ctx, requestID, prompt, maxTokens)
	if err != nil {
		return fmt.Errorf("streaming generation failed: %w", err)
	}

	return nil
}

// generateStreaming generates text and streams it via Wails events
func (c *LibraryChatService) generateStreaming(ctx context.Context, requestID string, prompt string, maxTokens int32) error {
	c.libService.chatMutex.Lock()
	defer c.libService.chatMutex.Unlock()

	if c.libService.chatModel == 0 || c.libService.chatContext == 0 {
		return fmt.Errorf("chat model not loaded")
	}

	modelPath := c.libService.chatModelPath
	chatID := fmt.Sprintf("chatcmpl-%d", time.Now().Unix())

	// Send initial metadata
	c.app.Event.Emit(fmt.Sprintf("stream:%s:meta", requestID), map[string]interface{}{
		"status":     200,
		"statusText": "OK",
		"headers": map[string]string{
			"Content-Type": "text/event-stream",
		},
	})

	// Tokenize prompt (add BOS for proper prompt processing)
	tokens := llama.Tokenize(c.libService.chatVocab, prompt, true, true)
	if len(tokens) == 0 {
		return fmt.Errorf("failed to tokenize prompt")
	}

	// Reset sampler state before new generation
	llama.SamplerReset(c.libService.chatSampler)

	// Decode prompt tokens to initialize context
	batch := llama.BatchGetOne(tokens)
	if llama.Decode(c.libService.chatContext, batch) != 0 {
		return fmt.Errorf("failed to decode prompt")
	}

	// Send first chunk with role
	firstChunk := ChatCompletionChunk{
		ID:      chatID,
		Object:  "chat.completion.chunk",
		Created: time.Now().Unix(),
		Model:   modelPath,
		Choices: []ChatCompletionChunkChoice{
			{
				Index: 0,
				Delta: ChatMessageDelta{
					Role: "assistant",
				},
			},
		},
	}
	c.emitSSEChunk(requestID, firstChunk)

	// Generate tokens and stream
	var utf8Buffer []byte // Buffer for incomplete UTF-8 sequences across tokens

	for nGenerated := int32(0); nGenerated < maxTokens; nGenerated++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Sample next token
		token := llama.SamplerSample(c.libService.chatSampler, c.libService.chatContext, -1)

		// Check for end of generation
		if llama.VocabIsEOG(c.libService.chatVocab, token) {
			// Send final chunk with finish_reason
			finalChunk := ChatCompletionChunk{
				ID:      chatID,
				Object:  "chat.completion.chunk",
				Created: time.Now().Unix(),
				Model:   modelPath,
				Choices: []ChatCompletionChunkChoice{
					{
						Index:        0,
						Delta:        ChatMessageDelta{},
						FinishReason: "stop",
					},
				},
			}
			c.emitSSEChunk(requestID, finalChunk)
			break
		}

		// Convert token to raw bytes
		buf := make([]byte, 256)
		length := llama.TokenToPiece(c.libService.chatVocab, token, buf, 0, false)
		tokenBytes := buf[:length]

		// Append to UTF-8 buffer (accumulate bytes from multiple tokens)
		// This is crucial for handling multibyte UTF-8 chars (emoji, etc) that span multiple tokens
		utf8Buffer = append(utf8Buffer, tokenBytes...)

		// Try to decode valid UTF-8 string from buffer
		content := ""
		validBytes := 0

		// Find the longest valid UTF-8 prefix in buffer
		for i := len(utf8Buffer); i > 0; i-- {
			if utf8.Valid(utf8Buffer[:i]) {
				// Found valid UTF-8 sequence up to position i
				content = string(utf8Buffer[:i])
				validBytes = i
				break
			}
		}

		// If we have valid content, emit it and remove from buffer
		if validBytes > 0 {
			// Keep incomplete bytes in buffer for next iteration
			utf8Buffer = utf8Buffer[validBytes:]
		} else {
			// No valid UTF-8 yet - incomplete multibyte sequence
			// Wait for more tokens to complete the sequence

			// Accept the token and prepare for next generation
			llama.SamplerAccept(c.libService.chatSampler, token)

			// Decode the new token to update context
			nextBatch := llama.BatchGetOne([]llama.Token{token})
			if llama.Decode(c.libService.chatContext, nextBatch) != 0 {
				return fmt.Errorf("failed to decode token")
			}

			// Don't send chunk yet, wait for complete UTF-8 sequence
			continue
		}

		// Send chunk with content (only if valid UTF-8)
		chunk := ChatCompletionChunk{
			ID:      chatID,
			Object:  "chat.completion.chunk",
			Created: time.Now().Unix(),
			Model:   modelPath,
			Choices: []ChatCompletionChunkChoice{
				{
					Index: 0,
					Delta: ChatMessageDelta{
						Content: content,
					},
				},
			},
		}
		c.emitSSEChunk(requestID, chunk)

		// Accept the token and prepare for next generation
		llama.SamplerAccept(c.libService.chatSampler, token)

		// Decode the new token to update context
		nextBatch := llama.BatchGetOne([]llama.Token{token})
		if llama.Decode(c.libService.chatContext, nextBatch) != 0 {
			return fmt.Errorf("failed to decode token")
		}

		// Small delay to prevent overwhelming the frontend
		time.Sleep(10 * time.Millisecond)
	}

	// Send done signal
	c.app.Event.Emit(fmt.Sprintf("stream:%s:data", requestID), "data: [DONE]\n\n")
	c.app.Event.Emit(fmt.Sprintf("stream:%s:end", requestID))

	return nil
}

// emitSSEChunk formats and emits a chunk as Server-Sent Event
func (c *LibraryChatService) emitSSEChunk(requestID string, chunk ChatCompletionChunk) {
	data, err := json.Marshal(chunk)
	if err != nil {
		log.Printf("Failed to marshal chunk: %v", err)
		return
	}

	// Format as SSE
	sseData := fmt.Sprintf("data: %s\n\n", string(data))
	c.app.Event.Emit(fmt.Sprintf("stream:%s:data", requestID), sseData)
}

// buildPrompt builds a prompt from chat messages
func (c *LibraryChatService) buildPrompt(messages []ChatMessage) string {
	// Try to get the chat template from the model
	var template string
	if c.libService.chatModel != 0 {
		template = llama.ModelChatTemplate(c.libService.chatModel, "")
	}

	// If no template, use a simple format
	if template == "" {
		var prompt strings.Builder
		for _, msg := range messages {
			switch msg.Role {
			case "system":
				prompt.WriteString(fmt.Sprintf("System: %s\n\n", msg.Content))
			case "user":
				prompt.WriteString(fmt.Sprintf("User: %s\n\n", msg.Content))
			case "assistant":
				prompt.WriteString(fmt.Sprintf("Assistant: %s\n\n", msg.Content))
			}
		}
		prompt.WriteString("Assistant:")
		return prompt.String()
	}

	// Use llama.cpp chat template
	llamaMessages := make([]llama.ChatMessage, len(messages))
	for i, msg := range messages {
		llamaMessages[i] = llama.NewChatMessage(msg.Role, msg.Content)
	}

	buf := make([]byte, 8192)
	length := llama.ChatApplyTemplate(template, llamaMessages, true, buf)
	return string(buf[:length])
}

// updateSampler updates the sampler parameters based on request
// This recreates the sampler with custom parameters from the request
func (c *LibraryChatService) updateSampler(req ChatCompletionRequest) {
	c.libService.chatMutex.Lock()
	defer c.libService.chatMutex.Unlock()

	// Only recreate if custom parameters are provided
	needsUpdate := req.Temperature > 0 || req.TopP > 0 || req.TopK > 0

	if !needsUpdate || c.libService.chatModel == 0 || c.libService.chatVocab == 0 {
		return
	}

	// Free existing sampler
	if c.libService.chatSampler != 0 {
		llama.SamplerFree(c.libService.chatSampler)
	}

	// Create new sampler chain with custom parameters
	c.libService.chatSampler = c.createCustomSampler(req)
}

// createCustomSampler creates a sampler chain with custom parameters
func (c *LibraryChatService) createCustomSampler(req ChatCompletionRequest) llama.Sampler {
	// Initialize sampler chain
	params := llama.SamplerChainDefaultParams()
	sampler := llama.SamplerChainInit(params)

	// Add penalties (always included for quality)
	penalties := llama.SamplerInitPenalties(
		64,  // penalty_last_n: last 64 tokens
		1.0, // penalty_repeat: 1.0 = disabled
		0.0, // penalty_freq: 0.0 = disabled
		0.0, // penalty_present: 0.0 = disabled
	)
	llama.SamplerChainAdd(sampler, penalties)

	// Add Top-K if specified
	if req.TopK > 0 {
		topK := llama.SamplerInitTopK(req.TopK)
		llama.SamplerChainAdd(sampler, topK)
	} else {
		// Default Top-K
		topK := llama.SamplerInitTopK(40)
		llama.SamplerChainAdd(sampler, topK)
	}

	// Add Top-P if specified
	if req.TopP > 0 {
		topP := llama.SamplerInitTopP(req.TopP, 0)
		llama.SamplerChainAdd(sampler, topP)
	} else {
		// Default Top-P
		topP := llama.SamplerInitTopP(0.95, 0)
		llama.SamplerChainAdd(sampler, topP)
	}

	// Add Min-P (always included for quality)
	minP := llama.SamplerInitMinP(0.05, 0)
	llama.SamplerChainAdd(sampler, minP)

	// Add Temperature if specified
	if req.Temperature > 0 {
		temp := llama.SamplerInitTempExt(req.Temperature, 0, 1.0)
		llama.SamplerChainAdd(sampler, temp)
	} else {
		// Default temperature
		temp := llama.SamplerInitTempExt(0.8, 0, 1.0)
		llama.SamplerChainAdd(sampler, temp)
	}

	// Always add distribution sampler last
	dist := llama.SamplerInitDist(llama.DefaultSeed)
	llama.SamplerChainAdd(sampler, dist)

	return sampler
}

// validateChatRequest validates a chat completion request
func validateChatRequest(req ChatCompletionRequest) error {
	// Validate messages
	if len(req.Messages) == 0 {
		return fmt.Errorf("messages cannot be empty")
	}

	// Validate each message
	for i, msg := range req.Messages {
		// Check role
		if msg.Role != "system" && msg.Role != "user" && msg.Role != "assistant" {
			return fmt.Errorf("invalid role '%s' in message %d: must be 'system', 'user', or 'assistant'", msg.Role, i)
		}

		// Check content
		if strings.TrimSpace(msg.Content) == "" {
			return fmt.Errorf("message content cannot be empty at index %d", i)
		}
	}

	// Validate max_tokens
	if req.MaxTokens < 0 {
		return fmt.Errorf("max_tokens must be positive, got: %d", req.MaxTokens)
	}

	// Validate temperature (typically 0-2)
	if req.Temperature < 0 || req.Temperature > 2 {
		return fmt.Errorf("temperature must be between 0 and 2, got: %f", req.Temperature)
	}

	// Validate top_p (0-1)
	if req.TopP < 0 || req.TopP > 1 {
		return fmt.Errorf("top_p must be between 0 and 1, got: %f", req.TopP)
	}

	// Validate top_k
	if req.TopK < 0 {
		return fmt.Errorf("top_k must be non-negative, got: %d", req.TopK)
	}

	return nil
}

// Note: LibraryProxyService has been removed as it's not used in main.go
// LibraryChatService can be used directly for chat functionality
// Example usage:
//   libService, _ := llama.NewLibraryService()
//   chatService := llama.NewLibraryChatService(libService, app)
//   app.Bind(chatService)
