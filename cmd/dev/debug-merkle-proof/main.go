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
	"github.com/kawai-network/veridium/pkg/store"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run cmd/dev/debug-merkle-proof/main.go <user_address>")
		os.Exit(1)
	}

	userAddress := os.Args[1]

	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		log.Fatal("Failed to initialize KV store:", err)
	}

	ctx := context.Background()

	fmt.Println("🔍 Debugging Merkle Proof")
	fmt.Println("═══════════════════════════════════════════════════════════════")
	fmt.Printf("User: %s\n", userAddress)
	fmt.Println()

	// Get proof from KV
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

	amount := new(big.Int)
	amount.SetString(proofRecord.Amount, 10)

	// Calculate leaf hash (with 32-byte padding)
	periodBytes := common.LeftPadBytes(new(big.Int).SetUint64(period).Bytes(), 32)
	addressBytes := common.HexToAddress(userAddress).Bytes()
	amountBytes := common.LeftPadBytes(amount.Bytes(), 32)

	leafHash := crypto.Keccak256(periodBytes, addressBytes, amountBytes)

	fmt.Println("📦 Proof Data:")
	fmt.Printf("   Period: %d\n", period)
	fmt.Printf("   Amount: %s wei\n", proofRecord.Amount)
	fmt.Printf("   Proof length: %d\n", len(proofRecord.Proof))
	fmt.Println()

	fmt.Println("🍃 Leaf Hash Calculation:")
	fmt.Printf("   Period (32 bytes): 0x%x\n", periodBytes)
	fmt.Printf("   Address (20 bytes): 0x%x\n", addressBytes)
	fmt.Printf("   Amount (32 bytes): 0x%x\n", amountBytes)
	fmt.Printf("   Leaf Hash: 0x%x\n", leafHash)
	fmt.Println()

	// Verify proof by reconstructing root
	currentHash := leafHash
	fmt.Println("🌳 Merkle Proof Verification:")
	fmt.Printf("   Starting with leaf: 0x%x\n", currentHash)

	for i, proofHex := range proofRecord.Proof {
		proofBytes := common.Hex2Bytes(proofHex[2:]) // Strip 0x

		fmt.Printf("\n   Step %d:\n", i+1)
		fmt.Printf("     Sibling: %s\n", proofHex)

		// Sort before hashing (OpenZeppelin requirement)
		var combined []byte
		if string(currentHash) < string(proofBytes) {
			combined = crypto.Keccak256(currentHash, proofBytes)
			fmt.Printf("     Hash(current, sibling)\n")
		} else {
			combined = crypto.Keccak256(proofBytes, currentHash)
			fmt.Printf("     Hash(sibling, current)\n")
		}

		fmt.Printf("     Result: 0x%x\n", combined)
		currentHash = combined
	}

	fmt.Println()
	fmt.Printf("   Final Root: 0x%x\n", currentHash)
	fmt.Println()

	// Get expected root from KV
	rootKey := fmt.Sprintf("cashback_period:%d:merkle_root", period)
	rootData, err := kv.GetCashbackData(ctx, rootKey)
	if err != nil {
		log.Fatalf("Failed to get merkle root from KV: %v", err)
	}

	var merkleRootHex string
	if err := json.Unmarshal(rootData, &merkleRootHex); err != nil {
		log.Fatalf("Failed to unmarshal merkle root: %v", err)
	}

	fmt.Println("🎯 Comparison:")
	fmt.Printf("   Calculated Root: 0x%x\n", currentHash)
	fmt.Printf("   Expected Root:   %s\n", merkleRootHex)
	fmt.Println()

	if fmt.Sprintf("0x%x", currentHash) == merkleRootHex {
		fmt.Println("✅ Proof is VALID!")
	} else {
		fmt.Println("❌ Proof is INVALID!")
		fmt.Println()
		fmt.Println("This means the proof stored in KV doesn't match the merkle root.")
		fmt.Println("The settlement code has a bug in proof generation.")
	}
}
