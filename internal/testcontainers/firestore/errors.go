package firestore

import "errors"

var (
	// ErrNilContainer is returned when an operation is performed on a nil Container.
	ErrNilContainer = errors.New("firestore container is nil")
	// ErrContainerNotRunning is returned when an operation is performed on a non-running container.
	ErrContainerNotRunning = errors.New("firestore container is not running")
)
