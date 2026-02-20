package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func main() {
	// Test data
	period := uint64(6)
	userAddress := common.HexToAddress("0xEbc28d5D68ab501cDC34b74fe408c879b0c34126")
	amount := new(big.Int)
	amount.SetString("27500000000000000000", 10)

	fmt.Println("🧪 Testing Cashback Leaf Hash")
	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Printf("Period: %d\n", period)
	fmt.Printf("User: %s\n", userAddress.Hex())
	fmt.Printf("Amount: %s\n", amount.String())
	fmt.Println()

	// Method 1: abi.encodePacked (no padding)
	periodBytes := new(big.Int).SetUint64(period).Bytes()
	addressBytes := userAddress.Bytes()
	amountBytes := amount.Bytes()

	fmt.Println("📊 Method 1: abi.encodePacked (no padding)")
	fmt.Printf("   Period bytes (%d): %x\n", len(periodBytes), periodBytes)
	fmt.Printf("   Address bytes (%d): %x\n", len(addressBytes), addressBytes)
	fmt.Printf("   Amount bytes (%d): %x\n", len(amountBytes), amountBytes)

	hash1 := crypto.Keccak256(periodBytes, addressBytes, amountBytes)
	fmt.Printf("   Leaf hash: 0x%x\n", hash1)
	fmt.Println()

	// Method 2: With 32-byte padding
	periodPadded := common.LeftPadBytes(new(big.Int).SetUint64(period).Bytes(), 32)
	amountPadded := common.LeftPadBytes(amount.Bytes(), 32)

	fmt.Println("📊 Method 2: With 32-byte padding")
	fmt.Printf("   Period bytes (%d): %x\n", len(periodPadded), periodPadded)
	fmt.Printf("   Address bytes (%d): %x\n", len(addressBytes), addressBytes)
	fmt.Printf("   Amount bytes (%d): %x\n", len(amountPadded), amountPadded)

	hash2 := crypto.Keccak256(periodPadded, addressBytes, amountPadded)
	fmt.Printf("   Leaf hash: 0x%x\n", hash2)
	fmt.Println()

	fmt.Println("═══════════════════════════════════════════════════════════")
	fmt.Println("✅ Use Method 1 hash to verify on-chain")
	fmt.Println()
	fmt.Println("Verify with cast:")
	fmt.Printf("cast keccak $(cast abi-encode 'f(uint256,address,uint256)' %d %s %s)\n",
		period, userAddress.Hex(), amount.String())
}
