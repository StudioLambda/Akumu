package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/studiolambda/akumu"
)

var (
	ErrUnexpected = errors.New("an unexpected error occurred")
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

		return ErrUnexpected
	})
}

func RecoverWith(handler http.Handler, handle func(value any) error) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				akumu.
					Failed(handle(err)).
					Handle(writer, request)
			}
		}()

		handler.ServeHTTP(writer, request)
	})
}
