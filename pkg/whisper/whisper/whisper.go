package whisper

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"unsafe"
)

// Dynamic library handle
var libWhisper uintptr

// Function declarations - will be bound to C functions
var (
	// Version
	Version func() *uint8

	// Context creation
	InitFromFileWithParams   func(path *uint8, params ContextParams) unsafe.Pointer
	InitFromBufferWithParams func(buffer unsafe.Pointer, bufferSize uintptr, params ContextParams) unsafe.Pointer
	ContextDefaultParams     func() ContextParams

	// Cleanup
	Free               func(ctx unsafe.Pointer)
	FreeParams         func(params *FullParams)
	FreeContextParams  func(params *ContextParams)

	// State
	InitState func(ctx unsafe.Pointer) unsafe.Pointer
	FreeState func(state unsafe.Pointer)

	// Audio processing
	PCMToMel           func(ctx unsafe.Pointer, samples *float32, nSamples int32, nThreads int32) int32
	PCMToMelWithState  func(ctx unsafe.Pointer, state unsafe.Pointer, samples *float32, nSamples int32, nThreads int32) int32
	SetMel             func(ctx unsafe.Pointer, data *float32, nLen int32, nMel int32) int32
	SetMelWithState    func(ctx unsafe.Pointer, state unsafe.Pointer, data *float32, nLen int32, nMel int32) int32

	// Encoding/Decoding
	Encode           func(ctx unsafe.Pointer, offset int32, nThreads int32) int32
	EncodeWithState  func(ctx unsafe.Pointer, state unsafe.Pointer, offset int32, nThreads int32) int32
	Decode           func(ctx unsafe.Pointer, tokens *Token, nTokens int32, nPast int32, nThreads int32) int32
	DecodeWithState  func(ctx unsafe.Pointer, state unsafe.Pointer, tokens *Token, nTokens int32, nPast int32, nThreads int32) int32

	// Full transcription
	Full             func(ctx unsafe.Pointer, params FullParams, samples *float32, nSamples int32) int32
	FullWithState    func(ctx unsafe.Pointer, state unsafe.Pointer, params FullParams, samples *float32, nSamples int32) int32
	FullParallel     func(ctx unsafe.Pointer, params FullParams, samples *float32, nSamples int32, nProcessors int32) int32
	FullDefaultParams func(strategy int32) FullParams

	// Results
	FullNSegments            func(ctx unsafe.Pointer) int32
	FullNSegmentsFromState   func(state unsafe.Pointer) int32
	FullLangID               func(ctx unsafe.Pointer) int32
	FullLangIDFromState      func(state unsafe.Pointer) int32
	FullGetSegmentT0         func(ctx unsafe.Pointer, iSegment int32) int64
	FullGetSegmentT0FromState func(state unsafe.Pointer, iSegment int32) int64
	FullGetSegmentT1         func(ctx unsafe.Pointer, iSegment int32) int64
	FullGetSegmentT1FromState func(state unsafe.Pointer, iSegment int32) int64
	FullGetSegmentText       func(ctx unsafe.Pointer, iSegment int32) *uint8
	FullGetSegmentTextFromState func(state unsafe.Pointer, iSegment int32) *uint8
	FullGetSegmentSpeakerTurnNext func(ctx unsafe.Pointer, iSegment int32) bool
	FullGetSegmentSpeakerTurnNextFromState func(state unsafe.Pointer, iSegment int32) bool

	// Tokenization
	Tokenize         func(ctx unsafe.Pointer, text *uint8, tokens *Token, nMaxTokens int32) int32
	TokenCount       func(ctx unsafe.Pointer, text *uint8) int32
	TokenToStr       func(ctx unsafe.Pointer, token Token) *uint8

	// Language
	LangMaxID    func() int32
	LangID       func(lang *uint8) int32
	LangStr      func(id int32) *uint8
	LangStrFull  func(id int32) *uint8
	LangAutoDetect func(ctx unsafe.Pointer, offsetMs int32, nThreads int32, langProbs *float32) int32
	LangAutoDetectWithState func(ctx unsafe.Pointer, state unsafe.Pointer, offsetMs int32, nThreads int32, langProbs *float32) int32

	// Model info
	NLen           func(ctx unsafe.Pointer) int32
	NLenFromState  func(state unsafe.Pointer) int32
	NVocab         func(ctx unsafe.Pointer) int32
	NTextCtx       func(ctx unsafe.Pointer) int32
	NAudioCtx      func(ctx unsafe.Pointer) int32
	IsMultilingual func(ctx unsafe.Pointer) int32
	ModelTypeReadable func(ctx unsafe.Pointer) *uint8

	// Special tokens
	TokenEOT  func(ctx unsafe.Pointer) Token
	TokenSOT  func(ctx unsafe.Pointer) Token
	TokenSOLM func(ctx unsafe.Pointer) Token
	TokenPREV func(ctx unsafe.Pointer) Token
	TokenNOSP func(ctx unsafe.Pointer) Token
	TokenNOT  func(ctx unsafe.Pointer) Token
	TokenBEG  func(ctx unsafe.Pointer) Token
	TokenLang func(ctx unsafe.Pointer, langID int32) Token
	TokenTranslate  func(ctx unsafe.Pointer) Token
	TokenTranscribe func(ctx unsafe.Pointer) Token

	// Timings
	GetTimings  func(ctx unsafe.Pointer) *Timings
	PrintTimings func(ctx unsafe.Pointer)
	ResetTimings func(ctx unsafe.Pointer)

	// System info
	PrintSystemInfo func() *uint8

	// Log set
	LogSet func(callback uintptr, userData unsafe.Pointer)
)

