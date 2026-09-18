package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const invalidMethodCode = 999

// methodCase is one HTTPMethod table row.
type methodCase struct {
	expected string
	method   HTTPMethod
	valid    bool
}

// allMethodCases returns name/validity coverage for HTTPMethod.
func allMethodCases() []methodCase {
	return []methodCase{
		{expected: "GET", method: GET, valid: true},
		{expected: "HEAD", method: HEAD, valid: true},
		{expected: "POST", method: POST, valid: true},
		{expected: "PUT", method: PUT, valid: true},
		{expected: "PATCH", method: PATCH, valid: true},
		{expected: "DELETE", method: DELETE, valid: true},
		{expected: "OPTIONS", method: OPTIONS, valid: true},
		{expected: "CONNECT", method: CONNECT, valid: true},
		{expected: "TRACE", method: TRACE, valid: true},
		{expected: "", method: HTTPMethod(invalidMethodCode), valid: false},
	}
}

// assertMethodCase checks one HTTPMethod table row.
func assertMethodCase(t *testing.T, tt methodCase) {
	t.Helper()
	assert.Equal(t, tt.valid, tt.method.IsValid())
	if !tt.valid {
		assert.Panics(t, func() { tt.method.GetName() })
		assert.Contains(t, tt.method.String(), "UNKNOWN")
		return
	}
	assert.Equal(t, tt.expected, tt.method.GetName())
	assert.Equal(t, tt.expected, tt.method.String())
}

// TestHTTPMethod covers valid and invalid HTTP method codes.
func TestHTTPMethod(t *testing.T) {
	for _, tt := range allMethodCases() {
		name := tt.expected
		if name == "" {
			name = "invalid"
		}
		t.Run(name, func(t *testing.T) {
			assertMethodCase(t, tt)
		})
	}
}

// TestSentinelErrors ensures public sentinel errors are exposed.
func TestSentinelErrors(t *testing.T) {
	assert.NotNil(t, ErrReservedPath)
	assert.NotNil(t, ErrDuplicateRoute)
	assert.NotNil(t, ErrNilRegistration)
	assert.NotNil(t, ErrNilAPI)
	assert.NotNil(t, ErrNilHandler)
	assert.NotNil(t, ErrInvalidMethod)
}
