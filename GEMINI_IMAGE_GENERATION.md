# Gemini Image Generation Implementation

## Overview

Image generation telah diupgrade dari Pollinations AI ke **Google Gemini API** untuk kualitas dan performa yang lebih baik.

## Perubahan Utama

### 1. Provider Baru: Google Gemini

Implementasi sekarang menggunakan [Gemini API](https://ai.google.dev/gemini-api/docs/image-generation#go) dengan dua model:

- **gemini-2.5-flash-image** (Nano Banana) - Cepat, resolusi 1024px
- **gemini-3-pro-image-preview** (Nano Banana Pro) - Kualitas tinggi, hingga 4K

### 2. File yang Dimodifikasi

**`internal/image/generation.go`**
- Fungsi `generateImageRemote()` sekarang menggunakan Gemini API
- Fungsi lama Pollinations dipindah ke `generateImageRemotePollinations()` sebagai fallback
- Import baru: `google.golang.org/genai` dan `github.com/kawai-network/x/constant`

### 3. Fitur Baru

#### Automatic Aspect Ratio Detection
Sistem secara otomatis mendeteksi aspect ratio dari width/height yang diberikan:
- 1:1 (square)
- 2:3, 3:2 (portrait/landscape)
- 3:4, 4:3
- 4:5, 5:4
- 9:16, 16:9 (mobile/widescreen)
- 21:9 (ultrawide)

#### Quality Settings
- **Standard**: `gemini-2.5-flash-image` dengan resolusi 1024px
- **HD**: `gemini-3-pro-image-preview` dengan resolusi 2K
- **4K**: `gemini-3-pro-image-preview` dengan resolusi 4K (set `Quality: "4k"`)

#### Model Selection
```go
// Default: gemini-2.5-flash-image (fast)
opts := GenerationOptions{
    Prompt: "a beautiful sunset",
}

// High quality: gemini-3-pro-image-preview
opts := GenerationOptions{
    Prompt: "a beautiful sunset",
    Model: "gemini-3-pro",
    Quality: "hd", // or "4k"
}
```

## API Key Management

API keys diambil secara random dari pool yang ada di `internal/constant/llm.go`:

```go
apiKey := constant.GetRandomGeminiApiKey()
```

Pool saat ini memiliki 5 Gemini API keys yang akan di-rotate secara otomatis untuk load balancing.

### Direct API Key Usage

API key dikirim langsung ke client tanpa menggunakan environment variable:

```go
clientConfig := &genai.ClientConfig{
    APIKey:  apiKey,
    Backend: genai.BackendGeminiAPI,
}

client, err := genai.NewClient(ctx, clientConfig)
```

**Benefits**:
- ✅ Thread-safe by design (no environment variable race conditions)
- ✅ No mutex needed
- ✅ Cleaner code
- ✅ Each request uses its own API key from the pool

## Aspect Ratio & Resolution Table

### Gemini 2.5 Flash Image (Nano Banana)

| Aspect Ratio | Resolution | Tokens |
|--------------|------------|--------|
| 1:1          | 1024x1024  | 1290   |
| 2:3          | 832x1248   | 1290   |
| 3:2          | 1248x832   | 1290   |
| 3:4          | 864x1184   | 1290   |
| 4:3          | 1184x864   | 1290   |
| 4:5          | 896x1152   | 1290   |
| 5:4          | 1152x896   | 1290   |
| 9:16         | 768x1344   | 1290   |
| 16:9         | 1344x768   | 1290   |
| 21:9         | 1536x672   | 1290   |

### Gemini 3 Pro Image Preview (Nano Banana Pro)

| Aspect Ratio | 1K Resolution | 2K Resolution | 4K Resolution |
|--------------|---------------|---------------|---------------|
| 1:1          | 1024x1024     | 2048x2048     | 4096x4096     |
| 2:3          | 848x1264      | 1696x2528     | 3392x5056     |
| 3:2          | 1264x848      | 2528x1696     | 5056x3392     |
| 3:4          | 896x1200      | 1792x2400     | 3584x4800     |
| 4:3          | 1200x896      | 2400x1792     | 4800x3584     |
| 4:5          | 928x1152      | 1856x2304     | 3712x4608     |
| 5:4          | 1152x928      | 2304x1856     | 4608x3712     |
| 9:16         | 768x1376      | 1536x2752     | 3072x5504     |
| 16:9         | 1376x768      | 2752x1536     | 5504x3072     |
| 21:9         | 1584x672      | 3168x1344     | 6336x2688     |

## Usage Example

```go
opts := image.GenerationOptions{
    Prompt:      "A futuristic city at sunset with flying cars",
    OutputPath:  "./output/image.png",
    Width:       1920,
    Height:      1080,
    Quality:     "hd",
    AspectRatio: "16:9",
}

err := stableDiffusion.generateImageRemote(opts)
if err != nil {
    log.Printf("Error: %v", err)
}
```

## Benefits

### 1. **Better Quality**
- Gemini menggunakan model terbaru dengan kualitas output yang lebih baik
- Support hingga 4K resolution

### 2. **More Reliable**
- API yang lebih stabil dari Google
- Built-in retry mechanism dengan multiple API keys

### 3. **Better Control**
- Precise aspect ratio control
- Multiple quality settings
- SynthID watermark untuk authenticity

### 4. **Cost Effective**
- Token-based pricing yang transparan
- Load balancing dengan multiple API keys

## Backward Compatibility

Fungsi Pollinations masih tersedia sebagai `generateImageRemotePollinations()` jika diperlukan sebagai fallback.

## Logging

Semua operasi Gemini di-log dengan prefix `[Gemini]`:

```
[Gemini] Using model: gemini-2.5-flash-image for prompt: a beautiful sunset
[Gemini] Aspect ratio: 16:9
[Gemini] Received image data: 245678 bytes
[Gemini] Image saved successfully to: ./output/image.png
```

## Error Handling

Implementasi mencakup error handling untuk:
- Missing API keys
- Failed API calls
- No image data returned
- File write errors
- Context timeout (120 seconds)

## Future Enhancements

Potential improvements:
1. Image editing (text-and-image-to-image)
2. Batch generation optimization
3. Caching mechanism
4. Automatic fallback to Pollinations on Gemini failure
5. Cost tracking and monitoring

## References

- [Gemini API Documentation](https://ai.google.dev/gemini-api/docs/image-generation#go)
- [Gemini Models Overview](https://ai.google.dev/gemini-api/docs/models/gemini)
- Original implementation: `internal/image/generation.go`

