package main

import (
	"fmt"
	"math"
	"os"
	"strings"

	"github.com/kawai-network/veridium/pkg/yzma/llama"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	if err := handleFlags(); err != nil {
		showUsage()
		return err
	}

	if err := llama.Load(*libPath); err != nil {
		fmt.Println("unable to load library", err.Error())
		os.Exit(1)
	}

	if !*verbose {
		llama.LogSet(llama.LogSilent())
	}

	llama.Init()
	defer llama.BackendFree()

	model, err := llama.ModelLoadFromFile(*modelFile, llama.ModelDefaultParams())
	if err != nil {
		return fmt.Errorf("unable to load model: %w", err)
	}
	defer llama.ModelFree(model)

	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = uint32(*contextSize)
	ctxParams.NBatch = uint32(*batchSize)
	ctxParams.PoolingType = poolingType
	ctxParams.Embeddings = 1

	lctx, err := llama.InitFromModel(model, ctxParams)
	if err != nil {
		return fmt.Errorf("unable to init context: %w", err)
	}
	defer llama.Free(lctx)

	// tokenize prompt
	vocab := llama.ModelGetVocab(model)
	tokens := llama.Tokenize(vocab, *prompt, true, true)

	// create batch and decode
	batch := llama.BatchGetOne(tokens)
	llama.Decode(lctx, batch)

	// get embeddings
	nEmbd := llama.ModelNEmbd(model)
	vec, err := llama.GetEmbeddingsSeq(lctx, 0, nEmbd)
	if err != nil {
		return fmt.Errorf("unable to get embeddings: %w", err)
	}

	// normalize embeddings
	var sum float64
	for _, v := range vec {
		sum += float64(v * v)
	}
	sum = math.Sqrt(sum)
	norm := float32(1.0 / sum)

	var b strings.Builder
	for i, v := range vec {
		if i > 0 {
			b.WriteString(" ")
		}
		b.WriteString(fmt.Sprintf("%f", v*norm))
	}
	fmt.Println(b.String())

	return nil
}
