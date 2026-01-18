package main

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/internal/generate/abi/miningdistributor"
	"github.com/kawai-network/veridium/pkg/store"
)

func main() {
	if err := testFixedMiningClaim(); err != nil {
		log.Fatalf("Failed to test fixed mining claim: %v", err)
	}
}

func testFixedMiningClaim() error {
	ctx := context.Background()

	fmt.Println("🧪 FIXED Mining Claim Test (Correct Proofs)")
	fmt.Println("═══════════════════════════════════════════════")
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
				if proof.RewardType == "kawai" && proof.PeriodID == 1768139780 {
					testProof = proof
					break
				}
			}
		}
	}

	if testProof == nil {
		return fmt.Errorf("FIXED test proof not found for period 1768139780")
	}

	fmt.Printf("📋 Found FIXED Test Proof:\n")
	fmt.Printf("   Settlement ID:     %d (KV storage)\n", testProof.PeriodID)
	fmt.Printf("   User Address:      %s\n", testProof.UserAddress)
	fmt.Printf("   Contributor Amount: %s KAWAI\n", testProof.ContributorAmount)
	fmt.Printf("   Developer Amount:   %s KAWAI\n", testProof.DeveloperAmount)
	fmt.Printf("   User Amount:        %s KAWAI\n", testProof.UserAmount)
	fmt.Printf("   Affiliator Amount:  %s KAWAI\n", testProof.AffiliatorAmount)
	fmt.Printf("   Proof Elements:    %d\n", len(testProof.Proof))
	for i, p := range testProof.Proof {
		fmt.Printf("     [%d]: %s\n", i, p)
	}
	fmt.Println()

	// Verify proof doesn't contain the leaf itself
	contractPeriod := uint64(7)

	// Generate the expected leaf
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

	expectedLeaf := crypto.Keccak256(
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

	fmt.Printf("🌿 Expected Leaf: 0x%x\n", expectedLeaf)

	// Check if proof contains the leaf
	leafInProof := false
	for i, p := range testProof.Proof {
		if p[:2] == "0x" {
			p = p[2:]
		}
		proofBytes := common.Hex2Bytes(p)
		if fmt.Sprintf("%x", proofBytes) == fmt.Sprintf("%x", expectedLeaf) {
			fmt.Printf("❌ ERROR: Proof element [%d] contains the leaf itself!\n", i)
			leafInProof = true
		}
	}

	if leafInProof {
		return fmt.Errorf("proof contains the leaf itself")
	}

	fmt.Printf("✅ Proof does not contain the leaf itself\n")
	fmt.Println()

	// Map settlement period to contract period
	contractPeriodMapped, err := mapSettlementPeriodToContractPeriod(testProof.PeriodID)
	if err != nil {
		return fmt.Errorf("failed to map period: %w", err)
	}

	fmt.Printf("🔄 Period Mapping: %d -> %d\n", testProof.PeriodID, contractPeriodMapped)
	fmt.Printf("🎯 Key Fix: Both settlement and contract use period %d for validation\n", contractPeriodMapped)
	fmt.Println()

	// Connect to Monad RPC
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to Monad: %w", err)
	}
	defer client.Close()

	// Load MiningRewardDistributor contract
	distributorAddr := common.HexToAddress(constant.MiningRewardDistributorAddr)
	distributor, err := miningdistributor.NewMiningRewardDistributor(distributorAddr, client)
	if err != nil {
		return fmt.Errorf("failed to load MiningRewardDistributor: %w", err)
	}

	// Check current contract period
	currentPeriod, err := distributor.CurrentPeriod(nil)
	if err != nil {
		return fmt.Errorf("failed to get current period: %w", err)
	}

	fmt.Printf("📊 Contract current period: %d\n", currentPeriod.Uint64())

	if currentPeriod.Uint64() != uint64(contractPeriodMapped) {
		return fmt.Errorf("contract period mismatch: expected %d, got %d", contractPeriodMapped, currentPeriod.Uint64())
	}

	// Check if already claimed
	claimed, err := distributor.HasClaimedPeriod(nil, big.NewInt(contractPeriodMapped), common.HexToAddress(testAddress))
	if err != nil {
		return fmt.Errorf("failed to check claim status: %w", err)
	}

	if claimed {
		fmt.Printf("⚠️  Already claimed for period %d\n", contractPeriodMapped)
		return nil
	}

	// Get private key for signing
	privateKeyHex := constant.GetAdminPrivateKey()
	if privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get chain ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}
	auth.Context = ctx

	// Convert proof strings to [32]byte array
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

	fmt.Printf("🚀 Submitting FIXED Mining Claim Transaction...\n")
	fmt.Printf("   Contract Period:    %d (matches Merkle leaf generation)\n", contractPeriodMapped)
	fmt.Printf("   Contributor Amount: %s\n", testProof.ContributorAmount)
	fmt.Printf("   Developer Amount:   %s\n", testProof.DeveloperAmount)
	fmt.Printf("   User Amount:        %s\n", testProof.UserAmount)
	fmt.Printf("   Affiliator Amount:  %s\n", testProof.AffiliatorAmount)
	fmt.Printf("   Developer Address:  %s\n", testProof.DeveloperAddress)
	fmt.Printf("   User Address:       %s\n", testProof.UserAddress)
	fmt.Printf("   Affiliator Address: %s\n", testProof.AffiliatorAddress)
	fmt.Printf("   Proof Elements:     %d\n", len(merkleProof))
	fmt.Println()

	// Submit claim transaction
	tx, err := distributor.ClaimReward(
		auth,
		big.NewInt(contractPeriodMapped), // Use contract period (7)
		contribAmt,
		devAmt,
		userAmt,
		affAmt,
		developer,
		user,
		affiliator,
		merkleProof,
	)
	if err != nil {
		fmt.Printf("❌ Mining claim transaction failed: %v\n", err)
		return fmt.Errorf("mining claim transaction failed: %w", err)
	}

	fmt.Printf("✅ Transaction submitted: %s\n", tx.Hash().Hex())
	fmt.Printf("⏳ Waiting for confirmation...\n")

	// Wait for confirmation
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for confirmation: %w", err)
	}

	if receipt.Status != 1 {
		fmt.Printf("❌ Transaction failed with status: %d\n", receipt.Status)
		return fmt.Errorf("transaction reverted")
	}

	fmt.Printf("🎉 FIXED Mining claim successful!\n")
	fmt.Printf("   Block Number: %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("   Gas Used:     %d\n", receipt.GasUsed)
	fmt.Printf("   TX Hash:      %s\n", tx.Hash().Hex())
	fmt.Println()

	fmt.Printf("✅ PROOF FIX VERIFIED!\n")
	fmt.Printf("🎯 Merkle proof validation worked with correct proofs\n")
	fmt.Printf("📝 Mining claims should now work through the UI\n")

	return nil
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
		1768137123: 7, // Previous CORRECT settlement -> Contract period 7
		1768137242: 7, // LATEST CORRECT settlement with matching periods -> Contract period 7
		1768139780: 7, // FIXED settlement with correct proofs -> Contract period 7
	}

	contractPeriod, exists := periodMapping[settlementPeriodID]
	if !exists {
		return 0, fmt.Errorf("unknown settlement period ID: %d", settlementPeriodID)
	}

	return contractPeriod, nil
}
