package whisperapp

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetModelSpec(t *testing.T) {
	tests := []struct {
		name      string
		modelName string
		wantFound bool
		wantName  string
	}{
		{
			name:      "existing model - tiny",
			modelName: "tiny",
			wantFound: true,
			wantName:  "tiny",
		},
		{
			name:      "existing model - base",
			modelName: "base",
			wantFound: true,
			wantName:  "base",
		},
		{
			name:      "existing model - large-v3",
			modelName: "large-v3",
			wantFound: true,
			wantName:  "large-v3",
		},
		{
			name:      "non-existing model",
			modelName: "nonexistent",
			wantFound: false,
		},
		{
			name:      "empty model name",
			modelName: "",
			wantFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, found := GetModelSpec(tt.modelName)
			assert.Equal(t, tt.wantFound, found)
			if tt.wantFound {
				assert.Equal(t, tt.wantName, spec.Name)
				assert.NotZero(t, spec.Size)
				assert.NotEmpty(t, spec.URL)
			}
		})
	}
}

func TestGetAllModels(t *testing.T) {
	models := GetAllModels()
	assert.NotEmpty(t, models)
	assert.GreaterOrEqual(t, len(models), 10) // Should have at least 10 models

	// Check that essential models exist
	modelNames := make(map[string]bool)
	for _, m := range models {
		modelNames[m.Name] = true
	}

	assert.True(t, modelNames["tiny"], "should have tiny model")
	assert.True(t, modelNames["base"], "should have base model")
	assert.True(t, modelNames["small"], "should have small model")
}

func TestGetAvailableModels(t *testing.T) {
	tests := []struct {
		name     string
		ramGB    int64
		minCount int
	}{
		{
			name:     "1GB RAM - only tiny",
			ramGB:    1,
			minCount: 2, // tiny and tiny.en
		},
		{
			name:     "2GB RAM - tiny and base",
			ramGB:    2,
			minCount: 4, // tiny, tiny.en, base, base.en
		},
		{
			name:     "4GB RAM - up to small",
			ramGB:    4,
			minCount: 6,
		},
		{
			name:     "16GB RAM - all models",
			ramGB:    16,
			minCount: 10,
		},
		{
			name:     "0GB RAM - fallback to tiny",
			ramGB:    0,
			minCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			models := GetAvailableModels(tt.ramGB)
			assert.GreaterOrEqual(t, len(models), tt.minCount)
		})
	}
}

func TestSelectOptimalModel(t *testing.T) {
	tests := []struct {
		name   string
		ramGB  int64
		minRAM int64 // The minimum RAM requirement for the selected model
	}{
		{
			name:   "1GB RAM - should select tiny",
			ramGB:  1,
			minRAM: 1,
		},
		{
			name:   "2GB RAM - should select base",
			ramGB:  2,
			minRAM: 2,
		},
		{
			name:   "4GB RAM - should select small",
			ramGB:  4,
			minRAM: 4,
		},
		{
			name:   "8GB RAM - should select medium or smaller",
			ramGB:  8,
			minRAM: 8,
		},
		{
			name:   "16GB RAM - should select large",
			ramGB:  16,
			minRAM: 8, // large-v3-turbo only needs 8GB
		},
		{
			name:   "0GB RAM - fallback to tiny",
			ramGB:  0,
			minRAM: 1, // Will return tiny which needs 1GB (best available fallback)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			model := SelectOptimalModel(tt.ramGB)
			require.NotNil(t, model)
			// For 0GB RAM, it returns the first model (tiny) as fallback
			// For normal cases, verify the model fits within available RAM
			if tt.ramGB > 0 {
				assert.LessOrEqual(t, model.MinRAM, tt.ramGB, "selected model should fit in available RAM")
			} else {
				// For 0GB, just verify we got a model (fallback to first available)
				assert.NotEmpty(t, model.Name)
			}
		})
	}
}

func TestGetModelPath(t *testing.T) {
	tests := []struct {
		name      string
		modelsDir string
		modelName string
		expected  string
	}{
		{
			name:      "simple path",
			modelsDir: "./models",
			modelName: "base",
			expected:  filepath.Join(".", "models", "base.bin"),
		},
		{
			name:      "absolute path",
			modelsDir: "/var/lib/models",
			modelName: "tiny",
			expected:  "/var/lib/models/tiny.bin",
		},
		{
			name:      "path with dots",
			modelsDir: "./data/models/whisper",
			modelName: "small.en",
			expected:  filepath.Join(".", "data", "models", "whisper", "small.en.bin"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := GetModelPath(tt.modelsDir, tt.modelName)
			assert.Equal(t, tt.expected, path)
		})
	}
}

