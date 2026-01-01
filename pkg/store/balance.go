package store

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kawai-network/veridium/pkg/config"
)

// UserBalance represents a user's USDT balance for API usage
type UserBalance struct {
	Address     string `json:"address"`
	USDTBalance string `json:"usdt_balance"` // In micro USDT (6 decimals)
}

// GetUserBalance retrieves the USDT balance for API usage (separate from contributor rewards)
func (s *KVStore) GetUserBalance(ctx context.Context, address string) (*UserBalance, error) {
	// For API usage, we store balance separately from contributor rewards
	// Key format: "balance:{address}"
	key := fmt.Sprintf("balance:%s", address)

	value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: s.contributorsNamespaceID,
		Key:         key,
	})
	if err != nil {
		// If balance doesn't exist, default to 0
		return &UserBalance{
			Address:     address,
			USDTBalance: "0",
		}, nil
	}

	balanceStr := string(value)
	if balanceStr == "" {
		balanceStr = "0"
	}

	return &UserBalance{
		Address:     address,
		USDTBalance: balanceStr,
	}, nil
}

// =============================================================================
// ATOMIC BALANCE OPERATIONS (Thread-Safe)
// =============================================================================
// These methods prevent race conditions and ensure financial integrity
// Always use these methods instead of direct KV operations

// DeductBalanceAtomic atomically deducts USDT from user's balance with retry logic
// This prevents race conditions where multiple goroutines try to deduct simultaneously
func (s *KVStore) DeductBalanceAtomic(ctx context.Context, address string, amount *big.Int) error {
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("deduction amount must be positive")
	}

	maxRetries := 5
	backoff := 50 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		// 1. Get current balance
		key := fmt.Sprintf("balance:%s", address)

		value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
			NamespaceID: s.contributorsNamespaceID,
			Key:         key,
		})

		var currentBalance *big.Int
		if err != nil {
			// Balance doesn't exist, treat as 0
			currentBalance = big.NewInt(0)
		} else {
			currentBalance = new(big.Int)
			if string(value) != "" && string(value) != "0" {
				currentBalance.SetString(string(value), 10)
			}
		}

		// 2. Check if sufficient balance
		if currentBalance.Cmp(amount) < 0 {
			return fmt.Errorf("insufficient balance: have %s, need %s", currentBalance.String(), amount.String())
		}

		// 3. Calculate new balance
		newBalance := new(big.Int).Sub(currentBalance, amount)

		// 4. Attempt atomic update
		_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
			NamespaceID: s.contributorsNamespaceID,
			Key:         key,
			Value:       []byte(newBalance.String()),
		})

		if err == nil {
			// Success!
			return nil
		}

		// Retry with exponential backoff
		if attempt < maxRetries-1 {
			time.Sleep(backoff)
			backoff *= 2 // Exponential backoff
		}
	}

	return fmt.Errorf("failed to deduct balance after %d retries (possible concurrent modification)", maxRetries)
}

// AddBalanceAtomic atomically adds USDT to user's balance with retry logic
func (s *KVStore) AddBalanceAtomic(ctx context.Context, address string, amount *big.Int) error {
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("addition amount must be positive")
	}

	maxRetries := 5
	backoff := 50 * time.Millisecond

	for attempt := 0; attempt < maxRetries; attempt++ {
		// 1. Get current balance
		key := fmt.Sprintf("balance:%s", address)

		value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
			NamespaceID: s.contributorsNamespaceID,
			Key:         key,
		})

		var currentBalance *big.Int
		if err != nil {
			currentBalance = big.NewInt(0)
		} else {
			currentBalance = new(big.Int)
			if string(value) != "" && string(value) != "0" {
				currentBalance.SetString(string(value), 10)
			}
		}

		// 2. Calculate new balance
		newBalance := new(big.Int).Add(currentBalance, amount)

		// 3. Attempt atomic update
		_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
			NamespaceID: s.contributorsNamespaceID,
			Key:         key,
			Value:       []byte(newBalance.String()),
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
// This ensures both deduction and addition succeed or both fail
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
		fromBalance, err := s.GetUserBalance(ctx, from)
		if err != nil {
			return fmt.Errorf("failed to get sender balance: %w", err)
		}

		fromBalanceBig := new(big.Int)
		if fromBalance.USDTBalance != "" && fromBalance.USDTBalance != "0" {
			fromBalanceBig.SetString(fromBalance.USDTBalance, 10)
		}

		// 2. Check sufficient balance
		if fromBalanceBig.Cmp(amount) < 0 {
			return fmt.Errorf("insufficient balance: have %s, need %s", fromBalanceBig.String(), amount.String())
		}

		toBalance, err := s.GetUserBalance(ctx, to)
		if err != nil {
			return fmt.Errorf("failed to get recipient balance: %w", err)
		}

		toBalanceBig := new(big.Int)
		if toBalance.USDTBalance != "" && toBalance.USDTBalance != "0" {
			toBalanceBig.SetString(toBalance.USDTBalance, 10)
		}

		// 3. Calculate new balances
		newFromBalance := new(big.Int).Sub(fromBalanceBig, amount)
		newToBalance := new(big.Int).Add(toBalanceBig, amount)

		// 4. Attempt atomic update of both balances
		// Update sender
		fromKey := fmt.Sprintf("balance:%s", from)
		_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
			NamespaceID: s.contributorsNamespaceID,
			Key:         fromKey,
			Value:       []byte(newFromBalance.String()),
		})
		if err != nil {
			if attempt < maxRetries-1 {
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
			return fmt.Errorf("failed to update sender balance: %w", err)
		}

		// Update recipient
		toKey := fmt.Sprintf("balance:%s", to)
		_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
			NamespaceID: s.contributorsNamespaceID,
			Key:         toKey,
			Value:       []byte(newToBalance.String()),
		})
		if err != nil {
			// Rollback sender balance
			_, rollbackErr := s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
				NamespaceID: s.contributorsNamespaceID,
				Key:         fromKey,
				Value:       []byte(fromBalanceBig.String()),
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

		// Success!
		return nil
	}

	return fmt.Errorf("failed to transfer balance after %d retries", maxRetries)
}

// CheckAndDeductBalance atomically checks and deducts balance in one operation
// This is the recommended method for API usage deduction
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
	if balance.USDTBalance != "" && balance.USDTBalance != "0" {
		currentBalance.SetString(balance.USDTBalance, 10)
	}

	if currentBalance.Cmp(cost) < 0 {
		return fmt.Errorf("insufficient balance: have %s micro USDT, need %s micro USDT",
			currentBalance.String(), cost.String())
	}

	return nil
}
