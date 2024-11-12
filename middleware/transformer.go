package middleware

import (
	"net/http"

	"github.com/studiolambda/akumu"
)

func Transformer(transformer akumu.Transformer) akumu.Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			req, err := akumu.TransformWith(request, transformer)

			if err != nil {
				akumu.
					Failed(err).
					Handle(writer, request)

				return
			}

			handler.ServeHTTP(writer, req)
		})
	}
}

func TransformerFunc(transformer akumu.TransformerFunc) akumu.Middleware {
	return Transformer(transformer)
}
