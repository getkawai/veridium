package main

import (
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kawai-network/veridium/pkg/merkle"
)

func main() {
	if err := debugTreeStructure(); err != nil {
		log.Fatalf("Failed to debug tree structure: %v", err)
	}
}

func debugTreeStructure() error {
	fmt.Println("🔍 Tree Structure Debug")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// Recreate the exact same leaves as in create-correct-mining-settlement
	testAddress := "0x0f3e75B9Bb3efcD87B1Ed15a30C8a7FBaABD204E"
	testDeveloperAddr := "0x94D5C06229811c4816107005ff05259f229Eb07b"
	testReferrerAddr := "0x2864Cd9a59f32b74f3f851B92973fD40883aD503"
	contractPeriod := uint64(7)

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
			developerAddr:     testDeveloperAddr,
			userAddr:          testAddress,
			affiliatorAddr:    testReferrerAddr,
		},
		// Entry 2: Another contributor
		{
			contributorAddr:   "0x1111111111111111111111111111111111111111",
			contributorAmount: new(big.Int).Mul(big.NewInt(100), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)), // 100 KAWAI
			developerAmount:   new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 5 KAWAI
			userAmount:        new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 5 KAWAI
			affiliatorAmount:  new(big.Int).Mul(big.NewInt(5), new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)),   // 5 KAWAI
			developerAddr:     testDeveloperAddr,
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
			developerAddr:     testDeveloperAddr,
			userAddr:          "0x4444444444444444444444444444444444444444",
			affiliatorAddr:    testReferrerAddr,
		},
	}

	// Generate leaves
	var leaves [][]byte
	for i, reward := range testRewards {
		leaf := generateCorrectMiningMerkleLeaf(
			contractPeriod,
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

		fmt.Printf("📋 Leaf %d (%s):\n", i, reward.contributorAddr)
		fmt.Printf("   Hash: 0x%x\n", leaf)
		fmt.Printf("   Contributor: %s\n", reward.contributorAddr)
		fmt.Printf("   Amount: %s\n", reward.contributorAmount.String())
		fmt.Println()
	}

	// Build tree
	tree := merkle.NewMerkleTree(leaves)

	fmt.Printf("🌳 Tree Structure:\n")
	fmt.Printf("   Root: 0x%x\n", tree.Root)
	fmt.Printf("   Layers: %d\n", len(tree.Layers))
	for i, layer := range tree.Layers {
		fmt.Printf("   Layer %d (%d nodes):\n", i, len(layer))
		for j, node := range layer {
			fmt.Printf("     [%d]: 0x%x\n", j, node)
		}
	}
	fmt.Println()

	// Generate proof for our test leaf (index 0)
	testLeaf := leaves[0]
	proof, ok := tree.GetProof(testLeaf)
	if !ok {
		return fmt.Errorf("failed to generate proof for test leaf")
	}

	fmt.Printf("🔍 Proof for Test Leaf (index 0):\n")
	fmt.Printf("   Leaf: 0x%x\n", testLeaf)
	fmt.Printf("   Proof (%d elements):\n", len(proof))
	for i, p := range proof {
		fmt.Printf("     [%d]: 0x%x\n", i, p)
	}
	fmt.Println()

	// Check if any proof element matches the leaf
	for i, p := range proof {
		if fmt.Sprintf("%x", p) == fmt.Sprintf("%x", testLeaf) {
			fmt.Printf("❌ ERROR: Proof element [%d] matches the leaf!\n", i)
			return fmt.Errorf("proof contains leaf")
		}
	}

	fmt.Printf("✅ Proof does not contain the leaf itself\n")
	fmt.Println()

	// Manual verification
	fmt.Printf("🧪 Manual Verification:\n")
	computedHash := testLeaf
	fmt.Printf("   Start: 0x%x\n", computedHash)

	for i, proofElement := range proof {
		if string(computedHash) <= string(proofElement) {
			computedHash = crypto.Keccak256(computedHash, proofElement)
		} else {
			computedHash = crypto.Keccak256(proofElement, computedHash)
		}
		fmt.Printf("   Step %d: 0x%x\n", i+1, computedHash)
	}

	fmt.Printf("   Final: 0x%x\n", computedHash)
	fmt.Printf("   Root:  0x%x\n", tree.Root)

	isValid := fmt.Sprintf("%x", computedHash) == fmt.Sprintf("%x", tree.Root)
	fmt.Printf("   Valid: %t\n", isValid)

	if !isValid {
		return fmt.Errorf("verification failed")
	}

	fmt.Printf("✅ Tree structure and proof generation are correct!\n")
	return nil
}

// generateCorrectMiningMerkleLeaf creates a 9-field Merkle leaf using CONTRACT PERIOD
func generateCorrectMiningMerkleLeaf(
	contractPeriod uint64,
	contributor common.Address,
	contributorAmt, developerAmt, userAmt, affiliatorAmt *big.Int,
	developer, user, affiliator common.Address,
) []byte {
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
