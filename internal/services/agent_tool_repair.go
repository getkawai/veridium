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

package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/kawai-network/veridium/pkg/fantasy"
	"github.com/kawai-network/veridium/pkg/fantasy/schema"
)

// RepairToolCall attempts to repair a malformed tool call.
// It handles common LLM mistakes like:
// - Malformed JSON (missing quotes, trailing commas, etc.)
// - Missing required fields (attempts to infer from context)
// - Wrong tool name (fuzzy matching)
//
// Returns the repaired tool call or an error if repair is not possible.
func RepairToolCall(ctx context.Context, opts fantasy.ToolCallRepairOptions) (*fantasy.ToolCallContent, error) {
	log.Printf("🔧 [REPAIR] Attempting to repair tool call: %s (error: %v)", opts.OriginalToolCall.ToolName, opts.ValidationError)

	repaired := opts.OriginalToolCall
	errorStr := opts.ValidationError.Error()

	// Strategy 1: Fix malformed JSON
	if strings.Contains(errorStr, "invalid JSON") || strings.Contains(errorStr, "unexpected") {
		fixedInput, err := repairJSON(opts.OriginalToolCall.Input)
		if err != nil {
			log.Printf("🔧 [REPAIR] JSON repair failed: %v", err)
		} else {
			log.Printf("🔧 [REPAIR] JSON repaired successfully")
			repaired.Input = fixedInput

			// Validate the repaired JSON
			var parsed map[string]any
			if err := json.Unmarshal([]byte(fixedInput), &parsed); err == nil {
				return &repaired, nil
			}
		}
	}

	// Strategy 2: Fix tool not found (fuzzy match)
	if strings.Contains(errorStr, "tool not found") {
		matchedTool := fuzzyMatchTool(opts.OriginalToolCall.ToolName, opts.AvailableTools)
		if matchedTool != "" {
			log.Printf("🔧 [REPAIR] Tool name corrected: %s -> %s", opts.OriginalToolCall.ToolName, matchedTool)
			repaired.ToolName = matchedTool
			return &repaired, nil
		}
	}

	// Strategy 3: Fix missing required fields
	if strings.Contains(errorStr, "missing required parameter") {
		fixedInput, err := fillMissingFields(opts.OriginalToolCall, opts.AvailableTools, opts.Messages)
		if err != nil {
			log.Printf("🔧 [REPAIR] Could not fill missing fields: %v", err)
		} else {
			log.Printf("🔧 [REPAIR] Missing fields filled")
			repaired.Input = fixedInput
			return &repaired, nil
		}
	}

	// Strategy 4: Combined repair - fix JSON first, then check fields
	fixedInput, jsonErr := repairJSON(opts.OriginalToolCall.Input)
	if jsonErr == nil {
		repaired.Input = fixedInput

		// Try to fill missing fields on the repaired JSON
		filledInput, fillErr := fillMissingFields(repaired, opts.AvailableTools, opts.Messages)
		if fillErr == nil {
			repaired.Input = filledInput
			log.Printf("🔧 [REPAIR] Combined repair successful")
			return &repaired, nil
		}

		// Even if filling fails, return the JSON-repaired version
		return &repaired, nil
	}

	log.Printf("🔧 [REPAIR] All repair strategies failed")
	return nil, fmt.Errorf("unable to repair tool call: %w", opts.ValidationError)
}

// repairJSON attempts to fix common JSON syntax errors using fantasy/schema package
func repairJSON(input string) (string, error) {
	if input == "" {
		return "{}", nil
	}

	// Use schema.ParsePartialJSON which handles repair internally
	obj, state, err := schema.ParsePartialJSON(input)
	if state == schema.ParseStateFailed {
		return "", fmt.Errorf("json repair failed: %w", err)
	}

	// Marshal back to clean JSON string
	result, err := json.Marshal(obj)
	if err != nil {
		return "", fmt.Errorf("failed to marshal repaired JSON: %w", err)
	}

	return string(result), nil
}

