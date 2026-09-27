package config

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)

// checkEndpointNonEmpty verifies endpoint string is non-empty.
func checkEndpointNonEmpty(endpoint string) error {
	if strings.TrimSpace(endpoint) == "" {
		return ErrEmptyEndpoint
	}
	return nil
}

// checkEndpointWhitespace verifies endpoint contains no whitespace.
func checkEndpointWhitespace(endpoint string) error {
	if strings.ContainsAny(endpoint, " \t\r\n") {
		return fmt.Errorf("%w: endpoint contains whitespace: %q", ErrInvalidEndpoint, endpoint)
	}
	return nil
}

// isPortInRange checks if port integer falls in [MinPort, MaxPort].
func isPortInRange(p int) bool {
	if p < MinPort {
		return false
	}
	return p <= MaxPort
}

// validatePortNumber validates that portStr parses to an integer within valid range.
func validatePortNumber(portStr string) error {
	p, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("%w: invalid port %q", ErrInvalidEndpoint, portStr)
	}
	if !isPortInRange(p) {
		return fmt.Errorf("%w: port %d out of range", ErrInvalidEndpoint, p)
	}
	return nil
}

// validateHostPortPair validates host and port when split succeeds.
func validateHostPortPair(host, port, endpoint string) error {
	if strings.TrimSpace(host) == "" {
		return fmt.Errorf("%w: endpoint has empty host: %q", ErrInvalidEndpoint, endpoint)
	}
	if port == "" {
		return nil
	}
	return validatePortNumber(port)
}

// validatePlainHost validates plain host format.
func validatePlainHost(endpoint string) error {
	u, err := url.Parse("http://" + endpoint)
	if err != nil {
		return fmt.Errorf("%w: invalid endpoint format: %q", ErrInvalidEndpoint, endpoint)
	}
	if u.Hostname() == "" {
		return fmt.Errorf("%w: invalid endpoint format: %q", ErrInvalidEndpoint, endpoint)
	}
	return nil
}

// validateEndpointHost validates endpoint formatted as host:port or plain host.
func validateEndpointHost(endpoint string) error {
	host, port, err := net.SplitHostPort(endpoint)
	if err == nil {
		return validateHostPortPair(host, port, endpoint)
	}
	return validatePlainHost(endpoint)
}
