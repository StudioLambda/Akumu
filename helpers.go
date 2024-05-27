package akumu

import (
	"encoding/json"
	"io"
)

func JSON[T any](reader io.Reader) (T, error) {
	result := *new(T)

	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}
