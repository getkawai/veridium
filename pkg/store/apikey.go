package store

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
	"github.com/kawai-network/veridium/internal/constant"
)

// APIKey represents an API key record
type APIKey struct {
	Key     string `json:"key"`
	Address string `json:"address"`
}

// CreateAPIKey generates a new API key for the given wallet address
func (s *KVStore) CreateAPIKey(ctx context.Context, address string) (*APIKey, error) {
	// Generate random API key (32 bytes = 64 hex chars)
	keyBytes := make([]byte, 32)
	if _, err := rand.Read(keyBytes); err != nil {
		return nil, fmt.Errorf("failed to generate random key: %w", err)
	}

	// Format as "vk-" prefix + hex string (similar to OpenAI's "sk-" format)
	apiKey := "vk-" + hex.EncodeToString(keyBytes)

	// Store in KV: key=apikey, value=address
	_, err := s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: constant.GetCfKvApikeyNamespaceId(),
		Key:         apiKey,
		Value:       []byte(address),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store API key: %w", err)
	}

	return &APIKey{
		Key:     apiKey,
		Address: address,
	}, nil
}

// GetAPIKey retrieves the wallet address associated with an API key
func (s *KVStore) GetAPIKey(ctx context.Context, apiKey string) (string, error) {
	value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: constant.GetCfKvApikeyNamespaceId(),
		Key:         apiKey,
	})
	if err != nil {
		return "", fmt.Errorf("API key not found: %w", err)
	}

	address := string(value)
	if address == "" {
		return "", fmt.Errorf("API key not found")
	}

	return address, nil
}

// ValidateAPIKey checks if an API key exists and returns the associated wallet address
func (s *KVStore) ValidateAPIKey(ctx context.Context, apiKey string) (string, error) {
	if apiKey == "" {
		return "", fmt.Errorf("API key is required")
	}

	// Check format (should start with "vk-")
	if len(apiKey) < 3 || apiKey[:3] != "vk-" {
		return "", fmt.Errorf("invalid API key format")
	}

	return s.GetAPIKey(ctx, apiKey)
}

// RevokeAPIKey removes an API key from the system
func (s *KVStore) RevokeAPIKey(ctx context.Context, apiKey string) error {
	_, err := s.client.DeleteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.DeleteWorkersKVEntryParams{
		NamespaceID: constant.GetCfKvApikeyNamespaceId(),
		Key:         apiKey,
	})
	if err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}

	return nil
}

// ListAPIKeys returns all API keys for a given wallet address (for dashboard)
func (s *KVStore) ListAPIKeys(ctx context.Context, address string) ([]string, error) {
	// Note: This is inefficient with KV (would need to scan all keys)
	// For now, we'll implement a simple version that requires scanning
	// In production, you might want to maintain a reverse index

	// This is a placeholder - KV doesn't support efficient reverse lookups
	// You'd need to either:
	// 1. Maintain a separate index (address -> []apikeys)
	// 2. Use a different storage pattern
	// 3. Accept the limitation for MVP

	return []string{}, fmt.Errorf("ListAPIKeys not implemented - use dashboard to track keys")
}
