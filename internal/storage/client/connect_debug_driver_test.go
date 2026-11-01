package client

import (
	"bytes"
	"context"
	"log/slog"
	"sync/atomic"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nawaphonOHM/whatever/v3/internal/storage/config"
)

func debugDriverTestConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Bucket = "test-bucket"
	cfg.Region = "us-east-1"
	cfg.Endpoint = "https://user:secret@custom.minio.internal:9000"
	return cfg
}

func assertDetectionAndCompletionLogs(t *testing.T, output string, entries []map[string]any) {
	t.Helper()
	detectLog := findStorageLog(entries, "evaluating storage provider detection")
	require.NotNil(t, detectLog)
	assert.Equal(t, "https://custom.minio.internal:9000", detectLog["raw_endpoint"])
	assert.NotContains(t, output, "secret")

	compLog := findStorageLog(entries, "storage client initialization complete")
	require.NotNil(t, compLog)
	assert.Equal(t, "when_required", compLog["checksum_mode"])
}

func assertCompatibilityLogs(t *testing.T, entries []map[string]any) {
	t.Helper()
	compatLog := findStorageLog(entries, "storage client compatibility configuration")
	require.NotNil(t, compatLog)
	assert.Equal(t, "WHEN_REQUIRED", compatLog["request_checksum_calculation"])
	assert.Equal(t, true, compatLog["has_driver_overrides"])
}

func assertDriverLogOutput(t *testing.T, output string, entries []map[string]any) {
	t.Helper()
	assertDetectionAndCompletionLogs(t, output, entries)
	assertCompatibilityLogs(t, entries)
}

func TestConnectWithConfig_DriverCallbackCount_DebugEnabled(t *testing.T) {
	isolateAWSEnv(t)
	var output bytes.Buffer
	setStorageTestLogger(t, slog.NewJSONHandler(&output, &slog.HandlerOptions{Level: slog.LevelDebug}))

	var callCount atomic.Int32
	driverOpt := WithDriverOptions(func(o *s3.Options) {
		callCount.Add(1)
		o.RequestChecksumCalculation = aws.RequestChecksumCalculationWhenRequired
	})

	_, err := ConnectWithConfig(context.Background(), debugDriverTestConfig(), driverOpt, WithSkipProbe(true))
	require.NoError(t, err)
	assert.Equal(t, int32(1), callCount.Load())
	assertDriverLogOutput(t, output.String(), parseStorageLogs(t, output.Bytes()))
}
