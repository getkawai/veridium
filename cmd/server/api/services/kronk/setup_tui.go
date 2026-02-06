package kronk

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kawai-network/veridium/internal/paths"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/pkg/hardware"
	"github.com/kawai-network/veridium/pkg/stablediffusion"
	"github.com/kawai-network/veridium/pkg/stablediffusion/download"
	"github.com/kawai-network/veridium/pkg/stablediffusion/modeldownloader"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/veridium/pkg/tools/defaults"
	"github.com/kawai-network/veridium/pkg/tools/libs"
	"github.com/kawai-network/veridium/pkg/tools/models"
	"github.com/kawai-network/veridium/pkg/whisper/model"
)

// Styles for the TUI
type styles struct {
	Title        lipgloss.Style
	Subtitle     lipgloss.Style
	Container    lipgloss.Style
	Success      lipgloss.Style
	Error        lipgloss.Style
	Warning      lipgloss.Style
	Info         lipgloss.Style
	Mnemonic     lipgloss.Style
	Help         lipgloss.Style
	ActiveStep   lipgloss.Style
	InactiveStep lipgloss.Style
}

func newStyles() *styles {
	return &styles{
		Title: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF6B9D")).
			MarginTop(1).
			MarginBottom(1),
		Subtitle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#A0A0A0")).
			MarginBottom(1),
		Container: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FF6B9D")).
			Padding(2).
			Width(80),
		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF00")).
			Bold(true),
		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true),
		Warning: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFA500")),
		Info: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00BFFF")),
		Mnemonic: lipgloss.NewStyle().
			Background(lipgloss.Color("#333333")).
			Foreground(lipgloss.Color("#FFFFFF")).
			Padding(1).
			Margin(1).
			Bold(true),
		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")),
		ActiveStep: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF6B9D")).
			Bold(true),
		InactiveStep: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#666666")),
	}
}

// Step represents a setup step
type step int

const (
	stepWelcome step = iota
	stepHardwareCheck
	stepWalletReplaceChoice
	stepWalletChoice
	stepWalletPassword
	stepWalletConfirmPassword
	stepWalletMnemonic
	stepWalletName
	stepImportKeystoreChoice
	stepImportKeystoreJSON
	stepImportKeystorePath
	stepImportPrivateKey
	stepLibraryDownload
	stepModelDownload
	stepLLMDownload
	stepSummary
	stepError
)

// Model represents the TUI state
type setupTUIModel struct {
	styles *styles
	step   step

	// Wallet setup
	walletService       *services.WalletService
	walletChoice        int
	walletReplaceChoice int // 0: skip, 1: replace
	keystoreChoice      int // 0: paste JSON, 1: file path
	passwordInput       textinput.Model
	confirmInput        textinput.Model
	nameInput           textinput.Model
	mnemonicInput       textinput.Model
	keystoreInput       textinput.Model // For pasting JSON
	keystorePathInput   textinput.Model // For file path
	privateKeyInput     textinput.Model
	generatedMnemonic   string
	walletAddress       string
	keystorePassword    string // Separate password for keystore import
	mnemonicCopied      bool
	copyError           string

	// Progress
	spinner     spinner.Model
	progressBar progress.Model

	// Hardware
	hwSpecs           *hardware.HardwareSpecs
	hardwarePassed    bool
	skipHardwareCheck bool

	// Downloads
	libraryProgress float64
	whisperProgress float64
	sdProgress      float64
	llmProgress     float64

	// Results
	result             *SetupResult
	errors             []error
	walletExists       bool
	walletAlreadySetup bool

	// Dimensions
	width  int
	height int

	// Communication
	progressChan chan float64
}

type progressMsg float64

func waitForProgress(c chan float64) tea.Cmd {
	return func() tea.Msg {
		return progressMsg(<-c)
	}
}

// Messages
type (
	hardwareCheckMsg struct {
		specs  *hardware.HardwareSpecs
		passed bool
		err    error
	}
	libraryDownloadMsg struct {
		progress float64
		done     bool
		err      error
	}
	modelDownloadMsg struct {
		modelType string
		progress  float64
		done      bool
		err       error
	}
	walletCreatedMsg struct {
		address string
		err     error
	}
	walletUnlockedMsg struct {
		address string
		err     error
	}
	mnemonicGeneratedMsg struct {
		mnemonic string
		err      error
	}
	keystoreReadMsg struct {
		content string
		err     error
	}
)

// Init initializes the TUI
func (m setupTUIModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		waitForProgress(m.progressChan),
	)
}

