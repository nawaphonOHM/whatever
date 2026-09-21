package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/nawaphonOHM/whatever/internal/mongodb/config"
)

// TestBuildClientOptions_WithExtraDriverOptions tests merging extra options.
func TestBuildClientOptions_WithExtraDriverOptions(t *testing.T) {
	cfg := config.DefaultConfig()
	extra1 := options.Client().SetDirect(true)
	extra2 := options.Client().SetAppName("override-app")

	clientOpts := BuildClientOptions(cfg, extra1, extra2)
	require.NotNil(t, clientOpts)

	assert.True(t, *clientOpts.Direct)
	assert.Equal(t, "override-app", *clientOpts.AppName)
}

// TestNewOptions_And_WithDriverOptions tests options container evaluation.
func TestNewOptions_And_WithDriverOptions(t *testing.T) {
	driverOpt1 := options.Client().SetDirect(true)
	driverOpt2 := options.Client().SetAppName("test-app")

	opts := NewOptions(WithDriverOptions(driverOpt1, nil, driverOpt2), nil)
	require.NotNil(t, opts)
	require.Len(t, opts.DriverOptions, 2)
	assert.Same(t, driverOpt1, opts.DriverOptions[0])
	assert.Same(t, driverOpt2, opts.DriverOptions[1])
}

// TestNewOptions_Empty tests NewOptions when no options are provided.
func TestNewOptions_Empty(t *testing.T) {
	opts := NewOptions()
	require.NotNil(t, opts)
	assert.Empty(t, opts.DriverOptions)
}
