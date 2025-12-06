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

// Package message provides types for creating, manipulating,
// and processing LLM messages in various formats.
package message

import "github.com/kawai-network/veridium/types"

// Role represents the role of a message sender.
type Role string

const (
	RoleSystem    Role = "system"
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleTool      Role = "tool"
)

// ContentType represents the type of content in a message part.
type ContentType string

const (
	ContentTypeText       ContentType = "text"
	ContentTypeReasoning  ContentType = "reasoning"
	ContentTypeFile       ContentType = "file"
	ContentTypeToolCall   ContentType = "tool-call"
	ContentTypeToolResult ContentType = "tool-result"
)

// Part represents a part of a message content.
type Part interface {
	GetType() ContentType
}

// TextPart represents text content in a message.
type TextPart struct {
	Text string `json:"text"`
}

func (t TextPart) GetType() ContentType {
	return ContentTypeText
}

// ReasoningPart represents reasoning/thinking content from the model.
type ReasoningPart struct {
	Text string `json:"text"`
}

func (r ReasoningPart) GetType() ContentType {
	return ContentTypeReasoning
}

// FilePart represents file content in a message.
type FilePart struct {
	Filename  string `json:"filename"`
	Data      []byte `json:"data"`
	MediaType string `json:"media_type"`
}

func (f FilePart) GetType() ContentType {
	return ContentTypeFile
}

// ToolCallPart represents a tool call in a message.
// Uses types.ToolCall for compatibility with existing tool system.
type ToolCallPart struct {
	ToolCall types.ToolCall `json:"tool_call"`
}

func (t ToolCallPart) GetType() ContentType {
	return ContentTypeToolCall
}

// ToolResultPart represents a tool result in a message.
type ToolResultPart struct {
	ToolCallID string `json:"tool_call_id"`
	ToolName   string `json:"tool_name"`
	Content    string `json:"content"`
	IsError    bool   `json:"is_error,omitempty"`
}

func (t ToolResultPart) GetType() ContentType {
	return ContentTypeToolResult
}

// Message represents a message in a conversation.
type Message struct {
	Role    Role   `json:"role"`
	Content []Part `json:"content"`
}

// GetRole returns the role as string (for template compatibility).
func (m Message) GetRole() string {
	return string(m.Role)
}

// GetContent returns the content as a map for template rendering.
// This maintains compatibility with jinja templates that expect map access.
func (m Message) GetContent() map[string]interface{} {
	result := make(map[string]interface{})

	// Extract text content
	var textContent string
	var toolCalls []map[string]interface{}

	for _, part := range m.Content {
		switch p := part.(type) {
		case TextPart:
			textContent += p.Text
		case ReasoningPart:
			// Reasoning is separate, could add to result if needed
		case ToolCallPart:
			tc := p.ToolCall
			toolCalls = append(toolCalls, map[string]interface{}{
				"id":   tc.ID,
				"type": tc.Type,
				"function": map[string]interface{}{
					"name":      tc.Function.Name,
					"arguments": tc.Function.Arguments,
				},
			})
		case ToolResultPart:
			result["name"] = p.ToolName
			result["content"] = p.Content
			result["tool_call_id"] = p.ToolCallID
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

// GetText returns concatenated text content from the message.
func (m Message) GetText() string {
	var text string
	for _, part := range m.Content {
		if p, ok := part.(TextPart); ok {
			text += p.Text
		}
	}
	return text
}

// GetToolCalls returns all tool calls from the message.
func (m Message) GetToolCalls() []types.ToolCall {
	var calls []types.ToolCall
	for _, part := range m.Content {
		if p, ok := part.(ToolCallPart); ok {
			calls = append(calls, p.ToolCall)
		}
	}
	return calls
}

// HasToolCalls returns true if the message contains tool calls.
func (m Message) HasToolCalls() bool {
	for _, part := range m.Content {
		if _, ok := part.(ToolCallPart); ok {
			return true
		}
	}
	return false
}

// Prompt represents a list of messages for the language model.
type Prompt []Message

// NewTextMessage creates a new message with text content.
func NewTextMessage(role Role, text string) Message {
	return Message{
		Role:    role,
		Content: []Part{TextPart{Text: text}},
	}
}

// NewUserMessage creates a new user message with text and optional files.
func NewUserMessage(text string, files ...FilePart) Message {
	content := []Part{TextPart{Text: text}}
	for _, f := range files {
		content = append(content, f)
	}
	return Message{
		Role:    RoleUser,
		Content: content,
	}
}

// NewSystemMessage creates a new system message.
func NewSystemMessage(text string) Message {
	return Message{
		Role:    RoleSystem,
		Content: []Part{TextPart{Text: text}},
	}
}

// NewAssistantMessage creates a new assistant message with text.
func NewAssistantMessage(text string) Message {
	return Message{
		Role:    RoleAssistant,
		Content: []Part{TextPart{Text: text}},
	}
}

// NewToolCallMessage creates a new assistant message with tool calls.
func NewToolCallMessage(toolCalls []types.ToolCall) Message {
	content := make([]Part, len(toolCalls))
	for i, tc := range toolCalls {
		content[i] = ToolCallPart{ToolCall: tc}
	}
	return Message{
		Role:    RoleAssistant,
		Content: content,
	}
}

// NewToolResultMessage creates a new tool result message.
func NewToolResultMessage(toolCallID, toolName, result string) Message {
	return Message{
		Role: RoleTool,
		Content: []Part{
			ToolResultPart{
				ToolCallID: toolCallID,
				ToolName:   toolName,
				Content:    result,
			},
		},
	}
}

// NewToolErrorMessage creates a new tool error message.
func NewToolErrorMessage(toolCallID, toolName, errorMsg string) Message {
	return Message{
		Role: RoleTool,
		Content: []Part{
			ToolResultPart{
				ToolCallID: toolCallID,
				ToolName:   toolName,
				Content:    errorMsg,
				IsError:    true,
			},
		},
	}
}