// Update handles messages and updates the model
func (m setupTUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case progressMsg:
		if m.step == stepLibraryDownload {
			m.libraryProgress = float64(msg)
		} else if m.step == stepModelDownload {
			// Basic heuristic: update whatever progress is active (whisper or SD)
			if m.result.WhisperReady {
				m.sdProgress = float64(msg)
			} else {
				m.whisperProgress = float64(msg)
			}
		} else if m.step == stepLLMDownload {
			m.llmProgress = float64(msg)
		}
		return m, waitForProgress(m.progressChan)

	case tea.KeyMsg:
		if m.step == stepWalletMnemonic && msg.String() == "c" {
			if m.generatedMnemonic != "" {
				if err := clipboard.WriteAll(m.generatedMnemonic); err != nil {
					m.copyError = "Check your clipboard setup!"
					m.mnemonicCopied = false
				} else {
					m.mnemonicCopied = true
					m.copyError = ""
				}
				return m, nil
			}
		}

		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			return m.handleEnter()
		case tea.KeyEsc:
			return m.handleBack()
		case tea.KeyUp:
			if m.step == stepWalletReplaceChoice && m.walletReplaceChoice > 0 {
				m.walletReplaceChoice--
			} else if m.step == stepWalletChoice && m.walletChoice > 0 {
				m.walletChoice--
			} else if m.step == stepImportKeystoreChoice && m.keystoreChoice > 0 {
				m.keystoreChoice--
			}
			return m, nil
		case tea.KeyDown:
			if m.step == stepWalletReplaceChoice && m.walletReplaceChoice < 1 {
				m.walletReplaceChoice++
			} else if m.step == stepWalletChoice && m.walletChoice < 3 {
				m.walletChoice++
			} else if m.step == stepImportKeystoreChoice && m.keystoreChoice < 1 {
				m.keystoreChoice++
			}
			return m, nil
		}

	case hardwareCheckMsg:
		if msg.err != nil {
			m.errors = append(m.errors, msg.err)
			m.hardwarePassed = false
		} else {
			m.hwSpecs = msg.specs
			m.hardwarePassed = msg.passed
		}
		return m, nil

	case libraryDownloadMsg:
		if msg.err != nil {
			m.errors = append(m.errors, fmt.Errorf("library download failed: %w", msg.err))
			m.step = stepError
		} else if msg.done {
			m.result.LibraryReady = true
			m.step = stepModelDownload
			return m, m.downloadModelsCmd()
		} else {
			m.libraryProgress = msg.progress
		}
		return m, nil

	case modelDownloadMsg:
		if msg.err != nil {
			m.errors = append(m.errors, fmt.Errorf("%s model download failed: %w", msg.modelType, msg.err))
			m.step = stepError
		} else if msg.done {
			switch msg.modelType {
			case "whisper":
				m.result.WhisperReady = true
			case "sd":
				m.result.StableDiffReady = true
				m.step = stepLLMDownload
				return m, m.downloadLLMCmd()
			case "llm":
				m.result.LLMReady = true
				m.step = stepSummary
			}
		} else {
			switch msg.modelType {
			case "whisper":
				m.whisperProgress = msg.progress
			case "sd":
				m.sdProgress = msg.progress
			case "llm":
				m.llmProgress = msg.progress
			}
		}
		return m, nil

	case walletCreatedMsg:
		if msg.err != nil {
			m.errors = append(m.errors, msg.err)
		} else {
			m.walletAddress = msg.address
			m.result.WalletCreated = true
			m.result.WalletAddress = msg.address
			m.step = stepHardwareCheck
			return m, m.checkHardwareCmd()
		}
		return m, nil

	case walletUnlockedMsg:
		if msg.err != nil {
			m.errors = append(m.errors, msg.err)
		} else {
			m.walletAddress = msg.address
			m.result.WalletAddress = msg.address
			m.step = stepHardwareCheck
			return m, m.checkHardwareCmd()
		}
		return m, nil

	case mnemonicGeneratedMsg:
		if msg.err != nil {
			m.errors = append(m.errors, fmt.Errorf("failed to generate mnemonic: %w", msg.err))
			m.step = stepError
		} else {
			m.generatedMnemonic = msg.mnemonic
		}
		return m, nil

	case keystoreReadMsg:
		if msg.err != nil {
			m.errors = append(m.errors, msg.err)
			m.step = stepError
		} else {
			m.keystoreInput.SetValue(msg.content)
		}
		return m, nil
	}

	// Update sub-components
	var cmds []tea.Cmd
	var cmd tea.Cmd

	m.spinner, cmd = m.spinner.Update(msg)
	cmds = append(cmds, cmd)

	progModel, progCmd := m.progressBar.Update(msg)
	if p, ok := progModel.(progress.Model); ok {
		m.progressBar = p
	}
	cmds = append(cmds, progCmd)

	m.passwordInput, cmd = m.passwordInput.Update(msg)
	cmds = append(cmds, cmd)

	m.confirmInput, cmd = m.confirmInput.Update(msg)
	cmds = append(cmds, cmd)

	m.nameInput, cmd = m.nameInput.Update(msg)
	cmds = append(cmds, cmd)

	m.mnemonicInput, cmd = m.mnemonicInput.Update(msg)
	cmds = append(cmds, cmd)

	m.keystoreInput, cmd = m.keystoreInput.Update(msg)
	cmds = append(cmds, cmd)

	m.keystorePathInput, cmd = m.keystorePathInput.Update(msg)
	cmds = append(cmds, cmd)

	m.privateKeyInput, cmd = m.privateKeyInput.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m setupTUIModel) View() string {
	var content string

	switch m.step {
	case stepWelcome:
		content = m.welcomeView()
	case stepWalletReplaceChoice:
		content = m.walletReplaceChoiceView()
	case stepWalletChoice:
		content = m.walletChoiceView()
	case stepWalletPassword:
		content = m.walletPasswordView()
	case stepWalletConfirmPassword:
		content = m.walletConfirmPasswordView()
	case stepWalletMnemonic:
		content = m.walletMnemonicView()
	case stepWalletName:
		content = m.walletNameView()
	case stepImportKeystoreChoice:
		content = m.importKeystoreChoiceView()
	case stepImportKeystoreJSON:
		content = m.importKeystoreJSONView()
	case stepImportKeystorePath:
		content = m.importKeystorePathView()
	case stepImportPrivateKey:
		content = m.importPrivateKeyView()
	case stepHardwareCheck:
		content = m.hardwareCheckView()
	case stepLibraryDownload:
		content = m.libraryDownloadView()
	case stepModelDownload:
		content = m.modelDownloadView()
	case stepLLMDownload:
		content = m.llmDownloadView()
	case stepSummary:
		content = m.summaryView()
	case stepError:
		content = m.errorView()
	}

	return m.styles.Container.Render(content)
}

