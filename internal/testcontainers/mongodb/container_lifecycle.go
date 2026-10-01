package mongodb

import (
	"context"
)

// Terminate stops and removes the running MongoDB container.
func (c *Container) Terminate(ctx context.Context) error {
	if err := c.checkRunning(); err != nil {
		return err
	}
	return c.rawContainer.Terminate(ctx)
}
