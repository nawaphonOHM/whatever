package server

import "testing"

const (
	envTrue          = "true"
	envFalse         = "false"
	overridePort     = 9090
	overrideBodySize = 67108864
	overrideHeaders  = 2097152
)

func setServerConfigOverrides(t *testing.T) {
	t.Helper()
	values := map[string]string{
		"OHM9996_SERVER_HOST":                      "127.0.0.1",
		"OHM9996_SERVER_PORT":                      "9090",
		"OHM9996_GIN_MODE":                         "test",
		"OHM9996_APP_VERSION":                      "2.0.0",
		"OHM9996_SERVER_DISPLAY_NAME":              "custom-app",
		"OHM9996_SERVER_PROFILE_PATH":              "/custom/pprof",
		"OHM9996_SERVER_READ_TIMEOUT":              "15s",
		"OHM9996_SERVER_WRITE_TIMEOUT":             "15s",
		"OHM9996_SERVER_IDLE_TIMEOUT":              "90s",
		"OHM9996_SERVER_SHUTDOWN_TIMEOUT":          "20s",
		"OHM9996_SERVER_TRUSTED_PROXIES":           "10.0.0.0/8,192.168.1.1",
		"OHM9996_SERVER_REMOTE_IP_HEADERS":         "X-Custom-IP,X-Real-IP",
		"OHM9996_SERVER_FORWARDED_BY_CLIENT_IP":    envFalse,
		"OHM9996_SERVER_REDIRECT_TRAILING_SLASH":   envFalse,
		"OHM9996_SERVER_REDIRECT_FIXED_PATH":       envTrue,
		"OHM9996_SERVER_HANDLE_METHOD_NOT_ALLOWED": envFalse,
		"OHM9996_SERVER_USE_RAW_PATH":              envTrue,
		"OHM9996_SERVER_UNESCAPE_PATH_VALUES":      envFalse,
		"OHM9996_SERVER_REMOVE_EXTRA_SLASH":        envTrue,
		"OHM9996_SERVER_MAX_BODY_SIZE":             "67108864",
		"OHM9996_SERVER_READ_HEADER_TIMEOUT":       "8s",
		"OHM9996_SERVER_MAX_HEADER_BYTES":          "2097152",
		"OHM9996_SERVER_ENABLE_ACCESS_LOG":         envFalse,
		"OHM9996_SERVER_ENABLE_METRICS":            envTrue,
		"OHM9996_SERVER_ENABLE_PROFILING":          envTrue,
	}
	for key, value := range values {
		t.Setenv(key, value)
	}
}
