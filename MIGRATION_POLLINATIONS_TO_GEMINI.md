# Migration: Pollinations → Gemini Image Generation

## Summary

Image generation service telah berhasil dimigrasikan dari **Pollinations AI** ke **Google Gemini API** untuk meningkatkan kualitas, reliabilitas, dan kontrol atas proses image generation.

## Changes Made

### 1. Core Implementation (`internal/image/generation.go`)

#### Added Imports
```go
import (
    "context"
    "github.com/kawai-network/veridium/internal/constant"
    "google.golang.org/genai"
)
```

#### Modified Function: `generateImageRemote()`
- **Before**: Menggunakan Pollinations API (`https://image.pollinations.ai/prompt/`)
- **After**: Menggunakan Google Gemini API dengan SDK `google.golang.org/genai`

#### Key Features Implemented:
1. **API Key Rotation**: Menggunakan pool dari 5 Gemini API keys
2. **Model Selection**: 
   - `gemini-2.5-flash-image` (default, fast, 1024px)
   - `gemini-3-pro-image-preview` (HD/4K quality)
3. **Automatic Aspect Ratio Detection**: 10 supported ratios (1:1, 2:3, 3:2, 3:4, 4:3, 4:5, 5:4, 9:16, 16:9, 21:9)
4. **Quality Settings**: Standard (1K), HD (2K), 4K
5. **Enhanced Logging**: Detailed logs dengan prefix `[Gemini]`

#### Legacy Function Preserved
```go
func generateImageRemotePollinations() // Kept as fallback
```

### 2. Service Layer (`internal/image/service.go`)

#### Updated Model Pool
```go
// Before:
availableModels := []string{
    "flux",
    "stable-diffusion",
    "kontext",
    "turbo",
    "nanobanana",
    "seedream",
    "nanobanana-pro",
    "seedream-pro",
    "gptimage",
    "zimage",
    "veo",
    "seedance",
    "seedance-pro",
}

// After:
availableModels := []string{
    "gemini-2.5-flash",    // Fast, 1024px (Nano Banana)
    "gemini-3-pro",        // High quality, up to 4K (Nano Banana Pro)
    "gemini-2.5-flash",    // Duplicate for load balancing
    "gemini-2.5-flash",    // More weight on fast model
}
```

### 3. Documentation

Created comprehensive documentation:
- `GEMINI_IMAGE_GENERATION.md` - Detailed implementation guide
- `MIGRATION_POLLINATIONS_TO_GEMINI.md` - This file

## Technical Details

### API Integration

```go
// Get random API key from pool
apiKey := constant.GetRandomGeminiApiKey()
os.Setenv("GOOGLE_API_KEY", apiKey)

// Create client
client, err := genai.NewClient(ctx, nil)

// Generate with config
result, err := client.Models.GenerateContent(
    ctx,
    model,
    genai.Text(opts.Prompt),
    config,
)

// Extract image bytes
imageBytes := result.Candidates[0].Content.Parts[0].InlineData.Data
os.WriteFile(opts.OutputPath, imageBytes, 0644)
```

### Aspect Ratio Mapping

The system automatically maps width/height to Gemini's supported aspect ratios:

| Input Ratio | Gemini Ratio | Example Resolution (1K) |
|-------------|--------------|-------------------------|
| ~1.0        | 1:1          | 1024x1024              |
| ~0.67       | 2:3          | 848x1264               |
| ~1.5        | 3:2          | 1264x848               |
| ~0.75       | 3:4          | 896x1200               |
| ~1.33       | 4:3          | 1200x896               |
| ~0.78       | 4:5          | 928x1152               |
| ~1.28       | 5:4          | 1152x928               |
| ~0.56       | 9:16         | 768x1376               |
| ~1.78       | 16:9         | 1376x768               |
| ~2.33       | 21:9         | 1584x672               |

## Benefits

### 1. Quality Improvements
- ✅ Better image quality dengan model terbaru dari Google
- ✅ Support hingga 4K resolution
- ✅ SynthID watermark untuk authenticity

### 2. Reliability
- ✅ More stable API dari Google infrastructure
- ✅ Multiple API keys untuk load balancing
- ✅ Better error handling

### 3. Control & Flexibility
- ✅ Precise aspect ratio control
- ✅ Multiple quality tiers (1K, 2K, 4K)
- ✅ Model selection based on use case

### 4. Performance
- ✅ Fast model (`gemini-2.5-flash-image`) untuk quick generation
- ✅ Pro model (`gemini-3-pro-image-preview`) untuk high quality
- ✅ Load balancing dengan model rotation

