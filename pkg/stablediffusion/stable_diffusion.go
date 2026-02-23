package stablediffusion

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	gosd "github.com/getkawai/stablediffusion"
	sd "github.com/kawai-network/stablediffusion"
)

// internalSDMu protects initialization state
var internalSDMu sync.RWMutex
var libraryInitialized bool

// InitLibrary initializes the stable diffusion library.
// This function is thread-safe and should be called once at application startup.
// Subsequent calls will reinitialize the library (not recommended).
func InitLibrary(libPath string) error {
	internalSDMu.Lock()
	defer internalSDMu.Unlock()

	if libPath == "" {
		return errors.New("library path is empty")
	}
	if _, err := os.Stat(libPath); err != nil {
		return fmt.Errorf("invalid library path: %w", err)
	}
	if err := gosd.Init(libPath); err != nil {
		return fmt.Errorf("failed to initialize gosd library: %w", err)
	}
	libraryInitialized = true
	return nil
}

// ============================================================================
// Default Generation Parameters
// ============================================================================
// These constants define sensible defaults for image and video generation.
// Values are chosen based on community best practices and model recommendations.
// ============================================================================

// Image Generation Defaults
const (
	// DefaultCfgScale is the default Classifier-Free Guidance scale.
	// Higher values (7-10) follow prompts more strictly but may reduce quality.
	// Lower values (3-5) allow more creative freedom. 5.0 is a balanced default.
	DefaultCfgScale = 5.0

	// DefaultImageCfgScale is the default image guidance scale for inpaint/instruct-pix2pix.
	// 1.0 means no additional image guidance beyond the base CfgScale.
	DefaultImageCfgScale = 1.0

	// DefaultDistilledGuidance is the default distilled guidance scale for models with guidance input.
	// 3.5 is recommended for distilled models like SDXL Turbo.
	DefaultDistilledGuidance = 3.5

	// DefaultSampleSteps is the default number of denoising steps.
	// 20 steps provides good quality for most samplers. Fewer steps (10-15) for faster generation,
	// more steps (30-50) for higher quality.
	DefaultSampleSteps = 20

	// DefaultStrength is the default noise strength for img2img.
	// 0.75 provides a good balance between following the prompt and preserving the original image.
	// Range: 0.0 (preserve original) to 1.0 (complete noise).
	DefaultStrength = 0.75

	// DefaultSeed is the default random seed.
	// 42 is a conventional choice (reference to "Hitchhiker's Guide to the Galaxy").
	// Use negative values or 0 with random seed generation for variation.
	DefaultSeed = 42

	// DefaultBatchCount is the default number of images to generate in one batch.
	DefaultBatchCount = 1

	// DefaultClipSkip specifies how many CLIP layers to skip.
	// -1 means no skip (use all layers). Some models benefit from skipping 1-2 layers.
	DefaultClipSkip = -1

	// DefaultEta is the default eta parameter for DDIM/TCD samplers.
	// 1.0 is the standard value for ancestral samplers.
	DefaultEta = 1.0

	// DefaultSkipLayerStart is when Skip Layer Guidance (SLG) starts (as fraction of total steps).
	DefaultSkipLayerStart = 0.01

	// DefaultSkipLayerEnd is when Skip Layer Guidance (SLG) ends (as fraction of total steps).
	DefaultSkipLayerEnd = 0.2

	// DefaultControlStrength is the default strength for ControlNet guidance.
	// 0.9 provides strong control while allowing some creative freedom.
	DefaultControlStrength = 0.9
)

// Video Generation Defaults
const (
	// DefaultVideoFrames is the default number of frames to generate.
	// 33 frames at 30fps = ~1.1 seconds of video.
	DefaultVideoFrames = 33

	// DefaultMOEBoundary is the timestep boundary for Wan2.2 MoE models.
	// 0.875 is the recommended value for MoE models.
	DefaultMOEBoundary = 0.875

	// DefaultVaceStrength is the default strength for Wan VACE (Video Attention Control Enhancement).
	DefaultVaceStrength = 1.0

	// DefaultHighNoiseCfgScale is the cfg scale for high noise diffusion models.
	DefaultHighNoiseCfgScale = 6.0

	// DefaultHighNoiseSampleSteps is the sample steps for high noise diffusion models.
	// -1 indicates auto-calculation based on main sample steps.
	DefaultHighNoiseSampleSteps = -1
)

