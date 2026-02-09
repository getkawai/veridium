//go:build darwin

package download

// LibraryName returns the name for the stable-diffusion.cpp library on macOS.
func LibraryName() string {
	return "libstable-diffusion.dylib"
}
