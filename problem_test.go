package akumu_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
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
	ErrSomethingBad = errors.Join(
		errors.New("first error"),
		errors.Join(errors.New("second error"), errors.New("third error")),
		fmt.Errorf("%w: failed", errors.New("last error")),
	)
)

func customProblemHandlerDirectErr(request *http.Request) error {
	return akumu.Failed(ErrSomethingBad)
}

func customProblemHandlerWithErr(request *http.Request) error {
	return akumu.Failed(
		ErrNotAuthenticated.WithError(ErrSomethingBad),
	)
}

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

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	request.Header.Add("Accept", "application/problem+json")
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

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	request.Header.Add("Accept", "application/problem+json")
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

func TestProblemErrorUnwraps(t *testing.T) {
	someErr := errors.New("some error")
	someOtherErr := errors.New("some other error")
	err := errors.Join(someErr, someOtherErr)
	problem := akumu.NewProblem(err, http.StatusBadRequest)

	if !errors.Is(problem, err) {
		t.Fatalf("%s should be %s", problem, err)
	}

	if !errors.Is(problem, someErr) {
		t.Fatalf("%s should be %s", problem, err)
	}

	if !errors.Is(problem, someOtherErr) {
		t.Fatalf("%s should be %s", problem, err)
	}
}

func TestProblemWithCustomError(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	request.Header.Add("Accept", "application/problem+json")

	response := akumu.Record(customProblemHandlerWithErr, request)

	data := make(map[string]any)

	if err := json.Unmarshal(response.Body.Bytes(), &data); err != nil {
		t.Fatalf("unable to deserialize response body: %s", err)
	}

	errors, ok := data["errors"]

	if !ok {
		t.Fatalf("expected errors to be in the problem response")
	}

	errs := errors.([]any)

	if expected := 4; len(errs) != expected {
		t.Fatalf("expected %d errors but got %d", expected, len(errs))
	}
}

func TestProblemWithCustomError2(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	response := akumu.Record(customProblemHandlerWithErr, request)
	res := string(response.Body.Bytes())

	if expect := "last error: failed"; !strings.Contains(res, expect) {
		t.Fatalf("error should contain '%s'", expect)
	}
}

func TestProblemDirectError(t *testing.T) {
	request, err := http.NewRequest("GET", "/", nil)

	if err != nil {
		t.Fatalf("unable to create request: %s", err)
	}

	request.Header.Add("Accept", "application/problem+json")

	response := akumu.Record(customProblemHandlerDirectErr, request)

	data := make(map[string]any)

	if err := json.Unmarshal(response.Body.Bytes(), &data); err != nil {
		t.Fatalf("unable to deserialize response body: %s", err)
	}

	errors, ok := data["errors"]

	if !ok {
		t.Fatalf("expected errors to be in the problem response")
	}

	errs := errors.([]any)

	if expected := 4; len(errs) != expected {
		t.Fatalf("expected %d errors but got %d", expected, len(errs))
	}
}
