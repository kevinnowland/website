package pkg

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

type StatusRecorder struct {
	http.ResponseWriter
	StatusCode int
}

func (w *StatusRecorder) WriteHeader(status int) {
	w.StatusCode = status
	w.ResponseWriter.WriteHeader(status)
}

func NewLoggingHandler(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(
				"Handling",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			ctx := context.WithValue(r.Context(), "StartTime", time.Now())

			recorder := &StatusRecorder{
				ResponseWriter: w,
				StatusCode:     200,
			}
			next.ServeHTTP(recorder, r)

			t := time.Now()
			elapsed := t.Sub(ctx.Value("StartTime").(time.Time))

			logger.Info(
				"FinishedHandling",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int64("duration_ns", elapsed.Nanoseconds()),
				slog.Int("status_code", recorder.StatusCode),
			)
		})
	}
}
