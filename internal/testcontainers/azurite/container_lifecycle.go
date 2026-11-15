package azurite

import (
	"context"
	"fmt"
)

// Terminate stops and removes the running Azurite container.
func (c *Container) Terminate(ctx context.Context) error {
	if c == nil {
		return ErrNilContainer
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if err := c.checkRunningLocked(); err != nil {
		return err
	}
	return c.terminateLocked(ctx)
}

func (c *Container) terminateLocked(ctx context.Context) error {
	if err := c.rawContainer.Terminate(ctx); err != nil {
		return fmt.Errorf("failed to terminate azurite container: %w", err)
	}
	c.terminated = true
	return nil
}
