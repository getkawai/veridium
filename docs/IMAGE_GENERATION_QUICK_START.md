# Image Generation Quick Start Guide

## 🚀 Quick Start

### Basic Usage

```go
import "github.com/kawai-network/veridium/internal/image"

// Create service
service := image.NewService(db, stableDiffusion, topicService)

// Generate image
req := image.CreateImageRequest{
    GenerationTopicId: "topic-123",
    Provider:          "gemini",
    Model:             "gemini-2.5-flash",
    ImageNum:          4,
    Params: image.RuntimeImageGenParams{
        Prompt:      "A beautiful sunset over mountains",
        Width:       ptr(1920),
        Height:      ptr(1080),
        AspectRatio: "16:9",
        Quality:     "standard",
    },
}

err := service.CreateImage(req)
```

## 🎨 Model Selection

### Fast Generation (Default)
```go
Model: "gemini-2.5-flash"  // 1024px, fast
```

### High Quality (Flux)
```go
Model: "@cf/black-forest-labs/flux-1-schnell"  // High quality, fast
Model: "@cf/black-forest-labs/flux-2-klein-9b" // Balanced quality/speed
```

### High Quality (Gemini)
```go
Model: "gemini-3-pro"      // Up to 4K, slower
Quality: "hd"              // 2K resolution
```

### Ultra High Quality
```go
Model: "gemini-3-pro"
Quality: "4k"              // 4K resolution
```

## 📐 Aspect Ratios

### Common Ratios
```go
AspectRatio: "1:1"    // Square - 1024x1024
AspectRatio: "16:9"   // Widescreen - 1376x768
AspectRatio: "9:16"   // Mobile - 768x1376
AspectRatio: "4:3"    // Classic - 1200x896
AspectRatio: "3:4"    // Portrait - 896x1200
```

### All Supported Ratios
- `1:1` - Square
- `2:3` / `3:2` - Photo
- `3:4` / `4:3` - Classic
- `4:5` / `5:4` - Portrait/Landscape
- `9:16` / `16:9` - Mobile/Widescreen
- `21:9` - Ultrawide

## 💡 Examples

### Example 1: Quick Social Media Image
```go
params := image.RuntimeImageGenParams{
    Prompt:      "Modern tech startup office, vibrant colors",
    AspectRatio: "1:1",  // Instagram square
}
```

### Example 2: HD Wallpaper
```go
params := image.RuntimeImageGenParams{
    Prompt:      "Cyberpunk city at night, neon lights",
    AspectRatio: "16:9",
    Quality:     "hd",
    Model:       "gemini-3-pro",
}
```

### Example 3: Mobile App Screenshot
```go
params := image.RuntimeImageGenParams{
    Prompt:      "Clean mobile app UI, minimalist design",
    AspectRatio: "9:16",
}
```

### Example 4: Ultra HD Print
```go
params := image.RuntimeImageGenParams{
    Prompt:      "Detailed landscape painting, photorealistic",
    AspectRatio: "3:2",
    Quality:     "4k",
    Model:       "gemini-3-pro",
}
```

## 🔧 Advanced Options

### Custom Dimensions
```go
// System will auto-detect aspect ratio
params := image.RuntimeImageGenParams{
    Prompt: "Custom sized image",
    Width:  ptr(1920),
    Height: ptr(1080),
}
```

### Multiple Images
```go
req := image.CreateImageRequest{
    ImageNum: 4,  // Generate 4 variations
    // ... other params
}
```

## 📊 Resolution Reference

### Gemini 2.5 Flash (Fast)

| Ratio | Resolution | Use Case |
|-------|-----------|----------|
| 1:1   | 1024x1024 | Social media |
| 16:9  | 1344x768  | Presentations |
| 9:16  | 768x1344  | Mobile screens |
| 4:3   | 1184x864  | Classic displays |

### Gemini 3 Pro (HD - 2K)

| Ratio | Resolution | Use Case |
|-------|-----------|----------|
| 1:1   | 2048x2048 | High-res social |
| 16:9  | 2752x1536 | HD wallpapers |
| 9:16  | 1536x2752 | Mobile HD |
| 3:2   | 2528x1696 | Photo prints |

