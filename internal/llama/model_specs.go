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
)

// QwenModelSpec represents a Qwen model specification for direct download
type QwenModelSpec struct {
	Name          string // Model name (e.g., "qwen2.5-0.5b-instruct-q4_k_m")
	URL           string // Direct download URL
	Quantization  string // Quantization type (Q4_K_M, Q5_K_M, etc.)
	Parameters    string // Parameter size (0.5b, 1.5b, etc.)
	MinRAM        int64  // Minimum RAM required in GB
	Size          int64  // Expected file size in bytes
	SHA256        string // Expected SHA256 checksum (optional)
	Description   string // Model description
	ProjectorURL  string // URL for the multimodal projector file (for VL models)
	ProjectorSize int64  // Expected size of the projector file
}

// Note: Service struct has been removed. All functions are now pure functions
// that accept modelsDir as a parameter or use GetModelsDirectory() directly.

// GetRecommendedVLModels returns recommended Qwen models for direct download
// Models are ordered from smallest to largest by MinRAM requirement
func GetRecommendedVLModels() []QwenModelSpec {
	return []QwenModelSpec{
		{
			Name:          "qwen3-vl-4b-instruct-q4_k_m",
			URL:           "https://huggingface.co/bartowski/Qwen_Qwen3-VL-4B-Instruct-GGUF/resolve/main/Qwen_Qwen3-VL-4B-Instruct-Q4_K_M.gguf",
			Quantization:  "Q4_K_M",
			Parameters:    "4b",
			MinRAM:        6,
			Size:          2500000000, // ~2.5 GB
			Description:   "Qwen3-VL 4B - Multimodal vision-language model, perfect for low-end hardware",
			ProjectorURL:  "https://huggingface.co/bartowski/Qwen_Qwen3-VL-4B-Instruct-GGUF/resolve/main/mmproj-Qwen_Qwen3-VL-4B-Instruct-f16.gguf",
			ProjectorSize: 600000000, // ~600 MB
		},
		{
			Name:          "qwen3-vl-8b-instruct-q4_k_m",
			URL:           "https://huggingface.co/bartowski/Qwen_Qwen3-VL-8B-Instruct-GGUF/resolve/main/Qwen_Qwen3-VL-8B-Instruct-Q4_K_M.gguf",
			Quantization:  "Q4_K_M",
			Parameters:    "8b",
			MinRAM:        12,
			Size:          5500000000, // ~5.5 GB
			Description:   "Qwen3-VL 8B - Advanced multimodal model with excellent vision capabilities",
			ProjectorURL:  "https://huggingface.co/bartowski/Qwen_Qwen3-VL-8B-Instruct-GGUF/resolve/main/mmproj-Qwen_Qwen3-VL-8B-Instruct-f16.gguf",
			ProjectorSize: 1159029920, // ~1.16 GB
		},
		{
			Name:          "qwen3-vl-32b-instruct-q4_k_m",
			URL:           "https://huggingface.co/bartowski/Qwen_Qwen3-VL-32B-Instruct-GGUF/resolve/main/Qwen_Qwen3-VL-32B-Instruct-Q4_K_M.gguf",
			Quantization:  "Q4_K_M",
			Parameters:    "32b",
			MinRAM:        24,
			Size:          22000000000, // ~22 GB
			Description:   "Qwen3-VL 32B - High-quality multimodal model for advanced vision-language tasks",
			ProjectorURL:  "https://huggingface.co/bartowski/Qwen_Qwen3-VL-32B-Instruct-GGUF/resolve/main/mmproj-Qwen_Qwen3-VL-32B-Instruct-f16.gguf",
			ProjectorSize: 1200000000, // ~1.2 GB
		},
	}
}

