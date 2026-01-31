// Package models provides model management utilities for stable-diffusion-go
package models

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// ModelType represents the type of model
type ModelType int

const (
	ModelTypeUnknown ModelType = iota
	ModelTypeDiffusion
	ModelTypeVAE
	ModelTypeLLM
	ModelTypeT5XXL
	ModelTypeCLIP
	ModelTypeEmbedding
	ModelTypeLoRA
	ModelTypeControlNet
	ModelTypeUpscaler
)

// String returns the string representation of the model type
func (t ModelType) String() string {
	switch t {
	case ModelTypeDiffusion:
		return "diffusion"
	case ModelTypeVAE:
		return "vae"
	case ModelTypeLLM:
		return "llm"
	case ModelTypeT5XXL:
		return "t5xxl"
	case ModelTypeCLIP:
		return "clip"
	case ModelTypeEmbedding:
		return "embedding"
	case ModelTypeLoRA:
		return "lora"
	case ModelTypeControlNet:
		return "controlnet"
	case ModelTypeUpscaler:
		return "upscaler"
	default:
		return "unknown"
	}
}

// ModelFormat represents the format of the model file
type ModelFormat int

const (
	FormatUnknown ModelFormat = iota
	FormatGGUF
	FormatSafetensors
	FormatCheckpoint
	FormatONNX
	FormatPT
	FormatBin
)

// String returns the string representation of the model format
func (f ModelFormat) String() string {
	switch f {
	case FormatGGUF:
		return "gguf"
	case FormatSafetensors:
		return "safetensors"
	case FormatCheckpoint:
		return "checkpoint"
	case FormatONNX:
		return "onnx"
	case FormatPT:
		return "pt"
	case FormatBin:
		return "bin"
	default:
		return "unknown"
	}
}

// ModelInfo contains metadata about a model
type ModelInfo struct {
	ID           string            `json:"id" yaml:"id"`
	Name         string            `json:"name" yaml:"name"`
	Type         ModelType         `json:"type" yaml:"type"`
	Path         string            `json:"path" yaml:"path"`
	Format       ModelFormat       `json:"format" yaml:"format"`
	Size         int64             `json:"size" yaml:"size"`
	Tags         []string          `json:"tags" yaml:"tags"`
	Metadata     map[string]string `json:"metadata" yaml:"metadata"`
	LastUsed     time.Time         `json:"last_used" yaml:"last_used"`
	UseCount     int               `json:"use_count" yaml:"use_count"`
	AddedAt      time.Time         `json:"added_at" yaml:"added_at"`
	Description  string            `json:"description" yaml:"description"`
	Source       string            `json:"source" yaml:"source"` // URL or source of the model
}

// HumanSize returns the size in human-readable format
func (m *ModelInfo) HumanSize() string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	size := m.Size
	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/GB)
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/MB)
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/KB)
	default:
		return fmt.Sprintf("%d B", size)
	}
}

// Exists checks if the model file exists
func (m *ModelInfo) Exists() bool {
	_, err := os.Stat(m.Path)
	return !os.IsNotExist(err)
}

// UpdateLastUsed updates the last used timestamp and increments use count
func (m *ModelInfo) UpdateLastUsed() {
	m.LastUsed = time.Now()
	m.UseCount++
}

// Registry manages a collection of models
type Registry struct {
	models map[string]*ModelInfo
	path   string
}

// NewRegistry creates a new model registry
func NewRegistry(registryPath string) *Registry {
	return &Registry{
		models: make(map[string]*ModelInfo),
		path:   registryPath,
	}
}

// Load loads the registry from disk
func (r *Registry) Load() error {
	data, err := os.ReadFile(r.path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // New registry
		}
		return fmt.Errorf("failed to read registry: %w", err)
	}

	var models []*ModelInfo
	if err := json.Unmarshal(data, &models); err != nil {
		return fmt.Errorf("failed to parse registry: %w", err)
	}

	r.models = make(map[string]*ModelInfo)
	for _, m := range models {
		r.models[m.ID] = m
	}

	return nil
}

