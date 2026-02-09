package remote

import (
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
	if gen.cloudflare == nil {
		t.Error("cloudflare generator not initialized")
	}
}

func TestGetAvailableModels(t *testing.T) {
	gen := NewGenerator()
	models := gen.GetAvailableModels()

	// Models available depend on API keys being configured
	// Test just checks the function doesn't panic
	_ = models
}

func TestIsAvailable(t *testing.T) {
	gen := NewGenerator()
	// Availability depends on API keys being configured
	// Test just checks the function doesn't panic
	_ = gen.IsAvailable()
}

func TestGetGenerator(t *testing.T) {
	gen := NewGenerator()

	tests := []struct {
		name      string
		genName   string
		wantError bool
	}{
		{
			name:      "gemini generator",
			genName:   "gemini",
			wantError: false, // May not be available but shouldn't panic
		},
		{
			name:      "cloudflare generator",
			genName:   "cloudflare",
			wantError: false, // May not be available but shouldn't panic
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
			if tt.wantError && err == nil {
				t.Error("expected error for unknown generator")
			}
			if !tt.wantError && err != nil && !strings.Contains(err.Error(), "not available") {
				t.Errorf("unexpected error: %v", err)
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
