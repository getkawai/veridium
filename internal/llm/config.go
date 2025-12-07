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

	"github.com/kawai-network/veridium/fantasy"
	"github.com/kawai-network/veridium/pkg/yzma/tools"
)

// ============================================================================
// LLM Provider Configuration Types
// ============================================================================

// ProviderType identifies the LLM provider type
type ProviderType string

const (
	ProviderLlama      ProviderType = "llama"      // Local llama.cpp (via yzma)
	ProviderOpenAI     ProviderType = "openai"     // OpenAI API
	ProviderOpenRouter ProviderType = "openrouter" // OpenRouter (OpenAI-compatible, multi-model)
	ProviderZhipuAI    ProviderType = "zhipu"      // Zhipu GLM (OpenAI-compatible, Chinese AI)
)

// ProviderConfig holds configuration for an LLM provider
type ProviderConfig struct {
	Type        ProviderType   `json:"type"`
	Name        string         `json:"name"`
	BaseURL     string         `json:"base_url"`
	APIKey      string         `json:"api_key"`
	Model       string         `json:"model"`
	MaxTokens   int            `json:"max_tokens"`
	Temperature float32        `json:"temperature"`
	Options     map[string]any `json:"options"`
}

// ProviderEndpoints contains default API endpoints for each provider
var ProviderEndpoints = map[ProviderType]string{
	ProviderOpenRouter: "https://openrouter.ai/api/v1",
	ProviderZhipuAI:    "https://api.z.ai/api/coding/paas/v4",
}

// DefaultModels contains recommended default models per provider
var DefaultModels = map[ProviderType]string{
	ProviderOpenRouter: "amazon/nova-2-lite-v1:free",
	ProviderZhipuAI:    "glm-4.6",
}

// ============================================================================
// Developer-configured provider settings (hardcoded for development)
// TODO: Move to config file or environment variables for production
// ============================================================================

// DevConfig holds all hardcoded development configurations
type DevConfig struct {
	Chat             ProviderConfig
	Title            ProviderConfig
	Summary          ProviderConfig
	OCRCleanup       ProviderConfig
	WhisperCleanup   ProviderConfig
	UseLocalFallback bool
}

// GetDefaultDevConfig returns the default development configuration
func GetDefaultDevConfig() DevConfig {
	return DevConfig{
		Chat: ProviderConfig{
			Type:      ProviderOpenRouter,
			Name:      "OpenRouter",
			APIKey:    "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",
			Model:     "amazon/nova-2-lite-v1:free",
			MaxTokens: 4096,
			Options:   map[string]any{"app_name": "Veridium"},
		},
		Title: ProviderConfig{
			Type:      ProviderZhipuAI,
			Name:      "Zhipu GLM",
			APIKey:    "a10854167085448cac33753523919ac9.D41CLq6KxXTY7g4u",
			Model:     "glm-4.6",
			MaxTokens: 256,
		},
		Summary: ProviderConfig{
			Type:      ProviderZhipuAI,
			Name:      "Zhipu GLM",
			APIKey:    "a10854167085448cac33753523919ac9.D41CLq6KxXTY7g4u",
			Model:     "glm-4.6",
			MaxTokens: 1024,
		},
		OCRCleanup: ProviderConfig{
			Type:      ProviderOpenRouter,
			Name:      "OpenRouter",
			APIKey:    "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",
			Model:     "amazon/nova-2-lite-v1:free",
			MaxTokens: 2048,
			Options:   map[string]any{"app_name": "Veridium"},
		},
		WhisperCleanup: ProviderConfig{
			Type:      ProviderOpenRouter,
			Name:      "OpenRouter",
			APIKey:    "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",
			Model:     "amazon/nova-2-lite-v1:free",
			MaxTokens: 16384,
			Options:   map[string]any{"app_name": "Veridium"},
		},
		UseLocalFallback: true,
	}
}

// BuildTaskRouter creates a TaskRouter from DevConfig
func BuildTaskRouter(config DevConfig, toolRegistry *tools.ToolRegistry, localModel fantasy.LanguageModel) *TaskRouter {
	factory := NewProviderFactory(toolRegistry)
	router := NewTaskRouter(nil)

	if config.UseLocalFallback && localModel != nil {
		router.setFallback(localModel)
		log.Printf("🔀 TaskRouter: Local Llama set as fallback model")
	}

	// Helper to configure a task model
	configureTask := func(task TaskType, cfg ProviderConfig, taskName string) {
		if cfg.Type == ProviderLlama {
			if localModel != nil {
				router.SetModel(task, localModel)
				log.Printf("🔀 TaskRouter: %s -> Local Llama", taskName)
			}
			return
		}
		if cfg.APIKey != "" {
			model, err := factory.CreateLanguageModel(cfg)
			if err != nil {
				log.Printf("⚠️  Failed to create %s model: %v", taskName, err)
				return
			}
			router.SetModel(task, model)
			log.Printf("🔀 TaskRouter: %s -> %s (%s)", taskName, cfg.Type, cfg.Model)
		} else if localModel != nil {
			router.SetModel(task, localModel)
			log.Printf("🔀 TaskRouter: %s -> Local Llama (no API key)", taskName)
		}
	}

	configureTask(TaskChat, config.Chat, "Chat")
	configureTask(TaskTitleGen, config.Title, "Title")
	configureTask(TaskSummaryGen, config.Summary, "Summary")
	configureTask(TaskOCRCleanup, config.OCRCleanup, "OCRCleanup")
	configureTask(TaskTranscriptCleanup, config.WhisperCleanup, "TranscriptCleanup")

	return router
}
