package main

import (
	"log/slog"
	"net/http"
)

func main() {
	logger := NewLogger()
	loggingHandler := NewLoggingHandler(logger)

	mux := http.NewServeMux()

	healthzHandler := http.HandlerFunc(Healthz)
	mux.Handle("/healthz", loggingHandler(healthzHandler))

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", loggingHandler(http.StripPrefix("/static/", fs)))

	indexHandler := http.HandlerFunc(Index)
	mux.Handle("/{$}", loggingHandler(indexHandler))

	genericTemplateHandler := http.HandlerFunc(GenericTemplate)
	mux.Handle("/", loggingHandler(genericTemplateHandler))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	logger.Info("Listening", slog.Int("port", 8080))
	err := server.ListenAndServe()
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
}
