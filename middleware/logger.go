package middleware

import (
	"log/slog"
	"net/http"

	"github.com/ProteaTech/bok"
)

// LogRequest is a middleware that logs the request
func Logger(logger *slog.Logger) bok.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			logger.InfoContext(
				r.Context(),
				"Request",
				"method", r.Method,
				"path", r.URL.RequestURI(),
				"remote", r.RemoteAddr,
				"user_agent", r.UserAgent(),
			)
			next(w, r)
		}
	}
}
