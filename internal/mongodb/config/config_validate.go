package config

import (
	"errors"
)

// validateConnectTimeouts verifies connect and server selection durations.
func (c *Config) validateConnectTimeouts() error {
	if c.ConnectTimeout < 0 {
		return errors.New("connect timeout cannot be negative")
	}
	if c.ServerSelectionTimeout < 0 {
		return errors.New("server selection timeout cannot be negative")
	}
	return nil
}

// validateSocketTimeouts verifies socket and idle duration settings.
func (c *Config) validateSocketTimeouts() error {
	if c.SocketTimeout < 0 {
		return errors.New("socket timeout cannot be negative")
	}
	if c.MaxConnIdleTime < 0 {
		return errors.New("max conn idle time cannot be negative")
	}
	return nil
}

// validateTimeouts checks that all duration settings are non-negative.
func (c *Config) validateTimeouts() error {
	if err := c.validateConnectTimeouts(); err != nil {
		return err
	}
	return c.validateSocketTimeouts()
}

// validatePoolSettings verifies connection pool boundary constraints.
func (c *Config) validatePoolSettings() error {
	if c.MaxPoolSize > 0 && c.MinPoolSize > c.MaxPoolSize {
		return errors.New("min pool size cannot be greater than max pool size")
	}
	return nil
}

// validateRequiredBaseFields checks that URI, Database, Username, and Password are present.
func (c *Config) validateRequiredBaseFields() error {
	checks := []struct {
		val string
		err string
	}{
		{val: c.URI, err: "mongodb uri cannot be empty"},
		{val: c.Database, err: "mongodb database cannot be empty"},
		{val: c.Username, err: "mongodb username cannot be empty"},
		{val: c.Password, err: "mongodb password cannot be empty"},
	}
	for _, chk := range checks {
		if chk.val == "" {
			return errors.New(chk.err)
		}
	}
	return nil
}

// validateFields checks URI, database, credentials, and config boundaries.
func (c *Config) validateFields() error {
	if err := c.validateRequiredBaseFields(); err != nil {
		return err
	}
	if err := c.validateTimeouts(); err != nil {
		return err
	}
	return c.validatePoolSettings()
}

// Validate checks that the configuration values are valid.
func (c *Config) Validate() error {
	if c == nil {
		return ErrNilConfig
	}
	return c.validateFields()
}
