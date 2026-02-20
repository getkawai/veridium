package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kawai-network/veridium/pkg/merkle"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/y/types"
	"github.com/kawai-network/x/constant"
)

func main() {
	if err := createProperTestMiningSettlement(); err != nil {
		log.Fatalf("Failed to create proper test mining settlement: %v", err)
	}
}

func createProperTestMiningSettlement() error {
	ctx := context.Background()

	fmt.Println("🧪 Creating Proper Test Mining Settlement (Multiple Leaves)")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()

	// Initialize KV store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		return fmt.Errorf("failed to initialize KV store: %w", err)
	}

	// Test addresses
	testUserAddr := "0x0f3e75B9Bb3efcD87B1Ed15a30C8a7FBaABD204E"     // The actual claiming address
	testDeveloperAddr := constant.GetAdminAddress()                  // Treasury address
	testReferrerAddr := "0x2864Cd9a59f32b74f3f851B92973fD40883aD503" // Another test address

	// Create test period (use current timestamp)
	period := uint64(time.Now().Unix())

	fmt.Printf("📋 Test Settlement Parameters:\n")
	fmt.Printf("   Period ID:         %d\n", period)
	fmt.Printf("   User Address:      %s\n", testUserAddr)
	fmt.Printf("   Developer:         %s\n", testDeveloperAddr)
	fmt.Printf("   Referrer:          %s\n", testReferrerAddr)
	fmt.Println()

	// Create test mining rewards (MULTIPLE entries to generate proper Merkle proofs)
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
			contributorAddr:   testUserAddr,
			contributorAmount: new(big.Int).Mul(big.NewInt(126), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), // 126 KAWAI
			developerAmount:   new(big.Int).Mul(big.NewInt(6), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 6 KAWAI
			userAmount:        new(big.Int).Mul(big.NewInt(6), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 6 KAWAI
			affiliatorAmount:  new(big.Int).Mul(big.NewInt(6), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 6 KAWAI
			developerAddr:     testDeveloperAddr,
			userAddr:          testUserAddr,
			affiliatorAddr:    testReferrerAddr,
		},
		// Entry 2: Another contributor (to create multiple leaves)
		{
			contributorAddr:   "0x1111111111111111111111111111111111111111",
			contributorAmount: new(big.Int).Mul(big.NewInt(100), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), // 100 KAWAI
			developerAmount:   new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 5 KAWAI
			userAmount:        new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 5 KAWAI
			affiliatorAmount:  new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 5 KAWAI
			developerAddr:     testDeveloperAddr,
			userAddr:          "0x2222222222222222222222222222222222222222", // Different user
			affiliatorAddr:    testReferrerAddr,
		},
		// Entry 3: Third contributor (to create even more leaves)
		{
			contributorAddr:   "0x3333333333333333333333333333333333333333",
			contributorAmount: new(big.Int).Mul(big.NewInt(80), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), // 80 KAWAI
			developerAmount:   new(big.Int).Mul(big.NewInt(4), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),  // 4 KAWAI
			userAmount:        new(big.Int).Mul(big.NewInt(4), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),  // 4 KAWAI
			affiliatorAmount:  new(big.Int).Mul(big.NewInt(4), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),  // 4 KAWAI
			developerAddr:     testDeveloperAddr,
			userAddr:          "0x4444444444444444444444444444444444444444", // Different user
			affiliatorAddr:    testReferrerAddr,
		},
	}

	// Generate Merkle leaves
	var leaves [][]byte
	proofs := make(map[string]*store.MerkleProofData)
	totalAmount := big.NewInt(0)

	for i, reward := range testRewards {
		// Generate 9-field Merkle leaf
		leaf := generateMiningMerkleLeaf(
			period,
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

		// Store proof data
		proofs[reward.contributorAddr] = &store.MerkleProofData{
			Index:             uint64(i),
			Amount:            reward.contributorAmount.String(),
			PeriodID:          int64(period),
			RewardType:        types.RewardTypeMining,
			ContributorAmount: reward.contributorAmount.String(),
			DeveloperAmount:   reward.developerAmount.String(),
			UserAmount:        reward.userAmount.String(),
			AffiliatorAmount:  reward.affiliatorAmount.String(),
			DeveloperAddress:  reward.developerAddr,
			UserAddress:       reward.userAddr,
			AffiliatorAddress: reward.affiliatorAddr,
			ClaimStatus:       store.ClaimStatusUnclaimed,
			CreatedAt:         time.Now(),
		}

		totalAmount.Add(totalAmount, reward.contributorAmount)
	}

	fmt.Printf("🌳 Generated %d Merkle leaves (multiple leaves = proper proofs)\n", len(leaves))

	// Build Merkle tree
	tree := merkle.NewMerkleTree(leaves)
	root := tree.Root

	fmt.Printf("📊 Merkle Root: 0x%x\n", root)
	fmt.Println()

	// Generate proofs for each leaf
	for i, leaf := range leaves {
		proof, ok := tree.GetProof(leaf)
		if !ok {
			return fmt.Errorf("failed to generate proof for leaf %d", i)
		}

		var proofHex []string
		for _, p := range proof {
			proofHex = append(proofHex, fmt.Sprintf("0x%x", p))
		}

		// Find corresponding contributor
		for contributorAddr, proofData := range proofs {
			if proofData.Index == uint64(i) {
				proofData.Proof = proofHex
				proofData.MerkleRoot = fmt.Sprintf("0x%x", root)
				fmt.Printf("✅ Generated proof for %s (index %d, %d elements)\n",
					contributorAddr, i, len(proofHex))
				break
			}
		}
	}

	// Save settlement
	fmt.Println()
	fmt.Printf("💾 Saving proper test settlement...\n")

	settlement, err := kv.PerformSettlementParallel(
		ctx,
		int64(period),
		fmt.Sprintf("0x%x", root),
		"kawai",
		proofs,
		1, // 1 worker for test
	)
	if err != nil {
		return fmt.Errorf("failed to save settlement: %w", err)
	}

	fmt.Printf("✅ Proper test settlement created successfully!\n")
	fmt.Println()
	fmt.Printf("📋 Settlement Details:\n")
	fmt.Printf("   Period ID:     %d\n", settlement.PeriodID)
	fmt.Printf("   Merkle Root:   %s\n", settlement.MerkleRoot)
	fmt.Printf("   Total Amount:  %s KAWAI\n", settlement.TotalAmount)
	fmt.Printf("   Contributors:  %d\n", settlement.ContributorCount)
	fmt.Printf("   Proofs Saved:  %d\n", settlement.ProofsSaved)
	fmt.Println()

	fmt.Printf("🎯 Key Improvements:\n")
	fmt.Printf("   ✅ Multiple leaves = proper Merkle proofs (not empty)\n")
	fmt.Printf("   ✅ Correct user address: %s\n", testUserAddr)
	fmt.Printf("   ✅ Valid proof elements for verification\n")
	fmt.Println()

	fmt.Printf("📝 Next Steps:\n")
	fmt.Printf("   1. Update period mapping to include period %d\n", period)
	fmt.Printf("   2. Upload this Merkle root to contract as period 6\n")
	fmt.Printf("   3. Test claiming - should work with proper proof validation\n")

	return nil
}

// generateMiningMerkleLeaf creates a 9-field Merkle leaf for MiningRewardDistributor
// Matches the Solidity keccak256(abi.encodePacked(...)) format
func generateMiningMerkleLeaf(
	period uint64,
	contributor common.Address,
	contributorAmt, developerAmt, userAmt, affiliatorAmt *big.Int,
	developer, user, affiliator common.Address,
) []byte {
	// Solidity abi.encodePacked packs values tightly without padding
	// For uint256, it uses 32 bytes; for address, it uses 20 bytes
	return crypto.Keccak256(
		common.LeftPadBytes(big.NewInt(int64(period)).Bytes(), 32), // uint256
		contributor.Bytes(),                             // address (20 bytes)
		common.LeftPadBytes(contributorAmt.Bytes(), 32), // uint256
		common.LeftPadBytes(developerAmt.Bytes(), 32),   // uint256
		common.LeftPadBytes(userAmt.Bytes(), 32),        // uint256
		common.LeftPadBytes(affiliatorAmt.Bytes(), 32),  // uint256
		developer.Bytes(),                               // address (20 bytes)
		user.Bytes(),                                    // address (20 bytes)
		affiliator.Bytes(),                              // address (20 bytes)
	)
}
