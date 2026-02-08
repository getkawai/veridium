package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kawai-network/veridium/pkg/tunnelkit"
	"github.com/kawai-network/x/env"
)

func main() {
	inputFile := flag.String("file", "tunnels.json", "Path to tunnels JSON file")
	tunnelID := flag.String("id", "", "Specific tunnel ID or hostname to decode (optional)")
	showAll := flag.Bool("all", false, "Show all decoded tokens")

	flag.Parse()

	if *tunnelID != "" {
		// Decode specific tunnel
		if err := decodeSpecificTunnel(*inputFile, *tunnelID); err != nil {
			log.Fatalf("Error: %v", err)
		}
	} else if *showAll {
		// Decode all tunnels
		if err := decodeAllTunnels(*inputFile); err != nil {
			log.Fatalf("Error: %v", err)
		}
	} else {
		// Default: show list without decoding
		if err := listTunnels(*inputFile); err != nil {
			log.Fatalf("Error: %v", err)
		}
	}
}

func listTunnels(inputFile string) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var tunnels []*tunnelkit.TunnelInfo
	if err := json.Unmarshal(data, &tunnels); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	fmt.Printf("Found %d tunnel(s) in %s:\n\n", len(tunnels), inputFile)

	for i, tunnel := range tunnels {
		fmt.Printf("[%d] Tunnel ID: %s\n", i+1, tunnel.TunnelID)
		fmt.Printf("    Hostname: %s\n", tunnel.Hostname)
		fmt.Printf("    Public URL: %s\n", tunnel.PublicURL)
		fmt.Printf("    Token: %s (obfuscated)\n\n", tunnel.TunnelToken[:min(40, len(tunnel.TunnelToken))]+"...")
	}

	fmt.Println("Use --all to decode all tokens, or --id <tunnel-id> to decode a specific token")
	return nil
}

func decodeAllTunnels(inputFile string) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var tunnels []*tunnelkit.TunnelInfo
	if err := json.Unmarshal(data, &tunnels); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	obf := env.New()

	fmt.Printf("Decoding %d tunnel token(s)...\n\n", len(tunnels))

	for i, tunnel := range tunnels {
		fmt.Printf("[%d] Tunnel ID: %s\n", i+1, tunnel.TunnelID)
		fmt.Printf("    Hostname: %s\n", tunnel.Hostname)
		fmt.Printf("    Public URL: %s\n", tunnel.PublicURL)

		if tunnel.TunnelToken != "" {
			decoded, err := obf.Decode(tunnel.TunnelToken)
			if err != nil {
				fmt.Printf("    Token: ERROR - %v\n\n", err)
				continue
			}
			fmt.Printf("    Token (Decoded): %s\n\n", decoded)
		} else {
			fmt.Printf("    Token: (empty)\n\n")
		}
	}

	return nil
}

func decodeSpecificTunnel(inputFile, identifier string) error {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	var tunnels []*tunnelkit.TunnelInfo
	if err := json.Unmarshal(data, &tunnels); err != nil {
		return fmt.Errorf("failed to parse JSON: %w", err)
	}

	obf := env.New()

	// Find tunnel by ID or hostname
	for _, tunnel := range tunnels {
		if tunnel.TunnelID == identifier || tunnel.Hostname == identifier {
			fmt.Printf("Tunnel ID: %s\n", tunnel.TunnelID)
			fmt.Printf("Hostname: %s\n", tunnel.Hostname)
			fmt.Printf("Public URL: %s\n", tunnel.PublicURL)

			if tunnel.TunnelToken == "" {
				return fmt.Errorf("tunnel token is empty")
			}

			decoded, err := obf.Decode(tunnel.TunnelToken)
			if err != nil {
				return fmt.Errorf("failed to decode token: %w", err)
			}

			fmt.Printf("\nObfuscated Token:\n%s\n", tunnel.TunnelToken)
			fmt.Printf("\nDecoded Token:\n%s\n", decoded)
			return nil
		}
	}

	return fmt.Errorf("tunnel not found: %s", identifier)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
