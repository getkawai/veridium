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

func testModel(modelPath, modelName string, useNoThink bool) (int, int, time.Duration, error) {
	fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Printf("Testing: %s\n", modelName)
	if useNoThink {
		fmt.Printf("Mode: With /no_think instruction\n")
	} else {
		fmt.Printf("Mode: Default (no special instruction)\n")
	}
	fmt.Printf(strings.Repeat("=", 80) + "\n\n")

	libService, err := llama.NewLibraryService()
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to create service: %w", err)
	}
	defer libService.Cleanup()

	time.Sleep(3 * time.Second)

	if err := libService.InitializeLibrary(); err != nil {
		return 0, 0, 0, fmt.Errorf("failed to initialize: %w", err)
	}

	if err := libService.LoadChatModel(modelPath); err != nil {
		return 0, 0, 0, fmt.Errorf("failed to load model: %w", err)
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
	
	// Add system message
	systemMsg := "You are a helpful AI assistant. Be concise and direct in your responses."
	if useNoThink {
		systemMsg += "\n\n/no_think\n\nIMPORTANT: Do NOT use <think> tags. Provide direct answers without showing your reasoning process."
	}
	
	history = append(history, &schema.Message{
		Role:    schema.System,
		Content: systemMsg,
	})

	successCount := 0
	totalTokensEstimate := 0
	totalTime := time.Duration(0)
	thinkTagCount := 0

	for i, question := range testQuestions {
		fmt.Printf("\n" + strings.Repeat("-", 80) + "\n")
		fmt.Printf("Turn %d/%d: %s\n", i+1, len(testQuestions), question)
		fmt.Printf(strings.Repeat("-", 80) + "\n")

		history = append(history, &schema.Message{
			Role:    schema.User,
			Content: question,
		})

		totalTokensEstimate += len(question) / 4

		// Generate response
		startTime := time.Now()
		response, err := llamaModel.Generate(ctx, history)
		duration := time.Since(startTime)
		totalTime += duration

		if err != nil {
			fmt.Printf("❌ ERROR: %v\n", err)
			break
		}

		history = append(history, &schema.Message{
			Role:    schema.Assistant,
			Content: response.Content,
		})

		totalTokensEstimate += len(response.Content) / 4

		// Check for think tags
		hasThinkTags := strings.Contains(response.Content, "<think>")
		if hasThinkTags {
			thinkTagCount++
		}

		// Display response
		displayContent := response.Content
		if len(displayContent) > 200 {
			displayContent = displayContent[:200] + "..."
		}

		fmt.Printf("🤖 Response (%d chars): %s\n", len(response.Content), displayContent)
		fmt.Printf("⏱️  Time: %v\n", duration)
		fmt.Printf("📈 Tokens: ~%d / 16,384 (%.1f%%)\n", totalTokensEstimate, float64(totalTokensEstimate)/16384.0*100)
		
		if hasThinkTags {
			fmt.Printf("⚠️  Has think tags\n")
		} else {
			fmt.Printf("✅ No think tags\n")
		}

		successCount++

		if totalTokensEstimate > 14000 {
			fmt.Printf("⚠️  Approaching context limit\n")
			break
		}
	}

	fmt.Printf("\n" + strings.Repeat("-", 80) + "\n")
	fmt.Printf("RESULTS for %s:\n", modelName)
	fmt.Printf(strings.Repeat("-", 80) + "\n")
	fmt.Printf("Successful turns: %d / %d\n", successCount, len(testQuestions))
	fmt.Printf("Total tokens: ~%d / 16,384 (%.1f%%)\n", totalTokensEstimate, float64(totalTokensEstimate)/16384.0*100)
	fmt.Printf("Average time: %v\n", totalTime/time.Duration(successCount))
	fmt.Printf("Total time: %v\n", totalTime)
	fmt.Printf("Think tags found: %d / %d responses\n", thinkTagCount, successCount)
	fmt.Printf(strings.Repeat("-", 80) + "\n")

	return successCount, totalTokensEstimate, totalTime, nil
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Failed to get home directory: %v", err)
	}

	modelsDir := filepath.Join(homeDir, ".llama-cpp", "models")

	fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Printf("COMPARISON: Llama 3.2 vs Qwen3 (with and without /no_think)\n")
	fmt.Printf(strings.Repeat("=", 80) + "\n")

	// Test 1: Llama 3.2 (no thinking capability)
	llama32Path := filepath.Join(modelsDir, "Llama-3.2-3B-Instruct-Q4_K_M.gguf")
	llama32Turns, llama32Tokens, llama32Time, err := testModel(llama32Path, "Llama 3.2 3B", false)
	if err != nil {
		log.Printf("❌ Llama 3.2 test failed: %v", err)
	}

	time.Sleep(2 * time.Second)

	// Test 2: Qwen3 with default behavior
	qwen3Path := filepath.Join(modelsDir, "Qwen_Qwen3-1.7B-Q4_K_M.gguf")
	qwen3DefaultTurns, qwen3DefaultTokens, qwen3DefaultTime, err := testModel(qwen3Path, "Qwen3 1.7B (Default)", false)
	if err != nil {
		log.Printf("❌ Qwen3 default test failed: %v", err)
	}

	time.Sleep(2 * time.Second)

	// Test 3: Qwen3 with /no_think
	qwen3NoThinkTurns, qwen3NoThinkTokens, qwen3NoThinkTime, err := testModel(qwen3Path, "Qwen3 1.7B (/no_think)", true)
	if err != nil {
		log.Printf("❌ Qwen3 /no_think test failed: %v", err)
	}

	// Final comparison
	fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Printf("FINAL COMPARISON\n")
	fmt.Printf(strings.Repeat("=", 80) + "\n\n")

	fmt.Printf("┌─────────────────────────────┬───────────┬──────────┬─────────────┬──────────┐\n")
	fmt.Printf("│ Model                       │ Turns     │ Tokens   │ Avg Time    │ Total    │\n")
	fmt.Printf("├─────────────────────────────┼───────────┼──────────┼─────────────┼──────────┤\n")
	
	if llama32Turns > 0 {
		fmt.Printf("│ Llama 3.2 3B                │ %2d / 10   │ ~%-6d  │ %8s    │ %7s  │\n", 
			llama32Turns, llama32Tokens, 
			(llama32Time/time.Duration(llama32Turns)).Round(time.Millisecond),
			llama32Time.Round(time.Millisecond))
	}
	
	if qwen3DefaultTurns > 0 {
		fmt.Printf("│ Qwen3 1.7B (Default)        │ %2d / 10   │ ~%-6d  │ %8s    │ %7s  │\n", 
			qwen3DefaultTurns, qwen3DefaultTokens,
			(qwen3DefaultTime/time.Duration(qwen3DefaultTurns)).Round(time.Millisecond),
			qwen3DefaultTime.Round(time.Millisecond))
	}
	
	if qwen3NoThinkTurns > 0 {
		fmt.Printf("│ Qwen3 1.7B (/no_think)      │ %2d / 10   │ ~%-6d  │ %8s    │ %7s  │\n", 
			qwen3NoThinkTurns, qwen3NoThinkTokens,
			(qwen3NoThinkTime/time.Duration(qwen3NoThinkTurns)).Round(time.Millisecond),
			qwen3NoThinkTime.Round(time.Millisecond))
	}
	
	fmt.Printf("└─────────────────────────────┴───────────┴──────────┴─────────────┴──────────┘\n\n")

	// Analysis
	fmt.Printf("📊 Analysis:\n\n")
	
	if llama32Turns > 0 && qwen3NoThinkTurns > 0 {
		turnsRatio := float64(llama32Turns) / float64(qwen3NoThinkTurns)
		tokensRatio := float64(llama32Tokens) / float64(qwen3NoThinkTokens)
		timeRatio := float64(llama32Time) / float64(qwen3NoThinkTime)
		
		fmt.Printf("Llama 3.2 vs Qwen3 (/no_think):\n")
		fmt.Printf("  - Turns: %.2fx\n", turnsRatio)
		fmt.Printf("  - Token efficiency: %.2fx\n", tokensRatio)
		fmt.Printf("  - Speed: %.2fx\n", timeRatio)
		fmt.Printf("\n")
	}
	
	if qwen3DefaultTurns > 0 && qwen3NoThinkTurns > 0 {
		turnsRatio := float64(qwen3NoThinkTurns) / float64(qwen3DefaultTurns)
		tokensRatio := float64(qwen3DefaultTokens) / float64(qwen3NoThinkTokens)
		timeRatio := float64(qwen3DefaultTime) / float64(qwen3NoThinkTime)
		
		fmt.Printf("Qwen3 /no_think vs Default:\n")
		fmt.Printf("  - Turns improvement: %.2fx\n", turnsRatio)
		fmt.Printf("  - Token reduction: %.2fx\n", tokensRatio)
		fmt.Printf("  - Speed improvement: %.2fx\n", timeRatio)
		fmt.Printf("\n")
	}

	fmt.Printf("🏆 Winner: ")
	if llama32Turns >= qwen3NoThinkTurns && llama32Turns >= qwen3DefaultTurns {
		fmt.Printf("Llama 3.2 3B (native non-thinking model)\n")
	} else if qwen3NoThinkTurns >= llama32Turns && qwen3NoThinkTurns >= qwen3DefaultTurns {
		fmt.Printf("Qwen3 1.7B with /no_think instruction\n")
	} else {
		fmt.Printf("Qwen3 1.7B default\n")
	}

	fmt.Printf(strings.Repeat("=", 80) + "\n")
}