// Image Dimension Defaults
const (
	// DefaultWidth is the default image width in pixels.
	DefaultWidth = 512

	// DefaultHeight is the default image height in pixels.
	DefaultHeight = 512
)

// Embedding embedding structure for defining model embeddings
type Embedding struct {
	Name string // Embedding name
	Path string // Embedding file path
}

// ContextParams context parameters structure for initializing Stable Diffusion context
type ContextParams struct {
	ModelPath                   string     // Full model path
	ClipLPath                   string     // CLIP-L text encoder path
	ClipGPath                   string     // CLIP-G text encoder path
	ClipVisionPath              string     // CLIP Vision encoder path
	T5XXLPath                   string     // T5-XXL text encoder path
	LLMPath                     string     // LLM text encoder path (e.g., qwenvl2.5 for qwen-image, mistral-small3.2 for flux2)
	LLMVisionPath               string     // LLM Vision encoder path
	DiffusionModelPath          string     // Standalone diffusion model path
	HighNoiseDiffusionModelPath string     // Standalone high noise diffusion model path
	VAEPath                     string     // VAE model path
	TAESDPath                   string     // TAE-SD model path, uses Tiny AutoEncoder for fast decoding (low quality)
	ControlNetPath              string     // ControlNet model path
	Embeddings                  *Embedding // Embedding information
	EmbeddingCount              uint32     // Number of embeddings
	PhotoMakerPath              string     // PhotoMaker model path
	TensorTypeRules             string     // Weight type rules per tensor pattern (e.g., "^vae\.=f16,model\.=q8_0")
	VAEDecodeOnly               bool       // Process VAE using only decode mode
	FreeParamsImmediately       bool       // Whether to free parameters immediately
	NThreads                    int32      // Number of threads to use for generation
	WType                       string     // Weight type (default: auto-detect from model file)
	RNGType                     string     // Random number generator type (default: "cuda")
	SamplerRNGType              string     // Sampler random number generator type (default: "cuda")
	Prediction                  string     // Prediction type override
	LoraApplyMode               string     // LoRA application mode (default: "auto")
	OffloadParamsToCPU          bool       // Keep weights in RAM to save VRAM, auto-load to VRAM when needed
	EnableMmap                  bool       // Whether to enable memory mapping
	KeepClipOnCPU               bool       // Keep CLIP on CPU (for low VRAM)
	KeepControlNetOnCPU         bool       // Keep ControlNet on CPU (for low VRAM)
	KeepVAEOnCPU                bool       // Keep VAE on CPU (for low VRAM)
	DiffusionFlashAttn          bool       // Use Flash attention in diffusion model (significantly reduces memory usage)
	TAEPreviewOnly              bool       // Prevent decoding final image with taesd (for preview="tae")
	DiffusionConvDirect         bool       // Use Conv2d direct in diffusion model
	VAEConvDirect               bool       // Use Conv2d direct in VAE model (should improve performance)
	CircularX                   bool       // Enable circular padding on X axis
	CircularY                   bool       // Enable circular padding on Y axis
	ForceSDXLVAConvScale        bool       // Force conv scale on SDXL VAE
	ChromaUseDitMask            bool       // Whether Chroma uses DiT mask
	ChromaUseT5Mask             bool       // Whether Chroma uses T5 mask
	ChromaT5MaskPad             int32      // Chroma T5 mask padding size
	QwenImageZeroCondT          bool       // Qwen-image zero condition T parameter
	FlowShift                   float32    // Shift value for Flow models (e.g., SD3.x or WAN)
}

// Lora LoRA structure for defining LoRA model parameters
type Lora struct {
	IsHighNoise bool    // Whether it's a high noise LoRA
	Multiplier  float32 // LoRA multiplier
	Path        string  // LoRA file path
}

