package akumu

// Problem represents a problem details for HTTP APIs.
// See https://tools.ietf.org/html/rfc7807 for more information.
type Problem struct {
	Type     string `json:"type"`
	Title    string `json:"title"`
	Detail   string `json:"detail"`
	Status   int    `json:"status"`
	Instance string `json:"instance"`
}

func (problem Problem) Error() string {
	return problem.Title
}
