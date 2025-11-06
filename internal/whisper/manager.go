package whisper

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Manager handles whisper-cpp installation and usage
type Manager struct {
	ModelsDir string // Directory where whisper models are stored
}

// NewManager creates a new whisper manager
func NewManager() *Manager {
	homeDir, _ := os.UserHomeDir()
	modelsDir := filepath.Join(homeDir, ".kawai-agent", "whisper-models")

	// Ensure models directory exists
	os.MkdirAll(modelsDir, 0755)

	return &Manager{
		ModelsDir: modelsDir,
	}
}

// IsWhisperInstalled checks if whisper-cpp is available
func (m *Manager) IsWhisperInstalled() bool {
	_, err := exec.LookPath("whisper-cpp")
	return err == nil
}

// GetAvailableModels returns a list of available whisper models
func (m *Manager) GetAvailableModels() []string {
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
func (m *Manager) DownloadModel(ctx context.Context, modelName string) error {
	modelPath := filepath.Join(m.ModelsDir, fmt.Sprintf("ggml-%s.bin", modelName))

	// Check if model already exists
	if _, err := os.Stat(modelPath); err == nil {
		log.Printf("Model %s already exists", modelName)
		return nil
	}

	log.Printf("Downloading whisper model: %s", modelName)

	// Download using curl
	url := fmt.Sprintf("https://huggingface.co/ggerganov/whisper.cpp/resolve/main/ggml-%s.bin", modelName)
	cmd := exec.CommandContext(ctx, "curl", "-L", "-o", modelPath, url)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to download model %s: %w\nOutput: %s", modelName, err, string(output))
	}

	log.Printf("Successfully downloaded model: %s", modelName)
	return nil
}

// GetModelPath returns the path to a specific model
func (m *Manager) GetModelPath(modelName string) string {
	return filepath.Join(m.ModelsDir, fmt.Sprintf("ggml-%s.bin", modelName))
}

// IsModelDownloaded checks if a model is downloaded
func (m *Manager) IsModelDownloaded(modelName string) bool {
	modelPath := m.GetModelPath(modelName)
	_, err := os.Stat(modelPath)
	return err == nil
}

// GetInstalledModels returns a list of downloaded models
func (m *Manager) GetInstalledModels() ([]string, error) {
	files, err := os.ReadDir(m.ModelsDir)
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
func (m *Manager) TranscribeFile(ctx context.Context, audioPath, modelName string, options map[string]interface{}) (string, error) {
	if !m.IsWhisperInstalled() {
		return "", fmt.Errorf("whisper-cpp is not installed")
	}

	// Check if model exists
	if !m.IsModelDownloaded(modelName) {
		return "", fmt.Errorf("model %s is not downloaded", modelName)
	}

	// Check if audio file exists
	if _, err := os.Stat(audioPath); err != nil {
		return "", fmt.Errorf("audio file not found: %w", err)
	}

	modelPath := m.GetModelPath(modelName)

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
	cmd := exec.CommandContext(ctx, "whisper-cpp", args...)
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
	return m.parseWhisperOutput(string(output)), nil
}

// parseWhisperOutput extracts transcription from whisper-cpp output
func (m *Manager) parseWhisperOutput(output string) string {
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
func (m *Manager) GetRecommendedModel() string {
	// For simplicity, recommend base model as good balance of speed/accuracy
	return "base"
}

// GetVersion returns the whisper-cpp version
func (m *Manager) GetVersion() string {
	if !m.IsWhisperInstalled() {
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

// Platform-specific installation methods are implemented in:
// - manager_darwin.go (macOS)
// - manager_linux.go (Linux)
// - manager_windows.go (Windows)
// Go compiler automatically selects the right file based on GOOS
