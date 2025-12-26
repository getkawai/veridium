package gateway

import (
	"fmt"
	"strings"

	"github.com/kawai-network/veridium/pkg/fantasy"
)

// RequestToPrompt converts an OpenAI ChatCompletionRequest to fantasy.Prompt.
func RequestToPrompt(req *ChatCompletionRequest) (fantasy.Prompt, error) {
	prompt := make(fantasy.Prompt, 0, len(req.Messages))

	for _, msg := range req.Messages {
		fantasyMsg, err := messageToFantasy(msg)
		if err != nil {
			return nil, fmt.Errorf("failed to convert message: %w", err)
		}
		prompt = append(prompt, fantasyMsg)
	}

	return prompt, nil
}

// messageToFantasy converts a single ChatMessage to fantasy.Message.
func messageToFantasy(msg ChatMessage) (fantasy.Message, error) {
	role := convertRole(msg.Role)

	var content []fantasy.MessagePart

	// Handle reasoning content (for assistant messages)
	if msg.ReasoningContent != "" && role == fantasy.MessageRoleAssistant {
		content = append(content, fantasy.ReasoningPart{
			Text: msg.ReasoningContent,
		})
	}

	// Handle tool calls (assistant message with tool_calls)
	if len(msg.ToolCalls) > 0 {
		for _, tc := range msg.ToolCalls {
			content = append(content, fantasy.ToolCallPart{
				ToolCallID: tc.ID,
				ToolName:   tc.Function.Name,
				Input:      tc.Function.Arguments,
			})
		}
	}

	// Handle tool result (tool message)
	if role == fantasy.MessageRoleTool && msg.ToolCallID != "" {
		text := msg.Content.GetText()
		content = append(content, fantasy.ToolResultPart{
			ToolCallID: msg.ToolCallID,
			Output:     fantasy.ToolResultOutputContentText{Text: text},
		})
		return fantasy.Message{Role: role, Content: content}, nil
	}

	// Handle content
	if msg.Content.Text != "" {
		// Simple text content
		content = append(content, fantasy.TextPart{Text: msg.Content.Text})
	} else if len(msg.Content.Parts) > 0 {
		// Multi-part content
		for _, part := range msg.Content.Parts {
			fantasyPart, err := contentPartToFantasy(part)
			if err != nil {
				// Skip unsupported parts with warning (don't fail)
				continue
			}
			if fantasyPart != nil {
				content = append(content, fantasyPart)
			}
		}
	}

	// Ensure at least empty text for messages without content
	if len(content) == 0 {
		content = append(content, fantasy.TextPart{Text: ""})
	}

	return fantasy.Message{
		Role:    role,
		Content: content,
	}, nil
}

// contentPartToFantasy converts a ContentPart to fantasy.MessagePart.
func contentPartToFantasy(part ContentPart) (fantasy.MessagePart, error) {
	switch part.Type {
	case "text":
		return fantasy.TextPart{Text: part.Text}, nil

	case "image_url":
		if part.ImageURL == nil {
			return nil, fmt.Errorf("image_url part missing image_url field")
		}
		// Extract image data from URL or data URI
		url := part.ImageURL.URL
		var data []byte
		var mediaType string

		if strings.HasPrefix(url, "data:") {
			// Parse data URL: data:image/png;base64,xxxx
			parts := strings.SplitN(url, ",", 2)
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid data URL format")
			}
			// Extract media type
			mediaInfo := strings.TrimPrefix(parts[0], "data:")
			mediaInfo = strings.TrimSuffix(mediaInfo, ";base64")
			mediaType = mediaInfo

			// Decode base64 - store as raw base64 string for now
			data = []byte(parts[1])
		} else {
			// External URL - store URL as data for later fetching
			mediaType = "image/url"
			data = []byte(url)
		}

		return fantasy.FilePart{
			MediaType: mediaType,
			Data:      data,
		}, nil

	case "input_audio":
		if part.InputAudio == nil {
			return nil, fmt.Errorf("input_audio part missing input_audio field")
		}
		mediaType := "audio/wav"
		if part.InputAudio.Format == "mp3" {
			mediaType = "audio/mpeg"
		}
		return fantasy.FilePart{
			MediaType: mediaType,
			Data:      []byte(part.InputAudio.Data),
		}, nil

	case "file":
		if part.File == nil {
			return nil, fmt.Errorf("file part missing file field")
		}
		// Handle file_id or file_data
		if part.File.FileID != "" {
			return fantasy.FilePart{
				MediaType: "application/pdf",
				Data:      []byte(part.File.FileID),
				Filename:  part.File.Filename,
			}, nil
		}
		if part.File.FileData != "" {
			// Parse data URL
			data := part.File.FileData
			mediaType := "application/octet-stream"
			if strings.HasPrefix(data, "data:") {
				parts := strings.SplitN(data, ",", 2)
				if len(parts) == 2 {
					mediaInfo := strings.TrimPrefix(parts[0], "data:")
					mediaInfo = strings.TrimSuffix(mediaInfo, ";base64")
					mediaType = mediaInfo
					data = parts[1]
				}
			}
			return fantasy.FilePart{
				MediaType: mediaType,
				Data:      []byte(data),
				Filename:  part.File.Filename,
			}, nil
		}
		return nil, fmt.Errorf("file part has no file_id or file_data")

	default:
		return nil, fmt.Errorf("unsupported content part type: %s", part.Type)
	}
}

