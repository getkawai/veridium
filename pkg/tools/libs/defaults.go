package libs

import (
	"path/filepath"

	"github.com/kawai-network/veridium/pkg/tools/defaults"
)

// Path returns the location for the libraries folder.
// If override is provided, it uses that path directly.
// Otherwise, it delegates to defaults.BaseDir to determine the base directory
// (which checks VERIDIUM_BASE_DIR env var, then home directory, then current directory).
func Path(override string) string {
	if override != "" {
		return override
	}

	return filepath.Join(defaults.BaseDir(""), localFolder)
}
