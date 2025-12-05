//go:build darwin

package whisper

import (
	"fmt"
	"log"
	"os/exec"
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
