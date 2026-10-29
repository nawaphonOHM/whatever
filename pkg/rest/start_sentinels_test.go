package rest

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/nawaphonOHM/whatever/v3/internal/rest/server"
)

func TestSentinelErrors(t *testing.T) {
	cases := []struct {
		err      error
		expected error
		message  string
	}{
		{
			err:      ErrMixedAPIGroups,
			expected: server.ErrMixedAPIGroups,
			message:  "cannot mix standard and raw API registration groups in blueprint",
		},
		{
			err:      ErrInvalidRegistrationGroup,
			expected: server.ErrInvalidRegistrationGroup,
			message:  "invalid API registration group",
		},
	}
	for _, tc := range cases {
		assert.NotNil(t, tc.err)
		assert.Same(t, tc.expected, tc.err)
		assert.Equal(t, tc.message, tc.err.Error())
	}
}
