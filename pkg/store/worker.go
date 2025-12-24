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

// WorkerData represents the data stored for a worker in KV
type WorkerData struct {
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
	SaveWorker(ctx context.Context, data *WorkerData) error
	GetWorker(ctx context.Context, address string) (*WorkerData, error)
	ListWorkers(ctx context.Context) ([]*WorkerData, error)
	GetOnlineWorkers(ctx context.Context) ([]*WorkerData, error)
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

// SaveWorker stores worker metadata and mining progress in Cloudflare KV.
// This data (mining results) is used weekly to calculate the 70/30 reward split.
func (s *KVStore) SaveWorker(ctx context.Context, data *WorkerData) error {
	key := data.WalletAddress
	value, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal worker data: %w", err)
	}

	_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.namespaceID,
		Key:         key,
		Value:       value,
	})
	if err != nil {
		return fmt.Errorf("failed to write to KV: %w", err)
	}

	log.Printf("[Store] Saved worker: %s", data.WalletAddress)
	return nil
}

func (s *KVStore) GetWorker(ctx context.Context, address string) (*WorkerData, error) {
	key := address

	value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: s.namespaceID,
		Key:         key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get from KV: %w", err)
	}

	var data WorkerData
	if err := json.Unmarshal(value, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal worker data: %w", err)
	}

	return &data, nil
}

func (s *KVStore) ListWorkers(ctx context.Context) ([]*WorkerData, error) {
	// List all keys in the worker namespace
	resp, err := s.client.ListWorkersKVKeys(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.ListWorkersKVsParams{
		NamespaceID: s.namespaceID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list keys: %w", err)
	}

	var workers []*WorkerData
	for _, key := range resp.Result {
		data, err := s.GetWorker(ctx, key.Name)
		if err != nil {
			log.Printf("[Warning] Failed to get worker data for %s: %v", key.Name, err)
			continue
		}
		workers = append(workers, data)
	}

	return workers, nil
}

func (s *KVStore) GetOnlineWorkers(ctx context.Context) ([]*WorkerData, error) {
	workers, err := s.ListWorkers(ctx)
	if err != nil {
		return nil, err
	}

	var online []*WorkerData
	expiration := time.Now().Add(-2 * time.Minute)

	for _, w := range workers {
		if w.LastSeen.After(expiration) {
			online = append(online, w)
		}
	}

	return online, nil
}

// UpdateHeartbeat updates the timestamp and online status for a worker.
// Regular heartbeats are used for real-time monitoring of node availability.
func (s *KVStore) UpdateHeartbeat(ctx context.Context, address string) error {
	worker, err := s.GetWorker(ctx, address)
	if err != nil {
		return err
	}

	worker.LastSeen = time.Now()
	worker.Status = "online"

	return s.SaveWorker(ctx, worker)
}

// RecordJobReward distributes rewards based on the 70/30 Rule.
// Supports two phases:
// - Phase 1 (Mining): Workers earn KAWAI tokens
// - Phase 2 (USDT): Workers earn USDT (when MAX_SUPPLY reached)
func (s *KVStore) RecordJobReward(ctx context.Context, workerAddress string, tokenUsage int64, adminAddress string, mode config.RewardMode) error {

	// Helper to update balance field
	updateBalance := func(addr string, amount *big.Int, field string) error {
		if amount.Cmp(big.NewInt(0)) == 0 {
			return nil
		}

		w, err := s.GetWorker(ctx, addr)
		if err != nil {
			return fmt.Errorf("failed to get account %s: %w", addr, err)
		}

		// Select correct balance field
		var currentBalStr string
		if field == "kawai" {
			currentBalStr = w.AccumulatedRewards
		} else {
			currentBalStr = w.AccumulatedUSDT
		}

		currentBal := new(big.Int)
		if currentBalStr != "" {
			currentBal.SetString(currentBalStr, 10)
		}

		currentBal.Add(currentBal, amount)

		if field == "kawai" {
			w.AccumulatedRewards = currentBal.String()
		} else {
			w.AccumulatedUSDT = currentBal.String()
		}

		return s.SaveWorker(ctx, w)
	}

	var workerShare, adminShare *big.Int
	var balanceField string

	if mode == config.ModeMining {
		// Phase 1: KAWAI Mining
		// Formula: (TokenUsage / 1,000,000) * KAWAI_RATE_PER_MILLION
		kawaiRate := config.GetKawaiRatePerMillion()
		baseRate := new(big.Int).Mul(big.NewInt(kawaiRate), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil))

		rewardAmount := new(big.Int).Mul(big.NewInt(tokenUsage), baseRate)
		rewardAmount.Div(rewardAmount, big.NewInt(1000000))

		workerShare = new(big.Int).Mul(rewardAmount, big.NewInt(70))
		workerShare.Div(workerShare, big.NewInt(100))
		adminShare = new(big.Int).Sub(rewardAmount, workerShare)
		balanceField = "kawai"

		log.Printf("[Phase 1 Mining] Job: %d Tokens -> %s KAWAI. Worker: %s | Admin: %s", tokenUsage, rewardAmount.String(), workerShare.String(), adminShare.String())
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

		workerShare = new(big.Int).Mul(rewardAmount, big.NewInt(70))
		workerShare.Div(workerShare, big.NewInt(100))
		adminShare = new(big.Int).Sub(rewardAmount, workerShare)
		balanceField = "usdt"

		log.Printf("[Phase 2 USDT] Job: %d Tokens -> %s USDT (micro). Worker: %s | Admin: %s", tokenUsage, rewardAmount.String(), workerShare.String(), adminShare.String())
		_ = usdtDecimals // silence unused var warning
	}

	// Update Worker Balance
	if err := updateBalance(workerAddress, workerShare, balanceField); err != nil {
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
