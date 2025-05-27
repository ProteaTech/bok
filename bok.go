package bok

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type appRouter struct {
	*http.ServeMux
	middleware []Middleware
	prefix     string
}

type Router interface {
	http.Handler
	GET(string, http.HandlerFunc)
	POST(string, http.HandlerFunc)
	PUT(string, http.HandlerFunc)
	DELETE(string, http.HandlerFunc)
	PATCH(string, http.HandlerFunc)
	OPTIONS(string, http.HandlerFunc)
	HEAD(string, http.HandlerFunc)
	Handle(string, http.Handler)
	HandleFunc(string, http.HandlerFunc)
	WithMiddleware(middleware ...Middleware) Router
}

func NewRouter() *appRouter {
	var router = &appRouter{
		ServeMux: http.NewServeMux(),
	}
	return router
}

func (r *appRouter) GET(path string, handler http.HandlerFunc) {
	path = getModifiedPath(r.prefix, path)
	r.Handle("GET "+path, handler)
}

func (r *appRouter) POST(path string, handler http.HandlerFunc) {
	path = getModifiedPath(r.prefix, path)
	r.Handle("POST "+path, handler)
}

func (r *appRouter) PUT(path string, handler http.HandlerFunc) {
	path = getModifiedPath(r.prefix, path)
	r.Handle("PUT "+path, handler)
}

func (r *appRouter) DELETE(path string, handler http.HandlerFunc) {
	path = getModifiedPath(r.prefix, path)
	r.Handle("DELETE "+path, handler)
}

func (r *appRouter) PATCH(path string, handler http.HandlerFunc) {
	path = getModifiedPath(r.prefix, path)
	r.Handle("PATCH "+path, handler)
}

func (r *appRouter) OPTIONS(path string, handler http.HandlerFunc) {
	path = getModifiedPath(r.prefix, path)
	r.Handle("OPTIONS "+path, handler)
}

func (r *appRouter) HEAD(path string, handler http.HandlerFunc) {
	path = getModifiedPath(r.prefix, path)
	handler = runMiddleware(r, handler)
	r.Handle("HEAD "+path, handler)
}

func (r *appRouter) Handle(path string, handler http.Handler) {
	path = getModifiedPath(r.prefix, path)
	handler = runMiddleware(r, handler.ServeHTTP)
	r.ServeMux.Handle(path, handler)
}

func (r *appRouter) HandleFunc(path string, handler http.HandlerFunc) {
	path = getModifiedPath(r.prefix, path)
	handler = runMiddleware(r, handler)
	r.ServeMux.HandleFunc(path, handler)
}

// WithMiddleware adds middleware to the router.
// middleware funcs are executed in the order they are passed.
func (r *appRouter) WithMiddleware(middleware ...Middleware) Router {
	var router = NewRouter()
	router.ServeMux = r.ServeMux
	router.prefix = r.prefix
	router.middleware = append(r.middleware, middleware...)
	return router
}

// SetPrefix sets a prefix for all routes in the router.
// For example, if the prefix is "/api", then a route
// registered with the path "/users" will be accessible at "/api/users".
func (r *appRouter) SetPrefix(prefix string) {
	r.prefix = prefix
	if len(prefix) > 0 && prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}
	r.ServeMux.Handle(prefix, http.StripPrefix(prefix, r.ServeMux))
}

// middleware needs to be ran in reverse order so that the first middleware
// passed is the first to be executed.
func runMiddleware(r *appRouter, handler http.HandlerFunc) http.HandlerFunc {
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}
	return handler
}

// getModifiedPath modifies the path such that it includes the prefix
// if one is set.
// If the prefix is empty, it returns the path as is.
// If the prefix does not end with a '/', it appends one,
// and then concatenates it with the path.
func getModifiedPath(prefix, path string) string {
	if prefix == "" {
		return path
	}
	if prefix[len(prefix)-1] != '/' {
		prefix += "/"
	}
	return prefix + path
}
