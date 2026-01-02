package pooled

import (
	"encoding/json"

	"github.com/kawai-network/veridium/pkg/cliproxy/sdk/cliproxy/executor"
	"github.com/kawai-network/veridium/pkg/fantasy"
)

// convertCallToRequest converts fantasy.Call to executor.Request.
func convertCallToRequest(call fantasy.Call) executor.Request {
	// Build a standard OpenAI-compatible payload
	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	
	type Payload struct {
		Messages    []Message `json:"messages"`
		MaxTokens   *int64    `json:"max_tokens,omitempty"`
		Temperature *float64  `json:"temperature,omitempty"`
		TopP        *float64  `json:"top_p,omitempty"`
	}

	// Convert fantasy.Prompt ([]Message) to payload messages
	messages := make([]Message, 0, len(call.Prompt))
	for _, msg := range call.Prompt {
		// Extract text from message content (MessagePart)
		text := ""
		for _, part := range msg.Content {
			// Use type switch to handle different content types
			switch c := part.(type) {
			case fantasy.TextPart:
				text += c.Text
			}
		}
		
		messages = append(messages, Message{
			Role:    string(msg.Role),
			Content: text,
		})
	}

	payload := Payload{
		Messages:    messages,
		MaxTokens:   call.MaxOutputTokens,
		Temperature: call.Temperature,
		TopP:        call.TopP,
	}

	// Serialize to JSON
	payloadJSON, _ := json.Marshal(payload)

	return executor.Request{
		Model:   "", // Model is handled by provider
		Payload: payloadJSON,
		Metadata: map[string]any{
			"fantasy_call": call, // Store original call for reference
		},
	}
}

// convertRequestToCall converts executor.Request to fantasy.Call.
func convertRequestToCall(req executor.Request) fantasy.Call {
	// Try to get original call from metadata
	if originalCall, ok := req.Metadata["fantasy_call"].(fantasy.Call); ok {
		return originalCall
	}

	// Otherwise parse from payload
	type Message struct {
		Role    string `json:"role"`
		Content string `json:"content"`
	}
	
	type Payload struct {
		Messages    []Message `json:"messages"`
		MaxTokens   *int64    `json:"max_tokens,omitempty"`
		Temperature *float64  `json:"temperature,omitempty"`
		TopP        *float64  `json:"top_p,omitempty"`
	}

	var payload Payload
	if err := json.Unmarshal(req.Payload, &payload); err != nil {
		// Return empty call if we can't parse
		return fantasy.Call{}
	}

	// Convert messages to fantasy.Prompt
	prompt := make(fantasy.Prompt, 0, len(payload.Messages))
	for _, msg := range payload.Messages {
		prompt = append(prompt, fantasy.Message{
			Role: fantasy.MessageRole(msg.Role),
			Content: []fantasy.MessagePart{
				fantasy.TextPart{Text: msg.Content},
			},
		})
	}

	return fantasy.Call{
		Prompt:          prompt,
		MaxOutputTokens: payload.MaxTokens,
		Temperature:     payload.Temperature,
		TopP:            payload.TopP,
	}
}

// convertResponseToFantasy converts executor.Response to fantasy.Response.
func convertResponseToFantasy(resp executor.Response) (*fantasy.Response, error) {
	// Extract response data from metadata
	contentStr, ok := resp.Metadata["content"].(string)
	if !ok {
		return &fantasy.Response{}, nil
	}

	// Parse the content as a simple text response
	// ResponseContent is []Content, so we create a slice with one TextContent
	content := fantasy.ResponseContent{
		fantasy.TextContent{Text: contentStr},
	}

	finishReasonStr := ""
	if fr, ok := resp.Metadata["finish_reason"].(string); ok {
		finishReasonStr = fr
	}

	return &fantasy.Response{
		Content:      content,
		FinishReason: fantasy.FinishReason(finishReasonStr),
		Usage: fantasy.Usage{
			InputTokens:  getIntFromMetadata(resp.Metadata, "prompt_tokens"),
			OutputTokens: getIntFromMetadata(resp.Metadata, "completion_tokens"),
			TotalTokens:  getIntFromMetadata(resp.Metadata, "total_tokens"),
		},
	}, nil
}

// convertFantasyToResponse converts fantasy.Response to executor.Response.
func convertFantasyToResponse(resp *fantasy.Response) executor.Response {
	// Extract text from ResponseContent
	// Handle both text content and reasoning content
	contentText := ""
	
	for _, content := range resp.Content {
		switch c := content.(type) {
		case fantasy.TextContent:
			contentText += c.Text
		case fantasy.ReasoningContent:
			// For reasoning content, include it as well
			contentText += c.Text
		}
	}
	
	// Fallback: use the Text() method if nothing extracted
	if contentText == "" {
		contentText = resp.Content.Text()
	}

	return executor.Response{
		Metadata: map[string]any{
			"content":           contentText,
			"finish_reason":     string(resp.FinishReason),
			"prompt_tokens":     resp.Usage.InputTokens,
			"completion_tokens": resp.Usage.OutputTokens,
			"total_tokens":      resp.Usage.TotalTokens,
		},
	}
}

// convertStreamToFantasy converts executor stream to fantasy.StreamResponse.
func convertStreamToFantasy(chunks <-chan executor.StreamChunk) fantasy.StreamResponse {
	return func(yield func(fantasy.StreamPart) bool) {
		for chunk := range chunks {
			part := convertChunkToFantasyPart(chunk)
			if !yield(part) {
				return
			}
		}
	}
}

// convertChunkToFantasyPart converts executor.StreamChunk to fantasy.StreamPart.
func convertChunkToFantasyPart(chunk executor.StreamChunk) fantasy.StreamPart {
	if chunk.Err != nil {
		return fantasy.StreamPart{
			Type:  fantasy.StreamPartTypeError,
			Error: chunk.Err,
		}
	}

	// Parse payload as delta text
	if len(chunk.Payload) > 0 {
		return fantasy.StreamPart{
			Type:  fantasy.StreamPartTypeTextDelta,
			Delta: string(chunk.Payload),
		}
	}

	return fantasy.StreamPart{
		Type: fantasy.StreamPartTypeTextDelta,
	}
}

// convertFantasyPartToChunk converts fantasy.StreamPart to executor.StreamChunk.
func convertFantasyPartToChunk(part fantasy.StreamPart) executor.StreamChunk {
	if part.Error != nil {
		return executor.StreamChunk{
			Err: part.Error,
		}
	}

	// Convert delta to payload
	return executor.StreamChunk{
		Payload: []byte(part.Delta),
	}
}

// getIntFromMetadata safely extracts an int from metadata.
func getIntFromMetadata(metadata map[string]any, key string) int64 {
	if val, ok := metadata[key].(int); ok {
		return int64(val)
	}
	if val, ok := metadata[key].(int64); ok {
		return val
	}
	if val, ok := metadata[key].(float64); ok {
		return int64(val)
	}
	return 0
}

