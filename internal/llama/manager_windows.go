//go:build windows

package llama

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// InstallLlamaCpp attempts to install llama.cpp on Windows using winget or scoop
func (lcm *LlamaCppReleaseManager) InstallLlamaCpp() error {
	if lcm.IsLlamaCppInstalled() {
		log.Println("llama.cpp is already installed")
		return nil
	}

	log.Println("Installing llama.cpp on Windows...")

	// Try winget first (built-in on Windows 10/11)
	if _, err := exec.LookPath("winget"); err == nil {
		log.Println("Trying winget installation...")
		cmd := exec.Command("winget", "install", "--id", "ggerganov.llama.cpp",
			"--silent",
			"--accept-package-agreements",
			"--accept-source-agreements",
			"--disable-interactivity")
		output, err := cmd.CombinedOutput()
		if err == nil {
			log.Println("llama.cpp installed successfully via winget")
			return nil
		}
		log.Printf("winget installation failed: %v\nOutput: %s", err, string(output))
	}

	// Try scoop as fallback (scoop is non-interactive by default)
	if _, err := exec.LookPath("scoop"); err == nil {
		log.Println("Trying scoop installation...")
		cmd := exec.Command("scoop", "install", "llama.cpp", "--no-cache")
		output, err := cmd.CombinedOutput()
		if err == nil {
			log.Println("llama.cpp installed successfully via scoop")
			return nil
		}
		log.Printf("scoop installation failed: %v\nOutput: %s", err, string(output))
	}

	// If both package managers are not available or failed
	return fmt.Errorf("no package manager found (winget/scoop). Will fallback to GitHub download")
}

// WindowsHardware represents detected Windows hardware capabilities
type WindowsHardware struct {
	GPU       string
	IsNVIDIA  bool
	IsAMD     bool
	IsIntel   bool
	HasCUDA   bool
	HasVulkan bool
	HasOpenCL bool
	HasAVX2   bool
}

// detectWindowsHardware detects Windows hardware capabilities
func (lcm *LlamaCppReleaseManager) detectWindowsHardware() *WindowsHardware {
	hardware := &WindowsHardware{}

	// Detect GPU
	hardware.GPU = lcm.detectWindowsGPU()
	hardware.IsNVIDIA = strings.Contains(strings.ToLower(hardware.GPU), "nvidia") || strings.Contains(strings.ToLower(hardware.GPU), "geforce") || strings.Contains(strings.ToLower(hardware.GPU), "rtx") || strings.Contains(strings.ToLower(hardware.GPU), "gtx")
	hardware.IsAMD = strings.Contains(strings.ToLower(hardware.GPU), "amd") || strings.Contains(strings.ToLower(hardware.GPU), "radeon")
	hardware.IsIntel = strings.Contains(strings.ToLower(hardware.GPU), "intel")

	// Detect CUDA (NVIDIA GPUs only)
	if hardware.IsNVIDIA {
		hardware.HasCUDA = lcm.detectCUDA()
	}

	// Detect Vulkan
	hardware.HasVulkan = lcm.detectVulkan()

	// Detect OpenCL
	hardware.HasOpenCL = lcm.detectOpenCL()

	// Detect AVX2 (assume modern CPUs have it)
	hardware.HasAVX2 = lcm.detectAVX2()

	return hardware
}

// detectWindowsGPU detects GPU on Windows
func (lcm *LlamaCppReleaseManager) detectWindowsGPU() string {
	// Try wmic first
	if out, err := exec.Command("wmic", "path", "win32_VideoController", "get", "name", "/value").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Name=") {
				gpu := strings.TrimPrefix(line, "Name=")
				gpu = strings.TrimSpace(gpu)
				if gpu != "" && !strings.Contains(gpu, "Microsoft") {
					return gpu
				}
			}
		}
	}

	// Try PowerShell as fallback
	if out, err := exec.Command("powershell", "-Command", "Get-WmiObject -Class Win32_VideoController | Select-Object -ExpandProperty Name").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" && !strings.Contains(line, "Microsoft") {
				return line
			}
		}
	}

	return "Unknown GPU"
}

