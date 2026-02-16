// Package whisperapp provides CLI helpers for setting up whisper with the new package.
package whisperapp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/kawai-network/veridium/internal/paths"
	"github.com/kawai-network/veridium/pkg/tools/downloader"
	whisper "github.com/kawai-network/whisper"
)

const whisperLibVersion = "v0.3.2"
const whisperLibReleaseURL = "https://api.github.com/repos/kawai-network/whisper/releases/latest"

type whisperRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
		Size               int64  `json:"size"`
	} `json:"assets"`
}

// SetupOptions contains options for whisper setup
type SetupOptions struct {
	ModelsDir      string
	LibDir         string
	DownloadModels []string
	AutoSelect     bool
	Interactive    bool
}

// DefaultSetupOptions returns default setup options
func DefaultSetupOptions() *SetupOptions {
	return &SetupOptions{
		ModelsDir:      paths.WhisperModels(),
		LibDir:         paths.WhisperLib(),
		DownloadModels: []string{},
		AutoSelect:     false,
		Interactive:    true,
	}
}

// SetupWhisper performs the complete setup process for whisper
func SetupWhisper(ctx context.Context, opts *SetupOptions) error {
	fmt.Println("=== Whisper Setup (kawai-network/whisper) ===")
	fmt.Println()

	// Step 1: Create directories
	if err := setupDirectories(opts); err != nil {
		return fmt.Errorf("failed to setup directories: %w", err)
	}

	// Step 2: Download whisper library
	if err := downloadWhisperLibrary(ctx, opts.LibDir); err != nil {
		return fmt.Errorf("failed to download whisper library: %w", err)
	}

	// Step 3: Download required models
	if err := downloadRequiredModels(ctx, opts); err != nil {
		return fmt.Errorf("failed to download models: %w", err)
	}

	// Step 4: Verify setup
	if err := verifySetup(opts); err != nil {
		return fmt.Errorf("setup verification failed: %w", err)
	}

	fmt.Println("\n=== Setup Complete! ===")
	fmt.Printf("Models directory: %s\n", opts.ModelsDir)
	fmt.Printf("Library directory: %s\n", opts.LibDir)

	return nil
}

// setupDirectories creates necessary directories
func setupDirectories(opts *SetupOptions) error {
	fmt.Println("Step 1: Creating directories...")

	dirs := []string{
		opts.ModelsDir,
		opts.LibDir,
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
		fmt.Printf("  ✓ Created: %s\n", dir)
	}

	return nil
}

// downloadWhisperLibrary downloads the gowhisper native library
func downloadWhisperLibrary(ctx context.Context, libDir string) error {
	fmt.Println("\nStep 2: Downloading whisper library...")

	libName := whisper.LibraryName(runtime.GOOS)
	libPath := fmt.Sprintf("%s/%s", libDir, libName)

	if _, err := os.Stat(libPath); err == nil {
		fmt.Printf("  ✓ Library already exists: %s\n", libName)
		return nil
	}

	release, err := getLatestWhisperRelease()
	if err != nil {
		return fmt.Errorf("failed to get latest release: %w", err)
	}

	fmt.Printf("  Latest release: %s\n", release.TagName)

	assetURL := ""
	var assetSize int64
	for _, asset := range release.Assets {
		if asset.Name == libName {
			assetURL = asset.BrowserDownloadURL
			assetSize = asset.Size
			break
		}
	}

	if assetURL == "" {
		return fmt.Errorf("no library found for platform %s (looking for %s)", runtime.GOOS, libName)
	}

	fmt.Printf("  Downloading %s (%s)...\n", libName, HumanSize(assetSize))

	logger := &DownloadProgressLogger{}
	progressFunc := func(src string, currentSize, totalSize int64, mibPerSec float64, complete bool) {
		logger.Log(currentSize, totalSize)
	}

	_, err = downloader.Download(ctx, assetURL, libPath, progressFunc, downloader.SizeIntervalMIB)
	if err != nil {
		return fmt.Errorf("failed to download library: %w", err)
	}

	fmt.Printf("  ✓ Downloaded %s\n", libName)
	return nil
}

