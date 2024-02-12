package akumu

import "net/http"

type Handler func(*http.Request) error

func Handle(handler Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		response := handler(request)

		if responder, ok := response.(Responder); ok {
			responder.Handle(writer, request)

			return
		}

		http.Error(writer, response.Error(), http.StatusInternalServerError)
	}
}
