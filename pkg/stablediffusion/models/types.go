package models

import "time"

// ModelSpec represents a Stable Diffusion model specification for catalog
type ModelSpec struct {
	Name             string // Model name/identifier
	URL              string // Download URL for the primary diffusion model
	Filename         string // Local filename for the primary diffusion model
	LLMURL           string // Optional download URL for LLM/text encoder model
	LLMFilename      string // Optional filename for LLM/text encoder model
	VAEURL           string // Optional download URL for VAE model
	VAEFilename      string // Optional filename for VAE model
	EditModelURL     string // Optional download URL for image edit model
	EditModelFile    string // Optional filename for image edit model
	EditFallbackURL  string // Optional fallback URL for image edit model
	EditFallbackFile string // Optional fallback filename for image edit model
	Size             int64  // Model file size in MB
	MinRAM           int64  // Minimum RAM required in GB
	RecommendedRAM   int64  // Recommended RAM in GB
	MinVRAM          int64  // Minimum VRAM required in GB (0 if CPU-only)
	RecommendedVRAM  int64  // Recommended VRAM in GB
	ModelType        string // Type of model (SD1.5, SDXL, etc.)
	Description      string // Model description
	Quantization     string // Quantization level (f16, q4_0, q8_0, etc.)
}

// HardwareSpecs represents system hardware specifications
type HardwareSpecs struct {
	TotalRAM     int64  // Total RAM in GB
	AvailableRAM int64  // Available RAM in GB
	CPU          string // CPU model
	CPUCores     int    // Number of CPU cores
	GPUMemory    int64  // GPU VRAM in GB (if available)
	GPUModel     string // GPU model
}

// ModelType represents the type of model (for detector/downloader)
type ModelType string

const (
	ModelTypeUnknown    ModelType = "unknown"
	ModelTypeDiffusion  ModelType = "diffusion"
	ModelTypeVAE        ModelType = "vae"
	ModelTypeLLM        ModelType = "llm"
	ModelTypeT5XXL      ModelType = "t5xxl"
	ModelTypeCLIP       ModelType = "clip"
	ModelTypeEmbedding  ModelType = "embedding"
	ModelTypeLoRA       ModelType = "lora"
	ModelTypeControlNet ModelType = "controlnet"
	ModelTypeUpscaler   ModelType = "upscaler"
)

// ModelFormat represents the file format of a model
type ModelFormat string

const (
	FormatUnknown     ModelFormat = "unknown"
	FormatGGUF        ModelFormat = "gguf"
	FormatSafetensors ModelFormat = "safetensors"
	FormatCheckpoint  ModelFormat = "checkpoint"
	FormatONNX        ModelFormat = "onnx"
	FormatBin         ModelFormat = "bin"
)

// ModelInfo represents detailed information about a detected model
type ModelInfo struct {
	ID          string            // Unique identifier
	Name        string            // Model name
	Type        ModelType         // Model type
	Path        string            // File path
	Format      ModelFormat       // File format
	Size        int64             // File size in bytes
	Tags        []string          // Tags extracted from filename
	Metadata    map[string]string // Additional metadata
	AddedAt     time.Time         // When model was added
	Description string            // Model description
	Source      string            // Download source URL
}

// Registry represents a model registry (forward declaration for detector)
type Registry interface {
	GetByPath(path string) (*ModelInfo, error)
	Register(model *ModelInfo) error
}
