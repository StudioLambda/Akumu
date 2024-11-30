package akumu

import (
	"errors"
	"fmt"
	"testing"
)

func TestUtilsStackTrace(t *testing.T) {
	err := errors.Join(
		errors.New("first"),
		fmt.Errorf("%w: foo", errors.New("second")),
		errors.Join(errors.New("third"), errors.New("last")),
	)

	trace := stackTrace(err)

	if expected := 4; len(trace) != expected {
		t.Fatalf("expected len '%d' but got '%d'", expected, len(trace))
	}
}
