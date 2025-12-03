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
	"log"
	"strings"
	"time"

	"github.com/kawai-network/veridium/pkg/yzma/llama"
	"github.com/kawai-network/veridium/pkg/yzma/message"
	"github.com/kawai-network/veridium/pkg/yzma/template"
	"github.com/kawai-network/veridium/pkg/yzma/tools"
)

// Default max iterations for agent loop
const DefaultMaxIterations = 10

// Default batch size for token processing (must match context NBatch parameter)
// This prevents GGML_ASSERT failures when processing long prompts
const DefaultBatchSize = 2048

// YzmaResponse represents a response from yzma model generation
type YzmaResponse struct {
	Content          string
	ToolCalls        []message.ToolCall
	FinishReason     string
	PromptTokens     int
	CompletionTokens int
	TotalTokens      int
}

// LlamaYzmaModel implements chat model with yzma tool calling
type LlamaYzmaModel struct {
	libService   *LibraryService
	toolRegistry *tools.ToolRegistry
	toolNames    []string // Requested tool names
	noTools      bool     // If true, disable all tools (for title/summary generation)
}

// NewLlamaYzmaModel creates a new yzma-based model adapter
func NewLlamaYzmaModel(libService *LibraryService, toolRegistry *tools.ToolRegistry) *LlamaYzmaModel {
	return &LlamaYzmaModel{
		libService:   libService,
		toolRegistry: toolRegistry,
		toolNames:    []string{}, // Empty = all tools
		noTools:      false,
	}
}

// WithTools sets the tools to use (empty = all enabled tools)
func (m *LlamaYzmaModel) WithTools(toolNames []string) *LlamaYzmaModel {
	return &LlamaYzmaModel{
		libService:   m.libService,
		toolRegistry: m.toolRegistry,
		toolNames:    toolNames,
		noTools:      false,
	}
}

// WithoutTools returns a model instance with tools completely disabled
// Use this for title/summary generation where tools should not be included
func (m *LlamaYzmaModel) WithoutTools() *LlamaYzmaModel {
	return &LlamaYzmaModel{
		libService:   m.libService,
		toolRegistry: m.toolRegistry,
		toolNames:    nil,
		noTools:      true,
	}
}

// Generate generates a response from yzma messages (native yzma interface)
func (m *LlamaYzmaModel) Generate(ctx context.Context, messages []message.Message) (*YzmaResponse, error) {
	// Ensure chat model is loaded
	if !m.libService.IsChatModelLoaded() {
		if err := m.libService.LoadChatModel(""); err != nil {
			return nil, fmt.Errorf("failed to load chat model: %w", err)
		}
	}

	// Enhance system prompt with tools if available
	enhancedMessages := m.enhanceWithTools(messages)

	// Apply chat template
	chatTemplate := llama.ModelChatTemplate(m.libService.chatModel, "")
	prompt, err := template.Apply(chatTemplate, enhancedMessages, true)
	if err != nil {
		return nil, fmt.Errorf("failed to apply template: %w", err)
	}

	// Generate response
	response, err := m.libService.Generate(prompt, 32768)
	if err != nil {
		return nil, fmt.Errorf("generation failed: %w", err)
	}

	// Parse tool calls from response
	toolCalls := tools.ParseToolCalls(response)

	if len(toolCalls) > 0 {
		log.Printf("🔧 Detected %d tool calls in response", len(toolCalls))
		for i, tc := range toolCalls {
			log.Printf("🔧 Tool call #%d: %s(%v)", i+1, tc.Function.Name, tc.Function.Arguments)
		}

		// Remove tool call tags from response for cleaner output
		response = m.cleanToolCallTags(response)
	}

	return &YzmaResponse{
		Content:          response,
		ToolCalls:        toolCalls,
		FinishReason:     "stop",
		PromptTokens:     len(prompt) / 4,
		CompletionTokens: len(response) / 4,
		TotalTokens:      (len(prompt) + len(response)) / 4,
	}, nil
}

