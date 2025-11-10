//go:build darwin

package llama

import (
	"encoding/json"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// detectHardwareCapabilities detects macOS system hardware capabilities
func (lcm *LlamaCppInstaller) detectHardwareCapabilities() *HardwareCapabilities {
	caps := &HardwareCapabilities{
		OS:        "darwin",
		Arch:      runtime.GOARCH,
		HasNVIDIA: false, // Modern Macs don't have NVIDIA GPUs
		HasAMD:    false, // Apple Silicon or Intel integrated
		HasIntel:  false, // Could be detected if needed
		HasCUDA:   false, // macOS doesn't support CUDA
		HasVulkan: false, // macOS uses Metal instead
		HasOpenCL: false, // macOS deprecated OpenCL
		HasAVX2:   true,  // Modern Macs support AVX2 (Intel) or equivalent (Apple Silicon)
	}

	log.Printf("Hardware capabilities: OS=%s, Arch=%s, NVIDIA=%v, CUDA=%v, Vulkan=%v",
		caps.OS, caps.Arch, caps.HasNVIDIA, caps.HasCUDA, caps.HasVulkan)

	return caps
}

// detectPlatformSpecs detects hardware specs on macOS
func (specs *HardwareSpecs) detectPlatformSpecs() {
	// Get total memory
	if out, err := exec.Command("sysctl", "-n", "hw.memsize").Output(); err == nil {
		if memBytes, err := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64); err == nil {
			specs.TotalRAM = memBytes / (1024 * 1024 * 1024) // Convert to GB
		}
	}

	// Estimate available memory (roughly 80% of total, accounting for OS usage)
	specs.AvailableRAM = int64(float64(specs.TotalRAM) * 0.8)

	// Get CPU model
	if out, err := exec.Command("sysctl", "-n", "machdep.cpu.brand_string").Output(); err == nil {
		specs.CPU = strings.TrimSpace(string(out))
	}

	// Try to detect GPU (basic detection)
	if out, err := exec.Command("system_profiler", "SPDisplaysDataType", "-json").Output(); err == nil {
		specs.parseGPUFromSystemProfiler(string(out))
	}
}

// parseGPUFromSystemProfiler parses GPU info from macOS system_profiler JSON output
func (specs *HardwareSpecs) parseGPUFromSystemProfiler(jsonStr string) {
	var data map[string]interface{}
	if err := json.Unmarshal([]byte(jsonStr), &data); err != nil {
		return
	}

	if displays, ok := data["SPDisplaysDataType"].([]interface{}); ok && len(displays) > 0 {
		if display, ok := displays[0].(map[string]interface{}); ok {
			if name, ok := display["sppci_model"].(string); ok {
				specs.GPUModel = name
			}
			// Note: macOS system_profiler doesn't easily provide VRAM info for all GPUs
			// For Apple Silicon, we can make educated guesses based on model
			if strings.Contains(specs.GPUModel, "Apple") {
				// Apple Silicon GPUs share system memory
				specs.GPUMemory = specs.TotalRAM / 4 // Conservative estimate
			}
		}
	}
}

// GetBinaryPath returns the path to a specific llama.cpp binary on macOS
// Priority: 1) Local binary path, 2) Homebrew paths, 3) System PATH
func (lcm *LlamaCppInstaller) GetBinaryPath(binaryName string) string {
	// First check local binary path
	localPath := filepath.Join(lcm.BinaryPath, binaryName)
	if _, err := os.Stat(localPath); err == nil {
		return localPath
	}

	// Check Homebrew paths (Apple Silicon and Intel)
	homebrewPaths := []string{
		"/opt/homebrew/bin/" + binaryName, // Apple Silicon
		"/usr/local/bin/" + binaryName,    // Intel Mac
	}

	for _, path := range homebrewPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Fallback to system PATH
	if systemPath, err := exec.LookPath(binaryName); err == nil {
		return systemPath
	}

	// Return local path as default (even if not exists) for error messages
	return localPath
}

// GetLlamaCLICacheDirectory returns the llama-cli cache directory on macOS
func GetLlamaCLICacheDirectory() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, "Library", "Caches", "llama.cpp")
}

// IsLlamaCppInstalled and VerifyInstalledBinary are now in installer.go (unified)
