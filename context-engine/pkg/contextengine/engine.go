package contextengine

import (
	"context"
	"fmt"

	"github.com/kawai-network/veridium/context-engine/internal/eino"
	"github.com/kawai-network/veridium/context-engine/internal/types"
)

// Engine is the main context engineering engine
type Engine struct {
	graphBuilder       *eino.GraphBuilder
	config             *types.ContextConfig
	customProcessors   []CustomProcessor
	disabledProcessors map[string]bool
}

// CustomProcessor represents a custom processor that can be added to the pipeline
type CustomProcessor struct {
	Name    string
	Process func(messages []*Message) ([]*Message, error)
	Order   int // Lower numbers run first
}

// New creates a new context engineering engine
func New(config Config) *Engine {
	internalConfig := &types.ContextConfig{
		SystemRole:         config.SystemRole,
		EnableHistoryCount: config.EnableHistoryCount,
		HistoryCount:       config.HistoryCount,
		HistorySummary:     config.HistorySummary,
		InputTemplate:      config.InputTemplate,
		Variables:          config.Variables,
		Tools:              convertTools(config.Tools),
		MessageContent: types.MessageContentConfig{
			FileContext: types.FileContextConfig{
				Enabled:        config.MessageContent.FileContext.Enabled,
				IncludeFileURL: config.MessageContent.FileContext.IncludeFileURL,
			},
			IsCanUseVideo:  config.MessageContent.IsCanUseVideo,
			IsCanUseVision: config.MessageContent.IsCanUseVision,
		},
		Model:             config.Model,
		Provider:          config.Provider,
		SessionID:         config.SessionID,
		IsWelcomeQuestion: config.IsWelcomeQuestion,
	}

	return &Engine{
		graphBuilder:       eino.NewGraphBuilder(internalConfig),
		config:             internalConfig,
		customProcessors:   []CustomProcessor{},
		disabledProcessors: make(map[string]bool),
	}
}

// Process processes messages through the context engineering pipeline
func (e *Engine) Process(ctx context.Context, messages []*Message) ([]*Message, error) {
	// Convert public Message to internal Message
	internalMessages := make([]*types.Message, len(messages))
	for i, msg := range messages {
		internalMessages[i] = &types.Message{
			ID:        msg.ID,
			Role:      msg.Role,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
			Meta:      msg.Meta,
		}
	}

	// Process through graph
	result, err := e.graphBuilder.ProcessMessages(ctx, internalMessages)
	if err != nil {
		return nil, err
	}

	// Convert back to public Message
	publicMessages := make([]*Message, len(result))
	for i, msg := range result {
		publicMessages[i] = &Message{
			ID:        msg.ID,
			Role:      msg.Role,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
			Meta:      msg.Meta,
		}
	}

	return publicMessages, nil
}

// convertTools converts public Tool to internal Tool
func convertTools(tools []Tool) []types.Tool {
	result := make([]types.Tool, len(tools))
	for i, tool := range tools {
		result[i] = types.Tool{
			ID:          tool.ID,
			Name:        tool.Name,
			Description: tool.Description,
			Parameters:  tool.Parameters,
			Manifest:    tool.Manifest,
		}
	}
	return result
}

// ============================================================================
// Pipeline Management APIs
// ============================================================================

// AddProcessor adds a custom processor to the pipeline
func (e *Engine) AddProcessor(processor CustomProcessor) *Engine {
	e.customProcessors = append(e.customProcessors, processor)
	return e
}

// RemoveProcessor removes a processor by name
func (e *Engine) RemoveProcessor(name string) *Engine {
	filtered := make([]CustomProcessor, 0, len(e.customProcessors))
	for _, p := range e.customProcessors {
		if p.Name != name {
			filtered = append(filtered, p)
		}
	}
	e.customProcessors = filtered
	return e
}

// DisableProcessor disables a built-in processor by name
func (e *Engine) DisableProcessor(name string) *Engine {
	e.disabledProcessors[name] = true
	return e
}

// EnableProcessor enables a previously disabled processor
func (e *Engine) EnableProcessor(name string) *Engine {
	delete(e.disabledProcessors, name)
	return e
}

// ClearCustomProcessors removes all custom processors
func (e *Engine) ClearCustomProcessors() *Engine {
	e.customProcessors = []CustomProcessor{}
	return e
}

