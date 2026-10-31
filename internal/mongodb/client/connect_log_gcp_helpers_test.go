package client

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/nawaphonOHM/whatever/v3/internal/mongodb/config"
)

const logFieldSeverity = "severity"

func assertMongoInitEntryGCP(t *testing.T, entry map[string]any) {
	t.Helper()
	require.NotNil(t, entry, "expected initializing mongodb client entry")
	assertGCPTimestamp(t, entry)
	assertGCPSeverity(t, entry, logLevelInfo)
	assert.Equal(t, "127.0.0.1", entry["host"])
	assert.Equal(t, float64(testLogMongoPort), entry["port"])
	assert.Equal(t, "mongodb", entry["protocol"])
	assert.Equal(t, "testdb", entry["database"])
}

func assertMongoOptsEntryGCP(t *testing.T, entry map[string]any) {
	t.Helper()
	require.NotNil(t, entry, "expected resolving mongodb connection options entry")
	assertGCPTimestamp(t, entry)
	assertGCPSeverity(t, entry, logLevelInfo)
	assert.Equal(t, "127.0.0.1", entry["host"])
	assert.Equal(t, float64(testLogMongoPort), entry["port"])
	assert.Equal(t, "admin", entry["auth_source"])
	assert.Equal(t, false, entry["tls"])
	assert.Equal(t, "standard", entry["uuid_representation"])
}

func assertMongoPingEntryGCP(t *testing.T, entry map[string]any) {
	t.Helper()
	require.NotNil(t, entry, "expected mongodb ping succeeded entry")
	assertGCPTimestamp(t, entry)
	assertGCPSeverity(t, entry, logLevelInfo)
}

func assertMongoCompleteEntryGCP(t *testing.T, entry map[string]any) {
	t.Helper()
	require.NotNil(t, entry, "expected mongodb client initialization complete entry")
	assertGCPTimestamp(t, entry)
	assertGCPSeverity(t, entry, logLevelInfo)
	assert.Equal(t, "testdb", entry["database"])
}

func gcpMongoTestConfig() *config.Config {
	cfg := config.DefaultConfig()
	cfg.Host = "127.0.0.1"
	cfg.Port = testLogMongoPort
	cfg.Database = "testdb"
	cfg.AuthSource = "admin"
	cfg.UUIDRepresentation = config.UUIDRepresentationStandard
	return cfg
}

func stubMongoPingAndProbe(t *testing.T) {
	t.Helper()
	rPing := SetMockPing(func(context.Context, *mongo.Client) error { return nil })
	rProbe := SetMockProbe(func(context.Context, *Client, string) error { return nil })
	t.Cleanup(func() {
		rProbe()
		rPing()
	})
}