// PMParams PhotoMaker parameters structure for defining PhotoMaker related parameters
type PMParams struct {
	IDImages      *sd.SDImage // ID images pointer
	IDImagesCount int32       // Number of ID images
	IDEmbedPath   string      // PhotoMaker v2 ID embedding path
	StyleStrength float32     // Strength to keep PhotoMaker input identity
}

// ImgGenParams image generation parameters structure for defining image generation related parameters
type ImgGenParams struct {
	Loras              *Lora             // LoRA parameters
	LoraCount          uint32            // Number of LoRAs
	Prompt             string            // Prompt to render
	NegativePrompt     string            // Negative prompt
	ClipSkip           int32             // Skip last layers of CLIP network (1 = no skip, 2 = skip one layer, <=0 = not specified)
	InitImagePath      string            // Initial image path for guidance
	RefImagesPath      []string          // Array of reference image paths for Flux Kontext models
	RefImagesCount     int32             // Number of reference images
	AutoResizeRefImage bool              // Whether to auto-resize reference images
	IncreaseRefIndex   bool              // Whether to auto-increase index based on reference image list order (starting from 1)
	MaskImagePath      string            // Inpainting mask image path
	Width              int32             // Image width (pixels)
	Height             int32             // Image height (pixels)
	CfgScale           float32           // Unconditional guidance scale.
	ImageCfgScale      float32           // Image guidance scale for inpaint or instruct-pix2pix models (default: same as `CfgScale`).
	DistilledGuidance  float32           // Distilled guidance scale for models with guidance input.
	SkipLayers         []int32           // Layers to skip for SLG steps (SLG will be enabled at step int([STEPS]x[START]) and disabled at int([STEPS]x[END])).
	SkipLayerStart     float32           // SLG enabling point.
	SkipLayerEnd       float32           // SLG disabling point.
	SlgScale           float32           // Skip layer guidance (SLG) scale, only for DiT models.
	Scheduler          string            // Denoiser sigma scheduler (default: discrete).
	SampleMethod       string            // Sampling method (default: euler for Flux/SD3/Wan, euler_a otherwise).
	SampleSteps        int32             // Number of sample steps.
	Eta                float32           // Eta in DDIM, only for DDIM and TCD.
	ShiftedTimestep    int32             // Shift timestep for NitroFusion models, default: 0, recommended N for NitroSD-Realism around 250 and 500 for NitroSD-Vibrant.
	CustomSigmas       []float32         // Custom sigma values for the sampler, comma-separated (e.g. "14.61,7.8,3.5,0.0").
	Strength           float32           // Noise/denoise strength (range [0.0, 1.0])
	Seed               int64             // RNG seed (< 0 for random seed)
	BatchCount         int32             // Number of images to generate
	ControlImagePath   string            // Control condition image path for ControlNet
	ControlStrength    float32           // Strength to apply ControlNet
	PMParams           *PMParams         // PhotoMaker parameters
	VAETilingParams    sd.SDTilingParams // VAE tiling parameters for reducing memory usage
	CacheParams        sd.SDCacheParams  // Cache parameters for DiT models
}

