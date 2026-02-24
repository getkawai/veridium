// Package main is the entry point for the Kronk Model Server.
package main

import (
	"fmt"
	"os"

	"github.com/kawai-network/veridium/cmd/server/api/services/kronk"
	"github.com/kawai-network/veridium/internal/paths"
)

func main() {
	// Use local data directory in development, user path in production.
	if os.Getenv("VERIDIUM_DEV") == "1" {
		paths.SetDataDir("data")
	} else {
		paths.SetDataDir(paths.UserDataDir())
	}

	// Check for subcommands first
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "setup":
			// Run setup command
			if err := kronk.SetupCommand(os.Args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		case "start":
			// Run start command (start the server)
			if err := kronk.StartCommand(os.Args[2:]); err != nil {
				fmt.Fprintf(os.Stderr, "ERROR: %v\n", err)
				os.Exit(1)
			}
			os.Exit(0)
		case "version":
			fmt.Println("Kronk Model Server")
			os.Exit(0)
		case "help", "--help", "-h":
			printHelp()
			os.Exit(0)
		}
	}

	// Default: show help
	printHelp()
	os.Exit(1)
}

func printHelp() {
	fmt.Println("Kawai Node - Kronk Model Server")
	fmt.Println()
	fmt.Println("Usage: kawai-contributor [COMMAND] [OPTIONS]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  setup    Setup wallet, libraries, and models")
	fmt.Println("  start    Start the model server")
	fmt.Println("  version  Show version information")
	fmt.Println("  help     Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  kawai-contributor setup         # Run interactive setup")
	fmt.Println("  kawai-contributor start         # Start the server")
}
