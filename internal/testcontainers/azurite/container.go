// Package azurite provides a reusable Testcontainers wrapper for Azurite.
package azurite

import (
	"context"
	"sync"

	tcazurite "github.com/testcontainers/testcontainers-go/modules/azure/azurite"
)

// Container wraps a running Azurite Testcontainer and provides service endpoints.
type Container struct {
	rawContainer *tcazurite.Container
	mu           sync.RWMutex
	terminated   bool
}

// NewContainer creates a Container wrapping an Azurite Testcontainer.
func NewContainer(raw *tcazurite.Container) *Container {
	return &Container{rawContainer: raw}
}

func (c *Container) checkRunningLocked() error {
	if c.terminated || !c.rawContainerIsRunning() {
		return ErrContainerUnavailable
	}
	return nil
}

func (c *Container) rawContainerIsRunning() bool {
	if c.rawContainer == nil {
		return false
	}
	if c.rawContainer.Container == nil {
		return false
	}
	return c.rawContainer.IsRunning()
}

func (c *Container) serviceURL(
	ctx context.Context,
	endpoint func(*tcazurite.Container, context.Context) (string, error),
) (string, error) {
	if c == nil {
		return "", ErrNilContainer
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	if err := c.checkRunningLocked(); err != nil {
		return "", err
	}
	return endpoint(c.rawContainer, ctx)
}

// BlobServiceURL returns the mapped Blob service endpoint.
func (c *Container) BlobServiceURL(ctx context.Context) (string, error) {
	return c.serviceURL(ctx, func(raw *tcazurite.Container, ctx context.Context) (string, error) {
		return raw.BlobServiceURL(ctx)
	})
}

// QueueServiceURL returns the mapped Queue service endpoint.
func (c *Container) QueueServiceURL(ctx context.Context) (string, error) {
	return c.serviceURL(ctx, func(raw *tcazurite.Container, ctx context.Context) (string, error) {
		return raw.QueueServiceURL(ctx)
	})
}

// TableServiceURL returns the mapped Table service endpoint.
func (c *Container) TableServiceURL(ctx context.Context) (string, error) {
	return c.serviceURL(ctx, func(raw *tcazurite.Container, ctx context.Context) (string, error) {
		return raw.TableServiceURL(ctx)
	})
}

// RawContainer returns the underlying Azurite Testcontainer.
func (c *Container) RawContainer() *tcazurite.Container {
	if c == nil {
		return nil
	}
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.rawContainer
}
