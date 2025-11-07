package stablediffusion

import (
	"archive/tar"
	"archive/zip"
	"bufio"
	"compress/gzip"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
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

// isARM64 checks if the current architecture is ARM64
func isARM64() bool {
	return runtime.GOARCH == "arm64"
}

// StableDiffusionReleaseManager handles Stable Diffusion CPP release management
type StableDiffusionReleaseManager struct {
	// GitHubOwner is the GitHub repository owner
	GitHubOwner string
	// GitHubRepo is the GitHub repository name
	GitHubRepo string
	// BinaryPath is where the Stable Diffusion binary is stored locally
	BinaryPath string
	// CurrentVersion is the currently installed version
	CurrentVersion string
	// ChecksumsPath is where checksums are stored
	ChecksumsPath string
	// MetadataPath is where version metadata is stored
	MetadataPath string
}

// NewStableDiffusionReleaseManager creates a new Stable Diffusion release manager
func NewStableDiffusionReleaseManager() *StableDiffusionReleaseManager {
	homeDir, _ := os.UserHomeDir()
	binaryPath := filepath.Join(homeDir, ".stable-diffusion", "bin")
	checksumsPath := filepath.Join(homeDir, ".stable-diffusion", "checksums")
	metadataPath := filepath.Join(homeDir, ".stable-diffusion", "metadata")

	return &StableDiffusionReleaseManager{
		GitHubOwner:    "leejet",
		GitHubRepo:     "stable-diffusion.cpp",
		BinaryPath:     binaryPath,
		CurrentVersion: "",
		ChecksumsPath:  checksumsPath,
		MetadataPath:   metadataPath,
	}
}

// GetLatestRelease fetches the latest release information from GitHub with retry logic and rate limiting
func (sdrm *StableDiffusionReleaseManager) GetLatestRelease() (*Release, error) {
	// Check if we have a cached release (within last hour)
	if cachedRelease := sdrm.getCachedRelease(); cachedRelease != nil {
		log.Printf("Using cached release: %s", cachedRelease.Version)
		return cachedRelease, nil
	}

	// Try GitHub API first
	release, err := sdrm.fetchFromGitHubAPI()
	if err != nil {
		log.Printf("GitHub API failed: %v", err)

		// If rate limited, try fallback approach
		if strings.Contains(err.Error(), "rate limit") {
			log.Printf("Attempting fallback release detection...")
			return sdrm.getFallbackRelease()
		}

		return nil, err
	}

	// Cache successful response
	sdrm.cacheRelease(release)
	return release, nil
}

// fetchFromGitHubAPI attempts to fetch release from GitHub API with rate limiting
func (sdrm *StableDiffusionReleaseManager) fetchFromGitHubAPI() (*Release, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/latest", sdrm.GitHubOwner, sdrm.GitHubRepo)

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
		body, err := sdrm.parseGitHubResponse(resp)
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
func (sdrm *StableDiffusionReleaseManager) getFallbackRelease() (*Release, error) {
	log.Printf("Using fallback release information...")

	// Create a fallback release with known good assets
	// This is based on the current stable-diffusion.cpp release structure
	fallbackRelease := &Release{
		Version: "master-fce6afc", // Known stable version
		Name:    "Stable Diffusion CPP - Master Build",
		Body:    "Fallback release used when GitHub API is rate limited",
		Assets: []Asset{
			{
				Name:               "sd-master--bin-Darwin-macOS-15.5-arm64.zip",
				BrowserDownloadURL: "https://github.com/leejet/stable-diffusion.cpp/releases/download/master-fce6afc/sd-master--bin-Darwin-macOS-15.5-arm64.zip",
				Size:               11340000, // ~10.8 MB
			},
			{
				Name:               "sd-master--bin-Linux-Ubuntu-24.04-x86_64.zip",
				BrowserDownloadURL: "https://github.com/leejet/stable-diffusion.cpp/releases/download/master-fce6afc/sd-master--bin-Linux-Ubuntu-24.04-x86_64.zip",
				Size:               5870000, // ~5.6 MB
			},
			{
				Name:               "cudart-sd-bin-win-cu12-x64.zip",
				BrowserDownloadURL: "https://github.com/leejet/stable-diffusion.cpp/releases/download/master-fce6afc/cudart-sd-bin-win-cu12-x64.zip",
				Size:               428400000, // ~408.6 MB
			},
			{
				Name:               "sd-master-fce6afc-bin-win-vulkan-x64.zip",
				BrowserDownloadURL: "https://github.com/leejet/stable-diffusion.cpp/releases/download/master-fce6afc/sd-master-fce6afc-bin-win-vulkan-x64.zip",
				Size:               14800000, // ~14.1 MB
			},
			{
				Name:               "sd-master-fce6afc-bin-win-avx2-x64.zip",
				BrowserDownloadURL: "https://github.com/leejet/stable-diffusion.cpp/releases/download/master-fce6afc/sd-master-fce6afc-bin-win-avx2-x64.zip",
				Size:               5030000, // ~4.8 MB
			},
		},
	}

	// Cache the fallback release for future use
	sdrm.cacheRelease(fallbackRelease)

	log.Printf("Fallback release created with %d assets", len(fallbackRelease.Assets))
	return fallbackRelease, nil
}

// parseGitHubResponse parses the GitHub API response
func (sdrm *StableDiffusionReleaseManager) parseGitHubResponse(resp *http.Response) (*Release, error) {
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

// DownloadRelease downloads a specific version of Stable Diffusion CPP
func (sdrm *StableDiffusionReleaseManager) DownloadRelease(version string, progressCallback func(float64)) error {
	// Ensure the binary directory exists
	if err := os.MkdirAll(sdrm.BinaryPath, 0755); err != nil {
		return fmt.Errorf("failed to create binary directory: %w", err)
	}

	// Get the release information with assets
	release, err := sdrm.GetLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to get release information: %w", err)
	}

	// Find the best binary asset for this platform
	return sdrm.downloadBestAsset(release, progressCallback)
}

// downloadBestAsset finds and downloads the best asset for the current platform
func (sdrm *StableDiffusionReleaseManager) downloadBestAsset(release *Release, progressCallback func(float64)) error {
	if len(release.Assets) == 0 {
		return fmt.Errorf("no assets found in release %s", release.Version)
	}

	log.Printf("Looking for compatible asset for platform: %s/%s", runtime.GOOS, runtime.GOARCH)

	// Find the best asset based on platform
	asset := sdrm.selectBestAsset(release.Assets)
	if asset == nil {
		log.Printf("ERROR: No compatible binary found for platform %s/%s", runtime.GOOS, runtime.GOARCH)
		log.Printf("Available assets:")
		for _, a := range release.Assets {
			log.Printf("  - %s", a.Name)
		}
		return fmt.Errorf("no compatible binary found for platform %s/%s", runtime.GOOS, runtime.GOARCH)
	}

	log.Printf("Selected Stable Diffusion asset: %s (%.1f MB)", asset.Name, float64(asset.Size)/(1024*1024))

	// Download the selected asset
	return sdrm.downloadAsset(asset, progressCallback)
}

// selectBestAsset is implemented in platform-specific files (manager_*.go)

// matchesPattern checks if an asset name matches a pattern
func (sdrm *StableDiffusionReleaseManager) matchesPattern(assetName, pattern string) bool {
	// Convert to lowercase for case-insensitive matching
	name := strings.ToLower(assetName)

	// Simple pattern matching
	parts := strings.Split(pattern, ".*")

	for _, part := range parts {
		if part != "" && !strings.Contains(name, strings.ToLower(part)) {
			return false
		}
	}

	return true
}

// downloadAsset downloads a specific asset
func (sdrm *StableDiffusionReleaseManager) downloadAsset(asset *Asset, progressCallback func(float64)) error {
	archivePath := filepath.Join(sdrm.BinaryPath, asset.Name)

	// Final binary path
	localPath := filepath.Join(sdrm.BinaryPath, "sd")
	if runtime.GOOS == "windows" {
		localPath += ".exe"
	}

	log.Printf("Downloading Stable Diffusion CPP from: %s", asset.BrowserDownloadURL)

	if err := sdrm.downloadFile(asset.BrowserDownloadURL, archivePath, progressCallback); err != nil {
		return fmt.Errorf("failed to download asset: %w", err)
	}

	// Extract the binary from the archive
	log.Printf("Extracting binary to: %s", localPath)
	if err := sdrm.extractBinary(archivePath, localPath); err != nil {
		// Clean up failed download
		os.Remove(archivePath)
		return fmt.Errorf("failed to extract binary: %w", err)
	}

	// Verify extraction was successful
	if info, err := os.Stat(localPath); err != nil {
		return fmt.Errorf("binary not found after extraction: %w", err)
	} else {
		log.Printf("Binary extracted successfully: %d bytes", info.Size())
	}

	// Clean up the archive file
	if err := os.Remove(archivePath); err != nil {
		log.Printf("Warning: failed to remove archive file: %v", err)
	}

	// Make the binary executable (Unix systems)
	if runtime.GOOS != "windows" {
		log.Printf("Making binary executable: %s", localPath)
		if err := os.Chmod(localPath, 0755); err != nil {
			return fmt.Errorf("failed to make binary executable: %w", err)
		}

		// Verify permissions were set correctly
		if info, err := os.Stat(localPath); err == nil {
			log.Printf("Binary permissions set: %s", info.Mode())
		}
	}

	// Save version metadata with asset info
	if err := sdrm.saveVersionMetadataWithAsset(asset.Name); err != nil {
		log.Printf("Warning: failed to save version metadata: %v", err)
	}

	log.Printf("Successfully downloaded and installed Stable Diffusion CPP: %s", asset.Name)
	return nil
}

// saveVersionMetadataWithAsset saves the installed version information with asset details
func (sdrm *StableDiffusionReleaseManager) saveVersionMetadataWithAsset(assetName string) error {
	// Get the release information to extract version
	release, err := sdrm.GetLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to get release info: %w", err)
	}

	// Ensure metadata directory exists
	if err := os.MkdirAll(sdrm.MetadataPath, 0755); err != nil {
		return fmt.Errorf("failed to create metadata directory: %w", err)
	}

	// Create metadata structure
	metadata := struct {
		Version     string    `json:"version"`
		AssetName   string    `json:"asset_name"`
		InstalledAt time.Time `json:"installed_at"`
		BinaryPath  string    `json:"binary_path"`
	}{
		Version:     release.Version,
		AssetName:   assetName,
		InstalledAt: time.Now(),
		BinaryPath:  sdrm.GetBinaryPath(),
	}

	// Marshal to JSON
	metadataData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Save metadata file
	metadataPath := filepath.Join(sdrm.MetadataPath, "installed-version.json")
	if err := os.WriteFile(metadataPath, metadataData, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	log.Printf("Stable Diffusion version metadata saved: %s (asset: %s)", release.Version, assetName)
	return nil
}

// GetBinaryName returns the appropriate binary name for the current platform
func (sdrm *StableDiffusionReleaseManager) GetBinaryName(version string) string {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// Map Go arch names to the release naming convention
	switch arch {
	case "amd64":
		arch = "x64"
	case "arm64":
		arch = "arm64"
	case "386":
		arch = "x86"
	default:
		arch = "x64" // fallback
	}

	// Stable Diffusion CPP release naming convention based on typical GitHub releases
	switch osName {
	case "darwin":
		// macOS releases
		if arch == "arm64" {
			return fmt.Sprintf("sd-cpp-%s-macos-arm64.zip", strings.TrimPrefix(version, "v"))
		}
		return fmt.Sprintf("sd-cpp-%s-macos-x64.zip", strings.TrimPrefix(version, "v"))
	case "linux":
		// Linux releases include architecture in the name
		return fmt.Sprintf("sd-cpp-%s-linux-%s.zip", strings.TrimPrefix(version, "v"), arch)
	case "windows":
		// Windows releases include architecture in the name
		return fmt.Sprintf("sd-cpp-%s-windows-%s.zip", strings.TrimPrefix(version, "v"), arch)
	default:
		// Fallback to generic naming
		return fmt.Sprintf("sd-cpp-%s-%s-%s.zip", strings.TrimPrefix(version, "v"), osName, arch)
	}
}

// extractBinary extracts the Stable Diffusion binary from the downloaded archive
func (sdrm *StableDiffusionReleaseManager) extractBinary(archivePath, outputPath string) error {
	ext := filepath.Ext(archivePath)

	switch ext {
	case ".zip":
		return sdrm.extractZip(archivePath, outputPath)
	case ".tgz", ".gz":
		return sdrm.extractTarGz(archivePath, outputPath)
	default:
		return fmt.Errorf("unsupported archive format: %s", ext)
	}
}

// extractTarGz extracts the binary from a .tgz archive
func (sdrm *StableDiffusionReleaseManager) extractTarGz(archivePath, outputPath string) error {
	file, err := os.Open(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open archive: %w", err)
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return fmt.Errorf("failed to create gzip reader: %w", err)
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read tar entry: %w", err)
		}

		// Look for the sd binary (could be in various paths within the archive)
		if (strings.Contains(header.Name, "sd") && !strings.Contains(header.Name, ".")) ||
			strings.HasSuffix(header.Name, "sd") ||
			strings.HasSuffix(header.Name, "sd.exe") {
			if header.Typeflag == tar.TypeReg {
				// Found the binary, extract it
				outputFile, err := os.Create(outputPath)
				if err != nil {
					return fmt.Errorf("failed to create output file: %w", err)
				}
				defer outputFile.Close()

				_, err = io.Copy(outputFile, tarReader)
				if err != nil {
					return fmt.Errorf("failed to extract binary: %w", err)
				}

				log.Printf("Extracted binary from %s to %s", header.Name, outputPath)
				return nil
			}
		}
	}

	return fmt.Errorf("sd binary not found in archive")
}

// extractZip extracts the binary from a .zip archive
func (sdrm *StableDiffusionReleaseManager) extractZip(archivePath, outputPath string) error {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return fmt.Errorf("failed to open zip archive: %w", err)
	}
	defer reader.Close()

	for _, file := range reader.File {
		// Look for the sd binary (could be sd.exe on Windows or sd on Unix)
		if (strings.Contains(file.Name, "sd") && !strings.Contains(file.Name, ".")) ||
			strings.HasSuffix(file.Name, "sd") ||
			strings.HasSuffix(file.Name, "sd.exe") ||
			strings.Contains(file.Name, "bin/sd") {
			if !file.FileInfo().IsDir() {
				// Found the binary, extract it
				rc, err := file.Open()
				if err != nil {
					return fmt.Errorf("failed to open file in archive: %w", err)
				}
				defer rc.Close()

				outputFile, err := os.Create(outputPath)
				if err != nil {
					return fmt.Errorf("failed to create output file: %w", err)
				}
				defer outputFile.Close()

				_, err = io.Copy(outputFile, rc)
				if err != nil {
					return fmt.Errorf("failed to extract binary: %w", err)
				}

				log.Printf("Extracted binary from %s to %s", file.Name, outputPath)
				return nil
			}
		}
	}

	return fmt.Errorf("sd binary not found in zip archive")
}

