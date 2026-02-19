package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kawai-network/veridium/pkg/store"
)

// HolderRegistry manages the registry of KAWAI token holders
type HolderRegistry struct {
	kvStore *store.KVStore
}

// NewHolderRegistry creates a new holder registry
func NewHolderRegistry(kvStore *store.KVStore) *HolderRegistry {
	return &HolderRegistry{
		kvStore: kvStore,
	}
}

// RegisterHolder registers a KAWAI holder address in the registry
// This is called automatically when:
// - User connects wallet (desktop app)
// - Contributor claims rewards (CLI)
// - User/contributor receives KAWAI transfer
func (hr *HolderRegistry) RegisterHolder(ctx context.Context, address common.Address, source string) error {
	addressHex := address.Hex()

	// Check if already registered
	existing, err := hr.kvStore.GetHolder(ctx, addressHex)
	if err == nil {
		// Already registered, just update lastSeen
		existing.LastSeen = time.Now().Unix()
		existing.Source = source // Update source if changed
		return hr.kvStore.SaveHolder(ctx, addressHex, existing)
	}

	// Check if error is something other than "not found"
	// If it's a real error (network, API, etc), we should not proceed
	if err != nil && !isNotFoundError(err) {
		return fmt.Errorf("failed to check existing holder: %w", err)
	}

	// Not found or acceptable error - register new holder
	holderInfo := &store.HolderInfo{
		Address:    addressHex,
		LastSeen:   time.Now().Unix(),
		Source:     source,
		Registered: time.Now().Unix(),
	}

	if err := hr.kvStore.SaveHolder(ctx, addressHex, holderInfo); err != nil {
		return fmt.Errorf("failed to register holder: %w", err)
	}

	return nil
}

// isNotFoundError checks if an error indicates a key was not found
// This helps distinguish between "not found" (expected) and real errors (network, API, etc)
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "key not found") ||
		strings.Contains(errMsg, "not found") ||
		strings.Contains(errMsg, "no such key") ||
		strings.Contains(errMsg, "does not exist")
}

// GetAllHolders returns all registered holder addresses from the registry
func (hr *HolderRegistry) GetAllHolders(ctx context.Context) ([]common.Address, error) {
	holders, err := hr.kvStore.ListHolders(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list holders: %w", err)
	}

	addresses := make([]common.Address, 0, len(holders))
	for _, holder := range holders {
		addresses = append(addresses, common.HexToAddress(holder.Address))
	}

	log.Printf("📊 [HOLDER REGISTRY] Found %d registered holders", len(addresses))
	return addresses, nil
}

// GetHolderInfo returns information about a specific holder
func (hr *HolderRegistry) GetHolderInfo(ctx context.Context, address common.Address) (*store.HolderInfo, error) {
	return hr.kvStore.GetHolder(ctx, address.Hex())
}

// GetHolderCount returns the total number of registered holders
func (hr *HolderRegistry) GetHolderCount(ctx context.Context) (int, error) {
	return hr.kvStore.GetHolderCount(ctx)
}

// RemoveHolder removes a holder from the registry (for cleanup/testing)
func (hr *HolderRegistry) RemoveHolder(ctx context.Context, address common.Address) error {
	return hr.kvStore.DeleteHolder(ctx, address.Hex())
}

// ExportHolders exports all holder addresses as JSON (for backup/debugging)
func (hr *HolderRegistry) ExportHolders(ctx context.Context) (string, error) {
	holders, err := hr.kvStore.ListHolders(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to list holders: %w", err)
	}

	data, err := json.MarshalIndent(holders, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal holders: %w", err)
	}

	return string(data), nil
}
