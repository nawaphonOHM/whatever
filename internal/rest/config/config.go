// Package config provides environment variable parsing and configuration.
package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// isNotExistError checks if error represents a file non-existence error.
func isNotExistError(err error) bool {
	if errors.Is(err, os.ErrNotExist) {
		return true
	}
	return os.IsNotExist(err)
}

// loadExplicitEnv loads the specified environment files if provided.
func loadExplicitEnv(filenames []string) error {
	err := godotenv.Load(filenames...)
	if err != nil && !isNotExistError(err) {
		return fmt.Errorf("failed to load env file: %w", err)
	}
	return nil
}

// tryLoadFile attempts to load a file if it exists on disk.
func tryLoadFile(path string) (bool, error) {
	if _, err := os.Stat(path); err == nil {
		if err := godotenv.Load(path); err != nil {
			return true, fmt.Errorf("failed to load %s file: %w", path, err)
		}
		return true, nil
	}
	return false, nil
}

// loadDefaultEnv tries loading default .env or configs/.env if present.
func loadDefaultEnv() error {
	for _, p := range []string{".env", "configs/.env"} {
		loaded, err := tryLoadFile(p)
		if loaded {
			return err
		}
	}
	return nil
}

// loadEnvFiles loads environment variables from files based on args.
func loadEnvFiles(filenames []string) error {
	if len(filenames) > 0 {
		return loadExplicitEnv(filenames)
	}
	return loadDefaultEnv()
}

type defaulter interface {
	SetDefaults()
}

// parseTarget applies defaults if available and parses environment variables.
func parseTarget[T any](target *T) error {
	if d, ok := any(target).(defaulter); ok {
		d.SetDefaults()
	}
	return env.Parse(target)
}

// Load parses environment variables into a struct of type T.
// It searches for .env files in filenames, or default locations if none.
// Values precedence: System Environment > .env file > Struct Default Tags.
func Load[T any](filenames ...string) (*T, error) {
	if err := loadEnvFiles(filenames); err != nil {
		return nil, err
	}

	var target T
	if err := parseTarget(&target); err != nil {
		return nil, fmt.Errorf("failed to parse environment config: %w", err)
	}

	return &target, nil
}
