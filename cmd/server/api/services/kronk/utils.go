package kronk

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
)

// validateKeystoreJSON validates keystore JSON format
func validateKeystoreJSON(jsonStr string) (bool, string) {
	// Check if JSON is empty
	if strings.TrimSpace(jsonStr) == "" {
		return false, "Keystore JSON cannot be empty"
	}

	// Parse JSON properly
	var ks struct {
		Address string `json:"address"`
		Crypto  struct {
			Kdf        string `json:"kdf"`
			Ciphertext string `json:"ciphertext"`
		} `json:"crypto"`
		Version int `json:"version"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &ks); err != nil {
		return false, fmt.Sprintf("Invalid JSON format: %v", err)
	}

	// Validate required fields
	if ks.Address == "" {
		return false, "Missing required field: address"
	}
	if ks.Version == 0 {
		return false, "Missing or invalid required field: version"
	}
	if ks.Crypto.Kdf == "" && ks.Crypto.Ciphertext == "" {
		return false, "Invalid keystore: missing crypto information (kdf or ciphertext)"
	}

	return true, ""
}

// calculatePasswordStrength calculates password strength score
func calculatePasswordStrength(password string) (int, string, string) {
	score := 0
	if len(password) >= 8 {
		score += 25
	}
	if len(password) >= 12 {
		score += 15
	}
	if len(password) >= 16 {
		score += 10
	}
	if len(password) > 0 {
		if strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
			score += 15
		}
		if strings.ContainsAny(password, "0123456789") {
			score += 15
		}
		if strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:,.<>?") {
			score += 20
		}
	}

	score = min(score, 100)

	var label, color string
	switch {
	case score < 40:
		label = "Weak - Add more characters"
		color = "#ff4d4f"
	case score < 60:
		label = "Fair - Add uppercase or numbers"
		color = "#faad14"
	case score < 80:
		label = "Good - Add special characters"
		color = "#52c41a"
	default:
		label = "Strong password"
		color = "#1890ff"
	}

	return score, label, color
}

// validatePrivateKey validates a hex private key
func validatePrivateKey(privateKey string) (bool, string) {
	// Clean private key
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privateKey = strings.TrimPrefix(privateKey, "0X")

	if len(privateKey) != 64 {
		return false, fmt.Sprintf("invalid length: %d/64 characters", len(privateKey))
	}

	if _, err := hex.DecodeString(privateKey); err != nil {
		return false, "invalid hex characters"
	}

	return true, ""
}
