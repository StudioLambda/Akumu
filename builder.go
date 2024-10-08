package akumu

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type BuilderHandler func(http.ResponseWriter, *http.Request, Builder)

type Builder struct {
	handler BuilderHandler
	status  int
	headers http.Header
	body    io.Reader
	err     error
	stream  <-chan []byte
	writer  func(writer http.ResponseWriter)
}

var (
	ErrWriterRequiresFlusher = errors.New("response writer requires a flusher")
)

func writeHeaders(writer http.ResponseWriter, builder Builder) bool {
	for key, values := range builder.headers {
		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}

	writer.WriteHeader(builder.status)

	return builder.status >= 500 && builder.status < 600
}

func DefaultResponderHandler(writer http.ResponseWriter, request *http.Request, builder Builder) {
	onError, hasOnError := request.Context().Value(OnErrorKey{}).(OnErrorHook)

	if builder.err != nil {
		parent := builder.WithoutError()
		handle(writer, request, builder.err, &parent)

		return
	}

	if builder.writer != nil {
		if writeHeaders(writer, builder) {
			if hasOnError {
				onError(ServerError{
					Code: builder.status,
					URL:  request.URL.String(),
					Text: http.StatusText(builder.status),
					Kind: ServerErrorWriter,
				})
			}
		}

		builder.writer(writer)
		return
	}

	if builder.body != nil {
		body, err := io.ReadAll(builder.body)

		if err != nil {
			NewProblem(err, http.StatusInternalServerError).
				Respond(request).
				Handle(writer, request)

			return
		}

		if writeHeaders(writer, builder) {
			if hasOnError {
				onError(ServerError{
					Code: builder.status,
					URL:  request.URL.String(),
					Text: http.StatusText(builder.status),
					Kind: ServerErrorBody,
					Body: string(body),
				})
			}
		}

		writer.Write(body)

		return
	}

	if builder.stream != nil {
		flusher, ok := writer.(http.Flusher)

		if !ok {
			NewProblem(ErrWriterRequiresFlusher, http.StatusInternalServerError).
				Respond(request).
				Handle(writer, request)

			return
		}

		if writeHeaders(writer, builder) {
			if hasOnError {
				onError(ServerError{
					Code: builder.status,
					URL:  request.URL.String(),
					Text: http.StatusText(builder.status),
					Kind: ServerErrorStream,
				})
			}
		}

		flusher.Flush()

		for {
			select {
			case <-request.Context().Done():
				return
			case message, ok := <-builder.stream:
				if !ok {
					return
				}

				_, _ = writer.Write(message)

				flusher.Flush()
			}
		}
	}

	if writeHeaders(writer, builder) {
		if hasOnError {
			onError(ServerError{
				Code: builder.status,
				URL:  request.URL.String(),
				Text: http.StatusText(builder.status),
				Kind: ServerErrorDefault,
			})
		}
	}
}

func Response(status int) Builder {
	return Builder{
		handler: DefaultResponderHandler,
		status:  status,
		headers: make(http.Header),
		body:    nil,
		stream:  nil,
	}
}

func Failed(err error) Builder {
	return Response(http.StatusInternalServerError).
		Failed(err)
}

func (builder Builder) Error() string {
	return http.StatusText(builder.status)
}

func (builder Builder) Status(status int) Builder {
	builder.status = status

	return builder
}

func (builder Builder) Headers(headers http.Header) Builder {
	builder.headers = headers

	return builder
}

func (builder Builder) Header(key, value string) Builder {
	builder.headers = builder.headers.Clone()
	builder.headers.Set(key, value)

	return builder
}

func (builder Builder) AppendHeader(key, value string) Builder {
	builder.headers = builder.headers.Clone()
	builder.headers.Add(key, value)

	return builder
}

func (builder Builder) Body(body []byte) Builder {
	return builder.BodyReader(bytes.NewReader(body))
}

func (builder Builder) BodyReader(body io.Reader) Builder {
	builder.body = body

	return builder
}

func (builder Builder) Stream(stream <-chan []byte) Builder {
	builder.stream = stream

	return builder.
		Header("Cache-Control", "no-cache").
		Header("Connection", "keep-alive")
}

func (builder Builder) SSE(stream <-chan []byte) Builder {
	return builder.
		Header("Content-Type", "text/event-stream").
		Stream(stream)
}

func (builder Builder) Cookie(cookie http.Cookie) Builder {
	if c := cookie.String(); c != "" {
		return builder.AppendHeader("Set-Cookie", c)
	}

	return builder
}

func (builder Builder) Failed(err error) Builder {
	builder.err = err

	return builder
}

func (builder Builder) Text(body string) Builder {
	return builder.
		Header("Content-Type", "text/plain").
		Body([]byte(body))
}

func (builder Builder) HTML(html string) Builder {
	return builder.
		Header("Content-Type", "text/html").
		Body([]byte(html))
}

func (builder Builder) JSON(body any) Builder {
	buffer := &bytes.Buffer{}

	if err := json.NewEncoder(buffer).Encode(body); err != nil {
		return builder.
			Status(http.StatusInternalServerError).
			Failed(err)
	}

	return builder.
		Header("Content-Type", "application/json").
		BodyReader(buffer)
}

func (builder Builder) BodyWriter(writer func(writer http.ResponseWriter)) Builder {
	builder.writer = writer

	return builder
}

func (builder Builder) Handler(handler BuilderHandler) Builder {
	builder.handler = handler

	return builder
}

func (builder Builder) Handle(response http.ResponseWriter, request *http.Request) {
	builder.handler(response, request, builder)
}

func (builder Builder) Respond(request *http.Request) Builder {
	return builder
}

func (builder Builder) WithoutError() Builder {
	builder.err = nil

	return builder
}

func (builder Builder) Merge(other Builder) Builder {
	if other.status != 0 {
		builder.status = other.status
	}

	if other.headers != nil {
		for key, values := range other.headers {
			for _, value := range values {
				builder.headers.Add(key, value)
			}
		}
	}

	if other.body != nil {
		builder.body = other.body
	}

	if other.stream != nil {
		builder.stream = other.stream
	}

	if other.err != nil {
		builder.err = other.err
	}

	return builder
}
