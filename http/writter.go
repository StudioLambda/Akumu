package http

import "net/http"

type Writer struct {
	http.ResponseWriter
}

func (writer Writer) Flush() {
	writer.ResponseWriter.(http.Flusher).Flush()
}
