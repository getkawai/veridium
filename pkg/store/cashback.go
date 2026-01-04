package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"
)

// CashbackRecord represents a single cashback entry
type CashbackRecord struct {
	UserAddress  string    `json:"user_address"`
	DepositTxHash string   `json:"deposit_tx_hash"`
	DepositAmount string   `json:"deposit_amount"` // USDT amount (6 decimals)
	CashbackAmount string  `json:"cashback_amount"` // KAWAI amount (18 decimals)
	Rate          uint64   `json:"rate"`            // Rate in basis points (e.g., 200 = 2%)
	Tier          uint64   `json:"tier"`            // Tier (1-5)
	IsFirstTime   bool     `json:"is_first_time"`   // First-time bonus
	CreatedAt     time.Time `json:"created_at"`
	Period        uint64   `json:"period"`          // Settlement period
	Claimed       bool     `json:"claimed"`
}

// CashbackStats represents user's cashback statistics
type CashbackStats struct {
	TotalCashback    string `json:"total_cashback"`     // Total KAWAI earned
	PendingCashback  string `json:"pending_cashback"`   // Pending (unclaimed)
	ClaimedCashback  string `json:"claimed_cashback"`   // Already claimed
	TotalDeposits    uint64 `json:"total_deposits"`     // Number of deposits
	FirstDepositAt   *time.Time `json:"first_deposit_at"`
	LastDepositAt    *time.Time `json:"last_deposit_at"`
}

// CalculateCashback calculates KAWAI cashback for a USDT deposit
// Returns: cashbackAmount (wei), rate (bps), tier, isFirstTime
func (s *KVStore) CalculateCashback(ctx context.Context, userAddress string, depositAmount *big.Int) (string, uint64, uint64, bool, error) {
	// Get user stats to check if first-time
	stats, err := s.GetCashbackStats(ctx, userAddress)
	if err != nil {
		log.Printf("⚠️  [Cashback] Failed to get stats for %s: %v", userAddress, err)
		// Assume not first-time on error
		stats = &CashbackStats{}
	}
	
	isFirstTime := stats.TotalDeposits == 0
	
	// Determine tier based on deposit amount (USDT has 6 decimals)
	// Tier 1: < 100 USDT
	// Tier 2: 100-500 USDT
	// Tier 3: 500-1000 USDT
	// Tier 4: 1000-5000 USDT
	// Tier 5: >= 5000 USDT
	
	usdtAmount := new(big.Int).Set(depositAmount) // Copy to avoid mutation
	tier := uint64(1)
	baseRate := uint64(100) // 1% in basis points
	
	// Convert to USDT (divide by 1e6)
	oneMillion := big.NewInt(1_000_000)
	usdtValue := new(big.Int).Div(usdtAmount, oneMillion)
	
	// Tier structure matches DepositCashbackDistributor.sol
	// Base rates: 1-2% (100-200 basis points)
	// Phase 1 bounds: 1.5%-2.5% (150-250 basis points)
	if usdtValue.Cmp(big.NewInt(5000)) >= 0 {
		tier = 5
		baseRate = 200 // 2% (capped by Phase 1 max: 2.5%)
	} else if usdtValue.Cmp(big.NewInt(1000)) >= 0 {
		tier = 4
		baseRate = 175 // 1.75%
	} else if usdtValue.Cmp(big.NewInt(500)) >= 0 {
		tier = 3
		baseRate = 150 // 1.5%
	} else if usdtValue.Cmp(big.NewInt(100)) >= 0 {
		tier = 2
		baseRate = 125 // 1.25%
	}
	
	// First-time bonus: 2.5% (Phase 1 max)
	rate := baseRate
	if isFirstTime {
		rate = 250 // 2.5% for first deposit (Phase 1 max bound)
	}
	
	// Calculate cashback: (depositAmount * rate * 1e18) / (10000 * 1e6)
	// This converts USDT (6 decimals) to KAWAI (18 decimals) with rate applied
	cashback := new(big.Int).Mul(depositAmount, big.NewInt(int64(rate)))
	cashback.Mul(cashback, big.NewInt(1e18))
	cashback.Div(cashback, big.NewInt(10000))
	cashback.Div(cashback, big.NewInt(1e6))
	
	// Apply tier caps (max KAWAI per deposit)
	// Tier 1: 5K KAWAI
	// Tier 2: 10K KAWAI
	// Tier 3: 15K KAWAI
	// Tier 4: 20K KAWAI
	// Tier 5: 20K KAWAI
	maxCashback := big.NewInt(5000) // Default 5K
	switch tier {
	case 2:
		maxCashback = big.NewInt(10000)
	case 3:
		maxCashback = big.NewInt(15000)
	case 4, 5:
		maxCashback = big.NewInt(20000)
	}
	maxCashback.Mul(maxCashback, big.NewInt(1e18)) // Convert to wei
	
	if cashback.Cmp(maxCashback) > 0 {
		cashback = maxCashback
	}
	
	return cashback.String(), rate, tier, isFirstTime, nil
}

