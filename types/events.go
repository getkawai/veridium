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
	"encoding/json"
	"fmt"

	"github.com/kawai-network/veridium/fantasy"
)

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

// GetToolCallArgumentsFromContent parses Input JSON string from ToolCallContent to map[string]string
func GetToolCallArgumentsFromContent(tc fantasy.ToolCallContent) map[string]string {
	if tc.Input == "" {
		return nil
	}
	// First try direct map[string]string
	var args map[string]string
	if err := json.Unmarshal([]byte(tc.Input), &args); err == nil {
		return args
	}
	// Fallback: parse as map[string]interface{} and convert
	var argsAny map[string]interface{}
	if err := json.Unmarshal([]byte(tc.Input), &argsAny); err != nil {
		return nil
	}
	args = make(map[string]string)
	for k, v := range argsAny {
		switch val := v.(type) {
		case string:
			args[k] = val
		case float64:
			args[k] = fmt.Sprintf("%v", val)
		case bool:
			args[k] = fmt.Sprintf("%v", val)
		default:
			// For complex types, marshal back to JSON
			if jsonBytes, err := json.Marshal(v); err == nil {
				args[k] = string(jsonBytes)
			} else {
				args[k] = fmt.Sprintf("%v", v)
			}
		}
	}
	return args
}

// ToolCallContentToToolCall converts fantasy.ToolCallContent to fantasy.ToolCall
func ToolCallContentToToolCall(tc fantasy.ToolCallContent) fantasy.ToolCall {
	return fantasy.ToolCall{
		ID:    tc.ToolCallID,
		Name:  tc.ToolName,
		Input: tc.Input,
	}
}

// NewToolCallMessageFromContent creates an assistant message with tool calls from ToolCallContent
func NewToolCallMessageFromContent(toolCalls []fantasy.ToolCallContent) fantasy.Message {
	content := make([]fantasy.MessagePart, 0, len(toolCalls))
	for _, tc := range toolCalls {
		content = append(content, fantasy.ToolCallPart{
			ToolCallID: tc.ToolCallID,
			ToolName:   tc.ToolName,
			Input:      tc.Input,
		})
	}
	return fantasy.Message{
		Role:    fantasy.MessageRoleAssistant,
		Content: content,
	}
}

// ============================================================================
// Event Types
// ============================================================================
// Note: Tool types moved to pkg/yzma/tools - use fantasy.AgentTool interface

// ChatStreamEvent represents the type of event during LLM interaction
// This is used for both streaming events and tool execution events
type ChatStreamEvent string

const (
	// Streaming events
	ChatEventStart     ChatStreamEvent = "start"     // Generation started
	ChatEventChunk     ChatStreamEvent = "chunk"     // Content chunk
	ChatEventReasoning ChatStreamEvent = "reasoning" // Reasoning content delta
	ChatEventComplete  ChatStreamEvent = "complete"  // Generation complete

	// Reasoning events (for models like DeepSeek R1, o1, etc.)
	ChatEventReasoningStart ChatStreamEvent = "reasoning_start" // Reasoning started
	ChatEventReasoningEnd   ChatStreamEvent = "reasoning_end"   // Reasoning finished

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
type ToolEventCallback func(eventType ChatStreamEvent, toolCall fantasy.ToolCall, result string)
