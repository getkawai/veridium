package services

import (
	"context"

	"github.com/kawai-network/veridium/pkg/store"
)

type APIKeyService struct {
	kvStore *store.KVStore
}

func NewAPIKeyService(kvStore *store.KVStore) *APIKeyService {
	return &APIKeyService{
		kvStore: kvStore,
	}
}

// CreateKey generates a new API key for the given wallet address
func (s *APIKeyService) CreateKey(address string) (*store.APIKey, error) {
	return s.kvStore.CreateAPIKey(context.Background(), address)
}

// ValidateKey checks if an API key exists and returns the associated wallet address
func (s *APIKeyService) ValidateKey(apiKey string) (string, error) {
	return s.kvStore.ValidateAPIKey(context.Background(), apiKey)
}

// RevokeKey deletes an API key
func (s *APIKeyService) RevokeKey(apiKey string) error {
	return s.kvStore.RevokeAPIKey(context.Background(), apiKey)
}

// GetKey retrieves the wallet address associated with an API key
func (s *APIKeyService) GetKey(apiKey string) (string, error) {
	return s.kvStore.GetAPIKey(context.Background(), apiKey)
}

// GetKeyByAddress retrieves the API key associated with a wallet address (reverse lookup)
func (s *APIKeyService) GetKeyByAddress(address string) (string, error) {
	return s.kvStore.GetAPIKeyByAddress(context.Background(), address)
}

// ListKeys returns all API keys for a given wallet address
func (s *APIKeyService) ListKeys(address string) ([]string, error) {
	return s.kvStore.ListAPIKeys(context.Background(), address)
}
