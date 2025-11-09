package memory

// SummaryBufferOption is a function for creating new summary buffer
// with other than the default values.
type SummaryBufferOption func(*ConversationSummaryBuffer)

// WithRecentMessagePairs sets how many recent conversation pairs to keep.
// Default is 3 pairs (6 messages).
func WithRecentMessagePairs(pairs int) SummaryBufferOption {
	return func(sb *ConversationSummaryBuffer) {
		if pairs < MinRecentMessagePairs {
			pairs = MinRecentMessagePairs
		}
		if pairs > MaxRecentMessagePairs {
			pairs = MaxRecentMessagePairs
		}
		sb.RecentMessagePairs = pairs
	}
}

// WithMaxTokensForRecent sets the maximum token budget for recent messages.
// Default is 1500 tokens.
func WithMaxTokensForRecent(tokens int) SummaryBufferOption {
	return func(sb *ConversationSummaryBuffer) {
		sb.MaxTokensForRecent = tokens
	}
}

// WithSummarizeThreshold sets when to trigger summarization (message pair count).
// Default is 6 pairs (12 messages).
func WithSummarizeThreshold(threshold int) SummaryBufferOption {
	return func(sb *ConversationSummaryBuffer) {
		sb.SummarizeThreshold = threshold
	}
}

// WithInitialSummary sets an initial summary (useful for restoring state).
func WithInitialSummary(summary string) SummaryBufferOption {
	return func(sb *ConversationSummaryBuffer) {
		sb.CurrentSummary = summary
	}
}
