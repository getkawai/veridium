package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/internal/generate/abi/miningdistributor"
	"github.com/kawai-network/veridium/pkg/store"
)

var (
	periodID     = flag.Int64("period", 0, "Period ID to upload (required)")
	rpcURL       = flag.String("rpc", "https://testnet-rpc.monad.xyz", "RPC URL")
	contractAddr = flag.String("contract", "0xa0dDC59DAcBA9201CC9Ef613707d287b77b2723F", "MiningRewardDistributor contract address")
	dryRun       = flag.Bool("dry-run", false, "Preview without uploading")
)

func main() {
	flag.Parse()

	// Load .env (optional)
	_ = godotenv.Load()

	if *periodID == 0 {
		fmt.Println("❌ Error: Period ID is required")
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("  go run cmd/upload-mining-root/main.go --period <PERIOD_ID>")
		fmt.Println("")
		fmt.Println("Example:")
		fmt.Println("  go run cmd/upload-mining-root/main.go --period 1767549424")
		fmt.Println("")
		fmt.Println("Options:")
		fmt.Println("  --period <ID>      Period ID to upload (required)")
		fmt.Println("  --rpc <URL>        RPC URL (default: Monad testnet)")
		fmt.Println("  --contract <ADDR>  Contract address (default: deployed address)")
		fmt.Println("  --dry-run          Preview without uploading")
		os.Exit(1)
	}

	ctx := context.Background()

	fmt.Println("🚀 Uploading Merkle Root to MiningRewardDistributor")
	fmt.Println("====================================================")
	fmt.Println("")

	// 1. Get settlement from KV store
	fmt.Println("📊 Step 1: Fetching settlement data from KV store...")
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatal("Failed to initialize KV store:", err)
	}

	settlement, err := kv.GetSettlementPeriod(ctx, *periodID)
	if err != nil {
		log.Fatal("Failed to get settlement:", err)
	}

	if settlement == nil {
		log.Fatalf("Settlement not found for period %d", *periodID)
	}

	fmt.Printf("  Period ID: %d\n", settlement.PeriodID)
	fmt.Printf("  Merkle Root: %s\n", settlement.MerkleRoot)
	fmt.Printf("  Contributors: %d\n", settlement.ContributorCount)
	fmt.Printf("  Total Amount: %s KAWAI\n", formatAmount(settlement.TotalAmount))
	fmt.Printf("  Status: %s\n", settlement.Status)
	fmt.Println("")

	if settlement.Status == store.SettlementStatusCompleted {
		// Settlement is completed, ready to upload
	} else {
		fmt.Printf("⚠️  Warning: Settlement status is '%s' (expected 'completed')\n", settlement.Status)
		fmt.Println("Continuing anyway...")
		fmt.Println("")
	}

	if *dryRun {
		fmt.Println("🔍 DRY RUN MODE - No transaction will be sent")
		fmt.Println("")
		fmt.Println("Would upload:")
		fmt.Printf("  Period: %d\n", settlement.PeriodID)
		fmt.Printf("  Root: %s\n", settlement.MerkleRoot)
		fmt.Println("")
		fmt.Println("To actually upload, run without --dry-run")
		return
	}

	// 2. Connect to blockchain
	fmt.Println("🔗 Step 2: Connecting to blockchain...")
	client, err := ethclient.Dial(*rpcURL)
	if err != nil {
		log.Fatal("Failed to connect to blockchain:", err)
	}
	defer client.Close()

	chainID, err := client.ChainID(ctx)
	if err != nil {
		log.Fatal("Failed to get chain ID:", err)
	}
	fmt.Printf("  Connected to chain ID: %s\n", chainID.String())
	fmt.Println("")

	// 3. Load private key
	fmt.Println("🔐 Step 3: Loading private key...")
	privateKeyHex := constant.GetObfuscatedTemp()
	if privateKeyHex == "" {
		log.Fatal("Private key not found in constant.GetObfuscatedTemp()")
	}

	// Remove 0x prefix if present
	if len(privateKeyHex) >= 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatal("Failed to parse private key:", err)
	}

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Failed to cast public key to ECDSA")
	}

	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	fmt.Printf("  Sender address: %s\n", fromAddress.Hex())
	fmt.Println("")

	// 4. Load contract
	fmt.Println("📜 Step 4: Loading MiningRewardDistributor contract...")
	contractAddress := common.HexToAddress(*contractAddr)
	contract, err := miningdistributor.NewMiningRewardDistributor(contractAddress, client)
	if err != nil {
		log.Fatal("Failed to load contract:", err)
	}
	fmt.Printf("  Contract address: %s\n", contractAddress.Hex())
	fmt.Println("")

	// 5. Get transaction options
	fmt.Println("⚙️  Step 5: Preparing transaction...")
	nonce, err := client.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		log.Fatal("Failed to get nonce:", err)
	}

	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		log.Fatal("Failed to get gas price:", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatal("Failed to create transactor:", err)
	}

	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(300000) // Set gas limit
	auth.GasPrice = gasPrice

	fmt.Printf("  Nonce: %d\n", nonce)
	fmt.Printf("  Gas Price: %s gwei\n", new(big.Int).Div(gasPrice, big.NewInt(1e9)).String())
	fmt.Printf("  Gas Limit: %d\n", auth.GasLimit)
	fmt.Println("")

	// 6. Parse Merkle root (convert to [32]byte)
	merkleRootHash := common.HexToHash(settlement.MerkleRoot)
	var merkleRoot [32]byte
	copy(merkleRoot[:], merkleRootHash[:])

	// 7. Send transaction
	fmt.Println("📤 Step 6: Sending transaction...")
	fmt.Printf("  setMerkleRoot(%s)\n", settlement.MerkleRoot)
	fmt.Printf("  Note: Contract uses currentPeriod internally\n")
	fmt.Println("")

	tx, err := contract.SetMerkleRoot(auth, merkleRoot)
	if err != nil {
		log.Fatal("Failed to send transaction:", err)
	}

	fmt.Printf("✅ Transaction sent!\n")
	fmt.Printf("  TX Hash: %s\n", tx.Hash().Hex())
	fmt.Println("")

	// 8. Wait for confirmation
	fmt.Println("⏳ Waiting for confirmation...")
	receipt, err := bind.WaitMined(ctx, client, tx)
	if err != nil {
		log.Fatal("Failed to wait for transaction:", err)
	}

	if receipt.Status == 1 {
		fmt.Println("✅ Transaction confirmed!")
		fmt.Printf("  Block: %d\n", receipt.BlockNumber.Uint64())
		fmt.Printf("  Gas Used: %d\n", receipt.GasUsed)
		fmt.Println("")

		// 9. Log success (settlement status already 'completed')
		fmt.Println("💾 Step 7: Upload complete")
		fmt.Printf("  TX Hash: %s\n", tx.Hash().Hex())
		fmt.Println("")

		fmt.Println("====================================================")
		fmt.Println("🎉 Merkle Root Upload Complete!")
		fmt.Println("")
		fmt.Println("Next steps:")
		fmt.Println("  1. Contributors can now claim rewards via frontend")
		fmt.Println("  2. Test claim flow in UI")
		fmt.Println("")
		fmt.Printf("View on explorer: https://testnet.monad.xyz/tx/%s\n", tx.Hash().Hex())
	} else {
		fmt.Println("❌ Transaction failed!")
		fmt.Printf("  TX Hash: %s\n", tx.Hash().Hex())
		os.Exit(1)
	}
}

func formatAmount(amountStr string) string {
	if amountStr == "" || amountStr == "0" {
		return "0"
	}

	// Convert wei to KAWAI (divide by 1e18)
	amount := new(big.Int)
	amount.SetString(amountStr, 10)

	// Divide by 1e18
	divisor := new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)
	result := new(big.Int).Div(amount, divisor)

	return result.String()
}
