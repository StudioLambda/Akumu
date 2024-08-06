package middleware

import (
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
		if ctx, ok := akumu.Context(request); ok {
			ctx.OnError(func(err akumu.ServerError, next akumu.OnErrorNext) {
				logger.ErrorContext(
					request.Context(),
					"server error",
					"code", err.Code,
					"text", err.Text,
					"url", err.URL,
					"kind", err.Kind,
				)

				next(err)
			})
		}

		handler.ServeHTTP(writer, request)
	})
}
