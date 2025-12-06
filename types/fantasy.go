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
	"fmt"
	"iter"
)

// ============================================================================
// Provider Options & Metadata (extensibility layer for provider-specific data)
// ============================================================================

// ProviderOptionsData is an interface for provider-specific options data.
// Implementations should also implement json.Marshaler and json.Unmarshaler.
type ProviderOptionsData interface {
	Options()
}

// ProviderMetadata represents additional provider-specific metadata in responses.
// Key is the provider name (e.g., "anthropic", "openai").
type ProviderMetadata map[string]ProviderOptionsData

// ProviderOptions represents additional provider-specific options in requests.
// Key is the provider name (e.g., "anthropic", "openai").
type ProviderOptions map[string]ProviderOptionsData

// ============================================================================
// Finish Reason
// ============================================================================

// FinishReason represents why a language model finished generating a response.
type FinishReason string

const (
	FinishReasonStop          FinishReason = "stop"           // model generated stop sequence
	FinishReasonLength        FinishReason = "length"         // model generated maximum number of tokens
	FinishReasonContentFilter FinishReason = "content-filter" // content filter violation stopped the model
	FinishReasonToolCalls     FinishReason = "tool-calls"     // model triggered tool calls
	FinishReasonError         FinishReason = "error"          // model stopped because of an error
	FinishReasonOther         FinishReason = "other"          // model stopped for other reasons
	FinishReasonUnknown       FinishReason = "unknown"        // the model has not transmitted a finish reason
)

// ============================================================================
// Usage Statistics
// ============================================================================

// Usage represents token usage statistics for a model call.
type Usage struct {
	InputTokens         int64 `json:"input_tokens"`
	OutputTokens        int64 `json:"output_tokens"`
	TotalTokens         int64 `json:"total_tokens"`
	ReasoningTokens     int64 `json:"reasoning_tokens"`
	CacheCreationTokens int64 `json:"cache_creation_tokens"`
	CacheReadTokens     int64 `json:"cache_read_tokens"`
}

func (u Usage) String() string {
	return fmt.Sprintf("Usage{Input: %d, Output: %d, Total: %d, Reasoning: %d, CacheCreation: %d, CacheRead: %d}",
		u.InputTokens, u.OutputTokens, u.TotalTokens, u.ReasoningTokens, u.CacheCreationTokens, u.CacheReadTokens)
}

// ============================================================================
// Tool Choice
// ============================================================================

// ModelToolChoice represents the tool choice preference for a model call.
type ModelToolChoice string

const (
	ModelToolChoiceNone     ModelToolChoice = "none"     // no tools should be used
	ModelToolChoiceAuto     ModelToolChoice = "auto"     // tools should be used automatically
	ModelToolChoiceRequired ModelToolChoice = "required" // tools are required
)

// SpecificToolChoice creates a tool choice for a specific tool name.
func SpecificToolChoice(name string) ModelToolChoice {
	return ModelToolChoice(name)
}

// ============================================================================
// Call Configuration (unified request parameters)
// ============================================================================

// ModelCall represents a call to a language model with all configuration.
type ModelCall struct {
	Prompt           Prompt           `json:"prompt"`
	MaxOutputTokens  *int64           `json:"max_output_tokens,omitempty"`
	Temperature      *float64         `json:"temperature,omitempty"`
	TopP             *float64         `json:"top_p,omitempty"`
	TopK             *int64           `json:"top_k,omitempty"`
	PresencePenalty  *float64         `json:"presence_penalty,omitempty"`
	FrequencyPenalty *float64         `json:"frequency_penalty,omitempty"`
	Tools            []ModelTool      `json:"tools,omitempty"`
	ToolChoice       *ModelToolChoice `json:"tool_choice,omitempty"`
	ProviderOptions  ProviderOptions  `json:"provider_options,omitempty"`
}

// ============================================================================
// Call Warnings
// ============================================================================

// CallWarningType represents the type of call warning.
type CallWarningType string

