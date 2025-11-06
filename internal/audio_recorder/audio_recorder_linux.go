//go:build linux

package audio_recorder

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// startPlatformRecording starts recording audio using platform-specific tools
func startPlatformRecording(outputPath string) (*exec.Cmd, error) {
	// Try multiple tools in order of preference

	// 1. Try arecord (ALSA - most common on Linux)
	if _, err := exec.LookPath("arecord"); err == nil {
		cmd := exec.Command("arecord", "-f", "S16_LE", "-r", "16000", "-c", "1", "-t", "wav", outputPath)
		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("failed to start recording with arecord: %w", err)
		}
		return cmd, nil
	}

	// 2. Try sox (powerful alternative)
	if _, err := exec.LookPath("sox"); err == nil {
		cmd := exec.Command("sox", "-d", "-r", "16000", "-c", "1", "-b", "16", "-e", "signed-integer", outputPath)
		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("failed to start recording with sox: %w", err)
		}
		return cmd, nil
	}

	// 3. Try ffmpeg (universal tool) - check system PATH
	if _, err := exec.LookPath("ffmpeg"); err == nil {
		cmd := exec.Command("ffmpeg", "-f", "alsa", "-i", "default", "-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le", outputPath)
		if err := cmd.Start(); err != nil {
			return nil, fmt.Errorf("failed to start recording with ffmpeg: %w", err)
		}
		return cmd, nil
	}

	// 4. Try ffmpeg in ~/.local/bin (auto-downloaded)
	homeDir, err := os.UserHomeDir()
	if err == nil {
		localFFmpeg := filepath.Join(homeDir, ".local", "bin", "ffmpeg")
		if _, err := os.Stat(localFFmpeg); err == nil {
			cmd := exec.Command(localFFmpeg, "-f", "alsa", "-i", "default", "-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le", outputPath)
			if err := cmd.Start(); err != nil {
				return nil, fmt.Errorf("failed to start recording with ffmpeg: %w", err)
			}
			return cmd, nil
		}
	}

	return nil, fmt.Errorf("no recording tool found. Please install one of: arecord, sox, or ffmpeg")
}

// stopPlatformRecording stops the recording process gracefully
func stopPlatformRecording(proc *exec.Cmd) error {
	if proc == nil || proc.Process == nil {
		return nil
	}

	// On Linux, send SIGINT (Ctrl+C) for graceful shutdown
	// This allows arecord to properly finalize the WAV file
	if err := proc.Process.Signal(os.Interrupt); err != nil {
		fmt.Printf("Failed to send SIGINT: %v, killing process\n", err)
		return proc.Process.Kill()
	}

	return nil
}

// checkPlatformRecordingTool checks if any recording tool is available
func checkPlatformRecordingTool() (string, error) {
	// Check for available tools in order of preference

	// 1. Check arecord (ALSA - most common)
	if _, err := exec.LookPath("arecord"); err == nil {
		return "arecord", nil
	}

	// 2. Check sox
	if _, err := exec.LookPath("sox"); err == nil {
		return "sox", nil
	}

	// 3. Check ffmpeg in PATH
	if _, err := exec.LookPath("ffmpeg"); err == nil {
		return "ffmpeg", nil
	}

	// 4. Check ffmpeg in ~/.local/bin (auto-downloaded)
	homeDir, err := os.UserHomeDir()
	if err == nil {
		localFFmpeg := filepath.Join(homeDir, ".local", "bin", "ffmpeg")
		if _, err := os.Stat(localFFmpeg); err == nil {
			return "ffmpeg (local)", nil
		}
	}

	return "", fmt.Errorf("no recording tool found")
}

