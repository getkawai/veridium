package main

import (
	"fmt"
	"os"

	"github.com/hybridgroup/yzma/pkg/llama"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <model_path> [text...]")
		fmt.Println("Example: go run main.go /path/to/model.gguf \"Hello world\"")
		return
	}

	modelPath := os.Args[1]
	
	// Load model
	mparams := llama.ModelDefaultParams()
	mparams.NGpuLayers = 35 // Enable GPU layers if available

	model, err := llama.ModelLoadFromFile(modelPath, mparams)
	if err != nil {
		fmt.Printf("Error loading model: %v\n", err)
		return
	}
	defer llama.ModelFree(model)

	vocab := llama.ModelGetVocab(model)
	
	fmt.Printf("=== YZMA TOKEN COUNTER ===\n")
	fmt.Printf("Model: %s\n", modelPath)
	fmt.Printf("Vocab Size: %d tokens\n", llama.VocabNTokens(vocab))
	fmt.Printf("\n")

	// If no text provided, show examples
	if len(os.Args) == 2 {
		showTokenExamples(vocab)
		return
	}

	// Count tokens for provided text
	text := os.Args[2]
	countTokens(vocab, text, "Input Text")
}

func countTokens(vocab llama.Vocab, text, label string) {
	// Tokenize with special tokens
	tokensWithSpecial := llama.Tokenize(vocab, text, true, true)
	
	// Tokenize without special tokens  
	tokensWithoutSpecial := llama.Tokenize(vocab, text, false, false)
	
	fmt.Printf("=== %s ===\n", label)
	fmt.Printf("Text: \"%s\"\n", text)
	fmt.Printf("With special tokens: %d tokens\n", len(tokensWithSpecial))
	fmt.Printf("Without special tokens: %d tokens\n", len(tokensWithoutSpecial))
	
	if len(tokensWithSpecial) <= 10 {
		fmt.Printf("Tokens: ")
		for i, token := range tokensWithSpecial {
			buf := make([]byte, 36)
			n := llama.TokenToPiece(vocab, token, buf, 0, true)
			tokenText := string(buf[:n])
			if i > 0 {
				fmt.Printf(", ")
			}
			fmt.Printf("'%s'", tokenText)
		}
		fmt.Println()
	}
	fmt.Println()
}

func showTokenExamples(vocab llama.Vocab) {
	examples := []string{
		"Hello world",
		"What is the capital of France?",
		"Artificial intelligence is transforming the world through machine learning algorithms.",
		"OpenAI's GPT models use transformers to understand and generate human-like text.",
		"The quick brown fox jumps over the lazy dog. This sentence contains every letter of the alphabet.",
		"",
		"a",
		"This is a very long text that contains multiple sentences. It should be tokenized properly by the language model. The tokenizer will break this down into individual tokens based on the vocabulary. Each token represents a piece of text that the model can understand and process efficiently.",
	}

	fmt.Println("=== TOKEN COUNTING EXAMPLES ===\n")
	
	for i, text := range examples {
		if text == "" {
			fmt.Printf("Example %d: [Empty String]\n", i+1)
		} else {
			fmt.Printf("Example %d: \"%s\"\n", i+1, text)
		}
		
		tokensWithSpecial := llama.Tokenize(vocab, text, true, true)
		tokensWithoutSpecial := llama.Tokenize(vocab, text, false, false)
		
		// Count characters
		charCount := len(text)
		
		fmt.Printf("  Characters: %d\n", charCount)
		fmt.Printf("  Tokens (with special): %d\n", len(tokensWithSpecial))
		fmt.Printf("  Tokens (without special): %d\n", len(tokensWithoutSpecial))
		fmt.Printf("  Characters per token (with special): %.2f\n", float64(charCount)/float64(len(tokensWithSpecial)))
		fmt.Println()
	}
	
	fmt.Println("=== SPECIAL TOKENS ===")
	fmt.Printf("BOS (Beginning of Sentence): %d\n", llama.VocabBOS(vocab))
	fmt.Printf("EOS (End of Sentence): %d\n", llama.VocabEOS(vocab))
	fmt.Printf("EOT (End of Turn): %d\n", llama.VocabEOT(vocab))
	fmt.Printf("SEP (Separator): %d\n", llama.VocabSEP(vocab))
	fmt.Printf("NL (New Line): %d\n", llama.VocabNL(vocab))
}
