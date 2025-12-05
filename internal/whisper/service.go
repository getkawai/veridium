package whisper

import (
	"context"
	"fmt"
	"log"
	"os"
)

// Service provides speech-to-text transcription using whisper-cpp CLI
// This is a Wails-compatible service wrapper around Manager
type Service struct {
	manager *Manager
}

// NewService creates a new Whisper service instance
// Automatically installs whisper-cpp and downloads recommended model in background
func NewService() (*Service, error) {
	manager := NewManager()

	service := &Service{
		manager: manager,
	}

	// Start background initialization
	go service.initializeInBackground()

	return service, nil
}

// initializeInBackground handles whisper-cpp installation and model download
func (s *Service) initializeInBackground() {
	ctx := context.Background()

	log.Println("🚀 Initializing Whisper STT in background...")

	// Step 1: Check and install FFmpeg (required by whisper-cpp for audio processing)
	if !s.manager.IsFFmpegInstalled() {
		log.Println("🔧 ffmpeg not found, attempting auto-installation...")
		log.Println("   (FFmpeg is required by whisper-cpp for audio format conversion)")
		if err := s.manager.InstallFFmpeg(); err != nil {
			log.Printf("⚠️  Failed to auto-install ffmpeg: %v", err)
			log.Printf("   Please install manually:")
			log.Printf("   - macOS: brew install ffmpeg")
			log.Printf("   - Linux: sudo apt install ffmpeg")
			log.Printf("   - Windows: winget install ffmpeg")
			// Continue anyway - whisper might still work with WAV files
		} else {
			log.Println("✅ ffmpeg installed successfully")
		}
	} else {
		ffmpegVersion := s.manager.GetFFmpegVersion()
		log.Printf("✅ ffmpeg is installed: %s", ffmpegVersion)
	}

	// Step 2: Check and install whisper-cpp if needed
	if !s.manager.IsWhisperInstalled() {
		log.Println("🔧 whisper-cpp not found, attempting auto-installation...")
		if err := s.manager.InstallWhisper(); err != nil {
			log.Printf("⚠️  Failed to auto-install whisper-cpp: %v", err)
			log.Printf("   Please install manually: brew install whisper-cpp")
			return
		}
		log.Println("✅ whisper-cpp installed successfully")
	} else {
		log.Println("✅ whisper-cli is installed and ready")
	}

	// Step 3: Check if any model is installed (this also validates and removes corrupted models)
	log.Println("🔍 Checking for installed Whisper models...")
	models, err := s.manager.GetInstalledModels()
	if err != nil {
		log.Printf("⚠️  Failed to check installed models: %v", err)
		return
	}

	// Step 4: Download recommended model if no models exist
	if len(models) == 0 {
		recommendedModel := s.manager.GetRecommendedModel()
		log.Printf("📥 No valid models found, downloading recommended model: %s", recommendedModel)
		log.Printf("   This will take a few minutes (142 MB download)...")
		log.Printf("   Download location: %s", s.manager.ModelsDir)

		if err := s.manager.DownloadModel(ctx, recommendedModel); err != nil {
			log.Printf("⚠️  Failed to download model %s: %v", recommendedModel, err)
			log.Printf("   You can download manually later from the UI")
			return
		}

		log.Printf("✅ Model %s downloaded successfully (141 MB)", recommendedModel)
		log.Println("🎉 Whisper STT is now ready to use!")
	} else {
		log.Printf("✅ Found %d valid model(s):", len(models))
		for _, model := range models {
			modelPath := s.manager.GetModelPath(model)
			if info, err := os.Stat(modelPath); err == nil {
				log.Printf("   - %s (%.2f MB)", model, float64(info.Size())/(1024*1024))
			}
		}
		log.Println("✅ Whisper STT is ready to use!")
	}
}

// GetModelsDirectory returns the path to the models directory
func (s *Service) GetModelsDirectory() string {
	return s.manager.ModelsDir
}

// ListModels returns all installed Whisper models
func (s *Service) ListModels() ([]string, error) {
	return s.manager.GetInstalledModels()
}

