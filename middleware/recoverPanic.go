package middleware

import (
	"log/slog"
	"net/http"

	"github.com/proteatech/bok"
)

// RecoverPanic is a middleware that recovers from panics
func RecoverPanic() bok.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					slog.Error("🚨 Recovered from panic",
						"panic", r,
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next(w, r)
		}
	}
}
