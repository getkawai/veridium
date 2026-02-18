// Package ttsapp provides the text-to-speech API endpoints using TTS.cpp
package ttsapp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	tts "github.com/kawai-network/TTS.cpp/bindings/go"
	"github.com/kawai-network/veridium/cmd/server/app/sdk/errs"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	"github.com/kawai-network/veridium/internal/paths"
)

// app represents the TTS application
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

// SpeechRequest represents an OpenAI-compatible text-to-speech request
type SpeechRequest struct {
	Model          string  `json:"model" binding:"required"`
	Input          string  `json:"input" binding:"required"`
	Voice          string  `json:"voice,omitempty"`
	ResponseFormat string  `json:"response_format,omitempty"` // mp3, opus, aac, flac, wav, pcm
	Speed          float32 `json:"speed,omitempty"`
}

// SpeechResponse represents the audio data response
type SpeechResponse struct {
	Data        []byte
	ContentType string
}

// Encode implements web.Encoder.
func (r SpeechResponse) Encode() ([]byte, string, error) {
	return r.Data, r.ContentType, nil
}

// generations handles POST /v1/audio/speech
func (a *app) generations(ctx context.Context, r *http.Request) web.Encoder {
	var req SpeechRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("failed to parse request: %w", err))
	}

	// Validate required fields
	if req.Input == "" {
		return errs.New(errs.InvalidArgument, fmt.Errorf("input is required"))
	}

	if req.Model == "" {
		req.Model = "kokoro" // Default model
	}

	// Set defaults
	if req.ResponseFormat == "" {
		req.ResponseFormat = "mp3"
	}

	if req.Voice == "" {
		req.Voice = "af_sarah" // Default voice for Kokoro
	}

	a.log.Info(ctx, "tts request", "model", req.Model, "voice", req.Voice, "input_length", len(req.Input))

	// Get model path
	modelPath, err := a.getModelPath(req.Model)
	if err != nil {
		a.log.Error(ctx, "model lookup failed", "model", req.Model, "error", err)
		return errs.Errorf(errs.NotFound, "model %s not found. Please run 'kawai-contributor setup' to download TTS models.", req.Model)
	}

	// Check if model exists
	if _, statErr := os.Stat(modelPath); os.IsNotExist(statErr) {
		a.log.Error(ctx, "model not found", "model", req.Model, "path", modelPath)
		return errs.Errorf(errs.NotFound, "model %s not found at %s. Please run 'kawai-contributor setup' to download TTS models.", req.Model, modelPath)
	}

	// Generate speech
	audioData, err := a.generate(ctx, modelPath, req.Input, req.Voice)
	if err != nil {
		a.log.Error(ctx, "tts generation failed", "error", err)
		return errs.New(errs.Internal, err)
	}

	// Determine content type based on response format
	contentType := "audio/mpeg"
	switch req.ResponseFormat {
	case "wav":
		contentType = "audio/wav"
	case "opus":
		contentType = "audio/opus"
	case "aac":
		contentType = "audio/aac"
	case "flac":
		contentType = "audio/flac"
	case "pcm":
		contentType = "audio/pcm"
	}

	return &SpeechResponse{
		Data:        audioData,
		ContentType: contentType,
	}
}

// generate performs the actual text-to-speech generation using TTS.cpp
func (a *app) generate(ctx context.Context, modelPath, text, voice string) ([]byte, error) {
	// Create library config
	libConfig := tts.LibraryConfig{
		LibraryPath:  paths.TTSLib(), // Will auto-download if needed
		AutoDownload: true,
		Version:      "v0.1.4",
	}

	// Create TTS config
	config := tts.Config{
		Voice:             voice,
		TopK:              50,
		Temperature:       1.0,
		RepetitionPenalty: 1.0,
		TopP:              1.0,
		MaxTokens:         0,
		UseCrossAttention: true,
	}

	a.log.Info(ctx, "creating TTS runner", "model", modelPath, "voice", voice)

	// Create runner
	runner, err := tts.NewRunnerWithConfig(modelPath, 4, config, false, libConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create TTS runner: %w", err)
	}
	defer runner.Close()

	a.log.Info(ctx, "generating speech", "text_length", len(text))

	// Generate audio
	audio, err := runner.Generate(text)
	if err != nil {
		return nil, fmt.Errorf("speech generation failed: %w", err)
	}

	a.log.Info(ctx, "speech generation completed", "samples", len(audio.Samples), "sample_rate", audio.SampleRate)

	// Convert to WAV format (TTS.cpp outputs raw float samples)
	wavData, err := a.convertToWAV(audio.Samples, audio.SampleRate)
	if err != nil {
		return nil, fmt.Errorf("failed to convert audio: %w", err)
	}

	return wavData, nil
}

