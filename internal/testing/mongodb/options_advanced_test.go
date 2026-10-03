package mongodb

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func TestWithOptions_Timeouts(t *testing.T) {
	opts := NewOptions(
		WithConnectTimeout(3*time.Second),
		WithServerSelectionTimeout(2*time.Second),
		WithSocketTimeout(4*time.Second),
		WithMaxConnIdleTime(5*time.Minute),
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
	opts := NewOptions(
		WithPoolLimits(testMinPoolSize, testMaxPoolSize),
		WithPing(false),
		WithTLS(true),
		WithDirectConnection(true),
		WithUUIDRepresentation("standard"),
	)
	assert.Equal(t, uint64(testMinPoolSize), opts.MinPoolSize)
	assert.Equal(t, uint64(testMaxPoolSize), opts.MaxPoolSize)
	assert.False(t, opts.EnablePing)
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
