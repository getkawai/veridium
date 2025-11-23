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

	// Test cases - simplified to avoid buffer overflow
	testMessages := []string{
		"what is AI?",
		"tell me more",
	}

	// Models to test
	models := []struct {
		name string
		path string
	}{
		{
			name: "Qwen3-1.7B (Current)",
			path: filepath.Join(modelsDir, "Qwen_Qwen3-1.7B-Q4_K_M.gguf"),
		},
		{
			name: "Qwen3-4B (Alternative)",
			path: filepath.Join(modelsDir, "Qwen_Qwen3-4B-Q4_K_M.gguf"),
		},
	}

	for _, model := range models {
		fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
		fmt.Printf("Testing: %s\n", model.name)
		fmt.Printf("Path: %s\n", model.path)
		fmt.Printf(strings.Repeat("=", 80) + "\n\n")

		// Create new service for each model
		libService, err := llama.NewLibraryService()
		if err != nil {
			log.Printf("❌ Failed to create library service: %v\n", err)
			continue
		}

		// Wait for initialization
		time.Sleep(3 * time.Second)

		if err := libService.InitializeLibrary(); err != nil {
			log.Printf("❌ Failed to initialize library: %v\n", err)
			continue
		}

		if err := libService.LoadChatModel(model.path); err != nil {
			log.Printf("❌ Failed to load model: %v\n", err)
			continue
		}

		llamaModel := llama.NewLlamaEinoModel(libService)
		ctx := context.Background()

		// Simulate conversation - start fresh each time to avoid buffer issues
		issueCount := 0

		for i, msg := range testMessages {
			fmt.Printf("\n--- Turn %d ---\n", i+1)
			fmt.Printf("User: %s\n\n", msg)

			// Fresh history for each turn to avoid buffer overflow
			history := []*schema.Message{
				{
					Role:    schema.User,
					Content: msg,
				},
			}

			response, err := llamaModel.Generate(ctx, history)
			if err != nil {
				fmt.Printf("❌ Error: %v\n", err)
				continue
			}

			// Truncate long responses for display
			displayContent := response.Content
			if len(displayContent) > 500 {
				displayContent = displayContent[:500] + "... [truncated]"
			}
			fmt.Printf("Assistant: %s\n", displayContent)

			// Analyze response for think tag issues
			hasThinkOpen := strings.Contains(response.Content, "<think>")
			hasThinkClose := strings.Contains(response.Content, "</think>")

			if hasThinkOpen {
				if !hasThinkClose {
					fmt.Printf("\n⚠️  ISSUE: Unclosed <think> tag!\n")
					issueCount++
				} else {
					// Check if there's actual content after </think>
					thinkEnd := strings.Index(response.Content, "</think>")
					afterThink := strings.TrimSpace(response.Content[thinkEnd+8:])

					if afterThink == "" {
						fmt.Printf("\n⚠️  ISSUE: No content after </think> tag!\n")
						issueCount++
					} else {
						fmt.Printf("\n✅ Think tags properly formatted\n")
						fmt.Printf("   Response length after think: %d chars\n", len(afterThink))
					}
				}
			} else {
				fmt.Printf("\n✅ No think tags (clean response)\n")
			}
		}

		fmt.Printf("\n" + strings.Repeat("-", 80) + "\n")
		fmt.Printf("Summary for %s:\n", model.name)
		fmt.Printf("Total issues found: %d out of %d responses\n", issueCount, len(testMessages))
		if issueCount == 0 {
			fmt.Printf("✅ All responses properly formatted!\n")
		} else {
			fmt.Printf("⚠️  Model has formatting issues\n")
		}
		fmt.Printf(strings.Repeat("-", 80) + "\n")

		// Cleanup
		libService.Cleanup()

		// Wait before next model
		time.Sleep(2 * time.Second)
	}

	fmt.Printf("\n" + strings.Repeat("=", 80) + "\n")
	fmt.Println("Comparison Complete!")
	fmt.Printf(strings.Repeat("=", 80) + "\n")
}
