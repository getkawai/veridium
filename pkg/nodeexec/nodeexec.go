// Package nodeexec provides Node.js child_process equivalents for Go
package nodeexec

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"time"
)

// ExecResult represents the result of an exec operation
type ExecResult struct {
	Stdout  string
	Stderr  string
	Error   error
	Code    int
	Signal  string
	Success bool
	Command string
}

// ExecOptions represents options for exec operations
type ExecOptions struct {
	Cwd         string            // Current working directory
	Env         map[string]string // Environment variables
	Timeout     time.Duration     // Execution timeout
	MaxBuffer   int               // Maximum buffer size (not used in Go, but kept for compatibility)
	Shell       string            // Shell to use
	Uid         int               // User ID (not supported on Windows)
	Gid         int               // Group ID (not supported on Windows)
	WindowsHide bool              // Hide window on Windows
	KillSignal  string            // Signal to send on timeout/kill
}

// ExecSync executes a command synchronously
// Equivalent to: child_process.execSync(command[, options])
func ExecSync(command string, options ...*ExecOptions) (*ExecResult, error) {
	var opts *ExecOptions
	if len(options) > 0 {
		opts = options[0]
	} else {
		opts = &ExecOptions{}
	}

	// Set default shell if not specified
	shell := opts.Shell
	if shell == "" {
		if isWindows() {
			shell = "cmd"
		} else {
			shell = "/bin/sh"
		}
	}

	// Prepare command execution
	var cmd *exec.Cmd

	if isWindows() && shell == "cmd" {
		cmd = exec.Command("cmd", "/c", command)
	} else {
		cmd = exec.Command(shell, "-c", command)
	}

	// Set working directory
	if opts.Cwd != "" {
		cmd.Dir = opts.Cwd
	}

	// Set environment variables
	if opts.Env != nil {
		cmd.Env = os.Environ() // Start with current environment
		for key, value := range opts.Env {
			cmd.Env = append(cmd.Env, key+"="+value)
		}
	}

	// Set timeout if specified
	var ctx context.Context
	var cancel context.CancelFunc
	if opts.Timeout > 0 {
		ctx, cancel = context.WithTimeout(context.Background(), opts.Timeout)
		cmd = exec.CommandContext(ctx, cmd.Path, cmd.Args[1:]...)
		cmd.Dir = opts.Cwd
		if opts.Env != nil {
			cmd.Env = make([]string, 0, len(opts.Env)+len(os.Environ()))
			cmd.Env = append(cmd.Env, os.Environ()...)
			for key, value := range opts.Env {
				cmd.Env = append(cmd.Env, key+"="+value)
			}
		}
		defer cancel()
	}

	// Execute command and capture output
	output, err := cmd.CombinedOutput()

	result := &ExecResult{
		Command: command,
		Stdout:  string(output),
		Stderr:  "", // CombinedOutput includes both stdout and stderr
		Code:    0,
		Success: true,
	}

	if err != nil {
		result.Error = err
		result.Success = false

		// Try to get exit code
		if exitErr, ok := err.(*exec.ExitError); ok {
			result.Code = exitErr.ExitCode()
		} else {
			result.Code = -1
		}
	}

	return result, err
}

// Exec executes a command asynchronously (returns channel)
// Equivalent to: child_process.exec(command[, options], callback)
func Exec(command string, options ...*ExecOptions) (<-chan *ExecResult, error) {
	resultChan := make(chan *ExecResult, 1)

	go func() {
		defer close(resultChan)
		result, _ := ExecSync(command, options...)
		resultChan <- result
	}()

	return resultChan, nil
}

