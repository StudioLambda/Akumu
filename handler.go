package akumu

import "net/http"

type Handler func(*http.Request) error

func handleError(writer http.ResponseWriter, request *http.Request, err error) {
	if err == nil {
		Response(http.StatusOK).Handle(writer, request)

		return
	}

	if builder, ok := err.(Builder); ok {
		builder.Handle(writer, request)

		return
	}

	if responder, ok := err.(Responder); ok {
		responder.
			Respond(request).
			Handle(writer, request)

		return
	}

	NewProblem(err, http.StatusInternalServerError).
		Respond(request).
		Handle(writer, request)
}

func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handleError(writer, request, handler(request))
}

func HandlerFunc(handler Handler) http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handleError(writer, request, handler(request))
	})
}
