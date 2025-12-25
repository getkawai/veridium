package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/pkg/obfuscator"
	"github.com/kawai-network/veridium/pkg/tunnelkit"
)

const (
	// Hardcoded credentials as requested
	defaultAccountID = "ceab218751d33cd804878196ad7bef74"
	defaultAPIToken  = "OP8BZQhyeJxrovCPKt15eUOSC6i5LXTVECGRSMc1"

	defaultDomain      = "getkawai.com" // Default domain, can be changed
	defaultTunnelCount = 3
	defaultOutputFile  = "tunnels.json"
)

func main() {
	// Although user said "other flags not needed", keeping them as optional overrides
	// is good practice, but defaulting to the requested constants.
	accountID := flag.String("acc-id", defaultAccountID, "Cloudflare Account ID")
	apiToken := flag.String("token", defaultAPIToken, "Cloudflare API Token")
	domain := flag.String("domain", defaultDomain, "Base domain for hostnames")
	count := flag.Int("count", defaultTunnelCount, "Number of tunnels to create")
	output := flag.String("out", defaultOutputFile, "Output JSON file path")

	flag.Parse()

	if *accountID == "" || *apiToken == "" {
		log.Fatal("Account ID and API Token are required")
	}

	fmt.Printf("Starting creation of %d tunnels for domain %s...\n", *count, *domain)

	var tunnels []*tunnelkit.TunnelInfo

	// Initialize obfuscator for TunnelToken
	obf := obfuscator.New()

	// Create output directory if it doesn't exist
	outDir := filepath.Dir(*output)
	if outDir != "." {
		if err := os.MkdirAll(outDir, 0755); err != nil {
			log.Fatalf("Failed to create output directory: %v", err)
		}
	}

	for i := 0; i < *count; i++ {
		tunnelName := fmt.Sprintf("node-%d", i+1)
		hostname := fmt.Sprintf("%s.%s", tunnelName, *domain)

		fmt.Printf("[%d/%d] Creating tunnel '%s' (Host: %s)... ", i+1, *count, tunnelName, hostname)

		cfg := tunnelkit.Config{
			AccountID:  *accountID,
			APIToken:   *apiToken,
			TunnelName: tunnelName,
			Hostname:   hostname,
			LocalURL:   constant.LocalWorkerURL,
		}

		// Use the refactored function
		info, err := tunnelkit.GetOrCreateTunnelWithDNS(cfg)
		if err != nil {
			fmt.Printf("FAILED: %v\n", err)
			continue
		}

		// Obfuscate the TunnelToken before storing
		if info.TunnelToken != "" {
			info.TunnelToken = obf.Encode(info.TunnelToken)
		}

		tunnels = append(tunnels, info)
		fmt.Printf("SUCCESS (ID: %s)\n", info.TunnelID)

		// Small delay to avoid aggressive rate limiting if necessary
		// Cloudflare API rate limits are usually high, but let's be safe.
		// time.Sleep(100 * time.Millisecond)
	}

	// Write results to JSON
	file, err := os.Create(*output)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(tunnels); err != nil {
		log.Fatalf("Failed to write JSON: %v", err)
	}

	fmt.Printf("\nDone! Successfully created %d/%d tunnels.\n", len(tunnels), *count)
	fmt.Printf("Results saved to %s\n", *output)
}
