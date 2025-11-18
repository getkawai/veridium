package processors

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// MessageContentConfig holds configuration for message content processing
type MessageContentConfig struct {
	FileContext    FileContextConfig
	IsCanUseVideo  func(model, provider string) bool
	IsCanUseVision func(model, provider string) bool
	Model          string
	Provider       string
}

// FileContextConfig holds file context configuration
type FileContextConfig struct {
	Enabled        bool
	IncludeFileURL bool
}

// NewMessageContentLambda creates a lambda node for message content processing
// This is a simplified version - full implementation with file context and image/video processing
// will be added in later phases
func NewMessageContentLambda(config MessageContentConfig) *compose.Lambda {
	return compose.InvokableLambda(func(ctx context.Context, msgs []*schema.Message) ([]*schema.Message, error) {
		result := make([]*schema.Message, len(msgs))
		for i, msg := range msgs {
			switch msg.Role {
			case schema.User:
				result[i] = processUserMessageContent(msg, config)
			case schema.Assistant:
				result[i] = processAssistantMessageContent(msg, config)
			default:
				result[i] = msg
			}
		}
		return result, nil
	})
}

// processUserMessageContent processes user message content
func processUserMessageContent(msg *schema.Message, config MessageContentConfig) *schema.Message {
	// For now, return message as-is
	// Full implementation will handle:
	// - File context injection
	// - Image processing (convert to base64 if needed)
	// - Video processing
	// - Structured content parts
	return msg
}

// processAssistantMessageContent processes assistant message content
func processAssistantMessageContent(msg *schema.Message, config MessageContentConfig) *schema.Message {
	// Handle reasoning content if present
	if msg.ReasoningContent != "" {
		// Create structured content with thinking part
		// For now, keep reasoning in ReasoningContent field
		return msg
	}

	// For now, return message as-is
	// Full implementation will handle:
	// - Image content in assistant messages
	// - Structured content parts
	return msg
}

// formatFileContext formats file context prompt
func formatFileContext(fileList []string, imageList []string, videoList []string, includeURL bool) string {
	parts := []string{}

	if len(fileList) > 0 {
		fileContext := "Files attached:\n"
		for i, file := range fileList {
			if includeURL {
				fileContext += fmt.Sprintf("%d. %s\n", i+1, file)
			} else {
				fileContext += fmt.Sprintf("%d. [File attached]\n", i+1)
			}
		}
		parts = append(parts, fileContext)
	}

	if len(imageList) > 0 {
		imageContext := "Images attached:\n"
		for i, img := range imageList {
			if includeURL {
				imageContext += fmt.Sprintf("%d. %s\n", i+1, img)
			} else {
				imageContext += fmt.Sprintf("%d. [Image attached]\n", i+1)
			}
		}
		parts = append(parts, imageContext)
	}

	if len(videoList) > 0 {
		videoContext := "Videos attached:\n"
		for i, vid := range videoList {
			if includeURL {
				videoContext += fmt.Sprintf("%d. %s\n", i+1, vid)
			} else {
				videoContext += fmt.Sprintf("%d. [Video attached]\n", i+1)
			}
		}
		parts = append(parts, videoContext)
	}

	if len(parts) > 0 {
		return "\n\n" + fmt.Sprintf("<file_context>\n%s</file_context>", fmt.Sprintf("%s", parts[0]))
	}

	return ""
}

