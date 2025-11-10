package llama

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hybridgroup/yzma/pkg/download"
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

// GetRecommendedQwenModels returns recommended Qwen models for direct download
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

// SelectOptimalQwenModel selects the best model based on available RAM
func SelectOptimalQwenModel(availableRAM int64) QwenModelSpec {
	models := GetRecommendedQwenModels()

	// Select the largest model that fits in available RAM
	var selectedModel QwenModelSpec
	for _, model := range models {
		if model.MinRAM <= availableRAM {
			selectedModel = model
		} else {
			break
		}
	}

	// If no model fits, use the smallest one
	if selectedModel.Name == "" {
		selectedModel = models[0]
		log.Printf("⚠️  System has low RAM (%dGB), using smallest model", availableRAM)
	}

	log.Printf("📦 Selected model: %s (%s, %s) - requires %dGB RAM (system has %dGB)",
		selectedModel.Name, selectedModel.Parameters, selectedModel.Quantization,
		selectedModel.MinRAM, availableRAM)

	return selectedModel
}

// DownloadModel downloads a model directly from HuggingFace using yzma/pkg/download
func (s *Service) DownloadModel(modelSpec QwenModelSpec) error {
	modelsDir := s.manager.GetModelsDirectory()

	// Ensure models directory exists
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return fmt.Errorf("failed to create models directory: %w", err)
	}

	// Build model filename
	modelFileName := fmt.Sprintf("%s.gguf", modelSpec.Name)
	destModelPath := filepath.Join(modelsDir, modelFileName)
	tempModelPath := destModelPath + ".tmp"

	// Clean up any stale temporary files from previous interrupted downloads
	if _, err := os.Stat(tempModelPath); err == nil {
		log.Printf("🧹 Cleaning up stale temporary file from previous download: %s", filepath.Base(tempModelPath))
		os.Remove(tempModelPath)
	}

	// Check if model already exists
	if _, err := os.Stat(destModelPath); err == nil {
		log.Printf("✅ Model already exists: %s", modelFileName)

		// Verify existing model integrity if checksum is provided
		if modelSpec.SHA256 != "" {
			if err := s.verifyModelChecksum(destModelPath, modelSpec.SHA256); err != nil {
				log.Printf("⚠️  Existing model checksum invalid, re-downloading...")
				os.Remove(destModelPath)
			} else {
				return nil
			}
		} else {
			return nil
		}
	}

	log.Printf("📥 Downloading model: %s", modelSpec.Name)
	log.Printf("   URL: %s", modelSpec.URL)
	log.Printf("   Expected size: %.1f MB", float64(modelSpec.Size)/(1024*1024))
	log.Printf("   This may take several minutes depending on network speed...")

	// Download to temporary file first
	if err := download.GetModel(modelSpec.URL, tempModelPath); err != nil {
		// Clean up failed download
		os.Remove(tempModelPath)
		return fmt.Errorf("failed to download model: %w", err)
	}

	// Verify downloaded file
	fileInfo, err := os.Stat(tempModelPath)
	if err != nil {
		os.Remove(tempModelPath)
		return fmt.Errorf("failed to stat downloaded file: %w", err)
	}

	// Check file size (allow 5% variance for metadata)
	actualSize := fileInfo.Size()
	expectedSize := modelSpec.Size
	sizeVariance := float64(actualSize) / float64(expectedSize)

	if sizeVariance < 0.95 || sizeVariance > 1.05 {
		os.Remove(tempModelPath)
		return fmt.Errorf("downloaded file size mismatch: got %d bytes, expected ~%d bytes (%.1f%% of expected)",
			actualSize, expectedSize, sizeVariance*100)
	}

	// Verify checksum if provided
	if modelSpec.SHA256 != "" {
		log.Printf("🔒 Verifying model integrity...")
		if err := s.verifyModelChecksum(tempModelPath, modelSpec.SHA256); err != nil {
			os.Remove(tempModelPath)
			return fmt.Errorf("model integrity check failed: %w", err)
		}
		log.Printf("✅ Model integrity verified")
	}

	// Verify it's a valid GGUF file
	if err := s.validateGGUFFile(tempModelPath); err != nil {
		os.Remove(tempModelPath)
		return fmt.Errorf("invalid GGUF file: %w", err)
	}

	// Move temporary file to final destination
	if err := os.Rename(tempModelPath, destModelPath); err != nil {
		os.Remove(tempModelPath)
		return fmt.Errorf("failed to move downloaded file: %w", err)
	}

	sizeMB := float64(actualSize) / (1024 * 1024)
	log.Printf("✅ Model downloaded successfully: %s (%.1f MB)", modelFileName, sizeMB)

	return nil
}

// verifyModelChecksum verifies the SHA256 checksum of a file
func (s *Service) verifyModelChecksum(filePath, expectedChecksum string) error {
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

// validateGGUFFile performs basic validation on a GGUF file
func (s *Service) validateGGUFFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Read GGUF magic number (first 4 bytes should be "GGUF")
	magic := make([]byte, 4)
	if _, err := io.ReadFull(file, magic); err != nil {
		return fmt.Errorf("failed to read file header: %w", err)
	}

	if string(magic) != "GGUF" {
		return fmt.Errorf("invalid GGUF magic number: got %q, expected \"GGUF\"", string(magic))
	}

	return nil
}

// copyFile copies a file from src to dst
func (s *Service) copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}

// CleanupStaleTempFiles removes all stale temporary download files
// This should be called on service startup to clean up interrupted downloads
func (s *Service) CleanupStaleTempFiles() error {
	modelsDir := s.manager.GetModelsDirectory()

	entries, err := os.ReadDir(modelsDir)
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
			tmpPath := filepath.Join(modelsDir, entry.Name())

			// Get file info to check age (optional: only delete if older than X)
			info, err := entry.Info()
			if err != nil {
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

// AutoDownloadRecommendedModel automatically downloads the best model for the system
func (s *Service) AutoDownloadRecommendedModel() error {
	// Clean up any stale temp files from previous interrupted downloads
	if err := s.CleanupStaleTempFiles(); err != nil {
		log.Printf("⚠️  Failed to cleanup stale temp files: %v", err)
		// Don't fail, just log
	}
	// Check if any models already exist
	models, err := s.GetAvailableModels()
	if err != nil {
		return fmt.Errorf("failed to check existing models: %w", err)
	}

	if len(models) > 0 {
		log.Printf("✅ Models already available (%d found), skipping auto-download", len(models))
		return nil
	}

	log.Println("📦 No models found, starting auto-download...")

	// Detect hardware specs
	specs := DetectHardwareSpecs()

	// Select optimal model based on available RAM
	modelSpec := SelectOptimalQwenModel(specs.AvailableRAM)

	// Download the model
	if err := s.DownloadModel(modelSpec); err != nil {
		return fmt.Errorf("failed to download model: %w", err)
	}

	log.Println("🎉 Model download completed successfully!")
	return nil
}
