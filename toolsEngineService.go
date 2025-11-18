package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/kawai-network/veridium/pkg/toolsengine"
	"github.com/kawai-network/veridium/pkg/toolsengine/builtin"
)

// ToolsEngineService provides tools engineering capabilities
type ToolsEngineService struct {
	engine *toolsengine.ToolsEngine
}

// NewToolsEngineService creates a new tools engine service
func NewToolsEngineService() *ToolsEngineService {
	// Create engine (Eino-native)
	engine, err := toolsengine.NewToolsEngine(toolsengine.Config{})
	if err != nil {
		log.Printf("Warning: Failed to initialize ToolsEngine: %v", err)
		return &ToolsEngineService{engine: nil}
	}

	// Register builtin tools (Eino-native)
	if err := builtin.RegisterAllBuiltinTools(engine); err != nil {
		log.Printf("Warning: Failed to register builtin tools: %v", err)
	}

	log.Printf("ToolsEngineService initialized with Eino-native tools")
	return &ToolsEngineService{engine: engine}
}

// GenerateToolsRequest represents the request for tool generation
type GenerateToolsRequest struct {
	ToolIDs  []string               `json:"toolIds"`
	Model    string                 `json:"model"`
	Provider string                 `json:"provider"`
	Context  map[string]interface{} `json:"context,omitempty"`
}

// GenerateToolsResponse represents the response from tool generation
type GenerateToolsResponse struct {
	Tools []toolsengine.ChatCompletionTool `json:"tools,omitempty"`
	Error string                           `json:"error,omitempty"`
}

// GenerateTools generates tools for chat completion
func (s *ToolsEngineService) GenerateTools(request GenerateToolsRequest) GenerateToolsResponse {
	if s.engine == nil {
		return GenerateToolsResponse{
			Error: "ToolsEngine not initialized",
		}
	}

	log.Printf("GenerateTools called: model=%s, provider=%s, toolIds=%v",
		request.Model, request.Provider, request.ToolIDs)

	tools, err := s.engine.GenerateTools(toolsengine.GenerateToolsParams{
		ToolIDs:  request.ToolIDs,
		Model:    request.Model,
		Provider: request.Provider,
		Context:  request.Context,
	})

	if err != nil {
		log.Printf("Error generating tools: %v", err)
		return GenerateToolsResponse{
			Error: err.Error(),
		}
	}

	log.Printf("Generated %d tools", len(tools))

	return GenerateToolsResponse{
		Tools: tools,
	}
}

// GetAvailableToolsResponse represents response from getting available tools
type GetAvailableToolsResponse struct {
	Tools []string `json:"tools"`
	Error string   `json:"error,omitempty"`
}

// GetAvailableTools returns all available tool IDs
func (s *ToolsEngineService) GetAvailableTools() GetAvailableToolsResponse {
	if s.engine == nil {
		return GetAvailableToolsResponse{
			Error: "ToolsEngine not initialized",
		}
	}

	tools := s.engine.GetToolIDs()

	log.Printf("GetAvailableTools returned %d tools", len(tools))

	return GetAvailableToolsResponse{
		Tools: tools,
	}
}

// GetToolStatsResponse represents response from getting tool stats
type GetToolStatsResponse struct {
	TotalTools   int      `json:"totalTools"`
	EnabledTools int      `json:"enabledTools"`
	ToolIDs      []string `json:"toolIds"`
	Error        string   `json:"error,omitempty"`
}

// GetToolStats returns statistics about registered tools
func (s *ToolsEngineService) GetToolStats() GetToolStatsResponse {
	if s.engine == nil {
		return GetToolStatsResponse{
			Error: "ToolsEngine not initialized",
		}
	}

	allTools := s.engine.GetAllTools()
	enabledTools := s.engine.GetEnabledTools()
	toolIDs := s.engine.GetToolIDs()

	return GetToolStatsResponse{
		TotalTools:   len(allTools),
		EnabledTools: len(enabledTools),
		ToolIDs:      toolIDs,
	}
}

// ExecuteToolRequest represents request to execute a tool
type ExecuteToolRequest struct {
	ToolID string                 `json:"toolId"`
	Args   map[string]interface{} `json:"args"`
}

// ExecuteToolResponse represents response from tool execution
type ExecuteToolResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// ExecuteTool executes a tool
func (s *ToolsEngineService) ExecuteTool(request ExecuteToolRequest) ExecuteToolResponse {
	if s.engine == nil {
		return ExecuteToolResponse{
			Error: "ToolsEngine not initialized",
		}
	}

	log.Printf("ExecuteTool called for: %s", request.ToolID)

	// Convert args to JSON
	argsJSON, err := json.Marshal(request.Args)
	if err != nil {
		return ExecuteToolResponse{
			Error: fmt.Sprintf("failed to marshal args: %v", err),
		}
	}

	// Execute tool (Eino-native)
	resultJSON, err := s.engine.ExecuteTool(context.Background(), request.ToolID, string(argsJSON))
	if err != nil {
		log.Printf("Error executing tool: %v", err)
		return ExecuteToolResponse{
			Error: err.Error(),
		}
	}

	// Parse result
	var result interface{}
	if err := json.Unmarshal([]byte(resultJSON), &result); err != nil {
		return ExecuteToolResponse{
			Error: fmt.Sprintf("failed to parse result: %v", err),
		}
	}

	log.Printf("Tool executed successfully: %s", request.ToolID)

	return ExecuteToolResponse{
		Result: result,
	}
}

// HasToolRequest represents request to check if tool exists
type HasToolRequest struct {
	ToolID string `json:"toolId"`
}

// HasToolResponse represents response from tool check
type HasToolResponse struct {
	HasTool bool   `json:"hasTool"`
	Error   string `json:"error,omitempty"`
}

// HasTool checks if a tool exists
func (s *ToolsEngineService) HasTool(request HasToolRequest) HasToolResponse {
	if s.engine == nil {
		return HasToolResponse{
			HasTool: false,
			Error:   "ToolsEngine not initialized",
		}
	}

	hasTool := s.engine.HasTool(request.ToolID)

	return HasToolResponse{
		HasTool: hasTool,
	}
}
