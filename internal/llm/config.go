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

// DevConfig holds all hardcoded development configurations
// Uses types.ProviderConfig directly - no duplicate types
type DevConfig struct {
	Chat              types.ProviderConfig
	Title             types.ProviderConfig
	Summary           types.ProviderConfig
	OCRCleanup        types.ProviderConfig
	TranscriptCleanup types.ProviderConfig
	UseLocalFallback  bool
}

// GetDefaultDevConfig returns the default development configuration
func GetDefaultDevConfig() DevConfig {
	return DevConfig{
		Chat: types.ProviderConfig{
			Type:      types.ProviderOpenRouter,
			Name:      "OpenRouter",
			APIKey:    "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",
			Model:     "amazon/nova-2-lite-v1:free",
			MaxTokens: 4096,
			Options:   map[string]any{"app_name": "Veridium"},
		},
		Title: types.ProviderConfig{
			Type:      types.ProviderZhipuAI,
			Name:      "Zhipu GLM",
			APIKey:    "a10854167085448cac33753523919ac9.D41CLq6KxXTY7g4u",
			Model:     "glm-4.6",
			MaxTokens: 256,
		},
		Summary: types.ProviderConfig{
			Type:      types.ProviderZhipuAI,
			Name:      "Zhipu GLM",
			APIKey:    "a10854167085448cac33753523919ac9.D41CLq6KxXTY7g4u",
			Model:     "glm-4.6",
			MaxTokens: 1024,
		},
		OCRCleanup: types.ProviderConfig{
			Type:      types.ProviderZhipuAI,
			Name:      "Zhipu GLM",
			APIKey:    "a10854167085448cac33753523919ac9.D41CLq6KxXTY7g4u",
			Model:     "glm-4.6",
			MaxTokens: 2048,
		},
		TranscriptCleanup: types.ProviderConfig{
			Type:      types.ProviderZhipuAI,
			Name:      "Zhipu GLM",
			APIKey:    "a10854167085448cac33753523919ac9.D41CLq6KxXTY7g4u",
			Model:     "glm-4.6",
			MaxTokens: 16384,
		},
		UseLocalFallback: true,
	}
}

// BuildTaskRouter creates a TaskRouter from DevConfig
func BuildTaskRouter(config DevConfig, toolRegistry *tools.ToolRegistry, localProvider Provider) *TaskRouter {
	factory := NewProviderFactory(toolRegistry)
	router := NewTaskRouter(toolRegistry, nil)

	if config.UseLocalFallback && localProvider != nil {
		router.SetFallback(localProvider)
		log.Printf("🔀 TaskRouter: Local Llama set as fallback provider")
	}

	// Helper to configure a task provider
	configureTask := func(task TaskType, cfg types.ProviderConfig, taskName string) {
		if cfg.Type == types.ProviderLlama {
			if localProvider != nil {
				router.SetProvider(task, localProvider)
				log.Printf("🔀 TaskRouter: %s -> Local Llama", taskName)
			}
			return
		}
		if cfg.APIKey != "" {
			provider, err := factory.CreateProvider(cfg)
			if err != nil {
				log.Printf("⚠️  Failed to create %s provider: %v", taskName, err)
				return
			}
			router.SetProvider(task, provider)
			log.Printf("🔀 TaskRouter: %s -> %s (%s)", taskName, cfg.Type, cfg.Model)
		} else if localProvider != nil {
			router.SetProvider(task, localProvider)
			log.Printf("🔀 TaskRouter: %s -> Local Llama (no API key)", taskName)
		}
	}

	configureTask(TaskChat, config.Chat, "Chat")
	configureTask(TaskTitleGen, config.Title, "Title")
	configureTask(TaskSummaryGen, config.Summary, "Summary")
	configureTask(TaskOCRCleanup, config.OCRCleanup, "OCRCleanup")
	configureTask(TaskTranscriptCleanup, config.TranscriptCleanup, "TranscriptCleanup")

	return router
}

// UpdateProvider updates a task's provider configuration at runtime
func (r *TaskRouter) UpdateProvider(task TaskType, config types.ProviderConfig) error {
	factory := NewProviderFactory(r.toolRegistry)
	provider, err := factory.CreateProvider(config)
	if err != nil {
		return &RouterError{Message: "failed to create provider: " + err.Error()}
	}
	r.SetProvider(task, provider)
	return nil
}
