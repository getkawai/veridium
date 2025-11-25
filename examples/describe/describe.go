package main

import (
	"fmt"
	"path"
	"unsafe"

	"github.com/kawai-network/veridium/pkg/yzma/llama"
	"github.com/kawai-network/veridium/pkg/yzma/mtmd"
)

func describe(tmpFile string) {
	llama.Init()
	defer llama.BackendFree()

	mtmdCtxParams := mtmd.ContextParamsDefault()

	switch {
	case *verbose:
		fmt.Println("Using model", path.Join(*modelsDir, *modelFile))
	default:
		llama.LogSet(llama.LogSilent())
		mtmd.LogSet(llama.LogSilent())
	}

	model, err := llama.ModelLoadFromFile(path.Join(*modelsDir, *modelFile), llama.ModelDefaultParams())
	if err != nil {
		fmt.Println("unable to load model", err.Error())
		return
	}
	defer llama.ModelFree(model)

	ctxParams := llama.ContextDefaultParams()
	ctxParams.NCtx = 4096
	ctxParams.NBatch = 2048

	lctx, err := llama.InitFromModel(model, ctxParams)
	if err != nil {
		fmt.Println("unable to init context", err.Error())
		return
	}
	defer llama.Free(lctx)

	vocab := llama.ModelGetVocab(model)
	// Create sampler chain
	sampler := llama.SamplerChainInit(llama.SamplerChainDefaultParams())
	llama.SamplerChainAdd(sampler, llama.SamplerInitTopK(40))
	llama.SamplerChainAdd(sampler, llama.SamplerInitTopP(0.95, 1))
	llama.SamplerChainAdd(sampler, llama.SamplerInitTempExt(0.8, 0, 1.0))
	llama.SamplerChainAdd(sampler, llama.SamplerInitDist(llama.DefaultSeed))

	mtmdCtx, err := mtmd.InitFromFile(path.Join(*modelsDir, *projFile), model, mtmdCtxParams)
	if err != nil {
		fmt.Println("unable to init mtmd context", err.Error())
		return
	}
	defer mtmd.Free(mtmdCtx)

	template = llama.ModelChatTemplate(model, "")
	messages = []llama.ChatMessage{llama.NewChatMessage("user", *prompt+mtmd.DefaultMarker())}
	output := mtmd.InputChunksInit()
	input := mtmd.NewInputText(chatTemplate(true), true, true)

	bitmap := mtmd.BitmapInitFromFile(mtmdCtx, tmpFile)
	defer mtmd.BitmapFree(bitmap)

	mtmd.Tokenize(mtmdCtx, output, input, []mtmd.Bitmap{bitmap})

	var n llama.Pos
	mtmd.HelperEvalChunks(mtmdCtx, lctx, output, 0, 0, int32(ctxParams.NBatch), true, &n)

	var sz int32 = 1
	batch := llama.BatchInit(1, 0, 1)
	batch.NSeqId = &sz
	batch.NTokens = 1
	seqs := unsafe.SliceData([]llama.SeqId{0})
	batch.SeqId = &seqs

	fmt.Println()

	for i := 0; i < llama.MaxToken; i++ {
		token := llama.SamplerSample(sampler, lctx, -1)

		if llama.VocabIsEOG(vocab, token) {
			fmt.Println()
			break
		}

		buf := make([]byte, 128)
		l := llama.TokenToPiece(vocab, token, buf, 0, true)

		fmt.Print(string(buf[:l]))

		batch.Token = &token
		batch.Pos = &n

		llama.Decode(lctx, batch)
		n++
	}
}

func chatTemplate(add bool) string {
	buf := make([]byte, 1024)
	len := llama.ChatApplyTemplate(template, messages, add, buf)
	result := string(buf[:len])
	return result
}
