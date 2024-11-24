package akumu

import (
	"net/http"
)

// Handler is the akumu's equivalent of the [http.Handler].
//
// Is is the function that can take care of a request and
// handle a correct response for it.
type Handler func(*http.Request) error

// handleResponder is a helper that will handle specific [Responder]
// responses. It also takes care of any parent [Builder].
func handleResponder(writer http.ResponseWriter, request *http.Request, parent *Builder, responder Responder) {
	if parent != nil {
		parent.
			Merge(responder.Respond(request)).
			Handle(writer, request)

		return
	}

	responder.
		Respond(request).
		Handle(writer, request)
}

// handleNoError is called whenever there's a response that does not
// contain any error. For example, returning `nil` in a handler.
func handleNoError(writer http.ResponseWriter, request *http.Request, parent *Builder) {
	if parent != nil {
		parent.Handle(writer, request)

		return
	}

	Response(http.StatusOK).Handle(writer, request)
}

// handle takes care of responding to a given request.
func handle(writer http.ResponseWriter, request *http.Request, err error, parent *Builder) {
	if err == nil {
		handleNoError(writer, request, parent)
		return
	}

	if responder, ok := err.(Responder); ok {
		handleResponder(writer, request, parent, responder)
		return
	}

	if parent != nil {
		builder := NewProblem(err, parent.status).
			Respond(request)

		parent.
			Merge(builder).
			Handle(writer, request)

		return
	}

	NewProblem(err, http.StatusInternalServerError).
		Respond(request).
		Handle(writer, request)
}

// ServeHTTP implements the [http.Handler] interface to have
// compatibility with the http package.
func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handle(writer, request, handler(request), nil)
}

// HandlerFunc transforms the [Handler] into an [http.HandlerFunc].
func (handler Handler) HandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(writer, request)
	})
}

// HandlerFunc transforms a [Handler] into an [http.HandlerFunc].
//
// This function is a simple alias of [Handler.HandlerFunc].
func HandlerFunc(handler Handler) http.HandlerFunc {
	return handler.HandlerFunc()
}
