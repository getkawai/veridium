package gateway

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/pkg/fantasy/llamalib"
)

// Handler handles OpenAI-compatible API requests.
type Handler struct {
	llm             *llamalib.Service
	whisperExecutor *WhisperExecutor
	imageExecutor   ImageExecutor
}

// NewHandler creates a new Handler with the given services.
func NewHandler(llm *llamalib.Service, whisperExecutor *WhisperExecutor, imageExecutor ImageExecutor) *Handler {
	return &Handler{
		llm:             llm,
		whisperExecutor: whisperExecutor,
		imageExecutor:   imageExecutor,
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
	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = 2048 // Default
	}

	prompt := formatMessagesToPrompt(req.Messages)
	content, err := h.llm.Generate(prompt, int32(maxTokens))
	if err != nil {
		h.sendError(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	content = strings.TrimSpace(content)

	response := ChatCompletionResponse{
		ID:      "chatcmpl-" + uuid.New().String()[:8],
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   constant.KawaiAutoModel,
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
		maxTokens := req.MaxTokens
		if maxTokens == 0 {
			maxTokens = 2048 // Default
		}

		prompt := formatMessagesToPrompt(req.Messages)
		result, err := h.llm.Generate(prompt, int32(maxTokens))
		if err != nil {
			errChan <- err
			return
		}

		// Simulate streaming by sending word by word
		words := strings.Fields(strings.TrimSpace(result))
		for i, word := range words {
			if i > 0 {
				streamChan <- " "
			}
			streamChan <- word
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
				Model:   constant.KawaiAutoModel,
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

// formatMessagesToPrompt converts chat messages to a single prompt string.
func formatMessagesToPrompt(messages []ChatMessage) string {
	var sb strings.Builder
	for _, msg := range messages {
		switch msg.Role {
		case "system":
			sb.WriteString(fmt.Sprintf("System: %s\n", msg.Content))
		case "user":
			sb.WriteString(fmt.Sprintf("User: %s\n", msg.Content))
		case "assistant":
			sb.WriteString(fmt.Sprintf("Assistant: %s\n", msg.Content))
		default:
			sb.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
		}
	}
	sb.WriteString("Assistant:")
	return sb.String()
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"model":  constant.KawaiAutoModel,
	})
}

// MockGenerate is a simple mock for testing without llamalib.
func MockGenerate(prompt string, maxTokens int32) (string, error) {
	if prompt == "" {
		return "", io.EOF
	}
	return fmt.Sprintf("Mock response to: %s", prompt), nil
}
