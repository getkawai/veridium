//go:build windows

package whisper

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// InstallWhisper provides instructions for Windows installation
func (m *Manager) InstallWhisper() error {
	if m.IsWhisperInstalled() {
		log.Println("whisper-cpp is already installed")
		return nil
	}

	return fmt.Errorf("please install whisper-cpp manually on Windows:\n" +
		"1. Download pre-built binaries from: https://github.com/ggml-org/whisper.cpp/releases\n" +
		"2. Extract and add to PATH\n" +
		"Or build from source: https://github.com/ggml-org/whisper.cpp")
}

// InstallFFmpeg installs ffmpeg on Windows using winget or chocolatey
func (m *Manager) InstallFFmpeg() error {
	if m.IsFFmpegInstalled() {
		log.Println("ffmpeg is already installed")
		return nil
	}

	log.Println("🔧 ffmpeg not found, attempting auto-installation...")

	// Try winget first (Windows 10/11)
	if _, err := exec.LookPath("winget"); err == nil {
		log.Println("   Installing ffmpeg via winget...")
		cmd := exec.Command("winget", "install", "ffmpeg", "--accept-package-agreements", "--accept-source-agreements")
		output, err := cmd.CombinedOutput()
		if err == nil {
			log.Println("✅ ffmpeg installed successfully via winget")
			return nil
		}
		log.Printf("   winget installation failed: %v\nOutput: %s", err, string(output))
	}

	// Try chocolatey
	if _, err := exec.LookPath("choco"); err == nil {
		log.Println("   Installing ffmpeg via chocolatey...")
		cmd := exec.Command("choco", "install", "ffmpeg", "-y")
		output, err := cmd.CombinedOutput()
		if err == nil {
			log.Println("✅ ffmpeg installed successfully via chocolatey")
			return nil
		}
		log.Printf("   chocolatey installation failed: %v\nOutput: %s", err, string(output))
	}

	return fmt.Errorf("no supported package manager found (winget or chocolatey). Please install ffmpeg manually from: https://ffmpeg.org/download.html")
}

// isFFmpegInCommonPaths checks common Windows installation paths for ffmpeg
func (m *Manager) isFFmpegInCommonPaths() string {
	// Common installation paths on Windows
	commonPaths := []string{
		`C:\ffmpeg\bin\ffmpeg.exe`,
		`C:\Program Files\ffmpeg\bin\ffmpeg.exe`,
		`C:\Program Files (x86)\ffmpeg\bin\ffmpeg.exe`,
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path
		}
	}

	// Check user's local directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		localFFmpeg := filepath.Join(homeDir, ".local", "bin", "ffmpeg.exe")
		if _, err := os.Stat(localFFmpeg); err == nil {
			return localFFmpeg
		}
	}

	return ""
}
