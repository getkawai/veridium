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
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/internal/generate/abi/kawaitoken"
)

func main() {
	ctx := context.Background()

	fmt.Println("🪙 Minting Test KAWAI Tokens")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// Connect to Monad
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to Monad: %v", err)
	}
	fmt.Printf("✓ Connected to Monad RPC\n")
	fmt.Printf("  %s\n\n", constant.MonadRpcUrl)

	// Load KawaiToken contract
	tokenAddr := common.HexToAddress(constant.KawaiTokenAddress)
	token, err := kawaitoken.NewKawaiToken(tokenAddr, client)
	if err != nil {
		log.Fatalf("Failed to load KawaiToken: %v", err)
	}
	fmt.Printf("✓ Loaded KawaiToken\n")
	fmt.Printf("  %s\n\n", tokenAddr.Hex())

	// Get private key
	privateKeyHex := constant.GetAdminPrivateKey()
	// Remove 0x prefix if present
	if len(privateKeyHex) >= 2 && privateKeyHex[:2] == "0x" {
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
	auth.Context = ctx

	// Test addresses to mint to
	testAddresses := []struct {
		Address common.Address
		Amount  *big.Int
		Label   string
	}{
		{
			Address: common.HexToAddress(constant.GetAdminAddress()),
			Amount:  new(big.Int).Mul(big.NewInt(10000), big.NewInt(1e18)), // 10,000 KAWAI
			Label:   "Admin/Test Address 1",
		},
		{
			Address: common.HexToAddress("0x0f3e75B9Bb3efcD87B1Ed15a30C8a7FBaABD204E"),
			Amount:  new(big.Int).Mul(big.NewInt(5000), big.NewInt(1e18)), // 5,000 KAWAI
			Label:   "Test Address 2",
		},
		{
			Address: common.HexToAddress("0x9f152652004F133f64522ECE18D3Dc0eD531d2d7"),
			Amount:  new(big.Int).Mul(big.NewInt(3000), big.NewInt(1e18)), // 3,000 KAWAI
			Label:   "Test Address 3",
		},
	}

	fmt.Println("🎯 Minting to test addresses:")
	fmt.Println()

	totalMinted := big.NewInt(0)

	for i, test := range testAddresses {
		fmt.Printf("%d. %s\n", i+1, test.Label)
		fmt.Printf("   Address: %s\n", test.Address.Hex())
		fmt.Printf("   Amount:  %s KAWAI\n", formatKawai(test.Amount))

		// Mint tokens
		tx, err := token.Mint(auth, test.Address, test.Amount)
		if err != nil {
			log.Printf("   ❌ Failed to mint: %v\n\n", err)
			continue
		}

		fmt.Printf("   TX Hash: %s\n", tx.Hash().Hex())

		// Wait for confirmation
		receipt, err := bind.WaitMined(ctx, client, tx)
		if err != nil {
			log.Printf("   ❌ Failed to confirm: %v\n\n", err)
			continue
		}

		if receipt.Status != 1 {
			log.Printf("   ❌ Transaction failed\n\n")
			continue
		}

		fmt.Printf("   ✅ Minted successfully (Block: %d)\n\n", receipt.BlockNumber.Uint64())
		totalMinted.Add(totalMinted, test.Amount)
	}

	// Check total supply
	totalSupply, err := token.TotalSupply(nil)
	if err != nil {
		log.Fatalf("Failed to get total supply: %v", err)
	}

	fmt.Println("═══════════════════════════════════════")
	fmt.Printf("✅ Minting Complete!\n\n")
	fmt.Printf("Total Minted:  %s KAWAI\n", formatKawai(totalMinted))
	fmt.Printf("Total Supply:  %s KAWAI\n", formatKawai(totalSupply))
	fmt.Println("═══════════════════════════════════════")

	os.Exit(0)
}

func formatKawai(amount *big.Int) string {
	// Convert wei to KAWAI (18 decimals)
	kawai := new(big.Float).SetInt(amount)
	kawai.Quo(kawai, big.NewFloat(1e18))
	return kawai.Text('f', 4)
}
