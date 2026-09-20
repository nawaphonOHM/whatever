package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
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
		{name: "empty all", expected: "/"},
		{
			name:   "unversioned simple",
			prefix: "/items", path: "/list",
			expected: "/items/list",
		},
		{
			name:   "unversioned missing slashes",
			prefix: "items", path: "list",
			expected: "/items/list",
		},
		{
			name:   "unversioned trailing slash",
			prefix: "/items/", path: "/list/",
			expected: "/items/list",
		},
		{
			name: "versioned prefix", version: 1,
			prefix: "/users", path: "/:id",
			expected: "/v1/users/:id",
		},
		{
			name:    "explicit api in prefix is caller-owned",
			version: 1, prefix: "/api/users",
			expected: "/v1/api/users",
		},
		{
			name:    "versioned prefix already containing /v1",
			version: 1, prefix: "/v1/users",
			expected: "/v1/users",
		},
		{
			name:    "versioned prefix already containing /v2",
			version: 2, prefix: "/v2/items", path: "/list",
			expected: "/v2/items/list",
		},
		{
			name:    "versioned empty prefix with path",
			version: 1, path: "/items",
			expected: "/v1/items",
		},
		{
			name: "reserved health exact match",
			path: ReservedHealthPath, expected: ReservedHealthPath,
		},
		{
			name:   "reserved health in prefix",
			prefix: ReservedHealthPath, expected: ReservedHealthPath,
		},
		{
			name: "reserved ready exact match",
			path: ReservedReadyPath, expected: ReservedReadyPath,
		},
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
