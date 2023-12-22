package http

import (
	"net/http"
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

		for message := range messages {
			writer.Write(message)
			writer.Flush()
		}
	}
}

func sseHandler(messages <-chan SSEEvent) RawHandler {
	return func(request Request, response Response, writer Writer) {
		for key, values := range response.headers {
			for _, value := range values {
				writer.Header().Set(key, value)
			}
		}

		writer.WriteHeader(int(response.status))

		for message := range messages {
			writer.Write(message.Bytes())
			writer.Flush()
		}
	}
}
