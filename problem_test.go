package akumu_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/studiolambda/akumu"
)

var (
	ErrNotAuthenticated = akumu.Problem{
		Type:   "http://example.com/problems/unauthenticated",
		Title:  "Unauthenticated",
		Detail: "The supplied credentials are invalid",
		Status: http.StatusUnauthorized,
	}
)

func customProblemHandler(request *http.Request) error {
	return akumu.Failed(
		ErrNotAuthenticated.With("username", "foobar"),
	)
}

func customProblemHandler2(request *http.Request) error {
	return akumu.
		Failed(ErrNotAuthenticated).
		Header("X-Foo", "Bar")
}

func problemHandler(request *http.Request) error {
	return akumu.Failed(akumu.Problem{
		Type:     "http://example.com/problems/not-found",
		Title:    http.StatusText(http.StatusNotFound),
		Detail:   "The requested resource could not be found.",
		Status:   http.StatusNotFound,
		Instance: request.URL.String(),
	})
}

func TestItHandlesErrorsWithHeaders(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)
	request.Header.Add("Accept", "application/problem+json")

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	response := akumu.Record(customProblemHandler2, request)

	if response.Code != ErrNotAuthenticated.Status {
		t.Fatalf("unexpected status code: %d, expected %d", response.Code, ErrNotAuthenticated.Status)
	}

	if value := response.Header().Get("X-Foo"); value != "Bar" {
		t.Fatalf("unexpected header value: %s, expected %s", value, "Bar")
	}
}

func TestProblemHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	response := akumu.Record(problemHandler, request)

	if response.Code != http.StatusNotFound {
		t.Fatalf("unexpected status code: %d", response.Code)
	}
}

func TestCustomProblemHandler(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)
	request.Header.Add("Accept", "application/problem+json")

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	response := akumu.Record(customProblemHandler, request)

	data := make(map[string]any)

	if err := json.Unmarshal(response.Body.Bytes(), &data); err != nil {
		t.Fatalf("unable to deserialize response body: %s", err)
	}

	if response.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status code: %d", response.Code)
	}

	if username, ok := data["username"]; !ok || username != "foobar" {
		t.Fatalf("unexpected username: %s", username)
	}
}
