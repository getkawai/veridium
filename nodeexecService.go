package main

import (
	"time"

	"github.com/kawai-network/veridium/pkg/nodeexec"
)

// ExecOptions represents options for exec operations (frontend compatible)
type ExecOptions struct {
	Cwd         string            `json:"cwd,omitempty"`
	Env         map[string]string `json:"env,omitempty"`
	Timeout     int               `json:"timeout,omitempty"` // in milliseconds
	MaxBuffer   int               `json:"maxBuffer,omitempty"`
	Shell       string            `json:"shell,omitempty"`
	Uid         int               `json:"uid,omitempty"`
	Gid         int               `json:"gid,omitempty"`
	WindowsHide bool              `json:"windowsHide,omitempty"`
	KillSignal  string            `json:"killSignal,omitempty"`
}

// ExecResult represents the result of an exec operation (frontend compatible)
type ExecResult struct {
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
	Error   string `json:"error,omitempty"`
	Code    int    `json:"code"`
	Signal  string `json:"signal,omitempty"`
	Success bool   `json:"success"`
	Command string `json:"command"`
}

// NodeExecService provides Node.js child_process equivalents as a Wails service
type NodeExecService struct{}

// ExecSync executes a command synchronously
func (n *NodeExecService) ExecSync(command string, options *ExecOptions) ExecResult {
	var opts *nodeexec.ExecOptions
	if options != nil {
		opts = &nodeexec.ExecOptions{
			Cwd:        options.Cwd,
			Env:        options.Env,
			Timeout:    time.Duration(options.Timeout) * time.Millisecond,
			Shell:      options.Shell,
			KillSignal: options.KillSignal,
		}
	}

	result, err := nodeexec.ExecSync(command, opts)

	execResult := ExecResult{
		Stdout:  result.Stdout,
		Stderr:  result.Stderr,
		Code:    result.Code,
		Signal:  result.Signal,
		Success: result.Success,
		Command: result.Command,
	}

	if err != nil {
		execResult.Error = err.Error()
	}

	return execResult
}

// Exec executes a command asynchronously (returns result directly for Wails)
func (n *NodeExecService) Exec(command string, options *ExecOptions) ExecResult {
	// For Wails binding, we execute synchronously since Wails handles async differently
	return n.ExecSync(command, options)
}

// SpawnSync executes a command with arguments synchronously
func (n *NodeExecService) SpawnSync(command string, args []string, options *ExecOptions) ExecResult {
	var opts *nodeexec.ExecOptions
	if options != nil {
		opts = &nodeexec.ExecOptions{
			Cwd:        options.Cwd,
			Env:        options.Env,
			Timeout:    time.Duration(options.Timeout) * time.Millisecond,
			Shell:      options.Shell,
			KillSignal: options.KillSignal,
		}
	}

	result, err := nodeexec.SpawnSync(command, args, opts)

	execResult := ExecResult{
		Stdout:  result.Stdout,
		Stderr:  result.Stderr,
		Code:    result.Code,
		Signal:  result.Signal,
		Success: result.Success,
		Command: result.Command,
	}

	if err != nil {
		execResult.Error = err.Error()
	}

	return execResult
}

// Spawn executes a command with arguments (synchronously for Wails)
func (n *NodeExecService) Spawn(command string, args []string, options *ExecOptions) ExecResult {
	return n.SpawnSync(command, args, options)
}

// Which finds the path to an executable
func (n *NodeExecService) Which(command string) (string, error) {
	return nodeexec.Which(command)
}

// GetDefaultShell returns the default shell for the system
func (n *NodeExecService) GetDefaultShell() string {
	return nodeexec.GetDefaultShell()
}

// KillProcess kills a process by PID
func (n *NodeExecService) KillProcess(pid int, signal ...string) error {
	return nodeexec.KillProcess(pid, signal...)
}

// GetProcessInfo gets basic process information
func (n *NodeExecService) GetProcessInfo(pid int) (map[string]interface{}, error) {
	return nodeexec.GetProcessInfo(pid)
}

// Get constants
func (n *NodeExecService) SIGTERM() string {
	return nodeexec.SIGTERM
}

func (n *NodeExecService) SIGKILL() string {
	return nodeexec.SIGKILL
}

func (n *NodeExecService) SIGHUP() string {
	return nodeexec.SIGHUP
}

func (n *NodeExecService) SIGINT() string {
	return nodeexec.SIGINT
}
