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

	"github.com/Knetic/govaluate"
	"github.com/cloudwego/eino/schema"
	"github.com/google/uuid"
	"github.com/kawai-network/veridium/internal/database"
	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/pkg/toolsengine"
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

	// 3. Initialize Tools Engine & Calculator Tool
	fmt.Println("Initializing Tools Engine...")
	toolsEngine, err := toolsengine.NewToolsEngine(toolsengine.Config{})
	if err != nil {
		log.Fatalf("Failed to create tools engine: %v", err)
	}

	// Create Calculator Tool
	calcExecutor := func(ctx context.Context, args map[string]interface{}) (interface{}, error) {
		expression, ok := args["expression"].(string)
		if !ok {
			return nil, fmt.Errorf("expression argument is required")
		}

		fmt.Printf("🧮 Calculator Tool Invoked: %s\n", expression)

		expr, err := govaluate.NewEvaluableExpression(expression)
		if err != nil {
			return nil, fmt.Errorf("invalid expression: %w", err)
		}

		result, err := expr.Evaluate(nil)
		if err != nil {
			return nil, fmt.Errorf("evaluation failed: %w", err)
		}

		return fmt.Sprintf("%v", result), nil
	}

	calcTool, err := toolsengine.NewToolBuilder("calculator", "Calculator").
		WithDescription("Useful for performing mathematical calculations. Input should be a mathematical expression string.").
		WithParameter("expression", schema.String, "The mathematical expression to evaluate (e.g. '2 + 2', '3 * 4')", true).
		WithExecutor(calcExecutor).
		WithCategory("utility").
		Build()

	if err != nil {
		log.Fatalf("Failed to build calculator tool: %v", err)
	}

	if err := toolsEngine.RegisterTool(calcTool); err != nil {
		log.Fatalf("Failed to register calculator tool: %v", err)
	}

	toolsBridge := services.NewToolsEngineBridge(toolsEngine)

	// 4. Initialize AgentChatService
	fmt.Println("Creating AgentChatService...")
	agentService := services.NewAgentChatService(
		nil, // app
		dbService,
		libService,
		nil,         // kbService
		toolsBridge, // toolsBridge
		nil,         // contextBridge
		nil,         // threadService
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
			Tools:     []string{"calculator"}, // Enable calculator tool
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
