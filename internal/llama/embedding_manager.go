package llama

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// EmbeddingModel represents an embedding model configuration
type EmbeddingModel struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	URL          string   `json:"url"`
	Filename     string   `json:"filename"`
	Size         int64    `json:"size"`
	SHA256       string   `json:"sha256"`
	Quantization string   `json:"quantization"`
	Dimensions   int      `json:"dimensions"`
	Languages    []string `json:"languages"`
	Type         string   `json:"type"` // "embedding" or "llm"
}

// EmbeddingManager handles embedding model downloads and management
type EmbeddingManager struct {
	ModelsDir       string
	DownloadDir     string
	AvailableModels map[string]*EmbeddingModel
}

// NewEmbeddingManager creates a new embedding manager
func NewEmbeddingManager() *EmbeddingManager {
	homeDir, _ := os.UserHomeDir()
	modelsDir := filepath.Join(homeDir, ".veridium", "models", "embeddings")
	downloadDir := filepath.Join(homeDir, ".veridium", "downloads")

	// Ensure directories exist
	os.MkdirAll(modelsDir, 0755)
	os.MkdirAll(downloadDir, 0755)

	em := &EmbeddingManager{
		ModelsDir:       modelsDir,
		DownloadDir:     downloadDir,
		AvailableModels: make(map[string]*EmbeddingModel),
	}

	// Initialize available models
	em.initializeAvailableModels()

	return em
}

// validateGGUFFile validates a GGUF file by checking its header structure
func (em *EmbeddingManager) validateGGUFFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// GGUF files start with a magic number: "GGUF" (4 bytes)
	magic := make([]byte, 4)
	if _, err := io.ReadFull(file, magic); err != nil {
		return fmt.Errorf("failed to read magic bytes: %w", err)
	}

	if string(magic) != "GGUF" {
		return fmt.Errorf("invalid GGUF file: wrong magic bytes (expected 'GGUF', got '%s')", string(magic))
	}

	// Read version (4 bytes, little-endian uint32)
	var version uint32
	if err := binary.Read(file, binary.LittleEndian, &version); err != nil {
		return fmt.Errorf("failed to read version: %w", err)
	}

	// Check if version is reasonable (GGUF versions are typically 1, 2, or 3)
	if version == 0 || version > 10 {
		return fmt.Errorf("invalid GGUF file: unreasonable version %d", version)
	}

	// Read tensor count (8 bytes, little-endian uint64)
	var tensorCount uint64
	if err := binary.Read(file, binary.LittleEndian, &tensorCount); err != nil {
		return fmt.Errorf("failed to read tensor count: %w", err)
	}

	// Read metadata key-value count (8 bytes, little-endian uint64)
	var metadataCount uint64
	if err := binary.Read(file, binary.LittleEndian, &metadataCount); err != nil {
		return fmt.Errorf("failed to read metadata count: %w", err)
	}

	// Basic sanity checks
	if tensorCount == 0 {
		return fmt.Errorf("invalid GGUF file: no tensors found")
	}

	if tensorCount > 100000 { // Reasonable upper bound
		return fmt.Errorf("invalid GGUF file: unreasonable tensor count %d", tensorCount)
	}

	if metadataCount > 10000 { // Reasonable upper bound for metadata entries
		return fmt.Errorf("invalid GGUF file: unreasonable metadata count %d", metadataCount)
	}

	log.Printf("GGUF file validation passed: version=%d, tensors=%d, metadata_entries=%d",
		version, tensorCount, metadataCount)
	return nil
}

