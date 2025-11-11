package llama

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// Constants for GGUF file validation
const (
	ggufMagicNumber     = "GGUF"
	ggufMagicNumberSize = 4
	ggufVersionSize     = 4
	ggufTensorCountSize = 8
	ggufMetadataSize    = 8
	maxGGUFVersion      = 10
	maxTensorCount      = 100000
	maxMetadataCount    = 10000
	sizeVarianceMin     = 0.95
	sizeVarianceMax     = 1.05
)

// QwenModelSpec represents a Qwen model specification for direct download
type QwenModelSpec struct {
	Name         string // Model name (e.g., "qwen2.5-0.5b-instruct-q4_k_m")
	URL          string // Direct download URL
	Quantization string // Quantization type (Q4_K_M, Q5_K_M, etc.)
	Parameters   string // Parameter size (0.5b, 1.5b, etc.)
	MinRAM       int64  // Minimum RAM required in GB
	Size         int64  // Expected file size in bytes
	SHA256       string // Expected SHA256 checksum (optional)
	Description  string // Model description
}

// Note: Service struct has been removed. All functions are now pure functions
// that accept modelsDir as a parameter or use GetModelsDirectory() directly.

// GetRecommendedQwenModels returns recommended Qwen models for direct download
// Models are ordered from smallest to largest by MinRAM requirement
func GetRecommendedQwenModels() []QwenModelSpec {
	return []QwenModelSpec{
		{
			Name:         "qwen2.5-0.5b-instruct-q4_k_m",
			URL:          "https://huggingface.co/Qwen/Qwen2.5-0.5B-Instruct-GGUF/resolve/main/qwen2.5-0.5b-instruct-q4_k_m.gguf",
			Quantization: "Q4_K_M",
			Parameters:   "0.5b",
			MinRAM:       2,
			Size:         491520000, // ~468 MB
			Description:  "Smallest Qwen model, perfect for low-end hardware and testing",
		},
		{
			Name:         "qwen2.5-1.5b-instruct-q4_k_m",
			URL:          "https://huggingface.co/Qwen/Qwen2.5-1.5B-Instruct-GGUF/resolve/main/qwen2.5-1.5b-instruct-q4_k_m.gguf",
			Quantization: "Q4_K_M",
			Parameters:   "1.5b",
			MinRAM:       4,
			Size:         1100000000, // ~1.1 GB
			Description:  "Lightweight Qwen model, good balance of speed and quality",
		},
		{
			Name:         "qwen2.5-3b-instruct-q4_k_m",
			URL:          "https://huggingface.co/Qwen/Qwen2.5-3B-Instruct-GGUF/resolve/main/qwen2.5-3b-instruct-q4_k_m.gguf",
			Quantization: "Q4_K_M",
			Parameters:   "3b",
			MinRAM:       6,
			Size:         2100000000, // ~2.1 GB
			Description:  "Mid-size Qwen model, excellent for most tasks",
		},
		{
			Name:         "qwen2.5-7b-instruct-q4_k_m",
			URL:          "https://huggingface.co/Qwen/Qwen2.5-7B-Instruct-GGUF/resolve/main/qwen2.5-7b-instruct-q4_k_m.gguf",
			Quantization: "Q4_K_M",
			Parameters:   "7b",
			MinRAM:       10,
			Size:         4800000000, // ~4.8 GB
			Description:  "High-quality Qwen model, great for advanced tasks",
		},
	}
}

// SelectOptimalQwenModel selects the best model based on available RAM.
// It returns the largest model that fits within the available RAM.
// If no model fits, it returns the smallest model as a fallback.
// Note: Currently excludes 7B model as it requires multi-file download support
func SelectOptimalQwenModel(availableRAM int64) QwenModelSpec {
	models := GetRecommendedQwenModels()

	// Select the largest model that fits in available RAM
	// But skip 7B model for now (multi-file download not yet supported)
	var selectedModel QwenModelSpec
	for _, model := range models {
		// Skip 7B model (requires multi-file download)
		if model.Parameters == "7b" {
			continue
		}

		if model.MinRAM <= availableRAM {
			selectedModel = model
		} else {
			break // Models are ordered, so we can stop here
		}
	}

	// If no model fits, use the smallest one as fallback
	if selectedModel.Name == "" {
		selectedModel = models[0]
		log.Printf("⚠️  System has low RAM (%dGB), using smallest model", availableRAM)
	}

	log.Printf("📦 Selected model: %s (%s, %s) - requires %dGB RAM (system has %dGB)",
		selectedModel.Name, selectedModel.Parameters, selectedModel.Quantization,
		selectedModel.MinRAM, availableRAM)

	return selectedModel
}

// Note: DownloadModel, CleanupStaleTempFiles, and GetAvailableModels have been moved to installer.go
// Use LlamaCppInstaller methods instead:
//   - installer.DownloadChatModel(modelSpec)
//   - installer.CleanupStaleTempFiles()
//   - installer.GetAvailableChatModels()

