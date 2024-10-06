package akumu

import "unicode"

func lowercase(str string) string {
	result := ""

	for _, r := range str {
		result += string(unicode.ToLower(r))
	}

	return result
}