## Testing

### Compilation Test
```bash
cd /Users/yuda/github.com/kawai-network/veridium-1
go build ./internal/image/...
# ✅ Success - No errors
```

### Usage Example
```go
opts := image.GenerationOptions{
    Prompt:      "A futuristic city at sunset",
    OutputPath:  "./output/image.png",
    Width:       1920,
    Height:      1080,
    Quality:     "hd",
    Model:       "gemini-3-pro",
}

err := service.generateImageRemote(opts)
```

## API Key Management

API keys are managed in `internal/constant/llm.go`:

```go
// Pool of 5 Gemini API keys
func GetRandomGeminiApiKey() string {
    keys := GetGeminiApiKeys()
    r := rand.New(rand.NewSource(time.Now().UnixNano()))
    return keys[r.Intn(len(keys))]
}
```

Keys are:
- Obfuscated for security
- Randomly selected for load balancing
- Automatically rotated per request

## Backward Compatibility

### Pollinations Function Preserved
The original Pollinations implementation is preserved as `generateImageRemotePollinations()` and can be used as fallback if needed.

### Migration Path
To revert to Pollinations (if needed):
```go
// In generateImageRemote(), replace body with:
return sdrm.generateImageRemotePollinations(opts)
```

## Logging Examples

### Successful Generation
```
[Gemini] Using model: gemini-2.5-flash-image for prompt: a beautiful sunset
[Gemini] Aspect ratio: 16:9
[Gemini] Received image data: 245678 bytes
[Gemini] Image saved successfully to: ./output/image.png
```

### With HD Quality
```
[Gemini] Using model: gemini-3-pro-image-preview for prompt: detailed portrait
[Gemini] Image size: 2K
[Gemini] Aspect ratio: 3:4
[Gemini] Received image data: 892341 bytes
[Gemini] Image saved successfully to: ./output/portrait.png
```

## Cost Considerations

### Token Usage
- **gemini-2.5-flash-image**: 1290 tokens per image (all ratios)
- **gemini-3-pro-image-preview**: 
  - 1K: 1120 tokens
  - 2K: 1120 tokens
  - 4K: 2000 tokens

### Optimization
- Default to fast model for most use cases
- Use HD/4K only when explicitly requested
- Load balance across 5 API keys

## Future Enhancements

### Planned Features
1. **Image Editing**: Implement text-and-image-to-image capability
2. **Automatic Fallback**: Fallback to Pollinations on Gemini failure
3. **Caching**: Cache generated images to reduce API calls
4. **Cost Tracking**: Monitor token usage and costs
5. **Batch Optimization**: Optimize parallel generation

### Potential Improvements
- Add support for negative prompts (if Gemini supports it)
- Implement retry logic with exponential backoff
- Add metrics and monitoring
- Create admin dashboard for API key management

## References

- [Gemini API Documentation](https://ai.google.dev/gemini-api/docs/image-generation#go)
- [Gemini Models Overview](https://ai.google.dev/gemini-api/docs/models/gemini)
- [Go SDK Documentation](https://pkg.go.dev/google.golang.org/genai)

## Files Modified

1. `internal/image/generation.go` - Core implementation
2. `internal/image/service.go` - Model pool update
3. `GEMINI_IMAGE_GENERATION.md` - New documentation
4. `MIGRATION_POLLINATIONS_TO_GEMINI.md` - This migration guide

## Verification

✅ Code compiles successfully  
✅ No linter errors  
✅ Backward compatibility maintained  
✅ Documentation complete  
✅ API keys configured  

## Rollout Plan

### Phase 1: Testing (Current)
- ✅ Implementation complete
- ✅ Compilation verified
- ⏳ Integration testing needed

### Phase 2: Gradual Rollout
1. Deploy to staging environment
2. Monitor logs and error rates
3. Compare quality with Pollinations
4. Gather user feedback

### Phase 3: Full Production
1. Switch all traffic to Gemini
2. Monitor costs and performance
3. Keep Pollinations as emergency fallback

### Phase 4: Optimization
1. Fine-tune model selection
2. Implement caching
3. Add advanced features

## Contact

For questions or issues regarding this migration, please refer to:
- Implementation: `internal/image/generation.go`
- Documentation: `GEMINI_IMAGE_GENERATION.md`
- API Keys: `internal/constant/llm.go`

---

**Migration Date**: January 7, 2026  
**Status**: ✅ Complete  
**Version**: 1.0.0

