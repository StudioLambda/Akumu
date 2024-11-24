package akumu

import (
	"errors"
	"fmt"
	"net/http"
)

// ErrServer is the base error that akumu
// will return whenever there's any error
// in the >=500 - <600 range.
type ErrServer struct {

	// Code determines the response's status code.
	// This is used to understand what response
	// was sent due to that particular error.
	Code int

	// Request stores the actual http request that
	// failed to execute. This is very useful to
	// extract information such as URL or headers
	// to understand more about what caused this error.
	Request *http.Request
}

var (
	// ErrServerWriter determines that the
	// server error comes from executing logic
	// around [Builder.BodyWriter] or derivates
	// that use this method.
	ErrServerWriter = errors.New("builder is a body writer")

	// ErrServerBody determines that the
	// server error comes from executing logic
	// around [Builder.BodyReader] or derivates
	// that use this method.
	ErrServerBody = errors.New("builder is a body reader")

	// ErrServerStream determines that the
	// server error comes from executing logic
	// around [Builder.Stream] or derivates
	// that use this method.
	ErrServerStream = errors.New("builder is a body stream")

	// ErrServerDefault determines that the
	// server error comes from executing logic
	// around having no body, stream nor writer
	// involved in delivering a response.
	ErrServerDefault = errors.New("builder is a default no body")
)

// Error implements the error interface
// for a server error.
func (err ErrServer) Error() string {
	return fmt.Sprintf(
		"%d: %s",
		err.Code,
		http.StatusText(err.Code),
	)
}

// Is determines if the given target error
// is an [ErrServer]. This is used when dealing
// with errors using [errors.Is] and [errors.As].
func (err ErrServer) Is(target error) bool {
	_, ok := target.(ErrServer)

	return ok
}
