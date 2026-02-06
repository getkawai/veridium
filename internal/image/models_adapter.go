package image

import (
	"strings"

	"github.com/kawai-network/veridium/pkg/tools/models"
)

// CheckInstalledModelsUnified uses the unified model management system
// to check for installed Stable Diffusion models.
// This is a wrapper that provides backward compatibility.
func (sdrm *StableDiffusion) CheckInstalledModelsUnified() ([]string, error) {
	// Create unified models manager (flat structure)
	modelsManager, err := models.New()
	if err != nil {
		// Fallback to legacy method
		return sdrm.CheckInstalledModels()
	}

	// Load index
	index := modelsManager.LoadIndex()

	var modelIDs []string
	for id, path := range index {
		// Only include validated, downloaded, and diffusion-type models
		if path.Validated && path.Downloaded && path.Type == models.ModelTypeDiffusion {
			modelIDs = append(modelIDs, id)
		}
	}

	// If no models found in index, fallback to legacy scan
	if len(modelIDs) == 0 {
		return sdrm.CheckInstalledModels()
	}

	return modelIDs, nil
}

// hasStableDiffusionModelUnified checks if any SD model exists using unified system
func (sdrm *StableDiffusion) hasStableDiffusionModelUnified(installedModels []string) bool {
	for _, model := range installedModels {
		modelLower := strings.ToLower(model)
		if strings.Contains(modelLower, "stable-diffusion") ||
			strings.Contains(modelLower, "sd-v1") ||
			strings.Contains(modelLower, "sd-v2") ||
			strings.Contains(modelLower, "sdxl") ||
			strings.Contains(modelLower, "sd-turbo") ||
			strings.Contains(modelLower, "flux") {
			return true
		}
	}
	return false
}
