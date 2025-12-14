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
	"time"

	"github.com/google/uuid"
	"github.com/kawai-network/veridium/fantasy"
	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/pkg/xlog"
	"github.com/kawai-network/veridium/types"
)

// ============================================================================
// Stream Event Types for ChatMockStream and ChatRealStream
// ============================================================================

// ToolResultData holds content and state for a tool result
type ToolResultData struct {
	Content interface{}
	State   interface{}
}

// StreamEventPayload represents the payload for chat stream events.
// This is emitted via Wails events ('chat:stream') and consumed by frontend.
// Frontend can import this type from generated bindings.
type StreamEventPayload struct {
	Type      types.ChatStreamEvent `json:"type"`
	SessionID string                `json:"session_id"`
	MessageID string                `json:"message_id"`
	TopicID   string                `json:"topic_id,omitempty"`

	// Content fields
	Content     string `json:"content,omitempty"`
	FullContent string `json:"full_content,omitempty"` // Legacy support

	// Reasoning
	Reasoning *ModelReasoning `json:"reasoning,omitempty"`

	// Tools
	Tool  *ChatToolPayload  `json:"tool,omitempty"`  // For single tool call event (legacy/alternative)
	Tools []ChatToolPayload `json:"tools,omitempty"` // For tool_call event

	// Tool Result
	ToolCallID  string             `json:"tool_call_id,omitempty"`
	ToolMsgID   string             `json:"tool_msg_id,omitempty"`
	Plugin      *ChatPluginPayload `json:"plugin,omitempty"`
	PluginState interface{}        `json:"pluginState,omitempty"`

	// Complete event fields
	Search      *GroundingSearch  `json:"search,omitempty"`
	ChunksList  []ChatFileChunk   `json:"chunksList,omitempty"`
	ImageList   []ChatImageItem   `json:"imageList,omitempty"`
	Usage       *ModelUsage       `json:"usage,omitempty"`
	Error       *ChatMessageError `json:"error,omitempty"`
	Performance *ModelPerformance `json:"performance,omitempty"`
}

// ============================================================================
// Parameter Structs for Reusable Methods
// ============================================================================

// SessionSetupResult contains the result of session and topic setup
type SessionSetupResult struct {
	Session       *AgentSession
	TopicID       string
	SessionID     string
	UserMessageID string // ID of saved user message
	IsNew         bool
}

// SaveUserMessageParams contains parameters for saving a user message
type SaveUserMessageParams struct {
	MessageID string // Optional: Pre-generated message ID from frontend
	Content   string
	SessionID string
	TopicID   string
	ThreadID  string
}

// SaveAssistantMessageParams contains parameters for saving an assistant message
type SaveAssistantMessageParams struct {
	MessageID string // Optional: Pre-generated message ID from frontend
	Content   string
	SessionID string
	TopicID   string
	ThreadID  string

	// Optional fields for rich content
	Reasoning interface{} // map[string]interface{} with content, status
	Tools     interface{} // []map[string]interface{} with tool calls
	Search    interface{} // map[string]interface{} with citations, searchQueries
	Metadata  interface{} // map[string]interface{} with model, temperature, chunksList, etc.
	Error     interface{} // error information if any
}

// SaveToolMessageParams contains parameters for saving a tool message
type SaveToolMessageParams struct {
	ToolCallID string      // ID linking to tool call in assistant message
	Identifier string      // Plugin identifier (e.g., "lobe-web-browsing")
	APIName    string      // API name (e.g., "search", "crawlSinglePage")
	Arguments  string      // JSON string of tool arguments
	Content    interface{} // Tool result (will be JSON marshaled)
	State      interface{} // Plugin state for frontend (will be JSON marshaled)
	SessionID  string
	TopicID    string
	ThreadID   string
	TimeOffset int64 // Offset from base timestamp (for ordering)
}

