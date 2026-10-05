package firestore

import (
	tcfirestore "github.com/testcontainers/testcontainers-go/modules/gcloud/firestore"
)

// Container wraps a running Firestore Testcontainer and provides connection utilities.
type Container struct {
	rawContainer  *tcfirestore.Container
	projectID     string
	datastoreMode bool
}

// NewContainer creates a new Container wrapping a Firestore testcontainer with options.
func NewContainer(raw *tcfirestore.Container, opts *Options) *Container {
	resolved := opts
	if resolved == nil {
		resolved = DefaultOptions()
	}
	return &Container{
		rawContainer:  raw,
		projectID:     resolved.ProjectID,
		datastoreMode: resolved.DatastoreMode,
	}
}

// checkRunning verifies that the receiver and underlying container are not nil.
func (c *Container) checkRunning() error {
	if c == nil {
		return ErrNilContainer
	}
	if c.rawContainer == nil {
		return ErrContainerNotRunning
	}
	return nil
}

// ProjectID returns the configured project ID.
func (c *Container) ProjectID() string {
	if c == nil {
		return ""
	}
	if c.rawContainer != nil {
		return c.rawContainer.ProjectID()
	}
	return c.projectID
}

// DatastoreMode returns whether the container is running in datastore mode.
func (c *Container) DatastoreMode() bool {
	if c == nil {
		return false
	}
	return c.datastoreMode
}

// RawContainer returns the underlying Firestore container instance.
func (c *Container) RawContainer() *tcfirestore.Container {
	if c == nil {
		return nil
	}
	return c.rawContainer
}
