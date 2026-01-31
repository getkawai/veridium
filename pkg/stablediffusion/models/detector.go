package models

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Detector automatically detects model types from files
type Detector struct {
	rules []DetectionRule
}

// DetectionRule defines a rule for detecting model types
type DetectionRule struct {
	Type       ModelType
	Extensions []string
	Patterns   []string // Filename patterns to match
	Format     ModelFormat
}

// NewDetector creates a new model detector with default rules
func NewDetector() *Detector {
	return &Detector{
		rules: []DetectionRule{
			{
				Type:       ModelTypeDiffusion,
				Extensions: []string{".gguf", ".safetensors", ".ckpt", ".pt"},
				Patterns:   []string{"diffusion", "unet", "model"},
				Format:     FormatUnknown,
			},
			{
				Type:       ModelTypeVAE,
				Extensions: []string{".safetensors", ".pt", ".ckpt"},
				Patterns:   []string{"vae", "ae", "autoencoder"},
				Format:     FormatUnknown,
			},
			{
				Type:       ModelTypeLLM,
				Extensions: []string{".gguf", ".bin"},
				Patterns:   []string{"llm", "lm", "gpt", "qwen", "mistral"},
				Format:     FormatUnknown,
			},
			{
				Type:       ModelTypeT5XXL,
				Extensions: []string{".gguf", ".bin", ".safetensors"},
				Patterns:   []string{"t5", "t5xxl", "umt5"},
				Format:     FormatUnknown,
			},
			{
				Type:       ModelTypeCLIP,
				Extensions: []string{".safetensors", ".pt", ".bin"},
				Patterns:   []string{"clip", "clipl", "clipg", "text_encoder"},
				Format:     FormatUnknown,
			},
			{
				Type:       ModelTypeEmbedding,
				Extensions: []string{".safetensors", ".pt", ".bin"},
				Patterns:   []string{"embedding", "embed"},
				Format:     FormatUnknown,
			},
			{
				Type:       ModelTypeLoRA,
				Extensions: []string{".safetensors", ".pt", ".ckpt"},
				Patterns:   []string{"lora"},
				Format:     FormatUnknown,
			},
			{
				Type:       ModelTypeControlNet,
				Extensions: []string{".safetensors", ".pt", ".ckpt", ".pth"},
				Patterns:   []string{"control", "controlnet", "canny", "depth", "openpose"},
				Format:     FormatUnknown,
			},
			{
				Type:       ModelTypeUpscaler,
				Extensions: []string{".safetensors", ".pth"},
				Patterns:   []string{"upscale", "esrgan", "real-esrgan", "4x", "2x"},
				Format:     FormatUnknown,
			},
		},
	}
}

// Detect attempts to detect model info from a file path
func (d *Detector) Detect(path string) (*ModelInfo, error) {
	// Get file info
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to stat file: %w", err)
	}

	if info.IsDir() {
		return nil, fmt.Errorf("path is a directory: %s", path)
	}

	filename := filepath.Base(path)
	ext := strings.ToLower(filepath.Ext(filename))
	nameWithoutExt := strings.TrimSuffix(filename, ext)

	// Detect format
	format := detectFormat(ext)

	// Detect type based on rules
	modelType := d.detectType(filename, ext)

	// Extract tags from filename
	tags := extractTags(nameWithoutExt)

	return &ModelInfo{
		ID:          generateModelIDFromPath(path),
		Name:        nameWithoutExt,
		Type:        modelType,
		Path:        path,
		Format:      format,
		Size:        info.Size(),
		Tags:        tags,
		Metadata:    make(map[string]string),
		AddedAt:     info.ModTime(),
		Description: "",
	}, nil
}

// DetectDirectory scans a directory and detects all models
func (d *Detector) DetectDirectory(dir string, recursive bool) ([]*ModelInfo, error) {
	models := make([]*ModelInfo, 0)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if !recursive && path != dir {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file has a model extension
		ext := strings.ToLower(filepath.Ext(path))
		if !isModelExtension(ext) {
			return nil
		}

		model, err := d.Detect(path)
		if err != nil {
			// Log error but continue scanning
			return nil
		}

		models = append(models, model)
		return nil
	}

	if err := filepath.Walk(dir, walkFn); err != nil {
		return nil, fmt.Errorf("failed to scan directory: %w", err)
	}

	return models, nil
}

// detectType determines the model type based on filename and extension
func (d *Detector) detectType(filename, ext string) ModelType {
	filenameLower := strings.ToLower(filename)

	for _, rule := range d.rules {
		// Check extension
		hasExt := false
		for _, e := range rule.Extensions {
			if ext == e {
				hasExt = true
				break
			}
		}

		if !hasExt {
			continue
		}

		// Check patterns
		for _, pattern := range rule.Patterns {
			if strings.Contains(filenameLower, pattern) {
				return rule.Type
			}
		}
	}

	// Default type based on extension
	switch ext {
	case ".gguf":
		return ModelTypeDiffusion // Most likely
	case ".safetensors":
		return ModelTypeDiffusion
	case ".pt", ".pth":
		return ModelTypeDiffusion
	case ".bin":
		return ModelTypeLLM
	default:
		return ModelTypeUnknown
	}
}

// detectFormat determines the model format from the extension
func detectFormat(ext string) ModelFormat {
	switch ext {
	case ".gguf":
		return FormatGGUF
	case ".safetensors":
		return FormatSafetensors
	case ".ckpt", ".pt", ".pth":
		return FormatCheckpoint
	case ".onnx":
		return FormatONNX
	case ".bin":
		return FormatBin
	default:
		return FormatUnknown
	}
}

// isModelExtension checks if an extension is a known model format
func isModelExtension(ext string) bool {
	knownExts := []string{
		".gguf", ".safetensors", ".ckpt", ".pt", ".pth",
		".bin", ".onnx",
	}

	for _, known := range knownExts {
		if ext == known {
			return true
		}
	}

	return false
}

// extractTags extracts potential tags from a filename
func extractTags(filename string) []string {
	tags := make([]string, 0)
	filename = strings.ToLower(filename)

	// Common tag patterns
	tagPatterns := map[string]string{
		"turbo":     "turbo",
		"xl":        "xl",
		"sdxl":      "xl",
		"1.5":       "sd1.5",
		"2.1":       "sd2.1",
		"3":         "sd3",
		"flux":      "flux",
		"realistic": "realistic",
		"anime":     "anime",
		"cartoon":   "cartoon",
		"q4":        "quantized",
		"q5":        "quantized",
		"q8":        "quantized",
		"fp16":      "fp16",
		"fp32":      "fp32",
	}

	for pattern, tag := range tagPatterns {
		if strings.Contains(filename, pattern) {
			tags = append(tags, tag)
		}
	}

	return tags
}

// generateModelIDFromPath generates a model ID from a file path
func generateModelIDFromPath(path string) string {
	filename := filepath.Base(path)
	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	// Clean up the name
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, " ", "_")
	nameWithoutExt = strings.ReplaceAll(nameWithoutExt, "-", "_")
	return nameWithoutExt
}

// AutoRegister scans a directory and registers all detected models
func AutoRegister(registry *Registry, dir string, recursive bool) error {
	detector := NewDetector()
	models, err := detector.DetectDirectory(dir, recursive)
	if err != nil {
		return err
	}

	for _, model := range models {
		// Check if already registered
		if _, err := registry.GetByPath(model.Path); err == nil {
			continue // Already registered
		}

		if err := registry.Register(model); err != nil {
			// Log error but continue
			continue
		}
	}

	return nil
}