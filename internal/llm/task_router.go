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

package llm

import (
	"context"
	"log"
	"sync"

	"github.com/kawai-network/veridium/pkg/yzma/tools"
	"github.com/kawai-network/veridium/types"
	"github.com/kawai-network/veridium/types/message"
)

// TaskType defines the type of LLM task for routing
type TaskType string

const (
	TaskChat              TaskType = "chat"               // Main conversation with streaming & tools
	TaskTitleGen          TaskType = "title"              // Title generation (lightweight)
	TaskSummaryGen        TaskType = "summary"            // Summary generation (background)
	TaskImageDescribe     TaskType = "image_describe"     // Image description (VL model)
	TaskOCRCleanup        TaskType = "ocr_cleanup"        // OCR text cleanup (remote first, local fallback)
	TaskTranscriptCleanup TaskType = "transcript_cleanup" // Video transcript cleanup (remote first, local fallback)
)

// ============================================================================
// TASK ROUTING CONFIGURATION
// ============================================================================
//
// TaskRouter distributes different LLM tasks to appropriate providers based on
// their strengths and cost efficiency.
//
// CURRENT TASK ASSIGNMENTS:
//
// | Task              | Primary       | Model              | Fallback     | Notes                     |
// |-------------------|---------------|--------------------|--------------|---------------------------|
// | Chat              | OpenRouter    | amazon/nova-2-lite | Local Llama  | Main conversation         |
// | Title             | Zhipu AI      | glm-4.6            | Local Llama  | Fast title generation     |
// | Summary           | Zhipu AI      | glm-4.6            | Local Llama  | Topic summarization       |
// | ImageDescribe     | Local Qwen VL | qwen3-vl           | None         | Vision-language (async)   |
// | OCRCleanup        | Zhipu AI      | glm-4.6            | Local Llama  | OCR text cleanup & format |
// | TranscriptCleanup | Zhipu AI      | glm-4.6            | Local Llama  | Video transcript cleanup  |
//
// FALLBACK BEHAVIOR:
// - GenerateWithoutTools() automatically tries fallback if primary fails
// - Chat streaming has its own error handling (no auto-fallback yet)
// - ImageDescribe has no fallback (requires VL capability)
//
// IMAGE/VIDEO PROCESSING FLOW:
// 1. File uploaded → FileProcessorService.ProcessFileFromPath()
// 2. For images/videos → processImageDescriptionAsync() / processVideoDescriptionAsync()
// 3. Local Qwen VL generates description (async, ~60-90 seconds for images)
// 4. Description saved to documents table
// 5. LLM uses "lobe-image-describe__getImageDescription" tool to get description
//   - Tool polls DB for up to 2 minutes waiting for VL to complete
//   - Returns description content from documents table
//
// ============================================================================
type TaskRouter struct {
	providers    map[TaskType]Provider
	fallback     Provider
	toolRegistry *tools.ToolRegistry
	mu           sync.RWMutex
}

// NewTaskRouter creates a new task router with optional fallback provider
func NewTaskRouter(toolRegistry *tools.ToolRegistry, fallback Provider) *TaskRouter {
	return &TaskRouter{
		providers:    make(map[TaskType]Provider),
		fallback:     fallback,
		toolRegistry: toolRegistry,
	}
}

// SetProvider sets a provider for a specific task type
func (r *TaskRouter) SetProvider(task TaskType, provider Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.providers[task] = provider
	log.Printf("🔀 TaskRouter: Set provider for task '%s'", task)
}

// SetFallback sets the fallback provider used when no specific provider is set
func (r *TaskRouter) SetFallback(provider Provider) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fallback = provider
	log.Printf("🔀 TaskRouter: Set fallback provider")
}

// GetProvider returns the provider for a task, or fallback if not set
func (r *TaskRouter) GetProvider(task TaskType) Provider {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if provider, exists := r.providers[task]; exists {
		return provider
	}
	return r.fallback
}

// HasProvider checks if a specific provider is set for a task
func (r *TaskRouter) HasProvider(task TaskType) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.providers[task]
	return exists
}

// GetToolRegistry returns the tool registry
func (r *TaskRouter) GetToolRegistry() *tools.ToolRegistry {
	return r.toolRegistry
}

