package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kawai-network/veridium/pkg/yzma/message"
	"github.com/kawai-network/veridium/types"
)

// ToolNameMapping maps Yzma tool names to frontend-compatible identifier/apiName pairs
// This is critical for frontend to render tool results correctly
type ToolNameMapping struct {
	Identifier string // Frontend identifier (e.g., "lobe-web-browsing")
	APIName    string // Frontend API name (e.g., "search")
	Type       string // Tool type (e.g., "builtin")
}

// toolNameMappings maps Yzma tool names to frontend-compatible values
var toolNameMappings = map[string]ToolNameMapping{
	// Web search tools
	"web_search":                {Identifier: "lobe-web-browsing", APIName: "search", Type: "builtin"},
	"lobe-web-browsing__search": {Identifier: "lobe-web-browsing", APIName: "search", Type: "builtin"},

	// Web crawling tools
	"lobe-web-browsing__crawlSinglePage": {Identifier: "lobe-web-browsing", APIName: "crawlSinglePage", Type: "builtin"},
	"lobe-web-browsing__crawlMultiPages": {Identifier: "lobe-web-browsing", APIName: "crawlMultiPages", Type: "builtin"},

	// Local file system tools
	"lobe-local-system__listLocalFiles":   {Identifier: "lobe-local-system", APIName: "listLocalFiles", Type: "builtin"},
	"lobe-local-system__readLocalFile":    {Identifier: "lobe-local-system", APIName: "readLocalFile", Type: "builtin"},
	"lobe-local-system__searchLocalFiles": {Identifier: "lobe-local-system", APIName: "searchLocalFiles", Type: "builtin"},
	"lobe-local-system__writeLocalFile":   {Identifier: "lobe-local-system", APIName: "writeLocalFile", Type: "builtin"},
	"lobe-local-system__renameLocalFile":  {Identifier: "lobe-local-system", APIName: "renameLocalFile", Type: "builtin"},
	"lobe-local-system__moveLocalFiles":   {Identifier: "lobe-local-system", APIName: "moveLocalFiles", Type: "builtin"},

	// Image tools
	"lobe-image-designer__text2image": {Identifier: "lobe-image-designer", APIName: "text2image", Type: "builtin"},

	// Code interpreter
	"lobe-code-interpreter__python": {Identifier: "lobe-code-interpreter", APIName: "python", Type: "builtin"},

	// Calculator
	"calculator": {Identifier: "calculator", APIName: "calculate", Type: "builtin"},
}

// mapToolName maps a Yzma tool name to frontend-compatible identifier/apiName
// Returns the original name if no mapping exists
func mapToolName(yzmaToolName string) (identifier, apiName, toolType string) {
	if mapping, ok := toolNameMappings[yzmaToolName]; ok {
		return mapping.Identifier, mapping.APIName, mapping.Type
	}
	// Fallback: use tool name as both identifier and apiName
	return yzmaToolName, yzmaToolName, "builtin"
}

// minInt returns the minimum of two integers
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// stripToolCallTags removes <tool_call>...</tool_call> blocks from the text using regex
func stripToolCallTags(text string) string {
	// First, try to remove complete <tool_call>...</tool_call> blocks
	re := regexp.MustCompile(`(?s)<tool_call>.*?</tool_call>`)
	cleaned := re.ReplaceAllString(text, "")

	// If <tool_call> tag still exists (incomplete/unclosed), remove everything from <tool_call> onwards
	if strings.Contains(cleaned, "<tool_call>") {
		idx := strings.Index(cleaned, "<tool_call>")
		cleaned = cleaned[:idx]
	}

	return strings.TrimSpace(cleaned)
}

// ============================================
// ChatRealStream - Real LLM Event Streaming
// ============================================

