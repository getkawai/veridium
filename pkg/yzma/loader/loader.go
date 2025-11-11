package loader

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/jupiterrider/ffi"
)

// LoadLibrary loads a shared library from the specified path.
// The path should be a directory containing the library files.
// The lib should be the "short name" for the library, for example:
// ggml, ggml-base, llama, mtmd
//
// Example:
//
//	LoadLibrary("/usr/local/lib", "llama") -> loads /usr/local/lib/libllama.dylib (on macOS)
func LoadLibrary(path, lib string) (ffi.Lib, error) {
	if path == "" {
		return ffi.Lib{}, fmt.Errorf("library path cannot be empty")
	}

	var filename string
	switch runtime.GOOS {
	case "linux", "freebsd":
		filename = filepath.Join(path, fmt.Sprintf("lib%s.so", lib))
	case "windows":
		filename = filepath.Join(path, fmt.Sprintf("%s.dll", lib))
	case "darwin":
		filename = filepath.Join(path, fmt.Sprintf("lib%s.dylib", lib))
	default:
		return ffi.Lib{}, fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	return ffi.Load(filename)
}