// getLatestWhisperRelease fetches the latest release info from GitHub
func getLatestWhisperRelease() (*whisperRelease, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(whisperLibReleaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch release info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	var release whisperRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, fmt.Errorf("failed to decode release info: %w", err)
	}

	return &release, nil
}

// downloadRequiredModels downloads the required whisper models
func downloadRequiredModels(ctx context.Context, opts *SetupOptions) error {
	fmt.Println("\nStep 3: Downloading models...")

	// Determine which models to download
	modelsToDownload := opts.DownloadModels

	if len(modelsToDownload) == 0 && opts.Interactive {
		// Interactive mode: ask user which models to download
		scanner := bufio.NewScanner(os.Stdin)

		// Show available models
		fmt.Println("\n  Available models:")
		models := GetAllModels()
		for i, model := range models {
			fmt.Printf("    %d. %-20s (%s) - %s\n", i+1, model.Name, HumanSize(model.Size), model.Description)
		}

		fmt.Print("\n  Enter model numbers to download (comma-separated, or 'all'): ")
		scanner.Scan()
		response := strings.ToLower(strings.TrimSpace(scanner.Text()))

		if response == "all" {
			// Download all models
			for _, model := range models {
				modelsToDownload = append(modelsToDownload, model.Name)
			}
		} else if response != "" {
			// Parse model numbers
			numbers := strings.Split(response, ",")
			for _, numStr := range numbers {
				numStr = strings.TrimSpace(numStr)
				if numStr == "" {
					continue
				}

				var num int
				_, err := fmt.Sscanf(numStr, "%d", &num)
				if err != nil || num < 1 || num > len(models) {
					fmt.Printf("  Warning: Invalid model number '%s'\n", numStr)
					continue
				}

				modelName := models[num-1].Name
				if !contains(modelsToDownload, modelName) {
					modelsToDownload = append(modelsToDownload, modelName)
				}
			}
		}
	} else if opts.AutoSelect {
		// Auto-select optimal model
		spec := SelectOptimalModel(4) // Assume 4GB RAM
		modelsToDownload = append(modelsToDownload, spec.Name)
		fmt.Printf("  Auto-selected optimal model: %s\n", spec.Name)
	}

	// Download models
	if len(modelsToDownload) == 0 {
		fmt.Println("  No models to download")
		return nil
	}

	fmt.Printf("  Downloading %d models...\n", len(modelsToDownload))

	for _, modelName := range modelsToDownload {
		spec, exists := GetModelSpec(modelName)
		if !exists {
			fmt.Printf("  ✗ Unknown model: %s\n", modelName)
			continue
		}

		if IsModelDownloaded(modelName) {
			fmt.Printf("  ✓ %s already downloaded\n", modelName)
			continue
		}

		fmt.Printf("  Downloading %s (%s)...\n", modelName, spec.Description)
		logger := &DownloadProgressLogger{}

		if err := DownloadModelWithProgress(ctx, modelName, logger.Log); err != nil {
			fmt.Printf("  ✗ Failed to download %s: %v\n", modelName, err)
			continue
		}

		fmt.Printf("  ✓ Downloaded %s\n", modelName)
	}

	return nil
}

