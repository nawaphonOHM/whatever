package mongodb

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()
	require.NotNil(t, cfg)

	assert.Equal(t, "mongodb://localhost:27017", cfg.URI)
	assert.Equal(t, "", cfg.Database)
	assert.Equal(t, 10*time.Second, cfg.ConnectTimeout)
	assert.Equal(t, 5*time.Second, cfg.ServerSelectionTimeout)
	assert.Equal(t, 10*time.Second, cfg.SocketTimeout)
	assert.Equal(t, uint64(100), cfg.MaxPoolSize)
	assert.Equal(t, uint64(5), cfg.MinPoolSize)
	assert.Equal(t, 10*time.Minute, cfg.MaxConnIdleTime)
	assert.Equal(t, "", cfg.AppName)

	assert.NoError(t, cfg.Validate())
}

func TestConfig_Load_Defaults(t *testing.T) {
	t.Setenv("OHM9969_MONGODB_URI", "")
	t.Setenv("OHM9969_MONGODB_DATABASE", "")
	t.Setenv("OHM9969_MONGODB_CONNECT_TIMEOUT", "")
	t.Setenv("OHM9969_MONGODB_SERVER_SELECTION_TIMEOUT", "")
	t.Setenv("OHM9969_MONGODB_SOCKET_TIMEOUT", "")
	t.Setenv("OHM9969_MONGODB_MAX_POOL_SIZE", "")
	t.Setenv("OHM9969_MONGODB_MIN_POOL_SIZE", "")
	t.Setenv("OHM9969_MONGODB_MAX_CONN_IDLE_TIME", "")
	t.Setenv("OHM9969_MONGODB_APP_NAME", "")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "mongodb://localhost:27017", cfg.URI)
	assert.Equal(t, "", cfg.Database)
	assert.Equal(t, 10*time.Second, cfg.ConnectTimeout)
	assert.Equal(t, 5*time.Second, cfg.ServerSelectionTimeout)
	assert.Equal(t, 10*time.Second, cfg.SocketTimeout)
	assert.Equal(t, uint64(100), cfg.MaxPoolSize)
	assert.Equal(t, uint64(5), cfg.MinPoolSize)
	assert.Equal(t, 10*time.Minute, cfg.MaxConnIdleTime)
	assert.Equal(t, "", cfg.AppName)

	assert.NoError(t, cfg.Validate())
}

func TestConfig_Load_EnvOverrides(t *testing.T) {
	t.Setenv("OHM9969_MONGODB_URI", "mongodb://user:pass@remote:27018")
	t.Setenv("OHM9969_MONGODB_DATABASE", "analytics")
	t.Setenv("OHM9969_MONGODB_CONNECT_TIMEOUT", "20s")
	t.Setenv("OHM9969_MONGODB_SERVER_SELECTION_TIMEOUT", "15s")
	t.Setenv("OHM9969_MONGODB_SOCKET_TIMEOUT", "30s")
	t.Setenv("OHM9969_MONGODB_MAX_POOL_SIZE", "200")
	t.Setenv("OHM9969_MONGODB_MIN_POOL_SIZE", "10")
	t.Setenv("OHM9969_MONGODB_MAX_CONN_IDLE_TIME", "15m")
	t.Setenv("OHM9969_MONGODB_APP_NAME", "my-service")

	cfg, err := LoadConfig()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "mongodb://user:pass@remote:27018", cfg.URI)
	assert.Equal(t, "analytics", cfg.Database)
	assert.Equal(t, 20*time.Second, cfg.ConnectTimeout)
	assert.Equal(t, 15*time.Second, cfg.ServerSelectionTimeout)
	assert.Equal(t, 30*time.Second, cfg.SocketTimeout)
	assert.Equal(t, uint64(200), cfg.MaxPoolSize)
	assert.Equal(t, uint64(10), cfg.MinPoolSize)
	assert.Equal(t, 15*time.Minute, cfg.MaxConnIdleTime)
	assert.Equal(t, "my-service", cfg.AppName)

	assert.NoError(t, cfg.Validate())
}

func TestConfig_Load_InvalidTypeConversion(t *testing.T) {
	t.Setenv("OHM9969_MONGODB_MAX_POOL_SIZE", "invalid-uint")

	cfg, err := LoadConfig()
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "failed to load mongodb config")
}

