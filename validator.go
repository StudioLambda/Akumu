package akumu

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

type Validator interface {
	Validate(request *http.Request) error
}

type ValidatorFunc func(request *http.Request) error

var (
	ErrValidationFailed           = errors.New("validation failed")
	ErrInvalidValidatorInContext  = errors.New("invalid validator in context")
	ErrValidatorNotFoundInContext = errors.New("validator not found in context")
)

func (validator ValidatorFunc) Validate(request *http.Request) error {
	return validator(request)
}

func Validate[T Validator](request *http.Request) error {
	return ValidateWith(request, *new(T))
}

func ValidateFrom(request *http.Request, key any) error {
	value := request.Context().Value(key)

	if value != nil {
		return fmt.Errorf("%w: %s", ErrValidatorNotFoundInContext, reflect.TypeOf(key))
	}

	validator, ok := value.(Validator)

	if !ok {
		return fmt.Errorf("%w: %s", ErrInvalidValidatorInContext, reflect.TypeOf(key))
	}

	return ValidateWith(request, validator)
}

func ValidateWith[T Validator](request *http.Request, validator T) error {
	if err := validator.Validate(request); err != nil {
		return NewProblem(
			errors.Join(ErrValidationFailed, err),
			http.StatusUnprocessableEntity,
		)
	}

	return nil
}

func ValidateWithFunc(request *http.Request, validator ValidatorFunc) error {
	return ValidateWith(request, validator)
}
