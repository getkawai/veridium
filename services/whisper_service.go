package services

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	wav "github.com/go-audio/wav"
	whisper "github.com/mutablelogic/go-whisper"
	"github.com/mutablelogic/go-whisper/pkg/schema"
	"github.com/mutablelogic/go-whisper/pkg/task"
)

// WhisperService provides speech-to-text transcription using whisper.cpp
type WhisperService struct {
	whisper   *whisper.Whisper
	modelsDir string
}

// NewWhisperService creates a new Whisper service instance
// initWhisper initializes a whisper instance with the given models directory
func initWhisper(modelsDir string) (*whisper.Whisper, error) {
	// Create models directory if it doesn't exist
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create models directory: %w", err)
	}

	// Create whisper instance with GPU support
	// Note: GPU is enabled by default, use OptNoGPU() to disable
	w, err := whisper.New(modelsDir,
		whisper.OptMaxConcurrent(2), // Allow 2 concurrent transcriptions
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create whisper service: %w", err)
	}

	return w, nil
}

func NewWhisperService() (*WhisperService, error) {
	// Get user data directory for model storage
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		userConfigDir = "."
	}
	modelsDir := filepath.Join(userConfigDir, "veridium", "whisper-models")

	w, err := initWhisper(modelsDir)
	if err != nil {
		return nil, err
	}

	return &WhisperService{
		whisper:   w,
		modelsDir: modelsDir,
	}, nil
}

// Close releases all resources used by the Whisper service
func (s *WhisperService) Close() error {
	if s.whisper != nil {
		return s.whisper.Close()
	}
	return nil
}

// GetModelsDirectory returns the path to the models directory
func (s *WhisperService) GetModelsDirectory() string {
	return s.modelsDir
}

// ListModels returns all available Whisper models
func (s *WhisperService) ListModels() []*schema.Model {
	if s.whisper == nil {
		return nil
	}
	return s.whisper.ListModels()
}

// GetModel returns a model by its ID
func (s *WhisperService) GetModel(id string) *schema.Model {
	if s.whisper == nil {
		return nil
	}
	return s.whisper.GetModelById(id)
}

// DownloadModel downloads a Whisper model from HuggingFace
// modelName examples: "ggml-base.bin", "ggml-small.bin", "ggml-medium.bin"
func (s *WhisperService) DownloadModel(ctx context.Context, modelName string) error {
	if s.whisper == nil {
		return fmt.Errorf("whisper service not initialized")
	}

	_, err := s.whisper.DownloadModel(ctx, modelName, func(current, total uint64) {
		// Progress callback - could emit event to frontend
		if total > 0 {
			progress := float64(current) / float64(total) * 100
			fmt.Printf("Download progress: %.2f%% (%d/%d bytes)\n", progress, current, total)
		}
	})

	return err
}

// DeleteModel deletes a model by its ID
func (s *WhisperService) DeleteModel(id string) error {
	if s.whisper == nil {
		return fmt.Errorf("whisper service not initialized")
	}
	return s.whisper.DeleteModelById(id)
}

// Transcribe transcribes an audio file to text
// modelId: the ID of the model to use (e.g., "ggml-base")
// audioPath: path to the audio file
func (s *WhisperService) Transcribe(ctx context.Context, modelId, audioPath string) (string, error) {
	if s.whisper == nil {
		return "", fmt.Errorf("whisper service not initialized")
	}

	model := s.whisper.GetModelById(modelId)
	if model == nil {
		return "", fmt.Errorf("model not found: %s", modelId)
	}

	// Load audio samples using ffmpeg
	samples, err := s.loadAudioSamples(audioPath)
	if err != nil {
		return "", fmt.Errorf("failed to load audio: %w", err)
	}

	var result []*schema.Segment
	err = s.whisper.WithModel(model, func(t *task.Context) error {
		// Transcribe with callback to collect segments
		err := t.Transcribe(ctx, 0, samples, func(seg *schema.Segment) {
			result = append(result, seg)
		})
		if err != nil {
			return fmt.Errorf("failed to transcribe: %w", err)
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	// Combine all segment texts
	var fullText string
	for _, seg := range result {
		fullText += seg.Text
	}

	return fullText, nil
}

// TranscribeWithSegments returns transcription with timestamps
// Returns segments with start/end times for each text segment
func (s *WhisperService) TranscribeWithSegments(ctx context.Context, modelId, audioPath string) ([]*schema.Segment, error) {
	if s.whisper == nil {
		return nil, fmt.Errorf("whisper service not initialized")
	}

	model := s.whisper.GetModelById(modelId)
	if model == nil {
		return nil, fmt.Errorf("model not found: %s", modelId)
	}

	// Load audio samples using ffmpeg
	samples, err := s.loadAudioSamples(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load audio: %w", err)
	}

	var segments []*schema.Segment
	err = s.whisper.WithModel(model, func(t *task.Context) error {
		// Transcribe with callback to collect segments
		err := t.Transcribe(ctx, 0, samples, func(seg *schema.Segment) {
			segments = append(segments, seg)
		})
		if err != nil {
			return fmt.Errorf("failed to transcribe: %w", err)
		}

		return nil
	})

	return segments, err
}

// loadAudioSamples loads audio file and converts to float32 samples for whisper
// Currently only supports WAV files. For other formats, convert to WAV first.
func (s *WhisperService) loadAudioSamples(audioPath string) ([]float32, error) {
	// Open WAV file
	fh, err := os.Open(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open audio file: %w", err)
	}
	defer fh.Close()

	// Decode WAV
	decoder := wav.NewDecoder(fh)
	buf, err := decoder.FullPCMBuffer()
	if err != nil {
		return nil, fmt.Errorf("failed to decode WAV: %w", err)
	}

	// Convert to float32
	samples := buf.AsFloat32Buffer().Data
	return samples, nil
}

// GetAvailableModels returns a list of recommended models for download
func (s *WhisperService) GetAvailableModels() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"id":          "ggml-tiny.bin",
			"name":        "Tiny",
			"size":        "75 MB",
			"description": "Fastest, lowest accuracy",
		},
		{
			"id":          "ggml-base.bin",
			"name":        "Base",
			"size":        "142 MB",
			"description": "Good balance of speed and accuracy",
		},
		{
			"id":          "ggml-small.bin",
			"name":        "Small",
			"size":        "466 MB",
			"description": "Better accuracy, slower",
		},
		{
			"id":          "ggml-medium.bin",
			"name":        "Medium",
			"size":        "1.5 GB",
			"description": "High accuracy, much slower",
		},
		{
			"id":          "ggml-large-v3.bin",
			"name":        "Large V3",
			"size":        "3.1 GB",
			"description": "Best accuracy, very slow",
		},
	}
}
