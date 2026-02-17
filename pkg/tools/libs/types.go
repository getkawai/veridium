package libs

type LibraryType string

const (
	LibraryLlama           LibraryType = "llama"
	LibraryWhisper         LibraryType = "whisper"
	LibraryStableDiffusion LibraryType = "stablediffusion"
	LibraryTTS             LibraryType = "tts"
)

func AllLibraryTypes() []LibraryType {
	return []LibraryType{LibraryLlama, LibraryWhisper, LibraryStableDiffusion, LibraryTTS}
}

func (lt LibraryType) String() string {
	return string(lt)
}

func ParseLibraryType(s string) LibraryType {
	switch s {
	case "llama":
		return LibraryLlama
	case "whisper":
		return LibraryWhisper
	case "stablediffusion", "sd":
		return LibraryStableDiffusion
	case "tts":
		return LibraryTTS
	default:
		return LibraryLlama
	}
}

func (lt LibraryType) DisplayName() string {
	switch lt {
	case LibraryLlama:
		return "llama.cpp"
	case LibraryWhisper:
		return "whisper.cpp"
	case LibraryStableDiffusion:
		return "stable-diffusion.cpp"
	case LibraryTTS:
		return "tts.cpp"
	default:
		return string(lt)
	}
}

func (lt LibraryType) DefaultVersion() string {
	switch lt {
	case LibraryLlama:
		return ""
	case LibraryWhisper:
		return "v1.8.3"
	case LibraryStableDiffusion:
		return "master-487-43e829f"
	case LibraryTTS:
		return ""
	default:
		return ""
	}
}

func (lt LibraryType) Subfolder() string {
	switch lt {
	case LibraryLlama:
		return "llama"
	case LibraryWhisper:
		return "whisper"
	case LibraryStableDiffusion:
		return "stablediffusion"
	case LibraryTTS:
		return "tts"
	default:
		return string(lt)
	}
}
