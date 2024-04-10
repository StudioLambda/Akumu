package akumu

import "net/http"

type Handler func(*http.Request) error

func handleError(writer http.ResponseWriter, request *http.Request, err error, parent *Builder) {
	if err == nil {
		if parent != nil {
			parent.Handle(writer, request)

			return
		}

		Response(http.StatusOK).Handle(writer, request)

		return
	}

	if builder, ok := err.(Builder); ok {
		if parent != nil {
			builder.Merge(*parent).Handle(writer, request)

			return
		}

		builder.Handle(writer, request)

		return
	}

	if responder, ok := err.(Responder); ok {
		if parent != nil {
			responder.
				Respond(request).
				Merge(*parent).
				Handle(writer, request)

			return
		}

		responder.
			Respond(request).
			Handle(writer, request)

		return
	}

	if parent != nil {
		NewProblem(err, http.StatusInternalServerError).
			Respond(request).
			Merge(*parent).
			Handle(writer, request)

		return
	}

	NewProblem(err, http.StatusInternalServerError).
		Respond(request).
		Handle(writer, request)
}

func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	handleError(writer, request, handler(request), nil)
}

func (handler Handler) HandlerFunc() http.HandlerFunc {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(writer, request)
	})
}

func HandlerFunc(handler func(*http.Request) error) http.HandlerFunc {
	return Handler(handler).HandlerFunc()
}
