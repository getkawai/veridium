//go:build linux

package llama

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

// InstallLlamaCpp attempts to install llama.cpp on Linux
// Note: Package managers on Linux require sudo which needs user interaction
// So we skip package manager and go directly to GitHub download
func (lcm *LlamaCppReleaseManager) InstallLlamaCpp() error {
	if lcm.IsLlamaCppInstalled() {
		log.Println("llama.cpp is already installed")
		return nil
	}

	// Package managers on Linux require sudo (user interaction)
	// Return error to trigger GitHub download fallback
	return fmt.Errorf("package manager installation requires sudo. Will use GitHub download")
}

// hasNVIDIAGPU checks if system has NVIDIA GPU on Linux
func (lcm *LlamaCppReleaseManager) hasNVIDIAGPU() bool {
	// Check lspci for NVIDIA
	if out, err := exec.Command("lspci").Output(); err == nil {
		output := strings.ToLower(string(out))
		if strings.Contains(output, "nvidia") {
			return true
		}
	}

	// Check /proc/driver/nvidia
	if _, err := os.Stat("/proc/driver/nvidia"); err == nil {
		return true
	}

	return false
}

// detectLinuxCUDA detects CUDA availability on Linux
func (lcm *LlamaCppReleaseManager) detectLinuxCUDA() bool {
	// Check for nvidia-smi
	if _, err := exec.LookPath("nvidia-smi"); err == nil {
		if err := exec.Command("nvidia-smi").Run(); err == nil {
			return true
		}
	}

	// Check for CUDA libraries
	cudaPaths := []string{
		"/usr/local/cuda",
		"/opt/cuda",
		"/usr/lib/x86_64-linux-gnu/libcuda.so",
		"/usr/lib64/libcuda.so",
	}

	for _, path := range cudaPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

// detectHardwareCapabilities detects Linux system hardware capabilities
func (lcm *LlamaCppReleaseManager) detectHardwareCapabilities() *HardwareCapabilities {
	caps := &HardwareCapabilities{
		OS:        "linux",
		Arch:      runtime.GOARCH,
		HasNVIDIA: lcm.hasNVIDIAGPU(),
		HasAMD:    false, // Could be detected via lspci if needed
		HasIntel:  false, // Could be detected via lspci if needed
		HasCUDA:   lcm.detectLinuxCUDA(),
		HasVulkan: false, // Could be detected if needed
		HasOpenCL: false, // Could be detected if needed
		HasAVX2:   true,  // Assume modern Linux systems
	}

	log.Printf("Hardware capabilities: OS=%s, Arch=%s, NVIDIA=%v, CUDA=%v, Vulkan=%v",
		caps.OS, caps.Arch, caps.HasNVIDIA, caps.HasCUDA, caps.HasVulkan)

	return caps
}

// getAssetPatterns returns Linux-specific asset patterns in priority order
func (lcm *LlamaCppReleaseManager) getAssetPatterns(hardware *HardwareCapabilities) []string {
	var patterns []string

	// Linux priority: CUDA > Ubuntu > Generic
	if hardware.HasNVIDIA && hardware.HasCUDA {
		patterns = append(patterns, "cudart.*linux.*cuda.*x64")
	}
	patterns = append(patterns, ".*ubuntu.*x64")
	patterns = append(patterns, ".*linux.*x64")

	return patterns
}

// detectPlatformSpecs detects hardware specs on Linux
func (specs *HardwareSpecs) detectPlatformSpecs() {
	// Get memory info from /proc/meminfo
	if out, err := exec.Command("cat", "/proc/meminfo").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "MemTotal:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					if memKB, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
						specs.TotalRAM = memKB / (1024 * 1024) // Convert KB to GB
					}
				}
			} else if strings.HasPrefix(line, "MemAvailable:") {
				fields := strings.Fields(line)
				if len(fields) >= 2 {
					if memKB, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
						specs.AvailableRAM = memKB / (1024 * 1024) // Convert KB to GB
					}
				}
			}
		}
	}

	// If MemAvailable not found, estimate as 80% of total
	if specs.AvailableRAM == 0 {
		specs.AvailableRAM = int64(float64(specs.TotalRAM) * 0.8)
	}

	// Get CPU info
	if out, err := exec.Command("cat", "/proc/cpuinfo").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "model name") {
				parts := strings.Split(line, ":")
				if len(parts) >= 2 {
					specs.CPU = strings.TrimSpace(parts[1])
					break
				}
			}
		}
	}

	// Try to detect NVIDIA GPU
	if out, err := exec.Command("nvidia-smi", "--query-gpu=name,memory.total", "--format=csv,noheader,nounits").Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		if len(lines) > 0 {
			parts := strings.Split(lines[0], ",")
			if len(parts) >= 2 {
				specs.GPUModel = strings.TrimSpace(parts[0])
				if vramMB, err := strconv.ParseInt(strings.TrimSpace(parts[1]), 10, 64); err == nil {
					specs.GPUMemory = vramMB / 1024 // Convert MB to GB
				}
			}
		}
	}
}
