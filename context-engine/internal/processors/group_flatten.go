package processors

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/kawai-network/veridium/context-engine/internal/types"
)

// GroupMessageFlattenConfig holds configuration for GroupMessageFlatten processor
type GroupMessageFlattenConfig struct {
	// No specific configuration needed for this processor
}

// NewGroupMessageFlattenLambda creates a lambda node that flattens group messages
// into standard assistant + tool message sequences
func NewGroupMessageFlattenLambda(config GroupMessageFlattenConfig) *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, messages []*schema.Message) ([]*schema.Message, error) {
		result := make([]*schema.Message, 0, len(messages))

		groupMessagesFlattened := 0
		assistantMessagesCreated := 0
		toolMessagesCreated := 0

		for _, msg := range messages {
			// Check if this is a group message
			// In schema.Message, we need to check the Role and any custom metadata
			if string(msg.Role) == "group" {
				// Get children from message metadata
				// Since schema.Message doesn't have Children field, we need to use Extra
				if msg.Extra == nil {
					// No children, skip this message
					continue
				}

				childrenData, ok := msg.Extra["children"]
				if !ok {
					// No children, skip
					continue
				}

				children, ok := childrenData.([]types.GroupChild)
				if !ok || len(children) == 0 {
					// Invalid or empty children, skip
					continue
				}

				groupMessagesFlattened++

				// Flatten each child
				for _, child := range children {
					// 1. Create assistant message from child
					assistantMsg := &schema.Message{
						Role:    schema.Assistant,
						Content: child.Content,
						Extra:   make(map[string]interface{}),
					}

					// Add tools if present (excluding result)
					if len(child.Tools) > 0 {
						tools := make([]map[string]interface{}, 0, len(child.Tools))
						for _, tool := range child.Tools {
							tools = append(tools, map[string]interface{}{
								"id":         tool.ID,
								"type":       tool.Type,
								"apiName":    tool.APIName,
								"identifier": tool.Identifier,
								"arguments":  tool.Arguments,
							})
						}
						assistantMsg.Extra["tools"] = tools
					}

					// Preserve metadata from original message
					if msg.Extra != nil {
						if reasoning, ok := msg.Extra["reasoning"]; ok {
							assistantMsg.Extra["reasoning"] = reasoning
						}
						if parentID, ok := msg.Extra["parentId"]; ok {
							assistantMsg.Extra["parentId"] = parentID
						}
						if threadID, ok := msg.Extra["threadId"]; ok {
							assistantMsg.Extra["threadId"] = threadID
						}
						if groupID, ok := msg.Extra["groupId"]; ok {
							assistantMsg.Extra["groupId"] = groupID
						}
						if agentID, ok := msg.Extra["agentId"]; ok {
							assistantMsg.Extra["agentId"] = agentID
						}
						if targetID, ok := msg.Extra["targetId"]; ok {
							assistantMsg.Extra["targetId"] = targetID
						}
						if topicID, ok := msg.Extra["topicId"]; ok {
							assistantMsg.Extra["topicId"] = topicID
						}
					}

					result = append(result, assistantMsg)
					assistantMessagesCreated++

					// 2. Create tool messages for each tool that has a result
					for _, tool := range child.Tools {
						if tool.Result != nil {
							toolMsg := &schema.Message{
								Role:    schema.Tool,
								Content: tool.Result.Content,
								Extra:   make(map[string]interface{}),
							}

							// Add plugin information
							toolMsg.Extra["plugin"] = map[string]interface{}{
								"id":         tool.ID,
								"type":       tool.Type,
								"apiName":    tool.APIName,
								"identifier": tool.Identifier,
								"arguments":  tool.Arguments,
							}

							if tool.Result.Error != "" {
								toolMsg.Extra["pluginError"] = tool.Result.Error
							}
							if tool.Result.State != "" {
								toolMsg.Extra["pluginState"] = tool.Result.State
							}

							toolMsg.Extra["tool_call_id"] = tool.ID

							// Preserve parent message references
							if msg.Extra != nil {
								if parentID, ok := msg.Extra["parentId"]; ok {
									toolMsg.Extra["parentId"] = parentID
								}
								if threadID, ok := msg.Extra["threadId"]; ok {
									toolMsg.Extra["threadId"] = threadID
								}
								if groupID, ok := msg.Extra["groupId"]; ok {
									toolMsg.Extra["groupId"] = groupID
								}
								if topicID, ok := msg.Extra["topicId"]; ok {
									toolMsg.Extra["topicId"] = topicID
								}
							}

							result = append(result, toolMsg)
							toolMessagesCreated++
						}
					}
				}
			} else {
				// Non-group message, keep as-is
				result = append(result, msg)
			}
		}

		// Log statistics
		if groupMessagesFlattened > 0 {
			fmt.Printf("[GroupMessageFlatten] Flattened %d group messages, created %d assistant messages and %d tool messages\n",
				groupMessagesFlattened, assistantMessagesCreated, toolMessagesCreated)
		}

		return result, nil
	})
}
