/*
 * Copyright 2025 Veridium Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/kawai-network/veridium/internal/database"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/internal/llama"
	"github.com/wailsapp/wails/v3/pkg/application"
)

// AgentChatService provides agent-based chat with RAG, tools, and context awareness
// This replaces/complements the existing LibraryChatService with Eino-based capabilities
type AgentChatService struct {
	app         *application.App
	db          *database.Service
	libService  *llama.LibraryService
	llamaModel  *llama.LlamaEinoModel
	kbService   *KnowledgeBaseService
	ragWorkflow *RAGWorkflow

	// Phase 3: Integration bridges
	toolsBridge   *ToolsEngineBridge   // Bridge to existing tools engine
	contextBridge *ContextEngineBridge // Bridge to existing context engine

	// Phase 4: Thread management
	threadService *ThreadManagementService // Thread management service

	// Reasoning mode configuration
	reasoningConfig ReasoningConfig // Controls reasoning behavior

	// Utility models: use separate small & fast models for efficiency
	titleModelPath   string // Path to title generation model (optional, falls back to main model)
	summaryModelPath string // Path to summary generation model (optional, falls back to title/main model)

	// Session management (hybrid: DB + in-memory cache)
	sessions      map[string]*AgentSession // In-memory cache for active sessions
	sessionsMutex sync.RWMutex
}

// AgentSession represents an ongoing conversation with context
type AgentSession struct {
	SessionID       string
	UserID          string
	Agent           adk.Agent
	Messages        []*schema.Message
	KnowledgeBaseID string
	Tools           []tool.BaseTool
	Context         map[string]any
	CreatedAt       int64
	UpdatedAt       int64
	DBSession       *db.Session // Link to DB session

	// Phase 4: Topic & Thread support
	TopicID  string // Current topic ID (may be auto-created)
	ThreadID string // Current thread ID (for branching)
}

// ChatRequest represents a chat request with agent capabilities
type ChatRequest struct {
	// Identity
	SessionID string `json:"session_id"`
	UserID    string `json:"user_id"`

	// Message content
	Message string `json:"message"`

	// Context (Phase 4: Topic & Thread support)
	TopicID  string `json:"topic_id,omitempty"`  // Current topic (may be auto-created)
	ThreadID string `json:"thread_id,omitempty"` // Current thread (for branching)
	ParentID string `json:"parent_id,omitempty"` // Parent message (for threading)

	// Configuration
	KnowledgeBaseID string         `json:"knowledge_base_id,omitempty"`
	Tools           []string       `json:"tools,omitempty"`
	Context         map[string]any `json:"context,omitempty"`
	Temperature     float32        `json:"temperature,omitempty"`
	MaxTokens       int            `json:"max_tokens,omitempty"`
	Stream          bool           `json:"stream,omitempty"`
}

// ChatResponse represents the response from agent
type ChatResponse struct {
	// IDs (Phase 4: Return created/used IDs)
	MessageID string `json:"message_id"`          // Created assistant message ID
	SessionID string `json:"session_id"`          // Session ID
	TopicID   string `json:"topic_id,omitempty"`  // Topic ID (may be auto-created)
	ThreadID  string `json:"thread_id,omitempty"` // Thread ID (if in thread)

	// Content
	Message      string             `json:"message"`
	ToolCalls    []schema.ToolCall  `json:"tool_calls,omitempty"`
	Sources      []*schema.Document `json:"sources,omitempty"`
	FinishReason string             `json:"finish_reason"`
	Usage        *schema.TokenUsage `json:"usage,omitempty"`

	// Metadata (Phase 4)
	CreatedAt int64  `json:"created_at"`      // Timestamp
	Error     string `json:"error,omitempty"` // Error if any
}

// NewAgentChatService creates a new agent-based chat service
// Set toolsBridge, contextBridge, and/or threadService to nil to disable those features
func NewAgentChatService(
	app *application.App,
	db *database.Service,
	libService *llama.LibraryService,
	kbService *KnowledgeBaseService,
	toolsBridge *ToolsEngineBridge,
	contextBridge *ContextEngineBridge,
	threadService *ThreadManagementService,
) *AgentChatService {
	llamaModel := llama.NewLlamaEinoModel(libService)
	ragWorkflow := NewRAGWorkflow(kbService)

	// Auto-detect utility models for lightweight tasks
	titleModelPath := detectTitleGenerationModel(libService)
	summaryModelPath := detectSummaryGenerationModel(libService)

	service := &AgentChatService{
		app:              app,
		db:               db,
		libService:       libService,
		llamaModel:       llamaModel,
		kbService:        kbService,
		ragWorkflow:      ragWorkflow,
		toolsBridge:      toolsBridge,
		contextBridge:    contextBridge,
		threadService:    threadService,
		reasoningConfig:  DefaultReasoningConfig(), // Default: disabled (non-reasoning)
		titleModelPath:   titleModelPath,
		summaryModelPath: summaryModelPath,
		sessions:         make(map[string]*AgentSession),
	}

	if titleModelPath != "" {
		log.Printf("📝 Auto-detected title model: %s", filepath.Base(titleModelPath))
	} else {
		log.Printf("📝 No dedicated title model found, will use main chat model")
	}

	if summaryModelPath != "" {
		log.Printf("📋 Auto-detected summary model: %s", filepath.Base(summaryModelPath))
	} else {
		log.Printf("📋 No dedicated summary model found, will use main chat model")
	}

	return service
}

// detectTitleGenerationModel finds the smallest, fastest model for title generation
// Prefers utility models (Llama 3.2 1B/3B) which are optimized for quick tasks
// Reuses detectSummaryGenerationModel since both tasks need the same type of model
func detectTitleGenerationModel(libService *llama.LibraryService) string {
	return detectSummaryGenerationModel(libService)
}

// detectSummaryGenerationModel finds optimal model for summarization
// Strategy: Prefer small non-reasoning models (Llama 3.2 1B/3B, Mistral)
// These models are fast, don't generate <think> tags, and have good quality for utility tasks
func detectSummaryGenerationModel(libService *llama.LibraryService) string {
	models, err := libService.GetAvailableModels()
	if err != nil || len(models) == 0 {
		return ""
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	modelsDir := filepath.Join(homeDir, ".llama-cpp", "models")

	var bestModel string
	var bestScore int

	for _, modelName := range models {
		modelPath := filepath.Join(modelsDir, modelName)
		nameLower := strings.ToLower(modelName)

		// Skip embedding models
		if strings.Contains(nameLower, "embedding") || strings.Contains(nameLower, "embed") {
			continue
		}

		// Get file size
		info, err := os.Stat(modelPath)
		if err != nil {
			continue
		}

		score := 0
		sizeMB := info.Size() / (1024 * 1024)

		// Prefer SMALL but not TOO small (1B-1.5B ideal for summary)
		// Summary needs slightly better quality than title
		if sizeMB >= 500 && sizeMB < 1000 {
			score += 100 // 1B models - BEST for summary (Llama 3.2 1B)
		} else if sizeMB >= 1000 && sizeMB < 2000 {
			score += 90 // 1.5B-3B models - good balance (Llama 3.2 3B)
		} else if sizeMB < 500 {
			score += 70 // 0.5B models - may be too simple
		} else if sizeMB < 4500 {
			score += 50 // 3B-7B models - slower, unnecessary (Mistral 7B)
		} else {
			score += 10 // 7B+ models - too slow for background task
		}

		// CRITICAL: Prefer non-reasoning models (NO <think> tags)
		if strings.Contains(nameLower, "llama-3.2-1b") || strings.Contains(nameLower, "llama_3.2_1b") {
			score += 100 // Llama 3.2 1B is BEST - fast, no think tags, good quality
		} else if strings.Contains(nameLower, "llama-3.2-3b") || strings.Contains(nameLower, "llama_3.2_3b") {
			score += 90 // Llama 3.2 3B is excellent - better quality
		} else if strings.Contains(nameLower, "llama") && !strings.Contains(nameLower, "3.2") {
			score += 70 // Other Llama models are okay
		} else if strings.Contains(nameLower, "mistral") {
			score += 60 // Mistral is decent
		} else if strings.Contains(nameLower, "gemma") {
			score += 50 // Gemma is okay
		} else if strings.Contains(nameLower, "phi") {
			score += 40 // Phi is acceptable
		} else if strings.Contains(nameLower, "qwen") {
			score -= 200 // AVOID Qwen - generates <think> tags in summary!
		} else if strings.Contains(nameLower, "deepseek") {
			score -= 150 // AVOID DeepSeek - reasoning model
		}

		// Prefer instruct/chat models
		if strings.Contains(nameLower, "instruct") || strings.Contains(nameLower, "chat") {
			score += 30
		}

		// Prefer Q4 quantization (good balance of quality & speed)
		if strings.Contains(nameLower, "q4") {
			score += 20
		}

		if score > bestScore {
			bestScore = score
			bestModel = modelPath
		}
	}

	return bestModel
}

// Chat processes a chat request and returns a response
func (s *AgentChatService) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// Get or create session
	session, err := s.getOrCreateSession(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create session: %w", err)
	}

	// ✅ LOAD SUMMARY FROM DATABASE IF EXISTS
	// Summary will be injected into system prompt when creating agent
	if req.TopicID != "" {
		topic, err := s.db.Queries().GetTopic(ctx, db.GetTopicParams{
			ID:     req.TopicID,
			UserID: req.UserID,
		})

		if err == nil && topic.HistorySummary.Valid && topic.HistorySummary.String != "" {
			// Store summary in session context for agent creation
			if session.Context == nil {
				session.Context = make(map[string]any)
			}
			session.Context["history_summary"] = topic.HistorySummary.String
			log.Printf("📋 Loaded history summary for topic %s (%d chars)",
				req.TopicID, len(topic.HistorySummary.String))
		}
	}

	// Phase 4: Handle topic and thread loading
	// If ThreadID is provided, load thread context
	if req.ThreadID != "" && s.threadService != nil {
		threadMessages, err := s.threadService.GetThreadMessages(ctx, req.ThreadID, req.UserID)
		if err != nil {
			log.Printf("⚠️  Warning: Failed to load thread messages: %v", err)
		} else {
			// Convert thread messages to Eino format
			einoMessages := make([]*schema.Message, 0, len(threadMessages))
			for _, dbMsg := range threadMessages {
				einoMsg, err := convertDBMessageToEino(&dbMsg)
				if err == nil {
					einoMessages = append(einoMessages, einoMsg)
				}
			}
			session.Messages = einoMessages
			session.ThreadID = req.ThreadID
			log.Printf("📋 Loaded %d messages from thread %s", len(einoMessages), req.ThreadID)
		}
	}

	// Validate model for current reasoning mode and auto-switch if needed
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

	// Set TopicID from request (may be auto-created later)
	currentTopicID := req.TopicID
	if currentTopicID != "" {
		session.TopicID = currentTopicID
	}

	// Add user message to session
	userMsg := &schema.Message{
		Role:    schema.User,
		Content: req.Message,
	}
	session.Messages = append(session.Messages, userMsg)

	// Phase 4: Auto-create topic if needed BEFORE saving first message
	// Only create if this is first message (session.Messages == 1 after adding user msg)
	if currentTopicID == "" && len(session.Messages) == 1 {
		// Create topic synchronously to get topicID for message saving
		topicID, err := s.createTopicForSessionSync(ctx, session.SessionID, session.UserID)
		if err != nil {
			log.Printf("⚠️  Warning: Failed to create topic: %v", err)
			// Continue without topic - will retry after response
		} else {
			currentTopicID = topicID
			session.TopicID = topicID
			log.Printf("📝 Auto-created topic: %s", topicID)
		}
	}

	// Save user message to DB with topic and thread IDs
	userMsgID, err := s.saveMessageToDB(ctx, userMsg, session.SessionID, session.UserID, currentTopicID, req.ThreadID)
	if err != nil {
		log.Printf("⚠️  Warning: Failed to save user message to DB: %v", err)
	} else {
		log.Printf("💾 Saved user message: %s (topic: %s, thread: %s)", userMsgID, currentTopicID, req.ThreadID)
	}

	// Phase 3: Process messages through context engine (if available)
	messagesToAgent := session.Messages
	if s.contextBridge != nil {
		processedMessages, err := s.contextBridge.ProcessMessagesForAgent(ctx, session.Messages)
		if err != nil {
			log.Printf("⚠️  Warning: Context engine processing failed: %v", err)
			// Continue with original messages
		} else {
			messagesToAgent = processedMessages
			log.Printf("🔄 Context engine processed messages: %d → %d",
				len(session.Messages), len(processedMessages))
		}
	}

	// Prepare agent input
	agentInput := &adk.AgentInput{
		Messages: messagesToAgent,
	}

	// Generate message ID for streaming
	assistantMsgID := uuid.New().String()

	// Run agent with streaming support
	var finalMessage string
	var toolCalls []schema.ToolCall
	var finishReason string
	var usage *schema.TokenUsage

	if req.Stream && s.app != nil {
		// Use token-by-token streaming via direct llama generation
		finalMessage, err = s.generateWithTokenStreaming(ctx, session, assistantMsgID, messagesToAgent)
		if err != nil {
			return nil, fmt.Errorf("streaming generation failed: %w", err)
		}
		finishReason = "stop"
	} else {
		// Use standard Eino agent (non-streaming)
		iterator := session.Agent.Run(ctx, agentInput)

		for {
			event, ok := iterator.Next()
			if !ok {
				break
			}

			if event.Err != nil {
				return nil, fmt.Errorf("agent execution failed: %w", event.Err)
			}

			if event.Output != nil && event.Output.MessageOutput != nil {
				msgVariant := event.Output.MessageOutput
				if msgVariant.Message != nil {
					if msgVariant.Role == schema.Assistant {
						finalMessage += msgVariant.Message.Content

						if len(msgVariant.Message.ToolCalls) > 0 {
							toolCalls = append(toolCalls, msgVariant.Message.ToolCalls...)
						}
						if msgVariant.Message.ResponseMeta != nil {
							finishReason = msgVariant.Message.ResponseMeta.FinishReason
							usage = msgVariant.Message.ResponseMeta.Usage
						}
					}
				}
			}
		}
	}

	// Note: Think tag stripping removed - proper model selection should prevent think tags
	// If reasoning mode is disabled, use non-reasoning models (Llama 3.2)
	// If reasoning mode is enabled, use reasoning models with /no_think (Qwen3)

	// Add assistant response to session
	assistantMsg := &schema.Message{
		Role:      schema.Assistant,
		Content:   finalMessage,
		ToolCalls: toolCalls,
		ResponseMeta: &schema.ResponseMeta{
			FinishReason: finishReason,
			Usage:        usage,
		},
	}
	session.Messages = append(session.Messages, assistantMsg)

	// Monitor context usage and provide warnings/recommendations
	if usage != nil && usage.TotalTokens > 0 {
		contextSize := 16384 // Current context window size
		contextUsage := float64(usage.TotalTokens) / float64(contextSize)
		turnCount := len(session.Messages) / 2 // Approximate turn count (user + assistant pairs)

		log.Printf("📊 Context usage: %d/%d tokens (%.1f%%) after %d turns",
			usage.TotalTokens, contextSize, contextUsage*100, turnCount)

		// Warning at 50% usage
		if contextUsage > 0.5 && contextUsage <= 0.8 {
			log.Printf("⚠️  Context usage > 50%%. Consider switching to more efficient reasoning mode.")
			if s.reasoningConfig.Mode != ReasoningDisabled {
				log.Printf("💡 Recommendation: Switch to ReasoningDisabled (Llama 3.2) for longer conversations")
			}
		}

		// Critical warning at 80% usage
		if contextUsage > 0.8 {
			log.Printf("🚨 Context usage > 80%% Conversation may become unstable.")
			log.Printf("💡 Recommendations:")
			log.Printf("   1. Switch to ReasoningDisabled mode")
			log.Printf("   2. Summarize and truncate old messages")
			log.Printf("   3. Start a new conversation")
		}
	}

	// Phase 4: Generate topic title after first response
	// If we already created a placeholder topic, update it with LLM-generated title
	// We allow up to 4 messages (2 turns) to retry title generation if it failed or was skipped
	if len(session.Messages) >= 2 && len(session.Messages) <= 4 {
		log.Printf("📝 Checking if topic title needs update (msgs: %d, topic: %s)", len(session.Messages), currentTopicID)

		if currentTopicID != "" {
			// Update existing topic with LLM-generated title
			// This runs in background, so it won't block response
			err := s.updateTopicTitle(ctx, currentTopicID, session.UserID, session.Messages)
			if err != nil {
				log.Printf("⚠️  Warning: Failed to trigger topic title update: %v", err)
			} else {
				log.Printf("📝 Triggered background title update for topic %s", currentTopicID)
			}
		} else {
			// Topic was not created before - create it now with the actual conversation
			topicID, err := s.createTopicForSessionWithTitle(ctx, session.SessionID, session.UserID, session.Messages)
			if err != nil {
				log.Printf("⚠️  Warning: Failed to create topic with title: %v", err)
			} else {
				currentTopicID = topicID
				session.TopicID = topicID
				log.Printf("📝 Created topic with auto-generated title: %s", topicID)
			}
		}
	}

	// Save assistant message to DB with topic and thread IDs
	savedMsgID, err := s.saveMessageToDB(ctx, assistantMsg, session.SessionID, session.UserID, currentTopicID, req.ThreadID)
	if err != nil {
		log.Printf("⚠️  Warning: Failed to save assistant message to DB: %v", err)
		// Use the generated ID as fallback
		savedMsgID = assistantMsgID
	} else {
		log.Printf("💾 Saved assistant message: %s (topic: %s, thread: %s)", savedMsgID, currentTopicID, req.ThreadID)
		// Update assistantMsgID with the actual DB ID
		assistantMsgID = savedMsgID
	}

	// ✅ AUTO-TRIGGER SUMMARY (NON-BLOCKING)
	// This happens AFTER response is sent to user, so no blocking
	// Summary generation runs in background goroutine
	if currentTopicID != "" {
		go func() {
			bgCtx := context.Background()
			// First: Try initial summary creation (if no summary exists)
			s.autoSummarizeIfNeeded(bgCtx, session, currentTopicID, session.UserID)

			// Second: Try incremental summary update (if summary exists and enough new messages)
			s.incrementalSummarizeIfNeeded(bgCtx, session, currentTopicID, session.UserID)
		}()
	}

	// Emit streaming complete event (if streaming enabled)
	if req.Stream && s.app != nil {
		s.app.Event.Emit("chat:stream", map[string]interface{}{
			"type":       "complete",
			"session_id": session.SessionID,
			"message_id": assistantMsgID,
			"content":    finalMessage,
			"topic_id":   currentTopicID,
			"thread_id":  req.ThreadID,
		})
	}

	// Update session timestamp in DB
	if err := s.updateSessionTimestamp(ctx, session.SessionID, session.UserID); err != nil {
		log.Printf("⚠️  Warning: Failed to update session timestamp: %v", err)
	}

	// Get sources if KB was used
	var sources []*schema.Document
	if req.KnowledgeBaseID != "" {
		// Query KB to get sources that were used
		sources, _ = s.kbService.QueryKnowledgeBase(ctx, req.KnowledgeBaseID, req.Message, 3, req.UserID)
	}

	// Phase 4: Return with all relevant IDs
	now := time.Now().UnixMilli()
	return &ChatResponse{
		MessageID:    assistantMsgID,
		SessionID:    session.SessionID,
		TopicID:      currentTopicID,
		ThreadID:     req.ThreadID,
		Message:      finalMessage,
		ToolCalls:    toolCalls,
		Sources:      sources,
		FinishReason: finishReason,
		Usage:        usage,
		CreatedAt:    now,
	}, nil
}

// getOrCreateSession gets an existing session or creates a new one with DB persistence
func (s *AgentChatService) getOrCreateSession(ctx context.Context, req ChatRequest) (*AgentSession, error) {
	s.sessionsMutex.Lock()
	defer s.sessionsMutex.Unlock()

	// 1. Check in-memory cache first
	if session, exists := s.sessions[req.SessionID]; exists {
		log.Printf("♻️  Reusing cached session: %s", req.SessionID)
		return session, nil
	}

	// 2. Try to load from database
	dbSession, dbMessages, err := s.loadSessionFromDB(ctx, req.SessionID, req.UserID)
	if err == nil {
		// Session exists in DB - reconstruct from history
		log.Printf("📂 Loading session from DB: %s (%d messages)", req.SessionID, len(dbMessages))

		// Convert DB messages to Eino messages
		einoMessages := make([]*schema.Message, 0, len(dbMessages))
		for _, dbMsg := range dbMessages {
			einoMsg, err := convertDBMessageToEino(&dbMsg)
			if err == nil {
				einoMessages = append(einoMessages, einoMsg)
			}
		}

		// Create session with history
		session := &AgentSession{
			SessionID:       req.SessionID,
			UserID:          req.UserID,
			Messages:        einoMessages,
			KnowledgeBaseID: req.KnowledgeBaseID,
			Context:         req.Context,
			CreatedAt:       dbSession.CreatedAt,
			UpdatedAt:       dbSession.UpdatedAt,
			DBSession:       dbSession,
		}

		// Collect tools
		var tools []tool.BaseTool
		if req.KnowledgeBaseID != "" {
			if kbTool, err := s.createKBSearchTool(ctx, req.KnowledgeBaseID, req.UserID); err == nil {
				tools = append(tools, kbTool)
			}
		}

		// Phase 3: Add tools from tools engine bridge
		if s.toolsBridge != nil && len(req.Tools) > 0 {
			bridgeTools := s.toolsBridge.GetToolsForAgent(req.Tools)
			tools = append(tools, bridgeTools...)
			log.Printf("🔧 Reconstructed session: Added %d tools from tools engine", len(bridgeTools))
		}

		session.Tools = tools

		// Create agent with history
		agent, err := s.createAgent(ctx, session)
		if err != nil {
			return nil, fmt.Errorf("failed to create agent: %w", err)
		}

		session.Agent = agent
		s.sessions[req.SessionID] = session

		log.Printf("✅ Reconstructed session from DB: %s (KB: %s, History: %d msgs)",
			req.SessionID, req.KnowledgeBaseID, len(einoMessages))

		return session, nil
	}

	// 3. Session doesn't exist - create new one
	log.Printf("🆕 Creating new session: %s", req.SessionID)

	// Create session in DB (handle race condition with UNIQUE constraint)
	dbSession, err = s.createSessionInDB(ctx, req.SessionID, req.UserID, req.KnowledgeBaseID)
	if err != nil {
		// Check if error is due to UNIQUE constraint (race condition)
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			// Another process created the session - retry loading with exponential backoff
			log.Printf("⚠️  Race condition detected, retrying load from DB...")

			var dbMessages []db.Message
			maxRetries := 3
			for i := 0; i < maxRetries; i++ {
				// Wait a bit for the other process to commit
				time.Sleep(time.Duration(10*(i+1)) * time.Millisecond)

				dbSession, dbMessages, err = s.loadSessionFromDB(ctx, req.SessionID, req.UserID)
				if err == nil {
					log.Printf("✅ Successfully loaded session after %d retries", i+1)
					break
				}

				if i < maxRetries-1 {
					log.Printf("⏳ Retry %d/%d: session not yet available, waiting...", i+1, maxRetries)
				}
			}

			if err != nil {
				return nil, fmt.Errorf("failed to load session after race condition (tried %d times): %w", maxRetries, err)
			}

			// Convert DB messages to Eino messages
			einoMessages := make([]*schema.Message, 0, len(dbMessages))
			for _, dbMsg := range dbMessages {
				einoMsg, err := convertDBMessageToEino(&dbMsg)
				if err == nil {
					einoMessages = append(einoMessages, einoMsg)
				}
			}

			// Create session with history
			session := &AgentSession{
				SessionID:       req.SessionID,
				UserID:          req.UserID,
				Messages:        einoMessages,
				KnowledgeBaseID: req.KnowledgeBaseID,
				Context:         req.Context,
				CreatedAt:       dbSession.CreatedAt,
				UpdatedAt:       dbSession.UpdatedAt,
				DBSession:       dbSession,
			}

			// Collect tools
			var tools []tool.BaseTool
			if req.KnowledgeBaseID != "" {
				if kbTool, err := s.createKBSearchTool(ctx, req.KnowledgeBaseID, req.UserID); err == nil {
					tools = append(tools, kbTool)
				}
			}

			// Phase 3: Add tools from tools engine bridge
			if s.toolsBridge != nil && len(req.Tools) > 0 {
				bridgeTools := s.toolsBridge.GetToolsForAgent(req.Tools)
				tools = append(tools, bridgeTools...)
				log.Printf("🔧 Added %d tools from tools engine after race recovery", len(bridgeTools))
			}

			session.Tools = tools

			// Create agent with history
			agent, err := s.createAgent(ctx, session)
			if err != nil {
				return nil, fmt.Errorf("failed to create agent after race recovery: %w", err)
			}

			session.Agent = agent
			s.sessions[req.SessionID] = session

			log.Printf("✅ Recovered from race condition: %s", req.SessionID)
			return session, nil
		}

		// Other error - fail
		return nil, fmt.Errorf("failed to create session in DB: %w", err)
	}

	// Create in-memory session
	session := &AgentSession{
		SessionID:       req.SessionID,
		UserID:          req.UserID,
		Messages:        make([]*schema.Message, 0),
		KnowledgeBaseID: req.KnowledgeBaseID,
		Context:         req.Context,
		CreatedAt:       dbSession.CreatedAt,
		UpdatedAt:       dbSession.UpdatedAt,
		DBSession:       dbSession,
	}

	// Collect tools for the agent
	var tools []tool.BaseTool

	// Add KB search tool if KB is specified
	if req.KnowledgeBaseID != "" {
		kbTool, err := s.createKBSearchTool(ctx, req.KnowledgeBaseID, req.UserID)
		if err != nil {
			log.Printf("⚠️  Warning: Failed to create KB tool: %v", err)
		} else {
			tools = append(tools, kbTool)
		}
	}

	// Phase 3: Add tools from tools engine bridge
	if s.toolsBridge != nil && len(req.Tools) > 0 {
		bridgeTools := s.toolsBridge.GetToolsForAgent(req.Tools)
		tools = append(tools, bridgeTools...)
		log.Printf("🔧 Added %d tools from tools engine: %v", len(bridgeTools), req.Tools)
	}

	session.Tools = tools

	// Create agent for this session
	agent, err := s.createAgent(ctx, session)
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	session.Agent = agent
	s.sessions[req.SessionID] = session

	log.Printf("✅ Created new agent session: %s (KB: %s, Tools: %d)",
		req.SessionID, req.KnowledgeBaseID, len(tools))

	return session, nil
}

// createAgent creates a new ADK agent for the session
func (s *AgentChatService) createAgent(ctx context.Context, session *AgentSession) (adk.Agent, error) {
	// Build base instruction
	baseInstruction := "You are a helpful AI assistant. "

	// ✅ INJECT SUMMARY IF EXISTS
	if historySummary, ok := session.Context["history_summary"].(string); ok && historySummary != "" {
		summaryContext := fmt.Sprintf(`

<chat_history_summary>
<docstring>Previous conversation summary (older messages have been compressed):</docstring>
<summary>%s</summary>
</chat_history_summary>

`, historySummary)
		baseInstruction += summaryContext
		log.Printf("📋 Injected history summary into system prompt (%d chars)", len(historySummary))
	}

	if session.KnowledgeBaseID != "" {
		baseInstruction += "You have access to a knowledge base. Use the search tool to find relevant information before answering questions. "
	}

	// Apply reasoning mode to instruction
	instruction := s.reasoningConfig.GetSystemPrompt(baseInstruction)

	log.Printf("🧠 Creating agent with reasoning mode: %s", s.reasoningConfig.Mode)
	if log.Default().Writer() != nil {
		// Only log full prompt in debug mode (can be very long with summary)
		log.Printf("📝 System prompt length: %d chars", len(instruction))
	}

	// Create agent config
	config := &adk.ChatModelAgentConfig{
		Name:        fmt.Sprintf("agent_%s", session.SessionID),
		Description: "AI assistant with RAG and tool capabilities",
		Instruction: instruction,
		Model:       s.llamaModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: session.Tools,
			},
		},
		MaxIterations: 10,
	}

	return adk.NewChatModelAgent(ctx, config)
}

// createKBSearchTool creates a knowledge base search tool
func (s *AgentChatService) createKBSearchTool(ctx context.Context, kbID, userID string) (tool.BaseTool, error) {
	kb, err := s.kbService.GetKnowledgeBase(ctx, kbID, userID)
	if err != nil {
		return nil, err
	}

	return &kbSearchTool{
		name:        fmt.Sprintf("search_%s", kb.Name),
		description: fmt.Sprintf("Search the %s knowledge base for relevant information", kb.Name),
		kbService:   s.kbService,
		kbID:        kbID,
		userID:      userID,
	}, nil
}

// ClearSession removes a session from memory
func (s *AgentChatService) ClearSession(sessionID string) {
	s.sessionsMutex.Lock()
	defer s.sessionsMutex.Unlock()

	delete(s.sessions, sessionID)
	log.Printf("🗑️  Cleared session: %s", sessionID)
}

// GetSessionHistory returns the message history for a session
func (s *AgentChatService) GetSessionHistory(sessionID string) ([]*schema.Message, error) {
	s.sessionsMutex.RLock()
	defer s.sessionsMutex.RUnlock()

	session, exists := s.sessions[sessionID]
	if !exists {
		return nil, fmt.Errorf("session not found: %s", sessionID)
	}

	return session.Messages, nil
}

// SetTitleModel sets a specific model for title generation
// Use a small, fast model (e.g., Llama 3.2 1B) for efficiency
// If not set, falls back to main chat model
func (s *AgentChatService) SetTitleModel(modelPath string) {
	s.sessionsMutex.Lock()
	defer s.sessionsMutex.Unlock()
	s.titleModelPath = modelPath
	log.Printf("📝 Title generation model set to: %s", modelPath)
}

// generateTopicTitle generates a concise title for the conversation using LLM
// ALWAYS uses non-reasoning model to avoid <think> tags - no stripping needed!
// Uses separate title model if configured, otherwise falls back to main chat model
func (s *AgentChatService) generateTopicTitle(ctx context.Context, messages []*schema.Message, locale string) (string, error) {
	if len(messages) == 0 {
		return "New Conversation", nil
	}

	// Build summary prompt - simpler now since non-reasoning models don't need warnings about <think> tags
	systemPrompt := fmt.Sprintf(`You are a professional conversation summarizer. Generate a concise title that captures the essence of the conversation.

Rules:
- Maximum 10 words
- Maximum 50 characters
- No punctuation marks, quotes, or special characters
- Use the language specified by the locale code: %s
- Output ONLY the title text, nothing else

Example: Sleep Functions for Body and Mind`, locale)

	// Build conversation text (User messages only)
	var conversationText string
	for _, msg := range messages {
		if msg.Role == schema.User {
			conversationText += fmt.Sprintf("user: %s\n", msg.Content)
		}
	}

	// Fallback: if no user messages found (unlikely but possible), use all messages
	if conversationText == "" {
		for _, msg := range messages {
			conversationText += fmt.Sprintf("%s: %s\n", msg.Role, msg.Content)
		}
	}

	log.Printf("📝 Generating title for conversation (%d messages, %d chars)", len(messages), len(conversationText))

	// Create messages for title generation
	titleMessages := []*schema.Message{
		{Role: schema.System, Content: systemPrompt},
		{Role: schema.User, Content: conversationText},
	}

	// Use separate title model if configured (should be non-reasoning model)
	var response *schema.Message
	var err error

	if s.titleModelPath != "" {
		// Temporarily load title model for generation
		log.Printf("📝 Using dedicated non-reasoning title model: %s", filepath.Base(s.titleModelPath))

		// Save current model state
		currentModel := s.libService.GetLoadedChatModel()

		// Load title model (should be non-reasoning: Llama/Mistral)
		if err := s.libService.LoadChatModel(s.titleModelPath); err != nil {
			log.Printf("⚠️  Failed to load title model, falling back to main model: %v", err)
			// Fallback to main model (may be reasoning model, so strip <think> tags)
			response, err = s.llamaModel.Generate(ctx, titleMessages)
			if err == nil && strings.Contains(response.Content, "<think>") {
				log.Printf("⚠️  Main model generated <think> tags in title - this shouldn't happen with proper model selection")
				response.Content = stripThinkTags(response.Content)
			}
		} else {
			// Generate with title model (non-reasoning, no <think> tags expected)
			titleModel := llama.NewLlamaEinoModel(s.libService)
			response, err = titleModel.Generate(ctx, titleMessages)

			// Sanity check: if title model somehow generates <think> tags, warn and strip
			if err == nil && strings.Contains(response.Content, "<think>") {
				log.Printf("⚠️  WARNING: Non-reasoning title model generated <think> tags! Model: %s", filepath.Base(s.titleModelPath))
				response.Content = stripThinkTags(response.Content)
			}

			// Restore main chat model
			if currentModel != "" {
				if restoreErr := s.libService.LoadChatModel(currentModel); restoreErr != nil {
					log.Printf("⚠️  Failed to restore main model: %v", restoreErr)
				}
			}
		}
	} else {
		// No dedicated title model - use main chat model
		// This may be a reasoning model, so check for <think> tags
		log.Printf("📝 Using main chat model for title generation (may be reasoning model)")
		response, err = s.llamaModel.Generate(ctx, titleMessages)
		if err == nil && strings.Contains(response.Content, "<think>") {
			log.Printf("⚠️  Main model generated <think> tags in title (expected if using reasoning model)")
			response.Content = stripThinkTags(response.Content)
		}
	}

	if err != nil {
		return "", fmt.Errorf("failed to generate title: %w", err)
	}

	log.Printf("📝 Raw title response: %q", response.Content)

	// Clean up the title
	title := strings.TrimSpace(response.Content)

	// If title is empty, try to extract from the original response
	if title == "" {
		// Try to find text in quotes
		quotePattern := regexp.MustCompile(`["']([^"']+)["']`)
		matches := quotePattern.FindStringSubmatch(response.Content)
		if len(matches) > 1 {
			title = matches[1]
		} else {
			// Fallback: use first non-empty line
			lines := strings.Split(response.Content, "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					title = line
					break
				}
			}
		}
	}

	// Final fallback if still empty
	if title == "" {
		log.Printf("⚠️  Title generation failed (empty result), using default")
		title = "New Conversation"
	}

	// Truncate to 50 characters
	if len(title) > 50 {
		title = title[:50]
	}

	return title, nil
}

// stripThinkTags removes <think>...</think> blocks from the text using regex
func stripThinkTags(text string) string {
	// First, try to remove complete <think>...</think> blocks
	re := regexp.MustCompile(`(?s)<think>.*?</think>`)
	cleaned := re.ReplaceAllString(text, "")

	// If <think> tag still exists (incomplete/unclosed), remove everything from <think> onwards
	if strings.Contains(cleaned, "<think>") {
		idx := strings.Index(cleaned, "<think>")
		cleaned = cleaned[:idx]
	}

	return strings.TrimSpace(cleaned)
}

// createTopicForSessionSync creates a topic synchronously and returns the topicID
// Phase 4: Used to create topic before first message so we have topicID for message saving
func (s *AgentChatService) createTopicForSessionSync(ctx context.Context, sessionID, userID string) (string, error) {
	// Check if topic already exists for this session
	count, err := s.db.Queries().CountTopicsBySession(ctx, db.CountTopicsBySessionParams{
		SessionID: sql.NullString{String: sessionID, Valid: true},
		UserID:    userID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to check existing topics: %w", err)
	}

	if count > 0 {
		// Topic already exists - find and return it
		// TODO: Add GetTopicBySession query to SQLC if needed
		// For now, return empty string and let it be created with title later
		return "", fmt.Errorf("topic already exists for session")
	}

	// Create topic with placeholder title
	topicID := uuid.New().String()
	now := time.Now().UnixMilli()

	// Use sessionID as is (don't convert "inbox" to NULL)
	sessionIDForDB := sql.NullString{String: sessionID, Valid: true}

	_, err = s.db.Queries().CreateTopic(ctx, db.CreateTopicParams{
		ID:             topicID,
		Title:          sql.NullString{String: "New Conversation", Valid: true},
		Favorite:       0,
		SessionID:      sessionIDForDB,
		GroupID:        sql.NullString{},
		UserID:         userID,
		ClientID:       sql.NullString{},
		HistorySummary: sql.NullString{},
		Metadata:       sql.NullString{},
		CreatedAt:      now,
		UpdatedAt:      now,
	})

	if err != nil {
		return "", fmt.Errorf("failed to create topic: %w", err)
	}

	return topicID, nil
}

// updateTopicTitle updates an existing topic with LLM-generated title
// Phase 4 FIX: Used to update placeholder topic with meaningful title after first response
func (s *AgentChatService) updateTopicTitle(ctx context.Context, topicID, userID string, messages []*schema.Message) error {
	// Create a copy of messages to avoid race conditions
	messagesCopy := make([]*schema.Message, len(messages))
	copy(messagesCopy, messages)

	// Run in background
	go func() {
		// Add a small delay to ensure main chat request finishes and releases model resources
		time.Sleep(2 * time.Second)

		log.Printf("🔄 Generating title in background for topic %s...", topicID)

		// Create context with timeout
		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer cancel()

		// Generate title (default locale: en-US)
		title, err := s.generateTopicTitle(ctx, messagesCopy, "en-US")
		if err != nil {
			log.Printf("⚠️  Warning: Failed to generate topic title: %v", err)
			title = "New Conversation"
		}

		// Fetch existing topic first to preserve history_summary and metadata
		existingTopic, err := s.db.Queries().GetTopic(ctx, db.GetTopicParams{
			ID:     topicID,
			UserID: userID,
		})
		if err != nil {
			log.Printf("⚠️  Failed to fetch topic for update: %v", err)
			return
		}

		// Update topic title in database
		now := time.Now().UnixMilli()
		_, err = s.db.Queries().UpdateTopic(ctx, db.UpdateTopicParams{
			Title:          sql.NullString{String: title, Valid: true},
			HistorySummary: existingTopic.HistorySummary, // Preserve existing
			Metadata:       existingTopic.Metadata,       // Preserve existing
			UpdatedAt:      now,
			ID:             topicID,
			UserID:         userID,
		})

		if err != nil {
			log.Printf("⚠️  Failed to update topic title in DB: %v", err)
		} else {
			log.Printf("✅ Updated topic %s with title: %s", topicID, title)

			// Emit event to notify UI
			if s.app != nil {
				s.app.Event.Emit("chat:topic:updated", map[string]interface{}{
					"topic_id": topicID,
					"title":    title,
				})
			}
		}
	}()

	// Return immediately - title will be updated in background
	return nil
}

// createTopicForSessionWithTitle creates a topic with LLM-generated title and returns topicID
// Phase 4: Used after first response to create/update topic with meaningful title
func (s *AgentChatService) createTopicForSessionWithTitle(ctx context.Context, sessionID, userID string, messages []*schema.Message) (string, error) {
	// Check if topic already exists for this session
	count, err := s.db.Queries().CountTopicsBySession(ctx, db.CountTopicsBySessionParams{
		SessionID: sql.NullString{String: sessionID, Valid: true},
		UserID:    userID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to check existing topics: %w", err)
	}

	if count > 0 {
		// Topic already exists - could update title here if needed
		// For now, just return success
		return "", nil
	}

	// Create topic in database with placeholder title first (non-blocking)
	topicID := uuid.New().String()
	now := time.Now().UnixMilli()

	// Use sessionID as is (don't convert "inbox" to NULL)
	sessionIDForDB := sql.NullString{String: sessionID, Valid: true}

	_, err = s.db.Queries().CreateTopic(ctx, db.CreateTopicParams{
		ID:             topicID,
		Title:          sql.NullString{String: "New Conversation", Valid: true}, // Placeholder
		Favorite:       0,
		SessionID:      sessionIDForDB,
		GroupID:        sql.NullString{},
		UserID:         userID,
		ClientID:       sql.NullString{},
		HistorySummary: sql.NullString{},
		Metadata:       sql.NullString{},
		CreatedAt:      now,
		UpdatedAt:      now,
	})

	if err != nil {
		return "", fmt.Errorf("failed to create topic: %w", err)
	}

	// Generate title in BACKGROUND to avoid blocking (model loading can take 2-5 seconds)
	messagesCopy := make([]*schema.Message, len(messages))
	copy(messagesCopy, messages)

	topicIDCopy := topicID
	userIDCopy := userID

	go func() {
		log.Printf("🔄 Generating title in background for new topic %s...", topicIDCopy)

		// Generate title (default locale: en-US)
		title, err := s.generateTopicTitle(context.Background(), messagesCopy, "en-US")
		if err != nil {
			log.Printf("⚠️  Warning: Failed to generate topic title: %v", err)
			title = "New Conversation"
		}

		// Update topic title in database
		now := time.Now().UnixMilli()
		_, err = s.db.Queries().UpdateTopic(context.Background(), db.UpdateTopicParams{
			Title:          sql.NullString{String: title, Valid: true},
			HistorySummary: sql.NullString{}, // Keep existing
			Metadata:       sql.NullString{}, // Keep existing
			UpdatedAt:      now,
			ID:             topicIDCopy,
			UserID:         userIDCopy,
		})

		if err != nil {
			log.Printf("⚠️  Failed to update topic title in DB: %v", err)
		} else {
			log.Printf("✅ Updated new topic %s with title: %s", topicIDCopy, title)

			// Emit event to notify UI
			if s.app != nil {
				s.app.Event.Emit("chat:topic:updated", map[string]interface{}{
					"topic_id": topicIDCopy,
					"title":    title,
				})
			}
		}
	}()

	return topicID, nil
}

// kbSearchTool implements tool.BaseTool for knowledge base search
type kbSearchTool struct {
	name        string
	description string
	kbService   *KnowledgeBaseService
	kbID        string
	userID      string
}

func (t *kbSearchTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	params := map[string]*schema.ParameterInfo{
		"query": {
			Type: "string",
			Desc: "The search query to find relevant information",
		},
		"top_k": {
			Type: "integer",
			Desc: "Number of results to return (default: 5)",
		},
	}

	return &schema.ToolInfo{
		Name:        t.name,
		Desc:        t.description,
		ParamsOneOf: schema.NewParamsOneOfByParams(params),
	}, nil
}

func (t *kbSearchTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// Parse arguments
	var args struct {
		Query string `json:"query"`
		TopK  int    `json:"top_k"`
	}

	if err := json.Unmarshal([]byte(argumentsInJSON), &args); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	if args.TopK <= 0 {
		args.TopK = 5
	}

	// Query knowledge base
	docs, err := t.kbService.QueryKnowledgeBase(ctx, t.kbID, args.Query, args.TopK, t.userID)
	if err != nil {
		return "", fmt.Errorf("KB search failed: %w", err)
	}

	// Format results
	result := map[string]any{
		"found": len(docs),
		"results": func() []map[string]any {
			results := make([]map[string]any, len(docs))
			for i, doc := range docs {
				results[i] = map[string]any{
					"content":  doc.Content,
					"metadata": doc.MetaData,
				}
			}
			return results
		}(),
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	return string(resultJSON), nil
}

// ============================================================================
// Database Persistence Methods
// ============================================================================

// convertDBMessageToEino converts a database message to Eino schema message
func convertDBMessageToEino(dbMsg *db.Message) (*schema.Message, error) {
	msg := &schema.Message{
		Role: schema.RoleType(dbMsg.Role),
	}

	// Content
	if dbMsg.Content.Valid {
		msg.Content = dbMsg.Content.String
	}

	// Tool calls (stored as JSON in Tools field)
	if dbMsg.Tools.Valid && dbMsg.Tools.String != "" {
		var toolCalls []schema.ToolCall
		if err := json.Unmarshal([]byte(dbMsg.Tools.String), &toolCalls); err == nil {
			msg.ToolCalls = toolCalls
		}
	}

	// Note: Eino schema.Message doesn't have ResponseTo field
	// Parent relationship is handled at the agent level

	return msg, nil
}

// convertEinoMessageToDB converts Eino schema message to DB message params
func convertEinoMessageToDB(einoMsg *schema.Message, sessionID, userID string) db.CreateMessageParams {
	now := time.Now().UnixMilli()

	params := db.CreateMessageParams{
		ID:        uuid.New().String(),
		Role:      string(einoMsg.Role),
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Session ID
	if sessionID != "" {
		params.SessionID = sql.NullString{String: sessionID, Valid: true}
	}

	// Content
	if einoMsg.Content != "" {
		params.Content = sql.NullString{String: einoMsg.Content, Valid: true}
	}

	// Tool calls
	if len(einoMsg.ToolCalls) > 0 {
		if toolCallsJSON, err := json.Marshal(einoMsg.ToolCalls); err == nil {
			params.Tools = sql.NullString{String: string(toolCallsJSON), Valid: true}
		}
	}

	// Note: Parent ID handling can be added later if needed

	return params
}

// loadSessionFromDB loads a session from database
func (s *AgentChatService) loadSessionFromDB(ctx context.Context, sessionID, userID string) (*db.Session, []db.Message, error) {
	// Load session metadata (sessionID can be either ID or slug)
	dbSession, err := s.db.Queries().GetSessionByIdOrSlug(ctx, db.GetSessionByIdOrSlugParams{
		ID:     sessionID,
		Slug:   sessionID,
		UserID: userID,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("session not found in DB: %w", err)
	}

	// Load message history (use actual session ID from DB, not slug)
	dbMessages, err := s.db.Queries().ListMessagesBySession(ctx, db.ListMessagesBySessionParams{
		UserID:    userID,
		SessionID: sql.NullString{String: dbSession.ID, Valid: true},
		Limit:     1000, // Load up to 1000 messages
		Offset:    0,
	})
	if err != nil {
		// If messages query fails, still return session
		log.Printf("Warning: Failed to load messages for session %s: %v", sessionID, err)
		return &dbSession, nil, nil
	}

	return &dbSession, dbMessages, nil
}

// createSessionInDB creates a new session in database
func (s *AgentChatService) createSessionInDB(ctx context.Context, sessionID, userID, kbID string) (*db.Session, error) {
	now := time.Now().UnixMilli()

	// Generate slug from sessionID (take first 8 chars)
	slug := sessionID
	if len(slug) > 8 {
		slug = slug[:8]
	}

	params := db.CreateSessionParams{
		ID:        sessionID,
		Slug:      slug,
		UserID:    userID,
		Type:      sql.NullString{String: "agent", Valid: true},
		Pinned:    0,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// Set title based on KB if available
	if kbID != "" {
		kb, err := s.kbService.GetKnowledgeBase(ctx, kbID, userID)
		if err == nil {
			params.Title = sql.NullString{
				String: fmt.Sprintf("Chat with %s", kb.Name),
				Valid:  true,
			}
		}
	}

	if !params.Title.Valid {
		params.Title = sql.NullString{String: "New Conversation", Valid: true}
	}

	dbSession, err := s.db.Queries().CreateSession(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create session in DB: %w", err)
	}

	return &dbSession, nil
}

// saveMessageToDB saves a message to database
// saveMessageToDB saves a message to the database
// Phase 4: Now accepts topicID and threadID for proper message linking
func (s *AgentChatService) saveMessageToDB(ctx context.Context, msg *schema.Message, sessionID, userID, topicID, threadID string) (string, error) {
	params := convertEinoMessageToDB(msg, sessionID, userID)

	// Phase 4: Add topic and thread IDs if provided
	if topicID != "" {
		params.TopicID = sql.NullString{String: topicID, Valid: true}
	}
	if threadID != "" {
		params.ThreadID = sql.NullString{String: threadID, Valid: true}
	}

	dbMsg, err := s.db.Queries().CreateMessage(ctx, params)
	if err != nil {
		return "", fmt.Errorf("failed to save message to DB: %w", err)
	}

	return dbMsg.ID, nil
}

// updateSessionTimestamp updates the session's updated_at timestamp
func (s *AgentChatService) updateSessionTimestamp(ctx context.Context, sessionID, userID string) error {
	now := time.Now().UnixMilli()

	// Get current session
	dbSession, err := s.db.Queries().GetSession(ctx, db.GetSessionParams{
		ID:     sessionID,
		UserID: userID,
	})
	if err != nil {
		return fmt.Errorf("session not found: %w", err)
	}

	// Update timestamp
	_, err = s.db.Queries().UpdateSession(ctx, db.UpdateSessionParams{
		Title:           dbSession.Title,
		Description:     dbSession.Description,
		Avatar:          dbSession.Avatar,
		BackgroundColor: dbSession.BackgroundColor,
		GroupID:         dbSession.GroupID,
		Pinned:          dbSession.Pinned,
		UpdatedAt:       now,
		ID:              sessionID,
		UserID:          userID,
	})

	return err
}

// SetReasoningMode sets the reasoning mode for the service
// This affects all new conversations created after this call
// Validates hardware specs before allowing resource-intensive reasoning modes
func (s *AgentChatService) SetReasoningMode(mode ReasoningMode) error {
	s.sessionsMutex.Lock()
	defer s.sessionsMutex.Unlock()

	// Validate mode
	if mode != ReasoningDisabled && mode != ReasoningEnabled && mode != ReasoningVerbose {
		return fmt.Errorf("invalid reasoning mode: %s", mode)
	}

	// Get hardware specs from the installer
	var hardwareSpecs *llama.HardwareSpecs
	if s.libService != nil {
		hardwareSpecs = s.libService.GetHardwareSpecs()
	}

	// Create temp config to validate hardware
	tempConfig := s.reasoningConfig
	tempConfig.Mode = mode

	// Validate hardware requirements for reasoning modes
	if mode == ReasoningEnabled || mode == ReasoningVerbose {
		if valid, reason := tempConfig.ValidateHardware(hardwareSpecs); !valid {
			// Hardware insufficient - suggest alternative
			suggested := SuggestModeForHardware(hardwareSpecs)
			log.Printf("⚠️  Hardware validation failed for %s mode: %s", mode, reason)
			log.Printf("💡 Auto-switching to %s mode based on available hardware", suggested)

			// Auto-switch to suggested mode instead of failing
			mode = suggested
			tempConfig.Mode = suggested

			// If still trying reasoning mode after suggestion, validate again
			if mode != ReasoningDisabled {
				if valid, reason := tempConfig.ValidateHardware(hardwareSpecs); !valid {
					// Even suggested mode failed - force disable
					log.Printf("⚠️  Even suggested mode failed: %s. Forcing disabled mode.", reason)
					mode = ReasoningDisabled
				}
			}
		}
	}

	s.reasoningConfig.Mode = mode
	log.Printf("🧠 Reasoning mode changed to: %s (%s)", mode, s.reasoningConfig.GetModeDescription())

	// Log hardware requirements
	hwReq := s.reasoningConfig.GetHardwareRequirements()
	log.Printf("💻 Hardware requirements: %s", hwReq.Description)
	if hardwareSpecs != nil {
		log.Printf("📊 Current system: RAM=%dGB, Cores=%d, GPU=%s",
			hardwareSpecs.AvailableRAM, hardwareSpecs.CPUCores, hardwareSpecs.GPUModel)
	}

	// Log performance expectations
	perf := s.reasoningConfig.GetExpectedPerformance()
	log.Printf("📊 Expected performance:")
	log.Printf("   - Speed: %s", perf["speed"])
	log.Printf("   - Token efficiency: %s", perf["token_efficiency"])
	log.Printf("   - Max turns: %s", perf["max_turns"])
	log.Printf("   - Response size: %s", perf["response_size"])

	return nil
}

// GetReasoningMode returns the current reasoning mode
func (s *AgentChatService) GetReasoningMode() ReasoningMode {
	s.sessionsMutex.RLock()
	defer s.sessionsMutex.RUnlock()
	return s.reasoningConfig.Mode
}

// GetReasoningConfig returns the current reasoning configuration
func (s *AgentChatService) GetReasoningConfig() ReasoningConfig {
	s.sessionsMutex.RLock()
	defer s.sessionsMutex.RUnlock()
	return s.reasoningConfig
}

// ValidateModelForReasoningMode checks if the currently loaded model is appropriate
func (s *AgentChatService) ValidateModelForReasoningMode() error {
	modelPath := s.libService.GetLoadedChatModel()
	if modelPath == "" {
		return fmt.Errorf("no model loaded")
	}

	return s.reasoningConfig.ValidateModelForMode(modelPath)
}

// GetRecommendedModelForMode returns the recommended model for current reasoning mode
func (s *AgentChatService) GetRecommendedModelForMode() string {
	return s.reasoningConfig.GetRecommendedModel()
}

// SwitchToRecommendedModel loads the recommended model for current reasoning mode
func (s *AgentChatService) SwitchToRecommendedModel() error {
	recommendedModel := s.reasoningConfig.GetRecommendedModel()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	modelsDir := filepath.Join(homeDir, ".llama-cpp", "models")
	modelPath := filepath.Join(modelsDir, recommendedModel)

	log.Printf("🔄 Switching to recommended model for %s mode: %s", s.reasoningConfig.Mode, recommendedModel)

	if err := s.libService.LoadChatModel(modelPath); err != nil {
		return fmt.Errorf("failed to load recommended model: %w", err)
	}

	// Recreate llama model adapter
	s.llamaModel = llama.NewLlamaEinoModel(s.libService)

	log.Printf("✅ Successfully switched to %s", recommendedModel)
	return nil
}

// generateWithTokenStreaming generates response with token-by-token streaming
// This uses LibraryChatService for true token-by-token streaming
func (s *AgentChatService) generateWithTokenStreaming(ctx context.Context, session *AgentSession, messageID string, messages []*schema.Message) (string, error) {
	// Convert Eino messages to ChatMessage format
	chatMessages := make([]llama.ChatMessage, 0, len(messages))
	for _, msg := range messages {
		role := ""
		switch msg.Role {
		case schema.System:
			role = "system"
		case schema.User:
			role = "user"
		case schema.Assistant:
			role = "assistant"
		default:
			continue
		}
		chatMessages = append(chatMessages, llama.ChatMessage{
			Role:    role,
			Content: msg.Content,
		})
	}

	// Create LibraryChatService
	chatService := llama.NewLibraryChatService(s.libService, s.app)

	// Setup event forwarding from LibraryChatService to our format
	// LibraryChatService emits: stream:{requestID}:data with SSE format
	// We need to convert to: chat:stream:{sessionID} with our format

	streamRequestID := messageID
	eventName := fmt.Sprintf("stream:%s:data", streamRequestID)
	var fullContent strings.Builder

	// Throttling variables
	lastEmitTime := time.Now()           // Initialize to allow first chunk immediately
	const emitInterval = 1 * time.Second // Emit max every 1s

	// Emit start event
	s.app.Event.Emit("chat:stream", map[string]interface{}{
		"type":       "start",
		"session_id": session.SessionID,
		"message_id": messageID,
	})

	unsubscribe := s.app.Event.On(eventName, func(event *application.CustomEvent) {
		sseData, ok := event.Data.(string)
		if !ok {
			return
		}

		// Parse SSE data (format: "data: {...}\n\n")
		if !strings.HasPrefix(sseData, "data: ") {
			return
		}

		jsonData := strings.TrimPrefix(sseData, "data: ")
		jsonData = strings.TrimSpace(jsonData)

		if jsonData == "[DONE]" {
			// Streaming complete - always emit final state
			s.app.Event.Emit("chat:stream", map[string]interface{}{
				"type":       "complete",
				"session_id": session.SessionID,
				"message_id": messageID,
				"content":    fullContent.String(),
			})
			return
		}

		// Parse JSON chunk
		var chunk llama.ChatCompletionChunk
		if err := json.Unmarshal([]byte(jsonData), &chunk); err != nil {
			log.Printf("Failed to parse chunk: %v", err)
			return
		}

		// Extract content from chunk
		if len(chunk.Choices) > 0 {
			content := chunk.Choices[0].Delta.Content
			if content != "" {
				fullContent.WriteString(content)

				// Throttle chunk events - only emit if enough time has passed
				now := time.Now()
				if now.Sub(lastEmitTime) >= emitInterval {
					lastEmitTime = now

					// Emit our chunk event
					s.app.Event.Emit("chat:stream", map[string]interface{}{
						"type":         "chunk",
						"session_id":   session.SessionID,
						"message_id":   messageID,
						"content":      content,
						"full_content": fullContent.String(),
					})
				}
			}
		}
	})
	defer unsubscribe()

	// Start streaming generation
	err := chatService.ChatCompletionStream(ctx, streamRequestID, llama.ChatCompletionRequest{
		Messages:  chatMessages,
		MaxTokens: 2000,
		Stream:    true,
	})

	if err != nil {
		return "", fmt.Errorf("streaming generation failed: %w", err)
	}

	return fullContent.String(), nil
}

// buildPromptFromMessages builds a prompt string from Eino messages
func (s *AgentChatService) buildPromptFromMessages(messages []*schema.Message) (string, error) {
	// Convert Eino messages to simple prompt format
	var prompt strings.Builder
	for _, msg := range messages {
		switch msg.Role {
		case schema.System:
			prompt.WriteString(fmt.Sprintf("System: %s\n\n", msg.Content))
		case schema.User:
			prompt.WriteString(fmt.Sprintf("User: %s\n\n", msg.Content))
		case schema.Assistant:
			prompt.WriteString(fmt.Sprintf("Assistant: %s\n\n", msg.Content))
		}
	}
	prompt.WriteString("Assistant:")

	return prompt.String(), nil
}

// ============================================================================
// History Summary Functions
// ============================================================================

// autoSummarizeIfNeeded checks conditions and triggers summary automatically
// Runs in background goroutine, does not block chat response
func (s *AgentChatService) autoSummarizeIfNeeded(ctx context.Context, session *AgentSession, topicID, userID string) {
	if topicID == "" {
		return // No topic, no summary
	}

	// 1. Check if reasoning mode supports summarization
	threshold := s.reasoningConfig.GetSummaryThreshold()
	if threshold == 0 {
		return // Verbose mode doesn't need summary (3-5 turns only)
	}

	// 2. Check turn count
	turnCount := len(session.Messages) / 2
	if turnCount < threshold {
		return // Not enough turns yet
	}

	// 3. Check if we already summarized recently
	topic, err := s.db.Queries().GetTopic(ctx, db.GetTopicParams{
		ID:     topicID,
		UserID: userID,
	})
	if err != nil {
		log.Printf("⚠️  autoSummarize: Failed to get topic: %v", err)
		return
	}

	// If summary exists, this will be handled by incrementalSummarizeIfNeeded()
	if topic.HistorySummary.Valid && topic.HistorySummary.String != "" {
		return // Summary already exists, incremental will handle updates
	}

	// 4. Get messages from database (not from session, could be stale)
	messages, err := s.db.Queries().GetMessagesByTopicId(ctx, db.GetMessagesByTopicIdParams{
		TopicID: sql.NullString{String: topicID, Valid: true},
		UserID:  userID,
	})
	if err != nil || len(messages) < 4 {
		log.Printf("⚠️  autoSummarize: Not enough messages to summarize: %d", len(messages))
		return
	}

	// 5. Determine which messages to summarize
	// Keep recent messages based on reasoning mode
	keepCount := s.getKeepMessageCount()
	if len(messages) <= keepCount {
		return // Not enough messages to warrant summary
	}

	oldMessages := messages[:len(messages)-keepCount]

	// 6. Convert to Eino format
	einoMessages := make([]*schema.Message, 0, len(oldMessages))
	for _, dbMsg := range oldMessages {
		einoMsg, err := convertDBMessageToEino(&dbMsg)
		if err == nil {
			einoMessages = append(einoMessages, einoMsg)
		}
	}

	if len(einoMessages) < 2 {
		return // Need at least 2 messages (1 turn) to summarize
	}

	// 7. Generate summary (this is already in background goroutine)
	log.Printf("🔄 Auto-generating summary for topic %s (%d old messages, keeping %d recent)",
		topicID, len(einoMessages), keepCount)

	summary, err := s.generateHistorySummary(ctx, einoMessages)
	if err != nil {
		log.Printf("❌ Auto-summary failed: %v", err)
		return
	}

	// 8. Save summary to database
	now := time.Now().UnixMilli()
	metadata := fmt.Sprintf(`{"summarized_at":%d,"message_count":%d,"reasoning_mode":"%s"}`,
		now, len(einoMessages), s.reasoningConfig.Mode)

	_, err = s.db.Queries().UpdateTopic(ctx, db.UpdateTopicParams{
		Title:          topic.Title,
		HistorySummary: sql.NullString{String: summary, Valid: true},
		Metadata: sql.NullString{
			String: metadata,
			Valid:  true,
		},
		UpdatedAt: now,
		ID:        topicID,
		UserID:    userID,
	})

	if err != nil {
		log.Printf("❌ Failed to save auto-summary: %v", err)
	} else {
		log.Printf("✅ Auto-summary completed for topic %s (compressed %d messages into %d chars)",
			topicID, len(einoMessages), len(summary))

		// Emit subtle event to frontend (optional: show small toast or indicator)
		if s.app != nil {
			s.app.Event.Emit("chat:summary:auto-complete", map[string]interface{}{
				"topic_id":       topicID,
				"message_count":  len(einoMessages),
				"summary_length": len(summary),
			})
		}
	}
}

// getKeepMessageCount returns how many recent messages to keep (not summarize)
func (s *AgentChatService) getKeepMessageCount() int {
	switch s.reasoningConfig.Mode {
	case ReasoningDisabled:
		return 20 // Keep last 20 messages (10 turns)
	case ReasoningEnabled:
		return 12 // Keep last 12 messages (6 turns)
	case ReasoningVerbose:
		return 6 // Keep last 6 messages (3 turns) - but won't reach here
	default:
		return 16
	}
}

// generateHistorySummary generates summary using optimal model
// Uses 3-tier fallback: summary model → title model → main model
func (s *AgentChatService) generateHistorySummary(ctx context.Context, messages []*schema.Message) (string, error) {
	if len(messages) == 0 {
		return "", fmt.Errorf("no messages to summarize")
	}

	// Build summary prompt
	systemPrompt := `You're an assistant who's good at extracting key takeaways from conversations and summarizing them. 

Rules:
- Summarize the conversation in the user's original language
- Focus on key decisions, actions, and outcomes
- Maintain context for future conversation continuation
- Maximum 400 tokens
- Output ONLY the summary text, nothing else`

	// Build conversation text
	var conversationText string
	for _, msg := range messages {
		role := string(msg.Role)
		conversationText += fmt.Sprintf("%s: %s\n\n", role, msg.Content)
	}

	conversationContent := fmt.Sprintf(`<chat_history>
%s
</chat_history>

Please summarize the above conversation and retain key information. The summarized content will be used as context for subsequent prompts.`, conversationText)

	summaryMessages := []*schema.Message{
		{Role: schema.System, Content: systemPrompt},
		{Role: schema.User, Content: conversationContent},
	}

	// ✅ STRATEGY: Use dedicated summary model (preferred) → title model → main model
	var response *schema.Message
	var err error
	var modelUsed string

	currentModel := s.libService.GetLoadedChatModel()

	// Try 1: Dedicated summary model (BEST - optimized for this task)
	if s.summaryModelPath != "" && s.summaryModelPath != currentModel {
		log.Printf("📋 Using dedicated summary model: %s", filepath.Base(s.summaryModelPath))

		if loadErr := s.libService.LoadChatModel(s.summaryModelPath); loadErr != nil {
			log.Printf("⚠️  Failed to load summary model: %v", loadErr)
		} else {
			summaryModel := llama.NewLlamaEinoModel(s.libService)
			response, err = summaryModel.Generate(ctx, summaryMessages)
			modelUsed = "summary"

			// Restore main model
			if currentModel != "" {
				if restoreErr := s.libService.LoadChatModel(currentModel); restoreErr != nil {
					log.Printf("⚠️  Failed to restore main model: %v", restoreErr)
				}
			}
		}
	}

	// Try 2: Title model (GOOD - small and fast, already tested)
	if response == nil && s.titleModelPath != "" && s.titleModelPath != currentModel {
		log.Printf("📋 Falling back to title model: %s", filepath.Base(s.titleModelPath))

		if loadErr := s.libService.LoadChatModel(s.titleModelPath); loadErr != nil {
			log.Printf("⚠️  Failed to load title model: %v", loadErr)
		} else {
			titleModel := llama.NewLlamaEinoModel(s.libService)
			response, err = titleModel.Generate(ctx, summaryMessages)
			modelUsed = "title"

			// Restore main model
			if currentModel != "" {
				s.libService.LoadChatModel(currentModel)
			}
		}
	}

	// Try 3: Main chat model (FALLBACK - may be slow/reasoning model)
	if response == nil {
		log.Printf("📋 Using main chat model for summary (no dedicated model available)")
		response, err = s.llamaModel.Generate(ctx, summaryMessages)
		modelUsed = "main"
	}

	if err != nil {
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	// Clean up response
	summary := strings.TrimSpace(response.Content)

	// ⚠️ CRITICAL: Strip <think> tags if present
	// This can happen if main model is reasoning model (Qwen, DeepSeek)
	if strings.Contains(summary, "<think>") {
		log.Printf("⚠️  WARNING: Summary contains <think> tags (model: %s), stripping...", modelUsed)
		summary = stripThinkTags(summary)
	}

	if summary == "" {
		return "", fmt.Errorf("empty summary generated")
	}

	log.Printf("✅ Summary generated using %s model (%d chars)", modelUsed, len(summary))
	return summary, nil
}

// incrementalSummarizeIfNeeded checks if existing summary needs update with new messages
// Runs in background goroutine, does not block chat response
func (s *AgentChatService) incrementalSummarizeIfNeeded(ctx context.Context, session *AgentSession, topicID, userID string) {
	if topicID == "" {
		return // No topic, no summary
	}

	// 1. Check if reasoning mode supports incremental summarization
	threshold := s.reasoningConfig.GetIncrementalSummaryThreshold()
	if threshold == 0 {
		return // Verbose mode doesn't use incremental summary
	}

	// 2. Get topic with existing summary
	topic, err := s.db.Queries().GetTopic(ctx, db.GetTopicParams{
		ID:     topicID,
		UserID: userID,
	})
	if err != nil {
		return // Topic not found
	}

	// 3. Check if summary exists (incremental only works if base summary exists)
	if !topic.HistorySummary.Valid || topic.HistorySummary.String == "" {
		return // No existing summary, nothing to update
	}

	// 4. Parse metadata to get summarized message count
	var metadata struct {
		SummarizedMessageCount int   `json:"summarized_message_count"`
		SummaryVersion         int   `json:"summary_version"`
		InitialSummaryAt       int64 `json:"initial_summary_at"`
	}

	if topic.Metadata.Valid && topic.Metadata.String != "" {
		if err := json.Unmarshal([]byte(topic.Metadata.String), &metadata); err != nil {
			log.Printf("⚠️  incrementalSummarize: Failed to parse metadata: %v", err)
			return
		}
	}

	// 5. Get all messages from database
	allMessages, err := s.db.Queries().GetMessagesByTopicId(ctx, db.GetMessagesByTopicIdParams{
		TopicID: sql.NullString{String: topicID, Valid: true},
		UserID:  userID,
	})
	if err != nil {
		log.Printf("⚠️  incrementalSummarize: Failed to get messages: %v", err)
		return
	}

	// 6. Count new messages since last summary
	newMessageCount := len(allMessages) - metadata.SummarizedMessageCount
	if newMessageCount < threshold*2 { // *2 because 1 turn = 2 messages
		return // Not enough new messages yet
	}

	log.Printf("🔄 Re-summarizing topic %s (v%d → v%d, %d new messages)",
		topicID, metadata.SummaryVersion, metadata.SummaryVersion+1, newMessageCount)

	// 7. Load existing summary
	existingSummary := topic.HistorySummary.String

	// 8. Get new messages to incorporate
	newMessages := allMessages[metadata.SummarizedMessageCount:]
	keepCount := s.getKeepMessageCount()

	// Make sure we have messages to summarize
	if len(newMessages) <= keepCount {
		return // All new messages are in "keep" range, nothing to summarize
	}

	messagesToSummarize := newMessages[:len(newMessages)-keepCount]

	// 9. Convert DB messages to Eino messages
	einoMessages := make([]*schema.Message, 0, len(messagesToSummarize))
	for _, dbMsg := range messagesToSummarize {
		einoMsg, err := convertDBMessageToEino(&dbMsg)
		if err == nil {
			einoMessages = append(einoMessages, einoMsg)
		}
	}

	if len(einoMessages) == 0 {
		return // No valid messages to summarize
	}

	// 10. Generate merged summary
	mergedSummary, err := s.generateIncrementalSummary(ctx, existingSummary, einoMessages)
	if err != nil {
		log.Printf("❌ Incremental summary failed: %v", err)
		return
	}

	// 11. Update metadata
	if metadata.InitialSummaryAt == 0 {
		metadata.InitialSummaryAt = time.Now().UnixMilli()
	}

	newMetadata := map[string]interface{}{
		"summary_version":          metadata.SummaryVersion + 1,
		"last_summarized_at":       time.Now().UnixMilli(),
		"summarized_message_count": len(allMessages) - keepCount,
		"initial_summary_at":       metadata.InitialSummaryAt,
		"reasoning_mode":           string(s.reasoningConfig.Mode),
	}

	metadataJSON, err := json.Marshal(newMetadata)
	if err != nil {
		log.Printf("⚠️  Failed to marshal metadata: %v", err)
		metadataJSON = []byte("{}")
	}

	// 12. Save updated summary to database
	_, err = s.db.Queries().UpdateTopic(ctx, db.UpdateTopicParams{
		ID:             topicID,
		UserID:         userID,
		Title:          topic.Title,
		HistorySummary: sql.NullString{String: mergedSummary, Valid: true},
		Metadata:       sql.NullString{String: string(metadataJSON), Valid: true},
		UpdatedAt:      time.Now().UnixMilli(),
	})

	if err != nil {
		log.Printf("❌ Failed to save incremental summary: %v", err)
		return
	}

	log.Printf("✅ Incremental summary v%d completed (%d total messages compressed)",
		newMetadata["summary_version"], newMetadata["summarized_message_count"])
}

// generateIncrementalSummary creates updated summary by merging existing summary with new messages
func (s *AgentChatService) generateIncrementalSummary(ctx context.Context, existingSummary string, newMessages []*schema.Message) (string, error) {
	// Build conversation context from new messages
	var messagesText strings.Builder
	for i, msg := range newMessages {
		role := "User"
		if msg.Role == schema.Assistant {
			role = "Assistant"
		}
		messagesText.WriteString(fmt.Sprintf("%s: %s\n", role, msg.Content))
		if i < len(newMessages)-1 {
			messagesText.WriteString("\n")
		}
	}

	// Build prompt for incremental summarization
	prompt := fmt.Sprintf(`You are a conversation summarizer. Your task is to UPDATE an existing summary with new information.

EXISTING SUMMARY:
%s

NEW MESSAGES TO INCORPORATE:
%s

INSTRUCTIONS:
1. Read the existing summary carefully
2. Read the new messages to identify new topics and developments
3. Create an UPDATED summary that:
   - Preserves important information from the existing summary
   - Adds new topics and details from recent messages
   - Maintains chronological flow (what happened first, then next)
   - Stays concise (maximum 400 tokens)
4. Use the same language as the original messages
5. Focus on key topics, decisions, and important information
6. DO NOT use <think> tags or any XML/HTML markup

UPDATED SUMMARY:`, existingSummary, messagesText.String())

	// Use same 3-tier model selection as regular summary
	summaryMessages := []*schema.Message{
		{
			Role:    schema.User,
			Content: prompt,
		},
	}

	// Try 1: Summary model (BEST)
	var response *schema.Message
	var err error
	modelUsed := "unknown"

	currentModel := s.libService.GetLoadedChatModel()

	if s.summaryModelPath != "" && s.summaryModelPath != currentModel {
		log.Printf("📋 Using dedicated summary model for incremental: %s", filepath.Base(s.summaryModelPath))

		if loadErr := s.libService.LoadChatModel(s.summaryModelPath); loadErr != nil {
			log.Printf("⚠️  Failed to load summary model: %v", loadErr)
		} else {
			summaryModel := llama.NewLlamaEinoModel(s.libService)
			response, err = summaryModel.Generate(ctx, summaryMessages)
			modelUsed = "summary"

			// Restore main model
			if currentModel != "" {
				s.libService.LoadChatModel(currentModel)
			}
		}
	}

	// Try 2: Title model (GOOD)
	if response == nil && s.titleModelPath != "" && s.titleModelPath != currentModel {
		log.Printf("📋 Falling back to title model for incremental: %s", filepath.Base(s.titleModelPath))

		if loadErr := s.libService.LoadChatModel(s.titleModelPath); loadErr != nil {
			log.Printf("⚠️  Failed to load title model: %v", loadErr)
		} else {
			titleModel := llama.NewLlamaEinoModel(s.libService)
			response, err = titleModel.Generate(ctx, summaryMessages)
			modelUsed = "title"

			// Restore main model
			if currentModel != "" {
				s.libService.LoadChatModel(currentModel)
			}
		}
	}

	// Try 3: Main chat model (FALLBACK)
	if response == nil {
		log.Printf("📋 Using main chat model for incremental summary")
		response, err = s.llamaModel.Generate(ctx, summaryMessages)
		modelUsed = "main"
	}

	if err != nil {
		return "", fmt.Errorf("failed to generate incremental summary: %w", err)
	}

	// Clean up response
	summary := strings.TrimSpace(response.Content)

	// Strip <think> tags if present
	if strings.Contains(summary, "<think>") {
		log.Printf("⚠️  WARNING: Incremental summary contains <think> tags (model: %s), stripping...", modelUsed)
		summary = stripThinkTags(summary)
	}

	if summary == "" {
		return "", fmt.Errorf("empty incremental summary generated")
	}

	log.Printf("✅ Incremental summary generated using %s model (%d chars)", modelUsed, len(summary))
	return summary, nil
}
