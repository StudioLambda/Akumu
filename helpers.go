package akumu

import (
	"encoding/json"
	"net/http"
)

// JSON decodes the given request payload into `T`
//
// This is very usefull for cases where you want
// to quickly take care of decoding JSON payloads into
// specific types. It automatically disallows unknown
// fields and uses [json.Decoder] with the [http.Request.Body].
func JSON[T any](request *http.Request) (T, error) {
	result := *new(T)

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&result); err != nil {
		return result, err
	}

	return result, nil
}