// SpawnSync executes a command with arguments synchronously
// Equivalent to: child_process.spawnSync(command, args[, options])
func SpawnSync(command string, args []string, options ...*ExecOptions) (*ExecResult, error) {
	var opts *ExecOptions
	if len(options) > 0 {
		opts = options[0]
	} else {
		opts = &ExecOptions{}
	}

	cmd := exec.Command(command, args...)

	// Set working directory
	if opts.Cwd != "" {
		cmd.Dir = opts.Cwd
	}

	// Set environment variables
	if opts.Env != nil {
		cmd.Env = os.Environ()
		for key, value := range opts.Env {
			cmd.Env = append(cmd.Env, key+"="+value)
		}
	}

	// Set timeout if specified
	if opts.Timeout > 0 {
		ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
		defer cancel()
		cmd = exec.CommandContext(ctx, command, args...)
		cmd.Dir = opts.Cwd
		if opts.Env != nil {
			cmd.Env = make([]string, 0, len(opts.Env))
			for key, value := range opts.Env {
				cmd.Env = append(cmd.Env, key+"="+value)
			}
		}
	}

	// Execute command
	output, err := cmd.CombinedOutput()

	result := &ExecResult{
		Command: command + " " + strings.Join(args, " "),
		Stdout:  string(output),
		Stderr:  "",
		Code:    0,
		Success: true,
	}

	if err != nil {
		result.Error = err
		result.Success = false

		if exitErr, ok := err.(*exec.ExitError); ok {
			result.Code = exitErr.ExitCode()
		} else {
			result.Code = -1
		}
	}

	return result, err
}

// Spawn executes a command with arguments asynchronously
// Equivalent to: child_process.spawn(command, args[, options])
func Spawn(command string, args []string, options ...*ExecOptions) (<-chan *ExecResult, error) {
	resultChan := make(chan *ExecResult, 1)

	go func() {
		defer close(resultChan)
		result, _ := SpawnSync(command, args, options...)
		resultChan <- result
	}()

	return resultChan, nil
}

// Fork is not directly supported in Go (no process forking like Node.js)
// This is a stub that returns an error
func Fork(modulePath string, args []string, options ...*ExecOptions) (<-chan *ExecResult, error) {
	resultChan := make(chan *ExecResult, 1)

	go func() {
		defer close(resultChan)
		resultChan <- &ExecResult{
			Command: modulePath,
			Error:   &UnsupportedOperationError{"fork is not supported in Go"},
			Code:    -1,
			Success: false,
		}
	}()

	return resultChan, &UnsupportedOperationError{"fork is not supported in Go"}
}

// UnsupportedOperationError represents operations not supported in Go
type UnsupportedOperationError struct {
	Message string
}

func (e *UnsupportedOperationError) Error() string {
	return e.Message
}

// Helper function to detect Windows
func isWindows() bool {
	return strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") ||
		strings.Contains(strings.ToLower(os.Getenv("GOOS")), "windows")
}

// Which finds the path to an executable
// Equivalent to: which command on Unix systems
func Which(command string) (string, error) {
	return exec.LookPath(command)
}

// Helper function to promisify exec (though Go doesn't need promises)
// This function is provided for compatibility with Node.js patterns
func ExecPromise(command string, options ...*ExecOptions) (*ExecResult, error) {
	return ExecSync(command, options...)
}

// GetDefaultShell returns the default shell for the system
func GetDefaultShell() string {
	if isWindows() {
		return "cmd"
	}
	return "/bin/sh"
}

// KillProcess kills a process by PID
// Note: This is a simplified version and may not work on all systems
func KillProcess(pid int, signal ...string) error {
	sig := "TERM"
	if len(signal) > 0 {
		sig = signal[0]
	}

	// On Unix-like systems
	if !isWindows() {
		cmd := exec.Command("kill", "-"+sig, string(rune(pid+'0')))
		return cmd.Run()
	}

	// On Windows, this is more complex and would require additional libraries
	return &UnsupportedOperationError{"kill is not fully supported on Windows"}
}

// GetProcessInfo gets basic process information
// This is a simplified version
func GetProcessInfo(pid int) (map[string]interface{}, error) {
	// On Unix-like systems
	if !isWindows() {
		result, err := ExecSync("ps -p " + string(rune(pid+'0')) + " -o pid,ppid,cmd")
		if err != nil {
			return nil, err
		}

		return map[string]interface{}{
			"pid":    pid,
			"output": result.Stdout,
		}, nil
	}

	return nil, &UnsupportedOperationError{"process info not supported on Windows"}
}

// Constants for signals (simplified)
const (
	SIGTERM = "TERM"
	SIGKILL = "KILL"
	SIGHUP  = "HUP"
	SIGINT  = "INT"
)
