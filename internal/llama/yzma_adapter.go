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

	"github.com/cloudwego/eino/schema"
	"github.com/kawai-network/veridium/pkg/yzma/llama"
	"github.com/kawai-network/veridium/pkg/yzma/message"
	"github.com/kawai-network/veridium/pkg/yzma/template"
	"github.com/kawai-network/veridium/pkg/yzma/tools"
)

// LlamaYzmaModel implements chat model with yzma tool calling
type LlamaYzmaModel struct {
	libService   *LibraryService
	toolRegistry *tools.ToolRegistry
	toolNames    []string // Requested tool names
}

// NewLlamaYzmaModel creates a new yzma-based model adapter
func NewLlamaYzmaModel(libService *LibraryService, toolRegistry *tools.ToolRegistry) *LlamaYzmaModel {
	return &LlamaYzmaModel{
		libService:   libService,
		toolRegistry: toolRegistry,
		toolNames:    []string{}, // Empty = all tools
	}
}

// WithTools sets the tools to use (empty = all enabled tools)
func (m *LlamaYzmaModel) WithTools(toolNames []string) *LlamaYzmaModel {
	return &LlamaYzmaModel{
		libService:   m.libService,
		toolRegistry: m.toolRegistry,
		toolNames:    toolNames,
	}
}

// GenerateWithTools generates a response and handles tool calls
func (m *LlamaYzmaModel) GenerateWithTools(ctx context.Context, messages []*schema.Message) (*schema.Message, []schema.ToolCall, error) {
	// Ensure chat model is loaded
	if !m.libService.IsChatModelLoaded() {
		if err := m.libService.LoadChatModel(""); err != nil {
			return nil, nil, fmt.Errorf("failed to load chat model: %w", err)
		}
	}

	// Convert Eino messages to yzma messages
	yzmaMessages, err := m.convertToYzmaMessages(messages)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to convert messages: %w", err)
	}

	// Enhance system prompt with tools if available
	if m.toolRegistry != nil && len(yzmaMessages) > 0 {
		toolsJSON, err := m.toolRegistry.FormatForPrompt(m.toolNames)
		if err == nil && toolsJSON != "" {
			// Find or create system message
			if yzmaMessages[0].GetRole() == "system" {
				// Enhance existing system message
				systemContent := yzmaMessages[0].(message.Chat).Content
				enhancedContent := tools.BuildSystemPrompt(systemContent, toolsJSON)
				yzmaMessages[0] = message.Chat{
					Role:    "system",
					Content: enhancedContent,
				}
			} else {
				// Prepend new system message with tools
				systemContent := tools.BuildSystemPrompt("You are a helpful AI assistant.", toolsJSON)
				yzmaMessages = append([]message.Message{
					message.Chat{Role: "system", Content: systemContent},
				}, yzmaMessages...)
			}
			
			log.Printf("🔧 Enhanced system prompt with %d tools", len(m.toolRegistry.GetByNames(m.toolNames)))
		}
	}

	// Apply chat template
	chatTemplate := llama.ModelChatTemplate(m.libService.chatModel, "")
	prompt, err := template.Apply(chatTemplate, yzmaMessages, true)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to apply template: %w", err)
	}

	// Generate response
	response, err := m.libService.Generate(prompt, 32768)
	if err != nil {
		return nil, nil, fmt.Errorf("generation failed: %w", err)
	}

	// Parse tool calls from response
	var toolCalls []schema.ToolCall
	yzmaToolCalls := tools.ParseToolCalls(response)
	
	if len(yzmaToolCalls) > 0 {
		log.Printf("🔧 Detected %d tool calls in response", len(yzmaToolCalls))
		
		// Convert yzma tool calls to Eino format
		toolCalls = make([]schema.ToolCall, len(yzmaToolCalls))
		for i, ytc := range yzmaToolCalls {
			// Convert arguments map to JSON string
			argsJSON := "{"
			first := true
			for k, v := range ytc.Function.Arguments {
				if !first {
					argsJSON += ","
				}
				argsJSON += fmt.Sprintf(`"%s":"%s"`, k, v)
				first = false
			}
			argsJSON += "}"
			
			toolCalls[i] = schema.ToolCall{
				ID:   fmt.Sprintf("call_%d", i),
				Type: "function",
				Function: schema.FunctionCall{
					Name:      ytc.Function.Name,
					Arguments: argsJSON,
				},
			}
			
			log.Printf("🔧 Tool call #%d: %s(%s)", i+1, ytc.Function.Name, argsJSON)
		}
		
		// Remove tool call tags from response for cleaner output
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
		response = strings.TrimSpace(cleanResponse)
	}

	// Create response message
	responseMsg := &schema.Message{
		Role:      schema.Assistant,
		Content:   response,
		ToolCalls: toolCalls,
		ResponseMeta: &schema.ResponseMeta{
			FinishReason: "stop",
			Usage: &schema.TokenUsage{
				PromptTokens:     len(prompt) / 4,
				CompletionTokens: len(response) / 4,
				TotalTokens:      (len(prompt) + len(response)) / 4,
			},
		},
	}

	return responseMsg, toolCalls, nil
}

