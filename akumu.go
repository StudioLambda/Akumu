package akumu

import (
	netHttp "net/http"

	"github.com/go-chi/chi/v5"
	"github.com/studiolambda/akumu/http"
)

type Akumu struct {
	router chi.Router
}

func New() *Akumu {
	return &Akumu{
		router: chi.NewRouter(),
	}
}

func (akumu *Akumu) Get(path string, handler http.Handler) {
	akumu.router.Method(string(http.MethodGet), path, handler)
}

func (akumu *Akumu) Start() {
	netHttp.ListenAndServe(":3000", akumu.router)
}
