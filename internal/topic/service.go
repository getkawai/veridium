package topic

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/internal/database"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/types"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// TopicService handles topic-related operations including title generation
type TopicService struct {
	db         *database.Service
	app        *application.App
	titleModel fantasy.LanguageModel
}

// NewService creates a new Topic Service
func NewService(db *database.Service, app *application.App) *TopicService {
	return &TopicService{
		db:  db,
		app: app,
	}
}

// SetTitleModel sets the model used for title generation
func (s *TopicService) SetTitleModel(model fantasy.LanguageModel) {
	s.titleModel = model
	if model != nil {
		log.Printf("✅ TopicService: Title model set (%s/%s)", model.Provider(), model.Model())
	}
}

// GenerateTitle generates a concise title for the conversation using titleModel
func (s *TopicService) GenerateTitle(ctx context.Context, messages []fantasy.Message, locale string) (string, error) {
	if len(messages) == 0 {
		return "New Conversation", nil
	}

	// Check if title model is configured
	if s.titleModel == nil {
		log.Printf("⚠️  Title model not configured, using default title")
		return "New Conversation", nil
	}

	// Build summary prompt
	systemPrompt := fmt.Sprintf(`You are a professional conversation summarizer. Generate a concise title that captures the essence of the conversation.

Rules:
- Maximum 10 words
- Maximum 50 characters
- No punctuation marks, quotes, or special characters
- Use the language specified by the locale code: %s
- Output ONLY the title text, nothing else

Example: Sleep Functions for Body and Mind`, locale)

	// Build conversation text (User messages only)
	var conversationText string
	for _, msg := range messages {
		if types.GetMessageRole(msg) == "user" {
			conversationText += fmt.Sprintf("user: %s\n", types.GetMessageText(msg))
		}
	}

	// Fallback: if no user messages found (unlikely but possible), use all messages
	if conversationText == "" {
		for _, msg := range messages {
			conversationText += fmt.Sprintf("%s: %s\n", types.GetMessageRole(msg), types.GetMessageText(msg))
		}
	}

	log.Printf("📝 Generating title for conversation (%d messages, %d chars)", len(messages), len(conversationText))

	// Create messages for title generation
	titleMessages := fantasy.Prompt{
		fantasy.NewSystemMessage(systemPrompt),
		fantasy.NewUserMessage(conversationText),
	}

	// Use titleModel directly (ChainLanguageModel handles fallback)
	resp, err := s.titleModel.Generate(ctx, fantasy.Call{Prompt: titleMessages})
	if err != nil {
		log.Printf("⚠️  Title generation failed: %v, using default", err)
		return "New Conversation", nil
	}

	responseContent := resp.Content.Text()

	// Strip <think> tags if present
	if strings.Contains(responseContent, "<think>") {
		responseContent = s.stripThinkTags(responseContent)
	}

	log.Printf("📝 Raw title response: %q", responseContent)

	// Clean up the title
	title := strings.TrimSpace(responseContent)

	// If title is empty, try to extract from the original response
	if title == "" {
		// Try to find text in quotes
		quotePattern := regexp.MustCompile(`["']([^"']+)["']`)
		matches := quotePattern.FindStringSubmatch(responseContent)
		if len(matches) > 1 {
			title = matches[1]
		} else {
			// Fallback: use first non-empty line
			lines := strings.Split(responseContent, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					title = line
					break
				}
			}
		}
	}

	// Final fallback if still empty
	if title == "" {
		title = "New Conversation"
	}

	// Truncate to 50 characters
	if len(title) > 50 {
		title = title[:50]
	}

	return title, nil
}

// GenerateTitleFromPrompt generates a title from a single prompt (e.g., image generation prompt)
func (s *TopicService) GenerateTitleFromPrompt(ctx context.Context, prompt string, locale string) (string, error) {
	// Wrap prompt in a user message
	messages := []fantasy.Message{
		fantasy.NewUserMessage(prompt),
	}
	return s.GenerateTitle(ctx, messages, locale)
}

// UpdateTopicTitleFromPrompt updates an existing topic with LLM-generated title from a prompt
func (s *TopicService) UpdateTopicTitleFromPrompt(ctx context.Context, topicID string, prompt string) error {
	// Wrap prompt in a user message
	messages := []fantasy.Message{
		fantasy.NewUserMessage(prompt),
	}
	return s.UpdateTopicTitle(ctx, topicID, messages)
}

