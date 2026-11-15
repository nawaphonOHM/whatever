// Package azurite provides public Azurite Testcontainers utilities, type aliases, and functional options.
package azurite

import (
	internalazurite "github.com/nawaphonOHM/whatever/v3/internal/testcontainers/azurite"
	tcazurite "github.com/testcontainers/testcontainers-go/modules/azure/azurite"
)

// Container wraps a running Azurite Testcontainer and provides service endpoints.
type Container = internalazurite.Container

// Option defines a functional option for configuring an Azurite container.
type Option = internalazurite.Option

// Options holds configuration settings for running an Azurite container.
type Options = internalazurite.Options

// Default settings and credentials intended only for local emulator use.
const (
	// DefaultImage is the pinned default Azurite emulator image.
	DefaultImage = internalazurite.DefaultImage
	// DefaultAccountName is the default Azurite testing account name.
	DefaultAccountName = tcazurite.AccountName
	// DefaultAccountKey is the default Azurite testing account key.
	DefaultAccountKey = tcazurite.AccountKey
)

// Sentinel errors for Azurite Testcontainers operations.
var (
	// ErrNilContainer is returned when an operation is performed on a nil Container.
	ErrNilContainer = internalazurite.ErrNilContainer
	// ErrContainerUnavailable is returned when the container is not running or has been terminated.
	ErrContainerUnavailable = internalazurite.ErrContainerUnavailable
)
