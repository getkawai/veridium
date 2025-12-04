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

// ============================================================================
// Tool Types
// ============================================================================

// ToolFunction represents a function within a tool call.
type ToolFunction struct {
	Name      string
	Arguments map[string]string
}

// ToolCall represents a call to a tool function within a tool message.
type ToolCall struct {
	Type     string
	Function ToolFunction
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
