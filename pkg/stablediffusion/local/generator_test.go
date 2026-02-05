package local

import (
	"context"
	"os"
	"path/filepath"
	"testing"
)

// MockExecutor for testing
type MockExecutor struct {
	runCalled bool
	lastArgs  []string
	err       error
}

func (m *MockExecutor) Run(ctx context.Context, name string, args ...string) error {
	m.runCalled = true
	m.lastArgs = args
	return m.err
}

func TestNewGenerator(t *testing.T) {
	gen := NewGenerator("/path/to/sd", "/path/to/models")
	if gen == nil {
		t.Fatal("NewGenerator returned nil")
	}
	if gen.binaryPath != "/path/to/sd" {
		t.Errorf("binaryPath = %s, want /path/to/sd", gen.binaryPath)
	}
	if gen.modelsPath != "/path/to/models" {
		t.Errorf("modelsPath = %s, want /path/to/models", gen.modelsPath)
	}
	if gen.executor == nil {
		t.Error("executor not initialized")
	}
}

func TestNewGeneratorWithExecutor(t *testing.T) {
	mockExec := &MockExecutor{}
	gen := NewGeneratorWithExecutor("/path/to/sd", "/path/to/models", mockExec)

	if gen == nil {
		t.Fatal("NewGeneratorWithExecutor returned nil")
	}
	if gen.executor != mockExec {
		t.Error("custom executor not set")
	}
}

func TestGetBinaryPath(t *testing.T) {
	gen := NewGenerator("/test/sd", "/test/models")
	if gen.GetBinaryPath() != "/test/sd" {
		t.Errorf("GetBinaryPath() = %s, want /test/sd", gen.GetBinaryPath())
	}
}

func TestGetModelsPath(t *testing.T) {
	gen := NewGenerator("/test/sd", "/test/models")
	if gen.GetModelsPath() != "/test/models" {
		t.Errorf("GetModelsPath() = %s, want /test/models", gen.GetModelsPath())
	}
}

func TestApplyDefaults(t *testing.T) {
	gen := NewGenerator("/test/sd", "/test/models")

	tests := []struct {
		name     string
		input    GenerationOptions
		expected GenerationOptions
	}{
		{
			name:  "empty options",
			input: GenerationOptions{},
			expected: GenerationOptions{
				Width:  1024,
				Height: 1024,
				Steps:  20,
				Cfg:    7.0,
			},
		},
		{
			name: "partial options",
			input: GenerationOptions{
				Width: 512,
				Steps: 30,
			},
			expected: GenerationOptions{
				Width:  512,
				Height: 1024,
				Steps:  30,
				Cfg:    7.0,
			},
		},
		{
			name: "all options set",
			input: GenerationOptions{
				Width:  800,
				Height: 600,
				Steps:  25,
				Cfg:    8.5,
			},
			expected: GenerationOptions{
				Width:  800,
				Height: 600,
				Steps:  25,
				Cfg:    8.5,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := gen.applyDefaults(tt.input)
			if result.Width != tt.expected.Width {
				t.Errorf("Width = %d, want %d", result.Width, tt.expected.Width)
			}
			if result.Height != tt.expected.Height {
				t.Errorf("Height = %d, want %d", result.Height, tt.expected.Height)
			}
			if result.Steps != tt.expected.Steps {
				t.Errorf("Steps = %d, want %d", result.Steps, tt.expected.Steps)
			}
			if result.Cfg != tt.expected.Cfg {
				t.Errorf("Cfg = %f, want %f", result.Cfg, tt.expected.Cfg)
			}
		})
	}
}

func TestBuildArgs(t *testing.T) {
	gen := NewGenerator("/test/sd", "/test/models")

	seed := int64(42)
	imageUrl := "/path/to/input.png"

	opts := GenerationOptions{
		Prompt:         "test prompt",
		NegativePrompt: "bad quality",
		ModelPath:      "/path/to/model.gguf",
		OutputPath:     "/path/to/output.png",
		Width:          512,
		Height:         512,
		Steps:          20,
		Cfg:            7.5,
		Seed:           &seed,
		SamplerName:    "euler_a",
		Scheduler:      "karras",
		Strength:       0.75,
		ImageUrl:       &imageUrl,
	}

	args := gen.buildArgs(opts)

	// Check required args
	expectedArgs := map[string]bool{
		"-m":                true,
		"-p":                true,
		"-o":                true,
		"--width":           true,
		"--height":          true,
		"--steps":           true,
		"--cfg-scale":       true,
		"--seed":            true,
		"-n":                true,
		"--sampling-method": true,
		"--schedule":        true,
		"--strength":        true,
		"-i":                true,
	}

	for i := 0; i < len(args); i++ {
		if expectedArgs[args[i]] {
			delete(expectedArgs, args[i])
		}
	}

	if len(expectedArgs) > 0 {
		t.Errorf("Missing expected args: %v", expectedArgs)
	}
}