// View methods
func (m setupTUIModel) welcomeView() string {
	var content string
	if m.walletExists {
		content = fmt.Sprintf(
			"%s\n\n%s\n\n%s\n\n%s",
			m.styles.Title.Render("🌸 Kawai DeAI Network - Setup"),
			m.styles.Subtitle.Render("Welcome back! Wallet already configured."),
			"An existing wallet was detected.\n"+
				"You can choose to keep using it or configure a new one.\n\n"+
				"Next steps:\n"+
				"  • Wallet configuration (Skip/Replace)\n"+
				"  • Hardware requirements check\n"+
				"  • Required libraries & models setup\n",
			m.styles.Help.Render("Press Enter to continue or Ctrl+C to exit"),
		)
	} else {
		content = fmt.Sprintf(
			"%s\n\n%s\n\n%s\n\n%s",
			m.styles.Title.Render("🌸 Kawai DeAI Network - Setup"),
			m.styles.Subtitle.Render("Earn KAWAI tokens by providing AI compute power"),
			"This wizard will guide you through setting up:\n"+
				"  • Hardware requirements check\n"+
				"  • Wallet configuration\n"+
				"  • Required libraries\n"+
				"  • AI models (Whisper, Stable Diffusion, LLM)\n",
			m.styles.Help.Render("Press Enter to continue or Ctrl+C to exit"),
		)
	}
	return content
}

