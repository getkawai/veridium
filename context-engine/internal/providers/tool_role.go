package providers

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// ToolSystemRoleConfig holds configuration for tool system role injection
type ToolSystemRoleConfig struct {
	GetToolSystemRoles func(tools []interface{}) string
	IsCanUseFC         func(model, provider string) bool
	Model              string
	Provider           string
	Tools              []interface{}
}

// NewToolSystemRoleProviderLambda creates a lambda node for tool system role injection
func NewToolSystemRoleProviderLambda(config ToolSystemRoleConfig) *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, msgs []*schema.Message) ([]*schema.Message, error) {
		// Check if tool system role should be injected
		toolSystemRole := getToolSystemRole(config)
		if toolSystemRole == "" {
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
			if toolSystemRole != "" {
				parts = append(parts, toolSystemRole)
			}
			result[systemMsgIndex] = &schema.Message{
				Role:    schema.System,
				Content: strings.Join(parts, "\n\n"),
			}
		} else {
			// Create new system message
			systemMsg := schema.SystemMessage(toolSystemRole)
			newResult := make([]*schema.Message, 0, len(result)+1)
			newResult = append(newResult, systemMsg)
			newResult = append(newResult, result...)
			result = newResult
		}

		return result, nil
	})
}

// getToolSystemRole determines if tool system role should be injected
func getToolSystemRole(config ToolSystemRoleConfig) string {
	// Check if there are tools
	if len(config.Tools) == 0 {
		return ""
	}

	// Check if function calling is supported
	if !config.IsCanUseFC(config.Model, config.Provider) {
		return ""
	}

	// Get tool system role
	if config.GetToolSystemRoles == nil {
		return ""
	}

	return config.GetToolSystemRoles(config.Tools)
}

