// Package firestore provides public Google Cloud Firestore Testcontainers utilities,
// type aliases, and functional options for integration testing.
package firestore

import (
	"github.com/nawaphonOHM/whatever/v2/internal/testcontainers/firestore"
)

// Container wraps a running Firestore Testcontainer and provides connection utilities.
type Container = firestore.Container

// Option defines a functional option for configuring a Firestore container.
type Option = firestore.Option

// Options holds configuration settings for running a Firestore container.
type Options = firestore.Options

// Constants for Firestore Testcontainers default settings.
const (
	// DefaultImage is the default Google Cloud SDK image used for Firestore emulator.
	DefaultImage = firestore.DefaultImage
	// DefaultPort is the default Firestore emulator container port.
	DefaultPort = firestore.DefaultPort
	// DefaultProjectID is the default Google Cloud Project ID for the Firestore container.
	DefaultProjectID = firestore.DefaultProjectID
)

// Sentinel errors for Firestore Testcontainers operations.
var (
	// ErrNilContainer is returned when an operation is performed on a nil Container.
	ErrNilContainer = firestore.ErrNilContainer
	// ErrContainerNotRunning is returned when an operation is performed on a non-running container.
	ErrContainerNotRunning = firestore.ErrContainerNotRunning
)
