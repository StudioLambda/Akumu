package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

type Response struct {
	status  Status
	version Version
	headers Headers
	body    Body
	handler RawHandler
}

type SSEEvent struct {
	ID    string
	Event string
	Data  []byte
	Retry int
}

func (event SSEEvent) Bytes() []byte {
	var buffer bytes.Buffer

	if event.ID != "" {
		buffer.WriteString("id: ")
		buffer.WriteString(event.ID)
		buffer.WriteString("\n")
	}

	if event.Event != "" {
		buffer.WriteString("event: ")
		buffer.WriteString(event.Event)
		buffer.WriteString("\n")
	}

	if event.Retry > 0 {
		buffer.WriteString("retry: ")
		buffer.WriteString(fmt.Sprintf("%d", event.Retry))
		buffer.WriteString("\n")
	}

	if len(event.Data) > 0 {
		buffer.WriteString("data: ")
		buffer.Write(event.Data)
		buffer.WriteString("\n")
	}

	buffer.WriteString("\n")

	return buffer.Bytes()
}

func (event SSEEvent) String() string {
	return string(event.Bytes())
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

	if response.status == 0 {
		response.status = StatusOK
	}

	if response.version == "" {
		response.version = Version1_1
	}

	if response.handler == nil {
		response.handler = bodyHandler
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

func (response Response) Body(body Body) Response {
	response.body = body

	return response
}

func (response Response) Stream(messages <-chan []byte) Response {
	response.handler = streamHandler(messages)

	return response.
		Header("Access-Control-Allow-Origin", "*").
		Header("Cache-Control", "no-cache").
		Header("Connection", "keep-alive").
		Header("Content-Type", "text/event-stream")
}

func (response Response) SSE(messages <-chan SSEEvent) Response {
	response.handler = sseHandler(messages)

	return response.
		Header("Access-Control-Allow-Origin", "*").
		Header("Cache-Control", "no-cache").
		Header("Connection", "keep-alive").
		Header("Content-Type", "text/event-stream")
}

func (response Response) Handler(handler RawHandler) Response {
	response.handler = handler

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

func (response Response) HTML(html string) Response {
	return response.
		Header("Content-Type", "text/html").
		BodyBytes([]byte(html))
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
