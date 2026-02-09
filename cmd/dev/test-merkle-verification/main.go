package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/generate/abi/miningdistributor"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/veridium/pkg/types"
	"github.com/kawai-network/x/constant"
)

func main() {
	if err := testMerkleVerification(); err != nil {
		log.Fatalf("Failed to test Merkle verification: %v", err)
	}
}

func testMerkleVerification() error {
	ctx := context.Background()

	fmt.Println("🔍 Merkle Verification Test")
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
				if proof.RewardType == types.RewardTypeMining && proof.PeriodID == 1768137242 {
					testProof = proof
					break
				}
			}
		}
	}

	if testProof == nil {
		return fmt.Errorf("CORRECT test proof not found for period 1768137242")
	}

	// Connect to contract
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to Monad: %w", err)
	}
	defer client.Close()

	distributorAddr := common.HexToAddress(constant.MiningRewardDistributorAddress)
	distributor, err := miningdistributor.NewMiningRewardDistributor(distributorAddr, client)
	if err != nil {
		return fmt.Errorf("failed to load MiningRewardDistributor: %w", err)
	}

	// Check contract period and root
	contractPeriod := uint64(7)
	periodRoot, err := distributor.PeriodMerkleRoots(nil, big.NewInt(int64(contractPeriod)))
	if err != nil {
		return fmt.Errorf("failed to get period root: %w", err)
	}

	fmt.Printf("📋 Contract State:\n")
	fmt.Printf("   Period:     %d\n", contractPeriod)
	fmt.Printf("   Root:       0x%x\n", periodRoot)
	fmt.Printf("   Expected:   %s\n", testProof.MerkleRoot)
	fmt.Println()

	// Check if roots match
	expectedRootBytes := common.Hex2Bytes(testProof.MerkleRoot[2:])
	var expectedRoot [32]byte
	copy(expectedRoot[:], expectedRootBytes)

	if periodRoot != expectedRoot {
		fmt.Printf("❌ Root mismatch!\n")
		fmt.Printf("   Contract has: 0x%x\n", periodRoot)
		fmt.Printf("   Expected:     0x%x\n", expectedRoot)
		return fmt.Errorf("Merkle root mismatch")
	}

	fmt.Printf("✅ Merkle roots match!\n")
	fmt.Println()

	// Parse amounts and addresses
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
		common.LeftPadBytes(big.NewInt(int64(contractPeriod)).Bytes(), 32), // period
		msgSender.Bytes(), // msg.sender
		common.LeftPadBytes(contribAmt.Bytes(), 32), // contributorAmount
		common.LeftPadBytes(devAmt.Bytes(), 32),     // developerAmount
		common.LeftPadBytes(userAmt.Bytes(), 32),    // userAmount
		common.LeftPadBytes(affAmt.Bytes(), 32),     // affiliatorAmount
		developer.Bytes(),                           // developer
		user.Bytes(),                                // user
		affiliator.Bytes(),                          // affiliator
	)

	fmt.Printf("🌿 Generated Leaf: 0x%x\n", leaf)
	fmt.Println()

	// Convert proof to [32]byte format
	merkleProof := make([][32]byte, len(testProof.Proof))
	for i, p := range testProof.Proof {
		if p[:2] == "0x" {
			p = p[2:]
		}
		proofBytes := common.Hex2Bytes(p)
		if len(proofBytes) != 32 {
			return fmt.Errorf("invalid proof element at index %d: expected 32 bytes, got %d", i, len(proofBytes))
		}
		copy(merkleProof[i][:], proofBytes)
	}

	fmt.Printf("🔍 Proof Elements:\n")
	for i, p := range merkleProof {
		fmt.Printf("   [%d]: 0x%x\n", i, p)
	}
	fmt.Println()

	// Test manual verification (OpenZeppelin style)
	fmt.Printf("🧪 Manual OpenZeppelin-style Verification:\n")

	// OpenZeppelin MerkleProof.verify implementation
	computedHash := leaf
	for i := 0; i < len(merkleProof); i++ {
		proofElement := merkleProof[i]

		// Compare as bytes (lexicographic order)
		if string(computedHash) <= string(proofElement[:]) {
			// Hash(current computed hash + current element of the proof)
			computedHash = crypto.Keccak256(computedHash, proofElement[:])
		} else {
			// Hash(current element of the proof + current computed hash)
			computedHash = crypto.Keccak256(proofElement[:], computedHash)
		}

		fmt.Printf("   Step %d: 0x%x\n", i+1, computedHash)
	}

	var computedRoot [32]byte
	copy(computedRoot[:], computedHash)

	isValid := computedRoot == expectedRoot
	fmt.Printf("   Final Hash: 0x%x\n", computedHash)
	fmt.Printf("   Expected:   0x%x\n", expectedRoot)
	fmt.Printf("   Valid:      %t\n", isValid)
	fmt.Println()

	if !isValid {
		fmt.Printf("❌ Manual verification failed!\n")
		return fmt.Errorf("manual verification failed")
	}

	fmt.Printf("✅ Manual verification successful!\n")
	fmt.Printf("🎯 The proof should work with the contract\n")

	return nil
}

// Helper function to compare byte arrays as big-endian integers
func compareBytesLE(a, b [32]byte) bool {
	for i := 0; i < 32; i++ {
		if a[i] < b[i] {
			return true
		} else if a[i] > b[i] {
			return false
		}
	}
	return false // equal
}
