package akumu

import (
	"github.com/go-chi/chi/v5"
	"github.com/studiolambda/akumu/http"
)

func (akumu *Akumu) Handle(method http.Method, path string, handler http.Handler) {
	akumu.server.Handler.(chi.Router).Method(string(method), path, handler)
}

func (akumu *Akumu) Get(path string, handler http.Handler) {
	akumu.Handle(http.MethodGet, path, handler)
}

func (akumu *Akumu) Post(path string, handler http.Handler) {
	akumu.Handle(http.MethodPost, path, handler)
}

func (akumu *Akumu) Put(path string, handler http.Handler) {
	akumu.Handle(http.MethodPut, path, handler)
}

func (akumu *Akumu) Patch(path string, handler http.Handler) {
	akumu.Handle(http.MethodPatch, path, handler)
}

func (akumu *Akumu) Delete(path string, handler http.Handler) {
	akumu.Handle(http.MethodDelete, path, handler)
}

func (akumu *Akumu) Head(path string, handler http.Handler) {
	akumu.Handle(http.MethodHead, path, handler)
}

func (akumu *Akumu) Connect(path string, handler http.Handler) {
	akumu.Handle(http.MethodConnect, path, handler)
}

func (akumu *Akumu) Options(path string, handler http.Handler) {
	akumu.Handle(http.MethodOptions, path, handler)
}

func (akumu *Akumu) Trace(path string, handler http.Handler) {
	akumu.Handle(http.MethodTrace, path, handler)
}
