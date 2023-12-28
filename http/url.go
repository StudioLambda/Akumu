package http

import "net/url"

type URL struct {
	native url.URL
}

type MutableURL struct {
	URL
}

func NewURL(url *url.URL) URL {
	return URL{
		native: *url,
	}
}

func NewMutableURL(url *url.URL) MutableURL {
	return MutableURL{
		URL: NewURL(url),
	}
}

func ParseURL(raw string) (URL, error) {
	url, err := url.Parse(raw)

	if err != nil {
		return URL{}, err
	}

	return NewURL(url), nil
}

func (url URL) String() string {
	return url.native.String()
}

func (url URL) Host() string {
	return url.native.Host
}

func (url URL) Path() string {
	return url.native.Path
}

func (url URL) Fragment() string {
	return url.native.Fragment
}

func (url *MutableURL) SetFragment(fragment string) {
	url.native.Fragment = fragment
}

func (url URL) RawQuery() string {
	return url.native.RawQuery
}

func (url URL) HasQuery(key string) bool {
	return url.native.Query().Has(key)
}

func (url URL) Query(key string) string {
	return url.native.Query().Get(key)
}

func (url URL) QueryAll(key string) []string {
	return url.native.Query()[key]
}

func (url *MutableURL) SetQuery(key, value string) {
	query := url.native.Query()

	query.Set(key, value)

	url.native.RawQuery = query.Encode()
}

func (url *MutableURL) AppendQuery(key, value string) {
	query := url.native.Query()

	query.Add(key, value)

	url.native.RawQuery = query.Encode()
}
