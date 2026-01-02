package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math/big"
	"sync"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/pkg/config"
)

// contributorLocks provides per-address mutex for serializing balance updates
// This prevents race conditions when multiple goroutines update the same contributor
var contributorLocks sync.Map

// ContributorStatus represents the status of a contributor
type ContributorStatus string

const (
	ContributorStatusOnline  ContributorStatus = "online"
	ContributorStatusOffline ContributorStatus = "offline"
	ContributorStatusDeleted ContributorStatus = "deleted"
	ContributorStatusAdmin   ContributorStatus = "admin"
)

// ContributorData represents the data stored for a contributor in KV
type ContributorData struct {
	WalletAddress      string            `json:"wallet_address"`
	EndpointURL        string            `json:"endpoint_url"`
	HardwareSpecs      string            `json:"hardware_specs"`
	RegisteredAt       time.Time         `json:"registered_at"`
	LastSeen           time.Time         `json:"last_seen"`
	Status             ContributorStatus `json:"status"`
	AccumulatedRewards string            `json:"accumulated_rewards"`        // KAWAI (Phase 1)
	AccumulatedUSDT    string            `json:"accumulated_usdt,omitempty"` // USDT (Phase 2)
	IsActive           bool              `json:"is_active"`                  // Soft delete flag
	DeletedAt          time.Time         `json:"deleted_at,omitempty"`       // When soft deleted
	IsAdmin            bool              `json:"is_admin,omitempty"`         // Admin flag
	Version            int64             `json:"version"`                    // Optimistic locking version
}

// =============================================================================
// CONTRIBUTOR OPERATIONS (Contributors Namespace)
// =============================================================================

// SaveContributor stores contributor metadata and mining progress in Cloudflare KV.
// Key format: {address} (lowercase)
func (s *KVStore) SaveContributor(ctx context.Context, data *ContributorData) error {
	key := ContributorKey(data.WalletAddress)
	value, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal contributor data: %w", err)
	}

	_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.contributorsNamespaceID,
		Key:         key,
		Value:       value,
	})
	if err != nil {
		return fmt.Errorf("failed to write to KV: %w", err)
	}

	slog.Info("Saved contributor", "address", data.WalletAddress)
	return nil
}

// GetContributor retrieves contributor data from KV
// Key format: {address} (lowercase)
func (s *KVStore) GetContributor(ctx context.Context, address string) (*ContributorData, error) {
	key := ContributorKey(address)

	value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: s.contributorsNamespaceID,
		Key:         key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get from KV: %w", err)
	}

	var data ContributorData
	if err := json.Unmarshal(value, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal contributor data: %w", err)
	}

	return &data, nil
}

