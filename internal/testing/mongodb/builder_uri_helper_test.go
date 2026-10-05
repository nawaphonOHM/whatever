package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInjectTLSQueryParam_AddsQueryParam(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple URI without query",
			input:    "mongodb://localhost:27017",
			expected: "mongodb://localhost:27017/?tls=true",
		},
		{
			name:     "URI with database path",
			input:    "mongodb://localhost:27017/testdb",
			expected: "mongodb://localhost:27017/testdb?tls=true",
		},
		{
			name:     "URI with existing query params",
			input:    "mongodb://localhost:27017/testdb?authSource=admin",
			expected: "mongodb://localhost:27017/testdb?authSource=admin&tls=true",
		},
		{
			name:     "URI with credentials",
			input:    "mongodb://user:pass@localhost:27017/testdb",
			expected: "mongodb://user:pass@localhost:27017/testdb?tls=true",
		},
		{
			name:     "URI with existing tls=false overwrites",
			input:    "mongodb://localhost:27017?tls=false",
			expected: "mongodb://localhost:27017/?tls=true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := injectTLSQueryParam(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
