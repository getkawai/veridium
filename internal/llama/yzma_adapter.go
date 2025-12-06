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
	"github.com/kawai-network/veridium/pkg/yzma/template"
	"github.com/kawai-network/veridium/pkg/yzma/tools"
	"github.com/kawai-network/veridium/types"
)

// Default max iterations for agent loop
const DefaultMaxIterations = 10

// Default batch size for token processing (must match context NBatch parameter)
// This prevents GGML_ASSERT failures when processing long prompts
const DefaultBatchSize = 2048

// Llama32ToolTemplate is a custom chat template for Llama 3.2 models with tool support.
// The original Llama 3.2 template has raise_exception for multiple tool calls,
// but gonja doesn't support that function. This template provides tool call support
// without the restrictive validation, following the same structure as Qwen.
const Llama32ToolTemplate = `{%- if tools %}
{{- '<|begin_of_text|><|start_header_id|>system<|end_header_id|>

' }}
{%- if messages[0]['role'] == 'system' %}
{{- messages[0]['content'] }}
{%- else %}
{{- 'You are a helpful assistant with access to the following functions.' }}
{%- endif %}
{{- '

# Tools

You may call one or more functions to assist with the user query.

You are provided with function signatures within <tools></tools> XML tags:
<tools>' }}
{%- for tool in tools %}
{{- '
' }}
{{- tool | tojson }}
{%- endfor %}
{{- '
</tools>

For each function call, return a json object with function name and arguments within <tool_call></tool_call> XML tags:
<tool_call>
{"name": <function-name>, "arguments": <args-json-object>}
</tool_call><|eot_id|>
' }}
{%- else %}
{%- if messages[0]['role'] == 'system' %}
{{- '<|begin_of_text|><|start_header_id|>system<|end_header_id|>

' + messages[0]['content'] + '<|eot_id|>
' }}
{%- else %}
{{- '<|begin_of_text|><|start_header_id|>system<|end_header_id|>

You are a helpful assistant.<|eot_id|>
' }}
{%- endif %}
{%- endif %}
{%- for message in messages %}
{%- if (message.role == "user") or (message.role == "system" and not loop.first) or (message.role == "assistant" and not message.tool_calls) %}
{{- '<|start_header_id|>' + message.role + '<|end_header_id|>

' + message.content + '<|eot_id|>
' }}
{%- elif message.role == "assistant" %}
{{- '<|start_header_id|>' + message.role + '<|end_header_id|>
' }}
{%- if message.content %}
{{- message.content }}
{%- endif %}
{%- for tool_call in message.tool_calls %}
{%- if tool_call.function is defined %}
{%- set tool_call = tool_call.function %}
{%- endif %}
{{- '
<tool_call>
{"name": "' }}
{{- tool_call.name }}
{{- '", "arguments": ' }}
{{- tool_call.arguments | tojson }}
{{- '}
</tool_call>' }}
{%- endfor %}
{{- '<|eot_id|>
' }}
{%- elif message.role == "tool" %}
{%- if (loop.index0 == 0) or (messages[loop.index0 - 1].role != "tool") %}
{{- '<|start_header_id|>ipython<|end_header_id|>
' }}
{%- endif %}
{{- '
<tool_response>
' }}
{{- message.content }}
{{- '
</tool_response>' }}
{%- if loop.last or (messages[loop.index0 + 1].role != "tool") %}
{{- '<|eot_id|>
' }}
{%- endif %}
{%- endif %}
{%- endfor %}
{%- if add_generation_prompt %}
{{- '<|start_header_id|>assistant<|end_header_id|>

' }}
{%- endif %}`

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

// getChatTemplate returns the appropriate chat template for the current model
// For Llama 3.2 models, we use a custom template that supports tool calls properly
// because the embedded template has raise_exception which gonja doesn't support
func (m *LlamaYzmaModel) getChatTemplate() string {
	modelPath := m.libService.GetLoadedChatModel()
	modelPathLower := strings.ToLower(modelPath)

	// Check if this is a Llama 3.2 model
	if strings.Contains(modelPathLower, "llama-3.2") || strings.Contains(modelPathLower, "llama_3.2") || strings.Contains(modelPathLower, "llama3.2") {
		log.Printf("🔧 Using custom Llama 3.2 tool template")
		return Llama32ToolTemplate
	}

	// Use embedded template from model for other models
	return llama.ModelChatTemplate(m.libService.chatModel, "")
}

