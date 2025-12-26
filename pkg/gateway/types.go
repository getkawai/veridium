// Package gateway provides OpenAI-compatible API server for contributor nodes.
package gateway

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
)

// ChatCompletionRequest represents an OpenAI Chat Completion request.
// Ref: https://platform.openai.com/docs/api-reference/chat/create
type ChatCompletionRequest struct {
	Model               string          `json:"model,omitempty"`
	Messages            []ChatMessage   `json:"messages"`
	MaxTokens           int             `json:"max_tokens,omitempty"`
	MaxCompletionTokens int             `json:"max_completion_tokens,omitempty"`
	Temperature         *float64        `json:"temperature,omitempty"`
	TopP                *float64        `json:"top_p,omitempty"`
	N                   int             `json:"n,omitempty"`
	Stream              bool            `json:"stream,omitempty"`
	Stop                json.RawMessage `json:"stop,omitempty"`
	PresencePenalty     *float64        `json:"presence_penalty,omitempty"`
	FrequencyPenalty    *float64        `json:"frequency_penalty,omitempty"`
	User                string          `json:"user,omitempty"`
	Tools               []Tool          `json:"tools,omitempty"`
	ToolChoice          json.RawMessage `json:"tool_choice,omitempty"`
	ResponseFormat      *ResponseFormat `json:"response_format,omitempty"`
	ReasoningEffort     string          `json:"reasoning_effort,omitempty"`
	Seed                *int64          `json:"seed,omitempty"`
	StreamOptions       *StreamOptions  `json:"stream_options,omitempty"`
}

// GetMaxTokens returns the effective max tokens value.
func (r *ChatCompletionRequest) GetMaxTokens() int {
	if r.MaxCompletionTokens > 0 {
		return r.MaxCompletionTokens
	}
	if r.MaxTokens > 0 {
		return r.MaxTokens
	}
	return 2048 // Default
}

// StreamOptions configures streaming behavior.
type StreamOptions struct {
	IncludeUsage bool `json:"include_usage,omitempty"`
}

// ResponseFormat specifies the format of the response.
type ResponseFormat struct {
	Type string `json:"type,omitempty"` // "text" or "json_object"
}

// Tool represents a tool available to the model.
type Tool struct {
	Type     string       `json:"type"` // "function"
	Function ToolFunction `json:"function"`
}

// ToolFunction describes a function tool.
type ToolFunction struct {
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Parameters  map[string]any `json:"parameters,omitempty"`
}

// ChatMessage represents a message in a chat conversation.
// Content can be either a string or an array of content parts.
type ChatMessage struct {
	Role             string         `json:"role"`
	Content          MessageContent `json:"content,omitempty"`
	Name             string         `json:"name,omitempty"`
	ToolCalls        []ToolCall     `json:"tool_calls,omitempty"`
	ToolCallID       string         `json:"tool_call_id,omitempty"`
	ReasoningContent string         `json:"reasoning_content,omitempty"`
}

// ToolCall represents a tool call made by the assistant.
type ToolCall struct {
	ID       string           `json:"id"`
	Type     string           `json:"type"` // "function"
	Function ToolCallFunction `json:"function"`
}

// ToolCallFunction contains the function call details.
type ToolCallFunction struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// MessageContent can be either a string or array of content parts.
type MessageContent struct {
	Text  string        // Simple text content
	Parts []ContentPart // Array of content parts (for multimodal)
}

// MarshalJSON implements json.Marshaler for MessageContent.
func (m MessageContent) MarshalJSON() ([]byte, error) {
	if len(m.Parts) > 0 {
		return json.Marshal(m.Parts)
	}
	return json.Marshal(m.Text)
}

// UnmarshalJSON implements json.Unmarshaler for MessageContent.
func (m *MessageContent) UnmarshalJSON(data []byte) error {
	// Try string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		m.Text = s
		m.Parts = nil
		return nil
	}

	// Try array of content parts
	var parts []ContentPart
	if err := json.Unmarshal(data, &parts); err == nil {
		m.Parts = parts
		m.Text = ""
		return nil
	}

	return fmt.Errorf("content must be string or array of content parts")
}

// GetText returns the text content, extracting from parts if necessary.
func (m MessageContent) GetText() string {
	if m.Text != "" {
		return m.Text
	}
	var texts []string
	for _, part := range m.Parts {
		if part.Type == "text" && part.Text != "" {
			texts = append(texts, part.Text)
		}
	}
	return strings.Join(texts, "\n")
}

// HasImages returns true if the content contains image parts.
func (m MessageContent) HasImages() bool {
	for _, part := range m.Parts {
		if part.Type == "image_url" && part.ImageURL != nil {
			return true
		}
	}
	return false
}

