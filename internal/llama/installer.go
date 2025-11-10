package llama

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/hybridgroup/yzma/pkg/download"
)

// Release represents a GitHub release
type Release struct {
	Version string  `json:"tag_name"`
	Name    string  `json:"name"`
	Body    string  `json:"body"`
	Assets  []Asset `json:"assets"`
	URL     string  `json:"html_url"`
}

// Asset represents a release asset
type Asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
	Size               int64  `json:"size"`
}

// LlamaCppInstaller handles llama.cpp installation and management
type LlamaCppInstaller struct {
	// BinaryPath is where the llama.cpp library is stored locally
	BinaryPath string
	// MetadataPath is where version metadata and cache are stored
	MetadataPath string
}

// LlamaCppReleaseManager is an alias for backward compatibility
// Deprecated: Use LlamaCppInstaller instead
type LlamaCppReleaseManager = LlamaCppInstaller

// NewLlamaCppInstaller creates a new llama.cpp installer
func NewLlamaCppInstaller() *LlamaCppInstaller {
	homeDir, _ := os.UserHomeDir()
	basePath := filepath.Join(homeDir, ".llama-cpp")
	binaryPath := filepath.Join(basePath, "bin")
	metadataPath := filepath.Join(basePath, "metadata")

	return &LlamaCppInstaller{
		BinaryPath:   binaryPath,
		MetadataPath: metadataPath,
	}
}

// GetLatestRelease fetches the latest release information from GitHub with retry logic and rate limiting
func (lcm *LlamaCppInstaller) GetLatestRelease() (*Release, error) {
	// Check if we have a cached release (within last hour)
	if cachedRelease := lcm.getCachedRelease(); cachedRelease != nil {
		log.Printf("Using cached llama.cpp release: %s", cachedRelease.Version)
		return cachedRelease, nil
	}

	// Use the download package to get the latest version (with built-in retry logic)
	version, err := download.LlamaLatestVersion()
	if err != nil {
		return nil, fmt.Errorf("failed to fetch latest llama.cpp version: %w", err)
	}

	// Create a release object with the version
	release := &Release{
		Version: version,
		Name:    fmt.Sprintf("llama.cpp %s", version),
		Body:    fmt.Sprintf("Release %s", version),
	}

	// Cache successful response
	lcm.cacheRelease(release)
	return release, nil
}

// InstallLlamaCpp installs llama.cpp via direct download
// Works on all platforms (macOS, Linux, Windows)
func (lcm *LlamaCppInstaller) InstallLlamaCpp() error {
	if lcm.IsLlamaCppInstalled() {
		log.Println("llama.cpp is already installed")
		return nil
	}

	// Use DownloadRelease (handles all platforms automatically)
	return lcm.DownloadRelease("", nil)
}

// DownloadRelease downloads a specific version of llama.cpp pre-built binaries
// Note: progressCallback is currently ignored as go-getter handles downloads internally
func (lcm *LlamaCppInstaller) DownloadRelease(version string, progressCallback func(float64)) error {
	// Ensure the binary directory exists
	if err := os.MkdirAll(lcm.BinaryPath, 0755); err != nil {
		return fmt.Errorf("failed to create binary directory: %w", err)
	}

	// Progress callback not supported by download package
	if progressCallback != nil {
		log.Printf("Warning: Progress callback not supported by download package")
	}

	// Get latest version if not specified
	if version == "" {
		release, err := lcm.GetLatestRelease()
		if err != nil {
			return fmt.Errorf("failed to get latest release: %w", err)
		}
		version = release.Version
	}

	// Detect processor type
	processor := lcm.detectProcessor()

	log.Printf("Installing llama.cpp %s for %s/%s", version, runtime.GOOS, processor)

	// Use the download package (handles download, extraction, everything)
	if err := download.Get(runtime.GOOS, processor, version, lcm.BinaryPath); err != nil {
		return fmt.Errorf("failed to download llama.cpp: %w", err)
	}

	// Make binaries executable on Unix systems
	if err := lcm.makeExecutable(); err != nil {
		log.Printf("Warning: failed to make binaries executable: %v", err)
	}

	// Save version metadata
	if err := lcm.saveVersionMetadata(version); err != nil {
		log.Printf("Warning: failed to save version metadata: %v", err)
	}

	log.Printf("Successfully installed llama.cpp %s", version)
	return nil
}

