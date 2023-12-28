package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/studiolambda/golidate/format"
)

type RawHandler func(request Request, response Response, writer Writer)

type Handler func(request Request) (response Response)

func (handler Handler) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	request := NewRequest(req)
	response := handler(request).safe()

	response.handler(request, response, Writer{writer})
}

func bodyHandler(request Request, response Response, writer Writer) {
	for key, values := range response.headers {
		for _, value := range values {
			writer.Header().Set(key, value)
		}
	}

	bytes, err := response.body.Bytes()

	if err != nil {
		writer.WriteHeader(int(StatusInternalServerError))

		return
	}

	writer.WriteHeader(int(response.status))
	writer.Write(bytes)
}

func streamHandler(messages <-chan []byte) RawHandler {
	return func(request Request, response Response, writer Writer) {
		for key, values := range response.headers {
			for _, value := range values {
				writer.Header().Set(key, value)
			}
		}

		writer.WriteHeader(int(response.status))

		for {
			select {
			case <-request.Context().Done():
				return
			case message, ok := <-messages:
				if !ok {
					return
				}

				writer.Write(message)
				writer.Flush()
			}
		}
	}
}

func sseHandler(messages <-chan SSE) RawHandler {
	return func(request Request, response Response, writer Writer) {
		for key, values := range response.headers {
			for _, value := range values {
				writer.Header().Set(key, value)
			}
		}

		writer.WriteHeader(int(response.status))

		for {
			select {
			case <-request.Context().Done():
				return
			case message, ok := <-messages:
				if !ok {
					return
				}

				writer.Write(message.Bytes())
				writer.Flush()
			}
		}
	}
}

func errorHandler(err error) RawHandler {
	return func(request Request, response Response, writer Writer) {
		headers := response.headers
		status := response.status
		traces := stackTrace(err)
		fields := make(Fields)

		for i := range traces {
			if err, ok := traces[len(traces)-1-i].(ErrorStatus); ok {
				if s := err.ErrorStatus(); s != 0 {
					status = s
				}
			}

			if err, ok := traces[len(traces)-1-i].(ErrorHeaders); ok {
				headers = headers.Merge(err.ErrorHeaders())
			}

			if err, ok := traces[len(traces)-1-i].(ErrorFields); ok {
				fields = fields.Merge(err.ErrorFields())
			}
		}

		for key, values := range headers {
			for _, value := range values {
				writer.Header().Set(key, value)
			}
		}

		result := ErrorResponse{
			Message: format.Capitalize()(err.Error()),
			Fields:  fields,
		}

		if request.Headers().Contains("Accept", "application/json") {
			writer.Header().Set("Content-Type", "application/json")
			writer.WriteHeader(int(status))
			result, _ := json.Marshal(result)
			writer.Write(result)

			return
		}

		if request.Headers().Contains("Accept", "text/html") {
			writer.Header().Set("Content-Type", "text/html")
			writer.WriteHeader(int(status))
			writer.Write([]byte(
				fmt.Sprintf("<h1>%s - %d</h1><p>%s</p>", status, status, result.Message),
			))

			return
		}

		writer.WriteHeader(int(status))
		writer.Write([]byte(err.Error()))
	}
}