// Generate routes to the appropriate provider and generates a response
func (r *TaskRouter) Generate(ctx context.Context, task TaskType, messages []message.Message) (*types.LLMResponse, error) {
	provider := r.GetProvider(task)
	if provider == nil {
		log.Printf("⚠️  TaskRouter: No provider for task '%s' and no fallback set", task)
		return nil, ErrNoProvider
	}

	log.Printf("🔀 TaskRouter: Routing '%s' task to provider", task)
	return provider.Generate(ctx, messages)
}

// GenerateWithoutTools routes to provider with tools disabled (for utility tasks)
// If the primary provider fails, it will try the fallback provider
func (r *TaskRouter) GenerateWithoutTools(ctx context.Context, task TaskType, messages []message.Message) (*types.LLMResponse, error) {
	provider := r.GetProvider(task)
	if provider == nil {
		log.Printf("⚠️  TaskRouter: No provider for task '%s' and no fallback set", task)
		return nil, ErrNoProvider
	}

	log.Printf("🔀 TaskRouter: Routing '%s' task to provider (no tools)", task)
	resp, err := provider.WithoutTools().Generate(ctx, messages)

	// If primary provider failed and we have a fallback, try it
	if err != nil && r.fallback != nil && r.fallback != provider {
		log.Printf("⚠️  TaskRouter: Primary provider failed for '%s': %v, trying fallback", task, err)
		fallbackResp, fallbackErr := r.fallback.WithoutTools().Generate(ctx, messages)
		if fallbackErr == nil {
			log.Printf("✅ TaskRouter: Fallback succeeded for '%s'", task)
			return fallbackResp, nil
		}
		log.Printf("⚠️  TaskRouter: Fallback also failed for '%s': %v", task, fallbackErr)
		// Return original error
	}

	return resp, err
}

// Chat is a convenience method for chat task with full agent loop
func (r *TaskRouter) Chat(ctx context.Context, messages []message.Message, maxIterations int) (*types.LLMResponse, []message.Message, error) {
	provider := r.GetProvider(TaskChat)
	if provider == nil {
		return nil, nil, ErrNoProvider
	}
	return provider.RunAgentLoop(ctx, messages, maxIterations)
}

// ChatWithStreaming is a convenience method for streaming chat
func (r *TaskRouter) ChatWithStreaming(ctx context.Context, messages []message.Message, maxIterations int, streamCallback types.StreamCallback, toolCallback types.ToolEventCallback) (*types.LLMResponse, []message.Message, error) {
	provider := r.GetProvider(TaskChat)
	if provider == nil {
		return nil, nil, ErrNoProvider
	}
	return provider.RunAgentLoopWithStreaming(ctx, messages, maxIterations, streamCallback, toolCallback)
}

// ChatWithTools returns the chat provider with specific tools enabled
func (r *TaskRouter) ChatWithTools(toolNames []string) Provider {
	provider := r.GetProvider(TaskChat)
	if provider == nil {
		return nil
	}
	return provider.WithTools(toolNames)
}

// GenerateTitle is a convenience method for title generation
func (r *TaskRouter) GenerateTitle(ctx context.Context, messages []message.Message) (*types.LLMResponse, error) {
	return r.GenerateWithoutTools(ctx, TaskTitleGen, messages)
}

// GenerateSummary is a convenience method for summary generation
func (r *TaskRouter) GenerateSummary(ctx context.Context, messages []message.Message) (*types.LLMResponse, error) {
	return r.GenerateWithoutTools(ctx, TaskSummaryGen, messages)
}

// DescribeImage is a convenience method for image description (VL task)
func (r *TaskRouter) DescribeImage(ctx context.Context, messages []message.Message) (*types.LLMResponse, error) {
	return r.GenerateWithoutTools(ctx, TaskImageDescribe, messages)
}

// ListConfiguredTasks returns all task types that have providers set
func (r *TaskRouter) ListConfiguredTasks() []TaskType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]TaskType, 0, len(r.providers))
	for task := range r.providers {
		tasks = append(tasks, task)
	}
	return tasks
}

// Error types
var (
	ErrNoProvider = &RouterError{Message: "no provider configured for task"}
)

type RouterError struct {
	Message string
}

func (e *RouterError) Error() string {
	return e.Message
}
