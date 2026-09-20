// Package rest provides standardized API registration contracts and
// JSON response envelopes according to RFC 9457 Problem Details.
package rest

import (
	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
)

// ExportableAPI defines a single API route endpoint.
type ExportableAPI = contracts.ExportableAPI

// RRestAPIRegistration groups APIs under a common prefix and version.
type RRestAPIRegistration = contracts.RRestAPIRegistration
