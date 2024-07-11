package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/studiolambda/akumu"
)

func Logger(handler http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(writer, request.WithContext(
			context.WithValue(request.Context(), akumu.LoggerKey{}, logger),
		))
	})
}
