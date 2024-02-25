package middleware

import (
	"log/slog"
	"net/http"

	"github.com/ProteaTech/bok"
)

// LogRequest is a middleware that logs the request
func Logger() bok.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			slog.Info("🚸 Handling request from logging middleware",
				"method", r.Method,
				"path", r.URL.Path,
				"remote", r.RemoteAddr,
				"user-agent", r.UserAgent(),
			)
			next(w, r)
		}
	}
}
