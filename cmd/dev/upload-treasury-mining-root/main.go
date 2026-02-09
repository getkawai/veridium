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
	"github.com/kawai-network/veridium/internal/generate/abi/miningdistributor"
	"github.com/kawai-network/x/constant"
)

func main() {
	if err := uploadTreasuryMiningRoot(); err != nil {
		log.Fatalf("Failed to upload treasury mining root: %v", err)
	}
}

func uploadTreasuryMiningRoot() error {
	ctx := context.Background()

	fmt.Println("🏦 Uploading Treasury Mining Root (Period 8)")
	fmt.Println("═══════════════════════════════════════════════")
	fmt.Println()

	// Treasury settlement details (generated with matching msg.sender)
	treasuryMerkleRoot := "0xe5edb46d8eaecf74f3de52811cacdb920eaeebd38a68ac0726fab70c5d5bf2e8"
	settlementID := int64(1768141059)
	contractPeriod := int64(8) // Advance to period 8

	fmt.Printf("📋 Treasury Settlement Details:\n")
	fmt.Printf("   Settlement ID:     %d (KV storage key)\n", settlementID)
	fmt.Printf("   Contract Period:   %d (sequential period)\n", contractPeriod)
	fmt.Printf("   Merkle Root:       %s\n", treasuryMerkleRoot)
	fmt.Println()

	fmt.Printf("🎯 Key Fix:\n")
	fmt.Printf("   ✅ Merkle leaves generated using CONTRACT PERIOD: %d\n", contractPeriod-1)
	fmt.Printf("   ✅ Treasury address as contributor matches msg.sender\n")
	fmt.Printf("   ✅ Contract will validate using same period: %d\n", contractPeriod-1)
	fmt.Printf("   ✅ Address mismatch issue resolved!\n")
	fmt.Println()

	// Connect to Monad RPC
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to Monad: %w", err)
	}
	defer client.Close()

	// Load MiningRewardDistributor contract
	distributorAddr := common.HexToAddress(constant.MiningRewardDistributorAddress)
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
	merkleRootHex := treasuryMerkleRoot
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
	fmt.Printf("⚠️  About to advance contract to period %d with Treasury Merkle root\n", contractPeriod)
	fmt.Printf("🌳 Root: %s\n", treasuryMerkleRoot)
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

	fmt.Printf("✅ Treasury mining root uploaded successfully in block %d\n", receipt.BlockNumber.Uint64())
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
	fmt.Printf("   Settlement Period 1768137242 -> Contract Period 7\n")
	fmt.Printf("   Settlement Period %d -> Contract Period %d (TREASURY)\n", settlementID, finalPeriod.Uint64())
	fmt.Println()

	fmt.Printf("✅ Treasury mining root uploaded successfully!\n")
	fmt.Printf("🎯 Address mismatch issue resolved - Treasury can now claim!\n")
	fmt.Printf("📝 Next: Test claiming with treasury address as msg.sender\n")

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
