package http

import "net/http"

type Request struct {
	method  Method
	url     URL
	version Version
	body    *Body
}

func NewRequest(request *http.Request) Request {
	return Request{
		method:  Method(request.Method),
		url:     URL(*request.URL),
		version: Version(request.Proto),
		body:    NewBody(request.Body),
	}
}
