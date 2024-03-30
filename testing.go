package akumu

import (
	"net/http"
	"net/http/httptest"
)

func RecordServer(server *http.Server, request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	server.Handler.ServeHTTP(recorder, request)

	return recorder
}

func RecordHandler(handler http.Handler, request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	return recorder
}

func Record(handler Handler, request *http.Request) *httptest.ResponseRecorder {
	return RecordHandler(handler, request)
}
