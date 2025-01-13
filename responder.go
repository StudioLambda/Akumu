package akumu

import "net/http"

// Responder interface defines a common way for types
// to be able to respond to an HTTP request in a custom way.
//
// For example, this is used by the [Problem] type to define
// a custom response.
//
// Keep in mind that because it's often necesary to return this
// in a [Handler], it's likely needed to implement the error
// interface as well.
type Responder interface {
	Respond(request *http.Request) Builder
}

// RawBuilder is a raw response that can be used
// to return a raw [http.Builder] from a [Handler].
type RawResponder interface {
	http.Handler
}
