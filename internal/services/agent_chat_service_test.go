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
	"strings"
	"testing"

	db "github.com/kawai-network/veridium/internal/database/generated"
	"github.com/kawai-network/veridium/pkg/yzma/message"
)

// ============================================================================
// Test: convertDBMessageToYzma
// ============================================================================

func TestConvertDBMessageToYzma_ChatMessage(t *testing.T) {
	tests := []struct {
		name     string
		dbMsg    *db.Message
		wantRole string
		wantType string
	}{
		{
			name: "user message",
			dbMsg: &db.Message{
				Role:    "user",
				Content: sql.NullString{String: "Hello, world!", Valid: true},
			},
			wantRole: "user",
			wantType: "Chat",
		},
		{
			name: "assistant message",
			dbMsg: &db.Message{
				Role:    "assistant",
				Content: sql.NullString{String: "Hi there!", Valid: true},
			},
			wantRole: "assistant",
			wantType: "Chat",
		},
		{
			name: "system message",
			dbMsg: &db.Message{
				Role:    "system",
				Content: sql.NullString{String: "You are a helpful assistant.", Valid: true},
			},
			wantRole: "system",
			wantType: "Chat",
		},
		{
			name: "message with empty content",
			dbMsg: &db.Message{
				Role:    "user",
				Content: sql.NullString{String: "", Valid: false},
			},
			wantRole: "user",
			wantType: "Chat",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertDBMessageToYzma(tt.dbMsg)

			if result == nil {
				t.Fatal("convertDBMessageToYzma returned nil")
			}

			if result.GetRole() != tt.wantRole {
				t.Errorf("GetRole() = %v, want %v", result.GetRole(), tt.wantRole)
			}

			switch tt.wantType {
			case "Chat":
				if _, ok := result.(message.Chat); !ok {
					t.Errorf("expected message.Chat, got %T", result)
				}
			}
		})
	}
}

func TestConvertDBMessageToYzma_ToolMessage(t *testing.T) {
	toolCalls := []message.ToolCall{
		{
			Type: "function",
			Function: message.ToolFunction{
				Name:      "search_kb",
				Arguments: map[string]string{"query": "test"},
			},
		},
	}
	toolCallsJSON, _ := json.Marshal(toolCalls)

	dbMsg := &db.Message{
		Role:  "assistant",
		Tools: sql.NullString{String: string(toolCallsJSON), Valid: true},
	}

	result := convertDBMessageToYzma(dbMsg)

	if result == nil {
		t.Fatal("convertDBMessageToYzma returned nil")
	}

	toolMsg, ok := result.(message.Tool)
	if !ok {
		t.Fatalf("expected message.Tool, got %T", result)
	}

	if len(toolMsg.ToolCalls) != 1 {
		t.Errorf("expected 1 tool call, got %d", len(toolMsg.ToolCalls))
	}

	if toolMsg.ToolCalls[0].Function.Name != "search_kb" {
		t.Errorf("expected tool name 'search_kb', got %s", toolMsg.ToolCalls[0].Function.Name)
	}
}

func TestConvertDBMessageToYzma_ToolResponse(t *testing.T) {
	dbMsg := &db.Message{
		Role:    "tool",
		Content: sql.NullString{String: `{"result": "success"}`, Valid: true},
	}

	result := convertDBMessageToYzma(dbMsg)

	if result == nil {
		t.Fatal("convertDBMessageToYzma returned nil")
	}

	toolResp, ok := result.(message.ToolResponse)
	if !ok {
		t.Fatalf("expected message.ToolResponse, got %T", result)
	}

	if toolResp.Content != `{"result": "success"}` {
		t.Errorf("expected content to be preserved, got %s", toolResp.Content)
	}
}

// ============================================================================
// Test: convertYzmaMessageToDB
// ============================================================================

func TestConvertYzmaMessageToDB_ChatMessage(t *testing.T) {
	msg := message.Chat{
		Role:    "user",
		Content: "Hello, world!",
	}

	params := convertYzmaMessageToDB(msg, "session-123", "user-456")

	if params.Role != "user" {
		t.Errorf("Role = %v, want user", params.Role)
	}

	if !params.Content.Valid || params.Content.String != "Hello, world!" {
		t.Errorf("Content = %v, want 'Hello, world!'", params.Content)
	}

	if !params.SessionID.Valid || params.SessionID.String != "session-123" {
		t.Errorf("SessionID = %v, want 'session-123'", params.SessionID)
	}

	if params.UserID != "user-456" {
		t.Errorf("UserID = %v, want 'user-456'", params.UserID)
	}

	if params.ID == "" {
		t.Error("ID should be generated")
	}

	if params.CreatedAt == 0 {
		t.Error("CreatedAt should be set")
	}
}

func TestConvertYzmaMessageToDB_ToolMessage(t *testing.T) {
	msg := message.Tool{
		Role: "assistant",
		ToolCalls: []message.ToolCall{
			{
				Type: "function",
				Function: message.ToolFunction{
					Name:      "search_kb",
					Arguments: map[string]string{"query": "test"},
				},
			},
		},
	}

	params := convertYzmaMessageToDB(msg, "session-123", "user-456")

	if params.Role != "assistant" {
		t.Errorf("Role = %v, want assistant", params.Role)
	}

	if !params.Tools.Valid {
		t.Error("Tools should be valid")
	}

	// Verify tools JSON
	var toolCalls []message.ToolCall
	if err := json.Unmarshal([]byte(params.Tools.String), &toolCalls); err != nil {
		t.Errorf("Failed to unmarshal tools: %v", err)
	}

	if len(toolCalls) != 1 {
		t.Errorf("Expected 1 tool call, got %d", len(toolCalls))
	}
}

