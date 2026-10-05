package firestore

import (
	"context"
	"fmt"
)

// URI returns the URI string (host:port) for the running Firestore container.
func (c *Container) URI(ctx context.Context) (string, error) {
	if err := c.checkRunning(); err != nil {
		return "", err
	}
	if uri := c.rawContainer.URI(); uri != "" {
		return uri, nil
	}
	return c.rawContainer.PortEndpoint(ctx, DefaultPort, "")
}

// Host returns the host IP or hostname where the Firestore container is accessible.
func (c *Container) Host(ctx context.Context) (string, error) {
	if err := c.checkRunning(); err != nil {
		return "", err
	}
	return c.rawContainer.Host(ctx)
}

// Port returns the mapped external host port for the Firestore container.
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
