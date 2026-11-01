package client

import (
	"bytes"
	"context"
	"log/slog"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

func assertGCSDriverOverrideLogs(t *testing.T, entries []map[string]any) {
	t.Helper()
	initLog := findStorageLog(entries, "initializing storage client")
	require.NotNil(t, initLog)
	assert.Equal(t, "when_supported", initLog["checksum_mode"])
	compatLog := findStorageLog(entries, "storage client compatibility configuration")
	require.NotNil(t, compatLog)
	assert.Equal(t, BehaviorChecksumExplicit, compatLog["checksum_behavior_change"])
	assert.Equal(t, "WHEN_SUPPORTED", compatLog["request_checksum_calculation"])
}

func TestConnect_GCS_DriverOverrideWhenSupported_DebugAndInfo(t *testing.T) {
	isolateAWSEnv(t)
	var output bytes.Buffer
	setStorageTestLogger(t, slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug}))
	cfg := config.DefaultConfig()
	cfg.Bucket = "test-bucket"
	cfg.Endpoint = "https://storage.googleapis.com"
	driverOpt := WithDriverOptions(func(o *s3.Options) {
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenSupported
	})
	_, err := ConnectWithConfig(context.Background(), cfg, driverOpt, WithSkipProbe(true))
	require.NoError(t, err)
	assertGCSDriverOverrideLogs(t, parseStorageLogs(t, output.Bytes()))
}

func assertAWSEnvOverrideLogs(t *testing.T, entries []map[string]any) {
	t.Helper()
	initLog := findStorageLog(entries, "initializing storage client")
	require.NotNil(t, initLog)
	assert.Equal(t, "when_required", initLog["checksum_mode"])
	compLog := findStorageLog(entries, "storage client initialization complete")
	require.NotNil(t, compLog)
	assert.Equal(t, "when_required", compLog["checksum_mode"])
	compatLog := findStorageLog(entries, "storage client compatibility configuration")
	require.NotNil(t, compatLog)
	assert.Equal(t, BehaviorChecksumExplicit, compatLog["checksum_behavior_change"])
}

func TestConnect_AWS_EnvChecksumOverride_DebugAndInfo(t *testing.T) {
	isolateAWSEnv(t)
	t.Setenv("AWS_REQUEST_CHECKSUM_CALCULATION", "WHEN_REQUIRED")
	var output bytes.Buffer
	setStorageTestLogger(t, slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug}))
	cfg := config.DefaultConfig()
	cfg.Bucket = "test-bucket"
	cfg.Region = "us-east-1"
	_, err := ConnectWithConfig(context.Background(), cfg, WithSkipProbe(true))
	require.NoError(t, err)
	assertAWSEnvOverrideLogs(t, parseStorageLogs(t, output.Bytes()))
}

func TestConnect_DetectionFailure_NoDebugRecord(t *testing.T) {
	isolateAWSEnv(t)
	var output bytes.Buffer
	setStorageTestLogger(t, slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug}))
	cfg := config.DefaultConfig()
	cfg.Bucket = "test-bucket"
	cfg.Endpoint = "://invalid-url"
	_, err := ConnectWithConfig(context.Background(), cfg, WithSkipProbe(true))
	require.Error(t, err)
	entries := parseStorageLogs(t, output.Bytes())
	assert.Nil(t, findStorageLog(entries, "evaluating storage provider detection"))
}