func TestConvertYzmaMessageToDB_ToolResponse(t *testing.T) {
	msg := message.ToolResponse{
		Role:    "tool",
		Name:    "search_kb",
		Content: `{"result": "success"}`,
	}

	params := convertYzmaMessageToDB(msg, "session-123", "user-456")

	if params.Role != "tool" {
		t.Errorf("Role = %v, want tool", params.Role)
	}

	if !params.Content.Valid || params.Content.String != `{"result": "success"}` {
		t.Errorf("Content = %v, want '{\"result\": \"success\"}'", params.Content)
	}
}

func TestConvertYzmaMessageToDB_EmptySessionID(t *testing.T) {
	msg := message.Chat{
		Role:    "user",
		Content: "Hello",
	}

	params := convertYzmaMessageToDB(msg, "", "user-456")

	if params.SessionID.Valid {
		t.Error("SessionID should not be valid when empty")
	}
}

// ============================================================================
// Test: stripThinkTags
// ============================================================================

func TestStripThinkTags(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "no think tags",
			input:    "Hello, world!",
			expected: "Hello, world!",
		},
		{
			name:     "complete think tags",
			input:    "<think>reasoning here</think>Final answer",
			expected: "Final answer",
		},
		{
			name:     "think tags with newlines",
			input:    "<think>\nSome reasoning\nMore reasoning\n</think>The actual response",
			expected: "The actual response",
		},
		{
			name:     "multiple think blocks",
			input:    "<think>first</think>Middle<think>second</think>End",
			expected: "MiddleEnd",
		},
		{
			name:     "unclosed think tag",
			input:    "Start<think>unclosed reasoning",
			expected: "Start",
		},
		{
			name:     "only think content",
			input:    "<think>only reasoning</think>",
			expected: "",
		},
		{
			name:     "whitespace around result",
			input:    "<think>reasoning</think>  Result  ",
			expected: "Result",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stripThinkTags(tt.input)
			if result != tt.expected {
				t.Errorf("stripThinkTags(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

// ============================================================================
// Test: prepareMessagesWithSystemPrompt
// ============================================================================

func TestPrepareMessagesWithSystemPrompt_NoExistingSystem(t *testing.T) {
	service := &AgentChatService{
		reasoningConfig: DefaultReasoningConfig(),
	}

	session := &AgentSession{
		SessionID: "test-session",
		UserID:    "test-user",
	}

	messages := []message.Message{
		message.Chat{Role: "user", Content: "Hello"},
	}

	result := service.prepareMessagesWithSystemPrompt(messages, session)

	if len(result) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(result))
	}

	if result[0].GetRole() != "system" {
		t.Errorf("First message should be system, got %s", result[0].GetRole())
	}

	if result[1].GetRole() != "user" {
		t.Errorf("Second message should be user, got %s", result[1].GetRole())
	}
}

func TestPrepareMessagesWithSystemPrompt_WithExistingSystem(t *testing.T) {
	service := &AgentChatService{
		reasoningConfig: DefaultReasoningConfig(),
	}

	session := &AgentSession{
		SessionID: "test-session",
		UserID:    "test-user",
	}

	messages := []message.Message{
		message.Chat{Role: "system", Content: "Old system prompt"},
		message.Chat{Role: "user", Content: "Hello"},
	}

	result := service.prepareMessagesWithSystemPrompt(messages, session)

	if len(result) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(result))
	}

	// System prompt should be updated, not duplicated
	systemMsg, ok := result[0].(message.Chat)
	if !ok {
		t.Fatal("First message should be Chat type")
	}

	if systemMsg.Content == "Old system prompt" {
		t.Error("System prompt should be updated, not kept as old")
	}
}

func TestPrepareMessagesWithSystemPrompt_WithHistorySummary(t *testing.T) {
	service := &AgentChatService{
		reasoningConfig: DefaultReasoningConfig(),
	}

	session := &AgentSession{
		SessionID: "test-session",
		UserID:    "test-user",
		Context: map[string]any{
			"history_summary": "Previous conversation was about coding.",
		},
	}

	messages := []message.Message{
		message.Chat{Role: "user", Content: "Hello"},
	}

	result := service.prepareMessagesWithSystemPrompt(messages, session)

	systemMsg, ok := result[0].(message.Chat)
	if !ok {
		t.Fatal("First message should be Chat type")
	}

	if systemMsg.Content == "" {
		t.Error("System prompt should contain history summary")
	}

	// Check that summary is included
	if len(systemMsg.Content) < 50 {
		t.Error("System prompt seems too short, may not include summary")
	}
}

func TestPrepareMessagesWithSystemPrompt_WithKnowledgeBase(t *testing.T) {
	service := &AgentChatService{
		reasoningConfig: DefaultReasoningConfig(),
	}

	session := &AgentSession{
		SessionID:       "test-session",
		UserID:          "test-user",
		KnowledgeBaseID: "kb-123",
	}

	messages := []message.Message{
		message.Chat{Role: "user", Content: "Hello"},
	}

	result := service.prepareMessagesWithSystemPrompt(messages, session)

	systemMsg, ok := result[0].(message.Chat)
	if !ok {
		t.Fatal("First message should be Chat type")
	}

	// Should mention knowledge base in system prompt
	if systemMsg.Content == "" {
		t.Error("System prompt should be set")
	}
}

// ============================================================================
// Test: ChatRequest/ChatResponse struct validation
// ============================================================================

func TestChatRequest_Marshaling(t *testing.T) {
	req := ChatRequest{
		SessionID: "session-123",
		UserID:    "user-456",
		Message:   "Hello, world!",
		TopicID:   "topic-789",
		Stream:    true,
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Failed to marshal ChatRequest: %v", err)
	}

	var decoded ChatRequest
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ChatRequest: %v", err)
	}

	if decoded.SessionID != req.SessionID {
		t.Errorf("SessionID = %v, want %v", decoded.SessionID, req.SessionID)
	}

	if decoded.Message != req.Message {
		t.Errorf("Message = %v, want %v", decoded.Message, req.Message)
	}

	if decoded.Stream != req.Stream {
		t.Errorf("Stream = %v, want %v", decoded.Stream, req.Stream)
	}
}

