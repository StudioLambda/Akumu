package akumu

import (
	"net/http"
	"sync"
)

type ContextKey struct{}

type OnErrorNext func(ServerError)
type OnErrorCallback func(ServerError, OnErrorNext)

// Contextual is a concurrent-safe implementation
// to append akumu-specific context to the given request.
type Contextual struct {
	mutex   sync.Mutex
	onError []OnErrorCallback
}

func NewContext() *Contextual {
	return &Contextual{}
}

func Context(request *http.Request) (*Contextual, bool) {
	ctx, ok := request.Context().Value(ContextKey{}).(*Contextual)

	return ctx, ok
}

func (ctx *Contextual) OnError(callback OnErrorCallback) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	ctx.onError = append(ctx.onError, callback)
}

func (ctx *Contextual) handleError(err ServerError) {
	ctx.mutex.Lock()
	defer ctx.mutex.Unlock()

	var index int
	var next func(ServerError)

	next = func(err ServerError) {
		if index < len(ctx.onError) {
			callback := ctx.onError[index]
			index++
			callback(err, next)
		}
	}

	next(err)
}
