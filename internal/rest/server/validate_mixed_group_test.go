package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testMixedPrefix1 = "/api1"
	testMixedPrefix2 = "/api2"
	testMixedPath1   = "/endpoint1"
	testMixedPath2   = "/endpoint2"
)

// mixedGroupCase holds test parameters for mixed group rejection.
type mixedGroupCase struct {
	name       string
	expectIdx0 string
	expectIdx1 string
	regs       []*contracts.RRestAPIRegistration
}

func singleAPI(prefix, path contracts.Pathz, grp contracts.APIGroup) *contracts.RRestAPIRegistration {
	return &contracts.RRestAPIRegistration{
		Group: grp, Prefix: prefix,
		Apis: []*contracts.ExportableAPI{{Path: path, Method: contracts.GET, Handler: dummyHandler}},
	}
}

func mixedGroupCases() []*mixedGroupCase {
	std, raw := contracts.APIGroupStandard, contracts.APIGroupRaw
	return []*mixedGroupCase{
		{
			name: "standard followed by raw",
			regs: []*contracts.RRestAPIRegistration{
				singleAPI(testMixedPrefix1, testMixedPath1, std), singleAPI(testMixedPrefix2, testMixedPath2, raw),
			},
			expectIdx0: "registration at index 0 has group \"standard\"",
			expectIdx1: "registration at index 1 has group \"raw\"",
		},
		{
			name: "default empty followed by raw",
			regs: []*contracts.RRestAPIRegistration{
				singleAPI(testMixedPrefix1, testMixedPath1, ""), singleAPI(testMixedPrefix2, testMixedPath2, raw),
			},
			expectIdx0: "registration at index 0 has group \"standard\"",
			expectIdx1: "registration at index 1 has group \"raw\"",
		},
		{
			name: "raw followed by standard",
			regs: []*contracts.RRestAPIRegistration{
				singleAPI(testMixedPrefix1, testMixedPath1, raw), singleAPI(testMixedPrefix2, testMixedPath2, std),
			},
			expectIdx0: "registration at index 0 has group \"raw\"",
			expectIdx1: "registration at index 1 has group \"standard\"",
		},
		{
			name: "raw followed by default empty",
			regs: []*contracts.RRestAPIRegistration{
				singleAPI(testMixedPrefix1, testMixedPath1, raw), singleAPI(testMixedPrefix2, testMixedPath2, ""),
			},
			expectIdx0: "registration at index 0 has group \"raw\"",
			expectIdx1: "registration at index 1 has group \"standard\"",
		},
		{
			name: "multiple standard followed by raw at index 2",
			regs: []*contracts.RRestAPIRegistration{
				singleAPI("/api0", "/endpoint0", std),
				singleAPI(testMixedPrefix1, testMixedPath1, ""),
				singleAPI(testMixedPrefix2, testMixedPath2, raw),
			},
			expectIdx0: "registration at index 0 has group \"standard\"",
			expectIdx1: "registration at index 2 has group \"raw\"",
		},
	}
}

// TestValidateRegistrations_MixedAPIGroups rejects blueprints mixing standard and raw groups.
func TestValidateRegistrations_MixedAPIGroups(t *testing.T) {
	for _, tt := range mixedGroupCases() {
		t.Run(tt.name, func(t *testing.T) {
			_, err := validateRegistrations(tt.regs)
			require.Error(t, err)
			assert.ErrorIs(t, err, ErrMixedAPIGroups)
			assert.Contains(t, err.Error(), tt.expectIdx0)
			assert.Contains(t, err.Error(), tt.expectIdx1)
		})
	}
}
