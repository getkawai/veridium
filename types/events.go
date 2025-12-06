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

package types

import (
	"context"
	"encoding/json"

	"github.com/kawai-network/veridium/fantasy"
)

// ToolCall is an alias to fantasy.ToolCall for backward compatibility
type ToolCall = fantasy.ToolCall

// GetToolCallArguments parses Input JSON string to map[string]string
// Helper function since fantasy.ToolCall doesn't have this method
func GetToolCallArguments(tc fantasy.ToolCall) map[string]string {
	if tc.Input == "" {
		return nil
	}
	var args map[string]string
	if err := json.Unmarshal([]byte(tc.Input), &args); err != nil {
		return nil
	}
	return args
}

// ============================================================================
// Tool Types
// ============================================================================

// ToolExecutor is a function that executes a tool
type ToolExecutor func(ctx context.Context, args map[string]string) (string, error)

// ToolDefinition represents a tool definition for registration
type ToolDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"` // JSON Schema
}

// Tool represents a tool (definition + executor)
type Tool struct {
	Type       fantasy.ToolType `json:"type"`
	Definition ToolDefinition   `json:"definition"`
	Executor   ToolExecutor     `json:"-"`
	Enabled    bool             `json:"-"`
}



// ============================================================================
// LLM Response Types
// ============================================================================

// LLMResponse represents a response from any LLM provider
type LLMResponse struct {
	Content          string     `json:"content"`
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`
	FinishReason     string     `json:"finish_reason"`
	PromptTokens     int        `json:"prompt_tokens"`
	CompletionTokens int        `json:"completion_tokens"`
	TotalTokens      int        `json:"total_tokens"`
	ReasoningContent string     `json:"reasoning_content,omitempty"` // For reasoning models (Qwen3, DeepSeek R1)
}

// ============================================================================
// Event Types
// ============================================================================

// ChatStreamEvent represents the type of event during LLM interaction
// This is used for both streaming events and tool execution events
type ChatStreamEvent string

const (
	// Streaming events
	ChatEventStart     ChatStreamEvent = "start"     // Generation started
	ChatEventChunk     ChatStreamEvent = "chunk"     // Content chunk
	ChatEventReasoning ChatStreamEvent = "reasoning" // Reasoning content
	ChatEventComplete  ChatStreamEvent = "complete"  // Generation complete

	// Tool events
	ChatEventToolCall   ChatStreamEvent = "tool_call"   // Tool call initiated (before execution)
	ChatEventToolResult ChatStreamEvent = "tool_result" // Tool execution result (after execution)
)

// ============================================================================
// Callback Types
// ============================================================================

// StreamCallback is called for each generated token during streaming
type StreamCallback func(token string, isLast bool)

// ToolEventCallback is called when tool events occur during agent loop
// eventType: ChatEventToolCall (before execution) or ChatEventToolResult (after execution)
// toolCall: the tool call being processed
// result: tool execution result (only for ChatEventToolResult event)
type ToolEventCallback func(eventType ChatStreamEvent, toolCall ToolCall, result string)
