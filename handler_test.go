package akumu_test

import (
	"net/http"
	"testing"

	"github.com/studiolambda/akumu"
)

func exampleRawResponse(request *http.Request) error {
	return akumu.Raw(
		http.RedirectHandler("foo", http.StatusTemporaryRedirect),
	)
}

func TestRawHandler(t *testing.T) {
	handler := akumu.RawHandler(http.RedirectHandler("foo", http.StatusTemporaryRedirect))

	request, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatal("failed to create http request")
	}

	response := handler.Record(request)

	if expected := http.StatusTemporaryRedirect; response.Code != expected {
		t.Fatalf("expected status code %d but got %d", expected, response.Code)
	}
}

func TestRawResponse(t *testing.T) {
	handler := akumu.Handler(exampleRawResponse)

	request, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatal("failed to create http request")
	}

	response := handler.Record(request)

	if expected := http.StatusTemporaryRedirect; response.Code != expected {
		t.Fatalf("expected status code %d but got %d", expected, response.Code)
	}
}
