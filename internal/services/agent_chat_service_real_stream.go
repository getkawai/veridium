package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/pkg/xlog"
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
	// Web search tool
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

	// Image describe (VL description from uploaded images)
	"lobe-image-describe__getImageDescription": {Identifier: "lobe-image-describe", APIName: "getImageDescription", Type: "builtin"},
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
	xlog.Info("🚀 [REAL STREAM] Starting real LLM streaming", "session_id", req.SessionID)
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
	if err := s.validateModelForReasoningMode(); err != nil {
		xlog.Warn("⚠️  Model mismatch detected", "error", err)
		xlog.Info("🔄 Auto-switching to recommended model")
		if switchErr := s.switchToRecommendedModel(); switchErr != nil {
			xlog.Error("❌ Failed to switch model", "error", switchErr)
			xlog.Warn("⚠️  Continuing with current model, but expect suboptimal performance")
		} else {
			xlog.Info("✅ Successfully switched to recommended model")
		}
	}

	// 3. Emit START event
	emit(StreamEventPayload{
		Type:    types.ChatEventStart,
		TopicID: currentTopicID,
	})

	// 3.5. Perform semantic search on attached files (if any)
	var fileChunks []ChatFileChunk
	if len(req.FileIDs) > 0 && s.vectorSearch != nil {
		xlog.Info("📎 [REAL STREAM] Searching attached files for context", "files", len(req.FileIDs))
		searchResults, err := s.vectorSearch.SemanticSearch(ctx, req.Message, req.FileIDs, 10)
		if err != nil {
			xlog.Warn("⚠️  [REAL STREAM] File search failed", "error", err)
		} else if len(searchResults) > 0 {
			xlog.Info("📚 [REAL STREAM] Found relevant chunks from attached files", "chunks", len(searchResults))
			for _, result := range searchResults {
				fileChunks = append(fileChunks, ChatFileChunk{
					ID:         result.ID,
					FileID:     result.FileID,
					Filename:   result.FileName,
					Text:       result.Text,
					Similarity: float64(result.Similarity), // Convert float32 to float64
				})
			}
		}

		// Fallback: If RAG returned no results but we have file IDs, try to get document content directly
		// For images/videos, VL description may still be processing async - poll with timeout
		if len(fileChunks) == 0 {
			xlog.Warn("⚠️  [REAL STREAM] No RAG results, falling back to direct document fetch with polling")
			for _, fileID := range req.FileIDs {
				chunk := s.waitForDocumentContent(ctx, fileID, 60*time.Second)
				if chunk != nil {
					fileChunks = append(fileChunks, *chunk)
				}
			}
		}
	}

	// 4. Build hybrid context from memory (if available)
	var memoryContext string
	if s.memoryIntegration != nil {
		xlog.Info("🧠 [REAL STREAM] Fetching hybrid context memories", "query", req.Message)
		// We use nil for shortTermMessages here as they are already in the session history
		// and handled by the agent. We only need the memory text to inject into system prompt.
		memCtx, err := s.memoryIntegration.BuildHybridContext(ctx, req.Message, nil)
		if err != nil {
			xlog.Warn("⚠️  [REAL STREAM] Failed to build hybrid context", "error", err)
		} else if memCtx != "" {
			memoryContext = memCtx
			xlog.Info("🧠 [REAL STREAM] Hybrid context retrieved", "chars", len(memoryContext))
		}
	}

	// 5. Build system prompt and get history messages (optimized for fantasy.Agent)
	// Pass user message for language detection to respond in the same language
	systemPrompt := s.buildSystemPrompt(session, memoryContext, req.Message)
	historyMessages := s.getHistoryMessages(session)

	// 4.5. Build file context if we have relevant chunks
	var fileContext string
	if len(fileChunks) > 0 {
		fileContext = "\n\n## Relevant Context from Attached Files:\n"
		for i, chunk := range fileChunks {
			fileContext += fmt.Sprintf("\n### [%d] From: %s (file_id: %s, similarity: %.2f)\n%s\n", i+1, chunk.Filename, chunk.FileID, chunk.Similarity, chunk.Text)
		}
		fileContext += "\n---\nUse the above context to help answer the user's question. When using tools like getImageDescription, use the file_id provided above.\n"
		xlog.Info("📝 [REAL STREAM] Prepared file chunks for context", "chunks", len(fileChunks))
	}

	// Build user prompt with file context if available
	userPrompt := req.Message
	if fileContext != "" {
		userPrompt = req.Message + fileContext
	}

	// Use pre-generated assistant message ID from frontend, or generate new one
	assistantMsgID := req.MessageAssistantID
	if assistantMsgID == "" {
		assistantMsgID = uuid.New().String()
	}

	// State for streaming
	var finalContent strings.Builder
	var reasoningContent strings.Builder
	var toolCalls []fantasy.ToolCallContent
	var toolMessages []fantasy.Message
	var usage *ModelUsage
	var uiTools []ChatToolPayload
	var toolResultsData []ToolResultData
	var ttft int64 // Time to first token
	var inToolCallTag bool
	var toolCallIndex int
	var mu sync.Mutex // Protect concurrent access

	// Get LanguageModel from chatModel (set via SetChatModel, can be ChainLanguageModel for fallback)
	model := s.chatModel

	// Use fantasy.Agent if we have a LanguageModel
	if model != nil {
		// Convert tools from ToolRegistry to fantasy.AgentTool
		agentTools := s.toolRegistry.ToAgentTools(session.ToolNames)

		// 5. Create fantasy.Agent with tools, system prompt, and repair function
		// fantasy.Agent handles: system prompt injection, history, and current prompt internally
		agent := fantasy.NewAgent(model,
			fantasy.WithTools(agentTools...),
			fantasy.WithSystemPrompt(systemPrompt),
			fantasy.WithStopConditions(fantasy.StepCountIs(10)), // Max 10 iterations
			fantasy.WithRepairToolCall(RepairToolCall),          // Auto-repair malformed tool calls
		)

		// Run agent with streaming callbacks
		// fantasy.Agent.createPrompt will build: [system] + historyMessages + [user: userPrompt]
		// Disable agent-level retry when using Chain (Chain has its own fallback mechanism)
		var maxRetries *int
		if _, isChain := model.(*fantasy.ChainLanguageModel); isChain {
			zero := 0
			maxRetries = &zero
		}
		result, runErr := agent.Stream(ctx, fantasy.AgentStreamCall{
			Prompt:     userPrompt,
			Messages:   historyMessages,
			MaxRetries: maxRetries,

			// Text streaming callbacks
			OnTextDelta: func(id, text string) error {
				mu.Lock()
				defer mu.Unlock()

				// Measure TTFT
				if ttft == 0 && text != "" {
					ttft = time.Since(startTime).Milliseconds()
				}

				// Track tool_call tag state
				if strings.Contains(text, "<tool_call>") {
					inToolCallTag = true
					return nil
				}
				if strings.Contains(text, "</tool_call>") {
					inToolCallTag = false
					return nil
				}
				if inToolCallTag {
					return nil
				}

				finalContent.WriteString(text)

				// Clean content for display
				cleanContent := finalContent.String()
				cleanContent = strings.TrimPrefix(cleanContent, "</think>")
				cleanContent = strings.TrimSpace(cleanContent)

				if cleanContent != "" {
					emit(StreamEventPayload{
						Type:    types.ChatEventChunk,
						Content: cleanContent,
					})
				}
				return nil
			},

			// Reasoning callbacks - for models like DeepSeek R1, o1, etc.
			OnReasoningStart: func(id string, content fantasy.ReasoningContent) error {
				mu.Lock()
				defer mu.Unlock()

				emit(StreamEventPayload{
					Type: types.ChatEventReasoningStart,
					Reasoning: &ModelReasoning{
						Content: content.Text,
					},
				})
				return nil
			},

			OnReasoningDelta: func(id, text string) error {
				mu.Lock()
				defer mu.Unlock()

				reasoningContent.WriteString(text)
				emit(StreamEventPayload{
					Type: types.ChatEventReasoning,
					Reasoning: &ModelReasoning{
						Content: reasoningContent.String(),
					},
				})
				return nil
			},

			OnReasoningEnd: func(id string, content fantasy.ReasoningContent) error {
				mu.Lock()
				defer mu.Unlock()

				emit(StreamEventPayload{
					Type: types.ChatEventReasoningEnd,
					Reasoning: &ModelReasoning{
						Content: content.Text,
					},
				})
				return nil
			},

			// Tool call callback - when tool call is detected
			OnToolCall: func(tc fantasy.ToolCallContent) error {
				mu.Lock()
				defer mu.Unlock()

				identifier, apiName, toolType := mapToolName(tc.ToolName)
				toolCallID := fmt.Sprintf("%s_tool_%d", assistantMsgID, toolCallIndex)
				xlog.Info("🔧 [REAL STREAM] Tool call (loading)", "tool", tc.ToolName, "identifier", identifier, "api_name", apiName)

				tool := ChatToolPayload{
					ID:         toolCallID,
					APIName:    apiName,
					Identifier: identifier,
					Arguments:  tc.Input,
					Type:       toolType,
				}
				uiTools = append(uiTools, tool)
				toolCalls = append(toolCalls, tc)

				emit(StreamEventPayload{
					Type:  types.ChatEventToolCall,
					Tools: uiTools,
				})
				toolCallIndex++
				return nil
			},

			// Tool result callback - when tool execution completes
			OnToolResult: func(tr fantasy.ToolResultContent) error {
				mu.Lock()
				defer mu.Unlock()

				identifier, apiName, toolType := mapToolName(tr.ToolName)

				// Get result content as string
				resultContent := ""
				switch r := tr.Result.(type) {
				case fantasy.ToolResultOutputContentText:
					resultContent = r.Text
				case fantasy.ToolResultOutputContentError:
					if r.Error != nil {
						resultContent = r.Error.Error()
					}
				}

				if len(resultContent) > 50 {
					xlog.Info("🔧 [REAL STREAM] Tool result", "tool", tr.ToolName, "result", resultContent[:50]+"...")
				} else {
					xlog.Info("🔧 [REAL STREAM] Tool result", "tool", tr.ToolName, "result", resultContent)
				}

				// Parse result as JSON for state
				var toolState interface{}
				if err := json.Unmarshal([]byte(resultContent), &toolState); err == nil {
					// Successfully parsed as JSON
				}

				toolResultsData = append(toolResultsData, ToolResultData{
					Content: resultContent,
					State:   toolState,
				})

				resultIndex := len(toolResultsData) - 1
				toolCallID := fmt.Sprintf("%s_tool_%d", assistantMsgID, resultIndex)

				// Find the corresponding tool input
				toolInput := ""
				if resultIndex < len(uiTools) {
					toolInput = uiTools[resultIndex].Arguments
				}

				emit(StreamEventPayload{
					Type:       types.ChatEventToolResult,
					ToolCallID: toolCallID,
					ToolMsgID:  fmt.Sprintf("tool_msg_%s_%d", assistantMsgID, resultIndex+1),
					Plugin: &ChatPluginPayload{
						Identifier: identifier,
						APIName:    apiName,
						Arguments:  toolInput,
						Type:       toolType,
					},
					PluginState: toolState,
					Content:     resultContent,
				})
				return nil
			},

			// Step finish callback - build tool messages from step
			OnStepFinish: func(step fantasy.StepResult) error {
				mu.Lock()
				defer mu.Unlock()

				// Add step messages to toolMessages
				toolMessages = append(toolMessages, step.Messages...)
				return nil
			},

			OnError: func(err error) {
				xlog.Error("❌ [REAL STREAM] Agent stream error", "error", err)
			},
		})

		if runErr != nil {
			xlog.Error("❌ [REAL STREAM] Agent execution failed", "error", runErr)
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

		// Build usage from result
		if result != nil {
			usage = &ModelUsage{
				TotalInputTokens:  int(result.TotalUsage.InputTokens),
				TotalOutputTokens: int(result.TotalUsage.OutputTokens),
				TotalTokens:       int(result.TotalUsage.TotalTokens),
			}
			// Update toolCalls from final response if not captured via callbacks
			if len(toolCalls) == 0 {
				toolCalls = result.Response.Content.ToolCalls()
			}
		}
	} else {
		// No language model available
		xlog.Error("❌ [REAL STREAM] No language model available")
		emit(StreamEventPayload{
			Type:    types.ChatEventComplete,
			Content: "Error: No language model available",
			Error: &ChatMessageError{
				Type:    "ConfigError",
				Message: "No language model configured. Please configure a provider or load a local model.",
			},
		})
		return fmt.Errorf("no language model available")
	}

	// 6. Clean final content - remove both think tags and tool_call tags
	finalContentStr := finalContent.String()
	finalContentStr = stripThinkTags(finalContentStr)
	finalContentStr = stripToolCallTags(finalContentStr)
	finalContentStr = strings.TrimSpace(finalContentStr)

	// 8. Add messages to session history
	session.Messages = append(session.Messages, toolMessages...)
	if len(toolCalls) > 0 {
		session.Messages = append(session.Messages, types.NewToolCallMessageFromContent(toolCalls))
	} else {
		session.Messages = append(session.Messages, fantasy.Message{
			Role:    fantasy.MessageRoleAssistant,
			Content: []fantasy.MessagePart{fantasy.TextPart{Text: finalContentStr}},
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
	}

	savedMsgID, err := s.saveAssistantMessage(ctx, SaveAssistantMessageParams{
		MessageID: assistantMsgID,
		Content:   finalContentStr,
		SessionID: req.SessionID,
		TopicID:   currentTopicID,
		ThreadID:  req.ThreadID,
		Reasoning: reasoning,
		Tools:     uiTools,
		Metadata:  fullMetadata,
	})
	if err != nil {
		xlog.Warn("⚠️  Warning: Failed to save assistant message to DB", "error", err)
		savedMsgID = assistantMsgID
	} else {
		xlog.Info("💾 Saved assistant message", "id", savedMsgID, "topic", currentTopicID, "thread", req.ThreadID)
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
			TimeOffset: int64(i + 2),
		})
	}

	// 13. Generate topic title after first response (background)
	// Count user messages to determine turn count (tool messages don't count as turns)
	userMsgCount := 0
	for _, msg := range session.Messages {
		if types.GetMessageRole(msg) == "user" {
			userMsgCount++
		}
	}
	xlog.Debug("📌 [TITLE CHECK]", "message_count", len(session.Messages), "user_msg_count", userMsgCount, "topic_id", currentTopicID)

	// Generate title on first turn (1 user message = first conversation)
	if userMsgCount >= 1 && userMsgCount <= 2 {
		if currentTopicID != "" {
			xlog.Info("📌 [TITLE CHECK] Conditions met (first turn), calling updateTopicTitle")
			if s.topicService != nil {
				err := s.topicService.UpdateTopicTitle(ctx, currentTopicID, session.Messages)
				if err != nil {
					xlog.Warn("⚠️  Warning: Failed to trigger topic title update", "error", err)
				}
			} else {
				xlog.Warn("⚠️  TopicService not initialized, skipping title update")
			}
		} else {
			xlog.Debug("📌 [TITLE CHECK] Skipped - no topicID")
		}
	} else {
		xlog.Debug("📌 [TITLE CHECK] Skipped - not first turn", "userMsgCount", userMsgCount)
	}

	// 14. Update session timestamp
	if err := s.updateSessionTimestamp(ctx, session.SessionID); err != nil {
		xlog.Warn("⚠️  Warning: Failed to update session timestamp", "error", err)
	}

	// 15. Auto-summarize if needed (background)
	if currentTopicID != "" {
		go func() {
			bgCtx := context.Background()
			s.autoSummarizeIfNeeded(bgCtx, session, currentTopicID)
			s.incrementalSummarizeIfNeeded(bgCtx, session, currentTopicID)
		}()
	}

	// 15.5. Store conversation to memory (background) - MemGPT-style
	if s.memoryIntegration != nil && finalContentStr != "" {
		go func() {
			bgCtx := context.Background()
			if err := s.memoryIntegration.StoreConversationMemory(bgCtx, req.Message, finalContentStr); err != nil {
				xlog.Warn("⚠️  [Memory] Failed to store conversation", "error", err)
			}
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
		ChunksList:  fileChunks, // Include RAG chunks from file attachments
	})

	totalTokens := 0
	if usage != nil {
		totalTokens = usage.TotalTokens
	}
	xlog.Info("✅ [REAL STREAM] Complete", "session", req.SessionID, "duration_ms", duration, "tokens", totalTokens)

	return nil
}

