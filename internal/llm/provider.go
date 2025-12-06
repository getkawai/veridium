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

// Package llm provides interfaces and types for LLM providers
package llm

import (
	"context"

	"github.com/kawai-network/veridium/types"
	"github.com/kawai-network/veridium/types/message"
)

// Provider defines the interface for LLM generation operations
// This interface is implemented by all LLM providers (local llama, OpenAI-compatible APIs, etc.)
type Provider interface {
	// Generate generates a response from messages (single turn, no tool execution)
	Generate(ctx context.Context, messages []message.Message) (*types.LLMResponse, error)

	// RunAgentLoop runs the agent loop with tool execution
	RunAgentLoop(ctx context.Context, messages message.Prompt, maxIterations int) (*types.LLMResponse, message.Prompt, error)

	// RunAgentLoopWithStreaming runs the agent loop with streaming callback and tool event callback
	RunAgentLoopWithStreaming(ctx context.Context, messages message.Prompt, maxIterations int, streamCallback types.StreamCallback, toolCallback types.ToolEventCallback) (*types.LLMResponse, message.Prompt, error)

	// WithTools returns a new provider configured with specific tools
	WithTools(toolNames []string) Provider

	// WithoutTools returns a new provider with tools disabled
	WithoutTools() Provider
}
