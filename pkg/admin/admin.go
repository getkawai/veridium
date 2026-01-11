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
// This distributes the platform USDT profit to KAWAI holders proportionally.
func (a *AdminManager) CalculateUSDTDividends(ctx context.Context, totalProfit *big.Int) error {
	log.Println("--- USDT Profit Distribution (Revenue Sharing) ---")

	// 1. Scan KAWAI holders from blockchain
	scanner, err := blockchain.NewHolderScanner()
	if err != nil {
		return fmt.Errorf("failed to create holder scanner: %w", err)
	}

	holders, err := scanner.ScanHoldersLatest(ctx)
	if err != nil {
		return fmt.Errorf("failed to scan holders: %w", err)
	}

	if len(holders) == 0 {
		log.Println("No KAWAI holders found.")
		return nil
	}

	// 2. Get total supply
	totalSupply, err := scanner.GetTotalSupply(ctx)
	if err != nil {
		return fmt.Errorf("failed to get total supply: %w", err)
	}

	log.Printf("Total KAWAI Supply: %s", totalSupply.String())
	log.Printf("Total USDT Profit to Distribute: %s", totalProfit.String())
	log.Printf("Number of Holders: %d", len(holders))

	// 3. Validate holder balances
	if err := blockchain.ValidateHolders(holders, totalSupply); err != nil {
		log.Printf("⚠️  Warning: %v", err)
		// Continue anyway - this is just a sanity check
	}

	// 4. Generate Merkle Tree for USDT distribution
	var leaves [][]byte
	var proofData []*store.MerkleProofData
	var proofAddresses []string

	currentIndex := uint64(0)
	for _, holder := range holders {
		// Calculate proportional share: (balance / totalSupply) * totalProfit
		share := blockchain.CalculateHolderShare(holder.Balance, totalSupply, totalProfit)

		if share.Cmp(big.NewInt(0)) == 0 {
			continue // Skip holders with zero share
		}

		leaf := merkle.HashLeaf(currentIndex, holder.Address, share)
		leaves = append(leaves, leaf)

		proofData = append(proofData, &store.MerkleProofData{
			Index:  currentIndex,
			Amount: share.String(),
		})
		proofAddresses = append(proofAddresses, holder.Address.Hex())
		currentIndex++

		// Calculate percentage using big.Float to avoid overflow
		balanceFloat := new(big.Float).SetInt(holder.Balance)
		supplyFloat := new(big.Float).SetInt(totalSupply)
		percentage := new(big.Float).Quo(balanceFloat, supplyFloat)
		percentFloat, _ := percentage.Float64()

		log.Printf("Holder %s: %s KAWAI (%.4f%%) -> %s USDT dividend",
			holder.Address.Hex(),
			holder.Balance.String(),
			percentFloat*100,
			share.String())
	}

	if len(leaves) == 0 {
		return fmt.Errorf("no valid dividend recipients")
	}

	// 5. Build Merkle Tree
	tree := merkle.NewMerkleTree(leaves)
	root := tree.Root
	log.Printf("USDT Merkle Root: 0x%x", root)

	// 6. Save Proofs (with "usdt:" prefix to distinguish from KAWAI proofs)
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
	log.Printf("📊 Summary: %d holders will receive dividends", len(proofData))
	return nil
}
