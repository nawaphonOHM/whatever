package azurite

import "errors"

var (
	// ErrNilContainer is returned when an operation is performed on a nil Container.
	ErrNilContainer = errors.New("azurite container is nil")
	// ErrContainerUnavailable is returned when the container is not running or has been terminated.
	ErrContainerUnavailable = errors.New("azurite container is unavailable")
)
