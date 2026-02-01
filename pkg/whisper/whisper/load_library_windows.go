//go:build windows

package whisper

import (
	"fmt"
	"reflect"
	"syscall"
	"unsafe"
)

var (
	kernel32         = syscall.NewLazyDLL("kernel32.dll")
	procLoadLibraryW = kernel32.NewProc("LoadLibraryW")
	procFreeLibrary  = kernel32.NewProc("FreeLibrary")
	procGetProcAddr  = kernel32.NewProc("GetProcAddress")
)

// openLibrary opens a dynamic library on Windows using syscall
func openLibrary(name string) (uintptr, error) {
	// Convert to UTF16 for Windows
	namePtr, err := syscall.UTF16PtrFromString(name)
	if err != nil {
		return 0, err
	}

	// Call LoadLibraryW
	handle, _, err := procLoadLibraryW.Call(uintptr(unsafe.Pointer(namePtr)))
	if handle == 0 {
		return 0, fmt.Errorf("failed to load library %s: %v", name, err)
	}
	return handle, nil
}

// closeLibrary closes a dynamic library on Windows
func closeLibrary(handle uintptr) error {
	ret, _, err := procFreeLibrary.Call(handle)
	if ret == 0 {
		return fmt.Errorf("failed to free library: %v", err)
	}
	return nil
}

// getProcAddress gets a function address from a loaded library
func getProcAddress(handle uintptr, name string) (uintptr, error) {
	namePtr, err := syscall.BytePtrFromString(name)
	if err != nil {
		return 0, err
	}

	addr, _, err := procGetProcAddr.Call(handle, uintptr(unsafe.Pointer(namePtr)))
	if addr == 0 {
		return 0, fmt.Errorf("failed to find %s: %v", name, err)
	}
	return addr, nil
}

// registerLibFunc registers a library function using syscall (Windows)
// This is a simplified implementation that uses syscall to call the function
func registerLibFunc(fn interface{}, lib uintptr, name string) {
	// Get function address
	addr, err := getProcAddress(lib, name)
	if err != nil {
		panic(fmt.Sprintf("failed to find %s: %v", name, err))
	}

	// Use reflection to set the function pointer
	// This is a workaround since we can't use purego on Windows
	fnValue := reflect.ValueOf(fn).Elem()
	fnValue.Set(reflect.ValueOf(addr))
}
