package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cloudflare/cloudflare-go"
)

type MerkleProofData struct {
	Index  uint64   `json:"index"`
	Amount string   `json:"amount"` // BigInt as string
	Proof  []string `json:"proof"`  // Hex strings
}

// Extend Store interface
// Note: We can't easily extend the interface in worker.go without modifying it.
// So we will just add methods to KVStore and let the user update the interface definition later
// or I will update worker.go in a separate step if strict interface compliance is needed.
// For now, I'll add the methods to KVStore.

func (s *KVStore) SaveMerkleProof(ctx context.Context, address string, data *MerkleProofData) error {
	key := fmt.Sprintf("proof:%s", address)
	value, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal proof data: %w", err)
	}

	_, err = s.client.WriteWorkersKVEntry(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.WriteWorkersKVEntryParams{
		NamespaceID: s.namespaceID,
		Key:         key,
		Value:       value,
	})
	if err != nil {
		return fmt.Errorf("failed to write proof to KV: %w", err)
	}
	return nil
}

func (s *KVStore) GetMerkleProof(ctx context.Context, address string) (*MerkleProofData, error) {
	key := fmt.Sprintf("proof:%s", address)

	value, err := s.client.GetWorkersKV(ctx, cloudflare.AccountIdentifier(s.accountID), cloudflare.GetWorkersKVParams{
		NamespaceID: s.namespaceID,
		Key:         key,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get proof from KV: %w", err)
	}

	var data MerkleProofData
	if err := json.Unmarshal(value, &data); err != nil {
		return nil, fmt.Errorf("failed to unmarshal proof data: %w", err)
	}

	return &data, nil
}
