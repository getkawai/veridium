package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/y/types"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <address>")
		fmt.Println("Example: go run main.go 0x0f3e75B9Bb3efcD87B1Ed15a30C8a7FBaABD204E")
		os.Exit(1)
	}

	address := os.Args[1]

	if err := testClaimingData(address); err != nil {
		log.Fatalf("Failed to test claiming data: %v", err)
	}
}

func testClaimingData(address string) error {
	ctx := context.Background()

	log.Println("🎯 Testing Claiming Data")
	log.Println("═══════════════════════════════════════════════════════════")
	log.Printf("Address: %s", address)
	log.Println("")

	// Initialize KV store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		return fmt.Errorf("failed to initialize KV store: %w", err)
	}

	// Test 1: Get claimable rewards
	log.Println("📊 Test 1: Getting claimable rewards...")
	claimableData, err := kv.GetClaimableRewards(ctx, address)
	if err != nil {
		log.Printf("❌ No claimable rewards found: %v", err)
	} else {
		log.Printf("✅ Claimable rewards found!")

		// Print summary
		if totalKawai, ok := claimableData["total_kawai_claimable"].(string); ok && totalKawai != "0" {
			log.Printf("   💰 Total KAWAI: %s", totalKawai)
		}
		if totalUSDT, ok := claimableData["total_usdt_claimable"].(string); ok && totalUSDT != "0" {
			log.Printf("   💵 Total USDT: %s", totalUSDT)
		}

		// Print unclaimed proofs
		if unclaimedProofs, ok := claimableData["unclaimed_proofs"].([]interface{}); ok && len(unclaimedProofs) > 0 {
			log.Printf("   📋 Unclaimed proofs: %d", len(unclaimedProofs))
			for i, proof := range unclaimedProofs {
				if proofMap, ok := proof.(map[string]interface{}); ok {
					rewardType := proofMap["reward_type"]
					amount := proofMap["amount"]
					periodID := proofMap["period_id"]
					log.Printf("     %d. %s: %s (period %v)", i+1, rewardType, amount, periodID)
				}
			}
		}
	}
	log.Println("")

	// Test 2: List all Merkle proofs for this address
	log.Println("📊 Test 2: Listing all Merkle proofs...")
	proofs, err := kv.ListMerkleProofs(ctx, address)
	if err != nil {
		log.Printf("❌ Failed to list proofs: %v", err)
	} else if len(proofs) == 0 {
		log.Printf("⚠️  No Merkle proofs found for this address")
	} else {
		log.Printf("✅ Found %d Merkle proofs:", len(proofs))
		for i, proof := range proofs {
			log.Printf("   %d. Type: %s, Period: %d, Amount: %s, Index: %d",
				i+1, proof.RewardType, proof.PeriodID, proof.Amount, proof.Index)
		}
	}
	log.Println("")

	// Test 3: Check specific reward types
	log.Println("📊 Test 3: Checking specific reward types...")

	// Check mining rewards
	miningProofs := filterProofsByType(proofs, types.RewardTypeMining)
	log.Printf("   🔨 Mining rewards (KAWAI): %d proofs", len(miningProofs))

	// Check revenue sharing
	usdtProofs := filterProofsByType(proofs, types.RewardTypeRevenue)
	log.Printf("   💰 Revenue sharing (USDT): %d proofs", len(usdtProofs))

	// Check cashback
	cashbackProofs := filterProofsByType(proofs, types.RewardTypeCashback)
	log.Printf("   🎁 Cashback rewards: %d proofs", len(cashbackProofs))

	// Check referral
	referralProofs := filterProofsByType(proofs, types.RewardTypeReferral)
	log.Printf("   🤝 Referral rewards: %d proofs", len(referralProofs))

	log.Println("")

	// Test 4: Settlement periods
	log.Println("📊 Test 4: Checking settlement periods...")
	periods, err := kv.ListSettlementPeriods(ctx)
	if err != nil {
		log.Printf("❌ Failed to list settlement periods: %v", err)
	} else {
		log.Printf("✅ Found %d settlement periods:", len(periods))
		for i, period := range periods {
			log.Printf("   %d. Period %d (%s): %s total, %s root",
				i+1, period.PeriodID, period.RewardType, period.TotalAmount, period.MerkleRoot)
		}
	}

	log.Println("")
	log.Println("═══════════════════════════════════════════════════════════")

	if len(proofs) > 0 {
		log.Printf("✅ Address has %d claimable proofs ready!", len(proofs))
		log.Printf("✅ Ready for on-chain claiming (requires MON for gas)")
	} else {
		log.Printf("⚠️  Address has no claimable rewards yet")
		log.Printf("💡 Run settlement commands to generate rewards first")
	}

	return nil
}

func filterProofsByType(proofs []*store.MerkleProofData, rewardType types.RewardType) []*store.MerkleProofData {
	var filtered []*store.MerkleProofData
	for _, proof := range proofs {
		if proof.RewardType == rewardType {
			filtered = append(filtered, proof)
		}
	}
	return filtered
}
