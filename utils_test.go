package akumu

import (
	"testing"
)

func TestUtilsLowercase(t *testing.T) {
	word := "somEthing-IsNot_lowercased2"
	lowercased := lowercase(word)

	if expected := "something-isnot_lowercased2"; lowercased != expected {
		t.Fatalf("expected '%s' but got '%s'", expected, lowercased)
	}
}
