// Package nodefs provides Node.js fs module equivalents for Go
package nodefs

import (
	"os"
	"path/filepath"
)

// FileExistsSync checks if a file exists synchronously
// Equivalent to: fs.existsSync(path)
func FileExistsSync(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// ReadFileSync reads a file synchronously
// Equivalent to: fs.readFileSync(path, 'utf8')
func ReadFileSync(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// ReadFile reads a file asynchronously (but synchronously in Go)
// Equivalent to: fs/promises.readFile(path)
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// WriteFileSync writes data to a file synchronously
// Equivalent to: fs.writeFileSync(path, data)
func WriteFileSync(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// WriteFile writes data to a file (synchronously in Go)
// Equivalent to: fs/promises.writeFile(path, data)
func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}

// ReadDirSync reads directory contents synchronously
// Equivalent to: fs.readdirSync(path)
func ReadDirSync(path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

// StatSync gets file stats synchronously
// Equivalent to: fs.statSync(path)
func StatSync(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// Stat gets file stats
// Equivalent to: fs/promises.stat(path)
func Stat(path string) (os.FileInfo, error) {
	return os.Stat(path)
}

// MkdtempSync creates a temporary directory synchronously
// Equivalent to: fs.mkdtempSync(prefix)
func MkdtempSync(prefix string) (string, error) {
	return os.MkdirTemp("", prefix)
}

// RmSync removes a file or directory recursively
// Equivalent to: fs.rmSync(path, { recursive: true })
func RmSync(path string) error {
	return os.RemoveAll(path)
}

// MkdirSync creates a directory synchronously
// Equivalent to: fs.mkdirSync(path, { recursive: true })
func MkdirSync(path string) error {
	return os.MkdirAll(path, 0755)
}

// Constants for file modes (similar to Node.js constants)
const (
	// F_OK - File exists
	F_OK os.FileMode = 0
	// R_OK - Readable
	R_OK os.FileMode = 4
	// W_OK - Writable
	W_OK os.FileMode = 2
	// X_OK - Executable
	X_OK os.FileMode = 1
)

// AccessSync checks file access permissions
// Equivalent to: fs.accessSync(path, mode)
func AccessSync(path string, mode os.FileMode) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	file.Close()

	// For more detailed permission checking, we'd need to use syscall
	// This is a simplified version
	return nil
}

// UnlinkSync removes a file
// Equivalent to: fs.unlinkSync(path)
func UnlinkSync(path string) error {
	return os.Remove(path)
}

// RenameSync renames a file
// Equivalent to: fs.renameSync(oldPath, newPath)
func RenameSync(oldPath, newPath string) error {
	return os.Rename(oldPath, newPath)
}

// CopyFileSync copies a file
// Equivalent to: fs.copyFileSync(src, dest)
func CopyFileSync(src, dest string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = destFile.ReadFrom(sourceFile)
	return err
}

// AppendFileSync appends data to a file
// Equivalent to: fs.appendFileSync(path, data)
func AppendFileSync(path string, data []byte) error {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	return err
}

// TruncateSync truncates a file to specified length
// Equivalent to: fs.truncateSync(path, len)
func TruncateSync(path string, length int64) error {
	return os.Truncate(path, length)
}

// ChmodSync changes file permissions
// Equivalent to: fs.chmodSync(path, mode)
func ChmodSync(path string, mode os.FileMode) error {
	return os.Chmod(path, mode)
}

// LstatSync gets file stats (without following symlinks)
// Equivalent to: fs.lstatSync(path)
func LstatSync(path string) (os.FileInfo, error) {
	return os.Lstat(path)
}

// SymlinkSync creates a symbolic link
// Equivalent to: fs.symlinkSync(target, path)
func SymlinkSync(target, path string) error {
	return os.Symlink(target, path)
}

// ReadlinkSync reads the target of a symbolic link
// Equivalent to: fs.readlinkSync(path)
func ReadlinkSync(path string) (string, error) {
	return os.Readlink(path)
}

// RealpathSync resolves to an absolute path
// Equivalent to: fs.realpathSync(path)
func RealpathSync(path string) (string, error) {
	return filepath.Abs(path)
}

// Constants for open flags (similar to Node.js)
const (
	O_RDONLY int = 0    // Open for reading only
	O_WRONLY int = 1    // Open for writing only
	O_RDWR   int = 2    // Open for reading and writing
	O_APPEND int = 8    // Append mode
	O_CREATE int = 64   // Create if it doesn't exist
	O_EXCL   int = 128  // Exclusive mode
	O_SYNC   int = 4096 // Synchronous I/O
	O_TRUNC  int = 512  // Truncate file
)

// OpenSync opens a file
// Equivalent to: fs.openSync(path, flags, mode)
func OpenSync(path string, flags int, mode os.FileMode) (*os.File, error) {
	var osFlags int

	switch flags & 3 { // Check read/write mode
	case O_RDONLY:
		osFlags = os.O_RDONLY
	case O_WRONLY:
		osFlags = os.O_WRONLY
	case O_RDWR:
		osFlags = os.O_RDWR
	default:
		osFlags = os.O_RDONLY
	}

	if flags&O_CREATE != 0 {
		osFlags |= os.O_CREATE
	}
	if flags&O_APPEND != 0 {
		osFlags |= os.O_APPEND
	}
	if flags&O_EXCL != 0 {
		osFlags |= os.O_EXCL
	}
	if flags&O_TRUNC != 0 {
		osFlags |= os.O_TRUNC
	}
	if flags&O_SYNC != 0 {
		osFlags |= os.O_SYNC
	}

	return os.OpenFile(path, osFlags, mode)
}
