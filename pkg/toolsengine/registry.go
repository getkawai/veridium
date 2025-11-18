package toolsengine

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/cloudwego/eino/components/tool"
)

// ToolRegistry manages Eino tools
type ToolRegistry struct {
	tools map[string]*Tool
	mu    sync.RWMutex
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]*Tool),
	}
}

// Register registers a tool
func (r *ToolRegistry) Register(t *Tool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if t.ID == "" {
		return fmt.Errorf("tool ID cannot be empty")
	}

	r.tools[t.ID] = t
	log.Printf("Registered tool: %s", t.ID)
	return nil
}

// Get retrieves a tool by ID
func (r *ToolRegistry) Get(id string) (*Tool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	t, exists := r.tools[id]
	return t, exists
}

// GetAll returns all tools
func (r *ToolRegistry) GetAll() []*Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]*Tool, 0, len(r.tools))
	for _, t := range r.tools {
		tools = append(tools, t)
	}
	return tools
}

// GetEnabled returns all enabled tools
func (r *ToolRegistry) GetEnabled() []*Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]*Tool, 0)
	for _, t := range r.tools {
		if t.Enabled {
			tools = append(tools, t)
		}
	}
	return tools
}

// GetByIDs returns tools for given IDs
func (r *ToolRegistry) GetByIDs(ids []string) []*Tool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	tools := make([]*Tool, 0, len(ids))
	for _, id := range ids {
		if t, exists := r.tools[id]; exists {
			tools = append(tools, t)
		}
	}
	return tools
}

// Remove removes a tool
func (r *ToolRegistry) Remove(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tools[id]; exists {
		delete(r.tools, id)
		log.Printf("Removed tool: %s", id)
		return true
	}
	return false
}

// Has checks if a tool exists
func (r *ToolRegistry) Has(id string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, exists := r.tools[id]
	return exists
}

// GetIDs returns all tool IDs
func (r *ToolRegistry) GetIDs() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	ids := make([]string, 0, len(r.tools))
	for id := range r.tools {
		ids = append(ids, id)
	}
	return ids
}

// Enable enables a tool
func (r *ToolRegistry) Enable(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if t, exists := r.tools[id]; exists {
		t.Enabled = true
		log.Printf("Enabled tool: %s", id)
		return true
	}
	return false
}

// Disable disables a tool
func (r *ToolRegistry) Disable(id string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()

	if t, exists := r.tools[id]; exists {
		t.Enabled = false
		log.Printf("Disabled tool: %s", id)
		return true
	}
	return false
}

// GetEinoTools returns Eino InvokableTool interfaces for enabled tools
func (r *ToolRegistry) GetEinoTools() []tool.InvokableTool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	einoTools := make([]tool.InvokableTool, 0)
	for _, t := range r.tools {
		if t.Enabled {
			einoTools = append(einoTools, t.InvokableTool)
		}
	}
	return einoTools
}

// Execute executes a tool
func (r *ToolRegistry) Execute(ctx context.Context, id string, argsJSON string) (string, error) {
	t, exists := r.Get(id)
	if !exists {
		return "", fmt.Errorf("tool not found: %s", id)
	}

	if !t.Enabled {
		return "", fmt.Errorf("tool is disabled: %s", id)
	}

	return t.InvokableRun(ctx, argsJSON)
}
