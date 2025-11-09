package llama

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
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

// LlamaCppReleaseManager handles llama.cpp release management
type LlamaCppReleaseManager struct {
	// GitHubOwner is the GitHub repository owner
	GitHubOwner string
	// GitHubRepo is the GitHub repository name
	GitHubRepo string
	// BinaryPath is where the llama.cpp binaries are stored locally
	BinaryPath string
	// SourcePath is where the llama.cpp source code is stored
	SourcePath string
	// CurrentVersion is the currently installed version
	CurrentVersion string
	// ChecksumsPath is where checksums are stored
	ChecksumsPath string
	// MetadataPath is where version metadata is stored
	MetadataPath string
	// BuildPath is where build artifacts are stored
	BuildPath string
}

// NewLlamaCppReleaseManager creates a new llama.cpp release manager
func NewLlamaCppReleaseManager() *LlamaCppReleaseManager {
	homeDir, _ := os.UserHomeDir()
	basePath := filepath.Join(homeDir, ".llama-cpp")
	binaryPath := filepath.Join(basePath, "bin")
	sourcePath := filepath.Join(basePath, "source")
	buildPath := filepath.Join(basePath, "build")
	checksumsPath := filepath.Join(basePath, "checksums")
	metadataPath := filepath.Join(basePath, "metadata")

	return &LlamaCppReleaseManager{
		GitHubOwner:    "ggml-org",
		GitHubRepo:     "llama.cpp",
		BinaryPath:     binaryPath,
		SourcePath:     sourcePath,
		BuildPath:      buildPath,
		CurrentVersion: "",
		ChecksumsPath:  checksumsPath,
		MetadataPath:   metadataPath,
	}
}

// GetLatestRelease fetches the latest release information from GitHub with retry logic and rate limiting
func (lcm *LlamaCppReleaseManager) GetLatestRelease() (*Release, error) {
	// Check if we have a cached release (within last hour)
	if cachedRelease := lcm.getCachedRelease(); cachedRelease != nil {
		log.Printf("Using cached llama.cpp release: %s", cachedRelease.Version)
		return cachedRelease, nil
	}

	// Try GitHub API first
	release, err := lcm.fetchFromGitHubAPI()
	if err != nil {
		log.Printf("llama.cpp GitHub API failed: %v", err)

		// If rate limited, try fallback approach
		if strings.Contains(err.Error(), "rate limit") {
			log.Printf("Attempting llama.cpp fallback release detection...")
			return lcm.getFallbackRelease()
		}

		return nil, err
	}

	// Cache successful response
	lcm.cacheRelease(release)
	return release, nil
}

// fetchFromGitHubAPI attempts to fetch release from GitHub API with rate limiting
func (lcm *LlamaCppReleaseManager) fetchFromGitHubAPI() (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", lcm.GitHubOwner, lcm.GitHubRepo)

	var lastErr error
	maxRetries := 2 // Reduced retries to avoid long waits

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			// Longer exponential backoff for rate limiting
			waitTime := time.Duration(attempt*attempt) * 30 * time.Second
			log.Printf("Waiting %v before retry (attempt %d/%d)", waitTime, attempt+1, maxRetries)
			time.Sleep(waitTime)
		}

		// Create HTTP client with proper configuration
		client := &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				DisableKeepAlives: true, // Prevent connection reuse issues
			},
		}

		// Create request with proper headers
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %w", err)
			continue
		}

		// Add User-Agent header to avoid GitHub API restrictions
		req.Header.Set("User-Agent", "Kawai-Agent/1.0")
		req.Header.Set("Accept", "application/vnd.github.v3+json")

		resp, err := client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to fetch latest release: %w", err)
			continue
		}

		if resp.StatusCode == 403 {
			// Rate limit exceeded - don't retry, use fallback instead
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			lastErr = fmt.Errorf("GitHub API rate limit exceeded: %s", string(body))
			break // Exit retry loop immediately
		}

		if resp.StatusCode != http.StatusOK {
			// Read response body for error details
			body, _ := io.ReadAll(resp.Body)
			resp.Body.Close()
			lastErr = fmt.Errorf("failed to fetch latest release: status %d, body: %s", resp.StatusCode, string(body))

			// Don't retry on 404 or other client errors
			if resp.StatusCode >= 400 && resp.StatusCode < 500 {
				return nil, lastErr
			}
			continue
		}

		// Success - parse the response
		body, err := lcm.parseGitHubResponse(resp)
		resp.Body.Close()
		if err != nil {
			lastErr = err
			continue
		}

		return body, nil
	}

	return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

