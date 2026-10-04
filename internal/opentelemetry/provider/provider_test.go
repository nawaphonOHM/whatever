package provider

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v2/internal/opentelemetry/config"
)

const (
	invalidSampleRate = 1.5
	testServiceName   = "test-service"
	testEndpoint      = "localhost:4317"
)

func TestInitTracerProvider_NilConfig(t *testing.T) {
	ctx := context.Background()
	shutdown, err := InitTracerProvider(ctx, nil)
	require.NoError(t, err)
	require.NotNil(t, shutdown)

	err = shutdown(ctx)
	assert.NoError(t, err)
}

func TestInitTracerProvider_DisabledConfig(t *testing.T) {
	ctx := context.Background()
	cfg := &config.Config{
		Enabled: false,
	}

	shutdown, err := InitTracerProvider(ctx, cfg)
	require.NoError(t, err)
	require.NotNil(t, shutdown)

	err = shutdown(ctx)
	assert.NoError(t, err)
}

func TestInitTracerProvider_InvalidConfig(t *testing.T) {
	tests := []struct {
		expectedErr error
		cfg         *config.Config
		name        string
	}{
		{
			name: "empty service name",
			cfg: &config.Config{
				ServiceName: "",
				Endpoint:    testEndpoint,
				Protocol:    config.ProtocolGRPC,
				SampleRate:  1.0,
				Enabled:     true,
			},
			expectedErr: config.ErrEmptyServiceName,
		},
		{
			name: "empty endpoint",
			cfg: &config.Config{
				ServiceName: testServiceName,
				Endpoint:    "",
				Protocol:    config.ProtocolGRPC,
				SampleRate:  1.0,
				Enabled:     true,
			},
			expectedErr: config.ErrEmptyEndpoint,
		},
		{
			name: "invalid protocol",
			cfg: &config.Config{
				ServiceName: testServiceName,
				Endpoint:    testEndpoint,
				Protocol:    "tcp",
				SampleRate:  1.0,
				Enabled:     true,
			},
			expectedErr: config.ErrInvalidProtocol,
		},
		{
			name: "invalid sample rate",
			cfg: &config.Config{
				ServiceName: testServiceName,
				Endpoint:    testEndpoint,
				Protocol:    config.ProtocolGRPC,
				SampleRate:  invalidSampleRate,
				Enabled:     true,
			},
			expectedErr: config.ErrInvalidSampleRate,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			shutdown, err := InitTracerProvider(context.Background(), tt.cfg)
			require.Error(t, err)
			assert.Nil(t, shutdown)
			assert.True(t, errors.Is(err, tt.expectedErr))
		})
	}
}
