package akumu_test

import (
	"net/http"
	"testing"

	"github.com/studiolambda/akumu"
)

func handler(request *http.Request) error {
	return akumu.Response(http.StatusOK)
}

func TestHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	response := akumu.Record(handler, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d", response.Code)
	}
}

func TestHTTPHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	response := akumu.Record(handler, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d", response.Code)
	}
}

func TestHTTPServer(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	server := &http.Server{Handler: akumu.Handler(handler)}
	response := akumu.RecordServer(server, request)

	if response.Code != http.StatusOK {
		t.Fatalf("unexpected status code: %d", response.Code)
	}
}
