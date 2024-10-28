package router

import (
	"net/http"
	"net/http/httptest"
	"path"
	"slices"

	"github.com/studiolambda/akumu"
)

type Router struct {
	configuration Configuration
	native        *http.ServeMux
	pattern       string
	parent        *Router
	middlewares   []akumu.Middleware
}

type Configuration struct {
	Exact bool
}

func NewRouter() *Router {
	return NewRouterWith(Configuration{
		Exact: true,
	})
}

func NewRouterWith(configuration Configuration) *Router {
	return &Router{
		configuration: configuration,
		native:        http.NewServeMux(),
		pattern:       "",
		parent:        nil,
		middlewares:   make([]akumu.Middleware, 0),
	}
}

func (router *Router) Group(pattern string, subrouter func(*Router)) {
	subrouter(&Router{
		native:      nil, // parent's native will be used
		pattern:     path.Join(router.pattern, pattern),
		parent:      router,
		middlewares: slices.Clone(router.middlewares),
	})
}

func (router *Router) mux() *http.ServeMux {
	if router.parent != nil {
		return router.parent.mux()
	}

	return router.native
}

func (router *Router) wrap(handler http.Handler) http.Handler {
	for _, middleware := range router.middlewares {
		handler = middleware(handler)
	}

	return handler
}

func (router *Router) Use(middlewares ...akumu.Middleware) {
	router.middlewares = append(router.middlewares, middlewares...)
}

func (router *Router) Method(method string, pattern string, handler akumu.Handler) {
	pattern = method + " " + pattern

	if router.configuration.Exact {
		pattern += "{$}"
	}

	router.mux().Handle(pattern, router.wrap(handler))
}

func (router *Router) Get(pattern string, handler akumu.Handler) {
	router.Method(http.MethodGet, pattern, handler)
}

func (router *Router) Head(pattern string, handler akumu.Handler) {
	router.Method(http.MethodHead, pattern, handler)
}

func (router *Router) Post(pattern string, handler akumu.Handler) {
	router.Method(http.MethodPost, pattern, handler)
}

func (router *Router) Put(pattern string, handler akumu.Handler) {
	router.Method(http.MethodPut, pattern, handler)
}

func (router *Router) Patch(pattern string, handler akumu.Handler) {
	router.Method(http.MethodPatch, pattern, handler)
}

func (router *Router) Delete(pattern string, handler akumu.Handler) {
	router.Method(http.MethodDelete, pattern, handler)
}

func (router *Router) Connect(pattern string, handler akumu.Handler) {
	router.Method(http.MethodConnect, pattern, handler)
}

func (router *Router) Options(pattern string, handler akumu.Handler) {
	router.Method(http.MethodOptions, pattern, handler)
}

func (router *Router) Trace(pattern string, handler akumu.Handler) {
	router.Method(http.MethodTrace, pattern, handler)
}

func (router *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	router.native.ServeHTTP(writer, request)
}

func (router *Router) Has(method string, pattern string) bool {
	if request, err := http.NewRequest(method, pattern, nil); err == nil {
		return router.Matches(request)
	}

	return false
}

func (router *Router) Matches(request *http.Request) bool {
	_, ok := router.HandlerMatch(request)

	return ok
}

func (router *Router) Handler(method string, pattern string) (akumu.Handler, bool) {
	if request, err := http.NewRequest(method, pattern, nil); err == nil {
		return router.HandlerMatch(request)
	}

	return nil, false
}

func (router *Router) HandlerMatch(request *http.Request) (akumu.Handler, bool) {
	if handler, pattern := router.native.Handler(request); pattern != "" {
		if handler, ok := handler.(akumu.Handler); ok {
			return handler, true
		}
	}

	return nil, false
}

func (router *Router) Record(request *http.Request) *httptest.ResponseRecorder {
	return akumu.RecordHandler(router, request)
}