func (m setupTUIModel) walletChoiceView() string {
	choices := []string{
		"Generate new mnemonic (recommended)",
		"Import existing mnemonic",
		"Import keystore JSON (MetaMask, etc.)",
		"Import private key",
	}

	var b strings.Builder
	b.WriteString(m.styles.Title.Render("🔐 Wallet Setup"))
	b.WriteString("\n\n")
	b.WriteString("Choose your setup method:\n\n")

	for i, choice := range choices {
		if i == m.walletChoice {
			b.WriteString(m.styles.ActiveStep.Render(fmt.Sprintf("  ● %s", choice)))
		} else {
			b.WriteString(m.styles.InactiveStep.Render(fmt.Sprintf("  ○ %s", choice)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(m.styles.Help.Render("Use ↑/↓ to select, Enter to confirm, ESC to go back"))

	return b.String()
}

func (m setupTUIModel) walletReplaceChoiceView() string {
	choices := []string{
		"Skip wallet setup (use existing wallet)",
		"Replace existing wallet",
	}

	var b strings.Builder
	b.WriteString(m.styles.Title.Render("🔐 Existing Wallet Detected"))
	b.WriteString("\n\n")
	b.WriteString(m.styles.Info.Render("A wallet is already configured on this system."))
	b.WriteString("\n\n")
	b.WriteString("What would you like to do?\n\n")

	for i, choice := range choices {
		if i == m.walletReplaceChoice {
			b.WriteString(m.styles.ActiveStep.Render(fmt.Sprintf("  ● %s", choice)))
		} else {
			b.WriteString(m.styles.InactiveStep.Render(fmt.Sprintf("  ○ %s", choice)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(m.styles.Help.Render("Use ↑/↓ to select, Enter to confirm, ESC to go back"))

	return b.String()
}

func calculatePasswordStrength(password string) (int, string, string) {
	score := 0
	if len(password) >= 8 {
		score += 25
	}
	if len(password) >= 12 {
		score += 15
	}
	if len(password) >= 16 {
		score += 10
	}
	if len(password) > 0 {
		if strings.ContainsAny(password, "ABCDEFGHIJKLMNOPQRSTUVWXYZ") {
			score += 15
		}
		if strings.ContainsAny(password, "0123456789") {
			score += 15
		}
		if strings.ContainsAny(password, "!@#$%^&*()_+-=[]{}|;:,.<>?") {
			score += 20
		}
	}

	score = min(score, 100)

	var label, color string
	switch {
	case score < 40:
		label = "Weak - Add more characters"
		color = "#ff4d4f"
	case score < 60:
		label = "Fair - Add uppercase or numbers"
		color = "#faad14"
	case score < 80:
		label = "Good - Add special characters"
		color = "#52c41a"
	default:
		label = "Strong password"
		color = "#1890ff"
	}

	return score, label, color
}

func renderPasswordBar(percent int, color string) string {
	filled := percent / 10
	empty := 10 - filled
	bar := strings.Repeat("█", filled) + strings.Repeat("░", empty)
	return fmt.Sprintf("[%s] %d%%", bar, percent)
}

func (m setupTUIModel) walletPasswordView() string {
	password := m.passwordInput.Value()
	var strengthSection string

	if password != "" {
		score, label, color := calculatePasswordStrength(password)
		bar := renderPasswordBar(score, color)
		strengthSection = fmt.Sprintf("\n%s\n%s\n",
			lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(bar),
			lipgloss.NewStyle().Foreground(lipgloss.Color(color)).Render(label),
		)
	}

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s%s\n%s",
		m.styles.Title.Render("🔐 Set Password"),
		"Enter a secure password (min 8 characters):",
		m.passwordInput.View(),
		strengthSection,
		m.styles.Help.Render("Press Enter to continue"),
	)
}

func (m setupTUIModel) walletConfirmPasswordView() string {
	helpMsg := m.styles.Help.Render("Press Enter to continue")

	// Check for mismatch
	if m.confirmInput.Value() != "" && m.passwordInput.Value() != m.confirmInput.Value() {
		helpMsg = m.styles.Error.Render("❌ Passwords do not match")
	}

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s",
		m.styles.Title.Render("🔐 Confirm Password"),
		"Confirm your password:",
		m.confirmInput.View(),
		helpMsg,
	)
}

func (m setupTUIModel) walletMnemonicView() string {
	if m.walletChoice != 0 {
		return fmt.Sprintf(
			"%s\n\n%s\n\n%s\n\n%s",
			m.styles.Title.Render("🔐 Import Mnemonic"),
			"Enter your 12 or 24 word mnemonic phrase:",
			m.mnemonicInput.View(),
			m.styles.Help.Render("Press Enter to continue"),
		)
	}

	words := strings.Fields(m.generatedMnemonic)
	var rows []string
	cols := 3 // Use 3 columns for better aspect ratio
	for i := 0; i < len(words); i += cols {
		end := i + cols
		if end > len(words) {
			end = len(words)
		}
		var rowParts []string
		for j := i; j < end; j++ {
			// formatted: " 1. word           "
			rowParts = append(rowParts, fmt.Sprintf("%2d. %-14s", j+1, words[j]))
		}
		rows = append(rows, strings.Join(rowParts, "    "))
	}

	// Simple mnemonic display without box
	mnemonicContent := lipgloss.NewStyle().
		Padding(1, 3).                     // Internal padding
		MarginLeft(2).                     // Indent from screen edge
		Render(strings.Join(rows, "\n\n")) // Double spacing between rows

	status := ""
	if m.mnemonicCopied {
		status = m.styles.Success.Render("\n\n   ✓ Copied to clipboard!")
	} else if m.copyError != "" {
		status = m.styles.Error.Render(fmt.Sprintf("\n\n   ✗ %s", m.copyError))
	}

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s\n%s%s",
		m.styles.Title.Render("⚠️  SAVE YOUR MNEMONIC"),
		m.styles.Warning.Render("   Write these words down on paper and store securely!"),
		mnemonicContent,
		"   Anyone with these words can access your funds.",
		m.styles.Help.Render("   Press 'c' to copy, Enter to confirm"),
		status,
	)
}

func (m setupTUIModel) walletNameView() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s",
		m.styles.Title.Render("🔐 Wallet Name"),
		"Enter a name for your wallet:",
		m.nameInput.View(),
		m.styles.Help.Render("Press Enter to continue or skip for default"),
	)
}

func (m setupTUIModel) importKeystoreChoiceView() string {
	choices := []string{
		"Paste keystore JSON content",
		"Enter file path to keystore",
	}

	var b strings.Builder
	b.WriteString(m.styles.Title.Render("📁 Import Keystore"))
	b.WriteString("\n\n")
	b.WriteString("How do you want to import the keystore?\n\n")

	for i, choice := range choices {
		if i == m.keystoreChoice {
			b.WriteString(m.styles.ActiveStep.Render(fmt.Sprintf("  ● %s", choice)))
		} else {
			b.WriteString(m.styles.InactiveStep.Render(fmt.Sprintf("  ○ %s", choice)))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(m.styles.Help.Render("Use ↑/↓ to select, Enter to confirm"))

	return b.String()
}

func (m setupTUIModel) importKeystoreJSONView() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s\n\n%s",
		m.styles.Title.Render("📁 Paste Keystore JSON"),
		"Paste your keystore JSON content below:",
		"(Usually found in MetaMask: Account → Account Details → Export Private Key)",
		m.keystoreInput.View(),
		m.styles.Help.Render("Press Enter to continue"),
	)
}

func (m setupTUIModel) importKeystorePathView() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s\n\n%s",
		m.styles.Title.Render("📁 Keystore File Path"),
		"Enter the full path to your keystore file:",
		"Example: /Users/username/Downloads/keystore.json",
		m.keystorePathInput.View(),
		m.styles.Help.Render("Press Enter to continue"),
	)
}

func (m setupTUIModel) importPrivateKeyView() string {
	privateKey := m.privateKeyInput.Value()
	var validationMsg string

	if privateKey != "" {
		cleanKey := strings.TrimPrefix(privateKey, "0x")
		cleanKey = strings.TrimPrefix(cleanKey, "0X")

		validHex := true
		if _, err := hex.DecodeString(cleanKey); err != nil {
			validHex = false
		}

		if len(cleanKey) == 64 && validHex {
			validationMsg = m.styles.Success.Render("✓ Valid format (64 hex characters)")
		} else {
			if !validHex {
				validationMsg = m.styles.Warning.Render("⚠ Invalid hex characters")
			} else {
				validationMsg = m.styles.Warning.Render(fmt.Sprintf("⚠ %d/64 characters", len(cleanKey)))
			}
		}
	}

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s\n%s\n\n%s",
		m.styles.Title.Render("🔑 Import Private Key"),
		"Enter your 64-character hex private key:",
		"(Example: 0x1234... or 1234...)",
		m.privateKeyInput.View(),
		validationMsg,
		m.styles.Help.Render("Press Enter to continue"),
	)
}

