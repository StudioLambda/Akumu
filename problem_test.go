package akumu_test

import (
	"net/http"
	"testing"

	"github.com/studiolambda/akumu"
)

func problemHandler(request *http.Request) error {
	return akumu.Failed(akumu.Problem{
		Type:     "http://example.com/problems/not-found",
		Title:    http.StatusText(http.StatusNotFound),
		Detail:   "The requested resource could not be found.",
		Status:   http.StatusNotFound,
		Instance: request.URL.String(),
	})
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
