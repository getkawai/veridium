package gateway

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kawai-network/veridium/pkg/fantasy/llamalib"
)

// LlamaExecutor implements LLMExecutor using llamalib service.
type LlamaExecutor struct {
	service   *llamalib.Service
	maxTokens int32
}

// NewLlamaExecutor creates a new executor backed by llamalib.
// It initializes the llamalib service and loads the chat model.
func NewLlamaExecutor(ctx context.Context, modelPath string, maxTokens int32) (*LlamaExecutor, error) {
	service := llamalib.NewService()

	// Wait for library initialization
	initCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	if err := service.WaitForInitialization(initCtx); err != nil {
		return nil, fmt.Errorf("failed to initialize llamalib: %w", err)
	}

	// Load chat model (empty string = auto-select best model)
	if err := service.LoadChatModel(modelPath); err != nil {
		return nil, fmt.Errorf("failed to load chat model: %w", err)
	}

	return &LlamaExecutor{
		service:   service,
		maxTokens: maxTokens,
	}, nil
}

// Execute runs inference and returns the response content.
func (e *LlamaExecutor) Execute(messages []ChatMessage) (string, error) {
	// Convert chat messages to a single prompt
	prompt := e.formatMessages(messages)

	result, err := e.service.Generate(prompt, e.maxTokens)
	if err != nil {
		return "", fmt.Errorf("generation failed: %w", err)
	}

	return strings.TrimSpace(result), nil
}

// ExecuteStream runs inference with streaming response.
// Note: llamalib doesn't support streaming natively, so we simulate it.
func (e *LlamaExecutor) ExecuteStream(messages []ChatMessage, stream chan<- string) error {
	// Generate full response first (llamalib doesn't have native streaming)
	result, err := e.Execute(messages)
	if err != nil {
		return err
	}

	// Simulate streaming by sending chunks
	words := strings.Fields(result)
	for i, word := range words {
		if i > 0 {
			stream <- " "
		}
		stream <- word
	}

	return nil
}

// formatMessages converts chat messages to a single prompt string.
func (e *LlamaExecutor) formatMessages(messages []ChatMessage) string {
	var sb strings.Builder
	for _, msg := range messages {
		switch msg.Role {
		case "system":
			sb.WriteString(fmt.Sprintf("System: %s\n", msg.Content))
		case "user":
			sb.WriteString(fmt.Sprintf("User: %s\n", msg.Content))
		case "assistant":
			sb.WriteString(fmt.Sprintf("Assistant: %s\n", msg.Content))
		default:
			sb.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
		}
	}
	sb.WriteString("Assistant:")
	return sb.String()
}

// Service returns the underlying llamalib service for advanced usage.
func (e *LlamaExecutor) Service() *llamalib.Service {
	return e.service
}
