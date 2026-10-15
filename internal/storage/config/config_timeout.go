package config

import "time"

// dur shortens time.Duration lines for struct tags.
type dur = time.Duration

// TimeoutFields defines operation and lifecycle timeout settings for storage.
// All timeout values must be positive (> 0). PresignExpiry cannot exceed 7 days (MaxPresignExpiry).
type TimeoutFields struct {
	ConnectTimeout dur `env:"OHM9996_STORAGE_CONNECT_TIMEOUT" envDefault:"10s"`
	RequestTimeout dur `env:"OHM9996_STORAGE_REQUEST_TIMEOUT" envDefault:"30s"`
	PingTimeout    dur `env:"OHM9996_STORAGE_PING_TIMEOUT" envDefault:"5s"`
	PresignExpiry  dur `env:"OHM9996_STORAGE_PRESIGN_EXPIRY" envDefault:"15m"`
}
