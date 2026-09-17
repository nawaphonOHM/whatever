// Package mongodb provides internal configuration and connection management for MongoDB.
package mongodb

import (
	"errors"
	"fmt"
	"time"

	intcfg "github.com/nawaphonOHM/whatever/internal/rest/config"
)

// Config defines the configuration options for connecting to a MongoDB instance or cluster.
// All environment variable keys use the OHM9969_MONGODB_ prefix.
type Config struct {
	URI                    string        `env:"OHM9969_MONGODB_URI" envDefault:"mongodb://localhost:27017"`
	Database               string        `env:"OHM9969_MONGODB_DATABASE" envDefault:""`
	ConnectTimeout         time.Duration `env:"OHM9969_MONGODB_CONNECT_TIMEOUT" envDefault:"10s"`
	ServerSelectionTimeout time.Duration `env:"OHM9969_MONGODB_SERVER_SELECTION_TIMEOUT" envDefault:"5s"`
	SocketTimeout          time.Duration `env:"OHM9969_MONGODB_SOCKET_TIMEOUT" envDefault:"10s"`
	MaxPoolSize            uint64        `env:"OHM9969_MONGODB_MAX_POOL_SIZE" envDefault:"100"`
	MinPoolSize            uint64        `env:"OHM9969_MONGODB_MIN_POOL_SIZE" envDefault:"5"`
	MaxConnIdleTime        time.Duration `env:"OHM9969_MONGODB_MAX_CONN_IDLE_TIME" envDefault:"10m"`
	AppName                string        `env:"OHM9969_MONGODB_APP_NAME" envDefault:""`
}

// DefaultConfig returns MongoDB configuration with recommended production defaults.
func DefaultConfig() *Config {
	return &Config{
		URI:                    "mongodb://localhost:27017",
		Database:               "",
		ConnectTimeout:         10 * time.Second,
		ServerSelectionTimeout: 5 * time.Second,
		SocketTimeout:          10 * time.Second,
		MaxPoolSize:            100,
		MinPoolSize:            5,
		MaxConnIdleTime:        10 * time.Minute,
		AppName:                "",
	}
}

// Validate checks that the configuration values are valid.
func (c *Config) Validate() error {
	if c == nil {
		return errors.New("mongodb config cannot be nil")
	}
	if c.URI == "" {
		return errors.New("mongodb uri cannot be empty")
	}
	if c.ConnectTimeout < 0 {
		return errors.New("connect timeout cannot be negative")
	}
	if c.ServerSelectionTimeout < 0 {
		return errors.New("server selection timeout cannot be negative")
	}
	if c.SocketTimeout < 0 {
		return errors.New("socket timeout cannot be negative")
	}
	if c.MaxConnIdleTime < 0 {
		return errors.New("max conn idle time cannot be negative")
	}
	if c.MaxPoolSize > 0 && c.MinPoolSize > c.MaxPoolSize {
		return errors.New("min pool size cannot be greater than max pool size")
	}
	return nil
}

// LoadConfig loads MongoDB configuration from environment variables with OHM9969_MONGODB_* prefix
// and validates the resulting settings.
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
