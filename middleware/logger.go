package middleware

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/studiolambda/akumu"
)

// Logger middleware sets a [slog.Logger] instance
// as the logger for any http requests.
func Logger(logger *slog.Logger) akumu.Middleware {
	return func(handler http.Handler) http.Handler {
		return LoggerWith(handler, logger)
	}
}

// LoggerDefault middleware sets the [slog.Default] instance
// as the logger for any http requests.
func LoggerDefault() akumu.Middleware {
	return func(handler http.Handler) http.Handler {
		return LoggerWith(handler, slog.Default())
	}
}

// LoggerWith middleware sets a [slog.Logger] instance
// as the logger for any http requests but this time accepting
// the handler as a parameter.
func LoggerWith(handler http.Handler, logger *slog.Logger) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		parent, hasParent := request.Context().Value(akumu.OnErrorKey{}).(akumu.OnErrorHook)

		handler.ServeHTTP(writer, request.WithContext(
			context.WithValue(request.Context(), akumu.OnErrorKey{}, func(err akumu.ErrServer) {
				if hasParent && parent != nil {
					parent(err)
				}

				logger.ErrorContext(
					request.Context(),
					"server error",
					"code", err.Code,
					"text", http.StatusText(err.Code),
					"url", err.Request.URL,
				)
			}),
		))
	})
}
