package kronk

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kawai-network/veridium/internal/paths"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/pkg/stablediffusion"
	"github.com/kawai-network/veridium/pkg/stablediffusion/modeldownloader"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/veridium/pkg/tools/defaults"
	"github.com/kawai-network/veridium/pkg/tools/libs"
	"github.com/kawai-network/veridium/pkg/whisper/model"
)

// SetupResult contains the result of setup
type SetupResult struct {
	WalletCreated   bool
	WalletAddress   string
	LibraryReady    bool
	WhisperReady    bool
	StableDiffReady bool
	Errors          []error
}

// RunSetup performs full setup (wallet, library, models)
func RunSetup() (*SetupResult, error) {
	result := &SetupResult{
		Errors: make([]error, 0),
	}

	ctx := context.Background()

	printBanner()
	fmt.Println("🔧 Starting Kronk Server Setup...")
	fmt.Println()

	// 1. Setup Wallet (REQUIRED)
	if err := setupWallet(ctx, result); err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("wallet setup failed: %w", err))
		printError("Wallet setup failed: " + err.Error())
		printSetupSummary(result)
		return result, fmt.Errorf("wallet setup is required")
	}

	// 2. Setup Libraries (REQUIRED)
	if err := setupLibraries(ctx, result); err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("library setup failed: %w", err))
		printError("Library setup failed: " + err.Error())
		printSetupSummary(result)
		return result, fmt.Errorf("library setup is required")
	}

	// 3. Setup Models (REQUIRED)
	if err := setupModels(ctx, result); err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("model setup failed: %w", err))
		printError("Model setup failed: " + err.Error())
		printSetupSummary(result)
		return result, fmt.Errorf("model setup is required")
	}

	// Print summary
	printSetupSummary(result)

	return result, nil
}

