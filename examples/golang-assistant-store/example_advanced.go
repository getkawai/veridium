package main

import (
	"fmt"
	"log"
)

// ExampleAdvanced demonstrates advanced usage of AssistantStore
func ExampleAdvanced() {
	store := NewAssistantStore("")

	// 1. Get all categories with counts
	fmt.Println("=== Categories ===")
	categories, err := store.GetCategories("en-US")
	if err != nil {
		log.Fatal(err)
	}
	
	for category, count := range categories {
		fmt.Printf("- %s: %d agents\n", category, count)
	}
	fmt.Println()

	// 2. Search agents
	fmt.Println("=== Search Results for 'web' ===")
	searchResults, err := store.SearchAgents("en-US", "web")
	if err != nil {
		log.Fatal(err)
	}
	
	for i, agent := range searchResults {
		if i >= 3 {
			break
		}
		fmt.Printf("%d. %s\n", i+1, agent.Meta.Title)
		fmt.Printf("   Description: %s\n", agent.Meta.Description)
		fmt.Println()
	}

	// 3. Filter by category
	fmt.Println("=== Development Category ===")
	devAgents, err := store.FilterByCategory("en-US", "development")
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Found %d development agents\n", len(devAgents))
	for i, agent := range devAgents {
		if i >= 3 {
			break
		}
		fmt.Printf("- %s\n", agent.Meta.Title)
	}
	fmt.Println()

	// 4. Whitelist filtering
	fmt.Println("=== Whitelist Filtering ===")
	filter := &FilterOptions{
		Whitelist: []string{"web-development", "api-design", "code-review"},
	}
	
	whitelistedAgents, err := store.GetAgentIndexWithFilter("en-US", filter)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("Whitelisted agents: %d\n", len(whitelistedAgents))
	for _, agent := range whitelistedAgents {
		fmt.Printf("- %s (%s)\n", agent.Meta.Title, agent.Identifier)
	}
	fmt.Println()

	// 5. Fetch agent detail
	if len(whitelistedAgents) > 0 {
		fmt.Println("=== Agent Detail ===")
		identifier := whitelistedAgents[0].Identifier
		detail, err := store.GetAgent(identifier, "en-US")
		if err != nil {
			log.Printf("Failed to fetch detail: %v", err)
		} else {
			fmt.Printf("Identifier: %s\n", detail.Identifier)
			fmt.Printf("Title: %s\n", detail.Meta.Title)
			fmt.Printf("Author: %s\n", detail.Author)
			fmt.Printf("Tags: %v\n", detail.Meta.Tags)
			fmt.Printf("Plugins: %v\n", detail.Plugins)
			
			if len(detail.SystemRole) > 100 {
				fmt.Printf("System Role: %s...\n", detail.SystemRole[:100])
			} else {
				fmt.Printf("System Role: %s\n", detail.SystemRole)
			}
		}
	}
}

// ExampleMultiLocale demonstrates fetching agents in different locales
func ExampleMultiLocale() {
	store := NewAssistantStore("")

	locales := []string{"en-US", "zh-CN", "ja-JP"}

	fmt.Println("=== Multi-Locale Support ===")
	for _, locale := range locales {
		agents, err := store.GetAgentIndex(locale)
		if err != nil {
			log.Printf("Failed to fetch %s: %v", locale, err)
			continue
		}

		fmt.Printf("\n%s: %d agents\n", locale, len(agents))
		
		// Show first agent in each locale
		if len(agents) > 0 {
			fmt.Printf("  First agent: %s\n", agents[0].Meta.Title)
		}
	}
}

// ExampleCaching demonstrates cache behavior
func ExampleCaching() {
	store := NewAssistantStore("")

	fmt.Println("=== Cache Demonstration ===")
	
	// First fetch (from network)
	fmt.Println("First fetch (from network)...")
	agents1, err := store.GetAgentIndex("en-US")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Fetched %d agents\n", len(agents1))

	// Second fetch (from cache)
	fmt.Println("\nSecond fetch (from cache)...")
	agents2, err := store.GetAgentIndex("en-US")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Fetched %d agents (should be instant)\n", len(agents2))

	// Clear cache
	store.cache.Clear()
	fmt.Println("\nCache cleared")

	// Third fetch (from network again)
	fmt.Println("Third fetch (from network again)...")
	agents3, err := store.GetAgentIndex("en-US")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Fetched %d agents\n", len(agents3))
}

// ExampleCustomURL demonstrates using custom base URL
func ExampleCustomURL() {
	// You can use custom CDN or self-hosted registry
	customURL := "https://your-custom-cdn.com/agents"
	store := NewAssistantStore(customURL)

	fmt.Println("=== Custom URL ===")
	fmt.Printf("Base URL: %s\n", store.baseURL)
	fmt.Printf("Index URL: %s\n", store.GetAgentIndexURL("en-US"))
	fmt.Printf("Agent URL: %s\n", store.GetAgentURL("web-development", "en-US"))
}

