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
	"github.com/kawai-network/veridium/pkg/fantasy"
)

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
