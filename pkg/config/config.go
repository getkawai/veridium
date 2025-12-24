package config

import (
	"os"
	"strconv"
)

// Default constants (can be overridden by environment variables)
const (
	// KAWAI_RATE_PER_MILLION is the KAWAI reward per 1 Million Tokens (Phase 1: Mining Era)
	// Default: 100 KAWAI per 1M Tokens
	DefaultKawaiRatePerMillion = 100

	// COST_RATE_PER_MILLION is the USDT cost per 1 Million Tokens (Phase 2: Post-Mining Era)
	// Default: $1 USDT per 1M Tokens
	DefaultCostRatePerMillion = 1.0
)

// GetKawaiRatePerMillion returns the KAWAI mining rate per 1M tokens.
// Override with env: KAWAI_RATE_PER_MILLION
func GetKawaiRatePerMillion() int64 {
	if val := os.Getenv("KAWAI_RATE_PER_MILLION"); val != "" {
		if rate, err := strconv.ParseInt(val, 10, 64); err == nil {
			return rate
		}
	}
	return DefaultKawaiRatePerMillion
}

// GetCostRatePerMillion returns the USDT cost rate per 1M tokens for Phase 2.
// Override with env: COST_RATE_PER_MILLION
func GetCostRatePerMillion() float64 {
	if val := os.Getenv("COST_RATE_PER_MILLION"); val != "" {
		if rate, err := strconv.ParseFloat(val, 64); err == nil {
			return rate
		}
	}
	return DefaultCostRatePerMillion
}

// RewardMode represents the current economic phase
type RewardMode string

const (
	// ModeMining = Phase 1: Workers earn KAWAI tokens
	ModeMining RewardMode = "mining"
	// ModeUSDT = Phase 2: Workers earn USDT (post max supply)
	ModeUSDT RewardMode = "usdt"
)
