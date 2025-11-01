package main

import (
	"github.com/kawai-network/veridium/pkg/nodepath"
)

// NodePathService provides Node.js path module equivalents as a Wails service
type NodePathService struct{}

// Basename returns the last portion of a path
func (n *NodePathService) Basename(path string, ext ...string) string {
	if len(ext) > 0 {
		return nodepath.Basename(path, ext[0])
	}
	return nodepath.Basename(path)
}

// Dirname returns the directory name of a path
func (n *NodePathService) Dirname(path string) string {
	return nodepath.Dirname(path)
}

// Extname returns the extension of the path
func (n *NodePathService) Extname(path string) string {
	return nodepath.Extname(path)
}

// Parse parses a path into its components
func (n *NodePathService) Parse(path string) map[string]string {
	parsed := nodepath.Parse(path)
	return map[string]string{
		"root": parsed.Root,
		"dir":  parsed.Dir,
		"base": parsed.Base,
		"ext":  parsed.Ext,
		"name": parsed.Name,
		"full": parsed.Full,
	}
}

// Format formats a path object into a path string
func (n *NodePathService) Format(root, dir, base, ext, name string) string {
	obj := nodepath.PathObject{
		Root: root,
		Dir:  dir,
		Base: base,
		Ext:  ext,
		Name: name,
	}
	return nodepath.Format(obj)
}

// IsAbsolute determines whether path is an absolute path
func (n *NodePathService) IsAbsolute(path string) bool {
	return nodepath.IsAbsolute(path)
}

// Join joins path segments using the platform-specific separator
func (n *NodePathService) Join(paths ...string) string {
	return nodepath.Join(paths...)
}

// Resolve resolves a sequence of path segments into an absolute path
func (n *NodePathService) Resolve(paths ...string) string {
	return nodepath.Resolve(paths...)
}

// Relative returns the relative path from from to to
func (n *NodePathService) Relative(from, to string) (string, error) {
	return nodepath.Relative(from, to), nil
}

// Normalize normalizes the given path
func (n *NodePathService) Normalize(path string) string {
	return nodepath.Normalize(path)
}

// ToNamespacedPath converts a path to a namespaced path (Windows specific)
func (n *NodePathService) ToNamespacedPath(path string) string {
	return nodepath.ToNamespacedPath(path)
}

// Matches returns true if the pattern matches the path
func (n *NodePathService) Matches(pattern, path string) bool {
	return nodepath.Matches(pattern, path)
}

// Glob performs glob pattern matching
func (n *NodePathService) Glob(pattern string) ([]string, error) {
	return nodepath.Glob(pattern)
}

// Get constants
func (n *NodePathService) Sep() string {
	return nodepath.Sep
}

func (n *NodePathService) Delimiter() string {
	return nodepath.Delimiter
}
