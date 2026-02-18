package kronk

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kawai-network/veridium/cmd/server/app/domain/ttsapp"
	"github.com/kawai-network/veridium/cmd/server/app/domain/whisperapp"
	"github.com/kawai-network/veridium/internal/paths"
	"github.com/kawai-network/veridium/internal/services"
	"github.com/kawai-network/veridium/pkg/hardware"
	"github.com/kawai-network/veridium/pkg/stablediffusion/modeldownloader"
	"github.com/kawai-network/veridium/pkg/store"
	"github.com/kawai-network/veridium/pkg/tools/defaults"
	"github.com/kawai-network/veridium/pkg/tools/libs"
	"github.com/kawai-network/veridium/pkg/tools/models"
)

// Constants for configuration
const (
	// MinRequiredMemoryGB is the minimum required RAM/VRAM in GB for hardware check
	MinRequiredMemoryGB = 24

	// MinPasswordLength is the minimum password length
	MinPasswordLength = 8

	// MaxPasswordLength is the maximum password length
	MaxPasswordLength = 128

	// PrivateKeyLength is the expected length of a hex private key
	PrivateKeyLength = 64

	// LibraryDownloadTimeout is the timeout for library downloads
	LibraryDownloadTimeout = 10 * time.Minute

	// ModelDownloadTimeout is the timeout for model downloads
	ModelDownloadTimeout = 30 * time.Minute

	// LLMDownloadTimeout is the timeout for LLM downloads
	LLMDownloadTimeout = 2 * time.Hour
)

// Model configuration (could be moved to config file)
const (
	DefaultLLMOrg       = "unsloth"
	DefaultLLMRepo      = "Nemotron-3-Nano-30B-A3B-GGUF"
	DefaultLLMFile      = "Nemotron-3-Nano-30B-A3B-Q4_K_M.gguf"
	DefaultWhisperModel = "base"
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
	stepHelp
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
	mnemonicCopied      bool
	copyError           string

	// Progress
	spinner     spinner.Model
	progressBar progress.Model
	// Additional progress bars for multiple downloads
	whisperProgressBar progress.Model
	sdProgressBar      progress.Model
	ttsProgressBar     progress.Model

	// Hardware
	hwSpecs           *hardware.HardwareSpecs
	hardwarePassed    bool
	skipHardwareCheck bool

	// Downloads
	libraryProgress float64
	whisperProgress float64
	sdProgress      float64
	ttsProgress     float64
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
	ctx          context.Context
	cancel       context.CancelFunc

	// Help screen
	showHelp bool
}

type progressMsg float64

// tickMsg is sent periodically to check progress

type tickMsg time.Time

// waitForProgress checks progress channel non-blocking
func waitForProgress(c chan float64) tea.Cmd {
	return func() tea.Msg {
		select {
		case p := <-c:
			return progressMsg(p)
		default:
			// No progress available, schedule next check
			return tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
				return tickMsg(t)
			})
		}
	}
}

