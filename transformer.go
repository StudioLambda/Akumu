package akumu

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
)

type Transformer interface {
	Transform(request *http.Request) (*http.Request, error)
}

type TransformerFunc func(request *http.Request) (*http.Request, error)

var (
	ErrTransformFailed              = errors.New("transform failed")
	ErrInvalidTransformerInContext  = errors.New("invalid transformer in context")
	ErrTransformerNotFoundInContext = errors.New("transformer not found in context")
)

func (transformer TransformerFunc) Transform(request *http.Request) (*http.Request, error) {
	return transformer(request)
}

func Transform[T Transformer](request *http.Request) (*http.Request, error) {
	return TransformWith(request, *new(T))
}

func TransformFrom(request *http.Request, key any) (*http.Request, error) {
	value := request.Context().Value(key)

	if value != nil {
		return nil, fmt.Errorf("%w: %s", ErrTransformerNotFoundInContext, reflect.TypeOf(key))
	}

	transformer, ok := value.(Transformer)

	if !ok {
		return nil, fmt.Errorf("%w: %s", ErrInvalidTransformerInContext, reflect.TypeOf(key))
	}

	return TransformWith(request, transformer)
}

func TransformWith[T Transformer](request *http.Request, transformer T) (*http.Request, error) {
	req, err := transformer.Transform(request)

	if err != nil {
		return nil, NewProblem(
			errors.Join(ErrTransformFailed, err),
			http.StatusUnprocessableEntity,
		)
	}

	return req, nil
}

func TransformWithFunc(request *http.Request, transformer TransformerFunc) (*http.Request, error) {
	return TransformWith(request, transformer)
}
