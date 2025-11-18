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
	"encoding/json"
	"fmt"

	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/model"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

// RAGAgent is an AI agent with RAG capabilities
type RAGAgent struct {
	agent       adk.Agent
	kbService   *KnowledgeBaseService
	ragWorkflow *RAGWorkflow
}

// RAGAgentConfig holds configuration for RAG agent
type RAGAgentConfig struct {
	Name             string
	Description      string
	Model            model.ToolCallingChatModel
	KnowledgeBaseIDs []string
	UserID           string
	Instruction      string
	MaxIterations    int
}

// NewRAGAgent creates a new agent with RAG capabilities
func NewRAGAgent(ctx context.Context, config *RAGAgentConfig, kbService *KnowledgeBaseService) (*RAGAgent, error) {
	ragWorkflow := NewRAGWorkflow(kbService)

	// Set default values
	if config.MaxIterations <= 0 {
		config.MaxIterations = 10
	}

	if config.Description == "" {
		config.Description = "AI assistant with knowledge base access"
	}

	if config.Instruction == "" {
		config.Instruction = "You are a helpful assistant with access to knowledge bases. Use the search tools to find relevant information before answering questions."
	}

	// Create knowledge base search tools for each KB
	var tools []tool.BaseTool
	for _, kbID := range config.KnowledgeBaseIDs {
		kbTool, err := createKnowledgeBaseTool(ctx, kbID, config.UserID, kbService)
		if err != nil {
			return nil, fmt.Errorf("failed to create KB tool for %s: %w", kbID, err)
		}
		tools = append(tools, kbTool)
	}

	// Create agent with tools
	agent, err := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name:        config.Name,
		Description: config.Description,
		Instruction: config.Instruction,
		Model:       config.Model,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: tools,
			},
		},
		MaxIterations: config.MaxIterations,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create agent: %w", err)
	}

	return &RAGAgent{
		agent:       agent,
		kbService:   kbService,
		ragWorkflow: ragWorkflow,
	}, nil
}

// Run executes the agent with RAG capabilities
func (a *RAGAgent) Run(ctx context.Context, userMessage string) (string, error) {
	// Create input messages
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: userMessage,
		},
	}

	// Create agent input
	input := &adk.AgentInput{
		Messages: messages,
	}

	// Run agent
	iterator := a.agent.Run(ctx, input)

	// Collect all events and build final response
	var finalResponse string
	for {
		event, ok := iterator.Next()
		if !ok {
			break
		}

		if event.Err != nil {
			return "", fmt.Errorf("agent execution failed: %w", event.Err)
		}

		// Extract message content from event
		if event.Output != nil && event.Output.MessageOutput != nil {
			msgVariant := event.Output.MessageOutput
			if msgVariant.Message != nil && msgVariant.Role == schema.Assistant {
				finalResponse += msgVariant.Message.Content
			}
		}
	}

	return finalResponse, nil
}

// Stream executes the agent with streaming responses
func (a *RAGAgent) Stream(ctx context.Context, userMessage string) (*adk.AsyncIterator[*adk.AgentEvent], error) {
	// Create input messages
	messages := []*schema.Message{
		{
			Role:    schema.User,
			Content: userMessage,
		},
	}

	// Create agent input
	input := &adk.AgentInput{
		Messages: messages,
	}

	// Run agent and return iterator
	return a.agent.Run(ctx, input), nil
}

// createKnowledgeBaseTool creates a tool for querying a knowledge base
func createKnowledgeBaseTool(ctx context.Context, kbID, userID string, kbService *KnowledgeBaseService) (tool.BaseTool, error) {
	// Get KB info from database
	kb, err := kbService.GetKnowledgeBase(ctx, kbID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get knowledge base: %w", err)
	}

	// Create the search function
	searchFunc := func(ctx context.Context, req *searchRequest) (*searchResponse, error) {
		topK := req.TopK
		if topK <= 0 {
			topK = 5
		}

		docs, err := kbService.QueryKnowledgeBase(ctx, kbID, req.Query, topK, userID)
		if err != nil {
			return nil, err
		}

		var results []string
		var sources []string
		for _, doc := range docs {
			results = append(results, doc.Content)
			if source, ok := doc.MetaData["source_file"].(string); ok {
				sources = append(sources, source)
			}
		}

		return &searchResponse{
			Results: results,
			Sources: sources,
		}, nil
	}

	// Create tool using BaseTool interface
	kbTool := &knowledgeBaseTool{
		name:        fmt.Sprintf("search_kb_%s", kb.Name),
		description: fmt.Sprintf("Search knowledge base: %s", kb.Name),
		searchFunc:  searchFunc,
	}

	return kbTool, nil
}

// knowledgeBaseTool implements tool.BaseTool interface
type knowledgeBaseTool struct {
	name        string
	description string
	searchFunc  func(ctx context.Context, req *searchRequest) (*searchResponse, error)
}

type searchRequest struct {
	Query string `json:"query" jsonschema:"description=Search query for the knowledge base"`
	TopK  int    `json:"top_k,omitempty" jsonschema:"description=Number of results to return (default: 5)"`
}

type searchResponse struct {
	Results []string `json:"results"`
	Sources []string `json:"sources"`
}

func (t *knowledgeBaseTool) Info(ctx context.Context) (*schema.ToolInfo, error) {
	// Define parameters using Eino's ParameterInfo
	params := map[string]*schema.ParameterInfo{
		"query": {
			Type: "string",
			Desc: "Search query for the knowledge base",
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

func (t *knowledgeBaseTool) InvokableRun(ctx context.Context, argumentsInJSON string, opts ...tool.Option) (string, error) {
	// Parse arguments
	var req searchRequest
	if err := json.Unmarshal([]byte(argumentsInJSON), &req); err != nil {
		return "", fmt.Errorf("failed to parse arguments: %w", err)
	}

	// Execute search
	resp, err := t.searchFunc(ctx, &req)
	if err != nil {
		return "", fmt.Errorf("search failed: %w", err)
	}

	// Format response as JSON
	respJSON, err := json.Marshal(resp)
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	return string(respJSON), nil
}
