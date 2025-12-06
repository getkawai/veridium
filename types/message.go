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

// Package types provides shared types and helper functions for message handling.
// For core message types (Message, MessagePart, TextPart, etc.), import fantasy package directly.
package types

import (
	"encoding/json"
	"errors"

	"github.com/kawai-network/veridium/fantasy"
)

// ============================================================================
// Message Helper Functions
// ============================================================================

// NewToolCallMessage creates a new assistant message with tool calls.
func NewToolCallMessage(toolCalls []fantasy.ToolCall) fantasy.Message {
	content := make([]fantasy.MessagePart, len(toolCalls))
	for i, tc := range toolCalls {
		content[i] = fantasy.ToolCallPart{
			ToolCallID: tc.ID,
			ToolName:   tc.Name,
			Input:      tc.Input,
		}
	}
	return fantasy.Message{
		Role:    fantasy.MessageRoleAssistant,
		Content: content,
	}
}

// NewToolResultMessage creates a new tool result message.
func NewToolResultMessage(toolCallID, toolName, result string) fantasy.Message {
	return fantasy.Message{
		Role: fantasy.MessageRoleTool,
		Content: []fantasy.MessagePart{
			fantasy.ToolResultPart{
				ToolCallID: toolCallID,
				Output:     fantasy.ToolResultOutputContentText{Text: result},
			},
		},
	}
}

// NewToolErrorMessage creates a new tool error message.
func NewToolErrorMessage(toolCallID, toolName, errorMsg string) fantasy.Message {
	return fantasy.Message{
		Role: fantasy.MessageRoleTool,
		Content: []fantasy.MessagePart{
			fantasy.ToolResultPart{
				ToolCallID: toolCallID,
				Output:     fantasy.ToolResultOutputContentError{Error: errors.New(errorMsg)},
			},
		},
	}
}

// ============================================================================
// Message Accessor Helper Functions
// ============================================================================

// GetMessageText returns concatenated text content from the message.
func GetMessageText(m fantasy.Message) string {
	var text string
	for _, part := range m.Content {
		if p, ok := fantasy.AsMessagePart[fantasy.TextPart](part); ok {
			text += p.Text
		}
	}
	return text
}

// GetMessageToolCalls returns all tool calls from the message.
func GetMessageToolCalls(m fantasy.Message) []fantasy.ToolCall {
	var calls []fantasy.ToolCall
	for _, part := range m.Content {
		if p, ok := fantasy.AsMessagePart[fantasy.ToolCallPart](part); ok {
			calls = append(calls, fantasy.ToolCall{
				ID:    p.ToolCallID,
				Name:  p.ToolName,
				Input: p.Input,
			})
		}
	}
	return calls
}

// HasMessageToolCalls returns true if the message contains tool calls.
func HasMessageToolCalls(m fantasy.Message) bool {
	for _, part := range m.Content {
		if _, ok := fantasy.AsMessagePart[fantasy.ToolCallPart](part); ok {
			return true
		}
	}
	return false
}

// GetMessageRole returns the role as string (for template compatibility).
func GetMessageRole(m fantasy.Message) string {
	return string(m.Role)
}

// GetMessageContent returns the content as a map for template rendering.
// This maintains compatibility with jinja templates that expect map access.
func GetMessageContent(m fantasy.Message) map[string]interface{} {
	result := make(map[string]interface{})

	var textContent string
	var toolCalls []map[string]interface{}

	for _, part := range m.Content {
		switch p := part.(type) {
		case fantasy.TextPart:
			textContent += p.Text
		case fantasy.ToolCallPart:
			// Parse Input JSON string to arguments map for template compatibility
			var arguments interface{}
			if p.Input != "" {
				var parsed map[string]interface{}
				if err := json.Unmarshal([]byte(p.Input), &parsed); err == nil {
					arguments = parsed
				}
			}
			toolCalls = append(toolCalls, map[string]interface{}{
				"id":        p.ToolCallID,
				"name":      p.ToolName,
				"input":     p.Input,
				"arguments": arguments,
			})
		case fantasy.ToolResultPart:
			result["tool_call_id"] = p.ToolCallID
			if textOutput, ok := p.Output.(fantasy.ToolResultOutputContentText); ok {
				result["content"] = textOutput.Text
			}
		}
	}

	if textContent != "" {
		result["content"] = textContent
	}
	if len(toolCalls) > 0 {
		result["tool_calls"] = toolCalls
	}

	return result
}

// GetToolResultContent extracts the text content from a ToolResultPart.
func GetToolResultContent(trp fantasy.ToolResultPart) string {
	if textOutput, ok := trp.Output.(fantasy.ToolResultOutputContentText); ok {
		return textOutput.Text
	}
	if errOutput, ok := trp.Output.(fantasy.ToolResultOutputContentError); ok {
		if errOutput.Error != nil {
			return errOutput.Error.Error()
		}
	}
	return ""
}
