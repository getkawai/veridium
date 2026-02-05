# Local Stable Diffusion Generation

This package provides local image generation using Stable Diffusion binary (stable-diffusion.cpp).

## Features

- **Binary Execution**: Direct execution of stable-diffusion.cpp binary
- **Model Support**: .ckpt, .safetensors, .pt, .bin, .gguf formats
- **Cross-Platform**: Works on macOS, Linux, Windows
- **Customizable**: Full control over generation parameters
- **Context Support**: Proper timeout and cancellation handling

## Usage

### Basic Usage

```go
import (
    "context"
    "github.com/kawai-network/veridium/pkg/stablediffusion/local"
)

// Create generator
gen := local.NewGenerator(
    "/path/to/sd",           // Binary path
    "/path/to/models",       // Models directory
)

// Check if available
if !gen.IsAvailable() {
    log.Fatal("SD binary not found")
}

// Get first available model
modelPath := gen.GetFirstAvailableModel()
if modelPath == "" {
    log.Fatal("No models found")
}

// Generate image
opts := local.GenerationOptions{
    Prompt:     "A beautiful sunset over mountains",
    ModelPath:  modelPath,
    OutputPath: "/path/to/output.png",
    Width:      1024,
    Height:     1024,
    Steps:      20,
    Cfg:        7.0,
}

err := gen.Generate(context.Background(), opts)
if err != nil {
    log.Fatal(err)
}
```

### With Custom Executor

```go
// Create custom executor (useful for testing)
type MockExecutor struct{}

func (m *MockExecutor) Run(ctx context.Context, name string, args ...string) error {
    // Custom execution logic
    return nil
}

gen := local.NewGeneratorWithExecutor(
    "/path/to/sd",
    "/path/to/models",
    &MockExecutor{},
)
```

### Advanced Options

```go
seed := int64(42)
imageUrl := "/path/to/input.png"

opts := local.GenerationOptions{
    Prompt:         "A cat in space",
    NegativePrompt: "blurry, low quality",
    ModelPath:      modelPath,
    OutputPath:     "output.png",
    Width:          512,
    Height:         512,
    Steps:          30,
    Cfg:            7.5,
    Seed:           &seed,
    SamplerName:    "euler_a",
    Scheduler:      "karras",
    
    // For img2img
    ImageUrl:       &imageUrl,
    Strength:       0.75,
}

err := gen.Generate(ctx, opts)
```

## Supported Model Formats

- `.ckpt` - Checkpoint files
- `.safetensors` - SafeTensors format
- `.pt` - PyTorch files
- `.bin` - Binary format
- `.gguf` - GGUF quantized models

## Binary Setup

The SD binary should be installed via `pkg/stablediffusion/setup.go`:

```go
import "github.com/kawai-network/veridium/pkg/stablediffusion"

// Ensure binary is installed
if err := stablediffusion.EnsureLibrary(); err != nil {
    log.Fatal(err)
}

// Get binary path
binaryPath := stablediffusion.GetLibraryPath()
```

## Environment Variables

The package automatically sets up environment variables for dynamic library loading:

- **macOS**: `DYLD_LIBRARY_PATH`
- **Linux**: `LD_LIBRARY_PATH`
- **Windows**: `PATH`

## Error Handling

```go
err := gen.Generate(ctx, opts)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "binary not found"):
        // Binary not installed
    case strings.Contains(err.Error(), "output file was not created"):
        // Generation failed
    case strings.Contains(err.Error(), "command execution failed"):
        // Execution error
    default:
        // Other errors
    }
}
```

## Testing

```bash
go test ./pkg/stablediffusion/local/...
```

## Integration with internal/image

The old `internal/image/local.go` now wraps this package for backward compatibility:

```go
// OLD (still works)
import "github.com/kawai-network/veridium/internal/image"
gen := image.NewLocalGenerator(engine)

// NEW (recommended)
import "github.com/kawai-network/veridium/pkg/stablediffusion/local"
gen := local.NewGenerator(binaryPath, modelsPath)
```

## Thread Safety

The generator is thread-safe and can be used concurrently:

```go
gen := local.NewGenerator(binaryPath, modelsPath)

var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(index int) {
        defer wg.Done()
        opts := local.GenerationOptions{
            Prompt:     fmt.Sprintf("Image %d", index),
            ModelPath:  modelPath,
            OutputPath: fmt.Sprintf("output_%d.png", index),
        }
        gen.Generate(context.Background(), opts)
    }(i)
}
wg.Wait()
```

## Performance Tips

1. **Reuse Generator**: Create once, use multiple times
2. **Model Selection**: Use quantized models (.gguf) for faster generation
3. **Step Count**: Lower steps = faster but lower quality
4. **Resolution**: Lower resolution = faster generation

## License

MIT
