package pipeline

import (
	"fmt"
	"os"
	"path/filepath"

	stablediffusion "github.com/kawai-network/veridium/pkg/stablediffusion"
	"github.com/kawai-network/veridium/pkg/stablediffusion/sd"
)

// TextToImageStep generates an image from text
type TextToImageStep struct {
	Params     stablediffusion.ImgGenParams
	OutputName string
}

// Name returns the step name
func (s *TextToImageStep) Name() string {
	return StepTextToImage
}

// Validate validates the step
func (s *TextToImageStep) Validate() error {
	if s.Params.Prompt == "" {
		return fmt.Errorf("prompt is required")
	}
	if s.OutputName == "" {
		s.OutputName = "txt2img_output.png"
	}
	return nil
}

// Execute executes the step
func (s *TextToImageStep) Execute(ctx *Context) error {
	if ctx.SD == nil {
		return fmt.Errorf("stable diffusion context not set")
	}

	outputPath := filepath.Join(ctx.WorkingDir, s.OutputName)

	if err := ctx.SD.GenerateImage(&s.Params, outputPath); err != nil {
		return err
	}

	ctx.InputPath = outputPath
	if ctx.OutputPath == "" {
		ctx.OutputPath = outputPath
	}

	return nil
}

// ImageToImageStep transforms an image using img2img
type ImageToImageStep struct {
	Params     stablediffusion.ImgGenParams
	OutputName string
}

// Name returns the step name
func (s *ImageToImageStep) Name() string {
	return StepImageToImage
}

// Validate validates the step
func (s *ImageToImageStep) Validate() error {
	if s.OutputName == "" {
		s.OutputName = "img2img_output.png"
	}
	return nil
}

// Execute executes the step
func (s *ImageToImageStep) Execute(ctx *Context) error {
	if ctx.SD == nil {
		return fmt.Errorf("stable diffusion context not set")
	}

	// Set input image from previous step if not specified
	if s.Params.InitImagePath == "" && ctx.InputPath != "" {
		s.Params.InitImagePath = ctx.InputPath
	}

	if s.Params.InitImagePath == "" {
		return fmt.Errorf("no input image specified")
	}

	outputPath := filepath.Join(ctx.WorkingDir, s.OutputName)

	if err := ctx.SD.GenerateImage(&s.Params, outputPath); err != nil {
		return err
	}

	ctx.InputPath = outputPath
	ctx.OutputPath = outputPath

	return nil
}

// UpscaleStep upscales an image
type UpscaleStep struct {
	Params     stablediffusion.UpscalerParams
	Factor     uint32
	OutputName string
}

// Name returns the step name
func (s *UpscaleStep) Name() string {
	return StepUpscale
}

// Validate validates the step
func (s *UpscaleStep) Validate() error {
	if s.Params.EsrganPath == "" {
		return fmt.Errorf("ESRGAN model path is required")
	}
	if s.Factor == 0 {
		s.Factor = 4
	}
	if s.OutputName == "" {
		s.OutputName = "upscaled.png"
	}
	return nil
}

// Execute executes the step
func (s *UpscaleStep) Execute(ctx *Context) error {
	if ctx.InputPath == "" {
		return fmt.Errorf("no input image to upscale")
	}

	// Create upscaler
	upscaler, err := stablediffusion.NewUpscaler(&s.Params)
	if err != nil {
		return fmt.Errorf("failed to create upscaler: %w", err)
	}

	outputPath := filepath.Join(ctx.WorkingDir, s.OutputName)

	if err := upscaler.Upscale(ctx.InputPath, s.Factor, outputPath); err != nil {
		return err
	}

	ctx.InputPath = outputPath
	ctx.OutputPath = outputPath

	return nil
}

// PreprocessStep preprocesses an image (e.g., Canny edge detection)
type PreprocessStep struct {
	Type       PreprocessType
	Params     map[string]float32
	OutputName string
}

// PreprocessType represents the type of preprocessing
type PreprocessType int

const (
	PreprocessCanny PreprocessType = iota
	PreprocessDepth
	PreprocessOpenPose
	PreprocessSegmentation
)

