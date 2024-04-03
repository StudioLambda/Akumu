package akumu

import (
	"fmt"
	"net/http"
	"strings"
	"unicode"
)

// Problem represents a problem details for HTTP APIs.
// See https://tools.ietf.org/html/rfc7807 for more information.
type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Status   int    `json:"status"`
	Instance string `json:"instance"`
}

func NewProblem(err error, status int) Problem {
	return Problem{
		Detail: err.Error(),
		Status: status,
	}
}

func (problem Problem) Error() string {
	return problem.Title
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

	if strings.Contains(request.Header.Get("Accept"), "text/html") {
		return Response(problem.Status).
			HTML(fmt.Sprintf(
				`<style>.akumu.titlecase{text-transform:capitalize;}.akumu.uppercase-first::first-letter{text-transform:uppercase;}</style><h1 class="akumu titlecase">%s &mdash; %d</h1><h2 class="akumu uppercase-first">%s &mdash; %s</h2><a href=\"%s\">%s</a>`,
				problem.Title, problem.Status, problem.Detail, problem.Instance, problem.Type, problem.Type,
			))
	}

	return Response(problem.Status).
		Text(fmt.Sprintf("%s\n\n%s", problem.Title, problem.Detail))
}
