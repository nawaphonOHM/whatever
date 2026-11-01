package detector

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/urlutil"
)

const (
	minPort = 1
	maxPort = 65535
)

// checkTargetRaw validates raw target string before URL parsing.
func checkTargetRaw(target string) error {
	if target == "" {
		return errors.New("invalid URL: empty target")
	}
	if strings.ContainsAny(target, " \t\r\n") {
		return errors.New("invalid URL: contains whitespace")
	}
	if strings.HasPrefix(target, "/") || strings.HasPrefix(target, "\\") {
		return errors.New("invalid URL: leading slash")
	}
	if urlutil.HasSchemeSingleSlash(target) {
		return errors.New("invalid URL: invalid scheme format")
	}
	return nil
}

// normalizeTargetScheme ensures the raw target string has a scheme prefix for URL parsing.
func normalizeTargetScheme(target string) string {
	if !urlutil.HasSchemePrefix(target) {
		return "https://" + target
	}
	return target
}

// parseRawURL transforms and parses raw target into *url.URL.
func parseRawURL(target string) (*url.URL, error) {
	raw := normalizeTargetScheme(target)
	u, err := url.Parse(raw)
	if err != nil {
		return nil, errors.New("invalid URL format")
	}
	return u, nil
}

// validateTargetScheme ensures scheme is http or https.
func validateTargetScheme(scheme string) error {
	lower := strings.ToLower(scheme)
	if lower == "http" || lower == "https" {
		return nil
	}
	return fmt.Errorf("unsupported URL scheme %q: only http and https are supported", scheme)
}

// isPortOutOfRange checks if numeric port is beyond allowed boundaries.
func isPortOutOfRange(p int) bool { return p < minPort || p > maxPort }

// parsePortNumber parses and checks port integer value.
func parsePortNumber(portStr string) error {
	p, err := strconv.Atoi(portStr)
	if err != nil {
		return errors.New("invalid port in URL: must be numeric")
	}
	if isPortOutOfRange(p) {
		return fmt.Errorf("invalid port in URL %d: must be between %d and %d", p, minPort, maxPort)
	}
	return nil
}

// validateTargetPort ensures embedded URL port is within valid range (1-65535).
func validateTargetPort(portStr string) error {
	if portStr == "" {
		return nil
	}
	return parsePortNumber(portStr)
}

// validateTargetComponents checks scheme, hostname, and port of parsed URL.
func validateTargetComponents(u *url.URL) error {
	if err := validateTargetScheme(u.Scheme); err != nil {
		return err
	}
	if u.Hostname() == "" {
		return errors.New("invalid URL: missing host")
	}
	if u.Host == "http:" || u.Host == "https:" || u.Hostname() == "http" || u.Hostname() == "https" {
		return errors.New("invalid URL: invalid host")
	}
	// u.Host (not u.Hostname()) keeps IPv6 brackets, so "[::1]" is not mistaken for a host starting with ':'.
	if strings.HasSuffix(u.Host, ":") || strings.HasPrefix(u.Host, ":") {
		return errors.New("invalid URL: invalid host")
	}
	// "@" is legal in paths; only reject the ambiguous shape where a nested "scheme://userinfo@host" sits in the path.
	if strings.Contains(u.Path, "@") && strings.Contains(u.Path, "://") {
		return errors.New("invalid URL: invalid path credentials")
	}
	if u.Opaque != "" {
		return errors.New("invalid URL: invalid host")
	}
	return validateTargetPort(u.Port())
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
