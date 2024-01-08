package http

import "net/url"

type URL struct {
	*url.URL
}

func NewURL(url *url.URL) *URL {
	return &URL{
		URL: url,
	}
}

func ParseURL(raw string) (*URL, error) {
	url, err := url.Parse(raw)

	if err != nil {
		return nil, err
	}

	return NewURL(url), nil
}

func (url *URL) String() string {
	return url.URL.String()
}

func (url *URL) Host() string {
	return url.URL.Host
}

func (url *URL) Path() string {
	return url.URL.Path
}

func (url *URL) Fragment() string {
	return url.URL.Fragment
}

func (url *URL) SetFragment(fragment string) {
	url.URL.Fragment = fragment
}

func (url *URL) RawQuery() string {
	return url.URL.RawQuery
}

func (url *URL) HasQuery(key string) bool {
	return url.URL.Query().Has(key)
}

func (url *URL) Query(key string) string {
	return url.URL.Query().Get(key)
}

func (url *URL) QueryAll(key string) []string {
	return url.URL.Query()[key]
}

func (url *URL) SetQuery(key, value string) {
	query := url.URL.Query()

	query.Set(key, value)

	url.URL.RawQuery = query.Encode()
}

func (url *URL) AppendQuery(key, value string) {
	query := url.URL.Query()

	query.Add(key, value)

	url.URL.RawQuery = query.Encode()
}