// initializeAvailableModels sets up the catalog of available embedding models
func (em *EmbeddingManager) initializeAvailableModels() {
	// Granite Embedding 107M Multilingual models from bartowski
	// High-quality multilingual embedding model with multiple quantization options
	em.AvailableModels["granite-embedding-107m-multilingual-q6_k_l"] = &EmbeddingModel{
		Name:         "granite-embedding-107m-multilingual-q6_k_l",
		Description:  "Granite Embedding 107M Multilingual (Q6_K_L) - Very high quality, recommended",
		URL:          "https://huggingface.co/bartowski/granite-embedding-107m-multilingual-GGUF/resolve/main/granite-embedding-107m-multilingual-Q6_K_L.gguf",
		Filename:     "granite-embedding-107m-multilingual-Q6_K_L.gguf",
		Size:         120163008, // Actual size from successful download
		Quantization: "Q6_K_L",
		Dimensions:   384,
		Languages:    []string{"en", "de", "fr", "it", "pt", "es", "nl", "pl", "ru", "zh", "ja", "ko", "ar", "hi", "tr", "th", "vi", "id"},
		Type:         "embedding",
	}

	em.AvailableModels["granite-embedding-107m-multilingual-q4_k_m"] = &EmbeddingModel{
		Name:         "granite-embedding-107m-multilingual-q4_k_m",
		Description:  "Granite Embedding 107M Multilingual (Q4_K_M) - Good quality, smaller size",
		URL:          "https://huggingface.co/bartowski/granite-embedding-107m-multilingual-GGUF/resolve/main/granite-embedding-107m-multilingual-Q4_K_M.gguf",
		Filename:     "granite-embedding-107m-multilingual-Q4_K_M.gguf",
		Size:         123000000, // ~123MB
		Quantization: "Q4_K_M",
		Dimensions:   384,
		Languages:    []string{"en", "de", "fr", "it", "pt", "es", "nl", "pl", "ru", "zh", "ja", "ko", "ar", "hi", "tr", "th", "vi", "id"},
		Type:         "embedding",
	}

	em.AvailableModels["granite-embedding-107m-multilingual-f16"] = &EmbeddingModel{
		Name:         "granite-embedding-107m-multilingual-f16",
		Description:  "Granite Embedding 107M Multilingual (F16) - Full precision, highest quality",
		URL:          "https://huggingface.co/bartowski/granite-embedding-107m-multilingual-GGUF/resolve/main/granite-embedding-107m-multilingual-f16.gguf",
		Filename:     "granite-embedding-107m-multilingual-f16.gguf",
		Size:         236000000, // ~236MB
		Quantization: "F16",
		Dimensions:   384,
		Languages:    []string{"en", "de", "fr", "it", "pt", "es", "nl", "pl", "ru", "zh", "ja", "ko", "ar", "hi", "tr", "th", "vi", "id"},
		Type:         "embedding",
	}

	// Paraphrase Multilingual MiniLM L12 v2 from mykor
	em.AvailableModels["paraphrase-multilingual-minilm-l12-v2-f16"] = &EmbeddingModel{
		Name:         "paraphrase-multilingual-minilm-l12-v2-f16",
		Description:  "Paraphrase Multilingual MiniLM L12 v2 (F16) - Alternative multilingual model",
		URL:          "https://huggingface.co/mykor/paraphrase-multilingual-MiniLM-L12-v2.gguf/resolve/main/paraphrase-multilingual-MiniLM-L12-118M-v2-F16.gguf",
		Filename:     "paraphrase-multilingual-MiniLM-L12-118M-v2-F16.gguf",
		Size:         242358176, // ~242MB
		Quantization: "F16",
		Dimensions:   384,
		Languages:    []string{"en", "de", "fr", "it", "pt", "es", "nl", "pl", "ru", "zh", "ja", "ko"},
		Type:         "embedding",
	}

	// Nomic Embed Text v1.5 (Popular choice for embeddings)
	em.AvailableModels["nomic-embed-text-v1.5-q4_k_m"] = &EmbeddingModel{
		Name:         "nomic-embed-text-v1.5-q4_k_m",
		Description:  "Nomic Embed Text v1.5 (Q4_K_M) - Popular, high-quality embeddings",
		URL:          "https://huggingface.co/nomic-ai/nomic-embed-text-v1.5-GGUF/resolve/main/nomic-embed-text-v1.5.Q4_K_M.gguf",
		Filename:     "nomic-embed-text-v1.5.Q4_K_M.gguf",
		Size:         550000000, // ~550MB
		Quantization: "Q4_K_M",
		Dimensions:   768,
		Languages:    []string{"en", "de", "fr", "it", "pt", "es", "nl", "pl", "ru", "zh", "ja", "ko", "ar", "hi", "tr", "th", "vi", "id"},
		Type:         "embedding",
	}

	log.Printf("Initialized %d embedding models in catalog", len(em.AvailableModels))
	log.Printf("Recommended: granite-embedding-107m-multilingual-q6_k_l (120MB, very high quality)")
}

// GetAvailableModels returns the list of available models
func (em *EmbeddingManager) GetAvailableModels() map[string]*EmbeddingModel {
	return em.AvailableModels
}

// GetModel returns a specific model by name
func (em *EmbeddingManager) GetModel(name string) (*EmbeddingModel, bool) {
	model, exists := em.AvailableModels[name]
	return model, exists
}

// IsModelDownloaded checks if a model is already downloaded
func (em *EmbeddingManager) IsModelDownloaded(modelName string) bool {
	model, exists := em.AvailableModels[modelName]
	if !exists {
		return false
	}

	modelPath := filepath.Join(em.ModelsDir, model.Filename)
	if _, err := os.Stat(modelPath); err == nil {
		// Validate the GGUF file structure
		if err := em.validateGGUFFile(modelPath); err != nil {
			log.Printf("Model %s exists but failed validation: %v", modelName, err)
			return false
		}
		return true
	}
	return false
}

