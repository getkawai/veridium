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

package llama

import (
	"context"
	"iter"

	"github.com/kawai-network/veridium/fantasy"
)

// LlamaLanguageModel wraps LlamaYzmaModel to implement fantasy.LanguageModel interface
type LlamaLanguageModel struct {
	model     *LlamaYzmaModel
	modelName string
}

// NewLlamaLanguageModel creates a new LlamaLanguageModel wrapper
func NewLlamaLanguageModel(model *LlamaYzmaModel, modelName string) *LlamaLanguageModel {
	return &LlamaLanguageModel{
		model:     model,
		modelName: modelName,
	}
}

// Generate implements fantasy.LanguageModel.Generate
func (m *LlamaLanguageModel) Generate(ctx context.Context, call fantasy.Call) (*fantasy.Response, error) {
	// Convert Call.Prompt to []fantasy.Message
	messages := call.Prompt

	// Check if tools should be disabled (ToolChoiceNone)
	if call.ToolChoice != nil && *call.ToolChoice == fantasy.ToolChoiceNone {
		return m.model.WithoutTools().Generate(ctx, messages)
	}

	// If tools are provided, set them on the underlying model
	if len(call.Tools) > 0 {
		// Tools are already in the messages via system prompt enhancement
		// The LlamaYzmaModel.enhanceWithTools handles this
	}

	return m.model.Generate(ctx, messages)
}

// Stream implements fantasy.LanguageModel.Stream
func (m *LlamaLanguageModel) Stream(ctx context.Context, call fantasy.Call) (fantasy.StreamResponse, error) {
	messages := call.Prompt

	// Check if tools should be disabled (ToolChoiceNone)
	modelToUse := m.model
	if call.ToolChoice != nil && *call.ToolChoice == fantasy.ToolChoiceNone {
		modelToUse = m.model.WithoutTools()
	}

	// Create an iterator that yields stream parts
	streamFunc := func(yield func(fantasy.StreamPart) bool) {
		textID := "text_0"

		// Emit text start
		if !yield(fantasy.StreamPart{
			Type: fantasy.StreamPartTypeTextStart,
			ID:   textID,
		}) {
			return
		}

		// Streaming callback
		callback := func(token string, isLast bool) {
			if token != "" {
				yield(fantasy.StreamPart{
					Type:  fantasy.StreamPartTypeTextDelta,
					ID:    textID,
					Delta: token,
				})
			}
		}

		// Run streaming generation
		resp, err := modelToUse.Stream(ctx, messages, callback)

		// Emit text end
		yield(fantasy.StreamPart{
			Type: fantasy.StreamPartTypeTextEnd,
			ID:   textID,
		})

		if err != nil {
			yield(fantasy.StreamPart{
				Type:  fantasy.StreamPartTypeError,
				Error: err,
			})
			return
		}

		// Emit tool calls if any
		toolCalls := resp.Content.ToolCalls()
		for _, tc := range toolCalls {
			// Emit tool input start
			yield(fantasy.StreamPart{
				Type:         fantasy.StreamPartTypeToolInputStart,
				ID:           tc.ToolCallID,
				ToolCallName: tc.ToolName,
			})

			// Emit tool input delta (full input at once for local models)
			yield(fantasy.StreamPart{
				Type:  fantasy.StreamPartTypeToolInputDelta,
				ID:    tc.ToolCallID,
				Delta: tc.Input,
			})

			// Emit tool input end
			yield(fantasy.StreamPart{
				Type: fantasy.StreamPartTypeToolInputEnd,
				ID:   tc.ToolCallID,
			})

			// Emit tool call
			yield(fantasy.StreamPart{
				Type:          fantasy.StreamPartTypeToolCall,
				ID:            tc.ToolCallID,
				ToolCallName:  tc.ToolName,
				ToolCallInput: tc.Input,
			})
		}

		// Emit finish
		yield(fantasy.StreamPart{
			Type:         fantasy.StreamPartTypeFinish,
			FinishReason: resp.FinishReason,
			Usage:        resp.Usage,
		})
	}

	return iter.Seq[fantasy.StreamPart](streamFunc), nil
}

// GenerateObject implements fantasy.LanguageModel.GenerateObject
// Not supported for local llama models
func (m *LlamaLanguageModel) GenerateObject(ctx context.Context, call fantasy.ObjectCall) (*fantasy.ObjectResponse, error) {
	return nil, &fantasy.Error{
		Title:   "not supported",
		Message: "GenerateObject is not supported for local llama models",
	}
}

// StreamObject implements fantasy.LanguageModel.StreamObject
// Not supported for local llama models
func (m *LlamaLanguageModel) StreamObject(ctx context.Context, call fantasy.ObjectCall) (fantasy.ObjectStreamResponse, error) {
	return nil, &fantasy.Error{
		Title:   "not supported",
		Message: "StreamObject is not supported for local llama models",
	}
}

// Provider implements fantasy.LanguageModel.Provider
func (m *LlamaLanguageModel) Provider() string {
	return "llama-local"
}

// Model implements fantasy.LanguageModel.Model
func (m *LlamaLanguageModel) Model() string {
	return m.modelName
}

// Ensure LlamaLanguageModel implements fantasy.LanguageModel
var _ fantasy.LanguageModel = (*LlamaLanguageModel)(nil)
