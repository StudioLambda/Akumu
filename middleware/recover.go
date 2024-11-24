package middleware

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/studiolambda/akumu"
)

var (
	// ErrRecoverUnexpectedError is the default error that's passed to
	// the recover response when the error cannot be determined from the
	// given recover()'s value.
	ErrRecoverUnexpectedError = errors.New("an unexpected error occurred")
)

// Recover recovers any panics during a [akumu.Handler] execution.
//
// If the recover() value is of type error, this is directly passed to
// the [akumu.Failed] method.
//
// If it's of type string, a new error is created
// with that string.
//
// If it implements [fmt.Stringer] then that string
// will be used.
//
// If none matches, [ErrRecoverUnexpectedError] is returned.
func Recover() akumu.Middleware {
	return func(handler http.Handler) http.Handler {
		return RecoverWith(handler, func(value any) error {
			switch err := (value).(type) {
			case error:
				return err
			case string:
				return errors.New(err)
			case fmt.Stringer:
				return errors.New(err.String())
			}

			return ErrRecoverUnexpectedError
		})
	}
}

// RecoverWith allows a handler to decide what to do with the recover() value,
// allowing to customize the error that [akumu.Failed] receives.
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
