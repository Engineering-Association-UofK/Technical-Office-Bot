package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)

// Open endpoint anyone can call
func Basic(next http.Handler) http.Handler {
	return recoveryMiddleware(
		loggingMiddleware(
			next,
		),
	)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		slog.Debug("Logging middleware", "Method", r.Method, "Path", r.URL.Path, "Time took", time.Since(start))
	})
}

func recoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				msg := fmt.Sprintf("Recovered from panic: %v\n", err)
				slog.Error(msg)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}()
		next.ServeHTTP(w, r)
	})
}