// SaveRAGDataParams contains parameters for saving RAG-related data
type SaveRAGDataParams struct {
	MessageID string
	UserQuery string
	Files     []db.CreateFileParams
	Chunks    []RAGChunkParams
}

// RAGChunkParams contains parameters for a single chunk with its file link
type RAGChunkParams struct {
	ID         string
	FileIndex  int // Index into SaveRAGDataParams.Files array (-1 if no file link)
	Text       string
	ChunkIndex int64
	Type       string
	Similarity int64 // 0-100 similarity score
}

// ============================================================================
// Reusable Helper Methods
// ============================================================================

// setupSessionAndTopic gets or creates a session, adds user message, and auto-creates a topic if needed.
// This follows the same pattern as Chat() in agent_chat_service.go:
// 1. Get or create session
// 2. Load history summary from DB (if topic exists)
// 3. Load thread messages (if ThreadID provided)
// 4. Auto-create topic if needed (before saving user message)
// 5. Add user message to session (in-memory)
// 6. Save user message to DB
//
// Note: Model validation is NOT included here - caller should handle it if needed.
func (s *AgentChatService) setupSessionAndTopic(ctx context.Context, req ChatRequest) (*SessionSetupResult, error) {
	xlog.Info("🔧 SetupSessionAndTopic", "session_id", req.SessionID, "topic_id", req.TopicID, "thread_id", req.ThreadID)

	// 1. Get or create session
	session, err := s.getOrCreateSession(ctx, req, req.TopicID)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create session: %w", err)
	}

	result := &SessionSetupResult{
		Session:   session,
		SessionID: session.SessionID,
		TopicID:   req.TopicID,
		IsNew:     len(session.Messages) == 0,
	}

	// 2. Load history summary from DB if topic exists
	if req.TopicID != "" {
		topic, err := s.db.Queries().GetTopic(ctx, req.TopicID)
		if err == nil && topic.HistorySummary.Valid && topic.HistorySummary.String != "" {
			if session.Context == nil {
				session.Context = make(map[string]any)
			}
			session.Context["history_summary"] = topic.HistorySummary.String
			xlog.Info("📋 Loaded history summary for topic", "topic_id", req.TopicID, "chars", len(topic.HistorySummary.String))
		}
	}

	// 3. Load thread messages if ThreadID provided
	if req.ThreadID != "" && s.threadService != nil {
		threadMessages, err := s.threadService.GetThreadMessages(ctx, req.ThreadID)
		if err != nil {
			xlog.Warn("⚠️  Warning: Failed to load thread messages", "error", err)
		} else {
			// Convert thread messages to message format
			yzmaMessages := make([]fantasy.Message, 0, len(threadMessages))
			for _, dbMsg := range threadMessages {
				if msg, ok := convertDBMessageToYzma(&dbMsg); ok {
					yzmaMessages = append(yzmaMessages, msg)
				}
			}
			session.Messages = yzmaMessages
			session.ThreadID = req.ThreadID
			xlog.Info("📋 Loaded messages from thread", "count", len(yzmaMessages), "thread_id", req.ThreadID)
		}
	}

	// 4. Set TopicID from request
	if req.TopicID != "" {
		session.TopicID = req.TopicID
	}

	// 5. Auto-create topic if needed (BEFORE saving user message so we have topicID)
	if result.TopicID == "" {
		topicID, err := s.createTopicForSessionSync(ctx, session.SessionID)
		if err != nil {
			xlog.Warn("⚠️  Warning: Failed to create topic", "error", err)
		} else {
			result.TopicID = topicID
			session.TopicID = topicID
			xlog.Info("📝 Auto-created topic", "topic_id", topicID)
		}
	}

	// 6. Add user message to session (in-memory for LLM context)
	session.Messages = append(session.Messages, fantasy.NewUserMessage(req.Message))

	// 7. Save user message to DB
	userMsgID, err := s.saveUserMessage(ctx, SaveUserMessageParams{
		MessageID: req.MessageUserID, // Use pre-generated ID from frontend
		Content:   req.Message,
		SessionID: session.SessionID,
		TopicID:   result.TopicID,
		ThreadID:  req.ThreadID,
	})
	if err != nil {
		xlog.Warn("⚠️  Warning: Failed to save user message to DB", "error", err)
	} else {
		result.UserMessageID = userMsgID
		xlog.Info("💾 Saved user message", "message_id", userMsgID, "topic_id", result.TopicID, "thread_id", req.ThreadID)
	}

	return result, nil
}

