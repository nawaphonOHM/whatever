package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// allowedAdjacentCase describes a non-conflicting path test scenario.
type allowedAdjacentCase struct {
	name     string
	group    contracts.APIGroup
	prefix   contracts.Pathz
	path     contracts.Pathz
	expected string
}

func allowedAdjacentCases() []*allowedAdjacentCase {
	return []*allowedAdjacentCase{
		{
			name:     "raw root param",
			group:    contracts.APIGroupRaw,
			path:     "/:name",
			expected: "/:name",
		},
		{
			name:     "raw subpath catchall",
			group:    contracts.APIGroupRaw,
			prefix:   "/users",
			path:     testWildcardAny,
			expected: "/users/*any",
		},
		{
			name:     "raw probe-like prefix healthy",
			group:    contracts.APIGroupRaw,
			path:     "/healthy",
			expected: "/healthy",
		},
		{
			name:     "raw probe-like prefix readiness",
			group:    contracts.APIGroupRaw,
			path:     "/readiness",
			expected: "/readiness",
		},
		{
			name:     "standard root catchall gets api prefixed",
			group:    contracts.APIGroupStandard,
			path:     testWildcardAny,
			expected: "/api/v1/*any",
		},
	}
}

// TestValidateRegistrations_AllowedProbeAdjacentPaths confirms valid non-conflicting paths are accepted.
func TestValidateRegistrations_AllowedProbeAdjacentPaths(t *testing.T) {
	for _, tt := range allowedAdjacentCases() {
		t.Run(tt.name, func(t *testing.T) {
			regs := []*contracts.RRestAPIRegistration{{
				Group:   tt.group,
				Prefix:  tt.prefix,
				Version: 1,
				Apis: []*contracts.ExportableAPI{{
					Path:    tt.path,
					Method:  contracts.GET,
					Handler: dummyHandler,
				}},
			}}

			routes, err := validateRegistrations(regs)
			require.NoError(t, err)
			require.Len(t, routes, 1)
			assert.Equal(t, tt.expected, routes[0].Path)
		})
	}
}
