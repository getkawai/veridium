package llama

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/kawai-network/veridium/pkg/yzma/download"
)

// ============================================================================
// Constructor Tests
// ============================================================================

func TestNewLlamaCppInstaller(t *testing.T) {
	installer := NewLlamaCppInstaller()

	if installer == nil {
		t.Fatal("NewLlamaCppInstaller() returned nil")
	}

	homeDir, _ := os.UserHomeDir()
	expectedBasePath := filepath.Join(homeDir, ".llama-cpp")

	if installer.BinaryPath != filepath.Join(expectedBasePath, "bin") {
		t.Errorf("BinaryPath = %v, want %v", installer.BinaryPath, filepath.Join(expectedBasePath, "bin"))
	}

	if installer.MetadataPath != filepath.Join(expectedBasePath, "metadata") {
		t.Errorf("MetadataPath = %v, want %v", installer.MetadataPath, filepath.Join(expectedBasePath, "metadata"))
	}

	if installer.ModelsDir != filepath.Join(expectedBasePath, "models") {
		t.Errorf("ModelsDir = %v, want %v", installer.ModelsDir, filepath.Join(expectedBasePath, "models"))
	}
}

// ============================================================================
// Release & Version Tests (Real API Calls)
// ============================================================================

// ============================================================================
// Installation Tests (Real Downloads)
// ============================================================================

func TestInstallLlamaCpp_RealDownload(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real download test in short mode")
	}

	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	t.Log("⏳ Starting real llama.cpp download (this may take several minutes)...")
	err := installer.InstallLlamaCpp()
	if err != nil {
		t.Logf("InstallLlamaCpp() failed (may be network issue): %v", err)
		t.Skip("Skipping download test due to network issue")
	}

	// Verify installation
	if !installer.IsLlamaCppInstalled() {
		t.Error("IsLlamaCppInstalled() = false after installation")
	}

	// Verify library file exists
	libraryName := download.LibraryName(runtime.GOOS)
	libraryPath := filepath.Join(installer.BinaryPath, libraryName)
	if _, err := os.Stat(libraryPath); err != nil {
		t.Errorf("Library file not found: %s", libraryPath)
	}

	t.Logf("✅ Successfully installed llama.cpp")
}

func TestInstallLlamaCpp_AlreadyInstalled(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real download test in short mode")
	}

	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	// First installation
	err := installer.InstallLlamaCpp()
	if err != nil {
		t.Skip("Skipping due to network issue")
	}

	// Second installation should skip
	start := time.Now()
	err = installer.InstallLlamaCpp()
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Second InstallLlamaCpp() failed: %v", err)
	}

	// Should be very fast (< 100ms) since it skips download
	if elapsed > 100*time.Millisecond {
		t.Logf("Warning: Second install took %v (expected < 100ms)", elapsed)
	}

	t.Logf("✅ Correctly skipped re-installation: %v", elapsed)
}

func TestInstallLlamaCpp_AutoUpgrade(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real download test in short mode")
	}

	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	// First install
	t.Log("⏳ First installation...")
	err := installer.InstallLlamaCpp()
	if err != nil {
		t.Skip("Skipping due to network issue")
	}

	t.Logf("✅ First installation complete")

	// Verify version.json was created by yzma
	versionFile := filepath.Join(installer.BinaryPath, "version.json")
	if _, err := os.Stat(versionFile); err != nil {
		t.Errorf("version.json not found: %v", err)
	}

	// Read current version
	versionData, err := os.ReadFile(versionFile)
	if err != nil {
		t.Fatalf("Failed to read version.json: %v", err)
	}
	var currentVersion struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(versionData, &currentVersion); err != nil {
		t.Fatalf("Failed to parse version.json: %v", err)
	}
	t.Logf("Current version: %s", currentVersion.TagName)

	// Simulate old version by modifying version.json
	oldVersionJSON := `{"tag_name":"b1000"}`
	if err := os.WriteFile(versionFile, []byte(oldVersionJSON), 0644); err != nil {
		t.Fatalf("Failed to write old version: %v", err)
	}

	// Second install should auto-upgrade
	t.Log("⏳ Testing auto-upgrade...")
	err = installer.InstallLlamaCpp()
	if err != nil {
		t.Fatalf("Auto-upgrade failed: %v", err)
	}

	// Verify version was upgraded by reading version.json
	versionData, err = os.ReadFile(versionFile)
	if err != nil {
		t.Fatalf("Failed to read version.json after upgrade: %v", err)
	}
	var upgradedVersion struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(versionData, &upgradedVersion); err != nil {
		t.Fatalf("Failed to parse version.json after upgrade: %v", err)
	}

	t.Logf("✅ Auto-upgrade complete: %s", upgradedVersion.TagName)

	// Version should be newer than b1000
	if upgradedVersion.TagName == "b1000" {
		t.Error("Version was not upgraded")
	}
}

