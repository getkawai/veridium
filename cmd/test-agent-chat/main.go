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

	// Load a model (optional but good for real testing)
	// We'll try to find a model or let the service handle it?
	// AgentChatService uses LlamaEinoModel which needs a loaded model in libService?
	// Actually LlamaEinoModel uses libService.
	// We should probably load a model if we want it to work.
	homeDir, _ := os.UserHomeDir()
	modelsDir := filepath.Join(homeDir, ".llama-cpp", "models")
	// Try to find a model
	modelPath := ""
	entries, _ := os.ReadDir(modelsDir)
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".gguf") {
			modelPath = filepath.Join(modelsDir, entry.Name())
			// Prefer Qwen or Llama
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
	// We pass nil for optional services for now
	fmt.Println("Creating AgentChatService...")
	agentService := services.NewAgentChatService(
		nil, // app
		dbService,
		libService,
		nil, // kbService
		nil, // toolsBridge
		nil, // contextBridge
		nil, // threadService
	)

	// 4. Start Chat Loop
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
			Stream:    true, // Enable streaming to test that path (though we won't see real-time stream in this simple CLI without callback)
			// Actually, if we pass Stream: true, the service might try to emit events to `app`.
			// Since `app` is nil, we should check if that causes panic.
			// The code checks `if req.Stream && s.app != nil`. So it should be safe.
			// But if app is nil, it falls back to non-streaming or just returns the full response?
			// Let's check the code:
			// if req.Stream && s.app != nil { ... } else { ... }
			// So if app is nil, it goes to else block (standard Eino agent).
		}

		// If we want to test streaming logic, we need to mock app or just test non-streaming.
		// Let's test non-streaming first to be safe, or let it fall back.
		req.Stream = false

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

		if resp.TopicID != "" {
			fmt.Printf("[Topic: %s]\n", resp.TopicID)
		}
	}
}
