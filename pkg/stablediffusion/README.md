# stable-diffusion-go

A pure Golang binding library for `stable-diffusion.cpp` based on `github.com/ebitengine/purego`, **no cgo dependency required**, supporting cross-platform operation.

## 🌟 Project Features

- **Pure Go Implementation**: Based on the purego library, calls C++ dynamic libraries without cgo
- **Cross-platform Support**: Supports Windows, Linux, macOS, and other mainstream operating systems
- **Complete Functionality**: Implements the main APIs of stable-diffusion.cpp, including text-to-image, image-to-image, video generation, etc.
- **Simple and Easy to Use**: Provides a concise Go language API for easy integration into existing projects
- **High Performance**: Supports performance optimization features like FlashAttention and model quantization
- **Flexible Library Support**: Compatible with multiple backend options including CPU (AVX/AVX2/AVX512), CUDA, ROCm, and Vulkan
- **Configuration Management**: YAML/JSON config files for easy setup and sharing
- **Batch Generation**: Generate multiple images with parallel processing
- **Progress Tracking**: Real-time progress callbacks with percentage and ETA
- **Model Management**: Registry system for organizing and managing models
- **Pipeline Operations**: Chain multiple operations (txt2img → upscale → save)

## 📁 Project Structure

```
stable-diffusion-go/
├── configs/            # Configuration file examples
│   └── example.yaml    # Example configuration file
├── examples/           # Example programs directory
│   ├── txt2img.go      # Text-to-image generation example
│   └── txt2vid.go      # Text-to-video generation example
├── lib/                # Dynamic library reference directory
│   ├── ggml.txt        # GGML library version reference
│   ├── stable-diffusion.cpp.txt  # Stable Diffusion C++ library reference
│   └── version.txt     # Required library version (master-453-4ff2c8c)
├── pkg/                # Go package directory
│   ├── sd/             # Core binding library
│   │   ├── load_library_unix.go   # Unix platform dynamic library loading
│   │   ├── load_library_windows.go # Windows platform dynamic library loading
│   │   ├── stable_diffusion.go    # Core functionality implementation
│   │   └── utils.go               # Auxiliary utility functions
│   ├── config/         # Configuration management
│   │   ├── types.go    # Config structures
│   │   ├── loader.go   # Config load/save
│   │   └── converter.go # Config converters
│   ├── batch/          # Batch generation
│   │   ├── types.go    # Batch types and generator
│   │   └── builder.go  # Batch builder
│   ├── models/         # Model management
│   │   ├── types.go    # Model registry
│   │   ├── detector.go # Model auto-detection
│   │   └── downloader.go # Model downloader
│   ├── pipeline/       # Pipeline operations
│   │   ├── types.go    # Pipeline types
│   │   └── steps.go    # Pipeline steps
│   └── progress/       # Progress tracking
│       ├── types.go    # Progress types
│       └── wrapper.go  # Progress wrapper
├── .gitignore          # Git ignore file configuration
├── LICENSE             # License file
├── README.md           # Project documentation (English)
├── README-ZH.md        # Project documentation (Chinese)
├── go.mod              # Go module file
├── go.sum              # Go dependency checksum file
├── output_demo.png     # Demo output image
└── stable_diffusion.go # Root directory entry file
```

**Important**: The dynamic libraries are **NOT included** in this repository. You must download them separately based on your platform and requirements.

## 🚀 Quick Start

### 1. Install Dependencies

```bash
go get github.com/orangelang/stable-diffusion-go
```

### 2. Prepare Model Files

Model files need to be prepared before use, supporting multiple formats:
- Diffusion models: `.gguf` format (e.g., z_image_turbo-Q4_K_M.gguf)
- LLM models: `.gguf` format (e.g., Qwen3-4B-Instruct-2507-Q4_K_M.gguf)
- VAE models: `.safetensors` format (e.g., diffusion_pytorch_model.safetensors)

### 3. Setup Library

**Option A: Automatic Setup (Recommended)**

The library will be automatically downloaded on first use. Just call `EnsureLibrary()` in your application:

```go
import "github.com/kawai-network/veridium/pkg/stablediffusion"

func main() {
    // Automatically download library if not present
    if err := stablediffusion.EnsureLibrary(); err != nil {
        log.Fatal("Failed to setup library:", err)
    }
    
    // Now you can use stable diffusion
    sd, err := stablediffusion.NewStableDiffusion(&stablediffusion.ContextParams{
        // ... your params
    })
}
```

