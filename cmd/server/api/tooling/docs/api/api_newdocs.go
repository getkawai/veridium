// Package api provides additional API documentation for transcriptions and images.
package api

// =============================================================================

func transcriptionsDoc() apiDoc {
	return apiDoc{
		Name:        "Speech-to-Text API",
		Description: "Transcribe audio to text or translate audio to English. Compatible with the OpenAI Audio API.",
		Filename:    "DocsAPITranscriptions.tsx",
		Component:   "DocsAPITranscriptions",
		Groups: []endpointGroup{
			{
				Name:        "Transcriptions",
				Description: "Transcribe audio to text in the original language.",
				Endpoints: []endpoint{
					{
						Method:      "POST",
						Path:        "/audio/transcriptions",
						Description: "Transcribes audio into the input language. Supports multiple response formats including verbose JSON with timestamps.",
						Auth:        "Required when auth is enabled. Token must have 'audio-transcriptions' endpoint access.",
						Headers: []header{
							{Name: "Authorization", Description: "Bearer token for authentication", Required: true},
							{Name: "Content-Type", Description: "Must be multipart/form-data", Required: true},
						},
						RequestBody: &requestBody{
							ContentType: "multipart/form-data",
							Fields: []field{
								{Name: "file", Type: "binary", Required: true, Description: "Audio file (flac, mp3, mp4, mpeg, mpga, m4a, ogg, wav, webm)"},
								{Name: "model", Type: "string", Required: true, Description: "Transcription model (e.g., 'tiny', 'base', 'small', 'medium', 'large')"},
								{Name: "language", Type: "string", Required: false, Description: "Language code (ISO-639-1). Auto-detected if not provided."},
								{Name: "prompt", Type: "string", Required: false, Description: "Optional text to guide style or continue previous segment"},
								{Name: "response_format", Type: "string", Required: false, Description: "Format: json, text, srt, vtt, verbose_json (default: json)"},
								{Name: "temperature", Type: "number", Required: false, Description: "Sampling temperature 0-1 (default: 0)"},
							},
						},
						Response: &response{
							ContentType: "application/json or text",
							Description: "Returns transcription text. Verbose JSON includes segments, timestamps, and language detection.",
						},
						Examples: []example{
							{
								Description: "Basic transcription:",
								Code: `curl -X POST https://api.getkawai.com/v1/audio/transcriptions \
  -H "Authorization: Bearer $API_KEY" \
  -F "file=@audio.mp3" \
  -F "model=base" \
  -F "language=en"`,
							},
							{
								Description: "Verbose JSON with timestamps:",
								Code: `curl -X POST https://api.getkawai.com/v1/audio/transcriptions \
  -H "Authorization: Bearer $API_KEY" \
  -F "file=@audio.mp3" \
  -F "model=base" \
  -F "response_format=verbose_json"`,
							},
						},
					},
				},
			},
			{
				Name:        "Translations",
				Description: "Translate audio from any language to English.",
				Endpoints: []endpoint{
					{
						Method:      "POST",
						Path:        "/audio/translations",
						Description: "Translates audio into English. The source language is automatically detected.",
						Auth:        "Required when auth is enabled. Token must have 'audio-translations' endpoint access.",
						Headers: []header{
							{Name: "Authorization", Description: "Bearer token for authentication", Required: true},
							{Name: "Content-Type", Description: "Must be multipart/form-data", Required: true},
						},
						RequestBody: &requestBody{
							ContentType: "multipart/form-data",
							Fields: []field{
								{Name: "file", Type: "binary", Required: true, Description: "Audio file (flac, mp3, mp4, mpeg, mpga, m4a, ogg, wav, webm)"},
								{Name: "model", Type: "string", Required: true, Description: "Translation model (e.g., 'tiny', 'base', 'small', 'medium', 'large')"},
								{Name: "prompt", Type: "string", Required: false, Description: "Optional text to guide style"},
								{Name: "response_format", Type: "string", Required: false, Description: "Format: json, text, srt, vtt, verbose_json (default: json)"},
								{Name: "temperature", Type: "number", Required: false, Description: "Sampling temperature 0-1 (default: 0)"},
							},
						},
						Response: &response{
							ContentType: "application/json or text",
							Description: "Returns English translation text. Verbose JSON includes segments and timestamps.",
						},
						Examples: []example{
							{
								Description: "Translate Spanish audio to English:",
								Code: `curl -X POST https://api.getkawai.com/v1/audio/translations \
  -H "Authorization: Bearer $API_KEY" \
  -F "file=@spanish-audio.mp3" \
  -F "model=base"`,
							},
							{
								Description: "Translate with verbose output:",
								Code: `curl -X POST https://api.getkawai.com/v1/audio/translations \
  -H "Authorization: Bearer $API_KEY" \
  -F "file=@french-audio.mp3" \
  -F "model=base" \
  -F "response_format=verbose_json"`,
							},
						},
					},
				},
			},
			{
				Name:        "Response Formats",
				Description: "Available response formats for transcription and translation.",
				Endpoints: []endpoint{
					{
						Method:      "",
						Path:        "JSON (default)",
						Description: "Simple JSON response with only the transcribed/translated text.",
						Examples: []example{
							{
								Code: `{
  "text": "Hello, this is the transcribed text."
}`,
							},
						},
					},
					{
						Method:      "",
						Path:        "Verbose JSON",
						Description: "Detailed JSON response with language detection, duration, segments, and word-level timestamps.",
						Examples: []example{
							{
								Code: `{
  "task": "transcribe",
  "language": "en",
  "duration": 5.2,
  "text": "Hello, this is the transcribed text.",
  "segments": [
    {
      "id": 0,
      "start": 0.0,
      "end": 2.5,
      "text": "Hello, this is",
      "tokens": [123, 456, 789]
    }
  ],
  "words": [
    {"word": "Hello", "start": 0.0, "end": 0.5}
  ]
}`,
							},
						},
					},
				},
			},
			{
				Name:        "Supported Models",
				Description: "Whisper models available for transcription and translation.",
				Endpoints: []endpoint{
					{
						Method:      "",
						Path:        "Whisper Models",
						Description: "OpenAI Whisper models for speech recognition.",
						Examples: []example{
							{
								Description: "Available models:",
								Code: `tiny    - 39M parameters, fastest
base    - 74M parameters, good balance
small   - 244M parameters, better accuracy
medium  - 769M parameters, high accuracy
large   - 1550M parameters, best accuracy`,
							},
						},
					},
				},
			},
		},
	}
}

