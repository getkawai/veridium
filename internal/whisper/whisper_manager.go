package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// WhisperManager handles whisper-cpp installation and usage
type WhisperManager struct {
	ModelsDir string // Directory where whisper models are stored
}

// NewWhisperManager creates a new whisper manager
func NewWhisperManager() *WhisperManager {
	homeDir, _ := os.UserHomeDir()
	modelsDir := filepath.Join(homeDir, ".kawai-agent", "whisper-models")

	// Ensure models directory exists
	os.MkdirAll(modelsDir, 0755)

	return &WhisperManager{
		ModelsDir: modelsDir,
	}
}

// IsWhisperInstalled checks if whisper-cpp is available
func (wm *WhisperManager) IsWhisperInstalled() bool {
	_, err := exec.LookPath("whisper-cpp")
	return err == nil
}

// InstallWhisper installs whisper-cpp using the appropriate package manager
func (wm *WhisperManager) InstallWhisper() error {
	if wm.IsWhisperInstalled() {
		log.Println("whisper-cpp is already installed")
		return nil
	}

	switch runtime.GOOS {
	case "darwin":
		return wm.installWithHomebrew()
	case "linux":
		return wm.installOnLinux()
	case "windows":
		return wm.installOnWindows()
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

// installWithHomebrew installs whisper-cpp using Homebrew on macOS
func (wm *WhisperManager) installWithHomebrew() error {
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

// installOnLinux provides instructions for Linux installation
func (wm *WhisperManager) installOnLinux() error {
	return fmt.Errorf("please install whisper-cpp manually on Linux:\n" +
		"1. Install Homebrew for Linux: /bin/bash -c \"$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\"\n" +
		"2. Run: brew install whisper-cpp\n" +
		"Or build from source: https://github.com/ggml-org/whisper.cpp")
}

// installOnWindows provides instructions for Windows installation
func (wm *WhisperManager) installOnWindows() error {
	return fmt.Errorf("please install whisper-cpp manually on Windows:\n" +
		"1. Download pre-built binaries from: https://github.com/ggml-org/whisper.cpp/releases\n" +
		"2. Extract and add to PATH\n" +
		"Or build from source: https://github.com/ggml-org/whisper.cpp")
}

// GetAvailableModels returns a list of available whisper models
func (wm *WhisperManager) GetAvailableModels() []string {
	return []string{
		"tiny",      // ~39 MB, fastest
		"tiny.en",   // ~39 MB, English only
		"base",      // ~142 MB
		"base.en",   // ~142 MB, English only
		"small",     // ~466 MB
		"small.en",  // ~466 MB, English only
		"medium",    // ~1.5 GB
		"medium.en", // ~1.5 GB, English only
		"large-v1",  // ~2.9 GB
		"large-v2",  // ~2.9 GB
		"large-v3",  // ~2.9 GB, latest and most accurate
	}
}

// DownloadModel downloads a whisper model if it doesn't exist
func (wm *WhisperManager) DownloadModel(modelName string) error {
	modelPath := filepath.Join(wm.ModelsDir, fmt.Sprintf("ggml-%s.bin", modelName))

	// Check if model already exists
	if _, err := os.Stat(modelPath); err == nil {
		log.Printf("Model %s already exists", modelName)
		return nil
	}

	log.Printf("Downloading whisper model: %s", modelName)

	// Download using curl
	url := fmt.Sprintf("https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-%s.bin", modelName)
	cmd := exec.Command("curl", "-L", "-o", modelPath, url)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to download model %s: %w\nOutput: %s", modelName, err, string(output))
	}

	log.Printf("Successfully downloaded model: %s", modelName)
	return nil
}

// GetModelPath returns the path to a specific model
func (wm *WhisperManager) GetModelPath(modelName string) string {
	return filepath.Join(wm.ModelsDir, fmt.Sprintf("ggml-%s.bin", modelName))
}

// IsModelDownloaded checks if a model is downloaded
func (wm *WhisperManager) IsModelDownloaded(modelName string) bool {
	modelPath := wm.GetModelPath(modelName)
	_, err := os.Stat(modelPath)
	return err == nil
}

// GetInstalledModels returns a list of downloaded models
func (wm *WhisperManager) GetInstalledModels() ([]string, error) {
	files, err := os.ReadDir(wm.ModelsDir)
	if err != nil {
		return nil, err
	}

	var models []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "ggml-") && strings.HasSuffix(file.Name(), ".bin") {
			// Extract model name from filename
			name := strings.TrimPrefix(file.Name(), "ggml-")
			name = strings.TrimSuffix(name, ".bin")
			models = append(models, name)
		}
	}

	return models, nil
}

