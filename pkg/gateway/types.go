// Package gateway provides OpenAI-compatible API server for worker nodes.
package gateway

import "encoding/json"

// ChatCompletionRequest represents an OpenAI Chat Completion request.
// Ref: https://platform.openai.com/docs/api-reference/chat/create
type ChatCompletionRequest struct {
	Model            string          `json:"model,omitempty"`
	Messages         []ChatMessage   `json:"messages"`
	MaxTokens        int             `json:"max_tokens,omitempty"`
	Temperature      *float64        `json:"temperature,omitempty"`
	TopP             *float64        `json:"top_p,omitempty"`
	N                int             `json:"n,omitempty"`
	Stream           bool            `json:"stream,omitempty"`
	Stop             json.RawMessage `json:"stop,omitempty"`
	PresencePenalty  *float64        `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64        `json:"frequency_penalty,omitempty"`
	User             string          `json:"user,omitempty"`
}

// ChatMessage represents a message in a chat conversation.
type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

// ChatCompletionResponse represents an OpenAI Chat Completion response.
type ChatCompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             *Usage   `json:"usage,omitempty"`
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
}

// Choice represents a single completion choice.
type Choice struct {
	Index        int          `json:"index"`
	Message      *ChatMessage `json:"message,omitempty"`
	Delta        *ChatMessage `json:"delta,omitempty"`
	FinishReason string       `json:"finish_reason,omitempty"`
}

// Usage represents token usage information.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionChunk represents a streaming chunk response.
type ChatCompletionChunk struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
}

// ErrorResponse represents an OpenAI API error response.
type ErrorResponse struct {
	Error ErrorDetail `json:"error"`
}

// ErrorDetail contains error details.
type ErrorDetail struct {
	Message string  `json:"message"`
	Type    string  `json:"type"`
	Param   *string `json:"param,omitempty"`
	Code    *string `json:"code,omitempty"`
}