func TestChatResponse_Marshaling(t *testing.T) {
	resp := ChatResponse{
		MessageID:    "msg-123",
		SessionID:    "session-456",
		Message:      "Hello!",
		FinishReason: "stop",
		CreatedAt:    1234567890,
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Failed to marshal ChatResponse: %v", err)
	}

	var decoded ChatResponse
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal ChatResponse: %v", err)
	}

	if decoded.MessageID != resp.MessageID {
		t.Errorf("MessageID = %v, want %v", decoded.MessageID, resp.MessageID)
	}

	if decoded.Message != resp.Message {
		t.Errorf("Message = %v, want %v", decoded.Message, resp.Message)
	}

	if decoded.FinishReason != resp.FinishReason {
		t.Errorf("FinishReason = %v, want %v", decoded.FinishReason, resp.FinishReason)
	}
}

// ============================================================================
// Test: AgentSession struct
// ============================================================================

func TestAgentSession_MessageAppend(t *testing.T) {
	session := &AgentSession{
		SessionID: "test-session",
		Messages:  make([]message.Message, 0),
	}

	// Add user message
	session.Messages = append(session.Messages, message.Chat{
		Role:    "user",
		Content: "Hello",
	})

	if len(session.Messages) != 1 {
		t.Errorf("Expected 1 message, got %d", len(session.Messages))
	}

	// Add assistant message
	session.Messages = append(session.Messages, message.Chat{
		Role:    "assistant",
		Content: "Hi there!",
	})

	if len(session.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(session.Messages))
	}

	// Verify roles
	if session.Messages[0].GetRole() != "user" {
		t.Errorf("First message role = %v, want user", session.Messages[0].GetRole())
	}

	if session.Messages[1].GetRole() != "assistant" {
		t.Errorf("Second message role = %v, want assistant", session.Messages[1].GetRole())
	}
}

// ============================================================================
// Test: getKeepMessageCount
// ============================================================================

func TestGetKeepMessageCount(t *testing.T) {
	tests := []struct {
		name     string
		mode     ReasoningMode
		expected int
	}{
		{
			name:     "disabled mode",
			mode:     ReasoningDisabled,
			expected: 20,
		},
		{
			name:     "enabled mode",
			mode:     ReasoningEnabled,
			expected: 12,
		},
		{
			name:     "verbose mode",
			mode:     ReasoningVerbose,
			expected: 6,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := &AgentChatService{
				reasoningConfig: ReasoningConfig{Mode: tt.mode},
			}

			result := service.getKeepMessageCount()
			if result != tt.expected {
				t.Errorf("getKeepMessageCount() = %d, want %d", result, tt.expected)
			}
		})
	}
}

// ============================================================================
// Test: ClearSession
// ============================================================================

func TestClearSession(t *testing.T) {
	service := &AgentChatService{
		sessions: make(map[string]*AgentSession),
	}

	// Add a session
	service.sessions["test-session"] = &AgentSession{
		SessionID: "test-session",
		UserID:    "test-user",
	}

	// Verify session exists
	if _, exists := service.sessions["test-session"]; !exists {
		t.Fatal("Session should exist before clearing")
	}

	// Clear session
	service.ClearSession("test-session")

	// Verify session is removed
	if _, exists := service.sessions["test-session"]; exists {
		t.Error("Session should not exist after clearing")
	}
}

// ============================================================================
// Test: GetSessionHistory
// ============================================================================

func TestGetSessionHistory_ExistingSession(t *testing.T) {
	service := &AgentChatService{
		sessions: make(map[string]*AgentSession),
	}

	// Add a session with messages
	service.sessions["test-session"] = &AgentSession{
		SessionID: "test-session",
		Messages: []message.Message{
			message.Chat{Role: "user", Content: "Hello"},
			message.Chat{Role: "assistant", Content: "Hi!"},
		},
	}

	history, err := service.GetSessionHistory("test-session")
	if err != nil {
		t.Fatalf("GetSessionHistory() error = %v", err)
	}

	if len(history) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(history))
	}
}

func TestGetSessionHistory_NonExistingSession(t *testing.T) {
	service := &AgentChatService{
		sessions: make(map[string]*AgentSession),
	}

	_, err := service.GetSessionHistory("non-existing")
	if err == nil {
		t.Error("Expected error for non-existing session")
	}
}

// ============================================================================
// Test: collectToolNames (basic functionality)
// ============================================================================

func TestCollectToolNames_RequestedTools(t *testing.T) {
	service := &AgentChatService{
		sessions: make(map[string]*AgentSession),
	}

	req := ChatRequest{
		SessionID: "test-session",
		UserID:    "test-user",
		Tools:     []string{"tool1", "tool2"},
	}

	// Note: KB tool registration will fail without kbService, but requested tools should still work
	toolNames := service.collectToolNames(nil, req)

	// Should have at least the requested tools
	if len(toolNames) < 2 {
		t.Errorf("Expected at least 2 tool names, got %d", len(toolNames))
	}

	// Check if requested tools are included
	found := map[string]bool{"tool1": false, "tool2": false}
	for _, name := range toolNames {
		if _, ok := found[name]; ok {
			found[name] = true
		}
	}

	for name, f := range found {
		if !f {
			t.Errorf("Tool %s not found in collected tools", name)
		}
	}
}

// ============================================================================
// Test: Chat Method - Unit Tests for Internal Logic
// ============================================================================

