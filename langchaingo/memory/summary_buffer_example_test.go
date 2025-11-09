package memory_test

import (
	"context"
	"fmt"
	"log"

	"github.com/kawai-network/veridium/langchaingo/llms/ollama"
	"github.com/kawai-network/veridium/langchaingo/memory"
)

// ExampleConversationSummaryBuffer demonstrates how to use ConversationSummaryBuffer
// to automatically summarize long conversations and prevent LLM confusion.
func ExampleConversationSummaryBuffer() {
	ctx := context.Background()

	// Create an LLM (using Ollama as example)
	llm, err := ollama.New(
		ollama.WithServerURL("http://localhost:11434"),
		ollama.WithModel("llama2"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create a ConversationSummaryBuffer
	// It will:
	// - Keep only the last 3 conversation pairs (6 messages)
	// - Automatically summarize older messages when > 6 pairs
	// - Run summarization in background (non-blocking)
	summaryBuffer := memory.NewConversationSummaryBuffer(
		llm,
		memory.WithReturnMessages(true),
		memory.WithMemoryKey("history"),
	)

	// Simulate a long conversation
	conversations := []struct {
		user string
		ai   string
	}{
		{"What is machine learning?", "Machine learning is a subset of AI..."},
		{"How does it work?", "It works by training models on data..."},
		{"What are neural networks?", "Neural networks are computing systems..."},
		{"Tell me about deep learning", "Deep learning uses multiple layers..."},
		{"What is NLP?", "NLP stands for Natural Language Processing..."},
		{"How does GPT work?", "GPT uses transformer architecture..."},
		{"What about fine-tuning?", "Fine-tuning adapts pre-trained models..."},
		// After this point, old messages will be summarized
		{"Can you explain RAG?", "RAG stands for Retrieval Augmented Generation..."},
	}

	// Add all conversations
	for _, conv := range conversations {
		err := summaryBuffer.SaveContext(ctx,
			map[string]any{"input": conv.user},
			map[string]any{"output": conv.ai},
		)
		if err != nil {
			log.Printf("Error saving context: %v", err)
		}
	}

	// Load memory variables - will return summary + last 3 pairs
	memoryVars, err := summaryBuffer.LoadMemoryVariables(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	// The returned history will contain:
	// 1. A system message with the summary of old conversations
	// 2. The last 3 conversation pairs (6 messages)
	fmt.Printf("Memory contains summary + recent messages\n")
	fmt.Printf("Summary available: %v\n", summaryBuffer.GetSummary() != "")

	// This prevents LLM confusion by limiting context to relevant information
	_ = memoryVars
}

// ExampleNewConversationSummaryBufferWithOptions shows how to customize the buffer.
func ExampleNewConversationSummaryBufferWithOptions() {
	ctx := context.Background()

	llm, err := ollama.New(
		ollama.WithServerURL("http://localhost:11434"),
		ollama.WithModel("llama2"),
	)
	if err != nil {
		log.Fatal(err)
	}

	// Create with custom options
	summaryOpts := []memory.SummaryBufferOption{
		memory.WithRecentMessagePairs(5),    // Keep last 5 pairs instead of 3
		memory.WithMaxTokensForRecent(2000), // Allow up to 2K tokens for recent
		memory.WithSummarizeThreshold(8),    // Summarize after 8 pairs instead of 6
	}

	summaryBuffer := memory.NewConversationSummaryBufferWithOptions(
		llm,
		summaryOpts,
		memory.WithReturnMessages(true),
		memory.WithMemoryKey("history"),
	)

	// Use the buffer
	err = summaryBuffer.SaveContext(ctx,
		map[string]any{"input": "Hello"},
		map[string]any{"output": "Hi there!"},
	)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Custom summary buffer created")
}
