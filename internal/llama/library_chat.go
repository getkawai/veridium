package llama

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/hybridgroup/yzma/pkg/llama"
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
	ID      string                    `json:"id"`
	Object  string                    `json:"object"`
	Created int64                     `json:"created"`
	Model   string                    `json:"model"`
	Choices []ChatCompletionChunkChoice `json:"choices"`
}

// ChatCompletionChunkChoice represents a streaming choice
type ChatCompletionChunkChoice struct {
	Index        int               `json:"index"`
	Delta        ChatMessageDelta  `json:"delta"`
	FinishReason string            `json:"finish_reason,omitempty"`
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
	
	// Tokenize prompt
	tokens := llama.Tokenize(c.libService.chatVocab, prompt, true, true)
	if len(tokens) == 0 {
		return fmt.Errorf("failed to tokenize prompt")
	}
	
	// Create batch
	batch := llama.BatchGetOne(tokens)
	
	// Handle encoder models
	if llama.ModelHasEncoder(c.libService.chatModel) {
		llama.Encode(c.libService.chatContext, batch)
		start := llama.ModelDecoderStartToken(c.libService.chatModel)
		if start == llama.TokenNull {
			start := llama.VocabBOS(c.libService.chatVocab)
		batch = llama.BatchGetOne([]llama.Token{start})
	}
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
	for pos := int32(0); pos < maxTokens; pos += batch.NTokens {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		
		llama.Decode(c.libService.chatContext, batch)
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
		
		// Convert token to text
		buf := make([]byte, 256)
		length := llama.TokenToPiece(c.libService.chatVocab, token, buf, 0, false)
		content := string(buf[:length])
		
		// Send chunk with content
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
		
		// Prepare next batch
		batch = llama.BatchGetOne([]llama.Token{token})
		
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
func (c *LibraryChatService) updateSampler(req ChatCompletionRequest) {
	c.libService.chatMutex.Lock()
	defer c.libService.chatMutex.Unlock()
	
	if c.libService.chatSampler == 0 {
		return
	}
	
	// Recreate sampler with new parameters
	// Note: In a production system, you might want to cache samplers
	// or have a more sophisticated parameter update mechanism
	
	// For now, we'll keep the existing sampler
	// In a full implementation, you'd recreate it with custom parameters
	// based on req.Temperature, req.TopP, req.TopK
}

// Proxy provides a compatibility layer for the existing proxy service
// This allows gradual migration from binary to library approach
type LibraryProxyService struct {
	libService *LibraryService
	chatService *LibraryChatService
	app        *application.App
}

// NewLibraryProxyService creates a new library-based proxy service
func NewLibraryProxyService(libService *LibraryService, app *application.App) *LibraryProxyService {
	return &LibraryProxyService{
		libService:  libService,
		chatService: NewLibraryChatService(libService, app),
		app:         app,
	}
}

// Fetch handles non-streaming requests (compatible with existing Fetch interface)
func (p *LibraryProxyService) Fetch(ctx context.Context, request ProxyRequest) (*ProxyResponse, error) {
	log.Printf("📍 [LibraryProxyService.Fetch] method=%s, path=%s", request.Method, request.Path)
	
	// Parse request based on path
	if strings.HasSuffix(request.Path, "/chat/completions") || strings.HasSuffix(request.Path, "/v1/chat/completions") {
		// Chat completion request
		var req ChatCompletionRequest
		if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
			return nil, fmt.Errorf("failed to parse request: %w", err)
		}
		
		// Force non-streaming for Fetch
		req.Stream = false
		
		// Handle the request
		resp, err := p.chatService.ChatCompletion(ctx, req)
		if err != nil {
			return nil, err
		}
		
		// Marshal response
		respData, err := json.Marshal(resp)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal response: %w", err)
		}
		
		return &ProxyResponse{
			Status:     200,
			StatusText: "OK",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: string(respData),
		}, nil
	}
	
	// For other endpoints, return error or forward to original service
	return nil, fmt.Errorf("endpoint not supported by library service: %s", request.Path)
}

// StreamFetch handles streaming requests (compatible with existing StreamFetch interface)
func (p *LibraryProxyService) StreamFetch(ctx context.Context, requestID string, request ProxyRequest) error {
	log.Printf("[LibraryProxyService.StreamFetch] Starting stream for request ID: %s", requestID)
	
	// Parse request based on path
	if strings.HasSuffix(request.Path, "/chat/completions") || strings.HasSuffix(request.Path, "/v1/chat/completions") {
		// Chat completion request
		var req ChatCompletionRequest
		if err := json.Unmarshal([]byte(request.Body), &req); err != nil {
			return fmt.Errorf("failed to parse request: %w", err)
		}
		
		// Handle streaming
		return p.chatService.ChatCompletionStream(ctx, requestID, req)
	}
	
	return fmt.Errorf("endpoint not supported by library service: %s", request.Path)
}

