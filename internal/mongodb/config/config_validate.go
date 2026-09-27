package config

import (
	"errors"
)

// Port boundary constants for MongoDB connection.
const (
	MinPort = 1
	MaxPort = 65535
)

// validateHost checks that Host is non-empty.
func (c *Config) validateHost() error {
	if c.Host == "" {
		return errors.New("mongodb host cannot be empty")
	}
	return nil
}

// credentialMismatchError generates an error when only one credential field is set.
func (c *Config) credentialMismatchError() error {
	if c.Username != "" {
		return errors.New("mongodb password cannot be empty when username is provided")
	}
	return errors.New("mongodb username cannot be empty when password is provided")
}

// validateCredentials checks that username and password are provided as a consistent pair.
func (c *Config) validateCredentials() error {
	hasUser := c.Username != ""
	hasPass := c.Password != ""
	if hasUser != hasPass {
		return c.credentialMismatchError()
	}
	return nil
}

// checkPortRange checks that the port is within the valid range (1-65535).
func (c *Config) checkPortRange() error {
	if c.Port < MinPort || c.Port > MaxPort {
		return errors.New("mongodb port must be between 1 and 65535")
	}
	return nil
}

// isSrvZeroPort checks if protocol is mongodb+srv with port omitted.
func (c *Config) isSrvZeroPort() bool {
	if c.Protocol != ProtocolMongoDBSrv {
		return false
	}
	return c.Port == 0
}

// validatePort checks that the port is within the valid range (1-65535).
// When Protocol is mongodb+srv, Port is optional (0 is allowed).
func (c *Config) validatePort() error {
	if c.isSrvZeroPort() {
		return nil
	}
	return c.checkPortRange()
}

// validateIdentity checks host and credentials configuration.
func (c *Config) validateIdentity() error {
	if err := c.validateHost(); err != nil {
		return err
	}
	return c.validateCredentials()
}

// Validate checks that the configuration values are valid.
func (c *Config) Validate() error {
	if c == nil {
		return ErrNilConfig
	}
	if err := c.validateIdentity(); err != nil {
		return err
	}
	return c.validateSettings()
}
