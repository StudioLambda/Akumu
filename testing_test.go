package akumu_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/studiolambda/akumu"
)

func handler(request *http.Request) error {
	return akumu.Response(http.StatusOK)
}

func handler2(request *http.Request) error {
	return akumu.Failed(errors.New("failure"))
}

func handler3(request *http.Request) error {
	return akumu.Response(http.StatusNotImplemented).Failed(errors.New("failure"))
}

func handler4(request *http.Request) error {
	return errors.New("failure")
}

func TestHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	response := akumu.Record(handler, request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"unexpected status code: expected %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}
}

func TestHTTPHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	response := akumu.RecordHandler(akumu.Handler(handler), request)

	if response.Code != http.StatusOK {
		t.Fatalf(
			"unexpected status code: expected %d, got %d",
			http.StatusOK,
			response.Code,
		)
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
		t.Fatalf(
			"unexpected status code: expected %d, got %d",
			http.StatusOK,
			response.Code,
		)
	}
}

func TestHandler2(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	response := akumu.Record(handler2, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf(
			"unexpected status code: expected %d, got %d",
			http.StatusInternalServerError,
			response.Code,
		)
	}
}

func TestHandler3(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	response := akumu.Record(handler3, request)

	if response.Code != http.StatusNotImplemented {
		t.Fatalf(
			"unexpected status code: expected %d, got %d",
			http.StatusNotImplemented,
			response.Code,
		)
	}
}

func TestHandler4(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	response := akumu.Record(handler4, request)

	if response.Code != http.StatusInternalServerError {
		t.Fatalf(
			"unexpected status code: expected %d, got %d",
			http.StatusInternalServerError,
			response.Code,
		)
	}
}
