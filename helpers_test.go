package akumu_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/studiolambda/akumu"
)

type JsonHandlerTestPayload struct {
	Foo int    `json:"foo"`
	Bar string `json:"bar"`
}

func TestJSON(t *testing.T) {
	request, err := http.NewRequest("GET", "/", strings.NewReader(`{"foo":10,"bar":"hello"}`))

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	payload, err := akumu.JSON[JsonHandlerTestPayload](request)

	if err != nil {
		t.Fatalf("unable to decode payload from request: %s", err)
	}

	if expected := 10; payload.Foo != expected {
		t.Fatalf("unexpected foo: %d, expected %d", payload.Foo, expected)
	}

	if expected := "hello"; payload.Bar != expected {
		t.Fatalf("unexpected bar: %s, expected %s", payload.Bar, expected)
	}
}
