package pkg

import (
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	var logLevel slog.Level
	switch l := os.Getenv("WEBSITE_LOG_LEVEL"); l {
	case "DEBUG":
		logLevel = slog.LevelDebug
	case "INFO":
		logLevel = slog.LevelInfo
	case "WARN":
		logLevel = slog.LevelWarn
	case "ERROR":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: false,
	}
	handler := (slog.NewJSONHandler(os.Stdout, opts))
	logger := slog.New(handler)

	return logger
}
