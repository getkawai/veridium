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

package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/kawai-network/veridium/pkg/yzma/tools"
	"github.com/kawai-network/veridium/types"
)

// Note: Provider implements llm.Provider interface
// Interface compliance is verified at factory.go level to avoid import cycle

// Provider implements llm.Provider for OpenAI-compatible APIs (OpenRouter, Zhipu GLM)
type Provider struct {
	client       *Client
	config       types.ProviderConfig
	toolRegistry *tools.ToolRegistry
	toolNames    []string // Specific tools to use (empty = all)
	noTools      bool     // If true, disable all tools
}

// NewProvider creates a new OpenAI-compatible provider
func NewProvider(config types.ProviderConfig, toolRegistry *tools.ToolRegistry) *Provider {
	// Set default model if not specified
	if config.Model == "" {
		if defaultModel, ok := types.DefaultModels[config.Type]; ok {
			config.Model = defaultModel
		}
	}

	// Set default max tokens if not specified
	if config.MaxTokens == 0 {
		config.MaxTokens = 4096
	}

	return &Provider{
		client:       NewClient(config),
		config:       config,
		toolRegistry: toolRegistry,
		toolNames:    []string{},
		noTools:      false,
	}
}

// WithTools returns a new provider configured with specific tools
func (p *Provider) WithTools(toolNames []string) *Provider {
	return &Provider{
		client:       p.client,
		config:       p.config,
		toolRegistry: p.toolRegistry,
		toolNames:    toolNames,
		noTools:      false,
	}
}

// WithoutTools returns a new provider with tools disabled
func (p *Provider) WithoutTools() *Provider {
	return &Provider{
		client:       p.client,
		config:       p.config,
		toolRegistry: p.toolRegistry,
		toolNames:    nil,
		noTools:      true,
	}
}

// Generate generates a response from messages (single turn)
func (p *Provider) Generate(ctx context.Context, messages types.Prompt) (*types.LLMResponse, error) {
	// Convert yzma messages to API format
	apiMessages := p.convertMessages(messages)

	// Build request
	req := types.ChatCompletionRequest{
		Model:     p.config.Model,
		Messages:  apiMessages,
		MaxTokens: &p.config.MaxTokens,
	}

	// Disable thinking/reasoning mode for Zhipu to get faster responses
	// This prevents the model from using tokens for reasoning_content
	if p.config.Type == types.ProviderZhipuAI {
		req.Thinking = &types.ThinkingConfig{Type: "disabled"}
	}

	// Add tools if enabled
	if !p.noTools && p.toolRegistry != nil {
		req.Tools = p.getToolDefinitions()
		if len(req.Tools) > 0 {
			req.ToolChoice = "auto"
		}
	}

	// Add temperature if specified
	if p.config.Temperature > 0 {
		req.Temperature = &p.config.Temperature
	}

	// Send request
	resp, err := p.client.ChatCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	// Convert response
	return p.convertResponse(resp), nil
}

// RunAgentLoop runs the agent loop with tool execution
func (p *Provider) RunAgentLoop(ctx context.Context, messages types.Prompt, maxIterations int) (*types.LLMResponse, types.Prompt, error) {
	if maxIterations <= 0 {
		maxIterations = 10
	}

	allMessages := make(types.Prompt, len(messages))
	copy(allMessages, messages)

	var finalResponse *types.LLMResponse
	var allToolMessages types.Prompt
	var allToolCalls []types.ToolCall

	for i := 0; i < maxIterations; i++ {
		log.Printf("🔄 [%s] Agent loop iteration %d/%d", p.config.Type, i+1, maxIterations)

		// Generate response
		resp, err := p.Generate(ctx, allMessages)
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
			log.Printf("✅ [%s] Agent loop completed after %d iterations", p.config.Type, i+1)
			finalResponse.ToolCalls = allToolCalls
			return finalResponse, allToolMessages, nil
		}

		// Execute tool calls
		toolMessages, err := p.executeToolCalls(ctx, resp.ToolCalls)
		if err != nil {
			log.Printf("⚠️  Tool execution error: %v", err)
		}

		// Add tool messages to history
		allMessages = append(allMessages, toolMessages...)
		allToolMessages = append(allToolMessages, toolMessages...)
	}

	log.Printf("⚠️  [%s] Agent loop reached max iterations (%d)", p.config.Type, maxIterations)
	if finalResponse != nil {
		finalResponse.ToolCalls = allToolCalls
	}
	return finalResponse, allToolMessages, nil
}

