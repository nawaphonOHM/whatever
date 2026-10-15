package config

import (
	"errors"
	"fmt"
	"strings"
)

// validateBucketNameChars verifies bucket name does not contain forbidden characters.
func validateBucketNameChars(bucket string) error {
	if strings.ContainsAny(bucket, " \t\r\n/\\") {
		return fmt.Errorf("bucket name contains invalid characters: %q", bucket)
	}
	return nil
}

// validateBucketLength verifies bucket name length is within limits.
func validateBucketLength(bucket string) error {
	if len(bucket) < MinBucketLength || len(bucket) > MaxBucketLength {
		return fmt.Errorf(
			"bucket name length must be between %d and %d characters: %q",
			MinBucketLength,
			MaxBucketLength,
			bucket,
		)
	}
	return nil
}

// validateBucket checks bucket name presence, length, and valid characters.
func (c *Config) validateBucket() error {
	if c.Bucket == "" {
		return errors.New("storage bucket cannot be empty")
	}
	if err := validateBucketLength(c.Bucket); err != nil {
		return err
	}
	return validateBucketNameChars(c.Bucket)
}

// credentialMismatchError generates an error when only one credential field is set.
func (c *Config) credentialMismatchError() error {
	if c.AccessKeyID != "" {
		return errors.New("storage secret access key cannot be empty when access key id is provided")
	}
	return errors.New("storage access key id cannot be empty when secret access key is provided")
}

// validateCredentials checks that AccessKeyID and SecretAccessKey are provided as a consistent pair.
func (c *Config) validateCredentials() error {
	hasKey := c.AccessKeyID != ""
	hasSecret := c.SecretAccessKey != ""
	if hasKey != hasSecret {
		return c.credentialMismatchError()
	}
	return nil
}

// validateIdentity checks bucket and credentials settings.
func (c *Config) validateIdentity() error {
	if err := c.validateBucket(); err != nil {
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