// setupWallet handles wallet creation/import
func setupWallet(ctx context.Context, result *SetupResult) error {
	printInfo("Setting up wallet...")

	// Initialize KV Store
	kv, kvErr := store.NewMultiNamespaceKVStore()
	if kvErr != nil {
		return fmt.Errorf("failed to connect to KV: %w", kvErr)
	}

	wallet := services.NewWalletService("", kv)

	// Check if wallet already exists
	if wallet.HasWallet() {
		wallets := wallet.GetWallets()
		if len(wallets) > 0 {
			printInfo("Wallet found!")

			if len(wallets) > 1 {
				// Multiple wallets - let user choose
				fmt.Println("\nAvailable wallets:")
				for i, w := range wallets {
					active := ""
					if w.IsActive {
						active = " (active)"
					}
					fmt.Printf("  %d. %s - %s%s\n", i+1, w.Description, w.Address[:10]+"...", active)
				}

				choice, err := promptChoice("\nSelect wallet:", func() []string {
					options := make([]string, len(wallets))
					for i, w := range wallets {
						options[i] = fmt.Sprintf("%s (%s...)", w.Description, w.Address[:10])
					}
					return options
				}())
				if err != nil {
					return fmt.Errorf("failed to select wallet: %w", err)
				}

				selectedWallet := wallets[choice].Address
				password, err := promptPassword("Enter password: ")
				if err != nil {
					return fmt.Errorf("failed to read password: %w", err)
				}

				walletAddress, err := wallet.SwitchWallet(selectedWallet, password)
				if err != nil {
					return fmt.Errorf("failed to switch wallet: %w", err)
				}

				result.WalletCreated = false
				result.WalletAddress = walletAddress
				printSuccess(fmt.Sprintf("Wallet unlocked: %s", walletAddress))
				return nil
			}

			// Single wallet - unlock it
			password, err := promptPassword("Enter password to unlock: ")
			if err != nil {
				return fmt.Errorf("failed to read password: %w", err)
			}

			walletAddress, err := wallet.UnlockWallet(password)
			if err != nil {
				return fmt.Errorf("invalid password: %w", err)
			}

			result.WalletCreated = false
			result.WalletAddress = walletAddress
			printSuccess(fmt.Sprintf("Wallet unlocked: %s", walletAddress))
			return nil
		}
	}

	// No wallet exists - create new one
	printInfo("No wallet found. Let's create one!")

	var password, mnemonic string
	var err error

	// Interactive wallet setup
	choice, err := promptChoice("Choose setup method:", []string{
		"Generate new mnemonic (recommended)",
		"Import existing mnemonic",
	})
	if err != nil {
		return fmt.Errorf("failed to get user choice: %w", err)
	}

	// Get password
	password, err = promptPassword("Enter password (min 8 characters): ")
	if err != nil {
		return fmt.Errorf("failed to read password: %w", err)
	}
	if err := validatePassword(password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	confirmPassword, err := promptPassword("Confirm password: ")
	if err != nil {
		return fmt.Errorf("failed to read password confirmation: %w", err)
	}
	if password != confirmPassword {
		return fmt.Errorf("passwords do not match")
	}

	if choice == 0 {
		// Generate new mnemonic
		mnemonic, err = wallet.GenerateMnemonic()
		if err != nil {
			return fmt.Errorf("failed to generate mnemonic: %w", err)
		}
		printMnemonic(mnemonic)

		if !promptYesNo("Have you written down your mnemonic?") {
			return fmt.Errorf("please write down your mnemonic before continuing")
		}
	} else {
		// Import existing mnemonic
		printInfo("Enter your 12 or 24 word mnemonic phrase")
		mnemonic, err = promptPassword("Mnemonic (hidden): ")
		if err != nil {
			return fmt.Errorf("failed to read mnemonic: %w", err)
		}
		// CRITICAL: Trim whitespace to ensure consistent wallet generation
		// Copy-paste can introduce trailing spaces/newlines that would generate
		// a different wallet address than the standard mnemonic
		mnemonic = strings.Join(strings.Fields(mnemonic), " ")
		if err := validateMnemonic(mnemonic); err != nil {
			return fmt.Errorf("invalid mnemonic: %w", err)
		}
	}

	// Get wallet name
	walletName, err := promptInput("Wallet name (e.g. My Contributor Wallet): ")
	if err != nil {
		return fmt.Errorf("failed to read wallet name: %w", err)
	}
	if walletName == "" {
		walletName = "Kronk Contributor"
	}

	address, err := wallet.CreateWallet(password, mnemonic, walletName)
	if err != nil {
		return fmt.Errorf("failed to create wallet: %w", err)
	}

	result.WalletCreated = true
	result.WalletAddress = address
	printSuccess(fmt.Sprintf("Wallet created: %s", address))

	return nil
}

// setupLibraries downloads required libraries
func setupLibraries(ctx context.Context, result *SetupResult) error {
	printInfo("Setting up libraries...")

	// Setup llama.cpp
	fmt.Println("  📦 Downloading llama.cpp library...")

	// Auto-detect platform
	arch, err := defaults.Arch("")
	if err != nil {
		return err
	}

	opSys, err := defaults.OS("")
	if err != nil {
		return err
	}

	processor, err := defaults.Processor("")
	if err != nil {
		return err
	}

	libMgr, err := libs.New(
		libs.WithBasePath(paths.NodeLibraries()),
		libs.WithArch(arch),
		libs.WithOS(opSys),
		libs.WithProcessor(processor),
		libs.WithAllowUpgrade(true),
		libs.WithVersion(defaults.LibVersion("")),
	)
	if err != nil {
		return fmt.Errorf("unable to create libs api: %w", err)
	}

	downloadCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()

	version, err := libMgr.Download(downloadCtx, func(ctx context.Context, msg string, args ...any) {
		fmt.Printf("    %s\n", fmt.Sprintf(msg, args...))
	})
	if err != nil {
		return fmt.Errorf("failed to download llama.cpp: %w", err)
	}

	fmt.Printf("  ✅ llama.cpp installed: %s\n", version.Version)

	// Setup Stable Diffusion library
	fmt.Println("  📦 Setting up Stable Diffusion library...")
	if err := stablediffusion.EnsureLibrary(); err != nil {
		return fmt.Errorf("failed to setup SD library: %w", err)
	}
	fmt.Println("  ✅ Stable Diffusion library ready")

	result.LibraryReady = true
	return nil
}

// setupModels downloads required models
func setupModels(ctx context.Context, result *SetupResult) error {
	printInfo("Setting up models...")

	// Setup Whisper model
	fmt.Println("  🎙️  Setting up Whisper model...")
	whisperModelsDir := filepath.Join(paths.Node(), "whisper-models")

	if err := os.MkdirAll(whisperModelsDir, 0755); err != nil {
		return fmt.Errorf("failed to create whisper models directory: %w", err)
	}

	existingModels, _ := model.ListDownloadedModels(whisperModelsDir)
	if len(existingModels) > 0 {
		fmt.Printf("  ℹ️  Whisper model already exists: %s\n", existingModels[0])
	} else {
		fmt.Println("  📥 Downloading Whisper base model (~141MB)...")
		if err := model.DownloadModel("base", whisperModelsDir, func(downloaded, total int64) {
			if total > 0 {
				percent := float64(downloaded) / float64(total) * 100
				fmt.Printf("\r    Progress: %.1f%%", percent)
			}
		}); err != nil {
			fmt.Println()
			return fmt.Errorf("failed to download whisper model: %w", err)
		}
		fmt.Println("\n  ✅ Whisper model downloaded")
	}

	result.WhisperReady = true

	// Setup Stable Diffusion model
	fmt.Println("  🎨 Setting up Stable Diffusion model...")
	modelsPath := filepath.Join(paths.Node(), "models")
	downloader := modeldownloader.New(modelsPath)

	modelFile, err := downloader.DiscoverModel()
	if err != nil {
		return fmt.Errorf("error discovering SD models: %w", err)
	}

	if modelFile != "" {
		fmt.Printf("  ℹ️  SD model already exists: %s\n", filepath.Base(modelFile))
		result.StableDiffReady = true
	} else {
		fmt.Println("  📥 Downloading Stable Diffusion model (~4GB)...")
		fmt.Println("     This may take a while depending on your connection...")

		modelFile, err = downloader.DownloadModelSimple(ctx, modeldownloader.DefaultModelURL)
		if err != nil {
			return fmt.Errorf("failed to download SD model: %w", err)
		}
		fmt.Printf("  ✅ SD model downloaded: %s\n", filepath.Base(modelFile))
		result.StableDiffReady = true
	}

	return nil
}

// printSetupSummary prints the setup summary
func printSetupSummary(result *SetupResult) {
	fmt.Println()
	fmt.Println("╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                   Setup Summary                           ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()

	if result.WalletCreated {
		fmt.Printf("  ✅ Wallet created: %s\n", result.WalletAddress)
	} else if result.WalletAddress != "" {
		fmt.Printf("  ℹ️  Wallet exists: %s\n", result.WalletAddress)
	} else {
		fmt.Println("  ❌ Wallet not configured")
	}

	if result.LibraryReady {
		fmt.Println("  ✅ Libraries ready (llama.cpp, SD)")
	} else {
		fmt.Println("  ❌ Libraries not ready")
	}

	if result.WhisperReady {
		fmt.Println("  ✅ Whisper model ready")
	} else {
		fmt.Println("  ❌ Whisper model not ready")
	}

	if result.StableDiffReady {
		fmt.Println("  ✅ Stable Diffusion model ready")
	} else {
		fmt.Println("  ❌ Stable Diffusion model not ready")
	}

	if len(result.Errors) > 0 {
		fmt.Println()
		fmt.Println("  ⚠️  Errors encountered:")
		for _, err := range result.Errors {
			fmt.Printf("     - %v\n", err)
		}
	}

	fmt.Println()
	fmt.Println("You can now start the server with: ./server")
	fmt.Println()
}

// SetupCommand returns the setup command for CLI integration
func SetupCommand(args []string) error {
	// Parse flags (only --help for now)
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--help", "-h":
			printSetupHelp()
			return nil
		}
	}

	_, err := RunSetup()
	return err
}

func printSetupHelp() {
	fmt.Println("Usage: ./server setup [OPTIONS]")
	fmt.Println()
	fmt.Println("Setup Kronk server (wallet, libraries, and models)")
	fmt.Println("All components are REQUIRED and cannot be skipped.")
	fmt.Println("Setup runs in interactive mode.")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --help, -h             Show this help")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  ./server setup         # Interactive setup")
}