// detectProcessor detects the best processor type for this system
func (lcm *LlamaCppInstaller) detectProcessor() string {
	hardware := lcm.detectHardwareCapabilities()

	// Priority: CUDA > Vulkan > Metal > CPU
	if hardware.HasCUDA {
		return "cuda"
	}
	if hardware.HasVulkan {
		return "vulkan"
	}
	if runtime.GOOS == "darwin" {
		return "metal" // macOS always supports Metal
	}
	return "cpu"
}

// HardwareCapabilities represents detected hardware capabilities
type HardwareCapabilities struct {
	OS        string
	Arch      string
	HasNVIDIA bool
	HasAMD    bool
	HasIntel  bool
	HasCUDA   bool
	HasVulkan bool
	HasOpenCL bool
	HasAVX2   bool
}

// makeExecutable makes all binaries in the binary directory executable (Unix systems)
func (lcm *LlamaCppInstaller) makeExecutable() error {
	if runtime.GOOS == "windows" {
		return nil // No need to set executable permissions on Windows
	}

	entries, err := os.ReadDir(lcm.BinaryPath)
	if err != nil {
		return fmt.Errorf("failed to read binary directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(lcm.BinaryPath, entry.Name())
			if err := os.Chmod(filePath, 0755); err != nil {
				log.Printf("Warning: failed to make %s executable: %v", entry.Name(), err)
			}
		}
	}

	return nil
}

