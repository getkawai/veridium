package services

import (
	"context"
	"fmt"
	"log"

	"github.com/kawai-network/veridium/pkg/fantasy"
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
	var model fantasy.LanguageModel

	// TODO: Phase 2 - Look up model from LibService dynamically based on req.Model/req.Provider
	// For now, we use the pre-configured chatModel or summaryModel
	if s.chatModel != nil {
		model = s.chatModel
	} else if s.summaryModel != nil {
		model = s.summaryModel
	} else {
		return "", fmt.Errorf("no language model available for completion")
	}

	// 2. Convert messages to fantasy format
	msgs := make([]fantasy.Message, len(req.Messages))
	for i, m := range req.Messages {
		role := fantasy.MessageRoleUser
		if m.Role == "system" {
			role = fantasy.MessageRoleSystem
		} else if m.Role == "assistant" {
			role = fantasy.MessageRoleAssistant
		}

		msgs[i] = fantasy.Message{
			Role:    role,
			Content: []fantasy.MessagePart{fantasy.TextPart{Text: m.Content}},
		}
	}

	// 3. Execute generation
	// We use the model directly instead of Agent since we don't need tools/loop
	resp, err := model.Generate(ctx, fantasy.Call{
		Prompt: msgs,
	})
	if err != nil {
		log.Printf("❌ [ChatCompletion] Generation failed: %v", err)
		return "", err
	}

	return resp.Content.Text(), nil
}