The library will be stored in:
- **Development**: `./data/lib/`
- **Packaged app**: Platform-specific user data directory
  - macOS: `~/Library/Application Support/Kawai/lib/`
  - Windows: `%APPDATA%\Kawai\lib\`
  - Linux: `~/.config/Kawai/lib/`

**Option B: Manual Download**

The project **does not include** precompiled dynamic libraries. You need to download them manually from the official stable-diffusion.cpp repository.

**Download Instructions:**

1. Visit the releases page: https://github.com/leejet/stable-diffusion.cpp/releases
2. Download the appropriate version for your platform:

**For Windows:**
- **CPU versions** (choose based on your CPU capabilities):
  - `sd-master-xxx-bin-win-avx-x64.zip` - AVX instruction set
  - `sd-master-xxx-bin-win-avx2-x64.zip` - AVX2 instruction set (recommended for modern CPUs)
  - `sd-master-xxx-bin-win-avx512-x64.zip` - AVX512 instruction set
  - `sd-master-xxx-bin-win-noavx-x64.zip` - No AVX (for older CPUs)
- **GPU versions**:
  - `sd-master-xxx-bin-win-cuda12-x64.zip` - NVIDIA GPUs with CUDA 12
  - `sd-master-xxx-bin-win-vulkan-x64.zip` - Vulkan support (NVIDIA/AMD/Intel)
  - `sd-master-xxx-bin-win-rocm-x64.zip` - AMD GPUs with ROCm

**For Linux:**
- Download `sd-master-xxx-bin-ubuntu-x64.zip` or compile from source

**For macOS:**
- Download `sd-master-xxx-bin-Darwin-macOS-15.7.3-arm64.zip` (Apple Silicon) or `sd-master-xxx-bin-Darwin-macOS-x64.zip` (Intel)

3. Extract the downloaded archive and place the dynamic library files in the appropriate directory:
   - Development: `./data/lib/`
   - Production: Managed by `internal/paths` package
4. The program will automatically detect and load the appropriate library based on your system

**Note**: If you're using AMD graphics cards on Windows, download the ROCm version. For non-NVIDIA GPUs, the Vulkan version is recommended for GPU acceleration.

### 4. Run Examples

#### Text-to-Image Generation

```bash
# Enter the examples directory
cd examples

# Run text-to-image example
go run txt2img.go
```

Example code:

```go
package main

import (
	"fmt"
	"log"
	stablediffusion "github.com/kawai-network/veridium/pkg/stablediffusion"
)

func main() {
	fmt.Println("Stable Diffusion Go - Text to Image Example")
	fmt.Println("===============================================")

	// Ensure library is downloaded (automatic on first run)
	if err := stablediffusion.EnsureLibrary(); err != nil {
		log.Fatal("Failed to setup library:", err)
	}

	// Create Stable Diffusion instance
	sd, err := stablediffusion.NewStableDiffusion(&stablediffusion.ContextParams{
		DiffusionModelPath: "path/to/diffusion_model.gguf",
		LLMPath:            "path/to/llm_model.gguf",
		VAEPath:            "path/to/vae_model.safetensors",
		DiffusionFlashAttn: true,
		OffloadParamsToCPU: true,
	})

	if err != nil {
		fmt.Println("Failed to create instance:", err)
		return
	}
	defer sd.Free()

	// Generate image
	err = sd.GenerateImage(&stablediffusion.ImgGenParams{
		Prompt:      "一位穿着明朝服饰的美女行走在花园中",
		Width:       512,
		Height:      512,
		SampleSteps: 10,
		CfgScale:    1.0,
	}, "output.png")

	if err != nil {
		fmt.Println("Failed to generate image:", err)
		return
	}

	fmt.Println("Image generated successfully!")
}
```
![](output_demo.png)

#### Text-to-Video Generation

```bash
# Run text-to-video example
go run txt2vid.go
```

## 📚 Core Features

### 1. Context Management

- Create and destroy Stable Diffusion contexts
- Support multiple model path configurations
- Provide rich performance optimization parameters

### 2. Text-to-Image Generation (txt2img)

- Generate high-quality images from text descriptions
- Support Chinese and English prompts
- Adjustable image dimensions, sampling steps, CFG scale, and other parameters
- Support random seed generation

### 3. Text-to-Video Generation (txt2vid)

- Generate videos from text prompts
- Support custom frame count and resolution
- Support Easycache optimization
- Integrate FFmpeg for video encoding

## 🆕 New Features (v2.0)

### 1. Configuration File Support

Manage your settings using YAML or JSON configuration files:

```go
import "github.com/orangelang/stable-diffusion-go/pkg/config"

// Load configuration from file
cfg, err := config.LoadConfig("config.yaml")
if err != nil {
    cfg = config.DefaultConfig()
}

