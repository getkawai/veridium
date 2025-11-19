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

	return &AgentChatService{
		app:           app,
		db:            db,
		libService:    libService,
		llamaModel:    llamaModel,
		kbService:     kbService,
		ragWorkflow:   ragWorkflow,
		toolsBridge:   toolsBridge,
		contextBridge: contextBridge,
		threadService: threadService,
		sessions:      make(map[string]*AgentSession),
	}
}

// Chat processes a chat request and returns a response
func (s *AgentChatService) Chat(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// Get or create session
	session, err := s.getOrCreateSession(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create session: %w", err)
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

	// Run agent
	iterator := session.Agent.Run(ctx, agentInput)

	// Collect response
	var finalMessage string
	var toolCalls []schema.ToolCall
	var finishReason string
	var usage *schema.TokenUsage

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

	// Phase 4: Generate topic title after first response if not already created
	if currentTopicID == "" && len(session.Messages) == 2 {
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

	// Save assistant message to DB with topic and thread IDs
	assistantMsgID, err := s.saveMessageToDB(ctx, assistantMsg, session.SessionID, session.UserID, currentTopicID, req.ThreadID)
	if err != nil {
		log.Printf("⚠️  Warning: Failed to save assistant message to DB: %v", err)
	} else {
		log.Printf("💾 Saved assistant message: %s (topic: %s, thread: %s)", assistantMsgID, currentTopicID, req.ThreadID)
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

// ChatStream processes a chat request with streaming response
// TODO: Implement streaming via Wails v3 events API once documented
func (s *AgentChatService) ChatStream(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// For now, use synchronous chat
	// TODO: Implement proper streaming when Wails v3 event API is stable
	return s.Chat(ctx, req)
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
	// Build instruction
	instruction := "You are a helpful AI assistant. "
	if session.KnowledgeBaseID != "" {
		instruction += "You have access to a knowledge base. Use the search tool to find relevant information before answering questions. "
	}
	instruction += "Be concise and accurate in your responses."

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

// generateTopicTitle generates a concise title for the conversation using LLM
func (s *AgentChatService) generateTopicTitle(ctx context.Context, messages []*schema.Message, locale string) (string, error) {
	if len(messages) == 0 {
		return "New Conversation", nil
	}

	// Build summary prompt similar to frontend chainSummaryTitle
	systemPrompt := fmt.Sprintf(`You are a professional conversation summarizer. Generate a concise title that captures the essence of the conversation.

Rules:
- Output ONLY the title text, no explanations or additional context
- Maximum 10 words
- Maximum 50 characters
- No punctuation marks, quotes, or special characters
- Do not wrap the title in quotation marks or any other delimiters
- Use the language specified by the locale code: %s
- The title should accurately reflect the main topic of the conversation
- Keep it short and to the point

Example output format: Sleep Functions for Body and Mind`, locale)

	// Build conversation text
	var conversationText string
	for _, msg := range messages {
		role := "user"
		if msg.Role == schema.Assistant {
			role = "assistant"
		}
		conversationText += fmt.Sprintf("%s: %s\n", role, msg.Content)
	}

	// Create messages for title generation
	titleMessages := []*schema.Message{
		{Role: schema.System, Content: systemPrompt},
		{Role: schema.User, Content: conversationText},
	}

	// Generate title using LlamaEinoModel
	response, err := s.llamaModel.Generate(ctx, titleMessages)
	if err != nil {
		return "", fmt.Errorf("failed to generate title: %w", err)
	}

	// Clean up the title
	title := strings.TrimSpace(response.Content)
	if len(title) > 50 {
		title = title[:50]
	}

	return title, nil
}

// createTopicForSession creates a topic for a session after first response
func (s *AgentChatService) createTopicForSession(ctx context.Context, sessionID, userID string, messages []*schema.Message) error {
	// Check if topic already exists for this session
	count, err := s.db.Queries().CountTopicsBySession(ctx, db.CountTopicsBySessionParams{
		SessionID: sql.NullString{String: sessionID, Valid: true},
		UserID:    userID,
	})
	if err != nil {
		return fmt.Errorf("failed to check existing topics: %w", err)
	}

	if count > 0 {
		// Topic already exists, skip
		return nil
	}

	// Generate title (default locale: en-US)
	title, err := s.generateTopicTitle(ctx, messages, "en-US")
	if err != nil {
		log.Printf("⚠️  Warning: Failed to generate topic title: %v", err)
		title = "New Conversation"
	}

	// Create topic in database
	topicID := uuid.New().String()
	now := time.Now().UnixMilli()

	_, err = s.db.Queries().CreateTopic(ctx, db.CreateTopicParams{
		ID:             topicID,
		Title:          sql.NullString{String: title, Valid: true},
		Favorite:       0,
		SessionID:      sql.NullString{String: sessionID, Valid: true},
		GroupID:        sql.NullString{},
		UserID:         userID,
		ClientID:       sql.NullString{},
		HistorySummary: sql.NullString{},
		Metadata:       sql.NullString{},
		CreatedAt:      now,
		UpdatedAt:      now,
	})

	if err != nil {
		return fmt.Errorf("failed to create topic: %w", err)
	}

	log.Printf("📝 Created topic for session %s: %s", sessionID, title)
	return nil
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

	_, err = s.db.Queries().CreateTopic(ctx, db.CreateTopicParams{
		ID:             topicID,
		Title:          sql.NullString{String: "New Conversation", Valid: true},
		Favorite:       0,
		SessionID:      sql.NullString{String: sessionID, Valid: true},
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

	// Generate title (default locale: en-US)
	title, err := s.generateTopicTitle(ctx, messages, "en-US")
	if err != nil {
		log.Printf("⚠️  Warning: Failed to generate topic title: %v", err)
		title = "New Conversation"
	}

	// Create topic in database
	topicID := uuid.New().String()
	now := time.Now().UnixMilli()

	_, err = s.db.Queries().CreateTopic(ctx, db.CreateTopicParams{
		ID:             topicID,
		Title:          sql.NullString{String: title, Valid: true},
		Favorite:       0,
		SessionID:      sql.NullString{String: sessionID, Valid: true},
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