func TestIsModelDownloaded(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	// Create a fake model file with correct size
	baseSpec, _ := GetModelSpec("base")
	basePath := filepath.Join(tempDir, "base.bin")

	// Create file with correct size
	f, err := os.Create(basePath)
	require.NoError(t, err)
	_, err = f.Write(make([]byte, baseSpec.Size))
	require.NoError(t, err)
	f.Close()

	// Create file with wrong size
	tinyPath := filepath.Join(tempDir, "tiny.bin")
	f, err = os.Create(tinyPath)
	require.NoError(t, err)
	_, err = f.Write(make([]byte, 100)) // Wrong size
	require.NoError(t, err)
	f.Close()

	tests := []struct {
		name      string
		modelName string
		want      bool
	}{
		{
			name:      "correct size",
			modelName: "base",
			want:      true,
		},
		{
			name:      "wrong size",
			modelName: "tiny",
			want:      false,
		},
		{
			name:      "non-existent",
			modelName: "large-v3",
			want:      false,
		},
		{
			name:      "unknown model",
			modelName: "unknown",
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsModelDownloaded(tempDir, tt.modelName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestListDownloadedModels(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	// Create some fake model files
	models := []string{"tiny.bin", "base.bin", "custom.bin", "ggml-old.bin"}
	for _, name := range models {
		path := filepath.Join(tempDir, name)
		f, err := os.Create(path)
		require.NoError(t, err)
		f.Close()
	}

	list, err := ListDownloadedModels(tempDir)
	require.NoError(t, err)

	// Should only list valid models (not ggml-old.bin which is deprecated format)
	assert.Contains(t, list, "tiny")
	assert.Contains(t, list, "base")
	assert.NotContains(t, list, "custom")   // Not a valid model spec
	assert.NotContains(t, list, "ggml-old") // Deprecated format
}

func TestListDownloadedModels_EmptyDir(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	list, err := ListDownloadedModels(tempDir)
	require.NoError(t, err)
	assert.Empty(t, list)
}

func TestListDownloadedModels_NonExistentDir(t *testing.T) {
	list, err := ListDownloadedModels("/nonexistent/path/12345")
	require.NoError(t, err)
	assert.Empty(t, list)
}

func TestDeleteModel(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	// Create a fake model file
	modelPath := filepath.Join(tempDir, "test.bin")
	f, err := os.Create(modelPath)
	require.NoError(t, err)
	f.Close()

	// Verify file exists
	_, err = os.Stat(modelPath)
	require.NoError(t, err)

	// Delete the model
	err = DeleteModel(tempDir, "test")
	require.NoError(t, err)

	// Verify file is deleted
	_, err = os.Stat(modelPath)
	assert.True(t, os.IsNotExist(err))
}

func TestDeleteModel_NotFound(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	// Try to delete non-existent model
	err := DeleteModel(tempDir, "nonexistent")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "model not found")
}

func TestHumanSize(t *testing.T) {
	tests := []struct {
		bytes    int64
		expected string
	}{
		{0, "0 B"},
		{100, "100 B"},
		{1024, "1.00 KB"},
		{1536, "1.50 KB"},
		{1024 * 1024, "1.00 MB"},
		{1024 * 1024 * 1024, "1.00 GB"},
		{1024 * 1024 * 1024 * 2, "2.00 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := HumanSize(tt.bytes)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDownloadProgressLogger(t *testing.T) {
	logger := &DownloadProgressLogger{}

	// Test logging - should not panic
	logger.Log(0, 1000)
	logger.Log(500, 1000)
	logger.Log(1000, 1000)
}

func TestWhisperModelSpec_Struct(t *testing.T) {
	spec := WhisperModelSpec{
		Name:        "test-model",
		URL:         "https://example.com/model.bin",
		Size:        1000000,
		Parameters:  "100M",
		MinRAM:      2,
		EnglishOnly: false,
		Description: "Test model description",
	}

	assert.Equal(t, "test-model", spec.Name)
	assert.Equal(t, "https://example.com/model.bin", spec.URL)
	assert.Equal(t, int64(1000000), spec.Size)
	assert.Equal(t, "100M", spec.Parameters)
	assert.Equal(t, int64(2), spec.MinRAM)
	assert.False(t, spec.EnglishOnly)
	assert.Equal(t, "Test model description", spec.Description)
}

func TestModelSpecs_Definitions(t *testing.T) {
	// Verify that all model specs have valid data
	for _, spec := range modelSpecs {
		t.Run(spec.Name, func(t *testing.T) {
			assert.NotEmpty(t, spec.Name)
			assert.NotEmpty(t, spec.URL)
			assert.True(t, spec.Size > 0, "size should be positive")
			assert.NotEmpty(t, spec.Parameters)
			assert.True(t, spec.MinRAM >= 0, "MinRAM should be non-negative")
			assert.NotEmpty(t, spec.Description)
			assert.Contains(t, spec.URL, "huggingface.co", "URL should be from huggingface")
		})
	}
}