// MessageImageData holds extracted image information from message content.
type MessageImageData struct {
	URL       string
	MediaType string
	Data      []byte
}

// GetImages returns all image data from content parts.
func (m MessageContent) GetImages() []MessageImageData {
	var images []MessageImageData
	for _, part := range m.Parts {
		if part.Type == "image_url" && part.ImageURL != nil {
			img := MessageImageData{URL: part.ImageURL.URL}
			// Try to extract base64 data if it's a data URL
			if strings.HasPrefix(part.ImageURL.URL, "data:") {
				parts := strings.SplitN(part.ImageURL.URL, ",", 2)
				if len(parts) == 2 {
					// Extract media type from "data:image/png;base64"
					mediaInfo := strings.TrimPrefix(parts[0], "data:")
					mediaInfo = strings.TrimSuffix(mediaInfo, ";base64")
					img.MediaType = mediaInfo
					// Decode base64
					if decoded, err := base64.StdEncoding.DecodeString(parts[1]); err == nil {
						img.Data = decoded
					}
				}
			}
			images = append(images, img)
		}
	}
	return images
}

// ContentPart represents a part of message content.
type ContentPart struct {
	Type       string        `json:"type"` // "text", "image_url", "input_audio", "file"
	Text       string        `json:"text,omitempty"`
	ImageURL   *ImageURLPart `json:"image_url,omitempty"`
	InputAudio *AudioPart    `json:"input_audio,omitempty"`
	File       *FilePart     `json:"file,omitempty"`
}

// ImageURLPart represents an image URL in a content part.
type ImageURLPart struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"` // "auto", "low", "high"
}

// AudioPart represents audio content in a message.
type AudioPart struct {
	Data   string `json:"data"`   // base64 encoded
	Format string `json:"format"` // "wav", "mp3"
}

// FilePart represents file content in a message.
type FilePart struct {
	FileID   string `json:"file_id,omitempty"`
	FileData string `json:"file_data,omitempty"` // data URL
	Filename string `json:"filename,omitempty"`
}

// ResponseMessage represents the message in a response.
type ResponseMessage struct {
	Role             string     `json:"role"`
	Content          string     `json:"content"`
	ReasoningContent string     `json:"reasoning_content,omitempty"`
	ToolCalls        []ToolCall `json:"tool_calls,omitempty"`
	Refusal          *string    `json:"refusal,omitempty"`
}

// ChatCompletionResponse represents an OpenAI Chat Completion response.
type ChatCompletionResponse struct {
	ID                string           `json:"id"`
	Object            string           `json:"object"`
	Created           int64            `json:"created"`
	Model             string           `json:"model"`
	Choices           []ResponseChoice `json:"choices"`
	Usage             *Usage           `json:"usage,omitempty"`
	SystemFingerprint string           `json:"system_fingerprint,omitempty"`
}

// ResponseChoice represents a single completion choice in response.
type ResponseChoice struct {
	Index        int              `json:"index"`
	Message      *ResponseMessage `json:"message,omitempty"`
	FinishReason string           `json:"finish_reason,omitempty"`
}

// DeltaMessage represents a delta update in streaming.
type DeltaMessage struct {
	Role             string          `json:"role,omitempty"`
	Content          string          `json:"content,omitempty"`
	ReasoningContent string          `json:"reasoning_content,omitempty"`
	ToolCalls        []ToolCallDelta `json:"tool_calls,omitempty"`
}

// ToolCallDelta represents a partial tool call in streaming.
type ToolCallDelta struct {
	Index    int                    `json:"index"`
	ID       string                 `json:"id,omitempty"`
	Type     string                 `json:"type,omitempty"`
	Function *ToolCallFunctionDelta `json:"function,omitempty"`
}

// ToolCallFunctionDelta represents partial function call data.
type ToolCallFunctionDelta struct {
	Name      string `json:"name,omitempty"`
	Arguments string `json:"arguments,omitempty"`
}

// StreamChoice represents a choice in streaming response.
type StreamChoice struct {
	Index        int           `json:"index"`
	Delta        *DeltaMessage `json:"delta,omitempty"`
	FinishReason string        `json:"finish_reason,omitempty"`
}

// Usage represents token usage information.
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ChatCompletionChunk represents a streaming chunk response.
type ChatCompletionChunk struct {
	ID                string         `json:"id"`
	Object            string         `json:"object"`
	Created           int64          `json:"created"`
	Model             string         `json:"model"`
	Choices           []StreamChoice `json:"choices"`
	Usage             *Usage         `json:"usage,omitempty"`
	SystemFingerprint string         `json:"system_fingerprint,omitempty"`
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

// Note: TranscriptionResponse is defined in whisper_types.go
// Note: ImageGenerationRequest, ImageGenerationResponse, ImageData are defined in image_types.go
