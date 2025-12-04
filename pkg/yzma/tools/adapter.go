package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kawai-network/veridium/types"
)

// ToolRegistry manages yzma tools
type ToolRegistry struct {
	tools map[string]*types.Tool
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]*types.Tool),
	}
}

// Register registers a tool
func (r *ToolRegistry) Register(tool *types.Tool) error {
	if tool.Function.Name == "" {
		return fmt.Errorf("tool name is required")
	}
	if tool.Executor == nil {
		return fmt.Errorf("tool executor is required")
	}

	r.tools[tool.Function.Name] = tool
	return nil
}

// Get retrieves a tool by name
func (r *ToolRegistry) Get(name string) (*types.Tool, bool) {
	tool, ok := r.tools[name]
	return tool, ok
}

// GetAll returns all tools
func (r *ToolRegistry) GetAll() []*types.Tool {
	tools := make([]*types.Tool, 0, len(r.tools))
	for _, t := range r.tools {
		tools = append(tools, t)
	}
	return tools
}

// GetEnabled returns all enabled tools
func (r *ToolRegistry) GetEnabled() []*types.Tool {
	tools := make([]*types.Tool, 0, len(r.tools))
	for _, t := range r.tools {
		if t.Enabled {
			tools = append(tools, t)
		}
	}
	return tools
}

// GetByNames returns tools by names (or all enabled if empty)
func (r *ToolRegistry) GetByNames(names []string) []*types.Tool {
	if len(names) == 0 {
		return r.GetEnabled()
	}

	tools := make([]*types.Tool, 0, len(names))
	for _, name := range names {
		if tool, ok := r.tools[name]; ok && tool.Enabled {
			tools = append(tools, tool)
		}
	}
	return tools
}

// Execute executes a tool
func (r *ToolRegistry) Execute(ctx context.Context, name string, args map[string]string) (string, error) {
	tool, ok := r.tools[name]
	if !ok {
		return "", fmt.Errorf("tool not found: %s", name)
	}

	if !tool.Enabled {
		return "", fmt.Errorf("tool is disabled: %s", name)
	}

	return tool.Executor(ctx, args)
}

// Enable enables a tool
func (r *ToolRegistry) Enable(name string) bool {
	if tool, ok := r.tools[name]; ok {
		tool.Enabled = true
		return true
	}
	return false
}

// Disable disables a tool
func (r *ToolRegistry) Disable(name string) bool {
	if tool, ok := r.tools[name]; ok {
		tool.Enabled = false
		return true
	}
	return false
}

// FormatForPrompt formats tools as JSON for system prompt
func (r *ToolRegistry) FormatForPrompt(toolNames []string) (string, error) {
	tools := r.GetByNames(toolNames)
	if len(tools) == 0 {
		return "", nil
	}

	simplified := make([]map[string]interface{}, len(tools))
	for i, t := range tools {
		simplified[i] = map[string]interface{}{
			"type": t.Type,
			"function": map[string]interface{}{
				"name":        t.Function.Name,
				"description": t.Function.Description,
				"parameters":  t.Function.Parameters,
			},
		}
	}

	toolsJSON, err := json.MarshalIndent(simplified, "", "  ")
	if err != nil {
		return "", err
	}
	return string(toolsJSON), nil
}

// ParseToolCalls parses tool calls from LLM response
// Supports multiple formats:
// 1. <tool_call>{"name": "...", "arguments": {...}}</tool_call>
// 2. <tool_name>JSON</tool_name> or <tool_name attr="value">content</tool_name>
// 3. <tool_name {JSON}> or <tool_name {"key": "value"}> (no closing tag)
// 4. {"name": "tool_name", "parameters": {...}} (pure JSON format)
func ParseToolCalls(response string) []types.ToolCall {
	var calls []types.ToolCall

	// Try format 1: <tool_call>{"name": "...", "arguments": {...}}</tool_call>
	calls = parseToolCallFormat(response)
	if len(calls) > 0 {
		return calls
	}

	// Try format 2: <tool_name>...</tool_name> with closing tag
	calls = parseXMLToolFormat(response)
	if len(calls) > 0 {
		return calls
	}

	// Try format 3: <tool_name {JSON}> without closing tag
	calls = parseInlineJSONToolFormat(response)
	if len(calls) > 0 {
		return calls
	}

	// Try format 4: {"name": "tool_name", "parameters": {...}} pure JSON
	calls = parsePureJSONToolFormat(response)
	return calls
}

