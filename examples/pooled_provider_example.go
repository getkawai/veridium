// Package main demonstrates how to use the pooled provider with CLIProxyAPI fallback mechanism.
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/kawai-network/veridium/internal/constant"
	"github.com/kawai-network/veridium/pkg/fantasy"
	"github.com/kawai-network/veridium/pkg/fantasy/providers/openrouter"
	"github.com/kawai-network/veridium/pkg/fantasy/providers/pooled"
)

func main() {
	ctx := context.Background()

	// Example 1: Create a pooled OpenRouter provider
	fmt.Println("=== Example 1: Pooled OpenRouter Provider ===")
	pooledOpenRouter := createPooledOpenRouter(ctx)
	if pooledOpenRouter != nil {
		testProvider(ctx, pooledOpenRouter, "OpenRouter")
	}

	// Example 2: Monitor account status
	fmt.Println("\n=== Example 2: Monitor Account Status ===")
	if pooledOpenRouter != nil {
		monitorAccounts(pooledOpenRouter)
	}

	// Example 3: Use in a chain with fallback
	fmt.Println("\n=== Example 3: Chain with Fallback ===")
	chain := createChainWithPooling(ctx)
	if chain != nil {
		testProvider(ctx, chain, "Chain")
	}
}

// createPooledOpenRouter creates a pooled provider with multiple OpenRouter API keys.
func createPooledOpenRouter(ctx context.Context) *pooled.PooledProvider {
	// Get all OpenRouter API keys
	keys := constant.GetOpenRouterApiKeys()
	if len(keys) == 0 {
		log.Println("No OpenRouter API keys available")
		return nil
	}

	fmt.Printf("Creating pooled provider with %d API keys\n", len(keys))

	// Create pooled provider
	provider, err := pooled.New(pooled.Config{
		ProviderName: "openrouter",
		BaseURL:      "https://openrouter.ai/api/v1",
		ModelName:    "auto",
		APIKeys:      keys,
		CreateClient: func(apiKey string) (fantasy.LanguageModel, error) {
			// Create OpenRouter provider with this API key
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
		log.Printf("Failed to create pooled provider: %v", err)
		return nil
	}

	fmt.Println("✅ Pooled provider created successfully")
	return provider
}

// testProvider tests a provider with a simple request.
func testProvider(ctx context.Context, provider fantasy.LanguageModel, name string) {
	fmt.Printf("\nTesting %s provider...\n", name)

	// Create a simple test call
	call := fantasy.Call{
		Prompt: fantasy.Prompt{
			fantasy.NewUserMessage("Say 'Hello from pooled provider!' in one sentence."),
		},
	}

	// Make the request
	resp, err := provider.Generate(ctx, call)
	if err != nil {
		log.Printf("❌ Request failed: %v", err)
		return
	}

	// Print response
	fmt.Printf("✅ Response received:\n")
	for _, content := range resp.Content {
		if content.GetType() == fantasy.ContentTypeText {
			if textContent, ok := content.(*fantasy.TextContent); ok {
				fmt.Printf("   %s\n", textContent.Text)
			}
		}
	}
	fmt.Printf("   Tokens: %d input, %d output, %d total\n",
		resp.Usage.InputTokens, resp.Usage.OutputTokens, resp.Usage.TotalTokens)
}

// monitorAccounts shows how to monitor account status.
func monitorAccounts(provider *pooled.PooledProvider) {
	// Get the underlying auth manager
	manager := provider.GetManager()

	// List all accounts
	accounts := manager.List()

	fmt.Printf("\nMonitoring %d accounts:\n", len(accounts))
	for i, account := range accounts {
		fmt.Printf("\n[Account %d] %s\n", i+1, account.Label)
		fmt.Printf("  Provider: %s\n", account.Provider)
		fmt.Printf("  Status: %s\n", account.Status)
		fmt.Printf("  Disabled: %v\n", account.Disabled)
		fmt.Printf("  Unavailable: %v\n", account.Unavailable)

		if account.Quota.Exceeded {
			fmt.Printf("  ⚠️  Quota Exceeded: true\n")
			fmt.Printf("     Reason: %s\n", account.Quota.Reason)
			if !account.Quota.NextRecoverAt.IsZero() {
				fmt.Printf("     Recover At: %s\n", account.Quota.NextRecoverAt)
			}
			fmt.Printf("     Backoff Level: %d\n", account.Quota.BackoffLevel)
		}

		if !account.NextRetryAfter.IsZero() {
			fmt.Printf("  Next Retry: %s\n", account.NextRetryAfter)
		}

		if account.LastError != nil {
			fmt.Printf("  Last Error: %s\n", account.LastError.Message)
		}

		// Per-model state
		if len(account.ModelStates) > 0 {
			fmt.Printf("  Model States:\n")
			for model, state := range account.ModelStates {
				if state.Unavailable {
					fmt.Printf("    - %s: Unavailable until %s\n", model, state.NextRetryAfter)
				} else {
					fmt.Printf("    - %s: Available\n", model)
				}
			}
		}
	}
}

// createChainWithPooling creates a fantasy chain with pooled providers.
func createChainWithPooling(ctx context.Context) fantasy.LanguageModel {
	var models []fantasy.LanguageModel

	// Add pooled OpenRouter
	pooledOR := createPooledOpenRouter(ctx)
	if pooledOR != nil {
		models = append(models, pooledOR)
	}

	// Add other providers...
	// (In real usage, you'd add more providers here)

	if len(models) == 0 {
		log.Println("No models available for chain")
		return nil
	}

	// Create chain with circuit breaker
	chain, err := fantasy.NewChain(models, fantasy.WithCircuitBreaker(1, 0))
	if err != nil {
		log.Printf("Failed to create chain: %v", err)
		return nil
	}

	fmt.Printf("✅ Chain created with %d models\n", len(models))
	return chain
}

