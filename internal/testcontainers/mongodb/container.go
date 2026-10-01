package mongodb

import (
	tcmongodb "github.com/testcontainers/testcontainers-go/modules/mongodb"
)

// Container wraps a running MongoDB Testcontainer and provides connection utilities.
type Container struct {
	rawContainer    *tcmongodb.MongoDBContainer
	defaultDatabase string
	username        string
	password        string
	replicaSet      string
}

// NewContainer creates a new Container wrapping a MongoDB testcontainer with options.
func NewContainer(raw *tcmongodb.MongoDBContainer, opts *Options) *Container {
	resolved := opts
	if resolved == nil {
		resolved = DefaultOptions()
	}
	return &Container{
		rawContainer:    raw,
		defaultDatabase: resolved.Database,
		username:        resolved.Username,
		password:        resolved.Password,
		replicaSet:      resolved.ReplicaSet,
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

// Database returns the configured default database name.
func (c *Container) Database() string {
	if c == nil {
		return ""
	}
	return c.defaultDatabase
}

// Username returns the configured root username.
func (c *Container) Username() string {
	if c == nil {
		return ""
	}
	return c.username
}

// Password returns the configured root password.
func (c *Container) Password() string {
	if c == nil {
		return ""
	}
	return c.password
}

// ReplicaSet returns the configured replica set name.
func (c *Container) ReplicaSet() string {
	if c == nil {
		return ""
	}
	return c.replicaSet
}

// RawContainer returns the underlying MongoDBContainer instance.
func (c *Container) RawContainer() *tcmongodb.MongoDBContainer {
	if c == nil {
		return nil
	}
	return c.rawContainer
}
