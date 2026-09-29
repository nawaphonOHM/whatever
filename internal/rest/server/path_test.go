package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
)

const (
	testAPIItemsList = "/api/items/list"
	testAPIV1Users   = "/api/v1/users"
	testListPath     = "/list"
)

// pathCase is one CalculateFullPath table row.
type pathCase struct {
	name     string
	prefix   contracts.Pathz
	path     contracts.Pathz
	expected string
	version  contracts.APIVersioning
}

// pathCases returns CalculateFullPath coverage rows.
func pathCases() []pathCase {
	return []pathCase{
		{name: "empty all", expected: "/api"},
		{name: "empty all version 1", version: 1, expected: "/api/v1"},
		{name: "unversioned simple", prefix: "/items", path: testListPath, expected: testAPIItemsList},
		{name: "unversioned missing slashes", prefix: "items", path: "list", expected: testAPIItemsList},
		{name: "unversioned trailing slash", prefix: "/items/", path: "/list/", expected: testAPIItemsList},
		{name: "unversioned multiple slashes", prefix: "///items///", path: "///list///", expected: testAPIItemsList},
		{name: "unversioned with api in prefix", prefix: "/api/items", path: testListPath, expected: testAPIItemsList},
		{name: "unversioned prefix api only with path", prefix: "/api", path: "/items", expected: "/api/items"},
		{name: "versioned prefix", version: 1, prefix: "/users", path: "/:id", expected: "/api/v1/users/:id"},
		{
			name:    "explicit api in prefix is deduplicated",
			version: 1, prefix: "/api/users", path: "/:id", expected: "/api/v1/users/:id",
		},
		{name: "versioned prefix already containing /v1", version: 1, prefix: "/v1/users", expected: testAPIV1Users},
		{
			name:    "versioned prefix already containing /api/v1",
			version: 1, prefix: "/api/v1/users", expected: testAPIV1Users,
		},
		{
			name:    "versioned prefix already containing /v1/api",
			version: 1, prefix: "/v1/api/users", expected: testAPIV1Users,
		},
		{
			name:    "versioned prefix already containing /v2",
			version: 2, prefix: "/v2/items", path: testListPath, expected: "/api/v2/items/list",
		},
		{
			name:    "versioned prefix already containing /api/v2",
			version: 2, prefix: "/api/v2/items", path: testListPath, expected: "/api/v2/items/list",
		},
		{name: "versioned empty prefix with path", version: 1, path: "/items", expected: "/api/v1/items"},
		{name: "versioned root ping path", version: 1, path: "/ping", expected: "/api/v1/ping"},
		{
			name:    "versioned prefix api with ping path",
			version: 1, prefix: "/api", path: "/ping", expected: "/api/v1/ping",
		},
		{
			name:    "versioned pre-prefixed orders",
			version: 1, prefix: "/api/orders", path: "/status", expected: "/api/v1/orders/status",
		},
		{name: "reserved health exact match", path: ReservedHealthPath, expected: ReservedHealthPath},
		{name: "reserved health in prefix", prefix: ReservedHealthPath, expected: ReservedHealthPath},
		{name: "reserved ready exact match", path: ReservedReadyPath, expected: ReservedReadyPath},
		{name: "reserved ready in prefix", prefix: ReservedReadyPath, expected: ReservedReadyPath},
	}
}

// TestCalculateFullPath covers path normalization cases.
func TestCalculateFullPath(t *testing.T) {
	for _, tt := range pathCases() {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateFullPath(tt.version, tt.prefix, tt.path)
			assert.Equal(t, tt.expected, got)
		})
	}
}
