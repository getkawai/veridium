package hardware

import (
	"runtime"

	"github.com/kawai-network/veridium/pkg/xlog"
)

// HardwareSpecs represents system hardware specifications
type HardwareSpecs struct {
	TotalRAM     int64  // Total RAM in GB
	AvailableRAM int64  // Available RAM in GB
	CPU          string // CPU model
	CPUCores     int    // Number of CPU cores
	GPUMemory    int64  // GPU VRAM in GB (if available)
	GPUModel     string // GPU model
}

// DetectHardwareSpecs detects the current system's hardware specifications
func DetectHardwareSpecs() *HardwareSpecs {
	specs := &HardwareSpecs{}
	specs.CPUCores = runtime.NumCPU()

	// Platform-specific detection is implemented in platform files
	specs.detectPlatformSpecs()

	xlog.Info("Detected hardware", "ram_gb", specs.TotalRAM, "available_gb", specs.AvailableRAM, "cpu_cores", specs.CPUCores, "gpu", specs.GPUModel, "vram_gb", specs.GPUMemory)

	return specs
}
