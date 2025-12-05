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

// Package openai provides an OpenAI-compatible API client that works with
// OpenRouter and Zhipu GLM providers.
package openai

import (
	"github.com/kawai-network/veridium/types"
)

// ProviderHeaders returns provider-specific HTTP headers
func ProviderHeaders(providerType types.ProviderType, apiKey string, options map[string]any) map[string]string {
	headers := map[string]string{
		"Content-Type":  "application/json",
		"Authorization": "Bearer " + apiKey,
	}

	switch providerType {
	case types.ProviderOpenRouter:
		// OpenRouter specific headers
		if appName, ok := options["app_name"].(string); ok {
			headers["X-Title"] = appName
		} else {
			headers["X-Title"] = "Veridium"
		}
		if siteURL, ok := options["site_url"].(string); ok {
			headers["HTTP-Referer"] = siteURL
		}

	case types.ProviderZhipuAI:
		// Zhipu uses standard Bearer token for v4 API
		// No additional headers needed
	}

	return headers
}
