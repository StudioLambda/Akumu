package utils

import (
	"mime"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type acceptPair struct {
	media   string
	quality float64
}

type Accept struct {
	values []acceptPair
}

func ParseAccept(request *http.Request) Accept {
	accept := Accept{
		values: make([]acceptPair, 0),
	}

	for _, header := range request.Header.Values("Accept") {
		for _, line := range strings.Split(header, ",") {
			media, parameters, err := mime.ParseMediaType(line)

			if err != nil {
				continue
			}

			quality := 1.0

			if param, ok := parameters["q"]; ok {
				if q, err := strconv.ParseFloat(param, 64); err == nil {
					quality = q
				}
			}

			accept.values = append(accept.values, acceptPair{
				media:   media,
				quality: quality,
			})
		}
	}

	return accept
}

func (accept Accept) find(media string) (acceptPair, bool) {
	for _, pair := range accept.values {
		if media == pair.media {
			return pair, true
		}

		// Test for wildcard in media type
		if strings.Contains(media, "/*") {
			// Compare only the first part.
			if trimmed := strings.TrimSuffix(media, "/*"); strings.HasPrefix(pair.media, trimmed) {
				return pair, true
			}
		}

		// Test for wildcard in accept media type
		if strings.Contains(pair.media, "/*") {
			// Compare only the first part.
			if trimmed := strings.TrimSuffix(pair.media, "/*"); strings.HasPrefix(media, trimmed) {
				return pair, true
			}
		}
	}

	return acceptPair{}, false
}

func (accept Accept) Accepts(media string) bool {
	_, found := accept.find(media)

	return found
}

func (accept Accept) Quality(media string) float64 {
	if pair, found := accept.find(media); found {
		return pair.quality
	}

	return 0
}

func (accept Accept) Order() []string {
	keys := make([]string, len(accept.values))

	for i, pair := range accept.values {
		keys[i] = pair.media
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return accept.values[i].quality > accept.values[j].quality
	})

	return keys
}
