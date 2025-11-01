package main

import (
	"fmt"
	"log"
	"time"

	"github.com/kawai-network/veridium/pkg/nodeexec"
)

func main() {
	fmt.Println("=== Node.js child_process equivalents ===")

	// Example 1: Simple command execution
	fmt.Println("\n1. Simple command execution:")
	result, err := nodeexec.ExecSync("echo 'Hello from Go!'")
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Success: %v\n", result.Success)
		fmt.Printf("Output: %s", result.Stdout)
		fmt.Printf("Exit code: %d\n", result.Code)
	}

	// Example 2: Command with options
	fmt.Println("\n2. Command with options:")
	result2, err := nodeexec.ExecSync("ls -la", &nodeexec.ExecOptions{
		Cwd:     ".", // Current directory
		Timeout: 5 * time.Second,
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Success: %v\n", result2.Success)
		fmt.Printf("Output length: %d bytes\n", len(result2.Stdout))
	}

	// Example 3: Environment variables
	fmt.Println("\n3. Command with custom environment:")
	result3, err := nodeexec.ExecSync("env | grep -E '(USER|HOME|SHELL)' | head -3", &nodeexec.ExecOptions{
		Env: map[string]string{
			"CUSTOM_VAR": "Hello from custom env",
			"PATH":       "/usr/local/bin:/usr/bin:/bin", // Simplified PATH
		},
	})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Environment output:\n%s", result3.Stdout)
	}

	// Example 4: Spawn process with arguments
	fmt.Println("\n4. Spawn process with arguments:")
	result4, err := nodeexec.SpawnSync("echo", []string{"Spawned", "process", "with", "args"})
	if err != nil {
		log.Printf("Error: %v", err)
	} else {
		fmt.Printf("Spawn result: %s", result4.Stdout)
	}

	// Example 5: Asynchronous execution
	fmt.Println("\n5. Asynchronous execution:")
	resultChan, err := nodeexec.Exec("sleep 1 && echo 'Async execution completed'")
	if err != nil {
		log.Printf("Error starting async exec: %v", err)
	} else {
		fmt.Println("Async command started, waiting for result...")
		asyncResult := <-resultChan
		fmt.Printf("Async result: %s", asyncResult.Stdout)
	}

	// Example 6: Which command (find executable)
	fmt.Println("\n6. Finding executables:")
	executables := []string{"go", "node", "python3", "ls"}
	for _, exe := range executables {
		path, err := nodeexec.Which(exe)
		if err != nil {
			fmt.Printf("%s: not found\n", exe)
		} else {
			fmt.Printf("%s: %s\n", exe, path)
		}
	}

	// Example 7: Error handling
	fmt.Println("\n7. Error handling:")
	result7, err := nodeexec.ExecSync("nonexistent_command_that_should_fail")
	if err != nil {
		fmt.Printf("Expected error: %v\n", err)
		fmt.Printf("Result success: %v\n", result7.Success)
		fmt.Printf("Exit code: %d\n", result7.Code)
	}

	// Example 8: Timeout handling
	fmt.Println("\n8. Timeout handling:")
	result8, err := nodeexec.ExecSync("sleep 10", &nodeexec.ExecOptions{
		Timeout: 2 * time.Second, // 2 second timeout
	})
	if err != nil {
		fmt.Printf("Timeout error (expected): %v\n", err)
		fmt.Printf("Command timed out: %v\n", result8.Code != 0)
	}

	fmt.Println("\n=== Examples completed ===")
}
