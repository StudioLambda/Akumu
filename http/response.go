package http

import (
	"bytes"
	"encoding/json"
	"io"
)

type Response struct {
	status  Status
	version Version
	headers Headers
	body    *Body
}

func (response *Response) safeHeaders() *Headers {
	if response.headers == nil {
		response.headers = make(Headers)
	}

	return &response.headers
}

func (response Response) safe() Response {
	if response.headers == nil {
		response.headers = make(Headers)
	}

	if response.body == nil {
		response.body = NewBody(bytes.NewReader([]byte{}))
	}

	if response.status == 0 {
		response.status = StatusOK
	}

	if response.version == "" {
		response.version = Version1_1
	}

	return response
}

func (response Response) Header(key, value string) Response {
	response.safeHeaders().Insert(key, value)

	return response
}

func (response Response) Status(status Status) Response {
	response.status = status

	return response
}

func (response Response) Version(version Version) Response {
	response.version = version

	return response
}

func (response Response) Body(body *Body) Response {
	response.body = body

	return response
}

func (response Response) BodyReader(reader io.Reader) Response {
	return response.Body(NewBody(reader))
}

func (response Response) BodyBytes(body []byte) Response {
	return response.BodyReader(bytes.NewReader(body))
}

func (response Response) Error(err error) Response {
	body := []byte(err.Error())
	response.body = NewBody(bytes.NewReader(body))

	return response
}

func (response Response) ErrorJSON(err error) Response {
	body, _ := json.Marshal(map[string]string{
		"error": err.Error(),
	})

	return response.
		Header("Content-Type", "application/json").
		BodyBytes(body)
}

func (response Response) JSON(value any) Response {
	encoded, err := json.Marshal(value)

	if err != nil {
		return response.ErrorJSON(err)
	}

	return response.
		Header("Content-Type", "application/json").
		BodyBytes(encoded)
}

func (response Response) IsJson() bool {
	return response.safeHeaders().Contains("Content-Type", "application/json")
}
