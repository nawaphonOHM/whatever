package firestore

import (
	"context"
)

// Terminate stops and removes the running Firestore container.
func (c *Container) Terminate(ctx context.Context) error {
	if err := c.checkRunning(); err != nil {
		return err
	}
	return c.rawContainer.Terminate(ctx)
}
