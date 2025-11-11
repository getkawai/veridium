package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/pkg/yzma/download"
	yzma "github.com/kawai-network/veridium/pkg/yzma/llama"
)

func main() {
	fmt.Println("🧪 Testing Library Path Management")
	fmt.Println("===================================\n")

	// Step 1: Create installer
	fmt.Println("📦 Step 1: Creating installer...")
	installer := llama.NewLlamaCppInstaller()
	fmt.Printf("   Binary Path: %s\n", installer.BinaryPath)
	fmt.Printf("   Metadata Path: %s\n", installer.MetadataPath)
	fmt.Printf("   Models Dir: %s\n\n", installer.ModelsDir)

	// Step 2: Check if llama.cpp is installed
	fmt.Println("🔍 Step 2: Checking installation status...")
	isInstalled := installer.IsLlamaCppInstalled()
	fmt.Printf("   Installed: %v\n", isInstalled)

	if !isInstalled {
		fmt.Println("\n📥 llama.cpp not installed, downloading...")

		// Get latest release
		release, err := installer.GetLatestRelease()
		if err != nil {
			log.Fatalf("❌ Failed to get latest release: %v", err)
		}
		fmt.Printf("   Latest version: %s\n", release.Version)

		// Download
		if err := installer.DownloadRelease(release.Version, nil); err != nil {
			log.Fatalf("❌ Failed to download: %v", err)
		}
		fmt.Println("   ✅ Download completed!")
	} else {
		version := installer.GetInstalledVersion()
		fmt.Printf("   Installed version: %s\n", version)
	}

	// Step 3: Verify all libraries exist
	fmt.Println("\n🔍 Step 3: Verifying all required libraries...")
	allExist := installer.VerifyAllLibrariesExist()
	fmt.Printf("   All libraries present: %v\n", allExist)

	if !allExist {
		fmt.Println("\n❌ Not all libraries found!")
		if err := installer.VerifyInstalledBinary(); err != nil {
			fmt.Printf("   Error: %v\n", err)
		}
		os.Exit(1)
	}

	// Step 4: Get library paths
	fmt.Println("\n📂 Step 4: Getting library paths...")
	libPath := installer.GetLibraryPath()
	fmt.Printf("   Library Directory: %s\n", libPath)

	libFilePath := installer.GetLibraryFilePath()
	fmt.Printf("   Main Library File: %s\n", libFilePath)

	requiredPaths := installer.GetRequiredLibraryPaths()
	fmt.Printf("   Required Libraries (%d):\n", len(requiredPaths))
	for i, path := range requiredPaths {
		exists := "✅"
		if _, err := os.Stat(path); err != nil {
			exists = "❌"
		}
		fmt.Printf("     %d. %s %s\n", i+1, exists, path)
	}

	// Step 5: Verify paths match download.RequiredLibraries()
	fmt.Println("\n🔍 Step 5: Verifying paths match download package...")
	expectedLibs := download.RequiredLibraries(runtime.GOOS)
	fmt.Printf("   Expected libraries from download package: %v\n", expectedLibs)

	if len(requiredPaths) != len(expectedLibs) {
		fmt.Printf("   ❌ Mismatch! Expected %d, got %d\n", len(expectedLibs), len(requiredPaths))
		os.Exit(1)
	}

	for i, expectedLib := range expectedLibs {
		expectedPath := filepath.Join(libPath, expectedLib)
		actualPath := requiredPaths[i]

		if expectedPath != actualPath {
			fmt.Printf("   ❌ Path mismatch at index %d:\n", i)
			fmt.Printf("      Expected: %s\n", expectedPath)
			fmt.Printf("      Actual:   %s\n", actualPath)
			os.Exit(1)
		}
		fmt.Printf("   ✅ Match: %s\n", filepath.Base(actualPath))
	}

	// Step 6: Test library loading with yzma
	fmt.Println("\n🔧 Step 6: Testing library loading with yzma...")
	fmt.Printf("   Loading from: %s\n", libPath)

	if err := yzma.Load(libPath); err != nil {
		log.Fatalf("❌ Failed to load library: %v", err)
	}
	fmt.Println("   ✅ Libraries loaded successfully!")

	// Step 7: Initialize llama.cpp
	fmt.Println("\n🚀 Step 7: Initializing llama.cpp backend...")
	yzma.Init()
	fmt.Println("   ✅ Backend initialized successfully!")

	// Step 8: Test with LibraryService
	fmt.Println("\n🧪 Step 8: Testing with LibraryService...")
	service, err := llama.NewLibraryService()
	if err != nil {
		log.Fatalf("❌ Failed to create service: %v", err)
	}
	defer service.Cleanup()

	fmt.Println("   ✅ LibraryService created successfully!")

	// Step 9: Verify service paths match installer
	fmt.Println("\n🔍 Step 9: Verifying service uses correct paths...")
	serviceModelsDir := service.GetModelsDirectory()
	installerModelsDir := installer.GetModelsDirectory()

	if serviceModelsDir != installerModelsDir {
		fmt.Printf("   ❌ Models directory mismatch:\n")
		fmt.Printf("      Service:   %s\n", serviceModelsDir)
		fmt.Printf("      Installer: %s\n", installerModelsDir)
		os.Exit(1)
	}
	fmt.Printf("   ✅ Models directory matches: %s\n", serviceModelsDir)

	// Step 10: Summary
	fmt.Println("\n" + strings.Repeat("=", 50))
	fmt.Println("✅ ALL TESTS PASSED!")
	fmt.Println(strings.Repeat("=", 50))
	fmt.Println("\n📊 Summary:")
	fmt.Printf("   Platform: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("   Library Path: %s\n", libPath)
	fmt.Printf("   Required Libraries: %d\n", len(requiredPaths))
	fmt.Printf("   All Libraries Present: ✅\n")
	fmt.Printf("   Library Loading: ✅\n")
	fmt.Printf("   Backend Initialization: ✅\n")
	fmt.Printf("   LibraryService Integration: ✅\n")
	fmt.Println("\n🎉 Library path management is working correctly!")
}
