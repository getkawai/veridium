package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/kawai-network/veridium/internal/constant"
)

type TelegramAlert struct {
	BotToken string
	ChatID   string
	Client   *http.Client
}

// NewTelegramAlert creates a new alert service
func NewTelegramAlert() *TelegramAlert {
	return &TelegramAlert{
		BotToken: constant.GetTelegramBotToken(),
		ChatID:   constant.GetTelegramChatId(),
		Client:   &http.Client{Timeout: 10 * time.Second},
	}
}

// SendMessage sends a text message to the configured Telegram chat
func (t *TelegramAlert) SendMessage(text string) error {
	if t.BotToken == "" || t.ChatID == "" {
		return fmt.Errorf("telegram credentials not configured")
	}

	apiURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.BotToken)

	payload := map[string]string{
		"chat_id":    t.ChatID,
		"text":       text,
		"parse_mode": "Markdown", // Optional: allows bold/italic
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram api error: status %d", resp.StatusCode)
	}

	return nil
}

// SendAlert formats and sends a standard alert message
func (t *TelegramAlert) SendAlert(level, source, message string) {
	if t.BotToken == "" || t.ChatID == "" {
		return // Silent return if not configured
	}

	// Format: 🚨 [ERROR] ServiceName: Something went wrong
	icon := "ℹ️"
	switch level {
	case "ERROR":
		icon = "🚨"
	case "WARNING":
		icon = "⚠️"
	case "SUCCESS":
		icon = "✅"
	}

	text := fmt.Sprintf("%s *[%s]* %s:\n%s", icon, level, source, message)

	// Determine hostname for context
	hostname, _ := os.Hostname()
	if hostname != "" {
		text += fmt.Sprintf("\n\n_Host: %s_", hostname)
	}

	// Send concurrently to avoid blocking the main thread
	go func() {
		if err := t.SendMessage(text); err != nil {
			slog.Error("Failed to send telegram alert", "error", err)
		}
	}()
}
