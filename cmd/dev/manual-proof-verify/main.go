package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kawai-network/x/store"
	"github.com/kawai-network/y/types"
)

func main() {
	if err := manualProofVerify(); err != nil {
		log.Fatalf("Failed to manually verify proof: %v", err)
	}
}

func manualProofVerify() error {
	ctx := context.Background()

	fmt.Println("🔍 Manual Proof Verification (OpenZeppelin Style)")
	fmt.Println("═══════════════════════════════════════════════════")
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

	// Find the FIXED settlement (Period 1768139780)
	var testProof *store.MerkleProofData
	if unclaimedRaw, ok := claimableData["unclaimed_proofs"]; ok {
		if unclaimedList, ok := unclaimedRaw.([]*store.MerkleProofData); ok {
			for _, proof := range unclaimedList {
				if proof.RewardType == types.RewardTypeMining && proof.PeriodID == 1768139780 {
					testProof = proof
					break
				}
			}
		}
	}

	if testProof == nil {
		return fmt.Errorf("FIXED test proof not found for period 1768139780")
	}

	fmt.Printf("📋 Test Proof Data:\n")
	fmt.Printf("   Settlement ID:     %d\n", testProof.PeriodID)
	fmt.Printf("   Merkle Root:       %s\n", testProof.MerkleRoot)
	fmt.Printf("   Proof Elements:    %d\n", len(testProof.Proof))
	for i, p := range testProof.Proof {
		fmt.Printf("     [%d]: %s\n", i, p)
	}
	fmt.Println()

	// Parse amounts and addresses
	contractPeriod := uint64(7)
	contribAmt := new(big.Int)
	contribAmt.SetString(testProof.ContributorAmount, 10)
	devAmt := new(big.Int)
	devAmt.SetString(testProof.DeveloperAmount, 10)
	userAmt := new(big.Int)
	userAmt.SetString(testProof.UserAmount, 10)
	affAmt := new(big.Int)
	affAmt.SetString(testProof.AffiliatorAmount, 10)

	msgSender := common.HexToAddress(testAddress)
	developer := common.HexToAddress(testProof.DeveloperAddress)
	user := common.HexToAddress(testProof.UserAddress)
	affiliator := common.HexToAddress(testProof.AffiliatorAddress)

	// Generate leaf exactly as contract does
	leaf := crypto.Keccak256(
		common.LeftPadBytes(big.NewInt(int64(contractPeriod)).Bytes(), 32),
		msgSender.Bytes(),
		common.LeftPadBytes(contribAmt.Bytes(), 32),
		common.LeftPadBytes(devAmt.Bytes(), 32),
		common.LeftPadBytes(userAmt.Bytes(), 32),
		common.LeftPadBytes(affAmt.Bytes(), 32),
		developer.Bytes(),
		user.Bytes(),
		affiliator.Bytes(),
	)

	fmt.Printf("🌿 Generated Leaf: 0x%x\n", leaf)
	fmt.Println()

	// Parse expected root
	expectedRootHex := testProof.MerkleRoot
	if expectedRootHex[:2] == "0x" {
		expectedRootHex = expectedRootHex[2:]
	}
	expectedRootBytes := common.Hex2Bytes(expectedRootHex)
	var expectedRoot [32]byte
	copy(expectedRoot[:], expectedRootBytes)

	fmt.Printf("🎯 Expected Root: 0x%x\n", expectedRoot)
	fmt.Println()

	// Convert proof to bytes
	proofBytes := make([][]byte, len(testProof.Proof))
	for i, p := range testProof.Proof {
		if p[:2] == "0x" {
			p = p[2:]
		}
		proofBytes[i] = common.Hex2Bytes(p)
	}

	fmt.Printf("🔍 Proof Elements (bytes):\n")
	for i, p := range proofBytes {
		fmt.Printf("   [%d]: 0x%x\n", i, p)
	}
	fmt.Println()

	// Manual OpenZeppelin-style verification
	fmt.Printf("🧪 OpenZeppelin MerkleProof.verify Simulation:\n")

	computedHash := leaf
	fmt.Printf("   Start: 0x%x (leaf)\n", computedHash)

	for i, proofElement := range proofBytes {
		fmt.Printf("   Step %d:\n", i+1)
		fmt.Printf("     Current:  0x%x\n", computedHash)
		fmt.Printf("     Proof[%d]: 0x%x\n", i, proofElement)

		// OpenZeppelin comparison: bytes32(a) < bytes32(b)
		// This is lexicographic comparison of the byte arrays
		if bytesLessThan(computedHash, proofElement) {
			fmt.Printf("     Order:    current < proof[%d] -> hash(current, proof[%d])\n", i, i)
			computedHash = crypto.Keccak256(computedHash, proofElement)
		} else {
			fmt.Printf("     Order:    current >= proof[%d] -> hash(proof[%d], current)\n", i, i)
			computedHash = crypto.Keccak256(proofElement, computedHash)
		}
		fmt.Printf("     Result:   0x%x\n", computedHash)
		fmt.Println()
	}

	fmt.Printf("🎯 Final Verification:\n")
	fmt.Printf("   Computed: 0x%x\n", computedHash)
	fmt.Printf("   Expected: 0x%x\n", expectedRoot)

	var computedRoot [32]byte
	copy(computedRoot[:], computedHash)

	isValid := computedRoot == expectedRoot
	fmt.Printf("   Valid:    %t\n", isValid)

	if !isValid {
		fmt.Printf("❌ Manual verification failed!\n")
		return fmt.Errorf("manual verification failed")
	}

	fmt.Printf("✅ Manual verification successful!\n")
	fmt.Printf("🎯 The proof should work with the contract\n")

	return nil
}

// bytesLessThan compares two byte slices lexicographically
// This matches Solidity's bytes32(a) < bytes32(b) comparison
func bytesLessThan(a, b []byte) bool {
	// Pad to 32 bytes if needed
	aPadded := make([]byte, 32)
	bPadded := make([]byte, 32)

	copy(aPadded[32-len(a):], a)
	copy(bPadded[32-len(b):], b)

	for i := 0; i < 32; i++ {
		if aPadded[i] < bPadded[i] {
			return true
		} else if aPadded[i] > bPadded[i] {
			return false
		}
	}
	return false // equal
}
