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
	"fmt"
	"sync"
	"time"

	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/pkg/yzma/message"
)

// ============================================================================
// Mock LLM Generator
// ============================================================================

// MockLLMGenerator implements LLMGenerator for testing
type MockLLMGenerator struct {
	// Response to return from Generate/RunAgentLoop
	Response *llama.YzmaResponse

	// Tool messages to return from RunAgentLoop
	ToolMessages []message.Message

	// Error to return (if set, overrides Response)
	Error error

	// Callback for inspecting calls
	OnGenerate      func(ctx context.Context, messages []message.Message)
	OnRunAgentLoop  func(ctx context.Context, messages []message.Message, maxIterations int)

	// Track calls for assertions
	GenerateCalls      [][]message.Message
	RunAgentLoopCalls  []RunAgentLoopCall
	StreamingCallbacks []func(token string, isLast bool)

	// Tools configuration
	toolNames []string

	mu sync.Mutex
}

// RunAgentLoopCall records a call to RunAgentLoop
type RunAgentLoopCall struct {
	Messages      []message.Message
	MaxIterations int
}

// NewMockLLMGenerator creates a new mock with default success response
func NewMockLLMGenerator() *MockLLMGenerator {
	return &MockLLMGenerator{
		Response: &llama.YzmaResponse{
			Content:      "Mock response",
			FinishReason: "stop",
		},
		GenerateCalls:     make([][]message.Message, 0),
		RunAgentLoopCalls: make([]RunAgentLoopCall, 0),
	}
}

// Generate implements LLMGenerator.Generate
func (m *MockLLMGenerator) Generate(ctx context.Context, messages []message.Message) (*llama.YzmaResponse, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.GenerateCalls = append(m.GenerateCalls, messages)

	if m.OnGenerate != nil {
		m.OnGenerate(ctx, messages)
	}

	if m.Error != nil {
		return nil, m.Error
	}

	return m.Response, nil
}

// RunAgentLoop implements LLMGenerator.RunAgentLoop
func (m *MockLLMGenerator) RunAgentLoop(ctx context.Context, messages []message.Message, maxIterations int) (*llama.YzmaResponse, []message.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.RunAgentLoopCalls = append(m.RunAgentLoopCalls, RunAgentLoopCall{
		Messages:      messages,
		MaxIterations: maxIterations,
	})

	if m.OnRunAgentLoop != nil {
		m.OnRunAgentLoop(ctx, messages, maxIterations)
	}

	if m.Error != nil {
		return nil, nil, m.Error
	}

	return m.Response, m.ToolMessages, nil
}

// RunAgentLoopWithStreaming implements LLMGenerator.RunAgentLoopWithStreaming
func (m *MockLLMGenerator) RunAgentLoopWithStreaming(ctx context.Context, messages []message.Message, maxIterations int, callback func(token string, isLast bool)) (*llama.YzmaResponse, []message.Message, error) {
	m.mu.Lock()
	m.StreamingCallbacks = append(m.StreamingCallbacks, callback)
	m.mu.Unlock()

	// Simulate streaming by calling callback with tokens
	if callback != nil && m.Response != nil {
		tokens := []string{}
		content := m.Response.Content
		// Split into words for simulation
		for i := 0; i < len(content); i += 5 {
			end := i + 5
			if end > len(content) {
				end = len(content)
			}
			tokens = append(tokens, content[i:end])
		}

		for i, token := range tokens {
			callback(token, i == len(tokens)-1)
		}
	}

	return m.RunAgentLoop(ctx, messages, maxIterations)
}

// WithTools implements LLMGenerator.WithTools
func (m *MockLLMGenerator) WithTools(toolNames []string) LLMGenerator {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Return a new mock with the same config but different tool names
	newMock := &MockLLMGenerator{
		Response:          m.Response,
		ToolMessages:      m.ToolMessages,
		Error:             m.Error,
		OnGenerate:        m.OnGenerate,
		OnRunAgentLoop:    m.OnRunAgentLoop,
		GenerateCalls:     m.GenerateCalls,     // Share call tracking
		RunAgentLoopCalls: m.RunAgentLoopCalls, // Share call tracking
		toolNames:         toolNames,
	}
	return newMock
}

// SetResponse sets the response to return
func (m *MockLLMGenerator) SetResponse(content string, finishReason string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Response = &llama.YzmaResponse{
		Content:      content,
		FinishReason: finishReason,
	}
}