// Name returns the step name
func (s *PreprocessStep) Name() string {
	return StepPreprocess
}

// Validate validates the step
func (s *PreprocessStep) Validate() error {
	if s.OutputName == "" {
		s.OutputName = "preprocessed.png"
	}
	return nil
}

// Execute executes the step
func (s *PreprocessStep) Execute(ctx *Context) error {
	if ctx.InputPath == "" {
		return fmt.Errorf("no input image to preprocess")
	}

	// Load input image
	inputImage, err := sd.LoadImage(ctx.InputPath)
	if err != nil {
		return fmt.Errorf("failed to load image: %w", err)
	}

	var result bool

	switch s.Type {
	case PreprocessCanny:
		highThreshold := s.Params["high_threshold"]
		if highThreshold == 0 {
			highThreshold = 100
		}
		lowThreshold := s.Params["low_threshold"]
		if lowThreshold == 0 {
			lowThreshold = 50
		}
		weak := s.Params["weak"]
		if weak == 0 {
			weak = 0.5
		}
		strong := s.Params["strong"]
		if strong == 0 {
			strong = 1.0
		}
		inverse := s.Params["inverse"] != 0

		result = sd.PreprocessCanny(inputImage, highThreshold, lowThreshold, weak, strong, inverse)
	default:
		return fmt.Errorf("unsupported preprocess type: %v", s.Type)
	}

	if !result {
		return fmt.Errorf("preprocessing failed")
	}

	// Save preprocessed image
	outputPath := filepath.Join(ctx.WorkingDir, s.OutputName)
	if err := sd.SaveImage(&inputImage, outputPath); err != nil {
		return fmt.Errorf("failed to save preprocessed image: %w", err)
	}

	ctx.InputPath = outputPath
	ctx.OutputPath = outputPath

	return nil
}

// ConvertStep converts a model to GGUF format
type ConvertStep struct {
	InputPath       string
	VAEPath         string
	OutputPath      string
	OutputType      string
	TensorTypeRules string
	ConvertName     bool
}

// Name returns the step name
func (s *ConvertStep) Name() string {
	return StepConvert
}

// Validate validates the step
func (s *ConvertStep) Validate() error {
	if s.InputPath == "" {
		return fmt.Errorf("input path is required")
	}
	if s.OutputPath == "" {
		return fmt.Errorf("output path is required")
	}
	return nil
}

// Execute executes the step
func (s *ConvertStep) Execute(ctx *Context) error {
	if err := stablediffusion.Convert(
		s.InputPath,
		s.VAEPath,
		s.OutputPath,
		s.OutputType,
		s.TensorTypeRules,
		s.ConvertName,
	); err != nil {
		return err
	}

	ctx.InputPath = s.OutputPath
	ctx.OutputPath = s.OutputPath

	return nil
}

// SaveStep saves the current image to a specific location
type SaveStep struct {
	Destination string
}

// Name returns the step name
func (s *SaveStep) Name() string {
	return StepSave
}

// Validate validates the step
func (s *SaveStep) Validate() error {
	if s.Destination == "" {
		return fmt.Errorf("destination is required")
	}
	return nil
}