// TrackCashback records a cashback entry for a deposit
func (s *KVStore) TrackCashback(ctx context.Context, userAddress, txHash string, depositAmount *big.Int, period uint64) error {
	// Calculate cashback
	cashbackAmount, rate, tier, isFirstTime, err := s.CalculateCashback(ctx, userAddress, depositAmount)
	if err != nil {
		return fmt.Errorf("failed to calculate cashback: %w", err)
	}
	
	// Create record
	record := CashbackRecord{
		UserAddress:    userAddress,
		DepositTxHash:  txHash,
		DepositAmount:  depositAmount.String(),
		CashbackAmount: cashbackAmount,
		Rate:           rate,
		Tier:           tier,
		IsFirstTime:    isFirstTime,
		CreatedAt:      time.Now(),
		Period:         period,
		Claimed:        false,
	}
	
	// Store record
	key := fmt.Sprintf("cashback:%s:%s", userAddress, txHash)
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("failed to marshal cashback record: %w", err)
	}
	
	if err := s.StoreMarketplaceData(ctx, key, data); err != nil {
		return fmt.Errorf("failed to store cashback record: %w", err)
	}
	
	// Track user for this period (for settlement)
	if err := s.trackUserForPeriod(ctx, period, userAddress); err != nil {
		log.Printf("⚠️  [Cashback] Failed to track user for period: %v", err)
		// Don't fail - record is already stored
	}
	
	// Update stats
	if err := s.updateCashbackStats(ctx, userAddress, cashbackAmount); err != nil {
		log.Printf("⚠️  [Cashback] Failed to update stats: %v", err)
		// Don't fail - record is already stored
	}
	
	log.Printf("✅ [Cashback] Tracked: user=%s, deposit=%s USDT, cashback=%s KAWAI, rate=%d bps, tier=%d, first=%v",
		userAddress, depositAmount.String(), cashbackAmount, rate, tier, isFirstTime)
	
	return nil
}

// GetCashbackStats retrieves cashback statistics for a user
func (s *KVStore) GetCashbackStats(ctx context.Context, userAddress string) (*CashbackStats, error) {
	key := fmt.Sprintf("cashback_stats:%s", userAddress)
	data, err := s.GetMarketplaceData(ctx, key)
	if err != nil {
		return &CashbackStats{
			TotalCashback:   "0",
			PendingCashback: "0",
			ClaimedCashback: "0",
			TotalDeposits:   0,
		}, nil
	}
	
	var stats CashbackStats
	if err := json.Unmarshal(data, &stats); err != nil {
		return nil, fmt.Errorf("failed to unmarshal cashback stats: %w", err)
	}
	
	return &stats, nil
}

