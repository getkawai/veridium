package llama

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/hybridgroup/yzma/pkg/download"
)

// TestNewLlamaCppInstaller tests the constructor
func TestNewLlamaCppInstaller(t *testing.T) {
	installer := NewLlamaCppInstaller()

	if installer == nil {
		t.Fatal("NewLlamaCppInstaller returned nil")
	}

	if installer.BinaryPath == "" {
		t.Error("BinaryPath should not be empty")
	}

	if installer.MetadataPath == "" {
		t.Error("MetadataPath should not be empty")
	}

	// Verify paths are under home directory
	homeDir, _ := os.UserHomeDir()
	expectedBase := filepath.Join(homeDir, ".llama-cpp")

	if !filepath.HasPrefix(installer.BinaryPath, expectedBase) {
		t.Errorf("BinaryPath should be under %s, got %s", expectedBase, installer.BinaryPath)
	}

	if !filepath.HasPrefix(installer.MetadataPath, expectedBase) {
		t.Errorf("MetadataPath should be under %s, got %s", expectedBase, installer.MetadataPath)
	}
}

// TestIsLlamaCppInstalled tests library detection
func TestIsLlamaCppInstalled(t *testing.T) {
	// Create temporary directory for testing
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   tmpDir,
		MetadataPath: filepath.Join(tmpDir, "metadata"),
	}

	// Test 1: Library not installed
	if installer.IsLlamaCppInstalled() {
		t.Error("IsLlamaCppInstalled should return false when library is not present")
	}

	// Test 2: Library installed
	libraryName := download.LibraryName(runtime.GOOS)
	if libraryName == "unknown" {
		t.Skip("Skipping test on unsupported platform")
	}

	libraryPath := filepath.Join(tmpDir, libraryName)
	file, err := os.Create(libraryPath)
	if err != nil {
		t.Fatalf("Failed to create test library file: %v", err)
	}
	file.Close()

	if !installer.IsLlamaCppInstalled() {
		t.Error("IsLlamaCppInstalled should return true when library is present")
	}
}

// TestVerifyInstalledBinary tests library verification
func TestVerifyInstalledBinary(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   tmpDir,
		MetadataPath: filepath.Join(tmpDir, "metadata"),
	}

	libraryName := download.LibraryName(runtime.GOOS)
	if libraryName == "unknown" {
		t.Skip("Skipping test on unsupported platform")
	}

	// Test 1: Library not found
	err := installer.VerifyInstalledBinary()
	if err == nil {
		t.Error("VerifyInstalledBinary should return error when library is not found")
	}

	// Test 2: Library exists
	libraryPath := filepath.Join(tmpDir, libraryName)
	file, err := os.Create(libraryPath)
	if err != nil {
		t.Fatalf("Failed to create test library file: %v", err)
	}
	file.Close()

	err = installer.VerifyInstalledBinary()
	if err != nil {
		t.Errorf("VerifyInstalledBinary should not return error when library exists: %v", err)
	}
}

// TestGetLatestRelease tests release fetching
func TestGetLatestRelease(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
	}

	// Ensure metadata directory exists
	os.MkdirAll(installer.MetadataPath, 0755)

	// Test 1: Fetch latest release (may hit network, but should work)
	release, err := installer.GetLatestRelease()
	if err != nil {
		t.Logf("GetLatestRelease failed (may be network issue): %v", err)
		t.Skip("Skipping test due to network/API issue")
	}

	if release == nil {
		t.Fatal("GetLatestRelease returned nil release")
	}

	if release.Version == "" {
		t.Error("Release version should not be empty")
	}

	if release.Name == "" {
		t.Error("Release name should not be empty")
	}

	// Test 2: Should use cache on second call
	release2, err := installer.GetLatestRelease()
	if err != nil {
		t.Fatalf("GetLatestRelease failed on second call: %v", err)
	}

	if release2.Version != release.Version {
		t.Errorf("Cached release version mismatch: got %s, expected %s", release2.Version, release.Version)
	}
}

// TestGetInstalledVersion tests version detection
func TestGetInstalledVersion(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
	}

	// Test 1: No version installed
	version := installer.GetInstalledVersion()
	if version != "" {
		t.Errorf("GetInstalledVersion should return empty string when no version installed, got: %s", version)
	}

	// Test 2: Version from metadata
	os.MkdirAll(installer.MetadataPath, 0755)
	testVersion := "b6924"
	err := installer.saveVersionMetadata(testVersion)
	if err != nil {
		t.Fatalf("Failed to save test metadata: %v", err)
	}

	version = installer.GetInstalledVersion()
	if version != testVersion {
		t.Errorf("GetInstalledVersion should return %s, got %s", testVersion, version)
	}
}

