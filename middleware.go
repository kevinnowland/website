// All of the Middleware functions assume that we are using a context
// that has a "LogAttrs" key with value of type []slog.Attr
//
// For the most part we assume that anything we save in the context we want to log
package main

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type StatusRecorder struct {
	http.ResponseWriter
	StatusCode *int
}

func (w StatusRecorder) WriteHeader(status int) {
	w.StatusCode = &status
	w.ResponseWriter.WriteHeader(status)
}

type Middleware func(http.Handler) http.Handler

type MiddlewareChain []Middleware

func ApplyMiddlewareChain(handler http.Handler, chain MiddlewareChain) http.Handler {
	length := len(chain)
	if length == 0 {
		return handler
	}

	return ApplyMiddlewareChain(chain[length-1](handler), chain[:length-1])
}

func RequestInfoMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		attrs := ctx.Value(LogAttrs).([]slog.Attr)
		attrs = append(attrs, slog.String("method", r.Method))
		attrs = append(attrs, slog.String("path", r.URL.Path))
		ctx = context.WithValue(ctx, LogAttrs, attrs)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func StatusRecorderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		recorder := StatusRecorder{w, nil}
		next.ServeHTTP(recorder, r)
	})
}

func DurationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		elapsed := time.Now().Sub(start)

		ctx := r.Context()
		attrs := ctx.Value(LogAttrs).([]slog.Attr)
		attrs = append(attrs, slog.Int64("duration_ns", elapsed.Nanoseconds()))
		ctx = context.WithValue(ctx, LogAttrs, attrs)
		req := r.WithContext(ctx)
		*r = *req
	})
}

// TODO: pull out the status handler
func NewLoggingMiddleware(logger *slog.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			requestId := ctx.Value(RequestId).(uuid.UUID)
			attrs := ctx.Value(LogAttrs).([]slog.Attr)
			attrs = append(attrs, slog.String("requestId", requestId.String()))
			logger.LogAttrs(
				r.Context(),
				slog.LevelInfo,
				"Handling",
				attrs...,
			)

			ctx = context.WithValue(ctx, LogAttrs, attrs)
			req := r.WithContext(ctx)
			*r = *req
			next.ServeHTTP(w, r)

			ctx = r.Context()
			attrs = ctx.Value(LogAttrs).([]slog.Attr)
			logger.LogAttrs(
				r.Context(),
				slog.LevelInfo,
				"FinishedHandling",
				attrs...,
			)
		})
	}
}
