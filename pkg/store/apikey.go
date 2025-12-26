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

	// Store in both namespaces:
	// 1. Apikey namespace: apikey -> address
	_, err := s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: constant.GetCfKvApikeyNamespaceId(),
		Key:         apiKey,
		Value:       []byte(address),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to store API key: %w", err)
	}

	// 2. Authz namespace: address -> apikey (reverse index)
	_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.authzNamespaceID,
		Key:         address,
		Value:       []byte(apiKey),
	})
	if err != nil {
		// Rollback: delete from apikey namespace if authz write fails
		s.client.DeleteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.DeleteWorkersKVEntryParams{
			NamespaceID: constant.GetCfKvApikeyNamespaceId(),
			Key:         apiKey,
		})
		return nil, fmt.Errorf("failed to store reverse index: %w", err)
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
	// First, get the address associated with this API key
	address, err := s.GetAPIKey(ctx, apiKey)
	if err != nil {
		return fmt.Errorf("failed to get address for API key: %w", err)
	}

	// Delete from both namespaces:
	// 1. Delete from apikey namespace
	_, err = s.client.DeleteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.DeleteWorkersKVEntryParams{
		NamespaceID: constant.GetCfKvApikeyNamespaceId(),
		Key:         apiKey,
	})
	if err != nil {
		return fmt.Errorf("failed to revoke API key: %w", err)
	}

	// 2. Delete from authz namespace (reverse index)
	_, err = s.client.DeleteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.DeleteWorkersKVEntryParams{
		NamespaceID: s.authzNamespaceID,
		Key:         address,
	})
	if err != nil {
		// Note: We don't rollback here since the main apikey is already deleted
		// This is acceptable for the reverse index
		return fmt.Errorf("failed to remove reverse index: %w", err)
	}

	return nil
}

// GetAPIKeyByAddress retrieves the API key associated with a wallet address (using reverse index)
func (s *KVStore) GetAPIKeyByAddress(ctx context.Context, address string) (string, error) {
	value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: s.authzNamespaceID,
		Key:         address,
	})
	if err != nil {
		return "", fmt.Errorf("API key not found for address: %w", err)
	}

	apiKey := string(value)
	if apiKey == "" {
		return "", fmt.Errorf("no API key found for address")
	}

	return apiKey, nil
}

// ListAPIKeys returns all API keys for a given wallet address (for dashboard)
func (s *KVStore) ListAPIKeys(ctx context.Context, address string) ([]string, error) {
	// With the reverse index, we can now efficiently get the API key for an address
	apiKey, err := s.GetAPIKeyByAddress(ctx, address)
	if err != nil {
		// No API key found for this address
		return []string{}, nil
	}

	// Currently, we only support one API key per address
	// In the future, this could be extended to support multiple keys
	return []string{apiKey}, nil
}
