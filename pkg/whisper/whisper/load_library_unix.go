//go:build darwin || linux

package whisper

import (
	"github.com/ebitengine/purego"
)

// openLibrary opens a dynamic library on Unix platforms
func openLibrary(name string) (uintptr, error) {
	return purego.Dlopen(name, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}

// closeLibrary closes a dynamic library on Unix platforms
func closeLibrary(handle uintptr) error {
	return purego.Dlclose(handle)
}

// registerLibFunc registers a library function using purego (Unix)
func registerLibFunc(fn interface{}, lib uintptr, name string) {
	purego.RegisterLibFunc(fn, lib, name)
}