// VidGenParams video generation parameters structure for defining video generation related parameters
type VidGenParams struct {
	Loras             *Lora    // LoRA parameters
	LoraCount         uint32   // Number of LoRAs
	Prompt            string   // Prompt to render
	NegativePrompt    string   // Negative prompt
	ClipSkip          int32    // Skip last layers of CLIP network (1 = no skip, 2 = skip one layer, <=0 = not specified)
	InitImagePath     string   // Initial image path for starting generation
	EndImagePath      string   // End image path for ending generation (required for flf2v)
	ControlFramesPath []string // Array of control frame image paths for video
	ControlFramesSize int32    // Control frame size
	Width             int32    // Video width (pixels)
	Height            int32    // Video height (pixels)

	CfgScale          float32   // Unconditional guidance scale.
	ImageCfgScale     float32   // Image guidance scale for inpaint or instruct-pix2pix models (default: same as `CfgScale`).
	DistilledGuidance float32   // Distilled guidance scale for models with guidance input.
	SkipLayers        []int32   // Layers to skip for SLG steps (SLG will be enabled at step int([STEPS]x[START]) and disabled at int([STEPS]x[END])).
	SkipLayerStart    float32   // SLG enabling point.
	SkipLayerEnd      float32   // SLG disabling point.
	SlgScale          float32   // Skip layer guidance (SLG) scale, only for DiT models.
	Scheduler         string    // Denoiser sigma scheduler (default: discrete).
	SampleMethod      string    // Sampling method (default: euler for Flux/SD3/Wan, euler_a otherwise).
	SampleSteps       int32     // Number of sample steps.
	Eta               float32   // Eta in DDIM, only for DDIM and TCD.
	ShiftedTimestep   int32     // Shift timestep for NitroFusion models, default: 0, recommended N for NitroSD-Realism around 250 and 500 for NitroSD-Vibrant.
	CustomSigmas      []float32 // Custom sigma values for the sampler, comma-separated (e.g. "14.61,7.8,3.5,0.0").

	HighNoiseCfgScale          float32   // High noise diffusion model equivalent of `cfg_scale`.
	HighNoiseImageCfgScale     float32   // High noise diffusion model equivalent of `image_cfg_scale`.
	HighNoiseDistilledGuidance float32   // High noise diffusion model equivalent of `guidance`.
	HighNoiseSkipLayers        []int32   // High noise diffusion model equivalent of `skip_layers`.
	HighNoiseSkipLayerStart    float32   // High noise diffusion model equivalent of `skip_layer_start`.
	HighNoiseSkipLayerEnd      float32   // High noise diffusion model equivalent of `skip_layer_end`.
	HighNoiseSlgScale          float32   // High noise diffusion model equivalent of `slg_scale`.
	HighNoiseScheduler         string    // High noise diffusion model equivalent of `scheduler`.
	HighNoiseSampleMethod      string    // High noise diffusion model equivalent of `sample_method`.
	HighNoiseSampleSteps       int32     // High noise diffusion model equivalent of `sample_steps` (default: -1 = auto).
	HighNoiseEta               float32   // High noise diffusion model equivalent of `eta`.
	HighNoiseShiftedTimestep   int32     // Shift timestep for NitroFusion models, default: 0, recommended N for NitroSD-Realism around 250 and 500 for NitroSD-Vibrant.
	HighNoiseCustomSigmas      []float32 // Custom sigma values for the sampler, comma-separated (e.g. "14.61,7.8,3.5,0.0").

	MOEBoundary  float32          // Timestep boundary for Wan2.2 MoE models
	Strength     float32          // Noise/denoise strength (range [0.0, 1.0])
	Seed         int64            // RNG seed (< 0 for random seed)
	VideoFrames  int32            // Number of video frames to generate
	VaceStrength float32          // Wan VACE strength
	CacheParams  sd.SDCacheParams // Cache parameters for DiT models
}

// StableDiffusion is the main structure for interacting with the Stable Diffusion library.
// It wraps a SDContext pointer and provides methods for image and video generation.
//
// # Lifecycle Management
//
// The StableDiffusion instance must be explicitly freed when no longer needed by calling Free().
// Generation methods (GenerateImage, GenerateVideo) do NOT automatically free the context,
// allowing you to reuse the same instance for multiple generations.
//
// # Usage Pattern
//
//	sd, err := NewStableDiffusion(params)
//	if err != nil {
//	    // handle error
//	}
//	defer sd.Free() // Ensure cleanup when done
//
//	// Generate multiple images with the same instance
//	err = sd.GenerateImage(params1, "output1.png")
//	err = sd.GenerateImage(params2, "output2.png")
//
// # Thread Safety
//
// StableDiffusion instances are NOT thread-safe. Do not call generation methods
// from multiple goroutines simultaneously on the same instance.
type StableDiffusion struct {
	engine *gosd.Engine
}

