/*
 * Copyright 2025 Veridium Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/kawai-network/veridium/types"
)

// Client is an HTTP client for OpenAI-compatible APIs (OpenRouter, Zhipu GLM)
type Client struct {
	httpClient   *http.Client
	baseURL      string
	apiKey       string
	providerType types.ProviderType
	options      map[string]any
}

// NewClient creates a new OpenAI-compatible API client
func NewClient(config types.ProviderConfig) *Client {
	// Determine base URL
	baseURL := config.BaseURL
	if baseURL == "" {
		if endpoint, ok := types.ProviderEndpoints[config.Type]; ok {
			baseURL = endpoint
		}
	}

	// Remove trailing slash
	baseURL = strings.TrimRight(baseURL, "/")

	return &Client{
		httpClient: &http.Client{
			Timeout: 5 * time.Minute, // Long timeout for streaming
		},
		baseURL:      baseURL,
		apiKey:       config.APIKey,
		providerType: config.Type,
		options:      config.Options,
	}
}

// ChatCompletion sends a chat completion request
func (c *Client) ChatCompletion(ctx context.Context, req types.ChatCompletionRequest) (*types.ChatCompletionResponse, error) {
	url := c.baseURL + "/chat/completions"

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range ProviderHeaders(c.providerType, c.apiKey, c.options) {
		httpReq.Header.Set(key, value)
	}

	log.Printf("🌐 [%s] Sending chat completion request to %s (model: %s, messages: %d)",
		c.providerType, url, req.Model, len(req.Messages))

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	// Read body for debugging
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var result types.ChatCompletionResponse
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	if result.Usage != nil {
		log.Printf("✅ [%s] Chat completion successful (id: %s, tokens: %d)",
			c.providerType, result.ID, result.Usage.TotalTokens)
	} else {
		log.Printf("✅ [%s] Chat completion successful (id: %s)", c.providerType, result.ID)
	}

	// Debug: Log content if choices exist
	if len(result.Choices) > 0 {
		contentPreview := ""
		if content, ok := result.Choices[0].Message.Content.(string); ok {
			if len(content) > 100 {
				contentPreview = content[:100] + "..."
			} else if content == "" {
				contentPreview = "(empty string)"
				// Log raw body when content is empty for debugging
				log.Printf("⚠️  [%s] Empty content - raw response: %s", c.providerType, string(bodyBytes))
			} else {
				contentPreview = content
			}
		} else if result.Choices[0].Message.Content == nil {
			contentPreview = "(nil content)"
			log.Printf("⚠️  [%s] Nil content - raw response: %s", c.providerType, string(bodyBytes))
		} else {
			contentPreview = fmt.Sprintf("(non-string: %T)", result.Choices[0].Message.Content)
		}
		log.Printf("📝 [%s] Response content: %s", c.providerType, contentPreview)
	} else {
		log.Printf("⚠️  [%s] No choices in response - raw response: %s", c.providerType, string(bodyBytes))
	}

	return &result, nil
}

// ChatCompletionStream sends a streaming chat completion request
func (c *Client) ChatCompletionStream(ctx context.Context, req types.ChatCompletionRequest, callback func(chunk *types.ChatCompletionStreamResponse) error) (*types.ChatCompletionResponse, error) {
	req.Stream = true

	url := c.baseURL + "/chat/completions"

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	for key, value := range ProviderHeaders(c.providerType, c.apiKey, c.options) {
		httpReq.Header.Set(key, value)
	}
	httpReq.Header.Set("Accept", "text/event-stream")

	log.Printf("🌐 [%s] Sending streaming request to %s (model: %s)",
		c.providerType, url, req.Model)

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp)
	}

	// Process SSE stream
	return c.processSSEStream(resp.Body, callback)
}

// processSSEStream processes Server-Sent Events stream
func (c *Client) processSSEStream(body io.Reader, callback func(chunk *types.ChatCompletionStreamResponse) error) (*types.ChatCompletionResponse, error) {
	scanner := bufio.NewScanner(body)

	// Increase buffer size for large responses
	buf := make([]byte, 64*1024)
	scanner.Buffer(buf, 1024*1024)

	var finalResponse types.ChatCompletionResponse
	var fullContent strings.Builder
	var toolCalls []types.APIToolCall
	var usage *types.APIUsage

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, ":") {
			continue
		}

		// Parse SSE data
		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")

		// Check for stream end
		if data == "[DONE]" {
			break
		}

		var chunk types.ChatCompletionStreamResponse
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			log.Printf("⚠️  Failed to parse stream chunk: %v (data: %s)", err, data)
			continue
		}

		// Update final response metadata
		if finalResponse.ID == "" {
			finalResponse.ID = chunk.ID
			finalResponse.Object = chunk.Object
			finalResponse.Created = chunk.Created
			finalResponse.Model = chunk.Model
		}

		// Process choices
		for _, choice := range chunk.Choices {
			// Accumulate content (handle both string and interface{})
			if content, ok := choice.Delta.Content.(string); ok && content != "" {
				fullContent.WriteString(content)
			}

			// Accumulate tool calls
			if len(choice.Delta.ToolCalls) > 0 {
				toolCalls = mergeToolCalls(toolCalls, choice.Delta.ToolCalls)
			}

			// Capture finish reason
			if choice.FinishReason != "" {
				if len(finalResponse.Choices) == 0 {
					finalResponse.Choices = []types.ChatCompletionChoice{{Index: choice.Index}}
				}
				finalResponse.Choices[0].FinishReason = choice.FinishReason
			}
		}

		// Capture usage (some providers send in final chunk)
		if chunk.Usage != nil {
			usage = chunk.Usage
		}

		// Call callback with chunk
		if callback != nil {
			if err := callback(&chunk); err != nil {
				return nil, fmt.Errorf("callback error: %w", err)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("stream scan error: %w", err)
	}

	// Build final response
	if len(finalResponse.Choices) == 0 {
		finalResponse.Choices = []types.ChatCompletionChoice{{Index: 0}}
	}

	finalResponse.Choices[0].Message = types.ChatCompletionMsg{
		Role:      "assistant",
		Content:   fullContent.String(),
		ToolCalls: toolCalls,
	}
	finalResponse.Usage = usage

	return &finalResponse, nil
}

// mergeToolCalls merges streamed tool call deltas
func mergeToolCalls(existing []types.APIToolCall, deltas []types.APIToolCall) []types.APIToolCall {
	for _, delta := range deltas {
		// Find existing tool call by index
		found := false
		for i, tc := range existing {
			if tc.ID == delta.ID || (delta.ID == "" && i < len(deltas)) {
				// Append to existing arguments
				existing[i].Function.Arguments += delta.Function.Arguments
				if delta.Function.Name != "" {
					existing[i].Function.Name = delta.Function.Name
				}
				found = true
				break
			}
		}

		if !found {
			// New tool call
			existing = append(existing, delta)
		}
	}
	return existing
}

// parseError parses an API error response
func (c *Client) parseError(resp *http.Response) error {
	body, _ := io.ReadAll(resp.Body)

	var apiErr types.APIError
	if err := json.Unmarshal(body, &apiErr); err == nil && apiErr.Error != nil {
		return fmt.Errorf("[%s] API error (status %d): %s - %s",
			c.providerType, resp.StatusCode, apiErr.Error.Type, apiErr.Error.Message)
	}

	return fmt.Errorf("[%s] API error (status %d): %s",
		c.providerType, resp.StatusCode, string(body))
}

// GetProviderType returns the provider type
func (c *Client) GetProviderType() types.ProviderType {
	return c.providerType
}

// GetBaseURL returns the base URL
func (c *Client) GetBaseURL() string {
	return c.baseURL
}
