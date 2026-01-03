package store

import (
	"context"
	"fmt"
	"math/big"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kawai-network/veridium/pkg/config"
)

// =============================================================================
// FREE TRIAL SYSTEM
// =============================================================================

// HasClaimedTrial checks if the user or machine has already claimed the free trial
func (s *KVStore) HasClaimedTrial(ctx context.Context, address string, machineID string) (bool, error) {
	// 1. Check Address
	keyAddr := fmt.Sprintf("trial_claimed:%s", address)
	valAddr, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: s.contributorsNamespaceID,
		Key:         keyAddr,
	})
	if err == nil && string(valAddr) == "true" {
		return true, nil
	}

	// 2. Check Machine ID
	if machineID != "" {
		keyMachine := fmt.Sprintf("trial_machine:%s", machineID)
		valMachine, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
			NamespaceID: s.contributorsNamespaceID,
			Key:         keyMachine,
		})
		if err == nil && string(valMachine) == "true" {
			return true, nil
		}
	}

	return false, nil
}

// ClaimFreeTrial claims the free trial for a user (adds USDT credits)
func (s *KVStore) ClaimFreeTrial(ctx context.Context, address string, machineID string) error {
	// 1. Check if already claimed (by address or machine)
	claimed, err := s.HasClaimedTrial(ctx, address, machineID)
	if err != nil {
		return fmt.Errorf("failed to check trial status: %w", err)
	}
	if claimed {
		return fmt.Errorf("free trial already claimed (by address or device)")
	}

	// 2. Calculate amount
	amountFloat := config.GetFreeTrialAmountUSDT()
	amountMicro := int64(amountFloat * 1_000_000)
	amountBig := big.NewInt(amountMicro)

	// 3. Add balance atomically
	if err := s.AddBalanceAtomic(ctx, address, amountBig); err != nil {
		return fmt.Errorf("failed to add trial balance: %w", err)
	}

	// 4. Mark Address as claimed
	keyAddr := fmt.Sprintf("trial_claimed:%s", address)
	_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.contributorsNamespaceID,
		Key:         keyAddr,
		Value:       []byte("true"),
	})
	if err != nil {
		return fmt.Errorf("failed to mark address as claimed: %w", err)
	}

	// 5. Mark Machine ID as claimed
	if machineID != "" {
		keyMachine := fmt.Sprintf("trial_machine:%s", machineID)
		_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
			NamespaceID: s.contributorsNamespaceID,
			Key:         keyMachine,
			Value:       []byte("true"),
		})
		if err != nil {
			// Log error but don't fail, address is already marked
			fmt.Printf("Warning: failed to mark machine ID %s as claimed: %v\n", machineID, err)
		}
	}

	return nil
}
