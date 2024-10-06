package akumu

import (
	"net/http"
)

type Handler func(*http.Request) error

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

func handleNoError(writer http.ResponseWriter, request *http.Request, parent *Builder) {
	if parent != nil {
		parent.Handle(writer, request)

		return
	}

	Response(http.StatusOK).Handle(writer, request)
}

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

func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handle(writer, request, handler(request), nil)
}

func (handler Handler) HandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(writer, request)
	})
}

func HandlerFunc(handler func(*http.Request) error) http.HandlerFunc {
	return Handler(handler).HandlerFunc()
}
