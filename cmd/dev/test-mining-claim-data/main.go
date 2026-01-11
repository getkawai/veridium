package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kawai-network/veridium/pkg/store"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <address>")
		fmt.Println("Example: go run main.go 0x0f3e75B9Bb3efcD87B1Ed15a30C8a7FBaABD204E")
		os.Exit(1)
	}

	address := os.Args[1]

	if err := testMiningClaimData(address); err != nil {
		log.Fatalf("Failed to test mining claim data: %v", err)
	}
}

func testMiningClaimData(address string) error {
	ctx := context.Background()

	fmt.Println("🔍 Testing Mining Claim Data")
	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("Address: %s\n", address)
	fmt.Println()

	// Initialize KV store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		return fmt.Errorf("failed to initialize KV store: %w", err)
	}

	// Get all Merkle proofs for this address
	proofs, err := kv.ListMerkleProofs(ctx, address)
	if err != nil {
		return fmt.Errorf("failed to list proofs: %v", err)
	}

	if len(proofs) == 0 {
		fmt.Println("❌ No proofs found for this address")
		return nil
	}

	fmt.Printf("✅ Found %d proofs:\n", len(proofs))
	fmt.Println()

	for i, proof := range proofs {
		if proof.RewardType != "kawai" {
			continue // Skip non-mining rewards
		}

		fmt.Printf("🔨 Mining Proof %d:\n", i+1)
		fmt.Printf("   Period ID:         %d\n", proof.PeriodID)
		fmt.Printf("   Index:             %d\n", proof.Index)
		fmt.Printf("   Reward Type:       %s\n", proof.RewardType)
		fmt.Printf("   Merkle Root:       %s\n", proof.MerkleRoot)
		fmt.Println()

		// Mining-specific fields (9-field format)
		fmt.Printf("   Contributor Amount: %s KAWAI\n", proof.ContributorAmount)
		fmt.Printf("   Developer Amount:   %s KAWAI\n", proof.DeveloperAmount)
		fmt.Printf("   User Amount:        %s KAWAI\n", proof.UserAmount)
		fmt.Printf("   Affiliator Amount:  %s KAWAI\n", proof.AffiliatorAmount)
		fmt.Println()

		fmt.Printf("   Developer Address:  %s\n", proof.DeveloperAddress)
		fmt.Printf("   User Address:       %s\n", proof.UserAddress)
		fmt.Printf("   Affiliator Address: %s\n", proof.AffiliatorAddress)
		fmt.Println()

		fmt.Printf("   Merkle Proof (%d elements):\n", len(proof.Proof))
		for j, p := range proof.Proof {
			fmt.Printf("     [%d] %s\n", j, p)
		}
		fmt.Println()

		fmt.Printf("   Claim Status:      %s\n", proof.ClaimStatus)
		fmt.Printf("   Created At:        %s\n", proof.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Println()

		// Show what would be needed for ClaimMiningReward call
		fmt.Printf("📋 ClaimMiningReward Parameters:\n")
		fmt.Printf("   period:             %d\n", proof.PeriodID)
		fmt.Printf("   contributorAmount:  %s\n", proof.ContributorAmount)
		fmt.Printf("   developerAmount:    %s\n", proof.DeveloperAmount)
		fmt.Printf("   userAmount:         %s\n", proof.UserAmount)
		fmt.Printf("   affiliatorAmount:   %s\n", proof.AffiliatorAmount)
		fmt.Printf("   developerAddress:   %s\n", proof.DeveloperAddress)
		fmt.Printf("   userAddress:        %s\n", proof.UserAddress)
		fmt.Printf("   affiliatorAddress:  %s\n", proof.AffiliatorAddress)
		fmt.Printf("   proof:              [%d elements]\n", len(proof.Proof))

		fmt.Println("═══════════════════════════════════════")
	}

	return nil
}
