package main

import (
	"log/slog"
	"net/http"
)

func main() {
	logger := NewLogger()
	loggingMiddleware := NewLoggingMiddleware(logger)

	mux := http.NewServeMux()

	healthzHandler := http.HandlerFunc(Healthz)
	mux.Handle("/healthz", loggingMiddleware(healthzHandler))

	indexHandler := http.HandlerFunc(Index)
	mux.Handle("/{$}", loggingMiddleware(indexHandler))

	genericTemplateHandler := http.HandlerFunc(GenericTemplate)
	mux.Handle("/", loggingMiddleware(genericTemplateHandler))

	server := &http.Server{
		Addr:        ":8080",
		Handler:     mux,
		ConnContext: ConnContext,
	}

	logger.Info("Listening", slog.Int("port", 8080))
	err := server.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
}
