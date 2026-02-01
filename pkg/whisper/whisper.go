// Package whisper provides a pure Go binding for whisper.cpp using purego.
// No CGO required - uses dynamic loading of whisper.cpp shared library.
package whisper

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/kawai-network/veridium/pkg/whisper/audio"
	w "github.com/kawai-network/veridium/pkg/whisper/whisper"
)

// Whisper represents a whisper context for transcription
type Whisper struct {
	ctx        unsafe.Pointer
	modelPath  string
	converter  *audio.Converter
}

// Config represents configuration options for Whisper
type Config struct {
	ModelPath   string
	UseGPU      bool
	FlashAttn   bool
	GPUDevice   int32
	FFmpegPath  string
}

// New creates a new Whisper instance with the specified configuration
func New(cfg Config) (*Whisper, error) {
	// Ensure library is loaded
	if err := EnsureLibrary(); err != nil {
		return nil, fmt.Errorf("failed to ensure library: %w", err)
	}

	// Find FFmpeg if path not provided
	ffmpegPath := cfg.FFmpegPath
	if ffmpegPath == "" {
		path, err := audio.FindFFmpeg()
		if err != nil {
			return nil, fmt.Errorf("ffmpeg not found: %w", err)
		}
		ffmpegPath = path
	}

	// Create context params
	params := w.ContextDefaultParams()
	params.UseGPU = cfg.UseGPU
	params.FlashAttn = cfg.FlashAttn
	params.GPUDevice = cfg.GPUDevice

	// Initialize context
	modelPath := w.CString(cfg.ModelPath)
	ctx := w.InitFromFileWithParams(modelPath, params)
	if ctx == nil {
		return nil, errors.New("failed to initialize whisper context")
	}

	return &Whisper{
		ctx:        ctx,
		modelPath:  cfg.ModelPath,
		converter:  audio.NewConverter(ffmpegPath),
	}, nil
}

// Free releases the whisper context
func (wh *Whisper) Free() {
	if wh.ctx != nil {
		w.Free(wh.ctx)
		wh.ctx = nil
	}
}

// Transcribe transcribes an audio file to text
func (wh *Whisper) Transcribe(audioPath string, opts ...TranscribeOption) (*Result, error) {
	if wh.ctx == nil {
		return nil, errors.New("whisper context not initialized")
	}

	// Convert audio to PCM
	samples, err := wh.converter.ConvertToPCM16kHzMono(audioPath)
	if err != nil {
		return nil, fmt.Errorf("failed to convert audio: %w", err)
	}

	return wh.TranscribeSamples(samples, opts...)
}

// TranscribeSamples transcribes raw float32 samples
func (wh *Whisper) TranscribeSamples(samples []float32, opts ...TranscribeOption) (*Result, error) {
	if wh.ctx == nil {
		return nil, errors.New("whisper context not initialized")
	}

	// Apply options
	cfg := &transcribeConfig{
		language:   "auto",
		threads:    4,
		translate:  false,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	// Create full params
	params := w.FullDefaultParams(int32(w.SamplingGreedy))
	params.NThreads = int32(cfg.threads)
	params.Translate = cfg.translate
	params.NoTimestamps = false
	params.PrintProgress = false
	params.PrintRealtime = false
	params.PrintTimestamps = false

	// Set language
	if cfg.language != "auto" {
		langStr := w.CString(cfg.language)
		params.Language = langStr
	}

	// Run transcription
	ret := w.Full(wh.ctx, params, &samples[0], int32(len(samples)))
	if ret != 0 {
		return nil, fmt.Errorf("transcription failed with code: %d", ret)
	}

	// Extract results
	return wh.extractResult(), nil
}

// TranscribeWithContext transcribes with context support for cancellation
func (wh *Whisper) TranscribeWithContext(ctx context.Context, audioPath string, opts ...TranscribeOption) (*Result, error) {
	// For now, just call Transcribe
	// In full implementation, we'd check ctx.Done() periodically
	return wh.Transcribe(audioPath, opts...)
}

// extractResult extracts transcription result from whisper context
func (wh *Whisper) extractResult() *Result {
	nSegments := w.FullNSegments(wh.ctx)
	
	segments := make([]Segment, nSegments)
	var fullText strings.Builder

	for i := int32(0); i < nSegments; i++ {
		segIdx := int32(i)
		
		textPtr := w.FullGetSegmentText(wh.ctx, segIdx)
		text := w.GoString(textPtr)
		
		t0 := w.FullGetSegmentT0(wh.ctx, segIdx)
		t1 := w.FullGetSegmentT1(wh.ctx, segIdx)
		
		segments[i] = Segment{
			Text:  text,
			Start: time.Duration(t0) * time.Millisecond,
			End:   time.Duration(t1) * time.Millisecond,
		}
		
		fullText.WriteString(text)
	}

	// Get detected language
	langID := w.FullLangID(wh.ctx)
	langStr := w.GoString(w.LangStr(langID))

	return &Result{
		Text:     fullText.String(),
		Segments: segments,
		Language: langStr,
	}
}

// Result represents a transcription result
type Result struct {
	Text     string
	Segments []Segment
	Language string
}

// Segment represents a single transcription segment
type Segment struct {
	Text  string
	Start time.Duration
	End   time.Duration
}

// transcribeConfig holds transcription options
type transcribeConfig struct {
	language  string
	threads   int
	translate bool
}

// TranscribeOption is a functional option for transcription
type TranscribeOption func(*transcribeConfig)

// WithLanguage sets the language for transcription
func WithLanguage(lang string) TranscribeOption {
	return func(cfg *transcribeConfig) {
		cfg.language = lang
	}
}

// WithThreads sets the number of threads for transcription
func WithThreads(n int) TranscribeOption {
	return func(cfg *transcribeConfig) {
		cfg.threads = n
	}
}

// WithTranslate enables translation to English
func WithTranslate() TranscribeOption {
	return func(cfg *transcribeConfig) {
		cfg.translate = true
	}
}

// IsMultilingual returns true if the loaded model supports multiple languages
func (wh *Whisper) IsMultilingual() bool {
	if wh.ctx == nil {
		return false
	}
	return w.IsMultilingual(wh.ctx) != 0
}

// GetModelType returns the type of the loaded model
func (wh *Whisper) GetModelType() string {
	if wh.ctx == nil {
		return ""
	}
	return w.GoString(w.ModelTypeReadable(wh.ctx))
}

// GetSystemInfo returns system information from whisper
func GetSystemInfo() string {
	return w.GoString(w.PrintSystemInfo())
}

// GetVersion returns the whisper library version
func GetVersion() string {
	return w.GoString(w.Version())
}

// LangID returns the language ID for a language code
func LangID(lang string) int32 {
	return w.LangID(w.CString(lang))
}

// LangStr returns the language code for a language ID
func LangStr(id int) string {
	return w.GoString(w.LangStr(int32(id)))
}

// LangMaxID returns the maximum language ID
func LangMaxID() int32 {
	return w.LangMaxID()
}
