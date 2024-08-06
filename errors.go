package akumu

import "fmt"

type ServerErrorKind string

type ServerError struct {
	Code int
	URL  string
	Text string
	Kind ServerErrorKind
	Body string
}

var (
	ServerErrorWriter  ServerErrorKind = "writer"
	ServerErrorBody    ServerErrorKind = "body"
	ServerErrorStream  ServerErrorKind = "stream"
	ServerErrorDefault ServerErrorKind = "default"
)

func (err ServerError) Error() string {
	return fmt.Sprintf(
		"%d: %s",
		err.Code,
		err.Text,
	)
}
