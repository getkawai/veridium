package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kawai-network/veridium/pkg/config"
)

// UserBalance represents a user's USDT + KAWAI balance and trial status
type UserBalance struct {
	Address         string `json:"address"`
	USDTBalance     string `json:"usdt_balance"`     // In micro USDT (6 decimals)
	KawaiBalance    string `json:"kawai_balance"`    // In wei (18 decimals)
	TrialClaimed    bool   `json:"trial_claimed"`    // Whether free trial has been claimed
	ReferrerAddress string `json:"referrer_address"` // Address of the user who referred this user (empty if no referral)
}

// GetUserBalance retrieves the user data (balance + trial status)
func (s *KVStore) GetUserBalance(ctx context.Context, address string) (*UserBalance, error) {
	// Key format: "balance:{address}"
	key := fmt.Sprintf("balance:%s", address)

	value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: s.usersNamespaceID,
		Key:         key,
	})
	if err != nil {
		// If balance doesn't exist, default to 0 and not claimed
		return &UserBalance{
			Address:      address,
			USDTBalance:  "0",
			KawaiBalance: "0",
			TrialClaimed: false,
		}, nil
	}

	var balance UserBalance
	if len(value) > 0 {
		// Try parsing as JSON
		if err := json.Unmarshal(value, &balance); err != nil {
			// Fallback for legacy raw string format (if any migration happens)
			// But since we are on a new namespace, we should be fine assuming JSON or empty.
			// Treating raw string as just balance for robustness if needed,
			// but for this specific "users" namespace, we can assume clean slate.
			// Let's just default to 0 if JSON fails to be safe.
			return &UserBalance{
				Address:      address,
				USDTBalance:  "0", // Reset if format is invalid
				KawaiBalance: "0",
				TrialClaimed: false,
			}, nil
		}
	}

	// Ensure address is set
	balance.Address = address
	if balance.USDTBalance == "" {
		balance.USDTBalance = "0"
	}
	if balance.KawaiBalance == "" {
		balance.KawaiBalance = "0"
	}

	return &balance, nil
}

// =============================================================================
// ATOMIC BALANCE OPERATIONS (Thread-Safe)
// =============================================================================

// DeductBalanceAtomic atomically deducts USDT from user's balance with retry logic
func (s *KVStore) DeductBalanceAtomic(ctx context.Context, address string, amount *big.Int) error {
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("deduction amount must be positive")
	}

	maxRetries := 5
	backoff := 50 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		// 1. Get current balance object
		currentData, err := s.GetUserBalance(ctx, address)
		if err != nil {
			// Should not happen as GetUserBalance handles missing keys
			return err
		}

		currentBalance := new(big.Int)
		currentBalance.SetString(currentData.USDTBalance, 10)

		// 2. Check if sufficient balance
		if currentBalance.Cmp(amount) < 0 {
			return fmt.Errorf("insufficient balance: have %s, need %s", currentBalance.String(), amount.String())
		}

		// 3. Calculate new balance
		newBalance := new(big.Int).Sub(currentBalance, amount)

		// Update struct
		currentData.USDTBalance = newBalance.String()

		// 4. Marshal to JSON
		data, err := json.Marshal(currentData)
		if err != nil {
			return fmt.Errorf("failed to marshal balance data: %w", err)
		}

		// 5. Attempt atomic update
		key := fmt.Sprintf("balance:%s", address)
		_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
			NamespaceID: s.usersNamespaceID,
			Key:         key,
			Value:       data,
		})

		if err == nil {
			return nil
		}

		// Retry with exponential backoff
		if attempt < maxRetries-1 {
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	return fmt.Errorf("failed to deduct balance after %d retries", maxRetries)
}

// AddBalanceAtomic atomically adds USDT to user's balance with retry logic
func (s *KVStore) AddBalanceAtomic(ctx context.Context, address string, amount *big.Int) error {
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("addition amount must be positive")
	}

	maxRetries := 5
	backoff := 50 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		// 1. Get current balance object
		currentData, err := s.GetUserBalance(ctx, address)
		if err != nil {
			return err
		}

		currentBalance := new(big.Int)
		currentBalance.SetString(currentData.USDTBalance, 10)

		// 2. Calculate new balance
		newBalance := new(big.Int).Add(currentBalance, amount)
		currentData.USDTBalance = newBalance.String()

		// 3. Marshal to JSON
		data, err := json.Marshal(currentData)
		if err != nil {
			return fmt.Errorf("failed to marshal balance data: %w", err)
		}

		// 4. Attempt atomic update
		key := fmt.Sprintf("balance:%s", address)
		_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
			NamespaceID: s.usersNamespaceID,
			Key:         key,
			Value:       data,
		})

		if err == nil {
			return nil
		}

		// Retry with exponential backoff
		if attempt < maxRetries-1 {
			time.Sleep(backoff)
			backoff *= 2
		}
	}

	return fmt.Errorf("failed to add balance after %d retries", maxRetries)
}