// TranscribeFile transcribes an audio file using whisper-cpp
func (wm *WhisperManager) TranscribeFile(audioPath, modelName string, options map[string]interface{}) (string, error) {
	if !wm.IsWhisperInstalled() {
		return "", fmt.Errorf("whisper-cpp is not installed")
	}

	// Check if model exists
	if !wm.IsModelDownloaded(modelName) {
		return "", fmt.Errorf("model %s is not downloaded", modelName)
	}

	// Check if audio file exists
	if _, err := os.Stat(audioPath); err != nil {
		return "", fmt.Errorf("audio file not found: %w", err)
	}

	modelPath := wm.GetModelPath(modelName)

	// Build command arguments
	args := []string{
		"-m", modelPath,
		"-f", audioPath,
		"--output-txt", // Output as text
	}

	// Add optional parameters
	if language, ok := options["language"].(string); ok && language != "" {
		args = append(args, "-l", language)
	}
	if threads, ok := options["threads"].(int); ok && threads > 0 {
		args = append(args, "-t", fmt.Sprintf("%d", threads))
	}
	if translate, ok := options["translate"].(bool); ok && translate {
		args = append(args, "--translate")
	}

	log.Printf("Running whisper-cpp with args: %v", args)

	// Execute whisper-cpp
	cmd := exec.Command("whisper-cpp", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("transcription failed: %w\nOutput: %s", err, string(output))
	}

	// whisper-cpp outputs to a .txt file, read it
	txtFile := strings.TrimSuffix(audioPath, filepath.Ext(audioPath)) + ".txt"
	if _, err := os.Stat(txtFile); err == nil {
		content, err := os.ReadFile(txtFile)
		if err == nil {
			// Clean up the generated txt file
			os.Remove(txtFile)
			return strings.TrimSpace(string(content)), nil
		}
	}

	// Fallback: parse output directly
	return wm.parseWhisperOutput(string(output)), nil
}

// parseWhisperOutput extracts transcription from whisper-cpp output
func (wm *WhisperManager) parseWhisperOutput(output string) string {
	lines := strings.Split(output, "\n")
	var transcription strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and lines that look like timestamps or metadata
		if line == "" || strings.HasPrefix(line, "[") || strings.Contains(line, "whisper_") {
			continue
		}

		// Skip lines that look like processing information
		if strings.Contains(line, "processing") || strings.Contains(line, "ms") {
			continue
		}

		transcription.WriteString(line)
		transcription.WriteString(" ")
	}

	return strings.TrimSpace(transcription.String())
}

// GetRecommendedModel returns a recommended model based on system resources
func (wm *WhisperManager) GetRecommendedModel() string {
	// For simplicity, recommend base model as good balance of speed/accuracy
	return "base"
}

// GetVersion returns the whisper-cpp version
func (wm *WhisperManager) GetVersion() string {
	if !wm.IsWhisperInstalled() {
		return "not installed"
	}

	cmd := exec.Command("whisper-cpp", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}

	// Parse version from help output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "whisper.cpp") {
			return strings.TrimSpace(line)
		}
	}

	return "installed"
}