// TestIsUpdateAvailable tests update detection
func TestIsUpdateAvailable(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
	}

	os.MkdirAll(installer.MetadataPath, 0755)

	// Test 1: No version installed
	available, latest, err := installer.IsUpdateAvailable()
	if err != nil {
		t.Logf("IsUpdateAvailable failed (may be network issue): %v", err)
		t.Skip("Skipping test due to network/API issue")
	}

	if !available {
		t.Error("Update should be available when no version is installed")
	}

	if latest == "" {
		t.Error("Latest version should not be empty")
	}

	// Test 2: Version installed - check against real latest version
	testVersion := "b6924"
	installer.saveVersionMetadata(testVersion)

	// Get real latest version from API (no mock)
	available, latest, err = installer.IsUpdateAvailable()
	if err != nil {
		t.Logf("IsUpdateAvailable failed (may be network issue): %v", err)
		t.Skip("Skipping test due to network/API issue")
	}

	// Log results for debugging
	t.Logf("Update available: %v, Latest: %s, Current: %s", available, latest, testVersion)

	// Verify we got a real version from API
	if latest == "" {
		t.Error("Latest version should not be empty")
	}

	// If versions match, update should not be available
	// If versions differ, update should be available
	if testVersion == latest && available {
		t.Logf("Note: Current version %s matches latest %s, but update available is true", testVersion, latest)
	}
}

// TestDownloadRelease tests download functionality
func TestDownloadRelease(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
	}

	// Test 1: Download with empty version (should get latest)
	// Note: This will actually download, so we'll skip if network is unavailable
	err := installer.DownloadRelease("", nil)
	if err != nil {
		t.Logf("DownloadRelease failed (may be network issue): %v", err)
		t.Skip("Skipping download test due to network issue")
	}

	// Verify library file exists
	libraryName := download.LibraryName(runtime.GOOS)
	if libraryName != "unknown" {
		libraryPath := filepath.Join(installer.BinaryPath, libraryName)
		if _, err := os.Stat(libraryPath); os.IsNotExist(err) {
			t.Errorf("Library file should exist after download: %s", libraryPath)
		}
	}

	// Verify metadata was saved
	version := installer.GetInstalledVersion()
	if version == "" {
		t.Error("Version should be saved after download")
	}
}

// TestInstallLlamaCpp tests installation flow
func TestInstallLlamaCpp(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
	}

	// Test 1: Install when not installed
	// Note: This will actually download, so we'll skip if network is unavailable
	err := installer.InstallLlamaCpp()
	if err != nil {
		t.Logf("InstallLlamaCpp failed (may be network issue): %v", err)
		t.Skip("Skipping install test due to network issue")
	}

	// Verify installation
	if !installer.IsLlamaCppInstalled() {
		t.Error("llama.cpp should be installed after InstallLlamaCpp")
	}

	// Test 2: Install when already installed (should return early)
	err = installer.InstallLlamaCpp()
	if err != nil {
		t.Errorf("InstallLlamaCpp should not fail when already installed: %v", err)
	}
}

// TestCleanupPartialDownloads tests cleanup functionality
func TestCleanupPartialDownloads(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
	}

	os.MkdirAll(installer.BinaryPath, 0755)
	os.MkdirAll(installer.MetadataPath, 0755)

	// Test 1: No corrupted files (should succeed)
	err := installer.CleanupPartialDownloads()
	if err != nil {
		t.Errorf("CleanupPartialDownloads should not fail when no corrupted files: %v", err)
	}

	// Test 2: Corrupted library (missing library file but metadata exists)
	installer.saveVersionMetadata("b6924")
	// Don't create library file - simulate corruption

	err = installer.CleanupPartialDownloads()
	if err != nil {
		t.Errorf("CleanupPartialDownloads should handle corrupted library gracefully: %v", err)
	}

	// Verify metadata was cleared
	version := installer.GetInstalledVersion()
	if version != "" {
		t.Error("Version metadata should be cleared after cleanup")
	}
}

