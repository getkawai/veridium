package processors

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// ToolCallConfig holds configuration for tool call processing
type ToolCallConfig struct {
	GenToolCallingName func(identifier, apiName, toolType string) string
	IsCanUseFC         func(model, provider string) bool
	Model              string
	Provider           string
}

// DefaultGenToolCallingName generates default tool calling name
func DefaultGenToolCallingName(identifier, apiName, toolType string) string {
	return fmt.Sprintf("%s.%s", identifier, apiName)
}

// NewToolCallLambda creates a lambda node for tool call processing
func NewToolCallLambda(config ToolCallConfig) *compose.Lambda {
	genName := config.GenToolCallingName
	if genName == nil {
		genName = DefaultGenToolCallingName
	}

	return compose.InvokableLambda(func(ctx context.Context, msgs []*schema.Message) ([]*schema.Message, error) {
		supportTools := true
		if config.IsCanUseFC != nil {
			supportTools = config.IsCanUseFC(config.Model, config.Provider)
		}

		result := make([]*schema.Message, len(msgs))
		for i, msg := range msgs {
			switch msg.Role {
			case schema.Assistant:
				result[i] = processAssistantToolCalls(msg, supportTools, genName)
			case schema.Tool:
				result[i] = processToolMessage(msg, supportTools, genName)
			default:
				result[i] = msg
			}
		}

		return result, nil
	})
}

// processAssistantToolCalls converts tools to tool_calls format for assistant messages
func processAssistantToolCalls(msg *schema.Message, supportTools bool, genName func(string, string, string) string) *schema.Message {
	// If tools are not supported or message has no tools, return message without tool-related fields
	// Note: We can't easily remove fields in Go, so we create a new message with only needed fields
	if !supportTools {
		return &schema.Message{
			Role:                  msg.Role,
			Content:               msg.Content,
			ResponseMeta:          msg.ResponseMeta,
			UserInputMultiContent: msg.UserInputMultiContent,
			AssistantGenMultiContent: msg.AssistantGenMultiContent,
		}
	}

	// For now, return message as-is since Eino schema.Message already has ToolCalls field
	// The actual conversion from custom tools format to ToolCalls should happen before
	// messages enter the pipeline, or we need to handle it based on Extra field
	return msg
}

// processToolMessage processes tool messages
func processToolMessage(msg *schema.Message, supportTools bool, genName func(string, string, string) string) *schema.Message {
	if !supportTools {
		// Convert tool message to user message if tools not supported
		return &schema.Message{
			Role:    schema.User,
			Content: msg.Content,
		}
	}

	// Tool message is already in correct format
	return msg
}

