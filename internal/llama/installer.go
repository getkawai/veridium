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

	"github.com/kawai-network/veridium/pkg/yzma/download"
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

// LlamaCppInstaller handles llama.cpp installation and model downloads
type LlamaCppInstaller struct {
	// BinaryPath is where the llama.cpp library is stored locally
	BinaryPath string
	// MetadataPath is where version metadata and cache are stored
	MetadataPath string
	// ModelsDir is where GGUF models (chat & embedding) are stored
	ModelsDir string
}

// NewLlamaCppInstaller creates a new llama.cpp installer
// Automatically cleans up any stale temporary files from previous failed downloads
func NewLlamaCppInstaller() *LlamaCppInstaller {
	homeDir, _ := os.UserHomeDir()
	basePath := filepath.Join(homeDir, ".llama-cpp")
	binaryPath := filepath.Join(basePath, "bin")
	metadataPath := filepath.Join(basePath, "metadata")
	modelsDir := filepath.Join(basePath, "models")

	installer := &LlamaCppInstaller{
		BinaryPath:   binaryPath,
		MetadataPath: metadataPath,
		ModelsDir:    modelsDir,
	}

	// Clean up any stale temp files from previous sessions
	// This handles the case where app was closed during download
	if err := installer.CleanupStaleTempFiles(); err != nil {
		log.Printf("⚠️  Failed to cleanup stale temp files on startup: %v", err)
	}

	return installer
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
// Uses pkg/yzma/download which now uses grab internally and handles:
// - Platform-specific binary URLs
// - ZIP download with grab (resume support!)
// - Automatic ZIP extraction
// - Built-in retry logic
// Note: progressCallback is currently ignored as grab handles downloads internally
func (lcm *LlamaCppInstaller) DownloadRelease(version string, progressCallback func(float64)) error {
	// Ensure the binary directory exists
	if err := os.MkdirAll(lcm.BinaryPath, 0755); err != nil {
		return fmt.Errorf("failed to create binary directory: %w", err)
	}

	// Progress callback not supported yet (grab handles progress internally)
	if progressCallback != nil {
		log.Printf("Warning: Progress callback not supported yet (grab handles progress internally)")
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

	// Download llama.cpp binaries using pkg/yzma/download
	// This now uses grab internally for ZIP download (with resume support!)
	// Then extracts the ZIP and handles platform-specific URLs
	// For model downloads, we use grab directly (see downloader.go)
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

// IsLlamaCppInstalled checks if all required llama.cpp libraries are installed
// Verifies that libggml, libggml-base, and libllama all exist
func (lcm *LlamaCppInstaller) IsLlamaCppInstalled() bool {
	return lcm.VerifyAllLibrariesExist()
}

// VerifyInstalledBinary verifies all required libraries are installed
// Returns an error if any required library (libggml, libggml-base, libllama) is missing
func (lcm *LlamaCppInstaller) VerifyInstalledBinary() error {
	requiredLibs := download.RequiredLibraries(runtime.GOOS)
	if len(requiredLibs) == 0 {
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	missingLibs := []string{}
	for _, lib := range requiredLibs {
		libPath := filepath.Join(lcm.BinaryPath, lib)
		if _, err := os.Stat(libPath); err != nil {
			missingLibs = append(missingLibs, lib)
		}
	}

	if len(missingLibs) > 0 {
		return fmt.Errorf("missing required libraries: %s", strings.Join(missingLibs, ", "))
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
	return lcm.ModelsDir
}

// GetLibraryPath returns the path to the llama.cpp library directory
// This is the directory containing all required libraries (libggml, libggml-base, libllama)
// This path should be passed to llama.Load() for library-based usage
func (lcm *LlamaCppInstaller) GetLibraryPath() string {
	return lcm.BinaryPath
}

// GetLibraryFilePath returns the full path to the main llama.cpp library file
// Returns the platform-specific library file (libllama.so, libllama.dylib, or llama.dll)
func (lcm *LlamaCppInstaller) GetLibraryFilePath() string {
	libraryName := download.LibraryName(runtime.GOOS)
	return filepath.Join(lcm.BinaryPath, libraryName)
}

// GetRequiredLibraryPaths returns full paths to all required library files
// Returns paths to: libggml, libggml-base, and libllama (platform-specific extensions)
// Use this to verify all required libraries are present before loading
func (lcm *LlamaCppInstaller) GetRequiredLibraryPaths() []string {
	requiredLibs := download.RequiredLibraries(runtime.GOOS)
	paths := make([]string, len(requiredLibs))
	for i, lib := range requiredLibs {
		paths[i] = filepath.Join(lcm.BinaryPath, lib)
	}
	return paths
}

// VerifyAllLibrariesExist checks if all required llama.cpp libraries are present
// Returns true only if libggml, libggml-base, and libllama all exist
func (lcm *LlamaCppInstaller) VerifyAllLibrariesExist() bool {
	requiredPaths := lcm.GetRequiredLibraryPaths()
	for _, path := range requiredPaths {
		if _, err := os.Stat(path); err != nil {
			return false
		}
	}
	return len(requiredPaths) > 0
}

// ============================================================================
// Model Download Methods
// ============================================================================

// DownloadChatModel downloads a chat model (Qwen) using model specs from model_specs.go
// Features:
// - Downloads to temporary file first (.tmp) to prevent corruption
// - Retries up to 3 times on network failure with exponential backoff (2s, 4s, 6s)
// - Validates file size, checksum (if provided), and GGUF format
// - Automatically cleans up partial downloads on failure
// - Only moves to final destination after successful validation
// - Skips download if model already exists and is valid
// - Handles app closure during download (temp files cleaned on next startup)
func (lcm *LlamaCppInstaller) DownloadChatModel(modelSpec QwenModelSpec) error {
	// Ensure models directory exists
	if err := os.MkdirAll(lcm.ModelsDir, 0755); err != nil {
		return fmt.Errorf("failed to create models directory: %w", err)
	}

	// Build model filename and paths
	modelFileName := fmt.Sprintf("%s.gguf", modelSpec.Name)
	destModelPath := filepath.Join(lcm.ModelsDir, modelFileName)
	tempModelPath := destModelPath + ".tmp"

	// Clean up any stale temporary files
	if err := lcm.cleanupTempFile(tempModelPath); err != nil {
		log.Printf("⚠️  Failed to cleanup stale temp file: %v", err)
	}

	// Check if model already exists
	if _, err := os.Stat(destModelPath); err == nil {
		log.Printf("✅ Model already exists: %s", modelFileName)
		// Verify existing model integrity if checksum is provided
		if modelSpec.SHA256 != "" {
			if err := verifyModelChecksum(destModelPath, modelSpec.SHA256); err != nil {
				log.Printf("⚠️  Existing model checksum invalid, re-downloading...")
				if removeErr := os.Remove(destModelPath); removeErr != nil {
					log.Printf("⚠️  Failed to remove invalid model: %v", removeErr)
				}
			} else {
				return nil // Model exists and is valid
			}
		} else {
			return nil // Model exists, no checksum to verify
		}
	}

	log.Printf("📥 Downloading chat model: %s", modelSpec.Name)
	log.Printf("   URL: %s", modelSpec.URL)
	log.Printf("   Expected size: %.1f MB", float64(modelSpec.Size)/(1024*1024))
	log.Printf("   This may take several minutes depending on network speed...")

	// Download using grab with automatic retry, resume, and progress tracking
	opts := download.DefaultDownloadOptions()
	if err := download.GetWithProgress(modelSpec.URL, tempModelPath, opts); err != nil {
		lcm.cleanupTempFile(tempModelPath)
		return fmt.Errorf("failed to download model: %w", err)
	}

	// Verify downloaded file
	if err := validateDownloadedFile(tempModelPath, modelSpec); err != nil {
		lcm.cleanupTempFile(tempModelPath)
		return err
	}

	// Move temporary file to final destination
	if err := os.Rename(tempModelPath, destModelPath); err != nil {
		lcm.cleanupTempFile(tempModelPath)
		return fmt.Errorf("failed to move downloaded file: %w", err)
	}

	fileInfo, _ := os.Stat(destModelPath)
	sizeMB := float64(fileInfo.Size()) / (1024 * 1024)
	log.Printf("✅ Model downloaded successfully: %s (%.1f MB)", modelFileName, sizeMB)

	return nil
}

// DownloadEmbeddingModel downloads an embedding model with automatic retry and cleanup
// Features:
// - Downloads to temporary file first (.tmp) to prevent corruption
// - Retries up to 3 times on network failure with exponential backoff (2s, 4s, 6s)
// - Validates GGUF file structure after download
// - Automatically cleans up partial downloads on failure
// - Only moves to final destination after successful validation
// - Skips download if model already exists
// - Handles app closure during download (temp files cleaned on next startup)
func (lcm *LlamaCppInstaller) DownloadEmbeddingModel(model *EmbeddingModel) error {
	// Ensure models directory exists
	if err := os.MkdirAll(lcm.ModelsDir, 0755); err != nil {
		return fmt.Errorf("failed to create models directory: %w", err)
	}

	finalPath := filepath.Join(lcm.ModelsDir, model.Filename)

	// Check if already downloaded
	if _, err := os.Stat(finalPath); err == nil {
		log.Printf("✅ Embedding model already exists: %s", model.Name)
		return nil
	}

	log.Printf("📥 Downloading embedding model: %s", model.Name)
	log.Printf("   URL: %s", model.URL)
	log.Printf("   Size: %.2f MB", float64(model.Size)/1024/1024)

	// Download using grab with automatic retry, resume, and progress tracking
	tempPath := finalPath + ".tmp"
	lcm.cleanupTempFile(tempPath) // Clean any stale temp file

	opts := download.DefaultDownloadOptions()
	if err := download.GetWithProgress(model.URL, tempPath, opts); err != nil {
		lcm.cleanupTempFile(tempPath)
		return fmt.Errorf("failed to download model: %w", err)
	}

	// Move temp file to final destination
	if err := os.Rename(tempPath, finalPath); err != nil {
		lcm.cleanupTempFile(tempPath)
		return fmt.Errorf("failed to move downloaded file: %w", err)
	}

	// Validate GGUF file structure
	if err := validateEmbeddingGGUFFile(finalPath); err != nil {
		lcm.cleanupTempFile(finalPath)
		return fmt.Errorf("downloaded file failed GGUF validation: %w", err)
	}

	log.Printf("✅ Embedding model downloaded successfully: %s", model.Name)
	return nil
}

// AutoDownloadRecommendedChatModel automatically downloads the best chat model for the system
func (lcm *LlamaCppInstaller) AutoDownloadRecommendedChatModel() error {
	// Clean up any stale temp files
	if err := lcm.CleanupStaleTempFiles(); err != nil {
		log.Printf("⚠️  Failed to cleanup stale temp files: %v", err)
	}

	// Check if any models already exist
	models, err := lcm.GetAvailableChatModels()
	if err != nil {
		return fmt.Errorf("failed to check existing models: %w", err)
	}

	if len(models) > 0 {
		log.Printf("✅ Chat models already available (%d found), skipping auto-download", len(models))
		return nil
	}

	log.Println("📦 No chat models found, starting auto-download...")

	// Detect hardware specs
	specs := DetectHardwareSpecs()

	// Select optimal model based on available RAM
	modelSpec := SelectOptimalQwenModel(specs.AvailableRAM)

	// Download the model
	if err := lcm.DownloadChatModel(modelSpec); err != nil {
		return fmt.Errorf("failed to download model: %w", err)
	}

	log.Println("🎉 Chat model download completed successfully!")
	return nil
}

// AutoDownloadRecommendedEmbeddingModel automatically downloads the recommended embedding model
func (lcm *LlamaCppInstaller) AutoDownloadRecommendedEmbeddingModel() error {
	downloaded := lcm.GetDownloadedEmbeddingModels()
	if len(downloaded) > 0 {
		log.Printf("✅ Embedding models already available (%d found), skipping auto-download", len(downloaded))
		return nil
	}

	log.Println("📦 No embedding models found, starting auto-download...")

	modelName := GetRecommendedEmbeddingModel()
	model, exists := GetEmbeddingModel(modelName)
	if !exists {
		return fmt.Errorf("recommended model not found: %s", modelName)
	}

	if err := lcm.DownloadEmbeddingModel(model); err != nil {
		return fmt.Errorf("failed to download embedding model: %w", err)
	}

	log.Println("🎉 Embedding model download completed successfully!")
	return nil
}

// GetAvailableChatModels returns a list of available chat model file paths
func (lcm *LlamaCppInstaller) GetAvailableChatModels() ([]string, error) {
	if err := os.MkdirAll(lcm.ModelsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create models directory: %w", err)
	}

	entries, err := os.ReadDir(lcm.ModelsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read models directory: %w", err)
	}

	var models []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		if strings.HasSuffix(strings.ToLower(name), ".gguf") {
			// Check if it's a chat model (not embedding model)
			if !lcm.isEmbeddingModel(name) {
				models = append(models, filepath.Join(lcm.ModelsDir, name))
			}
		}
	}

	return models, nil
}

// GetDownloadedEmbeddingModels returns a list of downloaded embedding models
func (lcm *LlamaCppInstaller) GetDownloadedEmbeddingModels() []*EmbeddingModel {
	var downloaded []*EmbeddingModel
	catalog := GetAvailableEmbeddingModels()

	for _, model := range catalog {
		modelPath := filepath.Join(lcm.ModelsDir, model.Filename)
		if _, err := os.Stat(modelPath); err == nil {
			// Validate the GGUF file structure
			if err := validateEmbeddingGGUFFile(modelPath); err == nil {
				downloaded = append(downloaded, model)
			}
		}
	}

	return downloaded
}

// CleanupStaleTempFiles removes all stale temporary download files (.tmp)
func (lcm *LlamaCppInstaller) CleanupStaleTempFiles() error {
	entries, err := os.ReadDir(lcm.ModelsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Directory doesn't exist yet, nothing to clean
		}
		return fmt.Errorf("failed to read models directory: %w", err)
	}

	cleaned := 0
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Check if it's a temporary file (.tmp)
		if strings.HasSuffix(entry.Name(), ".tmp") {
			tmpPath := filepath.Join(lcm.ModelsDir, entry.Name())

			info, err := entry.Info()
			if err != nil {
				log.Printf("⚠️  Failed to get info for %s: %v", entry.Name(), err)
				continue
			}

			log.Printf("🧹 Removing stale temporary file: %s (size: %.1f MB)",
				entry.Name(), float64(info.Size())/(1024*1024))

			if err := os.Remove(tmpPath); err != nil {
				log.Printf("⚠️  Failed to remove stale temp file %s: %v", entry.Name(), err)
			} else {
				cleaned++
			}
		}
	}

	if cleaned > 0 {
		log.Printf("✅ Cleaned up %d stale temporary file(s)", cleaned)
	}

	return nil
}

// Helper methods for model downloads

func (lcm *LlamaCppInstaller) cleanupTempFile(filePath string) error {
	if err := os.Remove(filePath); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (lcm *LlamaCppInstaller) isEmbeddingModel(filename string) bool {
	catalog := GetAvailableEmbeddingModels()
	for _, model := range catalog {
		if model.Filename == filename {
			return true
		}
	}
	return false
}
