package processors

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// MessageCleanupLambda removes unnecessary fields from messages
func MessageCleanupLambda() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, msgs []*schema.Message) ([]*schema.Message, error) {
		cleaned := make([]*schema.Message, len(msgs))
		for i, msg := range msgs {
			cleaned[i] = &schema.Message{
				Role:    msg.Role,
				Content: msg.Content,
				// Keep only essential fields
				Name:         msg.Name,
				ToolCalls:    msg.ToolCalls,
				ToolCallID:   msg.ToolCallID,
				ToolName:     msg.ToolName,
				ResponseMeta: msg.ResponseMeta,
			}

			// Preserve multimodal content if present
			if len(msg.UserInputMultiContent) > 0 {
				cleaned[i].UserInputMultiContent = msg.UserInputMultiContent
			}
			if len(msg.AssistantGenMultiContent) > 0 {
				cleaned[i].AssistantGenMultiContent = msg.AssistantGenMultiContent
			}
		}
		return cleaned, nil
	})
}

