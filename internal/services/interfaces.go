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

package services

import (
	"context"

	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/internal/llm"
	"github.com/kawai-network/veridium/pkg/yzma/message"
	"github.com/kawai-network/veridium/types"
)

// LlamaProviderAdapter wraps LlamaYzmaModel to implement llm.Provider interface
type LlamaProviderAdapter struct {
	model *llama.LlamaYzmaModel
}

// NewLlamaProviderAdapter creates a new adapter wrapping LlamaYzmaModel
func NewLlamaProviderAdapter(model *llama.LlamaYzmaModel) *LlamaProviderAdapter {
	return &LlamaProviderAdapter{model: model}
}

// Generate implements llm.Provider.Generate
func (a *LlamaProviderAdapter) Generate(ctx context.Context, messages []message.Message) (*types.LLMResponse, error) {
	return a.model.Generate(ctx, messages)
}

// RunAgentLoop implements llm.Provider.RunAgentLoop
func (a *LlamaProviderAdapter) RunAgentLoop(ctx context.Context, messages []message.Message, maxIterations int) (*types.LLMResponse, []message.Message, error) {
	return a.model.RunAgentLoop(ctx, messages, maxIterations)
}

// RunAgentLoopWithStreaming implements llm.Provider.RunAgentLoopWithStreaming
func (a *LlamaProviderAdapter) RunAgentLoopWithStreaming(ctx context.Context, messages []message.Message, maxIterations int, streamCallback types.StreamCallback, toolCallback types.ToolEventCallback) (*types.LLMResponse, []message.Message, error) {
	return a.model.RunAgentLoopWithStreaming(ctx, messages, maxIterations, streamCallback, toolCallback)
}

// WithTools implements llm.Provider.WithTools
func (a *LlamaProviderAdapter) WithTools(toolNames []string) llm.Provider {
	return &LlamaProviderAdapter{model: a.model.WithTools(toolNames)}
}

// WithoutTools implements llm.Provider.WithoutTools
func (a *LlamaProviderAdapter) WithoutTools() llm.Provider {
	return &LlamaProviderAdapter{model: a.model.WithoutTools()}
}

// GetModel returns the underlying LlamaYzmaModel (for cases where direct access is needed)
func (a *LlamaProviderAdapter) GetModel() *llama.LlamaYzmaModel {
	return a.model
}
