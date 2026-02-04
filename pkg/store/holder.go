package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
)

// HolderInfo represents information about a KAWAI holder
type HolderInfo struct {
	Address    string `json:"address"`
	LastSeen   int64  `json:"lastSeen"`
	Source     string `json:"source"` // "desktop", "cli", "transfer"
	Registered int64  `json:"registered"`
}

// SaveHolder saves holder information to KV store
// Key format: holder:{address}
func (s *KVStore) SaveHolder(ctx context.Context, address string, info *HolderInfo) error {
	key := fmt.Sprintf("holder:%s", address)

	data, err := json.Marshal(info)
	if err != nil {
		return fmt.Errorf("failed to marshal holder info: %w", err)
	}

	err = s.client.SetValue(ctx, s.holderNamespaceID, key, data)
	if err != nil {
		return fmt.Errorf("failed to write holder data to KV: %w", err)
	}
	return nil
}

// GetHolder retrieves holder information from KV store
func (s *KVStore) GetHolder(ctx context.Context, address string) (*HolderInfo, error) {
	key := fmt.Sprintf("holder:%s", address)

	value, err := s.client.GetValue(ctx, s.holderNamespaceID, key)
	if err != nil {
		return nil, fmt.Errorf("failed to get holder data from KV: %w", err)
	}

	var info HolderInfo
	if err := json.Unmarshal(value, &info); err != nil {
		return nil, fmt.Errorf("failed to unmarshal holder info: %w", err)
	}

	return &info, nil
}

// ListHolders retrieves all registered KAWAI holders
func (s *KVStore) ListHolders(ctx context.Context) ([]*HolderInfo, error) {
	// List all keys with prefix "holder:" with pagination support
	prefix := "holder:"
	var holders []*HolderInfo
	cursor := ""

	for {
		result, err := s.client.ListKeys(ctx, s.holderNamespaceID, prefix, cursor)
		if err != nil {
			return nil, fmt.Errorf("failed to list holder keys: %w", err)
		}

		// Process this page
		for _, keyInfo := range result.Result {
			value, err := s.client.GetValue(ctx, s.holderNamespaceID, keyInfo.Name)
			if err != nil {
				slog.Warn("Failed to get holder data", "key", keyInfo.Name, "error", err)
				continue
			}

			var info HolderInfo
			if err := json.Unmarshal(value, &info); err != nil {
				slog.Warn("Failed to unmarshal holder data", "key", keyInfo.Name, "error", err)
				continue
			}
			holders = append(holders, &info)
		}

		// Check if there are more pages
		if result.ResultInfo.Cursor == "" {
			break
		}
		cursor = result.ResultInfo.Cursor
	}

	return holders, nil
}

// DeleteHolder removes a holder from the registry
func (s *KVStore) DeleteHolder(ctx context.Context, address string) error {
	key := fmt.Sprintf("holder:%s", address)

	err := s.client.DeleteValue(ctx, s.holderNamespaceID, key)
	if err != nil {
		return fmt.Errorf("failed to delete holder data from KV: %w", err)
	}
	return nil
}

// GetHolderCount returns the total number of registered holders
func (s *KVStore) GetHolderCount(ctx context.Context) (int, error) {
	prefix := "holder:"
	totalCount := 0
	cursor := ""

	for {
		result, err := s.client.ListKeys(ctx, s.holderNamespaceID, prefix, cursor)
		if err != nil {
			return 0, fmt.Errorf("failed to count holders: %w", err)
		}

		totalCount += result.ResultInfo.Count

		// Check if there are more pages
		if result.ResultInfo.Cursor == "" {
			break
		}
		cursor = result.ResultInfo.Cursor
	}

	return totalCount, nil
}
