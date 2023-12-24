package http

import "errors"

type ErrorHeaders interface {
	ErrorHeaders() Headers
}

type ErrorStatus interface {
	ErrorStatus() Status
}

type Error struct {
	error   error
	status  Status
	headers Headers
}

func NewError(err error) Error {
	return Error{
		error:   err,
		status:  0,
		headers: make(Headers),
	}
}

func (err Error) Status(status Status) Error {
	err.status = status

	return err
}

func (err Error) Headers(headers Headers) Error {
	err.headers = headers

	return err
}

func (err Error) Header(key, value string) Error {
	err.headers.Insert(key, value)

	return err
}

func (err Error) Error() string {
	if err.error == nil {
		return ""
	}

	return err.error.Error()
}

func (err Error) ErrorHeaders() Headers {
	return err.headers
}

func (err Error) ErrorStatus() Status {
	return err.status
}

func unwrapErrors(err error) []error {
	result := []error{err}

	if e, ok := err.(interface {
		Unwrap() []error
	}); ok {
		for _, err := range e.Unwrap() {
			result = append(result, unwrapErrors(err)...)
		}

		return result
	}

	if other := errors.Unwrap(err); other != nil {
		return append(result, unwrapErrors(other)...)
	}

	return result
}
