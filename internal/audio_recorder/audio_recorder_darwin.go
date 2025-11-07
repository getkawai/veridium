//go:build darwin

package audio_recorder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// checkAvailableRecordingTools returns list of available recording tools on macOS
func checkAvailableRecordingTools() []string {
	var tools []string

	// Check for sox (preferred)
	if _, err := exec.LookPath("sox"); err == nil {
		tools = append(tools, "sox")
	}

	// Check for ffmpeg in PATH
	if _, err := exec.LookPath("ffmpeg"); err == nil {
		tools = append(tools, "ffmpeg")
	}

	// Check for ffmpeg in ~/.local/bin (auto-downloaded)
	homeDir, err := os.UserHomeDir()
	if err == nil {
		localFFmpeg := filepath.Join(homeDir, ".local", "bin", "ffmpeg")
		if _, err := os.Stat(localFFmpeg); err == nil {
			tools = append(tools, "ffmpeg-local")
		}
	}

	return tools
}

// startPlatformRecording starts recording audio using the specified tool
func startPlatformRecording(tool string, outputPath string) (*exec.Cmd, error) {
	var cmd *exec.Cmd

	switch tool {
	case "sox":
		// Use sox on macOS for recording
		// Format: WAV, 16-bit signed integer, 16kHz, mono (optimal for Whisper)
		cmd = exec.Command("sox", "-d", "-r", "16000", "-c", "1", "-b", "16", "-e", "signed-integer", outputPath)

	case "ffmpeg":
		// Use ffmpeg from PATH
		cmd = exec.Command("ffmpeg", "-f", "avfoundation", "-i", ":0", "-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le", outputPath)

	case "ffmpeg-local":
		// Use ffmpeg from ~/.local/bin
		homeDir, _ := os.UserHomeDir()
		localFFmpeg := filepath.Join(homeDir, ".local", "bin", "ffmpeg")
		cmd = exec.Command(localFFmpeg, "-f", "avfoundation", "-i", ":0", "-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le", outputPath)

	default:
		return nil, fmt.Errorf("unsupported recording tool: %s", tool)
	}

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start recording with %s: %w", tool, err)
	}

	return cmd, nil
}

// stopPlatformRecording stops the recording process gracefully
func stopPlatformRecording(proc *exec.Cmd) error {
	if proc == nil || proc.Process == nil {
		return nil
	}

	// On macOS, send SIGINT (Ctrl+C) for graceful shutdown
	// This allows sox to properly finalize the WAV file
	if err := proc.Process.Signal(os.Interrupt); err != nil {
		fmt.Printf("Failed to send SIGINT: %v, killing process\n", err)
		return proc.Process.Kill()
	}

	return nil
}

// checkPlatformRecordingTool checks if the recording tool is available (deprecated)
// Use checkAvailableRecordingTools() instead
func checkPlatformRecordingTool() (string, error) {
	tools := checkAvailableRecordingTools()
	if len(tools) == 0 {
		return "", fmt.Errorf("no recording tools found")
	}
	return tools[0], nil
}

// installPlatformRecordingTool installs the recording tool
func installPlatformRecordingTool() error {
	log.Println("🔧 sox not found, attempting auto-installation...")

	// Check if Homebrew is installed
	if _, err := exec.LookPath("brew"); err != nil {
		return fmt.Errorf("homebrew is not installed. Please install Homebrew first: https://brew.sh")
	}

	// Install sox
	log.Println("   Installing sox via Homebrew...")
	cmd := exec.Command("brew", "install", "sox")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to install sox: %w\nOutput: %s", err, string(output))
	}

	log.Println("✅ sox installed successfully")
	return nil
}