// IsReady reports whether the Stable Diffusion backend engine is loaded.
func (sDiffusion *StableDiffusion) IsReady() bool {
	return sDiffusion != nil && sDiffusion.engine != nil
}

// Free releases all resources held by the StableDiffusion context.
// After calling Free(), the instance cannot be used for further generations.
// It is safe to call Free() multiple times - subsequent calls are no-ops.
//
// IMPORTANT: Always call Free() when done with the instance to prevent memory leaks.
// The recommended pattern is to use defer:
//
//	sd, err := NewStableDiffusion(params)
//	if err != nil {
//	    return err
//	}
//	defer sd.Free()
func (sDiffusion *StableDiffusion) Free() {
	if sDiffusion.engine != nil {
		sDiffusion.engine.Close()
		sDiffusion.engine = nil
	}
}

// NewStableDiffusion creates a stable diffusion instance
func NewStableDiffusion(ctxParams *ContextParams) (*StableDiffusion, error) {
	internalSDMu.RLock()
	initialized := libraryInitialized
	internalSDMu.RUnlock()
	if !initialized {
		return nil, errors.New("library not initialized, call InitLibrary first")
	}
	if ctxParams == nil {
		return nil, errors.New("context params cannot be nil")
	}

	modelFile := ctxParams.DiffusionModelPath
	if modelFile == "" {
		modelFile = ctxParams.ModelPath
	}
	if modelFile == "" {
		return nil, errors.New("missing model path: set DiffusionModelPath or ModelPath")
	}

	pathsForBase := []string{modelFile}
	if ctxParams.LLMPath != "" {
		pathsForBase = append(pathsForBase, ctxParams.LLMPath)
	}
	if ctxParams.VAEPath != "" {
		pathsForBase = append(pathsForBase, ctxParams.VAEPath)
	}
	baseDir := commonBaseDir(pathsForBase)

	relModel, err := filepath.Rel(baseDir, modelFile)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve model path: %w", err)
	}
	if strings.HasPrefix(relModel, "..") {
		return nil, fmt.Errorf("model path %q is outside base directory %q", modelFile, baseDir)
	}

	opts := []string{"diffusion_model"}
	appendPathOpt := func(name, value string) error {
		if value == "" {
			return nil
		}
		rel, err := filepath.Rel(baseDir, value)
		if err != nil {
			return fmt.Errorf("failed to resolve %s: %w", name, err)
		}
		if strings.HasPrefix(rel, "..") {
			return fmt.Errorf("%s path %q is outside base directory %q", name, value, baseDir)
		}
		opts = append(opts, fmt.Sprintf("%s:%s", name, rel))
		return nil
	}

	if err := appendPathOpt("llm_path", ctxParams.LLMPath); err != nil {
		return nil, err
	}
	if err := appendPathOpt("vae_path", ctxParams.VAEPath); err != nil {
		return nil, err
	}

	if ctxParams.OffloadParamsToCPU {
		opts = append(opts, "offload_params_to_cpu:true")
	}
	if ctxParams.DiffusionFlashAttn {
		opts = append(opts, "diffusion_flash_attn:true")
	}
	if ctxParams.FlowShift != 0 {
		opts = append(opts, fmt.Sprintf("flow_shift:%g", ctxParams.FlowShift))
	}

	engine := &gosd.Engine{}
	if err := engine.Load(gosd.ModelOptions{
		Threads:   ctxParams.NThreads,
		ModelPath: baseDir,
		ModelFile: relModel,
		Options:   opts,
		CFGScale:  float32(DefaultCfgScale),
	}); err != nil {
		return nil, fmt.Errorf("failed to load stable diffusion model: %w", err)
	}

	return &StableDiffusion{engine: engine}, nil
}

