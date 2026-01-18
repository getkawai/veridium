package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/internal/generate/abi/miningdistributor"
)

func main() {
	if err := uploadProperTestRoot(); err != nil {
		log.Fatalf("Failed to upload proper test root: %v", err)
	}
}

func uploadProperTestRoot() error {
	ctx := context.Background()

	fmt.Println("🚀 Uploading Proper Test Mining Root (With Valid Proofs)")
	fmt.Println("═══════════════════════════════════════════════════════")
	fmt.Println()

	// Proper test settlement details
	testMerkleRoot := "0xfb4b540bc1c0e330ec197047890e4f32dc1ae99e3146618decc31e2744cb22f2"
	testPeriodID := int64(1768136095)

	fmt.Printf("📋 Proper Test Settlement Details:\n")
	fmt.Printf("   Period ID:     %d\n", testPeriodID)
	fmt.Printf("   Merkle Root:   %s\n", testMerkleRoot)
	fmt.Printf("   User Address:  0x0f3e75B9Bb3efcD87B1Ed15a30C8a7FBaABD204E\n")
	fmt.Printf("   Leaves:        3 (multiple leaves = valid proofs)\n")
	fmt.Printf("   Proof Elements: 2 (proper Merkle proof validation)\n")
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

	// Check current period
	currentPeriod, err := distributor.CurrentPeriod(nil)
	if err != nil {
		return fmt.Errorf("failed to get current period: %w", err)
	}

	fmt.Printf("🎯 Contract current period: %d\n", currentPeriod.Uint64())
	fmt.Printf("📝 Will advance to period: %d\n", currentPeriod.Uint64()+1)
	fmt.Println()

	// Get private key and setup auth
	privateKeyHex := constant.GetAdminPrivateKey()
	if strings.HasPrefix(privateKeyHex, "0x") {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return fmt.Errorf("failed to create transactor: %w", err)
	}
	auth.Context = ctx

	// Parse Merkle root
	merkleRootHex := testMerkleRoot
	if strings.HasPrefix(merkleRootHex, "0x") {
		merkleRootHex = merkleRootHex[2:]
	}
	merkleRootBytes := common.Hex2Bytes(merkleRootHex)
	if len(merkleRootBytes) != 32 {
		return fmt.Errorf("invalid Merkle root length: expected 32 bytes, got %d", len(merkleRootBytes))
	}
	var merkleRoot [32]byte
	copy(merkleRoot[:], merkleRootBytes)

	// Confirm before advancing
	fmt.Printf("⚠️  About to advance contract to period %d with proper test Merkle root\n", currentPeriod.Uint64()+1)
	if !confirm("Continue with upload?") {
		return fmt.Errorf("operation cancelled by user")
	}

	// Advance period with new Merkle root
	tx, err := distributor.AdvancePeriod(auth, merkleRoot)
	if err != nil {
		return fmt.Errorf("failed to advance period: %w", err)
	}

	fmt.Printf("✅ AdvancePeriod transaction sent: %s\n", tx.Hash().Hex())
	fmt.Printf("⏳ Waiting for confirmation...\n")

	// Wait for confirmation
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		return fmt.Errorf("failed to wait for confirmation: %w", err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed with status: %d", receipt.Status)
	}

	fmt.Printf("✅ Proper test mining root uploaded successfully in block %d\n", receipt.BlockNumber.Uint64())
	fmt.Println()

	// Verify final state
	finalPeriod, err := distributor.CurrentPeriod(nil)
	if err != nil {
		return fmt.Errorf("failed to get final period: %w", err)
	}

	fmt.Printf("🎉 Contract advanced to period %d\n", finalPeriod.Uint64())
	fmt.Println()

	fmt.Printf("📋 Final Period Mapping:\n")
	fmt.Printf("   Settlement Period 1767549424 -> Contract Period 1\n")
	fmt.Printf("   Settlement Period 1767557168 -> Contract Period 2\n")
	fmt.Printf("   Settlement Period 1767650263 -> Contract Period 3\n")
	fmt.Printf("   Settlement Period 1768130418 -> Contract Period 4\n")
	fmt.Printf("   Settlement Period 1768135359 -> Contract Period 5 (empty proof)\n")
	fmt.Printf("   Settlement Period %d -> Contract Period %d (PROPER TEST)\n", testPeriodID, finalPeriod.Uint64())
	fmt.Println()

	fmt.Printf("✅ Proper test mining root uploaded successfully!\n")
	fmt.Printf("🎯 This settlement should work because:\n")
	fmt.Printf("   ✅ Correct user address in Merkle proof\n")
	fmt.Printf("   ✅ Valid Merkle proof elements (2 elements)\n")
	fmt.Printf("   ✅ Proper period mapping to contract\n")
	fmt.Printf("   ✅ Multiple leaves = proper proof validation\n")
	fmt.Println()
	fmt.Printf("📝 Next: Test claiming with period %d\n", testPeriodID)

	return nil
}

// Helper function for user confirmation
func confirm(prompt string) bool {
	fmt.Printf("%s (y/n): ", prompt)
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}
