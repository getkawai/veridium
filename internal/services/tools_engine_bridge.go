package services

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino/components/tool"
	"github.com/kawai-network/veridium/pkg/toolsengine"
)

// ToolsEngineBridge bridges the existing ToolsEngine with Eino's agent system
type ToolsEngineBridge struct {
	toolsEngine *toolsengine.ToolsEngine
}

// NewToolsEngineBridge creates a new tools engine bridge
func NewToolsEngineBridge(toolsEngine *toolsengine.ToolsEngine) *ToolsEngineBridge {
	return &ToolsEngineBridge{
		toolsEngine: toolsEngine,
	}
}

// GetToolsForAgent returns Eino tools for the agent
// If toolIDs is empty, returns all enabled tools
func (b *ToolsEngineBridge) GetToolsForAgent(toolIDs []string) []tool.BaseTool {
	if len(toolIDs) == 0 {
		// Get all enabled tools
		return b.getAllEnabledTools()
	}

	// Get specific tools by ID
	einoTools := b.toolsEngine.GetEinoTools(toolIDs)
	baseTools := make([]tool.BaseTool, len(einoTools))
	for i, t := range einoTools {
		baseTools[i] = t
	}

	log.Printf("📦 ToolsEngineBridge: Loaded %d tools for agent: %v", len(baseTools), toolIDs)
	return baseTools
}

// getAllEnabledTools returns all enabled tools as Eino BaseTool interfaces
func (b *ToolsEngineBridge) getAllEnabledTools() []tool.BaseTool {
	enabledTools := b.toolsEngine.GetEnabledTools()
	baseTools := make([]tool.BaseTool, 0, len(enabledTools))

	for _, t := range enabledTools {
		if t.Enabled {
			baseTools = append(baseTools, t.InvokableTool)
		}
	}

	log.Printf("📦 ToolsEngineBridge: Loaded %d enabled tools for agent", len(baseTools))
	return baseTools
}

// GetToolByID returns a single tool by ID
func (b *ToolsEngineBridge) GetToolByID(toolID string) (tool.BaseTool, error) {
	t, exists := b.toolsEngine.GetTool(toolID)
	if !exists {
		return nil, fmt.Errorf("tool not found: %s", toolID)
	}

	if !t.Enabled {
		return nil, fmt.Errorf("tool is disabled: %s", toolID)
	}

	return t.InvokableTool, nil
}

// GetAvailableToolIDs returns all available tool IDs (enabled)
func (b *ToolsEngineBridge) GetAvailableToolIDs() []string {
	enabledTools := b.toolsEngine.GetEnabledTools()
	ids := make([]string, len(enabledTools))
	for i, t := range enabledTools {
		ids[i] = t.ID
	}
	return ids
}

// GetAvailableToolNames returns all available tool names with descriptions
func (b *ToolsEngineBridge) GetAvailableToolNames(ctx context.Context) ([]ToolDescription, error) {
	enabledTools := b.toolsEngine.GetEnabledTools()
	descriptions := make([]ToolDescription, 0, len(enabledTools))

	for _, t := range enabledTools {
		info, err := t.InvokableTool.Info(ctx)
		if err != nil {
			log.Printf("⚠️  Warning: Failed to get info for tool %s: %v", t.ID, err)
			continue
		}

		descriptions = append(descriptions, ToolDescription{
			ID:          t.ID,
			Name:        info.Name,
			Description: info.Desc,
			Category:    t.Category,
		})
	}

	return descriptions, nil
}

// ExecuteTool executes a tool directly via the tools engine
func (b *ToolsEngineBridge) ExecuteTool(ctx context.Context, toolID string, argsJSON string) (string, error) {
	return b.toolsEngine.ExecuteTool(ctx, toolID, argsJSON)
}

// EnableTool enables a tool by ID
func (b *ToolsEngineBridge) EnableTool(toolID string) bool {
	return b.toolsEngine.EnableTool(toolID)
}

// DisableTool disables a tool by ID
func (b *ToolsEngineBridge) DisableTool(toolID string) bool {
	return b.toolsEngine.DisableTool(toolID)
}

// HasTool checks if a tool exists
func (b *ToolsEngineBridge) HasTool(toolID string) bool {
	return b.toolsEngine.HasTool(toolID)
}

// ToolDescription describes a tool
type ToolDescription struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}