// convertRole converts OpenAI role string to fantasy.MessageRole.
func convertRole(role string) fantasy.MessageRole {
	switch role {
	case "system":
		return fantasy.MessageRoleSystem
	case "user":
		return fantasy.MessageRoleUser
	case "assistant":
		return fantasy.MessageRoleAssistant
	case "tool":
		return fantasy.MessageRoleTool
	default:
		return fantasy.MessageRoleUser
	}
}

// RequestToTools converts OpenAI tools to fantasy.Tool slice.
func RequestToTools(tools []Tool) []fantasy.Tool {
	if len(tools) == 0 {
		return nil
	}

	result := make([]fantasy.Tool, 0, len(tools))
	for _, t := range tools {
		if t.Type != "function" {
			continue
		}
		result = append(result, fantasy.FunctionTool{
			Name:        t.Function.Name,
			Description: t.Function.Description,
			InputSchema: t.Function.Parameters,
		})
	}
	return result
}

// ResponseToOpenAI converts fantasy.Response to OpenAI response message.
func ResponseToOpenAI(resp *fantasy.Response) *ResponseMessage {
	msg := &ResponseMessage{
		Role: "assistant",
	}

	var textContent []string
	var toolCalls []ToolCall

	for _, content := range resp.Content {
		switch c := content.(type) {
		case fantasy.TextContent:
			textContent = append(textContent, c.Text)
		case fantasy.ReasoningContent:
			msg.ReasoningContent = c.Text
		case fantasy.ToolCallContent:
			toolCalls = append(toolCalls, ToolCall{
				ID:   c.ToolCallID,
				Type: "function",
				Function: ToolCallFunction{
					Name:      c.ToolName,
					Arguments: c.Input,
				},
			})
		}
	}

	msg.Content = strings.Join(textContent, "")
	if len(toolCalls) > 0 {
		msg.ToolCalls = toolCalls
	}

	return msg
}

// ConvertFinishReason converts fantasy.FinishReason to OpenAI finish_reason.
func ConvertFinishReason(reason fantasy.FinishReason) string {
	switch reason {
	case fantasy.FinishReasonStop:
		return "stop"
	case fantasy.FinishReasonLength:
		return "length"
	case fantasy.FinishReasonContentFilter:
		return "content_filter"
	case fantasy.FinishReasonToolCalls:
		return "tool_calls"
	default:
		return "stop"
	}
}

// PromptToText converts fantasy.Prompt to a simple text prompt for LLM.
// This is used when the LLM only supports text input.
func PromptToText(prompt fantasy.Prompt) string {
	var sb strings.Builder

	for _, msg := range prompt {
		role := string(msg.Role)
		text := fantasy.GetMessageText(msg)

		switch msg.Role {
		case fantasy.MessageRoleSystem:
			sb.WriteString(fmt.Sprintf("System: %s\n", text))
		case fantasy.MessageRoleUser:
			sb.WriteString(fmt.Sprintf("User: %s\n", text))
		case fantasy.MessageRoleAssistant:
			sb.WriteString(fmt.Sprintf("Assistant: %s\n", text))
		case fantasy.MessageRoleTool:
			// Format tool results
			for _, part := range msg.Content {
				if tr, ok := part.(fantasy.ToolResultPart); ok {
					if textOutput, ok := tr.Output.(fantasy.ToolResultOutputContentText); ok {
						sb.WriteString(fmt.Sprintf("Tool Result: %s\n", textOutput.Text))
					}
				}
			}
		default:
			sb.WriteString(fmt.Sprintf("%s: %s\n", role, text))
		}
	}

	sb.WriteString("Assistant:")
	return sb.String()
}

// RequestToLlamaMessages converts OpenAI request messages to llamalib Message format.
func RequestToLlamaMessages(messages []ChatMessage) []LlamaMessage {
	result := make([]LlamaMessage, len(messages))
	for i, msg := range messages {
		lm := LlamaMessage{
			Role:       msg.Role,
			Content:    msg.Content.GetText(),
			ToolCallID: msg.ToolCallID,
		}

		// Convert tool calls
		for _, tc := range msg.ToolCalls {
			lm.ToolCalls = append(lm.ToolCalls, LlamaToolCall{
				ID:        tc.ID,
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			})
		}

		result[i] = lm
	}
	return result
}

// RequestToLlamaTools converts OpenAI tools to llamalib Tool format.
func RequestToLlamaTools(tools []Tool) []LlamaTool {
	if len(tools) == 0 {
		return nil
	}

	result := make([]LlamaTool, 0, len(tools))
	for _, t := range tools {
		if t.Type != "function" {
			continue
		}
		result = append(result, LlamaTool{
			Name:        t.Function.Name,
			Description: t.Function.Description,
			Parameters:  t.Function.Parameters,
		})
	}
	return result
}

// LlamaMessage represents a message for llamalib.
type LlamaMessage struct {
	Role       string
	Content    string
	ToolCalls  []LlamaToolCall
	ToolCallID string
}

// LlamaToolCall represents a tool call result.
type LlamaToolCall struct {
	ID        string
	Name      string
	Arguments string
}

// LlamaTool represents a tool definition for llamalib.
type LlamaTool struct {
	Name        string
	Description string
	Parameters  map[string]any
}

