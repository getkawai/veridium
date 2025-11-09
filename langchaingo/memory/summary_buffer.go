package memory

import (
	"context"
	"sync"

	"github.com/kawai-network/veridium/langchaingo/llms"
	"github.com/kawai-network/veridium/langchaingo/schema"
)

const (
	// Message retention constants
	DefaultRecentMessagePairs = 3  // Keep last 3 conversation pairs (6 messages)
	MinRecentMessagePairs     = 2  // Minimum 2 pairs for context
	MaxRecentMessagePairs     = 10 // Maximum 10 pairs to prevent bloat

	// Token budgets
	DefaultMaxTokensForRecent = 1500 // ~1.5K tokens for recent messages
	SummaryTokenBudget        = 500  // ~500 tokens for summary
	TotalContextBudget        = 2000 // Total budget (safe for most models)

	// Summarization triggers
	DefaultSummarizeThreshold = 6 // Summarize when > 6 pairs (12 messages)
)

// ConversationSummaryBuffer maintains a summary of old messages
// and keeps only recent messages for context.
// This prevents LLM confusion from very long conversations.
type ConversationSummaryBuffer struct {
	ConversationBuffer

	// LLM for generating summaries
	LLM llms.Model

	// How many recent conversation pairs to keep
	RecentMessagePairs int

	// Maximum tokens for recent messages
	MaxTokensForRecent int

	// When to trigger summarization (message pair count)
	SummarizeThreshold int

	// Current summary of old conversations
	CurrentSummary string

	// Mutex for thread-safe operations
	mu sync.RWMutex

	// Flag to prevent concurrent summarization
	summarizing bool
}

// Statically assert that ConversationSummaryBuffer implement the memory interface.
var _ schema.Memory = &ConversationSummaryBuffer{}

// NewConversationSummaryBuffer creates a new summary buffer memory.
func NewConversationSummaryBuffer(
	llm llms.Model,
	bufferOptions ...ConversationBufferOption,
) *ConversationSummaryBuffer {
	sb := &ConversationSummaryBuffer{
		LLM:                llm,
		RecentMessagePairs: DefaultRecentMessagePairs,
		MaxTokensForRecent: DefaultMaxTokensForRecent,
		SummarizeThreshold: DefaultSummarizeThreshold,
		CurrentSummary:     "",
		ConversationBuffer: *applyBufferOptions(bufferOptions...),
	}

	return sb
}

// NewConversationSummaryBufferWithOptions creates a new summary buffer memory with custom options.
func NewConversationSummaryBufferWithOptions(
	llm llms.Model,
	summaryOptions []SummaryBufferOption,
	bufferOptions ...ConversationBufferOption,
) *ConversationSummaryBuffer {
	sb := &ConversationSummaryBuffer{
		LLM:                llm,
		RecentMessagePairs: DefaultRecentMessagePairs,
		MaxTokensForRecent: DefaultMaxTokensForRecent,
		SummarizeThreshold: DefaultSummarizeThreshold,
		CurrentSummary:     "",
		ConversationBuffer: *applyBufferOptions(bufferOptions...),
	}

	// Apply summary-specific options
	for _, opt := range summaryOptions {
		opt(sb)
	}

	return sb
}

// MemoryVariables uses ConversationBuffer method for memory variables.
func (sb *ConversationSummaryBuffer) MemoryVariables(ctx context.Context) []string {
	return sb.ConversationBuffer.MemoryVariables(ctx)
}

// LoadMemoryVariables returns summary + recent messages.
// Format: [System: Summary] + [Recent N pairs of messages]
// This prevents LLM confusion by limiting context to relevant information.
func (sb *ConversationSummaryBuffer) LoadMemoryVariables(
	ctx context.Context,
	inputs map[string]any,
) (map[string]any, error) {
	sb.mu.RLock()
	defer sb.mu.RUnlock()

	// Get all messages
	allMessages, err := sb.ChatHistory.Messages(ctx)
	if err != nil {
		return nil, err
	}

	// If few messages, return all (no summarization needed)
	totalPairs := len(allMessages) / 2
	if totalPairs <= sb.RecentMessagePairs {
		return sb.ConversationBuffer.LoadMemoryVariables(ctx, inputs)
	}

	// Get recent messages
	recentMessages := sb.getRecentMessages(allMessages)

	// Build context with summary + recent messages
	var contextMessages []llms.ChatMessage

	// Add summary as system message if exists
	if sb.CurrentSummary != "" {
		contextMessages = append(contextMessages, llms.SystemChatMessage{
			Content: "Previous conversation summary:\n" + sb.CurrentSummary,
		})
	}

	// Add recent messages
	contextMessages = append(contextMessages, recentMessages...)

	if sb.ReturnMessages {
		return map[string]any{
			sb.MemoryKey: contextMessages,
		}, nil
	}

	// Convert to buffer string
	bufferString, err := llms.GetBufferString(
		contextMessages,
		sb.HumanPrefix,
		sb.AIPrefix,
	)
	if err != nil {
		return nil, err
	}

	return map[string]any{
		sb.MemoryKey: bufferString,
	}, nil
}

