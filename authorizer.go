package akumu

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

type Authorizer interface {
	Authorize(request *http.Request) error
}

type AuthorizerFunc func(request *http.Request) error

var (
	ErrAuthorizationFailed         = errors.New("authorization failed")
	ErrInvalidAuthorizerInContext  = errors.New("invalid authorizer in context")
	ErrAuthorizerNotFoundInContext = errors.New("authorizer not found in context")
)

func (authorizer AuthorizerFunc) Authorize(request *http.Request) error {
	return authorizer(request)
}

func Authorize[T Authorizer](request *http.Request) error {
	return AuthorizeWith(request, *new(T))
}

func AuthorizeFrom(request *http.Request, key any) error {
	value := request.Context().Value(key)

	if value != nil {
		return fmt.Errorf("%w: %s", ErrAuthorizerNotFoundInContext, reflect.TypeOf(key))
	}

	authorizer, ok := value.(Authorizer)

	if !ok {
		return fmt.Errorf("%w: %s", ErrInvalidAuthorizerInContext, reflect.TypeOf(key))
	}

	return AuthorizeWith(request, authorizer)
}

func AuthorizeWith[T Authorizer](request *http.Request, authorizer T) error {
	if err := authorizer.Authorize(request); err != nil {
		return NewProblem(
			errors.Join(ErrAuthorizationFailed, err),
			http.StatusForbidden,
		)
	}

	return nil
}

func AuthorizeWithFunc(request *http.Request, authorizer AuthorizerFunc) error {
	return AuthorizeWith(request, authorizer)
}