// GetModelPath returns the local path to a downloaded model
func (em *EmbeddingManager) GetModelPath(modelName string) (string, error) {
	model, exists := em.AvailableModels[modelName]
	if !exists {
		return "", fmt.Errorf("model %s not found in catalog", modelName)
	}

	modelPath := filepath.Join(em.ModelsDir, model.Filename)
	if _, err := os.Stat(modelPath); err != nil {
		return "", fmt.Errorf("model %s not downloaded: %w", modelName, err)
	}

	return modelPath, nil
}

// DownloadModel downloads an embedding model
func (em *EmbeddingManager) DownloadModel(modelName string, progressCallback func(downloaded, total int64)) error {
	model, exists := em.AvailableModels[modelName]
	if !exists {
		return fmt.Errorf("model %s not found in catalog", modelName)
	}

	// Check if already downloaded
	if em.IsModelDownloaded(modelName) {
		log.Printf("Model %s already downloaded", modelName)
		return nil
	}

	log.Printf("Starting download of embedding model: %s", model.Name)
	log.Printf("URL: %s", model.URL)
	log.Printf("Size: %.2f MB", float64(model.Size)/1024/1024)

	// Create temporary download path
	tempPath := filepath.Join(em.DownloadDir, model.Filename+".tmp")
	finalPath := filepath.Join(em.ModelsDir, model.Filename)

	// Create HTTP request
	req, err := http.NewRequest("GET", model.URL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set user agent
	req.Header.Set("User-Agent", fmt.Sprintf("Veridium/1.0 (%s %s)", runtime.GOOS, runtime.GOARCH))

	// Make request
	client := &http.Client{
		Timeout: 30 * time.Minute, // Long timeout for large models
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download model: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d %s", resp.StatusCode, resp.Status)
	}

	// Create temporary file
	tempFile, err := os.Create(tempPath)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer tempFile.Close()
	defer os.Remove(tempPath) // Clean up temp file on error

	// Get content length
	contentLength := resp.ContentLength
	if contentLength <= 0 && model.Size > 0 {
		contentLength = model.Size
	}

	// Download with progress tracking
	var downloaded int64
	buffer := make([]byte, 32*1024) // 32KB buffer
	hasher := sha256.New()

	for {
		n, err := resp.Body.Read(buffer)
		if n > 0 {
			// Write to file
			if _, writeErr := tempFile.Write(buffer[:n]); writeErr != nil {
				return fmt.Errorf("failed to write to temp file: %w", writeErr)
			}

			// Update hash
			hasher.Write(buffer[:n])

			// Update progress
			downloaded += int64(n)
			if progressCallback != nil {
				progressCallback(downloaded, contentLength)
			}
		}

		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}
	}

	// Close temp file before moving
	tempFile.Close()

	// Verify download completeness
	if contentLength > 0 && downloaded != contentLength {
		os.Remove(tempPath)
		return fmt.Errorf("download incomplete: expected %d bytes, got %d bytes", contentLength, downloaded)
	}

	// Validate GGUF file structure
	if err := em.validateGGUFFile(tempPath); err != nil {
		os.Remove(tempPath)
		return fmt.Errorf("downloaded file failed GGUF validation: %w", err)
	}

	log.Printf("Downloaded model %s passed GGUF validation", model.Name)

	// Verify SHA256 if provided
	if model.SHA256 != "" {
		calculatedHash := fmt.Sprintf("%x", hasher.Sum(nil))
		if !strings.EqualFold(calculatedHash, model.SHA256) {
			return fmt.Errorf("SHA256 mismatch: expected %s, got %s", model.SHA256, calculatedHash)
		}
		log.Printf("SHA256 verification passed for %s", model.Name)
	}

	// Move temp file to final location
	if err := os.Rename(tempPath, finalPath); err != nil {
		return fmt.Errorf("failed to move downloaded file: %w", err)
	}

	log.Printf("Successfully downloaded embedding model: %s to %s", model.Name, finalPath)
	return nil
}

// DeleteModel removes a downloaded model
func (em *EmbeddingManager) DeleteModel(modelName string) error {
	model, exists := em.AvailableModels[modelName]
	if !exists {
		return fmt.Errorf("model %s not found in catalog", modelName)
	}

	modelPath := filepath.Join(em.ModelsDir, model.Filename)
	if err := os.Remove(modelPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("model %s is not downloaded", modelName)
		}
		return fmt.Errorf("failed to delete model: %w", err)
	}

	log.Printf("Successfully deleted embedding model: %s", model.Name)
	return nil
}