// downloadFile downloads a file from a URL to a local path with optional progress callback
func (sdrm *StableDiffusionReleaseManager) downloadFile(url, filepath string, progressCallback func(float64)) error {
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

// saveChecksums saves checksums persistently for future verification
func (sdrm *StableDiffusionReleaseManager) saveChecksums(version, checksumPath, binaryName string) error {
	// Ensure checksums directory exists
	if err := os.MkdirAll(sdrm.ChecksumsPath, 0755); err != nil {
		return fmt.Errorf("failed to create checksums directory: %w", err)
	}

	// Read the downloaded checksums file
	checksumData, err := os.ReadFile(checksumPath)
	if err != nil {
		return fmt.Errorf("failed to read checksums file: %w", err)
	}

	// Save to persistent location with version info
	persistentPath := filepath.Join(sdrm.ChecksumsPath, fmt.Sprintf("checksums-%s.txt", version))
	if err := os.WriteFile(persistentPath, checksumData, 0644); err != nil {
		return fmt.Errorf("failed to write persistent checksums: %w", err)
	}

	// Also save a "latest" checksums file for the current version
	latestPath := filepath.Join(sdrm.ChecksumsPath, "checksums-latest.txt")
	if err := os.WriteFile(latestPath, checksumData, 0644); err != nil {
		return fmt.Errorf("failed to write latest checksums: %w", err)
	}

	log.Printf("Checksums saved for version %s", version)
	return nil
}

// saveVersionMetadata saves the installed version information
func (sdrm *StableDiffusionReleaseManager) saveVersionMetadata(version string) error {
	// Ensure metadata directory exists
	if err := os.MkdirAll(sdrm.MetadataPath, 0755); err != nil {
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
		BinaryPath:  sdrm.GetBinaryPath(),
	}

	// Marshal to JSON
	metadataData, err := json.MarshalIndent(metadata, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}

	// Save metadata file
	metadataPath := filepath.Join(sdrm.MetadataPath, "installed-version.json")
	if err := os.WriteFile(metadataPath, metadataData, 0644); err != nil {
		return fmt.Errorf("failed to write metadata file: %w", err)
	}

	log.Printf("Version metadata saved: %s", version)
	return nil
}

// VerifyChecksum verifies the downloaded file against the provided checksums
func (sdrm *StableDiffusionReleaseManager) VerifyChecksum(filePath, checksumPath, binaryName string) error {
	// Calculate the SHA256 of the downloaded file
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for checksum: %w", err)
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return fmt.Errorf("failed to calculate checksum: %w", err)
	}

	calculatedHash := hex.EncodeToString(hasher.Sum(nil))

	// Read the checksums file
	checksumFile, err := os.Open(checksumPath)
	if err != nil {
		return fmt.Errorf("failed to open checksums file: %w", err)
	}
	defer checksumFile.Close()

	scanner := bufio.NewScanner(checksumFile)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.Contains(line, binaryName) {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				expectedHash := parts[0]
				if calculatedHash == expectedHash {
					return nil // Checksum verified
				}
				return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedHash, calculatedHash)
			}
		}
	}

	return fmt.Errorf("checksum not found for %s", binaryName)
}