// SetError sets an error to return
func (m *MockLLMGenerator) SetError(err error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Error = err
}

// SetToolCalls sets tool calls in the response
func (m *MockLLMGenerator) SetToolCalls(toolCalls []message.ToolCall) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.Response == nil {
		m.Response = &llama.YzmaResponse{}
	}
	m.Response.ToolCalls = toolCalls
	m.Response.FinishReason = "tool_calls"
}

// GetGenerateCallCount returns the number of Generate calls
func (m *MockLLMGenerator) GetGenerateCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.GenerateCalls)
}

// GetRunAgentLoopCallCount returns the number of RunAgentLoop calls
func (m *MockLLMGenerator) GetRunAgentLoopCallCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.RunAgentLoopCalls)
}

// ============================================================================
// Mock Database Queries
// ============================================================================

// MockDBQueries implements a subset of db.Queries for testing
type MockDBQueries struct {
	// Sessions storage
	Sessions map[string]db.Session
	Messages map[string][]db.Message
	Topics   map[string]db.Topic

	// Errors to return
	GetSessionError    error
	CreateSessionError error
	CreateMessageError error
	GetTopicError      error
	CreateTopicError   error

	// Call tracking
	CreateSessionCalls []db.CreateSessionParams
	CreateMessageCalls []db.CreateMessageParams
	CreateTopicCalls   []db.CreateTopicParams

	mu sync.Mutex
}

// NewMockDBQueries creates a new mock database
func NewMockDBQueries() *MockDBQueries {
	return &MockDBQueries{
		Sessions:           make(map[string]db.Session),
		Messages:           make(map[string][]db.Message),
		Topics:             make(map[string]db.Topic),
		CreateSessionCalls: make([]db.CreateSessionParams, 0),
		CreateMessageCalls: make([]db.CreateMessageParams, 0),
		CreateTopicCalls:   make([]db.CreateTopicParams, 0),
	}
}

// GetSession returns a session by ID
func (m *MockDBQueries) GetSession(ctx context.Context, params db.GetSessionParams) (db.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.GetSessionError != nil {
		return db.Session{}, m.GetSessionError
	}

	session, exists := m.Sessions[params.ID]
	if !exists {
		return db.Session{}, sql.ErrNoRows
	}
	return session, nil
}

// GetSessionByIdOrSlug returns a session by ID or slug
func (m *MockDBQueries) GetSessionByIdOrSlug(ctx context.Context, params db.GetSessionByIdOrSlugParams) (db.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.GetSessionError != nil {
		return db.Session{}, m.GetSessionError
	}

	// Try by ID first
	if session, exists := m.Sessions[params.ID]; exists {
		return session, nil
	}

	// Try by slug
	for _, session := range m.Sessions {
		if session.Slug == params.Slug {
			return session, nil
		}
	}

	return db.Session{}, sql.ErrNoRows
}

// CreateSession creates a new session
func (m *MockDBQueries) CreateSession(ctx context.Context, params db.CreateSessionParams) (db.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CreateSessionCalls = append(m.CreateSessionCalls, params)

	if m.CreateSessionError != nil {
		return db.Session{}, m.CreateSessionError
	}

	session := db.Session{
		ID:        params.ID,
		Slug:      params.Slug,
		Title:     params.Title,
		UserID:    params.UserID,
		Type:      params.Type,
		CreatedAt: params.CreatedAt,
		UpdatedAt: params.UpdatedAt,
	}
	m.Sessions[params.ID] = session
	return session, nil
}

// ListMessagesBySession returns messages for a session
func (m *MockDBQueries) ListMessagesBySession(ctx context.Context, params db.ListMessagesBySessionParams) ([]db.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if params.SessionID.Valid {
		messages, exists := m.Messages[params.SessionID.String]
		if exists {
			return messages, nil
		}
	}
	return []db.Message{}, nil
}

// CreateMessage creates a new message
func (m *MockDBQueries) CreateMessage(ctx context.Context, params db.CreateMessageParams) (db.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CreateMessageCalls = append(m.CreateMessageCalls, params)

	if m.CreateMessageError != nil {
		return db.Message{}, m.CreateMessageError
	}

	msg := db.Message{
		ID:        params.ID,
		Role:      params.Role,
		Content:   params.Content,
		Tools:     params.Tools,
		SessionID: params.SessionID,
		TopicID:   params.TopicID,
		ThreadID:  params.ThreadID,
		UserID:    params.UserID,
		CreatedAt: params.CreatedAt,
		UpdatedAt: params.UpdatedAt,
	}

	if params.SessionID.Valid {
		m.Messages[params.SessionID.String] = append(m.Messages[params.SessionID.String], msg)
	}

	return msg, nil
}

