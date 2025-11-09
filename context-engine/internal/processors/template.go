package processors

import (
	"context"
	"regexp"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// InputTemplateConfig holds configuration for input template processing
type InputTemplateConfig struct {
	InputTemplate string
}

// NewInputTemplateLambda creates a lambda node for input template processing
func NewInputTemplateLambda(config InputTemplateConfig) *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, msgs []*schema.Message) ([]*schema.Message, error) {
		// Skip processing if no template is configured
		if config.InputTemplate == "" {
			return msgs, nil
		}

		result := make([]*schema.Message, len(msgs))
		copy(result, msgs)

		// Compile template pattern: {{text}}
		templatePattern := regexp.MustCompile(`{{\s*text\s*}}`)

		// Process each user message
		for i, msg := range result {
			if msg.Role == schema.User {
				// Handle string content
				if msg.Content != "" {
					originalContent := msg.Content
					processedContent := templatePattern.ReplaceAllStringFunc(config.InputTemplate, func(match string) string {
						return originalContent
					})

					if processedContent != originalContent {
						result[i] = &schema.Message{
							Role:    msg.Role,
							Content: processedContent,
							Name:    msg.Name,
						}
						// Preserve multimodal content if present
						if len(msg.UserInputMultiContent) > 0 {
							result[i].UserInputMultiContent = msg.UserInputMultiContent
						}
					}
				}
			}
		}

		return result, nil
	})
}