// installPlatformRecordingTool attempts to auto-download ffmpeg static binary
// Falls back to manual installation instructions if download fails
func installPlatformRecordingTool() error {
	log.Println("🔧 No recording tool found")
	log.Println("   Attempting to download ffmpeg static binary...")

	// Try to download static ffmpeg binary
	if err := downloadStaticFFmpeg(); err != nil {
		log.Printf("⚠️  Auto-download failed: %v", err)
		log.Println("")
		log.Println("   Please install one of these tools manually:")

		// Detect package manager and provide appropriate instructions
		if _, err := exec.LookPath("apt-get"); err == nil {
			log.Println("   $ sudo apt-get install -y alsa-utils    # Recommended (arecord)")
			log.Println("   $ sudo apt-get install -y sox           # Alternative")
			log.Println("   $ sudo apt-get install -y ffmpeg        # Alternative")
			return fmt.Errorf("no recording tool installed. Run: sudo apt-get install -y alsa-utils")
		}

		if _, err := exec.LookPath("yum"); err == nil {
			log.Println("   $ sudo yum install -y alsa-utils        # Recommended (arecord)")
			log.Println("   $ sudo yum install -y sox               # Alternative")
			log.Println("   $ sudo yum install -y ffmpeg            # Alternative")
			return fmt.Errorf("no recording tool installed. Run: sudo yum install -y alsa-utils")
		}

		if _, err := exec.LookPath("dnf"); err == nil {
			log.Println("   $ sudo dnf install -y alsa-utils        # Recommended (arecord)")
			log.Println("   $ sudo dnf install -y sox               # Alternative")
			log.Println("   $ sudo dnf install -y ffmpeg            # Alternative")
			return fmt.Errorf("no recording tool installed. Run: sudo dnf install -y alsa-utils")
		}

		if _, err := exec.LookPath("pacman"); err == nil {
			log.Println("   $ sudo pacman -S alsa-utils             # Recommended (arecord)")
			log.Println("   $ sudo pacman -S sox                    # Alternative")
			log.Println("   $ sudo pacman -S ffmpeg                 # Alternative")
			return fmt.Errorf("no recording tool installed. Run: sudo pacman -S alsa-utils")
		}

		if _, err := exec.LookPath("zypper"); err == nil {
			log.Println("   $ sudo zypper install alsa-utils        # Recommended (arecord)")
			log.Println("   $ sudo zypper install sox               # Alternative")
			log.Println("   $ sudo zypper install ffmpeg            # Alternative")
			return fmt.Errorf("no recording tool installed. Run: sudo zypper install alsa-utils")
		}

		// Generic instructions
		log.Println("   - alsa-utils (for arecord)")
		log.Println("   - sox")
		log.Println("   - ffmpeg")
		return fmt.Errorf("no recording tool installed. Please install alsa-utils, sox, or ffmpeg")
	}

	log.Println("✅ ffmpeg static binary downloaded successfully")
	return nil
}

// downloadStaticFFmpeg downloads a static ffmpeg binary
func downloadStaticFFmpeg() error {
	// Determine architecture
	arch := runtime.GOARCH
	var downloadURL string

	// Use johnvansickle's static ffmpeg builds (widely trusted)
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
		return fmt.Errorf("unsupported architecture: %s", arch)
	}

	// Create directory for local binaries
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	binDir := filepath.Join(homeDir, ".local", "bin")
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return fmt.Errorf("failed to create bin directory: %w", err)
	}

	ffmpegPath := filepath.Join(binDir, "ffmpeg")

	// Check if already downloaded
	if _, err := os.Stat(ffmpegPath); err == nil {
		log.Println("   ffmpeg binary already exists, skipping download")
		return nil
	}

	log.Printf("   Downloading ffmpeg for %s...", arch)
	log.Println("   This may take a few minutes (~50 MB)...")

	// Download the archive
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("failed to download ffmpeg: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Save to temp file
	tempFile := filepath.Join(os.TempDir(), "ffmpeg-static.tar.xz")
	out, err := os.Create(tempFile)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile)

	_, err = io.Copy(out, resp.Body)
	out.Close()
	if err != nil {
		return fmt.Errorf("failed to save download: %w", err)
	}

	// Extract ffmpeg binary
	log.Println("   Extracting ffmpeg...")
	extractDir := filepath.Join(os.TempDir(), "ffmpeg-extract")
	os.RemoveAll(extractDir)
	os.MkdirAll(extractDir, 0755)
	defer os.RemoveAll(extractDir)

	// Extract using tar command
	cmd := exec.Command("tar", "-xf", tempFile, "-C", extractDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract archive: %w", err)
	}

	// Find ffmpeg binary in extracted directory
	var ffmpegBinary string
	err = filepath.Walk(extractDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == "ffmpeg" {
			ffmpegBinary = path
			return filepath.SkipDir
		}
		return nil
	})

	if err != nil || ffmpegBinary == "" {
		return fmt.Errorf("failed to find ffmpeg binary in archive")
	}

	// Copy to bin directory
	if err := copyFile(ffmpegBinary, ffmpegPath); err != nil {
		return fmt.Errorf("failed to copy ffmpeg: %w", err)
	}

	// Make executable
	if err := os.Chmod(ffmpegPath, 0755); err != nil {
		return fmt.Errorf("failed to make ffmpeg executable: %w", err)
	}

	log.Printf("   ffmpeg installed to: %s", ffmpegPath)
	log.Println("   Note: You may need to add ~/.local/bin to your PATH")

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
