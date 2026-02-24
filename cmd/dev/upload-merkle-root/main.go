package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/pkg/jarvis/binding"
	"github.com/kawai-network/veridium/pkg/jarvis/util/reader"
	"github.com/kawai-network/x/constant"
	"github.com/kawai-network/contracts"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <contract_type> <merkle_root>")
		fmt.Println("")
		fmt.Println("Contract Types:")
		fmt.Println("  mining   - Upload to MiningRewardDistributor")
		fmt.Println("  cashback - Upload to DepositCashbackDistributor")
		fmt.Println("  referral - Upload to ReferralRewardDistributor")
		fmt.Println("")
		fmt.Println("Example: go run main.go mining 0x1234...")
		os.Exit(1)
	}

	contractType := os.Args[1]
	merkleRoot := os.Args[2]

	// Validate merkle root format
	if len(merkleRoot) != 66 || merkleRoot[:2] != "0x" {
		log.Fatalf("Invalid merkle root format (expected 0x + 64 hex chars): %s", merkleRoot)
	}

	// Get private key from obfuscated constant
	privateKey := constant.GetAdminPrivateKey()

	// Remove 0x prefix if present
	if len(privateKey) > 2 && privateKey[:2] == "0x" {
		privateKey = privateKey[2:]
	}

	// Parse private key
	key, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Connect to Monad RPC
	client, err := ethclient.Dial(contracts.MonadRpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to Monad RPC: %v", err)
	}
	defer client.Close()

	// Create reader for contracts
	nodes := map[string]string{"monad": contracts.MonadRpcUrl}
	ethReader := reader.NewEthReaderGeneric(nodes, nil)

	// Get chain ID
	ctx := context.Background()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	// Create transact opts
	auth, err := bind.NewKeyedTransactorWithChainID(key, chainID)
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}

	// Set gas limit (optional, will be estimated if not set)
	auth.GasLimit = uint64(300000)

	rootBytes := common.HexToHash(merkleRoot)

	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("Contract Type: %s\n", contractType)
	fmt.Printf("Merkle Root: %s\n", merkleRoot)
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("")
	fmt.Println("🚀 Uploading Merkle root to contract...")

	var tx *types.Transaction
	var txHash string

	switch contractType {
	case "mining":
		distributor, err := binding.MiningRewardDistributor("MiningRewardDistributor", ethReader)
		if err != nil {
			log.Fatalf("Failed to load mining distributor: %v", err)
		}
		tx, err = distributor.SetMerkleRoot(auth, rootBytes)
		if err != nil {
			log.Fatalf("Failed to set merkle root: %v", err)
		}
		txHash = tx.Hash().Hex()

	case "cashback":
		distributor, err := binding.CashbackDistributor("CashbackDistributor", ethReader)
		if err != nil {
			log.Fatalf("Failed to load cashback distributor: %v", err)
		}
		tx, err = distributor.SetMerkleRoot(auth, rootBytes)
		if err != nil {
			log.Fatalf("Failed to set merkle root: %v", err)
		}
		txHash = tx.Hash().Hex()

	case "referral":
		distributor, err := binding.ReferralRewardDistributor("ReferralRewardDistributor", ethReader)
		if err != nil {
			log.Fatalf("Failed to load referral distributor: %v", err)
		}
		tx, err = distributor.SetMerkleRoot(auth, rootBytes)
		if err != nil {
			log.Fatalf("Failed to set merkle root: %v", err)
		}
		txHash = tx.Hash().Hex()

	default:
		log.Fatalf("Unknown contract type: %s (use: mining, cashback, or referral)", contractType)
	}

	fmt.Println("")
	fmt.Println("✅ Transaction submitted!")
	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("Transaction Hash: %s\n", txHash)
	fmt.Printf("Explorer: https://explorer.monad.xyz/tx/%s\n", txHash)
	fmt.Println("═══════════════════════════════════════")
	fmt.Println("")
	fmt.Println("⏳ Waiting for confirmation...")

	// Wait for transaction to be mined
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		log.Fatalf("Failed to wait for transaction: %v", err)
	}

	if receipt.Status == 1 {
		fmt.Println("")
		fmt.Println("✅ Merkle root uploaded successfully!")
		fmt.Println("═══════════════════════════════════════")
		fmt.Printf("Block Number: %s\n", receipt.BlockNumber.String())
		fmt.Printf("Gas Used: %d\n", receipt.GasUsed)
		fmt.Println("═══════════════════════════════════════")
	} else {
		fmt.Println("")
		fmt.Println("❌ Transaction reverted on-chain!")
		fmt.Println("═══════════════════════════════════════")
		fmt.Printf("Block Number: %s\n", receipt.BlockNumber.String())
		fmt.Printf("Gas Used: %d\n", receipt.GasUsed)
		fmt.Println("═══════════════════════════════════════")
		os.Exit(1)
	}
}
