# Remote Image Generation

This package provides remote image generation capabilities using various AI APIs.

## Features

- **Gemini API**: Google's Gemini 2.5 Flash for high-quality image generation
- **Cloudflare Workers AI**: Fast and high-quality generation using Flux models
- **Pollinations AI**: Free fallback service with multiple models
- **Automatic Fallback**: Seamlessly switches between providers (Gemini -> Cloudflare -> Pollinations)
- **Aspect Ratio Support**: Intelligent aspect ratio calculation
- **Context Support**: Proper timeout and cancellation handling

## Usage

### Basic Usage

```go
import (
    "context"
    "github.com/kawai-network/veridium/pkg/stablediffusion/remote"
)

// Create generator
gen := remote.NewGenerator()

// Generate image
opts := remote.GenerationOptions{
    Prompt:     "A beautiful sunset over mountains",
    OutputPath: "/path/to/output.png",
    Width:      1024,
    Height:     1024,
}

err := gen.Generate(context.Background(), opts)
if err != nil {
    log.Fatal(err)
}
```

### Using Specific Generators

```go
// Use Gemini specifically
gemini := remote.NewGeminiGenerator()
err := gemini.Generate(ctx, opts)

// Use Cloudflare specifically
cloudflare := remote.NewCloudflareGenerator()
err := cloudflare.Generate(ctx, opts)

// Use Pollinations specifically
pollinations := remote.NewPollinationsGenerator()
err := pollinations.Generate(ctx, opts)
```

### With Fallback

```go
gen := remote.NewGenerator()

// Automatically tries Gemini first, falls back to Cloudflare, then Pollinations
err := gen.GenerateWithFallback(ctx, opts)
```

## Available Models

### Gemini Models
- `gemini-2.5-flash` - Fast generation, 1024px
- `gemini-2.5-flash-image` - Explicit image model

### Cloudflare Models
- `@cf/black-forest-labs/flux-1-schnell` - High quality, fast
- `@cf/black-forest-labs/flux-2-klein-9b` - Balanced quality and speed

### Pollinations Models
- `flux` - Default model
- `flux-realism` - Realistic images
- `flux-anime` - Anime style
- `flux-3d` - 3D rendered style
- `any-dark` - Dark themed
- `turbo` - Fast generation

## Configuration

### Aspect Ratios

The package automatically calculates aspect ratios from width/height:

- `1:1` - Square (1024x1024)
- `16:9` - Widescreen (1920x1080)
- `9:16` - Portrait (1080x1920)
- `4:3` - Standard (1024x768)
- `3:4` - Portrait (768x1024)
- And more...

You can also specify aspect ratio directly:

```go
opts := remote.GenerationOptions{
    Prompt:      "A landscape",
    AspectRatio: "16:9",
    OutputPath:  "output.png",
}
```

## API Keys

### Gemini API

Gemini API keys are managed through `internal/constant/llm.go`. The package automatically:
- Selects a random API key from the pool
- Handles rate limiting
- Falls back to Pollinations if no keys available

### Cloudflare API
Cloudflare credentials are managed through `internal/constant/llm.go` via `GetRandomCloudflareApiKey()`.
- Uses a pool of account IDs and API tokens
- Format: `ACCOUNT_ID:API_TOKEN`
- Automatically rotates between available keys

### Pollinations
No API key required - it's a free service.

## Error Handling

```go
err := gen.Generate(ctx, opts)
if err != nil {
    // Handle specific errors
    switch {
    case strings.Contains(err.Error(), "no Gemini API key"):
        // No API key available
    case strings.Contains(err.Error(), "rate limit"):
        // Rate limited
    case strings.Contains(err.Error(), "timeout"):
        // Request timed out
    default:
        // Other errors
    }
}
```

## Migration from internal/image

If you're migrating from `internal/image/remote.go`:

```go
// OLD
import "github.com/kawai-network/veridium/internal/image"
gen := image.NewRemoteGenerator()

// NEW
import "github.com/kawai-network/veridium/pkg/stablediffusion/remote"
gen := remote.NewGenerator()
```

The API is compatible, but the new package provides better separation of concerns.

## Thread Safety

All generators are thread-safe and can be used concurrently:

```go
gen := remote.NewGenerator()

var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)
    go func(index int) {
        defer wg.Done()
        opts := remote.GenerationOptions{
            Prompt:     fmt.Sprintf("Image %d", index),
            OutputPath: fmt.Sprintf("output_%d.png", index),
        }
        gen.Generate(context.Background(), opts)
    }(i)
}
wg.Wait()
```

## Testing

```bash
go test ./pkg/stablediffusion/remote/...
```

## License

MIT