// updateCashbackStats atomically updates user's cashback statistics
func (s *KVStore) updateCashbackStats(ctx context.Context, userAddress, cashbackAmount string) error {
	key := fmt.Sprintf("cashback_stats:%s", userAddress)
	
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		// Read current stats
		stats, err := s.GetCashbackStats(ctx, userAddress)
		if err != nil {
			return err
		}
		
	// Update stats
	totalCashback := new(big.Int)
	if stats.TotalCashback != "" && stats.TotalCashback != "0" {
		if _, ok := totalCashback.SetString(stats.TotalCashback, 10); !ok {
			return fmt.Errorf("invalid TotalCashback value: %s", stats.TotalCashback)
		}
	}
	
	pendingCashback := new(big.Int)
	if stats.PendingCashback != "" && stats.PendingCashback != "0" {
		if _, ok := pendingCashback.SetString(stats.PendingCashback, 10); !ok {
			return fmt.Errorf("invalid PendingCashback value: %s", stats.PendingCashback)
		}
	}
	
	newCashback := new(big.Int)
	if _, ok := newCashback.SetString(cashbackAmount, 10); !ok {
		return fmt.Errorf("invalid cashbackAmount value: %s", cashbackAmount)
	}
		
		totalCashback.Add(totalCashback, newCashback)
		pendingCashback.Add(pendingCashback, newCashback)
		
		now := time.Now()
		stats.TotalCashback = totalCashback.String()
		stats.PendingCashback = pendingCashback.String()
		stats.TotalDeposits++
		stats.LastDepositAt = &now
		if stats.FirstDepositAt == nil {
			stats.FirstDepositAt = &now
		}
		
		// Write back
		data, err := json.Marshal(stats)
		if err != nil {
			return fmt.Errorf("failed to marshal stats: %w", err)
		}
		
		if err := s.StoreMarketplaceData(ctx, key, data); err != nil {
			if i == maxRetries-1 {
				return fmt.Errorf("failed to write stats after %d retries: %w", maxRetries, err)
			}
			time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
			continue
		}
		
		return nil
	}
	
	return fmt.Errorf("failed to update stats after %d retries", maxRetries)
}

// trackUserForPeriod adds a user to the period's user list for settlement
func (s *KVStore) trackUserForPeriod(ctx context.Context, period uint64, userAddress string) error {
	periodKey := fmt.Sprintf("cashback_period:%d:users", period)
	
	// Get existing users
	data, err := s.GetMarketplaceData(ctx, periodKey)
	var users []string
	if err == nil && len(data) > 0 {
		if err := json.Unmarshal(data, &users); err != nil {
			return fmt.Errorf("failed to unmarshal users: %w", err)
		}
	}
	
	// Check if user already tracked
	for _, u := range users {
		if u == userAddress {
			return nil // Already tracked
		}
	}
	
	// Add user
	users = append(users, userAddress)
	data, err = json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshal users: %w", err)
	}
	
	if err := s.StoreMarketplaceData(ctx, periodKey, data); err != nil {
		return fmt.Errorf("failed to store users: %w", err)
	}
	
	return nil
}

// GetCurrentPeriod returns the current settlement period
// Period 1 starts at a configurable date, increments weekly
func (s *KVStore) GetCurrentPeriod() uint64 {
	// Default to 2025-01-01 if not configured
	// Can be overridden via CASHBACK_PERIOD_START env var (RFC3339 format)
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	if startEnv := os.Getenv("CASHBACK_PERIOD_START"); startEnv != "" {
		if parsed, err := time.Parse(time.RFC3339, startEnv); err == nil {
			startDate = parsed
		} else {
			log.Printf("Warning: Invalid CASHBACK_PERIOD_START format, using default: %v", err)
		}
	}
	
	now := time.Now().UTC()
	
	// Calculate weeks since start
	duration := now.Sub(startDate)
	weeks := duration.Hours() / (24 * 7)
	
	return uint64(weeks) + 1 // Period 1, 2, 3, ...
}

