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
	status  int
	headers http.Header
	body    io.Reader
	err     error
	stream  <-chan []byte
}

func HandleError(writer http.ResponseWriter, request *http.Request, err error, status int) {
	if responder, ok := err.(Responder); ok {
		responder.Handle(writer, request)

		return
	}

	if request.Header.Get("Accept") == "application/problem+json" {
		problem := Problem{
			Type:     "about:blank",
			Title:    http.StatusText(status),
			Detail:   err.Error(),
			Status:   status,
			Instance: request.URL.String(),
		}

		Failed(problem).Handle(writer, request)

		return
	}

	http.Error(writer, err.Error(), status)
}

func DefaultResponderHandler(writer http.ResponseWriter, request *http.Request, responder Responder) {
	for key, values := range responder.headers {
		for _, value := range values {
			writer.Header().Add(key, value)
		}
	}

	if responder.err != nil {
		HandleError(writer, request, responder.err, responder.status)

		return
	}

	if responder.body != nil {
		body, err := io.ReadAll(responder.body)

		if err != nil {
			HandleError(writer, request, err, http.StatusInternalServerError)

			return
		}

		writer.WriteHeader(responder.status)
		writer.Write(body)

		return
	}

	if responder.stream != nil {
		writer.WriteHeader(responder.status)

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

	writer.WriteHeader(responder.status)
}

func Response(status int) Responder {
	return Responder{
		handler: DefaultResponderHandler,
		status:  status,
		headers: make(http.Header),
		body:    nil,
		err:     nil,
		stream:  nil,
	}
}

func Failed(err error) Responder {
	return Response(http.StatusInternalServerError).Failed(err)
}

func (responder Responder) Error() string {
	return http.StatusText(responder.status)
}

func (responder Responder) Status(status int) Responder {
	responder.status = status

	return responder
}

func (responder Responder) Headers(headers http.Header) Responder {
	responder.headers = headers

	return responder
}

func (responder Responder) Header(key, value string) Responder {
	responder.headers = responder.headers.Clone()
	responder.headers.Set(key, value)

	return responder
}

func (responder Responder) AppendHeader(key, value string) Responder {
	responder.headers = responder.headers.Clone()
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
	if problem, ok := err.(Problem); ok {
		if problem.Type == "" {
			problem.Type = "about:blank"
		}

		if problem.Status == 0 {
			problem.Status = responder.status
		}

		if problem.Title == "" {
			problem.Title = http.StatusText(problem.Status)
		}

		return responder.
			Status(problem.Status).
			JSON(problem).
			Header("Content-Type", "application/problem+json")
	}

	responder.err = err

	return responder
}

func (responder Responder) Text(body string) Responder {
	return responder.
		Header("Content-Type", "text/plain").
		Body([]byte(body))
}

func (responder Responder) HTML(html string) Responder {
	return responder.
		Header("Content-Type", "text/html").
		Body([]byte(html))
}

func (responder Responder) JSON(body any) Responder {
	buffer := &bytes.Buffer{}

	if err := json.NewEncoder(buffer).Encode(body); err != nil {
		return responder.
			Status(http.StatusInternalServerError).
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
