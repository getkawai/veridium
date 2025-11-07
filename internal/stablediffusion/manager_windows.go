//go:build windows

package stablediffusion

import "log"

// selectBestAsset selects the best asset for Windows
func (sdrm *StableDiffusionReleaseManager) selectBestAsset(assets []Asset) *Asset {
	var patterns []string

	// Windows releases - prioritize CUDA > Vulkan > AVX2 > AVX > no-AVX
	// Actual formats: cudart-sd-bin-win-cu12-x64.zip, sd-master-fce6afc-bin-win-cuda12-x64.zip, etc.
	if isARM64() {
		patterns = append(patterns, ".*win.*arm64")
		patterns = append(patterns, ".*windows.*arm64")
	} else {
		// Priority order for Windows x64
		patterns = append(patterns, "cudart.*win.*cu12.*x64") // CUDA runtime
		patterns = append(patterns, ".*win.*cuda12.*x64")     // CUDA 12
		patterns = append(patterns, ".*win.*cuda.*x64")       // Any CUDA
		patterns = append(patterns, ".*win.*vulkan.*x64")     // Vulkan
		patterns = append(patterns, ".*win.*avx512.*x64")     // AVX-512
		patterns = append(patterns, ".*win.*avx2.*x64")       // AVX2
		patterns = append(patterns, ".*win.*avx.*x64")        // AVX
		patterns = append(patterns, ".*win.*noavx.*x64")      // No AVX (fallback)
		patterns = append(patterns, ".*win.*x64")             // Generic Windows x64
	}

	log.Printf("Trying %d patterns for Windows asset selection:", len(patterns))
	for i, pattern := range patterns {
		log.Printf("  Pattern %d: %s", i+1, pattern)
	}

	// Try each pattern in priority order
	for _, pattern := range patterns {
		log.Printf("Trying pattern: %s", pattern)
		for _, asset := range assets {
			if sdrm.matchesPattern(asset.Name, pattern) {
				log.Printf("SUCCESS: Matched pattern '%s' with asset '%s'", pattern, asset.Name)
				return &asset
			}
		}
		log.Printf("No match found for pattern: %s", pattern)
	}

	log.Printf("ERROR: No patterns matched any available assets for Windows")
	return nil
}
