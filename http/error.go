package http

type Fields map[string][]string

type ErrorHeaders interface {
	ErrorHeaders() Headers
}

type ErrorStatus interface {
	ErrorStatus() Status
}

type ErrorFields interface {
	ErrorFields() Fields
}

type Error struct {
	error   error
	status  Status
	headers Headers
	fields  Fields
}

type ErrorResponse struct {
	Message string `json:"message"`
	Fields  Fields `json:"fields"`
}

func NewError(err error) Error {
	return Error{
		error:   err,
		status:  0,
		headers: make(Headers),
		fields:  make(Fields),
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

func (err Error) Field(key string, values ...string) Error {
	err.fields[key] = append(err.fields[key], values...)

	return err
}

func (err Error) Fields(fields Fields) Error {
	err.fields = fields

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

func (err Error) ErrorFields() Fields {
	return err.fields
}

func (fields Fields) Merge(other Fields) Fields {
	for key, values := range other {
		fields[key] = append(fields[key], values...)
	}

	return fields
}

func stackTrace(err error) []error {
	result := make([]error, 0)

	// Unwrap joined errors and ignore the join itself.
	if e, ok := err.(interface {
		Unwrap() []error
	}); ok {
		for _, err := range e.Unwrap() {
			result = append(result, stackTrace(err)...)
		}

		return result
	}

	// We can ignore the wrapped error, as it's contained
	// in the fmt.Errorf string.
	return append(result, err)
}