// GenerateImage generates an image from a text prompt or existing image (img2img).
//
// Parameters:
//   - imgGenParams: Image generation parameters including prompt, dimensions, steps, etc.
//   - newImagePath: Output path for the generated image (PNG format). Directories are created automatically.
//
// Returns:
//   - error: nil on success, or an error describing what went wrong.
//
// # Resource Management
//
// This method does NOT free the StableDiffusion context. You can call GenerateImage
// multiple times on the same instance to generate multiple images efficiently.
// Always call Free() when you're completely done with the instance.
//
// # Examples
//
//	// Text-to-image
//	err := sd.GenerateImage(&ImgGenParams{
//	    Prompt: "a beautiful landscape",
//	    Width: 512,
//	    Height: 512,
//	}, "output.png")
//
//	// Image-to-image
//	err := sd.GenerateImage(&ImgGenParams{
//	    Prompt:        "make it look like a painting",
//	    InitImagePath: "input.png",
//	    Strength:      0.75,
//	}, "output.png")
//
//	// Multiple generations with same instance
//	for i, prompt := range prompts {
//	    err := sd.GenerateImage(&ImgGenParams{Prompt: prompt}, fmt.Sprintf("img_%d.png", i))
//	    if err != nil {
//	        log.Printf("Failed to generate image %d: %v", i, err)
//	    }
//	}
func (sDiffusion *StableDiffusion) GenerateImage(imgGenParams *ImgGenParams, newImagePath string) error {
	if sDiffusion == nil || sDiffusion.engine == nil {
		return errors.New("image generation engine is not available")
	}
	if imgGenParams == nil {
		return errors.New("image generation params cannot be nil")
	}

	dir := filepath.Dir(newImagePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	width := imgGenParams.Width
	height := imgGenParams.Height
	if width == 0 {
		width = DefaultWidth
	}
	if height == 0 {
		height = DefaultHeight
	}

	steps := imgGenParams.SampleSteps
	if steps == 0 {
		steps = DefaultSampleSteps
	}

	cfg := float32(imgGenParams.CfgScale)
	if cfg == 0 {
		cfg = float32(DefaultCfgScale)
	}

	seed := imgGenParams.Seed
	if seed == 0 {
		seed = DefaultSeed
	}

	strength := imgGenParams.Strength
	if imgGenParams.InitImagePath != "" && strength == 0 {
		strength = DefaultStrength
	}

	negative := imgGenParams.NegativePrompt
	if negative == "" {
		negative = "blurry, low quality, distorted"
	}

	enableParams := ""
	if imgGenParams.MaskImagePath != "" {
		enableParams = "mask:" + imgGenParams.MaskImagePath
	}

	return sDiffusion.engine.GenerateImage(gosd.GenerateImageOptions{
		PositivePrompt:   imgGenParams.Prompt,
		NegativePrompt:   negative,
		Dst:              newImagePath,
		Src:              imgGenParams.InitImagePath,
		EnableParameters: enableParams,
		RefImages:        imgGenParams.RefImagesPath,
		Width:            width,
		Height:           height,
		Seed:             seed,
		Step:             steps,
		Strength:         strength,
		CFGScale:         cfg,
	})
}

// GenerateVideo generates a video from a text prompt or existing images.
//
// Parameters:
//   - vidGenParams: Video generation parameters including prompt, dimensions, frames, etc.
//   - newVideoPath: Output path for the generated video (MP4 format). Requires FFmpeg.
//
// Returns:
//   - error: nil on success, or an error describing what went wrong.
//
// # Resource Management
//
// This method does NOT free the StableDiffusion context. You can call GenerateVideo
// multiple times on the same instance to generate multiple videos efficiently.
// Always call Free() when you're completely done with the instance.
//
// # Requirements
//
// FFmpeg must be installed and available in PATH for video encoding.
//
// # Examples
//
//	// Text-to-video
//	err := sd.GenerateVideo(&VidGenParams{
//	    Prompt: "a beautiful sunset over mountains",
//	    Width: 512,
//	    Height: 512,
//	    VideoFrames: 33,
//	}, "output.mp4")
//
//	// Image-to-video
//	err := sd.GenerateVideo(&VidGenParams{
//	    Prompt: "animate this scene",
//	    InitImagePath: "start.png",
//	    VideoFrames: 33,
//	}, "output.mp4")
func (sDiffusion *StableDiffusion) GenerateVideo(vidGenParams *VidGenParams, newVideoPath string) error {
	return errors.New("GenerateVideo is not supported by the current backend")
}

type UpscalerParams struct {
	EsrganPath         string // ESRGAN model path
	OffloadParamsToCPU bool   // Whether to save parameters to CPU
	Direct             bool   // Whether to use direct mode
	NThreads           int    // Number of threads to use
	TileSize           int    // Tile size
}

type Upscaler struct {
}

var errUnsupportedFeature = errors.New("feature not supported by the current backend")

func commonBaseDir(paths []string) string {
	if len(paths) == 0 {
		return "."
	}
	base := filepath.Clean(filepath.Dir(paths[0]))
	for i := 1; i < len(paths); i++ {
		dir := filepath.Clean(filepath.Dir(paths[i]))
		for {
			rel, err := filepath.Rel(base, dir)
			if err == nil && rel != ".." && !strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
				break
			}
			parent := filepath.Dir(base)
			if parent == base {
				return filepath.Clean(filepath.Dir(paths[0]))
			}
			base = parent
		}
	}
	return base
}