// verifySetup verifies that the setup is complete and working
func verifySetup(opts *SetupOptions) error {
	fmt.Println("\nStep 4: Verifying setup...")

	// Check models directory
	if _, err := os.Stat(opts.ModelsDir); os.IsNotExist(err) {
		return fmt.Errorf("models directory does not exist: %s", opts.ModelsDir)
	}
	fmt.Println("  ✓ Models directory exists")

	// Check library directory
	if _, err := os.Stat(opts.LibDir); os.IsNotExist(err) {
		return fmt.Errorf("library directory does not exist: %s", opts.LibDir)
	}
	fmt.Println("  ✓ Library directory exists")

	// Check for library file
	libName := whisper.LibraryName(runtime.GOOS)
	libPath := fmt.Sprintf("%s/%s", opts.LibDir, libName)
	if _, err := os.Stat(libPath); os.IsNotExist(err) {
		fmt.Printf("  ⚠ Library not found: %s\n", libName)
	} else {
		fmt.Printf("  ✓ Library found: %s\n", libName)
	}

	// Check for downloaded models
	downloadedModels, err := ListDownloadedModels()
	if err != nil {
		return fmt.Errorf("failed to list downloaded models: %w", err)
	}

	if len(downloadedModels) == 0 {
		fmt.Println("  ⚠ No models downloaded yet")
	} else {
		fmt.Printf("  ✓ Found %d downloaded models\n", len(downloadedModels))
		for _, model := range downloadedModels {
			spec, exists := GetModelSpec(model)
			if exists {
				fmt.Printf("    - %s (%s)\n", model, HumanSize(spec.Size))
			} else {
				fmt.Printf("    - %s (unknown size)\n", model)
			}
		}
	}

	return nil
}

// QuickSetup performs a quick non-interactive setup
func QuickSetup(ctx context.Context, modelName string) error {
	opts := DefaultSetupOptions()
	opts.Interactive = false
	opts.DownloadModels = []string{modelName}
	return SetupWhisper(ctx, opts)
}

// FullSetup performs a full interactive setup
func FullSetup(ctx context.Context) error {
	opts := DefaultSetupOptions()
	opts.Interactive = true
	return SetupWhisper(ctx, opts)
}

// SetupForProduction performs setup optimized for production use
// Uses paths.WhisperModels() and paths.WhisperLib() for consistent path resolution
func SetupForProduction(ctx context.Context, options ...func(*SetupOptions)) error {
	opts := &SetupOptions{
		ModelsDir:      paths.WhisperModels(),
		LibDir:         paths.WhisperLib(),
		DownloadModels: []string{"base"}, // Default to base model
		AutoSelect:     false,
		Interactive:    false,
	}

	// Apply options
	for _, opt := range options {
		opt(opts)
	}

	return SetupWhisper(ctx, opts)
}

// SetupWithModels sets up whisper with specific models
func SetupWithModels(ctx context.Context, models []string) error {
	opts := DefaultSetupOptions()
	opts.Interactive = false
	opts.DownloadModels = models
	return SetupWhisper(ctx, opts)
}

// contains checks if a string is in a slice
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// PrintSetupInstructions prints instructions for manual setup
func PrintSetupInstructions() {
	modelsDir := paths.WhisperModels()
	libDir := paths.WhisperLib()

	fmt.Println("=== Whisper Setup Instructions ===")
	fmt.Println()
	fmt.Printf("Paths (managed by internal/paths):\n")
	fmt.Printf("  Models: %s\n", modelsDir)
	fmt.Printf("  Library: %s\n\n", libDir)
	fmt.Println("1. Create directories:")
	fmt.Printf("   mkdir -p %s\n", modelsDir)
	fmt.Printf("   mkdir -p %s\n", libDir)
	fmt.Println()
	fmt.Println("2. Download models (choose one):")
	fmt.Println("   # Tiny (fastest, lowest quality)")
	fmt.Printf("   curl -L https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-tiny.bin \\\n")
	fmt.Printf("     -o %s/tiny.bin\n", modelsDir)
	fmt.Println()
	fmt.Println("   # Base (recommended, balanced)")
	fmt.Printf("   curl -L https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-base.bin \\\n")
	fmt.Printf("     -o %s/base.bin\n", modelsDir)
	fmt.Println()
	fmt.Println("   # Small (better quality, slower)")
	fmt.Printf("   curl -L https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-small.bin \\\n")
	fmt.Printf("     -o %s/small.bin\n", modelsDir)
	fmt.Println()
	fmt.Println("3. Download the whisper library:")
	libName := whisper.LibraryName(runtime.GOOS)
	fmt.Printf("   curl -L https://github.com/kawai-network/whisper/releases/latest/download/%s \\\n", libName)
	fmt.Printf("     -o %s/%s\n", libDir, libName)
	fmt.Println()
	fmt.Println("4. For production, consider:")
	fmt.Println("   - Using larger models (medium, large-v3)")
	fmt.Println("   - Setting up proper permissions")
	fmt.Println("   - Using absolute paths")
	fmt.Println("   - Monitoring disk space")
	fmt.Println("===============================")
	fmt.Println()
}