func (m setupTUIModel) hardwareCheckView() string {
	if m.hwSpecs == nil {
		return fmt.Sprintf(
			"%s\n\n%s %s",
			m.styles.Title.Render("🔍 Hardware Check"),
			m.spinner.View(),
			"Checking hardware requirements...",
		)
	}

	var b strings.Builder
	b.WriteString(m.styles.Title.Render("🔍 Hardware Check"))
	b.WriteString("\n\n")

	totalMemory := m.hwSpecs.AvailableRAM + m.hwSpecs.GPUMemory
	b.WriteString(fmt.Sprintf("Detected Hardware:\n"))
	b.WriteString(fmt.Sprintf("  CPU: %s (%d cores)\n", m.hwSpecs.CPU, m.hwSpecs.CPUCores))
	b.WriteString(fmt.Sprintf("  RAM: %dGB\n", m.hwSpecs.TotalRAM))
	if m.hwSpecs.GPUMemory > 0 {
		b.WriteString(fmt.Sprintf("  GPU: %s (%dGB VRAM)\n", m.hwSpecs.GPUModel, m.hwSpecs.GPUMemory))
	}
	b.WriteString(fmt.Sprintf("  Total Available: %dGB\n\n", totalMemory))

	if m.hardwarePassed {
		b.WriteString(m.styles.Success.Render("✓ Hardware check passed!"))
		b.WriteString("\n\n")
		b.WriteString(m.styles.Help.Render("Press Enter to continue"))
	} else {
		b.WriteString(m.styles.Error.Render("✗ Insufficient memory"))
		b.WriteString(fmt.Sprintf("\n\nRequired: 24GB RAM/VRAM\nAvailable: %dGB\n\n", totalMemory))
		b.WriteString(m.styles.Warning.Render("This server requires high-end hardware."))
	}

	return b.String()
}

func (m setupTUIModel) libraryDownloadView() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s %.1f%%\n\n%s",
		m.styles.Title.Render("📦 Downloading Libraries"),
		"Downloading llama.cpp and Stable Diffusion libraries...",
		m.progressBar.ViewAs(m.libraryProgress),
		m.libraryProgress*100,
		m.spinner.View()+" Please wait...",
	)
}

func (m setupTUIModel) modelDownloadView() string {
	var b strings.Builder
	b.WriteString(m.styles.Title.Render("🤖 Downloading Models"))
	b.WriteString("\n\n")

	b.WriteString("Whisper (Speech-to-Text):\n")
	b.WriteString(m.progressBar.ViewAs(m.whisperProgress))
	b.WriteString(fmt.Sprintf(" %.1f%%\n\n", m.whisperProgress*100))

	b.WriteString("Stable Diffusion (Image Generation):\n")
	b.WriteString(m.progressBar.ViewAs(m.sdProgress))
	b.WriteString(fmt.Sprintf(" %.1f%%\n\n", m.sdProgress*100))

	if m.whisperProgress < 1.0 || m.sdProgress < 1.0 {
		b.WriteString(m.spinner.View() + " Downloading...")
	}

	return b.String()
}

func (m setupTUIModel) llmDownloadView() string {
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s %.1f%%\n\n%s",
		m.styles.Title.Render("🧠 Downloading LLM"),
		"Downloading Nemotron 3 Nano (~18GB)",
		"This may take 10-60 minutes depending on your connection...",
		m.progressBar.ViewAs(m.llmProgress),
		m.llmProgress*100,
		m.spinner.View()+" Downloading...",
	)
}

func (m setupTUIModel) summaryView() string {
	var b strings.Builder
	b.WriteString(m.styles.Title.Render("✨ Setup Complete!"))
	b.WriteString("\n\n")

	if m.result.WalletCreated {
		b.WriteString(m.styles.Success.Render("✓ Wallet created: " + m.result.WalletAddress))
	} else {
		b.WriteString(m.styles.Success.Render("✓ Wallet: " + m.result.WalletAddress))
	}
	b.WriteString("\n")

	if m.result.LibraryReady {
		b.WriteString(m.styles.Success.Render("✓ Libraries ready"))
	} else {
		b.WriteString(m.styles.Error.Render("✗ Libraries failed"))
	}
	b.WriteString("\n")

	if m.result.WhisperReady {
		b.WriteString(m.styles.Success.Render("✓ Whisper model ready"))
	} else {
		b.WriteString(m.styles.Error.Render("✗ Whisper model failed"))
	}
	b.WriteString("\n")

	if m.result.StableDiffReady {
		b.WriteString(m.styles.Success.Render("✓ Stable Diffusion ready"))
	} else {
		b.WriteString(m.styles.Error.Render("✗ Stable Diffusion failed"))
	}
	b.WriteString("\n")

	if m.result.LLMReady {
		b.WriteString(m.styles.Success.Render("✓ LLM model ready (Nemotron 3 Nano)"))
	} else {
		b.WriteString(m.styles.Error.Render("✗ LLM model failed"))
	}
	b.WriteString("\n")

	if len(m.errors) > 0 {
		b.WriteString("\n")
		b.WriteString(m.styles.Warning.Render("Warnings:"))
		b.WriteString("\n")
		for _, err := range m.errors {
			b.WriteString(fmt.Sprintf("  • %v\n", err))
		}
	}

	b.WriteString("\n")
	b.WriteString(m.styles.Info.Render("Start the server with: ./server start"))
	b.WriteString("\n\n")
	b.WriteString(m.styles.Help.Render("Press Enter to exit"))

	return b.String()
}

func (m setupTUIModel) errorView() string {
	var b strings.Builder
	b.WriteString(m.styles.Title.Render("❌ Setup Failed"))
	b.WriteString("\n\n")

	b.WriteString(m.styles.Error.Render("The following errors occurred:"))
	b.WriteString("\n\n")

	for i, err := range m.errors {
		b.WriteString(fmt.Sprintf("%d. %v\n\n", i+1, err))
	}

	b.WriteString(m.styles.Help.Render("Press Enter to exit"))

	return b.String()
}

