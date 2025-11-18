package providers

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// InboxGuideConfig holds configuration for inbox guide injection
type InboxGuideConfig struct {
	InboxGuideSystemRole string
	InboxSessionID       string
	IsWelcomeQuestion    bool
	SessionID            string
}

// NewInboxGuideProviderLambda creates a lambda node for inbox guide injection
func NewInboxGuideProviderLambda(config InboxGuideConfig) *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, msgs []*schema.Message) ([]*schema.Message, error) {
		// Check if inbox guide should be injected
		shouldInject := config.IsWelcomeQuestion &&
			config.SessionID == config.InboxSessionID &&
			config.InboxGuideSystemRole != ""

		if !shouldInject {
			return msgs, nil
		}

		// Find existing system message
		var systemMsgIndex int = -1
		for i, msg := range msgs {
			if msg.Role == schema.System {
				systemMsgIndex = i
				break
			}
		}

		result := make([]*schema.Message, len(msgs))
		copy(result, msgs)

		if systemMsgIndex >= 0 {
			// Merge to existing system message
			existingMsg := result[systemMsgIndex]
			parts := []string{}
			if existingMsg.Content != "" {
				parts = append(parts, existingMsg.Content)
			}
			if config.InboxGuideSystemRole != "" {
				parts = append(parts, config.InboxGuideSystemRole)
			}
			result[systemMsgIndex] = &schema.Message{
				Role:    schema.System,
				Content: strings.Join(parts, "\n\n"),
			}
		} else {
			// Create new system message
			systemMsg := schema.SystemMessage(config.InboxGuideSystemRole)
			newResult := make([]*schema.Message, 0, len(result)+1)
			newResult = append(newResult, systemMsg)
			newResult = append(newResult, result...)
			result = newResult
		}

		return result, nil
	})
}

