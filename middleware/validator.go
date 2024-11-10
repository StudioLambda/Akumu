package middleware

import (
	"net/http"

	"github.com/studiolambda/akumu"
)

func Validator(validator akumu.Validator) akumu.Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if err := akumu.ValidateWith(request, validator); err != nil {
				akumu.
					Failed(err).
					Handle(writer, request)

				return
			}

			handler.ServeHTTP(writer, request)
		})
	}
}

func ValidatorFunc(validator akumu.ValidatorFunc) akumu.Middleware {
	return Validator(validator)
}
