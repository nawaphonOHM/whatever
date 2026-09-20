// Package rest provides standardized API registration contracts and
// JSON response envelopes according to RFC 9457 Problem Details.
package rest

import (
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
)

// Handler defines route handlers returning a standardized Response.
type Handler = contracts.Handler

// Middleware defines the function signature for route middlewares.
type Middleware = contracts.Middleware
