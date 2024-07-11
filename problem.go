package akumu

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"
	"strings"
	"unicode"
)

// Problem represents a problem details for HTTP APIs.
// See https://tools.ietf.org/html/rfc7807 for more information.
type Problem struct {
	additional map[string]any
	Type       string
	Title      string
	Detail     string
	Status     int
	Instance   string
}

func NewProblemFromError(err error, status int) Problem {
	return Problem{
		additional: make(map[string]any),
		Detail:     err.Error(),
		Status:     status,
	}
}

func (problem Problem) With(key string, value any) Problem {
	if problem.additional == nil {
		problem.additional = map[string]any{key: value}

		return problem
	}

	problem.additional = maps.Clone(problem.additional)
	problem.additional[key] = value

	return problem
}

func (problem Problem) Without(key string) Problem {
	if problem.additional == nil {
		return problem
	}

	problem.additional = maps.Clone(problem.additional)
	delete(problem.additional, key)

	return problem
}

func (problem Problem) Error() string {
	return problem.Title
}

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

func (problem Problem) Respond(request *http.Request) Builder {
	lowercase := func(str string) string {
		result := ""

		for _, r := range str {
			result += string(unicode.ToLower(r))
		}

		return result
	}

	if problem.Type == "" {
		problem.Type = "about:blank"
	}

	if problem.Status == 0 {
		problem.Status = http.StatusInternalServerError
	}

	if problem.Title == "" {
		problem.Title = http.StatusText(problem.Status)
	}

	if problem.Instance == "" {
		problem.Instance = request.URL.String()
	}

	problem.Title = lowercase(problem.Title)
	problem.Detail = lowercase(problem.Detail)

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
