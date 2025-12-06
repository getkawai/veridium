//go:build linux

package whisper

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

// InstallWhisper provides instructions for Linux installation
func (m *Manager) InstallWhisper() error {
	if m.IsWhisperInstalled() {
		log.Println("whisper-cpp is already installed")
		return nil
	}

	return fmt.Errorf("please install whisper-cpp manually on Linux:\n" +
		"1. Install Homebrew for Linux: /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"\n" +
		"2. Run: brew install whisper-cpp\n" +
		"Or build from source: https://github.com/ggml-org/whisper.cpp")
}

// InstallFFmpeg downloads and installs ffmpeg static binary on Linux
func (m *Manager) InstallFFmpeg() error {
	if m.IsFFmpegInstalled() {
		log.Println("ffmpeg is already installed")
		return nil
	}

	log.Println("🔧 ffmpeg not found, attempting auto-installation...")
	log.Println("   Downloading ffmpeg static binary...")

	// Determine architecture
	arch := runtime.GOARCH
	var downloadURL string

	switch arch {
	case "amd64":
		downloadURL = "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-amd64-static.tar.xz"
	case "arm64":
		downloadURL = "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-arm64-static.tar.xz"
	case "arm":
		downloadURL = "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-armhf-static.tar.xz"
	case "386":
		downloadURL = "https://johnvansickle.com/ffmpeg/releases/ffmpeg-release-i686-static.tar.xz"
	default:
		return fmt.Errorf("unsupported architecture for auto-download: %s. Please install ffmpeg manually", arch)
	}

	// Create ~/.local/bin directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	binDir := filepath.Join(homeDir, ".local", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	ffmpegPath := filepath.Join(binDir, "ffmpeg")

	// Check if already exists
	if _, err := os.Stat(ffmpegPath); err == nil {
		log.Println("   ffmpeg binary already exists, skipping download")
		return nil
	}

	log.Printf("   Downloading ffmpeg for %s...", arch)

	// Download the archive
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download ffmpeg: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download ffmpeg: HTTP %d", resp.StatusCode)
	}

	// Save to temp file
	tempFile := filepath.Join(os.TempDir(), "ffmpeg-static.tar.xz")
	out, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		return fmt.Errorf("failed to save ffmpeg archive: %w", err)
	}
	defer os.Remove(tempFile)

	// Extract ffmpeg binary
	log.Println("   Extracting ffmpeg...")
	extractDir := filepath.Join(os.TempDir(), "ffmpeg-extract")
	os.RemoveAll(extractDir)
	os.MkdirAll(extractDir, 0755)
	defer os.RemoveAll(extractDir)

	// Use tar to extract
	cmd := exec.Command("tar", "-xf", tempFile, "-C", extractDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract ffmpeg: %w", err)
	}

	// Find ffmpeg binary in extracted directory
	var ffmpegBinary string
	filepath.Walk(extractDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.IsDir() && info.Name() == "ffmpeg" {
			ffmpegBinary = path
			return filepath.SkipDir
		}
		return nil
	})

	if ffmpegBinary == "" {
		return fmt.Errorf("failed to find ffmpeg binary in archive")
	}

	// Copy to destination
	if err := copyFile(ffmpegBinary, ffmpegPath); err != nil {
		return fmt.Errorf("failed to copy ffmpeg: %w", err)
	}

	// Make executable
	if err := os.Chmod(ffmpegPath, 0755); err != nil {
		return fmt.Errorf("failed to make ffmpeg executable: %w", err)
	}

	log.Printf("✅ ffmpeg installed to: %s", ffmpegPath)
	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}

// detectAvailableRAM detects available RAM in GB on Linux
func (m *Manager) detectAvailableRAM() int64 {
	// Read from /proc/meminfo
	out, err := exec.Command("cat", "/proc/meminfo").Output()
	if err != nil {
		log.Printf("⚠️  Failed to detect RAM: %v, defaulting to 8GB", err)
		return 8
	}

	var totalRAM, availableRAM int64
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "MemTotal:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if memKB, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
					totalRAM = memKB / (1024 * 1024) // Convert KB to GB
				}
			}
		} else if strings.HasPrefix(line, "MemAvailable:") {
			fields := strings.Fields(line)
			if len(fields) >= 2 {
				if memKB, err := strconv.ParseInt(fields[1], 10, 64); err == nil {
					availableRAM = memKB / (1024 * 1024) // Convert KB to GB
				}
			}
		}
	}

	// If MemAvailable not found, estimate as 80% of total
	if availableRAM == 0 && totalRAM > 0 {
		availableRAM = int64(float64(totalRAM) * 0.8)
	}

	if availableRAM == 0 {
		log.Printf("⚠️  Could not parse RAM info, defaulting to 8GB")
		return 8
	}

	log.Printf("📊 Detected RAM: %dGB total, %dGB available for Whisper model selection", totalRAM, availableRAM)
	return availableRAM
}
