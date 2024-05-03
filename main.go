package main

import (
	"context"
	"log/slog"
	"net"
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
		Addr:    ":8080",
		Handler: mux,
		ConnContext: func(ctx context.Context, c net.Conn) context.Context {
			ctx = context.WithValue(ctx, "Errors", []error{})
			ctx = context.WithValue(ctx, "LogStrings", map[string]string{})
			return ctx
		},
	}

	logger.Info("Listening", slog.Int("port", 8080))
	err := server.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
}
