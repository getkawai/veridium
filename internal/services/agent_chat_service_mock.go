package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/kawai-network/veridium/internal/search"
	"github.com/kawai-network/veridium/pkg/localfs"
	"github.com/kawai-network/veridium/pkg/yzma/tools/builtin"
)

// CitationItem represents a citation from search results
type CitationItem struct {
	Favicon string `json:"favicon,omitempty"`
	ID      string `json:"id,omitempty"`
	Title   string `json:"title,omitempty"`
	URL     string `json:"url"`
}

// GroundingSearch represents search grounding data for messages
type GroundingSearch struct {
	Citations     []CitationItem `json:"citations,omitempty"`
	SearchQueries []string       `json:"searchQueries,omitempty"`
}

// ChatToolPayload represents a tool call payload
type ChatToolPayload struct {
	APIName    string `json:"apiName"`
	Arguments  string `json:"arguments"`
	ID         string `json:"id"`
	Identifier string `json:"identifier"`
	Type       string `json:"type"` // "builtin" or other tool types
}

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

	// 1. Setup session, topic, and save user message using reusable helper
	// SetupSessionAndTopic now handles:
	// - Get/create session
	// - Load history summary
	// - Load thread messages
	// - Auto-create topic
	// - Add user message to session (in-memory)
	// - Save user message to DB
	setup, err := s.setupSessionAndTopic(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to setup session/topic: %w", err)
	}
	currentTopicID := setup.TopicID

	// 2. Create assistant message with full mock data
	var assistantMsgID string
	mockContent := fmt.Sprintf(
		"This is a mock response to: \"%s\"\n\nI'm simulating the AI response to test the UI flow without calling the backend.",
		req.Message,
	)

	// Mock reasoning
	reasoning := map[string]interface{}{
		"content": "Let me think about this step by step:\n1. First, I need to understand the question\n2. Then, I will formulate a response\n3. Finally, I will provide a clear answer",
		"status":  "complete",
	}

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
	tool4Args := localfs.ListLocalFileParams{
		Path: "/home/user/documents",
	}
	tool4ArgsJSON, _ := json.Marshal(tool4Args)

	// Tool 5: Local System - readLocalFile
	loc := [2]int{0, 100}
	tool5Args := localfs.LocalReadFileParams{
		Path: "/home/user/documents/readme.md",
		Loc:  &loc,
	}
	tool5ArgsJSON, _ := json.Marshal(tool5Args)

	// Tool 6: Local System - searchLocalFiles
	tool6Args := localfs.LocalSearchFilesParams{
		Keywords:  "important document",
		Directory: "/home/user/documents",
	}
	tool6ArgsJSON, _ := json.Marshal(tool6Args)

	// Tool 7: Local System - writeLocalFile
	tool7Args := localfs.WriteLocalFileParams{
		Path:    "/home/user/documents/new_file.txt",
		Content: "Hello, this is a new file content.",
	}
	tool7ArgsJSON, _ := json.Marshal(tool7Args)

	// Tool 8: Local System - renameLocalFile
	tool8Args := localfs.RenameLocalFileParams{
		Path:    "/home/user/documents/old_name.txt",
		NewName: "new_name.txt",
	}
	tool8ArgsJSON, _ := json.Marshal(tool8Args)

	// Tool 9: Local System - moveLocalFiles
	tool9Args := localfs.MoveLocalFilesParams{
		Items: []localfs.MoveLocalFileParams{
			{OldPath: "/home/user/documents/file1.txt", NewPath: "/home/user/backup/file1.txt"},
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
	tools := []ChatToolPayload{
		// Web Browsing tools
		{
			ID:         "tool_1",
			Identifier: "lobe-web-browsing",
			APIName:    "search",
			Arguments:  string(tool1ArgsJSON),
			Type:       "builtin",
		},
		{
			ID:         "tool_2",
			Identifier: "lobe-web-browsing",
			APIName:    "crawlSinglePage",
			Arguments:  string(tool2ArgsJSON),
			Type:       "builtin",
		},
		{
			ID:         "tool_3",
			Identifier: "lobe-web-browsing",
			APIName:    "crawlMultiPages",
			Arguments:  string(tool3ArgsJSON),
			Type:       "builtin",
		},
		// Local System tools
		{
			ID:         "tool_4",
			Identifier: "lobe-local-system",
			APIName:    "listLocalFiles",
			Arguments:  string(tool4ArgsJSON),
			Type:       "builtin",
		},
		{
			ID:         "tool_5",
			Identifier: "lobe-local-system",
			APIName:    "readLocalFile",
			Arguments:  string(tool5ArgsJSON),
			Type:       "builtin",
		},
		{
			ID:         "tool_6",
			Identifier: "lobe-local-system",
			APIName:    "searchLocalFiles",
			Arguments:  string(tool6ArgsJSON),
			Type:       "builtin",
		},
		{
			ID:         "tool_7",
			Identifier: "lobe-local-system",
			APIName:    "writeLocalFile",
			Arguments:  string(tool7ArgsJSON),
			Type:       "builtin",
		},
		{
			ID:         "tool_8",
			Identifier: "lobe-local-system",
			APIName:    "renameLocalFile",
			Arguments:  string(tool8ArgsJSON),
			Type:       "builtin",
		},
		{
			ID:         "tool_9",
			Identifier: "lobe-local-system",
			APIName:    "moveLocalFiles",
			Arguments:  string(tool9ArgsJSON),
			Type:       "builtin",
		},
		// DALL-E Image Designer
		{
			ID:         "tool_10",
			Identifier: "lobe-image-designer",
			APIName:    "text2image",
			Arguments:  string(tool10ArgsJSON),
			Type:       "builtin",
		},
		// Code Interpreter
		{
			ID:         "tool_11",
			Identifier: "lobe-code-interpreter",
			APIName:    "python",
			Arguments:  string(tool11ArgsJSON),
			Type:       "builtin",
		},
	}

	// Mock search grounding
	searchGrounding := &GroundingSearch{
		Citations: []CitationItem{
			{
				ID:    "citation_1",
				Title: "Wikipedia - Example Article",
				URL:   "https://en.wikipedia.org/wiki/Example",
			},
			{
				ID:    "citation_2",
				Title: "GitHub Documentation",
				URL:   "https://docs.github.com/en",
			},
		},
		SearchQueries: []string{"test query", "related query"},
	}

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

	// 2b. Save assistant message using reusable helper
	assistantMsgID, err = s.saveAssistantMessage(ctx, SaveAssistantMessageParams{
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
		return nil, fmt.Errorf("failed to save assistant message: %w", err)
	}

	// 3. Save RAG data using reusable helper
	file1ID := uuid.New().String()
	file2ID := uuid.New().String()
	chunk1ID := uuid.New().String()
	chunk2ID := uuid.New().String()

	err = s.saveRAGData(ctx, SaveRAGDataParams{
		MessageID: assistantMsgID,
		UserQuery: req.Message,
		UserID:    req.UserID,
		Files: []RAGFileParams{
			{ID: file1ID, Name: "document.pdf", FileType: "application/pdf", URL: "", Size: 1024000},
			{ID: file2ID, Name: "guide.md", FileType: "text/markdown", URL: "", Size: 2048},
		},
		Chunks: []RAGChunkParams{
			{ID: chunk1ID, FileID: file1ID, Text: "This is a sample chunk from the knowledge base. It contains relevant information about the topic.", ChunkIndex: 0, Type: "text", Similarity: 95},
			{ID: chunk2ID, FileID: file2ID, Text: "Another chunk with more detailed information that was retrieved from the RAG system.", ChunkIndex: 0, Type: "text", Similarity: 87},
		},
	})
	if err != nil {
		log.Printf("⚠️  Failed to save RAG data: %v", err)
	}

	// ============================================
	// 4. Create tool messages using reusable helper
	// ============================================

	// Helper closure that wraps SaveToolMessage with common params
	saveToolMsg := func(toolID, identifier, apiName string, argsJSON []byte, content interface{}, state interface{}, timeOffset int64) error {
		_, err := s.saveToolMessage(ctx, SaveToolMessageParams{
			ToolCallID: toolID,
			Identifier: identifier,
			APIName:    apiName,
			Arguments:  string(argsJSON),
			Content:    content,
			State:      state,
			SessionID:  req.SessionID,
			TopicID:    currentTopicID,
			ThreadID:   req.ThreadID,
			UserID:     req.UserID,
			TimeOffset: timeOffset,
		})
		return err
	}

	// Tool 1: Web Browsing - search
	// Using proper types from searchpkg.UniformSearchResponse
	searchResponse := &search.UniformSearchResponse{
		Query:         "What is the weather today?",
		ResultNumbers: 2,
		CostTime:      150,
		Results: []search.UniformSearchResult{
			{
				Title:     "Weather Today - Current Conditions",
				URL:       "https://weather.com/today",
				Content:   "Current weather conditions and forecast.",
				ParsedUrl: "weather.com",
				Engines:   []string{"google"},
				Score:     0.95,
			},
			{
				Title:     "Weather Report - Bing",
				URL:       "https://www.bing.com/weather",
				Content:   "Detailed weather information.",
				ParsedUrl: "bing.com",
				Engines:   []string{"bing"},
				Score:     0.87,
			},
		},
	}
	err = saveToolMsg("tool_1", "lobe-web-browsing", "search", tool1ArgsJSON,
		searchResponse, nil, 2)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 2: Web Browsing - crawlSinglePage
	// Using proper types from builtin.CrawlPluginState
	crawlSingleState := &builtin.CrawlPluginState{
		Results: []builtin.CrawlResult{
			{
				OriginalUrl: "https://example.com/article",
				Crawler:     "jina",
				Data: builtin.CrawlData{
					Content:     "# Example Article\n\nThis is the full content of the crawled article. It contains detailed information about the topic discussed.\n\n## Section 1\nSome content here...\n\n## Section 2\nMore content here...",
					URL:         "https://example.com/article",
					Title:       "Example Article - Full Content",
					Description: "A comprehensive article about the topic.",
				},
			},
		},
	}
	err = saveToolMsg("tool_2", "lobe-web-browsing", "crawlSinglePage", tool2ArgsJSON,
		nil,              // content
		crawlSingleState, // pluginState
		3)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 3: Web Browsing - crawlMultiPages
	// Using proper types from builtin.CrawlPluginState
	crawlMultiState := &builtin.CrawlPluginState{
		Results: []builtin.CrawlResult{
			{
				OriginalUrl: "https://example.com/page1",
				Crawler:     "jina",
				Data: builtin.CrawlData{
					Content:     "# Page 1\n\nThis is the content from page 1. It provides information about topic A.\n\n## Overview\nDetailed overview...",
					URL:         "https://example.com/page1",
					Title:       "Page 1 - Topic A",
					Description: "Information about topic A.",
				},
			},
			{
				OriginalUrl: "https://example.com/page2",
				Crawler:     "jina",
				Data: builtin.CrawlData{
					Content:     "# Page 2\n\nThis is the content from page 2. It provides information about topic B.\n\n## Details\nMore details here...",
					URL:         "https://example.com/page2",
					Title:       "Page 2 - Topic B",
					Description: "Information about topic B.",
				},
			},
		},
	}
	err = saveToolMsg("tool_3", "lobe-web-browsing", "crawlMultiPages", tool3ArgsJSON,
		nil,             // content
		crawlMultiState, // pluginState
		4)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 4: Local System - listLocalFiles
	// Using proper types from builtin.LocalFileListState
	listFilesState := &builtin.LocalFileListState{
		ListResults: []localfs.LocalFileItem{
			{
				Name:         "document.pdf",
				Path:         "/home/user/documents/document.pdf",
				Size:         1024000,
				Type:         "file",
				IsDirectory:  false,
				ContentType:  "application/pdf",
				CreatedTime:  time.Now().Add(-24 * time.Hour),
				ModifiedTime: time.Now().Add(-24 * time.Hour),
			},
			{
				Name:         "images",
				Path:         "/home/user/documents/images",
				Size:         0,
				Type:         "directory",
				IsDirectory:  true,
				CreatedTime:  time.Now().Add(-48 * time.Hour),
				ModifiedTime: time.Now().Add(-48 * time.Hour),
			},
			{
				Name:         "notes.txt",
				Path:         "/home/user/documents/notes.txt",
				Size:         2048,
				Type:         "file",
				IsDirectory:  false,
				ContentType:  "text/plain",
				CreatedTime:  time.Now().Add(-12 * time.Hour),
				ModifiedTime: time.Now().Add(-12 * time.Hour),
			},
			{
				Name:         "readme.md",
				Path:         "/home/user/documents/readme.md",
				Size:         512,
				Type:         "file",
				IsDirectory:  false,
				ContentType:  "text/markdown",
				CreatedTime:  time.Now().Add(-6 * time.Hour),
				ModifiedTime: time.Now().Add(-6 * time.Hour),
			},
		},
	}
	err = saveToolMsg("tool_4", "lobe-local-system", "listLocalFiles", tool4ArgsJSON,
		nil,            // content
		listFilesState, // pluginState
		5)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 5: Local System - readLocalFile
	// Using proper types from builtin.LocalReadFileState
	readFileState := &builtin.LocalReadFileState{
		FileContent: localfs.LocalReadFileResult{
			Content:        "# README\n\nThis is a sample readme file.\n\n## Features\n- Feature 1\n- Feature 2\n- Feature 3",
			Filename:       "readme.md",
			FileType:       "text/markdown",
			CharCount:      120,
			LineCount:      8,
			TotalCharCount: 120,
			TotalLineCount: 8,
			Loc:            [2]int{0, 100},
			CreatedTime:    time.Now().Add(-6 * time.Hour),
			ModifiedTime:   time.Now().Add(-6 * time.Hour),
		},
	}
	err = saveToolMsg("tool_5", "lobe-local-system", "readLocalFile", tool5ArgsJSON,
		nil,           // content
		readFileState, // pluginState
		6)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 6: Local System - searchLocalFiles
	// Using proper types from builtin.LocalFileSearchState
	searchFilesState := &builtin.LocalFileSearchState{
		SearchResults: []localfs.LocalFileItem{
			{
				Name:         "important_doc.pdf",
				Path:         "/home/user/documents/important_doc.pdf",
				Size:         2048000,
				Type:         "file",
				IsDirectory:  false,
				ContentType:  "application/pdf",
				CreatedTime:  time.Now().Add(-72 * time.Hour),
				ModifiedTime: time.Now().Add(-72 * time.Hour),
			},
			{
				Name:         "important_notes.txt",
				Path:         "/home/user/documents/important_notes.txt",
				Size:         1024,
				Type:         "file",
				IsDirectory:  false,
				ContentType:  "text/plain",
				CreatedTime:  time.Now().Add(-36 * time.Hour),
				ModifiedTime: time.Now().Add(-36 * time.Hour),
			},
		},
	}
	err = saveToolMsg("tool_6", "lobe-local-system", "searchLocalFiles", tool6ArgsJSON,
		nil,              // content
		searchFilesState, // pluginState
		7)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 7: Local System - writeLocalFile
	// Using proper types from builtin.WriteFileResult
	writeFileResult := &localfs.WriteFileResult{
		Path:    "/home/user/documents/new_file.txt",
		Success: true,
		Message: "File written successfully",
	}
	err = saveToolMsg("tool_7", "lobe-local-system", "writeLocalFile", tool7ArgsJSON,
		writeFileResult, nil, 8)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 8: Local System - renameLocalFile
	// Using proper types from builtin.LocalRenameFileState
	renameFileState := &builtin.LocalRenameFileState{
		OldPath: "/home/user/documents/old_name.txt",
		NewPath: "/home/user/documents/new_name.txt",
		Success: true,
	}
	err = saveToolMsg("tool_8", "lobe-local-system", "renameLocalFile", tool8ArgsJSON,
		renameFileState, nil, 9)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 9: Local System - moveLocalFiles
	// Using proper types from builtin.LocalMoveFilesState
	moveFilesState := &builtin.LocalMoveFilesState{
		Results: []localfs.LocalMoveFilesResultItem{
			{
				SourcePath: "/home/user/documents/file1.txt",
				NewPath:    "/home/user/backup/file1.txt",
				Success:    true,
			},
		},
		SuccessCount: 1,
		TotalCount:   1,
	}
	err = saveToolMsg("tool_9", "lobe-local-system", "moveLocalFiles", tool9ArgsJSON,
		moveFilesState, nil, 10)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 10: DALL-E Image Designer - text2image
	// Using proper types from builtin.DallEImageItem
	// ImageId corresponds to actual files in data/files/ directory
	dalleImages := []builtin.DallEImageItem{
		{
			Prompt:     "A beautiful sunset over a calm ocean with vibrant orange and purple colors",
			PreviewUrl: "https://picsum.photos/seed/sunset/1024/1024",
			ImageId:    "501-1024x1024.jpg", // Actual file in data/files/
			Quality:    "hd",
			Size:       "1024x1024",
			Style:      "vivid",
		},
		{
			Prompt:     "A futuristic cityscape at night with neon lights and flying cars",
			PreviewUrl: "https://picsum.photos/seed/cityscape/1024/1024",
			ImageId:    "622-1024x1024.jpg", // Actual file in data/files/
			Quality:    "hd",
			Size:       "1024x1024",
			Style:      "vivid",
		},
	}
	err = saveToolMsg("tool_10", "lobe-image-designer", "text2image", tool10ArgsJSON,
		dalleImages, nil, 11)
	if err != nil {
		log.Printf("⚠️  %v", err)
	}

	// Tool 11: Code Interpreter - python
	// Using proper types from builtin.CodeInterpreterResponse
	codeInterpreterResponse := &builtin.CodeInterpreterResponse{
		Result: "{'a': 6, 'b': 15}",
		Output: []builtin.CodeInterpreterOutput{
			{
				Type: "stdout",
				Data: "              a         b\ncount  3.000000  3.000000\nmean   2.000000  5.000000\nstd    1.000000  1.000000\nmin    1.000000  4.000000\n25%    1.500000  4.500000\n50%    2.000000  5.000000\n75%    2.500000  5.500000\nmax    3.000000  6.000000\n",
			},
		},
		Files: []builtin.CodeInterpreterFileItem{},
	}
	err = saveToolMsg("tool_11", "lobe-code-interpreter", "python", tool11ArgsJSON,
		codeInterpreterResponse, // content
		nil,                     // pluginState (no error)
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
		CreatedAt:    time.Now().UnixMilli(),
	}, nil
}
