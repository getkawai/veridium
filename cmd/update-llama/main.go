package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kawai-network/veridium/pkg/fantasy/llamalib/download"
)

func main() {
	var (
		force       bool
		showVer     bool
		listVer     bool
		processor   string
		installPath string
		version     string
	)

	flag.BoolVar(&force, "force", false, "Force update")
	flag.BoolVar(&showVer, "version", false, "Show current version")
	flag.BoolVar(&listVer, "list", false, "List available versions")
	flag.StringVar(&processor, "processor", "cpu", "Processor type (cpu, cuda, vulkan, metal)")
	flag.StringVar(&installPath, "path", "", "Custom installation path")
	flag.Parse()

	// Default path logic matching common usage or user's environment
	if installPath == "" {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		installPath = filepath.Join(homeDir, ".llama-cpp", "bin")
	}

	if showVer {
		// logic to check installed version not easily available without running the binary or checking a manifest
		// For now just print incomplete implementation message or try to check files
		fmt.Println("Checking installed version is not fully implemented in this repair script.")
		return
	}

	targetVersion := version
	if targetVersion == "" {
		fmt.Println("Fetching latest version...")
		latest, err := download.LlamaLatestVersion()
		if err != nil {
			log.Fatalf("Failed to get latest version: %v", err)
		}
		targetVersion = latest
		fmt.Printf("Latest version: %s\n", targetVersion)
	}

	fmt.Printf("Installing version %s to %s with processor %s...\n", targetVersion, installPath, processor)

	// Determine OS and Arch
	goOS := runtime.GOOS
	goArch := runtime.GOARCH

	// Map runtime OS to download package OS
	// The download package expects "linux", "darwin", "windows" which match runtime.GOOS mostly

	err := download.Get(goArch, goOS, processor, targetVersion, installPath)
	if err != nil {
		log.Fatalf("Failed to download: %v", err)
	}

	// Also download the model if needed? No, this is just the binary.

	fmt.Println("Successfully updated llama.cpp!")
}
