package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// LLMExecutor defines the interface for executing LLM inference.
type LLMExecutor interface {
	// Execute runs inference and returns the response content.
	Execute(messages []ChatMessage) (string, error)
	// ExecuteStream runs inference with streaming response.
	ExecuteStream(messages []ChatMessage, stream chan<- string) error
}

// Handler handles OpenAI-compatible API requests.
type Handler struct {
	executor        LLMExecutor
	whisperExecutor *WhisperExecutor
	imageExecutor   ImageExecutor
	modelName       string
}

// NewHandler creates a new Handler with the given LLM and Whisper executors.
func NewHandler(executor LLMExecutor, whisperExecutor *WhisperExecutor, imageExecutor ImageExecutor, modelName string) *Handler {
	return &Handler{
		executor:        executor,
		whisperExecutor: whisperExecutor,
		imageExecutor:   imageExecutor,
		modelName:       modelName,
	}
}

// AudioTranscriptions handles POST /v1/audio/transcriptions
func (h *Handler) AudioTranscriptions(c *gin.Context) {
	if h.whisperExecutor == nil {
		h.sendError(c, http.StatusNotImplemented, "not_implemented", "Whisper service not available")
		return
	}

	// 1. Parse Multipart Form
	file, err := c.FormFile("file")
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "invalid_request_error", "file is required")
		return
	}

	model := c.PostForm("model")

	// 2. Transcribe
	text, err := h.whisperExecutor.Transcribe(c.Request.Context(), file, model)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	// 3. Return JSON Response
	c.JSON(http.StatusOK, TranscriptionResponse{
		Text: text,
	})
}

// ImageGenerations handles POST /v1/images/generations
func (h *Handler) ImageGenerations(c *gin.Context) {
	if h.imageExecutor == nil {
		h.sendError(c, http.StatusNotImplemented, "not_implemented", "Image generation service not available")
		return
	}

	var req ImageGenerationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "invalid_request_error", err.Error())
		return
	}

	// OpenAI spec: n defaults to 1 if not provided
	if req.N == 0 {
		req.N = 1
	}

	// 2. Generate
	data, err := h.imageExecutor.GenerateImage(c.Request.Context(), req)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	// 3. Return JSON Response
	c.JSON(http.StatusOK, ImageGenerationResponse{
		Created: time.Now().Unix(),
		Data:    data,
	})
}

// ChatCompletions handles POST /v1/chat/completions
func (h *Handler) ChatCompletions(c *gin.Context) {
	var req ChatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.sendError(c, http.StatusBadRequest, "invalid_request_error", err.Error())
		return
	}

	if len(req.Messages) == 0 {
		h.sendError(c, http.StatusBadRequest, "invalid_request_error", "messages is required")
		return
	}

	if req.Stream {
		h.handleStream(c, req)
	} else {
		h.handleNonStream(c, req)
	}
}

// handleNonStream handles non-streaming chat completion requests.
func (h *Handler) handleNonStream(c *gin.Context, req ChatCompletionRequest) {
	content, err := h.executor.Execute(req.Messages)
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}

	response := ChatCompletionResponse{
		ID:      "chatcmpl-" + uuid.New().String()[:8],
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   h.modelName,
		Choices: []Choice{
			{
				Index: 0,
				Message: &ChatMessage{
					Role:    "assistant",
					Content: content,
				},
				FinishReason: "stop",
			},
		},
		Usage: &Usage{
			PromptTokens:     estimateTokens(req.Messages),
			CompletionTokens: len(content) / 4,
			TotalTokens:      estimateTokens(req.Messages) + len(content)/4,
		},
	}

	c.JSON(http.StatusOK, response)
}

// handleStream handles streaming chat completion requests using SSE.
func (h *Handler) handleStream(c *gin.Context, req ChatCompletionRequest) {
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Transfer-Encoding", "chunked")

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		h.sendError(c, http.StatusInternalServerError, "internal_error", "streaming not supported")
		return
	}

	id := "chatcmpl-" + uuid.New().String()[:8]
	created := time.Now().Unix()

	streamChan := make(chan string, 100)
	errChan := make(chan error, 1)

	go func() {
		defer close(streamChan)
		if err := h.executor.ExecuteStream(req.Messages, streamChan); err != nil {
			errChan <- err
		}
	}()

	for {
		select {
		case content, ok := <-streamChan:
			if !ok {
				// Stream finished, send [DONE]
				fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
				flusher.Flush()
				return
			}

			chunk := ChatCompletionChunk{
				ID:      id,
				Object:  "chat.completion.chunk",
				Created: created,
				Model:   h.modelName,
				Choices: []Choice{
					{
						Index: 0,
						Delta: &ChatMessage{
							Content: content,
						},
					},
				},
			}

			data, _ := json.Marshal(chunk)
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			flusher.Flush()

		case err := <-errChan:
			errResp := ErrorResponse{
				Error: ErrorDetail{
					Message: err.Error(),
					Type:    "internal_error",
				},
			}
			data, _ := json.Marshal(errResp)
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			flusher.Flush()
			return

		case <-c.Request.Context().Done():
			return
		}
	}
}

// sendError sends an OpenAI-compatible error response.
func (h *Handler) sendError(c *gin.Context, status int, errType, message string) {
	c.JSON(status, ErrorResponse{
		Error: ErrorDetail{
			Message: message,
			Type:    errType,
		},
	})
}

// estimateTokens provides a rough estimate of token count.
func estimateTokens(messages []ChatMessage) int {
	total := 0
	for _, m := range messages {
		total += len(m.Content) / 4
	}
	return total
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"model":  h.modelName,
	})
}

// MockExecutor is a simple mock executor for testing.
type MockExecutor struct{}

func (m *MockExecutor) Execute(messages []ChatMessage) (string, error) {
	if len(messages) == 0 {
		return "", io.EOF
	}
	last := messages[len(messages)-1]
	return fmt.Sprintf("Mock response to: %s", last.Content), nil
}

func (m *MockExecutor) ExecuteStream(messages []ChatMessage, stream chan<- string) error {
	content, err := m.Execute(messages)
	if err != nil {
		return err
	}
	// Simulate streaming by sending word by word
	for i, r := range content {
		stream <- string(r)
		if i%10 == 0 {
			time.Sleep(10 * time.Millisecond)
		}
	}
	return nil
}