// getFallbackRelease provides a fallback release when GitHub API is rate limited
func (lcm *LlamaCppReleaseManager) getFallbackRelease() (*Release, error) {
	log.Printf("Using llama.cpp fallback release information...")

	// Create a fallback release with known good assets
	// This is based on the current llama.cpp release structure
	fallbackRelease := &Release{
		Version: "b6451", // Known stable version
		Name:    "llama.cpp b6451",
		Body:    "Fallback release used when GitHub API is rate limited",
		Assets: []Asset{
			{
				Name:               "llama-b6451-bin-macos-arm64.zip",
				BrowserDownloadURL: "https://github.com/ggml-org/llama.cpp/releases/download/b6451/llama-b6451-bin-macos-arm64.zip",
				Size:               3200000, // ~3.1 MB
			},
			{
				Name:               "llama-b6451-bin-macos-x64.zip",
				BrowserDownloadURL: "https://github.com/ggml-org/llama.cpp/releases/download/b6451/llama-b6451-bin-macos-x64.zip",
				Size:               3400000, // ~3.2 MB
			},
			{
				Name:               "llama-b6451-bin-ubuntu-x64.zip",
				BrowserDownloadURL: "https://github.com/ggml-org/llama.cpp/releases/download/b6451/llama-b6451-bin-ubuntu-x64.zip",
				Size:               4100000, // ~3.9 MB
			},
			{
				Name:               "cudart-llama-bin-win-cuda-12.4-x64.zip",
				BrowserDownloadURL: "https://github.com/ggml-org/llama.cpp/releases/download/b6451/cudart-llama-bin-win-cuda-12.4-x64.zip",
				Size:               391500000, // ~373.4 MB
			},
			{
				Name:               "llama-b6451-bin-win-vulkan-x64.zip",
				BrowserDownloadURL: "https://github.com/ggml-org/llama.cpp/releases/download/b6451/llama-b6451-bin-win-vulkan-x64.zip",
				Size:               15600000, // ~14.9 MB
			},
			{
				Name:               "llama-b6451-bin-win-avx2-x64.zip",
				BrowserDownloadURL: "https://github.com/ggml-org/llama.cpp/releases/download/b6451/llama-b6451-bin-win-avx2-x64.zip",
				Size:               4800000, // ~4.6 MB
			},
		},
	}

	// Cache the fallback release for future use
	lcm.cacheRelease(fallbackRelease)

	log.Printf("llama.cpp fallback release created with %d assets", len(fallbackRelease.Assets))
	return fallbackRelease, nil
}

// parseGitHubResponse parses the GitHub API response
func (lcm *LlamaCppReleaseManager) parseGitHubResponse(resp *http.Response) (*Release, error) {
	// Parse the JSON response properly
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	release := &Release{}
	if err := json.Unmarshal(body, release); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	// Validate the release data
	if release.Version == "" {
		return nil, fmt.Errorf("no version found in release data")
	}

	return release, nil
}

// DownloadRelease downloads a specific version of llama.cpp pre-built binaries
func (lcm *LlamaCppReleaseManager) DownloadRelease(version string, progressCallback func(float64)) error {
	// Ensure the binary directory exists
	if err := os.MkdirAll(lcm.BinaryPath, 0755); err != nil {
		return fmt.Errorf("failed to create binary directory: %w", err)
	}

	// Get the release information with assets
	release, err := lcm.GetLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to get release information: %w", err)
	}

	// Find the best binary asset for this platform
	return lcm.downloadBestAsset(release, progressCallback)
}

