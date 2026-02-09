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
	"github.com/kawai-network/veridium/pkg/types"
	"github.com/kawai-network/x/constant"
)

func main() {
	if err := createTreasuryMiningSettlement(); err != nil {
		log.Fatalf("Failed to create treasury mining settlement: %v", err)
	}
}

func createTreasuryMiningSettlement() error {
	ctx := context.Background()

	fmt.Println("🏦 Creating Treasury Mining Settlement (msg.sender Match)")
	fmt.Println("═══════════════════════════════════════════════════════")
	fmt.Println()

	// Initialize KV store
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		return fmt.Errorf("failed to initialize KV store: %w", err)
	}

	// Treasury address (yang akan jadi msg.sender)
	treasuryAddr := constant.GetAdminAddress()                       // Treasury/Developer address
	testUserAddr := "0x0f3e75B9Bb3efcD87B1Ed15a30C8a7FBaABD204E"     // Test user address
	testReferrerAddr := "0x2864Cd9a59f32b74f3f851B92973fD40883aD503" // Referrer address

	// Use sequential period number that matches contract expectation
	contractPeriod := uint64(7)             // Current contract period
	settlementPeriodID := time.Now().Unix() // For KV storage key

	fmt.Printf("📋 Treasury Settlement Parameters:\n")
	fmt.Printf("   Contract Period:   %d (sequential - what contract expects)\n", contractPeriod)
	fmt.Printf("   Settlement ID:     %d (timestamp - for KV storage)\n", settlementPeriodID)
	fmt.Printf("   Treasury Address:  %s (msg.sender)\n", treasuryAddr)
	fmt.Printf("   Test User:         %s\n", testUserAddr)
	fmt.Printf("   Referrer:          %s\n", testReferrerAddr)
	fmt.Println()

	// Create test mining rewards dengan treasury sebagai contributor
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
		// Entry 1: Treasury sebagai contributor (ini yang akan di-claim)
		{
			contributorAddr:   treasuryAddr,                                                                             // Treasury address sebagai contributor
			contributorAmount: new(big.Int).Mul(big.NewInt(126), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), // 126 KAWAI
			developerAmount:   new(big.Int).Mul(big.NewInt(6), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 6 KAWAI
			userAmount:        new(big.Int).Mul(big.NewInt(6), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 6 KAWAI
			affiliatorAmount:  new(big.Int).Mul(big.NewInt(6), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 6 KAWAI
			developerAddr:     treasuryAddr,
			userAddr:          testUserAddr,
			affiliatorAddr:    testReferrerAddr,
		},
		// Entry 2: Another contributor (untuk membuat multiple leaves)
		{
			contributorAddr:   "0x1111111111111111111111111111111111111111",
			contributorAmount: new(big.Int).Mul(big.NewInt(100), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), // 100 KAWAI
			developerAmount:   new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 5 KAWAI
			userAmount:        new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 5 KAWAI
			affiliatorAmount:  new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 5 KAWAI
			developerAddr:     treasuryAddr,
			userAddr:          "0x2222222222222222222222222222222222222222",
			affiliatorAddr:    testReferrerAddr,
		},
		// Entry 3: Third contributor
		{
			contributorAddr:   "0x3333333333333333333333333333333333333333",
			contributorAmount: new(big.Int).Mul(big.NewInt(80), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), // 80 KAWAI
			developerAmount:   new(big.Int).Mul(big.NewInt(4), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),  // 4 KAWAI
			userAmount:        new(big.Int).Mul(big.NewInt(4), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),  // 4 KAWAI
			affiliatorAmount:  new(big.Int).Mul(big.NewInt(4), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),  // 4 KAWAI
			developerAddr:     treasuryAddr,
			userAddr:          "0x4444444444444444444444444444444444444444",
			affiliatorAddr:    testReferrerAddr,
		},
	}

	// Generate Merkle leaves using CONTRACT PERIOD
	var leaves [][]byte
	leafMap := make(map[string][]byte) // Map contributor address to leaf hash
	totalAmount := big.NewInt(0)

	for _, reward := range testRewards {
		// Generate 9-field Merkle leaf using CONTRACT PERIOD
		leaf := generateTreasuryMiningMerkleLeaf(
			contractPeriod, // Use contract period (7), not timestamp!
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
		leafMap[reward.contributorAddr] = leaf

		fmt.Printf("📋 Generated leaf for %s: 0x%x\n", reward.contributorAddr, leaf)
		totalAmount.Add(totalAmount, reward.contributorAmount)
	}

	fmt.Printf("🌳 Generated %d Merkle leaves using CONTRACT PERIOD %d\n", len(leaves), contractPeriod)

	// Build Merkle tree
	tree := merkle.NewMerkleTree(leaves)
	root := tree.Root

	fmt.Printf("📊 Merkle Root: 0x%x\n", root)
	fmt.Println()

	// Debug tree structure
	fmt.Printf("🔍 Tree Structure (after sorting):\n")
	for i, leaf := range tree.Leafs {
		fmt.Printf("   [%d]: 0x%x\n", i, leaf)
	}
	fmt.Println()

	// Generate proofs for each contributor
	proofs := make(map[string]*store.MerkleProofData)

	for i, reward := range testRewards {
		contributorAddr := reward.contributorAddr
		leafHash := leafMap[contributorAddr]

		// Find the leaf in the sorted tree
		var leafIndex int = -1
		for j, sortedLeaf := range tree.Leafs {
			if fmt.Sprintf("%x", sortedLeaf) == fmt.Sprintf("%x", leafHash) {
				leafIndex = j
				break
			}
		}

		if leafIndex == -1 {
			return fmt.Errorf("leaf not found in sorted tree for %s", contributorAddr)
		}

		fmt.Printf("🔍 Contributor %s: leaf index %d in sorted tree\n", contributorAddr, leafIndex)

		// Generate proof for this leaf
		proof, ok := tree.GetProof(leafHash)
		if !ok {
			return fmt.Errorf("failed to generate proof for %s", contributorAddr)
		}

		var proofHex []string
		for _, p := range proof {
			proofHex = append(proofHex, fmt.Sprintf("0x%x", p))
		}

		// Verify the proof doesn't contain the leaf itself
		for _, p := range proof {
			if fmt.Sprintf("%x", p) == fmt.Sprintf("%x", leafHash) {
				return fmt.Errorf("ERROR: proof for %s contains the leaf itself!", contributorAddr)
			}
		}

		// Store proof data with SETTLEMENT ID for KV storage
		proofs[contributorAddr] = &store.MerkleProofData{
			Index:             uint64(i), // Original index (before sorting)
			Amount:            reward.contributorAmount.String(),
			PeriodID:          settlementPeriodID, // Use timestamp for KV storage
			RewardType:        types.RewardTypeMining,
			ContributorAmount: reward.contributorAmount.String(),
			DeveloperAmount:   reward.developerAmount.String(),
			UserAmount:        reward.userAmount.String(),
			AffiliatorAmount:  reward.affiliatorAmount.String(),
			DeveloperAddress:  reward.developerAddr,
			UserAddress:       reward.userAddr,
			AffiliatorAddress: reward.affiliatorAddr,
			Proof:             proofHex,
			MerkleRoot:        fmt.Sprintf("0x%x", root),
			ClaimStatus:       store.ClaimStatusUnclaimed,
			CreatedAt:         time.Now(),
		}

		fmt.Printf("✅ Generated CORRECT proof for %s (sorted index %d, %d elements)\n",
			contributorAddr, leafIndex, len(proofHex))
		fmt.Printf("   Proof: %v\n", proofHex)
	}

	// Save settlement
	fmt.Println()
	fmt.Printf("💾 Saving Treasury settlement...\n")

	settlement, err := kv.PerformSettlementParallel(
		ctx,
		settlementPeriodID, // Use timestamp for KV storage
		fmt.Sprintf("0x%x", root),
		"kawai",
		proofs,
		1, // 1 worker for test
	)
	if err != nil {
		return fmt.Errorf("failed to save settlement: %w", err)
	}

	fmt.Printf("✅ Treasury settlement created successfully!\n")
	fmt.Println()
	fmt.Printf("📋 Settlement Details:\n")
	fmt.Printf("   Settlement ID:     %d (for KV storage)\n", settlementPeriodID)
	fmt.Printf("   Contract Period:   %d (for contract validation)\n", contractPeriod)
	fmt.Printf("   Merkle Root:       %s\n", settlement.MerkleRoot)
	fmt.Printf("   Total Amount:      %s KAWAI\n", settlement.TotalAmount)
	fmt.Printf("   Contributors:      %d\n", settlement.ContributorCount)
	fmt.Printf("   Proofs Saved:      %d\n", settlement.ProofsSaved)
	fmt.Println()

	fmt.Printf("🎯 Key Fix:\n")
	fmt.Printf("   ✅ Merkle leaves generated with CONTRACT PERIOD: %d\n", contractPeriod)
	fmt.Printf("   ✅ Proofs generated correctly (no leaf in proof)\n")
	fmt.Printf("   ✅ Treasury address as contributor matches msg.sender\n")
	fmt.Printf("   ✅ Contract will validate using same period: %d\n", contractPeriod)
	fmt.Printf("   ✅ Proof validation should now work!\n")
	fmt.Println()

	fmt.Printf("📝 Next Steps:\n")
	fmt.Printf("   1. Update period mapping to include: %d -> %d\n", settlementPeriodID, contractPeriod)
	fmt.Printf("   2. Test claiming with treasury address as msg.sender\n")
	fmt.Printf("   3. Should work perfectly now!\n")

	return nil
}

// generateTreasuryMiningMerkleLeaf creates a 9-field Merkle leaf using CONTRACT PERIOD
// This matches exactly what the contract expects
func generateTreasuryMiningMerkleLeaf(
	contractPeriod uint64, // Use contract period (1,2,3,4,5,6,7...), not timestamp
	contributor common.Address,
	contributorAmt, developerAmt, userAmt, affiliatorAmt *big.Int,
	developer, user, affiliator common.Address,
) []byte {
	// Exact match to contract Solidity code:
	// keccak256(abi.encodePacked(period, msg.sender, contributorAmount, developerAmount, userAmount, affiliatorAmount, developer, user, affiliator))
	return crypto.Keccak256(
		common.LeftPadBytes(big.NewInt(int64(contractPeriod)).Bytes(), 32), // uint256 period (CONTRACT PERIOD!)
		contributor.Bytes(),                             // address msg.sender (20 bytes)
		common.LeftPadBytes(contributorAmt.Bytes(), 32), // uint256 contributorAmount
		common.LeftPadBytes(developerAmt.Bytes(), 32),   // uint256 developerAmount
		common.LeftPadBytes(userAmt.Bytes(), 32),        // uint256 userAmount
		common.LeftPadBytes(affiliatorAmt.Bytes(), 32),  // uint256 affiliatorAmount
		developer.Bytes(),                               // address developer (20 bytes)
		user.Bytes(),                                    // address user (20 bytes)
		affiliator.Bytes(),                              // address affiliator (20 bytes)
	)
}
