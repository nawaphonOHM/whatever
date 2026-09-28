package provider

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

func assertServiceAttributes(t *testing.T, res *resource.Resource, serviceName string) {
	attrs := res.Set()
	serviceNameVal, exists := attrs.Value(semconv.ServiceNameKey)
	require.True(t, exists)
	assert.Equal(t, serviceName, serviceNameVal.AsString())

	serviceVersionVal, exists := attrs.Value(semconv.ServiceVersionKey)
	require.True(t, exists)
	assert.Equal(t, DefaultServiceVersion, serviceVersionVal.AsString())
}

func assertTelemetryAttributes(t *testing.T, res *resource.Resource) {
	attrs := res.Set()
	telemetrySDKVal, exists := attrs.Value(semconv.TelemetrySDKNameKey)
	require.True(t, exists)
	assert.Equal(t, "opentelemetry", telemetrySDKVal.AsString())

	telemetryLangVal, exists := attrs.Value(semconv.TelemetrySDKLanguageKey)
	require.True(t, exists)
	assert.Equal(t, "go", telemetryLangVal.AsString())
}

func TestNewResource_Attributes(t *testing.T) {
	ctx := context.Background()
	serviceName := "test-resource-service"

	res, err := newResource(ctx, serviceName)
	require.NoError(t, err)
	require.NotNil(t, res)

	assertServiceAttributes(t, res, serviceName)
	assertTelemetryAttributes(t, res)
}

func TestIsFatalResourceError(t *testing.T) {
	assert.False(t, isFatalResourceError(nil))
	assert.False(t, isFatalResourceError(resource.ErrPartialResource))
	assert.True(t, isFatalResourceError(errors.New("generic error")))
}