// validateDownloadedFile performs all validation checks on a downloaded file
// Used by installer.go for model validation
func validateDownloadedFile(filePath string, modelSpec QwenModelSpec) error {
	// Check file size
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("failed to stat downloaded file: %w", err)
	}

	actualSize := fileInfo.Size()
	expectedSize := modelSpec.Size
	sizeVariance := float64(actualSize) / float64(expectedSize)

	if sizeVariance < sizeVarianceMin || sizeVariance > sizeVarianceMax {
		return fmt.Errorf("downloaded file size mismatch: got %d bytes, expected ~%d bytes (%.1f%% of expected)",
			actualSize, expectedSize, sizeVariance*100)
	}

	// Verify checksum if provided
	if modelSpec.SHA256 != "" {
		log.Printf("🔒 Verifying model integrity...")
		if err := verifyModelChecksum(filePath, modelSpec.SHA256); err != nil {
			return fmt.Errorf("model integrity check failed: %w", err)
		}
		log.Printf("✅ Model integrity verified")
	}

	// Verify it's a valid GGUF file
	if err := validateGGUFFile(filePath); err != nil {
		return fmt.Errorf("invalid GGUF file: %w", err)
	}

	return nil
}

// verifyModelChecksum verifies the SHA256 checksum of a file
// Used by installer.go for model validation
func verifyModelChecksum(filePath, expectedChecksum string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file for checksum: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return fmt.Errorf("failed to compute checksum: %w", err)
	}

	actualChecksum := hex.EncodeToString(hash.Sum(nil))
	if !strings.EqualFold(actualChecksum, expectedChecksum) {
		return fmt.Errorf("checksum mismatch: got %s, expected %s", actualChecksum, expectedChecksum)
	}

	return nil
}

// validateGGUFFile performs basic validation on a GGUF file by checking the magic number
// Used by installer.go for model validation
func validateGGUFFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read GGUF magic number (first 4 bytes should be "GGUF")
	magic := make([]byte, ggufMagicNumberSize)
	if _, err := io.ReadFull(file, magic); err != nil {
		return fmt.Errorf("failed to read file header: %w", err)
	}

	if string(magic) != ggufMagicNumber {
		return fmt.Errorf("invalid GGUF magic number: got %q, expected %q", string(magic), ggufMagicNumber)
	}

	return nil
}

// ============================================================================
// Embedding Model Functions
// ============================================================================

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
	Type         string   `json:"type"` // "embedding"
}

var embeddingModelsCatalog = map[string]*EmbeddingModel{
	"granite-embedding-107m-multilingual-q6_k_l": {
		Name:         "granite-embedding-107m-multilingual-q6_k_l",
		Description:  "Granite Embedding 107M Multilingual (Q6_K_L) - Very high quality, recommended",
		URL:          "https://huggingface.co/bartowski/granite-embedding-107m-multilingual-GGUF/resolve/main/granite-embedding-107m-multilingual-Q6_K_L.gguf",
		Filename:     "granite-embedding-107m-multilingual-Q6_K_L.gguf",
		Size:         120163008,
		Quantization: "Q6_K_L",
		Dimensions:   384,
		Languages:    []string{"en", "de", "fr", "it", "pt", "es", "nl", "pl", "ru", "zh", "ja", "ko", "ar", "hi", "tr", "th", "vi", "id"},
		Type:         "embedding",
	},
	"granite-embedding-107m-multilingual-q4_k_m": {
		Name:         "granite-embedding-107m-multilingual-q4_k_m",
		Description:  "Granite Embedding 107M Multilingual (Q4_K_M) - Good quality, smaller size",
		URL:          "https://huggingface.co/bartowski/granite-embedding-107m-multilingual-GGUF/resolve/main/granite-embedding-107m-multilingual-Q4_K_M.gguf",
		Filename:     "granite-embedding-107m-multilingual-Q4_K_M.gguf",
		Size:         123000000,
		Quantization: "Q4_K_M",
		Dimensions:   384,
		Languages:    []string{"en", "de", "fr", "it", "pt", "es", "nl", "pl", "ru", "zh", "ja", "ko", "ar", "hi", "tr", "th", "vi", "id"},
		Type:         "embedding",
	},
	"granite-embedding-107m-multilingual-f16": {
		Name:         "granite-embedding-107m-multilingual-f16",
		Description:  "Granite Embedding 107M Multilingual (F16) - Full precision, highest quality",
		URL:          "https://huggingface.co/bartowski/granite-embedding-107m-multilingual-GGUF/resolve/main/granite-embedding-107m-multilingual-f16.gguf",
		Filename:     "granite-embedding-107m-multilingual-f16.gguf",
		Size:         236000000,
		Quantization: "F16",
		Dimensions:   384,
		Languages:    []string{"en", "de", "fr", "it", "pt", "es", "nl", "pl", "ru", "zh", "ja", "ko", "ar", "hi", "tr", "th", "vi", "id"},
		Type:         "embedding",
	},
	"paraphrase-multilingual-minilm-l12-v2-f16": {
		Name:         "paraphrase-multilingual-minilm-l12-v2-f16",
		Description:  "Paraphrase Multilingual MiniLM L12 v2 (F16) - Alternative multilingual model",
		URL:          "https://huggingface.co/mykor/paraphrase-multilingual-MiniLM-L12-v2.gguf/resolve/main/paraphrase-multilingual-MiniLM-L12-118M-v2-F16.gguf",
		Filename:     "paraphrase-multilingual-MiniLM-L12-118M-v2-F16.gguf",
		Size:         242358176,
		Quantization: "F16",
		Dimensions:   384,
		Languages:    []string{"en", "de", "fr", "it", "pt", "es", "nl", "pl", "ru", "zh", "ja", "ko"},
		Type:         "embedding",
	},
	"nomic-embed-text-v1.5-q4_k_m": {
		Name:         "nomic-embed-text-v1.5-q4_k_m",
		Description:  "Nomic Embed Text v1.5 (Q4_K_M) - Popular, high-quality embeddings",
		URL:          "https://huggingface.co/nomic-ai/nomic-embed-text-v1.5-GGUF/resolve/main/nomic-embed-text-v1.5.Q4_K_M.gguf",
		Filename:     "nomic-embed-text-v1.5.Q4_K_M.gguf",
		Size:         550000000,
		Quantization: "Q4_K_M",
		Dimensions:   768,
		Languages:    []string{"en", "de", "fr", "it", "pt", "es", "nl", "pl", "ru", "zh", "ja", "ko", "ar", "hi", "tr", "th", "vi", "id"},
		Type:         "embedding",
	},
}

