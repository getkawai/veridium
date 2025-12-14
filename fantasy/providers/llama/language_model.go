package llama

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kawai-network/veridium/pkg/xlog"

	"github.com/google/uuid"
	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/fantasy/llamalib"
	"github.com/kawai-network/veridium/fantasy/llamalib/llama"
	"github.com/kawai-network/veridium/fantasy/llamalib/template"
	"github.com/kawai-network/veridium/fantasy/llamalib/tools"
)

// Default batch size for token processing
const defaultBatchSize = 2048

type languageModel struct {
	provider     string
	modelID      string
	service      *llamalib.Service
	toolRegistry *tools.ToolRegistry
}

func newLanguageModel(modelID string, provider string, service *llamalib.Service, toolRegistry *tools.ToolRegistry) *languageModel {
	return &languageModel{
		modelID:      modelID,
		provider:     provider,
		service:      service,
		toolRegistry: toolRegistry,
	}
}

// Model implements fantasy.LanguageModel.
func (l *languageModel) Model() string {
	return l.modelID
}

// Provider implements fantasy.LanguageModel.
func (l *languageModel) Provider() string {
	return l.provider
}

// Generate implements fantasy.LanguageModel.
func (l *languageModel) Generate(ctx context.Context, call fantasy.Call) (*fantasy.Response, error) {
	// Convert fantasy.Call to prompt string using chat template
	prompt, err := l.preparePrompt(call)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare prompt: %w", err)
	}

	// Get max tokens from call or use default
	maxTokens := int32(32768)
	if call.MaxOutputTokens != nil {
		maxTokens = int32(*call.MaxOutputTokens)
	}

	// Generate response
	response, promptTokens, completionTokens, err := l.generateWithTokenCounts(ctx, prompt, maxTokens)
	if err != nil {
		return nil, fmt.Errorf("generation failed: %w", err)
	}

	// Parse tool calls from response
	toolCalls := tools.ParseToolCalls(response)

	// Build response content
	content := make([]fantasy.Content, 0)

	// Clean tool call tags from response text
	cleanedResponse := l.cleanToolCallTags(response)
	if cleanedResponse != "" {
		content = append(content, fantasy.TextContent{Text: cleanedResponse})
	}

	// Add tool calls as content
	for _, tc := range toolCalls {
		content = append(content, fantasy.ToolCallContent{
			ToolCallID:       tc.ID,
			ToolName:         tc.Name,
			Input:            tc.Input,
			ProviderExecuted: false,
		})
	}

	finishReason := fantasy.FinishReasonStop
	if len(toolCalls) > 0 {
		finishReason = fantasy.FinishReasonToolCalls
	}

	return &fantasy.Response{
		Content:      content,
		FinishReason: finishReason,
		Usage: fantasy.Usage{
			InputTokens:  int64(promptTokens),
			OutputTokens: int64(completionTokens),
			TotalTokens:  int64(promptTokens + completionTokens),
		},
	}, nil
}

// Stream implements fantasy.LanguageModel.
func (l *languageModel) Stream(ctx context.Context, call fantasy.Call) (fantasy.StreamResponse, error) {
	prompt, err := l.preparePrompt(call)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare prompt: %w", err)
	}

	maxTokens := int32(32768)
	if call.MaxOutputTokens != nil {
		maxTokens = int32(*call.MaxOutputTokens)
	}

	return func(yield func(fantasy.StreamPart) bool) {
		var responseBuilder strings.Builder
		textStarted := false
		promptTokens := len(prompt) / 4 // Approximate

		// Start text stream
		if !yield(fantasy.StreamPart{
			Type: fantasy.StreamPartTypeTextStart,
			ID:   "0",
		}) {
			return
		}
		textStarted = true

		// Generate with streaming
		completionTokens, err := l.generateStreaming(ctx, prompt, maxTokens, func(token string, done bool) {
			if done {
				return
			}
			responseBuilder.WriteString(token)
			yield(fantasy.StreamPart{
				Type:  fantasy.StreamPartTypeTextDelta,
				ID:    "0",
				Delta: token,
			})
		})

		if err != nil {
			yield(fantasy.StreamPart{
				Type:  fantasy.StreamPartTypeError,
				Error: err,
			})
			return
		}

		// End text stream
		if textStarted {
			if !yield(fantasy.StreamPart{
				Type: fantasy.StreamPartTypeTextEnd,
				ID:   "0",
			}) {
				return
			}
		}

		// Parse tool calls from complete response
		response := responseBuilder.String()
		toolCalls := tools.ParseToolCalls(response)

		// Emit tool calls
		for _, tc := range toolCalls {
			if !yield(fantasy.StreamPart{
				Type:         fantasy.StreamPartTypeToolInputStart,
				ID:           tc.ID,
				ToolCallName: tc.Name,
			}) {
				return
			}
			if !yield(fantasy.StreamPart{
				Type:  fantasy.StreamPartTypeToolInputDelta,
				ID:    tc.ID,
				Delta: tc.Input,
			}) {
				return
			}
			if !yield(fantasy.StreamPart{
				Type: fantasy.StreamPartTypeToolInputEnd,
				ID:   tc.ID,
			}) {
				return
			}
			if !yield(fantasy.StreamPart{
				Type:          fantasy.StreamPartTypeToolCall,
				ID:            tc.ID,
				ToolCallName:  tc.Name,
				ToolCallInput: tc.Input,
			}) {
				return
			}
		}

		finishReason := fantasy.FinishReasonStop
		if len(toolCalls) > 0 {
			finishReason = fantasy.FinishReasonToolCalls
		}

		yield(fantasy.StreamPart{
			Type:         fantasy.StreamPartTypeFinish,
			FinishReason: finishReason,
			Usage: fantasy.Usage{
				InputTokens:  int64(promptTokens),
				OutputTokens: int64(completionTokens),
				TotalTokens:  int64(promptTokens + completionTokens),
			},
		})
	}, nil
}

