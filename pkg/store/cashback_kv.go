package store

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
)

// Cashback KV operations - using dedicated cashback namespace

// StoreCashbackData stores cashback data in the cashback namespace
func (s *KVStore) StoreCashbackData(ctx context.Context, key string, data []byte) error {
	_, err := s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.cashbackNamespaceID,
		Key:         key,
		Value:       data,
	})
	if err != nil {
		return fmt.Errorf("failed to write cashback data to KV: %w", err)
	}
	return nil
}

// StoreCashbackDataWithTTL stores cashback data with an expiration time (seconds)
func (s *KVStore) StoreCashbackDataWithTTL(ctx context.Context, key string, data []byte, ttl int) error {
	// Use Bulk Write API because Single Write API in this SDK version doesn't support TTL params
	// Value must be string
	payload := []*cloudflare.WorkersKVPair{
		{
			Key:           key,
			Value:         string(data),
			ExpirationTTL: ttl,
		},
	}

	_, err := s.client.WriteWorkersKVEntries(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntriesParams{
		NamespaceID: s.cashbackNamespaceID,
		KVs:         payload,
	})
	if err != nil {
		return fmt.Errorf("failed to write cashback data with TTL to KV: %w", err)
	}
	return nil
}

// GetCashbackData retrieves cashback data from the cashback namespace
func (s *KVStore) GetCashbackData(ctx context.Context, key string) ([]byte, error) {
	value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: s.cashbackNamespaceID,
		Key:         key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get cashback data from KV: %w", err)
	}
	return value, nil
}

// DeleteCashbackData deletes cashback data from the cashback namespace
func (s *KVStore) DeleteCashbackData(ctx context.Context, key string) error {
	_, err := s.client.DeleteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.DeleteWorkersKVEntryParams{
		NamespaceID: s.cashbackNamespaceID,
		Key:         key,
	})
	if err != nil {
		return fmt.Errorf("failed to delete cashback data from KV: %w", err)
	}
	return nil
}
