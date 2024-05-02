package main

import (
	"log/slog"
	"net/http"

	"website/pkg"
)

func main() {
	logger := pkg.NewLogger()
	loggingHandler := pkg.NewLoggingHandler(logger)

	mux := http.NewServeMux()

	healthzHandler := http.HandlerFunc(pkg.Healthz)
	mux.Handle("/healthz", loggingHandler(healthzHandler))

	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", loggingHandler(http.StripPrefix("/static/", fs)))

	indexHandler := http.HandlerFunc(pkg.Index)
	mux.Handle("/{$}", loggingHandler(indexHandler))

	genericTemplateHandler := http.HandlerFunc(pkg.GenericTemplate)
	mux.Handle("/", loggingHandler(genericTemplateHandler))

	logger.Info("Listening", slog.Int("port", 8080))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		logger.Error(err.Error())
		panic(err)
	}
}
