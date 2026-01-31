# Image Generation Architecture

This package provides image generation capabilities with support for both local (Stable Diffusion binary) and remote (API-based) generation.

## 📁 File Structure

```
internal/image/
├── service.go              # High-level orchestration service
├── local.go                # Local SD binary execution
├── remote.go               # Remote API generation (Gemini, Pollinations)
├── generation.go           # Common types and interfaces
├── manager.go              # Binary lifecycle management
├── models.go               # Model specifications
└── manager_*.go            # Platform-specific binary selection
```

## 🏗️ Architecture

### Separation of Concerns

The package is organized into three main layers:

1. **Service Layer** (`service.go`)
   - High-level orchestration
   - Database integration
   - Async processing with goroutines
   - Topic/WebSocket notifications
   - Batch management

2. **Generation Layer** (`local.go`, `remote.go`)
   - **LocalGenerator**: Executes local Stable Diffusion binary
   - **RemoteGenerator**: Calls remote APIs (Gemini, Pollinations)
   - Clean separation allows easy testing and provider switching

3. **Management Layer** (`manager.go`, `manager_*.go`)
   - Binary download from GitHub releases
   - Platform-specific binary selection
   - Process lifecycle management

## 🎯 Usage Examples

### Using Local Generation

```go
// Create local generator
engine := image.NewEngine()
localGen := image.NewLocalGenerator(engine)

// Check if binary is available
if !localGen.IsAvailable() {
    log.Fatal("SD binary not found")
}

// Generate image
opts := image.GenerationOptions{
    Prompt:     "a beautiful landscape",
    ModelPath:  localGen.GetFirstAvailableModel(),
    OutputPath: "output.png",
    Width:      1024,
    Height:     1024,
    Steps:      20,
    Cfg:        7.0,
}

err := localGen.Generate(context.Background(), opts)
```

### Using Remote Generation

```go
// Create remote generator
remoteGen := image.NewRemoteGenerator()

// Check if API keys are available
if !remoteGen.IsAvailable() {
    log.Fatal("No API keys configured")
}

// Generate image
opts := image.GenerationOptions{
    Prompt:      "a beautiful landscape",
    Model:       "gemini-2.5-flash",
    OutputPath:  "output.png",
    Width:       1024,
    Height:      1024,
    AspectRatio: "16:9",
}

err := remoteGen.Generate(context.Background(), opts)
```

### Using Service (Recommended for Web Apps)

```go
// Create service with database integration
db := database.NewService(...)
engine := image.NewEngine()
service := image.NewService(db, engine)

// Set topic service for real-time updates
service.SetTopicService(topicService)

// Create image request (async)
req := image.CreateImageRequest{
    GenerationTopicId: "topic-123",
    Provider:          "remote", // or "local"
    Model:            "gemini-2.5-flash",
    ImageNum:         4,
    Params: image.RuntimeImageGenParams{
        Prompt: "a beautiful landscape",
        Width:  ptr(1024),
        Height: ptr(1024),
    },
}

// Returns immediately, generation happens in background
err := service.CreateImage(req)

// Client receives updates via WebSocket/SSE
```

## 🔄 Migration from Old Code

### Before (monolithic generation.go)

```go
// Everything in one file
func (sdrm *StableDiffusion) generateImageRemote(opts GenerationOptions) error {
    // 200+ lines of Gemini API code
}

func (sdrm *StableDiffusion) generateImageRemotePollinations(opts GenerationOptions) error {
    // 60+ lines of Pollinations code
}

func (sdrm *StableDiffusion) createImageInternal(opts GenerationOptions) error {
    // 80+ lines of binary execution code
}
```

### After (separated concerns)

```go
// local.go - 150 lines
type LocalGenerator struct { ... }
func (lg *LocalGenerator) Generate(ctx context.Context, opts GenerationOptions) error { ... }

// remote.go - 250 lines
type RemoteGenerator struct { ... }
func (rg *RemoteGenerator) Generate(ctx context.Context, opts GenerationOptions) error { ... }

// generation.go - 100 lines (types + wrappers only)
type GenerationOptions struct { ... }
func (sdrm *StableDiffusion) CreateImageWithOptions(opts GenerationOptions) error {
    remoteGen := NewRemoteGenerator()
    return remoteGen.Generate(context.Background(), opts)
}
```

## ✅ Benefits

1. **Better Organization**
   - Each file has a single, clear responsibility
   - Easier to navigate and understand

2. **Easier Testing**
   - Can mock LocalGenerator or RemoteGenerator independently
   - Test each provider in isolation

3. **Scalability**
   - Easy to add new providers (OpenAI DALL-E, Midjourney, etc.)
   - Just create a new generator type

4. **Maintainability**
   - Changes to Gemini API don't affect local generation
   - Changes to SD binary don't affect remote APIs

5. **Clear Dependencies**
   - `local.go` only depends on binary execution
   - `remote.go` only depends on HTTP/API clients
   - No circular dependencies

## 🔧 Adding a New Provider

To add a new remote provider (e.g., OpenAI DALL-E):

1. Add method to `RemoteGenerator` in `remote.go`:

```go
func (rg *RemoteGenerator) generateWithDALLE(ctx context.Context, opts GenerationOptions) error {
    // DALL-E API implementation
}
```

2. Update `Generate()` method to route to new provider:

```go
func (rg *RemoteGenerator) Generate(ctx context.Context, opts GenerationOptions) error {
    switch opts.Model {
    case "dall-e-3":
        return rg.generateWithDALLE(ctx, opts)
    case "gemini-2.5-flash":
        return rg.generateWithGemini(ctx, opts)
    default:
        return rg.generateWithPollinations(ctx, opts)
    }
}
```

3. Add model to `GetAvailableModels()`:

```go
func (rg *RemoteGenerator) GetAvailableModels() []string {
    return []string{
        "dall-e-3",
        "gemini-2.5-flash",
        "flux",
        // ...
    }
}
```

## 📝 Notes

- **Backward Compatibility**: Old methods like `generateImageRemote()` are kept as wrappers for backward compatibility
- **Context Support**: All new methods accept `context.Context` for proper cancellation and timeout handling
- **Error Handling**: Each generator returns descriptive errors with provider-specific context
- **Logging**: Consistent logging format with provider prefixes (`[Gemini]`, `[LocalSD]`, etc.)
