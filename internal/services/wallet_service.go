package services

import (
	"github.com/kawai-network/x/jarvis"
	"github.com/kawai-network/x/store"
)

// WalletService is an alias to the shared implementation in x/jarvis.
type WalletService = jarvis.WalletService

// NewWalletService creates a wallet service backed by x/jarvis.
func NewWalletService(dataDir string, kvStore *store.KVStore) *WalletService {
	return jarvis.NewWalletService(dataDir, kvStore)
}
