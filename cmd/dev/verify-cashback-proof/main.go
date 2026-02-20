package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/kawai-network/contracts/cashbackdistributor"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/x/constant"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/dev/verify-cashback-proof/main.go <user_address>")
		os.Exit(1)
	}

	userAddress := os.Args[1]

	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatal("Failed to initialize KV store:", err)
	}

	ctx := context.Background()

	fmt.Println("🔍 Verifying Cashback Proof")
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

	fmt.Println("📦 Proof from KV:")
	fmt.Printf("   Period: %d\n", period)
	fmt.Printf("   Amount: %s wei\n", proofRecord.Amount)
	fmt.Printf("   Proof length: %d\n", len(proofRecord.Proof))
	fmt.Printf("   Claimed: %v\n", proofRecord.Claimed)
	fmt.Println()

	// 2. Get merkle root from KV
	rootKey := fmt.Sprintf("cashback_period:%d:merkle_root", period)
	rootData, err := kv.GetCashbackData(ctx, rootKey)
	if err != nil {
		log.Fatalf("Failed to get merkle root from KV: %v", err)
	}

	var merkleRootHex string
	if err := json.Unmarshal(rootData, &merkleRootHex); err != nil {
		log.Fatalf("Failed to unmarshal merkle root: %v", err)
	}

	fmt.Println("🌳 Merkle Root from KV:")
	fmt.Printf("   %s\n", merkleRootHex)
	fmt.Println()

	// 3. Get merkle root from contract
	client, err := ethclient.Dial(constant.MonadRpcUrl)
	if err != nil {
		log.Fatalf("Failed to connect to blockchain: %v", err)
	}
	defer client.Close()

	distributorAddr := common.HexToAddress(constant.CashbackDistributorAddress)
	distributor, err := cashbackdistributor.NewDepositCashbackDistributor(distributorAddr, client)
	if err != nil {
		log.Fatalf("Failed to load distributor contract: %v", err)
	}

	onChainRoot, err := distributor.PeriodMerkleRoots(nil, new(big.Int).SetUint64(period))
	if err != nil {
		log.Fatalf("Failed to get on-chain merkle root: %v", err)
	}

	fmt.Println("🌳 Merkle Root from Contract:")
	fmt.Printf("   0x%x\n", onChainRoot)
	fmt.Println()

	// 4. Compare roots
	if merkleRootHex != fmt.Sprintf("0x%x", onChainRoot) {
		fmt.Println("❌ MISMATCH: KV root != Contract root")
		fmt.Printf("   KV:       %s\n", merkleRootHex)
		fmt.Printf("   Contract: 0x%x\n", onChainRoot)
		fmt.Println()
		fmt.Println("This is the problem! Settlement stored wrong root.")
		os.Exit(1)
	}

	fmt.Println("✅ Merkle roots match!")
	fmt.Println()

	// 5. Reconstruct leaf hash
	amount := new(big.Int)
	amount.SetString(proofRecord.Amount, 10)

	// Hash: keccak256(abi.encodePacked(period, user, amount))
	periodBytes := new(big.Int).SetUint64(period).Bytes()
	addressBytes := common.HexToAddress(userAddress).Bytes()
	amountBytes := amount.Bytes()

	leafHash := crypto.Keccak256(periodBytes, addressBytes, amountBytes)

	fmt.Println("🍃 Leaf Hash (abi.encodePacked):")
	fmt.Printf("   Period bytes: 0x%x (%d bytes)\n", periodBytes, len(periodBytes))
	fmt.Printf("   Address bytes: 0x%x (%d bytes)\n", addressBytes, len(addressBytes))
	fmt.Printf("   Amount bytes: 0x%x (%d bytes)\n", amountBytes, len(amountBytes))
	fmt.Printf("   Leaf hash: 0x%x\n", leafHash)
	fmt.Println()

	// 6. For single-leaf tree, leaf hash should equal root
	if len(proofRecord.Proof) == 0 {
		fmt.Println("📝 Single-leaf tree (empty proof)")
		fmt.Println()

		if fmt.Sprintf("0x%x", leafHash) == merkleRootHex {
			fmt.Println("✅ Leaf hash matches merkle root (correct for single-leaf tree)")
		} else {
			fmt.Println("❌ MISMATCH: Leaf hash != Merkle root")
			fmt.Printf("   Leaf: 0x%x\n", leafHash)
			fmt.Printf("   Root: %s\n", merkleRootHex)
			fmt.Println()
			fmt.Println("This is the problem! Leaf encoding is wrong.")
			os.Exit(1)
		}
	}

	fmt.Println()
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Println("✅ All checks passed! Proof should be valid.")
	fmt.Println()
	fmt.Println("If claiming still fails, the issue is likely:")
	fmt.Println("  1. Frontend sending wrong parameters")
	fmt.Println("  2. Contract expecting different encoding")
	fmt.Println("  3. Period mismatch")
}
