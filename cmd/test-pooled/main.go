// Package main provides a CLI tool to test the pooled provider with CLIProxyAPI fallback.
package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/pkg/fantasy"
	"github.com/kawai-network/veridium/pkg/fantasy/providers/openaicompat"
	"github.com/kawai-network/veridium/pkg/fantasy/providers/openrouter"
	"github.com/kawai-network/veridium/pkg/fantasy/providers/pooled"
)

func main() {
	fmt.Println("🚀 Testing Pooled Provider with CLIProxyAPI Fallback")
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println()

	ctx := context.Background()

	// Test 1: OpenRouter Pooled
	fmt.Println("📊 Test 1: OpenRouter with Multiple API Keys")
	fmt.Println("-" + string(make([]byte, 60)))
	testOpenRouterPooled(ctx)

	fmt.Println()
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println()

	// Test 2: ZAI Pooled
	fmt.Println("📊 Test 2: ZAI with Multiple API Keys")
	fmt.Println("-" + string(make([]byte, 60)))
	testZAIPooled(ctx)

	fmt.Println()
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println()

	// Test 3: Chain with Fallback
	fmt.Println("📊 Test 3: Full Chain with Fallback")
	fmt.Println("-" + string(make([]byte, 60)))
	testChainWithFallback(ctx)

	fmt.Println()
	fmt.Println("✅ All tests completed!")
}

func testOpenRouterPooled(ctx context.Context) {
	// Get all OpenRouter API keys
	keys := constant.GetOpenRouterApiKeys()
	if len(keys) == 0 {
		log.Println("⚠️  No OpenRouter API keys available")
		return
	}

	fmt.Printf("🔑 Found %d OpenRouter API keys\n", len(keys))

	// Create pooled provider
	provider, err := pooled.New(pooled.Config{
		ProviderName: "openrouter",
		BaseURL:      "https://openrouter.ai/api/v1",
		ModelName:    "auto",
		APIKeys:      keys,
		CreateClient: func(apiKey string) (fantasy.LanguageModel, error) {
			p, err := openrouter.New(
				openrouter.WithAPIKey(apiKey),
				openrouter.WithModelSelection(openrouter.ModelSelectionCriteria{
					RequireReasoning:   false,
					RequireAttachments: false,
					MinContextWindow:   8000,
				}),
			)
			if err != nil {
				return nil, err
			}
			return p.LanguageModel(ctx, "")
		},
	})

	if err != nil {
		log.Printf("❌ Failed to create pooled provider: %v\n", err)
		return
	}

	fmt.Println("✅ Pooled provider created")

	// Test the provider
	testProvider(ctx, provider, "OpenRouter Pooled")

	// Show account status
	fmt.Println()
	showAccountStatus(provider)
}

func testZAIPooled(ctx context.Context) {
	// Get all ZAI API keys
	keys := constant.GetZaiApiKeys()
	if len(keys) == 0 {
		log.Println("⚠️  No ZAI API keys available")
		return
	}

	fmt.Printf("🔑 Found %d ZAI API keys\n", len(keys))

	// Create pooled provider
	provider, err := pooled.New(pooled.Config{
		ProviderName: "zai",
		BaseURL:      "https://api.z.ai/api/coding/paas/v4",
		ModelName:    "glm-4.7",
		APIKeys:      keys,
		CreateClient: func(apiKey string) (fantasy.LanguageModel, error) {
			p, err := openaicompat.New(
				openaicompat.WithName("zai"),
				openaicompat.WithBaseURL("https://api.z.ai/api/coding/paas/v4"),
				openaicompat.WithAPIKey(apiKey),
			)
			if err != nil {
				return nil, err
			}
			return p.LanguageModel(ctx, "glm-4.7")
		},
	})

	if err != nil {
		log.Printf("❌ Failed to create pooled provider: %v\n", err)
		return
	}

	fmt.Println("✅ Pooled provider created")

	// Test the provider
	testProvider(ctx, provider, "ZAI Pooled")

	// Show account status
	fmt.Println()
	showAccountStatus(provider)
}

