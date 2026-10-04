// Package mongodb provides public MongoDB Testcontainers utilities, type
// aliases, and functional options for integration testing.
package mongodb

import (
	"github.com/nawaphonOHM/whatever/v2/internal/testcontainers/mongodb"
)

// Container wraps a running MongoDB Testcontainer and provides connection utilities.
type Container = mongodb.Container

// Option defines a functional option for configuring a MongoDB container.
type Option = mongodb.Option

// Options holds configuration settings for running a MongoDB container.
type Options = mongodb.Options

// Constants for MongoDB Testcontainers default settings.
const (
	// DefaultImage is the default MongoDB container image.
	DefaultImage = mongodb.DefaultImage
	// DefaultPort is the default MongoDB container port.
	DefaultPort = mongodb.DefaultPort
)

// Sentinel errors for MongoDB Testcontainers operations.
var (
	// ErrNilContainer is returned when an operation is performed on a nil Container.
	ErrNilContainer = mongodb.ErrNilContainer
	// ErrContainerNotRunning is returned when an operation is performed on a non-running container.
	ErrContainerNotRunning = mongodb.ErrContainerNotRunning
)
