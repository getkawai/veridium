package main

import (
	"fmt"
	"os"

	"github.com/kawai-network/veridium/pkg/yzma/llama"
	"github.com/kawai-network/veridium/pkg/yzma/mtmd"
)

var messages []llama.ChatMessage

func main() {
	if err := handleFlags(); err != nil {
		showUsage()
		os.Exit(0)
	}

	if err := llama.Load(libPath); err != nil {
		fmt.Println("unable to load library", err.Error())
		os.Exit(1)
	}
	if err := mtmd.Load(libPath); err != nil {
		fmt.Println("unable to load library", err.Error())
		os.Exit(1)
	}

	tmpFile, err := obtainFile(imageFile)
	if err != nil {
		fmt.Println("unable to download image", err.Error())
		os.Exit(1)
	}
	defer os.Remove(tmpFile)

	describe(tmpFile)
}