func TestConfig_Load_InvalidValidation(t *testing.T) {
	t.Setenv("OHM9969_MONGODB_CONNECT_TIMEOUT", "-10s")

	cfg, err := LoadConfig()
	require.Error(t, err)
	assert.Nil(t, cfg)
	assert.Contains(t, err.Error(), "invalid mongodb config")
	assert.Contains(t, err.Error(), "connect timeout cannot be negative")
}

func TestConfig_Load_DotEnvFile(t *testing.T) {
	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env.mongodb")
	content := []byte(`
OHM9969_MONGODB_URI=mongodb://dotenv-host:27017
OHM9969_MONGODB_DATABASE=dotenv_db
OHM9969_MONGODB_MAX_POOL_SIZE=75
OHM9969_MONGODB_APP_NAME=dotenv-app
`)
	err := os.WriteFile(envPath, content, 0o600)
	require.NoError(t, err)

	cfg, err := LoadConfig(envPath)
	require.NoError(t, err)
	require.NotNil(t, cfg)

	assert.Equal(t, "mongodb://dotenv-host:27017", cfg.URI)
	assert.Equal(t, "dotenv_db", cfg.Database)
	assert.Equal(t, uint64(75), cfg.MaxPoolSize)
	assert.Equal(t, "dotenv-app", cfg.AppName)
	assert.NoError(t, cfg.Validate())
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name        string
		cfg         *Config
		expectedErr string
	}{
		{
			name:        "nil config",
			cfg:         nil,
			expectedErr: "mongodb config cannot be nil",
		},
		{
			name: "empty URI",
			cfg: &Config{
				URI: "",
			},
			expectedErr: "mongodb uri cannot be empty",
		},
		{
			name: "negative connect timeout",
			cfg: &Config{
				URI:            "mongodb://localhost:27017",
				ConnectTimeout: -1 * time.Second,
			},
			expectedErr: "connect timeout cannot be negative",
		},
		{
			name: "negative server selection timeout",
			cfg: &Config{
				URI:                    "mongodb://localhost:27017",
				ServerSelectionTimeout: -1 * time.Second,
			},
			expectedErr: "server selection timeout cannot be negative",
		},
		{
			name: "negative socket timeout",
			cfg: &Config{
				URI:           "mongodb://localhost:27017",
				SocketTimeout: -1 * time.Second,
			},
			expectedErr: "socket timeout cannot be negative",
		},
		{
			name: "negative max conn idle time",
			cfg: &Config{
				URI:             "mongodb://localhost:27017",
				MaxConnIdleTime: -1 * time.Second,
			},
			expectedErr: "max conn idle time cannot be negative",
		},
		{
			name: "min pool size greater than max pool size",
			cfg: &Config{
				URI:         "mongodb://localhost:27017",
				MaxPoolSize: 5,
				MinPoolSize: 10,
			},
			expectedErr: "min pool size cannot be greater than max pool size",
		},
		{
			name: "min pool size allowed when max pool size is zero (unlimited)",
			cfg: &Config{
				URI:         "mongodb://localhost:27017",
				MaxPoolSize: 0,
				MinPoolSize: 10,
			},
			expectedErr: "",
		},
		{
			name: "min pool size equals max pool size",
			cfg: &Config{
				URI:         "mongodb://localhost:27017",
				MaxPoolSize: 10,
				MinPoolSize: 10,
			},
			expectedErr: "",
		},
		{
			name: "zero timeouts valid",
			cfg: &Config{
				URI:                    "mongodb://localhost:27017",
				ConnectTimeout:         0,
				ServerSelectionTimeout: 0,
				SocketTimeout:          0,
				MaxConnIdleTime:        0,
				MaxPoolSize:            10,
				MinPoolSize:            0,
			},
			expectedErr: "",
		},
		{
			name: "valid config with all fields",
			cfg: &Config{
				URI:                    "mongodb://localhost:27017",
				Database:               "orders",
				ConnectTimeout:         5 * time.Second,
				ServerSelectionTimeout: 3 * time.Second,
				SocketTimeout:          5 * time.Second,
				MaxPoolSize:            50,
				MinPoolSize:            5,
				MaxConnIdleTime:        5 * time.Minute,
				AppName:                "order-svc",
			},
			expectedErr: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.expectedErr != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
