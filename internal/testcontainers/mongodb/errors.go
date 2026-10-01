package mongodb

import "errors"

var (
	// ErrNilContainer is returned when an operation is performed on a nil Container.
	ErrNilContainer = errors.New("mongodb container is nil")
	// ErrContainerNotRunning is returned when an operation is performed on a non-running container.
	ErrContainerNotRunning = errors.New("mongodb container is not running")
)
