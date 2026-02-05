package remote

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestNewGenerator(t *testing.T) {
	gen := NewGenerator()
	if gen == nil {
		t.Fatal("NewGenerator returned nil")
	}
	if gen.gemini == nil {
		t.Error("gemini generator not initialized")
	}
	if gen.pollinations == nil {
		t.Error("pollinations generator not initialized")
	}
}

func TestGetAvailableModels(t *testing.T) {
	gen := NewGenerator()
	models := gen.GetAvailableModels()

	if len(models) == 0 {
		t.Error("no models returned")
	}

	// Should have at least Pollinations models
	hasFlux := false
	for _, model := range models {
		if model == "flux" {
			hasFlux = true
			break
		}
	}
	if !hasFlux {
		t.Error("expected flux model in available models")
	}
}

func TestIsAvailable(t *testing.T) {
	gen := NewGenerator()
	// Should always be available (Pollinations doesn't require API key)
	if !gen.IsAvailable() {
		t.Error("generator should be available")
	}
}

func TestGetGenerator(t *testing.T) {
	gen := NewGenerator()

	tests := []struct {
		name      string
		genName   string
		wantError bool
	}{
		{
			name:      "pollinations generator",
			genName:   "pollinations",
			wantError: false,
		},
		{
			name:      "unknown generator",
			genName:   "unknown",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := gen.GetGenerator(tt.genName)
			if (err != nil) != tt.wantError {
				t.Errorf("GetGenerator() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

func TestGenerationOptions(t *testing.T) {
	opts := GenerationOptions{
		Prompt:     "test prompt",
		OutputPath: "/tmp/test.png",
		Width:      1024,
		Height:     1024,
	}

	if opts.Prompt != "test prompt" {
		t.Error("prompt not set correctly")
	}
	if opts.Width != 1024 {
		t.Error("width not set correctly")
	}
}

func TestCalculateAspectRatio(t *testing.T) {
	tests := []struct {
		name     string
		opts     GenerationOptions
		expected string
	}{
		{
			name: "square",
			opts: GenerationOptions{
				Width:  1024,
				Height: 1024,
			},
			expected: "1:1",
		},
		{
			name: "widescreen",
			opts: GenerationOptions{
				Width:  1920,
				Height: 1080,
			},
			expected: "16:9",
		},
		{
			name: "portrait",
			opts: GenerationOptions{
				Width:  1080,
				Height: 1920,
			},
			expected: "9:16",
		},
		{
			name: "explicit aspect ratio",
			opts: GenerationOptions{
				AspectRatio: "4:3",
			},
			expected: "4:3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateAspectRatio(tt.opts)
			if result != tt.expected {
				t.Errorf("calculateAspectRatio() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGeminiGenerator(t *testing.T) {
	gen := NewGeminiGenerator()
	if gen == nil {
		t.Fatal("NewGeminiGenerator returned nil")
	}

	models := gen.GetAvailableModels()
	if len(models) == 0 {
		t.Error("no Gemini models returned")
	}
}

func TestPollinationsGenerator(t *testing.T) {
	gen := NewPollinationsGenerator()
	if gen == nil {
		t.Fatal("NewPollinationsGenerator returned nil")
	}

	if !gen.IsAvailable() {
		t.Error("Pollinations should always be available")
	}

	models := gen.GetAvailableModels()
	if len(models) == 0 {
		t.Error("no Pollinations models returned")
	}
}

// mockTransport is a mock HTTP transport for testing
type mockTransport struct {
	response *http.Response
	err      error
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.response, nil
}

// TestGenerateWithInvalidPath tests error handling for invalid output path
func TestGenerateWithInvalidPath(t *testing.T) {
	// Mock HTTP transport to avoid real network calls
	oldTransport := http.DefaultClient.Transport
	http.DefaultClient.Transport = &mockTransport{
		response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewReader([]byte("fake image data"))),
			Header:     make(http.Header),
		},
	}
	defer func() {
		http.DefaultClient.Transport = oldTransport
	}()

	gen := NewPollinationsGenerator()

	opts := GenerationOptions{
		Prompt:     "test",
		OutputPath: "/invalid/path/that/does/not/exist/image.png",
		Width:      512,
		Height:     512,
	}

	ctx := context.Background()
	err := gen.Generate(ctx, opts)

	// Should fail because path doesn't exist (not because of network)
	if err == nil {
		t.Error("expected error for invalid path, got nil")
	}
	if err != nil && !strings.Contains(err.Error(), "failed to create output file") {
		t.Errorf("expected 'failed to create output file' error, got: %v", err)
	}
}
