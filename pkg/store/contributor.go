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

// ContributorData represents the data stored for a contributor in KV
type ContributorData struct {
	WalletAddress      string    `json:"wallet_address"`
	EndpointURL        string    `json:"endpoint_url"`
	HardwareSpecs      string    `json:"hardware_specs"`
	RegisteredAt       time.Time `json:"registered_at"`
	LastSeen           time.Time `json:"last_seen"`
	Status             string    `json:"status"`
	AccumulatedRewards string    `json:"accumulated_rewards"`        // KAWAI (Phase 1)
	AccumulatedUSDT    string    `json:"accumulated_usdt,omitempty"` // USDT (Phase 2)
}

type Store interface {
	SaveContributor(ctx context.Context, data *ContributorData) error
	GetContributor(ctx context.Context, address string) (*ContributorData, error)
	ListContributors(ctx context.Context) ([]*ContributorData, error)
	GetOnlineContributors(ctx context.Context) ([]*ContributorData, error)
	UpdateHeartbeat(ctx context.Context, address string) error
	SaveMerkleProof(ctx context.Context, address string, data *MerkleProofData) error
	GetMerkleProof(ctx context.Context, address string) (*MerkleProofData, error)
}

type KVStore struct {
	client      *cloudflare.API
	accountID   string
	namespaceID string
}

func NewKVStore(apiToken, accountID, namespaceID string) (*KVStore, error) {
	api, err := cloudflare.NewWithAPIToken(apiToken)
	if err != nil {
		return nil, fmt.Errorf("failed to create cloudflare client: %w", err)
	}

	return &KVStore{
		client:      api,
		accountID:   accountID,
		namespaceID: namespaceID,
	}, nil
}

// SaveContributor stores contributor metadata and mining progress in Cloudflare KV.
// This data (mining results) is used weekly to calculate the 70/30 reward split.
func (s *KVStore) SaveContributor(ctx context.Context, data *ContributorData) error {
	key := data.WalletAddress
	value, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal contributor data: %w", err)
	}

	_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.namespaceID,
		Key:         key,
		Value:       value,
	})
	if err != nil {
		return fmt.Errorf("failed to write to KV: %w", err)
	}

	log.Printf("[Store] Saved contributor: %s", data.WalletAddress)
	return nil
}

func (s *KVStore) GetContributor(ctx context.Context, address string) (*ContributorData, error) {
	key := address

	value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: s.namespaceID,
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

func (s *KVStore) ListContributors(ctx context.Context) ([]*ContributorData, error) {
	// List all keys in the contributor namespace
	resp, err := s.client.ListWorkersKVKeys(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.ListWorkersKVsParams{
		NamespaceID: s.namespaceID,
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
// Regular heartbeats are used for real-time monitoring of node availability.
func (s *KVStore) UpdateHeartbeat(ctx context.Context, address string) error {
	contributor, err := s.GetContributor(ctx, address)
	if err != nil {
		return err
	}

	contributor.LastSeen = time.Now()
	contributor.Status = "online"

	return s.SaveContributor(ctx, contributor)
}

// RecordJobReward distributes rewards based on the 70/30 Rule.
// Supports two phases:
// - Phase 1 (Mining): Contributors earn KAWAI tokens
// - Phase 2 (USDT): Contributors earn USDT (when MAX_SUPPLY reached)
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
		// Formula: (TokenUsage / 1,000,000) * KAWAI_RATE_PER_MILLION
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
		// Formula: (TokenUsage / 1,000,000) * COST_RATE_PER_MILLION
		usdtRate := config.GetCostRatePerMillion()
		// Store USDT in smallest unit (6 decimals for USDT)
		usdtDecimals := new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil)

		// Calculate: (tokenUsage * usdtRate * 1e6) / 1,000,000
		rewardAmount := new(big.Int).Mul(big.NewInt(tokenUsage), big.NewInt(int64(usdtRate*1000000)))
		rewardAmount.Div(rewardAmount, big.NewInt(1000000))
		// rewardAmount is now in USDT micro-units

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