// GetProcessors returns a list of all processor names (built-in + custom)
func (e *Engine) GetProcessors() []string {
	processors := []string{
		"GroupMessageFlatten",
		"HistoryTruncate",
		"SystemRoleInjector",
		"InboxGuide",
		"ToolSystemRole",
		"HistorySummary",
		"InputTemplate",
		"PlaceholderVariables",
		"MessageContent",
		"ToolCall",
		"ToolMessageReorder",
		"MessageCleanup",
	}

	// Add custom processors
	for _, p := range e.customProcessors {
		processors = append(processors, p.Name)
	}

	return processors
}

// Clone creates a deep copy of the engine with the same configuration
func (e *Engine) Clone() *Engine {
	// Clone config
	newConfig := &types.ContextConfig{
		SystemRole:         e.config.SystemRole,
		EnableHistoryCount: e.config.EnableHistoryCount,
		HistoryCount:       e.config.HistoryCount,
		HistorySummary:     e.config.HistorySummary,
		InputTemplate:      e.config.InputTemplate,
		Variables:          make(map[string]interface{}),
		Tools:              make([]types.Tool, len(e.config.Tools)),
		MessageContent:     e.config.MessageContent,
		Model:              e.config.Model,
		Provider:           e.config.Provider,
		SessionID:          e.config.SessionID,
		IsWelcomeQuestion:  e.config.IsWelcomeQuestion,
	}

	// Deep copy variables
	for k, v := range e.config.Variables {
		newConfig.Variables[k] = v
	}

	// Deep copy tools
	copy(newConfig.Tools, e.config.Tools)

	// Clone custom processors
	customProcessors := make([]CustomProcessor, len(e.customProcessors))
	copy(customProcessors, e.customProcessors)

	// Clone disabled processors
	disabledProcessors := make(map[string]bool)
	for k, v := range e.disabledProcessors {
		disabledProcessors[k] = v
	}

	return &Engine{
		graphBuilder:       eino.NewGraphBuilder(newConfig),
		config:             newConfig,
		customProcessors:   customProcessors,
		disabledProcessors: disabledProcessors,
	}
}

// ============================================================================
// Pipeline Validation API
// ============================================================================

// ValidationResult holds validation results
type ValidationResult struct {
	Valid  bool
	Errors []string
}

// Validate validates the pipeline configuration
func (e *Engine) Validate() ValidationResult {
	errors := []string{}

	// Check for duplicate processor names
	names := make(map[string]bool)
	for _, p := range e.customProcessors {
		if names[p.Name] {
			errors = append(errors, fmt.Sprintf("Duplicate processor name: %s", p.Name))
		}
		names[p.Name] = true
	}

	// Check if custom processors have valid names
	for _, p := range e.customProcessors {
		if p.Name == "" {
			errors = append(errors, "Processor missing name")
		}
		if p.Process == nil {
			errors = append(errors, fmt.Sprintf("Processor [%s] missing process function", p.Name))
		}
	}

	// Check configuration validity
	if e.config.EnableHistoryCount && e.config.HistoryCount < 0 {
		errors = append(errors, "HistoryCount must be non-negative")
	}

	return ValidationResult{
		Valid:  len(errors) == 0,
		Errors: errors,
	}
}

// ============================================================================
// Pipeline Statistics API
// ============================================================================

// Stats holds pipeline statistics
type Stats struct {
	ProcessorCount         int
	CustomProcessorCount   int
	DisabledProcessorCount int
	ProcessorNames         []string
	CustomProcessorNames   []string
	DisabledProcessorNames []string
}

// GetStats returns pipeline statistics
func (e *Engine) GetStats() Stats {
	builtInProcessors := []string{
		"GroupMessageFlatten",
		"HistoryTruncate",
		"SystemRoleInjector",
		"InboxGuide",
		"ToolSystemRole",
		"HistorySummary",
		"InputTemplate",
		"PlaceholderVariables",
		"MessageContent",
		"ToolCall",
		"ToolMessageReorder",
		"MessageCleanup",
	}

	customNames := make([]string, len(e.customProcessors))
	for i, p := range e.customProcessors {
		customNames[i] = p.Name
	}

	disabledNames := make([]string, 0, len(e.disabledProcessors))
	for name := range e.disabledProcessors {
		disabledNames = append(disabledNames, name)
	}

	allNames := append([]string{}, builtInProcessors...)
	allNames = append(allNames, customNames...)

	return Stats{
		ProcessorCount:         len(builtInProcessors) + len(e.customProcessors),
		CustomProcessorCount:   len(e.customProcessors),
		DisabledProcessorCount: len(e.disabledProcessors),
		ProcessorNames:         allNames,
		CustomProcessorNames:   customNames,
		DisabledProcessorNames: disabledNames,
	}
}
