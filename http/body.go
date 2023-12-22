package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
)

type Body struct {
	reader io.Reader
	cached []byte
}

var (
	ErrBodyBytes      = errors.New("unable to read body bytes")
	ErrBodyJSON       = errors.New("unable to read body json")
	ErrBodyStrictJSON = errors.New("unable to read body strict json")
	ErrBodyReader     = errors.New("unable to create body reader")
	ErrBodyNilReader  = errors.New("reader is nil")
)

func NewBody(reader io.Reader) Body {
	return Body{
		reader: reader,
		cached: nil,
	}
}

// Returns true if the body has been cached.
func (body *Body) IsCached() bool {
	return body.cached != nil
}

// Returns a reader that can be used to read the body.
// This method takes into account possible errors that
// may occur when reading the body.
func (body *Body) TryReader() (io.Reader, error) {
	b, err := body.Bytes()

	if err != nil {
		return nil, errors.Join(ErrBodyReader, err)
	}

	return bytes.NewReader(b), nil
}

// Returns a reader that can be used to read the body.
// This method does not take into account possible errors that
// may occur when reading the body. This method won't panic but
// it will return a reader with an empty body if an error occurs.
func (body *Body) Reader() io.Reader {
	b, _ := body.Bytes()

	return bytes.NewReader(b)
}

// Returns the body as a byte slice.
// This method takes into account possible errors that
// may occur when reading the body.
func (body *Body) Bytes() ([]byte, error) {
	if body.IsCached() {
		return body.cached, nil
	}

	if body.reader == nil {
		body.cached = []byte{}

		return body.cached, nil
	}

	bytes, err := io.ReadAll(body.reader)

	if err != nil {
		return nil, errors.Join(ErrBodyBytes, err)
	}

	body.cached = bytes

	return body.cached, nil
}

// Tries to convert the body to a string.
// This method takes into account possible errors that
// may occur when reading the body.
func (body *Body) TryString() (string, error) {
	bytes, err := body.Bytes()

	if err != nil {
		return "", errors.Join(ErrBodyBytes, err)
	}

	return string(bytes), nil
}

// Returns the body as a string.
// This method does not take into account possible errors that
// may occur when reading the body. This method won't panic but
// it will return an empty string if an error occurs.
//
// Mainly used to satisfy the fmt.Stringer interface.
func (body *Body) String() string {
	string, _ := body.TryString()

	return string
}

func (body *Body) JSON(result any) error {
	bytes, err := body.Bytes()

	if err != nil {
		return errors.Join(ErrBodyJSON, err)
	}

	if err := json.Unmarshal(bytes, result); err != nil {
		return errors.Join(ErrBodyJSON, err)
	}

	return nil
}

func (body *Body) StrictJSON(result any) error {
	reader, err := body.TryReader()

	if err != nil {
		return errors.Join(ErrBodyStrictJSON, err)
	}

	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(result); err != nil {
		return errors.Join(ErrBodyStrictJSON, err)
	}

	return nil
}

// Returns a binary representation of the body.
func (body *Body) MarshalBinary() ([]byte, error) {
	return body.Bytes()
}

// Sets the body to the given binary representation.
func (body *Body) UnmarshalBinary(data []byte) error {
	body.cached = data

	return nil
}