const (
	CallWarningTypeUnsupportedSetting CallWarningType = "unsupported-setting"
	CallWarningTypeUnsupportedTool    CallWarningType = "unsupported-tool"
	CallWarningTypeOther              CallWarningType = "other"
)

// CallWarning represents a warning from the model provider.
type CallWarning struct {
	Type    CallWarningType `json:"type"`
	Setting string          `json:"setting,omitempty"`
	Tool    ModelTool       `json:"tool,omitempty"`
	Details string          `json:"details,omitempty"`
	Message string          `json:"message,omitempty"`
}

// ============================================================================
// Response Types
// ============================================================================

// ModelResponse represents a response from a language model.
type ModelResponse struct {
	Content          ModelResponseContent `json:"content"`
	FinishReason     FinishReason         `json:"finish_reason"`
	Usage            Usage                `json:"usage"`
	Warnings         []CallWarning        `json:"warnings,omitempty"`
	ProviderMetadata ProviderMetadata     `json:"provider_metadata,omitempty"`
}

// ModelResponseContent represents the content of a model response.
type ModelResponseContent []ModelContent

// Text returns the concatenated text content of the response.
func (r ModelResponseContent) Text() string {
	var text string
	for _, c := range r {
		if c.GetContentType() == ModelContentTypeText {
			if tc, ok := c.(ModelTextContent); ok {
				text += tc.Text
			}
		}
	}
	return text
}

// ReasoningText returns all reasoning content as a concatenated string.
func (r ModelResponseContent) ReasoningText() string {
	var text string
	for _, c := range r {
		if c.GetContentType() == ModelContentTypeReasoning {
			if rc, ok := c.(ModelReasoningContent); ok {
				text += rc.Text
			}
		}
	}
	return text
}

// ToolCalls returns all tool call content parts.
func (r ModelResponseContent) ToolCalls() []ModelToolCallContent {
	var calls []ModelToolCallContent
	for _, c := range r {
		if c.GetContentType() == ModelContentTypeToolCall {
			if tc, ok := c.(ModelToolCallContent); ok {
				calls = append(calls, tc)
			}
		}
	}
	return calls
}

// HasToolCalls returns true if the response contains tool calls.
func (r ModelResponseContent) HasToolCalls() bool {
	for _, c := range r {
		if c.GetContentType() == ModelContentTypeToolCall {
			return true
		}
	}
	return false
}

// ============================================================================
// Streaming Types
// ============================================================================

// StreamPartType represents the type of a stream part.
type StreamPartType string

const (
	StreamPartTypeWarnings       StreamPartType = "warnings"
	StreamPartTypeTextStart      StreamPartType = "text_start"
	StreamPartTypeTextDelta      StreamPartType = "text_delta"
	StreamPartTypeTextEnd        StreamPartType = "text_end"
	StreamPartTypeReasoningStart StreamPartType = "reasoning_start"
	StreamPartTypeReasoningDelta StreamPartType = "reasoning_delta"
	StreamPartTypeReasoningEnd   StreamPartType = "reasoning_end"
	StreamPartTypeToolInputStart StreamPartType = "tool_input_start"
	StreamPartTypeToolInputDelta StreamPartType = "tool_input_delta"
	StreamPartTypeToolInputEnd   StreamPartType = "tool_input_end"
	StreamPartTypeToolCall       StreamPartType = "tool_call"
	StreamPartTypeToolResult     StreamPartType = "tool_result"
	StreamPartTypeSource         StreamPartType = "source"
	StreamPartTypeFinish         StreamPartType = "finish"
	StreamPartTypeError          StreamPartType = "error"
)

// StreamPart represents a part of a streaming response.
type StreamPart struct {
	Type             StreamPartType   `json:"type"`
	ID               string           `json:"id,omitempty"`
	ToolCallName     string           `json:"tool_call_name,omitempty"`
	ToolCallInput    string           `json:"tool_call_input,omitempty"`
	Delta            string           `json:"delta,omitempty"`
	ProviderExecuted bool             `json:"provider_executed,omitempty"`
	Usage            Usage            `json:"usage,omitempty"`
	FinishReason     FinishReason     `json:"finish_reason,omitempty"`
	Error            error            `json:"error,omitempty"`
	Warnings         []CallWarning    `json:"warnings,omitempty"`
	SourceType       SourceType       `json:"source_type,omitempty"`
	URL              string           `json:"url,omitempty"`
	Title            string           `json:"title,omitempty"`
	ProviderMetadata ProviderMetadata `json:"provider_metadata,omitempty"`
}

