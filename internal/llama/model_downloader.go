package llama

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// QwenModelSpec represents a Qwen model specification for llama-cli download
type QwenModelSpec struct {
	Repo         string // HuggingFace repo (e.g., "Qwen/Qwen2.5-0.5B-Instruct-GGUF")
	Quantization string // Quantization type (Q4_K_M, Q5_K_M, etc.)
	Parameters   string // Parameter size (0.5b, 1.5b, etc.)
	MinRAM       int64  // Minimum RAM required in GB
	Description  string // Model description
}

// GetRecommendedQwenModels returns recommended Qwen models for llama-cli download
func GetRecommendedQwenModels() []QwenModelSpec {
	return []QwenModelSpec{
		{
			Repo:         "Qwen/Qwen2.5-0.5B-Instruct-GGUF",
			Quantization: "Q4_K_M",
			Parameters:   "0.5b",
			MinRAM:       2,
			Description:  "Smallest Qwen model, perfect for low-end hardware and testing",
		},
		{
			Repo:         "Qwen/Qwen2.5-1.5B-Instruct-GGUF",
			Quantization: "Q4_K_M",
			Parameters:   "1.5b",
			MinRAM:       4,
			Description:  "Lightweight Qwen model, good balance of speed and quality",
		},
		{
			Repo:         "Qwen/Qwen2.5-3B-Instruct-GGUF",
			Quantization: "Q4_K_M",
			Parameters:   "3b",
			MinRAM:       6,
			Description:  "Mid-size Qwen model, excellent for most tasks",
		},
		{
			Repo:         "Qwen/Qwen2.5-7B-Instruct-GGUF",
			Quantization: "Q4_K_M",
			Parameters:   "7b",
			MinRAM:       10,
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
	if selectedModel.Repo == "" {
		selectedModel = models[0]
		log.Printf("⚠️  System has low RAM (%dGB), using smallest model", availableRAM)
	}

	log.Printf("📦 Selected model: %s (%s, %s) - requires %dGB RAM (system has %dGB)",
		selectedModel.Repo, selectedModel.Parameters, selectedModel.Quantization,
		selectedModel.MinRAM, availableRAM)

	return selectedModel
}

// DownloadModelWithLlamaCLI downloads a model using llama-cli's built-in HuggingFace integration
func (s *Service) DownloadModelWithLlamaCLI(modelSpec QwenModelSpec) error {
	// Get llama-cli path
	llamaCLI := s.manager.GetBinaryPath("llama-cli")
	if _, err := os.Stat(llamaCLI); err != nil {
		return fmt.Errorf("llama-cli not found at %s: %w", llamaCLI, err)
	}

	modelsDir := s.manager.GetModelsDirectory()

	// Ensure models directory exists
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return fmt.Errorf("failed to create models directory: %w", err)
	}

	// Build HuggingFace repo string (without quantization suffix for URL)
	hfRepo := modelSpec.Repo

	// Generate expected model filename from llama-cli cache
	// Cache location is platform-specific (see manager_*.go files)
	cacheDir := GetLlamaCLICacheDirectory()

	// Build cache filename: Repo_Name_model-quant.gguf
	repoName := strings.ReplaceAll(modelSpec.Repo, "/", "_")
	modelFileName := strings.ToLower(strings.ReplaceAll(filepath.Base(modelSpec.Repo), "-GGUF", ""))
	modelFileName = fmt.Sprintf("%s-%s.gguf", modelFileName, strings.ToLower(modelSpec.Quantization))
	cachedModelPath := filepath.Join(cacheDir, fmt.Sprintf("%s_%s", repoName, modelFileName))

	// Destination in our models directory
	destModelPath := filepath.Join(modelsDir, modelFileName)

	// Check if model already exists in our models directory
	if _, err := os.Stat(destModelPath); err == nil {
		// Validate existing model
		if err := s.validateGGUFFile(destModelPath); err != nil {
			log.Printf("⚠️  Existing model %s failed validation, will re-download: %v", modelFileName, err)
			os.Remove(destModelPath)
		} else {
			log.Printf("✅ Model already exists and is valid: %s", modelFileName)
		return nil
		}
	}

	// Check if model exists in llama-cli cache
	if _, err := os.Stat(cachedModelPath); err == nil {
		log.Printf("📦 Model found in llama-cli cache, validating...")
		// Validate cached model before copying
		if err := s.validateGGUFFile(cachedModelPath); err != nil {
			log.Printf("⚠️  Cached model failed validation, will re-download: %v", err)
			os.Remove(cachedModelPath)
		} else {
			log.Printf("📦 Copying validated model to models directory...")
		if err := s.copyFile(cachedModelPath, destModelPath); err != nil {
			return fmt.Errorf("failed to copy cached model: %w", err)
		}

			// Validate copied file
			if err := s.validateGGUFFile(destModelPath); err != nil {
				os.Remove(destModelPath)
				return fmt.Errorf("copied model failed validation: %w", err)
			}

		fileInfo, _ := os.Stat(destModelPath)
		sizeMB := float64(fileInfo.Size()) / (1024 * 1024)
			log.Printf("✅ Model copied and validated successfully: %s (%.1f MB)", modelFileName, sizeMB)
		return nil
		}
	}

	log.Printf("📥 Downloading model using llama-cli...")
	log.Printf("   Repository: %s", hfRepo)
	log.Printf("   Quantization: %s (auto-selected by llama-cli)", modelSpec.Quantization)
	log.Printf("   This may take several minutes depending on model size and network speed...")

	// Build llama-cli command
	// Let llama-cli download to its cache directory, then we'll copy it
	args := []string{
		"--hf-repo", hfRepo,
		"--prompt", "Hello", // Minimal prompt
		"-n", "1", // Generate only 1 token to exit quickly
		"--log-disable",       // Disable verbose llama.cpp logging
		"--no-display-prompt", // Don't display prompt
	}

	cmd := exec.Command(llamaCLI, args...)

	// Capture output for debugging
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("❌ llama-cli command failed")
		log.Printf("❌ llama-cli output:\n%s", string(output))
		return fmt.Errorf("failed to download model with llama-cli: %w", err)
	}

	// Log output for debugging (even if successful)
	if len(output) > 0 {
		log.Printf("📋 llama-cli output:\n%s", string(output))
	}

	// Verify model was downloaded to cache
	// First try the expected path
	if _, err := os.Stat(cachedModelPath); err != nil {
		// If not found at expected path, search for GGUF files in cache directory
		log.Printf("🔍 Model not found at expected path: %s", cachedModelPath)
		log.Printf("🔍 Searching for GGUF files in cache directory: %s", cacheDir)

		// Search for any GGUF files in cache directory
		entries, readErr := os.ReadDir(cacheDir)
		if readErr != nil {
			log.Printf("⚠️  Failed to read cache directory: %v", readErr)
			return fmt.Errorf("model file not found in cache after download (expected: %s): %w", cachedModelPath, err)
		}

		var foundGGUFFiles []string
		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".gguf") {
				foundPath := filepath.Join(cacheDir, entry.Name())
				foundGGUFFiles = append(foundGGUFFiles, foundPath)
			}
		}

		if len(foundGGUFFiles) == 0 {
			log.Printf("❌ No GGUF files found in cache directory")
			return fmt.Errorf("model file not found in cache after download. Expected: %s. No GGUF files found in cache directory: %s", cachedModelPath, cacheDir)
		}

		// Use the first found GGUF file (or try to match by repo name)
		var selectedFile string
		for _, file := range foundGGUFFiles {
			fileName := strings.ToLower(filepath.Base(file))
			// Try to match by repo name
			if strings.Contains(fileName, strings.ToLower(strings.ReplaceAll(modelSpec.Repo, "/", "_"))) ||
				strings.Contains(fileName, strings.ToLower(modelSpec.Parameters)) {
				selectedFile = file
				break
			}
		}

		// If no match found, use the first file
		if selectedFile == "" {
			selectedFile = foundGGUFFiles[0]
			log.Printf("⚠️  Using first found GGUF file: %s", selectedFile)
		} else {
			log.Printf("✅ Found matching GGUF file: %s", selectedFile)
		}

		cachedModelPath = selectedFile
	} else {
		log.Printf("✅ Model found at expected cache path: %s", cachedModelPath)
	}

	// Validate cached model before copying
	log.Printf("🔍 Validating cached model...")
	if err := s.validateGGUFFile(cachedModelPath); err != nil {
		// Remove corrupt cached file
		os.Remove(cachedModelPath)
		return fmt.Errorf("cached model failed validation: %w", err)
	}

	// Copy from cache to our models directory
	log.Printf("📦 Copying model from cache to models directory...")
	if err := s.copyFile(cachedModelPath, destModelPath); err != nil {
		return fmt.Errorf("failed to copy model from cache: %w", err)
	}

	// Validate copied file
	log.Printf("🔍 Validating copied model...")
	if err := s.validateGGUFFile(destModelPath); err != nil {
		// Remove corrupt copied file
		os.Remove(destModelPath)
		return fmt.Errorf("copied model failed validation: %w", err)
	}

	// Get file size for confirmation
	fileInfo, err := os.Stat(destModelPath)
	if err == nil {
		sizeMB := float64(fileInfo.Size()) / (1024 * 1024)
		log.Printf("✅ Model downloaded and validated successfully: %s (%.1f MB)", modelFileName, sizeMB)
	} else {
		log.Printf("✅ Model downloaded and validated successfully: %s", modelFileName)
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

// AutoDownloadRecommendedModel automatically downloads the best model for the system using llama-cli
func (s *Service) AutoDownloadRecommendedModel() error {
	// Check if any models already exist
	models, err := s.GetAvailableModels()
	if err != nil {
		return fmt.Errorf("failed to check existing models: %w", err)
	}

	if len(models) > 0 {
		log.Printf("✅ Models already available (%d found), skipping auto-download", len(models))
		return nil
	}

	log.Println("📦 No models found, starting auto-download with llama-cli...")

	// Detect hardware specs
	specs := DetectHardwareSpecs()

	// Select optimal model based on available RAM
	modelSpec := SelectOptimalQwenModel(specs.AvailableRAM)

	// Download the model using llama-cli
	if err := s.DownloadModelWithLlamaCLI(modelSpec); err != nil {
		return fmt.Errorf("failed to download model: %w", err)
	}

	log.Println("🎉 Model download completed successfully!")
	return nil
}
