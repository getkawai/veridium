package processors

import (
	"context"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// HistoryTruncateConfig holds configuration for history truncation
type HistoryTruncateConfig struct {
	EnableHistoryCount bool
	HistoryCount       int
}

// NewHistoryTruncateLambda creates a lambda node for history truncation
func NewHistoryTruncateLambda(config HistoryTruncateConfig) *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, msgs []*schema.Message) ([]*schema.Message, error) {
		// If history count is not enabled or count is invalid, return all messages
		if !config.EnableHistoryCount || config.HistoryCount <= 0 {
			return msgs, nil
		}

		// If count is greater than or equal to message count, return all messages
		if config.HistoryCount >= len(msgs) {
			return msgs, nil
		}

		// Return last N messages
		return msgs[len(msgs)-config.HistoryCount:], nil
	})
}

