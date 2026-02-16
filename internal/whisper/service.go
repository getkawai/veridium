package whisper

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/kawai-network/veridium/cmd/server/app/domain/whisperapp"
	"github.com/kawai-network/veridium/internal/paths"
	whisper "github.com/kawai-network/whisper"
)

type Service struct {
	libDir    string
	modelsDir string
}

func NewService() (*Service, error) {
	libDir := paths.WhisperLib()
	modelsDir := paths.WhisperModels()

	if err := os.MkdirAll(libDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create lib directory: %w", err)
	}
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create models directory: %w", err)
	}

	return &Service{
		libDir:    libDir,
		modelsDir: modelsDir,
	}, nil
}

func (s *Service) Transcribe(ctx context.Context, modelName, audioPath string) (string, error) {
	modelPath, err := whisperapp.GetModelFilePath(modelName)
	if err != nil {
		return "", fmt.Errorf("model not found: %w", err)
	}

	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		return "", fmt.Errorf("model file not found at %s. Run setup to download models", modelPath)
	}

	w, err := whisper.New(s.libDir)
	if err != nil {
		return "", fmt.Errorf("failed to create whisper instance: %w", err)
	}
	defer w.Close()

	if err := w.Load(modelPath); err != nil {
		return "", fmt.Errorf("failed to load model: %w", err)
	}

	opts := whisper.TranscriptionOptions{
		Threads:   4,
		Language:  "",
		Translate: false,
		Diarize:   false,
		Prompt:    "",
	}

	result, err := w.Transcribe(audioPath, opts)
	if err != nil {
		return "", fmt.Errorf("transcription failed: %w", err)
	}

	return result.Text, nil
}

func (s *Service) TranscribeWithOptions(ctx context.Context, modelName, audioPath string, options map[string]interface{}) (string, error) {
	opts := whisper.TranscriptionOptions{
		Threads: 4,
	}

	if lang, ok := options["language"].(string); ok {
		opts.Language = lang
	}
	if threads, ok := options["threads"].(int); ok {
		opts.Threads = uint32(threads)
	}
	if translate, ok := options["translate"].(bool); ok {
		opts.Translate = translate
	}
	if diarize, ok := options["diarize"].(bool); ok {
		opts.Diarize = diarize
	}
	if prompt, ok := options["prompt"].(string); ok {
		opts.Prompt = prompt
	}

	modelPath, err := whisperapp.GetModelFilePath(modelName)
	if err != nil {
		return "", fmt.Errorf("model not found: %w", err)
	}

	w, err := whisper.New(s.libDir)
	if err != nil {
		return "", fmt.Errorf("failed to create whisper instance: %w", err)
	}
	defer w.Close()

	if err := w.Load(modelPath); err != nil {
		return "", fmt.Errorf("failed to load model: %w", err)
	}

	result, err := w.Transcribe(audioPath, opts)
	if err != nil {
		return "", fmt.Errorf("transcription failed: %w", err)
	}

	return result.Text, nil
}

func (s *Service) ListModels() ([]string, error) {
	return whisperapp.ListDownloadedModels()
}

func (s *Service) IsModelDownloaded(modelName string) bool {
	return whisperapp.IsModelDownloaded(modelName)
}

func (s *Service) DownloadModel(ctx context.Context, modelName string) error {
	return whisperapp.DownloadModelWithProgress(ctx, modelName, nil)
}

func (s *Service) GetRecommendedModel() string {
	return "base"
}

func (s *Service) GetModelsDirectory() string {
	return s.modelsDir
}

func (s *Service) Close() error {
	return nil
}

func (s *Service) IsWhisperInstalled() bool {
	libName := whisper.LibraryName(runtime.GOOS)
	libFile := filepath.Join(s.libDir, libName)
	_, err := os.Stat(libFile)
	return err == nil
}

func (s *Service) InstallWhisper() error {
	ctx := context.Background()
	return whisperapp.QuickSetup(ctx, "base")
}

func (s *Service) GetVersion() string {
	return "gowhisper"
}

func (s *Service) IsFFmpegInstalled() bool {
	_, err := filepath.Abs("ffmpeg")
	if err != nil {
		_, err = filepath.Abs("/usr/bin/ffmpeg")
	}
	return err == nil
}

func (s *Service) InstallFFmpeg() error {
	log.Println("FFmpeg installation not supported. Please install manually: brew install ffmpeg (macOS) or sudo apt install ffmpeg (Linux)")
	return fmt.Errorf("ffmpeg installation not supported")
}

func (s *Service) GetFFmpegVersion() string {
	return "unknown"
}

func (s *Service) GetFFmpegPath() string {
	return ""
}

func (s *Service) GetAvailableModels() []string {
	models := whisperapp.GetAllModels()
	result := make([]string, len(models))
	for i, m := range models {
		result[i] = m.Name
	}
	return result
}

func (s *Service) GetDependencyStatus() map[string]interface{} {
	return map[string]interface{}{
		"whisper": map[string]interface{}{
			"installed": s.IsWhisperInstalled(),
			"version":   s.GetVersion(),
		},
		"ffmpeg": map[string]interface{}{
			"installed": s.IsFFmpegInstalled(),
			"version":   s.GetFFmpegVersion(),
			"path":      s.GetFFmpegPath(),
		},
	}
}
