package kronk

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/kawai-network/veridium/internal/paths"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/pkg/hardware"
	sdmodels "github.com/kawai-network/veridium/pkg/stablediffusion/models"
	"github.com/kawai-network/x/store"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/term"
)

// SimpleSetup provides a simple, reliable setup process without complex UI
type SimpleSetup struct {
	walletService *services.WalletService
	walletExists  bool
	input         *InputReader
}

// InputReader handles secure CLI input
type InputReader struct {
	scanner *bufio.Scanner
}

// NewInputReader creates a new InputReader
func NewInputReader() *InputReader {
	return &InputReader{
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// ReadLine reads a line of input with a prompt
func (r *InputReader) ReadLine(prompt string) (string, error) {
	fmt.Print(prompt)
	if r.scanner.Scan() {
		return r.scanner.Text(), nil
	}
	return "", r.scanner.Err()
}

// ReadInt reads an integer with validation and retry
func (r *InputReader) ReadInt(prompt string, min, max int) (int, error) {
	fmt.Print(prompt)
	line, err := r.ReadLine("")
	if err != nil {
		return 0, err
	}
	var val int
	_, err = fmt.Sscanf(line, "%d", &val)
	if err != nil || val < min || val > max {
		return 0, fmt.Errorf("invalid input: must be between %d and %d", min, max)
	}
	return val, nil
}

// ReadPassword reads a password with hidden input and confirmation
func (r *InputReader) ReadPassword(prompt string) (string, error) {
	for {
		fmt.Print(prompt)
		password, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", fmt.Errorf("failed to read password: %w", err)
		}
		fmt.Println() // newline after password input

		passStr := string(password)
		if len(passStr) < 8 {
			fmt.Println("Password must be at least 8 characters. Please try again.")
			continue
		}

		fmt.Print("Confirm password: ")
		confirm, err := term.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			return "", fmt.Errorf("failed to read password confirmation: %w", err)
		}
		fmt.Println() // newline after confirmation

		if passStr != string(confirm) {
			fmt.Println("Passwords do not match. Please try again.")
			continue
		}

		return passStr, nil
	}
}

// ReadSecureLine reads sensitive input with hidden display
func (r *InputReader) ReadSecureLine(prompt string) (string, error) {
	fmt.Print(prompt)
	data, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		return "", err
	}
	fmt.Println() // newline after input
	return string(data), nil
}

// NewSimpleSetup creates a new SimpleSetup instance
func NewSimpleSetup() (*SimpleSetup, error) {
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to KV: %w", err)
	}

	walletService := services.NewWalletService("", kv)

	return &SimpleSetup{
		walletService: walletService,
		walletExists:  walletService.HasWallet(),
		input:         NewInputReader(),
	}, nil
}

// createProgressBar creates a progress bar and returns a new DownloadService instance
// This avoids race conditions by creating separate service instances for each download
func (s *SimpleSetup) createProgressBar(description string) *DownloadService {
	newBar := func() *progressbar.ProgressBar {
		return progressbar.NewOptions64(
			-1,
			progressbar.OptionSetDescription(description),
			progressbar.OptionSetWidth(40),
			progressbar.OptionShowBytes(true),
			progressbar.OptionSetPredictTime(true),
			progressbar.OptionShowCount(),
			progressbar.OptionOnCompletion(func() {
				fmt.Println()
			}),
		)
	}

	bar := newBar()
	var lastCompleted int64
	var lastTotal int64
	var reachedComplete bool

	return NewDownloadService(
		WithBasePath(paths.Base()),
		WithModelsPath(paths.Models()),
		WithMaxRetries(3),
		WithRetryDelay(5*time.Second),
		WithTimeout(2*time.Hour),
		WithProgressCallback(func(completed, total int64, percent float64, mbps float64) {
			// progressbar/v3 can't be reused after hitting 100%, so create a new bar
			// when a subsequent file starts (progress counter drops/reset).
			if reachedComplete && (completed < lastCompleted || completed == 0) {
				bar = newBar()
				reachedComplete = false
			}

			if total > 0 {
				bar.ChangeMax64(total)
			}
			bar.Set64(completed)

			lastCompleted = completed
			lastTotal = total
			if lastTotal > 0 && completed >= lastTotal {
				reachedComplete = true
			}
		}),
	)
}

