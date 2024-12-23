package akumu

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"path"
	"slices"
	"strings"
)

// Router is the structure that handles
// http routing in an akumu application.
//
// This router is completly optional and
// uses [http.ServeMux] under the hood
// to register all the routes.
//
// It also handles some patterns automatically,
// such as {$}, that is appended on each route
// automatically, regardless of the pattern.
type Router struct {

	// native stores the actual [http.ServeMux]
	// that's used internally  to register the routes.
	native *http.ServeMux

	// pattern stores the current pattern that will be
	// used as a prefix to all the route registrations
	// on this router. This pattern is already joined with
	// the parent router's pattern if any.
	pattern string

	// parent stores the parent [Router] if any. This is
	// used to correctly resolve the [http.ServeMux] to
	// use by sub-routers so that they all register the
	// routes to the same [http.ServeMux].
	parent *Router

	// middlewares stores the actual middlewares that will
	// be applied to any route registration on the current
	// router. It already contains all the middlewares of
	// the parent's [Router] if any.
	middlewares []Middleware
}

// NewRouter creates a new [Router] instance and
// automatically creates all the needed components
// such as the middleware list or the native
// [http.ServeMux] that's used under the hood.
func NewRouter() *Router {
	return &Router{
		native:      http.NewServeMux(),
		pattern:     "",
		parent:      nil,
		middlewares: make([]Middleware, 0),
	}
}

// Group uses the given pattern to automatically
// mount a sub-router that has that pattern as a
// prefix.
//
// This means that any route registered with the
// sub-router will also have the given pattern suffixed.
//
// Keep in mind this can be nested as well, meaning that
// many sub-routers may be grouped, creating complex patterns.
func (router *Router) Group(pattern string, subrouter func(*Router)) {
	subrouter(&Router{
		native:      nil, // parent's native will be used
		pattern:     path.Join(router.pattern, pattern),
		parent:      router,
		middlewares: slices.Clone(router.middlewares),
	})
}

// With does create a new sub-router that automatically applies
// the given middlewares.
//
// This is very usefull when used to inline some middlewares to
// specific routes.
//
// In constrast to [Router.Use] method, it does create a new
// sub-router instead of modifying the current router.
func (router *Router) With(middlewares ...Middleware) *Router {
	return &Router{
		native:      nil, // parent's native will be used
		pattern:     router.pattern,
		parent:      router,
		middlewares: append(slices.Clone(router.middlewares), middlewares...),
	}
}

// mux returns the native [http.ServeMux] that is used
// internally by the router. This exists because sub-routers
// must use the same [http.ServeMux] and therefore, there's
// some recursivity involved to get the same [http.ServeMux].
func (router *Router) mux() *http.ServeMux {
	if router.parent != nil {
		return router.parent.mux()
	}

	return router.native
}

// wrap makes an [http.Handler] wrapped by the current routers'
// middlewares. This means that the resulting [http.Handler] is
// the same as first calling the router middlewares and then the
// provided [http.Handler].
func (router *Router) wrap(handler http.Handler) http.Handler {
	for i := len(router.middlewares) - 1; i >= 0; i-- {
		handler = router.middlewares[i](handler)
	}

	return handler
}

// Use appends to the current router the given middlewares.
//
// Subsequent route registrations will be wrapped with any previous
// middlewares that the router had defined, plus the new ones
// that are registered after this call.
//
// In constrats with the [Router.With] method, this one does modify
// the current router instead of returning a new sub-router.
func (router *Router) Use(middlewares ...Middleware) {
	router.middlewares = append(router.middlewares, middlewares...)
}

// register adds the given pattern and handler to the actual native
// router [http.ServeMux].
func (router *Router) register(pattern string, handler http.Handler) {
	router.
		mux().
		Handle(pattern, handler)
}

// Method registers a new handler to the router with the given
// method and pattern. This is usefull if you need to dynamically
// register a route to the router using a string as the method.
//
// A notable difference is that the patterns's ending slash "/" is
// not treated as an annonymous catch-all "{...}" and is instead treated
// as if it finished with "/{$}", making a specific route only.
//
// If the route does not finish in "/", one will be added automatically and
// then the paragraph above will apply unless the route finishes in a catch-all
// parameter "...}"
//
// Typically, the method string should be one of the following:
//   - [http.MethodGet]
//   - [http.MethodHead]
//   - [http.MethodPost]
//   - [http.MethodPut]
//   - [http.MethodPatch]
//   - [http.MethodDelete]
//   - [http.MethodConnect]
//   - [http.MethodOptions]
//   - [http.MethodTrace]
func (router *Router) Method(method string, pattern string, handler Handler) {
	pattern = path.Join(router.pattern, pattern)

	if pattern == "/" {
		router.register(
			fmt.Sprintf("%s %s{$}", method, pattern),
			router.wrap(handler),
		)

		return
	}

	if !strings.HasSuffix("/", pattern) && !strings.HasSuffix(pattern, "...}") {
		// The redirection route should ONLY be registered when:
		// - The path does not end in "/"
		// - The path does not end in "...}"
		router.register(
			fmt.Sprintf("%s %s/{$}", method, pattern),
			router.redirect(pattern),
		)
	}

	router.register(
		fmt.Sprintf("%s %s", method, pattern),
		router.wrap(handler),
	)
}

