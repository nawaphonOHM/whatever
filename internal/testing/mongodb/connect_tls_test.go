package mongodb

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v2/internal/mongodb/client"
)

func setMockTLSPing(t *testing.T, attempts *int) {
	origPing := pingClient
	origProbe := probeClient
	t.Cleanup(func() {
		pingClient = origPing
		probeClient = origProbe
	})
	pingClient = func(context.Context, *mongo.Client) error {
		*attempts++
		if *attempts == 1 {
			return errors.New("server requires TLS")
		}
		return nil
	}
	probeClient = func(context.Context, *client.Client, string) error {
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

func TestIsTLSError_NetworkAndResetPatterns(t *testing.T) {
	assert.True(t, isTLSError(errors.New("connection reset by peer")))
	assert.True(t, isTLSError(errors.New("incomplete read of full message")))
	assert.True(t, isTLSError(errors.New("broken pipe")))
	assert.True(t, isTLSError(errors.New("server selection error")))
	assert.True(t, isTLSError(errors.New("read: connection reset by peer")))
	assert.True(t, isTLSError(errors.New("incomplete read of full message: read tcp: connection reset by peer")))
	assert.True(t, isTLSError(errors.New("server selection error: context deadline exceeded")))
	assert.True(t, isTLSError(errors.New("write: broken pipe")))
}

func TestConnect_TLSFallbackSuccess(t *testing.T) {
	attempts := 0
	setMockTLSPing(t, &attempts)

	ctx := context.Background()
	tc, err := Connect(ctx, WithTLS(false))
	require.NoError(t, err)
	require.NotNil(t, tc)
	assert.Equal(t, 2, attempts)
}

func TestConnectURI_TLSFallbackSuccess(t *testing.T) {
	attempts := 0
	setMockTLSPing(t, &attempts)

	ctx := context.Background()
	uri := "mongodb://localhost:27017"
	tc, err := ConnectURI(ctx, uri, WithTLS(false))
	require.NoError(t, err)
	require.NotNil(t, tc)
	assert.Equal(t, 2, attempts)
}
