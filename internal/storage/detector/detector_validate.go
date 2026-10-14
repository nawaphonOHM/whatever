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
		return fmt.Errorf("invalid URL: contains whitespace: %q", target)
	}
	return nil
}

// normalizeTargetScheme ensures the raw target string has a scheme prefix for URL parsing.
func normalizeTargetScheme(target string) string {
	if !strings.Contains(target, "://") {
		return "https://" + target
	}
	return target
}

// parseRawURL transforms and parses raw target into *url.URL.
func parseRawURL(target string) (*url.URL, error) {
	raw := normalizeTargetScheme(target)
	u, err := url.Parse(raw)
	if err != nil {
		return nil, fmt.Errorf("invalid URL format: %w", err)
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
	if err != nil || isPortOutOfRange(p) {
		return fmt.Errorf("invalid port in URL %q: must be between %d and %d", portStr, minPort, maxPort)
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
func validateTargetComponents(u *url.URL, target string) error {
	if err := validateTargetScheme(u.Scheme); err != nil {
		return err
	}
	if u.Hostname() == "" {
		return fmt.Errorf("invalid URL: missing host: %q", target)
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
	return u, validateTargetComponents(u, target)
}