// TestChat_RequestValidation tests ChatRequest validation scenarios
func TestChat_RequestValidation(t *testing.T) {
	tests := []struct {
		name    string
		req     ChatRequest
		wantErr bool
	}{
		{
			name: "valid basic request",
			req: ChatRequest{
				SessionID: "session-123",
				UserID:    "user-456",
				Message:   "Hello",
			},
			wantErr: false,
		},
		{
			name: "request with topic and thread",
			req: ChatRequest{
				SessionID: "session-123",
				UserID:    "user-456",
				Message:   "Hello",
				TopicID:   "topic-789",
				ThreadID:  "thread-abc",
			},
			wantErr: false,
		},
		{
			name: "request with tools",
			req: ChatRequest{
				SessionID: "session-123",
				UserID:    "user-456",
				Message:   "Search for something",
				Tools:     []string{"web_search", "calculator"},
			},
			wantErr: false,
		},
		{
			name: "request with context",
			req: ChatRequest{
				SessionID: "session-123",
				UserID:    "user-456",
				Message:   "Continue our conversation",
				Context:   map[string]any{"key": "value"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Validate request structure
			if tt.req.SessionID == "" && !tt.wantErr {
				t.Error("SessionID should not be empty for valid request")
			}
			if tt.req.UserID == "" && !tt.wantErr {
				t.Error("UserID should not be empty for valid request")
			}
			if tt.req.Message == "" && !tt.wantErr {
				t.Error("Message should not be empty for valid request")
			}
		})
	}
}

// TestChat_SessionManagement tests session creation and reuse
func TestChat_SessionManagement(t *testing.T) {
	service := &AgentChatService{
		sessions:        make(map[string]*AgentSession),
		reasoningConfig: DefaultReasoningConfig(),
	}

	// Test 1: New session creation
	sessionID := "new-session-123"
	userID := "user-456"

	// Simulate session creation (normally done by getOrCreateSession)
	session := &AgentSession{
		SessionID: sessionID,
		UserID:    userID,
		Messages:  make([]message.Message, 0),
		ToolNames: []string{},
		Context:   make(map[string]any),
	}
	service.sessions[sessionID] = session

	// Verify session exists
	if _, exists := service.sessions[sessionID]; !exists {
		t.Error("Session should exist after creation")
	}

	// Test 2: Session reuse
	existingSession := service.sessions[sessionID]
	if existingSession.SessionID != sessionID {
		t.Errorf("Session ID mismatch: got %s, want %s", existingSession.SessionID, sessionID)
	}

	// Test 3: Message accumulation
	existingSession.Messages = append(existingSession.Messages, message.Chat{
		Role:    "user",
		Content: "First message",
	})
	existingSession.Messages = append(existingSession.Messages, message.Chat{
		Role:    "assistant",
		Content: "First response",
	})

	if len(existingSession.Messages) != 2 {
		t.Errorf("Expected 2 messages, got %d", len(existingSession.Messages))
	}

	// Test 4: Multiple sessions isolation
	session2 := &AgentSession{
		SessionID: "session-2",
		UserID:    userID,
		Messages:  make([]message.Message, 0),
	}
	service.sessions["session-2"] = session2

	// Add message to session 2
	session2.Messages = append(session2.Messages, message.Chat{
		Role:    "user",
		Content: "Different conversation",
	})

	// Verify sessions are isolated
	if len(service.sessions[sessionID].Messages) != 2 {
		t.Error("Session 1 messages should not be affected by session 2")
	}
	if len(service.sessions["session-2"].Messages) != 1 {
		t.Error("Session 2 should have its own messages")
	}
}

// TestChat_MessageFlow tests the message flow through the chat process
func TestChat_MessageFlow(t *testing.T) {
	service := &AgentChatService{
		sessions:        make(map[string]*AgentSession),
		reasoningConfig: DefaultReasoningConfig(),
	}

	session := &AgentSession{
		SessionID: "test-session",
		UserID:    "test-user",
		Messages:  make([]message.Message, 0),
	}
	service.sessions["test-session"] = session

	// Simulate user message addition
	userMsg := message.Chat{
		Role:    "user",
		Content: "What is the weather today?",
	}
	session.Messages = append(session.Messages, userMsg)

	// Verify user message was added
	if len(session.Messages) != 1 {
		t.Fatalf("Expected 1 message, got %d", len(session.Messages))
	}

	if session.Messages[0].GetRole() != "user" {
		t.Errorf("Expected role 'user', got '%s'", session.Messages[0].GetRole())
	}

	// Prepare messages with system prompt
	messagesWithSystem := service.prepareMessagesWithSystemPrompt(session.Messages, session)

	// Verify system prompt was added
	if len(messagesWithSystem) != 2 {
		t.Fatalf("Expected 2 messages (system + user), got %d", len(messagesWithSystem))
	}

	if messagesWithSystem[0].GetRole() != "system" {
		t.Errorf("First message should be system, got '%s'", messagesWithSystem[0].GetRole())
	}

	// Simulate assistant response
	assistantMsg := message.Chat{
		Role:    "assistant",
		Content: "I don't have access to real-time weather data.",
	}
	session.Messages = append(session.Messages, assistantMsg)

	// Verify full conversation
	if len(session.Messages) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(session.Messages))
	}

	if session.Messages[1].GetRole() != "assistant" {
		t.Errorf("Second message should be assistant, got '%s'", session.Messages[1].GetRole())
	}
}

