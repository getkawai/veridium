package main

import (
	"encoding/json"
	"fmt"

	"github.com/kawai-network/veridium/pkg/obfuscator"
	"github.com/kawai-network/veridium/pkg/tunnelkit"
)

func main() {
	fmt.Println("=== Tunnel Token Obfuscation Demo ===\n")

	// Simulate a tunnel token (this is a fake example for demo)
	originalToken := "eyJhIjoiY2VhYjIxODc1MWQzM2NkODA0ODc4MTk2YWQ3YmVmNzQiLCJ0IjoiYWJjMTIzNDU2Nzg5IiwicyI6IlNFQ1JFVF9LRVkifQ=="

	fmt.Println("1. Original Token (simulated):")
	fmt.Printf("   %s\n\n", originalToken)

	// Create obfuscator
	obf := obfuscator.New()

	// Obfuscate the token (this is what bulk-tunnels does)
	obfuscatedToken := obf.Encode(originalToken)

	fmt.Println("2. Obfuscated Token:")
	fmt.Printf("   %s\n\n", obfuscatedToken)

	fmt.Println("3. Comparison:")
	fmt.Printf("   Original length:    %d chars\n", len(originalToken))
	fmt.Printf("   Obfuscated length:  %d chars\n", len(obfuscatedToken))
	fmt.Printf("   Overhead:           +%d chars (%.1f%%)\n\n", 
		len(obfuscatedToken)-len(originalToken),
		float64(len(obfuscatedToken)-len(originalToken))/float64(len(originalToken))*100)

	// Create a sample tunnel info with obfuscated token
	tunnel := &tunnelkit.TunnelInfo{
		TunnelID:    "abc123-def456-ghi789",
		TunnelToken: obfuscatedToken,
		Hostname:    "node-1.getkawai.com",
		PublicURL:   "https://node-1.getkawai.com",
	}

	// Simulate saving to JSON (what bulk-tunnels writes)
	jsonData, _ := json.MarshalIndent(tunnel, "", "  ")
	fmt.Println("4. JSON Output (as saved by bulk-tunnels):")
	fmt.Printf("%s\n\n", string(jsonData))

	// Simulate reading and decoding (what decode-tunnel-token does)
	var loadedTunnel tunnelkit.TunnelInfo
	json.Unmarshal(jsonData, &loadedTunnel)

	decodedToken, err := obf.Decode(loadedTunnel.TunnelToken)
	if err != nil {
		fmt.Printf("Error decoding: %v\n", err)
		return
	}

	fmt.Println("5. Decoded Token:")
	fmt.Printf("   %s\n\n", decodedToken)

	// Verify
	fmt.Println("6. Verification:")
	if decodedToken == originalToken {
		fmt.Println("   ✓ SUCCESS: Decoded token matches original!")
	} else {
		fmt.Println("   ✗ FAILED: Tokens don't match")
	}

	// Show multiple tokens
	fmt.Println("\n=== Multiple Tokens Demo ===\n")
	
	tokens := []string{
		"token-1-short",
		"token-2-medium-length-example",
		"token-3-very-long-example-with-many-characters-to-demonstrate-obfuscation",
	}

	for i, token := range tokens {
		obfuscated := obf.Encode(token)
		fmt.Printf("[%d] Original:    %s\n", i+1, token)
		fmt.Printf("    Obfuscated:  %s\n", obfuscated)
		decoded, _ := obf.Decode(obfuscated)
		fmt.Printf("    Decoded:     %s\n", decoded)
		fmt.Printf("    Match:       %v\n\n", decoded == token)
	}

	fmt.Println("=== Benefits ===")
	fmt.Println("✓ Tokens are obfuscated in JSON files")
	fmt.Println("✓ Not easily readable or copy-pasteable")
	fmt.Println("✓ Can be decoded when needed using decode-tunnel-token")
	fmt.Println("✓ No key management required")
	fmt.Println("✓ Deterministic - same token always produces same obfuscation")
	fmt.Println("\n⚠️  Important: This is obfuscation, not encryption!")
	fmt.Println("    Anyone with the code can decode the tokens.")
	fmt.Println("    Use proper encryption for production secrets.")
}

