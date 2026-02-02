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
	"github.com/kawai-network/veridium/pkg/types"
)

type TelegramAlert struct {
	BotToken        string
	ChatID          string
	Client          *http.Client
	DiscordFallback *DiscordAlert // Fallback to Discord if Telegram fails
}

// NewTelegramAlert creates a new alert service with Discord fallback
func NewTelegramAlert() *TelegramAlert {
	return &TelegramAlert{
		BotToken:        constant.GetTelegramBotToken(),
		ChatID:          constant.GetTelegramChatId(),
		Client:          &http.Client{Timeout: 10 * time.Second},
		DiscordFallback: NewDiscordAlert(), // Initialize Discord fallback
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

			// Fallback to Discord if Telegram fails
			if t.DiscordFallback != nil {
				t.DiscordFallback.SendAlert(level, source, message)
			}
		}
	}()
}

// SendJobRewardLog sends a detailed job reward record to Telegram for audit trail
// This provides an immutable backup of all job rewards for verification during settlement
// Uses types.JobRewardRecord (same struct used by KV store) for easy verification
// Format: Machine-readable JSON for easy export and verification
func (t *TelegramAlert) SendJobRewardLog(record *types.JobRewardRecord) {
	if t.BotToken == "" || t.ChatID == "" {
		return // Silent return if not configured
	}

	// Format 1: Human-readable header
	text := "💰 *Job Reward*\n"
	text += fmt.Sprintf("`%s` | %d tokens | %s\n\n",
		record.Timestamp.Format("2006-01-02 15:04:05"),
		record.TokenUsage,
		record.RewardType)

	// Format 2: Machine-readable JSON (for easy export/verification)
	jsonData, err := json.Marshal(record)
	if err != nil {
		slog.Error("Failed to marshal job reward for telegram", "error", err)
		return
	}
	text += "```json\n" + string(jsonData) + "\n```\n\n"

	// Format 3: Quick summary
	splitType := "90/5/5"
	if record.HasReferrer {
		splitType = "85/5/5/5"
	}
	text += fmt.Sprintf("📊 Split: %s | C:%s U:%s D:%s",
		splitType,
		shortenAddress(record.ContributorAddress),
		shortenAddress(record.UserAddress),
		shortenAddress(record.DeveloperAddress))

	if record.HasReferrer {
		text += fmt.Sprintf(" A:%s", shortenAddress(record.ReferrerAddress))
	}

	// Send asynchronously to avoid blocking RecordJobReward
	go func() {
		if err := t.SendMessage(text); err != nil {
			slog.Error("Failed to send job reward log to Telegram", "error", err)

			// Fallback to Discord if Telegram fails
			if t.DiscordFallback != nil {
				discordMsg := fmt.Sprintf("💰 **Job Reward** (Telegram Fallback)\n```json\n%s\n```", string(jsonData))
				if discordErr := t.DiscordFallback.SendMessage(discordMsg); discordErr != nil {
					slog.Error("Failed to send job reward log to Discord fallback", "error", discordErr)
				} else {
					slog.Info("Job reward log sent to Discord fallback successfully")
				}
			}
		}
	}()
}

// SendCashbackLog sends a detailed cashback record to Telegram for audit trail
// This provides an immutable backup of all cashback records for verification during settlement
// Uses types.CashbackRecord (same struct used by KV store) for easy verification
// Format: Machine-readable JSON for easy export and verification
func (t *TelegramAlert) SendCashbackLog(record *types.CashbackRecord) {
	if t.BotToken == "" || t.ChatID == "" {
		return // Silent return if not configured
	}

	// Format 1: Human-readable header
	text := "💎 *Cashback Tracked*\n"
	text += fmt.Sprintf("`%s` | Period %d | Tier %d\n\n",
		record.Timestamp.Format("2006-01-02 15:04:05"),
		record.Period,
		record.Tier)

	// Format 2: Machine-readable JSON (for easy export/verification)
	jsonData, err := json.Marshal(record)
	if err != nil {
		slog.Error("Failed to marshal cashback record for telegram", "error", err)
		return
	}
	text += "```json\n" + string(jsonData) + "\n```\n\n"

	// Format 3: Quick summary
	firstTimeBonus := ""
	if record.IsFirstTime {
		firstTimeBonus = " 🎁 First-time"
	}
	text += fmt.Sprintf("📊 Rate: %.2f%% | User: %s%s",
		float64(record.RateBPS)/100.0,
		shortenAddress(record.UserAddress),
		firstTimeBonus)

	// Send asynchronously to avoid blocking TrackCashback
	go func() {
		if err := t.SendMessage(text); err != nil {
			slog.Error("Failed to send cashback log to Telegram", "error", err)

			// Fallback to Discord if Telegram fails
			if t.DiscordFallback != nil {
				discordMsg := fmt.Sprintf("💎 **Cashback Tracked** (Telegram Fallback)\n```json\n%s\n```", string(jsonData))
				if discordErr := t.DiscordFallback.SendMessage(discordMsg); discordErr != nil {
					slog.Error("Failed to send cashback log to Discord fallback", "error", discordErr)
				} else {
					slog.Info("Cashback log sent to Discord fallback successfully")
				}
			}
		}
	}()
}

