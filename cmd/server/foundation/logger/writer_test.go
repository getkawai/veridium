package logger

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWriter_CreatesFile(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	w := NewWriter(logPath)
	require.NotNil(t, w)

	// Write test data
	n, err := w.Write([]byte("test log\n"))
	require.NoError(t, err)
	assert.Equal(t, 9, n)

	// Verify file was created and contains data
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	assert.Equal(t, "test log\n", string(content))
}

func TestNewWriter_MultiWriter(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	w := NewWriter(logPath)
	require.NotNil(t, w)

	// Write test data (goes to both stdout and file)
	n, err := w.Write([]byte("test log\n"))
	require.NoError(t, err)
	assert.Equal(t, 9, n)

	// Verify file was created and contains data
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	assert.Equal(t, "test log\n", string(content))
}

func TestNewWriter_DirectoryCreation(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "nested", "dir", "test.log")

	w := NewWriter(logPath)
	require.NotNil(t, w)

	// Write test data
	n, err := w.Write([]byte("test\n"))
	require.NoError(t, err)
	assert.Equal(t, 5, n)

	// Verify directory was created
	_, err = os.Stat(filepath.Dir(logPath))
	require.NoError(t, err)

	// Verify file content
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	assert.Equal(t, "test\n", string(content))
}

func TestNewWriter_MultipleWrites(t *testing.T) {
	tmpDir := t.TempDir()
	logPath := filepath.Join(tmpDir, "test.log")

	w := NewWriter(logPath)
	require.NotNil(t, w)

	// Write multiple times
	w.Write([]byte("line 1\n"))
	w.Write([]byte("line 2\n"))
	w.Write([]byte("line 3\n"))

	// Verify all lines are in file
	content, err := os.ReadFile(logPath)
	require.NoError(t, err)
	assert.Equal(t, "line 1\nline 2\nline 3\n", string(content))
}
