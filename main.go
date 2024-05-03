package main

import (
	"log/slog"
	"net/http"
)

func main() {
	logger := NewLogger()

	middlewareChain := []Middleware{
		RequestInfoMiddleware,
		NewLoggingMiddleware(logger),
		DurationMiddleware,
	}

	mux := http.NewServeMux()

	healthzHandler := ApplyMiddlewareChain(
		http.HandlerFunc(Healthz),
		middlewareChain,
	)
	mux.Handle("/healthz", healthzHandler)

	indexHandler := ApplyMiddlewareChain(
		http.HandlerFunc(Index),
		middlewareChain,
	)
	mux.Handle("/{$}", indexHandler)

	genericTemplateHandler := ApplyMiddlewareChain(
		http.HandlerFunc(GenericTemplate),
		middlewareChain,
	)
	mux.Handle("/", genericTemplateHandler)

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
