package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/studiolambda/golidate"
)

type Request struct {
	method  Method
	url     URL
	headers Headers
	version Version
	body    Body
	ctx     context.Context
}

var (
	ErrRequestValidateJSON       = errors.New("unable to validate request json")
	ErrRequestValidate           = errors.New("unable to validate request")
	ErrRequestUnknownContentType = errors.New("unknown content type")
)

func NewRequest(request *http.Request) Request {
	return Request{
		method:  Method(request.Method),
		url:     NewURL(request.URL),
		version: Version(request.Proto),
		headers: Headers(request.Header.Clone()),
		body:    NewBody(request.Body),
		ctx:     request.Context(),
	}
}

func (request Request) Context() context.Context {
	return request.ctx
}

func (request Request) Method() Method {
	return request.method
}

func (request Request) URL() URL {
	return request.url
}

func (request Request) Version() Version {
	return request.version
}

func (request Request) Headers() Headers {
	return request.headers
}

func (request Request) Body() Body {
	return request.body
}

func (request Request) JSON(dest any) error {
	return request.body.JSON(dest)
}

func (request Request) ValidateJSON(dest golidate.Validator) error {
	if err := request.JSON(dest); err != nil {
		return errors.Join(ErrRequestValidateJSON, err)
	}

	if results := dest.Validate(request.ctx); !results.PassesAll() {
		err := fmt.Errorf("%w: validation failed", ErrRequestValidateJSON)

		return NewError(err).
			Status(StatusUnprocessableEntity)
	}

	return nil
}

func (request Request) Validate(dest golidate.Validator) error {
	if request.Headers().Contains("Content-Type", "application/json") {
		if err := request.ValidateJSON(dest); err != nil {
			return errors.Join(ErrRequestValidate, err)
		}

		return nil
	}

	return errors.Join(
		ErrRequestValidate,
		fmt.Errorf("%w: %s", ErrRequestUnknownContentType, request.Headers().First("Content-Type")),
	)
}
