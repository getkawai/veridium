package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/kawai-network/x/constant"
)

type DiscordAlert struct {
	WebhookURL string
	Client     *http.Client
}

// DiscordWebhookPayload represents the Discord webhook message structure
type DiscordWebhookPayload struct {
	Content string         `json:"content,omitempty"`
	Embeds  []DiscordEmbed `json:"embeds,omitempty"`
}

// DiscordEmbed represents a Discord embed object
type DiscordEmbed struct {
	Title       string              `json:"title,omitempty"`
	Description string              `json:"description,omitempty"`
	Color       int                 `json:"color,omitempty"`
	Fields      []DiscordEmbedField `json:"fields,omitempty"`
	Footer      *DiscordEmbedFooter `json:"footer,omitempty"`
	Timestamp   string              `json:"timestamp,omitempty"`
}

// DiscordEmbedField represents a field in a Discord embed
type DiscordEmbedField struct {
	Name   string `json:"name"`
	Value  string `json:"value"`
	Inline bool   `json:"inline,omitempty"`
}

// DiscordEmbedFooter represents the footer of a Discord embed
type DiscordEmbedFooter struct {
	Text string `json:"text"`
}

// NewDiscordAlert creates a new Discord alert service
func NewDiscordAlert() *DiscordAlert {
	return &DiscordAlert{
		WebhookURL: constant.GetDiscordWebhook(),
		Client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// SendMessage sends a text message to Discord webhook
func (d *DiscordAlert) SendMessage(text string) error {
	if d.WebhookURL == "" {
		return fmt.Errorf("discord webhook URL not configured")
	}

	payload := DiscordWebhookPayload{
		Content: text,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", d.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord webhook error: status %d", resp.StatusCode)
	}

	return nil
}

// SendEmbed sends an embed message to Discord webhook
func (d *DiscordAlert) SendEmbed(embed DiscordEmbed) error {
	if d.WebhookURL == "" {
		return fmt.Errorf("discord webhook URL not configured")
	}

	payload := DiscordWebhookPayload{
		Embeds: []DiscordEmbed{embed},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequest("POST", d.WebhookURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := d.Client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("discord webhook error: status %d", resp.StatusCode)
	}

	return nil
}

// SendAlert formats and sends a standard alert message to Discord
func (d *DiscordAlert) SendAlert(level, source, message string) {
	if d.WebhookURL == "" {
		return // Silent return if not configured
	}

	// Color codes for different levels
	color := 0x3498db // Blue (INFO)
	switch level {
	case "ERROR":
		color = 0xe74c3c // Red
	case "WARNING":
		color = 0xf39c12 // Orange
	case "SUCCESS":
		color = 0x2ecc71 // Green
	}

	// Determine hostname for context
	hostname, _ := os.Hostname()
	footer := &DiscordEmbedFooter{
		Text: fmt.Sprintf("Host: %s", hostname),
	}

	embed := DiscordEmbed{
		Title:       fmt.Sprintf("[%s] %s", level, source),
		Description: message,
		Color:       color,
		Footer:      footer,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	// Send asynchronously to avoid blocking
	go func() {
		if err := d.SendEmbed(embed); err != nil {
			slog.Error("Failed to send discord alert", "error", err)
		}
	}()
}
