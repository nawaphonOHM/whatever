package mongodb_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v2/pkg/testing/mongodb"
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

const testCustomPort = 27018

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
