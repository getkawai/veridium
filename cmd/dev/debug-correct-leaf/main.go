package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kawai-network/veridium/pkg/merkle"
	"github.com/kawai-network/veridium/pkg/store"
)

func main() {
	if err := debugCorrectLeaf(); err != nil {
		log.Fatalf("Failed to debug correct leaf: %v", err)
	}
}

func debugCorrectLeaf() error {
	ctx := context.Background()

	fmt.Println("🔍 CORRECT Leaf Generation Debug")
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

	// Find the CORRECT settlement (Period 1768137242)
	var testProof *store.MerkleProofData
	if unclaimedRaw, ok := claimableData["unclaimed_proofs"]; ok {
		if unclaimedList, ok := unclaimedRaw.([]*store.MerkleProofData); ok {
			for _, proof := range unclaimedList {
				if proof.RewardType == "kawai" && proof.PeriodID == 1768137242 {
					testProof = proof
					break
				}
			}
		}
	}

	if testProof == nil {
		return fmt.Errorf("CORRECT test proof not found for period 1768137242")
	}

	fmt.Printf("📋 CORRECT Test Proof Data:\n")
	fmt.Printf("   Settlement ID:      %d\n", testProof.PeriodID)
	fmt.Printf("   Index:              %d\n", testProof.Index)
	fmt.Printf("   User Address:       %s\n", testProof.UserAddress)
	fmt.Printf("   Contributor Amount: %s\n", testProof.ContributorAmount)
	fmt.Printf("   Developer Amount:   %s\n", testProof.DeveloperAmount)
	fmt.Printf("   User Amount:        %s\n", testProof.UserAmount)
	fmt.Printf("   Affiliator Amount:  %s\n", testProof.AffiliatorAmount)
	fmt.Printf("   Developer Address:  %s\n", testProof.DeveloperAddress)
	fmt.Printf("   Affiliator Address: %s\n", testProof.AffiliatorAddress)
	fmt.Printf("   Merkle Root:        %s\n", testProof.MerkleRoot)
	fmt.Printf("   Proof Elements:     %d\n", len(testProof.Proof))
	for i, p := range testProof.Proof {
		fmt.Printf("     [%d]: %s\n", i, p)
	}
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

	// Contract period (should be 7)
	contractPeriod := uint64(7)

	fmt.Printf("🧪 Contract Leaf Generation (What contract expects):\n")
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
		contractPeriod,
		msgSender,
		contribAmt, devAmt, userAmt, affAmt,
		developer, user, affiliator,
	)

	fmt.Printf("🌿 Contract Expected Leaf: 0x%x\n", contractLeaf)
	fmt.Println()

	// Now let's recreate the Merkle tree from the settlement to see what was actually generated
	fmt.Printf("🌳 Recreating Settlement Merkle Tree...\n")

	// Test data from create-correct-mining-settlement
	testRewards := []struct {
		contributorAddr   string
		contributorAmount *big.Int
		developerAmount   *big.Int
		userAmount        *big.Int
		affiliatorAmount  *big.Int
		developerAddr     string
		userAddr          string
		affiliatorAddr    string
	}{
		// Entry 1: Our test user as contributor
		{
			contributorAddr:   testAddress,
			contributorAmount: new(big.Int).Mul(big.NewInt(126), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), // 126 KAWAI
			developerAmount:   new(big.Int).Mul(big.NewInt(6), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 6 KAWAI
			userAmount:        new(big.Int).Mul(big.NewInt(6), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 6 KAWAI
			affiliatorAmount:  new(big.Int).Mul(big.NewInt(6), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 6 KAWAI
			developerAddr:     "0x94D5C06229811c4816107005ff05259f229Eb07b",
			userAddr:          testAddress,
			affiliatorAddr:    "0x2864Cd9a59f32b74f3f851B92973fD40883aD503",
		},
		// Entry 2: Another contributor
		{
			contributorAddr:   "0x1111111111111111111111111111111111111111",
			contributorAmount: new(big.Int).Mul(big.NewInt(100), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), // 100 KAWAI
			developerAmount:   new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 5 KAWAI
			userAmount:        new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 5 KAWAI
			affiliatorAmount:  new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 5 KAWAI
			developerAddr:     "0x94D5C06229811c4816107005ff05259f229Eb07b",
			userAddr:          "0x2222222222222222222222222222222222222222",
			affiliatorAddr:    "0x2864Cd9a59f32b74f3f851B92973fD40883aD503",
		},
		// Entry 3: Third contributor
		{
			contributorAddr:   "0x3333333333333333333333333333333333333333",
			contributorAmount: new(big.Int).Mul(big.NewInt(80), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), // 80 KAWAI
			developerAmount:   new(big.Int).Mul(big.NewInt(4), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),  // 4 KAWAI
			userAmount:        new(big.Int).Mul(big.NewInt(4), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),  // 4 KAWAI
			affiliatorAmount:  new(big.Int).Mul(big.NewInt(4), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),  // 4 KAWAI
			developerAddr:     "0x94D5C06229811c4816107005ff05259f229Eb07b",
			userAddr:          "0x4444444444444444444444444444444444444444",
			affiliatorAddr:    "0x2864Cd9a59f32b74f3f851B92973fD40883aD503",
		},
	}

	// Generate leaves using CONTRACT PERIOD
	var leaves [][]byte
	for i, reward := range testRewards {
		leaf := generateContractLeaf(
			contractPeriod, // Use contract period (7)
			common.HexToAddress(reward.contributorAddr),
			reward.contributorAmount,
			reward.developerAmount,
			reward.userAmount,
			reward.affiliatorAmount,
			common.HexToAddress(reward.developerAddr),
			common.HexToAddress(reward.userAddr),
			common.HexToAddress(reward.affiliatorAddr),
		)
		leaves = append(leaves, leaf)

		fmt.Printf("   Leaf %d (%s): 0x%x\n", i, reward.contributorAddr, leaf)
	}

	// Build Merkle tree
	tree := merkle.NewMerkleTree(leaves)
	root := tree.Root

	fmt.Printf("   Generated Root: 0x%x\n", root)
	fmt.Printf("   Stored Root:    %s\n", testProof.MerkleRoot)
	fmt.Println()

	// Check if roots match
	if fmt.Sprintf("0x%x", root) == testProof.MerkleRoot {
		fmt.Printf("✅ Merkle roots match!\n")
	} else {
		fmt.Printf("❌ Merkle root mismatch!\n")
		fmt.Printf("   Expected: %s\n", testProof.MerkleRoot)
		fmt.Printf("   Got:      0x%x\n", root)
	}

	// Generate proof for our test leaf (index 0)
	testLeaf := leaves[0] // Our test user is index 0
	proof, ok := tree.GetProof(testLeaf)
	if !ok {
		return fmt.Errorf("failed to generate proof for test leaf")
	}

	fmt.Printf("🔍 Generated Proof for Test User:\n")
	for i, p := range proof {
		fmt.Printf("   [%d]: 0x%x\n", i, p)
	}
	fmt.Println()

	// Compare with stored proof
	fmt.Printf("🔍 Stored Proof:\n")
	for i, p := range testProof.Proof {
		fmt.Printf("   [%d]: %s\n", i, p)
	}
	fmt.Println()

	// Verify proof manually
	fmt.Printf("🧪 Manual Proof Verification:\n")
	// Manual verification using tree
	currentHash := testLeaf
	for _, proofElement := range proof {
		// Determine order based on hash comparison
		if fmt.Sprintf("%x", currentHash) < fmt.Sprintf("%x", proofElement) {
			currentHash = crypto.Keccak256(currentHash, proofElement)
		} else {
			currentHash = crypto.Keccak256(proofElement, currentHash)
		}
	}
	isValid := fmt.Sprintf("%x", currentHash) == fmt.Sprintf("%x", root)
	fmt.Printf("   Proof Valid: %t\n", isValid)

	return nil
}

// generateContractLeaf creates a 9-field Merkle leaf exactly as the contract does
func generateContractLeaf(
	period uint64,
	msgSender common.Address,
	contributorAmt, developerAmt, userAmt, affiliatorAmt *big.Int,
	developer, user, affiliator common.Address,
) []byte {
	// Exact match to contract Solidity code:
	// keccak256(abi.encodePacked(period, msg.sender, contributorAmount, developerAmount, userAmount, affiliatorAmount, developer, user, affiliator))
	return crypto.Keccak256(
		common.LeftPadBytes(big.NewInt(int64(period)).Bytes(), 32), // uint256 period
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