// Save persists the registry to disk
func (r *Registry) Save() error {
	// Ensure directory exists
	dir := filepath.Dir(r.path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create registry directory: %w", err)
	}

	// Convert map to slice
	models := make([]*ModelInfo, 0, len(r.models))
	for _, m := range r.models {
		models = append(models, m)
	}

	data, err := json.MarshalIndent(models, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal registry: %w", err)
	}

	if err := os.WriteFile(r.path, data, 0644); err != nil {
		return fmt.Errorf("failed to write registry: %w", err)
	}

	return nil
}

// Register adds a model to the registry
func (r *Registry) Register(info *ModelInfo) error {
	if info.ID == "" {
		info.ID = generateModelID(info)
	}
	if info.AddedAt.IsZero() {
		info.AddedAt = time.Now()
	}

	r.models[info.ID] = info
	return r.Save()
}

// Unregister removes a model from the registry
func (r *Registry) Unregister(id string) error {
	delete(r.models, id)
	return r.Save()
}

// Get retrieves a model by ID
func (r *Registry) Get(id string) (*ModelInfo, error) {
	model, exists := r.models[id]
	if !exists {
		return nil, fmt.Errorf("model not found: %s", id)
	}
	return model, nil
}

// GetByPath retrieves a model by its file path
func (r *Registry) GetByPath(path string) (*ModelInfo, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	for _, model := range r.models {
		modelAbsPath, _ := filepath.Abs(model.Path)
		if modelAbsPath == absPath {
			return model, nil
		}
	}

	return nil, fmt.Errorf("model not found at path: %s", path)
}

// List returns all models, optionally filtered by type
func (r *Registry) List(filter ModelType) []*ModelInfo {
	result := make([]*ModelInfo, 0)
	for _, m := range r.models {
		if filter == ModelTypeUnknown || m.Type == filter {
			result = append(result, m)
		}
	}
	return result
}

// Search searches models by name, tags, or description
func (r *Registry) Search(query string) []*ModelInfo {
	query = strings.ToLower(query)
	result := make([]*ModelInfo, 0)

	for _, m := range r.models {
		if strings.Contains(strings.ToLower(m.Name), query) ||
			strings.Contains(strings.ToLower(m.Description), query) {
			result = append(result, m)
			continue
		}

		// Search in tags
		for _, tag := range m.Tags {
			if strings.Contains(strings.ToLower(tag), query) {
				result = append(result, m)
				break
			}
		}
	}

	return result
}

// FilterByTag returns models that have a specific tag
func (r *Registry) FilterByTag(tag string) []*ModelInfo {
	tag = strings.ToLower(tag)
	result := make([]*ModelInfo, 0)

	for _, m := range r.models {
		for _, t := range m.Tags {
			if strings.ToLower(t) == tag {
				result = append(result, m)
				break
			}
		}
	}

	return result
}

// GetStats returns statistics about the registry
func (r *Registry) GetStats() RegistryStats {
	stats := RegistryStats{
		TotalModels: len(r.models),
		ByType:      make(map[ModelType]int),
	}

	for _, m := range r.models {
		stats.ByType[m.Type]++
		stats.TotalSize += m.Size

		if stats.TotalSize > 0 {
			stats.AverageSize = stats.TotalSize / int64(stats.TotalModels)
		}
	}

	return stats
}

// Validate checks if all registered models exist on disk
func (r *Registry) Validate() []error {
	errors := make([]error, 0)

	for _, m := range r.models {
		if !m.Exists() {
			errors = append(errors, fmt.Errorf("model not found: %s (%s)", m.Name, m.Path))
		}
	}

	return errors
}

// Cleanup removes entries for models that no longer exist
func (r *Registry) Cleanup() (int, error) {
	removed := 0
	for id, m := range r.models {
		if !m.Exists() {
			delete(r.models, id)
			removed++
		}
	}

	if removed > 0 {
		return removed, r.Save()
	}

	return removed, nil
}

// RegistryStats contains statistics about the registry
type RegistryStats struct {
	TotalModels int
	TotalSize   int64
	AverageSize int64
	ByType      map[ModelType]int
}

// generateModelID generates a unique ID for a model
func generateModelID(info *ModelInfo) string {
	// Use filename + timestamp for uniqueness
	filename := filepath.Base(info.Path)
	nameWithoutExt := strings.TrimSuffix(filename, filepath.Ext(filename))
	timestamp := time.Now().UnixNano()
	return fmt.Sprintf("%s_%d", nameWithoutExt, timestamp)
}