// fuzzyMatchTool tries to find the closest matching tool name
func fuzzyMatchTool(name string, availableTools []fantasy.AgentTool) string {
	nameLower := strings.ToLower(name)

	// Exact match (case-insensitive)
	for _, tool := range availableTools {
		if strings.ToLower(tool.Info().Name) == nameLower {
			return tool.Info().Name
		}
	}

	// Partial match - check if tool name contains the input or vice versa
	for _, tool := range availableTools {
		toolName := tool.Info().Name
		toolNameLower := strings.ToLower(toolName)

		// Check contains
		if strings.Contains(toolNameLower, nameLower) || strings.Contains(nameLower, toolNameLower) {
			return toolName
		}

		// Check with underscores/dashes removed
		normalizedInput := strings.ReplaceAll(strings.ReplaceAll(nameLower, "_", ""), "-", "")
		normalizedTool := strings.ReplaceAll(strings.ReplaceAll(toolNameLower, "_", ""), "-", "")
		if normalizedInput == normalizedTool {
			return toolName
		}

		// Check suffix match (e.g., "search" matches "lobe-web-browsing__search")
		if strings.HasSuffix(toolNameLower, "__"+nameLower) {
			return toolName
		}
	}

	// Levenshtein distance for close matches (simple implementation)
	bestMatch := ""
	bestDistance := 999

	for _, tool := range availableTools {
		toolName := tool.Info().Name
		distance := levenshteinDistance(nameLower, strings.ToLower(toolName))

		// Only accept if distance is small relative to string length
		maxAcceptable := len(nameLower) / 3
		if maxAcceptable < 2 {
			maxAcceptable = 2
		}

		if distance < bestDistance && distance <= maxAcceptable {
			bestDistance = distance
			bestMatch = toolName
		}
	}

	return bestMatch
}

// fillMissingFields attempts to fill missing required fields with reasonable defaults
func fillMissingFields(toolCall fantasy.ToolCallContent, availableTools []fantasy.AgentTool, messages []fantasy.Message) (string, error) {
	// Find the tool
	var tool fantasy.AgentTool
	for _, t := range availableTools {
		if t.Info().Name == toolCall.ToolName {
			tool = t
			break
		}
	}

	if tool == nil {
		return "", fmt.Errorf("tool not found: %s", toolCall.ToolName)
	}

	// Parse current input
	var input map[string]any
	if err := json.Unmarshal([]byte(toolCall.Input), &input); err != nil {
		// Try to repair JSON first
		repaired, repairErr := repairJSON(toolCall.Input)
		if repairErr != nil {
			return "", fmt.Errorf("cannot parse input: %w", err)
		}
		if err := json.Unmarshal([]byte(repaired), &input); err != nil {
			return "", fmt.Errorf("cannot parse repaired input: %w", err)
		}
	}

	toolInfo := tool.Info()

	// Check each required field
	for _, required := range toolInfo.Required {
		if _, exists := input[required]; !exists {
			// Try to infer value from context
			inferredValue := inferFieldValue(required, toolInfo.Parameters, messages)
			if inferredValue != nil {
				input[required] = inferredValue
				log.Printf("🔧 [REPAIR] Inferred value for %s: %v", required, inferredValue)
			} else {
				// Use default based on type
				if paramSchema, ok := toolInfo.Parameters[required].(map[string]any); ok {
					defaultValue := getDefaultForType(paramSchema)
					if defaultValue != nil {
						input[required] = defaultValue
						log.Printf("🔧 [REPAIR] Using default for %s: %v", required, defaultValue)
					}
				}
			}
		}
	}

	// Marshal back to JSON
	result, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("cannot marshal repaired input: %w", err)
	}

	return string(result), nil
}

