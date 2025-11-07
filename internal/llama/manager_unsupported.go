//go:build !darwin && !linux && !windows

package llama

import "fmt"

// InstallLlamaCpp returns an error on unsupported platforms
func (lcm *LlamaCppReleaseManager) InstallLlamaCpp() error {
	return fmt.Errorf("llama.cpp installation not supported on this platform")
}
