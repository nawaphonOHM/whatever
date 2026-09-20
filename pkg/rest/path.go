// Package rest provides standardized API registration contracts and
// JSON response envelopes according to RFC 9457 Problem Details.
package rest

import (
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
)

// APIVersioning represents the API major version (0 = unversioned).
type APIVersioning = contracts.APIVersioning

// Pathz represents a URL path segment.
type Pathz = contracts.Pathz
