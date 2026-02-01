package whisper

import "unsafe"

// Token represents a whisper token
type Token int32

// Pos represents a position
type Pos int32

// SeqId represents a sequence ID
type SeqId int32

// SamplingStrategy represents the sampling strategy for decoding
type SamplingStrategy int32

const (
	SamplingGreedy SamplingStrategy = iota
	SamplingBeamSearch
)

// AlignmentHeadsPreset represents the alignment heads preset for DTW
type AlignmentHeadsPreset int32

const (
	AheadsNone AlignmentHeadsPreset = iota
	AheadsNTopMost
	AheadsCustom
	AheadsTinyEn
	AheadsTiny
	AheadsBaseEn
	AheadsBase
	AheadsSmallEn
	AheadsSmall
	AheadsMediumEn
	AheadsMedium
	AheadsLargeV1
	AheadsLargeV2
	AheadsLargeV3
	AheadsLargeV3Turbo
)

// Ahead represents a single alignment head
type Ahead struct {
	NTextLayer int32
	NHead      int32
}

// AheadList represents a list of alignment heads
type AheadList struct {
	NHeads int
	Heads  *Ahead
}

// ContextParams represents parameters for creating a whisper context
// Maps to struct whisper_context_params in whisper.h
type ContextParams struct {
	UseGPU             bool
	FlashAttn          bool
	GPUDevice          int32
	DTWTokenTimestamps bool
	DTWAheadsPreset    AlignmentHeadsPreset
	DTWNTop            int32
	DTWAheads          AheadList
	DTWMemSize         uintptr
}

// TokenData represents token-level data
// Maps to struct whisper_token_data in whisper.h
type TokenData struct {
	ID    Token
	TID   Token
	P     float32
	Plog  float32
	Pt    float32
	Ptsum float32
	T0    int64
	T1    int64
	TDtw  int64
	Vlen  float32
}

// VADParams represents Voice Activity Detection parameters
// Maps to struct whisper_vad_params in whisper.h
type VADParams struct {
	Threshold           float32
	MinSpeechDurationMs int32
	MinSilenceDurationMs int32
	MaxSpeechDurationS  float32
	SpeechPadMs         int32
	SamplesOverlap      float32
}

// FullParams represents parameters for the whisper_full() function
// Maps to struct whisper_full_params in whisper.h
type FullParams struct {
	Strategy           SamplingStrategy
	NThreads           int32
	NMaxTextCtx        int32
	OffsetMs           int32
	DurationMs         int32
	Translate          bool
	NoContext          bool
	NoTimestamps       bool
	SingleSegment      bool
	PrintSpecial       bool
	PrintProgress      bool
	PrintRealtime      bool
	PrintTimestamps    bool
	TokenTimestamps    bool
	TholdPt            float32
	TholdPtsum         float32
	MaxLen             int32
	SplitOnWord        bool
	MaxTokens          int32
	DebugMode          bool
	AudioCtx           int32
	TdrzEnable         bool
	SuppressRegex      *uint8
	InitialPrompt      *uint8
	CarryInitialPrompt bool
	PromptTokens       *Token
	PromptNTokens      int32
	Language           *uint8
	DetectLanguage     bool
	SuppressBlank      bool
	SuppressNst        bool
	Temperature        float32
	MaxInitialTs       float32
	LengthPenalty      float32
	TemperatureInc     float32
	EntropyThold       float32
	LogprobThold       float32
	NoSpeechThold      float32
	
	// Strategy-specific parameters
	Greedy struct {
		BestOf int32
	}
	
	BeamSearch struct {
		BeamSize int32
		Patience float32
	}
	
	// Callbacks (pointers - will be set via purego.NewCallback)
	NewSegmentCallback       uintptr
	NewSegmentCallbackUserData unsafe.Pointer
	ProgressCallback         uintptr
	ProgressCallbackUserData unsafe.Pointer
	EncoderBeginCallback     uintptr
	EncoderBeginCallbackUserData unsafe.Pointer
	AbortCallback            uintptr
	AbortCallbackUserData    unsafe.Pointer
	LogitsFilterCallback     uintptr
	LogitsFilterCallbackUserData unsafe.Pointer
	
	// Grammar
	GrammarRules   **GrammarElement
	NGrammarRules  uintptr
	IStartRule     uintptr
	GrammarPenalty float32
	
	// VAD
	VAD         bool
	VADModelPath *uint8
	VADParams   VADParams
}

// GrammarElement represents a grammar element
// Maps to struct whisper_grammar_element in whisper.h
type GrammarElement struct {
	Type  GrammarElementType
	Value uint32
}

// GrammarElementType represents the type of grammar element
type GrammarElementType int32

const (
	GrammarEnd GrammarElementType = iota
	GrammarAlt
	GrammarRuleRef
	GrammarChar
	GrammarCharNot
	GrammarCharRngUpper
	GrammarCharAlt
)

// Timings represents performance timing information
// Maps to struct whisper_timings in whisper.h
type Timings struct {
	SampleMs float32
	EncodeMs float32
	DecodeMs float32
	BatchdMs float32
	PromptMs float32
}

// LogLevel represents the log level for whisper
type LogLevel int32

const (
	LogLevelNone LogLevel = iota
	LogLevelDebug
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

// Constants from whisper.h
const (
	SampleRate = 16000
	NFFT       = 400
	HopLength  = 160
	ChunkSize  = 30
)
