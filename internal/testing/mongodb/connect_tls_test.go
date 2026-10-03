package mongodb

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func setMockTLSPing(t *testing.T, attempts *int) {
	origPing := pingClient
	t.Cleanup(func() { pingClient = origPing })
	pingClient = func(context.Context, *mongo.Client) error {
		*attempts++
		if *attempts == 1 {
			return errors.New("server requires TLS")
		}
		return nil
	}
}

func TestIsTLSError(t *testing.T) {
	assert.False(t, isTLSError(nil))
	assert.False(t, isTLSError(errors.New("connection refused")))
	assert.True(t, isTLSError(errors.New("server requires TLS")))
	assert.True(t, isTLSError(errors.New("failed ssl handshake")))
	assert.True(t, isTLSError(errors.New("connection closed")))
}

func TestConnect_TLSFallbackSuccess(t *testing.T) {
	attempts := 0
	setMockTLSPing(t, &attempts)

	ctx := context.Background()
	client, err := Connect(ctx, WithPing(true), WithTLS(false))
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, 2, attempts)
}

func TestConnectURI_TLSFallbackSuccess(t *testing.T) {
	attempts := 0
	setMockTLSPing(t, &attempts)

	ctx := context.Background()
	uri := "mongodb://localhost:27017"
	client, err := ConnectURI(ctx, uri, WithPing(true), WithTLS(false))
	require.NoError(t, err)
	require.NotNil(t, client)
	assert.Equal(t, 2, attempts)
}

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
