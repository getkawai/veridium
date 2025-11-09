package providers

import (
	"context"
	"strings"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// SystemRoleInjectorConfig holds configuration for system role injection
type SystemRoleInjectorConfig struct {
	SystemRole string
}

// NewSystemRoleInjectorLambda creates a lambda node for system role injection
func NewSystemRoleInjectorLambda(config SystemRoleInjectorConfig) *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, msgs []*schema.Message) ([]*schema.Message, error) {
		// Skip injection if no system role is configured
		if config.SystemRole == "" || strings.TrimSpace(config.SystemRole) == "" {
			return msgs, nil
		}

		// Check if system role already exists at the beginning
		if len(msgs) > 0 && msgs[0].Role == schema.System {
			return msgs, nil
		}

		// Inject system role at the beginning
		systemMsg := schema.SystemMessage(config.SystemRole)
		result := make([]*schema.Message, 0, len(msgs)+1)
		result = append(result, systemMsg)
		result = append(result, msgs...)

		return result, nil
	})
}

