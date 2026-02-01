package audio

import (
	"encoding/binary"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Converter handles audio conversion using FFmpeg binary
type Converter struct {
	ffmpegPath string
}

// NewConverter creates a new audio converter
func NewConverter(ffmpegPath string) *Converter {
	return &Converter{ffmpegPath: ffmpegPath}
}

// FindFFmpeg attempts to find FFmpeg binary in PATH or common locations
func FindFFmpeg() (string, error) {
	// Try PATH first
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return path, nil
	}

	// Try common locations
	commonPaths := []string{
		"/usr/local/bin/ffmpeg",
		"/usr/bin/ffmpeg",
		"/opt/homebrew/bin/ffmpeg", // macOS Homebrew on Apple Silicon
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("ffmpeg not found in PATH or common locations")
}

// ConvertToPCM16kHzMono converts any audio file to 16kHz mono float32 PCM
// This is the format required by whisper.cpp
func (c *Converter) ConvertToPCM16kHzMono(inputPath string) ([]float32, error) {
	if c.ffmpegPath == "" {
		return nil, fmt.Errorf("ffmpeg path not set")
	}

	// Create temp file for raw PCM output
	tempFile := filepath.Join(os.TempDir(), fmt.Sprintf("whisper_%d.raw", os.Getpid()))
	defer os.Remove(tempFile)

	// Run FFmpeg: convert to 16kHz, mono, 16-bit signed little-endian
	// whisper.cpp expects float32 samples, so we convert from int16
	cmd := exec.Command(c.ffmpegPath,
		"-i", inputPath,
		"-ar", "16000",      // Sample rate: 16kHz (whisper requirement)
		"-ac", "1",          // Channels: mono
		"-f", "s16le",       // Format: 16-bit signed little-endian PCM
		"-y",                // Overwrite output
		tempFile,
	)

	// Capture stderr for error reporting
	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("ffmpeg failed: %w\nOutput: %s", err, string(output))
	}

	// Read the raw PCM data
	data, err := os.ReadFile(tempFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read converted audio: %w", err)
	}

	// Convert int16 samples to float32
	// whisper.cpp expects samples in range [-1.0, 1.0]
	if len(data)%2 != 0 {
		return nil, fmt.Errorf("invalid PCM data length")
	}

	samples := make([]float32, len(data)/2)
	for i := 0; i < len(samples); i++ {
		// Read int16 (little-endian)
		val := int16(binary.LittleEndian.Uint16(data[i*2 : i*2+2]))
		// Convert to float32 in range [-1.0, 1.0]
		samples[i] = float32(val) / 32768.0
	}

	return samples, nil
}

// ConvertToPCM16kHzMonoFromReader converts audio from an io.Reader
// Note: FFmpeg needs a file, so we write to temp file first
func (c *Converter) ConvertToPCM16kHzMonoFromReader(inputData []byte, ext string) ([]float32, error) {
	// Write input to temp file
	tempInput := filepath.Join(os.TempDir(), fmt.Sprintf("whisper_input_%d%s", os.Getpid(), ext))
	defer os.Remove(tempInput)

	if err := os.WriteFile(tempInput, inputData, 0644); err != nil {
		return nil, fmt.Errorf("failed to write temp input: %w", err)
	}

	return c.ConvertToPCM16kHzMono(tempInput)
}

// GetAudioDuration returns the duration of an audio file in seconds
func (c *Converter) GetAudioDuration(inputPath string) (float64, error) {
	if c.ffmpegPath == "" {
		return 0, fmt.Errorf("ffmpeg path not set")
	}

	cmd := exec.Command(c.ffmpegPath,
		"-i", inputPath,
		"-f", "null",
		"-",
	)

	output, _ := cmd.CombinedOutput()
	
	// Parse duration from output
	// Duration: 00:05:30.50
	outputStr := string(output)
	
	// Try to find duration in output
	var hours, minutes, seconds float64
	_, err := fmt.Sscanf(outputStr, "Duration: %f:%f:%f", &hours, &minutes, &seconds)
	if err != nil {
		return 0, fmt.Errorf("failed to parse duration")
	}

	duration := hours*3600 + minutes*60 + seconds
	return duration, nil
}

// SupportedFormats returns list of audio formats supported by FFmpeg
func (c *Converter) SupportedFormats() []string {
	return []string{
		".mp3", ".wav", ".flac", ".aac", ".ogg", ".m4a", ".wma",
		".mp4", ".avi", ".mov", ".mkv", // Video formats (audio extracted)
	}
}
