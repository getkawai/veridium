package memory

// MetaData represents session metadata (matches frontend MetaData type)
type MetaData struct {
	Avatar          string   `json:"avatar,omitempty"`
	BackgroundColor string   `json:"backgroundColor,omitempty"`
	Description     string   `json:"description,omitempty"`
	Tags            []string `json:"tags,omitempty"`
	Title           string   `json:"title,omitempty"`
}

// WorkingModel represents a model configuration
type WorkingModel struct {
	Model    string `json:"model"`
	Provider string `json:"provider"`
}

// AgentChatConfig represents chat configuration (matches frontend LobeAgentChatConfig)
type AgentChatConfig struct {
	DisplayMode              string        `json:"displayMode,omitempty"`           // "chat" or "docs"
	EnableAutoCreateTopic    bool          `json:"enableAutoCreateTopic,omitempty"` // Auto create topic
	AutoCreateTopicThreshold int           `json:"autoCreateTopicThreshold"`        // Threshold for auto topic creation
	EnableMaxTokens          bool          `json:"enableMaxTokens,omitempty"`       // Enable max tokens
	EnableStreaming          bool          `json:"enableStreaming,omitempty"`       // Enable streaming
	EnableReasoning          bool          `json:"enableReasoning,omitempty"`       // Enable reasoning
	EnableReasoningEffort    bool          `json:"enableReasoningEffort,omitempty"` // Enable custom reasoning effort
	ReasoningBudgetToken     int           `json:"reasoningBudgetToken,omitempty"`  // Reasoning token budget
	ReasoningEffort          string        `json:"reasoningEffort,omitempty"`       // "low", "medium", "high"
	GPT5ReasoningEffort      string        `json:"gpt5ReasoningEffort,omitempty"`   // "minimal", "low", "medium", "high"
	TextVerbosity            string        `json:"textVerbosity,omitempty"`         // "low", "medium", "high"
	Thinking                 string        `json:"thinking,omitempty"`              // "disabled", "auto", "enabled"
	ThinkingBudget           int           `json:"thinkingBudget,omitempty"`        // Thinking budget
	DisableContextCaching    bool          `json:"disableContextCaching,omitempty"` // Disable context caching
	HistoryCount             int           `json:"historyCount,omitempty"`          // History message count
	EnableHistoryCount       bool          `json:"enableHistoryCount,omitempty"`    // Enable history count
	EnableCompressHistory    bool          `json:"enableCompressHistory,omitempty"` // Enable compress history
	InputTemplate            string        `json:"inputTemplate,omitempty"`         // Input template
	SearchMode               string        `json:"searchMode,omitempty"`            // "off", "on", "auto"
	SearchFCModel            *WorkingModel `json:"searchFCModel,omitempty"`         // Search function calling model
	URLContext               bool          `json:"urlContext,omitempty"`            // URL context
	UseModelBuiltinSearch    bool          `json:"useModelBuiltinSearch,omitempty"` // Use model builtin search
}

// TTSConfig represents TTS configuration (matches frontend LobeAgentTTSConfig)
type TTSConfig struct {
	ShowAllLocaleVoice bool              `json:"showAllLocaleVoice,omitempty"` // Show all locale voices
	STTLocale          string            `json:"sttLocale"`                    // STT locale ("auto" or language code)
	TTSService         string            `json:"ttsService"`                   // "browser", "openai", "edge", "microsoft"
	Voice              map[string]string `json:"voice"`                        // Voice settings per service
}

// LLMParams represents language model parameters
type LLMParams struct {
	FrequencyPenalty float64 `json:"frequency_penalty,omitempty"` // Frequency penalty
	MaxTokens        int     `json:"max_tokens,omitempty"`        // Max tokens
	PresencePenalty  float64 `json:"presence_penalty,omitempty"`  // Presence penalty
	Temperature      float64 `json:"temperature,omitempty"`       // Temperature
	TopP             float64 `json:"top_p,omitempty"`             // Top P
}

// FileItem represents a file reference
type FileItem struct {
	ID   string `json:"id"`
	Name string `json:"name,omitempty"`
	URL  string `json:"url,omitempty"`
}

// KnowledgeBaseItem represents a knowledge base reference
type KnowledgeBaseItem struct {
	ID      string `json:"id"`
	Name    string `json:"name,omitempty"`
	Enabled bool   `json:"enabled,omitempty"`
}

// FewShots represents few-shot examples
type FewShots []map[string]interface{}

// AgentConfig represents complete agent configuration (matches frontend LobeAgentConfig)
type AgentConfig struct {
	ChatConfig       AgentChatConfig     `json:"chatConfig"`                 // Chat configuration
	FewShots         FewShots            `json:"fewShots,omitempty"`         // Few-shot examples
	Files            []FileItem          `json:"files,omitempty"`            // Attached files
	ID               string              `json:"id,omitempty"`               // Config ID
	KnowledgeBases   []KnowledgeBaseItem `json:"knowledgeBases,omitempty"`   // Knowledge bases
	Model            string              `json:"model"`                      // Model name
	OpeningMessage   string              `json:"openingMessage,omitempty"`   // Opening message
	OpeningQuestions []string            `json:"openingQuestions,omitempty"` // Opening questions
	Params           LLMParams           `json:"params"`                     // LLM parameters
	Plugins          []string            `json:"plugins,omitempty"`          // Enabled plugins
	Provider         string              `json:"provider,omitempty"`         // Model provider
	SystemRole       string              `json:"systemRole"`                 // System role/prompt
	TTS              TTSConfig           `json:"tts"`                        // TTS configuration
}

// DefaultMetaData returns default metadata
func DefaultMetaData() MetaData {
	return MetaData{
		Avatar:          "🤖",
		BackgroundColor: "rgba(0,0,0,0)",
		Description:     "",
		Tags:            []string{},
	}
}

// DefaultAgentChatConfig returns default chat configuration
func DefaultAgentChatConfig() AgentChatConfig {
	return AgentChatConfig{
		AutoCreateTopicThreshold: 2,
		DisplayMode:              "chat",
		EnableAutoCreateTopic:    true,
		EnableCompressHistory:    true,
		EnableHistoryCount:       true,
		EnableReasoning:          false,
		EnableStreaming:          true,
		HistoryCount:             20,
		ReasoningBudgetToken:     1024,
		SearchMode:               "off",
		SearchFCModel: &WorkingModel{
			Model:    "gpt-4o-mini",
			Provider: "openai",
		},
	}
}

// DefaultTTSConfig returns default TTS configuration
func DefaultTTSConfig() TTSConfig {
	return TTSConfig{
		ShowAllLocaleVoice: false,
		STTLocale:          "auto",
		TTSService:         "openai",
		Voice: map[string]string{
			"openai": "alloy",
		},
	}
}

// DefaultLLMParams returns default LLM parameters
func DefaultLLMParams() LLMParams {
	return LLMParams{
		FrequencyPenalty: 0,
		PresencePenalty:  0,
		Temperature:      1,
		TopP:             1,
	}
}

// DefaultAgentConfig returns default agent configuration
func DefaultAgentConfig() AgentConfig {
	return AgentConfig{
		ChatConfig:       DefaultAgentChatConfig(),
		Model:            "gpt-4o-mini",
		OpeningQuestions: []string{},
		Params:           DefaultLLMParams(),
		Plugins:          []string{},
		Provider:         "openai",
		SystemRole:       "",
		TTS:              DefaultTTSConfig(),
	}
}