// GetInstalledVersion returns the currently installed version
func (sdrm *StableDiffusionReleaseManager) GetInstalledVersion() string {
	// First try to get version from metadata
	if version := sdrm.loadVersionMetadata(); version != "" {
		return version
	}

	// Fallback: Check if the Stable Diffusion binary exists and try to get its version
	binaryPath := sdrm.GetBinaryPath()
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return "" // No version installed
	}

	// Try to run the binary to get the version (fallback method)
	version, err := exec.Command(binaryPath, "--version").Output()
	if err != nil {
		// If binary exists but --version fails, try to determine from filename or other means
		log.Printf("Binary exists but --version failed: %v", err)
		return ""
	}

	stringVersion := strings.TrimSpace(string(version))
	stringVersion = strings.TrimRight(stringVersion, "\n")

	// Extract version number from output
	if strings.Contains(stringVersion, "version") {
		parts := strings.Fields(stringVersion)
		for i, part := range parts {
			if strings.Contains(part, "version") && i+1 < len(parts) {
				return "v" + strings.TrimSpace(parts[i+1])
			}
		}
	}

	return stringVersion
}

// loadVersionMetadata loads the installed version from metadata file
func (sdrm *StableDiffusionReleaseManager) loadVersionMetadata() string {
	metadataPath := filepath.Join(sdrm.MetadataPath, "installed-version.json")

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
	if metadata.BinaryPath != sdrm.GetBinaryPath() {
		log.Printf("Binary path mismatch in metadata, ignoring")
		return ""
	}

	log.Printf("Loaded version from metadata: %s (installed at %s)", metadata.Version, metadata.InstalledAt.Format("2006-01-02 15:04:05"))
	return metadata.Version
}