// saveVersionMetadata saves the installed version information
func (lcm *LlamaCppInstaller) saveVersionMetadata(version string) error {
	// Ensure metadata directory exists
	if err := os.MkdirAll(lcm.MetadataPath, 0755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	// Create metadata structure
	metadata := struct {
		Version     string    `json:"version"`
		InstalledAt time.Time `json:"installed_at"`
		BinaryPath  string    `json:"binary_path"`
	}{
		Version:     version,
		InstalledAt: time.Now(),
		BinaryPath:  lcm.BinaryPath,
	}

	// Marshal to JSON
	metadataData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Save metadata file
	metadataPath := filepath.Join(lcm.MetadataPath, "installed-version.json")
	if err := os.WriteFile(metadataPath, metadataData, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	log.Printf("Version metadata saved: %s", version)
	return nil
}

// GetInstalledVersion returns the currently installed version
func (lcm *LlamaCppInstaller) GetInstalledVersion() string {
	// Get version from metadata
	return lcm.loadVersionMetadata()
}

// loadVersionMetadata loads the installed version from metadata file
func (lcm *LlamaCppInstaller) loadVersionMetadata() string {
	metadataPath := filepath.Join(lcm.MetadataPath, "installed-version.json")

	// Check if metadata file exists
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		return ""
	}

	// Read metadata file
	metadataData, err := os.ReadFile(metadataPath)
	if err != nil {
		log.Printf("Failed to read metadata file: %v", err)
		return ""
	}

	// Parse metadata
	var metadata struct {
		Version     string    `json:"version"`
		InstalledAt time.Time `json:"installed_at"`
		BinaryPath  string    `json:"binary_path"`
	}

	if err := json.Unmarshal(metadataData, &metadata); err != nil {
		log.Printf("Failed to parse metadata file: %v", err)
		return ""
	}

	// Verify that the binary path in metadata matches current binary path
	if metadata.BinaryPath != lcm.BinaryPath {
		log.Printf("Binary path mismatch in metadata, ignoring")
		return ""
	}

	log.Printf("Loaded version from metadata: %s (installed at %s)", metadata.Version, metadata.InstalledAt.Format("2006-01-02 15:04:05"))
	return metadata.Version
}

// GetServerBinaryPath returns the path to the llama-server binary
// First checks system PATH (for package manager installations), then local binary path
func (lcm *LlamaCppInstaller) GetServerBinaryPath() string {
	binaryName := "llama-server"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	// First, check if llama-server is in system PATH (e.g., from Homebrew)
	if path, err := exec.LookPath(binaryName); err == nil {
		return path
	}

	// Fallback to local binary path (from GitHub release download)
	return filepath.Join(lcm.BinaryPath, binaryName)
}

// GetBinaryPath returns the path to a specific llama.cpp binary
// Platform-specific implementation in installer_*.go files

// IsLlamaCppInstalled checks if llama.cpp library is installed
func (lcm *LlamaCppInstaller) IsLlamaCppInstalled() bool {
	// Check if the library file exists (libllama.so, llama.dll, libllama.dylib)
	libraryName := download.LibraryName(runtime.GOOS)
	if libraryName == "unknown" {
		return false
	}

	libraryPath := filepath.Join(lcm.BinaryPath, libraryName)
	if _, err := os.Stat(libraryPath); err != nil {
		return false
	}

	return true
}

// VerifyInstalledBinary verifies the installed library
func (lcm *LlamaCppInstaller) VerifyInstalledBinary() error {
	// Check if the library file exists (libllama.so, llama.dll, libllama.dylib)
	libraryName := download.LibraryName(runtime.GOOS)
	if libraryName == "unknown" {
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	libraryPath := filepath.Join(lcm.BinaryPath, libraryName)
	if _, err := os.Stat(libraryPath); err != nil {
		return fmt.Errorf("library file not found: %s", libraryName)
	}

	return nil
}

// IsUpdateAvailable checks if an update is available
func (lcm *LlamaCppInstaller) IsUpdateAvailable() (bool, string, error) {
	log.Printf("IsUpdateAvailable: checking for updates...")

	latest, err := lcm.GetLatestRelease()
	if err != nil {
		log.Printf("IsUpdateAvailable: failed to get latest release: %v", err)
		return false, "", err
	}
	log.Printf("IsUpdateAvailable: latest release version: %s", latest.Version)

	current := lcm.GetInstalledVersion()
	log.Printf("IsUpdateAvailable: current installed version: %s", current)

	if current == "" {
		// No version installed, offer to download latest
		log.Printf("IsUpdateAvailable: no version installed, offering latest: %s", latest.Version)
		return true, latest.Version, nil
	}

	updateAvailable := latest.Version != current
	log.Printf("IsUpdateAvailable: update available: %v (latest: %s, current: %s)", updateAvailable, latest.Version, current)
	return updateAvailable, latest.Version, nil
}

// CleanupPartialDownloads removes any partial or corrupted downloads
func (lcm *LlamaCppInstaller) CleanupPartialDownloads() error {
	// Check library file
	if err := lcm.VerifyInstalledBinary(); err != nil {
		log.Printf("Found corrupted library, removing: %v", err)
		// Remove all files in the binary directory
		if removeErr := os.RemoveAll(lcm.BinaryPath); removeErr != nil {
			log.Printf("Failed to remove binary directory: %v", removeErr)
		}
		// Clear metadata since library is corrupted
		lcm.clearVersionMetadata()
	}

	return nil
}

// clearVersionMetadata clears the version metadata (used when binaries are corrupted or removed)
func (lcm *LlamaCppInstaller) clearVersionMetadata() {
	metadataPath := filepath.Join(lcm.MetadataPath, "installed-version.json")
	if err := os.Remove(metadataPath); err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to clear version metadata: %v", err)
	} else {
		log.Printf("Version metadata cleared")
	}
}

// GetAvailableBinaries returns a list of available llama.cpp binaries
func (lcm *LlamaCppInstaller) GetAvailableBinaries() ([]string, error) {
	entries, err := os.ReadDir(lcm.BinaryPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil // No binaries installed
		}
		return nil, fmt.Errorf("failed to read binary directory: %w", err)
	}

	var binaries []string
	for _, entry := range entries {
		if !entry.IsDir() {
			name := entry.Name()
			// Remove .exe extension for cross-platform consistency
			if runtime.GOOS == "windows" && strings.HasSuffix(name, ".exe") {
				name = strings.TrimSuffix(name, ".exe")
			}
			binaries = append(binaries, name)
		}
	}

	return binaries, nil
}

// HasBinary checks if a specific binary is installed
func (lcm *LlamaCppInstaller) HasBinary(binaryName string) bool {
	binaryPath := lcm.GetBinaryPath(binaryName)
	_, err := os.Stat(binaryPath)
	return err == nil
}

// RunBinary executes a llama.cpp binary with the given arguments
func (lcm *LlamaCppInstaller) RunBinary(binaryName string, args []string) (*exec.Cmd, error) {
	binaryPath := lcm.GetBinaryPath(binaryName)

	// Verify binary exists
	if _, err := os.Stat(binaryPath); err != nil {
		return nil, fmt.Errorf("binary %s not found: %w", binaryName, err)
	}

	// Create command
	cmd := exec.Command(binaryPath, args...)

	return cmd, nil
}

