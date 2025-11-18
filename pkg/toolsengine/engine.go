package toolsengine

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/components/tool"
)

// ToolsEngine handles tool management and execution
type ToolsEngine struct {
	registry *ToolRegistry
}

// Config holds configuration for the tools engine
type Config struct {
	// Reserved for future use
}

// NewToolsEngine creates a new tools engine
func NewToolsEngine(config Config) (*ToolsEngine, error) {
	registry := NewToolRegistry()

	engine := &ToolsEngine{
		registry: registry,
	}

	log.Printf("ToolsEngine initialized (Eino-native)")

	return engine, nil
}

// RegisterTool registers a tool
func (e *ToolsEngine) RegisterTool(t *Tool) error {
	return e.registry.Register(t)
}

// GetTool retrieves a tool by ID
func (e *ToolsEngine) GetTool(id string) (*Tool, bool) {
	return e.registry.Get(id)
}

// GetAllTools returns all tools
func (e *ToolsEngine) GetAllTools() []*Tool {
	return e.registry.GetAll()
}

// GetEnabledTools returns all enabled tools
func (e *ToolsEngine) GetEnabledTools() []*Tool {
	return e.registry.GetEnabled()
}

// GetToolIDs returns all tool IDs
func (e *ToolsEngine) GetToolIDs() []string {
	return e.registry.GetIDs()
}

// HasTool checks if a tool exists
func (e *ToolsEngine) HasTool(id string) bool {
	return e.registry.Has(id)
}

// RemoveTool removes a tool
func (e *ToolsEngine) RemoveTool(id string) bool {
	return e.registry.Remove(id)
}

// EnableTool enables a tool
func (e *ToolsEngine) EnableTool(id string) bool {
	return e.registry.Enable(id)
}

// DisableTool disables a tool
func (e *ToolsEngine) DisableTool(id string) bool {
	return e.registry.Disable(id)
}

// GenerateToolsParams holds parameters for tool generation
type GenerateToolsParams struct {
	ToolIDs  []string               `json:"toolIds"`
	Model    string                 `json:"model"`
	Provider string                 `json:"provider"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

// GenerateTools generates OpenAI-compatible tools
func (e *ToolsEngine) GenerateTools(params GenerateToolsParams) ([]ChatCompletionTool, error) {
	log.Printf("Generating tools for toolIds=%v", params.ToolIDs)

	// Get tools by IDs (or all if empty)
	var tools []*Tool
	if len(params.ToolIDs) == 0 {
		tools = e.registry.GetEnabled()
	} else {
		tools = e.registry.GetByIDs(params.ToolIDs)
	}

	if len(tools) == 0 {
		log.Printf("No tools found")
		return nil, nil
	}

	// Convert to Eino tools
	einoTools := make([]tool.InvokableTool, 0, len(tools))
	for _, t := range tools {
		if t.Enabled {
			einoTools = append(einoTools, t.InvokableTool)
		}
	}

	// Convert to OpenAI format
	openAITools, err := ConvertToOpenAI(context.Background(), einoTools)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to OpenAI format: %w", err)
	}

	log.Printf("Generated %d tools", len(openAITools))

	return openAITools, nil
}

// GetEinoTools returns Eino InvokableTool interfaces
func (e *ToolsEngine) GetEinoTools(toolIDs []string) []tool.InvokableTool {
	if len(toolIDs) == 0 {
		return e.registry.GetEinoTools()
	}

	tools := e.registry.GetByIDs(toolIDs)
	einoTools := make([]tool.InvokableTool, 0, len(tools))
	for _, t := range tools {
		if t.Enabled {
			einoTools = append(einoTools, t.InvokableTool)
		}
	}
	return einoTools
}

// ExecuteTool executes a tool
func (e *ToolsEngine) ExecuteTool(ctx context.Context, toolID string, argsJSON string) (string, error) {
	return e.registry.Execute(ctx, toolID, argsJSON)
}

// GetRegistry returns the underlying registry
func (e *ToolsEngine) GetRegistry() *ToolRegistry {
	return e.registry
}
