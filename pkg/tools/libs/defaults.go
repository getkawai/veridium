package libs

import (
	"path/filepath"

	"github.com/kawai-network/veridium/pkg/tools/defaults"
)

// Path returns the location for the libraries folder. It will check the
// KRONK_LIB_PATH env var first and then default to the home directory if
// one can be identified. Last resort it will choose the current directory.
func Path(override string) string {
	if override != "" {
		return override
	}

	return filepath.Join(defaults.BaseDir(""), localFolder)
}
