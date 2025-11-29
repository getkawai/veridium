package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kawai-network/veridium/pkg/yzma/message"
)

// YzmaTool represents a tool in yzma format
type YzmaTool struct {
	Type     string           `json:"type"`
	Function YzmaToolFunction `json:"function"`
	Executor ToolExecutor
	Enabled  bool
}

// YzmaToolFunction represents a function definition
type YzmaToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ToolExecutor is a function that executes a tool
type ToolExecutor func(ctx context.Context, args map[string]string) (string, error)

// ToolRegistry manages yzma tools
type ToolRegistry struct {
	tools map[string]*YzmaTool
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]*YzmaTool),
	}
}

// Register registers a tool
func (r *ToolRegistry) Register(tool *YzmaTool) error {
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
func (r *ToolRegistry) Get(name string) (*YzmaTool, bool) {
	tool, ok := r.tools[name]
	return tool, ok
}

// GetAll returns all tools
func (r *ToolRegistry) GetAll() []*YzmaTool {
	tools := make([]*YzmaTool, 0, len(r.tools))
	for _, t := range r.tools {
		tools = append(tools, t)
	}
	return tools
}

// GetEnabled returns all enabled tools
func (r *ToolRegistry) GetEnabled() []*YzmaTool {
	tools := make([]*YzmaTool, 0, len(r.tools))
	for _, t := range r.tools {
		if t.Enabled {
			tools = append(tools, t)
		}
	}
	return tools
}

// GetByNames returns tools by names (or all enabled if empty)
func (r *ToolRegistry) GetByNames(names []string) []*YzmaTool {
	if len(names) == 0 {
		return r.GetEnabled()
	}
	
	tools := make([]*YzmaTool, 0, len(names))
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
// Supports both <tool_call> and <tool_name> formats
func ParseToolCalls(response string) []message.ToolCall {
	var calls []message.ToolCall
	
	// Try format 1: <tool_call>{"name": "...", "arguments": {...}}</tool_call>
	calls = parseToolCallFormat(response)
	if len(calls) > 0 {
		return calls
	}
	
	// Try format 2: <tool_name>{"name": "...", "arguments": {...}}</tool_name>
	calls = parseXMLToolFormat(response)
	return calls
}

// parseToolCallFormat parses <tool_call> tags
func parseToolCallFormat(response string) []message.ToolCall {
	var calls []message.ToolCall
	
	start := strings.Index(response, "<tool_call>")
	end := strings.Index(response, "</tool_call>")
	
	for start != -1 && end != -1 && start < end {
		content := response[start+len("<tool_call>"):end]
		content = strings.TrimSpace(content)
		
		var parsed struct {
			Name      string            `json:"name"`
			Arguments map[string]string `json:"arguments"`
		}
		
		if err := json.Unmarshal([]byte(content), &parsed); err == nil {
			calls = append(calls, message.ToolCall{
				Type: "function",
				Function: message.ToolFunction{
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
func parseXMLToolFormat(response string) []message.ToolCall {
	var calls []message.ToolCall
	
	// Common tool names to look for
	toolNames := []string{"calculator", "web_search", "web-search", "search"}
	
	for _, toolName := range toolNames {
		openTag := "<" + toolName + ">"
		closeTag := "</" + toolName + ">"
		
		start := strings.Index(response, openTag)
		end := strings.Index(response, closeTag)
		
		for start != -1 && end != -1 && start < end {
			content := response[start+len(openTag):end]
			content = strings.TrimSpace(content)
			
			var parsed struct {
				Name      string            `json:"name"`
				Arguments map[string]string `json:"arguments"`
			}
			
			if err := json.Unmarshal([]byte(content), &parsed); err == nil {
				calls = append(calls, message.ToolCall{
					Type: "function",
					Function: message.ToolFunction{
						Name:      parsed.Name,
						Arguments: parsed.Arguments,
					},
				})
			}
			
			response = response[end+len(closeTag):]
			start = strings.Index(response, openTag)
			end = strings.Index(response, closeTag)
		}
	}
	
	return calls
}

// BuildSystemPrompt builds a system prompt with tool definitions
func BuildSystemPrompt(basePrompt string, toolsJSON string) string {
	if toolsJSON == "" {
		return basePrompt
	}
	
	return fmt.Sprintf(`%s

You have access to the following tools:

%s

When you need to use a tool, respond with a tool call in the following format:
<tool_call>
{"name": "function_name", "arguments": {"arg1": "value1", "arg2": "value2"}}
</tool_call>

After receiving tool results, provide a final answer to the user.`, basePrompt, toolsJSON)
}

