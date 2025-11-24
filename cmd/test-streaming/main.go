package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kawai-network/veridium/internal/database"
	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/wailsapp/wails/v3/pkg/application"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Println("🧪 Testing Agent Chat Streaming...")

	// Initialize database (uses ./data by default)
	dbService, err := database.NewService()
	if err != nil {
		log.Fatalf("Failed to init database: %v", err)
	}
	defer dbService.Close()

	// Initialize llama library service
	libService, err := llama.NewLibraryService()
	if err != nil {
		log.Fatalf("Failed to init llama library: %v", err)
	}
	defer libService.Cleanup()

	// Initialize library (auto-downloads if needed)
	log.Println("📚 Initializing llama.cpp library...")
	if err := libService.InitializeLibrary(); err != nil {
		log.Fatalf("Failed to initialize library: %v", err)
	}

	// Auto-load best model
	log.Println("📦 Auto-loading best available model...")
	models, err := libService.GetAvailableModels()
	if err != nil || len(models) == 0 {
		log.Fatal("No models found")
	}
	log.Printf("📦 Loading model: %s", models[0])
	if err := libService.LoadChatModel(models[0]); err != nil {
		log.Fatalf("Failed to load model: %v", err)
	}

	// Create mock Wails app for events
	app := application.New(application.Options{
		Name: "test-streaming",
	})

	// Create agent chat service
	agentService := services.NewAgentChatService(app, dbService, libService, nil, nil, nil, nil)

	// Test streaming
	ctx := context.Background()
	sessionID := "test-stream-session"
	userID := "DEFAULT_LOBE_CHAT_USER"

	// Setup event listener
	eventChannel := fmt.Sprintf("chat:stream:%s", sessionID)
	log.Printf("📡 Listening to events: %s", eventChannel)

	app.Event.On(eventChannel, func(event *application.CustomEvent) {
		data := event.Data.(map[string]interface{})
		eventType := data["type"].(string)

		switch eventType {
		case "start":
			log.Printf("✅ [START] message_id: %s", data["message_id"])
		case "chunk":
			content := data["content"].(string)
			fullContent := data["full_content"].(string)
			log.Printf("📦 [CHUNK] +%d chars, total: %d chars", len(content), len(fullContent))
			// Print first 50 chars of full content
			preview := fullContent
			if len(preview) > 50 {
				preview = preview[:50] + "..."
			}
			log.Printf("   Content: %s", preview)
		case "complete":
			log.Printf("✅ [COMPLETE] message_id: %s", data["message_id"])
			log.Printf("   Topic: %s", data["topic_id"])
		}
	})

	// Send message with streaming
	log.Println("\n🚀 Sending message with streaming enabled...")
	startTime := time.Now()

	response, err := agentService.Chat(ctx, services.ChatRequest{
		SessionID: sessionID,
		UserID:    userID,
		Message:   "What is artificial intelligence? Explain in 2-3 sentences.",
		Stream:    true, // Enable streaming
	})

	duration := time.Since(startTime)

	if err != nil {
		log.Fatalf("❌ Chat failed: %v", err)
	}

	// Print final response
	log.Println("\n📝 Final Response:")
	log.Printf("   Message ID: %s", response.MessageID)
	log.Printf("   Session ID: %s", response.SessionID)
	log.Printf("   Topic ID: %s", response.TopicID)
	log.Printf("   Content length: %d chars", len(response.Message))
	log.Printf("   Duration: %v", duration)
	log.Printf("\n   Content:\n%s\n", response.Message)

	log.Println("✅ Test completed!")
}

