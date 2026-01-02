package store

import (
	"context"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
)

// Marketplace operations - using dedicated p2p marketplace namespace

// StoreMarketplaceData stores marketplace data in the p2p marketplace namespace
func (s *KVStore) StoreMarketplaceData(ctx context.Context, key string, data []byte) error {
	_, err := s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.p2pMarketplaceNamespaceID,
		Key:         key,
		Value:       data,
	})
	// ... existing StoreMarketplaceData ...
	if err != nil {
		return fmt.Errorf("failed to write marketplace data to KV: %w", err)
	}
	return nil
}

// StoreMarketplaceDataWithTTL stores marketplace data with an expiration time (seconds)
func (s *KVStore) StoreMarketplaceDataWithTTL(ctx context.Context, key string, data []byte, ttl int) error {
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
		NamespaceID: s.p2pMarketplaceNamespaceID,
		KVs:         payload,
	})
	if err != nil {
		return fmt.Errorf("failed to write marketplace data with TTL to KV: %w", err)
	}
	return nil
}

// GetMarketplaceData retrieves marketplace data from the p2p marketplace namespace
func (s *KVStore) GetMarketplaceData(ctx context.Context, key string) ([]byte, error) {
	value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: s.p2pMarketplaceNamespaceID,
		Key:         key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get marketplace data from KV: %w", err)
	}
	return value, nil
}

// DeleteMarketplaceData deletes marketplace data from the p2p marketplace namespace
func (s *KVStore) DeleteMarketplaceData(ctx context.Context, key string) error {
	_, err := s.client.DeleteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.DeleteWorkersKVEntryParams{
		NamespaceID: s.p2pMarketplaceNamespaceID,
		Key:         key,
	})
	if err != nil {
		return fmt.Errorf("failed to delete marketplace data from KV: %w", err)
	}
	return nil
}
