package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/y/types"
)

func main() {
	if err := inspectKVData(); err != nil {
		log.Fatalf("Failed to inspect KV data: %v", err)
	}
}

func inspectKVData() error {
	ctx := context.Background()

	fmt.Println("🔍 KV Data Inspection")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// Test address
	testAddress := "0x0f3e75B9Bb3efcD87B1Ed15a30C8a7FBaABD204E"

	// Initialize KV store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		return fmt.Errorf("failed to initialize KV store: %w", err)
	}

	// Get claimable rewards
	claimableData, err := kv.GetClaimableRewards(ctx, testAddress)
	if err != nil {
		return fmt.Errorf("failed to get claimable rewards: %w", err)
	}

	fmt.Printf("📋 Claimable Data for %s:\n", testAddress)
	fmt.Printf("   Total KAWAI: %s\n", getStringFromMap(claimableData, "total_kawai_claimable", "0"))
	fmt.Printf("   Total USDT:  %s\n", getStringFromMap(claimableData, "total_usdt_claimable", "0"))
	fmt.Println()

	// Inspect all unclaimed proofs
	if unclaimedRaw, ok := claimableData["unclaimed_proofs"]; ok {
		if unclaimedList, ok := unclaimedRaw.([]*store.MerkleProofData); ok {
			fmt.Printf("📋 Unclaimed Proofs (%d total):\n", len(unclaimedList))
			for i, proof := range unclaimedList {
				fmt.Printf("\n   Proof %d:\n", i+1)
				fmt.Printf("     Period ID:         %d\n", proof.PeriodID)
				fmt.Printf("     Index:             %d\n", proof.Index)
				fmt.Printf("     Reward Type:       %s\n", proof.RewardType)
				fmt.Printf("     Amount:            %s\n", proof.Amount)
				fmt.Printf("     Contributor Amount: %s\n", proof.ContributorAmount)
				fmt.Printf("     Developer Amount:   %s\n", proof.DeveloperAmount)
				fmt.Printf("     User Amount:        %s\n", proof.UserAmount)
				fmt.Printf("     Affiliator Amount:  %s\n", proof.AffiliatorAmount)
				fmt.Printf("     Developer Address:  %s\n", proof.DeveloperAddress)
				fmt.Printf("     User Address:       %s\n", proof.UserAddress)
				fmt.Printf("     Affiliator Address: %s\n", proof.AffiliatorAddress)
				fmt.Printf("     Merkle Root:        %s\n", proof.MerkleRoot)
				fmt.Printf("     Proof Elements (%d):\n", len(proof.Proof))
				for j, p := range proof.Proof {
					fmt.Printf("       [%d]: %s\n", j, p)
				}
				fmt.Printf("     Claim Status:       %s\n", proof.ClaimStatus)
				fmt.Printf("     Created At:         %s\n", proof.CreatedAt.Format("2006-01-02 15:04:05"))
			}
		}
	}

	// Look specifically for the CORRECT settlement (1768137242)
	fmt.Println()
	fmt.Printf("🎯 Looking for CORRECT settlement (1768137242):\n")

	var correctProof *store.MerkleProofData
	if unclaimedRaw, ok := claimableData["unclaimed_proofs"]; ok {
		if unclaimedList, ok := unclaimedRaw.([]*store.MerkleProofData); ok {
			for _, proof := range unclaimedList {
				if proof.RewardType == types.RewardTypeMining && proof.PeriodID == 1768137242 {
					correctProof = proof
					break
				}
			}
		}
	}

	if correctProof == nil {
		fmt.Printf("❌ CORRECT settlement proof not found!\n")
		return fmt.Errorf("CORRECT settlement proof not found")
	}

	fmt.Printf("✅ Found CORRECT settlement proof:\n")
	fmt.Printf("   Period ID: %d\n", correctProof.PeriodID)
	fmt.Printf("   Index:     %d\n", correctProof.Index)
	fmt.Printf("   Root:      %s\n", correctProof.MerkleRoot)
	fmt.Printf("   Proof:     %v\n", correctProof.Proof)

	return nil
}

// getStringFromMap safely extracts a string from map[string]interface{}
func getStringFromMap(m map[string]interface{}, key string, defaultVal string) string {
	if val, ok := m[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return defaultVal
}
