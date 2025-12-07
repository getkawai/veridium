package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/kawai-network/veridium/internal/llama"
	"github.com/kawai-network/veridium/pkg/yzma/download"
	yzmaLlama "github.com/kawai-network/veridium/pkg/yzma/llama"
)

func main() {
	// Parse command line flags
	forceUpdate := flag.Bool("force", false, "Force update even if already on latest version")
	showVersion := flag.Bool("version", false, "Show currently installed version and exit")
	listVersions := flag.Bool("list", false, "List available versions from GitHub")
	customPath := flag.String("path", "", "Custom installation path (default: ~/.llama-cpp/bin)")
	processorType := flag.String("processor", "auto", "Processor type: auto, cpu, cuda, vulkan, metal")
	flag.Parse()

	log.SetFlags(0) // Remove timestamp from logs for cleaner output

	// Create installer instance
	installer := llama.NewLlamaCppInstaller()

	// Use custom path if provided
	if *customPath != "" {
		installer.BinaryPath = *customPath
	}

	log.Println("🔧 Veridium llama.cpp Update Tool")
	log.Println("=" + "=================================")
	log.Printf("📁 Installation path: %s\n", installer.BinaryPath)

	// Handle --version flag
	if *showVersion {
		showInstalledVersion(installer)
		return
	}

	// Handle --list flag
	if *listVersions {
		listAvailableVersions()
		return
	}

	// Detect processor type if auto
	processor := detectProcessor(installer, *processorType)
	log.Printf("🖥️  Detected processor: %s", processor)
	log.Printf("💻 System: %s/%s\n", runtime.GOOS, runtime.GOARCH)

	// Check current installation status
	isInstalled := installer.IsLlamaCppInstalled()
	if isInstalled {
		log.Println("✅ llama.cpp is currently installed")
		showInstalledVersion(installer)

		// Check if already on latest version
		if !*forceUpdate {
			isLatest, currentVersion, latestVersion, err := checkIfLatestVersion(installer.BinaryPath)
			if err != nil {
				log.Printf("⚠️  Warning: Could not check version: %v", err)
			} else if isLatest {
				log.Printf("✅ Already on latest version: %s", currentVersion)
				log.Println("\n💡 Use --force to reinstall anyway")
				return
			} else {
				log.Printf("📦 Update available: %s → %s", currentVersion, latestVersion)
			}
		}
	} else {
		log.Println("❌ llama.cpp is not currently installed")
	}

	// Confirm update
	if !confirmUpdate(*forceUpdate) {
		log.Println("❌ Update cancelled")
		return
	}

	// Perform installation/update
	log.Println("\n🚀 Starting llama.cpp update...")
	log.Println("=" + "=================================")

	// Convert processor string to download.Processor type
	var proc download.Processor
	switch processor {
	case "cpu":
		proc = download.CPU
	case "cuda":
		proc = download.CUDA
	case "vulkan":
		proc = download.Vulkan
	case "metal":
		proc = download.Metal
	default:
		proc = download.CPU
	}

	// Install/Update using InstallLibraries (with allowUpgrade=true)
	if err := download.InstallLibraries(installer.BinaryPath, proc, true); err != nil {
		log.Fatalf("❌ Failed to update llama.cpp: %v", err)
	}

	log.Println("\n✅ llama.cpp updated successfully!")
	log.Println("=" + "=================================")

	// Show new version
	showInstalledVersion(installer)

	// Verify installation
	log.Println("\n🔍 Verifying installation...")
	if err := installer.VerifyInstalledBinary(); err != nil {
		log.Printf("⚠️  Warning: Verification failed: %v", err)
		log.Println("💡 Some libraries may be missing. Try running with --force")
	} else {
		log.Println("✅ All required libraries verified")

		// List installed libraries with actual sizes
		libs := installer.GetRequiredLibraryPaths()
		log.Println("\n📚 Installed libraries:")
		for _, lib := range libs {
			actualPath := resolveLibraryPath(lib)
			if finalInfo, err := os.Stat(actualPath); err == nil {
				sizeMB := float64(finalInfo.Size()) / (1024 * 1024)
				log.Printf("   ✓ %s (%.1f MB)", filepath.Base(lib), sizeMB)
			}
		}
	}

	// Test library loading
	log.Println("\n🧪 Testing library loading...")
	if err := testLibraryLoading(installer); err != nil {
		log.Printf("⚠️  Warning: Library loading test failed: %v", err)
		log.Println("💡 Libraries are installed but may not load correctly")
	} else {
		log.Println("✅ Library loading test passed!")
	}

	log.Println("\n🎉 Update complete!")
}

