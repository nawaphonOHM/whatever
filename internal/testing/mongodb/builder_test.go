package mongodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	testMaxPoolLimit  = uint64(200)
	testConfigMaxPool = uint64(50)
	testConfigMinPool = uint64(10)
	testSampleURI     = "mongodb://localhost:27017/test_db"
)

func TestBuildURI_Defaults(t *testing.T) {
	opts := DefaultOptions()
	uri := BuildURI(opts, false)
	assert.Equal(t, "mongodb://localhost:27017/?uuidRepresentation=unspecified&tls=false", uri)
}

func TestBuildURI_FullOptions(t *testing.T) {
	opts := NewOptions(
		WithProtocol("mongodb+srv"),
		WithHost("cluster0.example.com"),
		WithPort(0),
		WithDatabase("analytics"),
		WithUsername("admin"),
		WithPassword("secret"),
		WithAuthSource("admin"),
		WithAppName("my-test-suite"),
		WithDirectConnection(true),
		WithUUIDRepresentation("csharpLegacy"),
	)
	uri := BuildURI(opts, true)
	expected := "mongodb+srv://admin:secret@cluster0.example.com/analytics" +
		"?uuidRepresentation=csharpLegacy&tls=true&authSource=admin&appName=my-test-suite&directConnection=true"
	assert.Equal(t, expected, uri)
}

func TestBuildURI_NilOptions(t *testing.T) {
	uri := BuildURI(nil, false)
	assert.Equal(t, "mongodb://localhost:27017/?uuidRepresentation=unspecified&tls=false", uri)
}

func TestBuildClientOptions_Configuration(t *testing.T) {
	extra := options.Client().SetMaxPoolSize(testMaxPoolLimit)
	opts := NewOptions(
		WithHost("127.0.0.1"),
		WithPort(testCustomPort),
		WithDatabase("test_db"),
		WithConnectTimeout(3*time.Second),
		WithServerSelectionTimeout(2*time.Second),
		WithSocketTimeout(4*time.Second),
		WithAppName("app-test"),
		WithPoolLimits(testConfigMaxPool, testConfigMinPool),
		WithDirectConnection(true),
	)

	driverOpts := BuildClientOptions(opts, extra)
	require.NotNil(t, driverOpts)
	assert.Equal(t, 3*time.Second, *driverOpts.ConnectTimeout)
	assert.Equal(t, 2*time.Second, *driverOpts.ServerSelectionTimeout)
	assert.Equal(t, 4*time.Second, *driverOpts.Timeout)
	assert.Equal(t, "app-test", *driverOpts.AppName)
	assert.True(t, *driverOpts.Direct)
	assert.Equal(t, testMaxPoolLimit, *driverOpts.MaxPoolSize)
}

func TestBuildClientOptionsWithTLS_Override(t *testing.T) {
	opts := NewOptions(WithHost("127.0.0.1"), WithTLS(false))
	driverOpts := BuildClientOptionsWithTLS(opts, true)
	require.NotNil(t, driverOpts)
}

func TestResolveURIFromOptions_Nil(t *testing.T) {
	assert.Equal(t, "", resolveURIFromOptions(nil))
}

func TestResolveURIFromOptions_NoTLS(t *testing.T) {
	opts := &Options{URI: testSampleURI, EnableTLS: false}
	assert.Equal(t, testSampleURI, resolveURIFromOptions(opts))
}

func TestResolveURIFromOptions_WithTLS(t *testing.T) {
	opts := &Options{URI: testSampleURI, EnableTLS: true}
	assert.Equal(t, testSampleURI+"?tls=true", resolveURIFromOptions(opts))

	optsWithQuery := &Options{URI: "mongodb://localhost:27017/test_db?replicaSet=rs0", EnableTLS: true}
	assert.Equal(t, "mongodb://localhost:27017/test_db?replicaSet=rs0&tls=true", resolveURIFromOptions(optsWithQuery))
}

func TestResolveURI(t *testing.T) {
	assert.Equal(t, "mongodb://localhost:27017/?uuidRepresentation=unspecified&tls=false", resolveURI(nil))

	optsWithURI := &Options{URI: testSampleURI}
	assert.Equal(t, testSampleURI, resolveURI(optsWithURI))

	optsWithURITLS := &Options{URI: testSampleURI, EnableTLS: true}
	assert.Equal(t, testSampleURI+"?tls=true", resolveURI(optsWithURITLS))

	optsWithoutURI := &Options{Host: "127.0.0.1", Port: DefaultPort, Database: "test_db", EnableTLS: true}
	expectedURI := "mongodb://127.0.0.1:27017/test_db?uuidRepresentation=unspecified&tls=true"
	assert.Equal(t, expectedURI, resolveURI(optsWithoutURI))
}

func TestResolveEnableTLS(t *testing.T) {
	assert.False(t, resolveEnableTLS(nil))
	assert.False(t, resolveEnableTLS(&Options{EnableTLS: false}))
	assert.True(t, resolveEnableTLS(&Options{EnableTLS: true}))
}
