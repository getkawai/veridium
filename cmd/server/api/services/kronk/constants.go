package kronk

import "time"

// Model configuration constants
const (
	DefaultLLMOrg       = "unsloth"
	DefaultLLMRepo      = "Nemotron-3-Nano-30B-A3B-GGUF"
	DefaultLLMFile      = "Nemotron-3-Nano-30B-A3B-Q4_K_M.gguf"
	DefaultWhisperModel = "base"
)

// Download timeout constants
const (
	LibraryDownloadTimeout = 10 * time.Minute
	ModelDownloadTimeout   = 30 * time.Minute
	LLMDownloadTimeout     = 2 * time.Hour
)