// =============================================================================

func imageDoc() apiDoc {
	return apiDoc{
		Name:        "Image Generation API",
		Description: "Generate, edit, and create variations of images. Compatible with the OpenAI Images API.",
		Filename:    "DocsAPIImage.tsx",
		Component:   "DocsAPIImage",
		Groups: []endpointGroup{
			{
				Name:        "Image Generation",
				Description: "Generate images from text prompts.",
				Endpoints: []endpoint{
					{
						Method:      "POST",
						Path:        "/images/generations",
						Description: "Generates images from a text prompt.",
						Auth:        "Required when auth is enabled. Token must have 'images-generations' endpoint access.",
						Headers: []header{
							{Name: "Authorization", Description: "Bearer token for authentication", Required: true},
							{Name: "Content-Type", Description: "Must be application/json", Required: true},
						},
						RequestBody: &requestBody{
							ContentType: "application/json",
							Fields: []field{
								{Name: "prompt", Type: "string", Required: true, Description: "Text description of the image to generate"},
								{Name: "model", Type: "string", Required: false, Description: "Model ID to use"},
								{Name: "n", Type: "integer", Required: false, Description: "Number of images to generate (1-10, default: 1)"},
								{Name: "size", Type: "string", Required: false, Description: "Image size (e.g., '1024x1024')"},
								{Name: "quality", Type: "string", Required: false, Description: "Image quality: 'standard' or 'hd'"},
								{Name: "response_format", Type: "string", Required: false, Description: "Response format: 'url' or 'b64_json'"},
							},
						},
						Response: &response{
							ContentType: "application/json",
							Description: "Returns generated image URLs or base64-encoded images.",
						},
						Examples: []example{
							{
								Description: "Generate a single image:",
								Code: `curl -X POST https://api.getkawai.com/v1/images/generations \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "A serene mountain landscape at sunset",
    "n": 1,
    "size": "1024x1024"
  }'`,
							},
						},
					},
				},
			},
			{
				Name:        "Image Editing",
				Description: "Edit existing images using masks and prompts.",
				Endpoints: []endpoint{
					{
						Method:      "POST",
						Path:        "/images/edits",
						Description: "Edits an image based on a text prompt and optional mask.",
						Auth:        "Required when auth is enabled. Token must have 'images-edits' endpoint access.",
						Headers: []header{
							{Name: "Authorization", Description: "Bearer token for authentication", Required: true},
							{Name: "Content-Type", Description: "Must be application/json", Required: true},
						},
						RequestBody: &requestBody{
							ContentType: "application/json",
							Fields: []field{
								{Name: "prompt", Type: "string", Required: true, Description: "Text description of the edit"},
								{Name: "image", Type: "string", Required: true, Description: "Base64-encoded image to edit"},
								{Name: "mask", Type: "string", Required: false, Description: "Base64-encoded mask (white=edit, black=keep)"},
								{Name: "n", Type: "integer", Required: false, Description: "Number of images (1-10, default: 1)"},
								{Name: "size", Type: "string", Required: false, Description: "Image size"},
								{Name: "response_format", Type: "string", Required: false, Description: "Response format: 'url' or 'b64_json'"},
							},
						},
						Response: &response{
							ContentType: "application/json",
							Description: "Returns edited image URLs or base64-encoded images.",
						},
						Examples: []example{
							{
								Description: "Edit image with mask:",
								Code: `curl -X POST https://api.getkawai.com/v1/images/edits \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "Add a beautiful sunset sky",
    "image": "data:image/png;base64,iVBORw0KG...",
    "mask": "data:image/png;base64,iVBORw0KG...",
    "n": 1
  }'`,
							},
						},
					},
				},
			},
			{
				Name:        "Image Variations",
				Description: "Create variations of existing images.",
				Endpoints: []endpoint{
					{
						Method:      "POST",
						Path:        "/images/variations",
						Description: "Creates variations of an input image.",
						Auth:        "Required when auth is enabled. Token must have 'images-variations' endpoint access.",
						Headers: []header{
							{Name: "Authorization", Description: "Bearer token for authentication", Required: true},
							{Name: "Content-Type", Description: "Must be application/json", Required: true},
						},
						RequestBody: &requestBody{
							ContentType: "application/json",
							Fields: []field{
								{Name: "image", Type: "string", Required: true, Description: "Base64-encoded image"},
								{Name: "n", Type: "integer", Required: false, Description: "Number of variations (1-10)"},
								{Name: "size", Type: "string", Required: false, Description: "Image size"},
								{Name: "response_format", Type: "string", Required: false, Description: "Response format: 'url' or 'b64_json'"},
							},
						},
						Response: &response{
							ContentType: "application/json",
							Description: "Returns variation image URLs or base64-encoded images.",
						},
						Examples: []example{
							{
								Description: "Generate 4 variations:",
								Code: `curl -X POST https://api.getkawai.com/v1/images/variations \
  -H "Authorization: Bearer $API_KEY" \
  -H "Content-Type: application/json" \
  -d '{
    "image": "data:image/png;base64,iVBORw0KG...",
    "n": 4
  }'`,
							},
						},
					},
				},
			},
		},
	}
}
