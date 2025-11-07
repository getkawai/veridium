//go:build !darwin && !linux && !windows

package stablediffusion

import "log"

// selectBestAsset stub for unsupported platforms
func (sdrm *StableDiffusionReleaseManager) selectBestAsset(assets []Asset) *Asset {
	log.Printf("ERROR: Stable Diffusion is not supported on this platform")
	return nil
}