// convertToWAV converts float32 samples to WAV format
func (a *app) convertToWAV(samples []float32, sampleRate uint32) ([]byte, error) {
	// Simple WAV header + PCM data
	// This is a basic implementation - for production, consider using a proper audio library

	numChannels := 1
	bitsPerSample := 16
	byteRate := int(sampleRate) * numChannels * bitsPerSample / 8
	blockAlign := numChannels * bitsPerSample / 8
	dataSize := len(samples) * bitsPerSample / 8

	// WAV header
	header := make([]byte, 44)

	// RIFF chunk
	copy(header[0:4], "RIFF")
	header[4] = byte((dataSize + 36) & 0xFF)
	header[5] = byte(((dataSize + 36) >> 8) & 0xFF)
	header[6] = byte(((dataSize + 36) >> 16) & 0xFF)
	header[7] = byte(((dataSize + 36) >> 24) & 0xFF)
	copy(header[8:12], "WAVE")

	// fmt chunk
	copy(header[12:16], "fmt ")
	header[16] = 16 // Subchunk1Size
	header[17] = 0
	header[18] = 0
	header[19] = 0
	header[20] = 1 // AudioFormat (PCM)
	header[21] = 0
	header[22] = byte(numChannels)
	header[23] = 0
	header[24] = byte(sampleRate & 0xFF)
	header[25] = byte((sampleRate >> 8) & 0xFF)
	header[26] = byte((sampleRate >> 16) & 0xFF)
	header[27] = byte((sampleRate >> 24) & 0xFF)
	header[28] = byte(byteRate & 0xFF)
	header[29] = byte((byteRate >> 8) & 0xFF)
	header[30] = byte((byteRate >> 16) & 0xFF)
	header[31] = byte((byteRate >> 24) & 0xFF)
	header[32] = byte(blockAlign)
	header[33] = 0
	header[34] = byte(bitsPerSample)
	header[35] = 0

	// data chunk
	copy(header[36:40], "data")
	header[40] = byte(dataSize & 0xFF)
	header[41] = byte((dataSize >> 8) & 0xFF)
	header[42] = byte((dataSize >> 16) & 0xFF)
	header[43] = byte((dataSize >> 24) & 0xFF)

	// Convert float32 samples to int16
	data := make([]byte, dataSize)
	for i, sample := range samples {
		// Clamp to [-1, 1]
		if sample > 1.0 {
			sample = 1.0
		} else if sample < -1.0 {
			sample = -1.0
		}

		// Convert to int16
		val := int16(sample * 32767)
		data[i*2] = byte(val & 0xFF)
		data[i*2+1] = byte((val >> 8) & 0xFF)
	}

	return append(header, data...), nil
}

// getModelPath returns the path to a TTS model
func (a *app) getModelPath(modelName string) (string, error) {
	modelsPath := paths.Models()
	ttsModelsPath := filepath.Join(modelsPath, "tts")

	// Map model names to files
	modelFiles := map[string]string{
		"kokoro": DefaultTTSModelName,
	}

	filename, ok := modelFiles[modelName]
	if !ok {
		// Try to use the model name directly as a filename
		filename = modelName + ".gguf"
	}

	return filepath.Join(ttsModelsPath, filename), nil
}
