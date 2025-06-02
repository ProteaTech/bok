package middleware

import (
	"log/slog"
	"net/http"

	"github.com/ProteaTech/bok"
)

// statusResponseWriter wraps http.ResponseWriter and
// records the status code for logging.
type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusResponseWriter) Write(b []byte) (int, error) {
	// if WriteHeader wasn’t called explicitly, the default is 200
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

// Logger logs the request *and* the response status code.
// Deprecated: middleware.Logger is deprecated and will be removed in a future
// version.
// Use router.SetLogger instead.
func Logger(logger *slog.Logger) bok.Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(origW http.ResponseWriter, r *http.Request) {
			// Wrap the ResponseWriter
			w := &statusResponseWriter{ResponseWriter: origW}

			// Call the next handler
			next(w, r)

			// Now w.statusCode holds whatever got written
			logger.InfoContext(
				r.Context(),
				"Request",
				"method", r.Method,
				"path", r.URL.RequestURI(),
				"remote", r.RemoteAddr,
				"status", w.statusCode,
				"user_agent", r.UserAgent(),
			)
		}
	}
}
