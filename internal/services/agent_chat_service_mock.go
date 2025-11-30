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

	// ============================================
	// Mock ALL builtin tools with their arguments
	// ============================================

	// Tool 1: Web Browsing - search
	tool1Args := map[string]interface{}{
		"query":         "What is the weather today?",
		"searchEngines": []string{"google", "bing"},
	}
	tool1ArgsJSON, _ := json.Marshal(tool1Args)

	// Tool 2: Web Browsing - crawlSinglePage
	tool2Args := map[string]interface{}{
		"url": "https://example.com/article",
	}
	tool2ArgsJSON, _ := json.Marshal(tool2Args)

	// Tool 3: Web Browsing - crawlMultiPages
	tool3Args := map[string]interface{}{
		"urls": []string{"https://example.com/page1", "https://example.com/page2"},
	}
	tool3ArgsJSON, _ := json.Marshal(tool3Args)

	// Tool 4: Local System - listLocalFiles
	tool4Args := map[string]interface{}{
		"path": "/home/user/documents",
	}
	tool4ArgsJSON, _ := json.Marshal(tool4Args)

	// Tool 5: Local System - readLocalFile
	tool5Args := map[string]interface{}{
		"path": "/home/user/documents/readme.md",
		"loc":  []int{0, 100},
	}
	tool5ArgsJSON, _ := json.Marshal(tool5Args)

	// Tool 6: Local System - searchLocalFiles
	tool6Args := map[string]interface{}{
		"keywords":  "important document",
		"directory": "/home/user/documents",
	}
	tool6ArgsJSON, _ := json.Marshal(tool6Args)

	// Tool 7: Local System - writeLocalFile
	tool7Args := map[string]interface{}{
		"path":    "/home/user/documents/new_file.txt",
		"content": "Hello, this is a new file content.",
	}
	tool7ArgsJSON, _ := json.Marshal(tool7Args)

	// Tool 8: Local System - renameLocalFile
	tool8Args := map[string]interface{}{
		"path":    "/home/user/documents/old_name.txt",
		"newName": "new_name.txt",
	}
	tool8ArgsJSON, _ := json.Marshal(tool8Args)

	// Tool 9: Local System - moveLocalFiles
	tool9Args := map[string]interface{}{
		"items": []map[string]interface{}{
			{"oldPath": "/home/user/documents/file1.txt", "newPath": "/home/user/backup/file1.txt"},
		},
	}
	tool9ArgsJSON, _ := json.Marshal(tool9Args)

	// Tool 10: DALL-E Image Designer - text2image
	tool10Args := map[string]interface{}{
		"prompts": []string{
			"A beautiful sunset over a calm ocean with vibrant orange and purple colors",
			"A futuristic cityscape at night with neon lights and flying cars",
		},
		"size":    "1024x1024",
		"quality": "hd",
		"style":   "vivid",
	}
	tool10ArgsJSON, _ := json.Marshal(tool10Args)

	// Tool 11: Code Interpreter - python
	tool11Args := map[string]interface{}{
		"code":     "import pandas as pd\nimport numpy as np\n\ndf = pd.DataFrame({'a': [1, 2, 3], 'b': [4, 5, 6]})\nprint(df.describe())\nresult = df.sum().to_dict()",
		"packages": []string{"pandas", "numpy"},
	}
	tool11ArgsJSON, _ := json.Marshal(tool11Args)

	// Mock tools array (will be matched with tool messages below)
	// Note: arguments must be JSON string for frontend compatibility
	tools := []map[string]interface{}{
		// Web Browsing tools
		{
			"id":         "tool_1",
			"identifier": "lobe-web-browsing",
			"apiName":    "search",
			"arguments":  string(tool1ArgsJSON),
			"type":       "builtin",
		},
		{
			"id":         "tool_2",
			"identifier": "lobe-web-browsing",
			"apiName":    "crawlSinglePage",
			"arguments":  string(tool2ArgsJSON),
			"type":       "builtin",
		},
		{
			"id":         "tool_3",
			"identifier": "lobe-web-browsing",
			"apiName":    "crawlMultiPages",
			"arguments":  string(tool3ArgsJSON),
			"type":       "builtin",
		},
		// Local System tools
		{
			"id":         "tool_4",
			"identifier": "lobe-local-system",
			"apiName":    "listLocalFiles",
			"arguments":  string(tool4ArgsJSON),
			"type":       "builtin",
		},
		{
			"id":         "tool_5",
			"identifier": "lobe-local-system",
			"apiName":    "readLocalFile",
			"arguments":  string(tool5ArgsJSON),
			"type":       "builtin",
		},
		{
			"id":         "tool_6",
			"identifier": "lobe-local-system",
			"apiName":    "searchLocalFiles",
			"arguments":  string(tool6ArgsJSON),
			"type":       "builtin",
		},
		{
			"id":         "tool_7",
			"identifier": "lobe-local-system",
			"apiName":    "writeLocalFile",
			"arguments":  string(tool7ArgsJSON),
			"type":       "builtin",
		},
		{
			"id":         "tool_8",
			"identifier": "lobe-local-system",
			"apiName":    "renameLocalFile",
			"arguments":  string(tool8ArgsJSON),
			"type":       "builtin",
		},
		{
			"id":         "tool_9",
			"identifier": "lobe-local-system",
			"apiName":    "moveLocalFiles",
			"arguments":  string(tool9ArgsJSON),
			"type":       "builtin",
		},
		// DALL-E Image Designer
		{
			"id":         "tool_10",
			"identifier": "lobe-image-designer",
			"apiName":    "text2image",
			"arguments":  string(tool10ArgsJSON),
			"type":       "builtin",
		},
		// Code Interpreter
		{
			"id":         "tool_11",
			"identifier": "lobe-code-interpreter",
			"apiName":    "python",
			"arguments":  string(tool11ArgsJSON),
			"type":       "builtin",
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

	// 3. Create mock RAG data (files, chunks, message_query_chunks)
	// Create mock files
	file1ID := uuid.New().String()
	file1Params := db.CreateFileParams{
		ID:        file1ID,
		Name:      "document.pdf",
		FileType:  "application/pdf",
		Url:       "",
		Size:      1024000,
		UserID:    req.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err = s.db.Queries().CreateFile(ctx, file1Params)
	if err != nil {
		log.Printf("⚠️  Failed to create mock file 1: %v", err)
	}

	file2ID := uuid.New().String()
	file2Params := db.CreateFileParams{
		ID:        file2ID,
		Name:      "guide.md",
		FileType:  "text/markdown",
		Url:       "",
		Size:      2048,
		UserID:    req.UserID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	_, err = s.db.Queries().CreateFile(ctx, file2Params)
	if err != nil {
		log.Printf("⚠️  Failed to create mock file 2: %v", err)
	}

	// Create mock chunks
	chunk1ID := uuid.New().String()
	chunk1Params := db.CreateChunkParams{
		ID:         chunk1ID,
		Text:       sql.NullString{String: "This is a sample chunk from the knowledge base. It contains relevant information about the topic.", Valid: true},
		ChunkIndex: sql.NullInt64{Int64: 0, Valid: true},
		Type:       sql.NullString{String: "text", Valid: true},
		UserID:     sql.NullString{String: req.UserID, Valid: true},
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	_, err = s.db.Queries().CreateChunk(ctx, chunk1Params)
	if err != nil {
		log.Printf("⚠️  Failed to create mock chunk 1: %v", err)
	}

	chunk2ID := uuid.New().String()
	chunk2Params := db.CreateChunkParams{
		ID:         chunk2ID,
		Text:       sql.NullString{String: "Another chunk with more detailed information that was retrieved from the RAG system.", Valid: true},
		ChunkIndex: sql.NullInt64{Int64: 0, Valid: true},
		Type:       sql.NullString{String: "text", Valid: true},
		UserID:     sql.NullString{String: req.UserID, Valid: true},
		CreatedAt:  now,
		UpdatedAt:  now,
	}
	_, err = s.db.Queries().CreateChunk(ctx, chunk2Params)
	if err != nil {
		log.Printf("⚠️  Failed to create mock chunk 2: %v", err)
	}

	// Link chunks to files
	err = s.db.Queries().LinkFileToChunk(ctx, db.LinkFileToChunkParams{
		FileID:    sql.NullString{String: file1ID, Valid: true},
		ChunkID:   sql.NullString{String: chunk1ID, Valid: true},
		CreatedAt: now,
		UserID:    req.UserID,
	})
	if err != nil {
		log.Printf("⚠️  Failed to link file 1 to chunk 1: %v", err)
	}

	err = s.db.Queries().LinkFileToChunk(ctx, db.LinkFileToChunkParams{
		FileID:    sql.NullString{String: file2ID, Valid: true},
		ChunkID:   sql.NullString{String: chunk2ID, Valid: true},
		CreatedAt: now,
		UserID:    req.UserID,
	})
	if err != nil {
		log.Printf("⚠️  Failed to link file 2 to chunk 2: %v", err)
	}

	// Create message query
	queryID := uuid.New().String()
	queryParams := db.CreateMessageQueryParams{
		ID:           queryID,
		MessageID:    assistantMsgID,
		UserQuery:    sql.NullString{String: req.Message, Valid: true},
		RewriteQuery: sql.NullString{String: req.Message, Valid: true},
		UserID:       req.UserID,
	}
	_, err = s.db.Queries().CreateMessageQuery(ctx, queryParams)
	if err != nil {
		log.Printf("⚠️  Failed to create message query: %v", err)
	}

	// Link message query to chunks
	err = s.db.Queries().LinkMessageQueryToChunk(ctx, db.LinkMessageQueryToChunkParams{
		MessageID:  sql.NullString{String: assistantMsgID, Valid: true},
		QueryID:    sql.NullString{String: queryID, Valid: true},
		ChunkID:    sql.NullString{String: chunk1ID, Valid: true},
		Similarity: sql.NullInt64{Int64: 95, Valid: true},
		UserID:     req.UserID,
	})
	if err != nil {
		log.Printf("⚠️  Failed to link query to chunk 1: %v", err)
	}

	err = s.db.Queries().LinkMessageQueryToChunk(ctx, db.LinkMessageQueryToChunkParams{
		MessageID:  sql.NullString{String: assistantMsgID, Valid: true},
		QueryID:    sql.NullString{String: queryID, Valid: true},
		ChunkID:    sql.NullString{String: chunk2ID, Valid: true},
		Similarity: sql.NullInt64{Int64: 87, Valid: true},
		UserID:     req.UserID,
	})
	if err != nil {
		log.Printf("⚠️  Failed to link query to chunk 2: %v", err)
	}

	log.Printf("💾 Created mock RAG data: 2 files, 2 chunks, 1 query")

	// ============================================
	// 4. Create tool messages (role='tool') with plugins for ALL builtin tools
	// ============================================

	// Helper to create tool message and plugin
	// result is stored in message.content, state is stored in plugin.state (for pluginState)
	createToolMessage := func(toolID, identifier, apiName string, argsJSON []byte, result interface{}, state interface{}, timeOffset int64) error {
		msgID := uuid.New().String()
		resultJSON, _ := json.Marshal(result)

		msgParams := db.CreateMessageParams{
			ID:        msgID,
			Role:      "tool",
			Content:   sql.NullString{String: string(resultJSON), Valid: true},
			SessionID: sql.NullString{String: req.SessionID, Valid: true},
			UserID:    req.UserID,
			CreatedAt: now + timeOffset,
			UpdatedAt: now + timeOffset,
		}
		if currentTopicID != "" {
			msgParams.TopicID = sql.NullString{String: currentTopicID, Valid: true}
		}
		if req.ThreadID != "" {
			msgParams.ThreadID = sql.NullString{String: req.ThreadID, Valid: true}
		}

		_, err := s.db.Queries().CreateMessage(ctx, msgParams)
		if err != nil {
			return fmt.Errorf("failed to save tool message %s: %w", toolID, err)
		}

		pluginParams := db.CreateMessagePluginParams{
			ID:         msgID,
			ToolCallID: sql.NullString{String: toolID, Valid: true},
			Type:       sql.NullString{String: "builtin", Valid: true},
			ApiName:    sql.NullString{String: apiName, Valid: true},
			Arguments:  sql.NullString{String: string(argsJSON), Valid: true},
			Identifier: sql.NullString{String: identifier, Valid: true},
			UserID:     req.UserID,
		}
		// Add state if provided (for pluginState in frontend)
		if state != nil {
			stateJSON, _ := json.Marshal(state)
			pluginParams.State = sql.NullString{String: string(stateJSON), Valid: true}
		}
		_, err = s.db.Queries().CreateMessagePlugin(ctx, pluginParams)
		if err != nil {
			return fmt.Errorf("failed to save tool plugin %s: %w", toolID, err)
		}

		log.Printf("💾 Saved tool message: %s (%s.%s)", toolID, identifier, apiName)
		return nil
	}

	// Tool 1: Web Browsing - search
	err = createToolMessage("tool_1", "lobe-web-browsing", "search", tool1ArgsJSON,
		map[string]interface{}{
			"results": []map[string]interface{}{
				{"title": "Weather Today - Current Conditions", "url": "https://weather.com/today", "description": "Current weather conditions and forecast."},
				{"title": "Weather Report - Bing", "url": "https://www.bing.com/weather", "description": "Detailed weather information."},
			},
		}, nil, 2)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 2: Web Browsing - crawlSinglePage
	// Format: pluginState.results = [{ originalUrl, data: { content, url, title } }]
	crawlSingleResult := []map[string]interface{}{
		{
			"originalUrl": "https://example.com/article",
			"crawler":     "jina",
			"data": map[string]interface{}{
				"content":     "# Example Article\n\nThis is the full content of the crawled article. It contains detailed information about the topic discussed.\n\n## Section 1\nSome content here...\n\n## Section 2\nMore content here...",
				"url":         "https://example.com/article",
				"title":       "Example Article - Full Content",
				"description": "A comprehensive article about the topic.",
			},
		},
	}
	err = createToolMessage("tool_2", "lobe-web-browsing", "crawlSinglePage", tool2ArgsJSON,
		nil, // content
		map[string]interface{}{"results": crawlSingleResult}, // pluginState
		3)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 3: Web Browsing - crawlMultiPages
	// Format: pluginState.results = [{ originalUrl, data: { content, url, title } }, ...]
	crawlMultiResult := []map[string]interface{}{
		{
			"originalUrl": "https://example.com/page1",
			"crawler":     "jina",
			"data": map[string]interface{}{
				"content":     "# Page 1\n\nThis is the content from page 1. It provides information about topic A.\n\n## Overview\nDetailed overview...",
				"url":         "https://example.com/page1",
				"title":       "Page 1 - Topic A",
				"description": "Information about topic A.",
			},
		},
		{
			"originalUrl": "https://example.com/page2",
			"crawler":     "jina",
			"data": map[string]interface{}{
				"content":     "# Page 2\n\nThis is the content from page 2. It provides information about topic B.\n\n## Details\nMore details here...",
				"url":         "https://example.com/page2",
				"title":       "Page 2 - Topic B",
				"description": "Information about topic B.",
			},
		},
	}
	err = createToolMessage("tool_3", "lobe-web-browsing", "crawlMultiPages", tool3ArgsJSON,
		nil, // content
		map[string]interface{}{"results": crawlMultiResult}, // pluginState
		4)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 4: Local System - listLocalFiles
	// Format: pluginState.listResults = [{ name, size, type, isDirectory, path, ... }]
	listFilesResult := []map[string]interface{}{
		{"name": "document.pdf", "size": 1024000, "type": "file", "isDirectory": false, "path": "/home/user/documents/document.pdf"},
		{"name": "images", "size": 0, "type": "directory", "isDirectory": true, "path": "/home/user/documents/images"},
		{"name": "notes.txt", "size": 2048, "type": "file", "isDirectory": false, "path": "/home/user/documents/notes.txt"},
		{"name": "readme.md", "size": 512, "type": "file", "isDirectory": false, "path": "/home/user/documents/readme.md"},
	}
	err = createToolMessage("tool_4", "lobe-local-system", "listLocalFiles", tool4ArgsJSON,
		nil, // content
		map[string]interface{}{"listResults": listFilesResult}, // pluginState
		5)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 5: Local System - readLocalFile
	// Format: pluginState.fileContent = { content, filename, fileType, charCount, lineCount, ... }
	readFileResult := map[string]interface{}{
		"content":        "# README\n\nThis is a sample readme file.\n\n## Features\n- Feature 1\n- Feature 2\n- Feature 3",
		"filename":       "readme.md",
		"fileType":       "text/markdown",
		"charCount":      120,
		"lineCount":      8,
		"totalCharCount": 120,
		"totalLineCount": 8,
	}
	err = createToolMessage("tool_5", "lobe-local-system", "readLocalFile", tool5ArgsJSON,
		nil, // content
		map[string]interface{}{"fileContent": readFileResult}, // pluginState
		6)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 6: Local System - searchLocalFiles
	// Format: pluginState.searchResults = [{ path, name, size, isDirectory, ... }]
	searchFilesResult := []map[string]interface{}{
		{"path": "/home/user/documents/important_doc.pdf", "name": "important_doc.pdf", "size": 2048000, "isDirectory": false, "type": "file"},
		{"path": "/home/user/documents/important_notes.txt", "name": "important_notes.txt", "size": 1024, "isDirectory": false, "type": "file"},
	}
	err = createToolMessage("tool_6", "lobe-local-system", "searchLocalFiles", tool6ArgsJSON,
		nil, // content
		map[string]interface{}{"searchResults": searchFilesResult}, // pluginState
		7)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 7: Local System - writeLocalFile
	err = createToolMessage("tool_7", "lobe-local-system", "writeLocalFile", tool7ArgsJSON,
		map[string]interface{}{
			"success": true,
			"path":    "/home/user/documents/new_file.txt",
			"message": "File written successfully",
		}, nil, 8)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 8: Local System - renameLocalFile
	err = createToolMessage("tool_8", "lobe-local-system", "renameLocalFile", tool8ArgsJSON,
		map[string]interface{}{
			"success": true,
			"oldPath": "/home/user/documents/old_name.txt",
			"newPath": "/home/user/documents/new_name.txt",
		}, nil, 9)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 9: Local System - moveLocalFiles
	err = createToolMessage("tool_9", "lobe-local-system", "moveLocalFiles", tool9ArgsJSON,
		map[string]interface{}{
			"results": []map[string]interface{}{
				{"sourcePath": "/home/user/documents/file1.txt", "newPath": "/home/user/backup/file1.txt", "success": true},
			},
		}, nil, 10)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 10: DALL-E Image Designer - text2image
	err = createToolMessage("tool_10", "lobe-image-designer", "text2image", tool10ArgsJSON,
		[]map[string]interface{}{
			{
				"prompt":        "A beautiful sunset over a calm ocean with vibrant orange and purple colors",
				"imageUrl":      "https://picsum.photos/1024",
				"revisedPrompt": "A breathtaking sunset scene over a tranquil ocean, with vibrant orange and deep purple hues painting the sky.",
			},
			{
				"prompt":        "A futuristic cityscape at night with neon lights and flying cars",
				"imageUrl":      "https://picsum.photos/1024",
				"revisedPrompt": "A stunning futuristic cityscape at night, illuminated by neon lights with sleek flying cars hovering above.",
			},
		}, nil, 11)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 11: Code Interpreter - python
	// Format: content = { result, output: [{ type, data }], files: [] }
	codeInterpreterContent := map[string]interface{}{
		"result": "{'a': 6, 'b': 15}",
		"output": []map[string]interface{}{
			{"type": "stdout", "data": "              a         b\ncount  3.000000  3.000000\nmean   2.000000  5.000000\nstd    1.000000  1.000000\nmin    1.000000  4.000000\n25%    1.500000  4.500000\n50%    2.000000  5.000000\n75%    2.500000  5.500000\nmax    3.000000  6.000000\n"},
		},
		"files": []interface{}{},
	}
	err = createToolMessage("tool_11", "lobe-code-interpreter", "python", tool11ArgsJSON,
		codeInterpreterContent, // content
		nil,                    // pluginState (no error)
		12)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	log.Printf("✅ [MOCK] Complete - saved %d messages (1 user, 1 assistant, 11 tools)", 13)

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