// getCachedRelease returns a cached release if it's still valid (within 1 hour)
func (lcm *LlamaCppInstaller) getCachedRelease() *Release {
	cachePath := filepath.Join(lcm.MetadataPath, "release-cache.json")

	// Check if cache file exists
	info, err := os.Stat(cachePath)
	if err != nil {
		return nil
	}

	// Check if cache is still valid (1 hour)
	if time.Since(info.ModTime()) > 1*time.Hour {
		return nil
	}

	// Read and parse cached release
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil
	}

	var release Release
	if err := json.Unmarshal(data, &release); err != nil {
		return nil
	}

	return &release
}

// cacheRelease saves a release to cache
func (lcm *LlamaCppInstaller) cacheRelease(release *Release) {
	// Ensure metadata directory exists
	os.MkdirAll(lcm.MetadataPath, 0755)

	cachePath := filepath.Join(lcm.MetadataPath, "release-cache.json")

	data, err := json.MarshalIndent(release, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal llama.cpp release cache: %v", err)
		return
	}

	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		log.Printf("Failed to write llama.cpp release cache: %v", err)
	}
}

// QuantizeModel quantizes a model using llama-quantize binary
func (lcm *LlamaCppInstaller) QuantizeModel(inputPath, outputPath, quantType string) error {
	// Verify input model exists
	if _, err := os.Stat(inputPath); err != nil {
		return fmt.Errorf("input model file not found: %w", err)
	}

	// Get quantize binary path
	quantizePath := lcm.GetBinaryPath("llama-quantize")
	if _, err := os.Stat(quantizePath); err != nil {
		return fmt.Errorf("llama-quantize binary not found: %w", err)
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Build command arguments
	args := []string{inputPath, outputPath, quantType}

	// Create and execute command
	cmd := exec.Command(quantizePath, args...)

	log.Printf("Starting model quantization: %s -> %s (type: %s)", inputPath, outputPath, quantType)

	// Run the command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("quantization failed: %w\nOutput: %s", err, string(output))
	}

	log.Printf("Model quantization completed successfully: %s", outputPath)
	return nil
}

// BenchmarkModel runs performance benchmarks on a model using llama-bench
func (lcm *LlamaCppInstaller) BenchmarkModel(modelPath string, options map[string]interface{}) (*BenchmarkResults, error) {
	// Verify model exists
	if _, err := os.Stat(modelPath); err != nil {
		return nil, fmt.Errorf("model file not found: %w", err)
	}

	// Get bench binary path
	benchPath := lcm.GetBinaryPath("llama-bench")
	if _, err := os.Stat(benchPath); err != nil {
		return nil, fmt.Errorf("llama-bench binary not found: %w", err)
	}

	// Build command arguments
	args := []string{"-m", modelPath}

	// Add optional parameters
	if threads, ok := options["threads"].(int); ok {
		args = append(args, "-t", fmt.Sprintf("%d", threads))
	}
	if ngl, ok := options["n_gpu_layers"].(int); ok {
		args = append(args, "-ngl", fmt.Sprintf("%d", ngl))
	}
	if batchSize, ok := options["batch_size"].(int); ok {
		args = append(args, "-b", fmt.Sprintf("%d", batchSize))
	}

	// Create and execute command
	cmd := exec.Command(benchPath, args...)

	log.Printf("Starting model benchmark: %s", modelPath)

	// Run the command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("benchmark failed: %w\nOutput: %s", err, string(output))
	}

	// Parse benchmark results
	results := lcm.parseBenchmarkOutput(string(output))
	log.Printf("Model benchmark completed: %s", modelPath)

	return results, nil
}

// BenchmarkResults represents the results of a model benchmark
type BenchmarkResults struct {
	ModelPath       string  `json:"model_path"`
	TokensPerSecond float64 `json:"tokens_per_second"`
	PromptTokens    int     `json:"prompt_tokens"`
	GeneratedTokens int     `json:"generated_tokens"`
	TotalTime       float64 `json:"total_time_seconds"`
	MemoryUsage     int64   `json:"memory_usage_mb"`
	ThreadsUsed     int     `json:"threads_used"`
	GPULayers       int     `json:"gpu_layers"`
}