// downloadBestAsset finds and downloads the best asset for the current platform
func (lcm *LlamaCppReleaseManager) downloadBestAsset(release *Release, progressCallback func(float64)) error {
	if len(release.Assets) == 0 {
		return fmt.Errorf("no assets found in release %s", release.Version)
	}

	// Detect hardware capabilities
	hardware := lcm.detectHardwareCapabilities()

	// Find the best asset based on hardware
	asset := lcm.selectBestAsset(release.Assets, hardware)
	if asset == nil {
		return fmt.Errorf("no compatible binary found for platform %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	log.Printf("Selected asset: %s (%.1f MB)", asset.Name, float64(asset.Size)/(1024*1024))

	// Download the selected asset
	return lcm.downloadAsset(asset, progressCallback)
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

// selectBestAsset selects the best asset based on hardware capabilities
func (lcm *LlamaCppReleaseManager) selectBestAsset(assets []Asset, hardware *HardwareCapabilities) *Asset {
	// Create priority list of preferred asset patterns
	patterns := lcm.getAssetPatterns(hardware)

	// Try each pattern in priority order
	for _, pattern := range patterns {
		for _, asset := range assets {
			if lcm.matchesPattern(asset.Name, pattern, hardware) {
				log.Printf("Matched pattern '%s' with asset '%s'", pattern, asset.Name)
				return &asset
			}
		}
	}

	return nil
}

// matchesPattern checks if an asset name matches a pattern for the given hardware
func (lcm *LlamaCppReleaseManager) matchesPattern(assetName, pattern string, hardware *HardwareCapabilities) bool {
	// Convert to lowercase for case-insensitive matching
	name := strings.ToLower(assetName)

	// Simple pattern matching (could be enhanced with regex if needed)
	parts := strings.Split(pattern, ".*")

	for _, part := range parts {
		if part != "" && !strings.Contains(name, strings.ToLower(part)) {
			return false
		}
	}

	// Additional checks for architecture compatibility
	if hardware.Arch == "arm64" && !strings.Contains(name, "arm64") && strings.Contains(name, "x64") {
		return false
	}

	if hardware.Arch == "amd64" && strings.Contains(name, "arm64") {
		return false
	}

	return true
}

// downloadAsset downloads a specific asset
func (lcm *LlamaCppReleaseManager) downloadAsset(asset *Asset, progressCallback func(float64)) error {
	archivePath := filepath.Join(lcm.BinaryPath, asset.Name)

	log.Printf("Downloading from: %s", asset.BrowserDownloadURL)

	if err := lcm.downloadFile(asset.BrowserDownloadURL, archivePath, progressCallback); err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}

	// Extract the binary archive
	if err := lcm.extractBinaries(archivePath, lcm.BinaryPath); err != nil {
		// Clean up failed download
		os.Remove(archivePath)
		return fmt.Errorf("failed to extract binaries: %w", err)
	}

	// Clean up the archive file
	if err := os.Remove(archivePath); err != nil {
		log.Printf("Warning: failed to remove archive file: %v", err)
	}

	// Make binaries executable on Unix systems
	if err := lcm.makeExecutable(); err != nil {
		log.Printf("Warning: failed to make binaries executable: %v", err)
	}

	// Save version metadata with asset info
	if err := lcm.saveVersionMetadataWithAsset(asset.Name); err != nil {
		log.Printf("Warning: failed to save version metadata: %v", err)
	}

	log.Printf("Successfully downloaded and installed: %s", asset.Name)
	return nil
}

// saveVersionMetadataWithAsset saves the installed version information with asset details
func (lcm *LlamaCppReleaseManager) saveVersionMetadataWithAsset(assetName string) error {
	// Get the release information to extract version
	release, err := lcm.GetLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to get release info: %w", err)
	}

	// Ensure metadata directory exists
	if err := os.MkdirAll(lcm.MetadataPath, 0755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	// Create metadata structure
	metadata := struct {
		Version     string    `json:"version"`
		AssetName   string    `json:"asset_name"`
		InstalledAt time.Time `json:"installed_at"`
		BinaryPath  string    `json:"binary_path"`
		SourcePath  string    `json:"source_path"`
		BuildPath   string    `json:"build_path"`
	}{
		Version:     release.Version,
		AssetName:   assetName,
		InstalledAt: time.Now(),
		BinaryPath:  lcm.BinaryPath,
		SourcePath:  lcm.SourcePath,
		BuildPath:   lcm.BuildPath,
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

	log.Printf("Version metadata saved: %s (asset: %s)", release.Version, assetName)
	return nil
}

// extractBinaries extracts the llama.cpp binaries from the downloaded archive
func (lcm *LlamaCppReleaseManager) extractBinaries(archivePath, destDir string) error {
	// llama.cpp releases are in ZIP format
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open zip archive: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		// Extract all files from the archive
		if file.FileInfo().IsDir() {
			// Create directory
			dirPath := filepath.Join(destDir, file.Name)
			if err := os.MkdirAll(dirPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		// Extract file
		rc, err := file.Open()
		if err != nil {
			return fmt.Errorf("failed to open file in archive: %w", err)
		}

		// Get the base filename (remove any directory structure)
		fileName := filepath.Base(file.Name)
		outputPath := filepath.Join(destDir, fileName)

		outputFile, err := os.Create(outputPath)
		if err != nil {
			rc.Close()
			return fmt.Errorf("failed to create output file: %w", err)
		}

		_, err = io.Copy(outputFile, rc)
		outputFile.Close()
		rc.Close()

		if err != nil {
			return fmt.Errorf("failed to extract file: %w", err)
		}

		log.Printf("Extracted binary: %s", fileName)
	}

	log.Printf("Successfully extracted binaries to %s", destDir)
	return nil
}

// makeExecutable makes all binaries in the binary directory executable (Unix systems)
func (lcm *LlamaCppReleaseManager) makeExecutable() error {
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

// progressReader wraps an io.Reader to report progress
type progressReader struct {
	Reader   io.Reader
	Total    int64
	Current  int64
	Callback func(float64)
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.Current += int64(n)
	if pr.Callback != nil && pr.Total > 0 {
		progress := float64(pr.Current) / float64(pr.Total) * 100
		pr.Callback(progress)
	}
	return n, err
}

// downloadFile downloads a file from a URL to a local path with optional progress callback
func (lcm *LlamaCppReleaseManager) downloadFile(url, filepath string, progressCallback func(float64)) error {
	// Create HTTP client with proper configuration for downloads
	client := &http.Client{
		Timeout: 300 * time.Second, // 5 minutes for large downloads
		Transport: &http.Transport{
			DisableKeepAlives: true, // Prevent connection reuse issues
		},
	}

	// Create request with proper headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}

	// Add User-Agent header
	req.Header.Set("User-Agent", "Kawai-Agent/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Create a progress reader if callback is provided
	var reader io.Reader = resp.Body
	if progressCallback != nil && resp.ContentLength > 0 {
		reader = &progressReader{
			Reader:   resp.Body,
			Total:    resp.ContentLength,
			Callback: progressCallback,
		}
	}

	_, err = io.Copy(out, reader)
	return err
}

// saveVersionMetadata saves the installed version information
func (lcm *LlamaCppReleaseManager) saveVersionMetadata(version string) error {
	// Ensure metadata directory exists
	if err := os.MkdirAll(lcm.MetadataPath, 0755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	// Create metadata structure
	metadata := struct {
		Version     string    `json:"version"`
		InstalledAt time.Time `json:"installed_at"`
		BinaryPath  string    `json:"binary_path"`
		SourcePath  string    `json:"source_path"`
		BuildPath   string    `json:"build_path"`
	}{
		Version:     version,
		InstalledAt: time.Now(),
		BinaryPath:  lcm.BinaryPath,
		SourcePath:  lcm.SourcePath,
		BuildPath:   lcm.BuildPath,
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
func (lcm *LlamaCppReleaseManager) GetInstalledVersion() string {
	// First try to get version from metadata
	if version := lcm.loadVersionMetadata(); version != "" {
		return version
	}

	// Fallback: Check if any llama.cpp binary exists and try to get its version
	binaryPath := lcm.GetMainBinaryPath()
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return "" // No version installed
	}

	// Try to run the binary to get the version (fallback method)
	version, err := exec.Command(binaryPath, "--version").Output()
	if err != nil {
		// If binary exists but --version fails, try to determine from other means
		log.Printf("Binary exists but --version failed: %v", err)
		return ""
	}

	stringVersion := strings.TrimSpace(string(version))
	stringVersion = strings.TrimRight(stringVersion, "\n")

	// Extract version number from output
	lines := strings.Split(stringVersion, "\n")
	for _, line := range lines {
		if strings.Contains(line, "version") || strings.Contains(line, "commit") {
			// Try to extract version-like pattern
			parts := strings.Fields(line)
			for _, part := range parts {
				if strings.HasPrefix(part, "v") || (len(part) > 5 && strings.Contains(part, ".")) {
					return part
				}
			}
		}
	}

	return stringVersion
}

// loadVersionMetadata loads the installed version from metadata file
func (lcm *LlamaCppReleaseManager) loadVersionMetadata() string {
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
		SourcePath  string    `json:"source_path"`
		BuildPath   string    `json:"build_path"`
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

// GetMainBinaryPath returns the path to the main llama.cpp binary (llama-cli)
func (lcm *LlamaCppReleaseManager) GetMainBinaryPath() string {
	binaryName := "llama-cli"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	return filepath.Join(lcm.BinaryPath, binaryName)
}

// GetServerBinaryPath returns the path to the llama-server binary
// First checks system PATH (for package manager installations), then local binary path
func (lcm *LlamaCppReleaseManager) GetServerBinaryPath() string {
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
// Platform-specific implementation in manager_*.go files

// IsUpdateAvailable checks if an update is available
func (lcm *LlamaCppReleaseManager) IsUpdateAvailable() (bool, string, error) {
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

// IsLlamaCppInstalled and VerifyInstalledBinary are implemented in platform-specific files:
// - manager_darwin.go for macOS (checks Homebrew paths)
// - manager_default.go for other platforms (Linux, Windows)

// CleanupPartialDownloads removes any partial or corrupted downloads
func (lcm *LlamaCppReleaseManager) CleanupPartialDownloads() error {
	// Check main binary
	mainBinaryPath := lcm.GetMainBinaryPath()
	if _, err := os.Stat(mainBinaryPath); err == nil {
		// Binary exists, verify it
		if verifyErr := lcm.VerifyInstalledBinary(); verifyErr != nil {
			log.Printf("Found corrupted binary, removing: %v", verifyErr)
			// Remove all binaries in the directory
			if removeErr := os.RemoveAll(lcm.BinaryPath); removeErr != nil {
				log.Printf("Failed to remove binary directory: %v", removeErr)
			}
			// Clear metadata since binaries are corrupted
			lcm.clearVersionMetadata()
		}
	}

	// Clean up any partial source downloads
	if entries, err := os.ReadDir(lcm.SourcePath); err == nil {
		for _, entry := range entries {
			if strings.HasSuffix(entry.Name(), ".tar.gz") {
				archivePath := filepath.Join(lcm.SourcePath, entry.Name())
				if removeErr := os.Remove(archivePath); removeErr != nil {
					log.Printf("Failed to remove partial archive: %v", removeErr)
				}
			}
		}
	}

	return nil
}

// clearVersionMetadata clears the version metadata (used when binaries are corrupted or removed)
func (lcm *LlamaCppReleaseManager) clearVersionMetadata() {
	metadataPath := filepath.Join(lcm.MetadataPath, "installed-version.json")
	if err := os.Remove(metadataPath); err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to clear version metadata: %v", err)
	} else {
		log.Printf("Version metadata cleared")
	}
}

// GetAvailableBinaries returns a list of available llama.cpp binaries
func (lcm *LlamaCppReleaseManager) GetAvailableBinaries() ([]string, error) {
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
func (lcm *LlamaCppReleaseManager) HasBinary(binaryName string) bool {
	binaryPath := lcm.GetBinaryPath(binaryName)
	_, err := os.Stat(binaryPath)
	return err == nil
}

// RunBinary executes a llama.cpp binary with the given arguments
func (lcm *LlamaCppReleaseManager) RunBinary(binaryName string, args []string) (*exec.Cmd, error) {
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
func (lcm *LlamaCppReleaseManager) getCachedRelease() *Release {
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
func (lcm *LlamaCppReleaseManager) cacheRelease(release *Release) {
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
func (lcm *LlamaCppReleaseManager) QuantizeModel(inputPath, outputPath, quantType string) error {
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
func (lcm *LlamaCppReleaseManager) BenchmarkModel(modelPath string, options map[string]interface{}) (*BenchmarkResults, error) {
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
func (lcm *LlamaCppReleaseManager) parseBenchmarkOutput(output string) *BenchmarkResults {
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
func (lcm *LlamaCppReleaseManager) ProcessImageWithText(imagePath, prompt, modelPath string) (string, error) {
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
func (lcm *LlamaCppReleaseManager) extractLlavaResponse(output string) string {
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
func (lcm *LlamaCppReleaseManager) GetAvailableQuantizationTypes() []string {
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
func (lcm *LlamaCppReleaseManager) GetModelsDirectory() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".llama-cpp", "models")
}
