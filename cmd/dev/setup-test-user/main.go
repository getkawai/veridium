package main

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/common"
	ethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/y/types"
	"github.com/kawai-network/contracts"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found")
	}

	fmt.Println("🚀 Setting up test user (generate wallet + send MON + inject mining data)")
	fmt.Println("═══════════════════════════════════════════════════════════════════════")
	fmt.Println()

	// STEP 1: Generate wallet
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		log.Fatal("Failed to generate private key:", err)
	}

	privateKeyBytes := crypto.FromECDSA(privateKey)
	address := crypto.PubkeyToAddress(privateKey.PublicKey)

	fmt.Println("✅ STEP 1: Wallet generated")
	fmt.Printf("   Address:     %s\n", address.Hex())
	fmt.Printf("   Private Key: 0x%x\n", privateKeyBytes)
	fmt.Println()

	// STEP 2: Send MON for gas
	fmt.Println("⏳ STEP 2: Sending 0.1 MON for gas...")
	if err := sendMON(address); err != nil {
		log.Fatal("Failed to send MON:", err)
	}
	fmt.Println("✅ STEP 2: MON sent successfully")
	fmt.Println()

	// STEP 3: Inject mining data
	fmt.Println("⏳ STEP 3: Injecting mining data (450 KAWAI)...")
	if err := injectMiningData(address.Hex()); err != nil {
		log.Fatal("Failed to inject mining data:", err)
	}
	fmt.Println("✅ STEP 3: Mining data injected")
	fmt.Println()

	// Summary
	fmt.Println("═══════════════════════════════════════════════════════════════════════")
	fmt.Println("✅ TEST USER SETUP COMPLETE!")
	fmt.Println()
	fmt.Println("📝 Test User Details:")
	fmt.Printf("   Address:     %s\n", address.Hex())
	fmt.Printf("   Private Key: 0x%x\n", privateKeyBytes)
	fmt.Println()
	fmt.Println("📝 Next Steps:")
	fmt.Println("   1. Run settlement: make settle-mining")
	fmt.Println("   2. Upload Merkle root: go run cmd/reward-settlement/main.go upload --type mining")
	fmt.Println("   3. Import private key to UI and test claiming")
	fmt.Println()
}

func sendMON(toAddress common.Address) error {
	// Connect to RPC
	client, err := ethclient.Dial(contracts.MonadRpcUrl)
	if err != nil {
		return fmt.Errorf("failed to connect to RPC: %w", err)
	}
	defer client.Close()

	// Load sender private key from env
	privateKeyHex := os.Getenv("PRIVATE_KEY")
	if privateKeyHex == "" {
		// Try loading from contracts/.env
		if err := godotenv.Load("contracts/.env"); err == nil {
			privateKeyHex = os.Getenv("PRIVATE_KEY")
		}
	}
	if privateKeyHex == "" {
		return fmt.Errorf("PRIVATE_KEY not set in .env or contracts/.env")
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex[2:]) // Remove 0x prefix
	if err != nil {
		return fmt.Errorf("invalid private key: %w", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("error casting public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)

	// Get nonce
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		return fmt.Errorf("failed to get nonce: %w", err)
	}

	// Amount: 0.1 MON
	amount := new(big.Int)
	amount.SetString("100000000000000000", 10) // 0.1 ETH in wei

	// Gas settings
	gasLimit := uint64(21000)
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get gas price: %w", err)
	}

	// Create transaction
	tx := ethtypes.NewTransaction(nonce, toAddress, amount, gasLimit, gasPrice, nil)

	// Get chain ID
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		return fmt.Errorf("failed to get chain ID: %w", err)
	}

	// Sign transaction
	signedTx, err := ethtypes.SignTx(tx, ethtypes.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Send transaction
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %w", err)
	}

	// Wait for confirmation
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	receipt, err := waitForReceipt(ctx, client, signedTx.Hash())
	if err != nil {
		return fmt.Errorf("failed to wait for receipt: %w", err)
	}

	if receipt.Status != 1 {
		return fmt.Errorf("transaction failed")
	}

	return nil
}

func waitForReceipt(ctx context.Context, client *ethclient.Client, txHash common.Hash) (*ethtypes.Receipt, error) {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-ticker.C:
			receipt, err := client.TransactionReceipt(ctx, txHash)
			if err == nil {
				return receipt, nil
			}
		}
	}
}

func injectMiningData(contributorAddress string) error {
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		return fmt.Errorf("failed to initialize KV store: %w", err)
	}

	ctx := context.Background()

	// Job 1: 100 KAWAI
	record1 := &store.JobRewardRecord{
		Timestamp:          time.Now(),
		ContributorAddress: contributorAddress,
		UserAddress:        "0x1111111111111111111111111111111111111111",
		ReferrerAddress:    "",
		DeveloperAddress:   contributorAddress,
		ContributorAmount:  "100000000000000000000",
		DeveloperAmount:    "5000000000000000000",
		UserAmount:         "5000000000000000000",
		AffiliatorAmount:   "0",
		TokenUsage:         1000,
		RewardType:         types.RewardTypeMining,
		HasReferrer:        false,
		IsSettled:          false,
	}
	if err := kv.SaveJobReward(ctx, record1); err != nil {
		return fmt.Errorf("failed to save job 1: %w", err)
	}

	time.Sleep(1 * time.Second)

	// Job 2: 150 KAWAI
	record2 := &store.JobRewardRecord{
		Timestamp:          time.Now(),
		ContributorAddress: contributorAddress,
		UserAddress:        "0x2222222222222222222222222222222222222222",
		ReferrerAddress:    "",
		DeveloperAddress:   contributorAddress,
		ContributorAmount:  "150000000000000000000",
		DeveloperAmount:    "5000000000000000000",
		UserAmount:         "5000000000000000000",
		AffiliatorAmount:   "0",
		TokenUsage:         1500,
		RewardType:         types.RewardTypeMining,
		HasReferrer:        false,
		IsSettled:          false,
	}
	if err := kv.SaveJobReward(ctx, record2); err != nil {
		return fmt.Errorf("failed to save job 2: %w", err)
	}

	time.Sleep(1 * time.Second)

	// Job 3: 200 KAWAI
	record3 := &store.JobRewardRecord{
		Timestamp:          time.Now(),
		ContributorAddress: contributorAddress,
		UserAddress:        "0x3333333333333333333333333333333333333333",
		ReferrerAddress:    "",
		DeveloperAddress:   contributorAddress,
		ContributorAmount:  "200000000000000000000",
		DeveloperAmount:    "5000000000000000000",
		UserAmount:         "5000000000000000000",
		AffiliatorAmount:   "0",
		TokenUsage:         2000,
		RewardType:         types.RewardTypeMining,
		HasReferrer:        false,
		IsSettled:          false,
	}
	if err := kv.SaveJobReward(ctx, record3); err != nil {
		return fmt.Errorf("failed to save job 3: %w", err)
	}

	return nil
}
