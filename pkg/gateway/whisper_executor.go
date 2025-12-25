package gateway

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/kawai-network/veridium/internal/whisper"
)

// WhisperExecutor implements execution using internal/whisper service.
type WhisperExecutor struct {
	service *whisper.Service
}

// NewWhisperExecutor creates a new executor backed by whisper service.
func NewWhisperExecutor(service *whisper.Service) *WhisperExecutor {
	return &WhisperExecutor{
		service: service,
	}
}

// Transcribe handles the actual transcription process using a multipart file.
// It creates a temporary file as required by the whisper CLI wrapper.
func (e *WhisperExecutor) Transcribe(ctx context.Context, file *multipart.FileHeader, modelName string) (string, error) {
	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	// Create temp file
	tempFile, err := os.CreateTemp("", "whisper-upload-*.wav")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name()) // Clean up
	defer tempFile.Close()

	// Copy content
	if _, err := io.Copy(tempFile, src); err != nil {
		return "", fmt.Errorf("failed to write temp file: %w", err)
	}

	// Determine model to use
	targetModel := modelName
	if targetModel == "" || targetModel == "whisper-1" {
		targetModel = "base" // Default to base model
	}

	// Transcribe
	// Note: whisper.Service.Transcribe takes (ctx, modelName, audioPath)
	return e.service.Transcribe(ctx, targetModel, tempFile.Name())
}
