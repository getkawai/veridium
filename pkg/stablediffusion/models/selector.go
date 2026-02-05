package models

// SelectOptimalModel selects the best Stable Diffusion model based on hardware specs
// Inspired by kronk's hardware-aware model selection pattern
func SelectOptimalModel(specs *HardwareSpecs) ModelSpec {
	models := GetAvailableModels()

	var selectedModel ModelSpec
	var bestScore int64

	for _, model := range models {
		ramOk := model.MinRAM <= specs.AvailableRAM
		vramOk := specs.GPUMemory == 0 || model.MinVRAM <= specs.GPUMemory

		if ramOk && vramOk {
			// Calculate score: prefer larger models that fit
			// Score = MinRAM + MinVRAM (higher is better quality)
			score := model.MinRAM + model.MinVRAM

			if score > bestScore {
				bestScore = score
				selectedModel = model
			}
		}
	}

	// Fallback to smallest model if none fit
	if selectedModel.Name == "" {
		selectedModel = models[0]
	}

	return selectedModel
}
