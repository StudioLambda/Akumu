package akumu

import "unicode"

func lowercase(str string) string {
	result := make([]byte, len(str))

	for i, r := range str {
		result[i] += byte(unicode.ToLower(r))
	}

	return string(result)
}
