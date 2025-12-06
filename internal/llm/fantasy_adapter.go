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

package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/pkg/yzma/tools"
	"github.com/kawai-network/veridium/types"
)

// FantasyProviderAdapter wraps a fantasy.LanguageModel to implement llm.Provider
type FantasyProviderAdapter struct {
	provider     fantasy.Provider
	model        fantasy.LanguageModel
	modelID      string
	toolRegistry *tools.ToolRegistry
	toolNames    []string
	noTools      bool
}

// NewFantasyProviderAdapter creates a new adapter
func NewFantasyProviderAdapter(provider fantasy.Provider, modelID string, toolRegistry *tools.ToolRegistry) (*FantasyProviderAdapter, error) {
	model, err := provider.LanguageModel(context.Background(), modelID)
	if err != nil {
		return nil, fmt.Errorf("failed to create language model: %w", err)
	}

	return &FantasyProviderAdapter{
		provider:     provider,
		model:        model,
		modelID:      modelID,
		toolRegistry: toolRegistry,
		toolNames:    []string{},
		noTools:      false,
	}, nil
}

// WithTools returns a new provider configured with specific tools
func (a *FantasyProviderAdapter) WithTools(toolNames []string) Provider {
	return &FantasyProviderAdapter{
		provider:     a.provider,
		model:        a.model,
		modelID:      a.modelID,
		toolRegistry: a.toolRegistry,
		toolNames:    toolNames,
		noTools:      false,
	}
}

// WithoutTools returns a new provider with tools disabled
func (a *FantasyProviderAdapter) WithoutTools() Provider {
	return &FantasyProviderAdapter{
		provider:     a.provider,
		model:        a.model,
		modelID:      a.modelID,
		toolRegistry: a.toolRegistry,
		toolNames:    nil,
		noTools:      true,
	}
}

// Generate generates a response from messages (single turn, no tool execution)
func (a *FantasyProviderAdapter) Generate(ctx context.Context, messages []fantasy.Message) (*types.LLMResponse, error) {
	call := a.buildCall(fantasy.Prompt(messages))

	resp, err := a.model.Generate(ctx, call)
	if err != nil {
		return nil, err
	}

	return a.convertResponse(resp), nil
}

// RunAgentLoop runs the agent loop with tool execution
func (a *FantasyProviderAdapter) RunAgentLoop(ctx context.Context, messages fantasy.Prompt, maxIterations int) (*types.LLMResponse, fantasy.Prompt, error) {
	if maxIterations <= 0 {
		maxIterations = 10
	}

	allMessages := make(fantasy.Prompt, len(messages))
	copy(allMessages, messages)

	var finalResponse *types.LLMResponse
	var allToolMessages fantasy.Prompt
	var allToolCalls []fantasy.ToolCall

	for i := 0; i < maxIterations; i++ {
		log.Printf("🔄 [fantasy] Agent loop iteration %d/%d", i+1, maxIterations)

		call := a.buildCall(allMessages)
		resp, err := a.model.Generate(ctx, call)
		if err != nil {
			return nil, allToolMessages, fmt.Errorf("generation failed at iteration %d: %w", i+1, err)
		}

		llmResp := a.convertResponse(resp)

		// Add assistant response to history
		if len(llmResp.ToolCalls) > 0 {
			allMessages = append(allMessages, types.NewToolCallMessage(llmResp.ToolCalls))
			allToolCalls = append(allToolCalls, llmResp.ToolCalls...)
		} else {
			allMessages = append(allMessages, fantasy.Message{
				Role:    fantasy.MessageRoleAssistant,
				Content: []fantasy.MessagePart{fantasy.TextPart{Text: llmResp.Content}},
			})
		}
		finalResponse = llmResp

		// Check if no tool calls - we're done
		if len(llmResp.ToolCalls) == 0 {
			log.Printf("✅ [fantasy] Agent loop completed after %d iterations", i+1)
			finalResponse.ToolCalls = allToolCalls
			return finalResponse, allToolMessages, nil
		}

		// Execute tool calls
		toolMessages, err := a.executeToolCalls(ctx, llmResp.ToolCalls)
		if err != nil {
			log.Printf("⚠️  Tool execution error: %v", err)
		}

		// Add tool messages to history
		allMessages = append(allMessages, toolMessages...)
		allToolMessages = append(allToolMessages, toolMessages...)
	}

	log.Printf("⚠️  [fantasy] Agent loop reached max iterations (%d)", maxIterations)
	if finalResponse != nil {
		finalResponse.ToolCalls = allToolCalls
	}
	return finalResponse, allToolMessages, nil
}

