package mongodb

import (
	"time"
)

// dur aliases time.Duration for concise struct definitions.
type dur = time.Duration

// BaseFields defines connection endpoints and metadata for MongoDB.
type BaseFields struct {
	URI      string `env:"URI" envDefault:"mongodb://localhost:27017"`
	Database string `env:"DATABASE" envDefault:""`
	AppName  string `env:"APP_NAME" envDefault:""`
}

// TimeoutFields defines timeout settings for MongoDB connections.
type TimeoutFields struct {
	ConnectTimeout         dur `env:"CONNECT_TIMEOUT" envDefault:"10s"`
	ServerSelectionTimeout dur `env:"SERVER_SELECTION_TIMEOUT" envDefault:"5s"`
	SocketTimeout          dur `env:"SOCKET_TIMEOUT" envDefault:"10s"`
	MaxConnIdleTime        dur `env:"MAX_CONN_IDLE_TIME" envDefault:"10m"`
}

// PoolFields defines connection pooling settings for MongoDB.
type PoolFields struct {
	MaxPoolSize uint64 `env:"MAX_POOL_SIZE" envDefault:"100"`
	MinPoolSize uint64 `env:"MIN_POOL_SIZE" envDefault:"5"`
}
