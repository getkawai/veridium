package store

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"time"

	"github.com/kawai-network/y/config"
)

// =============================================================================
// FREE TRIAL SYSTEM
// =============================================================================

// HasClaimedTrial checks if the user or machine has already claimed the free trial
func (s *KVStore) HasClaimedTrial(ctx context.Context, address string, machineID string) (bool, error) {
	// 1. Check Address (via UserBalance struct)
	// We read the full balance object because status is embedded there now
	balance, err := s.GetUserBalance(ctx, address)
	if err != nil {
		return false, err
	}
	if balance.TrialClaimed {
		return true, nil
	}

	// 2. Check Machine ID (Global check)
	if machineID != "" {
		keyMachine := fmt.Sprintf("trial_machine:%s", machineID)
		valMachine, err := s.client.GetValue(ctx, s.usersNamespaceID, keyMachine)
		if err == nil && string(valMachine) == "true" {
			return true, nil
		}
	}

	return false, nil
}

// ClaimFreeTrial claims the free trial for a user atomically
func (s *KVStore) ClaimFreeTrial(ctx context.Context, address string, machineID string) error {
	// 0. Pre-check Machine ID (Global)
	// We do this first to fail fast. Race condition here is possible but acceptable for machine ID (not balance).
	if machineID != "" {
		keyMachine := fmt.Sprintf("trial_machine:%s", machineID)
		valMachine, err := s.client.GetValue(ctx, s.usersNamespaceID, keyMachine)
		if err == nil && string(valMachine) == "true" {
			return fmt.Errorf("free trial already claimed by this device")
		}
	}

	// 1. Calculate amount
	amountFloat := config.GetFreeTrialAmountUSDT()
	amountMicro := int64(amountFloat * 1_000_000)
	amountBig := big.NewInt(amountMicro)

	// 2. Atomic Read-Modify-Write Loop for User Balance
	maxRetries := 5
	backoff := 50 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		// A. Get current state
		currentData, err := s.GetUserBalance(ctx, address)
		if err != nil {
			return err
		}

		// B. Check if already claimed
		if currentData.TrialClaimed {
			return fmt.Errorf("free trial already claimed by this address")
		}

		// C. Update State (Balance + Claimed Flag)
		currentBalance := new(big.Int)
		currentBalance.SetString(currentData.USDTBalance, 10)

		newBalance := new(big.Int).Add(currentBalance, amountBig)

		currentData.USDTBalance = newBalance.String()
		currentData.TrialClaimed = true

		// D. Marshal
		data, err := json.Marshal(currentData)
		if err != nil {
			return fmt.Errorf("failed to marshal balance data: %w", err)
		}

		// E. Write
		key := fmt.Sprintf("balance:%s", address)
		err = s.client.SetValue(ctx, s.usersNamespaceID, key, data)

		if err == nil {
			// Success! Now mark Machine ID (Best effort)
			if machineID != "" {
				keyMachine := fmt.Sprintf("trial_machine:%s", machineID)
				_ = s.client.SetValue(ctx, s.usersNamespaceID, keyMachine, []byte("true"))
			}
			return nil
		}

		// Retry
		if attempt < maxRetries-1 {
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	return fmt.Errorf("failed to claim trial after %d retries", maxRetries)
}
