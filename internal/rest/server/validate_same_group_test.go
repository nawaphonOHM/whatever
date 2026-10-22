package server

import (
	"testing"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/contracts"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type sameGroupCase struct {
	name          string
	expectedPaths []string
	regs          []*contracts.RRestAPIRegistration
}

func sameGroupCases() []*sameGroupCase {
	std, raw := contracts.APIGroupStandard, contracts.APIGroupRaw
	return []*sameGroupCase{
		{
			name: "two standard groups",
			regs: []*contracts.RRestAPIRegistration{
				singleAPI("/v1/users", "/profile", std),
				singleAPI("/v1/orders", "/list", std),
			},
			expectedPaths: []string{"/api/v1/users/profile", "/api/v1/orders/list"},
		},
		{
			name: "default empty and explicit standard",
			regs: []*contracts.RRestAPIRegistration{
				singleAPI("/items", "/all", ""),
				singleAPI("/orders", "/all", std),
			},
			expectedPaths: []string{"/api/items/all", "/api/orders/all"},
		},
		{
			name: "two raw groups",
			regs: []*contracts.RRestAPIRegistration{
				singleAPI("/custom/hooks", "/event", raw),
				singleAPI("/metrics", "/collector", raw),
			},
			expectedPaths: []string{"/custom/hooks/event", "/metrics/collector"},
		},
	}
}

// TestValidateRegistrations_SameGroupPasses ensures homogeneous groups pass validation.
func TestValidateRegistrations_SameGroupPasses(t *testing.T) {
	for _, tt := range sameGroupCases() {
		t.Run(tt.name, func(t *testing.T) {
			routes, err := validateRegistrations(tt.regs)
			require.NoError(t, err)
			require.Len(t, routes, len(tt.expectedPaths))
			for i, exp := range tt.expectedPaths {
				assert.Equal(t, exp, routes[i].Path)
			}
		})
	}
}
