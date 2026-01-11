package main

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kawai-network/veridium/pkg/merkle"
)

func main() {
	if err := testSimpleMerkle(); err != nil {
		log.Fatalf("Failed to test simple Merkle: %v", err)
	}
}

func testSimpleMerkle() error {
	fmt.Println("🧪 Simple Merkle Tree Test")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// Create simple test leaves
	leaf1 := crypto.Keccak256([]byte("leaf1"))
	leaf2 := crypto.Keccak256([]byte("leaf2"))
	leaf3 := crypto.Keccak256([]byte("leaf3"))

	leaves := [][]byte{leaf1, leaf2, leaf3}

	fmt.Printf("📋 Test Leaves:\n")
	for i, leaf := range leaves {
		fmt.Printf("   [%d]: 0x%x\n", i, leaf)
	}
	fmt.Println()

	// Build Merkle tree
	tree := merkle.NewMerkleTree(leaves)

	fmt.Printf("🌳 Tree Root: 0x%x\n", tree.Root)
	fmt.Println()

	// Generate proof for leaf1
	proof, ok := tree.GetProof(leaf1)
	if !ok {
		return fmt.Errorf("failed to generate proof for leaf1")
	}

	fmt.Printf("🔍 Proof for leaf1:\n")
	for i, p := range proof {
		fmt.Printf("   [%d]: 0x%x\n", i, p)
	}
	fmt.Println()

	// Check if proof contains the leaf itself
	for i, p := range proof {
		if fmt.Sprintf("%x", p) == fmt.Sprintf("%x", leaf1) {
			fmt.Printf("❌ ERROR: Proof element [%d] is the leaf itself!\n", i)
			return fmt.Errorf("proof contains leaf itself")
		}
	}

	fmt.Printf("✅ Proof does not contain the leaf itself\n")
	fmt.Println()

	// Manual verification (OpenZeppelin style)
	fmt.Printf("🧪 Manual Verification:\n")
	computedHash := leaf1
	fmt.Printf("   Start:  0x%x (leaf1)\n", computedHash)

	for i, proofElement := range proof {
		if string(computedHash) <= string(proofElement) {
			computedHash = crypto.Keccak256(computedHash, proofElement)
		} else {
			computedHash = crypto.Keccak256(proofElement, computedHash)
		}
		fmt.Printf("   Step %d: 0x%x\n", i+1, computedHash)
	}

	fmt.Printf("   Final:  0x%x\n", computedHash)
	fmt.Printf("   Root:   0x%x\n", tree.Root)

	isValid := fmt.Sprintf("%x", computedHash) == fmt.Sprintf("%x", tree.Root)
	fmt.Printf("   Valid:  %t\n", isValid)

	if !isValid {
		return fmt.Errorf("verification failed")
	}

	fmt.Printf("✅ Simple Merkle tree test passed!\n")
	return nil
}
