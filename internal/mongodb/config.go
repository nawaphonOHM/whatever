// Package mongodb provides internal configuration and connection management
// for MongoDB.
package mongodb

import (
	"fmt"
	"time"

	intcfg "github.com/nawaphonOHM/whatever/internal/rest/config"
)

// Default configuration constants.
const (
	DefaultURI             = "mongodb://localhost:27017"
	DefaultConnectTimeout  = 10 * time.Second
	DefaultServerSelection = 5 * time.Second
	DefaultSocketTimeout   = 10 * time.Second
	DefaultMaxPoolSize     = 100
	DefaultMinPoolSize     = 5
	DefaultMaxConnIdleTime = 10 * time.Minute
)

// Config defines configuration options for connecting to MongoDB.
// All environment variable keys use the OHM9996_MONGODB_ prefix.
type Config struct {
	BaseFields    `envPrefix:"OHM9996_MONGODB_"`
	TimeoutFields `envPrefix:"OHM9996_MONGODB_"`
	PoolFields    `envPrefix:"OHM9996_MONGODB_"`
}

// DefaultConfig returns MongoDB configuration with recommended production
// defaults.
func DefaultConfig() *Config {
	return &Config{
		BaseFields: BaseFields{
			URI:      DefaultURI,
			Database: "",
			AppName:  "",
		},
		TimeoutFields: TimeoutFields{
			ConnectTimeout:         DefaultConnectTimeout,
			ServerSelectionTimeout: DefaultServerSelection,
			SocketTimeout:          DefaultSocketTimeout,
			MaxConnIdleTime:        DefaultMaxConnIdleTime,
		},
		PoolFields: PoolFields{
			MaxPoolSize: DefaultMaxPoolSize,
			MinPoolSize: DefaultMinPoolSize,
		},
	}
}

// LoadConfig loads MongoDB configuration from environment variables with
// OHM9996_MONGODB_* prefix and validates the resulting settings.
func LoadConfig(filenames ...string) (*Config, error) {
	cfg, err := intcfg.Load[Config](filenames...)
	if err != nil {
		return nil, fmt.Errorf("failed to load mongodb config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid mongodb config: %w", err)
	}
	return cfg, nil
}
