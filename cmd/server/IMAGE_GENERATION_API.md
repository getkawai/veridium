# Image Generation API

## Overview

Kronk server sekarang mendukung OpenAI-compatible image generation API menggunakan Stable Diffusion.

## Endpoint

### POST /v1/images/generations

Generate images dari text prompt.

**Authentication**: Wallet-based (API key required)

**Request Body**:
```json
{
  "prompt": "a cute cat wearing sunglasses",
  "n": 1,
  "size": "1024x1024",
  "quality": "standard",
  "response_format": "url"
}
```

**Parameters**:
- `prompt` (required): Text description of the image
- `n` (optional): Number of images to generate (default: 1)
- `size` (optional): Image dimensions (default: "1024x1024")
- `quality` (optional): "standard" or "hd" (default: "standard")
  - `hd`: Uses 30 steps for higher quality
  - `standard`: Uses 20 steps
- `response_format` (optional): "url" or "b64_json" (default: "url")
  - `url`: Returns file URL
  - `b64_json`: Returns base64 encoded image
- `model` (optional): Specific model to use (uses first available if not specified)

**Response**:
```json
{
  "created": 1706745600,
  "data": [
    {
      "url": "/files/gen_abc123.png",
      "revised_prompt": "a cute cat wearing sunglasses"
    }
  ]
}
```

## File Server

Generated images are served via:

```
GET /files/{filename}
```

No authentication required for accessing generated images.

## Example Usage

### cURL

```bash
# Generate single image
curl -X POST http://localhost:42139/v1/images/generations \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "a beautiful sunset over mountains",
    "size": "1024x1024",
    "quality": "hd"
  }'

# Generate multiple images
curl -X POST http://localhost:42139/v1/images/generations \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "abstract art with vibrant colors",
    "n": 3,
    "response_format": "url"
  }'

# Get base64 encoded image
curl -X POST http://localhost:42139/v1/images/generations \
  -H "Authorization: Bearer YOUR_API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "futuristic city at night",
    "response_format": "b64_json"
  }'
```

### Python (OpenAI SDK)

```python
from openai import OpenAI

client = OpenAI(
    base_url="http://localhost:42139/v1",
    api_key="YOUR_API_KEY"
)

response = client.images.generate(
    prompt="a cute cat wearing sunglasses",
    n=1,
    size="1024x1024"
)

print(response.data[0].url)
```

### JavaScript

```javascript
const response = await fetch('http://localhost:42139/v1/images/generations', {
  method: 'POST',
  headers: {
    'Authorization': 'Bearer YOUR_API_KEY',
    'Content-Type': 'application/json'
  },
  body: JSON.stringify({
    prompt: 'a cute cat wearing sunglasses',
    n: 1,
    size: '1024x1024'
  })
});

const data = await response.json();
console.log(data.data[0].url);
```

## Implementation Details

### Architecture

```
POST /v1/images/generations
  ↓
imageapp.generations()
  ↓
generateImages() - Reused from pkg/gateway
  ↓
image.StableDiffusion.CreateImageWithOptions()
  ↓
Generated files saved to outputs/
  ↓
Response with /files/{filename} URL
```

### Code Reuse

Implementation reuses existing code:
- **Types**: `ImageGenerationRequest`, `ImageGenerationResponse`, `ImageData` (from `pkg/gateway/image_types.go`)
- **Logic**: Image generation logic adapted from `pkg/gateway/image_executor.go`
- **Engine**: Uses existing `internal/image.StableDiffusion` engine

### File Storage

- Generated images stored in: `{engine.GetOutputsPath()}/gen_{uuid}.png`
- Filename format: `gen_{uuid}.png`
- Files persist on disk for caching/debugging
- Served via `/files/*` route

### Security

- **Authentication**: Wallet-based API key required for generation
- **File Access**: Public access to generated files (no auth on `/files/*`)
- **Path Traversal**: Protected with `filepath.Clean()` and prefix validation

## Configuration

No additional configuration needed. Image generation is automatically enabled when:
1. Stable Diffusion engine is initialized
2. At least one SD model is available

## Error Handling

- `501 Not Implemented`: Image service not available (no engine or models)
- `400 Bad Request`: Invalid request (missing prompt, invalid JSON)
- `500 Internal Server Error`: Generation failed

## Performance

- **Generation Time**: Depends on:
  - Image size (1024x1024 is standard)
  - Quality setting (hd = 30 steps, standard = 20 steps)
  - Hardware (GPU vs CPU)
  - Number of images (n parameter)

- **Typical Times**:
  - Standard quality: 10-30 seconds (GPU)
  - HD quality: 20-60 seconds (GPU)
  - CPU: 5-10x slower

## Limitations

- Model selection not fully implemented (uses first available model)
- No streaming support
- No image editing/variations (only generation)
- No DALL-E specific features (style, user parameters ignored)

## Future Enhancements

- [ ] Model selection by name
- [ ] Image editing (inpainting/outpainting)
- [ ] Image variations
- [ ] Progress streaming
- [ ] Batch optimization
- [ ] Model caching
- [ ] Rate limiting per user
- [ ] Storage cleanup (old files)
