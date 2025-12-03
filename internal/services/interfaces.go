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
	"github.com/kawai-network/veridium/pkg/yzma/message"
)

// LLMGenerator defines the interface for LLM generation operations
// This allows mocking the LLM for testing
type LLMGenerator interface {
	// Generate generates a response from messages (single turn, no tool execution)
	Generate(ctx context.Context, messages []message.Message) (*llama.YzmaResponse, error)

	// RunAgentLoop runs the agent loop with tool execution
	RunAgentLoop(ctx context.Context, messages []message.Message, maxIterations int) (*llama.YzmaResponse, []message.Message, error)

	// RunAgentLoopWithStreaming runs the agent loop with streaming callback and tool event callback
	RunAgentLoopWithStreaming(ctx context.Context, messages []message.Message, maxIterations int, streamCallback llama.StreamCallback, toolCallback llama.ToolEventCallback) (*llama.YzmaResponse, []message.Message, error)

	// WithTools returns a new generator configured with specific tools
	WithTools(toolNames []string) LLMGenerator
}

// LLMGeneratorAdapter wraps LlamaYzmaModel to implement LLMGenerator interface
type LLMGeneratorAdapter struct {
	model *llama.LlamaYzmaModel
}

// NewLLMGeneratorAdapter creates a new adapter wrapping LlamaYzmaModel
func NewLLMGeneratorAdapter(model *llama.LlamaYzmaModel) *LLMGeneratorAdapter {
	return &LLMGeneratorAdapter{model: model}
}

// Generate implements LLMGenerator.Generate
func (a *LLMGeneratorAdapter) Generate(ctx context.Context, messages []message.Message) (*llama.YzmaResponse, error) {
	return a.model.Generate(ctx, messages)
}

// RunAgentLoop implements LLMGenerator.RunAgentLoop
func (a *LLMGeneratorAdapter) RunAgentLoop(ctx context.Context, messages []message.Message, maxIterations int) (*llama.YzmaResponse, []message.Message, error) {
	return a.model.RunAgentLoop(ctx, messages, maxIterations)
}

// RunAgentLoopWithStreaming implements LLMGenerator.RunAgentLoopWithStreaming
func (a *LLMGeneratorAdapter) RunAgentLoopWithStreaming(ctx context.Context, messages []message.Message, maxIterations int, streamCallback llama.StreamCallback, toolCallback llama.ToolEventCallback) (*llama.YzmaResponse, []message.Message, error) {
	return a.model.RunAgentLoopWithStreaming(ctx, messages, maxIterations, streamCallback, toolCallback)
}

// WithTools implements LLMGenerator.WithTools
func (a *LLMGeneratorAdapter) WithTools(toolNames []string) LLMGenerator {
	return &LLMGeneratorAdapter{model: a.model.WithTools(toolNames)}
}

// GetModel returns the underlying LlamaYzmaModel (for cases where direct access is needed)
func (a *LLMGeneratorAdapter) GetModel() *llama.LlamaYzmaModel {
	return a.model
}
