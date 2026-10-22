// Package rest provides standardized API registration contracts and
// JSON response envelopes according to RFC 9457 Problem Details.
package rest

import (
	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
)

// APIGroup represents the API registration group mode.
type APIGroup = contracts.APIGroup

const (
	// APIGroupStandard prefixes routes with /api and version segments (/api/v[N]).
	APIGroupStandard = contracts.APIGroupStandard
	// APIGroupRaw registers routes directly with prefix and path without /api or version segments.
	APIGroupRaw = contracts.APIGroupRaw
)

// ExportableAPI defines a single API route endpoint.
type ExportableAPI = contracts.ExportableAPI

// RRestAPIRegistration groups APIs under a common prefix and version.
type RRestAPIRegistration = contracts.RRestAPIRegistration