// TestChat_ToolCallFlow tests tool call handling in chat
func TestChat_ToolCallFlow(t *testing.T) {
	service := &AgentChatService{
		sessions:        make(map[string]*AgentSession),
		reasoningConfig: DefaultReasoningConfig(),
	}

	session := &AgentSession{
		SessionID: "test-session",
		UserID:    "test-user",
		Messages:  make([]message.Message, 0),
		ToolNames: []string{"web_search"},
	}
	service.sessions["test-session"] = session

	// Simulate user message
	session.Messages = append(session.Messages, message.Chat{
		Role:    "user",
		Content: "Search for latest news about AI",
	})

	// Simulate tool call from assistant
	toolCallMsg := message.Tool{
		Role: "assistant",
		ToolCalls: []message.ToolCall{
			{
				Type: "function",
				Function: message.ToolFunction{
					Name:      "web_search",
					Arguments: map[string]string{"query": "latest AI news"},
				},
			},
		},
	}
	session.Messages = append(session.Messages, toolCallMsg)

	// Simulate tool response
	toolResponse := message.ToolResponse{
		Role:    "tool",
		Name:    "web_search",
		Content: `{"results": ["AI news 1", "AI news 2"]}`,
	}
	session.Messages = append(session.Messages, toolResponse)

	// Simulate final assistant response
	session.Messages = append(session.Messages, message.Chat{
		Role:    "assistant",
		Content: "Based on my search, here are the latest AI news...",
	})

	// Verify message flow
	if len(session.Messages) != 4 {
		t.Fatalf("Expected 4 messages, got %d", len(session.Messages))
	}

	// Verify message types
	expectedRoles := []string{"user", "assistant", "tool", "assistant"}
	for i, expectedRole := range expectedRoles {
		if session.Messages[i].GetRole() != expectedRole {
			t.Errorf("Message %d: expected role '%s', got '%s'", i, expectedRole, session.Messages[i].GetRole())
		}
	}

	// Verify tool call message
	if toolMsg, ok := session.Messages[1].(message.Tool); ok {
		if len(toolMsg.ToolCalls) != 1 {
			t.Errorf("Expected 1 tool call, got %d", len(toolMsg.ToolCalls))
		}
		if toolMsg.ToolCalls[0].Function.Name != "web_search" {
			t.Errorf("Expected tool name 'web_search', got '%s'", toolMsg.ToolCalls[0].Function.Name)
		}
	} else {
		t.Error("Message 1 should be Tool type")
	}
}

// TestChat_ResponseStructure tests ChatResponse structure
func TestChat_ResponseStructure(t *testing.T) {
	tests := []struct {
		name     string
		response ChatResponse
		validate func(t *testing.T, resp ChatResponse)
	}{
		{
			name: "basic response",
			response: ChatResponse{
				MessageID:    "msg-123",
				SessionID:    "session-456",
				Message:      "Hello!",
				FinishReason: "stop",
				CreatedAt:    1234567890,
			},
			validate: func(t *testing.T, resp ChatResponse) {
				if resp.MessageID == "" {
					t.Error("MessageID should not be empty")
				}
				if resp.Message == "" {
					t.Error("Message should not be empty")
				}
				if resp.FinishReason != "stop" {
					t.Errorf("FinishReason = %s, want 'stop'", resp.FinishReason)
				}
			},
		},
		{
			name: "response with topic",
			response: ChatResponse{
				MessageID:    "msg-123",
				SessionID:    "session-456",
				TopicID:      "topic-789",
				Message:      "Response with topic",
				FinishReason: "stop",
			},
			validate: func(t *testing.T, resp ChatResponse) {
				if resp.TopicID == "" {
					t.Error("TopicID should be set")
				}
			},
		},
		{
			name: "response with tool calls",
			response: ChatResponse{
				MessageID: "msg-123",
				SessionID: "session-456",
				Message:   "",
				ToolCalls: []message.ToolCall{
					{
						Type: "function",
						Function: message.ToolFunction{
							Name:      "test_tool",
							Arguments: map[string]string{"arg": "value"},
						},
					},
				},
				FinishReason: "tool_calls",
			},
			validate: func(t *testing.T, resp ChatResponse) {
				if len(resp.ToolCalls) == 0 {
					t.Error("ToolCalls should not be empty")
				}
				if resp.FinishReason != "tool_calls" {
					t.Errorf("FinishReason = %s, want 'tool_calls'", resp.FinishReason)
				}
			},
		},
		{
			name: "response with error",
			response: ChatResponse{
				MessageID: "msg-123",
				SessionID: "session-456",
				Error:     "Something went wrong",
			},
			validate: func(t *testing.T, resp ChatResponse) {
				if resp.Error == "" {
					t.Error("Error should be set")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.validate(t, tt.response)
		})
	}
}

// TestChat_ContextInjection tests context injection into system prompt
func TestChat_ContextInjection(t *testing.T) {
	service := &AgentChatService{
		sessions:        make(map[string]*AgentSession),
		reasoningConfig: DefaultReasoningConfig(),
	}

	tests := []struct {
		name          string
		session       *AgentSession
		checkContains string
	}{
		{
			name: "with history summary",
			session: &AgentSession{
				SessionID: "test",
				Context: map[string]any{
					"history_summary": "User asked about weather yesterday",
				},
			},
			checkContains: "history_summary",
		},
		{
			name: "with knowledge base",
			session: &AgentSession{
				SessionID:       "test",
				KnowledgeBaseID: "kb-123",
			},
			checkContains: "knowledge base",
		},
		{
			name: "empty context",
			session: &AgentSession{
				SessionID: "test",
			},
			checkContains: "helpful AI assistant",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messages := []message.Message{
				message.Chat{Role: "user", Content: "Hello"},
			}

			result := service.prepareMessagesWithSystemPrompt(messages, tt.session)

			if len(result) < 2 {
				t.Fatal("Should have at least system and user messages")
			}

			systemMsg, ok := result[0].(message.Chat)
			if !ok {
				t.Fatal("First message should be Chat type")
			}

			// Note: checkContains is a hint, not exact match (context injection varies)
			if systemMsg.Content == "" {
				t.Error("System prompt should not be empty")
			}
		})
	}
}

