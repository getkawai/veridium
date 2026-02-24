package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/contracts/cashbackdistributor"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/contracts"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run cmd/dev/test-cashback-claim/main.go <user_address> <private_key>")
		os.Exit(1)
	}

	userAddress := os.Args[1]
	privateKeyHex := os.Args[2]

	// Strip 0x prefix
	if strings.HasPrefix(privateKeyHex, "0x") {
		privateKeyHex = privateKeyHex[2:]
	}

	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatal("Failed to initialize KV store:", err)
	}

	ctx := context.Background()

	fmt.Println("🧪 Testing Cashback Claim")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("User: %s\n", userAddress)
	fmt.Println()

	// 1. Get proof from KV
	period := uint64(1)
	proofKey := fmt.Sprintf("cashback_proof:%d:%s", period, userAddress)

	data, err := kv.GetCashbackData(ctx, proofKey)
	if err != nil {
		log.Fatalf("Failed to get proof from KV: %v", err)
	}

	var proofRecord struct {
		Proof   []string `json:"proof"`
		Amount  string   `json:"amount"`
		Claimed bool     `json:"claimed"`
	}
	if err := json.Unmarshal(data, &proofRecord); err != nil {
		log.Fatalf("Failed to unmarshal proof: %v", err)
	}

	fmt.Println("📦 Proof Data:")
	fmt.Printf("   Period: %d\n", period)
	fmt.Printf("   Amount: %s wei\n", proofRecord.Amount)
	fmt.Printf("   Proof length: %d\n", len(proofRecord.Proof))
	fmt.Println()

	// 2. Parse amount
	amount := new(big.Int)
	amount.SetString(proofRecord.Amount, 10)

	// 3. Convert proof to [][32]byte
	merkleProof := make([][32]byte, len(proofRecord.Proof))
	for i, p := range proofRecord.Proof {
		if strings.HasPrefix(p, "0x") {
			p = p[2:]
		}
		proofBytes := common.Hex2Bytes(p)
		if len(proofBytes) != 32 {
			log.Fatalf("Invalid proof element at index %d: expected 32 bytes, got %d", i, len(proofBytes))
		}
		copy(merkleProof[i][:], proofBytes)
	}

	// 4. Connect to blockchain
	client, err := ethclient.Dial(contracts.MonadRpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to blockchain: %v", err)
	}
	defer client.Close()

	// 5. Load contract
	distributorAddr := common.HexToAddress(contracts.CashbackDistributorAddress)
	distributor, err := cashbackdistributor.NewDepositCashbackDistributor(distributorAddr, client)
	if err != nil {
		log.Fatalf("Failed to load distributor contract: %v", err)
	}

	// 6. Parse private key
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// 7. Get chain ID
	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	// 8. Create transactor
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
	}
	auth.Context = ctx

	fmt.Println("📝 Transaction Parameters:")
	fmt.Printf("   Period: %d\n", period)
	fmt.Printf("   Amount: %s wei\n", amount.String())
	fmt.Printf("   Proof length: %d\n", len(merkleProof))
	fmt.Println()

	// 9. Submit transaction
	fmt.Println("⏳ Submitting claim transaction...")
	tx, err := distributor.ClaimCashback(auth, new(big.Int).SetUint64(period), amount, merkleProof)
	if err != nil {
		log.Fatalf("❌ Transaction failed: %v", err)
	}

	fmt.Printf("✅ Transaction submitted: %s\n", tx.Hash().Hex())
	fmt.Println()

	// 10. Wait for confirmation
	fmt.Println("⏳ Waiting for confirmation...")
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		log.Fatalf("Failed to wait for confirmation: %v", err)
	}

	if receipt.Status != 1 {
		log.Fatalf("❌ Transaction reverted (status: %d)", receipt.Status)
	}

	fmt.Printf("✅ Transaction confirmed in block %d\n", receipt.BlockNumber.Uint64())
	fmt.Printf("   Gas used: %d\n", receipt.GasUsed)
	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("✅ Cashback claim successful!")
}