// parseBenchmarkOutput parses the output from llama-bench
func (lcm *LlamaCppInstaller) parseBenchmarkOutput(output string) *BenchmarkResults {
	results := &BenchmarkResults{}

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Parse tokens per second
		if strings.Contains(line, "tokens per second") {
			parts := strings.Fields(line)
			for i, part := range parts {
				if part == "tokens" && i > 0 {
					if tps, err := strconv.ParseFloat(parts[i-1], 64); err == nil {
						results.TokensPerSecond = tps
					}
					break
				}
			}
		}

		// Parse total time
		if strings.Contains(line, "total time") {
			parts := strings.Fields(line)
			for i, part := range parts {
				if strings.Contains(part, "ms") && i >= 0 {
					timeStr := strings.TrimSuffix(part, "ms")
					if timeMs, err := strconv.ParseFloat(timeStr, 64); err == nil {
						results.TotalTime = timeMs / 1000.0 // Convert to seconds
					}
					break
				}
			}
		}

		// Parse memory usage
		if strings.Contains(line, "memory") && strings.Contains(line, "MB") {
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.Contains(part, "MB") {
					memStr := strings.TrimSuffix(part, "MB")
					if mem, err := strconv.ParseInt(memStr, 10, 64); err == nil {
						results.MemoryUsage = mem
					}
					break
				}
			}
		}
	}

	return results
}

// ProcessImageWithText processes an image with text using llama-llava-cli
func (lcm *LlamaCppInstaller) ProcessImageWithText(imagePath, prompt, modelPath string) (string, error) {
	// Verify files exist
	if _, err := os.Stat(imagePath); err != nil {
		return "", fmt.Errorf("image file not found: %w", err)
	}
	if _, err := os.Stat(modelPath); err != nil {
		return "", fmt.Errorf("model file not found: %w", err)
	}

	// Get llava binary path
	llavaPath := lcm.GetBinaryPath("llama-llava-cli")
	if _, err := os.Stat(llavaPath); err != nil {
		return "", fmt.Errorf("llama-llava-cli binary not found: %w", err)
	}

	// Build command arguments
	args := []string{
		"-m", modelPath,
		"--image", imagePath,
		"-p", prompt,
		"--temp", "0.1",
		"-n", "512",
	}

	// Create and execute command
	cmd := exec.Command(llavaPath, args...)

	log.Printf("Processing image with LLaVA: %s", imagePath)

	// Run the command and capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("image processing failed: %w\nOutput: %s", err, string(output))
	}

	// Extract the response from the output
	response := lcm.extractLlavaResponse(string(output))
	log.Printf("Image processing completed successfully")

	return response, nil
}

// extractLlavaResponse extracts the actual response from llava output
func (lcm *LlamaCppInstaller) extractLlavaResponse(output string) string {
	lines := strings.Split(output, "\n")
	var responseLines []string
	inResponse := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Look for the start of the actual response
		if strings.Contains(line, "assistant") || strings.Contains(line, ">") {
			inResponse = true
			continue
		}

		// Skip system messages and prompts
		if strings.HasPrefix(line, "system:") ||
			strings.HasPrefix(line, "user:") ||
			strings.HasPrefix(line, "llama_") ||
			strings.HasPrefix(line, "encode_") ||
			strings.Contains(line, "sampling parameters") {
			continue
		}

		// Collect response lines
		if inResponse && line != "" {
			responseLines = append(responseLines, line)
		}
	}

	response := strings.Join(responseLines, " ")
	return strings.TrimSpace(response)
}

// GetAvailableQuantizationTypes returns the supported quantization types
func (lcm *LlamaCppInstaller) GetAvailableQuantizationTypes() []string {
	return []string{
		"q4_0",   // 4-bit quantization (smallest)
		"q4_1",   // 4-bit quantization (better quality)
		"q5_0",   // 5-bit quantization
		"q5_1",   // 5-bit quantization (better quality)
		"q8_0",   // 8-bit quantization (good balance)
		"q2_k",   // 2-bit k-quantization
		"q3_k_s", // 3-bit k-quantization (small)
		"q3_k_m", // 3-bit k-quantization (medium)
		"q3_k_l", // 3-bit k-quantization (large)
		"q4_k_s", // 4-bit k-quantization (small)
		"q4_k_m", // 4-bit k-quantization (medium)
		"q5_k_s", // 5-bit k-quantization (small)
		"q5_k_m", // 5-bit k-quantization (medium)
		"q6_k",   // 6-bit k-quantization
		"f16",    // 16-bit float (highest quality)
		"f32",    // 32-bit float (original precision)
	}
}

// GetModelsDirectory returns the directory where models are stored
func (lcm *LlamaCppInstaller) GetModelsDirectory() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".llama-cpp", "models")
}
