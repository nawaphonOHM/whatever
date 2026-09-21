package config

import (
	"time"
)

// BaseFields defines connection endpoints and metadata for MongoDB.
type BaseFields struct {
	URI      string `env:"URI"`
	Database string `env:"DATABASE"`
	Username string `env:"USERNAME"`
	Password string `env:"PASSWORD"`
	AppName  string `env:"APP_NAME" envDefault:""`
}

// TimeoutFields defines timeout settings for MongoDB connections.
type TimeoutFields struct {
	ConnectTimeout         time.Duration `env:"CONNECT_TIMEOUT" envDefault:"10s"`
	ServerSelectionTimeout time.Duration `env:"SERVER_SELECTION_TIMEOUT" envDefault:"5s"`
	SocketTimeout          time.Duration `env:"SOCKET_TIMEOUT" envDefault:"10s"`
	MaxConnIdleTime        time.Duration `env:"MAX_CONN_IDLE_TIME" envDefault:"10m"`
}

// PoolFields defines connection pooling settings for MongoDB.
type PoolFields struct {
	MaxPoolSize uint64 `env:"MAX_POOL_SIZE" envDefault:"100"`
	MinPoolSize uint64 `env:"MIN_POOL_SIZE" envDefault:"5"`
}