// RunAgentLoopWithStreaming runs the agent loop with streaming
func (p *Provider) RunAgentLoopWithStreaming(ctx context.Context, messages types.Prompt, maxIterations int, streamCallback types.StreamCallback, toolCallback types.ToolEventCallback) (*types.LLMResponse, types.Prompt, error) {
	if maxIterations <= 0 {
		maxIterations = 10
	}

	allMessages := make(types.Prompt, len(messages))
	copy(allMessages, messages)

	var finalResponse *types.LLMResponse
	var allToolMessages types.Prompt
	var allToolCalls []types.ToolCall

	for i := 0; i < maxIterations; i++ {
		log.Printf("🔄 [%s] Agent loop (streaming) iteration %d/%d", p.config.Type, i+1, maxIterations)

		// Generate with streaming
		resp, err := p.generateWithStreaming(ctx, allMessages, streamCallback)
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

		// Check if no tool calls - we're done
		if len(resp.ToolCalls) == 0 {
			log.Printf("✅ [%s] Agent loop (streaming) completed after %d iterations", p.config.Type, i+1)
			finalResponse.ToolCalls = allToolCalls
			return finalResponse, allToolMessages, nil
		}

		// Emit tool_call events before execution
		if toolCallback != nil {
			for _, tc := range resp.ToolCalls {
				toolCallback(types.ChatEventToolCall, tc, "")
			}
		}

		// Execute tool calls
		toolMessages, err := p.executeToolCalls(ctx, resp.ToolCalls)
		if err != nil {
			log.Printf("⚠️  Tool execution error: %v", err)
		}

		// Emit tool_result events after execution
		if toolCallback != nil {
			for idx, tc := range resp.ToolCalls {
				if idx < len(toolMessages) {
					part := toolMessages[idx].Content[0].(types.ToolResultPart)
					toolCallback(types.ChatEventToolResult, tc, part.Content)
				}
			}
		}

		// Add tool messages to history
		allMessages = append(allMessages, toolMessages...)
		allToolMessages = append(allToolMessages, toolMessages...)
	}

	log.Printf("⚠️  [%s] Agent loop (streaming) reached max iterations (%d)", p.config.Type, maxIterations)
	if finalResponse != nil {
		finalResponse.ToolCalls = allToolCalls
	}
	return finalResponse, allToolMessages, nil
}

// generateWithStreaming generates a response with streaming
func (p *Provider) generateWithStreaming(ctx context.Context, messages types.Prompt, callback types.StreamCallback) (*types.LLMResponse, error) {
	// Convert messages
	apiMessages := p.convertMessages(messages)

	// Build request
	req := types.ChatCompletionRequest{
		Model:     p.config.Model,
		Messages:  apiMessages,
		MaxTokens: &p.config.MaxTokens,
		Stream:    true,
	}

	// Add tools if enabled
	if !p.noTools && p.toolRegistry != nil {
		req.Tools = p.getToolDefinitions()
		if len(req.Tools) > 0 {
			req.ToolChoice = "auto"
		}
	}

	// Stream callback adapter
	streamAdapter := func(chunk *types.ChatCompletionStreamResponse) error {
		if callback != nil && len(chunk.Choices) > 0 {
			content, _ := chunk.Choices[0].Delta.Content.(string)
			isLast := chunk.Choices[0].FinishReason != ""
			callback(content, isLast)
		}
		return nil
	}

	// Send streaming request
	resp, err := p.client.ChatCompletionStream(ctx, req, streamAdapter)
	if err != nil {
		return nil, err
	}

	return p.convertResponse(resp), nil
}

// convertMessages converts messages to API format
func (p *Provider) convertMessages(messages types.Prompt) []types.ChatCompletionMsg {
	result := make([]types.ChatCompletionMsg, 0, len(messages))

	for _, msg := range messages {
		apiMsg := types.ChatCompletionMsg{
			Role: msg.GetRole(),
		}

		// Check for tool calls
		if msg.HasToolCalls() {
			toolCalls := msg.GetToolCalls()
			apiToolCalls := make([]types.APIToolCall, len(toolCalls))
			for i, tc := range toolCalls {
				argsJSON, _ := json.Marshal(tc.Function.Arguments)
				apiToolCalls[i] = types.APIToolCall{
					ID:   tc.ID,
					Type: "function",
					Function: types.APIToolCallFunction{
						Name:      tc.Function.Name,
						Arguments: string(argsJSON),
					},
				}
			}
			apiMsg.ToolCalls = apiToolCalls
		}

		// Check for tool result
		for _, part := range msg.Content {
			switch p := part.(type) {
			case types.TextPart:
				apiMsg.Content = p.Text
			case types.ToolResultPart:
				apiMsg.Content = p.Content
				apiMsg.Name = p.ToolName
				apiMsg.ToolCallID = p.ToolCallID
			}
		}

		result = append(result, apiMsg)
	}

	return result
}