// Run executes the setup process with simple CLI progress
func (s *SimpleSetup) Run(ctx context.Context) (*SetupResult, error) {
	fmt.Println("🌸 Kawai DeAI Network - Setup")
	fmt.Println("================================")
	fmt.Println()

	result := &SetupResult{
		Errors:   make([]error, 0),
		Warnings: make([]error, 0),
	}

	// Step 1: Wallet Setup
	fmt.Println("Step 1/6: Wallet Configuration")
	fmt.Println("-------------------------------")

	walletResult, err := s.setupWallet()
	if err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("wallet setup failed: %w", err))
		fmt.Printf("❌ Wallet setup failed: %v\n", err)
	} else {
		result.WalletCreated = walletResult.Created
		result.WalletAddress = walletResult.Address
		fmt.Printf("✅ Wallet: %s\n", walletResult.Address)
	}
	fmt.Println()

	// Step 2: Libraries
	fmt.Println("Step 2/6: Downloading Libraries")
	fmt.Println("--------------------------------")
	fmt.Println("Downloading llama.cpp, whisper.cpp, Stable Diffusion, and TTS libraries...")
	fmt.Println()

	libSvc := s.createProgressBar("Downloading Libraries")
	libResults, err := libSvc.DownloadLibraries(ctx)
	if err != nil {
		result.Errors = append(result.Errors, fmt.Errorf("libraries download failed: %w", err))
		fmt.Printf("❌ Libraries download failed: %v\n", err)
	} else {
		allSuccess := true
		for i, libResult := range libResults {
			libNames := []string{"llama.cpp", "whisper.cpp", "Stable Diffusion", "TTS"}
			libName := "Unknown"
			if i < len(libNames) {
				libName = libNames[i]
			}

			if libResult.Success {
				fmt.Printf("✅ %s ready\n", libName)
			} else {
				fmt.Printf("❌ %s failed: %v\n", libName, libResult.Error)
				allSuccess = false
				result.Errors = append(result.Errors, libResult.Error)
			}
		}

		if allSuccess {
			result.LibraryReady = true
			fmt.Println()
			fmt.Println("✅ All libraries ready!")
		}
	}
	fmt.Println()

	// Step 3: Whisper Model
	fmt.Println("Step 3/6: Whisper Model (Speech-to-Text)")
	fmt.Println("-----------------------------------------")

	whisperSvc := s.createProgressBar("Downloading Whisper")
	whisperResult := whisperSvc.DownloadWhisperModel(ctx, "base")
	if whisperResult.Success {
		result.WhisperReady = true
		fmt.Printf("✅ Whisper model ready (%.2f MB)\n", float64(whisperResult.Bytes)/1024/1024)
	} else {
		result.Warnings = append(result.Warnings, fmt.Errorf("whisper: %w", whisperResult.Error))
		fmt.Printf("❌ Whisper failed: %v (optional)\n", whisperResult.Error)
	}
	fmt.Println()

	// Step 4: Stable Diffusion Model
	fmt.Println("Step 4/6: Stable Diffusion Model (Image Generation)")
	fmt.Println("---------------------------------------------------")

	hw := hardware.DetectHardwareSpecs()
	selectedModel := sdmodels.SelectOptimalModel(&sdmodels.HardwareSpecs{
		TotalRAM:     hw.TotalRAM,
		AvailableRAM: hw.AvailableRAM,
		CPU:          hw.CPU,
		CPUCores:     hw.CPUCores,
		GPUMemory:    hw.GPUMemory,
		GPUModel:     hw.GPUModel,
	})
	fmt.Printf("Hardware detected: RAM=%dGB, VRAM=%dGB, CPU cores=%d\n", hw.TotalRAM, hw.GPUMemory, hw.CPUCores)
	fmt.Printf("Selected SD model: %s (%s, ~%dMB)\n", selectedModel.Name, selectedModel.ModelType, selectedModel.Size)
	fmt.Println()

	sdSvc := s.createProgressBar("Downloading Stable Diffusion")
	sdResult := sdSvc.DownloadStableDiffusionModelSmart(ctx, selectedModel)
	if !sdResult.Success {
		fmt.Printf("⚠️  Selected model download failed (%s): %v\n", selectedModel.Name, sdResult.Error)
	}

	if sdResult.Success {
		result.StableDiffReady = true
		fmt.Printf("✅ Stable Diffusion model ready (%.2f MB)\n", float64(sdResult.Bytes)/1024/1024)
	} else {
		result.Errors = append(result.Errors, fmt.Errorf("stable diffusion: %w", sdResult.Error))
		fmt.Printf("❌ Stable Diffusion failed: %v\n", sdResult.Error)
	}
	fmt.Println()

	// Step 5: TTS Model
	fmt.Println("Step 5/6: TTS Model (Text-to-Speech)")
	fmt.Println("-------------------------------------")

	ttsSvc := s.createProgressBar("Downloading TTS")
	ttsResult := ttsSvc.DownloadTTSModel(ctx)
	if ttsResult.Success {
		result.TTSReady = true
		fmt.Printf("✅ TTS model ready (%.2f MB)\n", float64(ttsResult.Bytes)/1024/1024)
	} else {
		result.Warnings = append(result.Warnings, fmt.Errorf("TTS: %w", ttsResult.Error))
		fmt.Printf("❌ TTS failed: %v (optional)\n", ttsResult.Error)
	}
	fmt.Println()

	// Step 6: LLM Model
	fmt.Println("Step 6/6: LLM Model (Nemotron 3 Nano)")
	fmt.Println("--------------------------------------")
	fmt.Println("Downloading Nemotron 3 Nano (~18GB)")
	fmt.Println("This may take 10-60 minutes depending on your connection...")
	fmt.Println()

	// Check for context cancellation before starting long download
	select {
	case <-ctx.Done():
		return result, ctx.Err()
	default:
	}

	llmSvc := s.createProgressBar("Downloading LLM")
	llmResult := llmSvc.DownloadLLMModelDefault(ctx)
	if llmResult.Success {
		result.LLMReady = true
		fmt.Printf("✅ LLM model ready (%.2f GB)\n", float64(llmResult.Bytes)/1024/1024/1024)
	} else {
		result.Errors = append(result.Errors, fmt.Errorf("LLM: %w", llmResult.Error))
		fmt.Printf("❌ LLM failed: %v\n", llmResult.Error)
	}
	fmt.Println()

	// Summary
	fmt.Println("================================")
	fmt.Println("Setup Summary")
	fmt.Println("================================")
	if result.WalletAddress != "" {
		fmt.Printf("Wallet:           %s\n", result.WalletAddress)
	}
	fmt.Printf("Libraries:        %v\n", statusIcon(result.LibraryReady))
	fmt.Printf("Whisper:          %v\n", statusIcon(result.WhisperReady))
	fmt.Printf("Stable Diffusion: %v\n", statusIcon(result.StableDiffReady))
	fmt.Printf("TTS:              %v\n", statusIcon(result.TTSReady))
	fmt.Printf("LLM:              %v\n", statusIcon(result.LLMReady))

	if len(result.Errors) > 0 {
		fmt.Println()
		fmt.Println("Errors:")
		for _, err := range result.Errors {
			fmt.Printf("  • %v\n", err)
		}
	}

	if len(result.Warnings) > 0 {
		fmt.Println()
		fmt.Println("Warnings (optional components failed):")
		for _, err := range result.Warnings {
			fmt.Printf("  • %v\n", err)
		}
	}

	fmt.Println()
	fmt.Println("✅ Setup complete! Start the server with: ./server start")
	fmt.Println()

	return result, nil
}

