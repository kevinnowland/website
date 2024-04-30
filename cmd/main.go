package main

import (
	"log/slog"
	"net/http"
	"os"

	"website/pkg"
)

func main() {
	logger := getLogger()
	loggingHandler := pkg.NewLoggingHandler(logger)

	mux := http.NewServeMux()

	healthzHandler := http.HandlerFunc(pkg.Healthz)
	mux.Handle("/healthz", loggingHandler(healthzHandler))

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/", loggingHandler(fs))

	logger.Info("Listening", slog.Int("port", 8080))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
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

	opts := &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: false,
	}
	handler := (slog.NewJSONHandler(os.Stdout, opts))
	logger := slog.New(handler)

	return logger
}
