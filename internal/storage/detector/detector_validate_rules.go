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

func validateTargetScheme(scheme string) error {
	lower := strings.ToLower(scheme)
	if lower == "http" || lower == "https" {
		return nil
	}
	return fmt.Errorf("unsupported URL scheme %q: only http and https are supported", scheme)
}

func isPortOutOfRange(p int) bool { return p < minPort || p > maxPort }

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

func validateTargetPort(portStr string) error {
	if portStr == "" {
		return nil
	}
	return parsePortNumber(portStr)
}

func isTargetBareName(name string) bool {
	return name == "http" || name == "https"
}

func isBareTargetHost(u *url.URL) bool {
	if isTargetBareName(u.Hostname()) {
		return true
	}
	return u.Host == "http:" || u.Host == "https:"
}

func isColonTargetHost(host string) bool {
	// The caller passes u.Host to preserve IPv6 brackets, so [::1] is not
	// mistaken for a host beginning with a colon.
	return strings.HasSuffix(host, ":") || strings.HasPrefix(host, ":")
}

func isInvalidHostFormat(u *url.URL) bool {
	return isBareTargetHost(u) || isColonTargetHost(u.Host)
}

func checkHostFormat(u *url.URL) error {
	if u.Opaque != "" || isInvalidHostFormat(u) {
		return errors.New("invalid URL: invalid host")
	}
	return nil
}

func validateTargetHost(u *url.URL) error {
	if u.Hostname() == "" {
		return errors.New("invalid URL: missing host")
	}
	return checkHostFormat(u)
}

func validateTargetPath(path string) error {
	// A literal @ is valid in a path; reject the nested scheme://userinfo@host
	// shape, where parsing can otherwise hide credentials outside the authority.
	if strings.Contains(path, "@") && strings.Contains(path, "://") {
		return errors.New("invalid URL: invalid path credentials")
	}
	return nil
}

func validateTargetStructure(u *url.URL) error {
	if err := validateTargetScheme(u.Scheme); err != nil {
		return err
	}
	return validateTargetHost(u)
}

func validateTargetComponents(u *url.URL) error {
	if err := validateTargetStructure(u); err != nil {
		return err
	}
	if err := validateTargetPath(u.Path); err != nil {
		return err
	}
	return validateTargetPort(u.Port())
}
