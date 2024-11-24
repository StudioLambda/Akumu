package akumu

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

// BuilderHandler type is used to define how a request
// should be responded based on the given [Builder].
//
// This function is what ultimately gets called whenever
// akumu needs to send a response.
type BuilderHandler func(http.ResponseWriter, *http.Request, Builder)

// Builder is the type used to build and handle an
// incoming http request's response. That means that
// it contains the necesary utilities to handle the
// [http.Request].
//
// While the internals of this type are private, there's
// a lot of building methods attached to it, making it
// possible to quickly build Builder types.
//
// One should most likely not build a request from scratch
// and instead use the [Response] and [Failed] functions
// to create a new [Builder].
//
// The default [BuilderHandler] is [DefaultResponderHandler].
type Builder struct {

	// handler is the main response handler that will be executed
	// whenever there's the need to resolve this builder into
	// an actual http response. That means, writing the given result
	// of the builder into an [http.ResponseWriter].
	handler BuilderHandler

	// status is the http status code that the builder currently holds.
	//
	// Usually, the handler determines what to do with but, but the most
	// likely outcome is for it to be the actual response http code.
	status int

	// headers stores a [http.Header] of the actual http response.
	//
	// Although the handler determines what to do with it, the most
	// likely outcome is for it to be the actual response http headers.
	headers http.Header

	// err stores if the current builder is expected to fail. This is
	// useful as it contains the specific error that should be "thrown".
	//
	// There's specific behaviour defined in the [DefaultResponderHandler]
	// for handling errors that also implement [Responder].
	//
	// The [DefaultResponderHandler] has a specific priority, given there's also
	// the body, err, stream and writer possibilities.
	//
	// The handler determines what to do with it.
	err error

	// body is an [io.Reader] that will most likely be used as the response
	// body. That means that it will be read and written to the actual
	// [http.ResponseWriter].
	//
	// The [DefaultResponderHandler] has a specific priority, given there's also
	// the body, err, stream and writer possibilities.
	//
	// The handler determines what to do with it.
	body io.Reader

	// stream is a channel that can be used to directly stream parts of the
	// response, effectively using HTTP streaming using [http.ResponseWriter].
	//
	// The [DefaultResponderHandler] has a specific priority, given there's also
	// the body, err, stream and writer possibilities.
	//
	// The handler determines what to do with it.
	stream <-chan []byte

	// writer is a custom function that can be used directly to manipulate
	// the actual response that's sent. This is useful for custom streaming
	// or file downloads.
	//
	// The [DefaultResponderHandler] has a specific priority, given there's also
	// the body, err, stream and writer possibilities.
	//
	// The handler determines what to do with it.
	writer func(writer http.ResponseWriter)
}

var (
	// ErrWriterRequiresFlusher is an error that determines that
	// the given response writer needs a flusher in order to push
	// changes to the [http.ResponseWriter]. This should most likely
	// not happen due [http.ResponseWriter] already implementing [http.Flusher].
	ErrWriterRequiresFlusher = errors.New("response writer requires a flusher")
)

// writeHeaders writes the given response headers by calling
// [http.ResponseWriter]'s `WriteHeader` method.
//
// This function also returns true if the response status code is
// between [500, 599], making it possible to know if there was
// a server error.
func writeHeaders(writer http.ResponseWriter, builder Builder) bool {
	for key, values := range builder.headers {
		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}

	writer.WriteHeader(builder.status)

	return builder.status >= 500 && builder.status < 600
}

