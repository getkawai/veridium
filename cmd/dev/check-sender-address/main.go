package main

import (
	"crypto/ecdsa"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/kawai-network/x/constant"
)

func main() {
	if err := checkSenderAddress(); err != nil {
		log.Fatalf("Failed to check sender address: %v", err)
	}
}

func checkSenderAddress() error {
	fmt.Println("🔍 Sender Address Check")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// Get private key
	privateKeyHex := constant.GetAdminPrivateKey()
	if strings.HasPrefix(privateKeyHex, "0x") {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Get public key and address
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		return fmt.Errorf("failed to cast public key")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	fmt.Printf("📋 Transaction Sender:\n")
	fmt.Printf("   Private Key: %s...\n", privateKeyHex[:10])
	fmt.Printf("   Address:     %s\n", address.Hex())
	fmt.Println()

	// Expected test address
	testAddress := "0x0f3e75B9Bb3efcD87B1Ed15a30C8a7FBaABD204E"

	fmt.Printf("📋 Expected Test Address:\n")
	fmt.Printf("   Address:     %s\n", testAddress)
	fmt.Println()

	fmt.Printf("🎯 Address Comparison:\n")
	if strings.EqualFold(address.Hex(), testAddress) {
		fmt.Printf("   Match:       ✅ YES\n")
		fmt.Printf("   Status:      Transaction sender matches test address\n")
	} else {
		fmt.Printf("   Match:       ❌ NO\n")
		fmt.Printf("   Status:      Transaction sender does NOT match test address\n")
		fmt.Printf("   Issue:       msg.sender in contract will be %s, not %s\n", address.Hex(), testAddress)
		fmt.Printf("   Solution:    Need to use the correct private key or generate leaf for actual sender\n")
	}

	return nil
}