// NewUpscaler creates a new upscaler context
func NewUpscaler(params *UpscalerParams) (*Upscaler, error) {
	return nil, errUnsupportedFeature
}

// Upscale upscaling function
func (us *Upscaler) Upscale(inputImagePath string, upscaleFactor uint32, outputImagePath string) error {
	return errUnsupportedFeature
}

// Convert model conversion function, convert a model to gguf format.
// inputPath: Path to the input model.
// vaePath: Path to the vae.
// outputPath: Path to save the converted model.
// outputType: The weight type (default: auto).
// tensorTypeRules: Weight type per tensor pattern (example: "^vae\\\\.=f16,model\\\\.=q8_0")
func Convert(inputPath, vaePath, outputPath, outputType, tensorTypeRules string, convertName bool) error {
	return errUnsupportedFeature
}

// Re-export utility functions from external package for convenience

// LoadImage loads an image from file and converts to SDImage format
var LoadImage = sd.LoadImage

// SaveImage saves SDImage as PNG file
var SaveImage = sd.SaveImage

// SaveFrames saves all video frames as PNG files
var SaveFrames = sd.SaveFrames

// EncodeVideo encodes PNG frame sequence to video using FFmpeg
var EncodeVideo = sd.EncodeVideo

// PreprocessCanny preprocesses image with Canny edge detection
func PreprocessCanny(image sd.SDImage, highThreshold, lowThreshold, weak, strong float32, inverse bool) bool {
	return false
}

// generateImageFromPath generates SDImage from path (internal helper)
func generateImageFromPath(imagePath string) sd.SDImage {
	if imagePath == "" {
		return sd.SDImage{}
	}

	img, err := sd.LoadImage(imagePath)
	if err != nil {
		fmt.Println("Error loading image:", err)
		return sd.SDImage{}
	}
	return img
}

// generateImagesFromPaths generates multiple SDImages from paths (internal helper)
// Returns the slice of images to ensure the underlying array stays alive during C library calls
func generateImagesFromPaths(paths []string) []sd.SDImage {
	if len(paths) == 0 {
		return nil
	}

	// Create SDImage slice
	images := make([]sd.SDImage, 0, len(paths))

	// Iterate through all paths, generate SDImage
	for _, p := range paths {
		if p == "" {
			continue
		}

		img := generateImageFromPath(p)
		// Only add valid images
		if img.Data != nil {
			images = append(images, img)
		}
	}

	if len(images) == 0 {
		return nil
	}

	// Return the slice - caller must keep it alive during C library calls
	return images
}

// SetProgressCallback sets the progress callback for the internal SD instance
func SetProgressCallback(cb func(step int, steps int, time float32, data interface{}), data interface{}) {
	_ = cb
	_ = data
}

// CleanupTempDir cleans up temporary directory
func CleanupTempDir(tempDir string) error {
	return os.RemoveAll(tempDir)
}
