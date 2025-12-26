package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/pkg/config"
)

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
}

type Store interface {
	// Contributor operations
	SaveContributor(ctx context.Context, data *ContributorData) error
	GetContributor(ctx context.Context, address string) (*ContributorData, error)
	ListContributors(ctx context.Context) ([]*ContributorData, error)
	ListActiveContributors(ctx context.Context) ([]*ContributorData, error)
	ListContributorsWithBalance(ctx context.Context, rewardType string) ([]*ContributorData, error)
	GetOnlineContributors(ctx context.Context) ([]*ContributorData, error)
	UpdateHeartbeat(ctx context.Context, address string) error
	ResetAccumulatedRewards(ctx context.Context, address string, rewardType string) error
	RegisterContributor(ctx context.Context, address, endpointURL, hardwareSpecs string) (*ContributorData, error)
	SoftDeleteContributor(ctx context.Context, address string) error
	RestoreContributor(ctx context.Context, address string) error

	// Merkle proof operations (deprecated - use period-specific methods)
	SaveMerkleProof(ctx context.Context, address string, data *MerkleProofData) error
	GetMerkleProof(ctx context.Context, address string) (*MerkleProofData, error)

	// Period-specific Merkle proof operations
	SaveMerkleProofForPeriod(ctx context.Context, address string, periodID int64, data *MerkleProofData) error
	GetMerkleProofForPeriod(ctx context.Context, address string, periodID int64) (*MerkleProofData, error)
	ListMerkleProofs(ctx context.Context, address string) ([]*MerkleProofData, error)
	DeleteMerkleProof(ctx context.Context, address string, periodID int64) error

	// Claim status operations
	MarkClaimPending(ctx context.Context, address string, periodID int64, txHash string) error
	ConfirmClaim(ctx context.Context, address string, periodID int64) error
	MarkClaimFailed(ctx context.Context, address string, periodID int64, reason string) error
	RetryFailedClaim(ctx context.Context, address string, periodID int64) error
	GetPendingClaims(ctx context.Context) ([]*MerkleProofData, error)

	// Settlement operations
	GetSettlementSnapshots(ctx context.Context, rewardType string) ([]*SettlementSnapshot, error)
	PerformSettlement(ctx context.Context, periodID int64, merkleRoot string, rewardType string, proofs map[string]*MerkleProofData) (*SettlementPeriod, error)
	PerformSettlementWithConfig(ctx context.Context, periodID int64, merkleRoot string, rewardType string, proofs map[string]*MerkleProofData, config *SettlementConfig) (*SettlementPeriod, error)
	PerformSettlementParallel(ctx context.Context, periodID int64, merkleRoot string, rewardType string, proofs map[string]*MerkleProofData, workers int) (*SettlementPeriod, error)
	ResumeSettlement(ctx context.Context, periodID int64, proofs map[string]*MerkleProofData, config *SettlementConfig) (*SettlementPeriod, error)
	GetClaimableRewards(ctx context.Context, address string) (map[string]interface{}, error)

	// Settlement period operations
	SaveSettlementPeriod(ctx context.Context, period *SettlementPeriod) error
	GetSettlementPeriod(ctx context.Context, periodID int64) (*SettlementPeriod, error)
	ListSettlementPeriods(ctx context.Context) ([]*SettlementPeriod, error)

	// Admin operations
	EnsureAdminExists(ctx context.Context, adminAddress string) error
}

// KVStore implements Store interface with multiple namespaces
type KVStore struct {
	client    *cloudflare.API
	accountID string

	// Separate namespace IDs for different data types
	contributorsNamespaceID string
	proofsNamespaceID       string
	settlementsNamespaceID  string
	authzNamespaceID        string // Reverse index: address -> apikey
}

