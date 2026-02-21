package stablediffusion

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	sd "github.com/kawai-network/stablediffusion"
)

// internalSD is the internal StableDiffusion instance
// Protected by internalSDMu mutex for thread-safe access
var internalSD *sd.StableDiffusion

// internalSDMu protects access to internalSD
var internalSDMu sync.RWMutex

// InitLibrary initializes the stable diffusion library.
// This function is thread-safe and should be called once at application startup.
// Subsequent calls will reinitialize the library (not recommended).
func InitLibrary(libPath string) error {
	var err error

	internalSDMu.Lock()
	defer internalSDMu.Unlock()

	internalSD, err = sd.New(sd.LibraryConfig{LibPath: libPath})
	return err
}

// getInternalSD returns the internal SD instance with proper synchronization
// Returns nil if the library has not been initialized
func getInternalSD() *sd.StableDiffusion {
	internalSDMu.RLock()
	defer internalSDMu.RUnlock()
	return internalSD
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

// RNGTypeMap RNG type mapping
var RNGTypeMap = map[string]sd.RngType{
	"default":    sd.DefaultRNG,
	"cuda":       sd.CUDARNG, // Default
	"cpu":        sd.CPURNG,
	"type_count": sd.RNGTypeCount,
}

// SampleMethodMap sampling method mapping
var SampleMethodMap = map[string]sd.SampleMethod{
	"default":             -1, // Default
	"euler":               sd.EulerSampleMethod,
	"euler_a":             sd.EulerASampleMethod,
	"heun":                sd.HeunSampleMethod,
	"dpm2":                sd.DPM2SampleMethod,
	"dpm++2s_a":           sd.DPMPP2SASampleMethod,
	"dpm++2m":             sd.DPMPP2MSampleMethod,
	"dpm++2mv2":           sd.DPMPP2Mv2SampleMethod,
	"ipndm":               sd.IPNDMSampleMethod,
	"ipndm_v":             sd.IPNDMSampleMethodV,
	"lcm":                 sd.LCMSampleMethod,
	"ddim_trailing":       sd.DDIMTrailingSampleMethod,
	"tcd":                 sd.TCDSampleMethod,
	"sample_method_count": sd.SampleMethodCount,
}

// SchedulerMap scheduler mapping
var SchedulerMap = map[string]sd.Scheduler{
	"default":         -1, // Default
	"discrete":        sd.DiscreteScheduler,
	"karras":          sd.KarrasScheduler,
	"exponential":     sd.ExponentialScheduler,
	"ays":             sd.AYSScheduler,
	"gits":            sd.GITScheduler,
	"sgm_uniform":     sd.SGMUniformScheduler,
	"simple":          sd.SimpleScheduler,
	"smoothstep":      sd.SmoothstepScheduler,
	"kl_optimal":      sd.KLOptimalScheduler,
	"lcm":             sd.LCMScheduler,
	"scheduler_count": sd.SchedulerCount,
}

// PredictionMap prediction type mapping
var PredictionMap = map[string]sd.Prediction{
	"eps":        sd.EPSPred,
	"v":          sd.VPred,
	"edm_v":      sd.EDMVPred,
	"flow":       sd.FlowPred,
	"flux_flow":  sd.FluxFlowPred,
	"flux2_flow": sd.Flux2FlowPred,
	"default":    sd.PredictionCount, // Default
}

// SDTypeMap SDType mapping
var SDTypeMap = map[string]sd.SDType{
	"f32":  sd.SDTypeF32,
	"f16":  sd.SDTypeF16,
	"q4_0": sd.SDTypeQ4_0,
	"q4_1": sd.SDTypeQ4_1,
	"q5_0": sd.SDTypeQ5_0,
	"q5_1": sd.SDTypeQ5_1,
	"q8_0": sd.SDTypeQ8_0,
	"q8_1": sd.SDTypeQ8_1,
	// k-quantizations
	"q2_k":    sd.SDTypeQ2_K,
	"q3_k":    sd.SDTypeQ3_K,
	"q4_k":    sd.SDTypeQ4_K,
	"q5_k":    sd.SDTypeQ5_K,
	"q6_k":    sd.SDTypeQ6_K,
	"q8_k":    sd.SDTypeQ8_K,
	"iq2_xxs": sd.SDTypeIQ2_XXS,
	"iq2_xs":  sd.SDTypeIQ2_XS,
	"iq3_xxs": sd.SDTypeIQ3_XXS,
	"iq1_s":   sd.SDTypeIQ1_S,
	"iq4_nl":  sd.SDTypeIQ4_NL,
	"iq3_s":   sd.SDTypeIQ3_S,
	"iq2_s":   sd.SDTypeIQ2_S,
	"iq4_xs":  sd.SDTypeIQ4_XS,
	"i8":      sd.SDTypeI8,
	"i16":     sd.SDTypeI16,
	"i32":     sd.SDTypeI32,
	"i64":     sd.SDTypeI64,
	"f64":     sd.SDTypeF64,
	"iq1_m":   sd.SDTypeIQ1_M,
	"bf16":    sd.SDTypeBF16,
	// "q4_0_4_4": SD_TYPE_Q4_0_4_4,
	// "q4_0_4_8": SD_TYPE_Q4_0_4_8,
	// "q4_0_8_8": SD_TYPE_Q4_0_8_8,
	"tq1_0": sd.SDTypeTQ1_0,
	"tq2_0": sd.SDTypeTQ2_0,
	// "iq4_nl_4_4": SD_TYPE_IQ4_NL_4_4,
	// "iq4_nl_4_8": SD_TYPE_IQ4_NL_4_8,
	// "iq4_nl_8_8": SD_TYPE_IQ4_NL_8_8,
	"mxfp4":   sd.SDTypeMXFP4,
	"default": sd.SDTypeCount, // Default
}

// PreviewMap preview type mapping
var PreviewMap = map[string]sd.Preview{
	"none":          sd.PreviewNone, // Default
	"proj":          sd.PreviewProj,
	"tae":           sd.PreviewTAE,
	"vae":           sd.PreviewVAE,
	"preview_count": sd.PreviewCount,
}

// LoraApplyModeMap LoRA apply mode mapping
var LoraApplyModeMap = map[string]sd.LoraApplyMode{
	"auto":                  sd.LoraApplyAuto, // Default
	"immediately":           sd.LoraApplyImmediately,
	"at_runtime":            sd.LoraApplyAtRuntime,
	"lora_apply_mode_count": sd.LoraApplyModeCount,
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
	ctx *sd.SDContext // SD context pointer
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
	if sDiffusion.ctx != nil {
		sDiffusion.ctx.Free()
		sDiffusion.ctx = nil
	}
}

// NewStableDiffusion creates a stable diffusion instance
func NewStableDiffusion(ctxParams *ContextParams) (*StableDiffusion, error) {
	// Check if library is initialized
	lib := getInternalSD()
	if lib == nil {
		return nil, errors.New("library not initialized, call InitLibrary first")
	}

	// 1. Initialize context parameters
	var sdCtxParams sd.SDContextParams
	lib.ContextParamsInit(&sdCtxParams)

	if ctxParams.ModelPath != "" {
		sdCtxParams.ModelPath = sd.CString(ctxParams.ModelPath)
	}

	if ctxParams.ClipLPath != "" {
		sdCtxParams.ClipLPath = sd.CString(ctxParams.ClipLPath)
	}

	if ctxParams.ClipGPath != "" {
		sdCtxParams.ClipGPath = sd.CString(ctxParams.ClipGPath)
	}

	if ctxParams.ClipVisionPath != "" {
		sdCtxParams.ClipVisionPath = sd.CString(ctxParams.ClipVisionPath)
	}

	if ctxParams.T5XXLPath != "" {
		sdCtxParams.T5XXLPath = sd.CString(ctxParams.T5XXLPath)
	}

	if ctxParams.LLMPath != "" {
		sdCtxParams.LLMPath = sd.CString(ctxParams.LLMPath)
	}

	if ctxParams.LLMVisionPath != "" {
		sdCtxParams.LLMVisionPath = sd.CString(ctxParams.LLMVisionPath)
	}

	if ctxParams.DiffusionModelPath != "" {
		sdCtxParams.DiffusionModelPath = sd.CString(ctxParams.DiffusionModelPath)
	}

	if ctxParams.HighNoiseDiffusionModelPath != "" {
		sdCtxParams.HighNoiseDiffusionModelPath = sd.CString(ctxParams.HighNoiseDiffusionModelPath)
	}

	if ctxParams.VAEPath != "" {
		sdCtxParams.VAEPath = sd.CString(ctxParams.VAEPath)
	}

	if ctxParams.TAESDPath != "" {
		sdCtxParams.TAESDPath = sd.CString(ctxParams.TAESDPath)
	}

	if ctxParams.ControlNetPath != "" {
		sdCtxParams.ControlNetPath = sd.CString(ctxParams.ControlNetPath)
	}

	if ctxParams.Embeddings != nil {
		sdCtxParams.Embeddings = &sd.SDEmbedding{
			Name: sd.CString(ctxParams.Embeddings.Name),
			Path: sd.CString(ctxParams.Embeddings.Path),
		}
	}

	if ctxParams.EmbeddingCount > 0 {
		sdCtxParams.EmbeddingCount = ctxParams.EmbeddingCount
	}

	if ctxParams.PhotoMakerPath != "" {
		sdCtxParams.PhotoMakerPath = sd.CString(ctxParams.PhotoMakerPath)
	}

	if ctxParams.TensorTypeRules != "" {
		sdCtxParams.TensorTypeRules = sd.CString(ctxParams.TensorTypeRules)
	}

	sdCtxParams.VAEDecodeOnly = ctxParams.VAEDecodeOnly
	sdCtxParams.FreeParamsImmediately = ctxParams.FreeParamsImmediately

	if ctxParams.NThreads > 0 {
		sdCtxParams.NThreads = ctxParams.NThreads
	}

	if ctxParams.WType != "" {
		if WType, ok := SDTypeMap[ctxParams.WType]; ok {
			sdCtxParams.WType = WType
		} else {
			return nil, fmt.Errorf("invalid WType: %s", ctxParams.WType)
		}
	}

	if ctxParams.RNGType != "" {
		if RNGType, ok := RNGTypeMap[ctxParams.RNGType]; ok {
			sdCtxParams.RNGType = RNGType
		} else {
			return nil, fmt.Errorf("invalid RNG type: %s", ctxParams.RNGType)
		}
	}

	if ctxParams.SamplerRNGType != "" {
		if RNGType, ok := RNGTypeMap[ctxParams.SamplerRNGType]; ok {
			sdCtxParams.SamplerRNGType = RNGType
		} else {
			return nil, fmt.Errorf("invalid Sampler RNG type: %s", ctxParams.SamplerRNGType)
		}
	}

	if ctxParams.Prediction != "" {
		if Prediction, ok := PredictionMap[ctxParams.Prediction]; ok {
			sdCtxParams.Prediction = Prediction
		} else {
			return nil, fmt.Errorf("invalid Prediction: %s", ctxParams.Prediction)
		}
	}

	if ctxParams.LoraApplyMode != "" {
		if LoraApplyMode, ok := LoraApplyModeMap[ctxParams.LoraApplyMode]; ok {
			sdCtxParams.LoraApplyMode = LoraApplyMode
		} else {
			return nil, fmt.Errorf("invalid LoraApplyMode: %s", ctxParams.LoraApplyMode)
		}
	}

	sdCtxParams.OffloadParamsToCPU = ctxParams.OffloadParamsToCPU
	sdCtxParams.EnableMmap = ctxParams.EnableMmap
	sdCtxParams.KeepClipOnCPU = ctxParams.KeepClipOnCPU
	sdCtxParams.KeepControlNetOnCPU = ctxParams.KeepControlNetOnCPU
	sdCtxParams.KeepVAEOnCPU = ctxParams.KeepVAEOnCPU
	sdCtxParams.DiffusionFlashAttn = ctxParams.DiffusionFlashAttn
	sdCtxParams.TAEPreviewOnly = ctxParams.TAEPreviewOnly
	sdCtxParams.DiffusionConvDirect = ctxParams.DiffusionConvDirect
	sdCtxParams.VAEConvDirect = ctxParams.VAEConvDirect
	sdCtxParams.CircularX = ctxParams.CircularX
	sdCtxParams.CircularY = ctxParams.CircularY
	sdCtxParams.ForceSDXLVAConvScale = ctxParams.ForceSDXLVAConvScale
	sdCtxParams.ChromaUseDitMask = ctxParams.ChromaUseDitMask
	sdCtxParams.ChromaUseT5Mask = ctxParams.ChromaUseT5Mask

	if ctxParams.ChromaT5MaskPad != 0 {
		sdCtxParams.ChromaT5MaskPad = ctxParams.ChromaT5MaskPad
	}

	sdCtxParams.QwenImageZeroCondT = ctxParams.QwenImageZeroCondT

	if ctxParams.FlowShift != 0 {
		sdCtxParams.FlowShift = ctxParams.FlowShift
	}

	// 2. Create new context
	ctx, err := lib.NewContext(&sdCtxParams)
	if err != nil {
		return nil, fmt.Errorf("failed to create context: %w", err)
	}
	if ctx == nil {
		return nil, errors.New("failed to create context: returned nil")
	}

	return &StableDiffusion{
		ctx: ctx,
	}, nil
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
	// Extract the directory part of newImagePath and create it if it doesn't exist
	dir := filepath.Dir(newImagePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err = os.MkdirAll(dir, os.ModePerm)
		if err != nil {
			return errors.New("failed to create directory")
		}
	}

	// Initialize image generation parameters
	var sdImgGenParams sd.SDImgGenParams
	lib := getInternalSD()
	if lib == nil {
		return errors.New("library not initialized")
	}
	lib.ImgGenParamsInit(&sdImgGenParams)

	// Set generation parameters
	sdImgGenParams.Prompt = sd.CString(imgGenParams.Prompt)
	if imgGenParams.NegativePrompt == "" {
		imgGenParams.NegativePrompt = "blurry, low quality, distorted"
	}

	sdImgGenParams.NegativePrompt = sd.CString(imgGenParams.NegativePrompt)

	if imgGenParams.Loras == nil {
		sdImgGenParams.Loras = &sd.SDLora{
			IsHighNoise: false,
			Multiplier:  0,
			Path:        sd.CString(""),
		}
	} else {
		sdImgGenParams.Loras = &sd.SDLora{
			IsHighNoise: imgGenParams.Loras.IsHighNoise,
			Multiplier:  imgGenParams.Loras.Multiplier,
			Path:        sd.CString(imgGenParams.Loras.Path),
		}
	}

	sdImgGenParams.LoraCount = imgGenParams.LoraCount

	if imgGenParams.ClipSkip == 0 {
		imgGenParams.ClipSkip = DefaultClipSkip
	}
	sdImgGenParams.ClipSkip = imgGenParams.ClipSkip

	sdImgGenParams.InitImage = generateImageFromPath(imgGenParams.InitImagePath)

	// Process reference images - keep slice alive during C library call
	refImagesSlice := generateImagesFromPaths(imgGenParams.RefImagesPath)
	if len(refImagesSlice) > 0 {
		sdImgGenParams.RefImages = &refImagesSlice[0]
	} else {
		sdImgGenParams.RefImages = nil
	}

	// Set reference images count, use actual loaded count if user didn't provide
	sdImgGenParams.RefImagesCount = int32(len(imgGenParams.RefImagesPath))
	if imgGenParams.RefImagesCount > 0 {
		sdImgGenParams.RefImagesCount = imgGenParams.RefImagesCount
	}
	sdImgGenParams.AutoResizeRefImage = imgGenParams.AutoResizeRefImage
	sdImgGenParams.IncreaseRefIndex = imgGenParams.IncreaseRefIndex
	sdImgGenParams.MaskImage = generateImageFromPath(imgGenParams.MaskImagePath)

	// For img2img, use requested dimensions if explicitly provided.
	// Otherwise fallback to initial image dimensions.
	if imgGenParams.InitImagePath != "" {
		if imgGenParams.Width > 0 && imgGenParams.Height > 0 {
			sdImgGenParams.Width = imgGenParams.Width
			sdImgGenParams.Height = imgGenParams.Height
		} else {
			sdImgGenParams.Width = int32(sdImgGenParams.InitImage.Width)
			sdImgGenParams.Height = int32(sdImgGenParams.InitImage.Height)
		}
	} else {
		// Otherwise use default dimensions
		if imgGenParams.Width == 0 {
			imgGenParams.Width = DefaultWidth
		}
		if imgGenParams.Height == 0 {
			imgGenParams.Height = DefaultHeight
		}
		sdImgGenParams.Width = imgGenParams.Width
		sdImgGenParams.Height = imgGenParams.Height
	}

	if imgGenParams.CfgScale == 0 {
		imgGenParams.CfgScale = DefaultCfgScale
	}
	if imgGenParams.ImageCfgScale == 0 {
		imgGenParams.ImageCfgScale = DefaultImageCfgScale
	}

	if imgGenParams.DistilledGuidance == 0 {
		imgGenParams.DistilledGuidance = DefaultDistilledGuidance
	}

	var skipLayersPtr *int32
	if len(imgGenParams.SkipLayers) > 0 {
		skipLayersPtr = &imgGenParams.SkipLayers[0]
	} else {
		skipLayersPtr = nil
	}

	if imgGenParams.SkipLayerStart == 0 {
		imgGenParams.SkipLayerStart = DefaultSkipLayerStart
	}
	if imgGenParams.SkipLayerEnd == 0 {
		imgGenParams.SkipLayerEnd = DefaultSkipLayerEnd
	}

	var defaultSampleMethod sd.SampleMethod
	if imgGenParams.SampleMethod == "" || imgGenParams.SampleMethod == "default" {
		defaultSampleMethod = lib.GetDefaultSampleMethod(sDiffusion.ctx)
	} else {
		if sampleMethod, ok := SampleMethodMap[imgGenParams.SampleMethod]; ok {
			defaultSampleMethod = sampleMethod
		} else {
			return fmt.Errorf("invalid SampleMethod: %s", imgGenParams.SampleMethod)
		}
	}

	var defaultScheduler sd.Scheduler
	if imgGenParams.Scheduler == "" || imgGenParams.Scheduler == "default" {
		defaultScheduler = lib.GetDefaultScheduler(sDiffusion.ctx, defaultSampleMethod)
	} else {
		if scheduler, ok := SchedulerMap[imgGenParams.Scheduler]; ok {
			defaultScheduler = scheduler
		} else {
			return fmt.Errorf("invalid Scheduler: %s", imgGenParams.Scheduler)
		}
	}

	if imgGenParams.SampleSteps == 0 {
		imgGenParams.SampleSteps = DefaultSampleSteps
	}

	if imgGenParams.Eta == 0 {
		imgGenParams.Eta = DefaultEta
	}

	var customSigmasPtr *float32
	if len(imgGenParams.CustomSigmas) > 0 {
		customSigmasPtr = &imgGenParams.CustomSigmas[0]
	} else {
		customSigmasPtr = nil
	}

	sdImgGenParams.SampleParams = sd.SDSampleParams{
		Guidance: sd.SDGuidanceParams{
			TxtCfg:            imgGenParams.CfgScale,
			ImgCfg:            imgGenParams.ImageCfgScale,
			DistilledGuidance: imgGenParams.DistilledGuidance,
			SLG: sd.SDSLGParams{
				Layers:     skipLayersPtr,
				LayerCount: uintptr(len(imgGenParams.SkipLayers)),
				LayerStart: imgGenParams.SkipLayerStart,
				LayerEnd:   imgGenParams.SkipLayerEnd,
				Scale:      imgGenParams.SlgScale,
			},
		},
		SampleMethod:      defaultSampleMethod,
		Scheduler:         defaultScheduler,
		SampleSteps:       imgGenParams.SampleSteps,
		Eta:               imgGenParams.Eta,
		ShiftedTimestep:   imgGenParams.ShiftedTimestep,
		CustomSigmas:      customSigmasPtr,
		CustomSigmasCount: int32(len(imgGenParams.CustomSigmas)),
	}

	if imgGenParams.Strength == 0 {
		imgGenParams.Strength = DefaultStrength
	}
	sdImgGenParams.Strength = imgGenParams.Strength

	if imgGenParams.Seed == 0 {
		imgGenParams.Seed = DefaultSeed
	}
	sdImgGenParams.Seed = imgGenParams.Seed

	if imgGenParams.BatchCount == 0 {
		imgGenParams.BatchCount = DefaultBatchCount
	}
	sdImgGenParams.BatchCount = imgGenParams.BatchCount

	sdImgGenParams.ControlImage = generateImageFromPath(imgGenParams.ControlImagePath)

	if imgGenParams.ControlStrength == 0 {
		imgGenParams.ControlStrength = DefaultControlStrength
	}
	sdImgGenParams.ControlStrength = imgGenParams.ControlStrength

	if imgGenParams.PMParams != nil {
		sdImgGenParams.PMParams = sd.SDPMParams{
			IDImages:      imgGenParams.PMParams.IDImages,
			IDImagesCount: imgGenParams.PMParams.IDImagesCount,
			IDEmbedPath:   sd.CString(imgGenParams.PMParams.IDEmbedPath),
			StyleStrength: imgGenParams.PMParams.StyleStrength,
		}
	}

	if imgGenParams.VAETilingParams != (sd.SDTilingParams{}) {
		sdImgGenParams.VAETilingParams = imgGenParams.VAETilingParams
	} else {
		sdImgGenParams.VAETilingParams = sd.SDTilingParams{
			Enabled:       false,
			TileSizeX:     0,
			TileSizeY:     0,
			TargetOverlap: 0.5,
			RelSizeX:      0,
			RelSizeY:      0,
		}
	}

	// Initialize cache parameters
	var cacheParams sd.SDCacheParams
	lib.CacheParamsInit(&cacheParams)

	// If user provided cache parameters, use them
	if imgGenParams.CacheParams != (sd.SDCacheParams{}) {
		sdImgGenParams.Cache = imgGenParams.CacheParams
	} else {
		// Otherwise use default parameters
		sdImgGenParams.Cache = cacheParams
		// Set some reasonable default values
		sdImgGenParams.Cache.Mode = sd.SDCacheDisabled
		sdImgGenParams.Cache.ReuseThreshold = 0.2
		sdImgGenParams.Cache.StartPercent = 0.15
		sdImgGenParams.Cache.EndPercent = 0.95
	}

	// Generate image
	img := sDiffusion.ctx.GenerateImage(&sdImgGenParams)
	if img == nil {
		return errors.New("failed to generate image: context returned nil image")
	}

	// Save image
	if err := sd.SaveImage(img, newImagePath); err != nil {
		return fmt.Errorf("failed to save image: %w", err)
	}

	return nil
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
	// Extract the directory part of newVideoPath and create it if it doesn't exist
	dir := filepath.Dir(newVideoPath)
	tmpDir := filepath.Join(dir, "tmp")
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		err = os.MkdirAll(tmpDir, os.ModePerm)
		if err != nil {
			return errors.New("failed to create directory")
		}
	}

	// Initialize video generation parameters
	var sdVidGenParams sd.SDVidGenParams
	lib := getInternalSD()
	if lib == nil {
		return errors.New("library not initialized")
	}
	lib.VidGenParamsInit(&sdVidGenParams)

	// Set generation parameters
	sdVidGenParams.Prompt = sd.CString(vidGenParams.Prompt)
	sdVidGenParams.NegativePrompt = sd.CString(vidGenParams.NegativePrompt)

	if vidGenParams.Loras == nil {
		sdVidGenParams.Loras = &sd.SDLora{
			IsHighNoise: false,
			Multiplier:  0,
			Path:        sd.CString(""),
		}
	} else {
		sdVidGenParams.Loras = &sd.SDLora{
			IsHighNoise: vidGenParams.Loras.IsHighNoise,
			Multiplier:  vidGenParams.Loras.Multiplier,
			Path:        sd.CString(vidGenParams.Loras.Path),
		}
	}

	sdVidGenParams.LoraCount = vidGenParams.LoraCount

	if vidGenParams.ClipSkip == 0 {
		vidGenParams.ClipSkip = DefaultClipSkip
	}
	sdVidGenParams.ClipSkip = vidGenParams.ClipSkip

	sdVidGenParams.InitImage = generateImageFromPath(vidGenParams.InitImagePath)
	sdVidGenParams.EndImage = generateImageFromPath(vidGenParams.EndImagePath)

	// Process control frames
	var controlFrames []sd.SDImage
	if len(vidGenParams.ControlFramesPath) > 0 {
		controlFrames = make([]sd.SDImage, len(vidGenParams.ControlFramesPath))
		for i, path := range vidGenParams.ControlFramesPath {
			controlFrames[i] = generateImageFromPath(path)
		}
	}
	if len(controlFrames) > 0 {
		sdVidGenParams.ControlFrames = &controlFrames[0]
		sdVidGenParams.ControlFramesSize = int32(len(controlFrames))
	}

	if vidGenParams.Width == 0 {
		vidGenParams.Width = DefaultWidth
	}
	if vidGenParams.Height == 0 {
		vidGenParams.Height = DefaultHeight
	}
	sdVidGenParams.Width = vidGenParams.Width
	sdVidGenParams.Height = vidGenParams.Height

	if vidGenParams.CfgScale == 0 {
		vidGenParams.CfgScale = DefaultCfgScale
	}
	if vidGenParams.ImageCfgScale == 0 {
		vidGenParams.ImageCfgScale = DefaultImageCfgScale
	}

	if vidGenParams.DistilledGuidance == 0 {
		vidGenParams.DistilledGuidance = DefaultDistilledGuidance
	}

	var skipLayersPtr *int32
	if len(vidGenParams.SkipLayers) > 0 {
		skipLayersPtr = &vidGenParams.SkipLayers[0]
	} else {
		skipLayersPtr = nil
	}

	if vidGenParams.SkipLayerStart == 0 {
		vidGenParams.SkipLayerStart = DefaultSkipLayerStart
	}
	if vidGenParams.SkipLayerEnd == 0 {
		vidGenParams.SkipLayerEnd = DefaultSkipLayerEnd
	}

	var defaultSampleMethod sd.SampleMethod
	if vidGenParams.SampleMethod == "" || vidGenParams.SampleMethod == "default" {
		defaultSampleMethod = lib.GetDefaultSampleMethod(sDiffusion.ctx)
	} else {
		if sampleMethod, ok := SampleMethodMap[vidGenParams.SampleMethod]; ok {
			defaultSampleMethod = sampleMethod
		} else {
			return fmt.Errorf("invalid SampleMethod: %s", vidGenParams.SampleMethod)
		}
	}

	var defaultScheduler sd.Scheduler
	if vidGenParams.Scheduler == "" || vidGenParams.Scheduler == "default" {
		defaultScheduler = lib.GetDefaultScheduler(sDiffusion.ctx, defaultSampleMethod)
	} else {
		if scheduler, ok := SchedulerMap[vidGenParams.Scheduler]; ok {
			defaultScheduler = scheduler
		} else {
			return fmt.Errorf("invalid Scheduler: %s", vidGenParams.Scheduler)
		}
	}

	if vidGenParams.SampleSteps == 0 {
		vidGenParams.SampleSteps = DefaultSampleSteps
	}

	if vidGenParams.Eta == 0 {
		vidGenParams.Eta = DefaultEta
	}

	var customSigmasPtr *float32
	if len(vidGenParams.CustomSigmas) > 0 {
		customSigmasPtr = &vidGenParams.CustomSigmas[0]
	} else {
		customSigmasPtr = nil
	}

	sdVidGenParams.SampleParams = sd.SDSampleParams{
		Guidance: sd.SDGuidanceParams{
			TxtCfg:            vidGenParams.CfgScale,
			ImgCfg:            vidGenParams.ImageCfgScale,
			DistilledGuidance: vidGenParams.DistilledGuidance,
			SLG: sd.SDSLGParams{
				Layers:     skipLayersPtr,
				LayerCount: uintptr(len(vidGenParams.SkipLayers)),
				LayerStart: vidGenParams.SkipLayerStart,
				LayerEnd:   vidGenParams.SkipLayerEnd,
				Scale:      vidGenParams.SlgScale,
			},
		},
		SampleMethod:      defaultSampleMethod,
		Scheduler:         defaultScheduler,
		SampleSteps:       vidGenParams.SampleSteps,
		Eta:               vidGenParams.Eta,
		ShiftedTimestep:   vidGenParams.ShiftedTimestep,
		CustomSigmas:      customSigmasPtr,
		CustomSigmasCount: int32(len(vidGenParams.CustomSigmas)),
	}

	if vidGenParams.HighNoiseCfgScale == 0 {
		vidGenParams.HighNoiseCfgScale = DefaultHighNoiseCfgScale
	}
	if vidGenParams.HighNoiseImageCfgScale == 0 {
		vidGenParams.HighNoiseImageCfgScale = DefaultImageCfgScale
	}

	if vidGenParams.HighNoiseDistilledGuidance == 0 {
		vidGenParams.HighNoiseDistilledGuidance = DefaultDistilledGuidance
	}

	var highNoiseSkipLayersPtr *int32
	if len(vidGenParams.HighNoiseSkipLayers) > 0 {
		highNoiseSkipLayersPtr = &vidGenParams.HighNoiseSkipLayers[0]
	} else {
		highNoiseSkipLayersPtr = nil
	}

	if vidGenParams.HighNoiseSkipLayerStart == 0 {
		vidGenParams.HighNoiseSkipLayerStart = DefaultSkipLayerStart
	}
	if vidGenParams.HighNoiseSkipLayerEnd == 0 {
		vidGenParams.HighNoiseSkipLayerEnd = DefaultSkipLayerEnd
	}

	var defaultHighNoiseSampleMethod sd.SampleMethod
	if vidGenParams.HighNoiseSampleMethod == "" || vidGenParams.HighNoiseSampleMethod == "default" {
		defaultHighNoiseSampleMethod = lib.GetDefaultSampleMethod(sDiffusion.ctx)
	} else {
		if sampleMethod, ok := SampleMethodMap[vidGenParams.HighNoiseSampleMethod]; ok {
			defaultHighNoiseSampleMethod = sampleMethod
		} else {
			return fmt.Errorf("invalid SampleMethod: %s", vidGenParams.HighNoiseSampleMethod)
		}
	}

	var defaultHighNoiseScheduler sd.Scheduler
	if vidGenParams.HighNoiseScheduler == "" || vidGenParams.HighNoiseScheduler == "default" {
		defaultHighNoiseScheduler = lib.GetDefaultScheduler(sDiffusion.ctx, defaultHighNoiseSampleMethod)
	} else {
		if scheduler, ok := SchedulerMap[vidGenParams.Scheduler]; ok {
			defaultHighNoiseScheduler = scheduler
		} else {
			return fmt.Errorf("invalid Scheduler: %s", vidGenParams.HighNoiseScheduler)
		}
	}

	if vidGenParams.HighNoiseSampleSteps == 0 {
		vidGenParams.HighNoiseSampleSteps = DefaultHighNoiseSampleSteps
	}

	if vidGenParams.HighNoiseEta == 0 {
		vidGenParams.HighNoiseEta = DefaultEta
	}

	var highNoiseCustomSigmasPtr *float32
	if len(vidGenParams.HighNoiseCustomSigmas) > 0 {
		highNoiseCustomSigmasPtr = &vidGenParams.HighNoiseCustomSigmas[0]
	} else {
		highNoiseCustomSigmasPtr = nil
	}

	sdVidGenParams.HighNoiseSampleParams = sd.SDSampleParams{
		Guidance: sd.SDGuidanceParams{
			TxtCfg:            vidGenParams.HighNoiseCfgScale,
			ImgCfg:            vidGenParams.HighNoiseImageCfgScale,
			DistilledGuidance: vidGenParams.HighNoiseDistilledGuidance,
			SLG: sd.SDSLGParams{
				Layers:     highNoiseSkipLayersPtr,
				LayerCount: uintptr(len(vidGenParams.HighNoiseSkipLayers)),
				LayerStart: vidGenParams.HighNoiseSkipLayerStart,
				LayerEnd:   vidGenParams.HighNoiseSkipLayerEnd,
				Scale:      vidGenParams.HighNoiseSlgScale,
			},
		},
		SampleMethod:      defaultHighNoiseSampleMethod,
		Scheduler:         defaultHighNoiseScheduler,
		SampleSteps:       vidGenParams.HighNoiseSampleSteps,
		Eta:               vidGenParams.HighNoiseEta,
		ShiftedTimestep:   vidGenParams.HighNoiseShiftedTimestep,
		CustomSigmas:      highNoiseCustomSigmasPtr,
		CustomSigmasCount: int32(len(vidGenParams.HighNoiseCustomSigmas)),
	}

	if vidGenParams.MOEBoundary == 0 {
		vidGenParams.MOEBoundary = DefaultMOEBoundary
	}
	sdVidGenParams.MOEBoundary = vidGenParams.MOEBoundary

	if vidGenParams.Strength == 0 {
		vidGenParams.Strength = DefaultStrength
	}
	sdVidGenParams.Strength = vidGenParams.Strength

	if vidGenParams.Seed == 0 {
		vidGenParams.Seed = DefaultSeed
	}
	sdVidGenParams.Seed = vidGenParams.Seed

	if vidGenParams.VideoFrames == 0 {
		vidGenParams.VideoFrames = DefaultVideoFrames
	}
	sdVidGenParams.VideoFrames = vidGenParams.VideoFrames

	if vidGenParams.VaceStrength == 0 {
		vidGenParams.VaceStrength = DefaultVaceStrength
	}
	sdVidGenParams.VaceStrength = vidGenParams.VaceStrength

	// Initialize cache parameters
	var cacheParams sd.SDCacheParams
	lib.CacheParamsInit(&cacheParams)

	// If user provided cache parameters, use them
	if vidGenParams.CacheParams != (sd.SDCacheParams{}) {
		sdVidGenParams.Cache = vidGenParams.CacheParams
	} else {
		// Otherwise use default parameters
		sdVidGenParams.Cache = cacheParams
		// Set some reasonable default values
		sdVidGenParams.Cache.Mode = sd.SDCacheDisabled
		sdVidGenParams.Cache.ReuseThreshold = 0.2
		sdVidGenParams.Cache.StartPercent = 0.15
		sdVidGenParams.Cache.EndPercent = 0.95
	}

	// Generate video
	frames, numFrames := sDiffusion.ctx.GenerateVideo(&sdVidGenParams)
	if frames == nil || numFrames == 0 {
		return errors.New("failed to generate video: context returned no frames")
	}

	// Save video frames
	if err := sd.SaveFrames(frames, tmpDir); err != nil {
		return fmt.Errorf("failed to save video frames: %w", err)
	}
	if err := sd.EncodeVideo(tmpDir, newVideoPath, 30); err != nil {
		return fmt.Errorf("failed to encode video: %w", err)
	}
	if err := os.RemoveAll(tmpDir); err != nil {
		return fmt.Errorf("failed to cleanup temp directory: %w", err)
	}

	return nil
}

