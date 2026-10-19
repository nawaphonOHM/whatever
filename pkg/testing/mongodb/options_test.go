package mongodb_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/nawaphonOHM/whatever/v3/pkg/testing/mongodb"
)

const (
	testCustomPort  = 27018
	testTimeout3s   = 3 * time.Second
	testTimeout4s   = 4 * time.Second
	testIdleTime5m  = 5 * time.Minute
	testMinPoolSize = 10
	testMaxPoolSize = 50
)

func TestDefaultOptions(t *testing.T) {
	opts := mongodb.DefaultOptions()
	require.NotNil(t, opts)
	assert.Equal(t, mongodb.DefaultHost, opts.Host)
	assert.Equal(t, mongodb.DefaultPort, opts.Port)
	assert.Equal(t, mongodb.DefaultProtocol, opts.Protocol)
	assert.Equal(t, mongodb.DefaultUUIDRepresentation, opts.UUIDRepresentation)
	assert.Equal(t, mongodb.DefaultConnectTimeout, opts.ConnectTimeout)
	assert.False(t, opts.EnableTLS)
	assert.False(t, opts.DirectConnection)
}

func TestNewOptions_NilSafe(t *testing.T) {
	opts := mongodb.NewOptions(nil, nil)
	require.NotNil(t, opts)
	assert.Equal(t, mongodb.DefaultHost, opts.Host)
	assert.Equal(t, mongodb.DefaultPort, opts.Port)
	assert.Equal(t, mongodb.DefaultConnectTimeout, opts.ConnectTimeout)
}

func TestWithOptions_Endpoints(t *testing.T) {
	opts := mongodb.NewOptions(
		mongodb.WithHost("mongo.test.local"),
		mongodb.WithPort(testCustomPort),
		mongodb.WithProtocol("mongodb+srv"),
		mongodb.WithDatabase("integration_test"),
	)
	assert.Equal(t, "mongo.test.local", opts.Host)
	assert.Equal(t, testCustomPort, opts.Port)
	assert.Equal(t, "mongodb+srv", opts.Protocol)
	assert.Equal(t, "integration_test", opts.Database)
}

func TestWithOptions_AuthAndApp(t *testing.T) {
	opts := mongodb.NewOptions(
		mongodb.WithUsername("testuser"),
		mongodb.WithPassword("testpass"),
		mongodb.WithAuthSource("admin"),
		mongodb.WithAppName("test-suite"),
		mongodb.WithURI("mongodb://localhost:27017/test"),
	)
	assert.Equal(t, "testuser", opts.Username)
	assert.Equal(t, "testpass", opts.Password)
	assert.Equal(t, "admin", opts.AuthSource)
	assert.Equal(t, "test-suite", opts.AppName)
	assert.Equal(t, "mongodb://localhost:27017/test", opts.URI)
}

func TestWithOptions_Timeouts(t *testing.T) {
	opts := mongodb.NewOptions(
		mongodb.WithConnectTimeout(testTimeout3s),
		mongodb.WithServerSelectionTimeout(2*time.Second),
		mongodb.WithSocketTimeout(testTimeout4s),
		mongodb.WithMaxConnIdleTime(testIdleTime5m),
	)
	assert.Equal(t, testTimeout3s, opts.ConnectTimeout)
	assert.Equal(t, 2*time.Second, opts.ServerSelectionTimeout)
	assert.Equal(t, testTimeout4s, opts.SocketTimeout)
	assert.Equal(t, testIdleTime5m, opts.MaxConnIdleTime)
}

func TestWithOptions_PoolAndFlags(t *testing.T) {
	opts := mongodb.NewOptions(
		mongodb.WithPoolLimits(testMinPoolSize, testMaxPoolSize),
		mongodb.WithTLS(true),
		mongodb.WithDirectConnection(true),
		mongodb.WithUUIDRepresentation("standard"),
	)
	assert.Equal(t, uint64(testMinPoolSize), opts.MinPoolSize)
	assert.Equal(t, uint64(testMaxPoolSize), opts.MaxPoolSize)
	assert.True(t, opts.EnableTLS)
	assert.True(t, opts.DirectConnection)
	assert.Equal(t, "standard", opts.UUIDRepresentation)
}

func TestWithOptions_DriverOptions(t *testing.T) {
	d1 := options.Client().SetAppName("d1")
	d2 := options.Client().SetAppName("d2")
	opts := mongodb.NewOptions(
		mongodb.WithDriverOptions(d1, nil),
		mongodb.WithDriverOptions(d2),
	)
	require.Len(t, opts.DriverOptions, 2)
	assert.Equal(t, d1, opts.DriverOptions[0])
	assert.Equal(t, d2, opts.DriverOptions[1])
}
