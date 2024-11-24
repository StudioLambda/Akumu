package akumu

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"maps"
	"net/http"

	"github.com/studiolambda/akumu/utils"
)

// Problem represents a problem details for HTTP APIs.
// See https://datatracker.ietf.org/doc/html/rfc9457 for more information.
type Problem struct {

	// err stores the native error of the problem if
	// it happens to have one. This struct member can be
	// nil if the problem was created manually.
	//
	// Refer to [NewProblem] to automatically assign an error
	// to the new Problem.
	err error

	// additional stores additional metadata that will be
	// appended to the problem JSON representation. It is used
	// to serialize and de-serialize properties.
	additional map[string]any

	// Type member is a JSON string containing a URI reference
	// that identifies the problem type.
	//
	// Consumers MUST use the Type URI (after resolution,
	// if necessary) as the problem type's primary identifier.
	//
	// When this member is not present, its value is assumed
	// to be "about:blank".
	Type string

	// Title member is a string containing a short,
	// human-readable summary of the problem type.
	//
	// It SHOULD NOT change from occurrence to occurrence of
	// the problem, except for localization (e.g., using proactive
	// content negotiation.
	//
	// The Title string is advisory and is included only for users
	// who are unaware of and cannot discover the semantics of the
	// type URI (e.g., during offline log analysis).
	Title string

	// Detail member is a JSON string containing a human-readable
	// explanation specific to this occurrence of the problem.
	//
	// The Detail string, if present, ought to focus on helping
	// the client correct the problem, rather than giving debugging information.
	//
	// Consumers SHOULD NOT parse the Detail member for information.
	//
	// Extensions are more suitable and less error-prone ways to obtain
	// such information.
	Detail string

	// The Status member is a JSON number indicating the HTTP status
	// code generated by the origin server for this occurrence of the problem.
	//
	// The "status" member, if present, is only advisory; it conveys the HTTP
	// status code used for the convenience of the consumer.
	//
	// Generators MUST use the same status code in the actual HTTP response,
	// to assure that generic HTTP software that does not understand this format
	// still behaves correctly.
	//
	// Consumers can use the status member to determine what the original status
	// code used by the generator was when it has been changed
	// (e.g., by an intermediary or cache) and when a message's content is
	// persisted without HTTP information. Generic HTTP software will still
	// use the HTTP status code.
	Status int

	// 	The "instance" member is a JSON string containing a URI reference that
	// identifies the specific occurrence of the problem.
	//
	// When the "instance" URI is dereferenceable, the problem details object
	// can be fetched from it. It might also return information about the problem
	// occurrence in other formats through use of proactive content negotiation.
	//
	// When the "instance" URI is not dereferenceable, it serves as a unique identifier
	// for the problem occurrence that may be of significance to the server but is
	// opaque to the client.
	//
	// When "instance" contains a relative URI, it is resolved relative to the document's
	// base URI. However, using relative URIs can cause confusion, and they might not
	// be handled correctly by all implementations.
	//
	// For example, if the two resources "https://api.example.org/foo/bar/123"
	// and "https://api.example.org/widget/456" both respond with an "instance" equal
	// to the relative URI reference "example-instance", when resolved they will
	// identify different resources ("https://api.example.org/foo/bar/example-instance"
	// and "https://api.example.org/widget/example-instance", respectively).
	//
	// As a result, it is RECOMMENDED that absolute URIs be used in "instance" when possible,
	// and that when relative URIs are used, they include the full path (e.g., "/instances/123").
	Instance string
}

type ProblemControlsResolver[R any] func(problem Problem, request *http.Request) R

type ProblemControls struct {
	// Lowercase determines if the problem controls
	// should lowercase the errors found.
	Lowercase ProblemControlsResolver[bool]

	// DefaultStatus determines the default status code of a [Problem]
	// in case it does not have one defined.
	DefaultStatus ProblemControlsResolver[int]

	// DefaultType determines the default type of a [Problem]
	// in case it does not have one defined.
	DefaultType ProblemControlsResolver[string]

	// DefaultTitle determines the default title of a [Problem]
	// in case it does not have one defined.
	DefaultTitle ProblemControlsResolver[string]

	// DefaultInstance determines the default instance of a [Problem]
	// in case it does not have one defined.
	DefaultInstance ProblemControlsResolver[string]

	// Response allows customizing the actual Builder response
	// that a [Problem] should be resolved to.
	Response ProblemControlsResolver[Builder]
}

// ProblemsKey is the context key where the
// problem controls are stored in the request.
type ProblemsKey struct{}

func defaultedProblemControls(controls ProblemControls) ProblemControls {
	if controls.Lowercase == nil {
		controls.Lowercase = defaultProblemControlsLowercase
	}

	if controls.DefaultStatus == nil {
		controls.DefaultStatus = defaultProblemControlsStatus
	}

	if controls.DefaultType == nil {
		controls.DefaultType = defaultProblemControlsType
	}

	if controls.DefaultTitle == nil {
		controls.DefaultTitle = defaultProblemControlsTitle
	}

	if controls.DefaultInstance == nil {
		controls.DefaultInstance = defaultProblemControlsInstance
	}

	if controls.Response == nil {
		controls.Response = defaultProblemControlsResponse
	}

	return controls
}

