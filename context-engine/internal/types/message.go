package types

// Message represents a chat message in the context pipeline
type Message struct {
	ID        string
	Role      string                 // "user", "assistant", "system", "tool", "group"
	Content   interface{}            // string or []ContentPart
	CreatedAt int64
	UpdatedAt int64
	Meta      map[string]interface{} // Additional metadata
	
	// Group message fields
	Children  []GroupChild           // For role="group" messages
	
	// Additional fields for message context
	ParentID  string
	ThreadID  string
	GroupID   string
	AgentID   string
	TargetID  string
	TopicID   string
	Reasoning *Thinking              // For reasoning models
}

// GroupChild represents a child message in a group message
type GroupChild struct {
	ID      string
	Content string
	Tools   []ToolCall
}

// ToolCall represents a tool/function call
type ToolCall struct {
	ID         string
	Type       string // "function"
	APIName    string
	Identifier string
	Arguments  string
	Result     *ToolResult
}

// ToolResult represents the result of a tool call
type ToolResult struct {
	ID      string
	Content string
	Error   string
	State   string
}

// ContentPart represents a part of message content (text, image, video, thinking)
type ContentPart struct {
	Type     string   // "text", "image_url", "video_url", "thinking"
	Text     *string  // For text type
	ImageURL *ImageURL
	VideoURL *VideoURL
	Thinking *Thinking
}

// ImageURL represents an image in message content
type ImageURL struct {
	URL    string
	Detail string // "auto", "low", "high"
}

// VideoURL represents a video in message content
type VideoURL struct {
	URL string
}

// Thinking represents reasoning content (for reasoning models)
type Thinking struct {
	Signature string
	Content   string
}

// IsText returns true if content is plain text
func (m *Message) IsText() bool {
	_, ok := m.Content.(string)
	return ok
}

// IsStructured returns true if content is structured (array of ContentPart)
func (m *Message) IsStructured() bool {
	_, ok := m.Content.([]ContentPart)
	return ok
}

// GetTextContent returns text content if available
func (m *Message) GetTextContent() string {
	if text, ok := m.Content.(string); ok {
		return text
	}
	return ""
}

// GetStructuredContent returns structured content if available
func (m *Message) GetStructuredContent() []ContentPart {
	if parts, ok := m.Content.([]ContentPart); ok {
		return parts
	}
	return nil
}