// parsePureJSONToolFormat parses pure JSON tool call format
// Example: {"name": "web_search", "parameters": {"query": "AI news"}}
func parsePureJSONToolFormat(response string) []types.ToolCall {
	var calls []types.ToolCall

	// Known tool names
	toolNames := []string{"calculator", "web_search", "web-search", "search"}

	// Try to find and parse JSON objects
	remaining := strings.TrimSpace(response)

	// Check if the entire response is a JSON object
	if strings.HasPrefix(remaining, "{") {
		// Find matching closing brace
		braceCount := 0
		jsonEnd := -1
		for i, ch := range remaining {
			if ch == '{' {
				braceCount++
			} else if ch == '}' {
				braceCount--
				if braceCount == 0 {
					jsonEnd = i + 1
					break
				}
			}
		}

		if jsonEnd > 0 {
			jsonContent := remaining[:jsonEnd]

			// Try parsing as tool call with "name" and "parameters"/"arguments"
			var toolCall struct {
				Name       string                 `json:"name"`
				Parameters map[string]interface{} `json:"parameters"`
				Arguments  map[string]interface{} `json:"arguments"`
			}

			if err := json.Unmarshal([]byte(jsonContent), &toolCall); err == nil {
				// Check if it's a known tool
				isKnownTool := false
				for _, tn := range toolNames {
					if toolCall.Name == tn {
						isKnownTool = true
						break
					}
				}

				if isKnownTool && toolCall.Name != "" {
					// Use parameters or arguments
					params := toolCall.Parameters
					if params == nil {
						params = toolCall.Arguments
					}

					// Convert to map[string]string
					args := make(map[string]string)
					for k, v := range params {
						switch val := v.(type) {
						case string:
							args[k] = val
						case float64:
							args[k] = fmt.Sprintf("%v", val)
						case int:
							args[k] = fmt.Sprintf("%d", val)
						default:
							args[k] = fmt.Sprintf("%v", val)
						}
					}

					if len(args) > 0 {
						calls = append(calls, types.ToolCall{
							Type: "function",
							Function: types.ToolFunction{
								Name:      toolCall.Name,
								Arguments: args,
							},
						})
					}
				}
			}
		}
	}

	return calls
}

