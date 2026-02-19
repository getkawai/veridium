// Package whisperapp provides the audio transcription API endpoints using github.com/kawai-network/whisper.
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
	Model                 string   `json:"model" form:"model" binding:"required"`
	Language              string   `json:"language,omitempty" form:"language"`
	Prompt                string   `json:"prompt,omitempty" form:"prompt"`
	ResponseFormat        string   `json:"response_format,omitempty" form:"response_format"` // json, text, srt, verbose_json, vtt
	Temperature           float32  `json:"temperature,omitempty" form:"temperature"`
	TimestampGranularities []string `json:"timestamp_granularities,omitempty" form:"timestamp_granularities"` // word, segment
}

// TranscriptionResponse represents a basic OpenAI-compatible transcription response
type TranscriptionResponse struct {
	Text string `json:"text"`
}

// Encode implements web.Encoder.
func (r TranscriptionResponse) Encode() ([]byte, string, error) {
	data, err := json.Marshal(r)
	return data, "application/json", err
}

// VerboseTranscriptionResponse represents a verbose OpenAI-compatible transcription response
type VerboseTranscriptionResponse struct {
	Task     string  `json:"task"`
	Language string  `json:"language"`
	Duration float64 `json:"duration"`
	Text     string  `json:"text"`
	Words    []Word  `json:"words,omitempty"`
	Segments []Segment `json:"segments,omitempty"`
}

// Word represents a word with timestamp
type Word struct {
	Word  string  `json:"word"`
	Start float64 `json:"start"`
	End   float64 `json:"end"`
}

// Segment represents a segment with detailed information
type Segment struct {
	ID               int32   `json:"id"`
	Seek             int32   `json:"seek"`
	Start            float64 `json:"start"`
	End              float64 `json:"end"`
	Text             string  `json:"text"`
	Tokens           []int32 `json:"tokens"`
	Temperature      float32 `json:"temperature"`
	AvgLogprob       float64 `json:"avg_logprob"`
	CompressionRatio float64 `json:"compression_ratio"`
	NoSpeechProb     float64 `json:"no_speech_prob"`
}

// Encode implements web.Encoder.
func (r VerboseTranscriptionResponse) Encode() ([]byte, string, error) {
	data, err := json.Marshal(r)
	return data, "application/json", err
}

// TranslationResponse represents an OpenAI-compatible translation response
type TranslationResponse struct {
	Text string `json:"text"`
}

// Encode implements web.Encoder.
func (r TranslationResponse) Encode() ([]byte, string, error) {
	data, err := json.Marshal(r)
	return data, "application/json", err
}

// transcriptions handles POST /v1/audio/transcriptions
func (a *app) transcriptions(ctx context.Context, r *http.Request) web.Encoder {
	// Parse multipart form (for file upload)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
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
	prompt := r.FormValue("prompt")

	// Parse temperature
	var temperature float32 = 0
	if temp := r.FormValue("temperature"); temp != "" {
		fmt.Sscanf(temp, "%f", &temperature)
	}

	// Get model path using unified ModelPath structure
	modelPath, err := GetModelFilePath(modelName)
	if err != nil {
		a.log.Error(ctx, "model lookup failed", "model", modelName, "error", err)
		return errs.Errorf(errs.NotFound, "model %s not found. Please run 'kawai-contributor setup' to download whisper models.", modelName)
	}

	// Check if model exists
	if _, statErr := os.Stat(modelPath); os.IsNotExist(statErr) {
		a.log.Error(ctx, "model not found", "model", modelName, "path", modelPath)
		return errs.Errorf(errs.NotFound, "model %s not found at %s. Please run 'kawai-contributor setup' to download whisper models.", modelName, modelPath)
	}

	// Perform transcription with full result
	result, err := a.transcribeWithResult(ctx, modelPath, file, language, prompt, temperature)
	if err != nil {
		a.log.Error(ctx, "transcription failed", "error", err)
		return errs.New(errs.Internal, err)
	}

	// Format response based on requested format
	switch responseFormat {
	case "text":
		// Return plain text
		return &TextResponse{Text: result.Text}
	case "srt":
		// TODO: Implement SRT format
		return &TranscriptionResponse{Text: result.Text}
	case "vtt":
		// TODO: Implement VTT format
		return &TranscriptionResponse{Text: result.Text}
	case "verbose_json":
		// Return verbose JSON with segments and timestamps
		return a.buildVerboseResponse(ctx, result, language)
	default:
		// JSON format (default)
		return &TranscriptionResponse{Text: result.Text}
	}
}