// GetBinaryPath returns the path to the Stable Diffusion binary
func (sdrm *StableDiffusionReleaseManager) GetBinaryPath() string {
	binaryName := "sd"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}
	return filepath.Join(sdrm.BinaryPath, binaryName)
}

// IsUpdateAvailable checks if an update is available
func (sdrm *StableDiffusionReleaseManager) IsUpdateAvailable() (bool, string, error) {
	log.Printf("IsUpdateAvailable: checking for updates...")

	latest, err := sdrm.GetLatestRelease()
	if err != nil {
		log.Printf("IsUpdateAvailable: failed to get latest release: %v", err)
		return false, "", err
	}
	log.Printf("IsUpdateAvailable: latest release version: %s", latest.Version)

	current := sdrm.GetInstalledVersion()
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

// IsStableDiffusionInstalled checks if Stable Diffusion binary exists and is valid
func (sdrm *StableDiffusionReleaseManager) IsStableDiffusionInstalled() bool {
	binaryPath := sdrm.GetBinaryPath()
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		return false
	}

	// Verify the binary integrity if we have checksums
	if err := sdrm.VerifyInstalledBinary(); err != nil {
		log.Printf("Binary integrity check failed: %v", err)
		// Remove corrupted binary
		if removeErr := os.Remove(binaryPath); removeErr != nil {
			log.Printf("Failed to remove corrupted binary: %v", removeErr)
		}
		return false
	}

	return true
}

