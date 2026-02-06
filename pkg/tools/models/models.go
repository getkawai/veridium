// Package models provides support for tooling around model management.
package models

import (
	"context"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	"github.com/kawai-network/veridium/pkg/kronk/model"
	"github.com/kawai-network/veridium/pkg/tools/defaults"
	"go.yaml.in/yaml/v2"
)

var (
	localFolder = "models"
	indexFile   = ".index.yaml"
)

// ModelType represents the type of AI model
type ModelType string

const (
	ModelTypeLLM       ModelType = "llm"
	ModelTypeDiffusion ModelType = "diffusion"
	ModelTypeAudio     ModelType = "audio"
)

// Models manages the model system.
type Models struct {
	modelsPath string
	biMutex    sync.Mutex
}

// New constructs the models system using defaults paths.
func New() (*Models, error) {
	return NewWithPaths("")
}

// NewWithPaths constructs the models system, If the basePath is empty, the
// default location is used.
func NewWithPaths(basePath string) (*Models, error) {
	basePath = defaults.BaseDir(basePath)

	modelPath := filepath.Join(basePath, localFolder)

	if err := os.MkdirAll(modelPath, 0755); err != nil {
		return nil, fmt.Errorf("creating models directory: %w", err)
	}

	m := Models{
		modelsPath: modelPath,
	}

	return &m, nil
}

// NewLLMModels constructs a model manager for LLM models.
// Deprecated: Use New() instead. All models now use unified flat structure.
func NewLLMModels(basePath string) (*Models, error) {
	return NewWithPaths(basePath)
}

// NewDiffusionModels constructs a model manager for Stable Diffusion models.
// Deprecated: Use New() instead. All models now use unified flat structure.
func NewDiffusionModels(basePath string) (*Models, error) {
	return NewWithPaths(basePath)
}

// NewAudioModels constructs a model manager for audio models (Whisper, etc).
// Deprecated: Use New() instead. All models now use unified flat structure.
func NewAudioModels(basePath string) (*Models, error) {
	return NewWithPaths(basePath)
}

// Path returns the location of the models path.
func (m *Models) Path() string {
	return m.modelsPath
}

// LoadIndex loads the model index from disk.
// This is now exported for use by other packages.
func (m *Models) LoadIndex() map[string]Path {
	return m.loadIndex()
}

// BuildIndex builds the model index for fast model access.
func (m *Models) BuildIndex(log Logger) error {
	currentIndex := m.loadIndex()

	m.biMutex.Lock()
	defer m.biMutex.Unlock()

	if err := m.removeEmptyDirs(); err != nil {
		return fmt.Errorf("remove-empty-dirs: %w", err)
	}

	entries, err := os.ReadDir(m.modelsPath)
	if err != nil {
		return fmt.Errorf("list-models: reading models directory: %w", err)
	}

	index := make(map[string]Path)

	for _, orgEntry := range entries {
		if !orgEntry.IsDir() {
			continue
		}

		org := orgEntry.Name()

		modelEntries, err := os.ReadDir(fmt.Sprintf("%s/%s", m.modelsPath, org))
		if err != nil {
			continue
		}

		for _, modelEntry := range modelEntries {
			if !modelEntry.IsDir() {
				continue
			}

			modelFamily := modelEntry.Name()

			fileEntries, err := os.ReadDir(fmt.Sprintf("%s/%s/%s", m.modelsPath, org, modelFamily))
			if err != nil {
				continue
			}

			modelfiles := make(map[string][]string)
			projFiles := make(map[string]string)

			for _, fileEntry := range fileEntries {
				if fileEntry.IsDir() {
					continue
				}

				name := fileEntry.Name()

				if name == ".DS_Store" {
					continue
				}

				if strings.HasPrefix(name, "mmproj") {
					modelID := extractModelID(name[7:])
					projFiles[modelID] = filepath.Join(m.modelsPath, org, modelFamily, fileEntry.Name())
					continue
				}

				modelID := extractModelID(fileEntry.Name())
				filePath := filepath.Join(m.modelsPath, org, modelFamily, fileEntry.Name())
				modelfiles[modelID] = append(modelfiles[modelID], filePath)
			}

			ctx := context.Background()

			for modelID, files := range modelfiles {
				modelIDLower := strings.ToLower(modelID)
				isValidated := currentIndex[modelIDLower].Validated
				existingType := currentIndex[modelIDLower].Type
				modelValid := true

				log(ctx, "checking model", "modelID", modelID, "isValidated", isValidated)

				slices.Sort(files)

				mp := Path{
					ModelFiles: files,
					Downloaded: true,
					Type:       existingType, // Preserve existing type if available
				}

				// Only detect type if not already set
				if mp.Type == "" {
					mp.Type = detectModelType(modelID, files)
				}

				if projFile, exists := projFiles[modelID]; exists {
					mp.ProjFile = projFile
				}

				if !isValidated {
					for _, file := range files {
						log(ctx, "running check ", "model", path.Base(file))
						if err := model.CheckModel(file, true); err != nil {
							log(ctx, "running check ", "model", path.Base(file), "ERROR", err)
							modelValid = false
						}
					}

					if mp.ProjFile != "" {
						log(ctx, "running check ", "proj", path.Base(mp.ProjFile))
						if err := model.CheckModel(mp.ProjFile, true); err != nil {
							log(ctx, "running check ", "proj", path.Base(mp.ProjFile), "ERROR", err)
							modelValid = false
						}
					}
				} else {
					modelValid = isValidated
				}

				mp.Validated = modelValid

				index[modelIDLower] = mp
			}
		}
	}

	indexData, err := yaml.Marshal(&index)
	if err != nil {
		return fmt.Errorf("marshal index: %w", err)
	}

	indexPath := filepath.Join(m.modelsPath, indexFile)
	if err := os.WriteFile(indexPath, indexData, 0644); err != nil {
		return fmt.Errorf("write index file: %w", err)
	}

	return nil
}