// RunAgentLoopWithStreaming runs the agent loop with streaming callback
func (a *FantasyProviderAdapter) RunAgentLoopWithStreaming(ctx context.Context, messages fantasy.Prompt, maxIterations int, streamCallback types.StreamCallback, toolCallback types.ToolEventCallback) (*types.LLMResponse, fantasy.Prompt, error) {
	if maxIterations <= 0 {
		maxIterations = 10
	}

	allMessages := make(fantasy.Prompt, len(messages))
	copy(allMessages, messages)

	var finalResponse *types.LLMResponse
	var allToolMessages fantasy.Prompt
	var allToolCalls []fantasy.ToolCall

	for i := 0; i < maxIterations; i++ {
		log.Printf("🔄 [fantasy] Agent loop (streaming) iteration %d/%d", i+1, maxIterations)

		call := a.buildCall(allMessages)
		streamResp, err := a.model.Stream(ctx, call)
		if err != nil {
			return nil, allToolMessages, fmt.Errorf("streaming generation failed at iteration %d: %w", i+1, err)
		}

		llmResp := a.consumeStream(streamResp, streamCallback)

		// Add assistant response to history
		if len(llmResp.ToolCalls) > 0 {
			allMessages = append(allMessages, types.NewToolCallMessage(llmResp.ToolCalls))
			allToolCalls = append(allToolCalls, llmResp.ToolCalls...)
		} else {
			allMessages = append(allMessages, fantasy.Message{
				Role:    fantasy.MessageRoleAssistant,
				Content: []fantasy.MessagePart{fantasy.TextPart{Text: llmResp.Content}},
			})
		}
		finalResponse = llmResp

		// Check if no tool calls - we're done
		if len(llmResp.ToolCalls) == 0 {
			log.Printf("✅ [fantasy] Agent loop (streaming) completed after %d iterations", i+1)
			finalResponse.ToolCalls = allToolCalls
			return finalResponse, allToolMessages, nil
		}

		// Emit tool_call events before execution
		if toolCallback != nil {
			for _, tc := range llmResp.ToolCalls {
				toolCallback(types.ChatEventToolCall, tc, "")
			}
		}

		// Execute tool calls
		toolMessages, err := a.executeToolCalls(ctx, llmResp.ToolCalls)
		if err != nil {
			log.Printf("⚠️  Tool execution error: %v", err)
		}

		// Emit tool_result events after execution
		if toolCallback != nil {
			for idx, tc := range llmResp.ToolCalls {
				if idx < len(toolMessages) {
					part := toolMessages[idx].Content[0].(fantasy.ToolResultPart)
					toolCallback(types.ChatEventToolResult, tc, types.GetToolResultContent(part))
				}
			}
		}

		// Add tool messages to history
		allMessages = append(allMessages, toolMessages...)
		allToolMessages = append(allToolMessages, toolMessages...)
	}

	log.Printf("⚠️  [fantasy] Agent loop (streaming) reached max iterations (%d)", maxIterations)
	if finalResponse != nil {
		finalResponse.ToolCalls = allToolCalls
	}
	return finalResponse, allToolMessages, nil
}

// buildCall builds a fantasy.Call from messages
func (a *FantasyProviderAdapter) buildCall(messages fantasy.Prompt) fantasy.Call {
	prompt := a.convertToFantasyPrompt(messages)

	call := fantasy.Call{
		Prompt: prompt,
	}

	// Add tools if enabled
	if !a.noTools && a.toolRegistry != nil {
		call.Tools = a.getFantasyTools()
		if len(call.Tools) > 0 {
			toolChoice := fantasy.ToolChoiceAuto
			call.ToolChoice = &toolChoice
		}
	}

	return call
}

