// Package whisperapp provides the audio transcription API endpoints using go-whisper.
package whisperapp

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/kawai-network/veridium/cmd/server/app/sdk/errs"
	"github.com/kawai-network/veridium/cmd/server/foundation/logger"
	"github.com/kawai-network/veridium/cmd/server/foundation/web"
	"github.com/mutablelogic/go-whisper/pkg/schema"
	"github.com/mutablelogic/go-whisper/pkg/whisper"
)

// app represents the whisper application
type app struct {
	log     *logger.Logger
	manager *whisper.Manager
}

// Config contains all the mandatory systems required by handlers.
type Config struct {
	Log     *logger.Logger
	Manager *whisper.Manager
}

func newApp(cfg Config) *app {
	return &app{
		log:     cfg.Log,
		manager: cfg.Manager,
	}
}

// TranscriptionRequest represents an OpenAI-compatible audio transcription request
type TranscriptionRequest struct {
	Model       string `json:"model" form:"model" binding:"required"`
	Language    string `json:"language,omitempty" form:"language"`
	Prompt      string `json:"prompt,omitempty" form:"prompt"`
	ResponseFmt string `json:"response_format,omitempty" form:"response_format"` // json, text, srt, vtt
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
	if a.manager == nil {
		return errs.Errorf(errs.Unimplemented, "whisper service not available")
	}

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
	model := r.FormValue("model")
	if model == "" {
		model = "base" // Default model
	}

	responseFormat := r.FormValue("response_format")
	if responseFormat == "" {
		responseFormat = "json"
	}

	// Get the model
	whisperModel := a.manager.GetModelById(model)
	if whisperModel == nil {
		return errs.Errorf(errs.NotFound, "model %s not found", model)
	}

	// Perform transcription
	text, err := a.transcribe(ctx, whisperModel, file)
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
func (a *app) transcribe(ctx context.Context, model *schema.Model, file io.Reader) (string, error) {
	// Read audio data
	audioData, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to read audio data: %w", err)
	}

	var result string
	err = a.manager.WithModel(model, func(task *whisper.Task) error {
		// Create a segment callback to collect text
		var segments []string
		segmentCallback := func(seg *schema.Segment) {
			if seg != nil && seg.Text != "" {
				segments = append(segments, seg.Text)
			}
		}

		// Transcribe - convert []byte to []float32 if needed
		// go-whisper expects float32 audio samples
		samples := bytesToFloat32(audioData)
		
		err := task.Transcribe(ctx, 0, samples, segmentCallback)
		if err != nil {
			return err
		}

		// Join all segments
		result = joinSegments(segments)
		return nil
	})

	if err != nil {
		return "", err
	}

	return result, nil
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

// TextResponse represents a plain text response
type TextResponse struct {
	Text string
}

// Encode implements web.Encoder for plain text
func (r TextResponse) Encode() ([]byte, string, error) {
	return []byte(r.Text), "text/plain", nil
}
