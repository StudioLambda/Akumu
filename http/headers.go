package http

import (
	"net/http"
)

type Headers http.Header

func (headers Headers) Has(key string) bool {
	_, ok := headers[key]

	return ok
}

func (headers Headers) Contains(key, value string) bool {
	for _, val := range headers.All(key) {
		if val == value {
			return true
		}
	}

	return false
}

func (headers Headers) First(key string) string {
	return http.Header(headers).Get(key)
}

func (headers Headers) All(key string) []string {
	return http.Header(headers).Values(key)
}

func (headers Headers) Delete(key string) {
	http.Header(headers).Del(key)
}

func (headers Headers) Insert(key string, value string) {
	http.Header(headers).Set(key, value)
}

func (headers Headers) Append(key string, value string) {
	http.Header(headers).Add(key, value)
}

func (headers Headers) Clone() Headers {
	clone := http.Header(headers).Clone()

	return Headers(clone)
}
