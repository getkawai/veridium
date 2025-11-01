package main

import (
	"os"

	"github.com/kawai-network/veridium/pkg/nodefs"
)

// NodeFsService provides Node.js fs module equivalents as a Wails service
type NodeFsService struct{}

// FileExistsSync checks if a file exists synchronously
func (n *NodeFsService) FileExistsSync(path string) bool {
	return nodefs.FileExistsSync(path)
}

// ReadFileSync reads a file synchronously and returns as string
func (n *NodeFsService) ReadFileSync(path string) (string, error) {
	data, err := nodefs.ReadFileSync(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ReadFile reads a file (synchronously in Go)
func (n *NodeFsService) ReadFile(path string) (string, error) {
	data, err := nodefs.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// WriteFileSync writes data to a file synchronously
func (n *NodeFsService) WriteFileSync(path string, content string) error {
	return nodefs.WriteFileSync(path, []byte(content))
}

// WriteFile writes data to a file
func (n *NodeFsService) WriteFile(path string, content string) error {
	return nodefs.WriteFile(path, []byte(content))
}

// ReadDirSync reads directory contents synchronously
func (n *NodeFsService) ReadDirSync(path string) ([]string, error) {
	entries, err := nodefs.ReadDirSync(path)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(entries))
	for i, entry := range entries {
		names[i] = entry.Name()
	}
	return names, nil
}

// StatSync gets file stats synchronously
func (n *NodeFsService) StatSync(path string) (map[string]interface{}, error) {
	info, err := nodefs.StatSync(path)
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"name":    info.Name(),
		"size":    info.Size(),
		"mode":    info.Mode().String(),
		"modTime": info.ModTime().Unix(),
		"isDir":   info.IsDir(),
	}, nil
}

// MkdirSync creates a directory synchronously
func (n *NodeFsService) MkdirSync(path string) error {
	return nodefs.MkdirSync(path)
}

// RmSync removes a file or directory recursively
func (n *NodeFsService) RmSync(path string) error {
	return nodefs.RmSync(path)
}

// UnlinkSync removes a file
func (n *NodeFsService) UnlinkSync(path string) error {
	return nodefs.UnlinkSync(path)
}

// RenameSync renames a file
func (n *NodeFsService) RenameSync(oldPath, newPath string) error {
	return nodefs.RenameSync(oldPath, newPath)
}

// CopyFileSync copies a file
func (n *NodeFsService) CopyFileSync(src, dest string) error {
	return nodefs.CopyFileSync(src, dest)
}

// AppendFileSync appends data to a file
func (n *NodeFsService) AppendFileSync(path string, content string) error {
	return nodefs.AppendFileSync(path, []byte(content))
}

// TruncateSync truncates a file to specified length
func (n *NodeFsService) TruncateSync(path string, length int64) error {
	return nodefs.TruncateSync(path, length)
}

// ChmodSync changes file permissions
func (n *NodeFsService) ChmodSync(path string, mode uint32) error {
	return nodefs.ChmodSync(path, os.FileMode(mode))
}

// AccessSync checks file access permissions
func (n *NodeFsService) AccessSync(path string, mode int) error {
	return nodefs.AccessSync(path, os.FileMode(mode))
}

// MkdtempSync creates a temporary directory synchronously
func (n *NodeFsService) MkdtempSync(prefix string) (string, error) {
	return nodefs.MkdtempSync(prefix)
}
