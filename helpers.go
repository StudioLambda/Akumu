package akumu

import (
	"encoding/json"
	"net/http"
)

// JSON decodes the given request payload into `T`
func JSON[T any](request *http.Request) (T, error) {
	result := *new(T)

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}