// WalletResult represents wallet setup result
type WalletResult struct {
	Address string
	Created bool
}

// setupWallet handles wallet configuration with interactive prompts
func (s *SimpleSetup) setupWallet() (*WalletResult, error) {
	if s.walletExists {
		fmt.Println("🔐 Existing Wallet Detected")
		fmt.Println("----------------------------")
		fmt.Println("A wallet is already configured on this system.")
		fmt.Println()
		fmt.Println("What would you like to do?")
		fmt.Println("  1. Skip wallet setup (use existing wallet)")
		fmt.Println("  2. Replace existing wallet")
		fmt.Println()

		choice, err := s.input.ReadInt("Enter choice [1-2]: ", 1, 2)
		if err != nil {
			return nil, fmt.Errorf("failed to read choice: %w", err)
		}

		if choice == 1 {
			fmt.Println("Using existing wallet...")
			return &WalletResult{
				Address: "existing_wallet",
				Created: false,
			}, nil
		}
		fmt.Println()
	}

	fmt.Println("🔐 Wallet Setup")
	fmt.Println("---------------")
	fmt.Println("Choose your setup method:")
	fmt.Println("  1. Generate new mnemonic (recommended)")
	fmt.Println("  2. Import existing mnemonic")
	fmt.Println("  3. Import keystore JSON (MetaMask, etc.)")
	fmt.Println("  4. Import private key")
	fmt.Println()

	choice, err := s.input.ReadInt("Enter choice [1-4]: ", 1, 4)
	if err != nil {
		return nil, fmt.Errorf("failed to read choice: %w", err)
	}

	// Get password first (all methods need password)
	fmt.Println()
	password, err := s.input.ReadPassword("Password: ")
	if err != nil {
		return nil, fmt.Errorf("failed to get password: %w", err)
	}

	// Show password strength
	score, label, _ := calculatePasswordStrength(password)
	fmt.Printf("Password strength: %s (%d%%)\n", label, score)
	fmt.Println()

	walletName := "My Wallet"

	switch choice {
	case 1: // Generate new mnemonic
		return s.generateNewWallet(password, walletName)

	case 2: // Import existing mnemonic
		return s.importMnemonicWallet(password, walletName)

	case 3: // Import keystore JSON
		return s.importKeystoreWallet(password, walletName)

	case 4: // Import private key
		return s.importPrivateKeyWallet(password, walletName)

	default:
		return nil, fmt.Errorf("invalid choice: %d", choice)
	}
}

