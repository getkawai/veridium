package models

// GetAvailableModels returns all available Stable Diffusion models sorted by resource requirements
// Models are ordered from smallest to largest to ensure fallback logic works correctly
func GetAvailableModels() []ModelSpec {
	return []ModelSpec{
		{
			Name:            "qwen-image-2512-q4_k_m",
			URL:             "https://huggingface.co/unsloth/Qwen-Image-2512-GGUF/resolve/main/qwen-image-2512-Q4_K_M.gguf",
			Filename:        "qwen-image-2512-Q4_K_M.gguf",
			LLMURL:          "https://huggingface.co/unsloth/Qwen2.5-VL-7B-Instruct-GGUF/resolve/main/Qwen2.5-VL-7B-Instruct-UD-Q4_K_XL.gguf",
			LLMFilename:     "Qwen2.5-VL-7B-Instruct-UD-Q4_K_XL.gguf",
			VAEURL:          "https://huggingface.co/Comfy-Org/Qwen-Image_ComfyUI/resolve/main/split_files/vae/qwen_image_vae.safetensors",
			VAEFilename:     "qwen_image_vae.safetensors",
			EditModelURL:    "https://huggingface.co/unsloth/Qwen-Image-Edit-2511-GGUF/resolve/main/qwen-image-edit-2511-Q4_K_M.gguf",
			EditModelFile:   "qwen-image-edit-2511-Q4_K_M.gguf",
			Size:            13400,
			MinRAM:          14,
			RecommendedRAM:  24,
			MinVRAM:         8,
			RecommendedVRAM: 12,
			ModelType:       "Qwen-Image-2512",
			Description:     "Qwen-Image-2512 (Q4_K_M) bundle for stable-diffusion.cpp (diffusion + LLM + VAE)",
			Quantization:    "q4_k_m",
		},
		{
			Name:            "qwen-image-2512-q8_0",
			URL:             "https://huggingface.co/unsloth/Qwen-Image-2512-GGUF/resolve/main/qwen-image-2512-Q8_0.gguf",
			Filename:        "qwen-image-2512-Q8_0.gguf",
			LLMURL:          "https://huggingface.co/unsloth/Qwen2.5-VL-7B-Instruct-GGUF/resolve/main/Qwen2.5-VL-7B-Instruct-Q8_0.gguf",
			LLMFilename:     "Qwen2.5-VL-7B-Instruct-Q8_0.gguf",
			VAEURL:          "https://huggingface.co/Comfy-Org/Qwen-Image_ComfyUI/resolve/main/split_files/vae/qwen_image_vae.safetensors",
			VAEFilename:     "qwen_image_vae.safetensors",
			EditModelURL:    "https://huggingface.co/unsloth/Qwen-Image-Edit-2511-GGUF/resolve/main/qwen-image-edit-2511-Q8_0.gguf",
			EditModelFile:   "qwen-image-edit-2511-Q8_0.gguf",
			Size:            26200,
			MinRAM:          28,
			RecommendedRAM:  40,
			MinVRAM:         16,
			RecommendedVRAM: 24,
			ModelType:       "Qwen-Image-2512",
			Description:     "Qwen-Image-2512 (Q8_0) bundle for stable-diffusion.cpp (best quality, heavy resource usage)",
			Quantization:    "q8_0",
		},
	}
}