// getRecentMessages returns the last N conversation pairs.
func (sb *ConversationSummaryBuffer) getRecentMessages(
	messages []llms.ChatMessage,
) []llms.ChatMessage {
	totalMessages := len(messages)
	recentCount := sb.RecentMessagePairs * 2 // pairs = 2 messages each

	if totalMessages <= recentCount {
		return messages
	}

	// Get last N messages
	recentMessages := messages[totalMessages-recentCount:]

	// Token-based adjustment: if recent messages exceed token budget, reduce
	tokenCount := sb.estimateTokens(recentMessages)

	for tokenCount > sb.MaxTokensForRecent && len(recentMessages) > MinRecentMessagePairs*2 {
		// Remove oldest pair from recent messages
		recentMessages = recentMessages[2:]
		tokenCount = sb.estimateTokens(recentMessages)
	}

	return recentMessages
}

// estimateTokens estimates token count for messages.
func (sb *ConversationSummaryBuffer) estimateTokens(messages []llms.ChatMessage) int {
	bufferString, err := llms.GetBufferString(
		messages,
		sb.HumanPrefix,
		sb.AIPrefix,
	)
	if err != nil {
		return 0
	}

	return llms.CountTokens("", bufferString)
}

// SaveContext saves messages and triggers summarization if needed.
func (sb *ConversationSummaryBuffer) SaveContext(
	ctx context.Context,
	inputValues map[string]any,
	outputValues map[string]any,
) error {
	// Save messages normally
	err := sb.ConversationBuffer.SaveContext(ctx, inputValues, outputValues)
	if err != nil {
		return err
	}

	// Check if summarization is needed
	messages, err := sb.ChatHistory.Messages(ctx)
	if err != nil {
		return err
	}

	totalPairs := len(messages) / 2

	// Trigger background summarization if threshold exceeded
	if totalPairs > sb.SummarizeThreshold {
		go sb.summarizeInBackground(context.Background())
	}

	return nil
}

// summarizeInBackground performs summarization in a goroutine.
// This is non-blocking and won't slow down the chat response.
func (sb *ConversationSummaryBuffer) summarizeInBackground(ctx context.Context) {
	sb.mu.Lock()

	// Prevent concurrent summarization
	if sb.summarizing {
		sb.mu.Unlock()
		return
	}
	sb.summarizing = true
	sb.mu.Unlock()

	defer func() {
		sb.mu.Lock()
		sb.summarizing = false
		sb.mu.Unlock()
	}()

	// Get all messages
	allMessages, err := sb.ChatHistory.Messages(ctx)
	if err != nil {
		return
	}

	totalPairs := len(allMessages) / 2
	if totalPairs <= sb.RecentMessagePairs {
		return // No summarization needed
	}

	// Messages to summarize (all except recent)
	recentCount := sb.RecentMessagePairs * 2
	messagesToSummarize := allMessages[:len(allMessages)-recentCount]

	// Generate new summary
	newSummary, err := sb.generateSummary(ctx, messagesToSummarize)
	if err != nil {
		return
	}

	// Update summary
	sb.mu.Lock()
	sb.CurrentSummary = newSummary
	sb.mu.Unlock()
}

// generateSummary creates a summary using LLM.
func (sb *ConversationSummaryBuffer) generateSummary(
	ctx context.Context,
	messages []llms.ChatMessage,
) (string, error) {
	// Convert messages to text
	bufferString, err := llms.GetBufferString(
		messages,
		sb.HumanPrefix,
		sb.AIPrefix,
	)
	if err != nil {
		return "", err
	}

	// Create prompt for summarization
	prompt := `Summarize the following conversation concisely. 
Focus on key topics, decisions, and important context that would be useful for continuing the conversation.
Keep the summary under 400 tokens.

Conversation:
` + bufferString

	// If existing summary, refine it
	if sb.CurrentSummary != "" {
		prompt = `Existing summary:
` + sb.CurrentSummary + `

New conversation to incorporate:
` + bufferString + `

Provide an updated summary that incorporates the new information while maintaining key context from the existing summary.
Keep the summary under 400 tokens.`
	}

	// Generate summary using the LLM
	result, err := llms.GenerateFromSinglePrompt(
		ctx,
		sb.LLM,
		prompt,
		llms.WithMaxTokens(500),
		llms.WithTemperature(0.3), // Lower temperature for consistent summaries
	)

	return result, err
}

// Clear clears both summary and messages.
func (sb *ConversationSummaryBuffer) Clear(ctx context.Context) error {
	sb.mu.Lock()
	sb.CurrentSummary = ""
	sb.mu.Unlock()

	return sb.ConversationBuffer.Clear(ctx)
}

// GetMemoryKey returns the memory key.
func (sb *ConversationSummaryBuffer) GetMemoryKey(ctx context.Context) string {
	return sb.ConversationBuffer.GetMemoryKey(ctx)
}

// GetSummary returns the current summary (useful for debugging/UI).
func (sb *ConversationSummaryBuffer) GetSummary() string {
	sb.mu.RLock()
	defer sb.mu.RUnlock()
	return sb.CurrentSummary
}

// SetRecentMessagePairs allows dynamic adjustment of recent message count.
func (sb *ConversationSummaryBuffer) SetRecentMessagePairs(pairs int) {
	if pairs < MinRecentMessagePairs {
		pairs = MinRecentMessagePairs
	}
	if pairs > MaxRecentMessagePairs {
		pairs = MaxRecentMessagePairs
	}

	sb.mu.Lock()
	sb.RecentMessagePairs = pairs
	sb.mu.Unlock()
}

// SetSummarizeThreshold allows dynamic adjustment of summarization trigger.
func (sb *ConversationSummaryBuffer) SetSummarizeThreshold(threshold int) {
	sb.mu.Lock()
	sb.SummarizeThreshold = threshold
	sb.mu.Unlock()
}