func TestIsSupportedModelFormat(t *testing.T) {
	gen := NewGenerator("/test/sd", "/test/models")

	tests := []struct {
		filename string
		expected bool
	}{
		{"model.ckpt", true},
		{"model.safetensors", true},
		{"model.pt", true},
		{"model.bin", true},
		{"model.gguf", true},
		{"model.txt", false},
		{"model.json", false},
		{"model", false},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			result := gen.isSupportedModelFormat(tt.filename)
			if result != tt.expected {
				t.Errorf("isSupportedModelFormat(%s) = %v, want %v", tt.filename, result, tt.expected)
			}
		})
	}
}

func TestGetFirstAvailableModel(t *testing.T) {
	// Create temp directory with test models
	tmpDir := t.TempDir()

	// Create test files
	testFiles := []string{
		"readme.txt",
		"model1.gguf",
		"model2.safetensors",
	}

	for _, file := range testFiles {
		f, err := os.Create(filepath.Join(tmpDir, file))
		if err != nil {
			t.Fatal(err)
		}
		f.Close()
	}

	gen := NewGenerator("/test/sd", tmpDir)
	model := gen.GetFirstAvailableModel()

	if model == "" {
		t.Error("GetFirstAvailableModel returned empty string")
	}

	// Should return one of the model files
	if !filepath.IsAbs(model) {
		t.Error("returned path is not absolute")
	}
}

func TestIsAvailable(t *testing.T) {
	// Test with non-existent binary
	gen := NewGenerator("/nonexistent/sd", "/test/models")
	if gen.IsAvailable() {
		t.Error("IsAvailable() = true for non-existent binary")
	}

	// Test with existing file
	tmpFile := filepath.Join(t.TempDir(), "sd")
	f, err := os.Create(tmpFile)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	gen2 := NewGenerator(tmpFile, "/test/models")
	if !gen2.IsAvailable() {
		t.Error("IsAvailable() = false for existing binary")
	}
}

func TestGenerate_BinaryNotFound(t *testing.T) {
	gen := NewGenerator("/nonexistent/sd", "/test/models")

	opts := GenerationOptions{
		Prompt:     "test",
		ModelPath:  "/test/model.gguf",
		OutputPath: "/test/output.png",
	}

	err := gen.Generate(context.Background(), opts)
	if err == nil {
		t.Error("expected error for non-existent binary")
	}
}

func TestGenerate_WithMockExecutor(t *testing.T) {
	// Create temp binary file
	tmpDir := t.TempDir()
	binaryPath := filepath.Join(tmpDir, "sd")
	f, err := os.Create(binaryPath)
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	// Create temp output file (simulating successful generation)
	outputPath := filepath.Join(tmpDir, "output.png")

	mockExec := &MockExecutor{}
	gen := NewGeneratorWithExecutor(binaryPath, tmpDir, mockExec)

	// Mock executor should create output file
	mockExec.err = nil

	opts := GenerationOptions{
		Prompt:     "test prompt",
		ModelPath:  "/test/model.gguf",
		OutputPath: outputPath,
	}

	// Create output file to simulate successful generation
	outFile, err := os.Create(outputPath)
	if err != nil {
		t.Fatal(err)
	}
	outFile.Close()

	err = gen.Generate(context.Background(), opts)
	if err != nil {
		t.Errorf("Generate() error = %v", err)
	}

	if !mockExec.runCalled {
		t.Error("executor.Run() was not called")
	}
}

func TestGetBinaryName(t *testing.T) {
	name := GetBinaryName()
	if name == "" {
		t.Error("GetBinaryName() returned empty string")
	}
	// Should return "sd" or "sd.exe" depending on platform
	if name != "sd" && name != "sd.exe" {
		t.Errorf("GetBinaryName() = %s, want 'sd' or 'sd.exe'", name)
	}
}
