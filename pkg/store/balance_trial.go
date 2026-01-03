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

// HasClaimedTrial checks if the user has already claimed the free trial
func (s *KVStore) HasClaimedTrial(ctx context.Context, address string) (bool, error) {
	key := fmt.Sprintf("trial_claimed:%s", address)
	value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: s.contributorsNamespaceID,
		Key:         key,
	})

	if err != nil {
		// Key likely not found, meaning not claimed
		return false, nil
	}

	return string(value) == "true", nil
}

// ClaimFreeTrial claims the free trial for a user (adds USDT credits)
func (s *KVStore) ClaimFreeTrial(ctx context.Context, address string) error {
	// 1. Check if already claimed
	claimed, err := s.HasClaimedTrial(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to check trial status: %w", err)
	}
	if claimed {
		return fmt.Errorf("free trial already claimed")
	}

	// 2. Calculate amount (USDT -> Micro USDT)
	amountFloat := config.GetFreeTrialAmountUSDT()
	amountMicro := int64(amountFloat * 1_000_000)
	amountBig := big.NewInt(amountMicro)

	// 3. Add balance atomically
	if err := s.AddBalanceAtomic(ctx, address, amountBig); err != nil {
		return fmt.Errorf("failed to add trial balance: %w", err)
	}

	// 4. Mark as claimed
	key := fmt.Sprintf("trial_claimed:%s", address)
	_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.contributorsNamespaceID,
		Key:         key,
		Value:       []byte("true"),
	})

	if err != nil {
		return fmt.Errorf("failed to mark trial as claimed: %w", err)
	}

	return nil
}