// UpdateTopicTitle updates an existing topic with LLM-generated title
// It runs in the background
func (s *TopicService) UpdateTopicTitle(ctx context.Context, topicID string, messages []fantasy.Message) error {
	// Create a copy of messages to avoid race conditions
	messagesCopy := make([]fantasy.Message, len(messages))
	copy(messagesCopy, messages)

	log.Printf("📌 [TITLE] UpdateTopicTitle called for topic %s with %d messages", topicID, len(messagesCopy))

	// Run in background
	go func() {
		// Add a small delay
		time.Sleep(2 * time.Second)

		log.Printf("🔄 Generating title in background for topic %s...", topicID)

		// Create context with timeout
		bgCtx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		// Generate title (default locale: en-US)
		title, err := s.GenerateTitle(bgCtx, messagesCopy, "en-US")
		if err != nil {
			log.Printf("⚠️  Warning: Failed to generate topic title: %v", err)
			title = "New Conversation"
		}

		// Fetch existing topic first to preserve summary/metadata
		existingTopic, err := s.db.Queries().GetTopic(bgCtx, topicID)
		if err != nil {
			log.Printf("⚠️  Failed to fetch topic for update: %v", err)
			return
		}

		// Update topic title in database
		now := time.Now().UnixMilli()
		_, err = s.db.Queries().UpdateTopic(bgCtx, db.UpdateTopicParams{
			Title:          sql.NullString{String: title, Valid: true},
			HistorySummary: existingTopic.HistorySummary, // Preserve existing
			Metadata:       existingTopic.Metadata,       // Preserve existing
			UpdatedAt:      now,
			ID:             topicID,
		})

		if err != nil {
			log.Printf("⚠️  Failed to update topic title in DB: %v", err)
		} else {
			log.Printf("✅ Updated topic %s with title: %s", topicID, title)

			// Emit event to notify UI
			if s.app != nil {
				s.app.Event.Emit("chat:topic:updated", map[string]interface{}{
					"topic_id": topicID,
					"title":    title,
				})
			}
		}
	}()

	return nil
}

// UpdateGenerationTopicTitleFromPrompt updates an existing generation topic with LLM-generated title from a prompt
func (s *TopicService) UpdateGenerationTopicTitleFromPrompt(ctx context.Context, topicID string, prompt string) error {
	// Wrap prompt in a user message
	messages := []fantasy.Message{
		fantasy.NewUserMessage(prompt),
	}
	// Use UpdateGenerationTopicTitle instead of UpdateTopicTitle
	return s.UpdateGenerationTopicTitle(ctx, topicID, messages)
}

// UpdateGenerationTopicTitle updates an existing generation topic with LLM-generated title
// It runs in the background. It targets the 'generation_topics' table.
func (s *TopicService) UpdateGenerationTopicTitle(ctx context.Context, topicID string, messages []fantasy.Message) error {
	// Create a copy of messages to avoid race conditions
	messagesCopy := make([]fantasy.Message, len(messages))
	copy(messagesCopy, messages)

	log.Printf("📌 [TITLE] UpdateGenerationTopicTitle called for topic %s with %d messages", topicID, len(messagesCopy))

	// Run in background
	go func() {
		// Add a small delay
		time.Sleep(2 * time.Second)

		log.Printf("🔄 Generating title in background for generation topic %s...", topicID)

		// Create context with timeout
		bgCtx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
		defer cancel()

		// Generate title (default locale: en-US)
		title, err := s.GenerateTitle(bgCtx, messagesCopy, "en-US")
		if err != nil {
			log.Printf("⚠️  Warning: Failed to generate topic title: %v", err)
			title = "New Topic"
		}

		// Fetch existing generation topic first to preserve cover_url
		existingTopic, err := s.db.Queries().GetGenerationTopic(bgCtx, topicID)
		if err != nil {
			log.Printf("⚠️  Failed to fetch generation topic for update: %v", err)
			return
		}

		// Update topic title in database (generation_topics table)
		// We MUST preserve CoverUrl
		now := time.Now().UnixMilli()
		_, err = s.db.Queries().UpdateGenerationTopic(bgCtx, db.UpdateGenerationTopicParams{
			Title:     sql.NullString{String: title, Valid: true},
			CoverUrl:  existingTopic.CoverUrl, // Preserve existing cover URL
			UpdatedAt: now,
			ID:        topicID,
		})

		if err != nil {
			log.Printf("⚠️  Failed to update generation topic title in DB: %v", err)
		} else {
			log.Printf("✅ Updated generation topic %s with title: %s", topicID, title)

			// Emit event to notify UI
			if s.app != nil {
				s.app.Event.Emit("generation:topic:updated", map[string]interface{}{
					"topic_id": topicID,
					"title":    title,
				})
			}
		}
	}()

	return nil
}

// stripThinkTags removes <think>...</think> blocks from the text using regex
func (s *TopicService) stripThinkTags(text string) string {
	// First, try to remove complete <think>...</think> blocks
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	cleaned := re.ReplaceAllString(text, "")

	// If <think> tag still exists (incomplete/unclosed), remove everything from <think> onwards
	if strings.Contains(cleaned, "<think>") {
		idx := strings.Index(cleaned, "<think>")
		cleaned = cleaned[:idx]
	}

	return strings.TrimSpace(cleaned)
}
