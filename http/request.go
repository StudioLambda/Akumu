package http

import (
	"context"
	"net/http"
)

type Request struct {
	method  Method
	url     URL
	headers Headers
	version Version
	body    Body
	ctx     context.Context
}

func NewRequest(request *http.Request) Request {
	return Request{
		method:  Method(request.Method),
		url:     NewURL(request.URL),
		version: Version(request.Proto),
		headers: Headers(request.Header.Clone()),
		body:    NewBody(request.Body),
		ctx:     request.Context(),
	}
}

func (request Request) Context() context.Context {
	return request.ctx
}

func (request Request) Method() Method {
	return request.method
}

func (request Request) URL() URL {
	return request.url
}

func (request Request) Version() Version {
	return request.version
}

func (request Request) Headers() Headers {
	return request.headers
}

func (request Request) Body() Body {
	return request.body
}
