package akumu

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type ResponderHandler func(http.ResponseWriter, *http.Request, Responder)

type Responder struct {
	handler ResponderHandler
	code    int
	headers http.Header
	body    io.Reader
	err     error
	stream  <-chan []byte
}

func DefaultResponderHandler(writer http.ResponseWriter, request *http.Request, responder Responder) {
	for key, values := range responder.headers {
		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}

	if responder.err != nil {
		http.Error(writer, responder.err.Error(), responder.code)

		return
	}

	if responder.body != nil {
		body, err := io.ReadAll(responder.body)

		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)

			return
		}

		writer.WriteHeader(responder.code)
		writer.Write(body)

		return
	}

	if responder.stream != nil {
		writer.WriteHeader(responder.code)

		for {
			select {
			case <-request.Context().Done():
				return
			case message, ok := <-responder.stream:
				if !ok {
					return
				}

				writer.Write(message)
				writer.(http.Flusher).Flush()
			}
		}
	}

	writer.WriteHeader(responder.code)
}

func Response(code int) Responder {
	return Responder{
		handler: DefaultResponderHandler,
		code:    code,
		headers: make(http.Header),
		body:    nil,
		err:     nil,
		stream:  nil,
	}
}

func (responder Responder) Error() string {
	return http.StatusText(responder.code)
}

func (responder Responder) Code(code int) Responder {
	responder.code = code

	return responder
}

func (responder Responder) Headers(headers http.Header) Responder {
	responder.headers = headers

	return responder
}

func (responder Responder) Header(key, value string) Responder {
	responder.headers.Set(key, value)

	return responder
}

func (responder Responder) AppendHeader(key, value string) Responder {
	responder.headers.Add(key, value)

	return responder
}

func (responder Responder) Body(body []byte) Responder {
	return responder.BodyReader(bytes.NewReader(body))
}

func (responder Responder) BodyReader(body io.Reader) Responder {
	responder.body = body

	return responder
}

func (responder Responder) Stream(stream <-chan []byte) Responder {
	responder.stream = stream

	return responder.
		Header("Cache-Control", "no-cache").
		Header("Connection", "keep-alive")
}

func (responder Responder) SSE(stream <-chan []byte) Responder {
	return responder.
		Header("Content-Type", "text/event-stream").
		Stream(stream)
}

func (responder Responder) Failed(err error) Responder {
	responder.err = err

	return responder
}

func (responder Responder) JSON(body any) Responder {
	buffer := &bytes.Buffer{}

	if err := json.NewEncoder(buffer).Encode(body); err != nil {
		return responder.
			Code(http.StatusInternalServerError).
			Failed(err)
	}

	return responder.
		Header("Content-Type", "application/json").
		BodyReader(buffer)
}

func (responder Responder) Handler(handler ResponderHandler) Responder {
	responder.handler = handler

	return responder
}

func (responder Responder) Handle(response http.ResponseWriter, request *http.Request) {
	responder.handler(response, request, responder)
}
