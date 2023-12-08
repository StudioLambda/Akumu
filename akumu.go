package akumu

import (
	netHttp "net/http"

	"github.com/go-chi/chi/v5"
)

type Akumu struct {
	router chi.Router
}

func New() *Akumu {
	return &Akumu{
		router: chi.NewRouter(),
	}
}

func (akumu *Akumu) Start() {
	netHttp.ListenAndServe(":3000", akumu.router)
}