// inferFieldValue tries to infer a field value from conversation context
func inferFieldValue(fieldName string, params map[string]any, messages []fantasy.Message) any {
	// Common field inference based on field name
	fieldNameLower := strings.ToLower(fieldName)

	// Try to find relevant content in recent messages
	for i := len(messages) - 1; i >= 0 && i >= len(messages)-5; i-- {
		msg := messages[i]
		if msg.Role != fantasy.MessageRoleUser {
			continue
		}

		// Get text content
		for _, part := range msg.Content {
			if textPart, ok := part.(fantasy.TextPart); ok {
				text := textPart.Text

				// For query/search fields, use user message
				if fieldNameLower == "query" || fieldNameLower == "search" || fieldNameLower == "q" {
					// Use the user's message as the query
					if len(text) > 0 && len(text) < 500 {
						return text
					}
				}

				// For URL fields, try to extract URL
				if fieldNameLower == "url" || fieldNameLower == "link" {
					if url := extractURL(text); url != "" {
						return url
					}
				}

				// For path/file fields
				if fieldNameLower == "path" || fieldNameLower == "file" || fieldNameLower == "filepath" {
					if path := extractPath(text); path != "" {
						return path
					}
				}
			}
		}
	}

	return nil
}

// getDefaultForType returns a reasonable default value based on parameter type
func getDefaultForType(paramSchema map[string]any) any {
	paramType, _ := paramSchema["type"].(string)

	switch paramType {
	case "string":
		// Check for enum
		if enum, ok := paramSchema["enum"].([]any); ok && len(enum) > 0 {
			return enum[0]
		}
		// Check for default
		if def, ok := paramSchema["default"]; ok {
			return def
		}
		return ""
	case "integer", "number":
		if def, ok := paramSchema["default"]; ok {
			return def
		}
		if min, ok := paramSchema["minimum"]; ok {
			return min
		}
		return 0
	case "boolean":
		if def, ok := paramSchema["default"]; ok {
			return def
		}
		return false
	case "array":
		return []any{}
	case "object":
		return map[string]any{}
	default:
		return nil
	}
}

// extractURL extracts a URL from text
func extractURL(text string) string {
	// Simple URL extraction
	prefixes := []string{"https://", "http://"}
	for _, prefix := range prefixes {
		if idx := strings.Index(text, prefix); idx != -1 {
			end := idx
			for end < len(text) && !isURLTerminator(text[end]) {
				end++
			}
			if end > idx+len(prefix) {
				return text[idx:end]
			}
		}
	}
	return ""
}

// extractPath extracts a file path from text
func extractPath(text string) string {
	// Look for common path patterns
	// Unix-style paths
	if idx := strings.Index(text, "/"); idx != -1 {
		end := idx
		for end < len(text) && !isPathTerminator(text[end]) {
			end++
		}
		if end > idx+1 {
			return text[idx:end]
		}
	}

	// Windows-style paths
	for i := 0; i < len(text)-2; i++ {
		if text[i] >= 'A' && text[i] <= 'Z' && text[i+1] == ':' && text[i+2] == '\\' {
			end := i
			for end < len(text) && !isPathTerminator(text[end]) {
				end++
			}
			if end > i+3 {
				return text[i:end]
			}
		}
	}

	return ""
}

func isURLTerminator(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '"' || c == '\'' || c == '>' || c == ')' || c == ']'
}

func isPathTerminator(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r' || c == '"' || c == '\''
}

// levenshteinDistance calculates the edit distance between two strings
func levenshteinDistance(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	// Create matrix
	matrix := make([][]int, len(s1)+1)
	for i := range matrix {
		matrix[i] = make([]int, len(s2)+1)
		matrix[i][0] = i
	}
	for j := range matrix[0] {
		matrix[0][j] = j
	}

	// Fill matrix
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			cost := 1
			if s1[i-1] == s2[j-1] {
				cost = 0
			}

			matrix[i][j] = min(
				matrix[i-1][j]+1,      // deletion
				matrix[i][j-1]+1,      // insertion
				matrix[i-1][j-1]+cost, // substitution
			)
		}
	}

	return matrix[len(s1)][len(s2)]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}
	if b < c {
		return b
	}
	return c
}