// Command methods
func (m setupTUIModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.step {
	case stepWelcome:
		// If wallet exists, show replace choice, otherwise go to wallet choice
		if m.walletExists {
			m.step = stepWalletReplaceChoice
		} else {
			m.step = stepWalletChoice
		}
		return m, nil

	case stepWalletReplaceChoice:
		if m.walletReplaceChoice == 0 {
			// Skip wallet setup - go directly to hardware check
			m.step = stepHardwareCheck
			return m, m.checkHardwareCmd()
		} else {
			// Replace wallet - go to wallet choice
			m.step = stepWalletChoice
		}
		return m, nil

	case stepWalletChoice:
		// All choices proceed to password setup first
		m.step = stepWalletPassword
		m.passwordInput.Focus()
		return m, nil

	case stepWalletPassword:
		if len(m.passwordInput.Value()) >= 8 {
			m.step = stepWalletConfirmPassword
			m.passwordInput.Blur()
			m.confirmInput.Focus()
		}
		return m, nil

	case stepWalletConfirmPassword:
		if m.passwordInput.Value() == m.confirmInput.Value() {
			switch m.walletChoice {
			case 0: // Generate mnemonic
				m.step = stepWalletMnemonic
				if m.generatedMnemonic == "" {
					return m, m.generateMnemonicCmd()
				}
			case 1: // Import mnemonic
				m.step = stepWalletMnemonic
				m.mnemonicInput.Focus()
			case 2: // Import keystore
				m.step = stepImportKeystoreChoice
			case 3: // Import private key
				m.step = stepImportPrivateKey
				m.privateKeyInput.Focus()
			}
		}
		return m, nil

	case stepWalletMnemonic:
		m.step = stepWalletName
		m.nameInput.Focus()
		return m, nil

	case stepImportKeystoreChoice:
		if m.keystoreChoice == 0 {
			m.step = stepImportKeystoreJSON
			m.keystoreInput.Focus()
		} else {
			m.step = stepImportKeystorePath
			m.keystorePathInput.Focus()
		}
		return m, nil

	case stepImportKeystoreJSON:
		m.step = stepWalletName
		m.nameInput.Focus()
		return m, nil

	case stepImportKeystorePath:
		m.step = stepWalletName
		m.nameInput.Focus()
		return m, m.readKeystoreFileCmd()

	case stepImportPrivateKey:
		key := m.privateKeyInput.Value()
		cleanKey := strings.TrimPrefix(key, "0x")
		cleanKey = strings.TrimPrefix(cleanKey, "0X")
		if len(cleanKey) == 64 {
			if _, err := hex.DecodeString(cleanKey); err == nil {
				m.step = stepWalletName
				m.nameInput.Focus()
				return m, nil
			}
		}
		return m, nil

	case stepWalletName:
		return m, m.createWalletCmd()

	case stepHardwareCheck:
		if m.hardwarePassed {
			m.step = stepLibraryDownload
			return m, m.downloadLibrariesCmd()
		}
		return m, tea.Quit

	case stepSummary:
		return m, tea.Quit

	case stepError:
		return m, tea.Quit
	}

	return m, nil
}

func (m setupTUIModel) handleBack() (tea.Model, tea.Cmd) {
	switch m.step {
	// Can't go back from welcome or error
	case stepWelcome, stepError:
		return m, nil

	// During async operations, just cancel/quit
	case stepHardwareCheck, stepLibraryDownload, stepModelDownload, stepLLMDownload:
		return m, tea.Quit

	// Wallet setup flow - go back appropriately
	case stepWalletReplaceChoice:
		m.step = stepWelcome
		return m, nil

	case stepWalletChoice:
		// If wallet exists, go back to replace choice, otherwise to welcome
		if m.walletExists {
			m.step = stepWalletReplaceChoice
		} else {
			m.step = stepWelcome
		}
		return m, nil

	case stepWalletPassword:
		// Clear password when going back
		m.passwordInput.SetValue("")
		m.step = stepWalletChoice
		return m, nil

	case stepWalletConfirmPassword:
		m.confirmInput.SetValue("")
		m.confirmInput.Blur()
		m.step = stepWalletPassword
		m.passwordInput.Focus()
		return m, nil

	case stepWalletMnemonic:
		// Going back from mnemonic depends on wallet choice
		m.mnemonicInput.SetValue("")
		m.mnemonicInput.Blur()
		if m.walletChoice == 0 {
			// Generated mnemonic - go back to confirm password
			m.step = stepWalletConfirmPassword
			m.confirmInput.Focus()
		} else {
			// Import mnemonic - go back to confirm password
			m.step = stepWalletConfirmPassword
			m.confirmInput.Focus()
		}
		return m, nil

	case stepImportKeystoreChoice:
		m.step = stepWalletConfirmPassword
		m.confirmInput.Focus()
		return m, nil

	case stepImportKeystoreJSON:
		m.keystoreInput.SetValue("")
		m.keystoreInput.Blur()
		m.step = stepImportKeystoreChoice
		return m, nil

	case stepImportKeystorePath:
		m.keystorePathInput.SetValue("")
		m.keystorePathInput.Blur()
		m.step = stepImportKeystoreChoice
		return m, nil

	case stepImportPrivateKey:
		m.privateKeyInput.SetValue("")
		m.privateKeyInput.Blur()
		m.step = stepWalletConfirmPassword
		m.confirmInput.Focus()
		return m, nil

	case stepWalletName:
		m.nameInput.SetValue("")
		m.nameInput.Blur()
		// Go back based on wallet choice
		switch m.walletChoice {
		case 0, 1:
			m.step = stepWalletMnemonic
			if m.walletChoice == 1 {
				m.mnemonicInput.Focus()
			}
		case 2:
			m.step = stepImportKeystoreChoice
		case 3:
			m.step = stepImportPrivateKey
			m.privateKeyInput.Focus()
			m.step = stepImportPrivateKey
		}
		return m, nil

	// After wallet creation - can't go back
	case stepSummary:
		return m, nil
	}

	return m, nil
}

