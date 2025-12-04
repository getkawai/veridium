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
	"fmt"
	"strings"

	"github.com/kawai-network/veridium/types"
)

// RAGWorkflow implements Retrieval-Augmented Generation using Eino
type RAGWorkflow struct {
	kbService *KnowledgeBaseService
}

// RAGRequest represents a RAG query request
type RAGRequest struct {
	Query           string
	KnowledgeBaseID string
	UserID          string
	TopK            int
	IncludeSources  bool
}

// RAGResponse represents the RAG result
type RAGResponse struct {
	Context         string
	Sources         []*types.Document
	RetrievedChunks int
}

// NewRAGWorkflow creates a new RAG workflow
func NewRAGWorkflow(kbService *KnowledgeBaseService) *RAGWorkflow {
	return &RAGWorkflow{
		kbService: kbService,
	}
}

// BuildContext retrieves relevant documents and builds context for LLM
func (w *RAGWorkflow) BuildContext(ctx context.Context, req RAGRequest) (string, []*types.Document, error) {
	// 1. Get retriever
	retriever, err := w.kbService.GetRetriever(ctx, req.KnowledgeBaseID, req.UserID)
	if err != nil {
		return "", nil, fmt.Errorf("failed to get retriever: %w", err)
	}

	// 2. Retrieve relevant documents
	topK := req.TopK
	if topK <= 0 {
		topK = 5
	}

	docs, err := retriever(ctx, req.Query)
	if err != nil {
		return "", nil, fmt.Errorf("retrieval failed: %w", err)
	}

	// Limit to topK
	if len(docs) > topK {
		docs = docs[:topK]
	}

	// 3. Build context from documents
	var contextBuilder strings.Builder
	contextBuilder.WriteString("Here are relevant documents from the knowledge base:\n\n")

	for i, doc := range docs {
		contextBuilder.WriteString(fmt.Sprintf("Document %d:\n", i+1))
		contextBuilder.WriteString(doc.Content)
		contextBuilder.WriteString("\n\n")

		// Add source info if available
		if source, ok := doc.Metadata["source_file"].(string); ok {
			contextBuilder.WriteString(fmt.Sprintf("Source: %s\n", source))
		}
		contextBuilder.WriteString("---\n\n")
	}

	return contextBuilder.String(), docs, nil
}

// FormatContextForLLM formats retrieved documents for LLM context
func (w *RAGWorkflow) FormatContextForLLM(docs []*types.Document) string {
	var builder strings.Builder
	builder.WriteString("Here are relevant documents from the knowledge base:\n\n")

	for i, doc := range docs {
		builder.WriteString(fmt.Sprintf("Document %d:\n", i+1))
		builder.WriteString(doc.Content)
		builder.WriteString("\n\n")

		// Add source info if available
		if source, ok := doc.Metadata["source_file"].(string); ok {
			builder.WriteString(fmt.Sprintf("Source: %s\n", source))
		}
		builder.WriteString("---\n\n")
	}

	return builder.String()
}

// ExecuteRAG executes the RAG workflow and returns the response
func (w *RAGWorkflow) ExecuteRAG(ctx context.Context, req RAGRequest) (*RAGResponse, error) {
	// Build context using retriever
	contextStr, docs, err := w.BuildContext(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to build context: %w", err)
	}

	// Return response
	response := &RAGResponse{
		Context:         contextStr,
		Sources:         docs,
		RetrievedChunks: len(docs),
	}

	return response, nil
}
