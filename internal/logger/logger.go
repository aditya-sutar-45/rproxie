// Package logger
package logger

import (
	"log/slog"
	"os"
)

func ParseLogLevel(levelString string) slog.Level {
	var level slog.Level

	switch levelString {
	case "DEBUG":
		level = slog.LevelDebug
	case "WARN":
		level = slog.LevelWarn
	case "ERROR":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	return level
}

func New(logLevel string) *slog.Logger {
	level := ParseLogLevel(logLevel)

	opts := &slog.HandlerOptions{
		Level: level,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)

	return slog.New(handler)
}
