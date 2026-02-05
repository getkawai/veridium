package local

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

// DefaultExecutor implements CommandExecutor using os/exec
type DefaultExecutor struct {
	binDir string // Binary directory for library path setup
}

// NewDefaultExecutor creates a new default command executor
func NewDefaultExecutor(binDir string) *DefaultExecutor {
	return &DefaultExecutor{
		binDir: binDir,
	}
}

// Run executes a command with proper environment setup
func (e *DefaultExecutor) Run(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)

	// Setup environment for dynamic libraries
	env := e.buildEnvironment()
	cmd.Env = env

	// Capture output
	var stderr bytes.Buffer
	var stdout bytes.Buffer
	cmd.Stderr = &stderr
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("command execution failed: %w, stderr: %s", err, stderr.String())
	}

	return nil
}

// buildEnvironment creates environment variables for dynamic library loading
// Inherits current environment and prepends binary directory to library path
func (e *DefaultExecutor) buildEnvironment() []string {
	// Start with current environment to preserve PATH, HOME, TEMP, etc.
	env := os.Environ()

	// Prepend binary directory to library path based on OS
	// Important: Prepend instead of overwrite to preserve existing library paths
	switch runtime.GOOS {
	case "darwin":
		existingDyld := os.Getenv("DYLD_LIBRARY_PATH")
		if existingDyld != "" {
			env = append(env, fmt.Sprintf("DYLD_LIBRARY_PATH=%s%c%s", e.binDir, os.PathListSeparator, existingDyld))
		} else {
			env = append(env, fmt.Sprintf("DYLD_LIBRARY_PATH=%s", e.binDir))
		}
	case "linux":
		existingLd := os.Getenv("LD_LIBRARY_PATH")
		if existingLd != "" {
			env = append(env, fmt.Sprintf("LD_LIBRARY_PATH=%s%c%s", e.binDir, os.PathListSeparator, existingLd))
		} else {
			env = append(env, fmt.Sprintf("LD_LIBRARY_PATH=%s", e.binDir))
		}
	case "windows":
		// On Windows, prepend to PATH (not using %PATH% which doesn't expand)
		existingPath := os.Getenv("PATH")
		if existingPath != "" {
			env = append(env, fmt.Sprintf("PATH=%s%c%s", e.binDir, os.PathListSeparator, existingPath))
		} else {
			env = append(env, fmt.Sprintf("PATH=%s", e.binDir))
		}
	}

	return env
}