// translations handles POST /v1/audio/translations
func (a *app) translations(ctx context.Context, r *http.Request) web.Encoder {
	// Parse multipart form (for file upload)
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("failed to parse form: %w", err))
	}

	// Get the file
	file, header, err := r.FormFile("file")
	if err != nil {
		return errs.New(errs.InvalidArgument, fmt.Errorf("file is required: %w", err))
	}
	defer file.Close()

	a.log.Info(ctx, "translation request", "filename", header.Filename, "size", header.Size)

	// Parse form fields
	modelName := r.FormValue("model")
	if modelName == "" {
		modelName = "base" // Default model
	}

	responseFormat := r.FormValue("response_format")
	if responseFormat == "" {
		responseFormat = "json"
	}

	prompt := r.FormValue("prompt")

	// Parse temperature
	var temperature float32 = 0
	if temp := r.FormValue("temperature"); temp != "" {
		fmt.Sscanf(temp, "%f", &temperature)
	}

	// Get model path
	modelPath, err := GetModelFilePath(modelName)
	if err != nil {
		a.log.Error(ctx, "model lookup failed", "model", modelName, "error", err)
		return errs.Errorf(errs.NotFound, "model %s not found. Please run 'kawai-contributor setup' to download whisper models.", modelName)
	}

	// Check if model exists
	if _, statErr := os.Stat(modelPath); os.IsNotExist(statErr) {
		a.log.Error(ctx, "model not found", "model", modelName, "path", modelPath)
		return errs.Errorf(errs.NotFound, "model %s not found at %s. Please run 'kawai-contributor setup' to download whisper models.", modelName, modelPath)
	}

	// Perform translation (translate to English)
	result, err := a.translateWithResult(ctx, modelPath, file, prompt, temperature)
	if err != nil {
		a.log.Error(ctx, "translation failed", "error", err)
		return errs.New(errs.Internal, err)
	}

	// Format response based on requested format
	switch responseFormat {
	case "text":
		return &TextResponse{Text: result.Text}
	case "srt":
		// TODO: Implement SRT format
		return &TranslationResponse{Text: result.Text}
	case "vtt":
		// TODO: Implement VTT format
		return &TranslationResponse{Text: result.Text}
	case "verbose_json":
		// Return verbose JSON with segments
		return a.buildVerboseResponse(ctx, result, "en") // Translation is always to English
	default:
		// JSON format (default)
		return &TranslationResponse{Text: result.Text}
	}
}

// transcribeWithResult performs the actual transcription and returns full result
func (a *app) transcribeWithResult(ctx context.Context, modelPath string, file io.Reader, language, prompt string, temperature float32) (*whisper.TranscriptionResult, error) {
	// Save uploaded file to temp location
	tempFile, err := os.CreateTemp("", "whisper-*.audio")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy uploaded data to temp file
	if _, err := io.Copy(tempFile, file); err != nil {
		return nil, fmt.Errorf("failed to save audio file: %w", err)
	}
	tempFile.Close()

	// Get whisper library directory using centralized paths
	libDir := paths.WhisperLib()

	// Create whisper instance
	w, err := whisper.New(libDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create whisper instance: %w", err)
	}
	defer w.Close()

	// Load model
	if err := w.Load(modelPath); err != nil {
		return nil, fmt.Errorf("failed to load model from %s: %w", modelPath, err)
	}

	// Prepare transcription options
	opts := whisper.TranscriptionOptions{
		Threads:   defaultThreads,
		Language:  language,
		Translate: false,
		Diarize:   false,
		Prompt:    prompt,
	}

	// If language is empty or "auto", let whisper auto-detect
	if language == "" || language == "auto" {
		opts.Language = ""
	}

	a.log.Info(ctx, "starting transcription", "threads", opts.Threads, "language", opts.Language, "prompt", prompt)

	// Transcribe
	transcribeResult, err := w.Transcribe(tempFile.Name(), opts)
	if err != nil {
		return nil, fmt.Errorf("transcription failed: %w", err)
	}

	a.log.Info(ctx, "transcription completed", "text_length", len(transcribeResult.Text), "segments", len(transcribeResult.Segments))

	return &transcribeResult, nil
}