// =============================================================================

func (m *Models) removeEmptyDirs() error {
	var dirs []string

	err := filepath.WalkDir(m.modelsPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() && path != m.modelsPath {
			dirs = append(dirs, path)
		}

		return nil
	})

	if err != nil {
		return fmt.Errorf("walking directory tree: %w", err)
	}

	for i := len(dirs) - 1; i >= 0; i-- {
		entries, err := os.ReadDir(dirs[i])
		if err != nil {
			continue
		}

		if isDirEffectivelyEmpty(entries) {
			// Remove any .DS_Store before removing directory
			dsStore := filepath.Join(dirs[i], ".DS_Store")
			os.Remove(dsStore)
			os.Remove(dirs[i])
		}
	}

	return nil
}

// isDirEffectivelyEmpty returns true if directory only contains ignorable files like .DS_Store
func isDirEffectivelyEmpty(entries []os.DirEntry) bool {
	for _, e := range entries {
		if e.Name() != ".DS_Store" {
			return false
		}
	}

	return true
}

// detectModelType attempts to determine the model type based on model ID and file patterns
func detectModelType(modelID string, files []string) ModelType {
	modelIDLower := strings.ToLower(modelID)

	// Check for Whisper-specific patterns (must be exact match, not just "audio")
	if strings.Contains(modelIDLower, "whisper") {
		return ModelTypeAudio
	}

	// Check for Stable Diffusion patterns in model ID
	if strings.Contains(modelIDLower, "stable-diffusion") ||
		strings.Contains(modelIDLower, "sd-v1") ||
		strings.Contains(modelIDLower, "sd-v2") ||
		strings.Contains(modelIDLower, "sdxl") ||
		strings.Contains(modelIDLower, "sd-turbo") ||
		strings.Contains(modelIDLower, "flux") {
		return ModelTypeDiffusion
	}

	// Check file extensions and patterns
	hasGGUF := false
	hasSafetensors := false
	hasCkpt := false

	for _, file := range files {
		fileLower := strings.ToLower(file)

		// Whisper uses .bin files with specific naming
		if strings.HasSuffix(fileLower, ".bin") && strings.Contains(fileLower, "ggml") {
			return ModelTypeAudio
		}

		// Track file types
		if strings.HasSuffix(fileLower, ".gguf") {
			hasGGUF = true
		}
		if strings.HasSuffix(fileLower, ".safetensors") {
			hasSafetensors = true
		}
		if strings.HasSuffix(fileLower, ".ckpt") {
			hasCkpt = true
		}
	}

	// .ckpt files are primarily used for Stable Diffusion
	if hasCkpt {
		return ModelTypeDiffusion
	}

	// .gguf files are primarily used for LLMs (quantized models)
	if hasGGUF && !hasSafetensors && !hasCkpt {
		return ModelTypeLLM
	}

	// .safetensors files are primarily used for Stable Diffusion models
	// Only classify as LLM if there are clear LLM indicators
	if hasSafetensors {
		// Check for LLM-specific patterns
		for _, file := range files {
			fileLower := strings.ToLower(file)
			// Common LLM model name patterns
			if strings.Contains(fileLower, "llama") ||
				strings.Contains(fileLower, "mistral") ||
				strings.Contains(fileLower, "qwen") ||
				strings.Contains(fileLower, "gemma") ||
				strings.Contains(fileLower, "phi") ||
				strings.Contains(fileLower, "gpt") {
				return ModelTypeLLM
			}
		}
		// Default safetensors to diffusion (most common use case)
		return ModelTypeDiffusion
	}

	// Default to LLM for other cases
	return ModelTypeLLM
}

// NormalizeHuggingFaceDownloadURL converts short format to full HuggingFace download URLs.
// Input:  mradermacher/Qwen2-Audio-7B-GGUF/Qwen2-Audio-7B.Q8_0.gguf
// Output: https://huggingface.co/mradermacher/Qwen2-Audio-7B-GGUF/resolve/main/Qwen2-Audio-7B.Q8_0.gguf
func NormalizeHuggingFaceDownloadURL(url string) string {
	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		return url
	}

	parts := strings.Split(url, "/")
	if len(parts) >= 3 {
		org := parts[0]
		repo := parts[1]
		filename := strings.Join(parts[2:], "/")
		return fmt.Sprintf("https://huggingface.co/%s/%s/resolve/main/%s", org, repo, filename)
	}

	return url
}
