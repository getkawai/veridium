//go:build !darwin && !linux && !windows

package image

import "log"

// selectBestAsset stub for unsupported platforms
func (sdrm *StableDiffusion) selectBestAsset(assets []Asset) *Asset {
	log.Printf("ERROR: Stable Diffusion is not supported on this platform")
	return nil
}
