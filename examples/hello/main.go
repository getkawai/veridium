package main

import (
	"fmt"
	"os"

	"github.com/kawai-network/veridium/pkg/yzma/llama"
)

var (
	modelFile            = "./models/SmolLM-135M.Q2_K.gguf"
	prompt               = "Are you ready to go?"
	libPath              = os.Getenv("YZMA_LIB")
	responseLength int32 = 18
)

func main() {
	llama.Load(libPath)
	llama.Init()

	model, err := llama.ModelLoadFromFile(modelFile, llama.ModelDefaultParams())
	if err != nil {
		panic(err)
	}
	lctx, err := llama.InitFromModel(model, llama.ContextDefaultParams())
	if err != nil {
		panic(err)
	}

	vocab := llama.ModelGetVocab(model)

	// get tokens from the prompt
	tokens := llama.Tokenize(vocab, prompt, true, false)

	batch := llama.BatchGetOne(tokens)

	sampler := llama.SamplerChainInit(llama.SamplerChainDefaultParams())
	llama.SamplerChainAdd(sampler, llama.SamplerInitGreedy())

	for pos := int32(0); pos < responseLength; pos += batch.NTokens {
		llama.Decode(lctx, batch)
		token := llama.SamplerSample(sampler, lctx, -1)

		if llama.VocabIsEOG(vocab, token) {
			fmt.Println()
			break
		}

		buf := make([]byte, 36)
		len := llama.TokenToPiece(vocab, token, buf, 0, true)

		fmt.Print(string(buf[:len]))

		batch = llama.BatchGetOne([]llama.Token{token})
	}

	fmt.Println()
}
