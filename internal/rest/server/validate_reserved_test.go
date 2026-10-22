package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const testWildcardAny = "/*any"

func testReservedPaths(t *testing.T, grp contracts.APIGroup, reserved string, method contracts.HTTPMethod) {
	t.Run(string(grp)+"/path/"+reserved+"/"+method.GetName(), func(t *testing.T) {
		regs := []*contracts.RRestAPIRegistration{{
			Group: grp,
			Apis: []*contracts.ExportableAPI{{
				Path:    contracts.Pathz(reserved),
				Method:  method,
				Handler: dummyHandler,
			}},
		}}

		_, err := validateRegistrations(regs)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrReservedPath)
		assert.Contains(t, err.Error(), `is reserved for framework liveness/readiness probes`)
	})

	t.Run(string(grp)+"/prefix/"+reserved+"/"+method.GetName(), func(t *testing.T) {
		regs := []*contracts.RRestAPIRegistration{{
			Group:  grp,
			Prefix: contracts.Pathz(reserved),
			Apis: []*contracts.ExportableAPI{{
				Path:    "",
				Method:  method,
				Handler: dummyHandler,
			}},
		}}

		_, err := validateRegistrations(regs)
		require.Error(t, err)
		assert.ErrorIs(t, err, ErrReservedPath)
		assert.Contains(t, err.Error(), `is reserved for framework liveness/readiness probes`)
	})
}

type reservedCombo struct {
	grp      contracts.APIGroup
	reserved string
	method   contracts.HTTPMethod
}

func reservedCombos() []reservedCombo {
	return []reservedCombo{
		{contracts.APIGroupStandard, ReservedHealthPath, contracts.GET},
		{contracts.APIGroupStandard, ReservedHealthPath, contracts.POST},
		{contracts.APIGroupStandard, ReservedReadyPath, contracts.GET},
		{contracts.APIGroupStandard, ReservedReadyPath, contracts.POST},
		{contracts.APIGroupRaw, ReservedHealthPath, contracts.GET},
		{contracts.APIGroupRaw, ReservedHealthPath, contracts.POST},
		{contracts.APIGroupRaw, ReservedReadyPath, contracts.GET},
		{contracts.APIGroupRaw, ReservedReadyPath, contracts.POST},
		{"", ReservedHealthPath, contracts.GET},
		{"", ReservedHealthPath, contracts.POST},
		{"", ReservedReadyPath, contracts.GET},
		{"", ReservedReadyPath, contracts.POST},
	}
}

// TestValidateRegistrations_ReservedPath rejects /health and /ready across all groups and methods.
func TestValidateRegistrations_ReservedPath(t *testing.T) {
	for _, c := range reservedCombos() {
		testReservedPaths(t, c.grp, c.reserved, c.method)
	}
}
