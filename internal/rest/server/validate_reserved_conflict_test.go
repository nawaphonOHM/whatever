package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type probeConflictCase struct {
	name      string
	prefix    contracts.Pathz
	path      contracts.Pathz
	expectMsg string
}

func probeConflictCases() []*probeConflictCase {
	return []*probeConflictCase{
		{
			name:      "root catch-all any",
			path:      testWildcardAny,
			expectMsg: `root catch-all "/*any" conflicts with framework probe routes`,
		},
		{
			name:      "root catch-all wildcard prefix",
			prefix:    "/*filepath",
			expectMsg: `root catch-all "/*filepath" conflicts with framework probe routes`,
		},
		{
			name:      "health wildcard subpath",
			prefix:    "/health",
			path:      "/*x",
			expectMsg: `path "/health/*x" is reserved for framework liveness/readiness probes`,
		},
		{
			name:      "ready wildcard subpath",
			prefix:    "/ready",
			path:      testWildcardAny,
			expectMsg: `path "/ready/*any" is reserved for framework liveness/readiness probes`,
		},
		{
			name:      "health static subpath",
			path:      "/health/check",
			expectMsg: `path "/health/check" is reserved for framework liveness/readiness probes`,
		},
		{
			name:      "ready static subpath",
			path:      "/ready/status",
			expectMsg: `path "/ready/status" is reserved for framework liveness/readiness probes`,
		},
	}
}

// TestValidateRegistrations_RawCatchAllAndProbeSubpaths rejects router paths that conflict with probes.
func TestValidateRegistrations_RawCatchAllAndProbeSubpaths(t *testing.T) {
	for _, tt := range probeConflictCases() {
		t.Run("raw/"+tt.name, func(t *testing.T) {
			regs := []*contracts.RRestAPIRegistration{{
				Group:  contracts.APIGroupRaw,
				Prefix: tt.prefix,
				Apis: []*contracts.ExportableAPI{{
					Path:    tt.path,
					Method:  contracts.GET,
					Handler: dummyHandler,
				}},
			}}

			_, err := validateRegistrations(regs)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrReservedPath)
			assert.Contains(t, err.Error(), tt.expectMsg)
		})
	}
}
