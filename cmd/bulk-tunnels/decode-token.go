package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kawai-network/veridium/pkg/obfuscator"
	"github.com/kawai-network/veridium/pkg/tunnelkit"
)

// DecodeTunnelTokens reads a tunnels.json file and decodes all obfuscated TunnelTokens
func DecodeTunnelTokens(inputFile string) error {
	// Read the JSON file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON
	var tunnels []*tunnelkit.TunnelInfo
	if err := json.Unmarshal(data, &tunnels); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Initialize obfuscator
	obf := obfuscator.New()

	// Decode all tokens
	fmt.Println("Decoding tunnel tokens...")
	fmt.Println("=" + string(make([]byte, 80)) + "=")

	for i, tunnel := range tunnels {
		fmt.Printf("\n[%d] Tunnel ID: %s\n", i+1, tunnel.TunnelID)
		fmt.Printf("    Hostname: %s\n", tunnel.Hostname)
		fmt.Printf("    Public URL: %s\n", tunnel.PublicURL)

		if tunnel.TunnelToken != "" {
			decoded, err := obf.Decode(tunnel.TunnelToken)
			if err != nil {
				fmt.Printf("    Token (ERROR): Failed to decode - %v\n", err)
				continue
			}
			fmt.Printf("    Token (Obfuscated): %s...\n", tunnel.TunnelToken[:min(50, len(tunnel.TunnelToken))])
			fmt.Printf("    Token (Decoded): %s\n", decoded)
		} else {
			fmt.Printf("    Token: (empty)\n")
		}
	}

	return nil
}

// GetDecodedToken returns the decoded token for a specific tunnel by ID or hostname
func GetDecodedToken(inputFile, identifier string) (string, error) {
	// Read the JSON file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON
	var tunnels []*tunnelkit.TunnelInfo
	if err := json.Unmarshal(data, &tunnels); err != nil {
		return "", fmt.Errorf("failed to parse JSON: %w", err)
	}

	// Initialize obfuscator
	obf := obfuscator.New()

	// Find tunnel by ID or hostname
	for _, tunnel := range tunnels {
		if tunnel.TunnelID == identifier || tunnel.Hostname == identifier {
			if tunnel.TunnelToken == "" {
				return "", fmt.Errorf("tunnel token is empty")
			}

			decoded, err := obf.Decode(tunnel.TunnelToken)
			if err != nil {
				return "", fmt.Errorf("failed to decode token: %w", err)
			}

			return decoded, nil
		}
	}

	return "", fmt.Errorf("tunnel not found: %s", identifier)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
