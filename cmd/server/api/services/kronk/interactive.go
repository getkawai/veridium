package kronk

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/tyler-smith/go-bip39"
	"golang.org/x/term"
)

// promptPassword prompts for password input without echoing
func promptPassword(prompt string) (string, error) {
	fmt.Print(prompt)
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}
	return string(bytePassword), nil
}

// promptInput prompts for text input
func promptInput(prompt string) (string, error) {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// promptYesNo prompts for yes/no confirmation
func promptYesNo(prompt string) bool {
	for {
		input, err := promptInput(prompt + " (y/n): ")
		if err != nil {
			return false
		}
		input = strings.ToLower(input)
		if input == "y" || input == "yes" {
			return true
		}
		if input == "n" || input == "no" {
			return false
		}
		fmt.Println("Please enter 'y' or 'n'")
	}
}

// promptChoice prompts for a choice from options
func promptChoice(prompt string, options []string) (int, error) {
	fmt.Println(prompt)
	for i, opt := range options {
		fmt.Printf("  %d. %s\n", i+1, opt)
	}

	for {
		input, err := promptInput("Enter choice (1-" + fmt.Sprint(len(options)) + "): ")
		if err != nil {
			return 0, err
		}

		var choice int
		_, err = fmt.Sscanf(input, "%d", &choice)
		if err == nil && choice >= 1 && choice <= len(options) {
			return choice - 1, nil
		}
		fmt.Printf("Invalid choice. Please enter a number between 1 and %d\n", len(options))
	}
}

// printBanner prints welcome banner
func printBanner() {
	fmt.Println(`
╔═══════════════════════════════════════════════════════════╗
║                                                           ║
║   🌸 Kawai DeAI Network - Contributor Server 🌸          ║
║                                                           ║
║   Earn KAWAI tokens by providing AI compute power        ║
║                                                           ║
╚═══════════════════════════════════════════════════════════╝
`)
}

// printSuccess prints success message
func printSuccess(msg string) {
	fmt.Printf("\n✅ %s\n\n", msg)
}

// printError prints error message
func printError(msg string) {
	fmt.Printf("\n❌ %s\n\n", msg)
}

// printWarning prints warning message
func printWarning(msg string) {
	fmt.Printf("\n⚠️  %s\n\n", msg)
}

// printInfo prints info message
func printInfo(msg string) {
	fmt.Printf("\nℹ️  %s\n\n", msg)
}

// printMnemonic prints mnemonic with warning
func printMnemonic(mnemonic string) {
	fmt.Println("\n╔═══════════════════════════════════════════════════════════╗")
	fmt.Println("║                  ⚠️  SAVE YOUR MNEMONIC ⚠️                 ║")
	fmt.Println("╚═══════════════════════════════════════════════════════════╝")
	fmt.Println()
	fmt.Printf("  %s\n", mnemonic)
	fmt.Println()
	fmt.Println("⚠️  IMPORTANT:")
	fmt.Println("  • Write these words down on paper")
	fmt.Println("  • Store in a secure location")
	fmt.Println("  • NEVER share with anyone")
	fmt.Println("  • Anyone with these words can access your funds")
	fmt.Println()
}

// validatePassword validates password strength
func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters")
	}
	return nil
}

// validateMnemonic validates mnemonic phrase using BIP39
func validateMnemonic(mnemonic string) error {
	mnemonic = strings.TrimSpace(mnemonic)
	if !bip39.IsMnemonicValid(mnemonic) {
		return fmt.Errorf("invalid mnemonic phrase (BIP39 validation failed)")
	}
	return nil
}
