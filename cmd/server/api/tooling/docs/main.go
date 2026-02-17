package main

import (
	"fmt"
	"os"

	"github.com/kawai-network/veridium/cmd/server/api/tooling/docs/api"
	"github.com/kawai-network/veridium/cmd/server/api/tooling/docs/manual"
	"github.com/kawai-network/veridium/cmd/server/api/tooling/docs/sdk/examples"
	"github.com/kawai-network/veridium/cmd/server/api/tooling/docs/sdk/gofmt"
)

func main() {
	if err := run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run() error {
	// Generate Markdown for kawai-website/docs
	fmt.Println("Generating Markdown documentation for kawai-website/docs...")
	if err := api.RunMarkdown(); err != nil {
		return err
	}

	// Generate TSX components for BUI frontend (optional - may fail if directory doesn't exist)
	fmt.Println("Generating TSX components for BUI frontend...")
	if err := api.Run(); err != nil {
		fmt.Printf("Warning: TSX generation skipped (%v)\n", err)
	}

	if err := gofmt.Run(); err != nil {
		fmt.Printf("Warning: SDK docs skipped (%v)\n", err)
	}

	if err := examples.Run(); err != nil {
		fmt.Printf("Warning: Examples skipped (%v)\n", err)
	}

	if err := manual.Run(); err != nil {
		fmt.Printf("Warning: Manual skipped (%v)\n", err)
	}

	return nil
}
