package services

import (
	"context"
	"fmt"

	"github.com/kawai-network/veridium/internal/llm"
	"github.com/kawai-network/veridium/pkg/yzma/message"
)

// TaskRouterAdapter wraps TaskRouter to implement LLMProvider interface
type TaskRouterAdapter struct {
	router   *llm.TaskRouter
	taskType llm.TaskType
}

// NewTaskRouterAdapter creates a new adapter for TaskRouter
// Uses TaskOCRCleanup for OCR text cleanup (remote first, local fallback)
func NewTaskRouterAdapter(router *llm.TaskRouter) *TaskRouterAdapter {
	return &TaskRouterAdapter{
		router:   router,
		taskType: llm.TaskOCRCleanup, // Use OCR cleanup task (Zhipu -> Local fallback)
	}
}

// GenerateText implements LLMProvider interface
func (a *TaskRouterAdapter) GenerateText(ctx context.Context, prompt string) (string, error) {
	if a.router == nil {
		return "", fmt.Errorf("task router not initialized")
	}

	messages := []message.Message{
		message.Chat{
			Role:    "user",
			Content: prompt,
		},
	}

	resp, err := a.router.GenerateWithoutTools(ctx, a.taskType, messages)
	if err != nil {
		return "", fmt.Errorf("LLM generation failed: %w", err)
	}

	if resp == nil || resp.Content == "" {
		return "", fmt.Errorf("empty response from LLM")
	}

	return resp.Content, nil
}
