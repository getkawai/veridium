// Package whisperapp provides the audio transcription API endpoints using pkg/whisper.
package whisperapp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/kawai-network/veridium/cmd/server/app/sdk/errs"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	"github.com/kawai-network/veridium/pkg/whisper"
	"github.com/kawai-network/veridium/pkg/whisper/model"
)

// app represents the whisper application
type app struct {
	log *logger.Logger
}

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log *logger.Logger
}

func newApp(cfg Config) *app {
	return &app{
		log: cfg.Log,
	}
}

// TranscriptionRequest represents an OpenAI-compatible audio transcription request
type TranscriptionRequest struct {
	Model       string  `json:"model" form:"model" binding:"required"`
	Language    string  `json:"language,omitempty" form:"language"`
	Prompt      string  `json:"prompt,omitempty" form:"prompt"`
	ResponseFmt string  `json:"response_format,omitempty" form:"response_format"` // json, text, srt, vtt
	Temperature float32 `json:"temperature,omitempty" form:"temperature"`
}

// TranscriptionResponse represents an OpenAI-compatible transcription response
type TranscriptionResponse struct {
	Text string `json:"text"`
}

// Encode implements web.Encoder.
func (r TranscriptionResponse) Encode() ([]byte, string, error) {
	data, err := json.Marshal(r)
	return data, "application/json", err
}

// transcriptions handles POST /v1/audio/transcriptions
func (a *app) transcriptions(ctx context.Context, r *http.Request) web.Encoder {
	// Parse multipart form (for file upload)
	if err := r.ParseMultipartForm(32 << 20); err != nil { // 32MB max
		return errs.New(errs.InvalidArgument, fmt.Errorf("failed to parse form: %w", err))
	}

	// Get the file
	file, header, err := r.FormFile("file")
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("file is required: %w", err))
	}
	defer file.Close()

	a.log.Info(ctx, "transcription request", "filename", header.Filename, "size", header.Size)

	// Parse other form fields
	modelName := r.FormValue("model")
	if modelName == "" {
		modelName = "base" // Default model
	}

	responseFormat := r.FormValue("response_format")
	if responseFormat == "" {
		responseFormat = "json"
	}

	language := r.FormValue("language")

	// Get models directory from environment or use default
	modelsDir := os.Getenv("WHISPER_MODELS_DIR")
	if modelsDir == "" {
		modelsDir = "./data/models/whisper"
	}

	// Check if model exists, download if not
	modelPath := model.GetModelPath(modelsDir, modelName)
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		a.log.Info(ctx, "downloading whisper model", "model", modelName)
		if err := model.DownloadModel(modelName, modelsDir, nil); err != nil {
			a.log.Error(ctx, "failed to download model", "error", err)
			return errs.Errorf(errs.NotFound, "model %s not available and download failed", modelName)
		}
	}

	// Perform transcription
	text, err := a.transcribe(ctx, modelPath, file, language)
	if err != nil {
		a.log.Error(ctx, "transcription failed", "error", err)
		return errs.New(errs.Internal, err)
	}

	// Format response based on requested format
	switch responseFormat {
	case "text":
		// Return plain text
		return &TextResponse{Text: text}
	case "srt":
		// TODO: Implement SRT format
		return &TranscriptionResponse{Text: text}
	case "vtt":
		// TODO: Implement VTT format
		return &TranscriptionResponse{Text: text}
	default:
		// JSON format (default)
		return &TranscriptionResponse{Text: text}
	}
}

// transcribe performs the actual transcription
func (a *app) transcribe(ctx context.Context, modelPath string, file io.Reader, language string) (string, error) {
	// Save uploaded file to temp location
	tempFile, err := os.CreateTemp("", "whisper-*.audio")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy uploaded data to temp file
	if _, err := io.Copy(tempFile, file); err != nil {
		return "", fmt.Errorf("failed to save audio file: %w", err)
	}
	tempFile.Close()

	// Create whisper instance
	cfg := whisper.Config{
		ModelPath: modelPath,
		UseGPU:    true,
	}

	w, err := whisper.New(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to create whisper: %w", err)
	}
	defer w.Free()

	// Set transcription options
	opts := []whisper.TranscribeOption{
		whisper.WithThreads(4),
	}
	if language != "" && language != "auto" {
		opts = append(opts, whisper.WithLanguage(language))
	}

	// Transcribe
	result, err := w.Transcribe(tempFile.Name(), opts...)
	if err != nil {
		return "", fmt.Errorf("transcription failed: %w", err)
	}

	return result.Text, nil
}

// TextResponse represents a plain text response
type TextResponse struct {
	Text string
}

// Encode implements web.Encoder for plain text
func (r TextResponse) Encode() ([]byte, string, error) {
	return []byte(r.Text), "text/plain", nil
}

// bytesToFloat32 converts byte audio data to float32 samples
// This is a simplified conversion - proper implementation would decode audio format
func bytesToFloat32(data []byte) []float32 {
	// For 16-bit PCM: convert bytes to int16 then normalize to float32
	samples := make([]float32, 0, len(data)/2)
	for i := 0; i < len(data)-1; i += 2 {
		val := int16(data[i]) | int16(data[i+1])<<8
		samples = append(samples, float32(val)/32768.0)
	}
	return samples
}

// joinSegments joins transcription segments into a single text
func joinSegments(segments []string) string {
	return strings.Join(segments, " ")
}
