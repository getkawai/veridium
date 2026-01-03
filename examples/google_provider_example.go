package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"reflect"

	"github.com/kawai-network/veridium/pkg/fantasy"
	"github.com/kawai-network/veridium/pkg/fantasy/providers/google"
	"github.com/kawai-network/veridium/pkg/fantasy/schema"
)

func main() {
	// Ambil API key dari environment variable
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Fatal("GEMINI_API_KEY environment variable is required")
	}

	ctx := context.Background()

	// Contoh 1: Basic Text Generation
	fmt.Println("=== Example 1: Basic Text Generation ===")
	basicExample(ctx, apiKey)

	// Contoh 2: Streaming
	fmt.Println("\n=== Example 2: Streaming ===")
	streamingExample(ctx, apiKey)

	// Contoh 3: Tool Calling
	fmt.Println("\n=== Example 3: Tool Calling ===")
	toolCallingExample(ctx, apiKey)

	// Contoh 4: Object Generation
	fmt.Println("\n=== Example 4: Object Generation ===")
	objectGenerationExample(ctx, apiKey)
}

func basicExample(ctx context.Context, apiKey string) {
	provider, err := google.New(
		google.WithGeminiAPIKey(apiKey),
	)
	if err != nil {
		log.Fatal(err)
	}

	model, err := provider.LanguageModel(ctx, "gemini-1.5-flash")
	if err != nil {
		log.Fatal(err)
	}

	response, err := model.Generate(ctx, fantasy.Call{
		Prompt: fantasy.Prompt{
			{
				Role: fantasy.MessageRoleSystem,
				Content: []fantasy.MessagePart{
					fantasy.TextPart{Text: "You are a helpful assistant."},
				},
			},
			{
				Role: fantasy.MessageRoleUser,
				Content: []fantasy.MessagePart{
					fantasy.TextPart{Text: "What is the capital of Indonesia?"},
				},
			},
		},
		Temperature: fantasy.Opt(0.7),
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Response: %s\n", response.Content.Text())
	fmt.Printf("Usage: %+v\n", response.Usage)
}

func streamingExample(ctx context.Context, apiKey string) {
	provider, err := google.New(
		google.WithGeminiAPIKey(apiKey),
	)
	if err != nil {
		log.Fatal(err)
	}

	model, err := provider.LanguageModel(ctx, "gemini-1.5-flash")
	if err != nil {
		log.Fatal(err)
	}

	stream, err := model.Stream(ctx, fantasy.Call{
		Prompt: fantasy.Prompt{
			{
				Role: fantasy.MessageRoleUser,
				Content: []fantasy.MessagePart{
					fantasy.TextPart{Text: "Count from 1 to 5 slowly."},
				},
			},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("Streaming response: ")
	for part := range stream {
		switch part.Type {
		case fantasy.StreamPartTypeTextDelta:
			fmt.Print(part.Delta)
		case fantasy.StreamPartTypeFinish:
			fmt.Printf("\nUsage: %+v\n", part.Usage)
		case fantasy.StreamPartTypeError:
			log.Printf("Error: %v", part.Error)
		}
	}
	fmt.Println()
}

func toolCallingExample(ctx context.Context, apiKey string) {
	provider, err := google.New(
		google.WithGeminiAPIKey(apiKey),
	)
	if err != nil {
		log.Fatal(err)
	}

	model, err := provider.LanguageModel(ctx, "gemini-1.5-flash")
	if err != nil {
		log.Fatal(err)
	}

	tools := []fantasy.Tool{
		fantasy.FunctionTool{
			Name:        "get_weather",
			Description: "Get current weather for a location",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"location": map[string]any{
						"type":        "string",
						"description": "City name",
					},
					"unit": map[string]any{
						"type":        "string",
						"description": "Temperature unit (celsius or fahrenheit)",
						"enum":        []string{"celsius", "fahrenheit"},
					},
				},
				"required": []string{"location"},
			},
		},
	}

	response, err := model.Generate(ctx, fantasy.Call{
		Prompt: fantasy.Prompt{
			{
				Role: fantasy.MessageRoleUser,
				Content: []fantasy.MessagePart{
					fantasy.TextPart{Text: "What's the weather like in Jakarta?"},
				},
			},
		},
		Tools: tools,
	})
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Finish Reason: %s\n", response.FinishReason)
	for _, content := range response.Content {
		if content.GetType() == fantasy.ContentTypeToolCall {
			toolCall := content.(fantasy.ToolCallContent)
			fmt.Printf("Tool Call: %s\n", toolCall.ToolName)
			fmt.Printf("Arguments: %s\n", toolCall.Input)
		}
	}
}

func objectGenerationExample(ctx context.Context, apiKey string) {
	provider, err := google.New(
		google.WithGeminiAPIKey(apiKey),
	)
	if err != nil {
		log.Fatal(err)
	}

	model, err := provider.LanguageModel(ctx, "gemini-1.5-flash")
	if err != nil {
		log.Fatal(err)
	}

	type Person struct {
		Name string `json:"name" jsonschema:"description=Person's full name"`
		Age  int    `json:"age" jsonschema:"description=Person's age in years"`
		City string `json:"city" jsonschema:"description=City where person lives"`
	}

	objResponse, err := model.GenerateObject(ctx, fantasy.ObjectCall{
		Prompt: fantasy.Prompt{
			{
				Role: fantasy.MessageRoleUser,
				Content: []fantasy.MessagePart{
					fantasy.TextPart{
						Text: "Extract information: Budi is 25 years old and lives in Jakarta.",
					},
				},
			},
		},
		Schema: schema.Generate(reflect.TypeOf(Person{})),
	})
	if err != nil {
		log.Fatal(err)
	}

	person := objResponse.Object.(*Person)
	fmt.Printf("Extracted Person:\n")
	fmt.Printf("  Name: %s\n", person.Name)
	fmt.Printf("  Age: %d\n", person.Age)
	fmt.Printf("  City: %s\n", person.City)
	fmt.Printf("Usage: %+v\n", objResponse.Usage)
}
