package services

import (
	"context"
	"fmt"
	"time"
)

// ChatMock handles mock chat responses for testing UI flow without real AI backend
// This method returns a simple mock response that matches ChatResponse structure
//
// Usage from frontend:
//
//	const response = await AgentChatService.ChatMock(request);
//
// Note: This is a simplified mock that doesn't save to DB or handle streaming.
// For full mock with all UI components (reasoning, tools, chunks, etc.),
// the frontend currently handles mock data directly in generateAIChat.ts
func (s *AgentChatService) ChatMock(ctx context.Context, req ChatRequest) (*ChatResponse, error) {
	// Simulate processing delay
	time.Sleep(500 * time.Millisecond)

	// Get current timestamp
	now := time.Now().UnixMilli()

	// Build mock response content
	mockContent := fmt.Sprintf(
		"This is a mock response to: \"%s\"\n\n"+
			"I'm simulating the AI response to test the UI flow without calling the backend.\n\n"+
			"**Mock Data Includes:**\n"+
			"- Reasoning: Step-by-step thinking process\n"+
			"- RAG Chunks: Retrieved file chunks with similarity scores\n"+
			"- Tool Calls: Web browsing and local file system tools\n"+
			"- Search Grounding: Citations from web search\n"+
			"- Images: Sample placeholder images\n"+
			"- Usage: Token counts and performance metrics",
		req.Message,
	)

	// Generate message IDs
	messageID := fmt.Sprintf("mock-assistant-%d", now)
	topicID := req.TopicID
	if topicID == "" {
		topicID = fmt.Sprintf("mock-topic-%d", now)
	}

	// Return ChatResponse matching the expected structure
	return &ChatResponse{
		MessageID:    messageID,
		SessionID:    req.SessionID,
		TopicID:      topicID,
		ThreadID:     req.ThreadID,
		Message:      mockContent,
		FinishReason: "stop",
		CreatedAt:    now,
	}, nil
}
