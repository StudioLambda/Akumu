package middleware

import (
	"context"
	"net/http"

	"github.com/studiolambda/akumu"
)

func Problems(controls akumu.ProblemControls) akumu.Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			handler.ServeHTTP(writer, request.WithContext(
				context.WithValue(request.Context(), akumu.ProblemsKey{}, controls),
			))
		})
	}
}
