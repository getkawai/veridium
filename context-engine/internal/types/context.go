package types

// ContextConfig holds configuration for context engineering
type ContextConfig struct {
	// System role configuration
	SystemRole string

	// History configuration
	EnableHistoryCount bool
	HistoryCount       int
	HistorySummary     string

	// Input template configuration
	InputTemplate string

	// Placeholder variables
	Variables map[string]interface{}

	// Tool configuration
	Tools []Tool

	// Message content processing configuration
	MessageContent MessageContentConfig

	// Model and provider information
	Model    string
	Provider string

	// Session information
	SessionID        string
	IsWelcomeQuestion bool
}

// MessageContentConfig holds configuration for message content processing
type MessageContentConfig struct {
	FileContext FileContextConfig
	IsCanUseVideo func(model, provider string) bool
	IsCanUseVision func(model, provider string) bool
}

// FileContextConfig holds file context configuration
type FileContextConfig struct {
	Enabled     bool
	IncludeFileURL bool
}

// Tool represents a tool/function that can be called
type Tool struct {
	ID          string
	Name        string
	Description string
	Parameters  map[string]interface{}
	Manifest    map[string]interface{}
}

// PipelineContext represents the context flowing through the pipeline
type PipelineContext struct {
	// Messages being processed
	Messages []*Message

	// Metadata for processor communication
	Metadata map[string]interface{}

	// Abort control
	IsAborted  bool
	AbortReason string

	// Initial state (immutable)
	InitialState *ContextState
}

// ContextState holds the initial state
type ContextState struct {
	Messages []*Message
	Model    string
	Provider string
	SystemRole string
	Tools    []Tool
}

// Clone creates a deep copy of the context
func (pc *PipelineContext) Clone() *PipelineContext {
	cloned := &PipelineContext{
		Messages:     make([]*Message, len(pc.Messages)),
		Metadata:     make(map[string]interface{}),
		IsAborted:    pc.IsAborted,
		AbortReason:  pc.AbortReason,
		InitialState: pc.InitialState,
	}

	// Copy messages
	for i, msg := range pc.Messages {
		cloned.Messages[i] = &Message{
			ID:        msg.ID,
			Role:      msg.Role,
			Content:   msg.Content,
			CreatedAt: msg.CreatedAt,
			UpdatedAt: msg.UpdatedAt,
			Meta:      make(map[string]interface{}),
		}
		for k, v := range msg.Meta {
			cloned.Messages[i].Meta[k] = v
		}
	}

	// Copy metadata
	for k, v := range pc.Metadata {
		cloned.Metadata[k] = v
	}

	return cloned
}

// Abort marks the context as aborted with a reason
func (pc *PipelineContext) Abort(reason string) {
	pc.IsAborted = true
	pc.AbortReason = reason
}

