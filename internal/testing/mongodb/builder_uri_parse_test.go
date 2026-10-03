package mongodb

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtractDatabaseFromURI(t *testing.T) {
	testCases := []struct {
		name     string
		uri      string
		expected string
	}{
		{"empty string", "", ""},
		{"no path", "mongodb://localhost:27017", ""},
		{"root path", "mongodb://localhost:27017/", ""},
		{"with database", "mongodb://localhost:27017/my_db", "my_db"},
		{"with query params", "mongodb://localhost:27017/my_db?authSource=admin", "my_db"},
		{"escaped name", "mongodb://localhost:27017/my%20test%20db", "my test db"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := extractDatabaseFromURI(tc.uri)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func TestSanitizeURIPath_Invalid(t *testing.T) {
	assert.Equal(t, "", sanitizeURIPath(":invalid-url"))
}
