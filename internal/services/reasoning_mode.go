package services

import (
	"fmt"
	"strings"
)

// ReasoningMode defines the reasoning behavior of the chat model
type ReasoningMode string

const (
	// ReasoningDisabled uses non-reasoning models (Llama 3.2, Mistral, etc.)
	// - Fastest responses
	// - Most efficient token usage
	// - No internal reasoning shown
	// - Best for long conversations (50-100 turns)
	ReasoningDisabled ReasoningMode = "disabled"

	// ReasoningEnabled uses reasoning models (Qwen3, GPT-OSS, etc.) with /no_think
	// - Moderate speed
	// - Balanced token usage
	// - Minimal reasoning overhead
	// - Good for medium conversations (30-50 turns)
	ReasoningEnabled ReasoningMode = "enabled"

	// ReasoningVerbose uses reasoning models with full thinking process
	// - Slower responses
	// - High token usage
	// - Shows full reasoning process
	// - Best for single-turn Q&A (3-5 turns)
	ReasoningVerbose ReasoningMode = "verbose"
)

// ReasoningConfig holds configuration for reasoning behavior
type ReasoningConfig struct {
	Mode                  ReasoningMode
	PreferredNonReasoning string // e.g., "llama", "mistral"
	PreferredReasoning    string // e.g., "qwen", "gpt-oss"
	StripThinkTags        bool   // Whether to strip <think> tags from output
}

// DefaultReasoningConfig returns the default reasoning configuration
func DefaultReasoningConfig() ReasoningConfig {
	return ReasoningConfig{
		Mode:                  ReasoningDisabled, // Default: non-reasoning for efficiency
		PreferredNonReasoning: "llama",           // Llama 3.2 is most efficient
		PreferredReasoning:    "qwen",            // Qwen3 has good reasoning
		StripThinkTags:        false,             // Disabled: rely on proper model selection instead
	}
}

// GetSystemPrompt returns the appropriate system prompt based on reasoning mode
func (rc ReasoningConfig) GetSystemPrompt(basePrompt string) string {
	switch rc.Mode {
	case ReasoningDisabled:
		// Non-reasoning models: simple, direct instruction
		return basePrompt + "\n\nBe concise and direct in your responses."

	case ReasoningEnabled:
		// Reasoning models with /no_think: minimize reasoning overhead
		return basePrompt + "\n\n/no_think\n\nIMPORTANT: Provide concise answers. Use minimal internal reasoning."

	case ReasoningVerbose:
		// Reasoning models with full thinking: encourage detailed reasoning
		return basePrompt + "\n\nThink through your answer step by step. Show your reasoning process using <think> tags."

	default:
		return basePrompt
	}
}

// ShouldStripThinkTags returns whether think tags should be stripped for this mode
func (rc ReasoningConfig) ShouldStripThinkTags() bool {
	// Disabled: Think tag stripping removed - proper model selection should prevent think tags
	// - ReasoningDisabled: Use non-reasoning models (Llama 3.2) - no think tags generated
	// - ReasoningEnabled: Use reasoning models with /no_think (Qwen3) - minimal think tags
	// - ReasoningVerbose: Use reasoning models (Qwen3) - full think tags shown
	return false // Always return false - no stripping
}

// GetPreferredModelPattern returns the model name pattern to prefer
func (rc ReasoningConfig) GetPreferredModelPattern() string {
	switch rc.Mode {
	case ReasoningDisabled:
		return rc.PreferredNonReasoning
	case ReasoningEnabled, ReasoningVerbose:
		return rc.PreferredReasoning
	default:
		return rc.PreferredNonReasoning
	}
}

// IsReasoningModel checks if a model name indicates a reasoning model
func IsReasoningModel(modelName string) bool {
	nameLower := strings.ToLower(modelName)

	// Known reasoning models
	reasoningModels := []string{
		"qwen",     // Qwen3 series
		"gpt-oss",  // OpenAI GPT-OSS series
		"deepseek", // DeepSeek R1 series
		"o1",       // OpenAI O1 series
	}

	for _, pattern := range reasoningModels {
		if strings.Contains(nameLower, pattern) {
			return true
		}
	}

	return false
}

// IsNonReasoningModel checks if a model name indicates a non-reasoning model
func IsNonReasoningModel(modelName string) bool {
	nameLower := strings.ToLower(modelName)

	// Known non-reasoning models
	nonReasoningModels := []string{
		"llama",   // Llama 3.x series
		"mistral", // Mistral series
		"gemma",   // Google Gemma series
		"phi",     // Microsoft Phi series
	}

	for _, pattern := range nonReasoningModels {
		if strings.Contains(nameLower, pattern) {
			return true
		}
	}

	return false
}

// ValidateModelForMode checks if the loaded model is appropriate for the reasoning mode
func (rc ReasoningConfig) ValidateModelForMode(modelName string) error {
	isReasoning := IsReasoningModel(modelName)
	isNonReasoning := IsNonReasoningModel(modelName)

	switch rc.Mode {
	case ReasoningDisabled:
		if isReasoning && !isNonReasoning {
			return fmt.Errorf("reasoning mode is disabled but loaded model (%s) is a reasoning model. Consider using Llama 3.2 or Mistral", modelName)
		}
	case ReasoningEnabled, ReasoningVerbose:
		if isNonReasoning && !isReasoning {
			return fmt.Errorf("reasoning mode is enabled but loaded model (%s) is a non-reasoning model. Consider using Qwen3 or GPT-OSS", modelName)
		}
	}

	return nil
}

// GetRecommendedModel returns the recommended model for the current mode
func (rc ReasoningConfig) GetRecommendedModel() string {
	switch rc.Mode {
	case ReasoningDisabled:
		// Prefer Llama 3.2 for non-reasoning
		return "Llama-3.2-3B-Instruct-Q4_K_M.gguf"
	case ReasoningEnabled, ReasoningVerbose:
		// Prefer Qwen3 for reasoning
		return "Qwen_Qwen3-1.7B-Q4_K_M.gguf"
	default:
		return "Llama-3.2-3B-Instruct-Q4_K_M.gguf"
	}
}

// GetModeDescription returns a user-friendly description of the mode
func (rc ReasoningConfig) GetModeDescription() string {
	switch rc.Mode {
	case ReasoningDisabled:
		return "Fast & Efficient (No Reasoning) - Best for long conversations"
	case ReasoningEnabled:
		return "Balanced (Minimal Reasoning) - Good for most use cases"
	case ReasoningVerbose:
		return "Detailed (Full Reasoning) - Best for complex questions"
	default:
		return "Unknown mode"
	}
}

// GetExpectedPerformance returns performance expectations for the mode
func (rc ReasoningConfig) GetExpectedPerformance() map[string]string {
	switch rc.Mode {
	case ReasoningDisabled:
		return map[string]string{
			"speed":            "1.2s per turn",
			"token_efficiency": "~48 tokens/turn",
			"max_turns":        "50-100 turns",
			"response_size":    "~165 chars",
		}
	case ReasoningEnabled:
		return map[string]string{
			"speed":            "1.8s per turn",
			"token_efficiency": "~118 tokens/turn",
			"max_turns":        "30-50 turns",
			"response_size":    "~450 chars",
		}
	case ReasoningVerbose:
		return map[string]string{
			"speed":            "11s per turn",
			"token_efficiency": "~1000 tokens/turn",
			"max_turns":        "3-5 turns",
			"response_size":    "~3800 chars",
		}
	default:
		return map[string]string{}
	}
}
