package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)

// ============================================
// ChatMockStream - Event Streaming Mock
// ============================================

// StreamEventType represents the type of stream event
type StreamEventType string

const (
	StreamEventStart      StreamEventType = "start"       // Generation started
	StreamEventChunk      StreamEventType = "chunk"       // Content chunk
	StreamEventReasoning  StreamEventType = "reasoning"   // Reasoning content
	StreamEventToolCall   StreamEventType = "tool_call"   // Tool call initiated
	StreamEventToolResult StreamEventType = "tool_result" // Tool execution result
	StreamEventComplete   StreamEventType = "complete"    // Generation complete
)

// ChatMockStream handles mock chat with event streaming for realistic UI testing
// Instead of returning all messages at once, it emits events progressively:
// 1. start - Generation begins
// 2. reasoning - Thinking content (streamed)
// 3. chunk - Content chunks (streamed word by word)
// 4. tool_call - Tool call initiated
// 5. tool_result - Tool execution result (with pluginState)
// 6. complete - Generation finished
//
// Frontend listens to 'chat:stream' events via Events.On()
// Data is saved to DB at the end (same as ChatMock)
//
// Usage from frontend:
//
//	await AgentChatService.ChatRealStream(request);
//	// No return value - data comes via events
//	// Events.On('chat:stream', handler) receives all updates
func (s *AgentChatService) ChatRealStream(ctx context.Context, req ChatRequest) error {
	log.Printf("🎭 [REAL STREAM] Starting streaming real for session: %s", req.SessionID)

	// Helper to emit events with type safety
	emit := func(eventType StreamEventType, data interface{}) {
		if s.app == nil {
			return
		}

		// 1. Create base map with common fields
		payload := map[string]interface{}{
			"type":       string(eventType),
			"session_id": req.SessionID,
			"message_id": req.MessageAssistantID,
		}

		// 2. Merge data fields into payload
		// Using JSON round-trip to convert struct to map[string]interface{}
		// This is a simple way to merge without reflection complexity
		if data != nil {
			jsonData, _ := json.Marshal(data)
			var dataMap map[string]interface{}
			_ = json.Unmarshal(jsonData, &dataMap)

			for k, v := range dataMap {
				payload[k] = v
			}
		}

		s.app.Event.Emit("chat:stream", payload)
	}

	// 1. Setup session, topic, and save user message
	setup, err := s.setupSessionAndTopic(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to setup session/topic: %w", err)
	}
	currentTopicID := setup.TopicID

	// 2. Emit START event
	// Use UIChatMessage for consistency
	emit(StreamEventStart, &UIChatMessage{
		TopicID: currentTopicID,
	})
	time.Sleep(100 * time.Millisecond)

	// 3. Stream REASONING content
	reasoningContent := "Let me think about this step by step:\n1. First, I need to understand the question\n2. Then, I will formulate a response\n3. Finally, I will provide a clear answer"
	reasoningWords := splitIntoChunks(reasoningContent, 5) // 5 words per chunk

	for i, chunk := range reasoningWords {
		partialContent := joinChunks(reasoningWords[:i+1])
		// Use UIChatMessage with Reasoning field
		emit(StreamEventReasoning, &UIChatMessage{
			Reasoning: &ModelReasoning{
				Content: partialContent, // Frontend expects full content in reasoning.content
			},
			Content: chunk, // Optional: delta content
		})
		time.Sleep(50 * time.Millisecond)
	}
	time.Sleep(200 * time.Millisecond)

	// 4. Stream CONTENT
	mockContent := fmt.Sprintf(
		"This is a mock response to: \"%s\"\n\nI'm simulating the AI response to test the UI flow without calling the backend.",
		req.Message,
	)
	contentWords := splitIntoChunks(mockContent, 3) // 3 words per chunk

	for i, _ := range contentWords {
		partialContent := joinChunks(contentWords[:i+1])
		// Use UIChatMessage with Content field
		emit(StreamEventChunk, &UIChatMessage{
			Content: partialContent, // Frontend expects full content
		})
		time.Sleep(30 * time.Millisecond)
	}
	time.Sleep(200 * time.Millisecond)

	// 5. Prepare mock data (same as ChatMock)
	reasoning := &ModelReasoning{
		Content:  reasoningContent,
		Duration: 1500,
	}

	chunksList := []ChatFileChunk{
		{ID: "chunk_1", FileID: "file_1", Filename: "document.pdf", FileType: "application/pdf", FileURL: "/files/document.pdf", Text: "Sample chunk from knowledge base.", Similarity: 0.95},
		{ID: "chunk_2", FileID: "file_2", Filename: "guide.md", FileType: "text/markdown", FileURL: "/files/guide.md", Text: "Another chunk with detailed information.", Similarity: 0.87},
	}

	// Build tools array
	tools := s.buildMockTools()

	// 6. Stream TOOL_CALL and TOOL_RESULT for each tool
	toolResults := s.buildMockToolResults()

	for i, tool := range tools {
		// Emit tool_call
		// Use UIChatMessage with Tools field
		emit(StreamEventToolCall, &UIChatMessage{
			Tools: tools[:i+1], // All tools so far (for UI to render)
		})
		time.Sleep(300 * time.Millisecond) // Simulate tool execution

		// Get result for this tool
		result := toolResults[tool.ID]

		// Emit tool_result with pluginState for UI rendering
		// Let's go back to map for tool_result loop to avoid complexity,
		// BUT use ChatPluginPayload for the plugin field.
		toolResultPayload := map[string]interface{}{
			"tool_call_id": tool.ID,
			"tool_msg_id":  fmt.Sprintf("tool_msg_%s_%d", req.MessageAssistantID, i+1),
			"plugin": ChatPluginPayload{
				Identifier: tool.Identifier,
				APIName:    tool.APIName,
				Arguments:  tool.Arguments,
				Type:       tool.Type,
			},
			"pluginState": result.State,
			"content":     result.Content, // Content can be string or marshaled JSON
		}
		emit(StreamEventToolResult, toolResultPayload)

		time.Sleep(200 * time.Millisecond)
	}

	// 7. Build final metadata
	searchGrounding := &GroundingSearch{
		Citations: []CitationItem{
			{ID: "citation_1", Title: "Wikipedia - Example Article", URL: "https://en.wikipedia.org/wiki/Example"},
			{ID: "citation_2", Title: "GitHub Documentation", URL: "https://docs.github.com/en"},
		},
		SearchQueries: []string{"test query", "related query"},
	}

	imageList := []ChatImageItem{
		{ID: "img_1", URL: "https://via.placeholder.com/300x200", Alt: "Sample image 1"},
	}

	mockUsage := &ModelUsage{TotalInputTokens: 150, TotalOutputTokens: 80, TotalTokens: 230}
	mockPerformance := &ModelPerformance{Duration: 1500, Latency: 1800}

	fullMetadata := map[string]interface{}{
		"model": "mock-model", "temperature": 0.7,
		"chunksList": chunksList, "imageList": imageList,
		"usage": mockUsage, "performance": mockPerformance,
	}

	// 8. Save assistant message to DB
	assistantMsgID, err := s.saveAssistantMessage(ctx, SaveAssistantMessageParams{
		MessageID: req.MessageAssistantID,
		Content:   mockContent,
		SessionID: req.SessionID,
		TopicID:   currentTopicID,
		ThreadID:  req.ThreadID,
		UserID:    req.UserID,
		Reasoning: reasoning,
		Tools:     tools,
		Search:    searchGrounding,
		Metadata:  fullMetadata,
	})
	if err != nil {
		return fmt.Errorf("failed to save assistant message: %w", err)
	}

	// 9. Save RAG data
	file1ID := uuid.New().String()
	file2ID := uuid.New().String()
	chunk1ID := uuid.New().String()
	chunk2ID := uuid.New().String()
	_ = s.saveRAGData(ctx, SaveRAGDataParams{
		MessageID: assistantMsgID,
		UserQuery: req.Message,
		UserID:    req.UserID,
		Files: []RAGFileParams{
			{ID: file1ID, Name: "document.pdf", FileType: "application/pdf", URL: "", Size: 1024000},
			{ID: file2ID, Name: "guide.md", FileType: "text/markdown", URL: "", Size: 2048},
		},
		Chunks: []RAGChunkParams{
			{ID: chunk1ID, FileID: file1ID, Text: "Sample chunk.", ChunkIndex: 0, Type: "text", Similarity: 95},
			{ID: chunk2ID, FileID: file2ID, Text: "Another chunk.", ChunkIndex: 0, Type: "text", Similarity: 87},
		},
	})

	// 10. Save tool messages to DB
	for i, tool := range tools {
		result := toolResults[tool.ID]
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

	// 11. Emit COMPLETE event with final data
	// Use UIChatMessage for complete event
	emit(StreamEventComplete, &UIChatMessage{
		Content:     mockContent,
		TopicID:     currentTopicID,
		Reasoning:   reasoning,
		Search:      searchGrounding,
		ChunksList:  chunksList,
		ImageList:   imageList,
		Usage:       mockUsage,
		Performance: mockPerformance,
	})

	log.Printf("✅ [MOCK STREAM] Complete - emitted all events for session: %s", req.SessionID)
	return nil
}

// ToolResultData holds content and state for a tool result
type ToolResultData struct {
	Content interface{}
	State   interface{}
}
