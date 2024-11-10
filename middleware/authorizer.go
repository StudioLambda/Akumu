package middleware

import (
	"net/http"

	"github.com/studiolambda/akumu"
)

func Authorizer(authorizer akumu.Authorizer) akumu.Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if err := akumu.AuthorizeWith(request, authorizer); err != nil {
				akumu.
					Failed(err).
					Handle(writer, request)

				return
			}

			handler.ServeHTTP(writer, request)
		})
	}
}

func AuthorizerFunc(authorizer akumu.AuthorizerFunc) akumu.Middleware {
	return Authorizer(authorizer)
}