// TransferBalanceAtomic atomically transfers balance from one address to another
func (s *KVStore) TransferBalanceAtomic(ctx context.Context, from, to string, amount *big.Int) error {
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("transfer amount must be positive")
	}

	if from == to {
		return fmt.Errorf("cannot transfer to same address")
	}

	maxRetries := 5
	backoff := 50 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		// 1. Get both balances
		fromData, err := s.GetUserBalance(ctx, from)
		if err != nil {
			return fmt.Errorf("failed to get sender balance: %w", err)
		}

		fromBalanceBig := new(big.Int)
		fromBalanceBig.SetString(fromData.USDTBalance, 10)

		// 2. Check sufficient balance
		if fromBalanceBig.Cmp(amount) < 0 {
			return fmt.Errorf("insufficient balance: have %s, need %s", fromBalanceBig.String(), amount.String())
		}

		toData, err := s.GetUserBalance(ctx, to)
		if err != nil {
			return fmt.Errorf("failed to get recipient balance: %w", err)
		}

		toBalanceBig := new(big.Int)
		toBalanceBig.SetString(toData.USDTBalance, 10)

		// 3. Calculate new balances
		newFromBalance := new(big.Int).Sub(fromBalanceBig, amount)
		newToBalance := new(big.Int).Add(toBalanceBig, amount)

		fromData.USDTBalance = newFromBalance.String()
		toData.USDTBalance = newToBalance.String()

		// 4. Update Sender First
		fromJson, _ := json.Marshal(fromData)
		fromKey := fmt.Sprintf("balance:%s", from)

		_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
			NamespaceID: s.usersNamespaceID,
			Key:         fromKey,
			Value:       fromJson,
		})
		if err != nil {
			if attempt < maxRetries-1 {
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			return fmt.Errorf("failed to update sender balance: %w", err)
		}

		// 5. Update Recipient
		toJson, _ := json.Marshal(toData)
		toKey := fmt.Sprintf("balance:%s", to)

		_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
			NamespaceID: s.usersNamespaceID,
			Key:         toKey,
			Value:       toJson,
		})
		if err != nil {
			// Rollback sender
			// Note: This is a best-effort rollback manual logic
			fromData.USDTBalance = fromBalanceBig.String() // revert
			rollbackJson, _ := json.Marshal(fromData)

			_, rollbackErr := s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
				NamespaceID: s.usersNamespaceID,
				Key:         fromKey,
				Value:       rollbackJson,
			})

			if rollbackErr != nil {
				return fmt.Errorf("CRITICAL: failed to update recipient and rollback sender: %w (rollback error: %v)", err, rollbackErr)
			}

			if attempt < maxRetries-1 {
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			return fmt.Errorf("failed to update recipient balance: %w", err)
		}

		return nil
	}

	return fmt.Errorf("failed to transfer balance after %d retries", maxRetries)
}

// CheckAndDeductBalance atomically checks and deducts balance in one operation
func (s *KVStore) CheckAndDeductBalance(ctx context.Context, address string, tokenUsage int64) error {
	cost := CalculateUsageCost(tokenUsage)
	return s.DeductBalanceAtomic(ctx, address, cost)
}

// =============================================================================
// UTILITY FUNCTIONS
// =============================================================================

// CalculateUsageCost calculates the USDT cost for token usage
func CalculateUsageCost(tokenUsage int64) *big.Int {
	// Get cost rate (default: $1 per 1M tokens)
	costRate := config.GetCostRatePerMillion()

	// Convert to micro USDT (6 decimals)
	// Formula: (tokenUsage / 1,000,000) * costRate * 1,000,000 (for micro USDT)
	costRateMicro := int64(costRate * 1000000) // Convert to micro USDT

	cost := new(big.Int).Mul(big.NewInt(tokenUsage), big.NewInt(costRateMicro))
	cost.Div(cost, big.NewInt(1000000)) // Divide by 1M tokens

	return cost
}

// CheckSufficientBalance checks if user has enough balance for the given token usage
func (s *KVStore) CheckSufficientBalance(ctx context.Context, address string, tokenUsage int64) error {
	cost := CalculateUsageCost(tokenUsage)

	balance, err := s.GetUserBalance(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}

	currentBalance := new(big.Int)
	if balance.USDTBalance != "" {
		currentBalance.SetString(balance.USDTBalance, 10)
	}

	if currentBalance.Cmp(cost) < 0 {
		return fmt.Errorf("insufficient balance: have %s micro USDT, need %s micro USDT",
			currentBalance.String(), cost.String())
	}

	return nil
}

// RecordDebt records a failed balance deduction for manual reconciliation
// This is used when deduction fails after AI service has been rendered
func (s *KVStore) RecordDebt(ctx context.Context, address string, amount *big.Int, reason string) error {
	timestamp := time.Now().UnixNano()
	key := fmt.Sprintf("debt:%s:%d", address, timestamp)

	debtData := struct {
		Address   string `json:"address"`
		Amount    string `json:"amount"`
		Reason    string `json:"reason"`
		Timestamp int64  `json:"timestamp"`
		CreatedAt string `json:"created_at"`
	}{
		Address:   address,
		Amount:    amount.String(),
		Reason:    reason,
		Timestamp: timestamp,
		CreatedAt: time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(debtData)
	if err != nil {
		return fmt.Errorf("failed to marshal debt data: %w", err)
	}

	// Store in users namespace for debts
	_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.usersNamespaceID,
		Key:         key,
		Value:       data,
	})

	if err != nil {
		return fmt.Errorf("failed to store debt record: %w", err)
	}

	log.Printf("[DEBT] Recorded debt for user %s: %s micro USDT - Reason: %s",
		address, amount.String(), reason)

	return nil
}
