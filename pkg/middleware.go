package pkg

import (
	"context"
	"log/slog"
	"net/http"
	"time"
)

func NewLoggingHandler(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info(
				"Handling",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
			)

			ctx := context.WithValue(r.Context(), "StartTime", time.Now())

			next.ServeHTTP(w, r)

			t := time.Now()
			elapsed := t.Sub(ctx.Value("StartTime").(time.Time))

			logger.Info(
				"FinishedHAndling",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int64("duration_ns", elapsed.Nanoseconds()),
			)
		})
	}
}
