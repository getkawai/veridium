package model

import (
	"path/filepath"

	"github.com/kawai-network/veridium/pkg/tools/models"
)

// ListDownloadedModelsUnified uses the unified model management system
// to list downloaded Whisper models.
// This provides better performance through index caching.
func ListDownloadedModelsUnified(modelsDir string) ([]string, error) {
	// Get base path (modelsDir is typically {base}/models)
	basePath := filepath.Dir(modelsDir)

	// Create unified models manager with provided path
	modelsManager, err := models.NewWithPaths(basePath)
	if err != nil {
		// Fallback to legacy method
		return ListDownloadedModels(modelsDir)
	}

	// Load index
	index := modelsManager.LoadIndex()

	var modelNames []string
	for id, path := range index {
		// Only include validated, downloaded, and audio-type models
		if path.Validated && path.Downloaded && path.Type == models.ModelTypeAudio {
			modelNames = append(modelNames, id)
		}
	}

	// If no models found in index, fallback to legacy scan
	if len(modelNames) == 0 {
		return ListDownloadedModels(modelsDir)
	}

	return modelNames, nil
}

// BuildWhisperIndex builds the model index for Whisper models.
// This should be called during server startup for better performance.
func BuildWhisperIndex(modelsDir string, logger models.Logger) error {
	basePath := filepath.Dir(modelsDir)

	modelsManager, err := models.NewWithPaths(basePath)
	if err != nil {
		return err
	}

	return modelsManager.BuildIndex(logger)
}
