package xlog

import (
	"context"
	"log/slog"
	"os"
	"strings"
)

var logger *slog.Logger

func init() {
	var level slog.Level
	logLevelStr := os.Getenv("LOG_LEVEL")

	// Parse string log level (debug, info, warn, error)
	switch strings.ToLower(logLevelStr) {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "warn", "warning":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		// Default to Debug if not specified or invalid (legacy behavior seemed to default to Debug in previous code logic if conversion failed,
		// though logic was: var level = slog.LevelDebug; if v, err := strconv.Atoi(...); err == nil { level = slog.Level(v) }
		// We will stick to LevelDebug as default to match previous behavior for now, or LevelInfo is usually safer standard.
		// Previous code default was LevelDebug.
		level = slog.LevelDebug
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}
	handler := slog.NewJSONHandler(os.Stdout, opts)
	logger = slog.New(handler)
}

// Info logs at LevelInfo.
func Info(msg string, args ...any) {
	logger.Log(context.Background(), slog.LevelInfo, msg, args...)
}

// InfoContext logs at LevelInfo with context.
func InfoContext(ctx context.Context, msg string, args ...any) {
	logger.Log(ctx, slog.LevelInfo, msg, args...)
}

// Debug logs at LevelDebug.
func Debug(msg string, args ...any) {
	logger.Log(context.Background(), slog.LevelDebug, msg, args...)
}

// DebugContext logs at LevelDebug with context.
func DebugContext(ctx context.Context, msg string, args ...any) {
	logger.Log(ctx, slog.LevelDebug, msg, args...)
}

// Warn logs at LevelWarn.
func Warn(msg string, args ...any) {
	logger.Log(context.Background(), slog.LevelWarn, msg, args...)
}

// WarnContext logs at LevelWarn with context.
func WarnContext(ctx context.Context, msg string, args ...any) {
	logger.Log(ctx, slog.LevelWarn, msg, args...)
}

// Error logs at LevelError.
func Error(msg string, args ...any) {
	logger.Log(context.Background(), slog.LevelError, msg, args...)
}

// ErrorContext logs at LevelError with context.
func ErrorContext(ctx context.Context, msg string, args ...any) {
	logger.Log(ctx, slog.LevelError, msg, args...)
}
