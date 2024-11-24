package akumu_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/studiolambda/akumu"
)

func TestServerErrors(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/foo/bar", nil)

	if err != nil {
		t.Fatalf("unable to build http request: %s", err)
	}

	serverErr := akumu.ErrServer{
		Code:    http.StatusInternalServerError,
		Request: req,
	}

	if !errors.Is(serverErr, akumu.ErrServer{}) {
		t.Fatalf("server error cannot be resolved to its error type")
	}

	serr := akumu.ErrServer{}
	ok := errors.As(serverErr, &serr)

	if !ok {
		t.Fatalf("server error cannot be resolved to its error type")
	}

	if serr.Code != serverErr.Code {
		t.Fatalf("resolved server error code is different from original")
	}
}
