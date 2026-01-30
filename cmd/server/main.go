// Package main is the entry point for the Kronk Model Server.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/kawai-network/veridium/cmd/server/api/services/kronk"
)

func main() {
	// Parse command line flags
	showHelp := flag.Bool("help", false, "Show configuration help")
	flag.Parse()

	// Run the kronk model server
	if err := kronk.Run(*showHelp); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
		os.Exit(1)
	}
}
