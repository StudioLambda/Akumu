package akumu

import "net/http"

type Handler func(*http.Request) error

func (handler Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	HandleError(writer, request, handler(request), http.StatusInternalServerError)
}
