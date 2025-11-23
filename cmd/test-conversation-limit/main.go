package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cloudwego/eino/schema"
	"github.com/kawai-network/veridium/internal/llama"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	modelsDir := filepath.Join(homeDir, ".llama-cpp", "models")
	modelPath := filepath.Join(modelsDir, "Qwen_Qwen3-1.7B-Q4_K_M.gguf")

	fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Printf("Testing Conversation Length Limit\n")
	fmt.Printf("Context Window: 16,384 tokens\n")
	fmt.Printf(strings.Repeat("=", 80) + "\n\n")

	libService, err := llama.NewLibraryService()
	if err != nil {
		log.Fatalf("❌ Failed to create library service: %v\n", err)
	}
	defer libService.Cleanup()

	time.Sleep(3 * time.Second)

	if err := libService.InitializeLibrary(); err != nil {
		log.Fatalf("❌ Failed to initialize library: %v\n", err)
	}

	if err := libService.LoadChatModel(modelPath); err != nil {
		log.Fatalf("❌ Failed to load model: %v\n", err)
	}

	llamaModel := llama.NewLlamaEinoModel(libService)
	ctx := context.Background()

	// Test questions - progressively longer conversation
	testQuestions := []string{
		"What is AI?",
		"Tell me about machine learning",
		"Explain neural networks",
		"What are transformers?",
		"How does attention mechanism work?",
		"Explain GPT architecture",
		"What is BERT?",
		"Tell me about computer vision",
		"Explain reinforcement learning",
		"What is natural language processing?",
		"Tell me about deep learning",
		"Explain convolutional neural networks",
		"What are recurrent neural networks?",
		"Tell me about GANs",
		"Explain transfer learning",
		"What is fine-tuning?",
		"Tell me about prompt engineering",
		"Explain few-shot learning",
		"What is zero-shot learning?",
		"Tell me about model compression",
	}

	var history []*schema.Message
	successCount := 0
	totalTokensEstimate := 0

	for i, question := range testQuestions {
		fmt.Printf("\n" + strings.Repeat("-", 80) + "\n")
		fmt.Printf("Turn %d/%d: %s\n", i+1, len(testQuestions), question)
		fmt.Printf(strings.Repeat("-", 80) + "\n")

		// Add user message
		history = append(history, &schema.Message{
			Role:    schema.User,
			Content: question,
		})

		// Estimate tokens (rough: 1 token ≈ 4 chars)
		totalTokensEstimate += len(question) / 4

		fmt.Printf("📊 History: %d messages, ~%d tokens\n", len(history), totalTokensEstimate)

		// Generate response
		startTime := time.Now()
		response, err := llamaModel.Generate(ctx, history)
		duration := time.Since(startTime)

		if err != nil {
			fmt.Printf("❌ ERROR: %v\n", err)
			fmt.Printf("⚠️  Failed at turn %d with ~%d tokens in history\n", i+1, totalTokensEstimate)
			break
		}

		// Add assistant response to history
		history = append(history, &schema.Message{
			Role:    schema.Assistant,
			Content: response.Content,
		})

		totalTokensEstimate += len(response.Content) / 4

		// Check for think tag issues
		hasThinkOpen := strings.Contains(response.Content, "<think>")
		hasThinkClose := strings.Contains(response.Content, "</think>")
		
		status := "✅"
		if hasThinkOpen && !hasThinkClose {
			status = "❌ UNCLOSED THINK TAG"
		}

		// Display truncated response
		displayContent := response.Content
		if len(displayContent) > 200 {
			displayContent = displayContent[:200] + "..."
		}

		fmt.Printf("🤖 Response (%d chars): %s\n", len(response.Content), displayContent)
		fmt.Printf("⏱️  Generation time: %v\n", duration)
		fmt.Printf("📈 Total tokens estimate: ~%d / 16,384\n", totalTokensEstimate)
		fmt.Printf("📊 Context usage: %.1f%%\n", float64(totalTokensEstimate)/16384.0*100)
		fmt.Printf("%s Status\n", status)

		successCount++

		// Warning if approaching limit
		if totalTokensEstimate > 14000 {
			fmt.Printf("⚠️  WARNING: Approaching context limit!\n")
		}
	}

	fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Printf("FINAL RESULTS\n")
	fmt.Printf(strings.Repeat("=", 80) + "\n")
	fmt.Printf("Successful turns: %d / %d\n", successCount, len(testQuestions))
	fmt.Printf("Total messages: %d\n", len(history))
	fmt.Printf("Estimated tokens used: ~%d / 16,384\n", totalTokensEstimate)
	fmt.Printf("Context utilization: %.1f%%\n", float64(totalTokensEstimate)/16384.0*100)
	
	if successCount == len(testQuestions) {
		fmt.Printf("\n✅ SUCCESS: Handled all %d conversation turns!\n", len(testQuestions))
	} else {
		fmt.Printf("\n⚠️  Reached limit at turn %d\n", successCount)
	}
	fmt.Printf(strings.Repeat("=", 80) + "\n")
}