// enhanceWithTools adds tool definitions to system prompt
func (m *LlamaYzmaModel) enhanceWithTools(messages []message.Message) []message.Message {
	// Skip tool enhancement if noTools flag is set (for title/summary generation)
	if m.noTools {
		return messages
	}

	if m.toolRegistry == nil || len(messages) == 0 {
		return messages
	}

	toolsJSON, err := m.toolRegistry.FormatForPrompt(m.toolNames)
	if err != nil || toolsJSON == "" {
		return messages
	}

	// Make a copy to avoid modifying original
	result := make([]message.Message, len(messages))
	copy(result, messages)

	// Find or create system message
	if result[0].GetRole() == "system" {
		// Enhance existing system message
		if chat, ok := result[0].(message.Chat); ok {
			enhancedContent := tools.BuildSystemPrompt(chat.Content, toolsJSON)
			result[0] = message.Chat{
				Role:    "system",
				Content: enhancedContent,
			}
		}
	} else {
		// Prepend new system message with tools
		systemContent := tools.BuildSystemPrompt("You are a helpful AI assistant.", toolsJSON)
		result = append([]message.Message{
			message.Chat{Role: "system", Content: systemContent},
		}, result...)
	}

	log.Printf("🔧 Enhanced system prompt with %d tools", len(m.toolRegistry.GetByNames(m.toolNames)))
	return result
}

// cleanToolCallTags removes tool call XML tags from response
func (m *LlamaYzmaModel) cleanToolCallTags(response string) string {
	cleanResponse := response
	for strings.Contains(cleanResponse, "<tool_call>") {
		start := strings.Index(cleanResponse, "<tool_call>")
		end := strings.Index(cleanResponse, "</tool_call>")
		if start != -1 && end != -1 {
			cleanResponse = cleanResponse[:start] + cleanResponse[end+len("</tool_call>"):]
		} else {
			break
		}
	}
	return strings.TrimSpace(cleanResponse)
}

