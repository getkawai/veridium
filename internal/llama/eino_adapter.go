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

package llama

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/schema"
	"github.com/kawai-network/veridium/pkg/yzma/llama"
)

// LlamaEinoModel implements Eino's ToolCallingChatModel interface
// It wraps LibraryService to provide Eino-compatible chat model functionality
type LlamaEinoModel struct {
	libService *LibraryService
	tools      []*schema.ToolInfo
}

// NewLlamaEinoModel creates a new Eino model adapter wrapping LibraryService
func NewLlamaEinoModel(libService *LibraryService) *LlamaEinoModel {
	return &LlamaEinoModel{
		libService: libService,
		tools:      nil,
	}
}

// Generate implements model.BaseChatModel.Generate
// It converts Eino messages to llama format, generates response, and converts back
func (m *LlamaEinoModel) Generate(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.Message, error) {
	// Ensure chat model is loaded
	if !m.libService.IsChatModelLoaded() {
		if err := m.libService.LoadChatModel(""); err != nil {
			return nil, fmt.Errorf("failed to load chat model: %w", err)
		}
	}

	// Convert Eino messages to llama format
	prompt, err := m.convertMessagesToPrompt(input)
	if err != nil {
		return nil, fmt.Errorf("failed to convert messages: %w", err)
	}

	// Parse options
	options := model.GetCommonOptions(nil, opts...)

	// Set default max tokens (practically unlimited - 32768 tokens ~24k words)
	maxTokens := int32(32768)
	if options.MaxTokens != nil && *options.MaxTokens > 0 {
		maxTokens = int32(*options.MaxTokens)
	}

	// Update sampler with custom parameters
	m.updateSamplerFromConfig(options)

	// Generate response using LibraryService
	response, err := m.libService.Generate(prompt, maxTokens)
	if err != nil {
		return nil, fmt.Errorf("generation failed: %w", err)
	}

	// Convert response to Eino message
	responseMsg := &schema.Message{
		Role:    schema.Assistant,
		Content: response,
		ResponseMeta: &schema.ResponseMeta{
			FinishReason: "stop",
			Usage: &schema.TokenUsage{
				PromptTokens:     len(prompt) / 4, // Rough estimate
				CompletionTokens: len(response) / 4,
				TotalTokens:      (len(prompt) + len(response)) / 4,
			},
		},
	}

	return responseMsg, nil
}

// Stream implements model.BaseChatModel.Stream
// It generates streaming responses compatible with Eino's StreamReader
func (m *LlamaEinoModel) Stream(ctx context.Context, input []*schema.Message, opts ...model.Option) (*schema.StreamReader[*schema.Message], error) {
	// Ensure chat model is loaded
	if !m.libService.IsChatModelLoaded() {
		if err := m.libService.LoadChatModel(""); err != nil {
			return nil, fmt.Errorf("failed to load chat model: %w", err)
		}
	}

	// Convert Eino messages to llama format
	prompt, err := m.convertMessagesToPrompt(input)
	if err != nil {
		return nil, fmt.Errorf("failed to convert messages: %w", err)
	}

	// Parse options
	options := model.GetCommonOptions(nil, opts...)

	// Set default max tokens (practically unlimited - 32768 tokens ~24k words)
	maxTokens := int32(32768)
	if options.MaxTokens != nil && *options.MaxTokens > 0 {
		maxTokens = int32(*options.MaxTokens)
	}

	// Update sampler with custom parameters
	m.updateSamplerFromConfig(options)

	// Create stream reader
	reader, writer := schema.Pipe[*schema.Message](1)

	// Start generation in background
	go func() {
		defer writer.Close()

		err := m.generateStreaming(ctx, writer, prompt, maxTokens)
		if err != nil {
			writer.Send(&schema.Message{}, err)
		}
	}()

	return reader, nil
}

