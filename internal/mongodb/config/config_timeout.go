package config

import "time"

// dur shortens time.Duration lines for struct tags.
type dur = time.Duration

// TimeoutFields defines operation and lifecycle timeout settings for MongoDB.
type TimeoutFields struct {
	ConnectTimeout         dur `env:"OHM9996_MONGODB_CONNECT_TIMEOUT" envDefault:"10s"`
	ServerSelectionTimeout dur `env:"OHM9996_MONGODB_SERVER_SELECTION_TIMEOUT" envDefault:"5s"`
	SocketTimeout          dur `env:"OHM9996_MONGODB_SOCKET_TIMEOUT" envDefault:"10s"`
	MaxConnIdleTime        dur `env:"OHM9996_MONGODB_MAX_CONN_IDLE_TIME" envDefault:"10m"`
}