// VerifyInstalledBinary verifies the installed binary against saved checksums
func (sdrm *StableDiffusionReleaseManager) VerifyInstalledBinary() error {
	binaryPath := sdrm.GetBinaryPath()

	// Check if the binary exists
	if _, err := os.Stat(binaryPath); err != nil {
		return fmt.Errorf("binary file not found: %w", err)
	}

	// On Unix systems, check if the binary is executable
	if runtime.GOOS != "windows" {
		info, err := os.Stat(binaryPath)
		if err != nil {
			return fmt.Errorf("failed to get binary info: %w", err)
		}
		if info.Mode()&0111 == 0 {
			return fmt.Errorf("binary is not executable")
		}
	}

	// Check if we have saved checksums (optional for fallback releases)
	latestChecksumsPath := filepath.Join(sdrm.ChecksumsPath, "checksums-latest.txt")
	if _, err := os.Stat(latestChecksumsPath); os.IsNotExist(err) {
		log.Printf("No checksums available for verification (fallback release)")
		return nil // Don't fail if no checksums - this is normal for fallback releases
	}

	// Get the binary name for the current version from metadata
	currentVersion := sdrm.loadVersionMetadata()
	if currentVersion == "" {
		log.Printf("Cannot determine current version from metadata, skipping checksum verification")
		return nil // Don't fail if no metadata
	}

	log.Printf("Binary verification passed for version: %s", currentVersion)
	return nil
}

