package providers

import (
	"context"
	"fmt"
	"strings"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// HistorySummaryConfig holds configuration for history summary injection
type HistorySummaryConfig struct {
	HistorySummary      string
	FormatHistorySummary func(summary string) string
}

// DefaultHistorySummaryFormatter formats history summary with default template
func DefaultHistorySummaryFormatter(historySummary string) string {
	return fmt.Sprintf(`<chat_history_summary>
<docstring>Users may have lots of chat messages, here is the summary of the history:</docstring>
<summary>%s</summary>
</chat_history_summary>`, historySummary)
}

// NewHistorySummaryProviderLambda creates a lambda node for history summary injection
func NewHistorySummaryProviderLambda(config HistorySummaryConfig) *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, msgs []*schema.Message) ([]*schema.Message, error) {
		// Check if history summary exists
		if config.HistorySummary == "" {
			return msgs, nil
		}

		// Format history summary
		formatter := config.FormatHistorySummary
		if formatter == nil {
			formatter = DefaultHistorySummaryFormatter
		}
		formattedSummary := formatter(config.HistorySummary)

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
			if formattedSummary != "" {
				parts = append(parts, formattedSummary)
			}
			result[systemMsgIndex] = &schema.Message{
				Role:    schema.System,
				Content: strings.Join(parts, "\n\n"),
			}
		} else {
			// Create new system message
			systemMsg := schema.SystemMessage(formattedSummary)
			newResult := make([]*schema.Message, 0, len(result)+1)
			newResult = append(newResult, systemMsg)
			newResult = append(newResult, result...)
			result = newResult
		}

		return result, nil
	})
}