// Generate generates a response from messages (native interface)
func (m *LlamaYzmaModel) Generate(ctx context.Context, messages types.Prompt) (*types.LLMResponse, error) {
	// Ensure chat model is loaded
	if !m.libService.IsChatModelLoaded() {
		if err := m.libService.LoadChatModel(""); err != nil {
			return nil, fmt.Errorf("failed to load chat model: %w", err)
		}
	}

	// Enhance system prompt with tools if available
	enhancedMessages := m.enhanceWithTools(messages)

	// Apply chat template - use custom template for Llama 3.2
	chatTemplate := m.getChatTemplate()
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

	return &types.LLMResponse{
		Content:          response,
		ToolCalls:        toolCalls,
		FinishReason:     "stop",
		PromptTokens:     len(prompt) / 4,
		CompletionTokens: len(response) / 4,
		TotalTokens:      (len(prompt) + len(response)) / 4,
	}, nil
}

// enhanceWithTools adds tool definitions to system prompt
func (m *LlamaYzmaModel) enhanceWithTools(messages types.Prompt) types.Prompt {
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
	result := make(types.Prompt, len(messages))
	copy(result, messages)

	// Find or create system message
	if result[0].GetRole() == "system" {
		// Enhance existing system message
		enhancedContent := tools.BuildSystemPrompt(result[0].GetText(), toolsJSON)
		result[0] = types.NewSystemMessage(enhancedContent)
	} else {
		// Prepend new system message with tools
		systemContent := tools.BuildSystemPrompt("You are a helpful AI assistant.", toolsJSON)
		result = append(types.Prompt{types.NewSystemMessage(systemContent)}, result...)
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
func (m *LlamaYzmaModel) ExecuteToolCalls(ctx context.Context, toolCalls []types.ToolCall) (types.Prompt, error) {
	if m.toolRegistry == nil {
		return nil, fmt.Errorf("tool registry not available")
	}

	toolMessages := make(types.Prompt, 0, len(toolCalls))

	for _, tc := range toolCalls {
		log.Printf("🔧 Executing tool: %s", tc.Function.Name)

		// Execute tool with arguments directly (already map[string]string)
		result, err := m.toolRegistry.Execute(ctx, tc.Function.Name, tc.Function.Arguments)
		if err != nil {
			log.Printf("⚠️  Tool execution failed: %v", err)
			toolMessages = append(toolMessages, types.NewToolErrorMessage(tc.ID, tc.Function.Name, fmt.Sprintf("Error: %v", err)))
		} else {
			log.Printf("✅ Tool result: %s", result[:minInt(100, len(result))])
			toolMessages = append(toolMessages, types.NewToolResultMessage(tc.ID, tc.Function.Name, result))
		}
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
// This is the main entry point for agentic conversations (native interface)
// The returned types.LLMResponse contains ALL tool calls from all iterations (not just the last one)
func (m *LlamaYzmaModel) RunAgentLoop(ctx context.Context, messages types.Prompt, maxIterations int) (*types.LLMResponse, types.Prompt, error) {
	if maxIterations <= 0 {
		maxIterations = DefaultMaxIterations
	}

	// Track all messages for history
	allMessages := make(types.Prompt, len(messages))
	copy(allMessages, messages)

	var finalResponse *types.LLMResponse
	var allToolMessages types.Prompt
	var allToolCalls []types.ToolCall // Collect all tool calls from all iterations

	for i := 0; i < maxIterations; i++ {
		log.Printf("🔄 Agent loop iteration %d/%d", i+1, maxIterations)

		// Generate response with tools
		resp, err := m.Generate(ctx, allMessages)
		if err != nil {
			return nil, allToolMessages, fmt.Errorf("generation failed at iteration %d: %w", i+1, err)
		}

		// Add assistant response to history
		if len(resp.ToolCalls) > 0 {
			allMessages = append(allMessages, types.NewToolCallMessage(resp.ToolCalls))
			allToolCalls = append(allToolCalls, resp.ToolCalls...)
		} else {
			allMessages = append(allMessages, types.NewAssistantMessage(resp.Content))
		}
		finalResponse = resp

		// Check if no tool calls - we're done
		if len(resp.ToolCalls) == 0 {
			log.Printf("✅ Agent loop completed after %d iterations (no tool calls)", i+1)
			// Include all collected tool calls in the final response
			finalResponse.ToolCalls = allToolCalls
			return finalResponse, allToolMessages, nil
		}

		// Execute tool calls
		toolMessages, err := m.ExecuteToolCalls(ctx, resp.ToolCalls)
		if err != nil {
			log.Printf("⚠️  Tool execution error: %v", err)
			// Continue anyway with error message
		}

		// Add tool messages to history
		allMessages = append(allMessages, toolMessages...)
		allToolMessages = append(allToolMessages, toolMessages...)

		// Clear KV cache before regenerating (from yzma/examples/tooluse pattern)
		if err := m.clearKVCache(); err != nil {
			log.Printf("⚠️  Failed to clear KV cache: %v", err)
		}
	}

	log.Printf("⚠️  Agent loop reached max iterations (%d)", maxIterations)
	// Include all collected tool calls in the final response even when max iterations reached
	if finalResponse != nil {
		finalResponse.ToolCalls = allToolCalls
	}
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

// Stream generates a response with streaming (native interface)
// Returns the response and any error
func (m *LlamaYzmaModel) Stream(ctx context.Context, messages types.Prompt, callback types.StreamCallback) (*types.LLMResponse, error) {
	// Ensure chat model is loaded
	if !m.libService.IsChatModelLoaded() {
		if err := m.libService.LoadChatModel(""); err != nil {
			return nil, fmt.Errorf("failed to load chat model: %w", err)
		}
	}

	// Enhance system prompt with tools if available
	enhancedMessages := m.enhanceWithTools(messages)

	// Apply chat template - use custom template for Llama 3.2
	chatTemplate := m.getChatTemplate()
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

	return &types.LLMResponse{
		Content:          response,
		ToolCalls:        toolCalls,
		FinishReason:     "stop",
		PromptTokens:     len(prompt) / 4,
		CompletionTokens: len(response) / 4,
		TotalTokens:      (len(prompt) + len(response)) / 4,
	}, nil
}

// generateStreaming performs token-by-token generation with callback
func (m *LlamaYzmaModel) generateStreaming(ctx context.Context, prompt string, maxTokens int32, callback types.StreamCallback) (string, error) {
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

// RunAgentLoopWithStreaming runs agent loop with streaming support (native interface)
// Returns the final response, all tool messages, and any error
// The returned types.LLMResponse contains ALL tool calls from all iterations (not just the last one)
// streamCallback: called for each token during generation
// toolCallback: called for tool events (tool_call before execution, tool_result after)
func (m *LlamaYzmaModel) RunAgentLoopWithStreaming(ctx context.Context, messages types.Prompt, maxIterations int, streamCallback types.StreamCallback, toolCallback types.ToolEventCallback) (*types.LLMResponse, types.Prompt, error) {
	if maxIterations <= 0 {
		maxIterations = DefaultMaxIterations
	}

	allMessages := make(types.Prompt, len(messages))
	copy(allMessages, messages)

	var finalResponse *types.LLMResponse
	var allToolMessages types.Prompt
	var allToolCalls []types.ToolCall // Collect all tool calls from all iterations

	for i := 0; i < maxIterations; i++ {
		log.Printf("🔄 Agent loop (streaming) iteration %d/%d", i+1, maxIterations)

		// Generate response with streaming
		resp, err := m.Stream(ctx, allMessages, streamCallback)
		if err != nil {
			return nil, allToolMessages, fmt.Errorf("streaming generation failed at iteration %d: %w", i+1, err)
		}

		// Add assistant response to history
		if len(resp.ToolCalls) > 0 {
			allMessages = append(allMessages, types.NewToolCallMessage(resp.ToolCalls))
			allToolCalls = append(allToolCalls, resp.ToolCalls...)
		} else {
			allMessages = append(allMessages, types.NewAssistantMessage(resp.Content))
		}
		finalResponse = resp

		if len(resp.ToolCalls) == 0 {
			log.Printf("✅ Agent loop (streaming) completed after %d iterations", i+1)
			// Include all collected tool calls in the final response
			finalResponse.ToolCalls = allToolCalls
			return finalResponse, allToolMessages, nil
		}

		// Emit tool_call events BEFORE execution (for UI loading state)
		if toolCallback != nil {
			for _, tc := range resp.ToolCalls {
				toolCallback(types.ChatEventToolCall, tc, "")
			}
		}

		// Execute tool calls
		toolMessages, err := m.ExecuteToolCalls(ctx, resp.ToolCalls)
		if err != nil {
			log.Printf("⚠️  Tool execution error: %v", err)
		}

		// Emit tool_result events AFTER execution
		if toolCallback != nil {
			for idx, tc := range resp.ToolCalls {
				if idx < len(toolMessages) {
					part := toolMessages[idx].Content[0].(types.ToolResultPart)
					toolCallback(types.ChatEventToolResult, tc, part.Content)
				}
			}
		}

		allMessages = append(allMessages, toolMessages...)
		allToolMessages = append(allToolMessages, toolMessages...)

		// Clear KV cache before regenerating
		if err := m.clearKVCache(); err != nil {
			log.Printf("⚠️  Failed to clear KV cache: %v", err)
		}
	}

	log.Printf("⚠️  Agent loop (streaming) reached max iterations (%d)", maxIterations)
	// Include all collected tool calls in the final response even when max iterations reached
	if finalResponse != nil {
		finalResponse.ToolCalls = allToolCalls
	}
	return finalResponse, allToolMessages, nil
}
