package model

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetModelSpec(t *testing.T) {
	tests := []struct {
		name     string
		wantName string
		exists   bool
	}{
		{"base", "base", true},
		{"base.en", "base.en", true},
		{"tiny", "tiny", true},
		{"large-v3", "large-v3", true},
		{"nonexistent", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			spec, exists := GetModelSpec(tt.name)
			if exists != tt.exists {
				t.Errorf("GetModelSpec(%q) exists = %v, want %v", tt.name, exists, tt.exists)
			}
			if exists && spec.Name != tt.wantName {
				t.Errorf("GetModelSpec(%q) name = %q, want %q", tt.name, spec.Name, tt.wantName)
			}
		})
	}
}

func TestSelectOptimalModel(t *testing.T) {
	tests := []struct {
		ram      int64
		wantName string
	}{
		{1, "tiny.en"},      // 1GB RAM -> tiny.en (last model that fits)
		{2, "base.en"},      // 2GB RAM -> base.en (last model that fits)
		{4, "small.en"},     // 4GB RAM -> small.en (last model that fits)
		{8, "large-v3-turbo"}, // 8GB RAM -> large-v3-turbo (last model that fits, better than medium.en)
		{16, "large-v3-turbo"}, // 16GB RAM -> large-v3-turbo (last model that fits)
		{32, "large-v3-turbo"}, // 32GB RAM -> large-v3-turbo (last model that fits)
		{0, "tiny"},         // 0GB RAM -> tiny (fallback)
	}

	for _, tt := range tests {
		t.Run(tt.wantName, func(t *testing.T) {
			model := SelectOptimalModel(tt.ram)
			if model.Name != tt.wantName {
				t.Errorf("SelectOptimalModel(%d) = %q, want %q", tt.ram, model.Name, tt.wantName)
			}
		})
	}
}

func TestSelectOptimalModelEnglish(t *testing.T) {
	model := SelectOptimalModelEnglish(8)
	if model == nil {
		t.Fatal("SelectOptimalModelEnglish returned nil")
	}
	if !model.EnglishOnly {
		t.Errorf("Expected English-only model, got %q", model.Name)
	}
}

func TestGetAvailableModels(t *testing.T) {
	models := GetAvailableModels(4)
	if len(models) == 0 {
		t.Error("GetAvailableModels returned empty slice")
	}

	// All returned models should require <= 4GB RAM
	for _, m := range models {
		if m.MinRAM > 4 {
			t.Errorf("Model %q requires %dGB RAM, but only 4GB available", m.Name, m.MinRAM)
		}
	}
}

func TestHumanSize(t *testing.T) {
	tests := []struct {
		bytes int64
		want  string
	}{
		{100, "100 B"},
		{1024, "1.00 KB"},
		{1024 * 1024, "1.00 MB"},
		{77691776, "74.09 MB"},
		{1533121536, "1.43 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := HumanSize(tt.bytes)
			if got != tt.want {
				t.Errorf("HumanSize(%d) = %q, want %q", tt.bytes, got, tt.want)
			}
		})
	}
}

func TestGetModelPath(t *testing.T) {
	path := GetModelPath("/models", "base")
	expected := filepath.Join("/models", "ggml-base.bin")
	if path != expected {
		t.Errorf("GetModelPath = %q, want %q", path, expected)
	}
}

func TestIsModelDownloaded(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	// Should return false for non-existent model
	if IsModelDownloaded(tempDir, "base") {
		t.Error("IsModelDownloaded returned true for non-existent model")
	}

	// Create a dummy model file
	modelPath := filepath.Join(tempDir, "ggml-base.bin")
	if err := os.WriteFile(modelPath, []byte("dummy"), 0644); err != nil {
		t.Fatalf("Failed to create dummy model: %v", err)
	}

	// Should return true now
	if !IsModelDownloaded(tempDir, "base") {
		t.Error("IsModelDownloaded returned false for existing model")
	}
}

func TestListDownloadedModels(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	// Should return empty slice for empty directory
	models, err := ListDownloadedModels(tempDir)
	if err != nil {
		t.Fatalf("ListDownloadedModels failed: %v", err)
	}
	if len(models) != 0 {
		t.Errorf("Expected 0 models, got %d", len(models))
	}

	// Create some dummy model files
	files := []string{"ggml-base.bin", "ggml-small.bin", "ggml-tiny.en.bin"}
	for _, f := range files {
		path := filepath.Join(tempDir, f)
		if err := os.WriteFile(path, []byte("dummy"), 0644); err != nil {
			t.Fatalf("Failed to create %s: %v", f, err)
		}
	}

	// Should list the models
	models, err = ListDownloadedModels(tempDir)
	if err != nil {
		t.Fatalf("ListDownloadedModels failed: %v", err)
	}
	if len(models) != 3 {
		t.Errorf("Expected 3 models, got %d", len(models))
	}
}

func TestDeleteModel(t *testing.T) {
	// Create temp directory
	tempDir := t.TempDir()

	// Create a dummy model file
	modelPath := filepath.Join(tempDir, "ggml-base.bin")
	if err := os.WriteFile(modelPath, []byte("dummy"), 0644); err != nil {
		t.Fatalf("Failed to create dummy model: %v", err)
	}

	// Delete the model
	if err := DeleteModel(tempDir, "base"); err != nil {
		t.Fatalf("DeleteModel failed: %v", err)
	}

	// Should be deleted
	if _, err := os.Stat(modelPath); !os.IsNotExist(err) {
		t.Error("Model file still exists after deletion")
	}

	// Deleting non-existent model should return error
	if err := DeleteModel(tempDir, "base"); err == nil {
		t.Error("DeleteModel should return error for non-existent model")
	}
}