// CString converts a Go string to a C string (null-terminated)
func CString(s string) *uint8 {
	b := append([]byte(s), 0)
	return &b[0]
}

// GoString converts a C string to a Go string
func GoString(p *uint8) string {
	if p == nil {
		return ""
	}
	var b []byte
	for *p != 0 {
		b = append(b, byte(*p))
		p = (*uint8)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + 1))
	}
	return string(b)
}

// SetLibraryDirectory sets the directory where the whisper library is located
var libraryDirectory string

func SetLibraryDirectory(dir string) {
	libraryDirectory = dir
}

// Load loads the whisper library from the specified directory
func Load(dir string) error {
	if dir == "" {
		dir = libraryDirectory
	}
	if dir == "" {
		return fmt.Errorf("library directory not specified")
	}

	var libName string
	switch runtime.GOOS {
	case "darwin":
		libName = "libwhisper.dylib"
	case "windows":
		libName = "whisper.dll"
	case "linux":
		libName = "libwhisper.so"
	default:
		return fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	libPath := filepath.Join(dir, libName)
	
	// Check if library exists
	if _, err := os.Stat(libPath); err != nil {
		return fmt.Errorf("library not found at %s: %w", libPath, err)
	}

	var err error
	libWhisper, err = openLibrary(libPath)
	if err != nil {
		return fmt.Errorf("failed to load library: %w", err)
	}

	// Bind functions
	bindFunctions()

	return nil
}

// bindFunctions binds all whisper functions from the loaded library
func bindFunctions() {
	// Version
	registerFunc(&Version, "whisper_version")

	// Context creation
	registerFunc(&InitFromFileWithParams, "whisper_init_from_file_with_params")
	registerFunc(&InitFromBufferWithParams, "whisper_init_from_buffer_with_params")
	registerFunc(&ContextDefaultParams, "whisper_context_default_params")

	// Cleanup
	registerFunc(&Free, "whisper_free")
	registerFunc(&FreeParams, "whisper_free_params")
	registerFunc(&FreeContextParams, "whisper_free_context_params")

	// State
	registerFunc(&InitState, "whisper_init_state")
	registerFunc(&FreeState, "whisper_free_state")

	// Audio processing
	registerFunc(&PCMToMel, "whisper_pcm_to_mel")
	registerFunc(&PCMToMelWithState, "whisper_pcm_to_mel_with_state")
	registerFunc(&SetMel, "whisper_set_mel")
	registerFunc(&SetMelWithState, "whisper_set_mel_with_state")

	// Encoding/Decoding
	registerFunc(&Encode, "whisper_encode")
	registerFunc(&EncodeWithState, "whisper_encode_with_state")
	registerFunc(&Decode, "whisper_decode")
	registerFunc(&DecodeWithState, "whisper_decode_with_state")

	// Full transcription
	registerFunc(&Full, "whisper_full")
	registerFunc(&FullWithState, "whisper_full_with_state")
	registerFunc(&FullParallel, "whisper_full_parallel")
	registerFunc(&FullDefaultParams, "whisper_full_default_params")

	// Results
	registerFunc(&FullNSegments, "whisper_full_n_segments")
	registerFunc(&FullNSegmentsFromState, "whisper_full_n_segments_from_state")
	registerFunc(&FullLangID, "whisper_full_lang_id")
	registerFunc(&FullLangIDFromState, "whisper_full_lang_id_from_state")
	registerFunc(&FullGetSegmentT0, "whisper_full_get_segment_t0")
	registerFunc(&FullGetSegmentT0FromState, "whisper_full_get_segment_t0_from_state")
	registerFunc(&FullGetSegmentT1, "whisper_full_get_segment_t1")
	registerFunc(&FullGetSegmentT1FromState, "whisper_full_get_segment_t1_from_state")
	registerFunc(&FullGetSegmentText, "whisper_full_get_segment_text")
	registerFunc(&FullGetSegmentTextFromState, "whisper_full_get_segment_text_from_state")
	registerFunc(&FullGetSegmentSpeakerTurnNext, "whisper_full_get_segment_speaker_turn_next")
	registerFunc(&FullGetSegmentSpeakerTurnNextFromState, "whisper_full_get_segment_speaker_turn_next_from_state")

	// Tokenization
	registerFunc(&Tokenize, "whisper_tokenize")
	registerFunc(&TokenCount, "whisper_token_count")
	registerFunc(&TokenToStr, "whisper_token_to_str")

	// Language
	registerFunc(&LangMaxID, "whisper_lang_max_id")
	registerFunc(&LangID, "whisper_lang_id")
	registerFunc(&LangStr, "whisper_lang_str")
	registerFunc(&LangStrFull, "whisper_lang_str_full")
	registerFunc(&LangAutoDetect, "whisper_lang_auto_detect")
	registerFunc(&LangAutoDetectWithState, "whisper_lang_auto_detect_with_state")

	// Model info
	registerFunc(&NLen, "whisper_n_len")
	registerFunc(&NLenFromState, "whisper_n_len_from_state")
	registerFunc(&NVocab, "whisper_n_vocab")
	registerFunc(&NTextCtx, "whisper_n_text_ctx")
	registerFunc(&NAudioCtx, "whisper_n_audio_ctx")
	registerFunc(&IsMultilingual, "whisper_is_multilingual")
	registerFunc(&ModelTypeReadable, "whisper_model_type_readable")

	// Special tokens
	registerFunc(&TokenEOT, "whisper_token_eot")
	registerFunc(&TokenSOT, "whisper_token_sot")
	registerFunc(&TokenSOLM, "whisper_token_solm")
	registerFunc(&TokenPREV, "whisper_token_prev")
	registerFunc(&TokenNOSP, "whisper_token_nosp")
	registerFunc(&TokenNOT, "whisper_token_not")
	registerFunc(&TokenBEG, "whisper_token_beg")
	registerFunc(&TokenLang, "whisper_token_lang")
	registerFunc(&TokenTranslate, "whisper_token_translate")
	registerFunc(&TokenTranscribe, "whisper_token_transcribe")

	// Timings
	registerFunc(&GetTimings, "whisper_get_timings")
	registerFunc(&PrintTimings, "whisper_print_timings")
	registerFunc(&ResetTimings, "whisper_reset_timings")

	// System info
	registerFunc(&PrintSystemInfo, "whisper_print_system_info")

	// Log set
	registerFunc(&LogSet, "whisper_log_set")
}

// registerFunc is a helper to register a library function
// This is a platform-specific implementation
func registerFunc(fn interface{}, name string) {
	registerLibFunc(fn, libWhisper, name)
}
