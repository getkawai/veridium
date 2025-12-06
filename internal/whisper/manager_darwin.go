//go:build darwin

package whisper

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

// InstallWhisper installs whisper-cpp using Homebrew on macOS
func (m *Manager) InstallWhisper() error {
	if m.IsWhisperInstalled() {
		log.Println("whisper-cpp is already installed")
		return nil
	}

	log.Println("Installing whisper-cpp via Homebrew...")

	// Check if Homebrew is installed
	if _, err := exec.LookPath("brew"); err != nil {
		return fmt.Errorf("Homebrew is not installed. Please install Homebrew first: https://brew.sh")
	}

	// Install whisper-cpp
	cmd := exec.Command("brew", "install", "whisper-cpp")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install whisper-cpp: %w\nOutput: %s", err, string(output))
	}

	log.Println("whisper-cpp installed successfully via Homebrew")
	return nil
}

// InstallFFmpeg installs ffmpeg using Homebrew on macOS
func (m *Manager) InstallFFmpeg() error {
	if m.IsFFmpegInstalled() {
		log.Println("ffmpeg is already installed")
		return nil
	}

	log.Println("🔧 Installing ffmpeg via Homebrew...")

	// Check if Homebrew is installed
	if _, err := exec.LookPath("brew"); err != nil {
		return fmt.Errorf("Homebrew is not installed. Please install Homebrew first: https://brew.sh")
	}

	// Install ffmpeg
	cmd := exec.Command("brew", "install", "ffmpeg")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install ffmpeg: %w\nOutput: %s", err, string(output))
	}

	log.Println("✅ ffmpeg installed successfully via Homebrew")
	return nil
}

// detectAvailableRAM detects available RAM in GB on macOS
func (m *Manager) detectAvailableRAM() int64 {
	// Get total memory using sysctl
	out, err := exec.Command("sysctl", "-n", "hw.memsize").Output()
	if err != nil {
		log.Printf("⚠️  Failed to detect RAM: %v, defaulting to 8GB", err)
		return 8 // Default fallback
	}

	memBytes, err := strconv.ParseInt(strings.TrimSpace(string(out)), 10, 64)
	if err != nil {
		log.Printf("⚠️  Failed to parse RAM size: %v, defaulting to 8GB", err)
		return 8
	}

	totalRAM := memBytes / (1024 * 1024 * 1024) // Convert to GB
	// Estimate available as ~80% of total (conservative)
	availableRAM := int64(float64(totalRAM) * 0.8)

	log.Printf("📊 Detected RAM: %dGB total, ~%dGB available for Whisper model selection", totalRAM, availableRAM)
	return availableRAM
}
