package store

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/cloudflare/cloudflare-go"
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

	_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.holderNamespaceID,
		Key:         key,
		Value:       data,
	})
	if err != nil {
		return fmt.Errorf("failed to write holder data to KV: %w", err)
	}
	return nil
}

// GetHolder retrieves holder information from KV store
func (s *KVStore) GetHolder(ctx context.Context, address string) (*HolderInfo, error) {
	key := fmt.Sprintf("holder:%s", address)

	value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: s.holderNamespaceID,
		Key:         key,
	})
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
		params := cloudflare.ListWorkersKVsParams{
			NamespaceID: s.holderNamespaceID,
			Prefix:      prefix,
		}
		if cursor != "" {
			params.Cursor = cursor
		}

		resp, err := s.client.ListWorkersKVKeys(ctx, cloudflare.AccountIdentifier(s.accountID), params)
		if err != nil {
			return nil, fmt.Errorf("failed to list holder keys: %w", err)
		}

		// Process this page
		for _, keyInfo := range resp.Result {
			value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
				NamespaceID: s.holderNamespaceID,
				Key:         keyInfo.Name,
			})
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
		if resp.ResultInfo.Cursor == "" {
			break
		}
		cursor = resp.ResultInfo.Cursor
	}

	return holders, nil
}

// DeleteHolder removes a holder from the registry
func (s *KVStore) DeleteHolder(ctx context.Context, address string) error {
	key := fmt.Sprintf("holder:%s", address)

	_, err := s.client.DeleteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.DeleteWorkersKVEntryParams{
		NamespaceID: s.holderNamespaceID,
		Key:         key,
	})
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
		params := cloudflare.ListWorkersKVsParams{
			NamespaceID: s.holderNamespaceID,
			Prefix:      prefix,
		}
		if cursor != "" {
			params.Cursor = cursor
		}

		resp, err := s.client.ListWorkersKVKeys(ctx, cloudflare.AccountIdentifier(s.accountID), params)
		if err != nil {
			return 0, fmt.Errorf("failed to count holders: %w", err)
		}

		totalCount += len(resp.Result)

		// Check if there are more pages
		if resp.ResultInfo.Cursor == "" {
			break
		}
		cursor = resp.ResultInfo.Cursor
	}

	return totalCount, nil
}