// GetDownloadedModels returns a list of downloaded models
func (em *EmbeddingManager) GetDownloadedModels() []*EmbeddingModel {
	var downloaded []*EmbeddingModel

	for _, model := range em.AvailableModels {
		if em.IsModelDownloaded(model.Name) {
			downloaded = append(downloaded, model)
		}
	}

	return downloaded
}

// GetRecommendedModel returns the recommended embedding model
func (em *EmbeddingManager) GetRecommendedModel() string {
	// Recommend granite-embedding-107m-multilingual-q6_k_l as it offers:
	// - Very high quality (Q6_K_L quantization)
	// - Smaller size (120MB vs 550MB for nomic)
	// - Excellent multilingual support (18+ languages)
	// - From bartowski (trusted quantizer)
	return "granite-embedding-107m-multilingual-q6_k_l"
}

// GetModelInfo returns detailed information about a model
func (em *EmbeddingManager) GetModelInfo(modelName string) (map[string]interface{}, error) {
	model, exists := em.AvailableModels[modelName]
	if !exists {
		return nil, fmt.Errorf("model %s not found", modelName)
	}

	info := map[string]interface{}{
		"name":         model.Name,
		"description":  model.Description,
		"filename":     model.Filename,
		"size":         model.Size,
		"size_mb":      float64(model.Size) / 1024 / 1024,
		"quantization": model.Quantization,
		"dimensions":   model.Dimensions,
		"languages":    model.Languages,
		"type":         model.Type,
		"downloaded":   em.IsModelDownloaded(modelName),
	}

	if em.IsModelDownloaded(modelName) {
		modelPath := filepath.Join(em.ModelsDir, model.Filename)
		if stat, err := os.Stat(modelPath); err == nil {
			info["local_path"] = modelPath
			info["local_size"] = stat.Size()
			info["modified_time"] = stat.ModTime()
		}
	}

	return info, nil
}

// CleanupDownloads removes temporary download files
func (em *EmbeddingManager) CleanupDownloads() error {
	entries, err := os.ReadDir(em.DownloadDir)
	if err != nil {
		return fmt.Errorf("failed to read download directory: %w", err)
	}

	cleaned := 0
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".tmp") {
			tmpPath := filepath.Join(em.DownloadDir, entry.Name())
			if err := os.Remove(tmpPath); err != nil {
				log.Printf("Failed to remove temp file %s: %v", tmpPath, err)
			} else {
				cleaned++
			}
		}
	}

	if cleaned > 0 {
		log.Printf("Cleaned up %d temporary download files", cleaned)
	}

	return nil
}

// GetStorageUsage returns storage usage information
func (em *EmbeddingManager) GetStorageUsage() (map[string]interface{}, error) {
	var totalSize int64
	modelCount := 0

	entries, err := os.ReadDir(em.ModelsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read models directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".gguf") {
			if info, err := entry.Info(); err == nil {
				totalSize += info.Size()
				modelCount++
			}
		}
	}

	usage := map[string]interface{}{
		"models_directory": em.ModelsDir,
		"total_size":       totalSize,
		"total_size_mb":    float64(totalSize) / 1024 / 1024,
		"model_count":      modelCount,
		"available_models": len(em.AvailableModels),
	}

	return usage, nil
}

// AutoDownloadRecommendedModel automatically downloads the recommended embedding model
func (em *EmbeddingManager) AutoDownloadRecommendedModel() error {
	// Check if any models already exist
	downloaded := em.GetDownloadedModels()
	if len(downloaded) > 0 {
		log.Printf("✅ Embedding models already available (%d found), skipping auto-download", len(downloaded))
		return nil
	}

	log.Println("📦 No embedding models found, starting auto-download...")

	// Get recommended model
	modelName := em.GetRecommendedModel()
	model, exists := em.AvailableModels[modelName]
	if !exists {
		return fmt.Errorf("recommended model not found: %s", modelName)
	}

	log.Printf("📥 Downloading recommended embedding model: %s", model.Name)
	log.Printf("   Size: %.1f MB", float64(model.Size)/1024/1024)
	log.Printf("   Dimensions: %d", model.Dimensions)
	log.Printf("   Languages: %d supported", len(model.Languages))

	// Download with progress callback
	err := em.DownloadModel(modelName, func(downloaded, total int64) {
		if total > 0 {
			progress := float64(downloaded) / float64(total) * 100
			if downloaded%(5*1024*1024) == 0 { // Log every 5MB
				log.Printf("   Progress: %.1f%% (%.1f MB / %.1f MB)",
					progress,
					float64(downloaded)/1024/1024,
					float64(total)/1024/1024)
			}
		}
	})

	if err != nil {
		return fmt.Errorf("failed to download embedding model: %w", err)
	}

	log.Println("🎉 Embedding model download completed successfully!")
	return nil
}
