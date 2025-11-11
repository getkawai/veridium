package main

import (
	"fmt"
	"log"
	"os"

	"github.com/kawai-network/veridium/pkg/yzma/llama"
)

func main() {
	fmt.Println("🧪 Testing with Qwen2.5-0.5B Model")
	fmt.Println("===================================")

	// Load library
	libPath := "/opt/homebrew/lib"
	if err := llama.Load(libPath); err != nil {
		log.Fatal("Failed to load library:", err)
	}

	// Suppress verbose logs
	llama.LogSet(llama.LogSilent())
	llama.Init()
	defer llama.BackendFree()

	fmt.Println("✅ Library initialized")

	// Load small model
	modelPath := os.Getenv("HOME") + "/.llama-cpp/models/qwen2.5-0.5b-instruct-q4_k_m.gguf"
	fmt.Printf("📥 Loading model: %s\n", modelPath)

	mParams := llama.ModelDefaultParams()
	model := llama.ModelLoadFromFile(modelPath, mParams)
	if model == 0 {
		log.Fatal("Failed to load model")
	}
	defer llama.ModelFree(model)

	vocab := llama.ModelGetVocab(model)
	fmt.Println("✅ Model loaded")

	// Create context
	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = 2048
	ctxParams.NBatch = 512
	ctxParams.NThreads = 4

	ctx := llama.InitFromModel(model, ctxParams)
	if ctx == 0 {
		log.Fatal("Failed to create context")
	}
	defer llama.Free(ctx)

	fmt.Println("✅ Context created")

	// Create sampler - greedy for deterministic output
	sampler := llama.SamplerChainInit(llama.SamplerChainDefaultParams())
	llama.SamplerChainAdd(sampler, llama.SamplerInitTopK(40))
	llama.SamplerChainAdd(sampler, llama.SamplerInitTopP(0.9, 1))
	llama.SamplerChainAdd(sampler, llama.SamplerInitTempExt(0.8, 0, 1.0))
	llama.SamplerChainAdd(sampler, llama.SamplerInitDist(llama.DefaultSeed))

	// Test simple generation
	fmt.Println("\n💬 Test 1: Simple math")
	testGeneration(model, ctx, vocab, sampler, "What is 2+2? Answer:")

	fmt.Println("\n💬 Test 2: With chat template")
	testWithTemplate(model, ctx, vocab, sampler, "What is the capital of France?")

	fmt.Println("\n✅ All tests completed!")
}

func testGeneration(model llama.Model, ctx llama.Context, vocab llama.Vocab, sampler llama.Sampler, prompt string) {
	fmt.Printf("Prompt: %q\n", prompt)
	fmt.Print("Output: ")

	// Tokenize
	tokens := llama.Tokenize(vocab, prompt, true, false)
	batch := llama.BatchGetOne(tokens)

	// Generate
	for i := 0; i < 50; i++ {
		llama.Decode(ctx, batch)
		token := llama.SamplerSample(sampler, ctx, -1)

		if llama.VocabIsEOG(vocab, token) {
			break
		}

		buf := make([]byte, 256)
		length := llama.TokenToPiece(vocab, token, buf, 0, true)
		fmt.Print(string(buf[:length]))

		batch = llama.BatchGetOne([]llama.Token{token})
	}
	fmt.Println()
}

func testWithTemplate(model llama.Model, ctx llama.Context, vocab llama.Vocab, sampler llama.Sampler, prompt string) {
	template := llama.ModelChatTemplate(model, "")
	if template == "" {
		template = "chatml"
	}

	messages := []llama.ChatMessage{
		llama.NewChatMessage("user", prompt),
	}

	buf := make([]byte, 8192)
	length := llama.ChatApplyTemplate(template, messages, true, buf)
	formattedPrompt := string(buf[:length])

	fmt.Printf("Prompt: %q\n", prompt)
	fmt.Printf("Template: %s\n", template[:50]+"...")
	fmt.Print("Output: ")

	// Tokenize
	tokens := llama.Tokenize(vocab, formattedPrompt, true, false)
	batch := llama.BatchGetOne(tokens)

	// Generate
	for i := 0; i < 100; i++ {
		llama.Decode(ctx, batch)
		token := llama.SamplerSample(sampler, ctx, -1)

		if llama.VocabIsEOG(vocab, token) {
			break
		}

		buf := make([]byte, 256)
		length := llama.TokenToPiece(vocab, token, buf, 0, true)
		fmt.Print(string(buf[:length]))

		batch = llama.BatchGetOne([]llama.Token{token})
	}
	fmt.Println()
}
