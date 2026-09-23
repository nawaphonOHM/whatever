package server

import "time"

const (
	defaultMaxHeaderBytes = 1 << 20
	defaultMaxBodySize    = 32 << 20
	defaultReadHeaderSec  = 5
)

// ResourceFields defines server limits and resource options.
type ResourceFields struct {
	MaxBodySize       int64         `env:"OHM9996_SERVER_MAX_BODY_SIZE" envDefault:"33554432"`
	ReadHeaderTimeout time.Duration `env:"OHM9996_SERVER_READ_HEADER_TIMEOUT" envDefault:"5s"`
	MaxHeaderBytes    int           `env:"OHM9996_SERVER_MAX_HEADER_BYTES" envDefault:"1048576"`
}

func defaultResourceFields() ResourceFields {
	return ResourceFields{
		MaxBodySize:       defaultMaxBodySize,
		ReadHeaderTimeout: defaultReadHeaderSec * time.Second,
		MaxHeaderBytes:    defaultMaxHeaderBytes,
	}
}
