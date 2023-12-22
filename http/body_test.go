package http_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/studiolambda/akumu/http"
)

func TestBodyBytes(t *testing.T) {
	t.Run("CanRead", func(t *testing.T) {
		body := http.NewBody(strings.NewReader("Hello World"))

		bytes, err := body.Bytes()

		require.NoError(t, err)
		require.Len(t, bytes, 11)
		require.Equal(t, []byte("Hello World"), bytes)
	})

	t.Run("MultipleReads", func(t *testing.T) {
		body := http.NewBody(strings.NewReader("Hello World"))

		body.Bytes()
		body.Bytes()
		bytes, err := body.Bytes()

		require.NoError(t, err)
		require.Len(t, bytes, 11)
		require.Equal(t, []byte("Hello World"), bytes)
	})

	t.Run("FailsWhenInvalidReader", func(t *testing.T) {
		body := http.NewBody(nil)

		bytes, err := body.Bytes()

		require.NoError(t, err)
		require.Equal(t, []byte{}, bytes)
	})
}
