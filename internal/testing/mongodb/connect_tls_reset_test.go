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

func setMockPingWithErr(t *testing.T, attempts *int, initialErr error) {
	origPing := pingClient
	origProbe := probeClient
	t.Cleanup(func() {
		pingClient = origPing
		probeClient = origProbe
	})
	pingClient = func(context.Context, *mongo.Client) error {
		*attempts++
		if *attempts == 1 {
			return initialErr
		}
		return nil
	}
	probeClient = func(context.Context, *client.Client, string) error {
		return nil
	}
}

func runTestingFallbackTest(t *testing.T, initialErr error) {
	attempts := 0
	setMockPingWithErr(t, &attempts, initialErr)

	ctx := context.Background()
	tc, err := Connect(ctx, WithTLS(false))
	require.NoError(t, err)
	require.NotNil(t, tc)
	assert.Equal(t, 2, attempts)
}

func TestConnect_TLSFallback_ConnectionReset(t *testing.T) {
	runTestingFallbackTest(t, errors.New("read tcp 127.0.0.1:443: read: connection reset by peer"))
}

func TestConnect_TLSFallback_IncompleteRead(t *testing.T) {
	runTestingFallbackTest(t, errors.New("incomplete read of full message: connection reset by peer"))
}

func TestConnect_TLSFallback_BrokenPipe(t *testing.T) {
	runTestingFallbackTest(t, errors.New("write: broken pipe"))
}

func TestConnect_TLSFallback_ServerSelectionError(t *testing.T) {
	runTestingFallbackTest(t, errors.New("server selection error: context deadline exceeded"))
}

func TestConnectURI_TLSFallback_ConnectionReset(t *testing.T) {
	attempts := 0
	setMockPingWithErr(t, &attempts, errors.New("read: connection reset by peer"))

	ctx := context.Background()
	tc, err := ConnectURI(ctx, "mongodb://localhost:27017", WithTLS(false))
	require.NoError(t, err)
	require.NotNil(t, tc)
	assert.Equal(t, 2, attempts)
}