// WithTools implements model.ToolCallingChatModel.WithTools
// It returns a new instance with tools bound for tool calling
func (m *LlamaEinoModel) WithTools(tools []*schema.ToolInfo) (model.ToolCallingChatModel, error) {
	// Create new instance with tools
	newModel := &LlamaEinoModel{
		libService: m.libService,
		tools:      tools,
	}

	return newModel, nil
}

// convertMessagesToPrompt converts Eino messages to llama.cpp prompt format
func (m *LlamaEinoModel) convertMessagesToPrompt(messages []*schema.Message) (string, error) {
	m.libService.chatMutex.Lock()
	defer m.libService.chatMutex.Unlock()

	if m.libService.chatModel == 0 {
		return "", fmt.Errorf("chat model not loaded")
	}

	// Get chat template from model
	template := llama.ModelChatTemplate(m.libService.chatModel, "")

	// If no template, use simple format
	if template == "" {
		return m.buildSimplePrompt(messages), nil
	}

	// Convert to llama.ChatMessage format
	llamaMessages := make([]llama.ChatMessage, 0, len(messages))
	for _, msg := range messages {
		role := string(msg.Role)
		content := msg.Content

		// Handle tool calls in assistant messages
		if msg.Role == schema.Assistant && len(msg.ToolCalls) > 0 {
			// Append tool calls as formatted text
			var toolCallsStr strings.Builder
			toolCallsStr.WriteString(content)
			if content != "" {
				toolCallsStr.WriteString("\n")
			}
			for _, tc := range msg.ToolCalls {
				toolCallsStr.WriteString(fmt.Sprintf("\nTool Call: %s(%s)", tc.Function.Name, tc.Function.Arguments))
			}
			content = toolCallsStr.String()
		}

		// Handle tool messages
		if msg.Role == schema.Tool {
			role = "tool"
			// Include tool call ID and name in content
			content = fmt.Sprintf("Tool Result [%s/%s]: %s", msg.ToolCallID, msg.ToolName, msg.Content)
		}

		llamaMessages = append(llamaMessages, llama.NewChatMessage(role, content))
	}

	// Apply chat template
	buf := make([]byte, 16384) // Larger buffer for complex conversations
	length := llama.ChatApplyTemplate(template, llamaMessages, true, buf)

	return string(buf[:length]), nil
}

// buildSimplePrompt builds a simple prompt without template
func (m *LlamaEinoModel) buildSimplePrompt(messages []*schema.Message) string {
	var prompt strings.Builder

	for _, msg := range messages {
		switch msg.Role {
		case schema.System:
			prompt.WriteString(fmt.Sprintf("System: %s\n\n", msg.Content))
		case schema.User:
			prompt.WriteString(fmt.Sprintf("User: %s\n\n", msg.Content))
		case schema.Assistant:
			prompt.WriteString(fmt.Sprintf("Assistant: %s\n\n", msg.Content))
		case schema.Tool:
			prompt.WriteString(fmt.Sprintf("Tool [%s]: %s\n\n", msg.ToolName, msg.Content))
		}
	}

	prompt.WriteString("Assistant:")
	return prompt.String()
}