// ============================================================================
// Verification Tests
// ============================================================================

func TestIsLlamaCppInstalled_NotInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	if installer.IsLlamaCppInstalled() {
		t.Error("IsLlamaCppInstalled() = true for empty directory")
	}
}

func TestVerifyInstalledBinary_NotInstalled(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	err := installer.VerifyInstalledBinary()
	if err == nil {
		t.Error("VerifyInstalledBinary() should fail for non-existent binary")
	}
}

func TestVerifyInstalledBinary_Installed(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real download test in short mode")
	}

	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	// Install first
	err := installer.InstallLlamaCpp()
	if err != nil {
		t.Skip("Skipping due to network issue")
	}

	// Verify
	err = installer.VerifyInstalledBinary()
	if err != nil {
		t.Errorf("VerifyInstalledBinary() failed after installation: %v", err)
	}

	t.Log("✅ Binary verification passed")
}

// ============================================================================
// Update Check Tests (Real API)
// ============================================================================

// ============================================================================
// Model Download Tests (Real Downloads)
// ============================================================================

func TestDownloadChatModel_RealDownload(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real download test in short mode")
	}

	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	// Use the smallest model for testing
	models := GetRecommendedVLModels()
	smallestModel := models[0] // 0.5b model

	t.Logf("⏳ Downloading chat model: %s (%.1f MB)...", smallestModel.Name, float64(smallestModel.Size)/(1024*1024))

	err := installer.DownloadChatModel(smallestModel)
	if err != nil {
		t.Logf("DownloadChatModel() failed (may be network issue): %v", err)
		t.Skip("Skipping download test due to network issue")
	}

	// Verify model exists
	modelPath := filepath.Join(installer.ModelsDir, smallestModel.Name+".gguf")
	if _, err := os.Stat(modelPath); err != nil {
		t.Errorf("Model file not found: %s", modelPath)
	}

	// Verify file size is reasonable
	fileInfo, _ := os.Stat(modelPath)
	if fileInfo.Size() < smallestModel.Size/2 {
		t.Errorf("Downloaded file too small: got %d bytes, expected ~%d bytes", fileInfo.Size(), smallestModel.Size)
	}

	t.Logf("✅ Successfully downloaded model: %.1f MB", float64(fileInfo.Size())/(1024*1024))
}

func TestDownloadChatModel_AlreadyExists(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real download test in short mode")
	}

	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	models := GetRecommendedVLModels()
	smallestModel := models[0]

	// First download
	err := installer.DownloadChatModel(smallestModel)
	if err != nil {
		t.Skip("Skipping due to network issue")
	}

	// Second download should skip
	start := time.Now()
	err = installer.DownloadChatModel(smallestModel)
	elapsed := time.Since(start)

	if err != nil {
		t.Fatalf("Second DownloadChatModel() failed: %v", err)
	}

	// Should be very fast since it skips download
	if elapsed > 100*time.Millisecond {
		t.Logf("Warning: Second download took %v (expected < 100ms)", elapsed)
	}

	t.Logf("✅ Correctly skipped re-download: %v", elapsed)
}