// handleTick processes tick messages and continues polling if needed
func (m setupTUIModel) handleTick() (tea.Model, tea.Cmd) {
	// Check if we're in a download step and continue polling
	if m.step == stepLibraryDownload || m.step == stepModelDownload || m.step == stepLLMDownload {
		// Try to read progress
		select {
		case p := <-m.progressChan:
			if p < 0 {
				// Error signal - handle based on current step
				if m.step == stepLibraryDownload {
					m.step = stepError
					m.addError(fmt.Errorf("library download failed"))
					return m, nil
				}
				// For model downloads, mark as failed and continue to next
				if m.step == stepModelDownload {
					if !m.result.WhisperReady {
						m.addError(fmt.Errorf("whisper model download failed (optional)"))
						m.result.WhisperReady = false
						// Start SD download
						go m.downloadSDInternal()
						return m, tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
							return tickMsg(t)
						})
					} else if !m.result.StableDiffReady {
						m.addError(fmt.Errorf("stable diffusion model download failed"))
						m.result.StableDiffReady = false
						// Start TTS download
						go m.downloadTTSInternal()
						return m, tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
							return tickMsg(t)
						})
					} else if !m.result.TTSReady {
						m.addError(fmt.Errorf("tts model download failed (optional)"))
						m.result.TTSReady = false
						// Move to LLM step and start download
						m.step = stepLLMDownload
						return m, tea.Batch(m.downloadLLMCmd(), tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
							return tickMsg(t)
						}))
					}
				}
				if m.step == stepLLMDownload {
					m.addError(fmt.Errorf("llm model download failed"))
					m.result.LLMReady = false
					m.step = stepSummary
					return m, nil
				}
				return m, tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
					return tickMsg(t)
				})
			}
			if p >= 1.0 {
				// Download complete - handle state transitions
				switch m.step {
				case stepLibraryDownload:
					m.result.LibraryReady = true
					// Note: TTSReady is set after TTS model download, not library download
					m.step = stepModelDownload
					return m, tea.Batch(m.downloadModelsCmd(), tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
						return tickMsg(t)
					}))
				case stepModelDownload:
					// Check which model just completed and start next one
					if !m.result.WhisperReady {
						m.result.WhisperReady = true
						m.whisperProgress = 1.0
						// Start SD download
						go m.downloadSDInternal()
						// Continue polling
						return m, tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
							return tickMsg(t)
						})
					} else if !m.result.StableDiffReady {
						m.result.StableDiffReady = true
						m.sdProgress = 1.0
						// Start TTS download
						go m.downloadTTSInternal()
						// Continue polling
						return m, tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
							return tickMsg(t)
						})
					} else if !m.result.TTSReady {
						m.result.TTSReady = true
						m.ttsProgress = 1.0
						// Move to LLM step and start download
						m.step = stepLLMDownload
						return m, tea.Batch(m.downloadLLMCmd(), tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
							return tickMsg(t)
						}))
					}
				case stepLLMDownload:
					m.result.LLMReady = true
					m.llmProgress = 1.0
					m.step = stepSummary
					return m, nil
				}
			} else {
				// Update progress based on current step and what's downloading
				switch m.step {
				case stepLibraryDownload:
					m.libraryProgress = p
				case stepModelDownload:
					// Determine which model is downloading based on completion status
					if !m.result.WhisperReady {
						m.whisperProgress = p
					} else if !m.result.StableDiffReady {
						m.sdProgress = p
					} else if !m.result.TTSReady {
						m.ttsProgress = p
					}
				case stepLLMDownload:
					m.llmProgress = p
				}
			}
		default:
		}
		return m, tea.Tick(50*time.Millisecond, func(t time.Time) tea.Msg {
			return tickMsg(t)
		})
	}
	return m, nil
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
	// Typed progress messages for better tracking
	libraryProgressMsg float64
	whisperProgressMsg float64
	sdProgressMsg      float64
	ttsProgressMsg     float64
	llmProgressMsg     float64
)

// Constants for error handling
const (
	maxErrors  = 10
	maxRetries = 3
	retryDelay = 2 * time.Second
)

// Init initializes the TUI
func (m setupTUIModel) Init() tea.Cmd {
	return tea.Batch(
		m.spinner.Tick,
		m.progressBar.Init(),
		m.whisperProgressBar.Init(),
		m.sdProgressBar.Init(),
		m.ttsProgressBar.Init(),
	)
}

