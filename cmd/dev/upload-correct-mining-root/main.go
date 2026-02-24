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
	"github.com/kawai-network/contracts/miningdistributor"
	"github.com/kawai-network/x/constant"
	"github.com/kawai-network/contracts"
)

func main() {
	if err := uploadCorrectMiningRoot(); err != nil {
		log.Fatalf("Failed to upload correct mining root: %v", err)
	}
}

func uploadCorrectMiningRoot() error {
	ctx := context.Background()

	fmt.Println("🚀 Uploading CORRECT Mining Root (Period 7)")
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println()

	// CORRECT settlement details (generated with matching periods)
	correctMerkleRoot := "0xff77aeb0d8b803ac73709f80e3aab1ec566dbab8a0a3ea8182242062cf3ee19e"
	settlementID := int64(1768137242)
	contractPeriod := int64(7)

	fmt.Printf("📋 CORRECT Settlement Details:\n")
	fmt.Printf("   Settlement ID:     %d (KV storage key)\n", settlementID)
	fmt.Printf("   Contract Period:   %d (sequential period)\n", contractPeriod)
	fmt.Printf("   Merkle Root:       %s\n", correctMerkleRoot)
	fmt.Printf("   User Address:      0x0f3e75B9Bb3efcD87B1Ed15a30C8a7FBaABD204E\n")
	fmt.Println()

	fmt.Printf("🎯 Key Fix:\n")
	fmt.Printf("   ✅ Merkle leaves generated using CONTRACT PERIOD: %d\n", contractPeriod)
	fmt.Printf("   ✅ Contract will validate using same period: %d\n", contractPeriod)
	fmt.Printf("   ✅ Period mismatch issue resolved!\n")
	fmt.Println()

	// Connect to Monad RPC
	client, err := ethclient.Dial(contracts.MonadRpcUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to Monad: %w", err)
	}
	defer client.Close()

	// Load MiningRewardDistributor contract
	distributorAddr := common.HexToAddress(contracts.MiningRewardDistributorAddress)
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

	if currentPeriod.Uint64()+1 != uint64(contractPeriod) {
		return fmt.Errorf("contract period mismatch: expected to advance to %d, but would advance to %d", contractPeriod, currentPeriod.Uint64()+1)
	}
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
	merkleRootHex := correctMerkleRoot
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
	fmt.Printf("⚠️  About to advance contract to period %d with CORRECT Merkle root\n", contractPeriod)
	fmt.Printf("🌳 Root: %s\n", correctMerkleRoot)
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

	fmt.Printf("✅ CORRECT mining root uploaded successfully in block %d\n", receipt.BlockNumber.Uint64())
	fmt.Println()

	// Verify final state
	finalPeriod, err := distributor.CurrentPeriod(nil)
	if err != nil {
		return fmt.Errorf("failed to get final period: %w", err)
	}

	fmt.Printf("🎉 Contract advanced to period %d\n", finalPeriod.Uint64())
	fmt.Println()

	fmt.Printf("📋 Updated Period Mapping:\n")
	fmt.Printf("   Settlement Period 1767549424 -> Contract Period 1\n")
	fmt.Printf("   Settlement Period 1767557168 -> Contract Period 2\n")
	fmt.Printf("   Settlement Period 1767650263 -> Contract Period 3\n")
	fmt.Printf("   Settlement Period 1768130418 -> Contract Period 4\n")
	fmt.Printf("   Settlement Period 1768135359 -> Contract Period 5\n")
	fmt.Printf("   Settlement Period 1768136095 -> Contract Period 6\n")
	fmt.Printf("   Settlement Period %d -> Contract Period %d (CORRECT)\n", settlementID, finalPeriod.Uint64())
	fmt.Println()

	fmt.Printf("✅ CORRECT mining root uploaded successfully!\n")
	fmt.Printf("🎯 Period mismatch issue resolved - Merkle proofs should now work!\n")
	fmt.Printf("📝 Next: Test claiming with user address 0x0f3e75B9Bb3efcD87B1Ed15a30C8a7FBaABD204E\n")

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
