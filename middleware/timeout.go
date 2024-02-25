package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/ProteaTech/bok"
)

// Timeout is a middleware that sets a timeout for the request
func Timeout(t time.Duration) bok.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			// Create a new context with a timeout
			ctx, cancel := context.WithTimeout(r.Context(), t)
			defer func() {
				cancel()
				if ctx.Err() == context.DeadlineExceeded {
					w.WriteHeader(http.StatusGatewayTimeout)
				}
			}()
			r = r.WithContext(ctx) // Replace the request with the new context
			next(w, r)
		}
	}
}