// Update handles messages and updates the model
func (m setupTUIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Update progress bar widths
		progressWidth := msg.Width - 20
		if progressWidth > 60 {
			progressWidth = 60
		}
		if progressWidth < 20 {
			progressWidth = 20
		}
		m.progressBar.Width = progressWidth
		m.whisperProgressBar.Width = progressWidth
		m.sdProgressBar.Width = progressWidth
		m.ttsProgressBar.Width = progressWidth
		return m, nil

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		var cmds []tea.Cmd

		progressModel, cmd := m.progressBar.Update(msg)
		m.progressBar = progressModel.(progress.Model)
		cmds = append(cmds, cmd)

		whisperModel, cmd := m.whisperProgressBar.Update(msg)
		m.whisperProgressBar = whisperModel.(progress.Model)
		cmds = append(cmds, cmd)

		sdModel, cmd := m.sdProgressBar.Update(msg)
		m.sdProgressBar = sdModel.(progress.Model)
		cmds = append(cmds, cmd)

		ttsModel, cmd := m.ttsProgressBar.Update(msg)
		m.ttsProgressBar = ttsModel.(progress.Model)
		cmds = append(cmds, cmd)

		return m, tea.Batch(cmds...)

	case progressMsg:
		if m.step == stepLibraryDownload {
			m.libraryProgress = float64(msg)
		} else if m.step == stepModelDownload {
			// Basic heuristic: update whatever progress is active (whisper, SD, or TTS)
			if m.result.WhisperReady && m.result.StableDiffReady {
				m.ttsProgress = float64(msg)
			} else if m.result.WhisperReady {
				m.sdProgress = float64(msg)
			} else {
				m.whisperProgress = float64(msg)
			}
		} else if m.step == stepLLMDownload {
			m.llmProgress = float64(msg)
		}
		return m, waitForProgress(m.progressChan)

	case libraryProgressMsg:
		m.libraryProgress = float64(msg)
		return m, nil

	case whisperProgressMsg:
		m.whisperProgress = float64(msg)
		return m, nil

	case sdProgressMsg:
		m.sdProgress = float64(msg)
		return m, nil

	case ttsProgressMsg:
		m.ttsProgress = float64(msg)
		return m, nil

	case llmProgressMsg:
		m.llmProgress = float64(msg)
		return m, nil

	case tickMsg:
		return m.handleTick()

	case tea.KeyMsg:
		// F1 for help
		if msg.String() == "f1" {
			m.showHelp = !m.showHelp
			return m, nil
		}

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
			if m.cancel != nil {
				m.cancel()
			}
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
			m.addError(msg.err)
			m.hardwarePassed = false
		} else {
			m.hwSpecs = msg.specs
			m.hardwarePassed = msg.passed
		}
		return m, nil

	case libraryDownloadMsg:
		if msg.err != nil {
			// Check for timeout
			if m.ctx.Err() == context.DeadlineExceeded {
				m.addError(fmt.Errorf("library download timeout: %w", msg.err))
			} else {
				m.addError(fmt.Errorf("library download failed: %w", msg.err))
			}
			m.step = stepError
		} else if msg.done {
			m.result.LibraryReady = true
			// Note: TTSReady is set after TTS model download, not library download
			m.step = stepModelDownload
			return m, m.downloadModelsCmd()
		} else {
			m.libraryProgress = msg.progress
		}
		return m, nil

	case modelDownloadMsg:
		if msg.done {
			switch msg.modelType {
			case "whisper":
				if msg.err != nil {
					if m.ctx.Err() == context.DeadlineExceeded {
						m.addError(fmt.Errorf("whisper model download timeout (optional): %w", msg.err))
					} else {
						m.addError(fmt.Errorf("whisper model download failed (optional): %w", msg.err))
					}
					m.result.WhisperReady = false
				} else {
					m.result.WhisperReady = true
				}
				return m, m.downloadSDCmd()
			case "sd":
				if msg.err != nil {
					if m.ctx.Err() == context.DeadlineExceeded {
						m.addError(fmt.Errorf("stable diffusion model download timeout: %w", msg.err))
					} else {
						m.addError(fmt.Errorf("stable diffusion model download failed: %w", msg.err))
					}
					m.result.StableDiffReady = false
				} else {
					m.result.StableDiffReady = true
				}
				return m, m.downloadTTSCmd()
			case "tts":
				if msg.err != nil {
					if m.ctx.Err() == context.DeadlineExceeded {
						m.addError(fmt.Errorf("tts model download timeout (optional): %w", msg.err))
					} else {
						m.addError(fmt.Errorf("tts model download failed (optional): %w", msg.err))
					}
					m.result.TTSReady = false
				} else {
					m.result.TTSReady = true
				}
				m.step = stepLLMDownload
				return m, m.downloadLLMCmd()
			case "llm":
				if msg.err != nil {
					if m.ctx.Err() == context.DeadlineExceeded {
						m.addError(fmt.Errorf("llm model download timeout: %w", msg.err))
					} else {
						m.addError(fmt.Errorf("llm model download failed: %w", msg.err))
					}
					m.result.LLMReady = false
				} else {
					m.result.LLMReady = true
				}
				m.step = stepSummary
			}
		} else {
			switch msg.modelType {
			case "whisper":
				m.whisperProgress = msg.progress
			case "sd":
				m.sdProgress = msg.progress
			case "tts":
				m.ttsProgress = msg.progress
			case "llm":
				m.llmProgress = msg.progress
			}
		}
		return m, nil

	case walletCreatedMsg:
		if msg.err != nil {
			m.addError(msg.err)
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
			m.addError(msg.err)
		} else {
			m.walletAddress = msg.address
			m.result.WalletAddress = msg.address
			m.step = stepHardwareCheck
			return m, m.checkHardwareCmd()
		}
		return m, nil

	case mnemonicGeneratedMsg:
		if msg.err != nil {
			m.addError(fmt.Errorf("failed to generate mnemonic: %w", msg.err))
			m.step = stepError
		} else {
			m.generatedMnemonic = msg.mnemonic
		}
		return m, nil

	case keystoreReadMsg:
		if msg.err != nil {
			m.addError(msg.err)
			m.step = stepError
		} else {
			m.keystoreInput.SetValue(msg.content)
		}
		return m, nil
	}

	// Update sub-components - only update active input
	var cmds []tea.Cmd
	var cmd tea.Cmd

	// Only update the active input
	switch m.step {
	case stepWalletPassword:
		m.passwordInput, cmd = m.passwordInput.Update(msg)
		cmds = append(cmds, cmd)
	case stepWalletConfirmPassword:
		m.confirmInput, cmd = m.confirmInput.Update(msg)
		cmds = append(cmds, cmd)
	case stepWalletName:
		m.nameInput, cmd = m.nameInput.Update(msg)
		cmds = append(cmds, cmd)
	case stepWalletMnemonic:
		if m.walletChoice == 1 { // Import mnemonic
			m.mnemonicInput, cmd = m.mnemonicInput.Update(msg)
			cmds = append(cmds, cmd)
		}
	case stepImportKeystoreJSON:
		m.keystoreInput, cmd = m.keystoreInput.Update(msg)
		cmds = append(cmds, cmd)
	case stepImportKeystorePath:
		m.keystorePathInput, cmd = m.keystorePathInput.Update(msg)
		cmds = append(cmds, cmd)
	case stepImportPrivateKey:
		m.privateKeyInput, cmd = m.privateKeyInput.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

// View renders the UI
func (m setupTUIModel) View() string {
	if m.showHelp {
		return m.helpView()
	}

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

func validateKeystoreJSON(jsonStr string) (bool, string) {
	// Check if JSON is empty
	if strings.TrimSpace(jsonStr) == "" {
		return false, "Keystore JSON cannot be empty"
	}

	// Parse JSON properly
	var ks struct {
		Address string `json:"address"`
		Crypto  struct {
			Kdf        string `json:"kdf"`
			Ciphertext string `json:"ciphertext"`
		} `json:"crypto"`
		Version int `json:"version"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &ks); err != nil {
		return false, fmt.Sprintf("Invalid JSON format: %v", err)
	}

	// Validate required fields
	if ks.Address == "" {
		return false, "Missing required field: address"
	}
	if ks.Version == 0 {
		return false, "Missing or invalid required field: version"
	}
	if ks.Crypto.Kdf == "" && ks.Crypto.Ciphertext == "" {
		return false, "Invalid keystore: missing crypto information (kdf or ciphertext)"
	}

	return true, ""
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
	keystoreJSON := m.keystoreInput.Value()
	var validationMsg string

	if keystoreJSON != "" {
		valid, msg := validateKeystoreJSON(keystoreJSON)
		if valid {
			validationMsg = m.styles.Success.Render("✓ Valid keystore format")
		} else {
			validationMsg = m.styles.Error.Render("✗ " + msg)
		}
	}

	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%s\n%s\n\n%s",
		m.styles.Title.Render("📁 Paste Keystore JSON"),
		"Paste your keystore JSON content below:",
		"(Usually found in MetaMask: Account → Account Details → Export Private Key)",
		m.keystoreInput.View(),
		validationMsg,
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
	b.WriteString("Detected Hardware:\n")
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
	libPct := m.libraryProgress * 100
	return fmt.Sprintf(
		"%s\n\n%s\n\n%.1f%%\n%s\n\n%s",
		m.styles.Title.Render("📦 Downloading Libraries"),
		"Downloading llama.cpp, whisper.cpp, Stable Diffusion, and TTS libraries...",
		libPct,
		m.progressBar.ViewAs(m.libraryProgress),
		m.spinner.View()+" Please wait...",
	)
}

func (m setupTUIModel) modelDownloadView() string {
	var b strings.Builder
	b.WriteString(m.styles.Title.Render("🤖 Downloading Models"))
	b.WriteString("\n\n")

	whisperPct := m.whisperProgress * 100
	b.WriteString(fmt.Sprintf("Whisper (Speech-to-Text): %.1f%%\n", whisperPct))
	b.WriteString(m.whisperProgressBar.ViewAs(m.whisperProgress))
	b.WriteString("\n\n")

	sdPct := m.sdProgress * 100
	b.WriteString(fmt.Sprintf("Stable Diffusion (Image Generation): %.1f%%\n", sdPct))
	b.WriteString(m.sdProgressBar.ViewAs(m.sdProgress))
	b.WriteString("\n\n")

	ttsPct := m.ttsProgress * 100
	b.WriteString(fmt.Sprintf("TTS (Text-to-Speech): %.1f%%\n", ttsPct))
	b.WriteString(m.ttsProgressBar.ViewAs(m.ttsProgress))
	b.WriteString("\n\n")

	if m.whisperProgress < 1.0 || m.sdProgress < 1.0 || m.ttsProgress < 1.0 {
		b.WriteString(m.spinner.View() + " Downloading...")
	}

	return b.String()
}

func (m setupTUIModel) llmDownloadView() string {
	llmPct := m.llmProgress * 100
	return fmt.Sprintf(
		"%s\n\n%s\n\n%s\n\n%.1f%%\n%s\n\n%s",
		m.styles.Title.Render("🧠 Downloading LLM"),
		"Downloading Nemotron 3 Nano (~18GB)",
		"This may take 10-60 minutes depending on your connection...",
		llmPct,
		m.progressBar.ViewAs(m.llmProgress),
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

	if m.result.TTSReady {
		b.WriteString(m.styles.Success.Render("✓ TTS ready"))
	} else {
		b.WriteString(m.styles.Error.Render("✗ TTS failed"))
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

func (m setupTUIModel) helpView() string {
	var b strings.Builder
	b.WriteString(m.styles.Title.Render("📖 Help"))
	b.WriteString("\n\n")

	b.WriteString(m.styles.Info.Render("Keyboard Shortcuts:"))
	b.WriteString("\n\n")
	b.WriteString("  F1         - Toggle this help screen\n")
	b.WriteString("  Enter      - Confirm / Next step\n")
	b.WriteString("  ESC        - Go back / Cancel\n")
	b.WriteString("  Ctrl+C     - Exit setup\n")
	b.WriteString("  ↑/↓        - Navigate menu options\n")
	b.WriteString("  c          - Copy mnemonic (on mnemonic screen)\n")
	b.WriteString("\n")

	b.WriteString(m.styles.Info.Render("Setup Steps:"))
	b.WriteString("\n\n")
	b.WriteString("  1. Wallet Configuration\n")
	b.WriteString("     - Generate new or import existing wallet\n")
	b.WriteString("     - Set secure password (min 8 characters)\n")
	b.WriteString("\n")
	b.WriteString("  2. Hardware Check\n")
	b.WriteString("     - Verify system meets requirements (24GB RAM/VRAM)\n")
	b.WriteString("\n")
	b.WriteString("  3. Library Downloads\n")
	b.WriteString("     - llama.cpp, whisper.cpp, Stable Diffusion, TTS\n")
	b.WriteString("\n")
	b.WriteString("  4. Model Downloads\n")
	b.WriteString("     - Whisper (Speech-to-Text)\n")
	b.WriteString("     - Stable Diffusion (Image Generation)\n")
	b.WriteString("     - TTS (Text-to-Speech)\n")
	b.WriteString("     - LLM (Nemotron 3 Nano ~18GB)\n")
	b.WriteString("\n\n")

	b.WriteString(m.styles.Help.Render("Press F1 to close help"))

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
		m.blurAllInputs()
		m.passwordInput.Focus()
		return m, nil

	case stepWalletPassword:
		password := m.passwordInput.Value()
		if len(password) < MinPasswordLength {
			// Don't proceed if password too short
			return m, nil
		}
		score, _, _ := calculatePasswordStrength(password)
		if score < 40 {
			// Weak password, but allow (user's choice)
		}
		m.step = stepWalletConfirmPassword
		m.blurAllInputs()
		m.confirmInput.Focus()
		return m, nil

	case stepWalletConfirmPassword:
		if m.passwordInput.Value() == m.confirmInput.Value() {
			m.blurAllInputs()
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
		m.blurAllInputs()
		m.nameInput.Focus()
		return m, nil

	case stepImportKeystoreChoice:
		m.blurAllInputs()
		if m.keystoreChoice == 0 {
			m.step = stepImportKeystoreJSON
			m.keystoreInput.Focus()
		} else {
			m.step = stepImportKeystorePath
			m.keystorePathInput.Focus()
		}
		return m, nil

	case stepImportKeystoreJSON:
		// Validate keystore JSON before proceeding
		keystoreJSON := m.keystoreInput.Value()
		if keystoreJSON == "" {
			// Empty input, don't proceed
			return m, nil
		}
		valid, _ := validateKeystoreJSON(keystoreJSON)
		if !valid {
			// Invalid JSON, don't proceed
			return m, nil
		}
		// Valid JSON, proceed to wallet name
		m.step = stepWalletName
		m.blurAllInputs()
		m.nameInput.Focus()
		return m, nil

	case stepImportKeystorePath:
		m.step = stepWalletName
		m.blurAllInputs()
		m.nameInput.Focus()
		return m, m.readKeystoreFileCmd()

	case stepImportPrivateKey:
		key := m.privateKeyInput.Value()
		cleanKey := strings.TrimPrefix(key, "0x")
		cleanKey = strings.TrimPrefix(cleanKey, "0X")
		if len(cleanKey) == 64 {
			if _, err := hex.DecodeString(cleanKey); err == nil {
				m.step = stepWalletName
				m.blurAllInputs()
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
			// downloadLibrariesCmd already returns tickMsg to start tick-based polling
			// Do not add waitForProgress to avoid race condition with dual readers on progressChan
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
		m.blurAllInputs()
		m.step = stepWalletPassword
		m.passwordInput.Focus()
		return m, nil

	case stepWalletMnemonic:
		// Going back from mnemonic depends on wallet choice
		m.mnemonicInput.SetValue("")
		m.blurAllInputs()
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
		m.blurAllInputs()
		m.step = stepWalletConfirmPassword
		m.confirmInput.Focus()
		return m, nil

	case stepImportKeystoreJSON:
		m.keystoreInput.SetValue("")
		m.blurAllInputs()
		m.step = stepImportKeystoreChoice
		return m, nil

	case stepImportKeystorePath:
		m.keystorePathInput.SetValue("")
		m.blurAllInputs()
		m.step = stepImportKeystoreChoice
		return m, nil

	case stepImportPrivateKey:
		m.privateKeyInput.SetValue("")
		m.blurAllInputs()
		m.step = stepWalletConfirmPassword
		m.confirmInput.Focus()
		return m, nil

	case stepWalletName:
		m.nameInput.SetValue("")
		m.blurAllInputs()
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
		}
		return m, nil

	// After wallet creation - can't go back
	case stepSummary:
		return m, nil
	}

	return m, nil
}

// Helper methods
func (m *setupTUIModel) blurAllInputs() {
	m.passwordInput.Blur()
	m.confirmInput.Blur()
	m.nameInput.Blur()
	m.mnemonicInput.Blur()
	m.keystoreInput.Blur()
	m.keystorePathInput.Blur()
	m.privateKeyInput.Blur()
}

func (m *setupTUIModel) addError(err error) {
	m.errors = append(m.errors, err)
	if len(m.errors) > maxErrors {
		m.errors = m.errors[len(m.errors)-maxErrors:]
	}
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
			address, err = m.walletService.ImportKeystore(keystoreData, password, name)

		case 3: // Import private key
			privateKey := m.privateKeyInput.Value()
			privateKey = strings.TrimPrefix(privateKey, "0x")
			privateKey = strings.TrimPrefix(privateKey, "0X")
			address, err = m.walletService.ImportPrivateKey(privateKey, password, name)
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
	// Start download in goroutine and return immediately to not block UI
	go func() {
		// Auto-detect platform
		arch, err := defaults.Arch("")
		if err != nil {
			m.progressChan <- -1 // Signal error
			return
		}

		opSys, err := defaults.OS("")
		if err != nil {
			m.progressChan <- -1
			return
		}

		processor, err := defaults.Processor("")
		if err != nil {
			m.progressChan <- -1
			return
		}

		libraries := []libs.LibraryType{libs.LibraryLlama, libs.LibraryWhisper, libs.LibraryStableDiffusion, libs.LibraryTTS}
		totalLibs := float64(len(libraries))

		for i, libType := range libraries {
			// Check context cancellation
			if m.ctx.Err() != nil {
				m.progressChan <- -1
				return
			}

			// Create libs manager for each library type
			libMgr, err := libs.New(
				libs.WithBasePath(paths.Base()),
				libs.WithArch(arch),
				libs.WithOS(opSys),
				libs.WithProcessor(processor),
				libs.WithAllowUpgrade(true),
				libs.WithLibraryType(libType),
			)
			if err != nil {
				m.progressChan <- -1
				return
			}

			downloadCtx, cancel := context.WithTimeout(m.ctx, LibraryDownloadTimeout)

			baseProgress := float64(i) / totalLibs
			libWeight := 1.0 / totalLibs

			_, err = libMgr.DownloadWithProgress(downloadCtx, func(ctx context.Context, msg string, args ...any) {
				// Log message callback
			}, func(bytesComplete, totalBytes int64, mbps float64, done bool) {
				if totalBytes > 0 {
					percent := float64(bytesComplete) / float64(totalBytes)
					overallProgress := baseProgress + percent*libWeight
					select {
					case m.progressChan <- overallProgress:
					default:
					}
				}
			})
			cancel()

			if err != nil {
				m.progressChan <- -1
				return
			}

			// Update progress when library download completes
			completedProgress := float64(i+1) / totalLibs
			select {
			case m.progressChan <- completedProgress:
			default:
			}
		}

		// Signal completion
		m.progressChan <- 1.0
	}()

	// Return tick to start progress polling
	return func() tea.Msg {
		return tickMsg(time.Now())
	}
}

func (m setupTUIModel) downloadModelsCmd() tea.Cmd {
	// Start whisper download in goroutine
	go func() {
		// Check context cancellation
		if m.ctx.Err() != nil {
			m.progressChan <- -1
			return
		}

		downloadCtx, cancel := context.WithTimeout(m.ctx, ModelDownloadTimeout)
		defer cancel()

		// Setup Whisper model
		existingModels, _ := whisperapp.ListDownloadedModels()
		if len(existingModels) == 0 {
			// Download whisper base model with progress reporting
			progressCallback := func(currentBytes, totalBytes int64) {
				if totalBytes > 0 {
					percent := float64(currentBytes) / float64(totalBytes)
					select {
					case m.progressChan <- percent:
					default:
					}
				}
			}

			if err := whisperapp.DownloadModelWithProgress(downloadCtx, "base", progressCallback); err != nil {
				// Log error but continue
				fmt.Printf("Warning: failed to download whisper model: %v\n", err)
			}
		}

		// Signal whisper complete
		m.progressChan <- 1.0
	}()

	return func() tea.Msg {
		return tickMsg(time.Now())
	}
}

func (m setupTUIModel) downloadSDCmd() tea.Cmd {
	return func() tea.Msg {
		// Check context cancellation
		if m.ctx.Err() != nil {
			return modelDownloadMsg{modelType: "sd", progress: 0, done: false, err: m.ctx.Err()}
		}

		downloadCtx, cancel := context.WithTimeout(m.ctx, ModelDownloadTimeout)
		defer cancel()

		// Setup Stable Diffusion model
		modelsPath := paths.Models()
		downloader := modeldownloader.New(modelsPath)

		modelFile, err := downloader.DiscoverModel()
		if err != nil {
			return modelDownloadMsg{modelType: "sd", progress: 0, done: false, err: fmt.Errorf("error discovering SD models: %w", err)}
		}

		if modelFile == "" {
			// Download default SD model with progress reporting
			progressCallback := func(bytesComplete, totalBytes int64, mbps float64, done bool) {
				if totalBytes > 0 {
					percent := float64(bytesComplete) / float64(totalBytes)
					select {
					case m.progressChan <- percent:
					default:
					}
				}
			}

			modelFile, err = downloader.DownloadModel(downloadCtx, modeldownloader.DefaultModelURL, progressCallback)
			if err != nil {
				return modelDownloadMsg{modelType: "sd", progress: 0, done: false, err: fmt.Errorf("failed to download SD model: %w", err)}
			}
		}

		return modelDownloadMsg{modelType: "sd", progress: 1.0, done: true}
	}
}

func (m setupTUIModel) downloadTTSCmd() tea.Cmd {
	return func() tea.Msg {
		// Check context cancellation
		if m.ctx.Err() != nil {
			return modelDownloadMsg{modelType: "tts", progress: 0, done: false, err: m.ctx.Err()}
		}

		downloadCtx, cancel := context.WithTimeout(m.ctx, ModelDownloadTimeout)
		defer cancel()

		// Setup TTS model - uses standard path structure: models/{org}/{repo}/{filename}
		downloader := ttsapp.NewModelDownloader("")

		modelFile, err := downloader.DiscoverModel()
		if err != nil {
			return modelDownloadMsg{modelType: "tts", progress: 0, done: false, err: fmt.Errorf("error discovering TTS models: %w", err)}
		}

		if modelFile == "" {
			// Download default TTS model with progress reporting
			progressCallback := func(bytesComplete, totalBytes int64, mbps float64, done bool) {
				if totalBytes > 0 {
					percent := float64(bytesComplete) / float64(totalBytes)
					select {
					case m.progressChan <- percent:
					default:
					}
				}
			}

			modelFile, err = downloader.DownloadModel(downloadCtx, ttsapp.DefaultTTSModelURL, progressCallback)
			if err != nil {
				// TTS is optional, don't fail setup if download fails
				return modelDownloadMsg{modelType: "tts", progress: 1.0, done: true, err: fmt.Errorf("failed to download TTS model: %w", err)}
			}
		}

		return modelDownloadMsg{modelType: "tts", progress: 1.0, done: true}
	}
}

func (m setupTUIModel) downloadLLMCmd() tea.Cmd {
	// Start LLM download in goroutine
	go func() {
		// Check context cancellation
		if m.ctx.Err() != nil {
			m.progressChan <- -1
			return
		}

		// Nemotron 3 Nano model info
		modelOrg := DefaultLLMOrg
		modelRepo := DefaultLLMRepo
		modelFile := DefaultLLMFile
		modelsPath := paths.Models()
		modelPath := filepath.Join(modelsPath, modelOrg, modelRepo, modelFile)

		// Check if model already exists
		if _, err := os.Stat(modelPath); err == nil {
			m.progressChan <- 1.0
			return
		}

		// Download LLM model
		modelURL := fmt.Sprintf("https://huggingface.co/%s/%s/resolve/main/%s", modelOrg, modelRepo, modelFile)

		modelsManager, err := models.NewWithPaths(paths.Base())
		if err != nil {
			m.progressChan <- -1
			return
		}

		downloadCtx, cancel := context.WithTimeout(m.ctx, LLMDownloadTimeout)
		defer cancel()

		// Create a logger callback that parses progress from log messages
		_, err = modelsManager.Download(downloadCtx, func(ctx context.Context, msg string, args ...any) {
			// Parse progress from log messages
			if len(args) >= 4 {
				msgStr := fmt.Sprintf(msg, args...)
				var currentMiB, totalMiB int64
				n, _ := fmt.Sscanf(msgStr, "download-model: Downloading %s... %d MiB of %d MiB", new(string), &currentMiB, &totalMiB)
				if n == 3 && totalMiB > 0 {
					percent := float64(currentMiB) / float64(totalMiB)
					select {
					case m.progressChan <- percent:
					default:
					}
				}
			}
		}, modelURL, "")
		if err != nil {
			m.progressChan <- -1
			return
		}

		m.progressChan <- 1.0
	}()

	return func() tea.Msg {
		return tickMsg(time.Now())
	}
}

// downloadSDInternal downloads Stable Diffusion model in goroutine
func (m setupTUIModel) downloadSDInternal() {
	// Check context cancellation
	if m.ctx.Err() != nil {
		m.progressChan <- -1
		return
	}

	downloadCtx, cancel := context.WithTimeout(m.ctx, ModelDownloadTimeout)
	defer cancel()

	// Setup Stable Diffusion model
	modelsPath := paths.Models()
	downloader := modeldownloader.New(modelsPath)

	modelFile, err := downloader.DiscoverModel()
	if err != nil {
		m.progressChan <- -1
		return
	}

	if modelFile == "" {
		// Download default SD model with progress reporting
		progressCallback := func(bytesComplete, totalBytes int64, mbps float64, done bool) {
			if totalBytes > 0 {
				percent := float64(bytesComplete) / float64(totalBytes)
				select {
				case m.progressChan <- percent:
				default:
				}
			}
		}

		modelFile, err = downloader.DownloadModel(downloadCtx, modeldownloader.DefaultModelURL, progressCallback)
		if err != nil {
			m.progressChan <- -1
			return
		}
	}

	// Signal completion
	m.progressChan <- 1.0
}

// downloadTTSInternal downloads TTS model in goroutine
func (m setupTUIModel) downloadTTSInternal() {
	// Check context cancellation
	if m.ctx.Err() != nil {
		m.progressChan <- -1
		return
	}

	downloadCtx, cancel := context.WithTimeout(m.ctx, ModelDownloadTimeout)
	defer cancel()

	// Setup TTS model - uses standard path structure: models/{org}/{repo}/{filename}
	downloader := ttsapp.NewModelDownloader("")

	modelFile, err := downloader.DiscoverModel()
	if err != nil {
		m.progressChan <- -1
		return
	}

	if modelFile == "" {
		// Download default TTS model with progress reporting
		progressCallback := func(bytesComplete, totalBytes int64, mbps float64, done bool) {
			if totalBytes > 0 {
				percent := float64(bytesComplete) / float64(totalBytes)
				select {
				case m.progressChan <- percent:
				default:
				}
			}
		}

		modelFile, err = downloader.DownloadModel(downloadCtx, ttsapp.DefaultTTSModelURL, progressCallback)
		if err != nil {
			// TTS is optional, signal completion anyway
			m.progressChan <- 1.0
			return
		}
	}

	// Signal completion
	m.progressChan <- 1.0
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
	whisperProg := progress.New(progress.WithDefaultGradient())
	sdProg := progress.New(progress.WithDefaultGradient())
	ttsProg := progress.New(progress.WithDefaultGradient())

	// Check if wallet already exists
	walletExists := walletService.HasWallet()

	// Create context with cancel for cleanup
	ctx, cancel := context.WithCancel(context.Background())

	progressChan := make(chan float64, 100)

	model := setupTUIModel{
		styles:             s,
		step:               stepWelcome,
		walletService:      walletService,
		walletExists:       walletExists,
		passwordInput:      passwordInput,
		confirmInput:       confirmInput,
		nameInput:          nameInput,
		mnemonicInput:      mnemonicInput,
		keystoreInput:      keystoreInput,
		keystorePathInput:  keystorePathInput,
		privateKeyInput:    privateKeyInput,
		spinner:            sp,
		progressBar:        prog,
		whisperProgressBar: whisperProg,
		sdProgressBar:      sdProg,
		ttsProgressBar:     ttsProg,
		skipHardwareCheck:  skipHardwareCheck,
		progressChan:       progressChan,
		ctx:                ctx,
		cancel:             cancel,
		result: &SetupResult{
			Errors: make([]error, 0),
		},
	}

	p := tea.NewProgram(model, tea.WithAltScreen())
	finalModel, err := p.Run()

	// Cleanup
	cancel()
	close(progressChan)

	if err != nil {
		return nil, err
	}

	m := finalModel.(setupTUIModel)
	return m.result, nil
}
