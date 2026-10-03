package mongodb

import "errors"

var (
	// ErrNilClient is returned when an operation is performed on a nil TestClient.
	ErrNilClient = errors.New("mongodb test client is nil")
	// ErrNilConfig is returned when an operation receives nil options or configuration.
	ErrNilConfig = errors.New("mongodb config cannot be nil")
	// ErrEmptyURI is returned when attempting to connect with an empty URI string.
	ErrEmptyURI = errors.New("mongodb uri cannot be empty")
	// ErrInvalidPort is returned when the configured port is outside the valid range.
	ErrInvalidPort = errors.New("mongodb port must be between 1 and 65535")
	// ErrNilTestingTB is returned when a testing helper receives a nil testing.TB instance.
	ErrNilTestingTB = errors.New("testing.TB cannot be nil")
)
