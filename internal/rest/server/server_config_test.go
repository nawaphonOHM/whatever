package server

import (
	"os"
	"testing"

	intcfg "github.com/nawaphonOHM/whatever/internal/rest/config"
	"github.com/stretchr/testify/require"
)

var serverConfigEnvKeys = []string{
	"OHM9996_SERVER_HOST",
	"OHM9996_SERVER_PORT",
	"OHM9996_GIN_MODE",
	"OHM9996_APP_VERSION",
	"OHM9996_SERVER_DISPLAY_NAME",
	"OHM9996_SERVER_PROFILE_PATH",
	"OHM9996_SERVER_READ_TIMEOUT",
	"OHM9996_SERVER_WRITE_TIMEOUT",
	"OHM9996_SERVER_IDLE_TIMEOUT",
	"OHM9996_SERVER_SHUTDOWN_TIMEOUT",
	"OHM9996_SERVER_TRUSTED_PROXIES",
	"OHM9996_SERVER_REMOTE_IP_HEADERS",
	"OHM9996_SERVER_FORWARDED_BY_CLIENT_IP",
	"OHM9996_SERVER_REDIRECT_TRAILING_SLASH",
	"OHM9996_SERVER_REDIRECT_FIXED_PATH",
	"OHM9996_SERVER_HANDLE_METHOD_NOT_ALLOWED",
	"OHM9996_SERVER_USE_RAW_PATH",
	"OHM9996_SERVER_UNESCAPE_PATH_VALUES",
	"OHM9996_SERVER_REMOVE_EXTRA_SLASH",
	"OHM9996_SERVER_MAX_BODY_SIZE",
	"OHM9996_SERVER_READ_HEADER_TIMEOUT",
	"OHM9996_SERVER_MAX_HEADER_BYTES",
	"OHM9996_SERVER_ENABLE_ACCESS_LOG",
	"OHM9996_SERVER_ENABLE_METRICS",
	"OHM9996_SERVER_ENABLE_PROFILING",
}

func clearServerConfigEnv(t *testing.T) {
	t.Helper()
	for _, key := range serverConfigEnvKeys {
		value, exists := os.LookupEnv(key)
		require.NoError(t, os.Unsetenv(key))
		t.Cleanup(func() {
			if exists {
				require.NoError(t, os.Setenv(key, value))
				return
			}
			require.NoError(t, os.Unsetenv(key))
		})
	}
}

func loadServerConfig(t *testing.T) *Config {
	t.Helper()
	cfg, err := intcfg.Load[Config]()
	require.NoError(t, err)
	return cfg
}
