package akumu

import "unicode"

// lowercase makes the input string a lowercased
// string using [unicode.ToLower] on each character.
func lowercase(str string) string {
	result := make([]byte, len(str))

	for i, r := range str {
		result[i] += byte(unicode.ToLower(r))
	}

	return string(result)
}
