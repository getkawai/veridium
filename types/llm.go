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

package types

// ============================================================================
// LLM Provider Configuration Types
// ============================================================================
// Note: LLMProvider interface is defined in internal/llm/provider.go
// to avoid circular imports with pkg/yzma/message

// ProviderType identifies the LLM provider type
type ProviderType string

const (
	ProviderLlama      ProviderType = "llama"      // Local llama.cpp (via yzma)
	ProviderOpenRouter ProviderType = "openrouter" // OpenRouter (OpenAI-compatible, multi-model)
	ProviderZhipuAI    ProviderType = "zhipu"      // Zhipu GLM (OpenAI-compatible, Chinese AI)
)

// ProviderConfig holds configuration for an LLM provider
type ProviderConfig struct {
	Type        ProviderType   `json:"type"`
	Name        string         `json:"name"`        // Display name
	BaseURL     string         `json:"base_url"`    // API endpoint (empty for local llama)
	APIKey      string         `json:"api_key"`     // API key (empty for local llama)
	Model       string         `json:"model"`       // Model name/path
	MaxTokens   int            `json:"max_tokens"`  // Max tokens for generation
	Temperature float32        `json:"temperature"` // Temperature (0.0-2.0)
	Options     map[string]any `json:"options"`     // Provider-specific options
}

// ProviderEndpoints contains default API endpoints for each provider
var ProviderEndpoints = map[ProviderType]string{
	ProviderOpenRouter: "https://openrouter.ai/api/v1",
	ProviderZhipuAI:    "https://open.bigmodel.cn/api/paas/v4",
}

// DefaultModels contains recommended default models per provider
var DefaultModels = map[ProviderType]string{
	ProviderOpenRouter: "anthropic/claude-3.5-sonnet", // Best general-purpose model
	ProviderZhipuAI:    "glm-4-flash",                 // Fast and cost-effective
}
