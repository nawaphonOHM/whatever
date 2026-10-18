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