type UpscalerParams struct {
	EsrganPath         string // ESRGAN model path
	OffloadParamsToCPU bool   // Whether to save parameters to CPU
	Direct             bool   // Whether to use direct mode
	NThreads           int    // Number of threads to use
	TileSize           int    // Tile size
}

type Upscaler struct {
	ctx *sd.UpscalerContext
}

// NewUpscaler creates a new upscaler context
func NewUpscaler(params *UpscalerParams) (*Upscaler, error) {
	if params.NThreads == 0 {
		params.NThreads = -1
	}

	if params.TileSize == 0 {
		params.TileSize = 128
	}

	lib := getInternalSD()
	if lib == nil {
		return nil, errors.New("library not initialized")
	}

	ctx, err := lib.NewUpscalerContext(
		params.EsrganPath,
		params.OffloadParamsToCPU,
		params.Direct,
		int32(params.NThreads),
		int32(params.TileSize),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create upscaler context: %w", err)
	}
	return &Upscaler{ctx: ctx}, nil
}

// Upscale upscaling function
func (us *Upscaler) Upscale(inputImagePath string, upscaleFactor uint32, outputImagePath string) error {
	// Directly use LoadImage function to avoid dangling pointer issues
	inputSDImage := generateImageFromPath(inputImagePath)
	fmt.Printf("inputSDImage: %+v", inputSDImage)
	outputSDImage := us.ctx.Upscale(inputSDImage, upscaleFactor)
	fmt.Printf("outputSDImage: %+v", outputSDImage)

	defer us.ctx.Free()

	// Save image
	err := sd.SaveImage(&outputSDImage, outputImagePath)
	if err != nil {
		return fmt.Errorf("failed to save image: %v", err)
	}

	return nil
}

