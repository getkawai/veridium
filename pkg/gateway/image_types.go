package gateway

// ImageGenerationRequest represents an OpenAI-compatible image generation request
// https://platform.openai.com/docs/api-reference/images/create
type ImageGenerationRequest struct {
	Prompt         string `json:"prompt" binding:"required"`
	Model          string `json:"model"`
	N              int    `json:"n,omitempty"`               // Default 1
	Quality        string `json:"quality,omitempty"`         // standard, hd
	ResponseFormat string `json:"response_format,omitempty"` // url, b64_json
	Size           string `json:"size,omitempty"`            // 1024x1024 etc.
	Style          string `json:"style,omitempty"`           // vivid, natural
	User           string `json:"user,omitempty"`
}

// ImageGenerationResponse represents an OpenAI-compatible image generation response
type ImageGenerationResponse struct {
	Created int64       `json:"created"`
	Data    []ImageData `json:"data"`
}

// ImageData represents a single generated image
type ImageData struct {
	B64JSON       string `json:"b64_json,omitempty"`
	URL           string `json:"url,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}