// convertToFantasyPrompt converts fantasy.Prompt to normalized fantasy.Prompt
func (a *FantasyProviderAdapter) convertToFantasyPrompt(messages fantasy.Prompt) fantasy.Prompt {
	result := make(fantasy.Prompt, 0, len(messages))

	for _, msg := range messages {
		var fantasyMsg fantasy.Message

		switch msg.Role {
		case fantasy.MessageRoleSystem:
			fantasyMsg = fantasy.NewSystemMessage(types.GetMessageText(msg))
		case fantasy.MessageRoleUser:
			fantasyMsg = fantasy.NewUserMessage(types.GetMessageText(msg))
		case fantasy.MessageRoleAssistant:
			if types.HasMessageToolCalls(msg) {
				toolCalls := types.GetMessageToolCalls(msg)
				content := make([]fantasy.MessagePart, 0, len(toolCalls)+1)
				if text := types.GetMessageText(msg); text != "" {
					content = append(content, fantasy.TextPart{Text: text})
				}
				for _, tc := range toolCalls {
					content = append(content, fantasy.ToolCallPart{
						ToolCallID: tc.ID,
						ToolName:   tc.Name,
						Input:      tc.Input,
					})
				}
				fantasyMsg = fantasy.Message{
					Role:    fantasy.MessageRoleAssistant,
					Content: content,
				}
			} else {
				fantasyMsg = fantasy.Message{
					Role:    fantasy.MessageRoleAssistant,
					Content: []fantasy.MessagePart{fantasy.TextPart{Text: types.GetMessageText(msg)}},
				}
			}
		case fantasy.MessageRoleTool:
			for _, part := range msg.Content {
				if trp, ok := part.(fantasy.ToolResultPart); ok {
					fantasyMsg = fantasy.Message{
						Role: fantasy.MessageRoleTool,
						Content: []fantasy.MessagePart{
							fantasy.ToolResultPart{
								ToolCallID: trp.ToolCallID,
								Output:     fantasy.ToolResultOutputContentText{Text: types.GetToolResultContent(trp)},
							},
						},
					}
				}
			}
		}

		if len(fantasyMsg.Content) > 0 {
			result = append(result, fantasyMsg)
		}
	}

	return result
}

// getFantasyTools converts tool registry to fantasy tools
func (a *FantasyProviderAdapter) getFantasyTools() []fantasy.Tool {
	if a.toolRegistry == nil {
		return nil
	}

	registeredTools := a.toolRegistry.GetByNames(a.toolNames)
	fantasyTools := make([]fantasy.Tool, 0, len(registeredTools))

	for _, tool := range registeredTools {
		if !tool.Enabled {
			continue
		}

		fantasyTools = append(fantasyTools, fantasy.FunctionTool{
			Name:        tool.Definition.Name,
			Description: tool.Definition.Description,
			InputSchema: tool.Definition.Parameters,
		})
	}

	return fantasyTools
}

// convertResponse converts fantasy.Response to types.LLMResponse
func (a *FantasyProviderAdapter) convertResponse(resp *fantasy.Response) *types.LLMResponse {
	result := &types.LLMResponse{
		Content:      resp.Content.Text(),
		FinishReason: string(resp.FinishReason),
	}

	// Convert tool calls
	toolCalls := resp.Content.ToolCalls()
	if len(toolCalls) > 0 {
		result.ToolCalls = make([]fantasy.ToolCall, len(toolCalls))
		for i, tc := range toolCalls {
			var args map[string]string
			if err := json.Unmarshal([]byte(tc.Input), &args); err != nil {
				var argsAny map[string]interface{}
				if err := json.Unmarshal([]byte(tc.Input), &argsAny); err == nil {
					args = make(map[string]string)
					for k, v := range argsAny {
						args[k] = fmt.Sprintf("%v", v)
					}
				}
			}
			argsJSON, _ := json.Marshal(args)
			result.ToolCalls[i] = fantasy.ToolCall{
				ID:    tc.ToolCallID,
				Name:  tc.ToolName,
				Input: string(argsJSON),
			}
		}
	}

	// Set usage
	result.PromptTokens = int(resp.Usage.InputTokens)
	result.CompletionTokens = int(resp.Usage.OutputTokens)
	result.TotalTokens = int(resp.Usage.TotalTokens)

	return result
}