// ChatRealStream handles chat with REAL LLM calls using event streaming.
// This combines:
// - Real LLM logic from Chat() in agent_chat_service.go
// - Streaming architecture from ChatMockStream
//
// Flow:
// 1. start - Generation begins
// 2. reasoning - Real thinking content from LLM (if reasoning model)
// 3. chunk - Real content chunks from LLM
// 4. tool_call - Real tool call initiated by LLM
// 5. tool_result - Real tool execution result
// 6. complete - Generation finished
//
// Frontend listens to 'chat:stream' events via Events.On()
// Data is saved to DB at the end.
//
// Usage from frontend:
//
//	await AgentChatService.ChatRealStream(request);
//	// No return value - data comes via events
//	// Events.On('chat:stream', handler) receives all updates
func (s *AgentChatService) ChatRealStream(ctx context.Context, req ChatRequest) error {
	log.Printf("🚀 [REAL STREAM] Starting real LLM streaming for session: %s", req.SessionID)
	startTime := time.Now()

	// Helper to emit events with type safety using StreamEventPayload
	emit := func(payload StreamEventPayload) {
		if s.app == nil {
			return
		}
		// Set common fields
		payload.SessionID = req.SessionID
		payload.MessageID = req.MessageAssistantID
		s.app.Event.Emit("chat:stream", payload)
	}

	// 1. Setup session, topic, and save user message
	setup, err := s.setupSessionAndTopic(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to setup session/topic: %w", err)
	}
	session := setup.Session
	currentTopicID := setup.TopicID

	// 2. Validate model for current reasoning mode and auto-switch if needed
	if err := s.ValidateModelForReasoningMode(); err != nil {
		log.Printf("⚠️  Model mismatch detected: %v", err)
		log.Printf("🔄 Auto-switching to recommended model...")
		if switchErr := s.SwitchToRecommendedModel(); switchErr != nil {
			log.Printf("❌ Failed to switch model: %v", switchErr)
			log.Printf("⚠️  Continuing with current model, but expect suboptimal performance")
		} else {
			log.Printf("✅ Successfully switched to recommended model")
		}
	}

	// 3. Emit START event
	emit(StreamEventPayload{
		Type:    types.ChatEventStart,
		TopicID: currentTopicID,
	})

	// 4. Prepare messages for LLM
	messagesWithSystem := s.prepareMessagesWithSystemPrompt(session.Messages, session)

	// Use pre-generated assistant message ID from frontend, or generate new one
	assistantMsgID := req.MessageAssistantID
	if assistantMsgID == "" {
		assistantMsgID = uuid.New().String()
	}

	// Get LLM generator with requested tools
	llmWithTools := s.llmGenerator.WithTools(session.ToolNames)

	// 5. Run LLM with streaming + tool execution
	var finalContent strings.Builder
	var reasoningContent strings.Builder
	var toolCalls []types.ToolCall
	var toolMessages []message.Message
	var usage *ModelUsage
	var llmResp interface{}
	var uiTools []ChatToolPayload
	var toolResultsData []ToolResultData

	// Track timing
	var ttft int64 // Time to first token

	// Track state for filtering tool_call tags
	var inToolCallTag bool

	// Streaming callback that emits events to frontend
	streamCallback := func(token string, isLast bool) {
		if token == "" && !isLast {
			return
		}

		// Measure TTFT
		if ttft == 0 && token != "" {
			ttft = time.Since(startTime).Milliseconds()
		}

		// Detect if this is reasoning content (inside <think> tags)
		currentContent := finalContent.String() + token

		// Check for reasoning mode patterns
		isInThinkTag := strings.Contains(currentContent, "<think>") && !strings.Contains(currentContent, "</think>")
		hasThinkContent := strings.Contains(currentContent, "<think>")

		if isInThinkTag || (hasThinkContent && !strings.Contains(currentContent, "</think>")) {
			// Extract reasoning content from <think> tag
			if strings.Contains(token, "<think>") {
				// Start of thinking
				reasoningContent.Reset()
			} else if strings.Contains(token, "</think>") {
				// End of thinking - emit final reasoning
				emit(StreamEventPayload{
					Type: types.ChatEventReasoning,
					Reasoning: &ModelReasoning{
						Content: reasoningContent.String(),
					},
				})
			} else if isInThinkTag {
				// Inside thinking - accumulate and emit
				reasoningContent.WriteString(token)
				emit(StreamEventPayload{
					Type: types.ChatEventReasoning,
					Reasoning: &ModelReasoning{
						Content: reasoningContent.String(),
					},
				})
			}
		} else {
			// Regular content (not in <think> tag)

			// Track tool_call tag state - must check BEFORE appending to content
			if strings.Contains(token, "<tool_call>") {
				inToolCallTag = true
				return // Skip the opening tag token entirely
			}
			if strings.Contains(token, "</tool_call>") {
				inToolCallTag = false
				return // Skip the closing tag token entirely
			}

			// Skip ALL content if we're inside a tool_call tag
			if inToolCallTag {
				return
			}

			finalContent.WriteString(token)

			// Clean content for display (remove any remaining tags)
			cleanContent := finalContent.String()
			cleanContent = strings.TrimPrefix(cleanContent, "</think>")
			cleanContent = strings.TrimSpace(cleanContent)

			if cleanContent != "" {
				emit(StreamEventPayload{
					Type:    types.ChatEventChunk,
					Content: cleanContent,
				})
			}
		}

		if isLast {
			log.Printf("📝 [REAL STREAM] Streaming complete, final content length: %d", finalContent.Len())
		}
	}

	// Tool event callback - emits tool events to frontend in real-time
	var toolCallIndex int
	toolEventCallback := func(eventType types.ChatStreamEvent, tc types.ToolCall, result string) {
		argsJSON, _ := json.Marshal(tc.Function.Arguments)
		identifier, apiName, toolType := mapToolName(tc.Function.Name)

		if eventType == types.ChatEventToolCall {
			// Tool call initiated - emit loading state
			toolCallID := fmt.Sprintf("%s_tool_%d", assistantMsgID, toolCallIndex)
			log.Printf("🔧 [REAL STREAM] Tool call (loading): %s -> identifier=%s, apiName=%s", tc.Function.Name, identifier, apiName)

			tool := ChatToolPayload{
				ID:         toolCallID,
				APIName:    apiName,
				Identifier: identifier,
				Arguments:  string(argsJSON),
				Type:       toolType,
			}
			uiTools = append(uiTools, tool)

			emit(StreamEventPayload{
				Type:  types.ChatEventToolCall,
				Tools: uiTools,
			})
			toolCallIndex++
		} else if eventType == types.ChatEventToolResult {
			// Tool execution completed - emit result
			log.Printf("🔧 [REAL STREAM] Tool result: %s -> %s", tc.Function.Name, result[:minInt(50, len(result))])

			// Parse result as JSON for state
			var toolState interface{}
			if err := json.Unmarshal([]byte(result), &toolState); err == nil {
				// Successfully parsed as JSON
			}

			toolResultsData = append(toolResultsData, ToolResultData{
				Content: result,
				State:   toolState,
			})

			// Find the tool index for this result
			resultIndex := len(toolResultsData) - 1
			toolCallID := fmt.Sprintf("%s_tool_%d", assistantMsgID, resultIndex)

			emit(StreamEventPayload{
				Type:       types.ChatEventToolResult,
				ToolCallID: toolCallID,
				ToolMsgID:  fmt.Sprintf("tool_msg_%s_%d", assistantMsgID, resultIndex+1),
				Plugin: &ChatPluginPayload{
					Identifier: identifier,
					APIName:    apiName,
					Arguments:  string(argsJSON),
					Type:       toolType,
				},
				PluginState: toolState,
				Content:     result,
			})
		}
	}

	// Run agent loop with streaming and tool callbacks
	resp, tms, runErr := llmWithTools.RunAgentLoopWithStreaming(ctx, messagesWithSystem, 10, streamCallback, toolEventCallback)
	if runErr != nil {
		log.Printf("❌ [REAL STREAM] Agent execution failed: %v", runErr)
		emit(StreamEventPayload{
			Type:    types.ChatEventComplete,
			Content: fmt.Sprintf("Error: %v", runErr),
			Error: &ChatMessageError{
				Type:    "LLMError",
				Message: runErr.Error(),
			},
		})
		return fmt.Errorf("agent execution failed: %w", runErr)
	}

	llmResp = resp
	toolCalls = resp.ToolCalls
	toolMessages = tms

	// Build usage from response
	if resp != nil {
		usage = &ModelUsage{
			TotalInputTokens:  resp.PromptTokens,
			TotalOutputTokens: resp.CompletionTokens,
			TotalTokens:       resp.TotalTokens,
		}
	}

	// 6. Clean final content - remove both think tags and tool_call tags
	finalContentStr := finalContent.String()
	finalContentStr = stripThinkTags(finalContentStr)
	finalContentStr = stripToolCallTags(finalContentStr)
	finalContentStr = strings.TrimSpace(finalContentStr)

	// If final content is empty but we have response content, use that
	if finalContentStr == "" && resp != nil && resp.Content != "" {
		finalContentStr = stripThinkTags(resp.Content)
		finalContentStr = stripToolCallTags(resp.Content)
		finalContentStr = strings.TrimSpace(finalContentStr)
	}

	// 8. Add messages to session history
	session.Messages = append(session.Messages, toolMessages...)
	if len(toolCalls) > 0 {
		session.Messages = append(session.Messages, message.Tool{
			Role:      "assistant",
			ToolCalls: toolCalls,
		})
	} else {
		session.Messages = append(session.Messages, message.Chat{
			Role:    "assistant",
			Content: finalContentStr,
		})
	}

	// 9. Build performance metrics
	duration := time.Since(startTime).Milliseconds()
	performance := &ModelPerformance{
		Duration: duration,
		Latency:  duration,
		TTFT:     ttft,
	}
	if usage != nil && usage.TotalOutputTokens > 0 && duration > 0 {
		performance.TPS = float64(usage.TotalOutputTokens) / (float64(duration) / 1000.0)
	}

	// 10. Build reasoning data if present
	var reasoning *ModelReasoning
	if reasoningContent.Len() > 0 {
		reasoning = &ModelReasoning{
			Content:  reasoningContent.String(),
			Duration: duration,
		}
	}

	// 11. Save assistant message to DB
	fullMetadata := map[string]interface{}{
		"model":       s.libService.GetLoadedChatModel(),
		"usage":       usage,
		"performance": performance,
		"llm_resp":    llmResp,
	}

	savedMsgID, err := s.saveAssistantMessage(ctx, SaveAssistantMessageParams{
		MessageID: assistantMsgID,
		Content:   finalContentStr,
		SessionID: req.SessionID,
		TopicID:   currentTopicID,
		ThreadID:  req.ThreadID,
		UserID:    req.UserID,
		Reasoning: reasoning,
		Tools:     uiTools,
		Metadata:  fullMetadata,
	})
	if err != nil {
		log.Printf("⚠️  Warning: Failed to save assistant message to DB: %v", err)
		savedMsgID = assistantMsgID
	} else {
		log.Printf("💾 Saved assistant message: %s (topic: %s, thread: %s)", savedMsgID, currentTopicID, req.ThreadID)
	}

	// 12. Save tool messages to DB
	for i, tool := range uiTools {
		result := toolResultsData[i]
		_, _ = s.saveToolMessage(ctx, SaveToolMessageParams{
			ToolCallID: tool.ID,
			Identifier: tool.Identifier,
			APIName:    tool.APIName,
			Arguments:  tool.Arguments,
			Content:    result.Content,
			State:      result.State,
			SessionID:  req.SessionID,
			TopicID:    currentTopicID,
			ThreadID:   req.ThreadID,
			UserID:     req.UserID,
			TimeOffset: int64(i + 2),
		})
	}

	// 13. Generate topic title after first response (background)
	if len(session.Messages) >= 2 && len(session.Messages) <= 4 {
		if currentTopicID != "" {
			err := s.updateTopicTitle(ctx, currentTopicID, session.UserID, session.Messages)
			if err != nil {
				log.Printf("⚠️  Warning: Failed to trigger topic title update: %v", err)
			}
		}
	}

	// 14. Update session timestamp
	if err := s.updateSessionTimestamp(ctx, session.SessionID, session.UserID); err != nil {
		log.Printf("⚠️  Warning: Failed to update session timestamp: %v", err)
	}

	// 15. Auto-summarize if needed (background)
	if currentTopicID != "" {
		go func() {
			bgCtx := context.Background()
			s.autoSummarizeIfNeeded(bgCtx, session, currentTopicID, session.UserID)
			s.incrementalSummarizeIfNeeded(bgCtx, session, currentTopicID, session.UserID)
		}()
	}

	// 16. Emit COMPLETE event with final data
	emit(StreamEventPayload{
		Type:        types.ChatEventComplete,
		Content:     finalContentStr,
		TopicID:     currentTopicID,
		Reasoning:   reasoning,
		Usage:       usage,
		Performance: performance,
		Tools:       uiTools,
	})

	totalTokens := 0
	if usage != nil {
		totalTokens = usage.TotalTokens
	}
	log.Printf("✅ [REAL STREAM] Complete - session: %s, duration: %dms, tokens: %d",
		req.SessionID, duration, totalTokens)

	return nil
}
