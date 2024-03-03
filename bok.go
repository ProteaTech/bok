package bok

import (
	"net/http"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

type appRouter struct {
	*http.ServeMux
	middleware []Middleware
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
	r.Handle("GET "+path, handler)
}

func (r *appRouter) POST(path string, handler http.HandlerFunc) {
	r.Handle("POST "+path, handler)
}

func (r *appRouter) PUT(path string, handler http.HandlerFunc) {
	r.Handle("PUT "+path, handler)
}

func (r *appRouter) DELETE(path string, handler http.HandlerFunc) {
	r.Handle("DELETE "+path, handler)
}

func (r *appRouter) PATCH(path string, handler http.HandlerFunc) {
	r.Handle("PATCH "+path, handler)
}

func (r *appRouter) OPTIONS(path string, handler http.HandlerFunc) {
	r.Handle("OPTIONS "+path, handler)
}

func (r *appRouter) HEAD(path string, handler http.HandlerFunc) {
	handler = runMiddleware(r, handler)
	r.Handle("HEAD "+path, handler)
}

func (r *appRouter) Handle(path string, handler http.Handler) {
	handler = runMiddleware(r, handler.ServeHTTP)
	r.ServeMux.Handle(path, handler)
}

func (r *appRouter) HandleFunc(path string, handler http.HandlerFunc) {
	handler = runMiddleware(r, handler)
	r.ServeMux.HandleFunc(path, handler)
}

// WithMiddleware adds middleware to the router.
// middleware funcs are executed in the order they are passed.
func (r *appRouter) WithMiddleware(middleware ...Middleware) Router {
	var router = NewRouter()
	router.ServeMux = r.ServeMux
	router.middleware = append(r.middleware, middleware...)
	return router
}

// middleware needs to be ran in reverse order so that the first middleware
// passed is the first to be executed.
func runMiddleware(r *appRouter, handler http.HandlerFunc) http.HandlerFunc {
	for i := len(r.middleware) - 1; i >= 0; i-- {
		handler = r.middleware[i](handler)
	}
	return handler
}
