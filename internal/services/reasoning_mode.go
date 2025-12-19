package services

import (
	"fmt"
	"log"
	"strings"

	"github.com/kawai-network/veridium/fantasy/llamalib"
	"github.com/kawai-network/veridium/pkg/hardware"
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
// These models support <think> tags and native reasoning capabilities.
// Model patterns are dynamically loaded from llamalib.GetReasoningModelPatterns()
func IsReasoningModel(modelName string) bool {
	nameLower := strings.ToLower(modelName)

	// Get reasoning model patterns from model specs (no hardcoding)
	reasoningModels := llamalib.GetReasoningModelPatterns()

	for _, pattern := range reasoningModels {
		if strings.Contains(nameLower, pattern) {
			return true
		}
	}

	return false
}

// IsNonReasoningModel checks if a model name indicates a non-reasoning model
// These models do NOT support <think> tags and provide direct responses.
// Model patterns are dynamically loaded from llamalib.GetNonReasoningModelPatterns()
func IsNonReasoningModel(modelName string) bool {
	nameLower := strings.ToLower(modelName)

	// Get non-reasoning model patterns from model specs (no hardcoding)
	nonReasoningModels := llamalib.GetNonReasoningModelPatterns()

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

// HardwareRequirements defines minimum hardware specs for reasoning modes
type HardwareRequirements struct {
	MinRAM       int64  // Minimum RAM in GB
	MinCPUCores  int    // Minimum CPU cores
	RecommendGPU bool   // Whether GPU is recommended
	Description  string // Human-readable description
}

// GetHardwareRequirements returns hardware requirements for the reasoning mode
func (rc ReasoningConfig) GetHardwareRequirements() HardwareRequirements {
	switch rc.Mode {
	case ReasoningDisabled:
		// Non-reasoning models (Llama 3.2 3B) are lightweight
		return HardwareRequirements{
			MinRAM:       4,
			MinCPUCores:  2,
			RecommendGPU: false,
			Description:  "Lightweight - runs on most systems",
		}
	case ReasoningEnabled:
		// Reasoning models with /no_think need moderate resources
		return HardwareRequirements{
			MinRAM:       8,
			MinCPUCores:  4,
			RecommendGPU: true,
			Description:  "Moderate - 8GB RAM, 4+ cores recommended",
		}
	case ReasoningVerbose:
		// Full reasoning mode is resource-intensive
		return HardwareRequirements{
			MinRAM:       16,
			MinCPUCores:  6,
			RecommendGPU: true,
			Description:  "High-end - 16GB+ RAM, 6+ cores, GPU strongly recommended",
		}
	default:
		return HardwareRequirements{
			MinRAM:      4,
			MinCPUCores: 2,
		}
	}
}

// ValidateHardware checks if the system hardware is sufficient for the reasoning mode
// Returns true if hardware is sufficient, false otherwise with a reason
func (rc ReasoningConfig) ValidateHardware(specs *hardware.HardwareSpecs) (bool, string) {
	if specs == nil {
		// If we can't detect hardware, log warning but allow
		log.Printf("⚠️  Unable to detect hardware specs, proceeding with caution")
		return true, ""
	}

	requirements := rc.GetHardwareRequirements()

	// Check RAM
	if specs.AvailableRAM < requirements.MinRAM {
		return false, fmt.Sprintf(
			"Insufficient RAM: %dGB available, but %dGB required for %s mode. Consider using %s mode instead.",
			specs.AvailableRAM,
			requirements.MinRAM,
			rc.Mode,
			ReasoningDisabled,
		)
	}

	// Check CPU cores
	if specs.CPUCores < requirements.MinCPUCores {
		return false, fmt.Sprintf(
			"Insufficient CPU cores: %d cores available, but %d cores required for %s mode. Consider using %s mode instead.",
			specs.CPUCores,
			requirements.MinCPUCores,
			rc.Mode,
			ReasoningDisabled,
		)
	}

	// Warn if GPU is recommended but not available
	if requirements.RecommendGPU && specs.GPUMemory == 0 {
		log.Printf("⚠️  %s mode recommends GPU acceleration, but no GPU detected. Performance may be degraded.", rc.Mode)
		// Don't block - just warn
	}

	// Hardware is sufficient
	log.Printf("✅ Hardware validation passed for %s mode: RAM=%dGB (need %dGB), Cores=%d (need %d)",
		rc.Mode, specs.AvailableRAM, requirements.MinRAM, specs.CPUCores, requirements.MinCPUCores)
	return true, ""
}

// SuggestModeForHardware suggests the best reasoning mode for given hardware
func SuggestModeForHardware(specs *hardware.HardwareSpecs) ReasoningMode {
	if specs == nil {
		// Default to safest mode if we can't detect hardware
		return ReasoningDisabled
	}

	// Check from most demanding to least demanding
	verboseReq := ReasoningConfig{Mode: ReasoningVerbose}.GetHardwareRequirements()
	if specs.AvailableRAM >= verboseReq.MinRAM && specs.CPUCores >= verboseReq.MinCPUCores {
		return ReasoningVerbose
	}

	enabledReq := ReasoningConfig{Mode: ReasoningEnabled}.GetHardwareRequirements()
	if specs.AvailableRAM >= enabledReq.MinRAM && specs.CPUCores >= enabledReq.MinCPUCores {
		return ReasoningEnabled
	}

	// Default to disabled (non-reasoning) for lower-end hardware
	return ReasoningDisabled
}

// GetSummaryThreshold returns when to trigger auto-summarization (in turns)
// Returns 0 if mode doesn't support summarization
func (rc ReasoningConfig) GetSummaryThreshold() int {
	switch rc.Mode {
	case ReasoningDisabled:
		// Most efficient mode - can handle long conversations
		// Trigger summary after 10 turns (20 messages)
		// Token estimate: ~4,800 tokens (29% of 16K context) ✅
		return 10

	case ReasoningEnabled:
		// Balanced mode - moderate token usage
		// Trigger summary after 5 turns (10 messages)
		// Token estimate: ~2,400 tokens (15% of 16K context) ✅
		return 5

	case ReasoningVerbose:
		// High token usage - very short conversations
		// NO summary needed (3-5 turns max)
		return 0

	default:
		return 8
	}
}

// GetSummaryStrategy returns the summary strategy description
func (rc ReasoningConfig) GetSummaryStrategy() string {
	switch rc.Mode {
	case ReasoningDisabled:
		return "Auto-summarize after 10 turns (~4,800 tokens), keep last 20 messages"
	case ReasoningEnabled:
		return "Auto-summarize after 5 turns (~2,400 tokens), keep last 12 messages"
	case ReasoningVerbose:
		return "No summarization needed (short conversations only)"
	default:
		return "Auto-summarize after 8 turns (~3,840 tokens), keep last 16 messages"
	}
}

// GetIncrementalSummaryThreshold returns when to trigger incremental re-summarization
// This is used when a summary already exists and we want to update it with new messages
func (rc ReasoningConfig) GetIncrementalSummaryThreshold() int {
	switch rc.Mode {
	case ReasoningDisabled:
		// Re-summarize every 10 new turns (same as initial threshold)
		return 10

	case ReasoningEnabled:
		// Re-summarize every 5 new turns (same as initial threshold)
		return 5

	case ReasoningVerbose:
		// No incremental summary for verbose mode
		return 0

	default:
		return 8
	}
}

// ShouldSummarize checks if summarization should be triggered based on turn count
func (rc ReasoningConfig) ShouldSummarize(turnCount int) bool {
	threshold := rc.GetSummaryThreshold()
	if threshold == 0 {
		return false // Mode doesn't support summarization
	}
	return turnCount >= threshold
}