// getToolDefinitions builds tool definitions from registry
func (p *Provider) getToolDefinitions() []types.APIToolDefinition {
	if p.toolRegistry == nil {
		return nil
	}

	registeredTools := p.toolRegistry.GetByNames(p.toolNames)
	definitions := make([]types.APIToolDefinition, 0, len(registeredTools))

	for _, tool := range registeredTools {
		if !tool.Enabled {
			continue
		}

		definitions = append(definitions, types.APIToolDefinition{
			Type: "function",
			Function: types.APIToolFunction{
				Name:        tool.Function.Name,
				Description: tool.Function.Description,
				Parameters:  tool.Function.Parameters,
			},
		})
	}

	return definitions
}

// executeToolCalls executes tool calls and returns tool response messages
func (p *Provider) executeToolCalls(ctx context.Context, toolCalls []types.ToolCall) (types.Prompt, error) {
	if p.toolRegistry == nil {
		return nil, fmt.Errorf("tool registry not available")
	}

	toolMessages := make(types.Prompt, 0, len(toolCalls))

	for _, tc := range toolCalls {
		log.Printf("🔧 [%s] Executing tool: %s", p.config.Type, tc.Function.Name)

		result, err := p.toolRegistry.Execute(ctx, tc.Function.Name, tc.Function.Arguments)
		if err != nil {
			log.Printf("⚠️  Tool execution failed: %v", err)
			toolMessages = append(toolMessages, types.NewToolErrorMessage(tc.ID, tc.Function.Name, fmt.Sprintf("Error: %v", err)))
		} else {
			// Log truncated result
			displayResult := result
			if len(displayResult) > 100 {
				displayResult = displayResult[:100] + "..."
			}
			log.Printf("✅ Tool result: %s", displayResult)
			toolMessages = append(toolMessages, types.NewToolResultMessage(tc.ID, tc.Function.Name, result))
		}
	}

	return toolMessages, nil
}

// convertResponse converts API response to types.LLMResponse
func (p *Provider) convertResponse(resp *types.ChatCompletionResponse) *types.LLMResponse {
	if len(resp.Choices) == 0 {
		return &types.LLMResponse{
			FinishReason: "error",
		}
	}

	choice := resp.Choices[0]
	contentStr, ok := choice.Message.Content.(string)
	if !ok && choice.Message.Content != nil {
		// Try to handle other content types (e.g., []interface{} from some providers)
		log.Printf("⚠️  Content is not string: %T = %+v", choice.Message.Content, choice.Message.Content)
		// Try to extract text from content if it's a slice of content parts
		if parts, ok := choice.Message.Content.([]interface{}); ok {
			for _, part := range parts {
				if partMap, ok := part.(map[string]interface{}); ok {
					if text, ok := partMap["text"].(string); ok {
						contentStr += text
					}
				}
			}
		}
	}
	result := &types.LLMResponse{
		Content:      strings.TrimSpace(contentStr),
		FinishReason: choice.FinishReason,
	}

	// Convert API tool calls to types.ToolCall
	if len(choice.Message.ToolCalls) > 0 {
		result.ToolCalls = make([]types.ToolCall, len(choice.Message.ToolCalls))
		for i, tc := range choice.Message.ToolCalls {
			// Parse arguments from JSON string to map[string]string
			var args map[string]string
			if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
				// Try as map[string]interface{} and convert
				var argsAny map[string]interface{}
				if err := json.Unmarshal([]byte(tc.Function.Arguments), &argsAny); err == nil {
					args = make(map[string]string)
					for k, v := range argsAny {
						args[k] = fmt.Sprintf("%v", v)
					}
				}
			}

			result.ToolCalls[i] = types.ToolCall{
				ID:   tc.ID,
				Type: tc.Type,
				Function: types.ToolFunction{
					Name:      tc.Function.Name,
					Arguments: args,
				},
			}
		}
	}

	// Set usage
	if resp.Usage != nil {
		result.PromptTokens = resp.Usage.PromptTokens
		result.CompletionTokens = resp.Usage.CompletionTokens
		result.TotalTokens = resp.Usage.TotalTokens
	}

	return result
}

// GetConfig returns the provider configuration
func (p *Provider) GetConfig() types.ProviderConfig {
	return p.config
}

// SetModel changes the model
func (p *Provider) SetModel(model string) {
	p.config.Model = model
}
