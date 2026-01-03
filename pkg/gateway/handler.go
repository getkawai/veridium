package gateway

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/pkg/config"
	"github.com/kawai-network/veridium/pkg/fantasy/llamalib"
	"github.com/kawai-network/veridium/pkg/store"
)

// Handler handles OpenAI-compatible API requests.
type Handler struct {
	llm             *llamalib.Service
	whisperExecutor *WhisperExecutor
	imageExecutor   ImageExecutor
	kvStore         *store.KVStore
}

// NewHandler creates a new Handler with the given services.
func NewHandler(llm *llamalib.Service, whisperExecutor *WhisperExecutor, imageExecutor ImageExecutor, kvStore *store.KVStore) *Handler {
	return &Handler{
		llm:             llm,
		whisperExecutor: whisperExecutor,
		imageExecutor:   imageExecutor,
		kvStore:         kvStore,
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

	// Get user address from auth middleware
	userAddress, ok := GetUserAddress(c)
	if !ok {
		h.sendError(c, http.StatusUnauthorized, "unauthorized", "user not authenticated")
		return
	}

	if req.Stream {
		h.handleStream(c, req, userAddress)
	} else {
		h.handleNonStream(c, req, userAddress)
	}
}

// handleNonStream handles non-streaming chat completion requests.
func (h *Handler) handleNonStream(c *gin.Context, req ChatCompletionRequest, userAddress string) {
	maxTokens := req.GetMaxTokens()

	// Get user address for billing
	ctx := context.Background()

	// Estimate input tokens for pre-check
	inputTokens := estimateTokens(req.Messages)
	estimatedTotalTokens := int64(inputTokens + maxTokens) // Worst case estimate

	// Pre-check: Ensure user has sufficient balance for estimated usage
	err := h.kvStore.CheckSufficientBalance(ctx, userAddress, estimatedTotalTokens)
	if err != nil {
		h.sendError(c, http.StatusPaymentRequired, "insufficient_balance", err.Error())
		return
	}

	// Check if tools are provided
	hasTools := len(req.Tools) > 0

	// Check if images are present in the request
	imageData, textPrompt := extractImageAndText(req.Messages)
	hasImages := len(imageData) > 0

	var responseMsg *ResponseMessage
	var finishReason string
	var genErr error

	if hasImages && h.llm.IsVLModelLoaded() {
		// Use Vision-Language model for image processing
		log.Printf("Processing request with VL model (image size: %d bytes)", len(imageData))

		content, vlErr := h.llm.ProcessImageBytesWithText(c.Request.Context(), imageData, textPrompt, int32(maxTokens))
		if vlErr != nil {
			h.sendError(c, http.StatusInternalServerError, "internal_error", vlErr.Error())
			return
		}

		responseMsg = &ResponseMessage{
			Role:    "assistant",
			Content: strings.TrimSpace(content),
		}
		finishReason = "stop"
	} else if hasTools {
		// Use tool-aware generation
		messages := convertToLlamaMessages(req.Messages)
		tools := convertToLlamaTools(req.Tools)

		result, toolErr := h.llm.GenerateWithTools(c.Request.Context(), messages, tools, int32(maxTokens))
		if toolErr != nil {
			h.sendError(c, http.StatusInternalServerError, "internal_error", toolErr.Error())
			return
		}

		responseMsg = &ResponseMessage{
			Role:    "assistant",
			Content: result.Text,
		}

		if result.HasTools && len(result.ToolCalls) > 0 {
			finishReason = "tool_calls"
			responseMsg.ToolCalls = convertToolCallResults(result.ToolCalls)
		} else {
			finishReason = "stop"
		}
	} else {
		// Standard generation without tools
		prompt, convErr := RequestToPrompt(&req)
		if convErr != nil {
			h.sendError(c, http.StatusBadRequest, "invalid_request_error", convErr.Error())
			return
		}

		promptText := PromptToText(prompt)

		content, textErr := h.llm.Generate(promptText, int32(maxTokens))
		if textErr != nil {
			h.sendError(c, http.StatusInternalServerError, "internal_error", textErr.Error())
			return
		}

		responseMsg = &ResponseMessage{
			Role:    "assistant",
			Content: strings.TrimSpace(content),
		}
		finishReason = "stop"
	}

	// Calculate actual token usage and deduct from balance
	promptTokens := estimateTokens(req.Messages)
	completionTokens := len(responseMsg.Content) / 4
	totalTokens := int64(promptTokens + completionTokens)

	// Deduct actual usage from user balance (atomic operation to prevent race conditions)
	cost := store.CalculateUsageCost(totalTokens)
	err = h.kvStore.DeductBalanceAtomic(ctx, userAddress, cost)
	if err != nil {
		log.Printf("Failed to deduct balance for user %s: %v", userAddress, err)
		h.sendError(c, http.StatusInternalServerError, "billing_error", "Failed to process payment")
		return
	}

	log.Printf("User %s used %d tokens, charged %s micro USDT", userAddress, totalTokens, cost.String())

	// Build response
	response := ChatCompletionResponse{
		ID:      "chatcmpl-" + uuid.New().String()[:8],
		Object:  "chat.completion",
		Created: time.Now().Unix(),
		Model:   getModelName(req.Model),
		Choices: []ResponseChoice{
			{
				Index:        0,
				Message:      responseMsg,
				FinishReason: finishReason,
			},
		},
		Usage: &Usage{
			PromptTokens:     promptTokens,
			CompletionTokens: completionTokens,
			TotalTokens:      promptTokens + completionTokens,
		},
	}

	_ = genErr // Suppress unused variable warning
	c.JSON(http.StatusOK, response)
}

// convertToLlamaMessages converts gateway messages to llamalib format.
func convertToLlamaMessages(messages []ChatMessage) []llamalib.Message {
	result := make([]llamalib.Message, len(messages))
	for i, msg := range messages {
		lm := llamalib.Message{
			Role:       msg.Role,
			Content:    msg.Content.GetText(),
			ToolCallID: msg.ToolCallID,
		}

		// Convert tool calls
		for _, tc := range msg.ToolCalls {
			lm.ToolCalls = append(lm.ToolCalls, llamalib.ToolCallResult{
				ID:        tc.ID,
				Name:      tc.Function.Name,
				Arguments: tc.Function.Arguments,
			})
		}

		result[i] = lm
	}
	return result
}

// convertToLlamaTools converts gateway tools to llamalib format.
func convertToLlamaTools(tools []Tool) []llamalib.Tool {
	result := make([]llamalib.Tool, 0, len(tools))
	for _, t := range tools {
		if t.Type != "function" {
			continue
		}
		result = append(result, llamalib.Tool{
			Name:        t.Function.Name,
			Description: t.Function.Description,
			Parameters:  t.Function.Parameters,
		})
	}
	return result
}

// convertToolCallResults converts llamalib tool calls to gateway format.
func convertToolCallResults(toolCalls []llamalib.ToolCallResult) []ToolCall {
	result := make([]ToolCall, len(toolCalls))
	for i, tc := range toolCalls {
		result[i] = ToolCall{
			ID:   tc.ID,
			Type: "function",
			Function: ToolCallFunction{
				Name:      tc.Name,
				Arguments: tc.Arguments,
			},
		}
	}
	return result
}

// handleStream handles streaming chat completion requests using SSE.
func (h *Handler) handleStream(c *gin.Context, req ChatCompletionRequest, userAddress string) {
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

// extractImageAndText extracts the first image and combined text from messages.
// Returns image data (decoded from base64) and the text prompt.
func extractImageAndText(messages []ChatMessage) ([]byte, string) {
	var imageData []byte
	var textParts []string

	for _, msg := range messages {
		// Skip non-user messages for image extraction
		if msg.Role != "user" {
			continue
		}

		// Check for images in content parts
		images := msg.Content.GetImages()
		if len(images) > 0 && len(imageData) == 0 {
			// Take the first image
			img := images[0]
			if len(img.Data) > 0 {
				imageData = img.Data
			}
		}

		// Collect text
		text := msg.Content.GetText()
		if text != "" {
			textParts = append(textParts, text)
		}
	}

	return imageData, strings.Join(textParts, "\n")
}

// ClaimTrial handles POST /v1/user/claim-trial
func (h *Handler) ClaimTrial(c *gin.Context) {
	// Get user address from auth middleware
	userAddress, ok := GetUserAddress(c)
	if !ok {
		h.sendError(c, http.StatusUnauthorized, "unauthorized", "user not authenticated")
		return
	}

	// Claim trial
	ctx := c.Request.Context()
	err := h.kvStore.ClaimFreeTrial(ctx, userAddress)
	if err != nil {
		if strings.Contains(err.Error(), "already claimed") {
			h.sendError(c, http.StatusConflict, "already_claimed", "Free trial already claimed")
			return
		}
		h.sendError(c, http.StatusInternalServerError, "server_error", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":             "success",
		"message":            "Free trial claimed successfully",
		"balance_added_usdt": config.GetFreeTrialAmountUSDT(),
	})
}