// ExecuteToolCalls executes tool calls and returns tool response messages
func (m *LlamaYzmaModel) ExecuteToolCalls(ctx context.Context, toolCalls []message.ToolCall) ([]message.Message, error) {
	if m.toolRegistry == nil {
		return nil, fmt.Errorf("tool registry not available")
	}

	toolMessages := make([]message.Message, 0, len(toolCalls))

	for _, tc := range toolCalls {
		log.Printf("🔧 Executing tool: %s", tc.Function.Name)

		// Execute tool with arguments directly (already map[string]string)
		result, err := m.toolRegistry.Execute(ctx, tc.Function.Name, tc.Function.Arguments)
		if err != nil {
			log.Printf("⚠️  Tool execution failed: %v", err)
			result = fmt.Sprintf("Error: %v", err)
		} else {
			log.Printf("✅ Tool result: %s", result[:minInt(100, len(result))])
		}

		// Create tool response message
		toolMsg := message.ToolResponse{
			Role:    "tool",
			Name:    tc.Function.Name,
			Content: result,
		}

		toolMessages = append(toolMessages, toolMsg)
	}

	return toolMessages, nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// RunAgentLoop runs the agent loop with tool calling until no more tool calls or max iterations
// This is the main entry point for agentic conversations (native yzma interface)
func (m *LlamaYzmaModel) RunAgentLoop(ctx context.Context, messages []message.Message, maxIterations int) (*YzmaResponse, []message.Message, error) {
	if maxIterations <= 0 {
		maxIterations = DefaultMaxIterations
	}

	// Track all messages for history
	allMessages := make([]message.Message, len(messages))
	copy(allMessages, messages)

	var finalResponse *YzmaResponse
	var allToolMessages []message.Message

	for i := 0; i < maxIterations; i++ {
		log.Printf("🔄 Agent loop iteration %d/%d", i+1, maxIterations)

		// Generate response with tools
		resp, err := m.Generate(ctx, allMessages)
		if err != nil {
			return nil, allToolMessages, fmt.Errorf("generation failed at iteration %d: %w", i+1, err)
		}

		// Add assistant response to history
		if len(resp.ToolCalls) > 0 {
			// Assistant message with tool calls
			allMessages = append(allMessages, message.Tool{
				Role:      "assistant",
				ToolCalls: resp.ToolCalls,
			})
		} else {
			// Regular assistant message
			allMessages = append(allMessages, message.Chat{
				Role:    "assistant",
				Content: resp.Content,
			})
		}
		finalResponse = resp

		// Check if no tool calls - we're done
		if len(resp.ToolCalls) == 0 {
			log.Printf("✅ Agent loop completed after %d iterations (no tool calls)", i+1)
			return finalResponse, allToolMessages, nil
		}

		// Execute tool calls
		toolMessages, err := m.ExecuteToolCalls(ctx, resp.ToolCalls)
		if err != nil {
			log.Printf("⚠️  Tool execution error: %v", err)
			// Continue anyway with error message
		}

		// Add tool messages to history
		for _, tm := range toolMessages {
			allMessages = append(allMessages, tm)
			allToolMessages = append(allToolMessages, tm)
		}

		// Clear KV cache before regenerating (from yzma/examples/tooluse pattern)
		if err := m.clearKVCache(); err != nil {
			log.Printf("⚠️  Failed to clear KV cache: %v", err)
		}
	}

	log.Printf("⚠️  Agent loop reached max iterations (%d)", maxIterations)
	return finalResponse, allToolMessages, nil
}

// clearKVCache clears the KV cache for a fresh generation
func (m *LlamaYzmaModel) clearKVCache() error {
	m.libService.chatMutex.Lock()
	defer m.libService.chatMutex.Unlock()

	if m.libService.chatContext == 0 {
		return fmt.Errorf("chat context not initialized")
	}

	mem, err := llama.GetMemory(m.libService.chatContext)
	if err != nil {
		return fmt.Errorf("failed to get memory: %w", err)
	}

	if err := llama.MemoryClear(mem, true); err != nil {
		return fmt.Errorf("failed to clear memory: %w", err)
	}

	log.Printf("🧹 KV cache cleared")
	return nil
}

// StreamCallback is called for each generated token during streaming
type StreamCallback func(token string, isLast bool)

// Stream generates a response with streaming (native yzma interface)
// Returns the response and any error
func (m *LlamaYzmaModel) Stream(ctx context.Context, messages []message.Message, callback StreamCallback) (*YzmaResponse, error) {
	// Ensure chat model is loaded
	if !m.libService.IsChatModelLoaded() {
		if err := m.libService.LoadChatModel(""); err != nil {
			return nil, fmt.Errorf("failed to load chat model: %w", err)
		}
	}

	// Enhance system prompt with tools if available
	enhancedMessages := m.enhanceWithTools(messages)

	// Apply chat template
	chatTemplate := llama.ModelChatTemplate(m.libService.chatModel, "")
	prompt, err := template.Apply(chatTemplate, enhancedMessages, true)
	if err != nil {
		return nil, fmt.Errorf("failed to apply template: %w", err)
	}

	// Generate with streaming
	response, err := m.generateStreaming(ctx, prompt, 32768, callback)
	if err != nil {
		return nil, fmt.Errorf("streaming generation failed: %w", err)
	}

	// Parse tool calls from response
	toolCalls := tools.ParseToolCalls(response)

	if len(toolCalls) > 0 {
		log.Printf("🔧 Detected %d tool calls in streaming response", len(toolCalls))
		response = m.cleanToolCallTags(response)
	}

	return &YzmaResponse{
		Content:          response,
		ToolCalls:        toolCalls,
		FinishReason:     "stop",
		PromptTokens:     len(prompt) / 4,
		CompletionTokens: len(response) / 4,
		TotalTokens:      (len(prompt) + len(response)) / 4,
	}, nil
}

// generateStreaming performs token-by-token generation with callback
func (m *LlamaYzmaModel) generateStreaming(ctx context.Context, prompt string, maxTokens int32, callback StreamCallback) (string, error) {
	m.libService.chatMutex.Lock()
	defer m.libService.chatMutex.Unlock()

	if m.libService.chatModel == 0 || m.libService.chatContext == 0 {
		return "", fmt.Errorf("chat model not loaded")
	}

	// Tokenize prompt
	tokens := llama.Tokenize(m.libService.chatVocab, prompt, true, true)
	if len(tokens) == 0 {
		return "", fmt.Errorf("failed to tokenize prompt")
	}

	// Reset sampler state
	llama.SamplerReset(m.libService.chatSampler)

	// CRITICAL FIX: Process prompt tokens in batches to avoid GGML_ASSERT
	// When prompt is long (especially with tool definitions), total tokens can exceed
	// the batch size limit (2048), causing llama.cpp to crash with:
	// GGML_ASSERT(n_tokens_all <= cparams.n_batch) failed
	//
	// Solution: Split tokens into batches of DefaultBatchSize and process sequentially
	log.Printf("🔢 Processing %d prompt tokens in batches of %d", len(tokens), DefaultBatchSize)

	for i := 0; i < len(tokens); i += DefaultBatchSize {
		end := i + DefaultBatchSize
		if end > len(tokens) {
			end = len(tokens)
		}

		batchTokens := tokens[i:end]
		batch := llama.BatchGetOne(batchTokens)

		errCode, err := llama.Decode(m.libService.chatContext, batch)
		if err != nil || errCode != 0 {
			return "", fmt.Errorf("failed to decode prompt batch %d-%d: %w", i, end, err)
		}

		log.Printf("✅ Processed token batch %d-%d/%d", i, end, len(tokens))
	}

	var response strings.Builder

	// Generate tokens and stream
	for nGenerated := int32(0); nGenerated < maxTokens; nGenerated++ {
		// Check context cancellation
		select {
		case <-ctx.Done():
			return response.String(), ctx.Err()
		default:
		}

		// Sample next token
		token := llama.SamplerSample(m.libService.chatSampler, m.libService.chatContext, -1)

		// Check for end of generation
		if llama.VocabIsEOG(m.libService.chatVocab, token) {
			if callback != nil {
				callback("", true)
			}
			break
		}

		// Convert token to text
		buf := make([]byte, 256)
		length := llama.TokenToPiece(m.libService.chatVocab, token, buf, 0, false)
		content := string(buf[:length])

		response.WriteString(content)

		// Call streaming callback
		if callback != nil {
			callback(content, false)
		}

		// Accept token and prepare for next generation
		llama.SamplerAccept(m.libService.chatSampler, token)

		// Decode the new token
		nextBatch := llama.BatchGetOne([]llama.Token{token})
		errCode, err := llama.Decode(m.libService.chatContext, nextBatch)
		if err != nil || errCode != 0 {
			return response.String(), fmt.Errorf("failed to decode token: %w", err)
		}

		// Small delay to prevent overwhelming
		time.Sleep(5 * time.Millisecond)
	}

	return response.String(), nil
}

// RunAgentLoopWithStreaming runs agent loop with streaming support (native yzma interface)
func (m *LlamaYzmaModel) RunAgentLoopWithStreaming(ctx context.Context, messages []message.Message, maxIterations int, callback StreamCallback) (*YzmaResponse, []message.Message, error) {
	if maxIterations <= 0 {
		maxIterations = DefaultMaxIterations
	}

	allMessages := make([]message.Message, len(messages))
	copy(allMessages, messages)

	var finalResponse *YzmaResponse
	var allToolMessages []message.Message

	for i := 0; i < maxIterations; i++ {
		log.Printf("🔄 Agent loop (streaming) iteration %d/%d", i+1, maxIterations)

		// Generate response with streaming
		resp, err := m.Stream(ctx, allMessages, callback)
		if err != nil {
			return nil, allToolMessages, fmt.Errorf("streaming generation failed at iteration %d: %w", i+1, err)
		}

		// Add assistant response to history
		if len(resp.ToolCalls) > 0 {
			allMessages = append(allMessages, message.Tool{
				Role:      "assistant",
				ToolCalls: resp.ToolCalls,
			})
		} else {
			allMessages = append(allMessages, message.Chat{
				Role:    "assistant",
				Content: resp.Content,
			})
		}
		finalResponse = resp

		if len(resp.ToolCalls) == 0 {
			log.Printf("✅ Agent loop (streaming) completed after %d iterations", i+1)
			return finalResponse, allToolMessages, nil
		}

		// Execute tool calls
		toolMessages, err := m.ExecuteToolCalls(ctx, resp.ToolCalls)
		if err != nil {
			log.Printf("⚠️  Tool execution error: %v", err)
		}

		for _, tm := range toolMessages {
			allMessages = append(allMessages, tm)
			allToolMessages = append(allToolMessages, tm)
		}

		// Clear KV cache before regenerating
		if err := m.clearKVCache(); err != nil {
			log.Printf("⚠️  Failed to clear KV cache: %v", err)
		}
	}

	log.Printf("⚠️  Agent loop (streaming) reached max iterations (%d)", maxIterations)
	return finalResponse, allToolMessages, nil
}
