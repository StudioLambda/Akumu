package http

import (
	"encoding/json"
	"errors"
	"io"
)

type Body struct {
	reader io.Reader
	cached []byte
}

var (
	ErrBodyBytes     = errors.New("unable to read body bytes")
	ErrBodyJSON      = errors.New("unable to read body json")
	ErrBodyNilReader = errors.New("reader is nil")
)

func NewBody(reader io.Reader) *Body {
	return &Body{
		reader: reader,
	}
}

func (body *Body) IsCached() bool {
	return body.cached != nil
}

func (body *Body) Bytes() ([]byte, error) {
	if body.IsCached() {
		return body.cached, nil
	}

	if body.reader == nil {
		return nil, errors.Join(ErrBodyBytes, ErrBodyNilReader)
	}

	bytes, err := io.ReadAll(body.reader)

	if err != nil {
		return nil, errors.Join(ErrBodyBytes, err)
	}

	body.cached = bytes

	return bytes, nil
}

func (body *Body) String() (string, error) {
	bytes, err := body.Bytes()

	if err != nil {
		return "", errors.Join(ErrBodyBytes, err)
	}

	return string(bytes), nil
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
