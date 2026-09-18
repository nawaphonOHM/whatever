package server

import (
	server2 "github.com/nawaphonOHM/whatever/internal/rest/server"
)

// Sentinel validation errors.
var (
	ErrReservedPath    = server2.ErrReservedPath
	ErrDuplicateRoute  = server2.ErrDuplicateRoute
	ErrNilRegistration = server2.ErrNilRegistration
	ErrNilAPI          = server2.ErrNilAPI
	ErrNilHandler      = server2.ErrNilHandler
	ErrInvalidMethod   = server2.ErrInvalidMethod
)
