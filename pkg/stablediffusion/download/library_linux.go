//go:build linux

package download

// LibraryName returns the name for the stable-diffusion.cpp library on Linux.
func LibraryName() string {
	return "libstable-diffusion.so"
}
