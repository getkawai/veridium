package tools

import (
	"context"
	"encoding/json"

	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/types"
)

// ToolRegistryAdapter wraps types.Tool to implement fantasy.AgentTool
type ToolRegistryAdapter struct {
	tool            *types.Tool
	providerOptions fantasy.ProviderOptions
}

// NewToolRegistryAdapter creates an adapter from types.Tool to fantasy.AgentTool
func NewToolRegistryAdapter(tool *types.Tool) *ToolRegistryAdapter {
	return &ToolRegistryAdapter{
		tool: tool,
	}
}

// Info returns tool metadata
func (a *ToolRegistryAdapter) Info() fantasy.ToolInfo {
	required := make([]string, 0)
	if a.tool.Definition.Parameters != nil {
		if req, ok := a.tool.Definition.Parameters["required"].([]interface{}); ok {
			for _, r := range req {
				if s, ok := r.(string); ok {
					required = append(required, s)
				}
			}
		}
		if req, ok := a.tool.Definition.Parameters["required"].([]string); ok {
			required = req
		}
	}

	// Extract properties from parameters
	params := make(map[string]any)
	if a.tool.Definition.Parameters != nil {
		if props, ok := a.tool.Definition.Parameters["properties"].(map[string]interface{}); ok {
			params = props
		}
	}

	return fantasy.ToolInfo{
		Name:        a.tool.Definition.Name,
		Description: a.tool.Definition.Description,
		Parameters:  params,
		Required:    required,
		Parallel:    a.tool.Parallel, // Use tool's parallel setting
	}
}

// Run executes the tool
func (a *ToolRegistryAdapter) Run(ctx context.Context, call fantasy.ToolCall) (fantasy.ToolResponse, error) {
	// Parse input JSON to map[string]string
	args := make(map[string]string)
	if call.Input != "" {
		var argsAny map[string]interface{}
		if err := json.Unmarshal([]byte(call.Input), &argsAny); err == nil {
			for k, v := range argsAny {
				switch val := v.(type) {
				case string:
					args[k] = val
				case float64:
					args[k] = jsonNumber(val)
				case bool:
					if val {
						args[k] = "true"
					} else {
						args[k] = "false"
					}
				default:
					// For complex types, marshal back to JSON
					if jsonBytes, err := json.Marshal(v); err == nil {
						args[k] = string(jsonBytes)
					}
				}
			}
		}
	}

	// Execute the tool
	result, err := a.tool.Executor(ctx, args)
	if err != nil {
		return fantasy.NewTextErrorResponse(err.Error()), nil
	}

	return fantasy.NewTextResponse(result), nil
}

// ProviderOptions returns provider-specific options
func (a *ToolRegistryAdapter) ProviderOptions() fantasy.ProviderOptions {
	return a.providerOptions
}

// SetProviderOptions sets provider-specific options
func (a *ToolRegistryAdapter) SetProviderOptions(opts fantasy.ProviderOptions) {
	a.providerOptions = opts
}

// jsonNumber converts float64 to string without scientific notation for integers
func jsonNumber(f float64) string {
	if f == float64(int64(f)) {
		return json.Number(string(rune(int64(f)))).String()
	}
	b, _ := json.Marshal(f)
	return string(b)
}

// ToAgentTools converts a slice of types.Tool to fantasy.AgentTool
func ToAgentTools(tools []*types.Tool) []fantasy.AgentTool {
	agentTools := make([]fantasy.AgentTool, 0, len(tools))
	for _, t := range tools {
		if t != nil && t.Enabled {
			agentTools = append(agentTools, NewToolRegistryAdapter(t))
		}
	}
	return agentTools
}

// RegistryToAgentTools converts all enabled tools from ToolRegistry to fantasy.AgentTool
func (r *ToolRegistry) ToAgentTools(names []string) []fantasy.AgentTool {
	tools := r.GetByNames(names)
	return ToAgentTools(tools)
}
