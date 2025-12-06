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
func (a *FantasyProviderAdapter) Generate(ctx context.Context, messages []fantasy.Message) (*fantasy.Response, error) {
	call := a.buildCall(fantasy.Prompt(messages))

	resp, err := a.model.Generate(ctx, call)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// RunAgentLoop runs the agent loop with tool execution
func (a *FantasyProviderAdapter) RunAgentLoop(ctx context.Context, messages fantasy.Prompt, maxIterations int) (*fantasy.Response, fantasy.Prompt, error) {
	if maxIterations <= 0 {
		maxIterations = 10
	}

	allMessages := make(fantasy.Prompt, len(messages))
	copy(allMessages, messages)

	var finalResponse *fantasy.Response
	var allToolMessages fantasy.Prompt
	var allToolCalls []fantasy.ToolCallContent

	for i := 0; i < maxIterations; i++ {
		log.Printf("🔄 [fantasy] Agent loop iteration %d/%d", i+1, maxIterations)

		call := a.buildCall(allMessages)
		resp, err := a.model.Generate(ctx, call)
		if err != nil {
			return nil, allToolMessages, fmt.Errorf("generation failed at iteration %d: %w", i+1, err)
		}

		toolCalls := resp.Content.ToolCalls()

		// Add assistant response to history
		if len(toolCalls) > 0 {
			allMessages = append(allMessages, types.NewToolCallMessageFromContent(toolCalls))
			allToolCalls = append(allToolCalls, toolCalls...)
		} else {
			allMessages = append(allMessages, fantasy.Message{
				Role:    fantasy.MessageRoleAssistant,
				Content: []fantasy.MessagePart{fantasy.TextPart{Text: resp.Content.Text()}},
			})
		}
		finalResponse = resp

		// Check if no tool calls - we're done
		if len(toolCalls) == 0 {
			log.Printf("✅ [fantasy] Agent loop completed after %d iterations", i+1)
			finalResponse = a.appendToolCallsToResponse(finalResponse, allToolCalls)
			return finalResponse, allToolMessages, nil
		}

		// Execute tool calls
		toolMessages, err := a.executeToolCallsFromContent(ctx, toolCalls)
		if err != nil {
			log.Printf("⚠️  Tool execution error: %v", err)
		}

		// Add tool messages to history
		allMessages = append(allMessages, toolMessages...)
		allToolMessages = append(allToolMessages, toolMessages...)
	}

	log.Printf("⚠️  [fantasy] Agent loop reached max iterations (%d)", maxIterations)
	if finalResponse != nil {
		finalResponse = a.appendToolCallsToResponse(finalResponse, allToolCalls)
	}
	return finalResponse, allToolMessages, nil
}

// RunAgentLoopWithStreaming runs the agent loop with streaming callback
func (a *FantasyProviderAdapter) RunAgentLoopWithStreaming(ctx context.Context, messages fantasy.Prompt, maxIterations int, streamCallback types.StreamCallback, toolCallback types.ToolEventCallback) (*fantasy.Response, fantasy.Prompt, error) {
	if maxIterations <= 0 {
		maxIterations = 10
	}

	allMessages := make(fantasy.Prompt, len(messages))
	copy(allMessages, messages)

	var finalResponse *fantasy.Response
	var allToolMessages fantasy.Prompt
	var allToolCalls []fantasy.ToolCallContent

	for i := 0; i < maxIterations; i++ {
		log.Printf("🔄 [fantasy] Agent loop (streaming) iteration %d/%d", i+1, maxIterations)

		call := a.buildCall(allMessages)
		streamResp, err := a.model.Stream(ctx, call)
		if err != nil {
			return nil, allToolMessages, fmt.Errorf("streaming generation failed at iteration %d: %w", i+1, err)
		}

		resp := a.consumeStream(streamResp, streamCallback)
		toolCalls := resp.Content.ToolCalls()

		// Add assistant response to history
		if len(toolCalls) > 0 {
			allMessages = append(allMessages, types.NewToolCallMessageFromContent(toolCalls))
			allToolCalls = append(allToolCalls, toolCalls...)
		} else {
			allMessages = append(allMessages, fantasy.Message{
				Role:    fantasy.MessageRoleAssistant,
				Content: []fantasy.MessagePart{fantasy.TextPart{Text: resp.Content.Text()}},
			})
		}
		finalResponse = resp

		// Check if no tool calls - we're done
		if len(toolCalls) == 0 {
			log.Printf("✅ [fantasy] Agent loop (streaming) completed after %d iterations", i+1)
			finalResponse = a.appendToolCallsToResponse(finalResponse, allToolCalls)
			return finalResponse, allToolMessages, nil
		}

		// Emit tool_call events before execution
		if toolCallback != nil {
			for _, tc := range toolCalls {
				toolCallback(types.ChatEventToolCall, types.ToolCallContentToToolCall(tc), "")
			}
		}

		// Execute tool calls
		toolMessages, err := a.executeToolCallsFromContent(ctx, toolCalls)
		if err != nil {
			log.Printf("⚠️  Tool execution error: %v", err)
		}

		// Emit tool_result events after execution
		if toolCallback != nil {
			for idx, tc := range toolCalls {
				if idx < len(toolMessages) {
					part := toolMessages[idx].Content[0].(fantasy.ToolResultPart)
					toolCallback(types.ChatEventToolResult, types.ToolCallContentToToolCall(tc), types.GetToolResultContent(part))
				}
			}
		}

		// Add tool messages to history
		allMessages = append(allMessages, toolMessages...)
		allToolMessages = append(allToolMessages, toolMessages...)
	}

	log.Printf("⚠️  [fantasy] Agent loop (streaming) reached max iterations (%d)", maxIterations)
	if finalResponse != nil {
		finalResponse = a.appendToolCallsToResponse(finalResponse, allToolCalls)
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

// appendToolCallsToResponse creates a new Response replacing tool calls with the provided ones
// This replaces any existing tool calls in the response with allToolCalls to avoid duplication
func (a *FantasyProviderAdapter) appendToolCallsToResponse(resp *fantasy.Response, allToolCalls []fantasy.ToolCallContent) *fantasy.Response {
	if len(allToolCalls) == 0 {
		return resp
	}

	// Create new content with non-tool-call content from response + all collected tool calls
	newContent := make(fantasy.ResponseContent, 0)

	// Copy non-tool-call content (text, reasoning, etc.)
	for _, c := range resp.Content {
		if c.GetType() != fantasy.ContentTypeToolCall {
			newContent = append(newContent, c)
		}
	}

	// Add all collected tool calls (replacing any existing ones)
	for _, tc := range allToolCalls {
		newContent = append(newContent, tc)
	}

	return &fantasy.Response{
		Content:          newContent,
		FinishReason:     resp.FinishReason,
		Usage:            resp.Usage,
		Warnings:         resp.Warnings,
		ProviderMetadata: resp.ProviderMetadata,
	}
}

// executeToolCallsFromContent executes tool calls from ToolCallContent and returns tool response messages
func (a *FantasyProviderAdapter) executeToolCallsFromContent(ctx context.Context, toolCalls []fantasy.ToolCallContent) (fantasy.Prompt, error) {
	if a.toolRegistry == nil {
		return nil, fmt.Errorf("tool registry not available")
	}

	toolMessages := make(fantasy.Prompt, 0, len(toolCalls))

	for _, tc := range toolCalls {
		log.Printf("🔧 [fantasy] Executing tool: %s (input: %s)", tc.ToolName, tc.Input)

		// Parse arguments from Input JSON
		args := types.GetToolCallArgumentsFromContent(tc)
		log.Printf("🔧 [fantasy] Parsed args: %v", args)
		result, err := a.toolRegistry.Execute(ctx, tc.ToolName, args)
		if err != nil {
			log.Printf("⚠️  Tool execution failed: %v", err)
			toolMessages = append(toolMessages, types.NewToolErrorMessage(tc.ToolCallID, tc.ToolName, fmt.Sprintf("Error: %v", err)))
		} else {
			displayResult := result
			if len(displayResult) > 100 {
				displayResult = displayResult[:100] + "..."
			}
			log.Printf("✅ Tool result: %s", displayResult)
			toolMessages = append(toolMessages, types.NewToolResultMessage(tc.ToolCallID, tc.ToolName, result))
		}
	}

	return toolMessages, nil
}

// consumeStream consumes a stream response and calls the callback
func (a *FantasyProviderAdapter) consumeStream(stream fantasy.StreamResponse, callback types.StreamCallback) *fantasy.Response {
	var content strings.Builder
	var finishReason fantasy.FinishReason
	var usage fantasy.Usage
	toolCallMap := make(map[string]*fantasy.ToolCallContent)

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
			log.Printf("🔧 [fantasy] StreamPartTypeToolCall: ID=%s, Name=%s, Input=%s", part.ID, part.ToolCallName, part.ToolCallInput)
			tc := fantasy.ToolCallContent{
				ToolCallID: part.ID,
				ToolName:   part.ToolCallName,
				Input:      part.ToolCallInput,
			}
			toolCallMap[part.ID] = &tc
		case fantasy.StreamPartTypeFinish:
			finishReason = part.FinishReason
			usage = part.Usage
			if callback != nil {
				callback("", true)
			}
		case fantasy.StreamPartTypeError:
			log.Printf("⚠️  Stream error: %v", part.Error)
		}
	}

	// Build response content
	responseContent := make(fantasy.ResponseContent, 0)
	text := strings.TrimSpace(content.String())
	if text != "" {
		responseContent = append(responseContent, fantasy.TextContent{Text: text})
	}
	for _, tc := range toolCallMap {
		responseContent = append(responseContent, *tc)
	}

	return &fantasy.Response{
		Content:      responseContent,
		FinishReason: finishReason,
		Usage:        usage,
	}
}

// GetModel returns the model ID
func (a *FantasyProviderAdapter) GetModel() string {
	return a.modelID
}

// GetProviderName returns the provider name
func (a *FantasyProviderAdapter) GetProviderName() string {
	return a.model.Provider()
}
