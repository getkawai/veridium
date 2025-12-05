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
// OpenAI-Compatible API Types
// ============================================================================
// These types are used for OpenAI-compatible APIs (OpenRouter, Zhipu GLM, etc.)

// ChatCompletionRequest represents a chat completion API request
type ChatCompletionRequest struct {
	Model       string                `json:"model"`
	Messages    []ChatCompletionMsg   `json:"messages"`
	Temperature *float32              `json:"temperature,omitempty"`
	TopP        *float32              `json:"top_p,omitempty"`
	MaxTokens   *int                  `json:"max_tokens,omitempty"`
	Stream      bool                  `json:"stream,omitempty"`
	Stop        []string              `json:"stop,omitempty"`
	Tools       []APIToolDefinition   `json:"tools,omitempty"`
	ToolChoice  interface{}           `json:"tool_choice,omitempty"` // "auto", "none", or specific tool
}

// ChatCompletionMsg represents a message in API format
type ChatCompletionMsg struct {
	Role       string        `json:"role"` // system, user, assistant, tool
	Content    interface{}   `json:"content,omitempty"` // string or []ContentPart for multimodal
	Name       string        `json:"name,omitempty"`        // For tool messages
	ToolCalls  []APIToolCall `json:"tool_calls,omitempty"`  // For assistant messages with tool calls
	ToolCallID string        `json:"tool_call_id,omitempty"` // For tool response messages
}

// ContentPart represents a part of multimodal content
type ContentPart struct {
	Type     string    `json:"type"` // "text", "image_url", "video_url"
	Text     string    `json:"text,omitempty"`
	ImageURL *MediaURL `json:"image_url,omitempty"`
	VideoURL *MediaURL `json:"video_url,omitempty"`
}

// MediaURL represents an image or video URL for multimodal content
type MediaURL struct {
	URL    string `json:"url"` // Can be URL or data:mime;base64,... for inline
	Detail string `json:"detail,omitempty"` // "auto", "low", "high" for images
}

// APIToolDefinition defines a tool for the API
type APIToolDefinition struct {
	Type     string          `json:"type"` // "function"
	Function APIToolFunction `json:"function"`
}

// APIToolFunction defines the function schema for API
type APIToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters,omitempty"`
}

// APIToolCall represents a tool call in API format (arguments as JSON string)
type APIToolCall struct {
	ID       string              `json:"id"`
	Type     string              `json:"type"` // "function"
	Function APIToolCallFunction `json:"function"`
}

// APIToolCallFunction contains the function call details in API format
type APIToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON string
}

// ============================================================================
// API Response Types
// ============================================================================

// ChatCompletionResponse represents a chat completion API response
type ChatCompletionResponse struct {
	ID                string                 `json:"id"`
	Object            string                 `json:"object"`
	Created           int64                  `json:"created"`
	Model             string                 `json:"model"`
	Choices           []ChatCompletionChoice `json:"choices"`
	Usage             *APIUsage              `json:"usage,omitempty"`
	SystemFingerprint string                 `json:"system_fingerprint,omitempty"`
}

// ChatCompletionChoice represents a choice in the response
type ChatCompletionChoice struct {
	Index        int               `json:"index"`
	Message      ChatCompletionMsg `json:"message,omitempty"` // For non-streaming
	Delta        ChatCompletionMsg `json:"delta,omitempty"`   // For streaming
	FinishReason string            `json:"finish_reason,omitempty"`
}

// APIUsage represents token usage statistics from API
type APIUsage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionStreamResponse represents a streaming response chunk
type ChatCompletionStreamResponse struct {
	ID                string                 `json:"id"`
	Object            string                 `json:"object"`
	Created           int64                  `json:"created"`
	Model             string                 `json:"model"`
	Choices           []ChatCompletionChoice `json:"choices"`
	Usage             *APIUsage              `json:"usage,omitempty"` // Some providers return this in final chunk
	SystemFingerprint string                 `json:"system_fingerprint,omitempty"`
}

// ============================================================================
// API Error Types
// ============================================================================

// APIError represents an API error response
type APIError struct {
	Error *APIErrorDetail `json:"error"`
}

// APIErrorDetail contains error details
type APIErrorDetail struct {
	Message string `json:"message"`
	Type    string `json:"type"`
	Param   string `json:"param,omitempty"`
	Code    string `json:"code,omitempty"`
}
