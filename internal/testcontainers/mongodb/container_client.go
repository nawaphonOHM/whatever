package mongodb

import (
	"context"

	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/client"
	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/config"
)

// resolveHostPort retrieves both the mapped host and port from the container.
func (c *Container) resolveHostPort(ctx context.Context) (string, int, error) {
	host, err := c.Host(ctx)
	if err != nil {
		return "", 0, err
	}
	port, err := c.Port(ctx)
	return host, port, err
}

// populateConfig constructs a Config struct from container connection settings.
func (c *Container) populateConfig(host string, port int) *config.Config {
	cfg := config.DefaultConfig()
	cfg.Host = host
	cfg.Port = port
	cfg.Database = c.defaultDatabase
	cfg.Username = c.username
	cfg.Password = c.password
	if c.username != "" {
		cfg.AuthSource = "admin"
	}
	return cfg
}

// Config returns a populated MongoDB config struct initialized with container connection details.
func (c *Container) Config(ctx context.Context) (*config.Config, error) {
	if err := c.checkRunning(); err != nil {
		return nil, err
	}
	host, port, err := c.resolveHostPort(ctx)
	if err != nil {
		return nil, err
	}
	return c.populateConfig(host, port), nil
}

// Client connects to the running MongoDB container, verifies ping connectivity, and returns a Client.
func (c *Container) Client(ctx context.Context, opts ...client.Option) (*client.Client, error) {
	cfg, err := c.Config(ctx)
	if err != nil {
		return nil, err
	}
	return client.ConnectWithConfig(ctx, cfg, opts...)
}
