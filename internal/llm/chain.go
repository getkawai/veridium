// Package llm provides language model utilities including chain building.
package llm

import (
	"context"
	"log"

	"github.com/getkawai/unillm"
	googleprovider "github.com/getkawai/unillm/providers/google"
	"github.com/getkawai/unillm/providers/openaicompat"
	"github.com/getkawai/unillm/providers/openrouter"
	"github.com/kawai-network/x/constant"
)

// BuildModelChain creates a fallback chain of language models.
// It tries providers in order: Google Gemini → OpenRouter → Pollinations → ZAI → Local
func BuildModelChain(bgCtx context.Context, localModel unillm.LanguageModel, criteria openrouter.ModelSelectionCriteria, taskName string) []unillm.LanguageModel {
	var chain []unillm.LanguageModel

	// 1. Google Gemini 2.5 Flash-Lite (free tier with highest limits: 15 RPM, 1000 RPD)
	if apiKey := constant.GetRandomGeminiApiKey(); apiKey != "" {
		log.Printf("🔍 %s: Initializing Google Gemini 2.5 Flash-Lite...", taskName)
		if provider, err := googleprovider.New(googleprovider.WithGeminiAPIKey(apiKey)); err == nil {
			if geminiModel, err := provider.LanguageModel(bgCtx, "gemini-2.5-flash-lite"); err == nil {
				chain = append(chain, geminiModel)
				log.Printf("✅ %s: Added Google Gemini (gemini-2.5-flash-lite) to chain [15 RPM, 1000 RPD]", taskName)
			} else {
				log.Printf("❌ %s: Google Gemini provider initialized but failed to get model: %v", taskName, err)
			}
		} else {
			log.Printf("❌ %s: Failed to initialize Google Gemini provider: %v", taskName, err)
		}
	} else {
		log.Printf("ℹ️  %s: Skipping Google Gemini (no API key)", taskName)
	}

	// 2. OpenRouter (free tier)
	if apiKey := constant.GetRandomOpenRouterApiKey(); apiKey != "" {
		log.Printf("🔍 %s: Initializing OpenRouter...", taskName)
		if provider, err := openrouter.New(openrouter.WithAPIKey(apiKey), openrouter.WithModelSelection(criteria)); err == nil {
			if remoteModel, err := provider.LanguageModel(bgCtx, ""); err == nil {
				chain = append(chain, remoteModel)
				catalog := openrouter.GetCatalog()
				if selected := catalog.SelectFreeModel(criteria); selected != nil {
					log.Printf("✅ %s: Added OpenRouter (%s) to chain", taskName, selected.ID)
				} else {
					log.Printf("⚠️  %s: OpenRouter initialized but no free model matched criteria", taskName)
				}
			} else {
				log.Printf("❌ %s: OpenRouter provider initialized but failed to get model: %v", taskName, err)
			}
		} else {
			log.Printf("❌ %s: Failed to initialize OpenRouter provider: %v", taskName, err)
		}
	} else {
		log.Printf("ℹ️  %s: Skipping OpenRouter (no API key)", taskName)
	}

	// 3. Pollinations AI (fallback before local)
	if provider, err := openaicompat.New(
		openaicompat.WithName("pollinations"),
		openaicompat.WithBaseURL("https://text.pollinations.ai/openai"),
		openaicompat.WithAPIKey("dummy"), // Pollinations doesn't require API key, but SDK needs one
	); err == nil {
		if pollinationsModel, err := provider.LanguageModel(bgCtx, "openai"); err == nil {
			chain = append(chain, pollinationsModel)
			log.Printf("%s: Pollinations AI (openai)", taskName)
		} else {
			log.Printf("❌ %s: Pollinations provider initialized but failed to get model: %v", taskName, err)
		}
	} else {
		log.Printf("❌ %s: Failed to initialize Pollinations provider: %v", taskName, err)
	}

	// 4. ZAI GLM-4.7 (fallback before local)
	if apiKey := constant.GetRandomZaiApiKey(); apiKey != "" {
		if provider, err := openaicompat.New(
			openaicompat.WithName("zai"),
			openaicompat.WithBaseURL("https://api.z.ai/api/coding/paas/v4"),
			openaicompat.WithAPIKey(apiKey),
		); err == nil {
			if zaiModel, err := provider.LanguageModel(bgCtx, "glm-4.7"); err == nil {
				chain = append(chain, zaiModel)
				log.Printf("%s: ZAI (glm-4.7)", taskName)
			}
		}
	}

	// 5. Local model (final fallback)
	chain = append(chain, localModel)
	log.Printf("%s: Chain created with %d models (fallback: %s/%s)", taskName, len(chain), localModel.Provider(), localModel.Model())
	return chain
}