// parseInlineJSONToolFormat parses <tool_name {JSON}> format (no closing tag)
// Example: <web_search { "query": "AI news", "max_results": 10 }>
func parseInlineJSONToolFormat(response string) []types.ToolCall {
	var calls []types.ToolCall

	toolNames := []string{"calculator", "web_search", "web-search", "search"}

	for _, toolName := range toolNames {
		remaining := response
		for {
			// Find opening tag
			openTagStart := strings.Index(remaining, "<"+toolName)
			if openTagStart == -1 {
				break
			}

			// Find the JSON object after tool name
			afterToolName := remaining[openTagStart+len("<"+toolName):]
			afterToolName = strings.TrimSpace(afterToolName)

			// Check if it starts with { (JSON object)
			if !strings.HasPrefix(afterToolName, "{") {
				remaining = remaining[openTagStart+1:]
				continue
			}

			// Find the end of JSON object - match braces
			jsonStart := 0
			braceCount := 0
			jsonEnd := -1

			for i, ch := range afterToolName {
				if ch == '{' {
					braceCount++
				} else if ch == '}' {
					braceCount--
					if braceCount == 0 {
						jsonEnd = i + 1
						break
					}
				}
			}

			if jsonEnd == -1 {
				remaining = remaining[openTagStart+1:]
				continue
			}

			jsonContent := afterToolName[jsonStart:jsonEnd]

			// Parse JSON - try different argument formats
			var args map[string]string

			// Try parsing as map[string]interface{} first (handles numbers)
			var rawArgs map[string]interface{}
			if err := json.Unmarshal([]byte(jsonContent), &rawArgs); err == nil {
				args = make(map[string]string)
				for k, v := range rawArgs {
					switch val := v.(type) {
					case string:
						args[k] = val
					case float64:
						args[k] = fmt.Sprintf("%v", val)
					case int:
						args[k] = fmt.Sprintf("%d", val)
					default:
						args[k] = fmt.Sprintf("%v", val)
					}
				}
			}

			if len(args) > 0 {
				calls = append(calls, types.ToolCall{
					Type: "function",
					Function: types.ToolFunction{
						Name:      toolName,
						Arguments: args,
					},
				})
			}

			remaining = remaining[openTagStart+len("<"+toolName)+jsonEnd:]
		}
	}

	return calls
}

// parseToolCallFormat parses <tool_call> tags
func parseToolCallFormat(response string) []types.ToolCall {
	var calls []types.ToolCall

	start := strings.Index(response, "<tool_call>")
	end := strings.Index(response, "</tool_call>")

	for start != -1 && end != -1 && start < end {
		content := response[start+len("<tool_call>") : end]
		content = strings.TrimSpace(content)

		var parsed struct {
			Name      string            `json:"name"`
			Arguments map[string]string `json:"arguments"`
		}

		if err := json.Unmarshal([]byte(content), &parsed); err == nil {
			calls = append(calls, types.ToolCall{
				Type: "function",
				Function: types.ToolFunction{
					Name:      parsed.Name,
					Arguments: parsed.Arguments,
				},
			})
		}

		response = response[end+len("</tool_call>"):]
		start = strings.Index(response, "<tool_call>")
		end = strings.Index(response, "</tool_call>")
	}

	return calls
}

// parseXMLToolFormat parses <tool_name>...</tool_name> tags
// Supports both <tool_name>JSON</tool_name> and <tool_name attr="value">content</tool_name> formats
func parseXMLToolFormat(response string) []types.ToolCall {
	var calls []types.ToolCall

	// Common tool names to look for
	toolNames := []string{"calculator", "web_search", "web-search", "search"}

	for _, toolName := range toolNames {
		closeTag := "</" + toolName + ">"

		// Find all occurrences of this tool
		remaining := response
		for {
			// Find opening tag - could be <tool_name> or <tool_name attr="value">
			openTagStart := strings.Index(remaining, "<"+toolName)
			if openTagStart == -1 {
				break
			}

			// Find the end of opening tag (the closing >)
			openTagEnd := strings.Index(remaining[openTagStart:], ">")
			if openTagEnd == -1 {
				break
			}
			openTagEnd += openTagStart

			// Extract the full opening tag to parse attributes
			fullOpenTag := remaining[openTagStart : openTagEnd+1]

			// Find closing tag
			end := strings.Index(remaining[openTagEnd:], closeTag)
			if end == -1 {
				break
			}
			end += openTagEnd

			// Content between tags
			content := remaining[openTagEnd+1 : end]
			content = strings.TrimSpace(content)

			// Try to parse as JSON first (format: <tool_name>{"name": "...", "arguments": {...}}</tool_name>)
			var parsed struct {
				Name      string            `json:"name"`
				Arguments map[string]string `json:"arguments"`
			}

			if err := json.Unmarshal([]byte(content), &parsed); err == nil {
				calls = append(calls, types.ToolCall{
					Type: "function",
					Function: types.ToolFunction{
						Name:      parsed.Name,
						Arguments: parsed.Arguments,
					},
				})
			} else {
				// Try to parse attributes from opening tag (format: <tool_name query="..." max_results="5">)
				args := parseTagAttributes(fullOpenTag, toolName)

				// If content looks like JSON, try to parse it and merge with attributes
				if strings.HasPrefix(content, "{") {
					var jsonArgs map[string]string
					if err := json.Unmarshal([]byte(content), &jsonArgs); err == nil {
						for k, v := range jsonArgs {
							if _, exists := args[k]; !exists {
								args[k] = v
							}
						}
					}
				} else if content != "" && len(args) == 0 {
					// Content is plain text, use as query/expression
					if toolName == "calculator" {
						args["expression"] = content
					} else {
						args["query"] = content
					}
				}

				if len(args) > 0 {
					calls = append(calls, types.ToolCall{
						Type: "function",
						Function: types.ToolFunction{
							Name:      toolName,
							Arguments: args,
						},
					})
				}
			}

			remaining = remaining[end+len(closeTag):]
		}
	}

	return calls
}

