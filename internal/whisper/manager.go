package whisper

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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
	_ = os.MkdirAll(modelsDir, 0755)

	return &Manager{
		ModelsDir: modelsDir,
	}
}

// IsWhisperInstalled checks if whisper-cpp is available
func (m *Manager) IsWhisperInstalled() bool {
	_, err := exec.LookPath("whisper-cli")
	return err == nil
}

// IsFFmpegInstalled checks if ffmpeg is available in PATH or local directory
func (m *Manager) IsFFmpegInstalled() bool {
	// Check PATH first
	if _, err := exec.LookPath("ffmpeg"); err == nil {
		return true
	}

	// Check ~/.local/bin/ffmpeg
	homeDir, err := os.UserHomeDir()
	if err == nil {
		localFFmpeg := filepath.Join(homeDir, ".local", "bin", "ffmpeg")
		if _, err := os.Stat(localFFmpeg); err == nil {
			return true
		}
	}

	return false
}

// GetFFmpegPath returns the path to ffmpeg binary
func (m *Manager) GetFFmpegPath() string {
	// Check PATH first
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return path
	}

	// Check ~/.local/bin/ffmpeg
	homeDir, err := os.UserHomeDir()
	if err == nil {
		localFFmpeg := filepath.Join(homeDir, ".local", "bin", "ffmpeg")
		if _, err := os.Stat(localFFmpeg); err == nil {
			return localFFmpeg
		}
	}

	return ""
}

// GetFFmpegVersion returns the ffmpeg version string
func (m *Manager) GetFFmpegVersion() string {
	ffmpegPath := m.GetFFmpegPath()
	if ffmpegPath == "" {
		return "not installed"
	}

	cmd := exec.Command(ffmpegPath, "-version")
	output, err := cmd.Output()
	if err != nil {
		return "unknown"
	}

	// Parse first line for version
	lines := strings.Split(string(output), "\n")
	if len(lines) > 0 {
		return strings.TrimSpace(lines[0])
	}

	return "installed"
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

	// Verify the file was downloaded and has reasonable size
	info, err := os.Stat(modelPath)
	if err != nil {
		return fmt.Errorf("model file not found after download: %w", err)
	}

	// Models should be at least 10 MB (even tiny is ~39 MB)
	if info.Size() < 10*1024*1024 {
		_ = os.Remove(modelPath) // Remove incomplete file
		return fmt.Errorf("model download incomplete or corrupted (size: %d bytes, expected > 10 MB)", info.Size())
	}

	log.Printf("Successfully downloaded model: %s (size: %.2f MB)", modelName, float64(info.Size())/(1024*1024))
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
// Validates model files and removes corrupted ones
func (m *Manager) GetInstalledModels() ([]string, error) {
	files, err := os.ReadDir(m.ModelsDir)
	if err != nil {
		return nil, err
	}

	var models []string
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "ggml-") && strings.HasSuffix(file.Name(), ".bin") {
			modelPath := filepath.Join(m.ModelsDir, file.Name())

			// Validate model file size
			info, err := os.Stat(modelPath)
			if err != nil {
				log.Printf("⚠️  Warning: Cannot stat model file %s: %v", file.Name(), err)
				continue
			}

			// Check if model file is too small (corrupted/incomplete)
			if info.Size() < 10*1024*1024 { // Less than 10 MB
				log.Printf("⚠️  Removing corrupted model %s (size: %.2f MB, expected > 10 MB)",
					file.Name(), float64(info.Size())/(1024*1024))
				_ = os.Remove(modelPath)
				continue
			}

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
		"--output-txt", // Output as text file
		"--no-prints",  // Suppress debug output to stdout
	}

	// Add optional parameters
	// Note: whisper-cli defaults to "en" if no language is specified
	// We must explicitly pass "-l auto" for auto-detection to work properly
	if language, ok := options["language"].(string); ok && language != "" {
		args = append(args, "-l", language)
	} else {
		// Default to auto-detect if no language specified
		args = append(args, "-l", "auto")
	}
	if threads, ok := options["threads"].(int); ok && threads > 0 {
		args = append(args, "-t", fmt.Sprintf("%d", threads))
	}
	if translate, ok := options["translate"].(bool); ok && translate {
		args = append(args, "--translate")
	}

	log.Printf("Running whisper-cli with args: %v", args)

	// Execute whisper-cli
	cmd := exec.CommandContext(ctx, "whisper-cli", args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("transcription failed: %w\nOutput: %s", err, string(output))
	}

	// whisper-cli with --output-txt creates a .txt file next to the audio file
	txtFile := audioPath + ".txt"

	// Wait a moment for file to be written
	maxRetries := 10
	var content []byte
	for i := 0; i < maxRetries; i++ {
		if _, err := os.Stat(txtFile); err == nil {
			content, err = os.ReadFile(txtFile)
			if err == nil && len(content) > 0 {
				// Clean up the generated txt file
				_ = os.Remove(txtFile)
				transcription := strings.TrimSpace(string(content))
				log.Printf("Transcription result: %s", transcription)
				return transcription, nil
			}
		}
		if i < maxRetries-1 {
			// Wait 100ms before retry
			time.Sleep(100 * time.Millisecond)
		}
	}

	// If we couldn't read the file, return error
	if len(output) > 0 {
		log.Printf("Warning: Could not read output file %s, output was: %s", txtFile, string(output))
	}

	return "", fmt.Errorf("no transcription output generated")
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
// Model accuracy: large > medium > small > base > tiny
// IMPORTANT: Minimum recommended model is "small" for acceptable accuracy
// Model RAM requirements (approximate during inference):
//   - tiny:   ~1GB RAM   (~75MB download)  - NOT recommended
//   - base:   ~2GB RAM   (~150MB download) - NOT recommended
//   - small:  ~4GB RAM   (~500MB download) - MINIMUM recommended
//   - medium: ~8GB RAM   (~1.5GB download) - Good accuracy
//   - large:  ~16GB RAM  (~3GB download)   - Best accuracy
func (m *Manager) GetRecommendedModel() string {
	availableRAM := m.getAvailableRAMGB()

	// Select model based on available RAM
	// Minimum model is "small" for acceptable transcription quality
	switch {
	case availableRAM >= 16:
		return "medium" // Best practical choice for most use cases
	case availableRAM >= 8:
		return "medium" // Good accuracy
	default:
		return "small" // Minimum recommended for acceptable accuracy
	}
}

// getAvailableRAMGB returns available RAM in GB (platform-specific)
func (m *Manager) getAvailableRAMGB() int64 {
	// This will be overridden by platform-specific implementations
	return m.detectAvailableRAM()
}

// GetVersion returns the whisper-cpp version
func (m *Manager) GetVersion() string {
	if !m.IsWhisperInstalled() {
		return "not installed"
	}

	cmd := exec.Command("whisper-cli", "--help")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "unknown"
	}

	// Parse version from help output
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "whisper") {
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
