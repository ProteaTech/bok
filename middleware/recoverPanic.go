package middleware

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/ProteaTech/bok"
)

// RecoverPanic is a middleware that recovers from panics
func RecoverPanic() bok.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, req *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					opts := slog.HandlerOptions{
						AddSource: true,
					}
					logger := slog.New(slog.NewJSONHandler(os.Stdout, &opts))
					logger.Error(
						"🚨 Recovered from panic",
						"panic", r,
						"ctxErr", req.Context().Err().Error(),
						"url", req.URL,
						"method", req.Method,
					)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				}
			}()
			next(w, req)
		}
	}
}
