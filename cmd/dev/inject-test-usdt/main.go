package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/veridium/internal/generate/abi/mockstablecoin"
	"github.com/kawai-network/veridium/pkg/config"
	"github.com/kawai-network/x/constant"
)

func main() {
	// Safety check: Only allow on testnet
	if !config.IsTestnet() {
		log.Fatal("❌ ERROR: This tool is only available on testnet. On mainnet, you must acquire USDC through exchanges or bridges.")
	}

	ctx := context.Background()

	// Connect to Monad testnet
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to Monad: %v", err)
	}
	defer client.Close()

	// Load stablecoin contract (MockStablecoin on testnet)
	stablecoinAddr := common.HexToAddress(constant.StablecoinAddress)
	stablecoinContract, err := mockstablecoin.NewMockStablecoin(stablecoinAddr, client)
	if err != nil {
		log.Fatalf("Failed to load stablecoin contract: %v", err)
	}

	// Get private key from temp.go
	privateKeyHex := constant.GetAdminPrivateKey()
	// Strip 0x prefix if present
	if len(privateKeyHex) > 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Get chain ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	// Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}

	// Amount to inject: 1000 stablecoin (for testing)
	amount := new(big.Int)
	amount.SetString("1000000000", 10) // 1000 stablecoin (6 decimals)

	paymentVault := common.HexToAddress(constant.PaymentVaultAddress)

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("💵 Injecting Test Stablecoin to PaymentVault")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Printf("From:   %s\n", crypto.PubkeyToAddress(privateKey.PublicKey).Hex())
	fmt.Printf("To:     %s (PaymentVault)\n", paymentVault.Hex())
	fmt.Printf("Amount: 1000 stablecoin (MockStablecoin on testnet)\n")
	fmt.Println()
	fmt.Print("Continue? (y/n): ")

	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "Y" {
		fmt.Println("❌ Cancelled")
		os.Exit(0)
	}

	// Transfer stablecoin to PaymentVault
	tx, err := stablecoinContract.Transfer(auth, paymentVault, amount)
	if err != nil {
		log.Fatalf("Failed to transfer stablecoin: %v", err)
	}

	fmt.Println()
	fmt.Println("⏳ Waiting for transaction confirmation...")

	// Wait for transaction
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		log.Fatalf("Failed to wait for transaction: %v", err)
	}

	if receipt.Status == 0 {
		log.Fatalf("Transaction failed")
	}

	fmt.Println()
	fmt.Println("✅ Test stablecoin injected successfully!")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("Transaction Hash: %s\n", tx.Hash().Hex())
	fmt.Printf("Block Number:     %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("Gas Used:         %d\n", receipt.GasUsed)
	fmt.Printf("Explorer:         https://explorer.monad.xyz/tx/%s\n", tx.Hash().Hex())
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println()
	fmt.Println("✅ PaymentVault now has 1000 stablecoin for testing")
	fmt.Println("✅ Ready to run: make settle-revenue")
}