// redirect is a helper handler that takes care of redirecting
// to the specified path while maintaining the query string.
func (router *Router) redirect(path string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		target := path
		if r.URL.RawQuery != "" {
			target += "?" + r.URL.RawQuery
		}

		http.Redirect(w, r, target, http.StatusTemporaryRedirect)
	})
}

// generatePattern creates the actual pattern that will be
// registered to the mux handler.
// func (router *Router) generatePattern(method string, pattern string) string {
// 	pattern = path.Join(router.pattern, pattern)

// 	if strings.HasSuffix(pattern, "...}") {
// 		return fmt.Sprintf("%s %s", method, pattern)
// 	}

// 	if !strings.HasSuffix(pattern, "/") {
// 		return fmt.Sprintf("%s %s/{$}", method, pattern)
// 	}

// 	return fmt.Sprintf("%s %s{$}", method, pattern)
// }

// Get registers a new handler to the router using [Router.Method]
// and using the [http.MethodGet] as the method parameter.
func (router *Router) Get(pattern string, handler Handler) {
	router.Method(http.MethodGet, pattern, handler)
}

// Head registers a new handler to the router using [Router.Method]
// and using the [http.MethodHead] as the method parameter.
func (router *Router) Head(pattern string, handler Handler) {
	router.Method(http.MethodHead, pattern, handler)
}

// Post registers a new handler to the router using [Router.Method]
// and using the [http.MethodPost] as the method parameter.
func (router *Router) Post(pattern string, handler Handler) {
	router.Method(http.MethodPost, pattern, handler)
}

// Put registers a new handler to the router using [Router.Method]
// and using the [http.MethodPut] as the method parameter.
func (router *Router) Put(pattern string, handler Handler) {
	router.Method(http.MethodPut, pattern, handler)
}

// Patch registers a new handler to the router using [Router.Method]
// and using the [http.MethodPatch] as the method parameter.
func (router *Router) Patch(pattern string, handler Handler) {
	router.Method(http.MethodPatch, pattern, handler)
}

// Delete registers a new handler to the router using [Router.Method]
// and using the [http.MethodDelete] as the method parameter.
func (router *Router) Delete(pattern string, handler Handler) {
	router.Method(http.MethodDelete, pattern, handler)
}

// Connect registers a new handler to the router using [Router.Method]
// and using the [http.MethodConnect] as the method parameter.
func (router *Router) Connect(pattern string, handler Handler) {
	router.Method(http.MethodConnect, pattern, handler)
}

// Options registers a new handler to the router using [Router.Method]
// and using the [http.MethodOptions] as the method parameter.
func (router *Router) Options(pattern string, handler Handler) {
	router.Method(http.MethodOptions, pattern, handler)
}

// Trace registers a new handler to the router using [Router.Method]
// and using the [http.MethodTrace] as the method parameter.
func (router *Router) Trace(pattern string, handler Handler) {
	router.Method(http.MethodTrace, pattern, handler)
}

// ServeHTTP is the method that will make the router implement
// the [http.Handler] interface, making it possible to be used
// as a handler in places like [http.Server].
func (router *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	router.
		native.
		ServeHTTP(writer, request)
}

// Has reports whether the given pattern is registered in the router
// with the given method.
//
// Alternatively, check out the [Router.Matches] to use an [http.Request]
// as the parameter.
func (router *Router) Has(method string, pattern string) bool {
	if request, err := http.NewRequest(method, pattern, nil); err == nil {
		return router.Matches(request)
	}

	return false
}

// Matches reports whether the given [http.Request] match any registered
// route in the router.
//
// This means that, given the request method and the
// URL, a [Handler] can be resolved.
func (router *Router) Matches(request *http.Request) bool {
	_, ok := router.HandlerMatch(request)

	return ok
}

// Handler returns the [Handler] that matches the given method and pattern.
// The second return value determines if the [Handler] was found or not.
//
// For matching against an [http.Request] use the [Router.HandlerMatch] method.
func (router *Router) Handler(method string, pattern string) (http.Handler, bool) {
	if request, err := http.NewRequest(method, pattern, nil); err == nil {
		return router.HandlerMatch(request)
	}

	return nil, false
}

// HandlerMatch returns the [Handler] that matches the given [http.Request].
// The second return value determines if the [Handler] was found or not.
//
// For matching against a method and a pattern, use the [Router.Handler] method.
func (router *Router) HandlerMatch(request *http.Request) (http.Handler, bool) {
	if handler, pattern := router.native.Handler(request); pattern != "" {
		return handler, true
	}

	return nil, false
}

// Record returns a [httptest.ResponseRecorder] that can be used to inspect what
// the given http request would have returned as a response.
//
// This method is a shortcut of calling [RecordHandler] with the router as the
// [Handler] and the given request.
func (router *Router) Record(request *http.Request) *httptest.ResponseRecorder {
	return RecordHandler(router, request)
}
