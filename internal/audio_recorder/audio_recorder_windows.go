//go:build windows

package audio_recorder

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// checkAvailableRecordingTools returns list of available recording tools on Windows
func checkAvailableRecordingTools() []string {
	var tools []string

	// Check for ffmpeg in PATH
	if _, err := exec.LookPath("ffmpeg"); err == nil {
		tools = append(tools, "ffmpeg")
	}

	// Check for ffmpeg in common installation paths
	commonPaths := []string{
		`C:\ffmpeg\bin\ffmpeg.exe`,
		`C:\Program Files\ffmpeg\bin\ffmpeg.exe`,
		`C:\Program Files (x86)\ffmpeg\bin\ffmpeg.exe`,
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			tools = append(tools, "ffmpeg-"+path)
			break // Only add one custom path
		}
	}

	// Check for ffmpeg in user's local directory
	homeDir, err := os.UserHomeDir()
	if err == nil {
		localFFmpeg := filepath.Join(homeDir, ".local", "bin", "ffmpeg.exe")
		if _, err := os.Stat(localFFmpeg); err == nil {
			tools = append(tools, "ffmpeg-local")
		}
	}

	return tools
}

// startPlatformRecording starts recording audio using the specified tool
func startPlatformRecording(tool string, outputPath string) (*exec.Cmd, error) {
	var cmd *exec.Cmd

	switch {
	case tool == "ffmpeg":
		// ffmpeg from PATH
		cmd = exec.Command("ffmpeg", "-f", "dshow", "-i", "audio=", "-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le", outputPath)

	case tool == "ffmpeg-local":
		// ffmpeg from user's local directory
		homeDir, _ := os.UserHomeDir()
		localFFmpeg := filepath.Join(homeDir, ".local", "bin", "ffmpeg.exe")
		cmd = exec.Command(localFFmpeg, "-f", "dshow", "-i", "audio=", "-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le", outputPath)

	case len(tool) > 7 && tool[:7] == "ffmpeg-":
		// ffmpeg from custom path
		ffmpegPath := tool[7:] // Remove "ffmpeg-" prefix
		cmd = exec.Command(ffmpegPath, "-f", "dshow", "-i", "audio=", "-ar", "16000", "-ac", "1", "-c:a", "pcm_s16le", outputPath)

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

	// On Windows, kill the process directly
	return proc.Process.Kill()
}

// checkPlatformRecordingTool checks if the recording tool is available (deprecated)
// Use checkAvailableRecordingTools() instead
func checkPlatformRecordingTool() (string, error) {
	tools := checkAvailableRecordingTools()
	if len(tools) == 0 {
		return "", fmt.Errorf("ffmpeg not found")
	}
	return tools[0], nil
}

// installPlatformRecordingTool installs the recording tool
func installPlatformRecordingTool() error {
	log.Println("🔧 ffmpeg not found, attempting auto-installation...")

	// Check if winget is available (Windows 10+)
	if _, err := exec.LookPath("winget"); err == nil {
		log.Println("   Installing ffmpeg via winget...")
		cmd := exec.Command("winget", "install", "ffmpeg", "--accept-package-agreements", "--accept-source-agreements")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install ffmpeg: %w\nOutput: %s", err, string(output))
		}
		log.Println("✅ ffmpeg installed successfully")
		return nil
	}

	// Check if chocolatey is available
	if _, err := exec.LookPath("choco"); err == nil {
		log.Println("   Installing ffmpeg via chocolatey...")
		cmd := exec.Command("choco", "install", "ffmpeg", "-y")
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("failed to install ffmpeg: %w\nOutput: %s", err, string(output))
		}
		log.Println("✅ ffmpeg installed successfully")
		return nil
	}

	return fmt.Errorf("no supported package manager found (winget or chocolatey). Please install ffmpeg manually from: https://ffmpeg.org/download.html")
}
