// Package paths provides centralized data path configuration.
// All application data storage paths should be resolved through this package.
package paths

import (
	"os"
	"path/filepath"
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
		dataDir = "data" // Default for desktop app
	}
	os.MkdirAll(dataDir, 0755)
	initialized = true
}

// Base returns the base data directory
func Base() string {
	ensureInit()
	return dataDir
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
