package bok

import (
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
)

// statusResponseWriter is the same as before…
type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *statusResponseWriter) Write(b []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	return w.ResponseWriter.Write(b)
}

// Middleware wraps a http.HandlerFunc.
type Middleware func(http.HandlerFunc) http.HandlerFunc

// Router is our public interface.
type Router interface {
	http.Handler
	GET(path string, handler http.HandlerFunc)
	POST(path string, handler http.HandlerFunc)
	PUT(path string, handler http.HandlerFunc)
	DELETE(path string, handler http.HandlerFunc)
	PATCH(path string, handler http.HandlerFunc)
	OPTIONS(path string, handler http.HandlerFunc)
	HEAD(path string, handler http.HandlerFunc)
	WithMiddleware(middleware ...Middleware) Router
	SetPrefix(prefix string)
	SetLogger(logger *slog.Logger)
}

// appRouter is our implementation.
type appRouter struct {
	mux        *http.ServeMux
	prefix     string
	logger     *slog.Logger // optional, can be nil
	middleware []Middleware
}

// NewRouter constructs an empty router.
func NewRouter() *appRouter {
	return &appRouter{
		mux:    http.NewServeMux(),
		prefix: "",
		logger: nil,
	}
}

// ServeHTTP satisfies http.Handler.
func (r *appRouter) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	sw := &statusResponseWriter{ResponseWriter: w}

	r.mux.ServeHTTP(sw, req)

	if r.logger != nil {

		// if the request resulted in error, use logger.ErrorContext instead of
		// InfoContext
		if sw.statusCode >= 400 {
			r.logger.ErrorContext(
				req.Context(),
				"Request",
				"method", req.Method,
				"path", req.URL.RequestURI(),
				"remote", req.RemoteAddr,
				"status", sw.statusCode,
				"user_agent", req.UserAgent(),
			)
		} else {

			r.logger.InfoContext(
				req.Context(),
				"Request",
				"method", req.Method,
				"path", req.URL.RequestURI(),
				"remote", req.RemoteAddr,
				"status", sw.statusCode,
				"user_agent", req.UserAgent(),
			)
		}
	}
}

// SetLogger sets an optional logger for the router.
// If nil, no requests will be logged.
func (r *appRouter) SetLogger(logger *slog.Logger) {
	r.logger = logger
}

// SetPrefix just records the prefix; no mux.Handle inside here.
func (r *appRouter) SetPrefix(prefix string) {
	if prefix == "" {
		r.prefix = ""
		return
	}
	// Ensure leading slash, no trailing slash
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	r.prefix = strings.TrimRight(prefix, "/")
}

// WithMiddleware returns a shallow copy with extra middleware.
func (r *appRouter) WithMiddleware(mw ...Middleware) Router {
	// copy slice to avoid mutation
	all := slices.Clone(r.middleware)
	all = append(all, mw...)

	return &appRouter{
		mux:        r.mux,
		prefix:     r.prefix,
		middleware: all,
		logger:     r.logger,
	}
}

// helper to register a method-guarded route
func (r *appRouter) handle(method, path string, handler http.HandlerFunc) {
	full := joinPath(method, r.prefix, path)

	// apply middleware in reverse so first in slice is outermost
	h := applyMiddleware(r.middleware, handler)

	r.mux.HandleFunc(full, func(w http.ResponseWriter, req *http.Request) {
		if req.Method != method {
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
			return
		}
		h(w, req)
	})
}

func (r *appRouter) GET(path string, handler http.HandlerFunc) {
	r.handle(http.MethodGet, path, handler)
}
func (r *appRouter) POST(path string, handler http.HandlerFunc) {
	r.handle(http.MethodPost, path, handler)
}
func (r *appRouter) PUT(path string, handler http.HandlerFunc) {
	r.handle(http.MethodPut, path, handler)
}
func (r *appRouter) DELETE(path string, handler http.HandlerFunc) {
	r.handle(http.MethodDelete, path, handler)
}
func (r *appRouter) PATCH(path string, handler http.HandlerFunc) {
	r.handle(http.MethodPatch, path, handler)
}
func (r *appRouter) OPTIONS(path string, handler http.HandlerFunc) {
	r.handle(http.MethodOptions, path, handler)
}
func (r *appRouter) HEAD(path string, handler http.HandlerFunc) {
	r.handle(http.MethodHead, path, handler)
}

// joinPath safely concatenates Method + prefix + route, ensuring single slashes.
func joinPath(method, prefix, route string) string {

	if prefix != "" {
		// ensure the prefix starts with a slash
		if !strings.HasPrefix(prefix, "/") {
			prefix = "/" + prefix
		}

		// ensure the prefix does not end with a slash
		if strings.HasSuffix(prefix, "/") {
			prefix = strings.TrimSuffix(prefix, "/")
		}
	}

	// ensure leading slash on route
	if !strings.HasPrefix(route, "/") {
		route = "/" + route
	}

	// ensure no trailing slash on route
	if strings.HasSuffix(route, "/") {
		route = strings.TrimSuffix(route, "/")
	}

	// combine the method, prefix, and route
	return fmt.Sprintf("%s %s%s", method, prefix, route)
}

// applyMiddleware chains all given middleware around the handler.
func applyMiddleware(mw []Middleware, handler http.HandlerFunc) http.HandlerFunc {
	for i := len(mw) - 1; i >= 0; i-- {
		handler = mw[i](handler)
	}
	return handler
}
