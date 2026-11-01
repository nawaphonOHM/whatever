package detector

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
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
	return nil
}

// hasSchemePrefix reports whether target has a leading scheme prefix before delimiters.
func hasSchemePrefix(target string) bool {
	idx := strings.Index(target, "://")
	if idx <= 0 {
		return false
	}
	schemePart := target[:idx]
	if strings.ContainsAny(schemePart, "/?#@ \t\r\n") {
		return false
	}
	return true
}

// normalizeTargetScheme ensures the raw target string has a scheme prefix for URL parsing.
func normalizeTargetScheme(target string) string {
	if !hasSchemePrefix(target) {
		return "https://" + target
	}
	return target
}

// parseRawURL transforms and parses raw target into *url.URL.
func parseRawURL(target string) (*url.URL, error) {
	raw := normalizeTargetScheme(target)
	schemeIdx := strings.Index(raw, "://")
	afterScheme := raw[schemeIdx+3:]
	if atIdx := strings.Index(afterScheme, "@"); atIdx != -1 {
		termIdx := strings.IndexAny(afterScheme, "/?#")
		if termIdx != -1 && termIdx < atIdx {
			return nil, errors.New("invalid URL format")
		}
		if strings.Contains(afterScheme[atIdx+1:], "@") {
			return nil, errors.New("invalid URL format")
		}
	}

	u, err := url.Parse(raw)
	if err != nil {
		if strings.Contains(target, "@") {
			return nil, errors.New("invalid URL format")
		}
		var urlErr *url.Error
		if errors.As(err, &urlErr) {
			return nil, fmt.Errorf("invalid URL format: %w", urlErr.Err)
		}
		return nil, fmt.Errorf("invalid URL format: %w", err)
	}
	if strings.Contains(u.Path, "@") || strings.Contains(u.RawQuery, "@") || strings.Contains(u.Fragment, "@") {
		return nil, errors.New("invalid URL format")
	}
	return u, nil
}

// validateTargetScheme ensures scheme is http or https.
func validateTargetScheme(scheme, target string) error {
	lower := strings.ToLower(scheme)
	if lower == "http" || lower == "https" {
		return nil
	}
	if strings.Contains(target, "@") {
		return errors.New("unsupported URL scheme: only http and https are supported")
	}
	return fmt.Errorf("unsupported URL scheme %q: only http and https are supported", scheme)
}

// isPortOutOfRange checks if numeric port is beyond allowed boundaries.
func isPortOutOfRange(p int) bool { return p < minPort || p > maxPort }

// parsePortNumber parses and checks port integer value.
func parsePortNumber(portStr, target string) error {
	p, err := strconv.Atoi(portStr)
	if err != nil || isPortOutOfRange(p) {
		if strings.Contains(target, "@") {
			return fmt.Errorf("invalid port in URL: must be between %d and %d", minPort, maxPort)
		}
		return fmt.Errorf("invalid port in URL %q: must be between %d and %d", portStr, minPort, maxPort)
	}
	return nil
}

// validateTargetPort ensures embedded URL port is within valid range (1-65535).
func validateTargetPort(portStr, target string) error {
	if portStr == "" {
		return nil
	}
	return parsePortNumber(portStr, target)
}

// validateTargetComponents checks scheme, hostname, and port of parsed URL.
func validateTargetComponents(u *url.URL, target string) error {
	if err := validateTargetScheme(u.Scheme, target); err != nil {
		return err
	}
	if u.Hostname() == "" {
		return errors.New("invalid URL: missing host")
	}
	return validateTargetPort(u.Port(), target)
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
	return u, validateTargetComponents(u, target)
}