// StreamResponse represents a streaming response sequence.
type StreamResponse = iter.Seq[StreamPart]

// ============================================================================
// Content Types (for Response)
// ============================================================================

// ModelContentType represents the type of content in a response.
type ModelContentType string

const (
	ModelContentTypeText       ModelContentType = "text"
	ModelContentTypeReasoning  ModelContentType = "reasoning"
	ModelContentTypeFile       ModelContentType = "file"
	ModelContentTypeSource     ModelContentType = "source"
	ModelContentTypeToolCall   ModelContentType = "tool-call"
	ModelContentTypeToolResult ModelContentType = "tool-result"
)

// ModelContent represents generated content from the model.
type ModelContent interface {
	GetContentType() ModelContentType
}

// ModelTextContent represents text that the model has generated.
type ModelTextContent struct {
	Text             string           `json:"text"`
	ProviderMetadata ProviderMetadata `json:"provider_metadata,omitempty"`
}

func (t ModelTextContent) GetContentType() ModelContentType {
	return ModelContentTypeText
}

// ModelReasoningContent represents reasoning that the model has generated.
type ModelReasoningContent struct {
	Text             string           `json:"text"`
	ProviderMetadata ProviderMetadata `json:"provider_metadata,omitempty"`
}

func (r ModelReasoningContent) GetContentType() ModelContentType {
	return ModelContentTypeReasoning
}

// ModelFileContent represents a file that has been generated by the model.
type ModelFileContent struct {
	MediaType        string           `json:"media_type"`
	Data             []byte           `json:"data"`
	ProviderMetadata ProviderMetadata `json:"provider_metadata,omitempty"`
}

func (f ModelFileContent) GetContentType() ModelContentType {
	return ModelContentTypeFile
}

// SourceType represents the type of source.
type SourceType string

const (
	SourceTypeURL      SourceType = "url"
	SourceTypeDocument SourceType = "document"
)

// ModelSourceContent represents a source used as input to generate the response.
type ModelSourceContent struct {
	SourceType       SourceType       `json:"source_type"`
	ID               string           `json:"id"`
	URL              string           `json:"url,omitempty"`
	Title            string           `json:"title,omitempty"`
	MediaType        string           `json:"media_type,omitempty"`
	Filename         string           `json:"filename,omitempty"`
	ProviderMetadata ProviderMetadata `json:"provider_metadata,omitempty"`
}

func (s ModelSourceContent) GetContentType() ModelContentType {
	return ModelContentTypeSource
}

// ModelToolCallContent represents tool calls that the model has generated.
type ModelToolCallContent struct {
	ToolCallID       string           `json:"tool_call_id"`
	ToolName         string           `json:"tool_name"`
	Input            string           `json:"input"` // JSON string
	ProviderExecuted bool             `json:"provider_executed,omitempty"`
	ProviderMetadata ProviderMetadata `json:"provider_metadata,omitempty"`
	Invalid          bool             `json:"invalid,omitempty"`
	ValidationError  error            `json:"validation_error,omitempty"`
}

func (t ModelToolCallContent) GetContentType() ModelContentType {
	return ModelContentTypeToolCall
}

// ModelToolResultContent represents result of a tool call.
type ModelToolResultContent struct {
	ToolCallID       string                `json:"tool_call_id"`
	ToolName         string                `json:"tool_name"`
	Result           ModelToolResultOutput `json:"result"`
	ClientMetadata   string                `json:"client_metadata,omitempty"`
	ProviderExecuted bool                  `json:"provider_executed,omitempty"`
	ProviderMetadata ProviderMetadata      `json:"provider_metadata,omitempty"`
}

