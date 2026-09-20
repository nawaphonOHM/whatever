package contracts

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

type methodCase struct {
	name   string
	method HTTPMethod
}

func validMethodCases() []methodCase {
	return []methodCase{
		{name: http.MethodGet, method: GET},
		{name: http.MethodHead, method: HEAD},
		{name: http.MethodPost, method: POST},
		{name: http.MethodPut, method: PUT},
		{name: http.MethodPatch, method: PATCH},
		{name: http.MethodDelete, method: DELETE},
		{name: http.MethodOptions, method: OPTIONS},
		{name: http.MethodConnect, method: CONNECT},
		{name: http.MethodTrace, method: TRACE},
	}
}

func TestHTTPMethodGetNameValid(t *testing.T) {
	for _, tc := range validMethodCases() {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.name, tc.method.GetName())
		})
	}
}

func TestHTTPMethodGetNamePanic(t *testing.T) {
	assert.Panics(t, func() {
		HTTPMethod(999).GetName()
	})
}

func TestHTTPMethodString(t *testing.T) {
	for _, tc := range validMethodCases() {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.name, tc.method.String())
		})
	}
	t.Run("unknown", func(t *testing.T) {
		assert.Equal(t, "UNKNOWN(999)", HTTPMethod(999).String())
	})
}

func TestHTTPMethodIsValid(t *testing.T) {
	assert.True(t, GET.IsValid())
	assert.True(t, TRACE.IsValid())
	assert.False(t, HTTPMethod(-1).IsValid())
	assert.False(t, HTTPMethod(99).IsValid())
}