func testChainWithFallback(ctx context.Context) {
	var models []fantasy.LanguageModel

	// Add OpenRouter pooled
	openRouterKeys := constant.GetOpenRouterApiKeys()
	if len(openRouterKeys) > 0 {
		provider, err := pooled.New(pooled.Config{
			ProviderName: "openrouter",
			APIKeys:      openRouterKeys,
			CreateClient: func(apiKey string) (fantasy.LanguageModel, error) {
				p, err := openrouter.New(openrouter.WithAPIKey(apiKey))
				if err != nil {
					return nil, err
				}
				return p.LanguageModel(ctx, "")
			},
		})
		if err == nil {
			models = append(models, provider)
			fmt.Printf("✅ Added OpenRouter pooled (%d keys)\n", len(openRouterKeys))
		}
	}

	// Add Pollinations (free, no pooling needed)
	pollinations, err := openaicompat.New(
		openaicompat.WithName("pollinations"),
		openaicompat.WithBaseURL("https://text.pollinations.ai/openai"),
		openaicompat.WithAPIKey("dummy"),
	)
	if err == nil {
		if pollinationsModel, err := pollinations.LanguageModel(ctx, "openai"); err == nil {
			models = append(models, pollinationsModel)
			fmt.Println("✅ Added Pollinations AI")
		}
	}

	// Add ZAI pooled
	zaiKeys := constant.GetZaiApiKeys()
	if len(zaiKeys) > 0 {
		provider, err := pooled.New(pooled.Config{
			ProviderName: "zai",
			APIKeys:      zaiKeys,
			CreateClient: func(apiKey string) (fantasy.LanguageModel, error) {
				p, err := openaicompat.New(
					openaicompat.WithName("zai"),
					openaicompat.WithBaseURL("https://api.z.ai/api/coding/paas/v4"),
					openaicompat.WithAPIKey(apiKey),
				)
				if err != nil {
					return nil, err
				}
				return p.LanguageModel(ctx, "glm-4.7")
			},
		})
		if err == nil {
			models = append(models, provider)
			fmt.Printf("✅ Added ZAI pooled (%d keys)\n", len(zaiKeys))
		}
	}

	if len(models) == 0 {
		log.Println("❌ No models available for chain")
		return
	}

	// Create chain with circuit breaker
	chain, err := fantasy.NewChain(models, fantasy.WithCircuitBreaker(1, 0))
	if err != nil {
		log.Printf("❌ Failed to create chain: %v\n", err)
		return
	}

	fmt.Printf("✅ Chain created with %d models\n", len(models))
	fmt.Println()

	// Test the chain
	testProvider(ctx, chain, "Full Chain")
}

func testProvider(ctx context.Context, provider fantasy.LanguageModel, name string) {
	fmt.Printf("\n🧪 Testing %s...\n", name)

	// Create test prompt
	prompt := "What is recent news on AI? Give me a brief summary in 2-3 sentences."
	
	fmt.Printf("📝 Prompt: %s\n", prompt)
	fmt.Println()

	// Create call
	call := fantasy.Call{
		Prompt: fantasy.Prompt{
			fantasy.NewUserMessage(prompt),
		},
		MaxOutputTokens: ptr(int64(200)),
		Temperature:     ptr(0.7),
	}

	// Add timeout
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// Make request
	start := time.Now()
	resp, err := provider.Generate(ctx, call)
	duration := time.Since(start)

	if err != nil {
		log.Printf("❌ Request failed after %v: %v\n", duration, err)
		return
	}

	// Print response
	fmt.Printf("✅ Response received in %v\n", duration)
	fmt.Println()
	fmt.Println("📄 Response:")
	fmt.Println("---")

	// Use the Text() method which handles all content types
	fmt.Println(resp.Content.Text())

	fmt.Println("---")
	fmt.Println()
	fmt.Printf("📊 Usage: %d input tokens, %d output tokens, %d total\n",
		resp.Usage.InputTokens,
		resp.Usage.OutputTokens,
		resp.Usage.TotalTokens,
	)
}

func showAccountStatus(provider *pooled.PooledProvider) {
	fmt.Println("📈 Account Status:")
	fmt.Println()

	manager := provider.GetManager()
	accounts := manager.List()

	if len(accounts) == 0 {
		fmt.Println("  No accounts registered")
		return
	}

	for i, account := range accounts {
		fmt.Printf("  [%d] %s\n", i+1, account.Label)
		fmt.Printf("      Provider: %s\n", account.Provider)
		fmt.Printf("      Status: %s\n", account.Status)
		
		if account.Disabled {
			fmt.Println("      ⚠️  Disabled")
		}
		
		if account.Unavailable {
			fmt.Println("      ⚠️  Unavailable")
			if !account.NextRetryAfter.IsZero() {
				fmt.Printf("      Next Retry: %s\n", account.NextRetryAfter.Format(time.RFC3339))
			}
		}

		if account.Quota.Exceeded {
			fmt.Println("      ⚠️  Quota Exceeded")
			fmt.Printf("      Reason: %s\n", account.Quota.Reason)
			if !account.Quota.NextRecoverAt.IsZero() {
				fmt.Printf("      Recover At: %s\n", account.Quota.NextRecoverAt.Format(time.RFC3339))
			}
			fmt.Printf("      Backoff Level: %d\n", account.Quota.BackoffLevel)
		}

		if account.LastError != nil {
			fmt.Printf("      Last Error: %s\n", account.LastError.Message)
		}

		if len(account.ModelStates) > 0 {
			fmt.Println("      Model States:")
			for model, state := range account.ModelStates {
				status := "Available"
				if state.Unavailable {
					status = fmt.Sprintf("Unavailable until %s", state.NextRetryAfter.Format(time.RFC3339))
				}
				fmt.Printf("        - %s: %s\n", model, status)
			}
		}

		if i < len(accounts)-1 {
			fmt.Println()
		}
	}
}

func ptr[T any](v T) *T {
	return &v
}