// updateSamplerFromConfig updates sampler parameters from Eino options
func (m *LlamaEinoModel) updateSamplerFromConfig(options *model.Options) {
	m.libService.chatMutex.Lock()
	defer m.libService.chatMutex.Unlock()

	if m.libService.chatModel == 0 || m.libService.chatVocab == 0 {
		return
	}

	// Check if we need to update sampler
	needsUpdate := (options.Temperature != nil && *options.Temperature > 0) ||
		(options.TopP != nil && *options.TopP > 0)

	if !needsUpdate {
		return
	}

	// Free existing sampler
	if m.libService.chatSampler != 0 {
		llama.SamplerFree(m.libService.chatSampler)
	}

	// Create new sampler chain with custom parameters
	samplerParams := llama.SamplerChainDefaultParams()
	sampler := llama.SamplerChainInit(samplerParams)

	// Add penalties
	penalties := llama.SamplerInitPenalties(64, 1.0, 0.0, 0.0)
	llama.SamplerChainAdd(sampler, penalties)

	// Add Top-K
	topK := int32(40)
	llama.SamplerChainAdd(sampler, llama.SamplerInitTopK(topK))

	// Add Top-P
	topP := float32(0.95)
	if options.TopP != nil && *options.TopP > 0 {
		topP = float32(*options.TopP)
	}
	llama.SamplerChainAdd(sampler, llama.SamplerInitTopP(topP, 0))

	// Add Min-P
	llama.SamplerChainAdd(sampler, llama.SamplerInitMinP(0.05, 0))

	// Add Temperature
	temperature := float32(0.8)
	if options.Temperature != nil && *options.Temperature > 0 {
		temperature = float32(*options.Temperature)
	}
	llama.SamplerChainAdd(sampler, llama.SamplerInitTempExt(temperature, 0, 1.0))

	// Add distribution sampler
	llama.SamplerChainAdd(sampler, llama.SamplerInitDist(llama.DefaultSeed))

	m.libService.chatSampler = sampler
}

// generateStreaming performs streaming generation and writes to StreamWriter
func (m *LlamaEinoModel) generateStreaming(ctx context.Context, writer *schema.StreamWriter[*schema.Message], prompt string, maxTokens int32) error {
	m.libService.chatMutex.Lock()
	defer m.libService.chatMutex.Unlock()

	if m.libService.chatModel == 0 || m.libService.chatContext == 0 {
		return fmt.Errorf("chat model not loaded")
	}

	// Tokenize prompt
	tokens := llama.Tokenize(m.libService.chatVocab, prompt, true, true)
	if len(tokens) == 0 {
		return fmt.Errorf("failed to tokenize prompt")
	}

	// Reset sampler state
	llama.SamplerReset(m.libService.chatSampler)

	// Decode prompt tokens
	batch := llama.BatchGetOne(tokens)
	if llama.Decode(m.libService.chatContext, batch) != 0 {
		return fmt.Errorf("failed to decode prompt")
	}

	// Send first chunk with role
	firstChunk := &schema.Message{
		Role:    schema.Assistant,
		Content: "",
	}
	if closed := writer.Send(firstChunk, nil); closed {
		return fmt.Errorf("stream closed unexpectedly")
	}

	// Generate tokens and stream
	for nGenerated := int32(0); nGenerated < maxTokens; nGenerated++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// Sample next token
		token := llama.SamplerSample(m.libService.chatSampler, m.libService.chatContext, -1)

		// Check for end of generation
		if llama.VocabIsEOG(m.libService.chatVocab, token) {
			// Send final chunk with finish reason
			finalChunk := &schema.Message{
				Role:    schema.Assistant,
				Content: "",
				ResponseMeta: &schema.ResponseMeta{
					FinishReason: "stop",
				},
			}
			writer.Send(finalChunk, nil)
			break
		}

		// Convert token to text
		buf := make([]byte, 256)
		length := llama.TokenToPiece(m.libService.chatVocab, token, buf, 0, false)
		content := string(buf[:length])

		// Send chunk
		chunk := &schema.Message{
			Role:    schema.Assistant,
			Content: content,
		}
		if closed := writer.Send(chunk, nil); closed {
			return nil // Stream closed by receiver
		}

		// Accept token and prepare for next generation
		llama.SamplerAccept(m.libService.chatSampler, token)

		// Decode the new token
		nextBatch := llama.BatchGetOne([]llama.Token{token})
		if llama.Decode(m.libService.chatContext, nextBatch) != 0 {
			return fmt.Errorf("failed to decode token")
		}

		// Small delay to prevent overwhelming
		time.Sleep(5 * time.Millisecond)
	}

	return nil
}

// Compile-time check that LlamaEinoModel implements ToolCallingChatModel
var _ model.ToolCallingChatModel = (*LlamaEinoModel)(nil)
