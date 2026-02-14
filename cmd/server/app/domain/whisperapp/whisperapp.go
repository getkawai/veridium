// Package whisperapp provides the audio transcription API endpoints using github.com/kawai-network/whisper.
package whisperapp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/kawai-network/veridium/cmd/server/app/sdk/errs"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	"github.com/kawai-network/veridium/internal/paths"
	whisper "github.com/kawai-network/whisper"
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

	// Get models directory using centralized paths
	modelsDir := paths.WhisperModels()

	// Model path for standalone whisper (format: {name}.bin, not ggml-{name}.bin)
	modelPath := fmt.Sprintf("%s/%s.bin", modelsDir, modelName)

	// Check if model exists
	if _, err := os.Stat(modelPath); os.IsNotExist(err) {
		a.log.Error(ctx, "model not found", "model", modelName, "path", modelPath)
		return errs.Errorf(errs.NotFound, "model %s not found at %s. Please run 'kawai-contributor setup' to download whisper models.", modelName, modelPath)
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

// transcribe performs the actual transcription using github.com/kawai-network/whisper
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

	// Get whisper library directory using centralized paths
	libDir := paths.WhisperLib()

	// Create whisper instance (will auto-download library if not found)
	w, err := whisper.New(libDir)
	if err != nil {
		return "", fmt.Errorf("failed to create whisper instance: %w", err)
	}
	defer w.Close()

	// Load model
	if err := w.Load(modelPath); err != nil {
		return "", fmt.Errorf("failed to load model from %s: %w", modelPath, err)
	}

	// Prepare transcription options
	// Note: github.com/kawai-network/whisper uses a different options structure
	opts := whisper.TranscriptionOptions{
		Threads:   4, // Default to 4 threads
		Language:  language,
		Translate: false,
		Diarize:   false,
		Prompt:    "",
	}

	// If language is empty or "auto", let whisper auto-detect
	if language == "" || language == "auto" {
		opts.Language = ""
	}

	a.log.Info(ctx, "starting transcription", "threads", opts.Threads, "language", opts.Language)

	// Transcribe
	result, err := w.Transcribe(tempFile.Name(), opts)
	if err != nil {
		return "", fmt.Errorf("transcription failed: %w", err)
	}

	a.log.Info(ctx, "transcription completed", "text_length", len(result.Text))

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
