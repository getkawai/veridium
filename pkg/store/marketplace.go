package store

import (
	"context"
	"fmt"
)

// Marketplace operations - using dedicated p2p marketplace namespace

// StoreMarketplaceData stores marketplace data in the p2p marketplace namespace
func (s *KVStore) StoreMarketplaceData(ctx context.Context, key string, data []byte) error {
	err := s.client.SetValue(ctx, s.p2pMarketplaceNamespaceID, key, data)
	if err != nil {
		return fmt.Errorf("failed to write marketplace data to KV: %w", err)
	}
	return nil
}

// StoreMarketplaceDataWithTTL stores marketplace data with an expiration time (seconds)
func (s *KVStore) StoreMarketplaceDataWithTTL(ctx context.Context, key string, data []byte, ttl int) error {
	err := s.client.SetValueWithTTL(ctx, s.p2pMarketplaceNamespaceID, key, data, ttl)
	if err != nil {
		return fmt.Errorf("failed to write marketplace data with TTL to KV: %w", err)
	}
	return nil
}

// GetMarketplaceData retrieves marketplace data from the p2p marketplace namespace
func (s *KVStore) GetMarketplaceData(ctx context.Context, key string) ([]byte, error) {
	value, err := s.client.GetValue(ctx, s.p2pMarketplaceNamespaceID, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get marketplace data from KV: %w", err)
	}
	return value, nil
}

// DeleteMarketplaceData deletes marketplace data from the p2p marketplace namespace
func (s *KVStore) DeleteMarketplaceData(ctx context.Context, key string) error {
	err := s.client.DeleteValue(ctx, s.p2pMarketplaceNamespaceID, key)
	if err != nil {
		return fmt.Errorf("failed to delete marketplace data from KV: %w", err)
	}
	return nil
}
