package middleware

import (
	"context"
	"net/http"

	"github.com/studiolambda/akumu"
)

// Problems sets the given [akumu.ProblemControls] to the [http.Request].
//
// The [akumu.Problem] will respect those controls when a problem is handled.
func Problems(controls akumu.ProblemControls) akumu.Middleware {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			handler.ServeHTTP(writer, request.WithContext(
				context.WithValue(request.Context(), akumu.ProblemsKey{}, controls),
			))
		})
	}
}