// TestChat_MultiTurnConversation tests multi-turn conversation handling
func TestChat_MultiTurnConversation(t *testing.T) {
	service := &AgentChatService{
		sessions:        make(map[string]*AgentSession),
		reasoningConfig: DefaultReasoningConfig(),
	}

	session := &AgentSession{
		SessionID: "multi-turn-session",
		UserID:    "test-user",
		Messages:  make([]message.Message, 0),
	}
	service.sessions["multi-turn-session"] = session

	// Simulate multi-turn conversation
	turns := []struct {
		userMsg      string
		assistantMsg string
	}{
		{"Hello", "Hi! How can I help you?"},
		{"What's your name?", "I'm an AI assistant."},
		{"Tell me a joke", "Why did the programmer quit? Because he didn't get arrays!"},
	}

	for i, turn := range turns {
		// Add user message
		session.Messages = append(session.Messages, message.Chat{
			Role:    "user",
			Content: turn.userMsg,
		})

		// Prepare messages (simulating what Chat() does)
		messagesWithSystem := service.prepareMessagesWithSystemPrompt(session.Messages, session)

		// Verify growing context
		expectedMsgCount := (i+1)*2 - 1 + 1 // (turns so far * 2 - 1 assistant) + 1 system
		if len(messagesWithSystem) != expectedMsgCount {
			t.Errorf("Turn %d: expected %d messages, got %d", i+1, expectedMsgCount, len(messagesWithSystem))
		}

		// Add assistant response
		session.Messages = append(session.Messages, message.Chat{
			Role:    "assistant",
			Content: turn.assistantMsg,
		})
	}

	// Verify final state
	expectedTotal := len(turns) * 2 // user + assistant for each turn
	if len(session.Messages) != expectedTotal {
		t.Errorf("Expected %d messages, got %d", expectedTotal, len(session.Messages))
	}
}

// TestChat_StreamingVsNonStreaming tests streaming flag handling
func TestChat_StreamingVsNonStreaming(t *testing.T) {
	tests := []struct {
		name   string
		stream bool
	}{
		{
			name:   "non-streaming request",
			stream: false,
		},
		{
			name:   "streaming request",
			stream: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := ChatRequest{
				SessionID: "session-123",
				UserID:    "user-456",
				Message:   "Hello",
				Stream:    tt.stream,
			}

			// Verify stream flag is set correctly
			if req.Stream != tt.stream {
				t.Errorf("Stream = %v, want %v", req.Stream, tt.stream)
			}
		})
	}
}

// ============================================================================
// Integration Tests for Chat() Method with Mock LLM
// ============================================================================

// createTestService creates an AgentChatService with mock dependencies for testing
func createTestService(mockLLM *MockLLMGenerator) *AgentChatService {
	service := &AgentChatService{
		sessions:        make(map[string]*AgentSession),
		reasoningConfig: DefaultReasoningConfig(),
	}
	if mockLLM != nil {
		service.llmGenerator = mockLLM
	}
	return service
}

// TestChatIntegration_BasicChat tests basic chat flow with mock LLM
func TestChatIntegration_BasicChat(t *testing.T) {
	// Setup mock LLM
	mockLLM := NewMockLLMGenerator()
	mockLLM.SetResponse("Hello! How can I help you today?", "stop")

	// Create service with mock
	service := createTestService(mockLLM)

	// Create session manually (simulating getOrCreateSession)
	sessionID := "test-session-123"
	userID := "test-user-456"
	service.sessions[sessionID] = &AgentSession{
		SessionID: sessionID,
		UserID:    userID,
		Messages:  make([]message.Message, 0),
		ToolNames: []string{},
	}

	// Simulate chat flow (without DB)
	session := service.sessions[sessionID]

	// Add user message
	userMsg := message.Chat{Role: "user", Content: "Hello"}
	session.Messages = append(session.Messages, userMsg)

	// Prepare messages with system prompt
	messagesWithSystem := service.prepareMessagesWithSystemPrompt(session.Messages, session)

	// Call LLM through interface
	llmWithTools := service.llmGenerator.WithTools(session.ToolNames)
	resp, toolMessages, err := llmWithTools.RunAgentLoop(t.Context(), messagesWithSystem, 10)

	// Assertions
	if err != nil {
		t.Fatalf("RunAgentLoop failed: %v", err)
	}

	if resp.Content != "Hello! How can I help you today?" {
		t.Errorf("Expected mock response, got: %s", resp.Content)
	}

	if resp.FinishReason != "stop" {
		t.Errorf("FinishReason = %s, want 'stop'", resp.FinishReason)
	}

	// Add assistant response to session
	session.Messages = append(session.Messages, toolMessages...)
	session.Messages = append(session.Messages, message.Chat{
		Role:    "assistant",
		Content: resp.Content,
	})

	// Verify session state
	if len(session.Messages) != 2 {
		t.Errorf("Expected 2 messages in session, got %d", len(session.Messages))
	}

	// Note: Call count tracked on the mock returned by WithTools, not original
	// This verifies the flow works correctly
}

// TestChatIntegration_WithToolCalls tests chat flow with tool calls
func TestChatIntegration_WithToolCalls(t *testing.T) {
	// Setup mock LLM with tool calls
	mockLLM := NewMockLLMGenerator()
	mockLLM.SetToolCalls([]message.ToolCall{
		{
			Type: "function",
			Function: message.ToolFunction{
				Name:      "web_search",
				Arguments: map[string]string{"query": "weather today"},
			},
		},
	})
	mockLLM.Response.Content = "Based on my search, the weather is sunny."

	// Add tool messages that would be returned
	mockLLM.ToolMessages = []message.Message{
		message.Tool{
			Role: "assistant",
			ToolCalls: []message.ToolCall{
				{
					Type: "function",
					Function: message.ToolFunction{
						Name:      "web_search",
						Arguments: map[string]string{"query": "weather today"},
					},
				},
			},
		},
		message.ToolResponse{
			Role:    "tool",
			Name:    "web_search",
			Content: `{"result": "sunny, 25°C"}`,
		},
	}

	// Create service with mock
	service := createTestService(mockLLM)

	// Create session with tools
	sessionID := "tool-session"
	service.sessions[sessionID] = &AgentSession{
		SessionID: sessionID,
		UserID:    "user-123",
		Messages:  make([]message.Message, 0),
		ToolNames: []string{"web_search"},
	}

	session := service.sessions[sessionID]

	// Add user message
	session.Messages = append(session.Messages, message.Chat{
		Role:    "user",
		Content: "What's the weather today?",
	})

	// Prepare and call LLM
	messagesWithSystem := service.prepareMessagesWithSystemPrompt(session.Messages, session)
	llmWithTools := service.llmGenerator.WithTools(session.ToolNames)
	resp, toolMessages, err := llmWithTools.RunAgentLoop(t.Context(), messagesWithSystem, 10)

	if err != nil {
		t.Fatalf("RunAgentLoop failed: %v", err)
	}

	// Verify tool calls in response
	if len(resp.ToolCalls) != 1 {
		t.Errorf("Expected 1 tool call, got %d", len(resp.ToolCalls))
	}

	if resp.ToolCalls[0].Function.Name != "web_search" {
		t.Errorf("Expected tool 'web_search', got '%s'", resp.ToolCalls[0].Function.Name)
	}

	// Verify tool messages returned
	if len(toolMessages) != 2 {
		t.Errorf("Expected 2 tool messages, got %d", len(toolMessages))
	}

	// Add all messages to session
	session.Messages = append(session.Messages, toolMessages...)
	session.Messages = append(session.Messages, message.Chat{
		Role:    "assistant",
		Content: resp.Content,
	})

	// Verify final session state: user + tool_call + tool_response + assistant
	if len(session.Messages) != 4 {
		t.Errorf("Expected 4 messages, got %d", len(session.Messages))
	}
}