// GenerateObject implements fantasy.LanguageModel.
func (l *languageModel) GenerateObject(ctx context.Context, call fantasy.ObjectCall) (*fantasy.ObjectResponse, error) {
	// For now, use text-based object generation
	// TODO: Implement proper JSON mode if model supports it
	return nil, fmt.Errorf("GenerateObject not yet implemented for llama provider")
}

// StreamObject implements fantasy.LanguageModel.
func (l *languageModel) StreamObject(ctx context.Context, call fantasy.ObjectCall) (fantasy.ObjectStreamResponse, error) {
	return nil, fmt.Errorf("StreamObject not yet implemented for llama provider")
}

// preparePrompt converts fantasy.Call to a formatted prompt string
func (l *languageModel) preparePrompt(call fantasy.Call) (string, error) {
	// Convert fantasy.Prompt to template-compatible messages
	messages := l.convertToTemplateMessages(call.Prompt)

	// Enhance with tools if available and ToolChoice is not "none"
	// ToolChoiceNone disables all tools (useful for title/summary generation)
	shouldEnhanceTools := l.toolRegistry != nil && len(call.Tools) > 0
	if call.ToolChoice != nil && *call.ToolChoice == fantasy.ToolChoiceNone {
		shouldEnhanceTools = false
	}
	if shouldEnhanceTools {
		messages = l.enhanceWithTools(messages, call.Tools)
	}

	// Get chat template
	chatTemplate := l.getChatTemplate()

	// Apply chat template
	prompt, err := template.Apply(chatTemplate, messages, true)
	if err != nil {
		return "", fmt.Errorf("failed to apply template: %w", err)
	}

	return prompt, nil
}

// convertToTemplateMessages converts fantasy.Prompt to the format expected by template.Apply
func (l *languageModel) convertToTemplateMessages(prompt fantasy.Prompt) fantasy.Prompt {
	// fantasy.Prompt is already []fantasy.Message, which is compatible with template.Apply
	return prompt
}

// enhanceWithTools adds tool definitions to the prompt
func (l *languageModel) enhanceWithTools(messages fantasy.Prompt, callTools []fantasy.Tool) fantasy.Prompt {
	if len(callTools) == 0 {
		return messages
	}

	// Build tool definitions JSON
	var toolDefs []map[string]any
	for _, tool := range callTools {
		if ft, ok := tool.(fantasy.FunctionTool); ok {
			toolDefs = append(toolDefs, map[string]any{
				"type": "function",
				"function": map[string]any{
					"name":        ft.Name,
					"description": ft.Description,
					"parameters":  ft.InputSchema,
				},
			})
		}
	}

	if len(toolDefs) == 0 {
		return messages
	}

	// Format tools for prompt
	toolsJSON := tools.FormatToolsJSON(toolDefs)

	// Make a copy to avoid modifying original
	result := make(fantasy.Prompt, len(messages))
	copy(result, messages)

	// Find or create system message
	if len(result) > 0 && result[0].Role == fantasy.MessageRoleSystem {
		// Get existing system text
		var existingText string
		for _, part := range result[0].Content {
			if tp, ok := part.(fantasy.TextPart); ok {
				existingText = tp.Text
				break
			}
		}
		// Enhance existing system message
		enhancedContent := tools.BuildSystemPrompt(existingText, toolsJSON)
		result[0] = fantasy.NewSystemMessage(enhancedContent)
	} else {
		// Prepend new system message with tools
		systemContent := tools.BuildSystemPrompt("You are a helpful AI assistant.", toolsJSON)
		result = append(fantasy.Prompt{fantasy.NewSystemMessage(systemContent)}, result...)
	}

	return result
}

