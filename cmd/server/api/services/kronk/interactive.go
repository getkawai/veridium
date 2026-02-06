package kronk

import (
	"bufio"
	"fmt"
	"os"
	"strings"

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

// printInfo prints info message
func printInfo(msg string) {
	fmt.Printf("\nℹ️  %s\n\n", msg)
}
