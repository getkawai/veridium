package image

import (
	"github.com/kawai-network/veridium/pkg/stablediffusion/models"
)

// ModelSpec is a backward compatibility wrapper for pkg/stablediffusion/models.ModelSpec
type ModelSpec = models.ModelSpec

// HardwareSpecs is a backward compatibility wrapper for pkg/stablediffusion/models.HardwareSpecs
type HardwareSpecs = models.HardwareSpecs

// GetAvailableModels returns all available Stable Diffusion models
// Backward compatibility wrapper for pkg/stablediffusion/models.GetAvailableModels
func GetAvailableModels() []ModelSpec {
	return models.GetAvailableModels()
}

// SelectOptimalModel selects the best Stable Diffusion model based on hardware specs
// Backward compatibility wrapper for pkg/stablediffusion/models.SelectOptimalModel
func SelectOptimalModel(specs *HardwareSpecs) ModelSpec {
	return models.SelectOptimalModel(specs)
}
