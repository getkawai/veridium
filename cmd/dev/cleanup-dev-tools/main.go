package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	if err := cleanupDevTools(); err != nil {
		log.Fatalf("Failed to cleanup dev tools: %v", err)
	}
}

func cleanupDevTools() error {
	fmt.Println("🧹 Cleaning Up Development Tools")
	fmt.Println("═══════════════════════════════════════")
	fmt.Println()

	// List of development tools that can be cleaned up
	devTools := []string{
		"cmd/dev/debug-contract-leaf",
		"cmd/dev/debug-correct-leaf",
		"cmd/dev/debug-claim-call",
		"cmd/dev/manual-proof-verify",
		"cmd/dev/test-merkle-verification",
		"cmd/dev/verify-merkle-proof",
		"cmd/dev/check-sender-address",
		"cmd/dev/create-correct-mining-settlement",
		"cmd/dev/create-fixed-mining-settlement",
		"cmd/dev/create-treasury-mining-settlement",
		"cmd/dev/test-correct-mining-claim",
		"cmd/dev/test-fixed-mining-claim",
		"cmd/dev/test-treasury-mining-claim",
		"cmd/dev/upload-correct-mining-root",
		"cmd/dev/upload-treasury-mining-root",
		"cmd/dev/test-direct-mining-claim",
		"cmd/dev/create-test-mining-settlement",
		"cmd/dev/create-proper-test-mining-settlement",
		"cmd/dev/upload-test-mining-root",
		"cmd/dev/upload-proper-test-root",
	}

	// Keep these essential tools
	keepTools := []string{
		"cmd/dev/create-final-treasury-settlement",
		"cmd/dev/test-final-treasury-claim",
		"cmd/dev/update-period-8-root",
		"cmd/dev/check-contract-root",
		"cmd/dev/check-mining-contract-state",
		"cmd/dev/fix-mining-periods",
		"cmd/dev/upload-all-mining-roots",
	}

	fmt.Printf("📋 Development Tools Cleanup Plan:\n")
	fmt.Printf("   Tools to remove: %d\n", len(devTools))
	fmt.Printf("   Tools to keep:   %d\n", len(keepTools))
	fmt.Println()

	// Ask for confirmation
	fmt.Printf("⚠️  This will remove %d development tools that are no longer needed.\n", len(devTools))
	fmt.Printf("Essential tools will be kept for future maintenance.\n")
	fmt.Printf("Continue? (y/n): ")

	var response string
	fmt.Scanln(&response)

	if response != "y" && response != "yes" {
		fmt.Println("❌ Cleanup cancelled")
		return nil
	}

	// Remove development tools
	removedCount := 0
	for _, tool := range devTools {
		if _, err := os.Stat(tool); err == nil {
			fmt.Printf("🗑️  Removing %s\n", tool)
			if err := os.RemoveAll(tool); err != nil {
				fmt.Printf("   ❌ Failed to remove: %v\n", err)
			} else {
				removedCount++
			}
		}
	}

	fmt.Println()
	fmt.Printf("✅ Cleanup completed!\n")
	fmt.Printf("   Removed: %d tools\n", removedCount)
	fmt.Printf("   Kept:    %d essential tools\n", len(keepTools))
	fmt.Println()

	fmt.Printf("📋 Remaining Essential Tools:\n")
	for _, tool := range keepTools {
		if _, err := os.Stat(tool); err == nil {
			fmt.Printf("   ✅ %s\n", tool)
		}
	}

	return nil
}