// detectCUDA detects CUDA availability
func (lcm *LlamaCppReleaseManager) detectCUDA() bool {
	// Check for nvidia-smi
	if _, err := exec.LookPath("nvidia-smi"); err == nil {
		// Try to run nvidia-smi to verify CUDA is working
		if err := exec.Command("nvidia-smi").Run(); err == nil {
			return true
		}
	}

	// Check for CUDA installation paths
	cudaPaths := []string{
		"C:\\Program Files\\NVIDIA GPU Computing Toolkit\\CUDA",
		"C:\\Program Files (x86)\\NVIDIA GPU Computing Toolkit\\CUDA",
	}

	for _, path := range cudaPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

// detectVulkan detects Vulkan availability
func (lcm *LlamaCppReleaseManager) detectVulkan() bool {
	// Check Windows registry for Vulkan
	if out, err := exec.Command("reg", "query", "HKEY_LOCAL_MACHINE\\SOFTWARE\\Khronos\\Vulkan\\Drivers").Output(); err == nil {
		return len(strings.TrimSpace(string(out))) > 0
	}

	// Check for Vulkan DLLs in system
	vulkanPaths := []string{
		"C:\\Windows\\System32\\vulkan-1.dll",
		"C:\\Windows\\SysWOW64\\vulkan-1.dll",
	}

	for _, path := range vulkanPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

// detectOpenCL detects OpenCL availability
func (lcm *LlamaCppReleaseManager) detectOpenCL() bool {
	// Check for OpenCL DLLs
	openclPaths := []string{
		"C:\\Windows\\System32\\OpenCL.dll",
		"C:\\Windows\\SysWOW64\\OpenCL.dll",
	}

	for _, path := range openclPaths {
		if _, err := os.Stat(path); err == nil {
			return true
		}
	}

	return false
}

// detectAVX2 detects AVX2 CPU instruction support
func (lcm *LlamaCppReleaseManager) detectAVX2() bool {
	// Try to get CPU info via wmic
	if out, err := exec.Command("wmic", "cpu", "get", "name", "/value").Output(); err == nil {
		cpuInfo := strings.ToLower(string(out))
		// Most modern CPUs support AVX2, but we can check for known CPU families
		if strings.Contains(cpuInfo, "intel") || strings.Contains(cpuInfo, "amd") {
			// Assume modern Intel/AMD CPUs have AVX2 (safer assumption)
			return true
		}
	}

	// Default to true for modern systems
	return true
}

// detectHardwareCapabilities detects Windows system hardware capabilities
func (lcm *LlamaCppReleaseManager) detectHardwareCapabilities() *HardwareCapabilities {
	winHw := lcm.detectWindowsHardware()

	caps := &HardwareCapabilities{
		OS:        "windows",
		Arch:      runtime.GOARCH,
		HasNVIDIA: winHw.IsNVIDIA,
		HasAMD:    winHw.IsAMD,
		HasIntel:  winHw.IsIntel,
		HasCUDA:   winHw.HasCUDA,
		HasVulkan: winHw.HasVulkan,
		HasOpenCL: winHw.HasOpenCL,
		HasAVX2:   winHw.HasAVX2,
	}

	log.Printf("Hardware capabilities: OS=%s, Arch=%s, NVIDIA=%v, CUDA=%v, Vulkan=%v",
		caps.OS, caps.Arch, caps.HasNVIDIA, caps.HasCUDA, caps.HasVulkan)

	return caps
}

// getAssetPatterns returns Windows-specific asset patterns in priority order
func (lcm *LlamaCppReleaseManager) getAssetPatterns(hardware *HardwareCapabilities) []string {
	var patterns []string

	// Windows priority: CUDA > Vulkan > OpenCL > AVX2 > Generic
	if hardware.HasNVIDIA && hardware.HasCUDA {
		patterns = append(patterns, "cudart.*win.*cuda.*x64")
	}
	if hardware.HasVulkan {
		patterns = append(patterns, ".*win.*vulkan.*x64")
	}
	if hardware.HasOpenCL {
		patterns = append(patterns, ".*win.*opencl.*x64")
	}
	patterns = append(patterns, ".*win.*avx2.*x64")
	patterns = append(patterns, ".*win.*x64")

	return patterns
}

// detectPlatformSpecs detects hardware specs on Windows
func (specs *HardwareSpecs) detectPlatformSpecs() {
	// Get total memory using wmic
	if out, err := exec.Command("wmic", "computersystem", "get", "TotalPhysicalMemory", "/value").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "TotalPhysicalMemory=") {
				memStr := strings.TrimPrefix(line, "TotalPhysicalMemory=")
				memStr = strings.TrimSpace(memStr)
				if memBytes, err := strconv.ParseInt(memStr, 10, 64); err == nil {
					specs.TotalRAM = memBytes / (1024 * 1024 * 1024) // Convert to GB
				}
			}
		}
	}

	// Estimate available memory
	specs.AvailableRAM = int64(float64(specs.TotalRAM) * 0.8)

	// Get CPU info
	if out, err := exec.Command("wmic", "cpu", "get", "name", "/value").Output(); err == nil {
		lines := strings.Split(string(out), "\n")
		for _, line := range lines {
			if strings.HasPrefix(line, "Name=") {
				specs.CPU = strings.TrimPrefix(line, "Name=")
				specs.CPU = strings.TrimSpace(specs.CPU)
				break
			}
		}
	}
}

// GetBinaryPath returns the path to a specific llama.cpp binary on Windows
// Priority: 1) Local binary path, 2) System PATH
func (lcm *LlamaCppReleaseManager) GetBinaryPath(binaryName string) string {
	// Add .exe extension if not present
	if !strings.HasSuffix(binaryName, ".exe") {
		binaryName += ".exe"
	}

	// First check local binary path
	localPath := filepath.Join(lcm.BinaryPath, binaryName)
	if _, err := os.Stat(localPath); err == nil {
		return localPath
	}

	// Fallback to system PATH
	if systemPath, err := exec.LookPath(binaryName); err == nil {
		return systemPath
	}

	// Return local path as default (even if not exists) for error messages
	return localPath
}

// GetLlamaCLICacheDirectory returns the llama-cli cache directory on Windows
func GetLlamaCLICacheDirectory() string {
	// Use LOCALAPPDATA environment variable
	localAppData := os.Getenv("LOCALAPPDATA")
	if localAppData == "" {
		// Fallback to user home directory
		homeDir, _ := os.UserHomeDir()
		localAppData = filepath.Join(homeDir, "AppData", "Local")
	}
	return filepath.Join(localAppData, "llama.cpp", "cache")
}