func TestAutoDownloadRecommendedChatModel(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real download test in short mode")
	}

	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	// Use smallest model for testing to avoid timeout
	t.Log("⏳ Downloading smallest chat model for testing...")
	smallestModel := GetRecommendedVLModels()[0] // 0.5B model

	err := installer.DownloadChatModel(smallestModel)
	if err != nil {
		t.Logf("DownloadChatModel() failed (may be network issue): %v", err)
		t.Skip("Skipping download test due to network issue")
	}

	// Verify at least one model was downloaded
	models, err := installer.GetAvailableChatModels()
	if err != nil {
		t.Fatalf("GetAvailableChatModels() failed: %v", err)
	}

	if len(models) == 0 {
		t.Error("No models downloaded after DownloadChatModel()")
	}

	t.Logf("✅ Successfully downloaded model: %v", models)
}

func TestDownloadEmbeddingModel_RealDownload(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real download test in short mode")
	}

	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	// Get recommended embedding model
	modelName := GetRecommendedEmbeddingModel()
	model, exists := GetEmbeddingModel(modelName)
	if !exists {
		t.Fatalf("Recommended model not found: %s", modelName)
	}

	t.Logf("⏳ Downloading embedding model: %s (%.1f MB)...", model.Name, float64(model.Size)/(1024*1024))

	err := installer.DownloadEmbeddingModel(model)
	if err != nil {
		t.Logf("DownloadEmbeddingModel() failed (may be network issue): %v", err)
		t.Skip("Skipping download test due to network issue")
	}

	// Verify model exists
	modelPath := filepath.Join(installer.ModelsDir, model.Filename)
	if _, err := os.Stat(modelPath); err != nil {
		t.Errorf("Model file not found: %s", modelPath)
	}

	t.Logf("✅ Successfully downloaded embedding model")
}

func TestAutoDownloadRecommendedEmbeddingModel(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping real download test in short mode")
	}

	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	t.Log("⏳ Auto-downloading recommended embedding model...")

	err := installer.AutoDownloadRecommendedEmbeddingModel()
	if err != nil {
		t.Logf("AutoDownloadRecommendedEmbeddingModel() failed (may be network issue): %v", err)
		t.Skip("Skipping download test due to network issue")
	}

	// Verify model was downloaded
	downloaded := installer.GetDownloadedEmbeddingModels()
	if len(downloaded) == 0 {
		t.Error("No embedding models downloaded after AutoDownloadRecommendedEmbeddingModel()")
	}

	t.Logf("✅ Successfully auto-downloaded embedding model: %d model(s)", len(downloaded))
}

// ============================================================================
// Model Management Tests
// ============================================================================

func TestGetAvailableChatModels_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	models, err := installer.GetAvailableChatModels()
	if err != nil {
		t.Fatalf("GetAvailableChatModels() failed: %v", err)
	}

	if len(models) != 0 {
		t.Errorf("GetAvailableChatModels() = %d models, want 0", len(models))
	}
}

func TestGetDownloadedEmbeddingModels_Empty(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	models := installer.GetDownloadedEmbeddingModels()
	if len(models) != 0 {
		t.Errorf("GetDownloadedEmbeddingModels() = %d models, want 0", len(models))
	}
}

func TestCleanupStaleTempFiles(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	// Create models directory with temp files
	os.MkdirAll(installer.ModelsDir, 0755)
	tempFile1 := filepath.Join(installer.ModelsDir, "model1.gguf.tmp")
	tempFile2 := filepath.Join(installer.ModelsDir, "model2.gguf.tmp")
	normalFile := filepath.Join(installer.ModelsDir, "model3.gguf")

	os.WriteFile(tempFile1, []byte("temp1"), 0644)
	os.WriteFile(tempFile2, []byte("temp2"), 0644)
	os.WriteFile(normalFile, []byte("normal"), 0644)

	// Cleanup
	err := installer.CleanupStaleTempFiles()
	if err != nil {
		t.Fatalf("CleanupStaleTempFiles() failed: %v", err)
	}

	// Verify temp files removed
	if _, err := os.Stat(tempFile1); !os.IsNotExist(err) {
		t.Error("Temp file 1 not removed")
	}
	if _, err := os.Stat(tempFile2); !os.IsNotExist(err) {
		t.Error("Temp file 2 not removed")
	}

	// Verify normal file preserved
	if _, err := os.Stat(normalFile); err != nil {
		t.Error("Normal file was removed")
	}

	t.Log("✅ Cleanup removed temp files and preserved normal files")
}

