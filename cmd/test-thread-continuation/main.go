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

	// Test the exact scenario from backend-dev.log
	// where "more detail" in a thread continuation caused issues
	testScenario := []struct {
		msg      string
		response string // Expected pattern
	}{
		{msg: "what is AI?", response: ""},
		{msg: "elaborate", response: ""},
		{msg: "more detail", response: ""}, // This is where the issue occurred
	}

	modelPath := filepath.Join(modelsDir, "Qwen_Qwen3-1.7B-Q4_K_M.gguf")

	fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Printf("Testing Thread Continuation Scenario with Qwen3-1.7B\n")
	fmt.Printf("Simulating the exact flow from backend-dev.log\n")
	fmt.Printf(strings.Repeat("=", 80) + "\n\n")

	libService, err := llama.NewLibraryService()
	if err != nil {
		log.Fatalf("❌ Failed to create library service: %v\n", err)
	}
	defer libService.Cleanup()

	// Wait for initialization
	time.Sleep(3 * time.Second)

	if err := libService.InitializeLibrary(); err != nil {
		log.Fatalf("❌ Failed to initialize library: %v\n", err)
	}

	if err := libService.LoadChatModel(modelPath); err != nil {
		log.Fatalf("❌ Failed to load model: %v\n", err)
	}

	llamaModel := llama.NewLlamaEinoModel(libService)
	ctx := context.Background()

	// Build conversation history progressively (like in real chat)
	var history []*schema.Message
	issueCount := 0

	for i, turn := range testScenario {
		fmt.Printf("\n" + strings.Repeat("-", 80) + "\n")
		fmt.Printf("Turn %d: %s\n", i+1, turn.msg)
		fmt.Printf(strings.Repeat("-", 80) + "\n")

		// Add user message
		history = append(history, &schema.Message{
			Role:    schema.User,
			Content: turn.msg,
		})

		fmt.Printf("\n📝 Current history length: %d messages\n", len(history))

		// Generate response with accumulated history
		response, err := llamaModel.Generate(ctx, history)
		if err != nil {
			fmt.Printf("❌ Error: %v\n\n", err)
			continue
		}

		// Add assistant response to history
		history = append(history, &schema.Message{
			Role:    schema.Assistant,
			Content: response.Content,
		})

		// Display response (truncated)
		displayContent := response.Content
		if len(displayContent) > 300 {
			displayContent = displayContent[:300] + "... [truncated]"
		}
		fmt.Printf("\n🤖 Response:\n%s\n", displayContent)

		// Analyze for think tag issues
		hasThinkOpen := strings.Contains(response.Content, "<think>")
		hasThinkClose := strings.Contains(response.Content, "</think>")

		fmt.Printf("\n🔍 Analysis:\n")
		fmt.Printf("   - Has <think> tag: %v\n", hasThinkOpen)
		fmt.Printf("   - Has </think> tag: %v\n", hasThinkClose)
		fmt.Printf("   - Response length: %d chars\n", len(response.Content))

		if hasThinkOpen {
			if !hasThinkClose {
				fmt.Printf("\n❌ ISSUE DETECTED: Unclosed <think> tag!\n")
				fmt.Printf("   This matches the problem from backend-dev.log\n")
				issueCount++
			} else {
				thinkEnd := strings.Index(response.Content, "</think>")
				afterThink := strings.TrimSpace(response.Content[thinkEnd+8:])

				if afterThink == "" {
					fmt.Printf("\n⚠️  ISSUE: No content after </think> tag!\n")
					issueCount++
				} else {
					fmt.Printf("\n✅ Think tags properly formatted\n")
					fmt.Printf("   Content after </think>: %d chars\n", len(afterThink))
				}
			}
		} else {
			fmt.Printf("\n✅ No think tags (clean response)\n")
		}
	}

	fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Printf("FINAL RESULTS\n")
	fmt.Printf(strings.Repeat("=", 80) + "\n")
	fmt.Printf("Total turns: %d\n", len(testScenario))
	fmt.Printf("Issues found: %d\n", issueCount)
	
	if issueCount == 0 {
		fmt.Printf("\n✅ SUCCESS: No think tag issues detected!\n")
		fmt.Printf("   The model handled thread continuation properly.\n")
	} else {
		fmt.Printf("\n❌ FAILURE: Think tag issues detected!\n")
		fmt.Printf("   The model has the same problem as seen in backend-dev.log\n")
	}
	fmt.Printf(strings.Repeat("=", 80) + "\n")
}