### Gemini 3 Pro (4K)

| Ratio | Resolution | Use Case |
|-------|-----------|----------|
| 1:1   | 4096x4096 | Ultra HD prints |
| 16:9  | 5504x3072 | 4K displays |
| 21:9  | 6336x2688 | Cinema |

## ⚡ Performance Tips

### 1. Choose Right Model
- Use `gemini-2.5-flash` for quick iterations
- Use `gemini-3-pro` only when quality matters

### 2. Batch Generation
```go
ImageNum: 4  // Generate 4 at once for variations
```

### 3. Appropriate Quality
- `standard` for web/social media
- `hd` for presentations/prints
- `4k` only for professional use

## 🔍 Monitoring

### Check Logs
```bash
# Look for Gemini logs
grep "\[Gemini\]" logs/app.log

# Example output:
# [Gemini] Using model: gemini-2.5-flash-image
# [Gemini] Aspect ratio: 16:9
# [Gemini] Received image data: 245678 bytes
# [Gemini] Image saved successfully
```

### Error Handling
```go
err := service.CreateImage(req)
if err != nil {
    log.Printf("Generation failed: %v", err)
    // Errors are logged with details
    // Check database for async task status
}
```

## 📝 Prompt Tips

### Good Prompts
```go
✅ "A modern minimalist office with plants and natural light"
✅ "Cyberpunk city street at night with neon signs and rain"
✅ "Abstract geometric pattern in blue and gold colors"
```

### Avoid
```go
❌ "image"  // Too vague
❌ "nice picture"  // Not descriptive
❌ "something cool"  // Unclear
```

### Best Practices
1. Be specific and descriptive
2. Include style/mood keywords
3. Mention colors if important
4. Specify lighting/atmosphere
5. Keep it under 200 characters

## 🔐 API Keys

API keys are automatically managed:
- **Gemini**: Pool of 5 keys for load balancing.
- **Cloudflare**: Account IDs and Tokens managed with rotation.
- Configured in `internal/constant/llm.go`.

No manual key management needed! 🎉

## 🐛 Troubleshooting

### Issue: "No API key available" (Gemini/Cloudflare)
**Solution**: Check `internal/constant/llm.go` has valid keys configured.

### Issue: "No image data returned"
**Solution**: 
- Check prompt is valid.
- Verify API keys are working.
- Check logs for provider errors (e.g., `grep "[RemoteGen]" logs/app.log`).
### Issue: Generation takes too long
**Solution**:
- Use `gemini-2.5-flash` instead of `gemini-3-pro`
- Reduce `ImageNum` if generating multiple
- Check if API is rate-limited

### Issue: Low quality images
**Solution**:
- Switch to `gemini-3-pro` model
- Set `Quality: "hd"` or `"4k"`
- Improve prompt with more details

## 📚 More Resources

- [Full Documentation](../GEMINI_IMAGE_GENERATION.md)
- [Migration Guide](../MIGRATION_POLLINATIONS_TO_GEMINI.md)
- [Gemini API Docs](https://ai.google.dev/gemini-api/docs/image-generation)

## 🎯 Common Use Cases

### Social Media Content
```go
// Instagram post
AspectRatio: "1:1", Model: "gemini-2.5-flash"

// Instagram story
AspectRatio: "9:16", Model: "gemini-2.5-flash"

// Twitter header
AspectRatio: "3:1", Model: "gemini-2.5-flash"
```

### Web Design
```go
// Hero image
AspectRatio: "21:9", Quality: "hd"

// Blog thumbnail
AspectRatio: "16:9", Model: "gemini-2.5-flash"

// Product image
AspectRatio: "1:1", Quality: "hd"
```

### Print Materials
```go
// Poster
AspectRatio: "2:3", Quality: "4k", Model: "gemini-3-pro"

// Flyer
AspectRatio: "3:4", Quality: "hd", Model: "gemini-3-pro"

// Business card
AspectRatio: "16:9", Quality: "hd"
```

---

**Need Help?** Check the full documentation or logs for detailed error messages.

