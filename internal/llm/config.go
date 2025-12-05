/*
 * Copyright 2025 Veridium Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package llm

import (
	"log"

	"github.com/kawai-network/veridium/pkg/yzma/tools"
	"github.com/kawai-network/veridium/types"
)

// ============================================================================
// Developer-configured provider settings (hardcoded for development)
// TODO: Move to config file or environment variables for production
// ============================================================================

// TaskProviderConfig holds the provider configuration for a specific task
type TaskProviderConfig struct {
	ProviderType types.ProviderType
	APIKey       string
	Model        string
	MaxTokens    int
	Options      map[string]any
}

// DevConfig holds all hardcoded development configurations
// This is where developers configure which provider to use for each task
type DevConfig struct {
	// Chat task - main conversation (needs streaming, tool calling)
	Chat TaskProviderConfig

	// Title generation - lightweight, fast
	Title TaskProviderConfig

	// Summary generation - background task, can be slower
	Summary TaskProviderConfig

	// Image description - needs VL (Vision-Language) capability
	ImageDescribe TaskProviderConfig

	// UseLocalFallback - if true, use local llama as fallback for all tasks
	UseLocalFallback bool
}

// GetDefaultDevConfig returns the default development configuration
// Developers can modify this function to change provider assignments
func GetDefaultDevConfig() DevConfig {
	return DevConfig{
		// Chat: Use OpenRouter with free model for development
		Chat: TaskProviderConfig{
			ProviderType: types.ProviderOpenRouter,
			APIKey:       "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",
			Model:        "amazon/nova-2-lite-v1:free",
			MaxTokens:    4096,
			Options: map[string]any{
				"app_name": "Veridium",
			},
		},

		// Title: Use Zhipu GLM for fast, cheap title generation
		// GLM-4.6 is the latest model with good performance
		Title: TaskProviderConfig{
			ProviderType: types.ProviderZhipuAI,
			APIKey:       "a10854167085448cac33753523919ac9.D41CLq6KxXTY7g4u",
			Model:        "glm-4.6",
			MaxTokens:    256,
		},

		// Summary: Use local Llama for background summarization
		// No network latency, runs efficiently in background
		Summary: TaskProviderConfig{
			ProviderType: types.ProviderLlama,
			Model:        "", // Will use auto-detected model
			MaxTokens:    512,
		},

		// ImageDescribe: Use Zhipu GLM-4V for vision tasks
		ImageDescribe: TaskProviderConfig{
			ProviderType: types.ProviderZhipuAI,
			APIKey:       "a10854167085448cac33753523919ac9.D41CLq6KxXTY7g4u",
			Model:        "glm-4v-flash",
			MaxTokens:    1024,
		},

		// Use local Llama as fallback when remote providers fail
		UseLocalFallback: true,
	}
}

// BuildTaskRouter creates a TaskRouter from DevConfig
// This is the main entry point for setting up multi-provider routing
func BuildTaskRouter(
	config DevConfig,
	toolRegistry *tools.ToolRegistry,
	localProvider Provider, // Local llama provider as fallback
) *TaskRouter {
	factory := NewProviderFactory(toolRegistry)
	router := NewTaskRouter(toolRegistry, nil)

	// Set local provider as fallback if enabled
	if config.UseLocalFallback && localProvider != nil {
		router.SetFallback(localProvider)
		log.Printf("🔀 TaskRouter: Local Llama set as fallback provider")
	}

	// Configure Chat provider
	if config.Chat.APIKey != "" {
		chatProvider := createProviderFromConfig(factory, config.Chat)
		if chatProvider != nil {
			router.SetProvider(TaskChat, chatProvider)
			log.Printf("🔀 TaskRouter: Chat -> %s (%s)", config.Chat.ProviderType, config.Chat.Model)
		}
	} else if localProvider != nil {
		router.SetProvider(TaskChat, localProvider)
		log.Printf("🔀 TaskRouter: Chat -> Local Llama (no remote API key)")
	}

	// Configure Title provider
	if config.Title.APIKey != "" {
		titleProvider := createProviderFromConfig(factory, config.Title)
		if titleProvider != nil {
			router.SetProvider(TaskTitleGen, titleProvider)
			log.Printf("🔀 TaskRouter: Title -> %s (%s)", config.Title.ProviderType, config.Title.Model)
		}
	} else if localProvider != nil {
		// Use local provider for title if no API key
		router.SetProvider(TaskTitleGen, localProvider)
		log.Printf("🔀 TaskRouter: Title -> Local Llama (no remote API key)")
	}

	// Configure Summary provider
	if config.Summary.ProviderType == types.ProviderLlama && localProvider != nil {
		// Explicitly use local for summary
		router.SetProvider(TaskSummaryGen, localProvider)
		log.Printf("🔀 TaskRouter: Summary -> Local Llama (configured)")
	} else if config.Summary.APIKey != "" {
		summaryProvider := createProviderFromConfig(factory, config.Summary)
		if summaryProvider != nil {
			router.SetProvider(TaskSummaryGen, summaryProvider)
			log.Printf("🔀 TaskRouter: Summary -> %s (%s)", config.Summary.ProviderType, config.Summary.Model)
		}
	} else if localProvider != nil {
		router.SetProvider(TaskSummaryGen, localProvider)
		log.Printf("🔀 TaskRouter: Summary -> Local Llama (default)")
	}

	// Configure ImageDescribe provider
	if config.ImageDescribe.APIKey != "" {
		imageProvider := createProviderFromConfig(factory, config.ImageDescribe)
		if imageProvider != nil {
			router.SetProvider(TaskImageDescribe, imageProvider)
			log.Printf("🔀 TaskRouter: ImageDescribe -> %s (%s)", config.ImageDescribe.ProviderType, config.ImageDescribe.Model)
		}
	}
	// Note: No fallback for ImageDescribe - requires VL capability

	return router
}

// createProviderFromConfig creates a provider from TaskProviderConfig
func createProviderFromConfig(factory *ProviderFactory, config TaskProviderConfig) Provider {
	if config.ProviderType == types.ProviderLlama {
		// Local provider is handled separately
		return nil
	}

	providerConfig := types.ProviderConfig{
		Type:      config.ProviderType,
		APIKey:    config.APIKey,
		Model:     config.Model,
		MaxTokens: config.MaxTokens,
		Options:   config.Options,
	}

	// Set name based on type
	switch config.ProviderType {
	case types.ProviderOpenRouter:
		providerConfig.Name = "OpenRouter"
	case types.ProviderZhipuAI:
		providerConfig.Name = "Zhipu GLM"
	}

	provider, err := factory.CreateProvider(providerConfig)
	if err != nil {
		log.Printf("⚠️  Failed to create provider %s: %v", config.ProviderType, err)
		return nil
	}

	return provider
}

// ============================================================================
// Helper functions for updating config at runtime
// ============================================================================

// UpdateChatProvider updates the chat provider configuration
func (r *TaskRouter) UpdateChatProvider(config TaskProviderConfig) error {
	factory := NewProviderFactory(r.toolRegistry)
	provider := createProviderFromConfig(factory, config)
	if provider == nil {
		return &RouterError{Message: "failed to create chat provider"}
	}
	r.SetProvider(TaskChat, provider)
	return nil
}

// UpdateTitleProvider updates the title provider configuration
func (r *TaskRouter) UpdateTitleProvider(config TaskProviderConfig) error {
	factory := NewProviderFactory(r.toolRegistry)
	provider := createProviderFromConfig(factory, config)
	if provider == nil {
		return &RouterError{Message: "failed to create title provider"}
	}
	r.SetProvider(TaskTitleGen, provider)
	return nil
}

// UpdateSummaryProvider updates the summary provider configuration
func (r *TaskRouter) UpdateSummaryProvider(config TaskProviderConfig) error {
	factory := NewProviderFactory(r.toolRegistry)
	provider := createProviderFromConfig(factory, config)
	if provider == nil {
		return &RouterError{Message: "failed to create summary provider"}
	}
	r.SetProvider(TaskSummaryGen, provider)
	return nil
}