// Convert to ContextParams
ctxParams, _ := cfg.ToContextParams()
sd, _ := stablediffusion.NewStableDiffusion(ctxParams)
```

Example `config.yaml`:
```yaml
version: "1.0"

models:
  diffusion_model: "models/model.gguf"
  vae_model: "models/vae.safetensors"

generation:
  default_prompt: "masterpiece, best quality"
  default_width: 512
  default_height: 512
  default_sample_steps: 20

output:
  output_dir: "./outputs"
  naming_pattern: "sd_{timestamp}_{seed}.png"
```

### 2. Batch Image Generation

Generate multiple images efficiently:

```go
import "github.com/orangelang/stable-diffusion-go/pkg/batch"

// Create batch generator
batchGen := batch.NewGenerator(sd)

// Build batch parameters
params, _ := batch.NewBuilder().
    WithBaseParams(baseParams).
    AddSeedVariations("a beautiful landscape", 10).
    WithParallelism(4).
    Build()

// Generate with progress tracking
result, _ := batchGen.GenerateWithCallback(params, func(progress batch.BatchProgress) {
    fmt.Printf("Progress: %d/%d (%.1f%%)\n", 
        progress.Completed, progress.Total, progress.Percentage)
})
```

### 3. Progress Tracking

Monitor generation progress with detailed callbacks:

```go
import "github.com/orangelang/stable-diffusion-go/pkg/progress"

options := &progress.GenerationOptions{
    ProgressCallback: func(p progress.GenerationProgress) {
        fmt.Printf("\r[%s] %.1f%% - %s", 
            p.Stage, p.Percentage, p.Message)
    },
}

wrapper := progress.NewWrapper(sd, options)
err := wrapper.GenerateImage(params, "output.png")
```

### 4. Model Management

Organize and manage your models:

```go
import "github.com/orangelang/stable-diffusion-go/pkg/models"

// Create registry
registry := models.NewRegistry("models/registry.json")
registry.Load()

// Auto-detect models in directory
models.AutoRegister(registry, "./models", true)

// Search models
diffusionModels := registry.Search("turbo")

// Get model info
model, _ := registry.Get("model-id")
fmt.Printf("Model: %s, Size: %s\n", model.Name, model.HumanSize())
```

### 5. Pipeline Operations

Chain multiple operations together:

```go
import "github.com/orangelang/stable-diffusion-go/pkg/pipeline"

// Build pipeline
p := pipeline.New().
    Add(&pipeline.TextToImageStep{Params: imgParams}).
    Add(&pipeline.UpscaleStep{
        Params: upscalerParams, 
        Factor: 4,
    }).
    Add(&pipeline.SaveStep{Destination: "final.png"})

// Execute
ctx := pipeline.NewContext("./working")
ctx.SD = sd
err := p.Execute(ctx)
```

## 📝 Usage Guide

### Basic Usage

1. **Create Instance**: Use `NewStableDiffusion` to create a Stable Diffusion instance
2. **Configure Parameters**: Set context parameters and generation parameters
3. **Generate Content**: Call `GenerateImage` or `GenerateVideo` to generate content
4. **Release Resources**: Use `defer sd.Free()` to release resources

### Context Parameters Description

| Parameter Name | Type | Description |
|----------------|------|-------------|
| DiffusionModelPath | string | Diffusion model file path |
| LLMPath | string | LLM model file path |
| VAEPath | string | VAE model file path |
| NThreads | int | Number of threads |
| DiffusionFlashAttn | bool | Whether to enable FlashAttention |
| OffloadParamsToCPU | bool | Whether to offload some parameters to CPU |
| WType | SDType | Model quantization type |

### Image Generation Parameters Description

| Parameter Name | Type | Description |
|----------------|------|-------------|
| Prompt | string | Prompt text |
| NegativePrompt | string | Negative prompt text |
| Width | int | Image width |
| Height | int | Image height |
| Seed | int | Random seed |
| SampleSteps | int | Number of sampling steps |
| CfgScale | float64 | CFG scale |
| Strength | float64 | Initial image strength (img2img only) |

## 🔧 Performance Optimization

### 1. Adjust Thread Count

Adjust the `NThreads` parameter according to the number of CPU cores:

```go
ctxParams := &stablediffusion.ContextParams{
    // Other parameters...
    NThreads: 8, // Adjust according to CPU core count
}
```

### 2. Use Quantized Models

Using quantized models can improve performance and reduce memory usage:

```go
ctxParams := &stablediffusion.ContextParams{
    // Other parameters...
    WType: stablediffusion.SDTypeQ4_K, // Use Q4_K quantized model
}
```

### 3. Adjust Sampling Steps

Reducing the number of sampling steps can improve generation speed but may reduce image quality:

```go
imgGenParams := &stablediffusion.ImgGenParams{
    // Other parameters...
    SampleSteps: 10, // Reduce sampling steps
}
```

### 4. Enable FlashAttention

Enabling FlashAttention can accelerate the diffusion process:

```go
ctxParams := &stablediffusion.ContextParams{
    // Other parameters...
    DiffusionFlashAttn: true,
}
```

## ⚠️ Notes

1. **Dynamic Library Setup**: You must download the appropriate dynamic library for your platform (see section 3 above) and place it in an accessible location. The program will automatically detect and load the library based on your system configuration.
2. **Model Compatibility**: Ensure using model formats compatible with stable-diffusion.cpp
3. **Dependencies**: Install dependencies like CUDA or Vulkan as needed
4. **Video Generation**: Requires FFmpeg for video encoding
5. **Memory Usage**: Large models may require more memory, it is recommended to use quantized models
6. **About AMD Graphics Cards (Windows Platform)**: If using AMD graphics cards (including AMD integrated graphics), download the ROCm version from the releases page (see section 3 for download instructions)
7. **About Vulkan**: If using non-NVIDIA graphics cards (such as AMD or Intel graphics cards, including integrated graphics), you can use the Vulkan version to enable GPU acceleration

## 📦 Example Programs

### Text-to-Image Example

```go
package main