func (m setupTUIModel) createWalletCmd() tea.Cmd {
	return func() tea.Msg {
		password := m.passwordInput.Value()
		name := m.nameInput.Value()
		if name == "" {
			name = "My Wallet"
		}

		var address string
		var err error

		switch m.walletChoice {
		case 0, 1: // Generate or import mnemonic
			var mnemonic string
			if m.walletChoice == 0 {
				mnemonic = m.generatedMnemonic
			} else {
				mnemonic = m.mnemonicInput.Value()
			}
			address, err = m.walletService.CreateWallet(password, mnemonic, name)

		case 2: // Import keystore
			var keystoreData string
			if m.keystoreChoice == 0 {
				keystoreData = m.keystoreInput.Value()
			} else {
				keystoreData = m.keystoreInput.Value() // Already read from file
			}
			_ = keystoreData // Will be used when keystore import is implemented
			address, err = "", fmt.Errorf("keystore import not yet implemented")

		case 3: // Import private key
			privateKey := m.privateKeyInput.Value()
			privateKey = strings.TrimPrefix(privateKey, "0x")
			privateKey = strings.TrimPrefix(privateKey, "0X")
			// Try to import private key - need to check if walletService has this method
			address, err = "", fmt.Errorf("private key import not yet implemented")
		}

		return walletCreatedMsg{address: address, err: err}
	}
}

func (m setupTUIModel) readKeystoreFileCmd() tea.Cmd {
	return func() tea.Msg {
		path := m.keystorePathInput.Value()
		data, err := os.ReadFile(path)
		if err != nil {
			return keystoreReadMsg{content: "", err: fmt.Errorf("failed to read keystore file: %w", err)}
		}
		return keystoreReadMsg{content: string(data), err: nil}
	}
}

func (m setupTUIModel) checkHardwareCmd() tea.Cmd {
	return func() tea.Msg {
		specs := hardware.DetectHardwareSpecs()
		totalMemory := specs.AvailableRAM + specs.GPUMemory
		passed := totalMemory >= 24 || m.skipHardwareCheck
		return hardwareCheckMsg{specs: specs, passed: passed}
	}
}

func (m setupTUIModel) generateMnemonicCmd() tea.Cmd {
	return func() tea.Msg {
		mnemonic, err := m.walletService.GenerateMnemonic()
		return mnemonicGeneratedMsg{mnemonic: mnemonic, err: err}
	}
}

func (m setupTUIModel) downloadLibrariesCmd() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		// Auto-detect platform
		arch, err := defaults.Arch("")
		if err != nil {
			return libraryDownloadMsg{progress: 0, done: false, err: fmt.Errorf("failed to detect arch: %w", err)}
		}

		opSys, err := defaults.OS("")
		if err != nil {
			return libraryDownloadMsg{progress: 0, done: false, err: fmt.Errorf("failed to detect OS: %w", err)}
		}

		processor, err := defaults.Processor("")
		if err != nil {
			return libraryDownloadMsg{progress: 0, done: false, err: fmt.Errorf("failed to detect processor: %w", err)}
		}

		// Create libs manager
		libMgr, err := libs.New(
			libs.WithBasePath(paths.Libraries()),
			libs.WithArch(arch),
			libs.WithOS(opSys),
			libs.WithProcessor(processor),
			libs.WithAllowUpgrade(true),
			libs.WithVersion(defaults.LibVersion("")),
		)
		if err != nil {
			return libraryDownloadMsg{progress: 0, done: false, err: fmt.Errorf("failed to create libs manager: %w", err)}
		}

		// Download llama.cpp
		downloadCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()

		_, err = libMgr.Download(downloadCtx, func(ctx context.Context, msg string, args ...any) {
			// Progress callback - update progress bar via channel
			// Try to parse percentage from log message e.g. "Downloading: 14.6% ..."
			if strings.Contains(msg, "%") {
				parts := strings.Split(msg, " ")
				for _, p := range parts {
					if strings.Contains(p, "%") {
						val := strings.TrimRight(p, "%,") // Remove % and potential trailing comma
						val = strings.Trim(val, "()")     // Remove parens just in case
						if f, err := strconv.ParseFloat(val, 64); err == nil {
							// Non-blocking send
							select {
							case m.progressChan <- f / 100.0:
							default:
							}
						}
					}
				}
			}
		})
		if err != nil {
			return libraryDownloadMsg{progress: 0, done: false, err: fmt.Errorf("failed to download llama.cpp: %w", err)}
		}

		// Setup Stable Diffusion library
		// Use EnsureLibraryWithProgress to get updates and suppress stdout logs
		err = stablediffusion.EnsureLibraryWithProgress(download.ProgressCallback(func(url string, bytesComplete, totalBytes int64, mbps float64, done bool) {
			if totalBytes > 0 {
				percent := float64(bytesComplete) / float64(totalBytes)
				select {
				case m.progressChan <- percent:
				default:
				}
			}
		}))
		if err != nil {
			return libraryDownloadMsg{progress: 0, done: false, err: fmt.Errorf("failed to setup SD library: %w", err)}
		}

		return libraryDownloadMsg{progress: 1.0, done: true}
	}
}