// waitForDocumentContent polls for document content with VL description.
// For images/videos, VL processing is async and may take time.
// This polls every 2 seconds for up to maxWait duration.
// Returns nil if timeout or no meaningful content found.
func (s *AgentChatService) waitForDocumentContent(ctx context.Context, fileID string, maxWait time.Duration) *ChatFileChunk {
	const pollInterval = 2 * time.Second
	deadline := time.Now().Add(maxWait)

	for attempt := 1; time.Now().Before(deadline); attempt++ {
		doc, err := s.db.Queries().GetDocumentByFileID(ctx, sql.NullString{String: fileID, Valid: true})
		if err != nil {
			xlog.Warn("⚠️  [REAL STREAM] Failed to get document for file", "file_id", fileID, "error", err)
			return nil
		}

		// Check if document has meaningful content (VL description marker)
		if doc.Content.Valid && doc.Content.String != "" {
			content := doc.Content.String
			hasVLDescription := strings.Contains(content, "Image Description (AI Generated)") ||
				strings.Contains(content, "Video Description (AI Generated)")

			// If VL description is present, or content is substantial (>100 chars), use it
			if hasVLDescription || len(content) > 100 {
				filename := fileID
				if doc.Filename.Valid && doc.Filename.String != "" {
					filename = doc.Filename.String
				}
				xlog.Info("📄 [REAL STREAM] Got document content", "file", filename, "chars", len(content), "attempt", attempt)
				return &ChatFileChunk{
					ID:         doc.ID,
					FileID:     fileID,
					Filename:   filename,
					Text:       content,
					Similarity: 1.0,
				}
			}
		}

		// Check if we should continue polling
		if time.Now().Add(pollInterval).After(deadline) {
			break
		}

		xlog.Debug("⏳ [REAL STREAM] Waiting for VL description", "file_id", fileID, "attempt", attempt)
		select {
		case <-ctx.Done():
			xlog.Warn("⚠️  [REAL STREAM] Context cancelled while waiting for file", "file_id", fileID)
			return nil
		case <-time.After(pollInterval):
			// Continue polling
		}
	}

	xlog.Warn("⏰ [REAL STREAM] Timeout waiting for VL description", "file_id", fileID)
	return nil
}
