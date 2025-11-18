package contextengine

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
)

// ContextEngineService provides context engineering capabilities
type ContextEngineService struct{}

// NewContextEngineService creates a new context engine service
func NewContextEngineService() *ContextEngineService {
	return &ContextEngineService{}
}

// ContextEngineeringRequest represents the request payload for context engineering
type ContextEngineeringRequest struct {
	Messages           []Message              `json:"messages"`
	Tools              []string               `json:"tools,omitempty"`
	Model              string                 `json:"model"`
	Provider           string                 `json:"provider"`
	SystemRole         string                 `json:"systemRole,omitempty"`
	InputTemplate      string                 `json:"inputTemplate,omitempty"`
	EnableHistoryCount bool                   `json:"enableHistoryCount,omitempty"`
	HistoryCount       int                    `json:"historyCount,omitempty"`
	HistorySummary     string                 `json:"historySummary,omitempty"`
	SessionID          string                 `json:"sessionId,omitempty"`
	IsWelcomeQuestion  bool                   `json:"isWelcomeQuestion,omitempty"`
	Variables          map[string]interface{} `json:"variables,omitempty"`
}

// ContextEngineeringResponse represents the response from context engineering
type ContextEngineeringResponse struct {
	Messages []Message `json:"messages"`
	Error    string    `json:"error,omitempty"`
}

// ProcessMessages processes messages through the context engineering pipeline
func (s *ContextEngineService) ProcessMessages(request ContextEngineeringRequest) ContextEngineeringResponse {
	log.Printf("ProcessMessages called with %d messages", len(request.Messages))

	// Convert tools from string IDs to Tool objects
	// For now, we'll create minimal tool objects - in production, you'd fetch full tool definitions
	tools := make([]Tool, len(request.Tools))
	for i, toolID := range request.Tools {
		tools[i] = Tool{
			ID:   toolID,
			Name: toolID, // In production, fetch actual tool name from store
		}
	}

	// Create context engine config
	config := Config{
		SystemRole:         request.SystemRole,
		EnableHistoryCount: request.EnableHistoryCount,
		HistoryCount:       request.HistoryCount,
		HistorySummary:     request.HistorySummary,
		InputTemplate:      request.InputTemplate,
		Variables:          request.Variables,
		Tools:              tools,
		MessageContent: MessageContentConfig{
			FileContext: FileContextConfig{
				Enabled:        true,  // TODO: Make configurable
				IncludeFileURL: false, // TODO: Make configurable based on isDesktop
			},
			// TODO: Add isCanUseVideo and isCanUseVision functions
			IsCanUseVideo:  nil,
			IsCanUseVision: nil,
		},
		Model:             request.Model,
		Provider:          request.Provider,
		SessionID:         request.SessionID,
		IsWelcomeQuestion: request.IsWelcomeQuestion,
	}

	// Create engine
	engine := New(config)

	// Convert request messages to engine messages
	engineMessages := make([]*Message, len(request.Messages))
	for i, msg := range request.Messages {
		engineMessages[i] = &Message{
			ID:        msg.ID,
			Role:      msg.Role,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
			Meta:      msg.Meta,
		}
	}

	// Process messages
	result, err := engine.Process(context.Background(), engineMessages)
	if err != nil {
		log.Printf("Error processing messages: %v", err)
		return ContextEngineeringResponse{
			Error: err.Error(),
		}
	}

	// Convert result back to response messages
	responseMessages := make([]Message, len(result))
	for i, msg := range result {
		responseMessages[i] = Message{
			ID:        msg.ID,
			Role:      msg.Role,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
			Meta:      msg.Meta,
		}
	}

	log.Printf("ProcessMessages completed with %d messages", len(responseMessages))

	return ContextEngineeringResponse{
		Messages: responseMessages,
	}
}

// GetEngineStats returns statistics about the context engine
func (s *ContextEngineService) GetEngineStats() map[string]interface{} {
	// Create a sample engine to get stats
	engine := New(Config{})
	stats := engine.GetStats()

	return map[string]interface{}{
		"processorCount":         stats.ProcessorCount,
		"customProcessorCount":   stats.CustomProcessorCount,
		"disabledProcessorCount": stats.DisabledProcessorCount,
		"processorNames":         stats.ProcessorNames,
		"customProcessorNames":   stats.CustomProcessorNames,
		"disabledProcessorNames": stats.DisabledProcessorNames,
	}
}

// ValidateConfig validates a context engine configuration
func (s *ContextEngineService) ValidateConfig(configJSON string) map[string]interface{} {
	var config Config
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return map[string]interface{}{
			"valid":  false,
			"errors": []string{fmt.Sprintf("Invalid JSON: %v", err)},
		}
	}

	engine := New(config)
	result := engine.Validate()

	return map[string]interface{}{
		"valid":  result.Valid,
		"errors": result.Errors,
	}
}