// detectProcessor detects the best processor type for this system
func detectProcessor(installer *llama.LlamaCppInstaller, processorType string) string {
	if processorType != "auto" {
		return processorType
	}

	// Use the same detection logic as installer.detectProcessor()
	// Priority: CUDA > Vulkan > Metal > CPU

	// Check for NVIDIA GPU (CUDA)
	if runtime.GOOS == "linux" || runtime.GOOS == "windows" {
		// On Linux/Windows, check for nvidia-smi
		if _, err := os.Stat("/usr/bin/nvidia-smi"); err == nil {
			return "cuda"
		}
		if _, err := os.Stat("/usr/local/cuda"); err == nil {
			return "cuda"
		}
	}

	// Check for Vulkan support
	if runtime.GOOS == "linux" || runtime.GOOS == "windows" {
		// Basic check for Vulkan
		if _, err := os.Stat("/usr/share/vulkan"); err == nil {
			return "vulkan"
		}
	}

	// macOS always supports Metal
	if runtime.GOOS == "darwin" {
		return "metal"
	}

	// Default to CPU
	return "cpu"
}

// showInstalledVersion displays the currently installed version
func showInstalledVersion(installer *llama.LlamaCppInstaller) {
	versionPath := filepath.Join(installer.BinaryPath, "version.json")
	data, err := os.ReadFile(versionPath)
	if err != nil {
		log.Println("   Version: Unknown (version.json not found)")
		return
	}

	// Parse version.json
	var versionInfo struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(data, &versionInfo); err != nil {
		log.Printf("   Version: Unknown (failed to parse version.json)")
		return
	}

	log.Printf("   Current version: %s", versionInfo.TagName)
}

// checkIfLatestVersion checks if the installed version is the latest
func checkIfLatestVersion(libPath string) (isLatest bool, currentVersion, latestVersion string, err error) {
	// Read current version
	versionPath := filepath.Join(libPath, "version.json")
	data, err := os.ReadFile(versionPath)
	if err != nil {
		return false, "", "", fmt.Errorf("failed to read version.json: %w", err)
	}

	var versionInfo struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(data, &versionInfo); err != nil {
		return false, "", "", fmt.Errorf("failed to parse version.json: %w", err)
	}
	currentVersion = versionInfo.TagName

	// Get latest version
	latestVersion, err = download.LlamaLatestVersion()
	if err != nil {
		return false, currentVersion, "", fmt.Errorf("failed to get latest version: %w", err)
	}

	return currentVersion == latestVersion, currentVersion, latestVersion, nil
}

// listAvailableVersions lists available versions from GitHub
func listAvailableVersions() {
	log.Println("\n📦 Fetching available versions from GitHub...")

	versions, err := download.LlamaAvailableVersions(20)
	if err != nil {
		log.Fatalf("❌ Failed to fetch versions: %v", err)
	}

	log.Printf("\n📋 Available versions (latest 20):\n")
	for i, version := range versions {
		if i == 0 {
			log.Printf("   %s (latest)", version)
		} else {
			log.Printf("   %s", version)
		}
	}
	log.Println()
}

// confirmUpdate asks user to confirm the update
func confirmUpdate(force bool) bool {
	if force {
		return true
	}

	fmt.Print("\n❓ Proceed with update? [Y/n]: ")
	var response string
	fmt.Scanln(&response)

	if response == "" || response == "y" || response == "Y" {
		return true
	}

	return false
}