// Problems return the [ProblemControls] used to determine
// how [Problem] respond to http requests.
func Problems(request *http.Request) (ProblemControls, bool) {
	controls, ok := request.
		Context().
		Value(ProblemsKey{}).(ProblemControls)

	return controls, ok
}

func defaultProblemControlsLowercase(problem Problem, request *http.Request) bool {
	return true
}

func defaultProblemControlsStatus(problem Problem, request *http.Request) int {
	return http.StatusInternalServerError
}

func defaultProblemControlsType(problem Problem, request *http.Request) string {
	return "about:blank"
}

func defaultProblemControlsTitle(problem Problem, request *http.Request) string {
	return http.StatusText(problem.Status)
}

func defaultProblemControlsInstance(problem Problem, request *http.Request) string {
	return request.URL.String()
}

func ProblemControlsResponseFrom(responses map[string]Builder) ProblemControlsResolver[Builder] {
	return func(problem Problem, request *http.Request) Builder {
		accept := utils.ParseAccept(request)

		for _, media := range accept.Order() {
			if response, ok := responses[media]; ok {
				return response
			}
		}

		return Response(problem.Status).
			Text(fmt.Sprintf("%d %s\n\n%s", problem.Status, problem.Title, problem.Detail))
	}
}

func defaultProblemControlsResponse(problem Problem, request *http.Request) Builder {
	responses := map[string]Builder{
		"application/problem+json": Response(problem.Status).
			JSON(problem).
			Header("Content-Type", "application/problem+json"),
		"application/json": Response(problem.Status).
			JSON(problem).
			Header("Content-Type", "application/problem+json"),
	}

	return ProblemControlsResponseFrom(responses)(problem, request)
}

// NewProblem creates a new [Problem] from
// the given error and status code.
func NewProblem(err error, status int) Problem {
	return Problem{
		err:        err,
		additional: make(map[string]any),
		Detail:     err.Error(),
		Status:     status,
	}
}

// Additional returns the additional value of the given key.
//
// Use [Problem.With] to add additional values.
// Use [Problem.Without] to remove additional values.
func (problem Problem) Additional(key string) (any, bool) {
	additional, ok := problem.additional[key]

	return additional, ok
}

// With adds a new additional value to the given key.
func (problem Problem) With(key string, value any) Problem {
	if problem.additional == nil {
		problem.additional = map[string]any{key: value}

		return problem
	}

	problem.additional = maps.Clone(problem.additional)
	problem.additional[key] = value

	return problem
}

// Without removes an additional value to the given key.
func (problem Problem) Without(key string) Problem {
	if problem.additional == nil {
		return problem
	}

	problem.additional = maps.Clone(problem.additional)
	delete(problem.additional, key)

	return problem
}

// Error is the error-like string representation of a [Problem].
func (problem Problem) Error() string {
	if problem.err != nil {
		return fmt.Sprintf("%d %s: %s", problem.Status, http.StatusText(problem.Status), problem.err)
	}

	return problem.Title
}

func (problem Problem) Unwrap() error {
	return problem.err
}

// MarshalJSON replaces the default JSON encoding behaviour.
func (problem Problem) MarshalJSON() ([]byte, error) {
	mapped := make(map[string]any, len(problem.additional)+5)

	mapped["detail"] = problem.Detail
	mapped["instance"] = problem.Instance
	mapped["status"] = problem.Status
	mapped["title"] = problem.Title
	mapped["type"] = problem.Type

	for key, value := range problem.additional {
		mapped[key] = value
	}

	return json.Marshal(mapped)
}

// UnmarshalJSON replaces the default JSON decoding behaviour.
func (problem *Problem) UnmarshalJSON(data []byte) error {
	mapped := make(map[string]any)

	if err := json.Unmarshal(data, &mapped); err != nil {
		return err
	}

	if value, ok := mapped["detail"].(string); ok {
		problem.Detail = value
	}

	if value, ok := mapped["instance"].(string); ok {
		problem.Instance = value
	}

	if value, ok := mapped["status"].(float64); ok {
		problem.Status = int(value)
	}

	if value, ok := mapped["title"].(string); ok {
		problem.Title = value
	}

	if value, ok := mapped["type"].(string); ok {
		problem.Type = value
	}

	delete(mapped, "detail")
	delete(mapped, "instance")
	delete(mapped, "status")
	delete(mapped, "title")
	delete(mapped, "type")

	problem.additional = mapped

	return nil
}

func (problem Problem) controls(request *http.Request) ProblemControls {
	controls := ProblemControls{}

	if c, ok := Problems(request); ok {
		controls = c
	}

	return defaultedProblemControls(controls)
}

func (problem Problem) defaulted(request *http.Request) Problem {
	controls := problem.controls(request)

	if problem.Type == "" {
		problem.Type = controls.DefaultType(problem, request)
	}

	if problem.Status == 0 {
		problem.Status = controls.DefaultStatus(problem, request)
	}

	if problem.Title == "" {
		problem.Title = controls.DefaultTitle(problem, request)
	}

	if problem.Instance == "" {
		problem.Instance = controls.DefaultInstance(problem, request)
	}

	if controls.Lowercase(problem, request) {
		problem.Title = lowercase(problem.Title)
		problem.Detail = lowercase(problem.Detail)
	}

	return problem
}

// Respond implements [Responder] interface to implement
// how a problem responds to an http request.
func (problem Problem) Respond(request *http.Request) Builder {
	controls := problem.controls(request)

	return controls.Response(problem.defaulted(request), request)
}