// SendReferralTrialLog sends a detailed referral trial claim record to Telegram for audit trail
// This provides an immutable backup of all trial claims for verification and fraud detection
// Uses types.ReferralTrialRecord (same struct used by KV store) for easy verification
// Format: Machine-readable JSON for easy export and verification
func (t *TelegramAlert) SendReferralTrialLog(record *types.ReferralTrialRecord) {
	if t.BotToken == "" || t.ChatID == "" {
		return // Silent return if not configured
	}

	// Format 1: Human-readable header
	claimType := "Solo Trial"
	if record.IsReferral {
		claimType = "Referral Trial"
	}
	text := fmt.Sprintf("🎁 *%s Claimed*\n", claimType)
	text += fmt.Sprintf("`%s`\n\n", record.Timestamp.Format("2006-01-02 15:04:05"))

	// Format 2: Machine-readable JSON (for easy export/verification)
	jsonData, err := json.Marshal(record)
	if err != nil {
		slog.Error("Failed to marshal referral trial record for telegram", "error", err)
		return
	}
	text += "```json\n" + string(jsonData) + "\n```\n\n"

	// Format 3: Quick summary
	text += fmt.Sprintf("👤 User: %s\n", shortenAddress(record.UserAddress))
	if record.IsReferral {
		text += fmt.Sprintf("🤝 Referrer: %s (Code: %s)\n",
			shortenAddress(record.ReferrerAddress),
			record.ReferralCode)
	}
	// Safely truncate machine ID
	machineIDShort := record.MachineID
	if len(machineIDShort) > 16 {
		machineIDShort = machineIDShort[:16] + "..."
	}
	text += fmt.Sprintf("🔒 Machine: %s", machineIDShort)

	// Send asynchronously to avoid blocking ClaimFreeTrial
	go func() {
		if err := t.SendMessage(text); err != nil {
			slog.Error("Failed to send referral trial log to Telegram", "error", err)

			// Fallback to Discord if Telegram fails
			if t.DiscordFallback != nil {
				discordMsg := fmt.Sprintf("🎁 **%s Claimed** (Telegram Fallback)\n```json\n%s\n```", claimType, string(jsonData))
				if discordErr := t.DiscordFallback.SendMessage(discordMsg); discordErr != nil {
					slog.Error("Failed to send referral trial log to Discord fallback", "error", discordErr)
				} else {
					slog.Info("Referral trial log sent to Discord fallback successfully")
				}
			}
		}
	}()
}

// shortenAddress shortens Ethereum address for display (0x1234...5678)
func shortenAddress(addr string) string {
	if len(addr) < 10 {
		return addr
	}
	return addr[:6] + "..." + addr[len(addr)-4:]
}

// SendBalanceDeductionFailure sends an alert when balance deduction fails after AI service
// This is a critical error that requires manual intervention to reconcile the debt
func (t *TelegramAlert) SendBalanceDeductionFailure(userAddress string, amount string, err error) {
	if t.BotToken == "" || t.ChatID == "" {
		return // Silent return if not configured
	}

	text := "🚨 *BALANCE DEDUCTION FAILED*\n"
	text += fmt.Sprintf("`%s`\n\n", time.Now().Format("2006-01-02 15:04:05"))
	text += fmt.Sprintf("👤 User: %s\n", shortenAddress(userAddress))
	text += fmt.Sprintf("💰 Amount: %s micro USDT\n", amount)
	text += fmt.Sprintf("❌ Error: %v\n\n", err)
	text += "⚠️ *ACTION REQUIRED*: Manual debt reconciliation needed"

	// Send asynchronously to avoid blocking
	go func() {
		if sendErr := t.SendMessage(text); sendErr != nil {
			slog.Error("Failed to send balance deduction failure alert to Telegram", "error", sendErr)

			// Fallback to Discord if Telegram fails
			if t.DiscordFallback != nil {
				discordMsg := fmt.Sprintf("🚨 **BALANCE DEDUCTION FAILED** (Telegram Fallback)\nUser: %s\nAmount: %s micro USDT\nError: %v",
					userAddress, amount, err)
				if discordErr := t.DiscordFallback.SendMessage(discordMsg); discordErr != nil {
					slog.Error("Failed to send balance deduction failure alert to Discord fallback", "error", discordErr)
				}
			}
		}
	}()
}
