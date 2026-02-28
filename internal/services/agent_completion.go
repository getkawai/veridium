package services

import (
	"context"
	"fmt"
	"log"

	unillm "github.com/getkawai/unillm"
)

// ChatCompletionRequest represents a stateless chat completion request
// This is used for tasks like translation, summarization, etc.
type ChatCompletionRequest struct {
	Model    string                  `json:"model"`
	Provider string                  `json:"provider"`
	Messages []ChatCompletionMessage `json:"messages"`
}

// ChatCompletionMessage represents a message in the completion request
type ChatCompletionMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletion executes a stateless chat completion
// It does not use session history or agents tools, just simple LLM generation
func (s *AgentChatService) ChatCompletion(ctx context.Context, req ChatCompletionRequest) (string, error) {
	// 1. Select model
	// Priority:
	// 1. Use passed model/provider if available in LibService
	// 2. Fallback to s.chatModel
	var model unillm.LanguageModel

	// TODO: Phase 2 - Look up model from LibService dynamically based on req.Model/req.Provider
	// For now, we use the pre-configured chatModel or summaryModel
	if s.chatModel != nil {
		model = s.chatModel
	} else if s.summaryModel != nil {
		model = s.summaryModel
	} else {
		return "", fmt.Errorf("no language model available for completion")
	}

	// 2. Convert messages to unillm format
	msgs := make([]unillm.Message, len(req.Messages))
	for i, m := range req.Messages {
		role := unillm.MessageRoleUser
		switch m.Role {
		case "system":
			role = unillm.MessageRoleSystem
		case "assistant":
			role = unillm.MessageRoleAssistant
		}

		msgs[i] = unillm.Message{
			Role:    role,
			Content: []unillm.MessagePart{unillm.TextPart{Text: m.Content}},
		}
	}

	// 3. Execute generation
	// We use the model directly instead of Agent since we don't need tools/loop
	resp, err := model.Generate(ctx, unillm.Call{
		Prompt: msgs,
	})
	if err != nil {
		log.Printf("❌ [ChatCompletion] Generation failed: %v", err)
		return "", err
	}

	return resp.Content.Text(), nil
}
