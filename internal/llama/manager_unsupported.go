//go:build !darwin && !linux && !windows

package llama

import (
	"fmt"
	"runtime"
)

// InstallLlamaCpp returns an error on unsupported platforms
func (lcm *LlamaCppReleaseManager) InstallLlamaCpp() error {
	return fmt.Errorf("llama.cpp installation not supported on this platform")
}

// detectHardwareCapabilities stub for unsupported platforms
func (lcm *LlamaCppReleaseManager) detectHardwareCapabilities() *HardwareCapabilities {
	return &HardwareCapabilities{
		OS:        runtime.GOOS,
		Arch:      runtime.GOARCH,
		HasNVIDIA: false,
		HasAMD:    false,
		HasIntel:  false,
		HasCUDA:   false,
		HasVulkan: false,
		HasOpenCL: false,
		HasAVX2:   false,
	}
}

// getAssetPatterns stub for unsupported platforms
func (lcm *LlamaCppReleaseManager) getAssetPatterns(hardware *HardwareCapabilities) []string {
	// Return empty patterns for unsupported platforms
	return []string{}
}

// detectPlatformSpecs stub for unsupported platforms
func (specs *HardwareSpecs) detectPlatformSpecs() {
	// Fallback values for unsupported platforms
	specs.TotalRAM = 8
	specs.AvailableRAM = 6
	specs.CPU = "Unknown CPU"
	specs.GPUModel = "Unknown GPU"
	specs.GPUMemory = 0
}
