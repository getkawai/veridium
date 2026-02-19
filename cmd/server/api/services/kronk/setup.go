package kronk

import (
	"context"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
)

// SetupCommand returns the setup command for CLI integration
func SetupCommand(args []string) error {
	// Initialize Sentry for error tracking
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://709dabacc882a777ef059392d056e3da@o4510568649654272.ingest.us.sentry.io/4510568655290368",
		EnableTracing:    true,
		TracesSampleRate: 1.0,
		EnableLogs:       true,
		BeforeSendLog: func(log *sentry.Log) *sentry.Log {
			if log.Severity < sentry.LogSeverityWarning {
				return nil
			}
			return log
		},
	})
	if err != nil {
		// Silent fail - setup can continue without Sentry
	}
	defer sentry.Flush(2 * time.Second)

	// Parse flags
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--help", "-h":
			printSetupHelp()
			return nil
		}
	}

	// Use simple CLI setup mode
	ctx := context.Background()
	err = RunSetupSimple(ctx)
	if err != nil {
		// Capture error to Sentry before returning
		sentry.CaptureException(err)
	}
	return err
}

func printSetupHelp() {
	fmt.Println("Usage: ./server setup [OPTIONS]")
	fmt.Println()
	fmt.Println("Setup Kronk server (wallet, libraries, and models)")
	fmt.Println("All components are REQUIRED and cannot be skipped.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help, -h  Show this help")
	fmt.Println()
	fmt.Println("Hardware Requirements:")
	fmt.Println("  - RAM + VRAM: 24GB minimum")
	fmt.Println("  - Disk Space: 25GB free")
	fmt.Println("  - GPU: NVIDIA recommended")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ./server setup  # Interactive setup")
}
