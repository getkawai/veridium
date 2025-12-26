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
	maxTokens := req.GetMaxTokens()

	// Convert OpenAI messages to fantasy Prompt
	prompt, err := RequestToPrompt(&req)
	if err != nil {
		h.sendError(c, http.StatusBadRequest, "invalid_request_error", err.Error())
		return
	}

	// Convert prompt to text for LLM
	promptText := PromptToText(prompt)

	// Check if request has images (for future VL model support)
	hasImages := false
	for _, msg := range req.Messages {
		if msg.Content.HasImages() {
			hasImages = true
			break
		}
	}

	var content string
	if hasImages {
		// TODO: Use VL model when images are present
		// For now, just process text
		content, err = h.llm.Generate(promptText, int32(maxTokens))
	} else {
		content, err = h.llm.Generate(promptText, int32(maxTokens))
	}

	if err != nil {
		h.sendError(c, http.StatusInternalServerError, "internal_error", err.Error())
		return
	}
	content = strings.TrimSpace(content)

	// Build response
	response := ChatCompletionResponse{
		ID:      "chatcmpl-" + uuid.New().String()[:8],
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   getModelName(req.Model),
		Choices: []ResponseChoice{
			{
				Index: 0,
				Message: &ResponseMessage{
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
	c.Header("X-Accel-Buffering", "no") // Disable nginx buffering

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		h.sendError(c, http.StatusInternalServerError, "internal_error", "streaming not supported")
		return
	}

	id := "chatcmpl-" + uuid.New().String()[:8]
	created := time.Now().Unix()
	model := getModelName(req.Model)

	// Send initial chunk with role
	initialChunk := ChatCompletionChunk{
		ID:      id,
		Object:  "chat.completion.chunk",
		Created: created,
		Model:   model,
		Choices: []StreamChoice{
			{
				Index: 0,
				Delta: &DeltaMessage{
					Role: "assistant",
				},
			},
		},
	}
	data, _ := json.Marshal(initialChunk)
	fmt.Fprintf(c.Writer, "data: %s\n\n", data)
	flusher.Flush()

	// Convert OpenAI messages to fantasy Prompt
	prompt, err := RequestToPrompt(&req)
	if err != nil {
		h.sendStreamError(c, flusher, err)
		return
	}

	// Convert prompt to text for LLM
	promptText := PromptToText(prompt)
	maxTokens := req.GetMaxTokens()

	promptTokens := estimateTokens(req.Messages)
	completionTokens := 0

	// Use real streaming with callback
	streamErr := h.llm.GenerateStream(c.Request.Context(), promptText, int32(maxTokens), func(token string) bool {
		// Check if client disconnected
		select {
		case <-c.Request.Context().Done():
			return false
		default:
		}

		completionTokens += len(token) / 4

		// Send token chunk
		chunk := ChatCompletionChunk{
			ID:      id,
			Object:  "chat.completion.chunk",
			Created: created,
			Model:   model,
			Choices: []StreamChoice{
				{
					Index: 0,
					Delta: &DeltaMessage{
						Content: token,
					},
				},
			},
		}

		data, _ := json.Marshal(chunk)
		fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		flusher.Flush()

		return true // Continue generation
	})

	if streamErr != nil {
		// Only send error if it's not a context cancellation
		if streamErr != c.Request.Context().Err() {
			h.sendStreamError(c, flusher, streamErr)
		}
		return
	}

	// Send finish reason
	finishChunk := ChatCompletionChunk{
		ID:      id,
		Object:  "chat.completion.chunk",
		Created: created,
		Model:   model,
		Choices: []StreamChoice{
			{
				Index:        0,
				Delta:        &DeltaMessage{},
				FinishReason: "stop",
			},
		},
	}
	data, _ = json.Marshal(finishChunk)
	fmt.Fprintf(c.Writer, "data: %s\n\n", data)
	flusher.Flush()

	// Send usage if requested
	if req.StreamOptions != nil && req.StreamOptions.IncludeUsage {
		usageChunk := ChatCompletionChunk{
			ID:      id,
			Object:  "chat.completion.chunk",
			Created: created,
			Model:   model,
			Choices: []StreamChoice{},
			Usage: &Usage{
				PromptTokens:     promptTokens,
				CompletionTokens: completionTokens,
				TotalTokens:      promptTokens + completionTokens,
			},
		}
		data, _ = json.Marshal(usageChunk)
		fmt.Fprintf(c.Writer, "data: %s\n\n", data)
		flusher.Flush()
	}

	// Send [DONE]
	fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
	flusher.Flush()
}

// sendStreamError sends an error response during streaming.
func (h *Handler) sendStreamError(c *gin.Context, flusher http.Flusher, err error) {
	errResp := ErrorResponse{
		Error: ErrorDetail{
			Message: err.Error(),
			Type:    "internal_error",
		},
	}
	data, _ := json.Marshal(errResp)
	fmt.Fprintf(c.Writer, "data: %s\n\n", data)
	flusher.Flush()
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
		// Get text from MessageContent (handles both string and array formats)
		text := m.Content.GetText()
		total += len(text) / 4

		// Account for tool calls
		for _, tc := range m.ToolCalls {
			total += len(tc.Function.Name)/4 + len(tc.Function.Arguments)/4
		}

		// Account for reasoning content
		if m.ReasoningContent != "" {
			total += len(m.ReasoningContent) / 4
		}
	}
	return total
}

// getModelName returns the model name to use in responses.
func getModelName(requestedModel string) string {
	if requestedModel == "" {
		return constant.KawaiAutoModel
	}
	return requestedModel
}

// HealthCheck handles GET /health
func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"model":  constant.KawaiAutoModel,
	})
}

// Models handles GET /v1/models - returns available models
func (h *Handler) Models(c *gin.Context) {
	models := []map[string]any{
		{
			"id":       constant.KawaiAutoModel,
			"object":   "model",
			"created":  time.Now().Unix(),
			"owned_by": "kawai-network",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"object": "list",
		"data":   models,
	})
}

// MockGenerate is a simple mock for testing without llamalib.
func MockGenerate(prompt string, maxTokens int32) (string, error) {
	if prompt == "" {
		return "", io.EOF
	}
	return fmt.Sprintf("Mock response to: %s", prompt), nil
}
