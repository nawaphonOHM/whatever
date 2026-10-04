package mongodb_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/nawaphonOHM/whatever/v2/pkg/testing/mongodb"
)

func TestWithOptions_Timeouts(t *testing.T) {
	opts := mongodb.NewOptions(
		mongodb.WithConnectTimeout(3*time.Second),
		mongodb.WithServerSelectionTimeout(2*time.Second),
		mongodb.WithSocketTimeout(4*time.Second),
		mongodb.WithMaxConnIdleTime(5*time.Minute),
	)
	assert.Equal(t, 3*time.Second, opts.ConnectTimeout)
	assert.Equal(t, 2*time.Second, opts.ServerSelectionTimeout)
	assert.Equal(t, 4*time.Second, opts.SocketTimeout)
	assert.Equal(t, 5*time.Minute, opts.MaxConnIdleTime)
}

const (
	testMinPoolSize = 10
	testMaxPoolSize = 50
)

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
