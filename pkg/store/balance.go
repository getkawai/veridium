package store

import (
	"context"
	"fmt"
	"math/big"

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

// DeductBalance deducts USDT from user's balance for API usage
func (s *KVStore) DeductBalance(ctx context.Context, address string, amount *big.Int) error {
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("deduction amount must be positive")
	}

	// Get current balance
	balance, err := s.GetUserBalance(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}

	// Parse current balance
	currentBalance := new(big.Int)
	if balance.USDTBalance != "" && balance.USDTBalance != "0" {
		currentBalance.SetString(balance.USDTBalance, 10)
	}

	// Check if sufficient balance
	if currentBalance.Cmp(amount) < 0 {
		return fmt.Errorf("insufficient balance: have %s, need %s", currentBalance.String(), amount.String())
	}

	// Deduct amount
	newBalance := new(big.Int).Sub(currentBalance, amount)

	// Save new balance
	key := fmt.Sprintf("balance:%s", address)
	_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.contributorsNamespaceID,
		Key:         key,
		Value:       []byte(newBalance.String()),
	})
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}

// AddBalance adds USDT to user's balance (for deposits from PaymentVault)
func (s *KVStore) AddBalance(ctx context.Context, address string, amount *big.Int) error {
	if amount.Cmp(big.NewInt(0)) <= 0 {
		return fmt.Errorf("addition amount must be positive")
	}

	// Get current balance
	balance, err := s.GetUserBalance(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to get balance: %w", err)
	}

	// Parse current balance
	currentBalance := new(big.Int)
	if balance.USDTBalance != "" && balance.USDTBalance != "0" {
		currentBalance.SetString(balance.USDTBalance, 10)
	}

	// Add amount
	newBalance := new(big.Int).Add(currentBalance, amount)

	// Save new balance
	key := fmt.Sprintf("balance:%s", address)
	_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.contributorsNamespaceID,
		Key:         key,
		Value:       []byte(newBalance.String()),
	})
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	return nil
}

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