// TestDetectProcessor tests processor detection
func TestDetectProcessor(t *testing.T) {
	installer := NewLlamaCppInstaller()

	processor := installer.detectProcessor()

	validProcessors := []string{"cpu", "cuda", "vulkan", "metal"}
	isValid := false
	for _, p := range validProcessors {
		if processor == p {
			isValid = true
			break
		}
	}

	if !isValid {
		t.Errorf("detectProcessor returned invalid processor: %s", processor)
	}

	// macOS should return "metal"
	if runtime.GOOS == "darwin" && processor != "metal" {
		t.Logf("Warning: macOS detected processor as %s, expected metal (may be OK if CUDA/Vulkan detected)", processor)
	}
}

// TestLibraryName tests library name detection
func TestLibraryName(t *testing.T) {
	testCases := []struct {
		os       string
		expected string
	}{
		{"linux", "libllama.so"},
		{"darwin", "libllama.dylib"},
		{"windows", "llama.dll"},
		{"freebsd", "libllama.so"},
		{"unknown", "unknown"},
	}

	for _, tc := range testCases {
		result := download.LibraryName(tc.os)
		if result != tc.expected {
			t.Errorf("LibraryName(%s) = %s, expected %s", tc.os, result, tc.expected)
		}
	}
}

// TestCacheRelease tests caching functionality
func TestCacheRelease(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
	}

	release := &Release{
		Version: "b6924",
		Name:    "llama.cpp b6924",
		Body:    "Test release",
	}

	// Cache the release
	installer.cacheRelease(release)

	// Verify cache file exists
	cachePath := filepath.Join(installer.MetadataPath, "release-cache.json")
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		t.Error("Cache file should exist after caching release")
	}

	// Test getCachedRelease
	cached := installer.getCachedRelease()
	if cached == nil {
		t.Error("getCachedRelease should return cached release")
	}

	if cached.Version != release.Version {
		t.Errorf("Cached release version mismatch: got %s, expected %s", cached.Version, release.Version)
	}
}

// TestCacheExpiry tests cache expiry
func TestCacheExpiry(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
	}

	os.MkdirAll(installer.MetadataPath, 0755)

	release := &Release{
		Version: "b6924",
		Name:    "llama.cpp b6924",
		Body:    "Test release",
	}

	// Cache the release
	installer.cacheRelease(release)

	// Manually set cache file modification time to 2 hours ago
	cachePath := filepath.Join(installer.MetadataPath, "release-cache.json")
	oldTime := time.Now().Add(-2 * time.Hour)
	os.Chtimes(cachePath, oldTime, oldTime)

	// Should not return cached release (expired)
	cached := installer.getCachedRelease()
	if cached != nil {
		t.Error("getCachedRelease should return nil for expired cache")
	}
}

// TestMakeExecutable tests making binaries executable
func TestMakeExecutable(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("Skipping on Windows (no executable permissions)")
	}

	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   tmpDir,
		MetadataPath: filepath.Join(tmpDir, "metadata"),
	}

	os.MkdirAll(tmpDir, 0755)

	// Create a test file
	testFile := filepath.Join(tmpDir, "test-binary")
	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	file.Close()

	// Make executable
	err = installer.makeExecutable()
	if err != nil {
		t.Errorf("makeExecutable failed: %v", err)
	}

	// Verify file is executable
	info, err := os.Stat(testFile)
	if err != nil {
		t.Fatalf("Failed to stat test file: %v", err)
	}

	if info.Mode()&0111 == 0 {
		t.Error("File should be executable after makeExecutable")
	}
}

// TestSaveVersionMetadata tests metadata saving
func TestSaveVersionMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
	}

	testVersion := "b6924"
	err := installer.saveVersionMetadata(testVersion)
	if err != nil {
		t.Fatalf("saveVersionMetadata failed: %v", err)
	}

	// Verify metadata file exists
	metadataPath := filepath.Join(installer.MetadataPath, "installed-version.json")
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		t.Error("Metadata file should exist after saveVersionMetadata")
	}

	// Verify version can be loaded
	version := installer.GetInstalledVersion()
	if version != testVersion {
		t.Errorf("Loaded version mismatch: got %s, expected %s", version, testVersion)
	}
}

// TestClearVersionMetadata tests metadata clearing
func TestClearVersionMetadata(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
	}

	os.MkdirAll(installer.MetadataPath, 0755)

	// Save metadata first
	installer.saveVersionMetadata("b6924")

	// Clear metadata
	installer.clearVersionMetadata()

	// Verify metadata is cleared
	version := installer.GetInstalledVersion()
	if version != "" {
		t.Error("Version should be empty after clearVersionMetadata")
	}
}