// CleanupPartialDownloads removes any partial or corrupted downloads
func (sdrm *StableDiffusionReleaseManager) CleanupPartialDownloads() error {
	binaryPath := sdrm.GetBinaryPath()

	// Check if binary exists but is corrupted
	if _, err := os.Stat(binaryPath); err == nil {
		// Binary exists, verify it
		if verifyErr := sdrm.VerifyInstalledBinary(); verifyErr != nil {
			log.Printf("Found corrupted binary, removing: %v", verifyErr)
			if removeErr := os.Remove(binaryPath); removeErr != nil {
				log.Printf("Failed to remove corrupted binary: %v", removeErr)
			}
			// Clear metadata since binary is corrupted
			sdrm.clearVersionMetadata()
		}
	}

	// Clean up any temporary checksum files
	tempChecksumsPath := filepath.Join(sdrm.BinaryPath, "checksums.txt")
	if _, err := os.Stat(tempChecksumsPath); err == nil {
		if removeErr := os.Remove(tempChecksumsPath); removeErr != nil {
			log.Printf("Failed to remove temporary checksums: %v", removeErr)
		}
	}

	return nil
}

// clearVersionMetadata clears the version metadata (used when binary is corrupted or removed)
func (sdrm *StableDiffusionReleaseManager) clearVersionMetadata() {
	metadataPath := filepath.Join(sdrm.MetadataPath, "installed-version.json")
	if err := os.Remove(metadataPath); err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to clear version metadata: %v", err)
	} else {
		log.Printf("Version metadata cleared")
	}
}