// GetRecommendedModels returns recommended non-reasoning text models for direct download
// These models do NOT generate <think> tags and are suitable for general chat
// Models are ordered from smallest to largest by MinRAM requirement
// Using Llama 3.2 series - proven, stable, no reasoning overhead
func GetRecommendedModels() []QwenModelSpec {
	return []QwenModelSpec{
		// Llama 3.2 1B - Smallest, fastest
		{
			Name:         "llama-3.2-1b-instruct-q4_k_m",
			URL:          "https://huggingface.co/bartowski/Llama-3.2-1B-Instruct-GGUF/resolve/main/Llama-3.2-1B-Instruct-Q4_K_M.gguf",
			Quantization: "Q4_K_M",
			Parameters:   "1b",
			MinRAM:       2,
			Size:         697000000, // ~697 MB
			Description:  "Llama 3.2 1B - Ultra fast, no reasoning tags",
		},
		// Llama 3.2 3B - Good balance
		{
			Name:         "llama-3.2-3b-instruct-q4_k_m",
			URL:          "https://huggingface.co/bartowski/Llama-3.2-3B-Instruct-GGUF/resolve/main/Llama-3.2-3B-Instruct-Q4_K_M.gguf",
			Quantization: "Q4_K_M",
			Parameters:   "3b",
			MinRAM:       4,
			Size:         2019000000, // ~2.0 GB
			Description:  "Llama 3.2 3B - Fast, no reasoning tags, great for general chat",
		},
		// Llama 3.2 3B Q8 - Higher quality
		{
			Name:         "llama-3.2-3b-instruct-q8_0",
			URL:          "https://huggingface.co/bartowski/Llama-3.2-3B-Instruct-GGUF/resolve/main/Llama-3.2-3B-Instruct-Q8_0.gguf",
			Quantization: "Q8_0",
			Parameters:   "3b",
			MinRAM:       6,
			Size:         3420000000, // ~3.4 GB
			Description:  "Llama 3.2 3B Q8 - Better quality, no reasoning tags",
		},
	}
}

// GetRecommendedReasoningModels returns Qwen3 models that support reasoning (<think> tags)
// Use these only when reasoning mode is explicitly enabled
func GetRecommendedReasoningModels() []QwenModelSpec {
	return []QwenModelSpec{
		{
			Name:         "qwen3-4b-instruct-q4_k_m",
			URL:          "https://huggingface.co/bartowski/Qwen_Qwen3-4B-GGUF/resolve/main/Qwen_Qwen3-4B-Q4_K_M.gguf",
			Quantization: "Q4_K_M",
			Parameters:   "4b",
			MinRAM:       6,
			Size:         2497280960, // ~2.5 GB
			Description:  "Qwen3 4B - Reasoning model with <think> tags",
		},
		{
			Name:         "qwen3-8b-instruct-q4_k_m",
			URL:          "https://huggingface.co/bartowski/Qwen_Qwen3-8B-GGUF/resolve/main/Qwen_Qwen3-8B-Q4_K_M.gguf",
			Quantization: "Q4_K_M",
			Parameters:   "8b",
			MinRAM:       12,
			Size:         5027784224, // ~5.0 GB
			Description:  "Qwen3 8B - Reasoning model with <think> tags",
		},
		{
			Name:         "qwen3-14b-instruct-q4_k_m",
			URL:          "https://huggingface.co/bartowski/Qwen_Qwen3-14B-GGUF/resolve/main/Qwen_Qwen3-14B-Q4_K_M.gguf",
			Quantization: "Q4_K_M",
			Parameters:   "14b",
			MinRAM:       16,
			Size:         9001753632, // ~9.0 GB
			Description:  "Qwen3 14B - Reasoning model with <think> tags",
		},
		{
			Name:         "qwen3-32b-instruct-q4_k_m",
			URL:          "https://huggingface.co/bartowski/Qwen_Qwen3-32B-GGUF/resolve/main/Qwen_Qwen3-32B-Q4_K_M.gguf",
			Quantization: "Q4_K_M",
			Parameters:   "32b",
			MinRAM:       24,
			Size:         19762149696, // ~19.8 GB
			Description:  "Qwen3 32B - Reasoning model with <think> tags",
		},
	}
}

// SelectOptimalQwenModel selects the best VL model based on available RAM.
// It returns the largest model that fits within the available RAM.
// If no model fits, it returns the smallest model as a fallback.
// Note: Currently excludes 7B model as it requires multi-file download support
func SelectOptimalQwenModel(availableRAM int64) QwenModelSpec {
	models := GetRecommendedVLModels()

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

	log.Printf("📦 Selected VL model: %s (%s, %s) - requires %dGB RAM (system has %dGB)",
		selectedModel.Name, selectedModel.Parameters, selectedModel.Quantization,
		selectedModel.MinRAM, availableRAM)

	return selectedModel
}

