package main

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/studiolambda/akumu"
	"github.com/studiolambda/akumu/middleware"
)

// Define specific Problem errors that are easy to use
// and maintain. The framework handles very well.

var ErrInvalidCookieValue = akumu.Problem{
	Title:  "invalid cookie value",
	Detail: "the cookie existed but contained an invalid value",
	Status: http.StatusBadRequest,
}

func helloWorld(request *http.Request) error {

	// Responses can be any error, but using the built-in
	// response builder will quickly get the job done
	// in most situations.

	cookie, err := request.Cookie("foo")

	if err != nil {

		// Wraping errors with Failed creates a new response
		// builder, allowing appending or tweaking how a response
		// behaves.

		return akumu.
			Failed(err).
			Status(http.StatusConflict)
	}

	if cookie.Value != "hello-word" {

		// Returning Problem errors allows for quick and easy
		// errors that are not only mantainable but also use the RFC:
		// https://datatracker.ietf.org/doc/html/rfc9457

		return akumu.Failed(ErrInvalidCookieValue)
	}

	// The Builder contain many useful methods for quickly building
	// http responses, such as appending headers, cookies or transforming
	// body data into specific content types (HTML, JSON, etc).
	//
	// Not only does it help transforming but also use the apropiate
	// headers for the responses.
	//
	// For example, you can use HTTP Streaming with the SSE or Stream
	// methods and the HTTP headers will be set for you (although you can
	// always modify them after calling the method).

	return akumu.
		Response(http.StatusOK).
		Text("Hello, World")
}

func notify(handler http.Handler) http.Handler {

	// Implement any middleware, or use existing middlewares
	// from the github.com/studiolambda/akumu/middlewares pkg.

	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		slog.Info("new request has been made", "url", request.URL)

		handler.ServeHTTP(writer, request)
	})
}

func main() {

	// Create a new HTTP router or use any existing
	// routers such as http.ServeMux or Chi router.

	router := akumu.NewRouter()

	// Calling Use appends a new middleware to the current
	// router stack. Any route after this one, will
	// pass through this middleware.

	router.Use(middleware.Recover())

	// Registering a route is a simple operation that just requires
	// the pattern and the handler to register.

	router.Get("/", helloWorld)

	// Route gruping can be done quite easy by just calling the Group
	// method and providing a prefix to use ("" is also valid). The provided
	// sub-router inherits the middleware and prefix from the previous router.

	router.Group("/prefix", func(router *akumu.Router) {

		// This middleware will only be used by this sub-router, the previous
		// router will not get affected.
		//
		// In particular, the Logger middleware allows using an slog.Logger instance
		// to quickly log any error on the server error range [500, 600).

		router.Use(middleware.LoggerDefault())

		// The With method can be used to create in-line routers that are perfect
		// for using when middlewares need to be applied directly to the handlers.
		//
		// In this case, this route will be /prefix/foo and will use the 3 middlewares
		// in the chain (recover, logger and notify) in this order:
		// recover -> logger -> notify -> helloWorld.

		router.
			With(notify).
			Post("/foo", helloWorld)
	})

	if err := http.ListenAndServe(":8080", router); !errors.Is(err, http.ErrServerClosed) {
		slog.Error("failed serving http server", "err", err)
	}
}