func (t ModelToolResultContent) GetContentType() ModelContentType {
	return ModelContentTypeToolResult
}

// ============================================================================
// Tool Result Output Types
// ============================================================================

// ModelToolResultOutputType represents the type of tool result output.
type ModelToolResultOutputType string

const (
	ModelToolResultOutputTypeText  ModelToolResultOutputType = "text"
	ModelToolResultOutputTypeError ModelToolResultOutputType = "error"
	ModelToolResultOutputTypeMedia ModelToolResultOutputType = "media"
)

// ModelToolResultOutput represents the output content of a tool result.
type ModelToolResultOutput interface {
	GetOutputType() ModelToolResultOutputType
}

// ModelToolResultOutputText represents text output.
type ModelToolResultOutputText struct {
	Text string `json:"text"`
}

func (t ModelToolResultOutputText) GetOutputType() ModelToolResultOutputType {
	return ModelToolResultOutputTypeText
}

// ModelToolResultOutputError represents error output.
type ModelToolResultOutputError struct {
	Error error `json:"error"`
}

func (t ModelToolResultOutputError) GetOutputType() ModelToolResultOutputType {
	return ModelToolResultOutputTypeError
}

// ModelToolResultOutputMedia represents media output.
type ModelToolResultOutputMedia struct {
	Data      string `json:"data"`       // base64 encoded
	MediaType string `json:"media_type"`
}

func (t ModelToolResultOutputMedia) GetOutputType() ModelToolResultOutputType {
	return ModelToolResultOutputTypeMedia
}

// ============================================================================
// Tool Types (for Call configuration)
// ============================================================================

// ModelToolType represents the type of tool.
type ModelToolType string

const (
	ModelToolTypeFunction        ModelToolType = "function"
	ModelToolTypeProviderDefined ModelToolType = "provider-defined"
)

// ModelTool represents a tool that can be used by the model.
type ModelTool interface {
	GetToolType() ModelToolType
	GetToolName() string
}

// ModelFunctionTool represents a function tool.
type ModelFunctionTool struct {
	Name            string          `json:"name"`
	Description     string          `json:"description"`
	InputSchema     map[string]any  `json:"input_schema"`
	ProviderOptions ProviderOptions `json:"provider_options,omitempty"`
}

func (f ModelFunctionTool) GetToolType() ModelToolType {
	return ModelToolTypeFunction
}

func (f ModelFunctionTool) GetToolName() string {
	return f.Name
}

// ModelProviderDefinedTool represents a provider-defined tool.
type ModelProviderDefinedTool struct {
	ID   string         `json:"id"`
	Name string         `json:"name"`
	Args map[string]any `json:"args,omitempty"`
}

func (p ModelProviderDefinedTool) GetToolType() ModelToolType {
	return ModelToolTypeProviderDefined
}

func (p ModelProviderDefinedTool) GetToolName() string {
	return p.Name
}

// ============================================================================
// Provider & LanguageModel Interfaces
// ============================================================================

// ModelProvider represents a provider of language models.
type ModelProvider interface {
	Name() string
	LanguageModel(ctx context.Context, modelID string) (LanguageModel, error)
}

// LanguageModel represents a language model that can generate responses.
type LanguageModel interface {
	Generate(context.Context, ModelCall) (*ModelResponse, error)
	Stream(context.Context, ModelCall) (StreamResponse, error)
	Provider() string
	Model() string
}

// ============================================================================
// Type Assertion Helpers
// ============================================================================

// AsModelContent converts a ModelContent interface to a specific type.
func AsModelContent[T ModelContent](content ModelContent) (T, bool) {
	var zero T
	if content == nil {
		return zero, false
	}
	if v, ok := content.(T); ok {
		return v, true
	}
	return zero, false
}

// AsModelToolResultOutput converts a ModelToolResultOutput interface to a specific type.
func AsModelToolResultOutput[T ModelToolResultOutput](output ModelToolResultOutput) (T, bool) {
	var zero T
	if output == nil {
		return zero, false
	}
	if v, ok := output.(T); ok {
		return v, true
	}
	return zero, false
}
