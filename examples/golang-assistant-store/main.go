package main

import (
	"fmt"
	"log"
)

func main() {
	// Initialize AssistantStore
	store := NewAssistantStore("")

	// Fetch assistant list for English locale
	agents, err := store.GetAgentIndex("en-US")
	if err != nil {
		log.Fatalf("Failed to fetch agents: %v", err)
	}

	fmt.Printf("Found %d agents\n\n", len(agents))

	// Display first 5 agents
	for i, agent := range agents {
		if i >= 5 {
			break
		}
		fmt.Printf("%d. %s (%s)\n", i+1, agent.Meta.Title, agent.Identifier)
		fmt.Printf("   Category: %s\n", agent.Category)
		fmt.Printf("   Author: %s\n", agent.Author)
		fmt.Printf("   Tags: %v\n\n", agent.Meta.Tags)
	}

	// Fetch detail for specific agent
	if len(agents) > 0 {
		identifier := agents[0].Identifier
		fmt.Printf("\nFetching detail for: %s\n", identifier)

		detail, err := store.GetAgent(identifier, "en-US")
		if err != nil {
			log.Printf("Failed to fetch agent detail: %v", err)
		} else {
			fmt.Printf("Title: %s\n", detail.Meta.Title)
			fmt.Printf("Description: %s\n", detail.Meta.Description)
			if len(detail.SystemRole) > 100 {
				fmt.Printf("System Role: %s...\n", detail.SystemRole[:100])
			} else if len(detail.SystemRole) > 0 {
				fmt.Printf("System Role: %s\n", detail.SystemRole)
			} else {
				fmt.Printf("System Role: (empty)\n")
			}
		}
	}
}
