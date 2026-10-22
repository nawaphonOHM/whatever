package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type allowedNonGetCase struct {
	name     string
	path     contracts.Pathz
	expected string
	method   contracts.HTTPMethod
}

func allowedNonGetCases() []*allowedNonGetCase {
	return []*allowedNonGetCase{
		{
			name:     "raw post root catchall allowed",
			path:     testWildcardAny,
			method:   contracts.POST,
			expected: "/*any",
		},
		{
			name:     "raw put root catchall allowed",
			path:     testWildcardAny,
			method:   contracts.PUT,
			expected: "/*any",
		},
		{
			name:     "raw delete root catchall allowed",
			path:     testWildcardAny,
			method:   contracts.DELETE,
			expected: "/*any",
		},
		{
			name:     "raw post health subpath allowed",
			path:     "/health/check",
			method:   contracts.POST,
			expected: "/health/check",
		},
		{
			name:     "raw post ready subpath allowed",
			path:     "/ready/webhook",
			method:   contracts.POST,
			expected: "/ready/webhook",
		},
	}
}

// TestValidateRegistrations_AllowedNonGetCatchAll verifies non-GET root catch-alls and subpaths are accepted.
func TestValidateRegistrations_AllowedNonGetCatchAll(t *testing.T) {
	for _, tt := range allowedNonGetCases() {
		t.Run(tt.name, func(t *testing.T) {
			regs := []*contracts.RRestAPIRegistration{{
				Group: contracts.APIGroupRaw,
				Apis: []*contracts.ExportableAPI{{
					Path:    tt.path,
					Method:  tt.method,
					Handler: dummyHandler,
				}},
			}}

			routes, err := validateRegistrations(regs)
			require.NoError(t, err)
			require.Len(t, routes, 1)
			assert.Equal(t, tt.expected, routes[0].Path)
			assert.Equal(t, tt.method.GetName(), routes[0].Method)
		})
	}
}
