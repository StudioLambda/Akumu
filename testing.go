package akumu

import (
	"net/http"
	"net/http/httptest"
)

// RecordServer records what a given [http.Server] would give as a reponse to a [http.Request].
//
// The response is recorded using a [httptest.ResponseRecorder].
func RecordServer(server *http.Server, request *http.Request) *httptest.ResponseRecorder {
	return RecordHandler(server.Handler, request)
}

// Record records what a given [Handler] would give as a reponse to a [http.Request].
//
// The response is recorded using a [httptest.ResponseRecorder].
func Record(handler Handler, request *http.Request) *httptest.ResponseRecorder {
	return RecordHandler(handler, request)
}

// Record records what a given [http.Handler] would give as a reponse to a [http.Request].
//
// The response is recorded using a [httptest.ResponseRecorder].
func RecordHandler(handler http.Handler, request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)

	return recorder
}
