package logger

import (
	"io"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

// NewWriter creates an io.Writer that writes to both stdout and a rotating log file.
// File logging includes automatic rotation:
// - Max size: 100 MB
// - Max backups: 3 files
// - Max age: 28 days
// - Compression: gzip enabled
func NewWriter(logPath string) io.Writer {
	// Ensure directory exists
	dir := filepath.Dir(logPath)
	os.MkdirAll(dir, 0755)

	fileWriter := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    100, // MB
		MaxBackups: 3,
		MaxAge:     28, // days
		Compress:   true,
	}

	return io.MultiWriter(os.Stdout, fileWriter)
}
