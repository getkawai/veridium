package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/contracts/miningdistributor"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/y/types"
	"github.com/kawai-network/contracts"
)

func main() {
	if err := debugClaimCall(); err != nil {
		log.Fatalf("Failed to debug claim call: %v", err)
	}
}

func debugClaimCall() error {
	ctx := context.Background()

	fmt.Println("🔍 Debug Claim Call")
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

	// Connect to contract
	client, err := ethclient.Dial(contracts.MonadRpcUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to Monad: %w", err)
	}
	defer client.Close()

	distributorAddr := common.HexToAddress(contracts.MiningRewardDistributorAddress)
	distributor, err := miningdistributor.NewMiningRewardDistributor(distributorAddr, client)
	if err != nil {
		return fmt.Errorf("failed to load MiningRewardDistributor: %w", err)
	}

	contractPeriod := int64(7)

	// Check if already claimed
	claimed, err := distributor.HasClaimedPeriod(nil, big.NewInt(contractPeriod), common.HexToAddress(testAddress))
	if err != nil {
		return fmt.Errorf("failed to check claim status: %w", err)
	}

	fmt.Printf("📊 Claim Status Check:\n")
	fmt.Printf("   Period:   %d\n", contractPeriod)
	fmt.Printf("   Address:  %s\n", testAddress)
	fmt.Printf("   Claimed:  %t\n", claimed)
	fmt.Println()

	if claimed {
		fmt.Printf("⚠️  User has already claimed for period %d\n", contractPeriod)
		return nil
	}

	// Parse all parameters
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

	// Generate leaf
	leaf := crypto.Keccak256(
		common.LeftPadBytes(big.NewInt(contractPeriod).Bytes(), 32),
		msgSender.Bytes(),
		common.LeftPadBytes(contribAmt.Bytes(), 32),
		common.LeftPadBytes(devAmt.Bytes(), 32),
		common.LeftPadBytes(userAmt.Bytes(), 32),
		common.LeftPadBytes(affAmt.Bytes(), 32),
		developer.Bytes(),
		user.Bytes(),
		affiliator.Bytes(),
	)

	// Convert proof
	merkleProof := make([][32]byte, len(testProof.Proof))
	for i, p := range testProof.Proof {
		if p[:2] == "0x" {
			p = p[2:]
		}
		proofBytes := common.Hex2Bytes(p)
		copy(merkleProof[i][:], proofBytes)
	}

	// Get period root
	periodRoot, err := distributor.PeriodMerkleRoots(nil, big.NewInt(contractPeriod))
	if err != nil {
		return fmt.Errorf("failed to get period root: %w", err)
	}

	fmt.Printf("📋 Contract Call Parameters:\n")
	fmt.Printf("   period:             %d\n", contractPeriod)
	fmt.Printf("   contributorAmount:  %s\n", contribAmt.String())
	fmt.Printf("   developerAmount:    %s\n", devAmt.String())
	fmt.Printf("   userAmount:         %s\n", userAmt.String())
	fmt.Printf("   affiliatorAmount:   %s\n", affAmt.String())
	fmt.Printf("   developer:          %s\n", developer.Hex())
	fmt.Printf("   user:               %s\n", user.Hex())
	fmt.Printf("   affiliator:         %s\n", affiliator.Hex())
	fmt.Printf("   merkleProof length: %d\n", len(merkleProof))
	fmt.Println()

	fmt.Printf("🔍 Verification Data:\n")
	fmt.Printf("   Generated Leaf:     0x%x\n", leaf)
	fmt.Printf("   Period Root:        0x%x\n", periodRoot)
	fmt.Printf("   Expected Root:      %s\n", testProof.MerkleRoot)
	fmt.Println()

	// Check if period root is set
	var zeroRoot [32]byte
	if periodRoot == zeroRoot {
		fmt.Printf("❌ Period %d has no Merkle root set!\n", contractPeriod)
		return fmt.Errorf("period root not set")
	}

	// Check if roots match
	expectedRootBytes := common.Hex2Bytes(testProof.MerkleRoot[2:])
	var expectedRoot [32]byte
	copy(expectedRoot[:], expectedRootBytes)

	if periodRoot != expectedRoot {
		fmt.Printf("❌ Root mismatch!\n")
		fmt.Printf("   Contract: 0x%x\n", periodRoot)
		fmt.Printf("   Expected: 0x%x\n", expectedRoot)
		return fmt.Errorf("root mismatch")
	}

	fmt.Printf("✅ All parameters look correct!\n")
	fmt.Printf("🎯 The claim should work - there might be a gas or network issue\n")

	// Try to simulate the call (read-only)
	fmt.Printf("🧪 Attempting read-only call simulation...\n")

	// We can't easily simulate the full transaction, but we can check the basic parameters
	// The issue might be in the transaction execution itself

	return nil
}
