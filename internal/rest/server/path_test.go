package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
)

const (
	testAPIItemsList = "/api/items/list"
	testAPIV1Users   = "/api/v1/users"
	testListPath     = "/list"
	testItemsPrefix  = "/items"
	testAPIPrefix    = "/api"
	testRootPath     = "/"
)

// pathCase is one CalculateFullPath table row.
type pathCase struct {
	name     string
	prefix   contracts.Pathz
	path     contracts.Pathz
	expected string
	group    contracts.APIGroup
	version  contracts.APIVersioning
}

// pathCases returns CalculateFullPath coverage rows.
func pathCases() []*pathCase {
	cases := standardPathCases()
	return append(cases, rawPathCases()...)
}

// TestCalculateFullPath covers path normalization cases.
func TestCalculateFullPath(t *testing.T) {
	for _, tt := range pathCases() {
		t.Run(tt.name, func(t *testing.T) {
			got := CalculateFullPath(tt.group, tt.version, tt.prefix, tt.path)
			assert.Equal(t, tt.expected, got)
		})
	}
}
