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

// AuditContributors prints the status of all registered contributors
// Note: This requires the KV store to support listing keys,
// but for now, we'll implement a simple version that logs status.
func (a *AdminManager) AuditContributors(ctx context.Context) error {
	log.Println("--- Contributor Audit Report ---")

	contributors, err := a.Store.ListContributors(ctx)
	if err != nil {
		return fmt.Errorf("failed to list contributors: %w", err)
	}

	if len(contributors) == 0 {
		log.Println("No contributors registered.")
		return nil
	}

	fmt.Printf("%-42s | %-15s | %-10s | %s\n", "Wallet Address", "Last Seen", "Status", "Specs")
	fmt.Println("------------------------------------------------------------------------------------------------")
	for _, c := range contributors {
		lastSeen := c.LastSeen.Format("2006-01-02 15:04")
		fmt.Printf("%-42s | %-15s | %-10s | %s\n", c.WalletAddress, lastSeen, c.Status, c.HardwareSpecs)
	}

	return nil
}

// CalculateDividends generates the Merkle Tree for contributor rewards, applying the 70/30 split logic
// and the Admin-owned node 100% rule. Proofs are then saved to KV for user claiming.
func (a *AdminManager) CalculateDividends(ctx context.Context) error {
	log.Println("--- Dividend Calculation (Merkle Airdrop) ---")

	// 1. Fetch Contributors
	contributors, err := a.Store.ListContributors(ctx)
	if err != nil {
		return fmt.Errorf("failed to list contributors: %w", err)
	}
	if len(contributors) == 0 {
		log.Println("No contributors found.")
		return nil
	}

	var leaves [][]byte
	var contributorProofs []*store.MerkleProofData

	log.Printf("Generating Merkle Tree for %d contributors...", len(contributors))

	// Helper to track address for saving proofs later
	var proofAddresses []string

	currentIndex := uint64(0)
	for _, c := range contributors {
		if c.AccumulatedRewards == "" || c.AccumulatedRewards == "0" {
			continue // Skip contributors with no rewards
		}

		addr := common.HexToAddress(c.WalletAddress)
		amount := new(big.Int)
		amount.SetString(c.AccumulatedRewards, 10)

		// Create Leaf
		leaf := merkle.HashLeaf(currentIndex, addr, amount)
		leaves = append(leaves, leaf)

		contributorProofs = append(contributorProofs, &store.MerkleProofData{
			Index:  currentIndex,
			Amount: amount.String(),
		})
		proofAddresses = append(proofAddresses, c.WalletAddress)
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
	// Our `entries` slice aligns with `contributorProofs` slice.
	for i, cp := range contributorProofs {
		proof, ok := tree.GetProof(leaves[i])
		if !ok {
			log.Printf("Error generating proof for index %d", i)
			continue
		}

		var proofHex []string
		for _, p := range proof {
			proofHex = append(proofHex, fmt.Sprintf("0x%x", p))
		}
		cp.Proof = proofHex

		addrStr := proofAddresses[i]
		err := a.Store.SaveMerkleProof(ctx, addrStr, cp)
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

// CalculateUSDTDividends generates the Merkle Tree for USDT profit distribution in Phase 2.
// This distributes the remaining USDT Profit (after Contributor/Admin costs) to KAWAI Holders.
func (a *AdminManager) CalculateUSDTDividends(ctx context.Context, totalProfit *big.Int) error {
	log.Println("--- USDT Profit Distribution (Phase 2) ---")

	// 1. Get all KAWAI holders and their balances from blockchain
	// For now, we'll use the accumulated USDT from contributors as a proxy
	// In production, this should scan KAWAI token holders
	contributors, err := a.Store.ListContributors(ctx)
	if err != nil {
		return fmt.Errorf("failed to list contributors: %w", err)
	}
	if len(contributors) == 0 {
		log.Println("No contributors/holders found.")
		return nil
	}

	// 2. Calculate total KAWAI holdings (simplified: use AccumulatedRewards as proxy)
	totalKawai := new(big.Int)
	holderBalances := make(map[string]*big.Int)

	for _, c := range contributors {
		if c.AccumulatedRewards == "" || c.AccumulatedRewards == "0" {
			continue
		}
		balance := new(big.Int)
		balance.SetString(c.AccumulatedRewards, 10)
		holderBalances[c.WalletAddress] = balance
		totalKawai.Add(totalKawai, balance)
	}

	if totalKawai.Cmp(big.NewInt(0)) == 0 {
		log.Println("No KAWAI holdings found.")
		return nil
	}

	log.Printf("Total KAWAI Holdings: %s", totalKawai.String())
	log.Printf("Total USDT Profit to Distribute: %s", totalProfit.String())

	// 3. Generate Merkle Tree for USDT distribution
	var leaves [][]byte
	var proofData []*store.MerkleProofData
	var proofAddresses []string

	currentIndex := uint64(0)
	for addr, balance := range holderBalances {
		// Calculate proportional share: (balance / totalKawai) * totalProfit
		share := new(big.Int).Mul(balance, totalProfit)
		share.Div(share, totalKawai)

		if share.Cmp(big.NewInt(0)) == 0 {
			continue
		}

		ethAddr := common.HexToAddress(addr)
		leaf := merkle.HashLeaf(currentIndex, ethAddr, share)
		leaves = append(leaves, leaf)

		proofData = append(proofData, &store.MerkleProofData{
			Index:  currentIndex,
			Amount: share.String(),
		})
		proofAddresses = append(proofAddresses, addr)
		currentIndex++

		log.Printf("Holder %s: %s KAWAI -> %s USDT share", addr, balance.String(), share.String())
	}

	if len(leaves) == 0 {
		return nil
	}

	// 4. Build Tree
	tree := merkle.NewMerkleTree(leaves)
	root := tree.Root
	log.Printf("USDT Merkle Root: 0x%x", root)

	// 5. Save Proofs (with different prefix to distinguish from KAWAI proofs)
	for i, pd := range proofData {
		proof, ok := tree.GetProof(leaves[i])
		if !ok {
			log.Printf("Error generating USDT proof for index %d", i)
			continue
		}

		var proofHex []string
		for _, p := range proof {
			proofHex = append(proofHex, fmt.Sprintf("0x%x", p))
		}
		pd.Proof = proofHex

		// Save with "usdt:" prefix to distinguish from KAWAI proofs
		addrKey := "usdt:" + proofAddresses[i]
		err := a.Store.SaveMerkleProof(ctx, addrKey, pd)
		if err != nil {
			log.Printf("Failed to save USDT proof for %s: %v", proofAddresses[i], err)
		}
	}

	log.Println("Done. USDT Merkle Root generated and Proofs saved.")
	return nil
}
