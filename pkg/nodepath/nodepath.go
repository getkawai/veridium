// Package nodepath provides Node.js path module equivalents for Go
package nodepath

import (
	"path/filepath"
	"strings"
)

// PathObject represents a parsed path object
type PathObject struct {
	Root string
	Dir  string
	Base string
	Ext  string
	Name string
	Full string
}

// Parse parses a path into its components
// Equivalent to: path.parse(path)
func Parse(path string) PathObject {
	obj := PathObject{
		Full: path,
	}

	obj.Dir = filepath.Dir(path)
	obj.Base = filepath.Base(path)
	obj.Ext = filepath.Ext(path)
	obj.Name = strings.TrimSuffix(obj.Base, obj.Ext)

	// For Windows, root would be like "C:" or "\\server\share"
	// For Unix-like systems, root is "/"
	if len(path) > 0 {
		if path[0] == '/' || path[0] == '\\' {
			obj.Root = string(path[0])
		} else if len(path) >= 3 && path[1] == ':' && (path[2] == '/' || path[2] == '\\') {
			// Windows drive letter
			obj.Root = path[:3]
		}
	}

	return obj
}

// Basename returns the last portion of a path
// Equivalent to: path.basename(path[, ext])
func Basename(path string, ext ...string) string {
	base := filepath.Base(path)
	if len(ext) > 0 {
		return strings.TrimSuffix(base, ext[0])
	}
	return base
}

// Dirname returns the directory name of a path
// Equivalent to: path.dirname(path)
func Dirname(path string) string {
	return filepath.Dir(path)
}

// Extname returns the extension of the path
// Equivalent to: path.extname(path)
func Extname(path string) string {
	return filepath.Ext(path)
}

// Format formats a path object into a path string
// Equivalent to: path.format(pathObject)
func Format(obj PathObject) string {
	if obj.Full != "" {
		return obj.Full
	}

	result := obj.Base
	if obj.Dir != "" && obj.Dir != "." {
		result = filepath.Join(obj.Dir, result)
	}
	return result
}

// IsAbsolute determines whether path is an absolute path
// Equivalent to: path.isAbsolute(path)
func IsAbsolute(path string) bool {
	return filepath.IsAbs(path)
}

// Join joins path segments using the platform-specific separator
// Equivalent to: path.join([...paths])
func Join(paths ...string) string {
	return filepath.Join(paths...)
}

// Resolve resolves a sequence of path segments into an absolute path
// Equivalent to: path.resolve([...paths])
func Resolve(paths ...string) string {
	if len(paths) == 0 {
		cwd, _ := filepath.Abs(".")
		return cwd
	}

	var result string
	for _, path := range paths {
		if IsAbsolute(path) {
			result = path
		} else if result == "" {
			result = path
		} else {
			result = filepath.Join(result, path)
		}
	}

	if result == "" {
		cwd, _ := filepath.Abs(".")
		return cwd
	}

	abs, _ := filepath.Abs(result)
	return abs
}

// Relative returns the relative path from from to to
// Equivalent to: path.relative(from, to)
func Relative(from, to string) string {
	rel, err := filepath.Rel(from, to)
	if err != nil {
		return ""
	}
	return rel
}

// Normalize normalizes the given path
// Equivalent to: path.normalize(path)
func Normalize(path string) string {
	return filepath.Clean(path)
}

// Sep is the platform-specific path segment separator
// Equivalent to: path.sep
var Sep = string(filepath.Separator)

// Delimiter is the platform-specific path delimiter
// Equivalent to: path.delimiter (for PATH environment variable)
var Delimiter = string(filepath.ListSeparator)

// Posix provides POSIX-specific path methods
var Posix = &posixPath{}

type posixPath struct{}

func (p *posixPath) Join(paths ...string) string {
	return filepath.Join(paths...)
}

func (p *posixPath) Resolve(paths ...string) string {
	return Resolve(paths...)
}

func (p *posixPath) Normalize(path string) string {
	return filepath.Clean(path)
}

func (p *posixPath) IsAbsolute(path string) bool {
	return strings.HasPrefix(path, "/")
}

func (p *posixPath) Basename(path string) string {
	return Basename(path)
}

func (p *posixPath) Extname(path string) string {
	return Extname(path)
}

func (p *posixPath) Dirname(path string) string {
	return Dirname(path)
}

func (p *posixPath) Sep() string {
	return "/"
}

func (p *posixPath) Delimiter() string {
	return ":"
}

// Win32 provides Windows-specific path methods
var Win32 = &win32Path{}

type win32Path struct{}

func (w *win32Path) Join(paths ...string) string {
	return filepath.Join(paths...)
}

func (w *win32Path) Resolve(paths ...string) string {
	return Resolve(paths...)
}

func (w *win32Path) Normalize(path string) string {
	return filepath.Clean(path)
}

func (w *win32Path) IsAbsolute(path string) bool {
	return filepath.IsAbs(path)
}

func (w *win32Path) Basename(path string) string {
	return Basename(path)
}

func (w *win32Path) Extname(path string) string {
	return Extname(path)
}

func (w *win32Path) Dirname(path string) string {
	return Dirname(path)
}

func (w *win32Path) Sep() string {
	return "\\"
}

func (w *win32Path) Delimiter() string {
	return ";"
}

// ToNamespacedPath converts a path to a namespaced path (Windows specific)
// Equivalent to: path.toNamespacedPath(path)
func ToNamespacedPath(path string) string {
	// On non-Windows systems, this is a no-op
	return path
}

// Matches returns true if the pattern matches the path
// This is a simplified glob matching - for full glob support, use a dedicated library
func Matches(pattern, path string) bool {
	matched, err := filepath.Match(pattern, path)
	return err == nil && matched
}

// Glob performs glob pattern matching
// Note: This is a simplified version. For production use, consider using github.com/gobwas/glob
func Glob(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}
	return matches, nil
}
