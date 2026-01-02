package app

import (
	"context"
	"log"

	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/pkg/fantasy"
	llamaprovider "github.com/kawai-network/veridium/pkg/fantasy/providers/llama"
	"github.com/kawai-network/veridium/pkg/fantasy/providers/openaicompat"
	"github.com/kawai-network/veridium/pkg/fantasy/providers/openrouter"
	"github.com/kawai-network/veridium/pkg/fantasy/providers/pooled"
)

// buildModelChainV2 creates a model chain with CLIProxyAPI's fallback mechanism.
// This version uses account pooling and smart error handling.
func (ctx *Context) buildModelChainV2(bgCtx context.Context, localModel fantasy.LanguageModel, criteria openrouter.ModelSelectionCriteria, taskName string) []fantasy.LanguageModel {
	var chain []fantasy.LanguageModel

	// 1. OpenRouter with multiple API keys (pooled)
	openRouterKeys := constant.GetOpenRouterApiKeys() // Assume this returns []string
	if len(openRouterKeys) > 0 {
		pooledProvider, err := pooled.New(pooled.Config{
			ProviderName: "openrouter",
			BaseURL:      "https://openrouter.ai/api/v1",
			ModelName:    "auto", // Will be selected by criteria
			APIKeys:      openRouterKeys,
			CreateClient: func(apiKey string) (fantasy.LanguageModel, error) {
				provider, err := openrouter.New(
					openrouter.WithAPIKey(apiKey),
					openrouter.WithModelSelection(criteria),
				)
				if err != nil {
					return nil, err
				}
				return provider.LanguageModel(bgCtx, "")
			},
		})

		if err == nil {
			chain = append(chain, pooledProvider)
			catalog := openrouter.GetCatalog()
			if selected := catalog.SelectFreeModel(criteria); selected != nil {
				log.Printf("%s: OpenRouter Pooled (%s) with %d keys", taskName, selected.ID, len(openRouterKeys))
			}
		} else {
			log.Printf("Warning: Failed to create pooled OpenRouter: %v", err)
		}
	}

	// 2. Pollinations AI (no pooling needed, free service)
	if provider, err := openaicompat.New(
		openaicompat.WithName("pollinations"),
		openaicompat.WithBaseURL("https://text.pollinations.ai/openai"),
		openaicompat.WithAPIKey("dummy"),
	); err == nil {
		if pollinationsModel, err := provider.LanguageModel(bgCtx, "openai"); err == nil {
			chain = append(chain, pollinationsModel)
			log.Printf("%s: Pollinations AI (openai)", taskName)
		}
	}

	// 3. ZAI with multiple API keys (pooled)
	zaiKeys := constant.GetZaiApiKeys() // Assume this returns []string
	if len(zaiKeys) > 0 {
		pooledProvider, err := pooled.New(pooled.Config{
			ProviderName: "zai",
			BaseURL:      "https://api.z.ai/api/coding/paas/v4",
			ModelName:    "glm-4.7",
			APIKeys:      zaiKeys,
			CreateClient: func(apiKey string) (fantasy.LanguageModel, error) {
				provider, err := openaicompat.New(
					openaicompat.WithName("zai"),
					openaicompat.WithBaseURL("https://api.z.ai/api/coding/paas/v4"),
					openaicompat.WithAPIKey(apiKey),
				)
				if err != nil {
					return nil, err
				}
				return provider.LanguageModel(bgCtx, "glm-4.7")
			},
		})

		if err == nil {
			chain = append(chain, pooledProvider)
			log.Printf("%s: ZAI Pooled (glm-4.7) with %d keys", taskName, len(zaiKeys))
		} else {
			log.Printf("Warning: Failed to create pooled ZAI: %v", err)
		}
	}

	// 4. Local model (final fallback)
	chain = append(chain, localModel)
	log.Printf("%s: Chain created with %d models (fallback: %s/%s)", taskName, len(chain), localModel.Provider(), localModel.Model())
	return chain
}

// InitLanguageModelsV2 initializes language models with pooled providers.
func (ctx *Context) InitLanguageModelsV2() {
	if ctx.LibService == nil {
		return
	}

	bgCtx := context.Background()
	llamaProvider, err := llamaprovider.New(
		llamaprovider.WithService(ctx.LibService),
		llamaprovider.WithToolRegistry(ctx.ToolRegistry),
	)
	if err != nil {
		log.Printf("Warning: Llama provider failed: %v", err)
		return
	}

	localModel, err := llamaProvider.LanguageModel(bgCtx, "")
	if err != nil {
		log.Printf("Warning: Local LLM failed: %v", err)
		return
	}

	// Circuit breaker: skip rate-limited models until app restart
	circuitBreaker := fantasy.WithCircuitBreaker(1, 0)

	var err error
	ctx.ChatModel, err = fantasy.NewChain(ctx.buildModelChainV2(bgCtx, localModel, openrouter.ModelSelectionCriteria{
		RequireReasoning: true, RequireAttachments: true, MinContextWindow: 100000,
	}, "ChatModel"), circuitBreaker)
	if err != nil {
		log.Printf("Warning: ChatModel chain creation failed: %v", err)
	}

	ctx.TitleModel, err = fantasy.NewChain(ctx.buildModelChainV2(bgCtx, localModel, openrouter.ModelSelectionCriteria{}, "TitleModel"), circuitBreaker)
	if err != nil {
		log.Printf("Warning: TitleModel chain creation failed: %v", err)
	}

	ctx.SummaryModel, err = fantasy.NewChain(ctx.buildModelChainV2(bgCtx, localModel, openrouter.ModelSelectionCriteria{
		MinContextWindow: 50000,
	}, "SummaryModel"), circuitBreaker)
	if err != nil {
		log.Printf("Warning: SummaryModel chain creation failed: %v", err)
	}

	ctx.CleanupModel, err = fantasy.NewChain(ctx.buildModelChainV2(bgCtx, localModel, openrouter.ModelSelectionCriteria{}, "CleanupModel"), circuitBreaker)
	if err != nil {
		log.Printf("Warning: CleanupModel chain creation failed: %v", err)
	}

	log.Printf("Language models initialized with pooled providers")
}
