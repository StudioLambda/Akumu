package akumu

import "net/http"

// Middleware is just an alias for a function that
// takes a handler and returns another handler.
//
// Use this in places where the common pattern
// of middlewares is needed.
//
// To maintain compatibility with the ecosystem,
// the handlers used by middlewares are [http.Handler]
// instead of akumu's [Handler]. Use the [Builder.Handle]
// to help when dealing with akumu in middlewares.
type Middleware = func(http.Handler) http.Handler
