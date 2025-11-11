package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/pkg/yzma/download"
)

// Quick check to verify library paths without downloading
func main() {
	fmt.Println("🔍 Quick Library Path Check")
	fmt.Println("============================\n")

	// Create installer
	installer := llama.NewLlamaCppInstaller()
	
	// Get paths
	libPath := installer.GetLibraryPath()
	requiredPaths := installer.GetRequiredLibraryPaths()
	expectedLibs := download.RequiredLibraries(runtime.GOOS)
	
	// Display info
	fmt.Printf("Platform: %s/%s\n\n", runtime.GOOS, runtime.GOARCH)
	
	fmt.Printf("Library Directory:\n  %s\n\n", libPath)
	
	fmt.Printf("Expected Libraries from download.RequiredLibraries():\n")
	for i, lib := range expectedLibs {
		fmt.Printf("  %d. %s\n", i+1, lib)
	}
	
	fmt.Printf("\nActual Paths from installer.GetRequiredLibraryPaths():\n")
	for i, path := range requiredPaths {
		exists := "❌"
		if _, err := os.Stat(path); err == nil {
			exists = "✅"
		}
		fmt.Printf("  %d. %s %s\n", i+1, exists, path)
	}
	
	// Verify consistency
	fmt.Println("\nConsistency Check:")
	if len(requiredPaths) != len(expectedLibs) {
		fmt.Printf("  ❌ FAIL: Count mismatch (expected %d, got %d)\n", len(expectedLibs), len(requiredPaths))
		os.Exit(1)
	}
	fmt.Printf("  ✅ Count matches: %d libraries\n", len(requiredPaths))
	
	// Check if installed
	isInstalled := installer.IsLlamaCppInstalled()
	fmt.Printf("\nInstallation Status:\n")
	if isInstalled {
		version := installer.GetInstalledVersion()
		fmt.Printf("  ✅ Installed (version: %s)\n", version)
	} else {
		fmt.Printf("  ❌ Not installed\n")
		fmt.Printf("\nTo install, run:\n")
		fmt.Printf("  go run ./examples/installer/main.go\n")
		fmt.Printf("  OR\n")
		fmt.Printf("  go run ./cmd/test-library-paths/main.go\n")
	}
	
	fmt.Println("\n✅ Path configuration is correct!")
}

