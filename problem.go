package akumu

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"strings"
)

// Problem represents a problem details for HTTP APIs.
// See https://datatracker.ietf.org/doc/html/rfc9457 for more information.
type Problem struct {
	additional map[string]any
	Type       string
	Title      string
	Detail     string
	Status     int
	Instance   string
}

type ProblemControlsResolver[R any] func(problem Problem, request *http.Request) R

type ProblemControls struct {
	Lowercase       ProblemControlsResolver[bool]
	DefaultStatus   ProblemControlsResolver[int]
	DefaultType     ProblemControlsResolver[string]
	DefaultTitle    ProblemControlsResolver[string]
	DefaultInstance ProblemControlsResolver[string]
}

// ProblemsKey is the context key where the
// problem controls are stored in the request.
type ProblemsKey struct{}

// Problems return the [ProblemControls] used to determine
// how [Problem] respond to http requests.
func Problems(request *http.Request) (ProblemControls, bool) {
	controls, ok := request.
		Context().
		Value(ProblemsKey{}).(ProblemControls)

	return controls, ok
}

func defaultProblemControllerLowercase(problem Problem, request *http.Request) bool {
	return true
}

func defaultProblemControllerStatus(problem Problem, request *http.Request) int {
	return http.StatusInternalServerError
}

func defaultProblemControllerType(problem Problem, request *http.Request) string {
	return "about:blank"
}

func defaultProblemControllerTitle(problem Problem, request *http.Request) string {
	return http.StatusText(problem.Status)
}

func defaultProblemControllerInstance(problem Problem, request *http.Request) string {
	return request.URL.String()
}

func NewProblemControls() ProblemControls {
	return ProblemControls{
		Lowercase:       defaultProblemControllerLowercase,
		DefaultStatus:   defaultProblemControllerStatus,
		DefaultType:     defaultProblemControllerType,
		DefaultTitle:    defaultProblemControllerTitle,
		DefaultInstance: defaultProblemControllerInstance,
	}
}

// NewProblem creates a new [Problem] from
// the given error and status code.
func NewProblem(err error, status int) Problem {
	return Problem{
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
	return problem.Title
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
	if controls, ok := Problems(request); ok {
		return controls
	}

	return NewProblemControls()
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
	problem = problem.defaulted(request)

	if strings.Contains(request.Header.Get("Accept"), "application/problem+json") || strings.Contains(request.Header.Get("Accept"), "application/json") {
		return Response(problem.Status).
			JSON(problem).
			Header("Content-Type", "application/problem+json")
	}

	// todo: disabled due to improvement schedule
	// if strings.Contains(request.Header.Get("Accept"), "text/html") {
	// 	return Response(problem.Status).
	// 		HTML(fmt.Sprintf(
	// 			`<style>.akumu.titlecase{text-transform:capitalize;}.akumu.uppercase-first::first-letter{text-transform:uppercase;}</style><h1 class="akumu titlecase">%s &mdash; %d</h1><h2 class="akumu uppercase-first">%s &mdash; %s</h2><a href=\"%s\">%s</a>`,
	// 			problem.Title, problem.Status, problem.Detail, problem.Instance, problem.Type, problem.Type,
	// 		))
	// }

	return Response(problem.Status).
		Text(fmt.Sprintf("%s\n\n%s", problem.Title, problem.Detail))
}
