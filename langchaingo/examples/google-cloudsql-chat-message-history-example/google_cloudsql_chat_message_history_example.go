package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kawai-network/veridium/langchaingo/llms"
	"github.com/kawai-network/veridium/langchaingo/memory/sqlite3"
	_ "modernc.org/sqlite"
)

func printMessages(ctx context.Context, cmh *sqlite3.SqliteChatMessageHistory) {
	msgs, err := cmh.Messages(ctx)
	if err != nil {
		log.Fatal(err)
	}
	for _, msg := range msgs {
		fmt.Println("Message:", msg)
	}
}

func main() {
	// Create a temporary directory for the SQLite database
	tempDir := filepath.Join(os.TempDir(), "sqlite-chat-history-example")
	os.MkdirAll(tempDir, 0755)
	defer os.RemoveAll(tempDir) // Clean up after example

	dbPath := filepath.Join(tempDir, "chat_history.db")

	// Open SQLite database
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	ctx := context.Background()
	tableName := "chat_messages"
	sessionID := "example-session-123"

	// Creates a new Chat Message History with SQLite
	cmh := sqlite3.NewSqliteChatMessageHistory(
		sqlite3.WithDB(db),
		sqlite3.WithContext(ctx),
		sqlite3.WithTableName(tableName),
		sqlite3.WithSession(sessionID),
	)

	// Creates individual messages and adds them to the chat message history.
	aiMessage := llms.AIChatMessage{Content: "test AI message"}
	humanMessage := llms.HumanChatMessage{Content: "test HUMAN message"}
	// Adds a user message to the chat message history.
	err = cmh.AddUserMessage(ctx, aiMessage.GetContent())
	if err != nil {
		log.Fatal(err)
	}
	// Adds a user message to the chat message history.
	err = cmh.AddUserMessage(ctx, humanMessage.GetContent())
	if err != nil {
		log.Fatal(err)
	}

	printMessages(ctx, cmh)

	// Create multiple messages and store them in the chat message history at the same time.
	multipleMessages := []llms.ChatMessage{
		llms.AIChatMessage{Content: "first AI test message from AddMessages"},
		llms.AIChatMessage{Content: "second AI test message from AddMessages"},
		llms.HumanChatMessage{Content: "first HUMAN test message from AddMessages"},
	}

	// Adds multiple messages to the chat message history.
	err = cmh.AddMessages(ctx, multipleMessages)
	if err != nil {
		log.Fatal(err)
	}

	printMessages(ctx, cmh)

	// Create messages that will overwrite the existing ones
	overWrittingMessages := []llms.ChatMessage{
		llms.AIChatMessage{Content: "overwritten AI test message"},
		llms.HumanChatMessage{Content: "overwritten HUMAN test message"},
	}
	// Overwrites the existing messages with new ones.
	err = cmh.SetMessages(ctx, overWrittingMessages)
	if err != nil {
		log.Fatal(err)
	}

	printMessages(ctx, cmh)

	// Clear all the messages from the current session.
	err = cmh.Clear(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