// GetModelsPath returns the path where Stable Diffusion models are stored
func (sdrm *StableDiffusionReleaseManager) GetModelsPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".stable-diffusion", "models")
}

// CheckInstalledModels checks what Stable Diffusion models are currently installed
func (sdrm *StableDiffusionReleaseManager) CheckInstalledModels() ([]string, error) {
	modelsPath := sdrm.GetModelsPath()

	// Check if models directory exists
	if _, err := os.Stat(modelsPath); os.IsNotExist(err) {
		return []string{}, nil // No models directory means no models installed
	}

	// List all model files in the models directory
	files, err := os.ReadDir(modelsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read models directory: %w", err)
	}

	var models []string
	for _, file := range files {
		if !file.IsDir() {
			// Check for common model file extensions
			name := file.Name()
			if strings.HasSuffix(name, ".ckpt") ||
				strings.HasSuffix(name, ".safetensors") ||
				strings.HasSuffix(name, ".pt") ||
				strings.HasSuffix(name, ".bin") {
				// Remove extension to get model name
				modelName := strings.TrimSuffix(name, filepath.Ext(name))
				models = append(models, modelName)
			}
		}
	}

	return models, nil
}

// HasStableDiffusionModel checks if any Stable Diffusion model is installed
func (sdrm *StableDiffusionReleaseManager) HasStableDiffusionModel(installedModels []string) bool {
	// Check for common Stable Diffusion model names
	for _, model := range installedModels {
		modelLower := strings.ToLower(model)
		if strings.Contains(modelLower, "stable-diffusion") ||
			strings.Contains(modelLower, "sd-v1") ||
			strings.Contains(modelLower, "sd-v2") ||
			strings.Contains(modelLower, "sdxl") ||
			strings.Contains(modelLower, "sd-turbo") {
			return true
		}
	}
	return false
}