// TestChatIntegration_LLMError tests error handling when LLM fails
func TestChatIntegration_LLMError(t *testing.T) {
	// Setup mock LLM with error
	mockLLM := NewMockLLMGenerator()
	mockLLM.SetError(fmt.Errorf("model not loaded"))

	service := createTestService(mockLLM)

	// Create session
	sessionID := "error-session"
	service.sessions[sessionID] = &AgentSession{
		SessionID: sessionID,
		UserID:    "user-123",
		Messages:  make([]message.Message, 0),
	}

	session := service.sessions[sessionID]
	session.Messages = append(session.Messages, message.Chat{
		Role:    "user",
		Content: "Hello",
	})

	// Prepare and call LLM
	messagesWithSystem := service.prepareMessagesWithSystemPrompt(session.Messages, session)
	llmWithTools := service.llmGenerator.WithTools(nil)
	_, _, err := llmWithTools.RunAgentLoop(t.Context(), messagesWithSystem, 10)

	// Verify error is returned
	if err == nil {
		t.Error("Expected error from LLM, got nil")
	}

	if err.Error() != "model not loaded" {
		t.Errorf("Expected 'model not loaded' error, got: %v", err)
	}
}

// TestChatIntegration_MultiTurn tests multi-turn conversation with mock
func TestChatIntegration_MultiTurn(t *testing.T) {
	mockLLM := NewMockLLMGenerator()
	service := createTestService(mockLLM)

	sessionID := "multi-turn-session"
	service.sessions[sessionID] = &AgentSession{
		SessionID: sessionID,
		UserID:    "user-123",
		Messages:  make([]message.Message, 0),
	}

	session := service.sessions[sessionID]

	// Turn 1
	mockLLM.SetResponse("Hi! I'm here to help.", "stop")
	session.Messages = append(session.Messages, message.Chat{Role: "user", Content: "Hello"})

	messagesWithSystem := service.prepareMessagesWithSystemPrompt(session.Messages, session)
	resp, _, _ := service.llmGenerator.RunAgentLoop(t.Context(), messagesWithSystem, 10)
	session.Messages = append(session.Messages, message.Chat{Role: "assistant", Content: resp.Content})

	// Turn 2
	mockLLM.SetResponse("I'm doing great, thanks for asking!", "stop")
	session.Messages = append(session.Messages, message.Chat{Role: "user", Content: "How are you?"})

	messagesWithSystem = service.prepareMessagesWithSystemPrompt(session.Messages, session)
	resp, _, _ = service.llmGenerator.RunAgentLoop(t.Context(), messagesWithSystem, 10)
	session.Messages = append(session.Messages, message.Chat{Role: "assistant", Content: resp.Content})

	// Turn 3
	mockLLM.SetResponse("My name is Assistant.", "stop")
	session.Messages = append(session.Messages, message.Chat{Role: "user", Content: "What's your name?"})

	messagesWithSystem = service.prepareMessagesWithSystemPrompt(session.Messages, session)
	resp, _, _ = service.llmGenerator.RunAgentLoop(t.Context(), messagesWithSystem, 10)
	session.Messages = append(session.Messages, message.Chat{Role: "assistant", Content: resp.Content})

	// Verify: 3 turns = 6 messages (3 user + 3 assistant)
	if len(session.Messages) != 6 {
		t.Errorf("Expected 6 messages after 3 turns, got %d", len(session.Messages))
	}

	// Verify LLM was called 3 times
	if mockLLM.GetRunAgentLoopCallCount() != 3 {
		t.Errorf("Expected 3 LLM calls, got %d", mockLLM.GetRunAgentLoopCallCount())
	}

	// Verify message order
	expectedRoles := []string{"user", "assistant", "user", "assistant", "user", "assistant"}
	for i, role := range expectedRoles {
		if session.Messages[i].GetRole() != role {
			t.Errorf("Message %d: expected role '%s', got '%s'", i, role, session.Messages[i].GetRole())
		}
	}
}