// saveUserMessage saves a user message to the database
// Returns the message ID
func (s *AgentChatService) saveUserMessage(ctx context.Context, params SaveUserMessageParams) (string, error) {
	// Use pre-generated ID if provided, otherwise generate new one
	msgID := params.MessageID
	if msgID == "" {
		msgID = uuid.New().String()
	}
	now := time.Now().UnixMilli()

	dbParams := db.CreateMessageParams{
		ID:        msgID,
		Role:      "user",
		Content:   sql.NullString{String: params.Content, Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if params.SessionID != "" {
		dbParams.SessionID = sql.NullString{String: params.SessionID, Valid: true}
	}
	if params.TopicID != "" {
		dbParams.TopicID = sql.NullString{String: params.TopicID, Valid: true}
	}
	if params.ThreadID != "" {
		dbParams.ThreadID = sql.NullString{String: params.ThreadID, Valid: true}
	}

	_, err := s.db.Queries().CreateMessage(ctx, dbParams)
	if err != nil {
		return "", fmt.Errorf("failed to save user message: %w", err)
	}

	xlog.Info("💾 Saved user message", "message_id", msgID)
	return msgID, nil
}

// saveAssistantMessage saves an assistant message with all optional rich content
// Returns the message ID
func (s *AgentChatService) saveAssistantMessage(ctx context.Context, params SaveAssistantMessageParams) (string, error) {
	// Use pre-generated ID if provided, otherwise generate new one
	msgID := params.MessageID
	if msgID == "" {
		msgID = uuid.New().String()
	}
	now := time.Now().UnixMilli()

	dbParams := db.CreateMessageParams{
		ID:        msgID,
		Role:      "assistant",
		Content:   sql.NullString{String: params.Content, Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if params.SessionID != "" {
		dbParams.SessionID = sql.NullString{String: params.SessionID, Valid: true}
	}
	if params.TopicID != "" {
		dbParams.TopicID = sql.NullString{String: params.TopicID, Valid: true}
	}
	if params.ThreadID != "" {
		dbParams.ThreadID = sql.NullString{String: params.ThreadID, Valid: true}
	}

	// Marshal optional JSON fields
	if params.Reasoning != nil {
		if reasoningJSON, err := json.Marshal(params.Reasoning); err == nil {
			dbParams.Reasoning = sql.NullString{String: string(reasoningJSON), Valid: true}
		}
	}
	if params.Tools != nil {
		if toolsJSON, err := json.Marshal(params.Tools); err == nil {
			dbParams.Tools = sql.NullString{String: string(toolsJSON), Valid: true}
		}
	}
	if params.Search != nil {
		if searchJSON, err := json.Marshal(params.Search); err == nil {
			dbParams.Search = sql.NullString{String: string(searchJSON), Valid: true}
		}
	}
	if params.Metadata != nil {
		if metadataJSON, err := json.Marshal(params.Metadata); err == nil {
			dbParams.Metadata = sql.NullString{String: string(metadataJSON), Valid: true}
		}
	}
	if params.Error != nil {
		if errorJSON, err := json.Marshal(params.Error); err == nil {
			dbParams.Error = sql.NullString{String: string(errorJSON), Valid: true}
		}
	}

	_, err := s.db.Queries().CreateMessage(ctx, dbParams)
	if err != nil {
		return "", fmt.Errorf("failed to save assistant message: %w", err)
	}

	xlog.Info("💾 Saved assistant message", "message_id", msgID)
	return msgID, nil
}

// saveToolMessage saves a tool message with its plugin entry
// Returns the message ID
func (s *AgentChatService) saveToolMessage(ctx context.Context, params SaveToolMessageParams) (string, error) {
	msgID := uuid.New().String()
	now := time.Now().UnixMilli()

	// Marshal content
	var contentStr string
	if params.Content != nil {
		if contentJSON, err := json.Marshal(params.Content); err == nil {
			contentStr = string(contentJSON)
		}
	}

	// Create message
	msgParams := db.CreateMessageParams{
		ID:        msgID,
		Role:      "tool",
		Content:   sql.NullString{String: contentStr, Valid: contentStr != ""},
		CreatedAt: now + params.TimeOffset,
		UpdatedAt: now + params.TimeOffset,
	}

	if params.SessionID != "" {
		msgParams.SessionID = sql.NullString{String: params.SessionID, Valid: true}
	}
	if params.TopicID != "" {
		msgParams.TopicID = sql.NullString{String: params.TopicID, Valid: true}
	}
	if params.ThreadID != "" {
		msgParams.ThreadID = sql.NullString{String: params.ThreadID, Valid: true}
	}

	_, err := s.db.Queries().CreateMessage(ctx, msgParams)
	if err != nil {
		return "", fmt.Errorf("failed to save tool message %s: %w", params.ToolCallID, err)
	}

	// Create plugin entry
	pluginParams := db.CreateMessagePluginParams{
		ID:         msgID,
		ToolCallID: sql.NullString{String: params.ToolCallID, Valid: true},
		Type:       sql.NullString{String: "builtin", Valid: true},
		ApiName:    sql.NullString{String: params.APIName, Valid: true},
		Arguments:  sql.NullString{String: params.Arguments, Valid: params.Arguments != ""},
		Identifier: sql.NullString{String: params.Identifier, Valid: true},
	}

	// Add state if provided (for pluginState in frontend)
	if params.State != nil {
		if stateJSON, err := json.Marshal(params.State); err == nil {
			pluginParams.State = sql.NullString{String: string(stateJSON), Valid: true}
		}
	}

	_, err = s.db.Queries().CreateMessagePlugin(ctx, pluginParams)
	if err != nil {
		return "", fmt.Errorf("failed to save tool plugin %s: %w", params.ToolCallID, err)
	}

	xlog.Info("💾 Saved tool message", "tool_call_id", params.ToolCallID, "identifier", params.Identifier, "api_name", params.APIName)
	return msgID, nil
}

// saveRAGData saves RAG-related data (files, chunks, and links them to a message)
func (s *AgentChatService) saveRAGData(ctx context.Context, params SaveRAGDataParams) error {
	// 1. Create files and collect their IDs
	fileIDs := make([]string, len(params.Files))
	for i, file := range params.Files {
		// UserID removed from file params
		createdFile, err := s.db.Queries().CreateFile(ctx, file)
		if err != nil {
			xlog.Warn("⚠️  Failed to create file", "error", err)
			continue
		}
		fileIDs[i] = createdFile.ID
	}

	// 2. Create chunks and link to files using FileIndex
	for _, chunk := range params.Chunks {
		chunkParams := db.CreateChunkParams{
			ID:         chunk.ID,
			Text:       sql.NullString{String: chunk.Text, Valid: true},
			ChunkIndex: sql.NullInt64{Int64: chunk.ChunkIndex, Valid: true},
			Type:       sql.NullString{String: chunk.Type, Valid: true},
		}
		_, err := s.db.Queries().CreateChunk(ctx, chunkParams)
		if err != nil {
			xlog.Warn("⚠️  Failed to create chunk", "chunk_id", chunk.ID, "error", err)
			continue
		}

		// Link chunk to file using FileIndex
		if chunk.FileIndex >= 0 && chunk.FileIndex < len(fileIDs) && fileIDs[chunk.FileIndex] != "" {
			fileID := fileIDs[chunk.FileIndex]
			err = s.db.Queries().LinkFileToChunk(ctx, db.LinkFileToChunkParams{
				FileID:  sql.NullString{String: fileID, Valid: true},
				ChunkID: sql.NullString{String: chunk.ID, Valid: true},
			})
			if err != nil {
				xlog.Warn("⚠️  Failed to link file to chunk", "chunk_id", chunk.ID, "error", err)
			}
		}
	}

	// 3. Create message query and link chunks to message
	if params.MessageID != "" && len(params.Chunks) > 0 {
		queryID := uuid.New().String()
		queryParams := db.CreateMessageQueryParams{
			ID:           queryID,
			MessageID:    params.MessageID,
			UserQuery:    sql.NullString{String: params.UserQuery, Valid: true},
			RewriteQuery: sql.NullString{String: params.UserQuery, Valid: true},
		}
		_, err := s.db.Queries().CreateMessageQuery(ctx, queryParams)
		if err != nil {
			xlog.Warn("⚠️  Failed to create message query", "error", err)
		} else {
			// Link chunks to message query
			for _, chunk := range params.Chunks {
				err = s.db.Queries().LinkMessageQueryToChunk(ctx, db.LinkMessageQueryToChunkParams{
					MessageID:  sql.NullString{String: params.MessageID, Valid: true},
					QueryID:    sql.NullString{String: queryID, Valid: true},
					ChunkID:    sql.NullString{String: chunk.ID, Valid: true},
					Similarity: sql.NullInt64{Int64: chunk.Similarity, Valid: true},
				})
				if err != nil {
					xlog.Warn("⚠️  Failed to link query to chunk", "chunk_id", chunk.ID, "error", err)
				}
			}
		}
	}

	xlog.Info("💾 Saved RAG data", "files", len(params.Files), "chunks", len(params.Chunks))
	return nil
}

// linkMessageToChunks creates links between a message and chunks with similarity scores
// This is a simpler version when files/chunks already exist
func (s *AgentChatService) linkMessageToChunks(ctx context.Context, messageID, userQuery string, chunks []RAGChunkParams) error {
	if messageID == "" || len(chunks) == 0 {
		return nil
	}

	queryID := uuid.New().String()
	queryParams := db.CreateMessageQueryParams{
		ID:           queryID,
		MessageID:    messageID,
		UserQuery:    sql.NullString{String: userQuery, Valid: true},
		RewriteQuery: sql.NullString{String: userQuery, Valid: true},
	}

	_, err := s.db.Queries().CreateMessageQuery(ctx, queryParams)
	if err != nil {
		return fmt.Errorf("failed to create message query: %w", err)
	}

	for _, chunk := range chunks {
		err = s.db.Queries().LinkMessageQueryToChunk(ctx, db.LinkMessageQueryToChunkParams{
			MessageID:  sql.NullString{String: messageID, Valid: true},
			QueryID:    sql.NullString{String: queryID, Valid: true},
			ChunkID:    sql.NullString{String: chunk.ID, Valid: true},
			Similarity: sql.NullInt64{Int64: chunk.Similarity, Valid: true},
		})
		if err != nil {
			xlog.Warn("⚠️  Failed to link message to chunk", "chunk_id", chunk.ID, "error", err)
		}
	}

	xlog.Info("💾 Linked message to chunks", "message_id", messageID, "chunks", len(chunks))
	return nil
}

// GetStreamEventPayloadType returns an empty StreamEventPayload for type inference.
// This method exists to expose StreamEventPayload in Wails bindings.
// Frontend can use this type for handling 'chat:stream' events.
func (s *AgentChatService) GetStreamEventPayloadType() StreamEventPayload {
	return StreamEventPayload{}
}
