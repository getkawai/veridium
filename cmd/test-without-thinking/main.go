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
	fmt.Printf("Testing Qwen3 WITHOUT Think Tags\n")
	fmt.Printf("Using /no_think instruction\n")
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

	// Test questions
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
	}

	var history []*schema.Message
	
	// Add system message with /no_think instruction
	history = append(history, &schema.Message{
		Role:    schema.System,
		Content: "You are a helpful AI assistant. Be concise and direct in your responses.\n\n/no_think\n\nIMPORTANT: Do NOT use <think> tags. Provide direct answers without showing your reasoning process.",
	})

	successCount := 0
	totalTokensEstimate := 0
	totalTime := time.Duration(0)

	for i, question := range testQuestions {
		fmt.Printf("\n" + strings.Repeat("-", 80) + "\n")
		fmt.Printf("Turn %d/%d: %s\n", i+1, len(testQuestions), question)
		fmt.Printf(strings.Repeat("-", 80) + "\n")

		// Add user message
		history = append(history, &schema.Message{
			Role:    schema.User,
			Content: question,
		})

		totalTokensEstimate += len(question) / 4

		fmt.Printf("📊 History: %d messages, ~%d tokens\n", len(history), totalTokensEstimate)

		// Generate response
		startTime := time.Now()
		response, err := llamaModel.Generate(ctx, history)
		duration := time.Since(startTime)
		totalTime += duration

		if err != nil {
			fmt.Printf("❌ ERROR: %v\n", err)
			fmt.Printf("⚠️  Failed at turn %d\n", i+1)
			break
		}

		// Add assistant response to history
		history = append(history, &schema.Message{
			Role:    schema.Assistant,
			Content: response.Content,
		})

		totalTokensEstimate += len(response.Content) / 4

		// Check for think tags
		hasThinkTags := strings.Contains(response.Content, "<think>")
		
		status := "✅ No think tags"
		if hasThinkTags {
			status = "⚠️  Still has think tags"
		}

		// Display response
		displayContent := response.Content
		if len(displayContent) > 300 {
			displayContent = displayContent[:300] + "..."
		}

		fmt.Printf("🤖 Response (%d chars): %s\n", len(response.Content), displayContent)
		fmt.Printf("⏱️  Time: %v (avg: %v)\n", duration, totalTime/time.Duration(i+1))
		fmt.Printf("📈 Tokens: ~%d / 16,384 (%.1f%%)\n", totalTokensEstimate, float64(totalTokensEstimate)/16384.0*100)
		fmt.Printf("%s\n", status)

		successCount++

		// Check if approaching limit
		if totalTokensEstimate > 14000 {
			fmt.Printf("⚠️  Approaching context limit, stopping test\n")
			break
		}
	}

	fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Printf("FINAL RESULTS\n")
	fmt.Printf(strings.Repeat("=", 80) + "\n")
	fmt.Printf("Successful turns: %d / %d\n", successCount, len(testQuestions))
	fmt.Printf("Total messages: %d\n", len(history))
	fmt.Printf("Estimated tokens: ~%d / 16,384 (%.1f%%)\n", totalTokensEstimate, float64(totalTokensEstimate)/16384.0*100)
	fmt.Printf("Average time per turn: %v\n", totalTime/time.Duration(successCount))
	fmt.Printf("Total time: %v\n", totalTime)
	
	if successCount == len(testQuestions) {
		fmt.Printf("\n✅ SUCCESS: Completed all %d turns!\n", len(testQuestions))
	} else {
		fmt.Printf("\n⚠️  Completed %d/%d turns\n", successCount, len(testQuestions))
	}
	
	fmt.Printf("\n📊 Comparison with think tags:\n")
	fmt.Printf("   With think tags: 3 turns (~3,000 tokens, ~33s)\n")
	fmt.Printf("   Without think tags: %d turns (~%d tokens, ~%ds)\n", 
		successCount, totalTokensEstimate, int(totalTime.Seconds()))
	fmt.Printf("   Improvement: %.1fx more turns\n", float64(successCount)/3.0)
	fmt.Printf(strings.Repeat("=", 80) + "\n")
}