// DiagnoseSetup checks for common setup issues
func DiagnoseSetup(opts *SetupOptions) error {
	fmt.Println("=== Whisper Setup Diagnostics ===")
	fmt.Println()

	issues := 0

	// Check directories
	fmt.Println("Checking directories:")
	for _, dir := range []string{opts.ModelsDir, opts.LibDir} {
		if info, err := os.Stat(dir); err != nil {
			if os.IsNotExist(err) {
				fmt.Printf("  ✗ %s does not exist\n", dir)
				issues++
			} else {
				fmt.Printf("  ✗ %s: %v\n", dir, err)
				issues++
			}
		} else {
			if !info.IsDir() {
				fmt.Printf("  ✗ %s is not a directory\n", dir)
				issues++
			} else {
				// Check permissions
				if info.Mode().Perm()&0200 == 0 { // No write permission
					fmt.Printf("  ⚠ %s is not writable\n", dir)
					issues++
				} else {
					fmt.Printf("  ✓ %s exists and is writable\n", dir)
				}
			}
		}
	}

	// Check for models
	fmt.Println("\nChecking models:")
	downloadedModels, err := ListDownloadedModels()
	if err != nil {
		fmt.Printf("  ✗ Failed to list models: %v\n", err)
		issues++
	} else {
		if len(downloadedModels) == 0 {
			fmt.Println("  ⚠ No models downloaded")
			issues++
		} else {
			fmt.Printf("  ✓ Found %d models\n", len(downloadedModels))
			for _, model := range downloadedModels {
				spec, exists := GetModelSpec(model)
				if exists {
					path, pathErr := GetModelFilePath(model)
					if pathErr != nil {
						fmt.Printf("    ✗ %s: %v\n", model, pathErr)
						issues++
					} else {
						info, statErr := os.Stat(path)
						if statErr != nil {
							fmt.Printf("    ✗ %s: %v\n", model, statErr)
							issues++
						} else {
							if info.Size() != spec.Size {
								fmt.Printf("    ⚠ %s: size mismatch (expected %s, got %s)\n",
									model, HumanSize(spec.Size), HumanSize(info.Size()))
								issues++
							} else {
								fmt.Printf("    ✓ %s (%s)\n", model, HumanSize(spec.Size))
							}
						}
					}
				}
			}
		}
	}

	// Check FFmpeg
	fmt.Println("\nChecking dependencies:")
	ffmpegPath, err := exec.LookPath("ffmpeg")
	if err != nil {
		fmt.Println("  ✗ FFmpeg not found in PATH")
		fmt.Println("    Install with: brew install ffmpeg (macOS) or sudo apt install ffmpeg (Linux)")
		issues++
	} else {
		fmt.Printf("  ✓ FFmpeg found at %s\n", ffmpegPath)
	}

	// Check whisper library
	fmt.Println("\nChecking whisper library:")
	libName := whisper.LibraryName(runtime.GOOS)
	libPath := fmt.Sprintf("%s/%s", opts.LibDir, libName)
	if _, err := os.Stat(libPath); os.IsNotExist(err) {
		fmt.Printf("  ✗ Library not found: %s\n", libPath)
		issues++
	} else {
		fmt.Printf("  ✓ Library found: %s\n", libPath)
	}

	// Summary
	fmt.Println("\n=== Diagnostics Summary ===")
	if issues == 0 {
		fmt.Println("✓ No issues found! Whisper is properly configured.")
	} else {
		fmt.Printf("✗ Found %d issue(s) that need to be addressed.\n", issues)
		fmt.Println("\nRun the following to fix issues:")
		fmt.Println("  go run cmd/server/main.go setup whisper")
	}
	fmt.Println("===========================")
	fmt.Println()

	return nil
}
