package urlutil

import (
	"errors"
	"net/url"
	"strings"
)

func sanitizeParsedURL(u *url.URL) {
	u.User = nil
	u.RawQuery = ""
	u.Fragment = ""
	u.RawFragment = ""
}

func parseSchemelessURL(raw string) (*url.URL, error) {
	if strings.HasPrefix(raw, "/") && strings.Contains(raw, atChar) {
		return nil, errors.New("invalid schemeless url")
	}
	return url.Parse(defaultSchemePrefix + raw)
}

func stripDefaultScheme(res string) string {
	if after, found := strings.CutPrefix(res, defaultSchemePrefix); found {
		return after
	}
	return strings.TrimPrefix(res, "https:")
}

func sanitizeSchemelessURL(raw string) string {
	u, err := parseSchemelessURL(raw)
	if err != nil {
		return UnparseableURLPlaceholder
	}
	if isAmbiguousMalformed(raw, u) {
		return UnparseableURLPlaceholder
	}
	sanitizeParsedURL(u)
	return stripDefaultScheme(u.String())
}

func sanitizeAbsoluteURL(raw string) string {
	u, err := url.Parse(raw)
	if err != nil {
		return UnparseableURLPlaceholder
	}
	if isAmbiguousMalformed(raw, u) {
		return UnparseableURLPlaceholder
	}
	sanitizeParsedURL(u)
	return u.String()
}
