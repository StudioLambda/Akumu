package akumu

import (
	"net/http"
	"net/http/httptest"
)

func RecordServer(server *http.Server, request *http.Request) *httptest.ResponseRecorder {
	return RecordHandler(server.Handler, request)
}

func Record(handler Handler, request *http.Request) *httptest.ResponseRecorder {
	return RecordHandler(handler, request)
}

func RecordHandler(handler http.Handler, request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	return recorder
}
