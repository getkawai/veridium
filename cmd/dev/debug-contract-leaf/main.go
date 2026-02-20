package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/y/types"
)

func main() {
	if err := debugContractLeaf(); err != nil {
		log.Fatalf("Failed to debug contract leaf: %v", err)
	}
}

func debugContractLeaf() error {
	ctx := context.Background()

	fmt.Println("🔍 Contract Leaf Generation Debug")
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

	// Find the proper test settlement (Period 1768136095)
	var testProof *store.MerkleProofData
	if unclaimedRaw, ok := claimableData["unclaimed_proofs"]; ok {
		if unclaimedList, ok := unclaimedRaw.([]*store.MerkleProofData); ok {
			for _, proof := range unclaimedList {
				if proof.RewardType == types.RewardTypeMining && proof.PeriodID == 1768136095 {
					testProof = proof
					break
				}
			}
		}
	}

	if testProof == nil {
		return fmt.Errorf("test proof not found for period 1768136095")
	}

	fmt.Printf("📋 Test Proof Data:\n")
	fmt.Printf("   Period ID:         %d\n", testProof.PeriodID)
	fmt.Printf("   User Address:      %s\n", testProof.UserAddress)
	fmt.Printf("   Contributor Amount: %s\n", testProof.ContributorAmount)
	fmt.Printf("   Developer Amount:   %s\n", testProof.DeveloperAmount)
	fmt.Printf("   User Amount:        %s\n", testProof.UserAmount)
	fmt.Printf("   Affiliator Amount:  %s\n", testProof.AffiliatorAmount)
	fmt.Printf("   Developer Address:  %s\n", testProof.DeveloperAddress)
	fmt.Printf("   Affiliator Address: %s\n", testProof.AffiliatorAddress)
	fmt.Println()

	// Map to contract period
	contractPeriod, err := mapSettlementPeriodToContractPeriod(testProof.PeriodID)
	if err != nil {
		return fmt.Errorf("failed to map period: %w", err)
	}

	fmt.Printf("🔄 Period Mapping: %d -> %d\n", testProof.PeriodID, contractPeriod)
	fmt.Println()

	// Parse amounts
	contribAmt := new(big.Int)
	contribAmt.SetString(testProof.ContributorAmount, 10)

	devAmt := new(big.Int)
	devAmt.SetString(testProof.DeveloperAmount, 10)

	userAmt := new(big.Int)
	userAmt.SetString(testProof.UserAmount, 10)

	affAmt := new(big.Int)
	affAmt.SetString(testProof.AffiliatorAmount, 10)

	// Parse addresses
	msgSender := common.HexToAddress(testAddress) // This is msg.sender in contract
	developer := common.HexToAddress(testProof.DeveloperAddress)
	user := common.HexToAddress(testProof.UserAddress)
	affiliator := common.HexToAddress(testProof.AffiliatorAddress)

	fmt.Printf("🧪 Contract Leaf Generation (Exact Match):\n")
	fmt.Printf("   period:             %d\n", contractPeriod)
	fmt.Printf("   msg.sender:         %s\n", msgSender.Hex())
	fmt.Printf("   contributorAmount:  %s\n", contribAmt.String())
	fmt.Printf("   developerAmount:    %s\n", devAmt.String())
	fmt.Printf("   userAmount:         %s\n", userAmt.String())
	fmt.Printf("   affiliatorAmount:   %s\n", affAmt.String())
	fmt.Printf("   developer:          %s\n", developer.Hex())
	fmt.Printf("   user:               %s\n", user.Hex())
	fmt.Printf("   affiliator:         %s\n", affiliator.Hex())
	fmt.Println()

	// Generate leaf exactly as contract does
	contractLeaf := generateContractLeaf(
		uint64(contractPeriod), // Use contract period, not settlement period!
		msgSender,
		contribAmt, devAmt, userAmt, affAmt,
		developer, user, affiliator,
	)

	fmt.Printf("🌿 Contract Leaf Hash: 0x%x\n", contractLeaf)
	fmt.Println()

	// Also generate with settlement period for comparison
	settlementLeaf := generateContractLeaf(
		uint64(testProof.PeriodID), // Use settlement period
		msgSender,
		contribAmt, devAmt, userAmt, affAmt,
		developer, user, affiliator,
	)

	fmt.Printf("🌿 Settlement Leaf Hash: 0x%x\n", settlementLeaf)
	fmt.Println()

	// Compare with stored proof
	fmt.Printf("🔍 Comparison:\n")
	fmt.Printf("   Stored Merkle Root: %s\n", testProof.MerkleRoot)
	fmt.Printf("   Contract uses period: %d (mapped from %d)\n", contractPeriod, testProof.PeriodID)
	fmt.Printf("   Settlement used period: %d\n", testProof.PeriodID)
	fmt.Println()

	if contractPeriod != testProof.PeriodID {
		fmt.Printf("❌ PERIOD MISMATCH DETECTED!\n")
		fmt.Printf("   Contract expects period: %d\n", contractPeriod)
		fmt.Printf("   Settlement generated with period: %d\n", testProof.PeriodID)
		fmt.Printf("   This is why the proof fails!\n")
		fmt.Println()

		fmt.Printf("💡 Solution: Generate settlement using contract period numbers\n")
		fmt.Printf("   Instead of timestamp-based periods, use sequential 1,2,3,4,5,6\n")
	} else {
		fmt.Printf("✅ Period numbers match\n")
	}

	return nil
}

// generateContractLeaf generates leaf exactly as the contract does
func generateContractLeaf(
	period uint64,
	msgSender common.Address,
	contributorAmt, developerAmt, userAmt, affiliatorAmt *big.Int,
	developer, user, affiliator common.Address,
) []byte {
	// Exact match to contract Solidity code:
	// keccak256(abi.encodePacked(period, msg.sender, contributorAmount, developerAmount, userAmount, affiliatorAmount, developer, user, affiliator))
	return crypto.Keccak256(
		common.LeftPadBytes(new(big.Int).SetUint64(period).Bytes(), 32), // uint256 period
		msgSender.Bytes(), // address msg.sender (20 bytes)
		common.LeftPadBytes(contributorAmt.Bytes(), 32), // uint256 contributorAmount
		common.LeftPadBytes(developerAmt.Bytes(), 32),   // uint256 developerAmount
		common.LeftPadBytes(userAmt.Bytes(), 32),        // uint256 userAmount
		common.LeftPadBytes(affiliatorAmt.Bytes(), 32),  // uint256 affiliatorAmount
		developer.Bytes(),  // address developer (20 bytes)
		user.Bytes(),       // address user (20 bytes)
		affiliator.Bytes(), // address affiliator (20 bytes)
	)
}

// mapSettlementPeriodToContractPeriod maps settlement period IDs to sequential contract periods
func mapSettlementPeriodToContractPeriod(settlementPeriodID int64) (int64, error) {
	periodMapping := map[int64]int64{
		1767549424: 1, // Oldest settlement -> Contract period 1
		1767557168: 2, // Second oldest -> Contract period 2
		1767650263: 3, // Third oldest -> Contract period 3
		1768130418: 4, // Newest settlement -> Contract period 4
		1768135359: 5, // Test settlement with correct addresses -> Contract period 5
		1768136095: 6, // Proper test settlement with valid proofs -> Contract period 6
	}

	contractPeriod, exists := periodMapping[settlementPeriodID]
	if !exists {
		return 0, fmt.Errorf("unknown settlement period ID: %d", settlementPeriodID)
	}

	return contractPeriod, nil
}
