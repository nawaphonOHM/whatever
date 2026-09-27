package config

import (
	"errors"
	"fmt"
)

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

// validateConnectionTimeouts ensures connect and server selection timeouts are non-negative.
func (c *Config) validateConnectionTimeouts() error {
	if c.ConnectTimeout < 0 {
		return errors.New("mongodb connect timeout cannot be negative")
	}
	if c.ServerSelectionTimeout < 0 {
		return errors.New("mongodb server selection timeout cannot be negative")
	}
	return nil
}

// validateOperationTimeouts ensures socket and idle timeouts are non-negative.
func (c *Config) validateOperationTimeouts() error {
	if c.SocketTimeout < 0 {
		return errors.New("mongodb socket timeout cannot be negative")
	}
	if c.MaxConnIdleTime < 0 {
		return errors.New("mongodb max conn idle time cannot be negative")
	}
	return nil
}

// validateTimeouts ensures all configured timeout durations are non-negative.
func (c *Config) validateTimeouts() error {
	if err := c.validateConnectionTimeouts(); err != nil {
		return err
	}
	return c.validateOperationTimeouts()
}

// validatePool ensures connection pool capacity bounds are valid.
func (c *Config) validatePool() error {
	if c.MaxPoolSize == 0 {
		return errors.New("mongodb max pool size must be greater than 0")
	}
	if c.MinPoolSize > c.MaxPoolSize {
		return errors.New("mongodb min pool size cannot exceed max pool size")
	}
	return nil
}

// validateNetworkSettings verifies port, protocol, and UUID representation settings.
func (c *Config) validateNetworkSettings() error {
	if err := c.validatePort(); err != nil {
		return err
	}
	if err := c.validateProtocol(); err != nil {
		return err
	}
	return c.validateUUIDRepresentation()
}

// validateLifecycleSettings verifies timeout and pool limits.
func (c *Config) validateLifecycleSettings() error {
	if err := c.validateTimeouts(); err != nil {
		return err
	}
	return c.validatePool()
}

// validateSettings verifies port, protocol, UUID representation, timeout, and pool settings.
func (c *Config) validateSettings() error {
	if err := c.validateNetworkSettings(); err != nil {
		return err
	}
	return c.validateLifecycleSettings()
}