import (
	"fmt"
	stablediffusion "github.com/orangelang/stable-diffusion-go"
)

func main() {
	// Create instance
	sd, err := stablediffusion.NewStableDiffusion(&stablediffusion.ContextParams{
		DiffusionModelPath: "models/z_image_turbo-Q4_K_M.gguf",
		LLMPath:            "models/Qwen3-4B-Instruct-2507-Q4_K_M.gguf",
		VAEPath:            "models/diffusion_pytorch_model.safetensors",
		DiffusionFlashAttn: true,
	})
	if err != nil {
		fmt.Println("Failed to create instance:", err)
		return
	}
	defer sd.Free()

	// Generate image
	err = sd.GenerateImage(&stablediffusion.ImgGenParams{
		Prompt:      "A cute Corgi dog running on the grass",
		Width:       512,
		Height:      512,
		SampleSteps: 15,
		CfgScale:    2.0,
	}, "output_corgi.png")

	if err != nil {
		fmt.Println("Failed to generate image:", err)
		return
	}

	fmt.Println("Image generated successfully!")
}
```

### Text-to-Video Example

```go
package main

import (
	"fmt"
	stablediffusion "github.com/orangelang/stable-diffusion-go"
)

func main() {
	// Create instance
	sd, err := stablediffusion.NewStableDiffusion(&stablediffusion.ContextParams{
		DiffusionModelPath: "D:\\hf-mirror\\wan2.1\\wan2.1_t2v_1.3B_bf16.safetensors",
		T5XXLPath:          "D:\\hf-mirror\\wan2.1\\umt5-xxl-encoder-Q4_K_M.gguf",
		VAEPath:            "D:\\hf-mirror\\wan2.1\\wan_2.1_vae.safetensors",
		DiffusionFlashAttn: true,
		KeepClipOnCPU:      true,
		OffloadParamsToCPU: true,
		NThreads:           4,
		FlowShift:          3.0,
	})

	if err != nil {
		fmt.Println("Failed to create stable diffusion instance:", err)
		return
	}
	defer sd.Free()

	err = sd.GenerateVideo(&stablediffusion.VidGenParams{
		Prompt:      "一个在长满桃花树下拍照的美女",
		Width:       300,
		Height:      300,
		SampleSteps: 40,
		VideoFrames: 33,
		CfgScale:    6.0,
	}, "./output.mp4")

	if err != nil {
		fmt.Println("Failed to generate video:", err)
		return
	}

	fmt.Println("Video generated successfully!")
}
```

## 📄 License

MIT License

## 🤝 Contribution

Welcome to submit Issues and Pull Requests!

## 🔗 Related Projects

- [stable-diffusion.cpp](https://github.com/leejet/stable-diffusion.cpp): C++ implementation of Stable Diffusion model
- [purego](https://github.com/ebitengine/purego): Go language FFI library without cgo

## 📞 Support

If you encounter problems during use, please:
1. Check the example code
2. Check the dynamic library path and model files
3. Check project Issues
4. Submit a new Issue

---

Thank you for using stable-diffusion-go! If this project has helped you, please give us a Star ⭐️
