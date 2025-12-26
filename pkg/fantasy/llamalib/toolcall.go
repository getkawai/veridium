package llamalib

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/google/uuid"
)

// ToolCallResult represents a parsed tool call from LLM output.
type ToolCallResult struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON string
}

// ParseResult represents the parsed output from LLM, separating text and tool calls.
type ParseResult struct {
	Text      string           // Text content (before/between tool calls)
	ToolCalls []ToolCallResult // Extracted tool calls
	HasTools  bool             // Whether tool calls were found
}

var (
	// toolCallRegex matches <tool_call>...</tool_call> blocks
	toolCallRegex = regexp.MustCompile(`(?s)<tool_call>\s*(\{.*?\})\s*</tool_call>`)
)

// ParseToolCalls extracts tool calls from LLM output.
// It handles the <tool_call>{"name": "...", "arguments": {...}}</tool_call> format
// used by ChatML and Llama tool templates.
func ParseToolCalls(output string) ParseResult {
	result := ParseResult{}

	// Find all tool call matches
	matches := toolCallRegex.FindAllStringSubmatch(output, -1)
	if len(matches) == 0 {
		// No tool calls found, return entire output as text
		result.Text = strings.TrimSpace(output)
		return result
	}

	result.HasTools = true

	// Extract text before the first tool call
	firstMatch := toolCallRegex.FindStringIndex(output)
	if firstMatch != nil && firstMatch[0] > 0 {
		result.Text = strings.TrimSpace(output[:firstMatch[0]])
	}

	// Parse each tool call
	for _, match := range matches {
		if len(match) < 2 {
			continue
		}

		jsonStr := match[1]
		tc, err := parseToolCallJSON(jsonStr)
		if err != nil {
			continue // Skip invalid tool calls
		}

		// Generate ID if not provided
		if tc.ID == "" {
			tc.ID = "call_" + uuid.New().String()[:8]
		}

		result.ToolCalls = append(result.ToolCalls, tc)
	}

	return result
}

// parseToolCallJSON parses a tool call JSON object.
// Expected format: {"name": "function_name", "arguments": {...}}
func parseToolCallJSON(jsonStr string) (ToolCallResult, error) {
	var raw struct {
		Name      string          `json:"name"`
		Arguments json.RawMessage `json:"arguments"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &raw); err != nil {
		return ToolCallResult{}, err
	}

	tc := ToolCallResult{
		Name: raw.Name,
	}

	// Arguments can be object or already a string
	if len(raw.Arguments) > 0 {
		// If it's an object, marshal it back to string
		var obj map[string]interface{}
		if err := json.Unmarshal(raw.Arguments, &obj); err == nil {
			// It's an object, convert to string
			argBytes, _ := json.Marshal(obj)
			tc.Arguments = string(argBytes)
		} else {
			// It might be a string, try unquoting
			var str string
			if err := json.Unmarshal(raw.Arguments, &str); err == nil {
				tc.Arguments = str
			} else {
				tc.Arguments = string(raw.Arguments)
			}
		}
	} else {
		tc.Arguments = "{}"
	}

	return tc, nil
}

// StripToolCalls removes tool call tags from the output, leaving only text content.
func StripToolCalls(output string) string {
	// Remove all tool call blocks
	result := toolCallRegex.ReplaceAllString(output, "")
	// Clean up extra whitespace
	result = strings.TrimSpace(result)
	return result
}

