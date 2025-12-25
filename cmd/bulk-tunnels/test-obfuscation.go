package main

import (
	"encoding/json"
	"fmt"

	"github.com/kawai-network/veridium/pkg/obfuscator"
	"github.com/kawai-network/veridium/pkg/tunnelkit"
)

// TestObfuscation demonstrates how tunnel tokens are obfuscated and decoded
func TestObfuscation() {
	fmt.Println("=== Tunnel Token Obfuscation Test ===\n")

	// Simulate a tunnel token (this is a fake example)
	originalToken := "eyJhIjoiY2VhYjIxODc1MWQzM2NkODA0ODc4MTk2YWQ3YmVmNzQiLCJ0IjoiYWJjMTIzNDU2Nzg5IiwicyI6IlNFQ1JFVF9LRVkifQ=="

	fmt.Println("1. Original Token:")
	fmt.Printf("   %s\n\n", originalToken)

	// Create obfuscator
	obf := obfuscator.New()

	// Obfuscate the token (this is what bulk-tunnels does)
	obfuscatedToken := obf.Encode(originalToken)

	fmt.Println("2. Obfuscated Token:")
	fmt.Printf("   %s\n\n", obfuscatedToken)

	// Create a sample tunnel info with obfuscated token
	tunnel := &tunnelkit.TunnelInfo{
		TunnelID:    "abc123-def456-ghi789",
		TunnelToken: obfuscatedToken,
		Hostname:    "node-1.getkawai.com",
		PublicURL:   "https://node-1.getkawai.com",
	}

	// Simulate saving to JSON (what bulk-tunnels writes)
	jsonData, _ := json.MarshalIndent(tunnel, "", "  ")
	fmt.Println("3. JSON Output (as saved by bulk-tunnels):")
	fmt.Printf("%s\n\n", string(jsonData))

	// Simulate reading and decoding (what decode-tunnel-token does)
	var loadedTunnel tunnelkit.TunnelInfo
	json.Unmarshal(jsonData, &loadedTunnel)

	decodedToken, err := obf.Decode(loadedTunnel.TunnelToken)
	if err != nil {
		fmt.Printf("Error decoding: %v\n", err)
		return
	}

	fmt.Println("4. Decoded Token:")
	fmt.Printf("   %s\n\n", decodedToken)

	// Verify
	fmt.Println("5. Verification:")
	if decodedToken == originalToken {
		fmt.Println("   ✓ SUCCESS: Decoded token matches original!")
	} else {
		fmt.Println("   ✗ FAILED: Tokens don't match")
	}

	fmt.Println("\n=== Benefits ===")
	fmt.Println("✓ Token is obfuscated in JSON file")
	fmt.Println("✓ Not easily readable or copy-pasteable")
	fmt.Println("✓ Can be decoded when needed")
	fmt.Println("✓ No key management required")
	fmt.Println("\n⚠️  Remember: This is obfuscation, not encryption!")
}

