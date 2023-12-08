package http

import "net/http"

type Handler func(request Request) (response Response)

func (handler Handler) ServeHTTP(writter http.ResponseWriter, req *http.Request) {
	request := NewRequest(req)
	response := handler(request).safe()

	writter.WriteHeader(int(response.status))

	for key, values := range response.headers {
		for _, value := range values {
			writter.Header().Add(key, value)
		}
	}

	bytes, err := response.body.Bytes()

	if err != nil {
		writter.WriteHeader(int(StatusInternalServerError))

		return
	}

	writter.Write(bytes)
}
