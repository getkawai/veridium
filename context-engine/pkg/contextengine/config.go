package contextengine

// Config holds configuration for the context engine
type Config struct {
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
	SessionID         string
	IsWelcomeQuestion bool
}

// MessageContentConfig holds configuration for message content processing
type MessageContentConfig struct {
	FileContext    FileContextConfig
	IsCanUseVideo  func(model, provider string) bool
	IsCanUseVision func(model, provider string) bool
}

// FileContextConfig holds file context configuration
type FileContextConfig struct {
	Enabled        bool
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

