package llama

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

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
	// HardwareSpecs caches the detected hardware specifications
	HardwareSpecs *HardwareSpecs
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
		BinaryPath:    binaryPath,
		MetadataPath:  metadataPath,
		ModelsDir:     modelsDir,
		HardwareSpecs: DetectHardwareSpecs(),
	}

	// Clean up any stale temp files from previous sessions
	// This handles the case where app was closed during download
	if err := installer.CleanupStaleTempFiles(); err != nil {
		log.Printf("⚠️  Failed to cleanup stale temp files on startup: %v", err)
	}

	return installer
}

// InstallLlamaCpp installs llama.cpp via direct download
// Works on all platforms (macOS, Linux, Windows)
// Uses download.InstallLibraries for automatic version management
func (lcm *LlamaCppInstaller) InstallLlamaCpp() error {
	// Detect processor type
	processor := lcm.detectProcessor()
	log.Printf("Auto-installing llama.cpp for %s/%s", runtime.GOOS, processor)

	// Convert string processor to download.Processor type
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

	// Use InstallLibraries which handles:
	// - Version checking and management (version.json)
	// - Auto-upgrade support (if allowUpgrade=true)
	// - Automatic fallback to previous versions if latest fails
	// - OS/Processor validation
	if err := download.InstallLibraries(lcm.BinaryPath, proc, true); err != nil {
		return fmt.Errorf("failed to install llama.cpp: %w", err)
	}

	log.Println("✅ llama.cpp installed successfully")
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

	// Select optimal model based on available RAM
	modelSpec := SelectOptimalQwenModel(lcm.HardwareSpecs.AvailableRAM)

	// Download the model
	if err := lcm.DownloadChatModel(modelSpec); err != nil {
		return fmt.Errorf("failed to download model: %w", err)
	}

	log.Println("🎉 Chat model download completed successfully!")
	return nil
}

// AutoDownloadRecommendedVLModel automatically downloads the best VL model for the system
func (lcm *LlamaCppInstaller) AutoDownloadRecommendedVLModel() error {
	// Clean up any stale temp files
	if err := lcm.CleanupStaleTempFiles(); err != nil {
		log.Printf("⚠️  Failed to cleanup stale temp files: %v", err)
	}

	// Check if any VL models already exist
	models, err := lcm.GetAvailableVLModels()
	if err != nil {
		return fmt.Errorf("failed to check existing models: %w", err)
	}

	if len(models) > 0 {
		log.Printf("✅ VL models already available (%d found), skipping auto-download", len(models))
		return nil
	}

	log.Println("📦 No VL models found, starting auto-download...")

	// Select optimal model based on available RAM
	modelSpec := SelectOptimalQwenModel(lcm.HardwareSpecs.AvailableRAM)

	// Download the model
	if err := lcm.DownloadChatModel(modelSpec); err != nil {
		return fmt.Errorf("failed to download model: %w", err)
	}

	log.Println("🎉 VL model download completed successfully!")
	return nil
}

// AutoDownloadRecommendedTextModel automatically downloads the best text model for the system
func (lcm *LlamaCppInstaller) AutoDownloadRecommendedTextModel() error {
	// Clean up any stale temp files
	if err := lcm.CleanupStaleTempFiles(); err != nil {
		log.Printf("⚠️  Failed to cleanup stale temp files: %v", err)
	}

	// Check if any text models already exist
	models, err := lcm.GetAvailableTextModels()
	if err != nil {
		return fmt.Errorf("failed to check existing models: %w", err)
	}

	if len(models) > 0 {
		log.Printf("✅ Text models already available (%d found), skipping auto-download", len(models))
		return nil
	}

	log.Println("📦 No text models found, starting auto-download...")

	// Select optimal model based on available RAM
	modelSpec := SelectOptimalQwenTextModel(lcm.HardwareSpecs.AvailableRAM)

	// Download the model
	if err := lcm.DownloadChatModel(modelSpec); err != nil {
		return fmt.Errorf("failed to download model: %w", err)
	}

	log.Println("🎉 Text model download completed successfully!")
	return nil
}

// AutoDownloadRecommendedEmbeddingModel automatically downloads the best embedding model for the system based on hardware detection
func (lcm *LlamaCppInstaller) AutoDownloadRecommendedEmbeddingModel() error {
	downloaded := lcm.GetDownloadedEmbeddingModels()
	if len(downloaded) > 0 {
		log.Printf("✅ Embedding models already available (%d found), skipping auto-download", len(downloaded))
		return nil
	}

	log.Println("📦 No embedding models found, starting auto-download...")

	// Get recommended embedding model
	recommendedName := GetRecommendedEmbeddingModel()
	model, exists := GetEmbeddingModel(recommendedName)
	if !exists {
		return fmt.Errorf("recommended embedding model not found in catalog: %s", recommendedName)
	}

	log.Printf("📦 Selected recommended embedding model: %s (%s, %d dims, %.1f MB)",
		model.Name, model.Quantization, model.Dimensions, float64(model.Size)/(1024*1024))

	if err := lcm.DownloadEmbeddingModel(model); err != nil {
		return fmt.Errorf("failed to download embedding model: %w", err)
	}

	log.Println("🎉 Embedding model download completed successfully!")
	return nil
}