// GetTopic returns a topic by ID
func (m *MockDBQueries) GetTopic(ctx context.Context, params db.GetTopicParams) (db.Topic, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.GetTopicError != nil {
		return db.Topic{}, m.GetTopicError
	}

	topic, exists := m.Topics[params.ID]
	if !exists {
		return db.Topic{}, sql.ErrNoRows
	}
	return topic, nil
}

// CreateTopic creates a new topic
func (m *MockDBQueries) CreateTopic(ctx context.Context, params db.CreateTopicParams) (db.Topic, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.CreateTopicCalls = append(m.CreateTopicCalls, params)

	if m.CreateTopicError != nil {
		return db.Topic{}, m.CreateTopicError
	}

	topic := db.Topic{
		ID:             params.ID,
		Title:          params.Title,
		SessionID:      params.SessionID,
		UserID:         params.UserID,
		HistorySummary: params.HistorySummary,
		Metadata:       params.Metadata,
		CreatedAt:      params.CreatedAt,
		UpdatedAt:      params.UpdatedAt,
	}
	m.Topics[params.ID] = topic
	return topic, nil
}

// CountTopicsBySession counts topics for a session
func (m *MockDBQueries) CountTopicsBySession(ctx context.Context, params db.CountTopicsBySessionParams) (int64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	count := int64(0)
	for _, topic := range m.Topics {
		if topic.SessionID.Valid && topic.SessionID.String == params.SessionID.String {
			count++
		}
	}
	return count, nil
}

// UpdateTopic updates a topic
func (m *MockDBQueries) UpdateTopic(ctx context.Context, params db.UpdateTopicParams) (db.Topic, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	topic, exists := m.Topics[params.ID]
	if !exists {
		return db.Topic{}, sql.ErrNoRows
	}

	topic.Title = params.Title
	topic.HistorySummary = params.HistorySummary
	topic.Metadata = params.Metadata
	topic.UpdatedAt = params.UpdatedAt
	m.Topics[params.ID] = topic

	return topic, nil
}

// UpdateSession updates a session
func (m *MockDBQueries) UpdateSession(ctx context.Context, params db.UpdateSessionParams) (db.Session, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	session, exists := m.Sessions[params.ID]
	if !exists {
		return db.Session{}, sql.ErrNoRows
	}

	session.Title = params.Title
	session.UpdatedAt = params.UpdatedAt
	m.Sessions[params.ID] = session

	return session, nil
}

// GetMessagesByTopicId returns messages for a topic
func (m *MockDBQueries) GetMessagesByTopicId(ctx context.Context, params db.GetMessagesByTopicIdParams) ([]db.Message, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	var result []db.Message
	for _, messages := range m.Messages {
		for _, msg := range messages {
			if msg.TopicID.Valid && msg.TopicID.String == params.TopicID.String {
				result = append(result, msg)
			}
		}
	}
	return result, nil
}

// AddSession adds a session to the mock (for test setup)
func (m *MockDBQueries) AddSession(id, userID string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now().UnixMilli()
	m.Sessions[id] = db.Session{
		ID:        id,
		Slug:      id[:8],
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AddTopic adds a topic to the mock (for test setup)
func (m *MockDBQueries) AddTopic(id, sessionID, userID string, title string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now().UnixMilli()
	m.Topics[id] = db.Topic{
		ID:        id,
		Title:     sql.NullString{String: title, Valid: true},
		SessionID: sql.NullString{String: sessionID, Valid: true},
		UserID:    userID,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// AddMessage adds a message to the mock (for test setup)
func (m *MockDBQueries) AddMessage(sessionID, role, content string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	now := time.Now().UnixMilli()
	msg := db.Message{
		ID:        fmt.Sprintf("msg-%d", now),
		Role:      role,
		Content:   sql.NullString{String: content, Valid: true},
		SessionID: sql.NullString{String: sessionID, Valid: true},
		CreatedAt: now,
		UpdatedAt: now,
	}
	m.Messages[sessionID] = append(m.Messages[sessionID], msg)
}