// SelectOptimalQwenTextModel selects the best text model based on available RAM.
// It returns the largest model that fits within the available RAM.
// If no model fits, it returns the smallest model as a fallback.
func SelectOptimalQwenTextModel(availableRAM int64) QwenModelSpec {
	models := GetRecommendedModels()

	// Select the largest model that fits in available RAM
	var selectedModel QwenModelSpec
	for _, model := range models {
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

	log.Printf("📦 Selected text model: %s (%s, %s) - requires %dGB RAM (system has %dGB)",
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
// Utility/Lightweight Model Functions (for Summary, Title, etc.)
// ============================================================================

// GetRecommendedUtilityModels returns recommended small models for utility tasks
// These are optimized for:
// - Summary generation (compress old messages)
// - Title generation (create conversation titles)
// - Quick text processing tasks
// All models are non-reasoning (no <think> tags) and lightweight (<1GB)
func GetRecommendedUtilityModels() []QwenModelSpec {
	return []QwenModelSpec{
		{
			Name:         "llama-3.2-1b-instruct-q4_k_m",
			URL:          "https://huggingface.co/bartowski/Llama-3.2-1B-Instruct-GGUF/resolve/main/Llama-3.2-1B-Instruct-Q4_K_M.gguf",
			Quantization: "Q4_K_M",
			Parameters:   "1b",
			MinRAM:       2,
			Size:         697000000, // ~697 MB
			Description:  "Llama 3.2 1B - BEST for summary/title (fast, no <think> tags, good quality)",
		},
		{
			Name:         "llama-3.2-1b-instruct-q5_k_m",
			URL:          "https://huggingface.co/bartowski/Llama-3.2-1B-Instruct-GGUF/resolve/main/Llama-3.2-1B-Instruct-Q5_K_M.gguf",
			Quantization: "Q5_K_M",
			Parameters:   "1b",
			MinRAM:       2,
			Size:         810000000, // ~810 MB
			Description:  "Llama 3.2 1B Q5 - Better quality than Q4, slightly slower",
		},
		{
			Name:         "llama-3.2-3b-instruct-q4_k_m",
			URL:          "https://huggingface.co/bartowski/Llama-3.2-3B-Instruct-GGUF/resolve/main/Llama-3.2-3B-Instruct-Q4_K_M.gguf",
			Quantization: "Q4_K_M",
			Parameters:   "3b",
			MinRAM:       4,
			Size:         1900000000, // ~1.9 GB
			Description:  "Llama 3.2 3B - Higher quality, good for systems with more RAM",
		},
		{
			Name:         "mistral-7b-instruct-q4_k_m",
			URL:          "https://huggingface.co/bartowski/Mistral-7B-Instruct-v0.3-GGUF/resolve/main/Mistral-7B-Instruct-v0.3-Q4_K_M.gguf",
			Quantization: "Q4_K_M",
			Parameters:   "7b",
			MinRAM:       8,
			Size:         4370000000, // ~4.37 GB
			Description:  "Mistral 7B - Alternative, excellent quality (but larger)",
		},
	}
}

// SelectOptimalUtilityModel selects the best small model for utility tasks
// Prefers: Llama 3.2 1B (fastest, no reasoning overhead) > Llama 3B > Mistral
// Returns the smallest model that fits in RAM, or nil if none available
func SelectOptimalUtilityModel(availableRAM int64) *QwenModelSpec {
	models := GetRecommendedUtilityModels()

	// Select the smallest model that fits (utility tasks don't need large models)
	// Priority: 1B > 3B > 7B
	for _, model := range models {
		if model.MinRAM <= availableRAM {
			log.Printf("📦 Selected utility model: %s (%s, %s) - requires %dGB RAM (system has %dGB)",
				model.Name, model.Parameters, model.Quantization,
				model.MinRAM, availableRAM)
			return &model
		}
	}

	// If even smallest model doesn't fit, return nil
	log.Printf("⚠️  No utility model fits in %dGB RAM (minimum required: %dGB)",
		availableRAM, models[0].MinRAM)
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

// Embedding model catalog - contains verified embedding models with correct URLs
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