// consumeStream consumes a stream response and calls the callback
func (a *FantasyProviderAdapter) consumeStream(stream fantasy.StreamResponse, callback types.StreamCallback) *types.LLMResponse {
	var content strings.Builder
	var toolCalls []fantasy.ToolCall
	var finishReason string
	var usage fantasy.Usage
	toolCallMap := make(map[string]*fantasy.ToolCall)

	for part := range stream {
		switch part.Type {
		case fantasy.StreamPartTypeTextDelta:
			content.WriteString(part.Delta)
			if callback != nil {
				callback(part.Delta, false)
			}
		case fantasy.StreamPartTypeTextEnd:
			if callback != nil {
				callback("", false)
			}
		case fantasy.StreamPartTypeToolCall:
			var args map[string]string
			if err := json.Unmarshal([]byte(part.ToolCallInput), &args); err != nil {
				var argsAny map[string]interface{}
				if err := json.Unmarshal([]byte(part.ToolCallInput), &argsAny); err == nil {
					args = make(map[string]string)
					for k, v := range argsAny {
						args[k] = fmt.Sprintf("%v", v)
					}
				}
			}
			argsJSON, _ := json.Marshal(args)
			tc := fantasy.ToolCall{
				ID:    part.ID,
				Name:  part.ToolCallName,
				Input: string(argsJSON),
			}
			toolCallMap[part.ID] = &tc
		case fantasy.StreamPartTypeFinish:
			finishReason = string(part.FinishReason)
			usage = part.Usage
			if callback != nil {
				callback("", true)
			}
		case fantasy.StreamPartTypeError:
			log.Printf("⚠️  Stream error: %v", part.Error)
		}
	}

	for _, tc := range toolCallMap {
		toolCalls = append(toolCalls, *tc)
	}

	return &types.LLMResponse{
		Content:          strings.TrimSpace(content.String()),
		ToolCalls:        toolCalls,
		FinishReason:     finishReason,
		PromptTokens:     int(usage.InputTokens),
		CompletionTokens: int(usage.OutputTokens),
		TotalTokens:      int(usage.TotalTokens),
	}
}

// executeToolCalls executes tool calls and returns tool response messages
func (a *FantasyProviderAdapter) executeToolCalls(ctx context.Context, toolCalls []fantasy.ToolCall) (fantasy.Prompt, error) {
	if a.toolRegistry == nil {
		return nil, fmt.Errorf("tool registry not available")
	}

	toolMessages := make(fantasy.Prompt, 0, len(toolCalls))

	for _, tc := range toolCalls {
		log.Printf("🔧 [fantasy] Executing tool: %s", tc.Name)

		result, err := a.toolRegistry.Execute(ctx, tc.Name, types.GetToolCallArguments(tc))
		if err != nil {
			log.Printf("⚠️  Tool execution failed: %v", err)
			toolMessages = append(toolMessages, types.NewToolErrorMessage(tc.ID, tc.Name, fmt.Sprintf("Error: %v", err)))
		} else {
			displayResult := result
			if len(displayResult) > 100 {
				displayResult = displayResult[:100] + "..."
			}
			log.Printf("✅ Tool result: %s", displayResult)
			toolMessages = append(toolMessages, types.NewToolResultMessage(tc.ID, tc.Name, result))
		}
	}

	return toolMessages, nil
}

// GetModel returns the model ID
func (a *FantasyProviderAdapter) GetModel() string {
	return a.modelID
}

// GetProviderName returns the provider name
func (a *FantasyProviderAdapter) GetProviderName() string {
	return a.model.Provider()
}