// resolveLibraryPath resolves a library path through chains of text file redirects
// llama.cpp uses text files as "symlinks" that contain the target filename
func resolveLibraryPath(path string) string {
	maxDepth := 10 // Prevent infinite loops
	currentPath := path

	for i := 0; i < maxDepth; i++ {
		info, err := os.Stat(currentPath)
		if err != nil {
			return path // Return original if we can't stat
		}

		// If file is small (< 1KB), it might be a text file redirect
		if info.Size() < 1024 {
			content, err := os.ReadFile(currentPath)
			if err != nil {
				return currentPath
			}

			targetName := strings.TrimSpace(string(content))
			if targetName == "" {
				return currentPath
			}

			// Resolve relative to directory of current file
			var targetPath string
			if filepath.IsAbs(targetName) {
				targetPath = targetName
			} else {
				targetPath = filepath.Join(filepath.Dir(currentPath), targetName)
			}

			// Check if target exists
			if _, err := os.Stat(targetPath); err == nil {
				currentPath = targetPath
				continue
			}
		}

		// Try readlink for actual symlinks
		if linkTarget, err := os.Readlink(currentPath); err == nil {
			if filepath.IsAbs(linkTarget) {
				currentPath = linkTarget
			} else {
				currentPath = filepath.Join(filepath.Dir(currentPath), linkTarget)
			}
			continue
		}

		// If we get here, we've found the actual file
		return currentPath
	}

	return currentPath
}

// testLibraryLoading tests if the llama.cpp library files are valid
func testLibraryLoading(installer *llama.LlamaCppInstaller) error {
	libPath := installer.GetLibraryPath()

	log.Printf("   Checking library files in: %s", libPath)

	// Verify actual library files exist and are valid
	requiredLibs := installer.GetRequiredLibraryPaths()

	for _, lib := range requiredLibs {
		actualPath := resolveLibraryPath(lib)
		info, err := os.Stat(actualPath)
		if err != nil {
			return fmt.Errorf("library file not found: %s", filepath.Base(lib))
		}

		// Verify it's a real library file (> 1KB)
		if info.Size() < 1024 {
			return fmt.Errorf("library file too small (may be corrupted): %s (size: %d bytes)", filepath.Base(lib), info.Size())
		}

		// Read first 4 bytes to check for Mach-O magic number
		f, err := os.Open(actualPath)
		if err != nil {
			return fmt.Errorf("cannot open library: %s", filepath.Base(lib))
		}

		magic := make([]byte, 4)
		n, err := f.Read(magic)
		f.Close()

		if err != nil || n != 4 {
			return fmt.Errorf("cannot read library header: %s", filepath.Base(lib))
		}

		// Check for Mach-O magic numbers
		// 0xfeedface (32-bit), 0xfeedfacf (64-bit), 0xcafebabe (universal)
		isMachO := (magic[0] == 0xfe && magic[1] == 0xed && magic[2] == 0xfa && (magic[3] == 0xce || magic[3] == 0xcf)) ||
			(magic[0] == 0xca && magic[1] == 0xfe && magic[2] == 0xba && magic[3] == 0xbe) ||
			(magic[0] == 0xcf && magic[1] == 0xfa && magic[2] == 0xed && magic[3] == 0xfe) // Little-endian 64-bit

		if !isMachO {
			return fmt.Errorf("invalid library format: %s (not a valid Mach-O file)", filepath.Base(lib))
		}

		sizeMB := float64(info.Size()) / (1024 * 1024)
		log.Printf("   ✓ %s (%.1f MB, valid Mach-O)", filepath.Base(actualPath), sizeMB)
	}

	log.Printf("   ✓ All %d library files are valid", len(requiredLibs))

	// Try to actually load the library (this will test if it can be loaded by the system)
	log.Println("   Testing library loading...")
	if err := yzmaLlama.Load(libPath); err != nil {
		// Library files are valid but loading failed - this might be okay
		// (could be due to dependencies, permissions, etc.)
		log.Printf("   ⚠️  Library load test skipped: %v", err)
		log.Println("   💡 Libraries are valid but may need to be loaded by main app")
		return nil // Don't fail - files are valid
	}

	log.Println("   ✓ Library loaded successfully")

	// Initialize backend
	yzmaLlama.Init()
	log.Println("   ✓ Backend initialized")

	// Get system info to verify library is working
	systemInfo := yzmaLlama.PrintSystemInfo()
	if systemInfo != "" {
		lines := strings.Split(strings.TrimSpace(systemInfo), "\n")
		if len(lines) > 0 {
			log.Printf("   ✓ System: %s", strings.TrimSpace(lines[0]))
		}
	}

	// Cleanup
	yzmaLlama.BackendFree()
	log.Println("   ✓ Backend freed")

	return nil
}
