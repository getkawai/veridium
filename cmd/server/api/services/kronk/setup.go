package kronk

import (
	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/getsentry/sentry-go"
	"github.com/kawai-network/veridium/internal/paths"
)

// SetupResult contains the result of setup
type SetupResult struct {
	WalletCreated   bool
	WalletAddress   string
	LibraryReady    bool
	WhisperReady    bool
	StableDiffReady bool
	TTSReady        bool
	LLMReady        bool
	Errors          []error
}

// SetupCommand returns the setup command for CLI integration
func SetupCommand(args []string) error {
	f, err := tea.LogToFile(paths.ContributorLog(), "contributor-setup")
	if err != nil {
		// Silent fail - setup can continue without debug logging
	} else {
		defer f.Close()
	}

	// Initialize Sentry for error tracking (does not write to stdout)
	err = sentry.Init(sentry.ClientOptions{
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
		// We don't print to stdout to avoid interfering with TUI
	}
	defer sentry.Flush(2 * time.Second)

	// Parse flags
	skipHardwareCheck := false

	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--help", "-h":
			printSetupHelp()
			return nil
		case "--skip-hardware-check":
			// Only allow skip in development builds
			if tag == "develop" {
				skipHardwareCheck = true
			} else {
				return fmt.Errorf("--skip-hardware-check is only available in development builds")
			}
		}
	}

	// Always use TUI mode
	_, err = NewSetupTUI(skipHardwareCheck)
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
	fmt.Println("Setup runs in interactive TUI mode.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help, -h                Show this help")
	fmt.Println("  --skip-hardware-check     Skip hardware requirements check (testing only)")
	fmt.Println()
	fmt.Println("Hardware Requirements:")
	fmt.Println("  - RAM + VRAM: 24GB minimum")
	fmt.Println("  - Disk Space: 25GB free")
	fmt.Println("  - GPU: NVIDIA recommended")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ./server setup                        # Interactive TUI setup")
	fmt.Println("  ./server setup --skip-hardware-check  # Skip hardware check (testing)")
}