// getChatTemplate returns the appropriate chat template for the current model
func (l *languageModel) getChatTemplate() string {
	modelPath := l.service.GetLoadedChatModel()
	modelPathLower := strings.ToLower(modelPath)

	// Check if this is a Llama 3.2 model - use custom template
	if strings.Contains(modelPathLower, "llama-3.2") || strings.Contains(modelPathLower, "llama_3.2") || strings.Contains(modelPathLower, "llama3.2") {
		xlog.Info("🔧 Using custom Llama 3.2 tool template")
		return llamalib.Llama32ToolTemplate
	}

	// Use embedded template from model for other models
	return llama.ModelChatTemplate(l.service.GetChatModel(), "")
}

// generateWithTokenCounts generates response and returns token counts
func (l *languageModel) generateWithTokenCounts(ctx context.Context, prompt string, maxTokens int32) (string, int, int, error) {
	// Tokenize prompt
	vocab := l.service.GetChatVocab()
	tokens := llama.Tokenize(vocab, prompt, true, true)
	if len(tokens) == 0 {
		return "", 0, 0, fmt.Errorf("failed to tokenize prompt")
	}

	promptTokens := len(tokens)

	// Generate response
	response, err := l.service.Generate(prompt, maxTokens)
	if err != nil {
		return "", promptTokens, 0, err
	}

	// Estimate completion tokens
	completionTokens := len(response) / 4

	return response, promptTokens, completionTokens, nil
}

// generateStreaming performs token-by-token generation with callback
// Uses WithChatLock to ensure thread-safe access to chat model resources
func (l *languageModel) generateStreaming(ctx context.Context, prompt string, maxTokens int32, callback func(token string, done bool)) (int, error) {
	var completionTokens int
	var genErr error

	l.service.WithChatLock(func() {
		chatModel, chatContext, chatVocab, chatSampler := l.service.GetChatResourcesUnsafe()

		if chatModel == 0 || chatContext == 0 {
			genErr = fmt.Errorf("chat model not loaded")
			return
		}

		// Tokenize prompt
		tokens := llama.Tokenize(chatVocab, prompt, true, true)
		if len(tokens) == 0 {
			genErr = fmt.Errorf("failed to tokenize prompt")
			return
		}

		// Reset sampler state
		llama.SamplerReset(chatSampler)

		// Process prompt tokens in batches
		xlog.Debug("🔢 Processing prompt tokens in batches", "count", len(tokens), "batch_size", defaultBatchSize)

		for i := 0; i < len(tokens); i += defaultBatchSize {
			end := i + defaultBatchSize
			if end > len(tokens) {
				end = len(tokens)
			}

			batchTokens := tokens[i:end]
			batch := llama.BatchGetOne(batchTokens)

			errCode, err := llama.Decode(chatContext, batch)
			if err != nil || errCode != 0 {
				genErr = fmt.Errorf("failed to decode prompt batch %d-%d: %w", i, end, err)
				return
			}
		}

		// Generate tokens and stream
		for nGenerated := int32(0); nGenerated < maxTokens; nGenerated++ {
			// Check context cancellation
			select {
			case <-ctx.Done():
				callback("", true)
				genErr = ctx.Err()
				return
			default:
			}

			// Sample next token
			token := llama.SamplerSample(chatSampler, chatContext, -1)

			// Check for end of generation
			if llama.VocabIsEOG(chatVocab, token) {
				callback("", true)
				break
			}

			// Convert token to text
			buf := make([]byte, 256)
			length := llama.TokenToPiece(chatVocab, token, buf, 0, false)
			content := string(buf[:length])

			completionTokens++
			callback(content, false)

			// Accept token and prepare for next generation
			llama.SamplerAccept(chatSampler, token)

			// Decode the new token
			nextBatch := llama.BatchGetOne([]llama.Token{token})
			errCode, err := llama.Decode(chatContext, nextBatch)
			if err != nil || errCode != 0 {
				genErr = fmt.Errorf("failed to decode token: %w", err)
				return
			}

			// Small delay to prevent overwhelming
			time.Sleep(5 * time.Millisecond)
		}
	})

	return completionTokens, genErr
}

// cleanToolCallTags removes tool call XML tags from response
func (l *languageModel) cleanToolCallTags(response string) string {
	cleanResponse := response
	for strings.Contains(cleanResponse, "<tool_call>") {
		start := strings.Index(cleanResponse, "<tool_call>")
		end := strings.Index(cleanResponse, "</tool_call>")
		if start != -1 && end != -1 {
			cleanResponse = cleanResponse[:start] + cleanResponse[end+len("</tool_call>"):]
		} else {
			break
		}
	}
	return strings.TrimSpace(cleanResponse)
}

// generateToolCallID generates a unique ID for a tool call
func generateToolCallID() string {
	return "call_" + uuid.NewString()[:8]
}
