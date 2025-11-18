package toolsengine

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/schema"
)

// ToolExecutor is a function that executes a tool (simple, map-based)
type ToolExecutor func(ctx context.Context, args map[string]interface{}) (interface{}, error)

// Tool wraps an Eino tool with metadata
type Tool struct {
	tool.InvokableTool

	ID       string
	Category string
	Version  string
	Enabled  bool
	Metadata map[string]interface{}
}

// simpleTool implements Eino InvokableTool interface (for backward compatibility)
type simpleTool struct {
	id       string
	name     string
	desc     string
	params   map[string]*schema.ParameterInfo
	executor ToolExecutor
}

// NewTool creates a new Eino-compatible tool (simple, map-based)
// For backward compatibility with existing code
func NewTool(id, name, desc string, params map[string]*schema.ParameterInfo, executor ToolExecutor) *Tool {
	return &Tool{
		InvokableTool: &simpleTool{
			id:       id,
			name:     name,
			desc:     desc,
			params:   params,
			executor: executor,
		},
		ID:       id,
		Enabled:  true,
		Metadata: make(map[string]interface{}),
	}
}

// NewTypedTool creates a tool using Eino's utils.NewTool with Go generics
// This is the recommended way for type-safe tools
//
// Example:
//
//	type Request struct { Query string `json:"query"` }
//	type Response struct { Results []string `json:"results"` }
//	fn := func(ctx context.Context, req *Request) (*Response, error) { ... }
//	toolInfo := &schema.ToolInfo{Name: "search", Desc: "Search"}
//	tool := NewTypedTool("search", toolInfo, fn)
func NewTypedTool[T, D any](id string, toolInfo *schema.ToolInfo, fn func(context.Context, T) (D, error)) *Tool {
	einoTool := utils.NewTool(toolInfo, fn)
	return WrapEinoTool(id, einoTool)
}

// InferTool creates a tool using Eino's utils.InferTool (auto schema generation)
// This is the most convenient way - schema is automatically generated from Go types
//
// Example:
//
//	type Request struct {
//	    Query string `json:"query" jsonschema_description:"Search query"`
//	}
//	type Response struct { Results []string `json:"results"` }
//	fn := func(ctx context.Context, req *Request) (*Response, error) { ... }
//	tool := InferTool("search", "search", "Search the web", fn)
func InferTool[T, D any](id, name, desc string, fn func(context.Context, T) (D, error)) (*Tool, error) {
	einoTool, err := utils.InferTool(name, desc, fn)
	if err != nil {
		return nil, fmt.Errorf("failed to infer tool: %w", err)
	}
	return WrapEinoTool(id, einoTool), nil
}

// WrapEinoTool wraps an existing Eino tool with metadata
// This is useful for integrating Eino-ext tools or tools created with utils.InferTool/NewTool
func WrapEinoTool(id string, einoTool tool.InvokableTool) *Tool {
	return &Tool{
		InvokableTool: einoTool,
		ID:            id,
		Enabled:       true,
		Metadata:      make(map[string]interface{}),
	}
}

// Info implements tool.BaseTool
func (t *simpleTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	info := &schema.ToolInfo{
		Name: t.name,
		Desc: t.desc,
	}

	if t.params != nil && len(t.params) > 0 {
		info.ParamsOneOf = schema.NewParamsOneOfByParams(t.params)
	}

	return info, nil
}

// InvokableRun implements tool.InvokableTool
func (t *simpleTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// Parse arguments
	var args map[string]interface{}
	if argumentsInJSON != "" && argumentsInJSON != "{}" {
		if err := json.Unmarshal([]byte(argumentsInJSON), &args); err != nil {
			return "", fmt.Errorf("failed to parse arguments: %w", err)
		}
	} else {
		args = make(map[string]interface{})
	}

	// Execute tool
	result, err := t.executor(ctx, args)
	if err != nil {
		return "", fmt.Errorf("tool execution failed: %w", err)
	}

	// Convert result to JSON string
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(resultJSON), nil
}

// ToolBuilder helps build tools with fluent API
type ToolBuilder struct {
	id       string
	name     string
	desc     string
	params   map[string]*schema.ParameterInfo
	executor ToolExecutor
	category string
	version  string
	metadata map[string]interface{}
}

// NewToolBuilder creates a new tool builder
func NewToolBuilder(id, name string) *ToolBuilder {
	return &ToolBuilder{
		id:       id,
		name:     name,
		params:   make(map[string]*schema.ParameterInfo),
		metadata: make(map[string]interface{}),
	}
}

// WithDescription sets the description
func (b *ToolBuilder) WithDescription(desc string) *ToolBuilder {
	b.desc = desc
	return b
}

// WithParameter adds a parameter
func (b *ToolBuilder) WithParameter(name string, paramType schema.DataType, desc string, required bool) *ToolBuilder {
	b.params[name] = &schema.ParameterInfo{
		Type:     paramType,
		Desc:     desc,
		Required: required,
	}
	return b
}

// WithExecutor sets the executor
func (b *ToolBuilder) WithExecutor(executor ToolExecutor) *ToolBuilder {
	b.executor = executor
	return b
}

// WithCategory sets the category
func (b *ToolBuilder) WithCategory(category string) *ToolBuilder {
	b.category = category
	return b
}

// WithVersion sets the version
func (b *ToolBuilder) WithVersion(version string) *ToolBuilder {
	b.version = version
	return b
}

// Build creates the tool
func (b *ToolBuilder) Build() (*Tool, error) {
	if b.id == "" {
		return nil, fmt.Errorf("tool ID is required")
	}
	if b.name == "" {
		return nil, fmt.Errorf("tool name is required")
	}
	if b.executor == nil {
		return nil, fmt.Errorf("tool executor is required")
	}

	tool := NewTool(b.id, b.name, b.desc, b.params, b.executor)
	tool.Category = b.category
	tool.Version = b.version
	tool.Metadata = b.metadata

	return tool, nil
}

// ConvertToOpenAI converts Eino tools to OpenAI ChatCompletionTool format
func ConvertToOpenAI(ctx context.Context, einoTools []tool.InvokableTool) ([]ChatCompletionTool, error) {
	openAITools := make([]ChatCompletionTool, 0, len(einoTools))

	for _, einoTool := range einoTools {
		info, err := einoTool.Info(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get tool info: %w", err)
		}

		// Convert to JSON Schema
		var parameters map[string]interface{}
		if info.ParamsOneOf != nil {
			jsonSchema, err := info.ParamsOneOf.ToJSONSchema()
			if err != nil {
				return nil, fmt.Errorf("failed to convert params to JSON schema: %w", err)
			}

			if jsonSchema != nil {
				schemaJSON, err := json.Marshal(jsonSchema)
				if err != nil {
					return nil, fmt.Errorf("failed to marshal JSON schema: %w", err)
				}

				if err := json.Unmarshal(schemaJSON, &parameters); err != nil {
					return nil, fmt.Errorf("failed to unmarshal JSON schema: %w", err)
				}
			}
		}

		if parameters == nil {
			parameters = map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			}
		}

		openAITool := ChatCompletionTool{
			Type: "function",
			Function: FunctionDefinition{
				Name:        info.Name,
				Description: info.Desc,
				Parameters:  parameters,
			},
		}

		openAITools = append(openAITools, openAITool)
	}

	return openAITools, nil
}

// ChatCompletionTool represents an OpenAI-compatible tool definition
type ChatCompletionTool struct {
	Type     string             `json:"type"`
	Function FunctionDefinition `json:"function"`
}

// FunctionDefinition represents a function definition for OpenAI
type FunctionDefinition struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}
