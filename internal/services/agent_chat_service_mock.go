package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	db "github.com/kawai-network/veridium/internal/database/generated"
)

// ChatMock handles mock chat responses for testing UI flow without real AI backend
// This method saves complete mock messages to DB with all UI components:
// - Reasoning (step-by-step thinking)
// - RAG Chunks (file chunks with similarity scores)
// - Tool Calls (web browsing, local file system)
// - Tool Messages (separate role='tool' messages with results)
// - Search Grounding (citations from web search)
// - Images (placeholder images)
// - Usage & Performance metrics
//
// Usage from frontend:
//
//	const response = await AgentChatService.ChatMock(request);
//
// The frontend can then fetch messages from DB normally using internal_fetchMessages
func (s *AgentChatService) ChatMock(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	log.Printf("🎭 [MOCK] Starting mock chat for session: %s", req.SessionID)

	// Simulate processing delay
	time.Sleep(500 * time.Millisecond)

	// Get current timestamp
	now := time.Now().UnixMilli()

	// Get or create session
	session, err := s.getOrCreateSession(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get/create session: %w", err)
	}

	// Auto-create topic if needed
	currentTopicID := req.TopicID
	if currentTopicID == "" {
		topicID, err := s.createTopicForSessionSync(ctx, session.SessionID, session.UserID)
		if err != nil {
			log.Printf("⚠️  Warning: Failed to create topic: %v", err)
		} else {
			currentTopicID = topicID
			session.TopicID = topicID
			log.Printf("📝 Auto-created topic: %s", topicID)
		}
	}

	// 1. Save user message
	userMsgID := uuid.New().String()
	userParams := db.CreateMessageParams{
		ID:        userMsgID,
		Role:      "user",
		Content:   sql.NullString{String: req.Message, Valid: true},
		SessionID: sql.NullString{String: req.SessionID, Valid: true},
		UserID:    req.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if currentTopicID != "" {
		userParams.TopicID = sql.NullString{String: currentTopicID, Valid: true}
	}
	if req.ThreadID != "" {
		userParams.ThreadID = sql.NullString{String: req.ThreadID, Valid: true}
	}

	_, err = s.db.Queries().CreateMessage(ctx, userParams)
	if err != nil {
		return nil, fmt.Errorf("failed to save user message: %w", err)
	}
	log.Printf("💾 Saved mock user message: %s", userMsgID)

	// 2. Create assistant message with full mock data
	assistantMsgID := uuid.New().String()
	mockContent := fmt.Sprintf(
		"This is a mock response to: \"%s\"\n\nI'm simulating the AI response to test the UI flow without calling the backend.",
		req.Message,
	)

	// Mock reasoning
	reasoning := map[string]interface{}{
		"content": "Let me think about this step by step:\n1. First, I need to understand the question\n2. Then, I will formulate a response\n3. Finally, I will provide a clear answer",
		"status":  "complete",
	}
	reasoningJSON, _ := json.Marshal(reasoning)

	// Mock RAG chunks
	chunksList := []map[string]interface{}{
		{
			"id":         "chunk_1",
			"fileId":     "file_1",
			"filename":   "document.pdf",
			"fileType":   "application/pdf",
			"fileUrl":    "/files/document.pdf",
			"text":       "This is a sample chunk from the knowledge base. It contains relevant information about the topic.",
			"similarity": 0.95,
		},
		{
			"id":         "chunk_2",
			"fileId":     "file_2",
			"filename":   "guide.md",
			"fileType":   "text/markdown",
			"fileUrl":    "/files/guide.md",
			"text":       "Another chunk with more detailed information that was retrieved from the RAG system.",
			"similarity": 0.87,
		},
	}

	// Mock tools (will be matched with tool messages below)
	tools := []map[string]interface{}{
		{
			"id":         "tool_1",
			"identifier": "lobe-web-browsing",
			"apiName":    "search",
			"arguments": map[string]interface{}{
				"query":         "What is the weather today?",
				"searchEngines": []string{"google"},
			},
			"type": "builtin",
		},
		{
			"id":         "tool_2",
			"identifier": "lobe-local-system",
			"apiName":    "listLocalFiles",
			"arguments": map[string]interface{}{
				"path": "/home/user/documents",
			},
			"type": "builtin",
		},
	}
	toolsJSON, _ := json.Marshal(tools)

	// Mock search grounding
	search := map[string]interface{}{
		"citations": []map[string]interface{}{
			{
				"id":    "citation_1",
				"title": "Wikipedia - Example Article",
				"url":   "https://en.wikipedia.org/wiki/Example",
			},
			{
				"id":    "citation_2",
				"title": "GitHub Documentation",
				"url":   "https://docs.github.com/en",
			},
		},
		"searchQueries": []string{"test query", "related query"},
	}
	searchJSON, _ := json.Marshal(search)

	// Mock image list
	imageList := []map[string]interface{}{
		{
			"id":  "img_1",
			"url": "https://via.placeholder.com/300x200",
			"alt": "Sample image 1",
		},
	}

	// Mock usage
	usage := map[string]interface{}{
		"prompt_tokens":     150,
		"completion_tokens": 80,
		"total_tokens":      230,
	}

	// Mock performance
	performance := map[string]interface{}{
		"total_tokens": 230,
		"duration":     1500,
	}

	// Combine all metadata into one JSON object
	fullMetadata := map[string]interface{}{
		"model":       "mock-model",
		"temperature": 0.7,
		"chunksList":  chunksList,
		"imageList":   imageList,
		"usage":       usage,
		"performance": performance,
	}
	fullMetadataJSON, _ := json.Marshal(fullMetadata)

	// Save assistant message
	assistantParams := db.CreateMessageParams{
		ID:        assistantMsgID,
		Role:      "assistant",
		Content:   sql.NullString{String: mockContent, Valid: true},
		SessionID: sql.NullString{String: req.SessionID, Valid: true},
		UserID:    req.UserID,
		CreatedAt: now + 1,
		UpdatedAt: now + 1,
		Reasoning: sql.NullString{String: string(reasoningJSON), Valid: true},
		Tools:     sql.NullString{String: string(toolsJSON), Valid: true},
		Search:    sql.NullString{String: string(searchJSON), Valid: true},
		Metadata:  sql.NullString{String: string(fullMetadataJSON), Valid: true},
	}
	if currentTopicID != "" {
		assistantParams.TopicID = sql.NullString{String: currentTopicID, Valid: true}
	}
	if req.ThreadID != "" {
		assistantParams.ThreadID = sql.NullString{String: req.ThreadID, Valid: true}
	}

	_, err = s.db.Queries().CreateMessage(ctx, assistantParams)
	if err != nil {
		return nil, fmt.Errorf("failed to save assistant message: %w", err)
	}
	log.Printf("💾 Saved mock assistant message: %s", assistantMsgID)

	// 3. Create tool messages (role='tool') with plugins - these are separate messages
	// Tool message 1: Web browsing results
	toolMsg1ID := uuid.New().String()
	toolResult1 := map[string]interface{}{
		"results": []map[string]interface{}{
			{
				"title":       "Mock Search Result 1",
				"url":         "https://example.com/result1",
				"description": "This is a mock search result for testing purposes.",
			},
			{
				"title":       "Mock Search Result 2",
				"url":         "https://example.com/result2",
				"description": "Another mock search result with relevant information.",
			},
		},
	}
	toolResult1JSON, _ := json.Marshal(toolResult1)

	toolMsg1Params := db.CreateMessageParams{
		ID:        toolMsg1ID,
		Role:      "tool",
		Content:   sql.NullString{String: string(toolResult1JSON), Valid: true},
		SessionID: sql.NullString{String: req.SessionID, Valid: true},
		UserID:    req.UserID,
		CreatedAt: now + 2,
		UpdatedAt: now + 2,
	}
	if currentTopicID != "" {
		toolMsg1Params.TopicID = sql.NullString{String: currentTopicID, Valid: true}
	}
	if req.ThreadID != "" {
		toolMsg1Params.ThreadID = sql.NullString{String: req.ThreadID, Valid: true}
	}

	_, err = s.db.Queries().CreateMessage(ctx, toolMsg1Params)
	if err != nil {
		return nil, fmt.Errorf("failed to save tool message 1: %w", err)
	}

	// Create plugin entry for tool message 1
	plugin1Params := db.CreateMessagePluginParams{
		ID:         toolMsg1ID,
		ToolCallID: sql.NullString{String: "tool_1", Valid: true}, // Must match tool id
		Type:       sql.NullString{String: "builtin", Valid: true},
		ApiName:    sql.NullString{String: "search", Valid: true},
		Identifier: sql.NullString{String: "lobe-web-browsing", Valid: true},
		UserID:     req.UserID,
	}
	_, err = s.db.Queries().CreateMessagePlugin(ctx, plugin1Params)
	if err != nil {
		return nil, fmt.Errorf("failed to save tool plugin 1: %w", err)
	}
	log.Printf("💾 Saved mock tool message 1 with plugin: %s", toolMsg1ID)

	// Tool message 2: Local file system results
	toolMsg2ID := uuid.New().String()
	toolResult2 := map[string]interface{}{
		"files": []map[string]interface{}{
			{"name": "document.pdf", "size": 1024000, "type": "file"},
			{"name": "images", "size": 0, "type": "directory"},
			{"name": "notes.txt", "size": 2048, "type": "file"},
		},
	}
	toolResult2JSON, _ := json.Marshal(toolResult2)

	toolMsg2Params := db.CreateMessageParams{
		ID:        toolMsg2ID,
		Role:      "tool",
		Content:   sql.NullString{String: string(toolResult2JSON), Valid: true},
		SessionID: sql.NullString{String: req.SessionID, Valid: true},
		UserID:    req.UserID,
		CreatedAt: now + 3,
		UpdatedAt: now + 3,
	}
	if currentTopicID != "" {
		toolMsg2Params.TopicID = sql.NullString{String: currentTopicID, Valid: true}
	}
	if req.ThreadID != "" {
		toolMsg2Params.ThreadID = sql.NullString{String: req.ThreadID, Valid: true}
	}

	_, err = s.db.Queries().CreateMessage(ctx, toolMsg2Params)
	if err != nil {
		return nil, fmt.Errorf("failed to save tool message 2: %w", err)
	}

	// Create plugin entry for tool message 2
	plugin2Params := db.CreateMessagePluginParams{
		ID:         toolMsg2ID,
		ToolCallID: sql.NullString{String: "tool_2", Valid: true}, // Must match tool id
		Type:       sql.NullString{String: "builtin", Valid: true},
		ApiName:    sql.NullString{String: "listLocalFiles", Valid: true},
		Identifier: sql.NullString{String: "lobe-local-system", Valid: true},
		UserID:     req.UserID,
	}
	_, err = s.db.Queries().CreateMessagePlugin(ctx, plugin2Params)
	if err != nil {
		return nil, fmt.Errorf("failed to save tool plugin 2: %w", err)
	}
	log.Printf("💾 Saved mock tool message 2 with plugin: %s", toolMsg2ID)

	log.Printf("✅ [MOCK] Complete - saved 4 messages (1 user, 1 assistant, 2 tools)")

	// Return response
	return &ChatResponse{
		MessageID:    assistantMsgID,
		SessionID:    req.SessionID,
		TopicID:      currentTopicID,
		ThreadID:     req.ThreadID,
		Message:      mockContent,
		FinishReason: "stop",
		CreatedAt:    now,
	}, nil
}