// generateNewWallet generates a new wallet with mnemonic
func (s *SimpleSetup) generateNewWallet(password, walletName string) (*WalletResult, error) {
	fmt.Println("Generating new mnemonic...")
	mnemonic, err := s.walletService.GenerateMnemonic()
	if err != nil {
		return nil, fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	// Auto-copy to clipboard
	copyErr := clipboard.WriteAll(mnemonic)

	fmt.Println()
	fmt.Println("⚠️  IMPORTANT: Save your mnemonic phrase!")
	fmt.Println("Write these words down on paper and store them securely!")
	fmt.Println()
	fmt.Println(mnemonic)
	fmt.Println()

	if copyErr == nil {
		fmt.Println("✅ Mnemonic automatically copied to clipboard!")
	} else {
		fmt.Println("⚠️  Failed to copy to clipboard, please copy manually!")
	}

	fmt.Println()
	fmt.Println("Anyone with these words can access your funds!")
	fmt.Println()
	fmt.Print("Press Enter after you've saved the mnemonic...")
	_, _ = s.input.ReadLine("")

	// Confirm mnemonic saved
	fmt.Println()
	confirm, err := s.input.ReadLine("Have you saved the mnemonic? Type 'yes' to confirm: ")
	if err != nil {
		return nil, fmt.Errorf("failed to read confirmation: %w", err)
	}
	if strings.ToLower(strings.TrimSpace(confirm)) != "yes" {
		return nil, fmt.Errorf("mnemonic not confirmed")
	}

	address, err := s.walletService.CreateWallet(password, mnemonic, walletName)
	if err != nil {
		return nil, fmt.Errorf("failed to create wallet: %w", err)
	}

	return &WalletResult{
		Address: address,
		Created: true,
	}, nil
}

// importMnemonicWallet imports a wallet from mnemonic phrase
func (s *SimpleSetup) importMnemonicWallet(password, walletName string) (*WalletResult, error) {
	fmt.Println()
	fmt.Println("Enter your 12 or 24 word mnemonic phrase:")
	fmt.Println("(Input will be hidden for security)")
	mnemonic, err := s.input.ReadSecureLine("> ")
	if err != nil {
		return nil, fmt.Errorf("failed to read mnemonic: %w", err)
	}

	// Validate mnemonic (basic check)
	words := strings.Fields(mnemonic)
	if len(words) != 12 && len(words) != 24 {
		return nil, fmt.Errorf("invalid mnemonic: must be 12 or 24 words, got %d", len(words))
	}

	address, err := s.walletService.CreateWallet(password, mnemonic, walletName)
	if err != nil {
		return nil, fmt.Errorf("failed to import mnemonic: %w", err)
	}

	return &WalletResult{
		Address: address,
		Created: true,
	}, nil
}

// importKeystoreWallet imports a wallet from keystore JSON
func (s *SimpleSetup) importKeystoreWallet(password, walletName string) (*WalletResult, error) {
	fmt.Println()
	fmt.Println("Choose import method:")
	fmt.Println("  1. Paste keystore JSON content")
	fmt.Println("  2. Enter file path to keystore")
	fmt.Println()

	choice, err := s.input.ReadInt("Enter choice [1-2]: ", 1, 2)
	if err != nil {
		return nil, fmt.Errorf("failed to read choice: %w", err)
	}

	var keystoreData string
	if choice == 1 {
		fmt.Println()
		fmt.Println("⚠️  Warning: Pasting keystore in terminal may be logged in shell history.")
		fmt.Println("Consider using file import (option 2) for better security.")
		fmt.Println()
		fmt.Println("Paste your keystore JSON content:")
		keystoreData, err = s.input.ReadLine("> ")
		if err != nil {
			return nil, fmt.Errorf("failed to read keystore: %w", err)
		}
	} else {
		fmt.Println()
		fmt.Println("Enter the full path to your keystore file:")
		path, err := s.input.ReadLine("> ")
		if err != nil {
			return nil, fmt.Errorf("failed to read path: %w", err)
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read keystore file: %w", err)
		}
		keystoreData = string(data)
	}

	// Validate keystore JSON
	valid, errMsg := validateKeystoreJSON(keystoreData)
	if !valid {
		return nil, fmt.Errorf("invalid keystore: %s", errMsg)
	}

	address, err := s.walletService.ImportKeystore(keystoreData, password, walletName)
	if err != nil {
		return nil, fmt.Errorf("failed to import keystore: %w", err)
	}

	return &WalletResult{
		Address: address,
		Created: true,
	}, nil
}

// importPrivateKeyWallet imports a wallet from private key
func (s *SimpleSetup) importPrivateKeyWallet(password, walletName string) (*WalletResult, error) {
	fmt.Println()
	fmt.Println("⚠️  Security Warning: Never share your private key!")
	fmt.Println("Enter your 64-character hex private key (input will be hidden):")
	fmt.Println()

	privateKey, err := s.input.ReadSecureLine("Private key: ")
	if err != nil {
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}

	// Clean and validate private key
	privateKey = strings.TrimPrefix(privateKey, "0x")
	privateKey = strings.TrimPrefix(privateKey, "0X")

	valid, errMsg := validatePrivateKey(privateKey)
	if !valid {
		return nil, fmt.Errorf("invalid private key: %s", errMsg)
	}

	address, err := s.walletService.ImportPrivateKey(privateKey, password, walletName)
	if err != nil {
		return nil, fmt.Errorf("failed to import private key: %w", err)
	}

	return &WalletResult{
		Address: address,
		Created: true,
	}, nil
}

func statusIcon(ok bool) string {
	if ok {
		return "✅"
	}
	return "❌"
}

// DownloadWithProgress downloads a single file with progress bar
func DownloadWithProgress(ctx context.Context, url, dest string) error {
	bar := progressbar.NewOptions64(
		-1,
		progressbar.OptionSetDescription("Downloading"),
		progressbar.OptionSetWidth(40),
		progressbar.OptionShowBytes(true),
		progressbar.OptionSetPredictTime(true),
		progressbar.OptionShowCount(),
		progressbar.OptionOnCompletion(func() {
			fmt.Println()
		}),
	)

	svc := NewDownloadService(
		WithProgressCallback(func(completed, total int64, percent float64, mbps float64) {
			if total > 0 {
				bar.ChangeMax64(total)
			}
			bar.Set64(completed)
		}),
	)

	result := svc.DownloadWithRetry(ctx, url, dest)

	if !result.Success {
		return result.Error
	}

	return nil
}

// EnsureDirectoryExists creates necessary directories
func EnsureDirectoryExists() error {
	dirs := []string{
		paths.Base(),
		paths.Models(),
		filepath.Join(paths.Base(), "wallet"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return nil
}

// RunSetupSimple is the main entry point for simple setup
func RunSetupSimple(ctx context.Context) error {
	// Ensure directories exist
	if err := EnsureDirectoryExists(); err != nil {
		return fmt.Errorf("failed to setup directories: %w", err)
	}

	// Create setup instance
	setup, err := NewSimpleSetup()
	if err != nil {
		return fmt.Errorf("failed to create setup: %w", err)
	}

	// Run setup
	_, err = setup.Run(ctx)
	return err
}
