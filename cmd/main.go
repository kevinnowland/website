package main

import (
	"log/slog"
	"net/http"
	"os"
)

func main() {
	logger := getLogger()
	loggingHandler := newLoggingHandler(logger)

	mux := http.NewServeMux()

	healthzHandler := http.HandlerFunc(healthz)
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

func newLoggingHandler(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(
				"Handling",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			next.ServeHTTP(w, r)

			logger.Info(
				"FinishedHAndling",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)
		})
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
