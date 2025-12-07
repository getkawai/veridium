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

	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/pkg/yzma/tools"
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
	models       map[TaskType]fantasy.LanguageModel
	fallback     fantasy.LanguageModel
	toolRegistry *tools.ToolRegistry
	mu           sync.RWMutex
}

// NewTaskRouter creates a new task router with optional fallback model
func NewTaskRouter(toolRegistry *tools.ToolRegistry, fallback fantasy.LanguageModel) *TaskRouter {
	return &TaskRouter{
		models:       make(map[TaskType]fantasy.LanguageModel),
		fallback:     fallback,
		toolRegistry: toolRegistry,
	}
}

// SetModel sets a language model for a specific task type
func (r *TaskRouter) SetModel(task TaskType, model fantasy.LanguageModel) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.models[task] = model
	log.Printf("🔀 TaskRouter: Set model for task '%s'", task)
}

// SetFallback sets the fallback model used when no specific model is set
func (r *TaskRouter) SetFallback(model fantasy.LanguageModel) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fallback = model
	log.Printf("🔀 TaskRouter: Set fallback model")
}

// GetModel returns the model for a task, or fallback if not set
func (r *TaskRouter) GetModel(task TaskType) fantasy.LanguageModel {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if model, exists := r.models[task]; exists {
		return model
	}
	return r.fallback
}

// HasModel checks if a specific model is set for a task
func (r *TaskRouter) HasModel(task TaskType) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.models[task]
	return exists
}

// GetToolRegistry returns the tool registry
func (r *TaskRouter) GetToolRegistry() *tools.ToolRegistry {
	return r.toolRegistry
}

// Generate routes to the appropriate model and generates a response
func (r *TaskRouter) Generate(ctx context.Context, task TaskType, messages fantasy.Prompt) (*fantasy.Response, error) {
	model := r.GetModel(task)
	if model == nil {
		log.Printf("⚠️  TaskRouter: No model for task '%s' and no fallback set", task)
		return nil, ErrNoProvider
	}

	log.Printf("🔀 TaskRouter: Routing '%s' task to model %s", task, model.Model())
	return model.Generate(ctx, fantasy.Call{Prompt: messages})
}

// GenerateWithFallback routes to model and tries fallback if primary fails
// Used for utility tasks (title, summary, OCR cleanup, etc.)
func (r *TaskRouter) GenerateWithFallback(ctx context.Context, task TaskType, messages fantasy.Prompt) (*fantasy.Response, error) {
	model := r.GetModel(task)
	if model == nil {
		log.Printf("⚠️  TaskRouter: No model for task '%s' and no fallback set", task)
		return nil, ErrNoProvider
	}

	log.Printf("🔀 TaskRouter: Routing '%s' task to model %s", task, model.Model())
	resp, err := model.Generate(ctx, fantasy.Call{Prompt: messages})

	// If primary model failed and we have a fallback, try it
	if err != nil && r.fallback != nil && r.fallback != model {
		log.Printf("⚠️  TaskRouter: Primary model failed for '%s': %v, trying fallback", task, err)
		fallbackResp, fallbackErr := r.fallback.Generate(ctx, fantasy.Call{Prompt: messages})
		if fallbackErr == nil {
			log.Printf("✅ TaskRouter: Fallback succeeded for '%s'", task)
			return fallbackResp, nil
		}
		log.Printf("⚠️  TaskRouter: Fallback also failed for '%s': %v", task, fallbackErr)
	}

	return resp, err
}

// GenerateTitle is a convenience method for title generation
func (r *TaskRouter) GenerateTitle(ctx context.Context, messages fantasy.Prompt) (*fantasy.Response, error) {
	return r.GenerateWithFallback(ctx, TaskTitleGen, messages)
}

// GenerateSummary is a convenience method for summary generation
func (r *TaskRouter) GenerateSummary(ctx context.Context, messages fantasy.Prompt) (*fantasy.Response, error) {
	return r.GenerateWithFallback(ctx, TaskSummaryGen, messages)
}

// DescribeImage is a convenience method for image description (VL task)
func (r *TaskRouter) DescribeImage(ctx context.Context, messages fantasy.Prompt) (*fantasy.Response, error) {
	return r.GenerateWithFallback(ctx, TaskImageDescribe, messages)
}

// ListConfiguredTasks returns all task types that have models set
func (r *TaskRouter) ListConfiguredTasks() []TaskType {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tasks := make([]TaskType, 0, len(r.models))
	for task := range r.models {
		tasks = append(tasks, task)
	}
	return tasks
}

// GetChatModel returns the fantasy.LanguageModel for chat tasks
// This is used by fantasy.Agent for streaming
func (r *TaskRouter) GetChatModel() fantasy.LanguageModel {
	return r.GetModel(TaskChat)
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