// ExecuteToolCalls executes tool calls and returns tool messages
func (m *LlamaYzmaModel) ExecuteToolCalls(ctx context.Context, toolCalls []schema.ToolCall) ([]*schema.Message, error) {
	if m.toolRegistry == nil {
		return nil, fmt.Errorf("tool registry not available")
	}

	toolMessages := make([]*schema.Message, 0, len(toolCalls))
	
	for _, tc := range toolCalls {
		log.Printf("🔧 Executing tool: %s", tc.Function.Name)
		
		// Parse arguments from JSON string to map
		args := make(map[string]string)
		argsStr := strings.Trim(tc.Function.Arguments, "{}")
		if argsStr != "" {
			pairs := strings.Split(argsStr, ",")
			for _, pair := range pairs {
				parts := strings.SplitN(pair, ":", 2)
				if len(parts) == 2 {
					key := strings.Trim(parts[0], `"`)
					value := strings.Trim(parts[1], `"`)
					args[key] = value
				}
			}
		}
		
		// Execute tool
		result, err := m.toolRegistry.Execute(ctx, tc.Function.Name, args)
		if err != nil {
			log.Printf("⚠️  Tool execution failed: %v", err)
			result = fmt.Sprintf("Error: %v", err)
		} else {
			log.Printf("✅ Tool result: %s", result[:minInt(100, len(result))])
		}
		
		// Create tool message
		toolMsg := &schema.Message{
			Role:       schema.Tool,
			Content:    result,
			ToolCallID: tc.ID,
			ToolName:   tc.Function.Name,
		}
		
		toolMessages = append(toolMessages, toolMsg)
	}
	
	return toolMessages, nil
}

// convertToYzmaMessages converts Eino messages to yzma messages
func (m *LlamaYzmaModel) convertToYzmaMessages(messages []*schema.Message) ([]message.Message, error) {
	yzmaMessages := make([]message.Message, 0, len(messages))
	
	for _, msg := range messages {
		switch msg.Role {
		case schema.System:
			yzmaMessages = append(yzmaMessages, message.Chat{
				Role:    "system",
				Content: msg.Content,
			})
			
		case schema.User:
			yzmaMessages = append(yzmaMessages, message.Chat{
				Role:    "user",
				Content: msg.Content,
			})
			
		case schema.Assistant:
			if len(msg.ToolCalls) > 0 {
				// Convert tool calls to yzma format
				yzmaToolCalls := make([]message.ToolCall, len(msg.ToolCalls))
				for i, tc := range msg.ToolCalls {
					// Parse arguments
					args := make(map[string]string)
					argsStr := strings.Trim(tc.Function.Arguments, "{}")
					if argsStr != "" {
						pairs := strings.Split(argsStr, ",")
						for _, pair := range pairs {
							parts := strings.SplitN(pair, ":", 2)
							if len(parts) == 2 {
								key := strings.Trim(parts[0], `"`)
								value := strings.Trim(parts[1], `"`)
								args[key] = value
							}
						}
					}
					
					yzmaToolCalls[i] = message.ToolCall{
						Type: "function",
						Function: message.ToolFunction{
							Name:      tc.Function.Name,
							Arguments: args,
						},
					}
				}
				
				yzmaMessages = append(yzmaMessages, message.Tool{
					Role:      "assistant",
					ToolCalls: yzmaToolCalls,
				})
			} else {
				yzmaMessages = append(yzmaMessages, message.Chat{
					Role:    "assistant",
					Content: msg.Content,
				})
			}
			
		case schema.Tool:
			yzmaMessages = append(yzmaMessages, message.ToolResponse{
				Role:    "tool",
				Name:    msg.ToolName,
				Content: msg.Content,
			})
		}
	}
	
	return yzmaMessages, nil
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

