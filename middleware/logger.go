package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/studiolambda/akumu"
)

func Logger(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return LoggerWith(handler, logger)
	}
}

func LoggerDefault() func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return LoggerWith(handler, slog.Default())
	}
}

func LoggerWith(handler http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		handler.ServeHTTP(writer, request.WithContext(
			context.WithValue(request.Context(), akumu.LoggerKey{}, logger),
		))
	})
}
