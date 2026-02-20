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
	if err := verifyMerkleProof(); err != nil {
		log.Fatalf("Failed to verify Merkle proof: %v", err)
	}
}

func verifyMerkleProof() error {
	ctx := context.Background()

	fmt.Println("🔍 Merkle Proof Verification")
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

	fmt.Printf("📋 Test Proof Details:\n")
	fmt.Printf("   Period ID:         %d\n", testProof.PeriodID)
	fmt.Printf("   Index:             %d\n", testProof.Index)
	fmt.Printf("   Merkle Root:       %s\n", testProof.MerkleRoot)
	fmt.Printf("   User Address:      %s\n", testProof.UserAddress)
	fmt.Printf("   Contributor Amount: %s\n", testProof.ContributorAmount)
	fmt.Printf("   Developer Amount:   %s\n", testProof.DeveloperAmount)
	fmt.Printf("   User Amount:        %s\n", testProof.UserAmount)
	fmt.Printf("   Affiliator Amount:  %s\n", testProof.AffiliatorAmount)
	fmt.Printf("   Developer Address:  %s\n", testProof.DeveloperAddress)
	fmt.Printf("   Affiliator Address: %s\n", testProof.AffiliatorAddress)
	fmt.Printf("   Proof Elements:     %d\n", len(testProof.Proof))
	for i, p := range testProof.Proof {
		fmt.Printf("     [%d] %s\n", i, p)
	}
	fmt.Println()

	// Recreate the Merkle leaf using the same logic as settlement generation
	period := uint64(testProof.PeriodID)
	contributor := common.HexToAddress(testAddress)

	contribAmt := new(big.Int)
	contribAmt.SetString(testProof.ContributorAmount, 10)

	devAmt := new(big.Int)
	devAmt.SetString(testProof.DeveloperAmount, 10)

	userAmt := new(big.Int)
	userAmt.SetString(testProof.UserAmount, 10)

	affAmt := new(big.Int)
	affAmt.SetString(testProof.AffiliatorAmount, 10)

	developer := common.HexToAddress(testProof.DeveloperAddress)
	user := common.HexToAddress(testProof.UserAddress)
	affiliator := common.HexToAddress(testProof.AffiliatorAddress)

	// Generate the leaf hash
	leaf := generateMiningMerkleLeaf(
		period,
		contributor,
		contribAmt, devAmt, userAmt, affAmt,
		developer, user, affiliator,
	)

	fmt.Printf("🌿 Generated Leaf Hash: 0x%x\n", leaf)
	fmt.Println()

	// Verify the proof manually
	fmt.Printf("🔍 Manual Proof Verification:\n")

	// Parse Merkle root
	merkleRootHex := testProof.MerkleRoot
	if merkleRootHex[:2] == "0x" {
		merkleRootHex = merkleRootHex[2:]
	}
	expectedRoot := common.Hex2Bytes(merkleRootHex)

	// Start with the leaf
	currentHash := leaf
	fmt.Printf("   Starting with leaf: 0x%x\n", currentHash)

	// Apply each proof element
	for i, proofHex := range testProof.Proof {
		if proofHex[:2] == "0x" {
			proofHex = proofHex[2:]
		}
		proofBytes := common.Hex2Bytes(proofHex)

		// Hash with proof element (order matters in Merkle trees)
		// Try both orders to see which one works
		hash1 := crypto.Keccak256(currentHash, proofBytes)
		hash2 := crypto.Keccak256(proofBytes, currentHash)

		fmt.Printf("   Proof[%d]: 0x%x\n", i, proofBytes)
		fmt.Printf("     Option 1 (leaf+proof): 0x%x\n", hash1)
		fmt.Printf("     Option 2 (proof+leaf): 0x%x\n", hash2)

		// For now, let's use the standard order (smaller hash first)
		if string(currentHash) < string(proofBytes) {
			currentHash = hash1
		} else {
			currentHash = hash2
		}

		fmt.Printf("     Selected: 0x%x\n", currentHash)
	}

	fmt.Printf("   Final computed root: 0x%x\n", currentHash)
	fmt.Printf("   Expected root:       0x%x\n", expectedRoot)
	fmt.Println()

	// Check if they match
	if string(currentHash) == string(expectedRoot) {
		fmt.Printf("✅ Proof verification PASSED!\n")
	} else {
		fmt.Printf("❌ Proof verification FAILED!\n")
		fmt.Printf("   This explains why the contract rejects the proof\n")
	}
	fmt.Println()

	// Let's also check what the contract expects
	fmt.Printf("🔍 Contract Validation Check:\n")
	fmt.Printf("   The contract will:\n")
	fmt.Printf("   1. Recreate leaf using: keccak256(abi.encodePacked(period, contributor, amounts, addresses))\n")
	fmt.Printf("   2. Verify proof against stored Merkle root for period\n")
	fmt.Printf("   3. Check that user address in leaf matches msg.sender\n")
	fmt.Println()

	// Debug the leaf generation format
	fmt.Printf("🔍 Leaf Generation Debug:\n")
	fmt.Printf("   Period (uint256):      %d (0x%x)\n", period, period)
	fmt.Printf("   Contributor (address): %s\n", contributor.Hex())
	fmt.Printf("   Contrib Amount:        %s\n", contribAmt.String())
	fmt.Printf("   Developer Amount:      %s\n", devAmt.String())
	fmt.Printf("   User Amount:           %s\n", userAmt.String())
	fmt.Printf("   Affiliator Amount:     %s\n", affAmt.String())
	fmt.Printf("   Developer Address:     %s\n", developer.Hex())
	fmt.Printf("   User Address:          %s\n", user.Hex())
	fmt.Printf("   Affiliator Address:    %s\n", affiliator.Hex())

	return nil
}

// generateMiningMerkleLeaf creates a 9-field Merkle leaf for MiningRewardDistributor
// Matches the Solidity keccak256(abi.encodePacked(...)) format
// IMPORTANT: Contract uses msg.sender as contributor, not the contributor field from data
func generateMiningMerkleLeaf(
	period uint64,
	contributor common.Address, // This should be the claiming address (msg.sender)
	contributorAmt, developerAmt, userAmt, affiliatorAmt *big.Int,
	developer, user, affiliator common.Address,
) []byte {
	// Solidity abi.encodePacked packs values tightly without padding
	// For uint256, it uses 32 bytes; for address, it uses 20 bytes
	// Order must match contract: period, msg.sender, amounts, addresses
	return crypto.Keccak256(
		common.LeftPadBytes(big.NewInt(int64(period)).Bytes(), 32), // uint256 period
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
