package config

import (
	"fmt"
	"net/url"
	"strings"
)

// hasURLTarget checks if the parsed URL has a host or path target.
func hasURLTarget(u *url.URL) bool {
	if u.Host != "" {
		return true
	}
	return u.Path != ""
}

// validateURLTarget ensures scheme and host or path target are present.
func validateURLTarget(u *url.URL) error {
	if u.Scheme == "" {
		return fmt.Errorf("%w: endpoint URL missing scheme", ErrInvalidEndpoint)
	}
	if !hasURLTarget(u) {
		return fmt.Errorf("%w: endpoint URL missing host/target", ErrInvalidEndpoint)
	}
	return nil
}

// validateURLPort checks the port component of parsed URL if present.
func validateURLPort(u *url.URL) error {
	if u.Port() == "" {
		return nil
	}
	return validatePortNumber(u.Port())
}

// validateEndpointURL parses and validates URL-formatted endpoint.
func validateEndpointURL(endpoint string) error {
	u, err := url.Parse(endpoint)
	if err != nil {
		return fmt.Errorf("%w: failed to parse endpoint URL: %w", ErrInvalidEndpoint, err)
	}
	if err := validateURLTarget(u); err != nil {
		return err
	}
	return validateURLPort(u)
}

// checkEndpointFormat dispatches endpoint validation based on URL scheme presence.
func checkEndpointFormat(endpoint string) error {
	if strings.Contains(endpoint, "://") {
		return validateEndpointURL(endpoint)
	}
	return validateEndpointHost(endpoint)
}

// validateEndpoint validates that Endpoint is a valid host, host:port, or URL.
func (c *Config) validateEndpoint() error {
	trimmed := strings.TrimSpace(c.Endpoint)
	if err := checkEndpointNonEmpty(trimmed); err != nil {
		return err
	}
	if err := checkEndpointWhitespace(trimmed); err != nil {
		return err
	}
	return checkEndpointFormat(trimmed)
}
