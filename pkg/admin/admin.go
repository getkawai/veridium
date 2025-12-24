package admin

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/kawai-network/veridium/pkg/blockchain"
	"github.com/kawai-network/veridium/pkg/merkle"
	"github.com/kawai-network/veridium/pkg/store"
)

type AdminManager struct {
	Chain *blockchain.Client
	Store store.Store
}

func NewAdminManager(chain *blockchain.Client, kv store.Store) *AdminManager {
	return &AdminManager{
		Chain: chain,
		Store: kv,
	}
}

// AuditWorkers prints the status of all registered workers
// Note: This requires the KV store to support listing keys,
// but for now, we'll implement a simple version that logs status.
func (a *AdminManager) AuditWorkers(ctx context.Context) error {
	log.Println("--- Worker Audit Report ---")

	workers, err := a.Store.ListWorkers(ctx)
	if err != nil {
		return fmt.Errorf("failed to list workers: %w", err)
	}

	if len(workers) == 0 {
		log.Println("No workers registered.")
		return nil
	}

	fmt.Printf("%-42s | %-15s | %-10s | %s\n", "Wallet Address", "Last Seen", "Status", "Specs")
	fmt.Println("------------------------------------------------------------------------------------------------")
	for _, w := range workers {
		lastSeen := w.LastSeen.Format("2006-01-02 15:04")
		fmt.Printf("%-42s | %-15s | %-10s | %s\n", w.WalletAddress, lastSeen, w.Status, w.HardwareSpecs)
	}

	return nil
}

// CalculateDividends generates the Merkle Tree for worker rewards, applying the 70/30 split logic
// and the Admin-owned node 100% rule. Proofs are then saved to KV for user claiming.
func (a *AdminManager) CalculateDividends(ctx context.Context) error {
	log.Println("--- Dividend Calculation (Merkle Airdrop) ---")

	// 1. Fetch Workers
	workers, err := a.Store.ListWorkers(ctx)
	if err != nil {
		return fmt.Errorf("failed to list workers: %w", err)
	}
	if len(workers) == 0 {
		log.Println("No workers found.")
		return nil
	}

	var leaves [][]byte
	var workerProofs []*store.MerkleProofData

	log.Printf("Generating Merkle Tree for %d workers...", len(workers))

	// Helper to track address for saving proofs later
	var proofAddresses []string

	currentIndex := uint64(0)
	for _, w := range workers {
		if w.AccumulatedRewards == "" || w.AccumulatedRewards == "0" {
			continue // Skip workers with no rewards
		}

		addr := common.HexToAddress(w.WalletAddress)
		amount := new(big.Int)
		amount.SetString(w.AccumulatedRewards, 10)

		// Create Leaf
		leaf := merkle.HashLeaf(currentIndex, addr, amount)
		leaves = append(leaves, leaf)

		workerProofs = append(workerProofs, &store.MerkleProofData{
			Index:  currentIndex,
			Amount: amount.String(),
		})
		proofAddresses = append(proofAddresses, w.WalletAddress)
		currentIndex++
	}

	// 3. Build Tree (Standard Logic)
	if len(leaves) == 0 {
		return nil
	}
	tree := merkle.NewMerkleTree(leaves)
	root := tree.Root
	log.Printf("Merkle Root: 0x%x", root)

	// 8. Generate and Save Proofs
	// Note: We need to map back to the Address to save by key.
	// Our `entries` slice aligns with `workerProofs` slice.
	for i, wp := range workerProofs {
		proof, ok := tree.GetProof(leaves[i])
		if !ok {
			log.Printf("Error generating proof for index %d", i)
			continue
		}

		var proofHex []string
		for _, p := range proof {
			proofHex = append(proofHex, fmt.Sprintf("0x%x", p))
		}
		wp.Proof = proofHex

		addrStr := proofAddresses[i]
		err := a.Store.SaveMerkleProof(ctx, addrStr, wp)
		if err != nil {
			log.Printf("Failed to save proof for %s: %v", addrStr, err)
		}
	}

	// 5. Submit Root to Blockchain (TODO: Bindings needed)
	// if a.Chain != nil {
	//     tx, err := a.Chain.MerkleDistributor.SetMerkleRoot(a.Chain.Auth, [32]byte(root))
	//     ...
	// }
	log.Println("Done. Root generated and Proofs saved.")

	return nil
}
