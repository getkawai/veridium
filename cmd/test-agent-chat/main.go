package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/kawai-network/veridium/internal/database"
	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/internal/services"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Initializing Agent Chat Service Test...")

	// 1. Initialize Database
	fmt.Println("Initializing Database...")
	dbService, err := database.NewService()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer dbService.Close()

	// 2. Initialize Llama Library
	fmt.Println("Initializing Llama Library...")
	libService, err := llama.NewLibraryService()
	if err != nil {
		log.Fatalf("Failed to create library service: %v", err)
	}
	defer libService.Cleanup()

	if err := libService.InitializeLibrary(); err != nil {
		log.Fatalf("Failed to initialize library: %v", err)
	}

	// Load a model
	homeDir, _ := os.UserHomeDir()
	modelsDir := filepath.Join(homeDir, ".llama-cpp", "models")
	modelPath := ""
	entries, _ := os.ReadDir(modelsDir)
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".gguf") {
			modelPath = filepath.Join(modelsDir, entry.Name())
			if strings.Contains(strings.ToLower(entry.Name()), "qwen") || strings.Contains(strings.ToLower(entry.Name()), "llama") {
				break
			}
		}
	}

	if modelPath != "" {
		fmt.Printf("Loading model: %s\n", filepath.Base(modelPath))
		if err := libService.LoadChatModel(modelPath); err != nil {
			log.Printf("Warning: Failed to load model: %v", err)
		}
	} else {
		fmt.Println("Warning: No model found in .llama-cpp/models")
	}

	// 3. Initialize AgentChatService
	fmt.Println("Creating AgentChatService...")
	agentService := services.NewAgentChatService(
		nil, // app
		dbService,
		libService,
		nil, // kbService
		nil, // contextBridge
		nil, // threadService
	)

	// 5. Start Chat Loop
	reader := bufio.NewReader(os.Stdin)
	sessionID := uuid.New().String()
	userID := "DEFAULT_LOBE_CHAT_USER"

	fmt.Printf("\n=== Agent Chat Test (Session: %s) ===\n", sessionID)
	fmt.Println("Type 'exit' or 'quit' to stop.")

	ctx := context.Background()

	for {
		fmt.Print("\nYou: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "exit" || input == "quit" {
			break
		}

		if input == "" {
			continue
		}

		req := services.ChatRequest{
			SessionID: sessionID,
			UserID:    userID,
			Message:   input,
			Stream:    false,
			Tools:     []string{"web-search"}, // Enable search tool (calculator removed)
		}

		fmt.Print("Agent: ")
		startTime := time.Now()
		resp, err := agentService.Chat(ctx, req)
		duration := time.Since(startTime)

		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("%s\n", resp.Message)
		fmt.Printf("(Took %v)\n", duration)

		if len(resp.ToolCalls) > 0 {
			fmt.Printf("[Tool Calls: %d]\n", len(resp.ToolCalls))
			for _, tc := range resp.ToolCalls {
				fmt.Printf("  - %s(%s)\n", tc.Function.Name, tc.Function.Arguments)
			}
		}

		if resp.TopicID != "" {
			fmt.Printf("[Topic: %s]\n", resp.TopicID)
		}
	}
}