// Convert model conversion function, convert a model to gguf format.
// inputPath: Path to the input model.
// vaePath: Path to the vae.
// outputPath: Path to save the converted model.
// outputType: The weight type (default: auto).
// tensorTypeRules: Weight type per tensor pattern (example: "^vae\\\\.=f16,model\\\\.=q8_0")
func Convert(inputPath, vaePath, outputPath, outputType, tensorTypeRules string, convertName bool) error {
	if outputPath == "" {
		outputPath = "output.gguf"
	}

	if outputType == "" {
		outputType = "default"
	}

	var outputSDType sd.SDType
	if value, ok := SDTypeMap[outputType]; ok {
		outputSDType = value
	} else {
		return fmt.Errorf("invalid SDType: %s", outputType)
	}

	res, err := sd.Convert(inputPath, vaePath, outputPath, outputSDType, tensorTypeRules, convertName)
	if err != nil {
		return fmt.Errorf("failed to convert model: %w", err)
	}
	if !res {
		return errors.New("failed to convert model: conversion returned false")
	}
	return nil
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
	lib := getInternalSD()
	if lib == nil {
		return false
	}
	return lib.PreprocessCanny(image, highThreshold, lowThreshold, weak, strong, inverse)
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
	lib := getInternalSD()
	if lib != nil {
		lib.SetProgressCallback(cb, data)
	}
}

// CleanupTempDir cleans up temporary directory
func CleanupTempDir(tempDir string) error {
	return os.RemoveAll(tempDir)
}