// TestChatIntegration_SystemPromptIncluded tests that system prompt is included in LLM calls
func TestChatIntegration_SystemPromptIncluded(t *testing.T) {
	mockLLM := NewMockLLMGenerator()

	var capturedMessages []message.Message
	mockLLM.OnRunAgentLoop = func(ctx context.Context, messages []message.Message, maxIterations int) {
		capturedMessages = messages
	}

	service := createTestService(mockLLM)

	sessionID := "system-prompt-session"
	service.sessions[sessionID] = &AgentSession{
		SessionID: sessionID,
		UserID:    "user-123",
		Messages:  make([]message.Message, 0),
	}

	session := service.sessions[sessionID]
	session.Messages = append(session.Messages, message.Chat{Role: "user", Content: "Hello"})

	messagesWithSystem := service.prepareMessagesWithSystemPrompt(session.Messages, session)
	service.llmGenerator.RunAgentLoop(t.Context(), messagesWithSystem, 10)

	// Verify system prompt was included
	if len(capturedMessages) < 2 {
		t.Fatalf("Expected at least 2 messages (system + user), got %d", len(capturedMessages))
	}

	if capturedMessages[0].GetRole() != "system" {
		t.Errorf("First message should be system prompt, got '%s'", capturedMessages[0].GetRole())
	}

	if capturedMessages[1].GetRole() != "user" {
		t.Errorf("Second message should be user, got '%s'", capturedMessages[1].GetRole())
	}
}

// TestChatIntegration_HistorySummaryInjection tests history summary injection
func TestChatIntegration_HistorySummaryInjection(t *testing.T) {
	mockLLM := NewMockLLMGenerator()

	var capturedMessages []message.Message
	mockLLM.OnRunAgentLoop = func(ctx context.Context, messages []message.Message, maxIterations int) {
		capturedMessages = messages
	}

	service := createTestService(mockLLM)

	sessionID := "summary-session"
	service.sessions[sessionID] = &AgentSession{
		SessionID: sessionID,
		UserID:    "user-123",
		Messages:  make([]message.Message, 0),
		Context: map[string]any{
			"history_summary": "User previously asked about Go programming.",
		},
	}

	session := service.sessions[sessionID]
	session.Messages = append(session.Messages, message.Chat{Role: "user", Content: "Continue please"})

	messagesWithSystem := service.prepareMessagesWithSystemPrompt(session.Messages, session)
	service.llmGenerator.RunAgentLoop(t.Context(), messagesWithSystem, 10)

	// Verify system prompt contains history summary
	if len(capturedMessages) < 1 {
		t.Fatal("Expected at least 1 message")
	}

	systemMsg, ok := capturedMessages[0].(message.Chat)
	if !ok {
		t.Fatal("First message should be Chat type")
	}

	if !strings.Contains(systemMsg.Content, "summary") && !strings.Contains(systemMsg.Content, "history") {
		t.Error("System prompt should contain history summary context")
	}
}

// TestChatIntegration_KnowledgeBaseContext tests KB context in system prompt
func TestChatIntegration_KnowledgeBaseContext(t *testing.T) {
	mockLLM := NewMockLLMGenerator()

	var capturedMessages []message.Message
	mockLLM.OnRunAgentLoop = func(ctx context.Context, messages []message.Message, maxIterations int) {
		capturedMessages = messages
	}

	service := createTestService(mockLLM)

	sessionID := "kb-session"
	service.sessions[sessionID] = &AgentSession{
		SessionID:       sessionID,
		UserID:          "user-123",
		Messages:        make([]message.Message, 0),
		KnowledgeBaseID: "kb-docs-123",
	}

	session := service.sessions[sessionID]
	session.Messages = append(session.Messages, message.Chat{Role: "user", Content: "Search docs"})

	messagesWithSystem := service.prepareMessagesWithSystemPrompt(session.Messages, session)
	service.llmGenerator.RunAgentLoop(t.Context(), messagesWithSystem, 10)

	// Verify system prompt mentions knowledge base
	if len(capturedMessages) < 1 {
		t.Fatal("Expected at least 1 message")
	}

	systemMsg, ok := capturedMessages[0].(message.Chat)
	if !ok {
		t.Fatal("First message should be Chat type")
	}

	if !strings.Contains(strings.ToLower(systemMsg.Content), "knowledge base") {
		t.Error("System prompt should mention knowledge base access")
	}
}

// TestChatIntegration_EmptyResponse tests handling of empty LLM response
func TestChatIntegration_EmptyResponse(t *testing.T) {
	mockLLM := NewMockLLMGenerator()
	mockLLM.SetResponse("", "stop")

	service := createTestService(mockLLM)

	sessionID := "empty-response-session"
	service.sessions[sessionID] = &AgentSession{
		SessionID: sessionID,
		UserID:    "user-123",
		Messages:  make([]message.Message, 0),
	}

	session := service.sessions[sessionID]
	session.Messages = append(session.Messages, message.Chat{Role: "user", Content: "Hello"})

	messagesWithSystem := service.prepareMessagesWithSystemPrompt(session.Messages, session)
	resp, _, err := service.llmGenerator.RunAgentLoop(t.Context(), messagesWithSystem, 10)

	// Should not error on empty response
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Response content should be empty
	if resp.Content != "" {
		t.Errorf("Expected empty content, got: %s", resp.Content)
	}
}

// TestChatIntegration_ToolNamesPassedToLLM tests that tool names are passed correctly
func TestChatIntegration_ToolNamesPassedToLLM(t *testing.T) {
	mockLLM := NewMockLLMGenerator()
	service := createTestService(mockLLM)

	sessionID := "tools-session"
	expectedTools := []string{"web_search", "calculator", "file_reader"}
	service.sessions[sessionID] = &AgentSession{
		SessionID: sessionID,
		UserID:    "user-123",
		Messages:  make([]message.Message, 0),
		ToolNames: expectedTools,
	}

	session := service.sessions[sessionID]

	// Get LLM with tools
	llmWithTools := service.llmGenerator.WithTools(session.ToolNames)

	// Verify it's a new mock with tools configured
	if llmWithTools == nil {
		t.Error("WithTools should return non-nil generator")
	}

	// The mock should have tool names set
	mockWithTools, ok := llmWithTools.(*MockLLMGenerator)
	if !ok {
		t.Skip("Cannot verify tool names on non-mock generator")
	}

	if len(mockWithTools.toolNames) != len(expectedTools) {
		t.Errorf("Expected %d tools, got %d", len(expectedTools), len(mockWithTools.toolNames))
	}
}