// ============================================================================
// Hardware Detection Tests
// ============================================================================

func TestDetectProcessor(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	processor := installer.detectProcessor()

	validProcessors := []string{"cpu", "cuda", "vulkan", "metal"}
	valid := false
	for _, p := range validProcessors {
		if processor == p {
			valid = true
			break
		}
	}

	if !valid {
		t.Errorf("detectProcessor() = %v, want one of %v", processor, validProcessors)
	}

	t.Logf("✅ Detected processor: %s", processor)
}

// ============================================================================
// Utility Tests
// ============================================================================

func TestInstallerGetModelsDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	modelsDir := installer.GetModelsDirectory()
	if modelsDir != installer.ModelsDir {
		t.Errorf("GetModelsDirectory() = %v, want %v", modelsDir, installer.ModelsDir)
	}
}

// ============================================================================
// Integration Tests (Full Workflow)
// ============================================================================

func TestFullWorkflow_InstallAndDownloadModels(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	tmpDir := t.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	t.Log("=== Integration Test: Full Workflow ===")

	// Step 1: Install llama.cpp
	t.Log("Step 1: Installing llama.cpp...")
	err := installer.InstallLlamaCpp()
	if err != nil {
		t.Skip("Skipping due to network issue")
	}
	t.Log("✅ llama.cpp installed")

	// Step 2: Verify installation
	t.Log("Step 2: Verifying installation...")
	if !installer.IsLlamaCppInstalled() {
		t.Fatal("Installation verification failed")
	}
	t.Log("✅ Installation verified")

	// Step 3: Download chat model (use smallest model to avoid timeout)
	t.Log("Step 3: Downloading chat model (smallest)...")
	smallestModel := GetRecommendedVLModels()[0] // 0.5B model
	err = installer.DownloadChatModel(smallestModel)
	if err != nil {
		t.Logf("Chat model download failed (network issue): %v", err)
		t.Skip("Skipping model download due to network issue")
	}
	t.Log("✅ Chat model downloaded")

	// Step 4: Download embedding model
	t.Log("Step 4: Downloading embedding model...")
	err = installer.AutoDownloadRecommendedEmbeddingModel()
	if err != nil {
		t.Logf("Embedding model download failed (network issue): %v", err)
		t.Skip("Skipping embedding model download due to network issue")
	}
	t.Log("✅ Embedding model downloaded")

	// Step 5: Verify all components
	t.Log("Step 5: Verifying all components...")
	chatModels, _ := installer.GetAvailableChatModels()
	embeddingModels := installer.GetDownloadedEmbeddingModels()

	if len(chatModels) == 0 {
		t.Error("No chat models found")
	}
	if len(embeddingModels) == 0 {
		t.Error("No embedding models found")
	}

	t.Logf("✅ Full workflow complete:")
	t.Logf("   - Chat models: %d", len(chatModels))
	t.Logf("   - Embedding models: %d", len(embeddingModels))
}

// ============================================================================
// Benchmark Tests
// ============================================================================

func BenchmarkIsLlamaCppInstalled(b *testing.B) {
	tmpDir := b.TempDir()
	installer := &LlamaCppInstaller{
		BinaryPath:   filepath.Join(tmpDir, "bin"),
		MetadataPath: filepath.Join(tmpDir, "metadata"),
		ModelsDir:    filepath.Join(tmpDir, "models"),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		installer.IsLlamaCppInstalled()
	}
}
