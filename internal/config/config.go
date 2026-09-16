package config

import (
	"fmt"
	"os"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Load parses environment variables into a struct of type T.
// It searches for .env files in the provided filenames, or ".env" by default if none provided.
// If .env file is missing, it silently proceeds without error and reads from system environment.
// Values precedence: System Environment > .env file > Struct Default Tags.
func Load[T any](filenames ...string) (*T, error) {
	var target T

	// Attempt to load .env file if available
	if len(filenames) > 0 {
		_ = godotenv.Load(filenames...)
	} else {
		// Try default .env file in current directory or configs/.env
		if _, err := os.Stat(".env"); err == nil {
			_ = godotenv.Load(".env")
		} else if _, err := os.Stat("configs/.env"); err == nil {
			_ = godotenv.Load("configs/.env")
		}
	}

	// Parse environment variables with struct tags
	if err := env.Parse(&target); err != nil {
		return nil, fmt.Errorf("failed to parse environment config: %w", err)
	}

	return &target, nil
}
