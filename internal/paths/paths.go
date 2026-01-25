// Package paths provides centralized data path configuration.
// All application data storage paths should be resolved through this package.
//
// Platform-specific data directories:
//   - macOS:   ~/Library/Application Support/Kawai/
//   - Windows: %APPDATA%\Kawai\
//   - Linux:   ~/.config/Kawai/ (or $XDG_CONFIG_HOME/Kawai/)
//
// Development mode (running from terminal): ./data/
package paths

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	dataDir     string
	initialized bool
	mu          sync.RWMutex
)

// SetDataDir sets custom data directory. Must be called before any path access.
// For development, call SetDataDir("data") early in main().
func SetDataDir(dir string) {
	mu.Lock()
	defer mu.Unlock()
	dataDir = dir
	initialized = false // Reset to re-initialize
}

func ensureInit() {
	mu.Lock()
	defer mu.Unlock()
	if initialized {
		return
	}
	if dataDir == "" {
		if IsPackaged() {
			// Running from packaged app - use platform-specific user data directory
			homeDir, err := os.UserHomeDir()
			if err != nil {
				dataDir = "data" // Fallback
			} else {
				switch runtime.GOOS {
				case "darwin":
					dataDir = filepath.Join(homeDir, "Library", "Application Support", "Kawai")
				case "windows":
					// Prefer APPDATA, fallback to LOCALAPPDATA
					if appData := os.Getenv("APPDATA"); appData != "" {
						dataDir = filepath.Join(appData, "Kawai")
					} else if localAppData := os.Getenv("LOCALAPPDATA"); localAppData != "" {
						dataDir = filepath.Join(localAppData, "Kawai")
					} else {
						dataDir = filepath.Join(homeDir, "AppData", "Roaming", "Kawai")
					}
				case "linux":
					// Follow XDG Base Directory specification
					if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
						dataDir = filepath.Join(xdgConfig, "Kawai")
					} else {
						dataDir = filepath.Join(homeDir, ".config", "Kawai")
					}
				default:
					// Other Unix-like systems
					dataDir = filepath.Join(homeDir, ".config", "Kawai")
				}
			}
		} else {
			// Running from terminal or development - use relative path
			dataDir = "data"
		}
	}
	os.MkdirAll(dataDir, 0755)
	initialized = true
}

// Base returns the base data directory
func Base() string {
	ensureInit()
	return dataDir
}

// IsPackaged returns true if running from a packaged app bundle
func IsPackaged() bool {
	execPath, err := os.Executable()
	if err != nil {
		return false
	}
	execDir := filepath.Dir(execPath)

	switch runtime.GOOS {
	case "darwin":
		// macOS: Check for .app bundle structure
		// Path: /path/to/App.app/Contents/MacOS/binary
		if filepath.Base(filepath.Dir(filepath.Dir(execPath))) == "Contents" {
			return true
		}
	case "windows":
		// Windows: Check for resources directory or Program Files
		if filepath.Base(execDir) == "resources" || filepath.Base(filepath.Dir(execDir)) == "resources" {
			return true
		}
		// Also check if installed in Program Files (case-insensitive)
		execPathLower := strings.ToLower(execPath)
		if strings.Contains(execPathLower, "program files") {
			return true
		}
	case "linux":
		// Linux: Check for resources directory or standard install paths
		if filepath.Base(execDir) == "resources" || filepath.Base(filepath.Dir(execDir)) == "resources" {
			return true
		}
		// Check common Linux install paths (use strings.HasPrefix instead of deprecated filepath.HasPrefix)
		if strings.HasPrefix(execPath, "/usr/") || strings.HasPrefix(execPath, "/opt/") {
			return true
		}
	}

	return false
}

// =============================================================================
// Application-level paths
// =============================================================================

// Database returns path to main SQLite database
func Database() string { return filepath.Join(Base(), "veridium.db") }

// DuckDB returns path to DuckDB database
func DuckDB() string { return filepath.Join(Base(), "duckdb.db") }

// KBAssets returns path to knowledge base assets directory
func KBAssets() string { return filepath.Join(Base(), "kb-assets") }

// FileBase returns path to file uploads directory
func FileBase() string { return filepath.Join(Base(), "files") }

// =============================================================================
// Jarvis-specific paths (blockchain/wallet)
// =============================================================================

// Jarvis returns the base directory for Jarvis data
func Jarvis() string { return filepath.Join(Base(), "jarvis") }

// JarvisKeystores returns path to keystore directory
func JarvisKeystores() string { return filepath.Join(Jarvis(), "keystores") }

// JarvisNetworks returns path to custom networks directory
func JarvisNetworks() string { return filepath.Join(Jarvis(), "networks") }

// JarvisAddressBookDB returns path to DuckDB addressbook database
func JarvisAddressBookDB() string { return filepath.Join(Jarvis(), "addressbook.duckdb") }

// JarvisAddressBookHash returns path to addressbook hash file
func JarvisAddressBookHash() string { return filepath.Join(Jarvis(), "addressbook.hash") }

// JarvisCache returns path to cache file
func JarvisCache() string { return filepath.Join(Jarvis(), "cache.json") }

// JarvisAddressBook returns path to addresses.json
func JarvisAddressBook() string { return filepath.Join(Jarvis(), "addresses.json") }

// JarvisSecrets returns path to secrets.json
func JarvisSecrets() string { return filepath.Join(Jarvis(), "secrets.json") }
