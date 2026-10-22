package server

import "errors"

const (
	pathSeparator    = "/"
	defaultAPIPrefix = "/api"
)

// Reserved path constants for framework-managed health probes.
const (
	ReservedHealthPath = "/health"
	ReservedReadyPath  = "/ready"
)

// Sentinel validation errors.
var (
	ErrNilBluePrint = errors.New(
		"blueprint cannot be nil",
	)
	ErrReservedPath = errors.New(
		"route conflicts with reserved health/readiness endpoint",
	)
	ErrDuplicateRoute = errors.New(
		"duplicate route registration",
	)
	ErrNilRegistration = errors.New(
		"registration cannot be nil",
	)
	ErrNilAPI = errors.New(
		"exportable API entry cannot be nil",
	)
	ErrNilHandler = errors.New(
		"handler cannot be nil",
	)
	ErrInvalidMethod = errors.New(
		"invalid HTTP method",
	)
	ErrMixedAPIGroups = errors.New(
		"cannot mix standard and raw API registration groups in blueprint",
	)
	ErrAPIGroupMismatch = errors.New(
		"registration API group does not match OHM9996_SERVER_API_GROUP configuration",
	)
	ErrInvalidAPIGroup = errors.New(
		"invalid OHM9996_SERVER_API_GROUP configuration",
	)
	ErrInvalidRegistrationGroup = errors.New(
		"invalid API registration group",
	)
)