// Execute executes the step
func (s *SaveStep) Execute(ctx *Context) error {
	if ctx.InputPath == "" {
		return fmt.Errorf("no image to save")
	}

	// Ensure destination directory exists
	dir := filepath.Dir(s.Destination)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	// Copy file
	input, err := os.ReadFile(ctx.InputPath)
	if err != nil {
		return fmt.Errorf("failed to read input file: %w", err)
	}

	if err := os.WriteFile(s.Destination, input, 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	ctx.OutputPath = s.Destination

	return nil
}

// LoadStep loads an image from a specific location
type LoadStep struct {
	Source string
}

// Name returns the step name
func (s *LoadStep) Name() string {
	return StepLoad
}

// Validate validates the step
func (s *LoadStep) Validate() error {
	if s.Source == "" {
		return fmt.Errorf("source is required")
	}
	return nil
}

// Execute executes the step
func (s *LoadStep) Execute(ctx *Context) error {
	if _, err := os.Stat(s.Source); os.IsNotExist(err) {
		return fmt.Errorf("source file does not exist: %s", s.Source)
	}

	ctx.InputPath = s.Source

	return nil
}

// CustomStep allows defining custom step logic
type CustomStep struct {
	StepName   string
	ValidateFn func() error
	ExecuteFn  func(*Context) error
}

// Name returns the step name
func (s *CustomStep) Name() string {
	return s.StepName
}

// Validate validates the step
func (s *CustomStep) Validate() error {
	if s.ValidateFn != nil {
		return s.ValidateFn()
	}
	return nil
}

// Execute executes the step
func (s *CustomStep) Execute(ctx *Context) error {
	if s.ExecuteFn != nil {
		return s.ExecuteFn(ctx)
	}
	return nil
}

// PresetPipelines provides common pipeline presets

// Txt2ImgPipeline creates a simple text-to-image pipeline
func Txt2ImgPipeline(sd *stablediffusion.StableDiffusion, params stablediffusion.ImgGenParams, outputPath string) *Pipeline {
	ctx := NewContext(filepath.Dir(outputPath))
	ctx.SD = sd

	return New().
		Add(&TextToImageStep{
			Params:     params,
			OutputName: filepath.Base(outputPath),
		}).
		Add(&SaveStep{Destination: outputPath})
}

// Txt2ImgUpscalePipeline creates a text-to-image with upscaling pipeline
func Txt2ImgUpscalePipeline(
	sd *stablediffusion.StableDiffusion,
	imgParams stablediffusion.ImgGenParams,
	upscalerParams stablediffusion.UpscalerParams,
	upscaleFactor uint32,
	outputPath string,
) *Pipeline {
	workingDir := filepath.Dir(outputPath)
	ctx := NewContext(workingDir)
	ctx.SD = sd

	return New().
		Add(&TextToImageStep{
			Params:     imgParams,
			OutputName: "generated.png",
		}).
		Add(&UpscaleStep{
			Params:     upscalerParams,
			Factor:     upscaleFactor,
			OutputName: "upscaled.png",
		}).
		Add(&SaveStep{Destination: outputPath})
}

// Img2ImgPipeline creates an image-to-image pipeline
func Img2ImgPipeline(
	sd *stablediffusion.StableDiffusion,
	inputPath string,
	params stablediffusion.ImgGenParams,
	outputPath string,
) *Pipeline {
	workingDir := filepath.Dir(outputPath)
	ctx := NewContext(workingDir)
	ctx.SD = sd

	params.InitImagePath = inputPath

	return New().
		Add(&ImageToImageStep{
			Params:     params,
			OutputName: filepath.Base(outputPath),
		}).
		Add(&SaveStep{Destination: outputPath})
}

// ControlNetPipeline creates a ControlNet pipeline
func ControlNetPipeline(
	sd *stablediffusion.StableDiffusion,
	controlImagePath string,
	params stablediffusion.ImgGenParams,
	outputPath string,
) *Pipeline {
	workingDir := filepath.Dir(outputPath)
	ctx := NewContext(workingDir)
	ctx.SD = sd

	params.ControlImagePath = controlImagePath

	return New().
		Add(&TextToImageStep{
			Params:     params,
			OutputName: filepath.Base(outputPath),
		}).
		Add(&SaveStep{Destination: outputPath})
}

// PreprocessAndGeneratePipeline creates a pipeline that preprocesses and generates
func PreprocessAndGeneratePipeline(
	sd *stablediffusion.StableDiffusion,
	inputPath string,
	preprocessType PreprocessType,
	params stablediffusion.ImgGenParams,
	outputPath string,
) *Pipeline {
	workingDir := filepath.Dir(outputPath)
	ctx := NewContext(workingDir)
	ctx.SD = sd

	return New().
		Add(&LoadStep{Source: inputPath}).
		Add(&PreprocessStep{
			Type:       preprocessType,
			OutputName: "preprocessed.png",
		}).
		Add(&TextToImageStep{
			Params:     params,
			OutputName: filepath.Base(outputPath),
		}).
		Add(&SaveStep{Destination: outputPath})
}
