package middleware

import (
	"github.com/ProteaTech/bok"
	"net/http"
)

// StaticAssetsCache is a middleware that sets the cache control header for static assets
func StaticAssetsCache() bok.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Cache-Control", "public, max-age=31536000")
			next(w, r)
		}
	}
}
