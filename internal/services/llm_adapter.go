package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/internal/llm"
)

// TaskRouterAdapter wraps TaskRouter to implement LLMProvider interface
// It auto-detects whether the prompt is for OCR cleanup or transcript cleanup
type TaskRouterAdapter struct {
	router *llm.TaskRouter
}

// NewTaskRouterAdapter creates a new adapter for TaskRouter
func NewTaskRouterAdapter(router *llm.TaskRouter) *TaskRouterAdapter {
	return &TaskRouterAdapter{
		router: router,
	}
}

// System prompts for different cleanup tasks
const (
	// SystemPromptTranscriptCleanup is a global system prompt for transcript cleanup
	// This helps the model understand the task without complex user prompts
	SystemPromptTranscriptCleanup = `You are a transcription editor for Indonesian language content.
Your task is to fix common speech-to-text errors. Apply these corrections:
- terlion/terliunan → triliun
- stengah → setengah
- merogikan → merugikan
- ngontongnya → menguntungkan
- meleksaham → melek saham
- persubstritp → oversubscribe
- SOJK/SCUJK → SEOJK
- IHSK → IHSG
- rebu → ribu
- lokasinya → alokasinya
- timis → tipis

Keep segment markers (**[Segment X]**). Keep stock codes (RLCO, GOTO, PGHB).
Output ONLY the corrected text without explanations.`

	// SystemPromptOCRCleanup is a global system prompt for OCR text cleanup
	SystemPromptOCRCleanup = `You are an OCR text editor. Clean up raw OCR output:
- Fix obvious typos and character recognition errors
- Preserve original structure and formatting
- Format as proper markdown when appropriate
- Output ONLY the cleaned text without explanations.`
)

// GenerateText implements LLMProvider interface
// Auto-detects task type based on prompt content and adds appropriate system prompt
func (a *TaskRouterAdapter) GenerateText(ctx context.Context, prompt string) (string, error) {
	if a.router == nil {
		return "", fmt.Errorf("task router not initialized")
	}

	// Auto-detect task type based on prompt keywords
	taskType := a.detectTaskType(prompt)

	// Build messages with system prompt
	messages := fantasy.Prompt{}

	// Add system prompt based on task type
	switch taskType {
	case llm.TaskTranscriptCleanup:
		messages = append(messages, fantasy.NewSystemMessage(SystemPromptTranscriptCleanup))
	case llm.TaskOCRCleanup:
		messages = append(messages, fantasy.NewSystemMessage(SystemPromptOCRCleanup))
	}

	// Add user prompt
	messages = append(messages, fantasy.NewUserMessage(prompt))

	resp, err := a.router.GenerateWithFallback(ctx, taskType, messages)
	if err != nil {
		return "", fmt.Errorf("LLM generation failed: %w", err)
	}

	respText := resp.Content.Text()
	if resp == nil || respText == "" {
		return "", fmt.Errorf("empty response from LLM")
	}

	return respText, nil
}

// detectTaskType determines which task type to use based on prompt content
func (a *TaskRouterAdapter) detectTaskType(prompt string) llm.TaskType {
	promptLower := strings.ToLower(prompt)

	// Check for transcript cleanup indicators
	if strings.Contains(promptLower, "transcription") ||
		strings.Contains(promptLower, "whisper") ||
		strings.Contains(promptLower, "speech-to-text") ||
		strings.Contains(promptLower, "video transcription") {
		return llm.TaskTranscriptCleanup
	}

	// Default to OCR cleanup
	return llm.TaskOCRCleanup
}
