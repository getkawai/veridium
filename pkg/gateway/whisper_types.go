package gateway

// TranscriptionRequest represents an OpenAI Audio Transcription request.
// Only fields relevant for binding form data manually are defined here.
// In actual handler, we parse multipart/form-data.
type TranscriptionRequest struct {
	File           []byte  `form:"file" binding:"required"`
	Model          string  `form:"model" binding:"required"`
	Language       string  `form:"language"`
	Prompt         string  `form:"prompt"`
	ResponseFormat string  `form:"response_format"`
	Temperature    float32 `form:"temperature"`
}

// TranscriptionResponse represents an OpenAI Audio Transcription response.
type TranscriptionResponse struct {
	Text string `json:"text"`
}