// GetAvailableModels returns a list of models available for download
func (s *Service) GetAvailableModels() []map[string]interface{} {
	models := s.manager.GetAvailableModels()
	result := make([]map[string]interface{}, 0, len(models))

	// Map to user-friendly format
	sizeMap := map[string]string{
		"tiny":      "39 MB",
		"tiny.en":   "39 MB",
		"base":      "142 MB",
		"base.en":   "142 MB",
		"small":     "466 MB",
		"small.en":  "466 MB",
		"medium":    "1.5 GB",
		"medium.en": "1.5 GB",
		"large-v1":  "2.9 GB",
		"large-v2":  "2.9 GB",
		"large-v3":  "2.9 GB",
	}

	descMap := map[string]string{
		"tiny":      "Fastest, lowest accuracy",
		"tiny.en":   "Fastest, English only",
		"base":      "Good balance of speed and accuracy",
		"base.en":   "Good balance, English only",
		"small":     "Better accuracy, slower",
		"small.en":  "Better accuracy, English only",
		"medium":    "High accuracy, much slower",
		"medium.en": "High accuracy, English only",
		"large-v1":  "Best accuracy, very slow",
		"large-v2":  "Best accuracy, very slow",
		"large-v3":  "Best accuracy, very slow (latest)",
	}

	for _, model := range models {
		result = append(result, map[string]interface{}{
			"id":          fmt.Sprintf("ggml-%s.bin", model),
			"name":        model,
			"size":        sizeMap[model],
			"description": descMap[model],
		})
	}

	return result
}

// DownloadModel downloads a Whisper model
func (s *Service) DownloadModel(ctx context.Context, modelName string) error {
	return s.manager.DownloadModel(ctx, modelName)
}

// IsModelDownloaded checks if a model is downloaded
func (s *Service) IsModelDownloaded(modelName string) bool {
	return s.manager.IsModelDownloaded(modelName)
}

// Transcribe transcribes an audio file to text
// modelName: the name of the model to use (e.g., "base", "tiny")
// audioPath: path to the audio file
func (s *Service) Transcribe(ctx context.Context, modelName, audioPath string) (string, error) {
	// Default options
	options := map[string]interface{}{
		"language": "auto", // Auto-detect language
		"threads":  4,      // Use 4 threads
	}

	return s.manager.TranscribeFile(ctx, audioPath, modelName, options)
}

// TranscribeWithOptions transcribes with custom options
func (s *Service) TranscribeWithOptions(ctx context.Context, modelName, audioPath string, options map[string]interface{}) (string, error) {
	return s.manager.TranscribeFile(ctx, audioPath, modelName, options)
}

// IsWhisperInstalled checks if whisper-cpp is installed
func (s *Service) IsWhisperInstalled() bool {
	return s.manager.IsWhisperInstalled()
}

// InstallWhisper installs whisper-cpp (platform-specific)
func (s *Service) InstallWhisper() error {
	return s.manager.InstallWhisper()
}

// GetVersion returns the whisper-cpp version
func (s *Service) GetVersion() string {
	return s.manager.GetVersion()
}

// GetRecommendedModel returns a recommended model
func (s *Service) GetRecommendedModel() string {
	return s.manager.GetRecommendedModel()
}

// IsFFmpegInstalled checks if ffmpeg is installed
func (s *Service) IsFFmpegInstalled() bool {
	return s.manager.IsFFmpegInstalled()
}

// InstallFFmpeg installs ffmpeg (platform-specific)
func (s *Service) InstallFFmpeg() error {
	return s.manager.InstallFFmpeg()
}

// GetFFmpegVersion returns the ffmpeg version
func (s *Service) GetFFmpegVersion() string {
	return s.manager.GetFFmpegVersion()
}

// GetDependencyStatus returns the status of all dependencies
func (s *Service) GetDependencyStatus() map[string]interface{} {
	return map[string]interface{}{
		"whisper": map[string]interface{}{
			"installed": s.manager.IsWhisperInstalled(),
			"version":   s.manager.GetVersion(),
		},
		"ffmpeg": map[string]interface{}{
			"installed": s.manager.IsFFmpegInstalled(),
			"version":   s.manager.GetFFmpegVersion(),
			"path":      s.manager.GetFFmpegPath(),
		},
	}
}

// Close is a no-op for CLI-based service (kept for interface compatibility)
func (s *Service) Close() error {
	return nil
}