// GetAvailableEmbeddingModels returns the catalog of available embedding models
func GetAvailableEmbeddingModels() map[string]*EmbeddingModel {
	return embeddingModelsCatalog
}

// GetEmbeddingModel returns a specific embedding model by name
func GetEmbeddingModel(name string) (*EmbeddingModel, bool) {
	model, exists := embeddingModelsCatalog[name]
	return model, exists
}

// GetRecommendedEmbeddingModel returns the recommended embedding model name
func GetRecommendedEmbeddingModel() string {
	return "granite-embedding-107m-multilingual-q6_k_l"
}

// Note: The following functions have been moved to installer.go as methods:
//   - GetEmbeddingModelsDirectory() → installer.GetModelsDirectory()
//   - IsEmbeddingModelDownloaded() → used internally by installer
//   - GetDownloadedEmbeddingModels() → installer.GetDownloadedEmbeddingModels()
//   - DownloadEmbeddingModel() → installer.DownloadEmbeddingModel()
//   - AutoDownloadRecommendedEmbeddingModel() → installer.AutoDownloadRecommendedEmbeddingModel()
//
// This file now only contains:
//   - Model specifications and catalogs (data layer)
//   - Validation helpers (used by installer)

// validateEmbeddingGGUFFile validates a GGUF file by checking its header structure.
// It performs more thorough validation than validateGGUFFile by checking version, tensor count, and metadata count.
func validateEmbeddingGGUFFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// GGUF files start with a magic number: "GGUF" (4 bytes)
	magic := make([]byte, ggufMagicNumberSize)
	if _, err := io.ReadFull(file, magic); err != nil {
		return fmt.Errorf("failed to read magic bytes: %w", err)
	}

	if string(magic) != ggufMagicNumber {
		return fmt.Errorf("invalid GGUF file: wrong magic bytes (expected %q, got %q)", ggufMagicNumber, string(magic))
	}

	// Read version (4 bytes, little-endian uint32)
	var version uint32
	if err := binary.Read(file, binary.LittleEndian, &version); err != nil {
		return fmt.Errorf("failed to read version: %w", err)
	}

	// Check if version is reasonable
	if version == 0 || version > maxGGUFVersion {
		return fmt.Errorf("invalid GGUF file: unreasonable version %d (expected 1-%d)", version, maxGGUFVersion)
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

	if tensorCount > maxTensorCount {
		return fmt.Errorf("invalid GGUF file: unreasonable tensor count %d (max %d)", tensorCount, maxTensorCount)
	}

	if metadataCount > maxMetadataCount {
		return fmt.Errorf("invalid GGUF file: unreasonable metadata count %d (max %d)", metadataCount, maxMetadataCount)
	}

	log.Printf("GGUF file validation passed: version=%d, tensors=%d, metadata_entries=%d",
		version, tensorCount, metadataCount)
	return nil
}