// ListContributors returns all contributors (no prefix needed - entire namespace is contributors)
func (s *KVStore) ListContributors(ctx context.Context) ([]*ContributorData, error) {
	resp, err := s.client.ListWorkersKVKeys(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.ListWorkersKVsParams{
		NamespaceID: s.contributorsNamespaceID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	var contributors []*ContributorData
	for _, key := range resp.Result {
		data, err := s.GetContributor(ctx, key.Name)
		if err != nil {
			slog.Warn("Failed to get contributor data", "key", key.Name, "error", err)
			continue
		}
		contributors = append(contributors, data)
	}

	return contributors, nil
}

// GetOnlineContributors returns contributors with recent heartbeat
func (s *KVStore) GetOnlineContributors(ctx context.Context) ([]*ContributorData, error) {
	contributors, err := s.ListContributors(ctx)
	if err != nil {
		return nil, err
	}

	var online []*ContributorData
	expiration := time.Now().Add(-2 * time.Minute)

	for _, c := range contributors {
		if c.LastSeen.After(expiration) {
			online = append(online, c)
		}
	}

	return online, nil
}

// UpdateHeartbeat updates the timestamp and online status for a contributor.
func (s *KVStore) UpdateHeartbeat(ctx context.Context, address string) error {
	contributor, err := s.GetContributor(ctx, address)
	if err != nil {
		return err
	}

	contributor.LastSeen = time.Now()
	contributor.Status = ContributorStatusOnline

	return s.SaveContributor(ctx, contributor)
}

// ResetAccumulatedRewards resets the accumulated rewards for a contributor after settlement.
// rewardType can be "kawai" or "usdt"
func (s *KVStore) ResetAccumulatedRewards(ctx context.Context, address string, rewardType string) error {
	contributor, err := s.GetContributor(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to get contributor: %w", err)
	}

	switch rewardType {
	case "kawai":
		oldBalance := contributor.AccumulatedRewards
		contributor.AccumulatedRewards = "0"
		slog.Info("Reset KAWAI rewards", "address", address, "old_balance", oldBalance)
	case "usdt":
		oldBalance := contributor.AccumulatedUSDT
		contributor.AccumulatedUSDT = "0"
		slog.Info("Reset USDT rewards", "address", address, "old_balance", oldBalance)
	default:
		return fmt.Errorf("invalid reward type: %s (must be 'kawai' or 'usdt')", rewardType)
	}

	return s.SaveContributor(ctx, contributor)
}

// SoftDeleteContributor marks a contributor as inactive (soft delete)
func (s *KVStore) SoftDeleteContributor(ctx context.Context, address string) error {
	contributor, err := s.GetContributor(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to get contributor: %w", err)
	}

	contributor.IsActive = false
	contributor.DeletedAt = time.Now()
	contributor.Status = ContributorStatusDeleted

	if err := s.SaveContributor(ctx, contributor); err != nil {
		return fmt.Errorf("failed to soft delete contributor: %w", err)
	}

	slog.Info("Soft deleted contributor", "address", address)
	return nil
}

// RestoreContributor restores a soft-deleted contributor
func (s *KVStore) RestoreContributor(ctx context.Context, address string) error {
	contributor, err := s.GetContributor(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to get contributor: %w", err)
	}

	contributor.IsActive = true
	contributor.DeletedAt = time.Time{}
	contributor.Status = ContributorStatusOffline

	if err := s.SaveContributor(ctx, contributor); err != nil {
		return fmt.Errorf("failed to restore contributor: %w", err)
	}

	slog.Info("Restored contributor", "address", address)
	return nil
}

// ListActiveContributors returns only active (non-deleted) contributors
func (s *KVStore) ListActiveContributors(ctx context.Context) ([]*ContributorData, error) {
	contributors, err := s.ListContributors(ctx)
	if err != nil {
		return nil, err
	}

	active := make([]*ContributorData, 0)
	for _, c := range contributors {
		if c.IsActive {
			active = append(active, c)
		}
	}

	return active, nil
}

// ListContributorsWithBalance returns contributors with non-zero balance (for settlement)
func (s *KVStore) ListContributorsWithBalance(ctx context.Context, rewardType string) ([]*ContributorData, error) {
	contributors, err := s.ListContributors(ctx)
	if err != nil {
		return nil, err
	}

	withBalance := make([]*ContributorData, 0)
	for _, c := range contributors {
		var balance string
		if rewardType == "kawai" {
			balance = c.AccumulatedRewards
		} else {
			balance = c.AccumulatedUSDT
		}

		if balance != "" && balance != "0" {
			withBalance = append(withBalance, c)
		}
	}

	return withBalance, nil
}

// RegisterContributor registers a new contributor with proper initialization
func (s *KVStore) RegisterContributor(ctx context.Context, address, endpointURL, hardwareSpecs string) (*ContributorData, error) {
	// Check if already exists
	existing, err := s.GetContributor(ctx, address)
	if err == nil && existing != nil {
		// If soft deleted, restore
		if !existing.IsActive {
			existing.IsActive = true
			existing.DeletedAt = time.Time{}
			existing.EndpointURL = endpointURL
			existing.HardwareSpecs = hardwareSpecs
			existing.Status = ContributorStatusOnline
			existing.LastSeen = time.Now()

			if err := s.SaveContributor(ctx, existing); err != nil {
				return nil, fmt.Errorf("failed to restore contributor: %w", err)
			}
			slog.Info("Restored and updated contributor", "address", address)
			return existing, nil
		}

		// Already active
		return existing, nil
	}

	// Create new contributor
	contributor := &ContributorData{
		WalletAddress:      address,
		EndpointURL:        endpointURL,
		HardwareSpecs:      hardwareSpecs,
		RegisteredAt:       time.Now(),
		LastSeen:           time.Now(),
		Status:             ContributorStatusOnline,
		IsActive:           true,
		AccumulatedRewards: "0",
		AccumulatedUSDT:    "0",
	}

	if err := s.SaveContributor(ctx, contributor); err != nil {
		return nil, fmt.Errorf("failed to register contributor: %w", err)
	}

	slog.Info("Registered new contributor", "address", address)
	return contributor, nil
}

// RecordJobReward distributes rewards based on the 70/30 Rule.
// Admin address and reward mode are automatically determined.
// This method is thread-safe using per-address mutex to prevent race conditions.
//
//wails:ignore
func (s *KVStore) RecordJobReward(ctx context.Context, contributorAddress string, tokenUsage int64) error {
	// Get random admin address from treasury pool
	adminAddress := constant.GetRandomTreasuryAddress()

	// Check if we reached max supply (1B tokens)
	mode := config.ModeMining
	if s.supplyQuerier != nil {
		currentSupply, _ := s.supplyQuerier.GetTotalSupply(ctx)
		maxSupply, _ := s.supplyQuerier.GetMaxSupply(ctx)
		if currentSupply != nil && maxSupply != nil && currentSupply.Cmp(maxSupply) >= 0 {
			mode = config.ModeUSDT // Max supply reached, switch to USDT
		}
	}

	// Helper to update balance field with per-address locking
	// This ensures only ONE goroutine can update a specific address at a time
	updateBalance := func(addr string, amount *big.Int, field string) error {
		if amount.Cmp(big.NewInt(0)) == 0 {
			return nil
		}

		// Get or create mutex for this specific address
		lockInterface, _ := contributorLocks.LoadOrStore(addr, &sync.Mutex{})
		lock := lockInterface.(*sync.Mutex)

		// Acquire lock - blocks if another goroutine is updating this address
		lock.Lock()
		defer lock.Unlock()

		// Now we have exclusive access to this address's balance
		// Safe to do read-modify-write without race conditions

		// 1. READ - Get current contributor data
		c, err := s.GetContributor(ctx, addr)
		if err != nil {
			return fmt.Errorf("failed to get account %s: %w", addr, err)
		}

		// 2. MODIFY - Calculate new balance
		var currentBalStr string
		if field == "kawai" {
			currentBalStr = c.AccumulatedRewards
		} else {
			currentBalStr = c.AccumulatedUSDT
		}

		currentBal := new(big.Int)
		if currentBalStr != "" {
			currentBal.SetString(currentBalStr, 10)
		}

		newBal := new(big.Int).Add(currentBal, amount)

		if field == "kawai" {
			c.AccumulatedRewards = newBal.String()
		} else {
			c.AccumulatedUSDT = newBal.String()
		}

		// 3. WRITE - Save updated balance
		// No retry needed - mutex ensures no concurrent modification
		err = s.SaveContributor(ctx, c)
		if err != nil {
			return fmt.Errorf("failed to save contributor %s: %w", addr, err)
		}

		slog.Info("Balance updated",
			"address", addr,
			"field", field,
			"amount", amount.String(),
			"new_balance", newBal.String())

		return nil
	}

	var contributorShare, adminShare *big.Int
	var balanceField string

	if mode == config.ModeMining {
		// Phase 1: KAWAI Mining with Dynamic Difficulty (Halving)

		// Default Rate: 100 KAWAI per Million Tokens
		rateVal := int64(100)

		var currentSupply *big.Int
		if s.supplyQuerier != nil {
			var err error
			currentSupply, err = s.supplyQuerier.GetTotalSupply(ctx)
			if err != nil {
				slog.Warn("Failed to fetch total supply for halving logic, using default rate", "error", err)
			}
		}

		if currentSupply != nil {
			// Define Halving Thresholds (Tokens with 18 decimals)
			// Threshold 1: 50% Mined (500M)
			// Threshold 2: 75% Mined (750M)
			// Threshold 3: 87.5% Mined (875M)

			exp18 := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
			supply500M := new(big.Int).Mul(big.NewInt(500000000), exp18)
			supply750M := new(big.Int).Mul(big.NewInt(750000000), exp18)
			supply875M := new(big.Int).Mul(big.NewInt(875000000), exp18)

			if currentSupply.Cmp(supply875M) >= 0 {
				rateVal = 12 // Halving 3
			} else if currentSupply.Cmp(supply750M) >= 0 {
				rateVal = 25 // Halving 2
			} else if currentSupply.Cmp(supply500M) >= 0 {
				rateVal = 50 // Halving 1
			}
		}

		// baseRate = rate * 10^18 (KAWAI decimals)
		baseRate := new(big.Int).Mul(big.NewInt(rateVal), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

		// High-precision calculation: (tokens * baseRate * share) / (1,000,000 * 100)
		totalScaled := new(big.Int).Mul(big.NewInt(tokenUsage), baseRate)

		contributorShare = new(big.Int).Mul(totalScaled, big.NewInt(70))
		contributorShare.Div(contributorShare, big.NewInt(100000000)) // 1,000,000 * 100

		totalReward := new(big.Int).Div(totalScaled, big.NewInt(1000000))
		adminShare = new(big.Int).Sub(totalReward, contributorShare)

		balanceField = "kawai"
		slog.Info("Mining reward distributed", "tokens", tokenUsage, "rate", rateVal, "kawai_total", totalReward.String(), "contributor_share", contributorShare.String(), "admin_share", adminShare.String())
	} else {
		// Phase 2: USDT Payment
		usdtRate := config.GetCostRatePerMillion()
		// usdtRateUnits = rate * 10^6 (USDT decimals)
		usdtRateUnits := big.NewInt(int64(usdtRate * 1000000))

		// High-precision calculation: (tokens * rateUnits * share) / (1,000,000 * 100)
		totalScaled := new(big.Int).Mul(big.NewInt(tokenUsage), usdtRateUnits)

		contributorShare = new(big.Int).Mul(totalScaled, big.NewInt(70))
		contributorShare.Div(contributorShare, big.NewInt(100000000)) // 1,000,000 * 100

		totalReward := new(big.Int).Div(totalScaled, big.NewInt(1000000))
		adminShare = new(big.Int).Sub(totalReward, contributorShare)

		balanceField = "usdt"
		slog.Info("USDT reward distributed", "tokens", tokenUsage, "usdt_total", totalReward.String(), "contributor_share", contributorShare.String(), "admin_share", adminShare.String())
	}

	// Update Contributor Balance
	if err := updateBalance(contributorAddress, contributorShare, balanceField); err != nil {
		return err
	}

	// Update Admin Balance
	if adminShare.Cmp(big.NewInt(0)) > 0 {
		// Ensure admin account exists before updating balance
		if err := s.EnsureAdminExists(ctx, adminAddress); err != nil {
			return fmt.Errorf("failed to ensure admin exists: %w", err)
		}

		if err := updateBalance(adminAddress, adminShare, balanceField); err != nil {
			return fmt.Errorf("failed to update admin fee: %w", err)
		}
	}

	return nil
}
