package akumu_test

import (
	"net/http"
	"testing"

	"github.com/studiolambda/akumu"
)

func TestRouterHas(t *testing.T) {
	router := akumu.NewRouter()

	router.Get("/", func(request *http.Request) error {
		return akumu.Response(http.StatusOK)
	})

	t.Run("Exists", func(t *testing.T) {
		if !router.Has(http.MethodGet, "/") {
			t.Fatal("router should have the route")
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		if router.Has(http.MethodGet, "/not-found") {
			t.Fatal("router should not have the route")
		}
	})
}

func TestRouterMatches(t *testing.T) {
	router := akumu.NewRouter()

	router.Get("/", func(request *http.Request) error {
		return akumu.Response(http.StatusOK)
	})

	request, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatalf("unable to create http request: %v", err)
	}

	if !router.Matches(request) {
		t.Fatal("router does not have the route")
	}
}

func TestRouterHandler(t *testing.T) {
	router := akumu.NewRouter()

	router.Get("/", func(request *http.Request) error {
		return akumu.Response(http.StatusOK)
	})

	request, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatalf("failed to create http request: %v", err)
	}

	response := router.Record(request)

	if expected := http.StatusOK; response.Code != expected {
		t.Fatalf("expected status code %d but got %d", expected, response.Code)
	}
}
