package mongodb

import (
	"context"
	"fmt"
)

// ConnectionString returns the connection URI string for the running MongoDB container.
func (c *Container) ConnectionString(ctx context.Context) (string, error) {
	if err := c.checkRunning(); err != nil {
		return "", err
	}
	return c.rawContainer.ConnectionString(ctx)
}

// Host returns the host IP or hostname where the MongoDB container is accessible.
func (c *Container) Host(ctx context.Context) (string, error) {
	if err := c.checkRunning(); err != nil {
		return "", err
	}
	return c.rawContainer.Host(ctx)
}

// Port returns the mapped external host port for the MongoDB container.
func (c *Container) Port(ctx context.Context) (int, error) {
	if err := c.checkRunning(); err != nil {
		return 0, err
	}
	p, err := c.rawContainer.MappedPort(ctx, DefaultPort)
	if err != nil {
		return 0, fmt.Errorf("failed to get mapped port: %w", err)
	}
	return int(p.Num()), nil
}
