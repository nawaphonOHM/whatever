package config

import (
	"errors"
	"fmt"
)

// Port boundary constants for MongoDB connection.
const (
	MinPort = 1
	MaxPort = 65535
)

// validateHostAndUser checks that Host and Username are non-empty.
func (c *Config) validateHostAndUser() error {
	if c.Host == "" {
		return errors.New("mongodb host cannot be empty")
	}
	if c.Username == "" {
		return errors.New("mongodb username cannot be empty")
	}
	return nil
}

// validatePassword checks that Password is non-empty.
func (c *Config) validatePassword() error {
	if c.Password == "" {
		return errors.New("mongodb password cannot be empty")
	}
	return nil
}

// validateRequiredFields checks that Host, Username, and Password are non-empty.
func (c *Config) validateRequiredFields() error {
	if err := c.validateHostAndUser(); err != nil {
		return err
	}
	return c.validatePassword()
}

// validatePort checks that the port is within the valid range (1-65535).
func (c *Config) validatePort() error {
	if c.Port < MinPort || c.Port > MaxPort {
		return errors.New("mongodb port must be between 1 and 65535")
	}
	return nil
}

// validateProtocol verifies that the protocol is either mongodb or mongodb+srv.
func (c *Config) validateProtocol() error {
	switch c.Protocol {
	case ProtocolMongoDB, ProtocolMongoDBSrv:
		return nil
	default:
		return fmt.Errorf("invalid mongodb protocol: %s", c.Protocol)
	}
}

// validateUUIDRepresentation checks that the UUID representation is a supported format.
func (c *Config) validateUUIDRepresentation() error {
	switch c.UUIDRepresentation {
	case UUIDRepresentationUnspecified,
		UUIDRepresentationStandard,
		UUIDRepresentationCSharpLegacy,
		UUIDRepresentationJavaLegacy,
		UUIDRepresentationPythonLegacy:
		return nil
	default:
		return fmt.Errorf("invalid mongodb uuid representation: %s", c.UUIDRepresentation)
	}
}

// validateSettings verifies port, protocol, and UUID representation settings.
func (c *Config) validateSettings() error {
	if err := c.validatePort(); err != nil {
		return err
	}
	if err := c.validateProtocol(); err != nil {
		return err
	}
	return c.validateUUIDRepresentation()
}

// Validate checks that the configuration values are valid.
func (c *Config) Validate() error {
	if c == nil {
		return ErrNilConfig
	}
	if err := c.validateRequiredFields(); err != nil {
		return err
	}
	return c.validateSettings()
}
