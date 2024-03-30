package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/studiolambda/akumu"
)

func Recover(handler http.Handler) http.Handler {
	return RecoverWith(handler, func(value any) error {
		switch err := (value).(type) {
		case error:
			return err
		case string:
			return errors.New(err)
		case fmt.Stringer:
			return errors.New(err.String())
		}

		return errors.New("an unexpected error occurred")
	})
}

func RecoverWith(handler http.Handler, handle func(value any) error) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				akumu.HandleError(writer, request, handle(err), http.StatusInternalServerError)
			}
		}()

		handler.ServeHTTP(writer, request)
	})
}