// translateWithResult performs audio translation to English
func (a *app) translateWithResult(ctx context.Context, modelPath string, file io.Reader, prompt string, temperature float32) (*whisper.TranscriptionResult, error) {
	// Save uploaded file to temp location
	tempFile, err := os.CreateTemp("", "whisper-*.audio")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	// Copy uploaded data to temp file
	if _, err := io.Copy(tempFile, file); err != nil {
		return nil, fmt.Errorf("failed to save audio file: %w", err)
	}
	tempFile.Close()

	// Get whisper library directory
	libDir := paths.WhisperLib()

	// Create whisper instance
	w, err := whisper.New(libDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create whisper instance: %w", err)
	}
	defer w.Close()

	// Load model
	if err := w.Load(modelPath); err != nil {
		return nil, fmt.Errorf("failed to load model from %s: %w", modelPath, err)
	}

	// Prepare translation options (Translate=true for English output)
	opts := whisper.TranscriptionOptions{
		Threads:   defaultThreads,
		Language:  "", // Auto-detect source language
		Translate: true,
		Diarize:   false,
		Prompt:    prompt,
	}

	a.log.Info(ctx, "starting translation", "threads", opts.Threads)

	// Translate
	transcribeResult, err := w.Transcribe(tempFile.Name(), opts)
	if err != nil {
		return nil, fmt.Errorf("translation failed: %w", err)
	}

	a.log.Info(ctx, "translation completed", "text_length", len(transcribeResult.Text), "segments", len(transcribeResult.Segments))

	return &transcribeResult, nil
}

// buildVerboseResponse builds a verbose JSON response from transcription result
func (a *app) buildVerboseResponse(ctx context.Context, result *whisper.TranscriptionResult, language string) *VerboseTranscriptionResponse {
	// Calculate duration from last segment
	var duration float64
	if len(result.Segments) > 0 {
		lastSegment := result.Segments[len(result.Segments)-1]
		duration = float64(lastSegment.End) / 1000.0 // Convert ms to seconds
	}

	// Build segments
	segments := make([]Segment, len(result.Segments))
	for i, seg := range result.Segments {
		segments[i] = Segment{
			ID:               seg.Id,
			Seek:             0, // Not provided by whisper library
			Start:            float64(seg.Start) / 1000.0, // Convert ms to seconds
			End:              float64(seg.End) / 1000.0,
			Text:             seg.Text,
			Tokens:           seg.Tokens,
			Temperature:      0, // Not provided
			AvgLogprob:       0, // Not provided
			CompressionRatio: 0, // Not provided
			NoSpeechProb:     0, // Not provided
		}
	}

	// Build words from segments (simplified - extract words from segment text)
	var words []Word
	for _, seg := range segments {
		words = append(words, splitSegmentToWords(seg.Text, seg.Start, seg.End)...)
	}

	return &VerboseTranscriptionResponse{
		Task:     "transcribe",
		Language: language,
		Duration: duration,
		Text:     result.Text,
		Segments: segments,
		Words:    words,
	}
}

// splitSegmentToWords splits segment text into words with estimated timestamps
func splitSegmentToWords(text string, start, end float64) []Word {
	// Simple word splitting - in production, use whisper's word-level timestamps
	words := strings.Fields(text)
	if len(words) == 0 {
		return nil
	}

	duration := (end - start) / float64(len(words))
	result := make([]Word, len(words))

	for i, word := range words {
		result[i] = Word{
			Word:  word,
			Start: start + float64(i)*duration,
			End:   start + float64(i+1)*duration,
		}
	}

	return result
}

// TextResponse represents a plain text response
type TextResponse struct {
	Text string
}

// Encode implements web.Encoder for plain text
func (r TextResponse) Encode() ([]byte, string, error) {
	return []byte(r.Text), "text/plain", nil
}

// Constants for transcription request limits
const (
	// maxUploadSize is the maximum allowed upload file size (32MB)
	maxUploadSize = 32 << 20 // 32 MB
	// defaultThreads is the default number of CPU threads for transcription
	defaultThreads = 4
)
