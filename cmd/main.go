package main

import (
	"log/slog"
	"os"
)

func main() {
	logger := getLogger()

	logger.Info("hello, world!")
}

func getLogger() *slog.Logger {
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

	logOpts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: true,
	}
	handler := (slog.NewJSONHandler(os.Stdout, logOpts))
	logger := slog.New(handler)

	return logger
}
