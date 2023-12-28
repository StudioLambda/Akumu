package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/studiolambda/golidate"
	"github.com/studiolambda/golidate/format"
	"github.com/studiolambda/golidate/translate/language"
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
	ErrRequestValidate           = errors.New("request validation failed")
	ErrValidationFailed          = errors.New("validation failed")
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
		return NewError(errors.Join(ErrRequestValidate, err)).
			Status(StatusBadRequest)
	}

	if results := dest.Validate(request.ctx); !results.PassesAll() {
		fields := results.
			Failed().
			Translate(language.English).
			Group().
			Messages(format.Capitalize())

		return NewError(ErrValidationFailed).
			Status(StatusUnprocessableEntity).
			Fields(fields)
	}

	return nil
}

func (request Request) Validate(dest golidate.Validator) error {
	if request.Headers().Contains("Content-Type", "application/json") {
		if err := request.ValidateJSON(dest); err != nil {
			return err
		}

		return nil
	}

	err := errors.Join(
		ErrRequestValidate,
		fmt.Errorf("%w: %s", ErrRequestUnknownContentType, request.Headers().First("Content-Type")),
	)

	return NewError(err).
		Status(StatusBadRequest)
}