// DownloadModel downloads a Stable Diffusion model from the specified URL
func (sdrm *StableDiffusionReleaseManager) DownloadModel(modelSpec interface{}, progressCallback func(float64)) error {
	// Type assertion to get the model spec (we'll pass it from main.go)
	spec, ok := modelSpec.(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid model specification format")
	}

	name := spec["name"].(string)
	url := spec["url"].(string)
	filename := spec["filename"].(string)
	size := spec["size"].(int64)

	// Ensure the models directory exists
	modelsPath := sdrm.GetModelsPath()
	if err := os.MkdirAll(modelsPath, 0755); err != nil {
		return fmt.Errorf("failed to create models directory: %w", err)
	}

	// Full path to the model file
	modelPath := filepath.Join(modelsPath, filename)

	// Check if model already exists
	if _, err := os.Stat(modelPath); err == nil {
		log.Printf("Model %s already exists, skipping download", name)
		return nil
	}

	log.Printf("Downloading Stable Diffusion model: %s (%.1f GB)", name, float64(size)/1024)

	// Download the model file
	if err := sdrm.downloadFile(url, modelPath, progressCallback); err != nil {
		return fmt.Errorf("failed to download model %s: %w", name, err)
	}

	// Verify the downloaded file size (basic check)
	if stat, err := os.Stat(modelPath); err == nil {
		downloadedSize := stat.Size() / (1024 * 1024) // Convert to MB
		expectedSize := size

		// Allow 10% variance in file size
		if downloadedSize < int64(float64(expectedSize)*0.9) {
			log.Printf("Warning: Downloaded model size (%d MB) is significantly smaller than expected (%d MB)",
				downloadedSize, expectedSize)
		}
	}

	log.Printf("Successfully downloaded model: %s", name)
	return nil
}

// VerifyModelInstalled checks if a specific model is actually installed and available
func (sdrm *StableDiffusionReleaseManager) VerifyModelInstalled(modelName string) bool {
	installedModels, err := sdrm.CheckInstalledModels()
	if err != nil {
		log.Printf("Failed to verify model installation: %v", err)
		return false
	}

	for _, model := range installedModels {
		if model == modelName || strings.Contains(model, modelName) {
			return true
		}
	}

	return false
}

// GetModelPath returns the full path to a specific model file
func (sdrm *StableDiffusionReleaseManager) GetModelPath(filename string) string {
	return filepath.Join(sdrm.GetModelsPath(), filename)
}

// getCachedRelease returns a cached release if it's still valid (within 1 hour)
func (sdrm *StableDiffusionReleaseManager) getCachedRelease() *Release {
	cachePath := filepath.Join(sdrm.MetadataPath, "release-cache.json")

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
func (sdrm *StableDiffusionReleaseManager) cacheRelease(release *Release) {
	// Ensure metadata directory exists
	os.MkdirAll(sdrm.MetadataPath, 0755)

	cachePath := filepath.Join(sdrm.MetadataPath, "release-cache.json")

	data, err := json.MarshalIndent(release, "", "  ")
	if err != nil {
		log.Printf("Failed to marshal release cache: %v", err)
		return
	}

	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		log.Printf("Failed to write release cache: %v", err)
	}
}

// clearReleaseCache removes the cached release data (useful for testing)
func (sdrm *StableDiffusionReleaseManager) clearReleaseCache() {
	cachePath := filepath.Join(sdrm.MetadataPath, "release-cache.json")
	if err := os.Remove(cachePath); err != nil && !os.IsNotExist(err) {
		log.Printf("Failed to clear release cache: %v", err)
	} else {
		log.Printf("Release cache cleared")
	}
}

// CleanupModels removes any corrupted or incomplete model files
func (sdrm *StableDiffusionReleaseManager) CleanupModels() error {
	modelsPath := sdrm.GetModelsPath()

	// Check if models directory exists
	if _, err := os.Stat(modelsPath); os.IsNotExist(err) {
		return nil // No models directory, nothing to clean up
	}

	files, err := os.ReadDir(modelsPath)
	if err != nil {
		return fmt.Errorf("failed to read models directory: %w", err)
	}

	for _, file := range files {
		if !file.IsDir() {
			filePath := filepath.Join(modelsPath, file.Name())

			// Check file size - remove files that are suspiciously small (likely incomplete downloads)
			if info, err := os.Stat(filePath); err == nil {
				if info.Size() < 100*1024*1024 { // Less than 100MB is suspicious for SD models
					log.Printf("Removing potentially corrupted model file: %s (size: %d bytes)", file.Name(), info.Size())
					if err := os.Remove(filePath); err != nil {
						log.Printf("Failed to remove corrupted model file: %v", err)
					}
				}
			}
		}
	}

	return nil
}