// NewMultiNamespaceKVStore creates a new KVStore with separate namespaces
func NewMultiNamespaceKVStore() (*KVStore, error) {
	api, err := cloudflare.NewWithAPIToken(constant.GetCfApiToken())
	if err != nil {
		return nil, fmt.Errorf("failed to create cloudflare client: %w", err)
	}

	return &KVStore{
		client:                  api,
		accountID:               constant.GetCfAccountId(),
		contributorsNamespaceID: constant.GetCfKvContributorsNamespaceId(),
		proofsNamespaceID:       constant.GetCfKvProofsNamespaceId(),
		settlementsNamespaceID:  constant.GetCfKvSettlementsNamespaceId(),
		authzNamespaceID:        constant.GetCfKvAuthzNamespaceId(),
	}, nil
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

	log.Printf("[Store] Saved contributor: %s", data.WalletAddress)
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
			log.Printf("[Warning] Failed to get contributor data for %s: %v", key.Name, err)
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
		log.Printf("[Store] Reset KAWAI rewards for %s: %s -> 0", address, oldBalance)
	case "usdt":
		oldBalance := contributor.AccumulatedUSDT
		contributor.AccumulatedUSDT = "0"
		log.Printf("[Store] Reset USDT rewards for %s: %s -> 0", address, oldBalance)
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

	log.Printf("[Store] Soft deleted contributor: %s", address)
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

	log.Printf("[Store] Restored contributor: %s", address)
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
			log.Printf("[Store] Restored and updated contributor: %s", address)
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

	log.Printf("[Store] Registered new contributor: %s", address)
	return contributor, nil
}

// RecordJobReward distributes rewards based on the 70/30 Rule.
func (s *KVStore) RecordJobReward(ctx context.Context, contributorAddress string, tokenUsage int64, adminAddress string, mode config.RewardMode) error {

	// Helper to update balance field
	updateBalance := func(addr string, amount *big.Int, field string) error {
		if amount.Cmp(big.NewInt(0)) == 0 {
			return nil
		}

		c, err := s.GetContributor(ctx, addr)
		if err != nil {
			return fmt.Errorf("failed to get account %s: %w", addr, err)
		}

		// Select correct balance field
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

		currentBal.Add(currentBal, amount)

		if field == "kawai" {
			c.AccumulatedRewards = currentBal.String()
		} else {
			c.AccumulatedUSDT = currentBal.String()
		}

		return s.SaveContributor(ctx, c)
	}

	var contributorShare, adminShare *big.Int
	var balanceField string

	if mode == config.ModeMining {
		// Phase 1: KAWAI Mining
		kawaiRate := config.GetKawaiRatePerMillion()
		baseRate := new(big.Int).Mul(big.NewInt(kawaiRate), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

		rewardAmount := new(big.Int).Mul(big.NewInt(tokenUsage), baseRate)
		rewardAmount.Div(rewardAmount, big.NewInt(1000000))

		contributorShare = new(big.Int).Mul(rewardAmount, big.NewInt(70))
		contributorShare.Div(contributorShare, big.NewInt(100))
		adminShare = new(big.Int).Sub(rewardAmount, contributorShare)
		balanceField = "kawai"

		log.Printf("[Phase 1 Mining] Job: %d Tokens -> %s KAWAI. Contributor: %s | Admin: %s", tokenUsage, rewardAmount.String(), contributorShare.String(), adminShare.String())
	} else {
		// Phase 2: USDT Payment
		usdtRate := config.GetCostRatePerMillion()
		usdtDecimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil)

		rewardAmount := new(big.Int).Mul(big.NewInt(tokenUsage), big.NewInt(int64(usdtRate*1000000)))
		rewardAmount.Div(rewardAmount, big.NewInt(1000000))

		contributorShare = new(big.Int).Mul(rewardAmount, big.NewInt(70))
		contributorShare.Div(contributorShare, big.NewInt(100))
		adminShare = new(big.Int).Sub(rewardAmount, contributorShare)
		balanceField = "usdt"

		log.Printf("[Phase 2 USDT] Job: %d Tokens -> %s USDT (micro). Contributor: %s | Admin: %s", tokenUsage, rewardAmount.String(), contributorShare.String(), adminShare.String())
		_ = usdtDecimals // silence unused var warning
	}

	// Update Contributor Balance
	if err := updateBalance(contributorAddress, contributorShare, balanceField); err != nil {
		return err
	}

	// Update Admin Balance
	if adminShare.Cmp(big.NewInt(0)) > 0 {
		if err := updateBalance(adminAddress, adminShare, balanceField); err != nil {
			return fmt.Errorf("failed to update admin fee: %w", err)
		}
	}

	return nil
}
