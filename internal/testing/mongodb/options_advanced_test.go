package mongodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	testTimeout3s   = 3 * time.Second
	testTimeout4s   = 4 * time.Second
	testIdleTime5m  = 5 * time.Minute
	testMinPoolSize = 10
	testMaxPoolSize = 50
)

func TestWithOptions_Timeouts(t *testing.T) {
	opts := NewOptions(
		WithConnectTimeout(testTimeout3s),
		WithServerSelectionTimeout(2*time.Second),
		WithSocketTimeout(testTimeout4s),
		WithMaxConnIdleTime(testIdleTime5m),
	)
	assert.Equal(t, testTimeout3s, opts.ConnectTimeout)
	assert.Equal(t, 2*time.Second, opts.ServerSelectionTimeout)
	assert.Equal(t, testTimeout4s, opts.SocketTimeout)
	assert.Equal(t, testIdleTime5m, opts.MaxConnIdleTime)
}

func TestWithOptions_PoolAndFlags(t *testing.T) {
	opts := NewOptions(
		WithPoolLimits(testMinPoolSize, testMaxPoolSize),
		WithTLS(true),
		WithDirectConnection(true),
		WithUUIDRepresentation("standard"),
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
	opts := NewOptions(
		WithDriverOptions(d1, nil),
		WithDriverOptions(d2),
	)
	require.Len(t, opts.DriverOptions, 2)
	assert.Equal(t, d1, opts.DriverOptions[0])
	assert.Equal(t, d2, opts.DriverOptions[1])
}

func TestOptions_IsFirestore(t *testing.T) {
	var nilOpts *Options
	assert.False(t, nilOpts.IsFirestore())

	assert.False(t, (&Options{Host: "localhost"}).IsFirestore())
	assert.True(t, (&Options{Host: "test.firestore.goog"}).IsFirestore())
	assert.True(t, (&Options{URI: "mongodb://user:pass@test.firestore.goog:443/db"}).IsFirestore())
}

func TestOptions_BuildURI_Firestore(t *testing.T) {
	opts := &Options{
		Host:     "abc.asia-southeast1.firestore.goog",
		Database: "testdb",
		Username: "user",
		Password: "password",
	}
	uri := BuildURI(opts, false)
	assert.Contains(t, uri, "abc.asia-southeast1.firestore.goog:443")
	assert.Contains(t, uri, "tls=true")
	assert.Contains(t, uri, "loadBalanced=true")
	assert.Contains(t, uri, "retryWrites=false")
	assert.Contains(t, uri, "authMechanism=SCRAM-SHA-256")
}

func TestOptions_BuildClientOptions_Firestore(t *testing.T) {
	opts := &Options{
		Host: "abc.asia-southeast1.firestore.goog",
	}
	clientOpts := BuildClientOptions(opts)
	assert.NotNil(t, clientOpts)
	assert.NotNil(t, clientOpts.LoadBalanced)
	assert.True(t, *clientOpts.LoadBalanced)
	assert.NotNil(t, clientOpts.RetryWrites)
	assert.False(t, *clientOpts.RetryWrites)
}
