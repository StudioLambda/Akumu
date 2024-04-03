package akumu

import "net/http"

type Responder interface {
	Respond(request *http.Request) Builder
}