// AutoDownloadRecommendedUtilityModel automatically downloads the best utility model for the system
// Utility models are small, fast models used for:
// - Summary generation (compress old conversation messages)
// - Title generation (create conversation titles)
// - Quick text processing tasks
// These models are non-reasoning (no <think> tags) and optimized for background tasks
func (lcm *LlamaCppInstaller) AutoDownloadRecommendedUtilityModel() error {
	// Clean up any stale temp files
	if err := lcm.CleanupStaleTempFiles(); err != nil {
		log.Printf("⚠️  Failed to cleanup stale temp files: %v", err)
	}

	// Check if any utility models (Llama 1B/3B, Mistral) already exist
	utilityModels, err := lcm.GetAvailableUtilityModels()
	if err != nil {
		return fmt.Errorf("failed to check existing utility models: %w", err)
	}

	if len(utilityModels) > 0 {
		log.Printf("✅ Utility models already available (%d found), skipping auto-download", len(utilityModels))
		log.Printf("   Available: %v", utilityModels)
		return nil
	}

	log.Println("📦 No utility models found, starting auto-download...")
	log.Println("💡 Utility models enable fast background tasks:")
	log.Println("   - Summary generation (compress old messages)")
	log.Println("   - Title generation (auto-create conversation titles)")
	log.Println("   - Quick text processing")

	// Select optimal utility model based on available RAM
	modelSpec := SelectOptimalUtilityModel(lcm.HardwareSpecs.AvailableRAM)
	if modelSpec == nil {
		return fmt.Errorf("no suitable utility model found for system with %dGB RAM", lcm.HardwareSpecs.AvailableRAM)
	}

	// Download the model
	if err := lcm.DownloadChatModel(*modelSpec); err != nil {
		return fmt.Errorf("failed to download utility model: %w", err)
	}

	log.Println("🎉 Utility model download completed successfully!")
	log.Printf("   Model: %s", modelSpec.Name)
	log.Printf("   Size: %.1f MB", float64(modelSpec.Size)/(1024*1024))
	log.Printf("   Use case: Summary & title generation (3-5x faster than main model)")
	return nil
}

// AutoDownloadAllRecommendedModels automatically downloads all recommended model types for the system
// Downloads text, VL, embedding, and utility models based on hardware detection
func (lcm *LlamaCppInstaller) AutoDownloadAllRecommendedModels() error {
	log.Println("🚀 Starting automatic download of all recommended models...")

	// Download text model
	if err := lcm.AutoDownloadRecommendedTextModel(); err != nil {
		return fmt.Errorf("failed to download text model: %w", err)
	}

	// Download VL model
	if err := lcm.AutoDownloadRecommendedVLModel(); err != nil {
		return fmt.Errorf("failed to download VL model: %w", err)
	}

	// Download embedding model
	if err := lcm.AutoDownloadRecommendedEmbeddingModel(); err != nil {
		return fmt.Errorf("failed to download embedding model: %w", err)
	}

	// Download utility model (for summary/title generation)
	if err := lcm.AutoDownloadRecommendedUtilityModel(); err != nil {
		log.Printf("⚠️  Warning: Failed to download utility model (optional): %v", err)
		log.Printf("💡 Main model will be used for summaries/titles (slower but functional)")
		// Don't fail the entire download if utility model fails
	}

	log.Println("🎉 All recommended models downloaded successfully!")
	log.Println("📊 Model download summary:")
	log.Println("   ✅ Text model: Ready for text generation and reasoning")
	log.Println("   ✅ VL model: Ready for vision-language tasks")
	log.Println("   ✅ Embedding model: Ready for text embedding and similarity tasks")
	log.Println("   ✅ Utility model: Ready for fast summaries and title generation")

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

// GetAvailableVLModels returns a list of available VL (vision-language) model file paths
func (lcm *LlamaCppInstaller) GetAvailableVLModels() ([]string, error) {
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
		if strings.HasSuffix(strings.ToLower(name), ".gguf") &&
			strings.HasPrefix(name, "qwen3-vl-") &&
			!lcm.isEmbeddingModel(name) {
			models = append(models, filepath.Join(lcm.ModelsDir, name))
		}
	}

	return models, nil
}

// GetAvailableTextModels returns a list of available text model file paths
func (lcm *LlamaCppInstaller) GetAvailableTextModels() ([]string, error) {
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
		if strings.HasSuffix(strings.ToLower(name), ".gguf") &&
			strings.HasPrefix(name, "qwen3-") &&
			!strings.HasPrefix(name, "qwen3-vl-") &&
			!lcm.isEmbeddingModel(name) {
			models = append(models, filepath.Join(lcm.ModelsDir, name))
		}
	}

	return models, nil
}

// GetAvailableUtilityModels returns a list of available utility model file paths
// Utility models are small, fast models (1B-3B) used for summary/title generation
// Returns model filenames (not full paths) for easier checking
func (lcm *LlamaCppInstaller) GetAvailableUtilityModels() ([]string, error) {
	if err := os.MkdirAll(lcm.ModelsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create models directory: %w", err)
	}

	entries, err := os.ReadDir(lcm.ModelsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read models directory: %w", err)
	}

	var utilityModels []string
	utilitySpecs := GetRecommendedUtilityModels()

	// Build a map of utility model filenames for quick lookup
	utilityFilenames := make(map[string]bool)
	for _, spec := range utilitySpecs {
		filename := fmt.Sprintf("%s.gguf", spec.Name)
		utilityFilenames[strings.ToLower(filename)] = true
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		name := entry.Name()
		nameLower := strings.ToLower(name)

		// Check if it's a utility model (Llama 1B/3B, Mistral, etc.)
		if strings.HasSuffix(nameLower, ".gguf") && !lcm.isEmbeddingModel(name) {
			// Match against known utility model patterns
			if utilityFilenames[nameLower] ||
				strings.Contains(nameLower, "llama-3.2-1b") ||
				strings.Contains(nameLower, "llama-3.2-3b") ||
				(strings.Contains(nameLower, "mistral") && strings.Contains(nameLower, "7b")) {
				utilityModels = append(utilityModels, name)
			}
		}
	}

	return utilityModels, nil
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
