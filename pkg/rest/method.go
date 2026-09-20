package rest

import (
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
)

// HTTPMethod represents supported HTTP request methods.
type HTTPMethod = contracts.HTTPMethod

// Supported HTTP method constants.
const (
	GET     = contracts.GET
	HEAD    = contracts.HEAD
	POST    = contracts.POST
	PUT     = contracts.PUT
	PATCH   = contracts.PATCH
	DELETE  = contracts.DELETE
	OPTIONS = contracts.OPTIONS
	CONNECT = contracts.CONNECT
	TRACE   = contracts.TRACE
)
