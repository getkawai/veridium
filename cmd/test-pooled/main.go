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
	"github.com/kawai-network/veridium/pkg/fantasy/tools"
	"github.com/kawai-network/veridium/pkg/fantasy/tools/builtin"
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
	fmt.Println("=" + string(make([]byte, 60)))
	fmt.Println()

	// Test 4: Tool Calling
	fmt.Println("📊 Test 4: Tool Calling with Pooled Provider")
	fmt.Println("-" + string(make([]byte, 60)))
	testToolCalling(ctx)

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

func testToolCalling(ctx context.Context) {
	// Get OpenRouter API keys
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
	fmt.Println()

	// Create tool registry and register builtin tools
	registry := tools.NewToolRegistry()

	// Register calculator tool
	if err := builtin.RegisterCalculator(registry); err != nil {
		log.Printf("❌ Failed to register calculator: %v\n", err)
		return
	}

	// Register web search tool (using WebSearch instead of WebBrowsing)
	if err := builtin.RegisterWebSearch(registry); err != nil {
		log.Printf("⚠️  Failed to register web search: %v (skipping)\n", err)
		// Don't return, just skip this tool
	}

	fmt.Println("✅ Registered builtin tools")
	fmt.Println()

	// Convert AgentTools to FunctionTools
	calculatorAgentTool, ok := registry.Get("calculator")
	if !ok {
		log.Println("❌ Calculator tool not found")
		return
	}
	calculatorTool := agentToolToFunctionTool(calculatorAgentTool)

	// Try to get web search tool (might not be available)
	var searchTool *fantasy.FunctionTool
	if searchAgentTool, ok := registry.Get("web_search"); ok {
		tool := agentToolToFunctionTool(searchAgentTool)
		searchTool = &tool
	}

	// Test 1: Simple calculator expression
	fmt.Println("🧪 Test 1: Calculator Tool (Simple)")
	testPromptWithTools(ctx, provider,
		"Use the calculator tool with expression '156 + 844'",
		[]fantasy.Tool{calculatorTool})

	fmt.Println()
	fmt.Println("-" + string(make([]byte, 60)))
	fmt.Println()

	// Test 2: Calculator with sqrt
	fmt.Println("🧪 Test 2: Calculator Tool (sqrt)")
	testPromptWithTools(ctx, provider,
		"Call calculator tool with expression 'sqrt(144)'",
		[]fantasy.Tool{calculatorTool})

	fmt.Println()
	fmt.Println("-" + string(make([]byte, 60)))
	fmt.Println()

	// Test 3: Web search (if available)
	if searchTool != nil {
		fmt.Println("🧪 Test 3: Web Search Tool")
		testPromptWithTools(ctx, provider,
			"Search for 'latest AI news 2024' using the web search tool.",
			[]fantasy.Tool{*searchTool})

		fmt.Println()
		fmt.Println("-" + string(make([]byte, 60)))
		fmt.Println()

		// Test 4: Multiple tools
		fmt.Println("🧪 Test 4: Multiple Tools")
		testPromptWithTools(ctx, provider,
			"Calculate 50 * 20, then search for information about that number.",
			[]fantasy.Tool{calculatorTool, *searchTool})
	} else {
		fmt.Println("⚠️  Skipping web search tests (tool not available)")
	}
}

func testPromptWithTools(ctx context.Context, provider fantasy.LanguageModel, prompt string, tools []fantasy.Tool) {
	fmt.Printf("📝 Prompt: %s\n", prompt)
	fmt.Printf("🔧 Tools: %d available\n", len(tools))
	for _, tool := range tools {
		if ft, ok := tool.(fantasy.FunctionTool); ok {
			fmt.Printf("   - %s\n", ft.Name)
		}
	}
	fmt.Println()

	// Create call with tools
	call := fantasy.Call{
		Prompt: fantasy.Prompt{
			fantasy.NewUserMessage(prompt),
		},
		Tools:           tools,
		MaxOutputTokens: ptr(int64(500)),
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

	// Check finish reason
	fmt.Printf("🏁 Finish Reason: %s\n", resp.FinishReason)

	// Check if tools were called
	hasToolCalls := false
	toolCallCount := 0
	textContent := ""

	for _, part := range resp.Content {
		switch c := part.(type) {
		case fantasy.TextContent:
			textContent += c.Text
		case fantasy.ToolCallContent:
			hasToolCalls = true
			toolCallCount++
			fmt.Printf("\n🔧 Tool Call #%d:\n", toolCallCount)
			fmt.Printf("   Tool: %s\n", c.ToolName)
			fmt.Printf("   ID: %s\n", c.ToolCallID)
			fmt.Printf("   Arguments: %s\n", c.Input)
		}
	}

	if hasToolCalls {
		fmt.Printf("\n✅ SUCCESS: %d tool call(s) detected!\n", toolCallCount)
		if textContent != "" {
			fmt.Printf("\n💬 Additional Text: %s\n", textContent)
		}
	} else {
		fmt.Println("\n⚠️  NO TOOL CALLS - Model responded with text only")
		if textContent != "" {
			fmt.Printf("\n📄 Response:\n%s\n", textContent)
		}
	}

	fmt.Printf("\n📊 Usage: %d input, %d output, %d total tokens\n",
		resp.Usage.InputTokens,
		resp.Usage.OutputTokens,
		resp.Usage.TotalTokens,
	)
}

// agentToolToFunctionTool converts fantasy.AgentTool to fantasy.FunctionTool
func agentToolToFunctionTool(agentTool fantasy.AgentTool) fantasy.FunctionTool {
	info := agentTool.Info()
	return fantasy.FunctionTool{
		Name:        info.Name,
		Description: info.Description,
		InputSchema: map[string]any{
			"type":       "object",
			"properties": info.Parameters,
			"required":   info.Required,
		},
	}
}

func ptr[T any](v T) *T {
	return &v
}
