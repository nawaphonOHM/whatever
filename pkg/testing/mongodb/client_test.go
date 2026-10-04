package mongodb_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v2/pkg/testing/mongodb"
)

const (
	expectedDefaultPort     = 27017
	expectedDefaultMaxPool  = uint64(100)
	expectedDefaultMinPool  = uint64(5)
	expectedConnectTimeout  = 10 * time.Second
	expectedServerTimeout   = 5 * time.Second
	expectedSocketTimeout   = 10 * time.Second
	expectedIdleTimeout     = 10 * time.Minute
	expectedDefaultProtocol = "mongodb"
	expectedDefaultHost     = "localhost"
	expectedDefaultUUID     = "unspecified"
)

func TestConstants(t *testing.T) {
	assert.Equal(t, expectedDefaultHost, mongodb.DefaultHost)
	assert.Equal(t, expectedDefaultPort, mongodb.DefaultPort)
	assert.Equal(t, expectedDefaultProtocol, mongodb.DefaultProtocol)
	assert.Equal(t, expectedConnectTimeout, mongodb.DefaultConnectTimeout)
	assert.Equal(t, expectedServerTimeout, mongodb.DefaultServerSelectionTimeout)
	assert.Equal(t, expectedSocketTimeout, mongodb.DefaultSocketTimeout)
	assert.Equal(t, expectedIdleTimeout, mongodb.DefaultMaxConnIdleTime)
	assert.Equal(t, expectedDefaultMaxPool, mongodb.DefaultMaxPoolSize)
	assert.Equal(t, expectedDefaultMinPool, mongodb.DefaultMinPoolSize)
	assert.Equal(t, expectedDefaultUUID, mongodb.DefaultUUIDRepresentation)
}

func TestSentinelErrors(t *testing.T) {
	require.Error(t, mongodb.ErrNilClient)
	require.Error(t, mongodb.ErrNilConfig)
	require.Error(t, mongodb.ErrEmptyURI)
	require.Error(t, mongodb.ErrInvalidPort)
	require.Error(t, mongodb.ErrNilTestingTB)

	assert.True(t, errors.Is(mongodb.ErrNilClient, mongodb.ErrNilClient))
	assert.True(t, errors.Is(mongodb.ErrNilConfig, mongodb.ErrNilConfig))
	assert.True(t, errors.Is(mongodb.ErrEmptyURI, mongodb.ErrEmptyURI))
	assert.True(t, errors.Is(mongodb.ErrInvalidPort, mongodb.ErrInvalidPort))
	assert.True(t, errors.Is(mongodb.ErrNilTestingTB, mongodb.ErrNilTestingTB))
}

func assertNilClientFixtures(ctx context.Context, t *testing.T, c *mongodb.TestClient) {
	assert.ErrorIs(t, c.DropDatabase(ctx), mongodb.ErrNilClient)
	assert.ErrorIs(t, c.DropCollection(ctx, "users"), mongodb.ErrNilClient)
	_, err := c.ListCollectionNames(ctx)
	assert.ErrorIs(t, err, mongodb.ErrNilClient)
	assert.ErrorIs(t, c.TruncateCollections(ctx), mongodb.ErrNilClient)
}

func TestNilTestClient(t *testing.T) {
	var c *mongodb.TestClient
	ctx := context.Background()

	// False positive: calling methods on this nil receiver is intentional nil-safety coverage;
	// proof: TestNilTestClient in client_test.go.
	assert.Nil(t, c.Database())
	assert.Nil(t, c.Collection("users"))
	assert.Nil(t, c.RawClient())
	assert.Nil(t, c.Client())
	assert.ErrorIs(t, c.Ping(ctx), mongodb.ErrNilClient)
	assert.ErrorIs(t, c.Disconnect(ctx), mongodb.ErrNilClient)
	assert.ErrorIs(t, c.Close(), mongodb.ErrNilClient)

	assertNilClientFixtures(ctx, t, c)
}

func TestNewClient_NilRaw(t *testing.T) {
	c := mongodb.NewClient(nil, "testdb")
	require.NotNil(t, c)
	assert.Nil(t, c.Database())
	assert.Nil(t, c.Collection("users"))
	assert.Nil(t, c.RawClient())
	assert.ErrorIs(t, c.Ping(context.Background()), mongodb.ErrNilClient)
}
