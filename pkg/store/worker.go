package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/cloudflare/cloudflare-go"
)

// WorkerData represents the data stored for a worker in KV
type WorkerData struct {
	WalletAddress string    `json:"wallet_address"`
	EndpointURL   string    `json:"endpoint_url"`
	HardwareSpecs string    `json:"hardware_specs"`
	RegisteredAt  time.Time `json:"registered_at"`
	LastSeen      time.Time `json:"last_seen"`
	Status        string    `json:"status"`
}

type Store interface {
	SaveWorker(ctx context.Context, data *WorkerData) error
	GetWorker(ctx context.Context, address string) (*WorkerData, error)
	ListWorkers(ctx context.Context) ([]*WorkerData, error)
	GetOnlineWorkers(ctx context.Context) ([]*WorkerData, error)
	UpdateHeartbeat(ctx context.Context, address string) error
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
func (s *KVStore) UpdateHeartbeat(ctx context.Context, address string) error {
	worker, err := s.GetWorker(ctx, address)
	if err != nil {
		return err
	}

	worker.LastSeen = time.Now()
	worker.Status = "online"

	return s.SaveWorker(ctx, worker)
}
