package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHTTPMethod(t *testing.T) {
	tests := []struct {
		method   HTTPMethod
		expected string
		valid    bool
	}{
		{GET, "GET", true},
		{HEAD, "HEAD", true},
		{POST, "POST", true},
		{PUT, "PUT", true},
		{PATCH, "PATCH", true},
		{DELETE, "DELETE", true},
		{OPTIONS, "OPTIONS", true},
		{CONNECT, "CONNECT", true},
		{TRACE, "TRACE", true},
		{HTTPMethod(999), "", false},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.valid, tt.method.IsValid())
			if tt.valid {
				assert.Equal(t, tt.expected, tt.method.GetName())
				assert.Equal(t, tt.expected, tt.method.String())
			} else {
				assert.Panics(t, func() { tt.method.GetName() })
				assert.Contains(t, tt.method.String(), "UNKNOWN")
			}
		})
	}
}

func TestSentinelErrors(t *testing.T) {
	assert.NotNil(t, ErrReservedPath)
	assert.NotNil(t, ErrDuplicateRoute)
	assert.NotNil(t, ErrNilRegistration)
	assert.NotNil(t, ErrNilAPI)
	assert.NotNil(t, ErrNilHandler)
	assert.NotNil(t, ErrInvalidMethod)
}
