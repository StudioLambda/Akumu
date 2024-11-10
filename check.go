package akumu

import (
	"errors"
	"net/http"
)

type Validator interface {
	Validate(request *http.Request) error
}

type Authorizer interface {
	Authorize(request *http.Request) error
}

var (
	ErrValidationFailed    = errors.New("validation failed")
	ErrAuthorizationFailed = errors.New("authorization failed")
)

func Check[T any](request *http.Request) error {
	var checker any = *new(T)

	if authorizer, ok := checker.(Authorizer); ok {
		if err := authorizer.Authorize(request); err != nil {
			return NewProblem(
				errors.Join(ErrAuthorizationFailed, err),
				http.StatusForbidden,
			)
		}
	}

	if validator, ok := checker.(Validator); ok {
		if err := validator.Validate(request); err != nil {
			return NewProblem(
				errors.Join(ErrValidationFailed, err),
				http.StatusUnprocessableEntity,
			)
		}
	}

	return nil
}