func (m setupTUIModel) downloadModelsCmd() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		// Setup Whisper model
		whisperModelsDir := paths.Models()
		if err := os.MkdirAll(whisperModelsDir, 0755); err != nil {
			return modelDownloadMsg{modelType: "whisper", progress: 0, done: false, err: fmt.Errorf("failed to create whisper models dir: %w", err)}
		}

		existingModels, _ := model.ListDownloadedModels(whisperModelsDir)
		if len(existingModels) == 0 {
			// Download whisper base model
			if err := model.DownloadModel("base", whisperModelsDir, nil); err != nil {
				return modelDownloadMsg{modelType: "whisper", progress: 0, done: false, err: fmt.Errorf("failed to download whisper model: %w", err)}
			}
		}

		// Setup Stable Diffusion model
		modelsPath := paths.Models()
		downloader := modeldownloader.New(modelsPath)

		modelFile, err := downloader.DiscoverModel()
		if err != nil {
			return modelDownloadMsg{modelType: "sd", progress: 0, done: false, err: fmt.Errorf("error discovering SD models: %w", err)}
		}

		if modelFile == "" {
			// Download default SD model
			modelFile, err = downloader.DownloadModelSimple(ctx, modeldownloader.DefaultModelURL)
			if err != nil {
				return modelDownloadMsg{modelType: "sd", progress: 0, done: false, err: fmt.Errorf("failed to download SD model: %w", err)}
			}
		}

		return modelDownloadMsg{modelType: "sd", progress: 1.0, done: true}
	}
}

func (m setupTUIModel) downloadLLMCmd() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()

		// Nemotron 3 Nano model info
		modelOrg := "unsloth"
		modelRepo := "Nemotron-3-Nano-30B-A3B-GGUF"
		modelFile := "Nemotron-3-Nano-30B-A3B-Q4_K_XL.gguf"
		modelsPath := paths.Models()
		modelPath := filepath.Join(modelsPath, modelOrg, modelRepo, modelFile)

		// Check if model already exists
		if _, err := os.Stat(modelPath); err == nil {
			return modelDownloadMsg{modelType: "llm", progress: 1.0, done: true}
		}

		// Download LLM model
		modelURL := fmt.Sprintf("https://huggingface.co/%s/%s/resolve/main/%s", modelOrg, modelRepo, modelFile)

		modelsManager, err := models.NewWithPaths(paths.Base())
		if err != nil {
			return modelDownloadMsg{modelType: "llm", progress: 0, done: false, err: fmt.Errorf("failed to create models manager: %w", err)}
		}

		downloadCtx, cancel := context.WithTimeout(ctx, 2*time.Hour)
		defer cancel()

		_, err = modelsManager.Download(downloadCtx, func(ctx context.Context, msg string, args ...any) {
			// Progress callback
		}, modelURL, "")
		if err != nil {
			return modelDownloadMsg{modelType: "llm", progress: 0, done: false, err: fmt.Errorf("failed to download LLM model: %w", err)}
		}

		return modelDownloadMsg{modelType: "llm", progress: 1.0, done: true}
	}
}

// NewSetupTUI creates a new TUI setup program
func NewSetupTUI(skipHardwareCheck bool) (*SetupResult, error) {
	kv, err := store.NewMultiNamespaceKVStore()
	if err != nil {
		return nil, fmt.Errorf("failed to connect to KV: %w", err)
	}

	walletService := services.NewWalletService("", kv)

	s := newStyles()

	// Initialize inputs
	passwordInput := textinput.New()
	passwordInput.Placeholder = "Enter password (min 8 chars)"
	passwordInput.EchoMode = textinput.EchoPassword
	passwordInput.Focus()

	confirmInput := textinput.New()
	confirmInput.Placeholder = "Confirm password"
	confirmInput.EchoMode = textinput.EchoPassword

	nameInput := textinput.New()
	nameInput.Placeholder = "My Wallet"

	mnemonicInput := textinput.New()
	mnemonicInput.Placeholder = "Enter mnemonic phrase"

	keystoreInput := textinput.New()
	keystoreInput.Placeholder = "Paste keystore JSON here"

	keystorePathInput := textinput.New()
	keystorePathInput.Placeholder = "/path/to/keystore.json"

	privateKeyInput := textinput.New()
	privateKeyInput.Placeholder = "0x... (64 hex characters)"

	sp := spinner.New()
	sp.Spinner = spinner.Dot
	sp.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6B9D"))

	prog := progress.New(progress.WithDefaultGradient())

	// Check if wallet already exists
	walletExists := walletService.HasWallet()

	model := setupTUIModel{
		styles:            s,
		step:              stepWelcome,
		walletService:     walletService,
		walletExists:      walletExists,
		passwordInput:     passwordInput,
		confirmInput:      confirmInput,
		nameInput:         nameInput,
		mnemonicInput:     mnemonicInput,
		keystoreInput:     keystoreInput,
		keystorePathInput: keystorePathInput,
		privateKeyInput:   privateKeyInput,
		spinner:           sp,
		progressBar:       prog,
		skipHardwareCheck: skipHardwareCheck,
		progressChan:      make(chan float64, 100),
		result: &SetupResult{
			Errors: make([]error, 0),
		},
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()
	if err != nil {
		return nil, err
	}

	m := finalModel.(setupTUIModel)
	return m.result, nil
}
