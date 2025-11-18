package processors

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// ToolMessageReorderLambda reorders tool messages to ensure proper sequence
// Tool messages should come immediately after the assistant message with matching tool_calls
func ToolMessageReorderLambda() *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, msgs []*schema.Message) ([]*schema.Message, error) {
		if len(msgs) == 0 {
			return msgs, nil
		}

		// 1. Collect all valid tool_call_ids from assistant messages
		validToolCallIDs := make(map[string]bool)
		for _, msg := range msgs {
			if msg.Role == schema.Assistant && len(msg.ToolCalls) > 0 {
				for _, toolCall := range msg.ToolCalls {
					validToolCallIDs[toolCall.ID] = true
				}
			}
		}

		// 2. Collect all valid tool messages by tool_call_id
		toolMessages := make(map[string]*schema.Message)
		for _, msg := range msgs {
			if msg.Role == schema.Tool && msg.ToolCallID != "" {
				if validToolCallIDs[msg.ToolCallID] {
					toolMessages[msg.ToolCallID] = msg
				}
			}
		}

		// 3. Reorder messages
		result := make([]*schema.Message, 0, len(msgs))
		addedToolCallIDs := make(map[string]bool)

		for _, msg := range msgs {
			// Skip invalid tool messages
			if msg.Role == schema.Tool {
				if msg.ToolCallID == "" || !validToolCallIDs[msg.ToolCallID] {
					continue // Skip invalid tool message
				}
				// Skip if already added
				if addedToolCallIDs[msg.ToolCallID] {
					continue
				}
			}

			// Add the message
			result = append(result, msg)

			// If this is an assistant message with tool calls, add corresponding tool messages immediately after
			if msg.Role == schema.Assistant && len(msg.ToolCalls) > 0 {
				for _, toolCall := range msg.ToolCalls {
					if toolMsg, exists := toolMessages[toolCall.ID]; exists {
						result = append(result, toolMsg)
						addedToolCallIDs[toolCall.ID] = true
					}
				}
			} else if msg.Role == schema.Tool && msg.ToolCallID != "" {
				addedToolCallIDs[msg.ToolCallID] = true
			}
		}

		return result, nil
	})
}

