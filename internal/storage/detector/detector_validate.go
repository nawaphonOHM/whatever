package detector

import (
	"errors"
	"net/url"
	"strings"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/urlutil"
)

func checkTargetRawBasic(target string) error {
	if target == "" {
		return errors.New("invalid URL: empty target")
	}
	if strings.ContainsAny(target, " \t\r\n") {
		return errors.New("invalid URL: contains whitespace")
	}
	return nil
}

func hasLeadingSlash(target string) bool {
	return strings.HasPrefix(target, "/") || strings.HasPrefix(target, "\\")
}

func checkTargetRawPrefixes(target string) error {
	if hasLeadingSlash(target) {
		return errors.New("invalid URL: leading slash")
	}
	if urlutil.HasSchemeSingleSlash(target) {
		return errors.New("invalid URL: invalid scheme format")
	}
	return nil
}

func checkTargetRaw(target string) error {
	if err := checkTargetRawBasic(target); err != nil {
		return err
	}
	return checkTargetRawPrefixes(target)
}

func normalizeTargetScheme(target string) string {
	if !urlutil.HasSchemePrefix(target) {
		return "https://" + target
	}
	return target
}

func parseRawURL(target string) (*url.URL, error) {
	raw := normalizeTargetScheme(target)
	u, err := url.Parse(raw)
	if err != nil {
		return nil, errors.New("invalid URL format")
	}
	return u, nil
}

// ParseTarget validates raw target string and returns parsed URL with non-empty host and valid scheme and port.
func ParseTarget(target string) (*url.URL, error) {
	if err := checkTargetRaw(target); err != nil {
		return nil, err
	}
	u, err := parseRawURL(target)
	if err != nil {
		return nil, err
	}
	return u, validateTargetComponents(u)
}