// parseTagAttributes extracts attributes from XML-like opening tag
// e.g., <web_search query="AI news" max_results="5"> -> {"query": "AI news", "max_results": "5"}
func parseTagAttributes(tag string, toolName string) map[string]string {
	args := make(map[string]string)

	// Remove < and > and tool name
	inner := strings.TrimPrefix(tag, "<"+toolName)
	inner = strings.TrimSuffix(inner, ">")
	inner = strings.TrimSpace(inner)

	if inner == "" {
		return args
	}

	// Parse attributes using regex-like approach
	// Supports: attr="value" or attr='value'
	for len(inner) > 0 {
		inner = strings.TrimSpace(inner)
		if inner == "" {
			break
		}

		// Find attribute name (until = or space)
		eqIdx := strings.Index(inner, "=")
		if eqIdx == -1 {
			break
		}

		attrName := strings.TrimSpace(inner[:eqIdx])
		inner = inner[eqIdx+1:]
		inner = strings.TrimSpace(inner)

		if len(inner) == 0 {
			break
		}

		// Find attribute value (quoted)
		quote := inner[0]
		if quote != '"' && quote != '\'' {
			// Try to find space-separated value
			spaceIdx := strings.Index(inner, " ")
			if spaceIdx == -1 {
				args[attrName] = inner
				break
			}
			args[attrName] = inner[:spaceIdx]
			inner = inner[spaceIdx:]
			continue
		}

		// Find closing quote
		inner = inner[1:] // skip opening quote
		closeIdx := strings.Index(inner, string(quote))
		if closeIdx == -1 {
			break
		}

		args[attrName] = inner[:closeIdx]
		inner = inner[closeIdx+1:]
	}

	return args
}

// BuildSystemPrompt builds a system prompt with tool definitions
func BuildSystemPrompt(basePrompt string, toolsJSON string) string {
	if toolsJSON == "" {
		return basePrompt
	}

	return fmt.Sprintf(`%s

# Available Tools

You have access to the following tools:

%s

# Tool Usage Instructions

IMPORTANT: When you need to use a tool, you MUST use EXACTLY this format:

<tool_call>
{"name": "TOOL_NAME", "arguments": {"param1": "value1", "param2": "value2"}}
</tool_call>

Example for web_search:
<tool_call>
{"name": "web_search", "arguments": {"query": "latest AI news", "max_results": "5"}}
</tool_call>

Example for calculator:
<tool_call>
{"name": "calculator", "arguments": {"expression": "2 + 2 * 3"}}
</tool_call>

Rules:
1. ALWAYS use <tool_call> tags - other formats will NOT work
2. Use "arguments" (not "parameters") for the parameter object
3. All argument values must be strings (e.g., "5" not 5)
4. Wait for tool results before providing your final answer
5. After receiving results, synthesize them into a helpful response`, basePrompt, toolsJSON)
}
