//go:build windows

package download

// LibraryName returns the name for the stable-diffusion.cpp library on Windows.
func LibraryName() string {
	return "libgosd-fallback.dll"
}