// DefaultResponderHandler is the default handler that's used in
// a [Builder]. This handler does most of the things expected
// for akumu to handle http responses, although it can be customized
// if needed.
//
// By default, this handler does make use of the [OnErrorHook] that is
// found in the request's context, on the [OnErrorKey] key if server error happens.
//
// By default, this handler does handle the [Builder] in the following order of priority:
//  1. errors
//  2. writer
//  3. body
//  4. stream
//  5. default (no body)
//
// This means that if a [Builder] contain more than one possible response type, only the
// first one defined, following the order above, will be executed.
func DefaultResponderHandler(writer http.ResponseWriter, request *http.Request, builder Builder) {
	onError, hasOnError := request.Context().Value(OnErrorKey{}).(OnErrorHook)

	if builder.err != nil {
		parent := builder.WithoutError()
		handle(writer, request, builder.err, &parent)

		return
	}

	if builder.writer != nil {
		if writeHeaders(writer, builder) && hasOnError {
			serverErr := ErrServer{
				Code:    builder.status,
				Request: request,
			}

			onError(errors.Join(serverErr, ErrServerWriter))
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

		if writeHeaders(writer, builder) && hasOnError {
			serverErr := ErrServer{
				Code:    builder.status,
				Request: request,
			}

			onError(errors.Join(serverErr, ErrServerBody))
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

		if writeHeaders(writer, builder) && hasOnError {
			serverErr := ErrServer{
				Code:    builder.status,
				Request: request,
			}

			onError(errors.Join(serverErr, ErrServerStream))
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

	if writeHeaders(writer, builder) && hasOnError {
		serverErr := ErrServer{
			Code:    builder.status,
			Request: request,
		}

		onError(errors.Join(serverErr, ErrServerDefault))
	}
}

// Response is used to create a builder with the given
// HTTP status.
//
// This starts a new builder that can directly be returned
// in a [Handler], as a [Builder] implements the error interface.
func Response(status int) Builder {
	return Builder{
		handler: DefaultResponderHandler,
		status:  status,
		headers: make(http.Header),
		body:    nil,
		stream:  nil,
	}
}

// Failed is used to create a builder with the given
// error as the main factor. It should only be used
// by either custom errors that modify responses (such
// as the [Problem] errors) or errors that are often
// server errors (such as database errors) that are
// controled but out of the nature of the request.
//
// This is basically an alias for a [Response] with
// a status code [http.StatusInternalServerError] and
// a failed error.
//
// This starts a new builder that can directly be returned
// in a [Handler], as a [Builder] implements the error interface.
func Failed(err error) Builder {
	return Response(http.StatusInternalServerError).
		Failed(err)
}

// Error implements the error interface for any [Builder],
// making it possible to be used in a [Handler].
func (builder Builder) Error() string {
	return http.StatusText(builder.status)
}

// Status sets the HTTP status of the [Builder] into
// the desired one, leading the the response http status code.
func (builder Builder) Status(status int) Builder {
	builder.status = status

	return builder
}

// Headers sets the given HTTP headers to the builder.
//
// This is an override operation and does not merge previous
// headers.
func (builder Builder) Headers(headers http.Header) Builder {
	builder.headers = headers

	return builder
}

// Header sets a new header to the [Builder].
//
// This overrides any previous headers with the same key.
//
// Because of immutability, this method creates a copy
// of the headers first, so previous [Builder] instances
// are not affected.
//
// You may look into [Builder.AppendHeader] if you prefer
// to append rather than override.
func (builder Builder) Header(key, value string) Builder {
	builder.headers = builder.headers.Clone()
	builder.headers.Set(key, value)

	return builder
}

// AppendHeader appends a new header to the [Builder].
//
// This does not overrides any previous headers with the same key.
//
// Because of immutability, this method creates a copy
// of the headers first, so previous [Builder] instances
// are not affected.
//
// You may look into [Builder.Header] if you prefer
// to override rather than append.
func (builder Builder) AppendHeader(key, value string) Builder {
	builder.headers = builder.headers.Clone()
	builder.headers.Add(key, value)

	return builder
}

// Body sets the [Builder]'s body reader to a new
// [bytes.Reader] with the given []byte.
//
// If you already have a reader use the [Builder.BodyReader].
func (builder Builder) Body(body []byte) Builder {
	return builder.BodyReader(bytes.NewReader(body))
}

// Body sets the [Builder]'s body reader.
//
// If you don't have a reader use the [Builder.Body].
func (builder Builder) BodyReader(body io.Reader) Builder {
	builder.body = body

	return builder
}

// Stream sets the [Builder] to stream from the given channel.
//
// The streaming will end as soon as the channel is closed or
// whenever the request's context is canceled.
//
// It also sets the aproppiate request headers to stream, most
// notably, the Cache-Control to `no-cache` and Connection to `keep-alive`.
func (builder Builder) Stream(stream <-chan []byte) Builder {
	builder.stream = stream

	return builder.
		Header("Cache-Control", "no-cache").
		Header("Connection", "keep-alive")
}

// Stream sets the [Builder] to stream from the given channel.
//
// The streaming will end as soon as the channel is closed or
// whenever the request's context is canceled.
//
// A part from the [Stream] headers, this method additionally
// sets the Content-Type to `text/event-stream`.
func (builder Builder) SSE(stream <-chan []byte) Builder {
	return builder.
		Header("Content-Type", "text/event-stream").
		Stream(stream)
}

// Cookie sets a new [http.Cookie] to the [Builder]'s response.
//
// This effectively appends a new "Set-Cookie" header with the
// cokkies' value.
func (builder Builder) Cookie(cookie http.Cookie) Builder {
	if c := cookie.String(); c != "" {
		return builder.AppendHeader("Set-Cookie", c)
	}

	return builder
}

// Failed indicates that the current [Builder] is
// intended to fail with the given error.
//
// The status code is not changed, but you may use
// [Builder.Status] or use [Failed] directly.
func (builder Builder) Failed(err error) Builder {
	builder.err = err

	return builder
}

// Text sets the response body text to the given
// string and also makes sure the Content-Type
// is set to "text/plain".
func (builder Builder) Text(body string) Builder {
	return builder.
		Header("Content-Type", "text/plain").
		Body([]byte(body))
}

// HTML sets the response body text to the given
// string and also makes sure the Content-Type
// is set to "text/html".
func (builder Builder) HTML(html string) Builder {
	return builder.
		Header("Content-Type", "text/html").
		Body([]byte(html))
}

// JSON encodes the given body variable into the
// request's body and also makes sure the Content-Type
// is set to "application/json".
//
// If the encoding fails, the [Builder] is set to a
// status of [http.StatusInternalServerError] with the
// failed error.
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

// BodyWriter marks the [Builder] as a custom writer function
// making the handler execute the logic passed here when a response
// is written.
//
// This is useful when you need a lot of control over how a response
// is written to the actual [http.ResponseWriter] such as for file downloads.
func (builder Builder) BodyWriter(writer func(writer http.ResponseWriter)) Builder {
	builder.writer = writer

	return builder
}

// Handler sets a custom [BuilderHandler] into the current [Builder].
//
// This is not likely something that is needed as the [DefaultResponderHandler]
// takes care of most needed things already.
func (builder Builder) Handler(handler BuilderHandler) Builder {
	builder.handler = handler

	return builder
}

// Handle executes the given [Builder] with the given response writer and request.
//
// It internally passes it to the [Builder]'s Handler.
func (builder Builder) Handle(response http.ResponseWriter, request *http.Request) {
	builder.handler(response, request, builder)
}

// Respond implements [Responder] interface on a [Builder].
func (builder Builder) Respond(request *http.Request) Builder {
	return builder
}

// WithoutError is used to remove any errors from the current [Builder].
func (builder Builder) WithoutError() Builder {
	builder.err = nil

	return builder
}

// Merge merges another builder with this one.
//
// The [Builder] passed as a parameter takes precedence.
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
