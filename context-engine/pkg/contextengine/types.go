package contextengine

// Message represents a chat message
type Message struct {
	ID        string
	Role      string                 // "user", "assistant", "system", "tool"
	Content   interface{}            // string or []ContentPart
	CreatedAt int64
	UpdatedAt int64
	Meta      map[string]interface{} // Additional metadata
}

// ContentPart represents a part of message content
type ContentPart struct {
	Type     string
	Text     *string
	ImageURL *ImageURL
	VideoURL *VideoURL
	Thinking *Thinking
}

// ImageURL represents an image in message content
type ImageURL struct {
	URL    string
	Detail string
}

// VideoURL represents a video in message content
type VideoURL struct {
	URL string
}

// Thinking represents reasoning content
type Thinking struct {
	Signature string
	Content   string
}

